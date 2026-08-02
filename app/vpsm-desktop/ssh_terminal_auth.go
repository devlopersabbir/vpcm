package main

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func buildSSHAuthMethods(params SSHConnectionParams) []ssh.AuthMethod {
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

		if signer, err := tryParseKeyFile(relPath, params.AuthSecret); err == nil {
			addSigner(signer)
			return
		}

		if home, err := os.UserHomeDir(); err == nil {
			candidates := []string{
				filepath.Join(home, ".ssh", relPath),
				filepath.Join(home, relPath),
				filepath.Join(home, "Downloads", relPath),
				filepath.Join(home, "Desktop", relPath),
			}
			for _, cand := range candidates {
				if signer, err := tryParseKeyFile(cand, params.AuthSecret); err == nil {
					addSigner(signer)
					return
				}
			}
		}
	}

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

	if params.AuthSecret != "" && params.AuthType != "password" {
		addKeyPathSigner(params.AuthSecret)
		if signer, err := ssh.ParsePrivateKey([]byte(params.AuthSecret)); err == nil {
			addSigner(signer)
		}
		pwd := params.AuthSecret
		authMethods = append(authMethods, ssh.Password(pwd))
	}

	if homeDir, err := os.UserHomeDir(); err == nil {
		sshDir := filepath.Join(homeDir, ".ssh")
		if entries, err := os.ReadDir(sshDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				name := entry.Name()
				if name == "known_hosts" || name == "known_hosts.old" || name == "config" || strings.HasSuffix(name, ".pub") {
					continue
				}
				keyPath := filepath.Join(sshDir, name)
				if signer, err := tryParseKeyFile(keyPath, params.AuthSecret); err == nil {
					addSigner(signer)
				}
			}
		}
	}

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

	return authMethods
}
