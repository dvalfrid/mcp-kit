// Package dotenv reads settings the way every server here does: the
// process environment first, then a .env file, then whatever else the
// server chains after them (for example its own settings file).
package dotenv

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Read parses a .env file: KEY=VALUE lines, '#' comments, an optional
// "export " prefix, optional single or double quotes around the value. A
// missing file is not an error: the settings may come from the
// environment instead. An error never echoes a line, which could hold a
// password.
func Read(path string) (map[string]string, error) {
	f, err := os.Open(path) //nolint:gosec // G304: the path is the operator's own config file
	if errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer func() { _ = f.Close() }() // read-only file

	vars := map[string]string{}
	sc := bufio.NewScanner(f)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		key = strings.TrimSpace(strings.TrimPrefix(key, "export "))
		if !ok || key == "" {
			return nil, fmt.Errorf("%s line %d: expected KEY=VALUE", path, n)
		}
		vars[key] = unquote(strings.TrimSpace(value))
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return vars, nil
}

func unquote(v string) string {
	if len(v) >= 2 && (v[0] == '"' || v[0] == '\'') && v[len(v)-1] == v[0] {
		return v[1 : len(v)-1]
	}
	return v
}

// Lookup returns a lookup function where the process environment wins
// over values from a .env file.
func Lookup(dotenv map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		if v, ok := os.LookupEnv(key); ok {
			return v, true
		}
		v, ok := dotenv[key]
		return v, ok
	}
}

// FromMap is a lookup function over a fixed map.
func FromMap(m map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}

// Chain returns a lookup function that asks each lookup in turn and
// returns the first value found.
func Chain(lookups ...func(string) (string, bool)) func(string) (string, bool) {
	return func(key string) (string, bool) {
		for _, l := range lookups {
			if v, ok := l(key); ok {
				return v, true
			}
		}
		return "", false
	}
}
