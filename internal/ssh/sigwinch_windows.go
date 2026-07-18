//go:build windows

package ssh

import "os"

func notifyWindowChange(sigChan chan<- os.Signal) {
	// SIGWINCH is not supported on Windows
}
