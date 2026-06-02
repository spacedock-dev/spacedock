// ABOUTME: unit table for EmitPythonJSON / escapeNonASCII — the ensure_ascii
// ABOUTME: parity escaping that makes native JSON byte-match Python json.dumps.
package claudeteam

import (
	"bytes"
	"strings"
	"testing"
)

// TestEscapeNonASCIIMatchesPythonEnsureASCII drives the escape helper over the
// chars json.dumps escapes that Go's encoder does not — DEL and every non-ASCII
// rune (BMP accents/dashes/arrows + a surrogate-pair astral emoji) — plus the
// chars both leave alone (plain ASCII, the encoder's own \t/\n escapes, and the
// HTML-significant < > & that ensure_ascii does NOT touch). Each input's native
// emission must byte-match the frozen json.dumps(ensure_ascii=True) output
// captured from python3 at retirement time, so the table pins the certified
// behavior without execing an interpreter. Inputs and wants are written with Go
// escapes so this source file stays pure ASCII; each want is the literal
// backslash-u form json.dumps emits.
func TestEscapeNonASCIIMatchesPythonEnsureASCII(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"plain", "\"plain\""},
		{"em—dash", "\"em\\u2014dash\""},                 // em-dash U+2014
		{"smart“q”", "\"smart\\u201cq\\u201d\""},         // curly quotes
		{"arrow→there", "\"arrow\\u2192there\""},         // rightwards arrow
		{"café", "\"caf\\u00e9\""},                       // accented e
		{"rocket\U0001f680", "\"rocket\\ud83d\\ude80\""}, // astral emoji -> surrogate pair
		{"del\x7fx", "\"del\\u007fx\""},                  // DEL: json.dumps escapes it, Go's encoder does not
		{"nbsp x", "\"nbsp\\u00a0x\""},                   // non-breaking space
		{"tab\tnl\n", "\"tab\\tnl\\n\""},                 // encoder already escapes these; helper leaves them
		{"less<more>amp&", "\"less<more>amp&\""},         // HTML chars stay raw under ensure_ascii
	}
	for _, tc := range cases {
		got := emitJSONString(t, tc.in)
		if got != tc.want {
			t.Errorf("emit %q\n  native = %s\n  want   = %s", tc.in, got, tc.want)
		}
	}
}

// emitJSONString runs a bare string value through EmitPythonJSON and returns the
// emitted JSON with the trailing newline trimmed, so the comparison matches the
// json.dumps form (no indent on a scalar, no trailing newline).
func emitJSONString(t *testing.T, s string) string {
	t.Helper()
	var buf bytes.Buffer
	if EmitPythonJSON(&buf, s) != 0 {
		t.Fatalf("EmitPythonJSON returned non-zero for %q", s)
	}
	return strings.TrimRight(buf.String(), "\n")
}
