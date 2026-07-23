package ssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/devlopersabbir/vpcm/internal/events"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type sshService struct {
	timeout time.Duration
}

func NewService(timeout time.Duration) SSHService {
	return &sshService{timeout: timeout}
}

type realSSHClient struct {
	client *ssh.Client
}

func (s *sshService) Connect(ctx context.Context, host string, port int, username string, authType string, authSecret string) (Client, error) {
	var authMethods []ssh.AuthMethod

	if authType == "key" || authType == "keyfile" {
		var signer ssh.Signer
		var parseErr error

		trimmedSecret := strings.TrimSpace(authSecret)
		if trimmedSecret == "" {
			return nil, fmt.Errorf("no SSH key file path or key content provided")
		}

		if strings.HasPrefix(trimmedSecret, "-----BEGIN") {
			signer, parseErr = ssh.ParsePrivateKey([]byte(authSecret))
			if parseErr != nil {
				return nil, fmt.Errorf("failed to parse private key: %w", parseErr)
			}
		} else {
			keyBytes, readErr := os.ReadFile(authSecret)
			if readErr != nil {
				if os.IsNotExist(readErr) {
					return nil, fmt.Errorf("SSH key file not found: %s", authSecret)
				}
				return nil, fmt.Errorf("failed to read SSH key file %s: %w", authSecret, readErr)
			}
			signer, parseErr = ssh.ParsePrivateKey(keyBytes)
			if parseErr != nil {
				return nil, fmt.Errorf("failed to parse private key from file %s: %w", authSecret, parseErr)
			}
		}

		if signer == nil {
			return nil, fmt.Errorf("failed to parse private key: empty signer")
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else {
		pwd := authSecret
		authMethods = append(authMethods, ssh.Password(pwd))
		authMethods = append(authMethods, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			answers := make([]string, len(questions))
			for i := range questions {
				answers[i] = pwd
			}
			return answers, nil
		}))
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         s.timeout,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, s.timeout)
	if err != nil {
		return nil, err
	}

	sshConn, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, err
	}

	client := ssh.NewClient(sshConn, chans, reqs)
	
	events.Publish(events.Event{
		Type:    "SSHConnected",
		Payload: addr,
	})

	return &realSSHClient{client: client}, nil
}

func (c *realSSHClient) RunCommand(ctx context.Context, cmd string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	out, err := session.CombinedOutput(cmd)
	return string(out), err
}

func (c *realSSHClient) Shell(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdin = stdin
	session.Stdout = stdout
	session.Stderr = stderr

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	fd := int(os.Stdin.Fd())
	width := 80
	height := 40

	if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		fd = int(tty.Fd())
		defer tty.Close()
	}

	if term.IsTerminal(fd) {
		w, h, err := term.GetSize(fd)
		if err == nil {
			width = w
			height = h
		}
		oldState, err := term.MakeRaw(fd)
		if err == nil {
			defer term.Restore(fd, oldState)
		}

		sigChan := make(chan os.Signal, 1)
		notifyWindowChange(sigChan)
		go func() {
			for range sigChan {
				if w, h, err := term.GetSize(fd); err == nil {
					_ = session.WindowChange(h, w)
				}
			}
		}()
		defer func() {
			signal.Stop(sigChan)
			close(sigChan)
		}()
	}

	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		return err
	}

	if err := session.Shell(); err != nil {
		return err
	}

	return session.Wait()
}

func (c *realSSHClient) Close() error {
	return c.client.Close()
}
