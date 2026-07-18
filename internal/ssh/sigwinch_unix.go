//go:build !windows

package ssh

import (
	"os"
	"os/signal"
	"syscall"
)

func notifyWindowChange(sigChan chan<- os.Signal) {
	signal.Notify(sigChan, syscall.SIGWINCH)
}
