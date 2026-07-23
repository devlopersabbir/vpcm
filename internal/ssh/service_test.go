package ssh

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func generateTestPrivateKey(t *testing.T) ([]byte, string) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}

	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	keyBytes := pem.EncodeToMemory(pemBlock)

	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "id_rsa")
	if err := os.WriteFile(keyPath, keyBytes, 0600); err != nil {
		t.Fatalf("failed to write key file: %v", err)
	}

	return keyBytes, keyPath
}

func TestConnectKeyAuthParsing(t *testing.T) {
	keyBytes, keyPath := generateTestPrivateKey(t)
	svc := NewService(1 * time.Second)

	// Test invalid key file path
	_, err := svc.Connect(context.Background(), "127.0.0.1", 22222, "root", "keyfile", "/nonexistent/path/key.pem")
	if err == nil {
		t.Fatal("expected error for non-existent key file path, got nil")
	}
	if !strings.Contains(err.Error(), "SSH key file not found") {
		t.Fatalf("expected error to contain 'SSH key file not found', got: %v", err)
	}

	// Test raw key string parsing (should fail at connection dial step, not key parse step)
	_, err = svc.Connect(context.Background(), "127.0.0.1", 22222, "root", "key", string(keyBytes))
	if err == nil {
		// Connection will fail because 127.0.0.1:22222 is not listening, but error shouldn't be "failed to parse private key"
	} else if err.Error() == "failed to parse private key" {
		t.Fatalf("unexpected key parse error for valid key string: %v", err)
	}

	// Test keyfile path parsing (should fail at connection dial step, not key parse step)
	_, err = svc.Connect(context.Background(), "127.0.0.1", 22222, "root", "keyfile", keyPath)
	if err != nil && err.Error() == "failed to parse private key" {
		t.Fatalf("unexpected key parse error for valid key path: %v", err)
	}
}
