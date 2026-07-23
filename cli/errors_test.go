package cli

import (
	"errors"
	"os"

	"testing"
)

func TestPrintError(t *testing.T) {
	// Ensure PrintError does not panic with nil or non-nil errors
	PrintError(nil)

	errKeyNotFound := errors.New("SSH key connection failed: SSH key file not found: christian.pem")
	PrintError(errKeyNotFound)

	errKeyParse := errors.New("failed to parse private key: ssh: no key found")
	PrintError(errKeyParse)

	errArgs := errors.New("accepts 1 arg(s), received 0")
	PrintError(errArgs)
}

func TestSSHKeyFileNotFound(t *testing.T) {
	oldIdentity := identityFile
	defer func() { identityFile = oldIdentity }()

	identityFile = "/nonexistent/path/christian.pem"
	_, statErr := os.Stat(identityFile)
	if statErr == nil {
		t.Fatal("expected non-existent file path")
	}
}
