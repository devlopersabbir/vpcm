package ssh

import (
	"context"
	"io"
)

type Client interface {
	RunCommand(ctx context.Context, cmd string) (string, error)
	Shell(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error
	Close() error
}

type SSHService interface {
	Connect(ctx context.Context, host string, port int, username string, authType string, authSecret string) (Client, error)
}
