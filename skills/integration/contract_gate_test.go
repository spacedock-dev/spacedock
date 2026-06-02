// ABOUTME: AC-2 bracketing test over the vendored FO Startup contract — the
// ABOUTME: embedded contract-range literal must bracket the binary's CONTRACT_VERSION.
package integration

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// foSharedCore reads the vendored first-officer shared core contract text.
func foSharedCore(t *testing.T) string {
	t.Helper()
	p := filepath.Join(skillsRoot(t), "first-officer", "references", "first-officer-shared-core.md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read FO shared core: %v", err)
	}
	return string(b)
}

// embeddedRangeRe matches the half-open contract range literal embedded in the
// Startup step-0 prose (e.g. `>=1,<2`).
var embeddedRangeRe = regexp.MustCompile(`>=\s*\d+\s*,\s*<\s*\d+`)

// TestStartupEmbeddedRangeBracketsContractVersion locks the embedded-range
// bracketing invariant: the range literal embedded in the FO Startup prose
// brackets CONTRACT_VERSION. Both surfaces live in spacedock-v1, so a single go
// test closes the 4th-source-of-truth drift (the FO contract embeds its own
// expected range as a literal).
//
// Oracle: contract.ParseRange (the half-open-range parser) + the compiled
// contract.CONTRACT_VERSION constant. This is NOT bare prose-grep — it does not
// assert "the prose says X"; it parses the embedded literal and checks it
// brackets the binary's real contract version, catching FO/binary range drift.
// The contract-gate ORDERING behavior (gate runs before discover/boot) is owned
// behaviorally by internal/contract/gate_test.go, which drives a real spacedock
// stub --version and observes discover invoked 0×/1×.
func TestStartupEmbeddedRangeBracketsContractVersion(t *testing.T) {
	startup := sectionAfter(foSharedCore(t), "## Startup")
	raw := embeddedRangeRe.FindString(startup)
	if raw == "" {
		t.Fatalf("Startup section has no embedded contract range literal (>=N,<M)")
	}
	lo, hi, err := contract.ParseRange(raw)
	if err != nil {
		t.Fatalf("embedded range %q does not parse: %v", raw, err)
	}
	if !(lo <= contract.CONTRACT_VERSION && contract.CONTRACT_VERSION < hi) {
		t.Fatalf("embedded Startup range %s does not bracket CONTRACT_VERSION=%d", raw, contract.CONTRACT_VERSION)
	}
	// Guard against a stray literal: the embedded range must be a single
	// occurrence in the Startup section (one source of truth, not several).
	if got := len(embeddedRangeRe.FindAllString(startup, -1)); got != 1 {
		t.Fatalf("Startup section has %d embedded range literals, want exactly 1", got)
	}
}

// startupStep1 returns the text of Startup step 1 (the contract version gate),
// scoped from its `1. **Contract version gate` opener up to the start of the
// `2.` step. The two-class abort split lives entirely inside this step, so the
// AC-1/AC-2 string relationships are asserted against this span, not the whole
// Startup section (steps 2-7 are unrelated and must not leak into the checks).
func startupStep1(t *testing.T) string {
	t.Helper()
	startup := sectionAfter(foSharedCore(t), "## Startup")
	const opener = "1. **Contract version gate"
	i := strings.Index(startup, opener)
	if i < 0 {
		t.Fatalf("Startup section has no step-1 opener %q", opener)
	}
	rest := startup[i:]
	// Step 1 ends where step 2 begins. Numbered steps open at column 0, so a
	// newline immediately followed by `2.` is the boundary.
	if j := strings.Index(rest, "\n2."); j >= 0 {
		return rest[:j]
	}
	return rest
}

// TestStartupAbortSplitsByBinaryPresence locks the binary-absent FO-startup
// hardening (entity binary-absent-fo-bootstrap):
//
//   - AC-1 — the binary-absent / non-executable abort class carries both
//     runnable install lines verbatim (released brew lane + source go-build lane)
//     and routes to NO `spacedock doctor` within that class, since `doctor` is the
//     same missing binary.
//   - AC-2 — `spacedock doctor` survives, attached to the binary-present /
//     contract-out-of-range class (Class B), so the fix is a per-class split and
//     not a blanket deletion of the doctor route.
//
// Oracle: the abort-clause text of Startup step 1, isolated via startupStep1 and
// split by the two class markers. This is the legitimate doc-as-deliverable case
// recorded in the entity's AC-1 framing note: the binary is absent by definition
// of the failure mode, so no `spacedock` command can run to emit the install
// hint — the contract prose the FO loads IS the only artifact present at failure.
// The check is therefore a presence test (both install lines) plus a
// banned-token absence test (no doctor route inside Class A) over the real file,
// which is proof at the claim's own level, not a behavioral claim a code gate
// could enforce here.
func TestStartupAbortSplitsByBinaryPresence(t *testing.T) {
	step1 := startupStep1(t)

	const (
		classAMarker = "**Binary absent or non-executable**"
		classBMarker = "**Binary present but contract out of range**"
		brewLine     = "brew install spacedock-dev/homebrew-tap/spacedock"
		goBuildLine  = "go build -o spacedock ./cmd/spacedock"
		doctor       = "spacedock doctor"
	)

	aStart := strings.Index(step1, classAMarker)
	if aStart < 0 {
		t.Fatalf("step 1 has no binary-absent class marker %q", classAMarker)
	}
	bStart := strings.Index(step1, classBMarker)
	if bStart < 0 {
		t.Fatalf("step 1 has no binary-present-out-of-range class marker %q", classBMarker)
	}
	if !(aStart < bStart) {
		t.Fatalf("class markers out of order: Class A at %d, Class B at %d (Class A must precede Class B)", aStart, bStart)
	}

	// Class A spans from its marker up to where Class B begins.
	classA := step1[aStart:bStart]
	// Class B spans from its marker to the end of step 1.
	classB := step1[bStart:]

	// AC-1: both install lines present verbatim, inside the binary-absent class.
	for _, line := range []string{brewLine, goBuildLine} {
		if !strings.Contains(classA, line) {
			t.Errorf("binary-absent class is missing runnable install line %q", line)
		}
	}
	// AC-1: Class A must carry the no-doctor PROHIBITION and must not route to
	// doctor. The prohibition string itself names `spacedock doctor` once; any
	// further mention would be a route. Require the prohibition is present (so
	// deleting the guidance fails, not skips) AND that the lone doctor mention is
	// exactly that prohibition (so an added route fails). `doctor` is the same
	// missing binary, so it can never be a live remedy in Class A.
	const doctorProhibition = "Do NOT run `spacedock doctor`"
	if !strings.Contains(classA, doctorProhibition) {
		t.Errorf("binary-absent class is missing the no-doctor prohibition %q — Class A must carry the guidance that `doctor` is the same missing binary", doctorProhibition)
	}
	if n := strings.Count(classA, doctor); n != 1 {
		t.Errorf("binary-absent class mentions %q %d time(s); want exactly 1 (the prohibition only) — any extra mention is a route, and `doctor` must not be a route in Class A", doctor, n)
	}

	// AC-2: doctor survives as a LIVE ROUTE in the binary-present class. A bare
	// mention is too weak — a disclaimer ("Historically we suggested spacedock
	// doctor but no longer.") would satisfy it. Assert the active routing phrasing
	// so a gutted route phrased as a disclaimer fails: the class must instruct to
	// `run `spacedock doctor`` for `the per-class remedy`.
	const (
		doctorRoute  = "run `spacedock doctor`"
		perClassVerb = "for the per-class remedy"
	)
	if !strings.Contains(classB, doctorRoute) || !strings.Contains(classB, perClassVerb) {
		t.Errorf("binary-present-out-of-range class no longer carries the live doctor route (%q + %q) — a blanket deletion or a disclaimer-only mention would hit this", doctorRoute, perClassVerb)
	}
}

// TestStartupGateGuidanceHasSingleProseSource locks AC-3: the startup-gate abort
// guidance lives in exactly ONE prose file (first-officer-shared-core.md), and
// agents/first-officer.md continues to delegate to the shared core rather than
// restating the gate. A second prose copy would let the two surfaces drift, so a
// single-file edit would silently fail to harden the mirror.
//
// Oracle: a walk of the skills/ and agents/ trees for the gate markers over
// markdown prose files only. The integration .go test files legitimately mention
// `spacedock doctor` in a comment (marketplace_manifest_test.go) and assert it
// here (this file); those are Go sources, not a prose mirror of the gate, so the
// walk is scoped to .md files. This is NOT a bare prose-grep AC: it does not
// assert "the prose says X"; it asserts a structural single-source invariant —
// that the gate guidance is NOT duplicated across prose files — which a code
// change can violate and this test would catch.
func TestStartupGateGuidanceHasSingleProseSource(t *testing.T) {
	root := repoRoot(t)
	// Markers unique to the startup-gate abort prose. A second .md file carrying
	// any of these would be a drift-prone mirror of the single source of truth.
	markers := []string{
		"Contract version gate",
		"per-class remedy",
		"spacedock doctor",
	}

	var proseSources []string
	for _, tree := range []string{"skills", "agents"} {
		base := filepath.Join(root, tree)
		err := filepath.Walk(base, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || filepath.Ext(p) != ".md" {
				return nil
			}
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			text := string(b)
			for _, m := range markers {
				if strings.Contains(text, m) {
					rel, _ := filepath.Rel(root, p)
					proseSources = append(proseSources, rel)
					return nil
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", tree, err)
		}
	}

	const sole = "skills/first-officer/references/first-officer-shared-core.md"
	if len(proseSources) != 1 || proseSources[0] != sole {
		t.Errorf("startup-gate prose has sources %v, want exactly [%s] (single source of truth)", proseSources, sole)
	}

	// agents/first-officer.md must still delegate, not mirror the gate prose.
	agentPath := filepath.Join(root, "agents", "first-officer.md")
	b, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("read %s: %v", agentPath, err)
	}
	agent := string(b)
	const delegation = "begin the Startup procedure from the shared core"
	if !strings.Contains(agent, delegation) {
		t.Errorf("agents/first-officer.md missing delegation line %q", delegation)
	}
	for _, m := range []string{"per-class remedy", "spacedock doctor"} {
		if strings.Contains(agent, m) {
			t.Errorf("agents/first-officer.md now mirrors gate prose marker %q (must delegate, not restate)", m)
		}
	}
}
