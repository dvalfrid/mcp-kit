package keystore

import (
	"errors"
	"testing"

	"github.com/zalando/go-keyring"
)

var _ Store = OS{}
var _ Store = Memory{}

func TestMemory(t *testing.T) {
	t.Parallel()
	m := Memory{}
	if _, err := m.Get("k"); !errors.Is(err, ErrNotFound) {
		t.Errorf("empty: %v", err)
	}
	if err := m.Set("k", "v"); err != nil {
		t.Fatal(err)
	}
	if v, err := m.Get("k"); err != nil || v.Reveal() != "v" {
		t.Errorf("got %v, %v", v, err)
	}
	if err := m.Delete("k"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Get("k"); !errors.Is(err, ErrNotFound) {
		t.Errorf("after delete: %v", err)
	}
}

// The OS store against go-keyring's in-memory mock: values are kept per
// service, so one server never reads another's secrets.
func TestOSKeepsServicesApart(t *testing.T) {
	keyring.MockInit()
	a, b := OS{Service: "a-mcp"}, OS{Service: "b-mcp"}
	if err := a.Set("token", "secret-a"); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Get("token"); !errors.Is(err, ErrNotFound) {
		t.Errorf("another service read it: %v", err)
	}
	if v, err := a.Get("token"); err != nil || v.Reveal() != "secret-a" {
		t.Errorf("got %v, %v", v, err)
	}
	if err := a.Delete("token"); err != nil {
		t.Fatal(err)
	}
	if err := a.Delete("token"); err != nil {
		t.Errorf("deleting nothing: %v", err)
	}
}

func TestOSNeedsAService(t *testing.T) {
	keyring.MockInit()
	defer func() {
		if recover() == nil {
			t.Error("no panic for an empty service name")
		}
	}()
	_, _ = OS{}.Get("k")
}
