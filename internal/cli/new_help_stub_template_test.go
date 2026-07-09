// ABOUTME: AC-2 — `spacedock new --help` surfaces the copyable entity stub
// ABOUTME: template and the piped example, so a filing FO is handed the shape.
package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// TestNewHelpShowsStubTemplate asserts `spacedock new --help` renders the shared
// stub template (frontmatter skeleton + id-omitted + commit next-step) and the
// pipe-able printf example — the surface a filing FO consults before filing. It
// drives the real help renderer, not a source substring.
func TestNewHelpShowsStubTemplate(t *testing.T) {
	var stdout bytes.Buffer
	Run([]string{"new", "--help"}, &stdout, &bytes.Buffer{})
	out := stdout.String()

	for _, want := range []string{"title:", "status:", "id OMITTED", "state commit"} {
		if !strings.Contains(out, want) {
			t.Fatalf("new --help missing stub-template token %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, status.EntityStubTemplate) {
		t.Fatalf("new --help does not embed the shared EntityStubTemplate verbatim:\n%s", out)
	}
	// The piped example hands the FO a one-shot fill-and-file command.
	if !strings.Contains(out, "| spacedock new") {
		t.Fatalf("new --help missing the piped fill-and-file example:\n%s", out)
	}
}
