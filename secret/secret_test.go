package secret

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

// Guardrail: a Secret never shows its value, however it is printed.
func TestSecretNeverShows(t *testing.T) {
	t.Parallel()
	s := Secret("hunter2")
	var logs bytes.Buffer
	slog.New(slog.NewTextHandler(&logs, nil)).Info("x", "password", s)
	js, err := json.Marshal(struct{ P Secret }{s})
	if err != nil {
		t.Fatal(err)
	}
	for name, out := range map[string]string{
		"%v": fmt.Sprintf("%v", s), "Sprint": fmt.Sprint(s), "%#v": fmt.Sprintf("%#v", s),
		"%+v struct": fmt.Sprintf("%+v", struct{ P Secret }{s}), "slog": logs.String(), "json": string(js),
	} {
		if strings.Contains(out, "hunter2") {
			t.Errorf("%s shows the secret: %q", name, out)
		}
	}
	if s.Reveal() != "hunter2" {
		t.Error("Reveal")
	}
}
