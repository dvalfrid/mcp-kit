// Package secret has the Secret type: a string that never shows its value
// when printed, logged or marshalled. Passwords, tokens and key
// passphrases are held as a Secret from the moment they are read, so a
// stray log line or error message can't leak them.
package secret

import "log/slog"

const redacted = "[redacted]"

// Secret is a string that never shows its value when printed, logged or
// marshalled. Use Reveal at the single place the value is needed.
type Secret string

// Reveal returns the secret's value.
func (s Secret) Reveal() string { return string(s) }

// String implements fmt.Stringer.
func (s Secret) String() string { return redacted }

// GoString implements fmt.GoStringer (%#v).
func (s Secret) GoString() string { return redacted }

// LogValue implements slog.LogValuer.
func (s Secret) LogValue() slog.Value { return slog.StringValue(redacted) }

// MarshalText implements encoding.TextMarshaler, so JSON and similar
// encoders print the placeholder.
func (s Secret) MarshalText() ([]byte, error) { return []byte(redacted), nil }
