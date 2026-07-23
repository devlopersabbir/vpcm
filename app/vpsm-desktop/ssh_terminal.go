package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
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

func tryParseKeyFile(path string, passphrase string) (ssh.Signer, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	// Expand ~ to user home
	if len(path) > 1 && path[0] == '~' && (path[1] == '/' || path[1] == '\\') {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, path[2:])
		}
	}

	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err == nil {
		return signer, nil
	}

	if passphrase != "" {
		if signerPass, err := ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(passphrase)); err == nil {
			return signerPass, nil
		}
	}

	return nil, err
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
	addedSignerFingerprints := make(map[string]bool)

	addSigner := func(signer ssh.Signer) {
		if signer == nil {
			return
		}
		fp := string(signer.PublicKey().Marshal())
		if !addedSignerFingerprints[fp] {
			addedSignerFingerprints[fp] = true
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}

	addKeyPathSigner := func(relPath string) {
		if relPath == "" {
			return
		}

		// Try exact path
		if signer, err := tryParseKeyFile(relPath, params.AuthSecret); err == nil {
			addSigner(signer)
			return
		}

		// Try candidate locations if path is relative
		if home, err := os.UserHomeDir(); err == nil {
			candidates := []string{
				filepath.Join(home, ".ssh", relPath),
				filepath.Join(home, relPath),
				filepath.Join(home, "Downloads", relPath),
				filepath.Join(home, "Desktop", relPath),
				filepath.Join("/Users/sabbir/own/vpcm", relPath),
			}
			for _, cand := range candidates {
				if signer, err := tryParseKeyFile(cand, params.AuthSecret); err == nil {
					addSigner(signer)
					return
				}
			}
		}
	}

	// 1. Primary Auth Method based on params.AuthType
	if params.AuthType == "password" && params.AuthSecret != "" {
		pwd := params.AuthSecret
		authMethods = append(authMethods, ssh.Password(pwd))
		authMethods = append(authMethods, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			answers := make([]string, len(questions))
			for i := range questions {
				answers[i] = pwd
			}
			return answers, nil
		}))
	}

	if (params.AuthType == "key" || params.AuthType == "keyfile") && params.AuthSecret != "" {
		addKeyPathSigner(params.AuthSecret)
		if signer, err := ssh.ParsePrivateKey([]byte(params.AuthSecret)); err == nil {
			addSigner(signer)
		}
	}

	// 2. Fallback: Process explicit AuthSecret as both key and password/keyboard-interactive
	if params.AuthSecret != "" && params.AuthType != "password" {
		addKeyPathSigner(params.AuthSecret)

		// Try as raw key string
		if signer, err := ssh.ParsePrivateKey([]byte(params.AuthSecret)); err == nil {
			addSigner(signer)
		}

		// Try as password & keyboard-interactive
		pwd := params.AuthSecret
		authMethods = append(authMethods, ssh.Password(pwd))
		authMethods = append(authMethods, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			answers := make([]string, len(questions))
			for i := range questions {
				answers[i] = pwd
			}
			return answers, nil
		}))
	}

	// 2. Auto-scan private key files in ~/.ssh/
	if homeDir, err := os.UserHomeDir(); err == nil {
		sshDir := filepath.Join(homeDir, ".ssh")
		if entries, err := os.ReadDir(sshDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				name := entry.Name()
				if name == "known_hosts" || name == "known_hosts.old" || name == "config" || (len(name) > 4 && name[len(name)-4:] == ".pub") {
					continue
				}
				keyPath := filepath.Join(sshDir, name)
				if signer, err := tryParseKeyFile(keyPath, params.AuthSecret); err == nil {
					addSigner(signer)
				}
			}
		}
	}

	// 3. Auto-scan .pem / .key files in working directory & home
	if homeDir, err := os.UserHomeDir(); err == nil {
		for _, dir := range []string{homeDir, "/Users/sabbir/own/vpcm"} {
			if entries, err := os.ReadDir(dir); err == nil {
				for _, entry := range entries {
					if entry.IsDir() {
						continue
					}
					name := entry.Name()
					if strings.HasSuffix(name, ".pem") || strings.HasSuffix(name, ".key") {
						if signer, err := tryParseKeyFile(filepath.Join(dir, name), params.AuthSecret); err == nil {
							addSigner(signer)
						}
					}
				}
			}
		}
	}

	// 4. Connect to local SSH Agent if available (SSH_AUTH_SOCK)
	if agentSock := os.Getenv("SSH_AUTH_SOCK"); agentSock != "" {
		if agentConn, err := net.DialTimeout("unix", agentSock, 2*time.Second); err == nil {
			ag := agent.NewClient(agentConn)
			if signers, err := ag.Signers(); err == nil && len(signers) > 0 {
				for _, s := range signers {
					addSigner(s)
				}
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
