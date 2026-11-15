// Package keystore keeps a server's secrets in the operating system's
// credential store: Windows Credential Manager, the macOS Keychain, or the
// Secret Service on Linux (through github.com/zalando/go-keyring, no cgo).
// Each server uses its own service name, so one server never reads
// another's secrets.
//
// Known weakness on macOS: go-keyring creates items through
// /usr/bin/security, so any program running as the same user can read the
// item without a prompt. That is no worse than a file only the user can
// read; each server states it in its SECURITY.md.
package keystore

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"

	"github.com/dvalfrid/mcp-kit/secret"
)

// ErrNotFound means nothing is stored under the key.
var ErrNotFound = errors.New("nothing stored in the credential store")

// Store keeps one secret per key.
type Store interface {
	Get(key string) (secret.Secret, error)
	Set(key string, value secret.Secret) error
	Delete(key string) error
}

// OS is the operating system's credential store, under one service name.
type OS struct {
	// Service names the server, e.g. "proton-calendar-mcp". Required.
	Service string
}

// Get returns the stored secret, or ErrNotFound.
func (o OS) Get(key string) (secret.Secret, error) {
	v, err := keyring.Get(o.service(), key)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("read the credential store: %w", err)
	}
	return secret.Secret(v), nil
}

// Set stores the secret, replacing any earlier one.
func (o OS) Set(key string, value secret.Secret) error {
	if err := keyring.Set(o.service(), key, value.Reveal()); err != nil {
		return fmt.Errorf("write the credential store: %w", err)
	}
	return nil
}

// Delete removes the secret. Nothing stored is not an error.
func (o OS) Delete(key string) error {
	err := keyring.Delete(o.service(), key)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return fmt.Errorf("clear the credential store: %w", err)
	}
	return nil
}

// DeleteAll removes every secret stored under the service: for
// uninstalling.
func (o OS) DeleteAll() error {
	if err := keyring.DeleteAll(o.service()); err != nil {
		return fmt.Errorf("clear the credential store: %w", err)
	}
	return nil
}

func (o OS) service() string {
	if o.Service == "" {
		// A programming error: secrets under an empty name would be
		// shared by every server that made the same mistake.
		panic("keystore: OS.Service is empty")
	}
	return o.Service
}

// Memory is a Store for tests.
type Memory map[string]secret.Secret

// Get returns the secret, or ErrNotFound.
func (m Memory) Get(key string) (secret.Secret, error) {
	v, ok := m[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

// Set stores the secret.
func (m Memory) Set(key string, value secret.Secret) error {
	m[key] = value
	return nil
}

// Delete removes the secret.
func (m Memory) Delete(key string) error {
	delete(m, key)
	return nil
}
