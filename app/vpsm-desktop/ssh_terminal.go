package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type SSHTerminalSession struct {
	ID          string
	Client      *ssh.Client
	Session     *ssh.Session
	StdinWriter io.WriteCloser
	CloseOnce   sync.Once
}

var (
	activeSessions   = make(map[string]*SSHTerminalSession)
	activeSessionsMu sync.Mutex
)

// SSHConnectionParams holds details required for direct SSH dial
type SSHConnectionParams struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	AuthType   string `json:"auth_type"`   // "password", "key", or "keyfile"
	AuthSecret string `json:"auth_secret"` // raw password, raw private key, or key filepath
	Rows       int    `json:"rows"`
	Cols       int    `json:"cols"`
}

// StartSSHTerminal initiates a direct SSH PTY session to a target host
func (a *App) StartSSHTerminal(params SSHConnectionParams) (string, error) {
	if params.Port == 0 {
		params.Port = 22
	}
	if params.Username == "" {
		params.Username = "root"
	}
	if params.Rows <= 0 {
		params.Rows = 24
	}
	if params.Cols <= 0 {
		params.Cols = 80
	}

	var authMethods []ssh.AuthMethod

	switch params.AuthType {
	case "key":
		if params.AuthSecret != "" {
			// Try as raw key string
			if signer, err := ssh.ParsePrivateKey([]byte(params.AuthSecret)); err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			} else {
				// Try as key file path
				if keyBytes, err := os.ReadFile(params.AuthSecret); err == nil {
					if signer, err := ssh.ParsePrivateKey(keyBytes); err == nil {
						authMethods = append(authMethods, ssh.PublicKeys(signer))
					}
				}
			}
		}
	case "keyfile":
		if params.AuthSecret != "" {
			if keyBytes, err := os.ReadFile(params.AuthSecret); err == nil {
				if signer, err := ssh.ParsePrivateKey(keyBytes); err == nil {
					authMethods = append(authMethods, ssh.PublicKeys(signer))
				}
			}
		}
	case "password":
		if params.AuthSecret != "" {
			authMethods = append(authMethods, ssh.Password(params.AuthSecret))
		}
	}

	// Fallback: If AuthSecret provided, test it both as key string/filepath AND password
	if params.AuthSecret != "" {
		if signer, err := ssh.ParsePrivateKey([]byte(params.AuthSecret)); err == nil {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
		if keyBytes, err := os.ReadFile(params.AuthSecret); err == nil {
			if signer, err := ssh.ParsePrivateKey(keyBytes); err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
		authMethods = append(authMethods, ssh.Password(params.AuthSecret))
	}

	// Auto-scan ALL private key files in ~/.ssh/ directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		sshDir := homeDir + "/.ssh"
		if entries, err := os.ReadDir(sshDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				name := entry.Name()
				// Ignore known hosts, config, public keys, sockets, etc.
				if name == "known_hosts" || name == "known_hosts.old" || name == "config" || (len(name) > 4 && name[len(name)-4:] == ".pub") {
					continue
				}
				keyPath := sshDir + "/" + name
				keyBytes, err := os.ReadFile(keyPath)
				if err != nil {
					continue
				}

				// Attempt unencrypted private key parse
				signer, err := ssh.ParsePrivateKey(keyBytes)
				if err == nil {
					authMethods = append(authMethods, ssh.PublicKeys(signer))
				} else {
					// If key is encrypted with passphrase and AuthSecret is provided, try parsing with passphrase
					if params.AuthSecret != "" {
						if signerWithPass, err := ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(params.AuthSecret)); err == nil {
							authMethods = append(authMethods, ssh.PublicKeys(signerWithPass))
						}
					}
				}
			}
		}
	}

	// Connect to local SSH Agent if available (SSH_AUTH_SOCK)
	if agentSock := os.Getenv("SSH_AUTH_SOCK"); agentSock != "" {
		if agentConn, err := net.DialTimeout("unix", agentSock, 2*time.Second); err == nil {
			ag := agent.NewClient(agentConn)
			if signers, err := ag.Signers(); err == nil && len(signers) > 0 {
				authMethods = append(authMethods, ssh.PublicKeys(signers...))
			}
		}
	}

	if len(authMethods) == 0 {
		return "", fmt.Errorf("no authentication method provided")
	}

	sshConfig := &ssh.ClientConfig{
		User:            params.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	targetAddr := fmt.Sprintf("%s:%d", params.Host, params.Port)

	// Direct TCP dial to the remote SSH server
	conn, err := net.DialTimeout("tcp", targetAddr, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to connect to %s: %w", targetAddr, err)
	}

	sshConn, chans, reqs, err := ssh.NewClientConn(conn, targetAddr, sshConfig)
	if err != nil {
		conn.Close()
		return "", fmt.Errorf("SSH handshake failed with %s: %w", targetAddr, err)
	}

	client := ssh.NewClient(sshConn, chans, reqs)

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}

	// Request pseudo terminal (PTY)
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", params.Rows, params.Cols, modes); err != nil {
		session.Close()
		client.Close()
		return "", fmt.Errorf("failed to request PTY: %w", err)
	}

	stdinWriter, err := session.StdinPipe()
	if err != nil {
		session.Close()
		client.Close()
		return "", fmt.Errorf("failed to open stdin pipe: %w", err)
	}

	stdoutReader, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		client.Close()
		return "", fmt.Errorf("failed to open stdout pipe: %w", err)
	}

	stderrReader, err := session.StderrPipe()
	if err != nil {
		session.Close()
		client.Close()
		return "", fmt.Errorf("failed to open stderr pipe: %w", err)
	}

	sessionID := uuid.New().String()

	termSession := &SSHTerminalSession{
		ID:          sessionID,
		Client:      client,
		Session:     session,
		StdinWriter: stdinWriter,
	}

	activeSessionsMu.Lock()
	activeSessions[sessionID] = termSession
	activeSessionsMu.Unlock()

	// Start remote shell
	if err := session.Shell(); err != nil {
		a.CloseSSHTerminal(sessionID)
		return "", fmt.Errorf("failed to start shell: %w", err)
	}

	// Stream stdout to Wails frontend event
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stdoutReader.Read(buf)
			if n > 0 {
				runtime.EventsEmit(a.ctx, fmt.Sprintf("ssh:data:%s", sessionID), string(buf[:n]))
			}
			if err != nil {
				if err != io.EOF {
					runtime.EventsEmit(a.ctx, fmt.Sprintf("ssh:error:%s", sessionID), err.Error())
				}
				break
			}
		}
		a.CloseSSHTerminal(sessionID)
	}()

	// Stream stderr to Wails frontend event
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stderrReader.Read(buf)
			if n > 0 {
				runtime.EventsEmit(a.ctx, fmt.Sprintf("ssh:data:%s", sessionID), string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	return sessionID, nil
}

// SendSSHTerminalInput sends keystrokes/data to the active SSH session
func (a *App) SendSSHTerminalInput(sessionID string, data string) error {
	activeSessionsMu.Lock()
	session, ok := activeSessions[sessionID]
	activeSessionsMu.Unlock()

	if !ok || session.StdinWriter == nil {
		return fmt.Errorf("session not found or closed")
	}

	_, err := session.StdinWriter.Write([]byte(data))
	return err
}

// ResizeSSHTerminal changes the window dimension (rows/cols) of the PTY session
func (a *App) ResizeSSHTerminal(sessionID string, rows int, cols int) error {
	activeSessionsMu.Lock()
	session, ok := activeSessions[sessionID]
	activeSessionsMu.Unlock()

	if !ok || session.Session == nil {
		return fmt.Errorf("session not found")
	}

	return session.Session.WindowChange(rows, cols)
}

// CloseSSHTerminal gracefully closes the SSH session
func (a *App) CloseSSHTerminal(sessionID string) {
	activeSessionsMu.Lock()
	session, ok := activeSessions[sessionID]
	if ok {
		delete(activeSessions, sessionID)
	}
	activeSessionsMu.Unlock()

	if ok && session != nil {
		session.CloseOnce.Do(func() {
			runtime.EventsEmit(a.ctx, fmt.Sprintf("ssh:closed:%s", sessionID), true)
			if session.StdinWriter != nil {
				session.StdinWriter.Close()
			}
			if session.Session != nil {
				session.Session.Close()
			}
			if session.Client != nil {
				session.Client.Close()
			}
		})
	}
}
