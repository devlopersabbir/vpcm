package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
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
