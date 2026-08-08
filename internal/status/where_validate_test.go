// ABOUTME: AC-1 (#314) --where robustness — Check A (compound-in-one-string
// ABOUTME: operator-count syntax guard) and Check B (unknown-field-name guard).
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// compoundWhereShapes are the four ways a single --where argument can carry two
// clauses instead of one. All four must be rejected: today all four return
// 156-of-156 rows at exit 0 on a live corpus (silent match-all for three of
// them, silent empty-result for the "sprint=A sprint-readiness=ready" case).
// Only a multi-operator syntax check catches all four — a shape like
// "a!=1 b!=2" cuts to the legitimate field name "a" with a garbage value, so
// field-name validation alone (Check B) does not see anything wrong with it.
var compoundWhereShapes = []struct {
	name string
	arg  string
}{
	{"eq-then-neq", "a=1 b!=2"},
	{"neq-then-neq", "a!=1 b!=2"},
	{"eq-then-eq", "a=1 b=2"},
	{"neq-then-eq", "a!=1 b=2"},
}

// TestWhereCompoundShapesRejected locks Check A: every compound-in-one-string
// shape exits 1 with a message naming the repeated-flag fix (repeat --where),
// not a silent match-all/match-none at exit 0.
func TestWhereCompoundShapesRejected(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	for _, tc := range compoundWhereShapes {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--where", tc.arg)
			if code != 1 {
				t.Fatalf("--where %q exit=%d (out=%q err=%q), want 1", tc.arg, code, out, errOut)
			}
			if out != "" {
				t.Fatalf("--where %q: stdout must be empty on rejection, got %q", tc.arg, out)
			}
			if !strings.Contains(errOut, "repeat --where") {
				t.Fatalf("--where %q: stderr = %q, want it to name the repeated-flag fix", tc.arg, errOut)
			}
		})
	}
}

// TestWhereUnknownFieldRejected locks Check B: a field name absent from both the
// scanned corpus and the canonical schema is a loud error listing the known
// fields, not a silent match-all (!=) or match-none (=).
func TestWhereUnknownFieldRejected(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	cases := []string{"nosuchfield!=x", "spint=v"}
	for _, arg := range cases {
		arg := arg
		t.Run(arg, func(t *testing.T) {
			out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--where", arg)
			if code != 1 {
				t.Fatalf("--where %q exit=%d (out=%q err=%q), want 1", arg, code, out, errOut)
			}
			if out != "" {
				t.Fatalf("--where %q: stdout must be empty on rejection, got %q", arg, out)
			}
			if !strings.Contains(errOut, "unknown field") {
				t.Fatalf("--where %q: stderr = %q, want it to say unknown field", arg, errOut)
			}
			// The schema half of the union must be listed even though this
			// active-only corpus has no entity carrying these canonical keys.
			for _, want := range []string{"verdict", "archived", "id", "status", "title"} {
				if !strings.Contains(errOut, want) {
					t.Fatalf("--where %q: stderr = %q, want known-fields list to include %q", arg, errOut, want)
				}
			}
		})
	}
}

// TestWhereNoOverRejection locks the counter-baseline: a known-but-unpopulated
// canonical field (verdict is schema-declared but absent from every active
// entity in this fixture) must not error. Losing the schema half of the union
// would make this exit 1, breaking a legitimate query on a fresh active read.
func TestWhereNoOverRejection(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	_, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--where", "verdict!=")
	if code != 0 {
		t.Fatalf("--where 'verdict!=' exit=%d (err=%q), want 0 (no over-rejection)", code, errOut)
	}
}

// TestWhereDerivedNamesAccepted locks that the two computed field names
// (gate-readiness, next-suppressed-by) are accepted even though their materializers only
// populate them in e.fields when something in the same invocation already
// references them — the static derived list must not depend on that ordering.
func TestWhereDerivedNamesAccepted(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "suppress-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	cases := []string{"next-suppressed-by = concurrency-full", "gate-readiness=validating"}
	for _, arg := range cases {
		arg := arg
		t.Run(arg, func(t *testing.T) {
			_, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--where", arg)
			if code != 0 {
				t.Fatalf("--where %q exit=%d (err=%q), want 0", arg, code, errOut)
			}
		})
	}
}

func TestWhereRejectsRemovedEligibilityProjection(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "suppress-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	_, errOut, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--where", "gate-eligible=true")
	if code != 1 || !strings.Contains(errOut, "unknown field \"gate-eligible\"") {
		t.Fatalf("removed gate-eligible filter exit=%d stderr=%q", code, errOut)
	}
}

// TestWhereEmptyWorkflowSkipsFieldValidation locks the zero-entity guard: with
// nothing scanned, Check B must not fire (every result is empty either way), so
// a fresh workflow's first --where does not error before any entity exists to
// populate the corpus half of the known-field union.
func TestWhereEmptyWorkflowSkipsFieldValidation(t *testing.T) {
	root := t.TempDir()
	env := pinnedEnv(t)

	_, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--where", "sprint=X")
	if code != 0 {
		t.Fatalf("--where on empty workflow exit=%d (err=%q), want 0", code, errOut)
	}
}
