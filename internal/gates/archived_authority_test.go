package gates

import (
	"path/filepath"
	"strings"
	"testing"
)

// The fixture is a real retained provider Result. Its annotations carry
// "target", "kind", and "body" — fields the Annotation struct no longer models —
// alongside "selectors", which it never modeled. Frozen authority must keep
// decoding and verifying through that tolerance.
func TestArchivedProviderResultDecodesWithoutAnnotationProseFields(t *testing.T) {
	raw := mustReadBytes(t, filepath.Join("testdata", "archived-provider-result.json"))
	for _, unmodeled := range []string{`"target":`, `"kind":`, `"body":`, `"selectors":`} {
		if !strings.Contains(string(raw), unmodeled) {
			t.Fatalf("fixture no longer carries %s; it cannot prove decode tolerance", unmodeled)
		}
	}
	result, err := decodeProviderResult(raw)
	if err != nil {
		t.Fatalf("decode archived Result: %v", err)
	}
	if len(result.Annotations) != 2 || result.Resolution.Decision != "revise" ||
		len(result.Resolution.Includes) != 2 {
		t.Fatalf("archived Result lost modeled authority: %#v", result)
	}
	if err := verifyProviderResolution(result); err != nil {
		t.Fatalf("verify archived Resolution: %v", err)
	}
}
