package ssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/devlopersabbir/vpcm/internal/events"
	"golang.org/x/crypto/ssh"
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
	var authMethod ssh.AuthMethod

	if authType == "key" {
		signer, err := ssh.ParsePrivateKey([]byte(authSecret))
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		authMethod = ssh.Password(authSecret)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For v0.0.1 simplicity
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

	if err := session.RequestPty("xterm-256color", 80, 40, modes); err != nil {
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
