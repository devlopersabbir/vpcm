package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
)

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

	authMethods := buildSSHAuthMethods(params)
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

	if err := session.Shell(); err != nil {
		a.CloseSSHTerminal(sessionID)
		return "", fmt.Errorf("failed to start shell: %w", err)
	}

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
