package dotenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), ".env")
	body := "# comment\n\nA=1\nexport B = \"two words\"\nC='x'\nD=\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"A": "1", "B": "two words", "C": "x", "D": ""}
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
}

func TestReadMissingFile(t *testing.T) {
	t.Parallel()
	got, err := Read(filepath.Join(t.TempDir(), "none"))
	if err != nil || len(got) != 0 {
		t.Errorf("got %v, %v", got, err)
	}
}

// Guardrail: a bad line is reported by number, never echoed: it could
// hold a password.
func TestReadNeverEchoesALine(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("hunter2\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := Read(path)
	if err == nil || strings.Contains(err.Error(), "hunter2") || !strings.Contains(err.Error(), "line 1") {
		t.Errorf("err = %v", err)
	}
}

func TestChain(t *testing.T) {
	t.Parallel()
	get := Chain(FromMap(map[string]string{"A": "first"}), FromMap(map[string]string{"A": "second", "B": "b"}))
	if v, ok := get("A"); !ok || v != "first" {
		t.Errorf("A = %q, %v", v, ok)
	}
	if v, ok := get("B"); !ok || v != "b" {
		t.Errorf("B = %q, %v", v, ok)
	}
	if _, ok := get("C"); ok {
		t.Error("C found")
	}
}

func TestLookupPrefersTheEnvironment(t *testing.T) {
	t.Setenv("MCP_KIT_TEST_LOOKUP", "env")
	get := Lookup(map[string]string{"MCP_KIT_TEST_LOOKUP": "file", "ONLY_FILE": "f"})
	if v, _ := get("MCP_KIT_TEST_LOOKUP"); v != "env" {
		t.Errorf("got %q", v)
	}
	if v, _ := get("ONLY_FILE"); v != "f" {
		t.Errorf("got %q", v)
	}
}
