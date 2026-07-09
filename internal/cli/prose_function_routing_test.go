// ABOUTME: AC-1/AC-4 — bind each prose-function core's notation to the binary's command
// ABOUTME: routing: shipped backtick verbs route, un-shipped guillemet targets are rejected.
package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// proseFunctionCores are the restructured FO contract cores carrying the «fn» prose-function
// declarations. The walk reads them by repo-relative path (the same filesystem-oracle idiom
// boot_resident_closure_test.go uses), so the check drives the on-disk contract, not a copy.
var proseFunctionCores = []string{
	filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md"),
	filepath.Join("skills", "fo-dispatch-core", "SKILL.md"),
	filepath.Join("skills", "fo-merge-core", "SKILL.md"),
}

// migrationTargetRe matches a prose-function declaration's `→` migration-target line and
// captures (1) the notation word — `shipped` (a backtick, a live verb) or `prose` (a
// guillemet, hand-followed) — and (2) the FIRST backticked `spacedock …` verb token on the
// line, if any. A shipped line naming `spacedock state ready` yields ("shipped",
// "spacedock state ready"); a prose line naming `spacedock dispatch next-action`
// line yields ("prose", "spacedock dispatch next-action"); a bare `→ **prose**` line with no
// `spacedock …` token (the `«feedback.route»` skill-is-the-body form) yields ("prose", "").
var migrationTargetRe = regexp.MustCompile("→ \\*\\*(shipped|prose)\\*\\*[^\\n]*?`(spacedock [^`]+)`")

// migrationBareRe matches a `→ **shipped**` / `→ **prose**` line that carries NO backticked
// `spacedock …` token, so a declaration whose notation says "shipped" but names no verb is
// caught as a malformed shipped-target rather than silently skipped.
var migrationBareRe = regexp.MustCompile("→ \\*\\*(shipped|prose)\\*\\*")

// verbBinding is one extracted migration target: the verb argv the contract names, the
// notation (shipped → must route; prose → must be un-routed when a target is named), and the
// source line for a useful failure message.
type verbBinding struct {
	notation string   // "shipped" or "prose"
	argv     []string // the verb tokens after "spacedock", e.g. {"state","ready"}
	verb     string   // the full `spacedock …` token, for messages
	line     string   // the source line
}

// extractVerbBindings parses one core's text and returns every prose-function migration-target
// binding it declares. It scopes extraction to the `→` migration-target LINES of the
// declaration blocks — NOT every inline backtick in the prose — so the deterministic extract
// sub-calls a judgment body names (`status --read --checklist`) and the negated mention of a
// never-shipping binary do not masquerade as migration targets. A `shipped` line MUST carry a
// verb token (malformed otherwise); a `prose` line carries one only when it names a future
// `becomes` target (the pure-judgment bodies name none).
func extractVerbBindings(body string) ([]verbBinding, error) {
	var out []verbBinding
	for _, line := range strings.Split(body, "\n") {
		if !migrationBareRe.MatchString(line) {
			continue
		}
		m := migrationTargetRe.FindStringSubmatch(line)
		if m == nil {
			// A `→` line with a notation word but no `spacedock …` token: legitimate for a
			// pure-judgment `prose` body (the skill is the body); a `shipped` line that names
			// no verb is malformed.
			if strings.Contains(line, "**shipped**") {
				return out, malformedShippedLine(line)
			}
			continue
		}
		notation, verb := m[1], m[2]
		fields := strings.Fields(verb)
		out = append(out, verbBinding{
			notation: notation,
			argv:     fields[1:], // drop the leading "spacedock"
			verb:     verb,
			line:     strings.TrimSpace(line),
		})
	}
	return out, nil
}

// malformedShippedLine is returned when a `→ **shipped**` line carries no `spacedock …` verb
// token, so the parser surfaces the defect instead of silently passing.
type malformedShippedLine string

func (m malformedShippedLine) Error() string {
	return "→ **shipped** line names no `spacedock …` verb token: " + string(m)
}

// verbRejected is the routing ORACLE: it invokes the verb's argv against the compiled root
// (run → newRootCommand) from a workflow-free tempdir and reports whether the routing REJECTED
// it. Rejection is recognized across ALL THREE marker forms at their three distinct sites:
// `unknown command: {cmd}` at the cobra root (an absent top-level command, e.g. `gate`);
// `unknown subcommand` at a parent's `switch args[0]` (e.g. `merge teleport`, `state warp`);
// and `unknown dispatch subcommand` from dispatch's own handler. A shipped verb routes into its
// handler (which may then complain about a missing flag/arg — that is acceptance, not
// rejection). This is the function the production check AND the RED control both drive, so
// mutating it (stubbing it to never reject) reds the control — what makes the control
// load-bearing rather than a re-implementation.
func verbRejected(t *testing.T, argv []string) bool {
	t.Helper()
	env := []string{"PATH=" + os.Getenv("PATH")}
	var out, errBuf bytes.Buffer
	run(context.Background(), argv, env, t.TempDir(), strings.NewReader(""), &out, &errBuf, &status.NativeRunner{}, nil)
	combined := out.String() + errBuf.String()
	for _, marker := range []string{"unknown command:", "unknown subcommand", "unknown dispatch subcommand"} {
		if strings.Contains(combined, marker) {
			return true
		}
	}
	return false
}

// TestProseFunctionNotationBindsToRouting (AC-1/AC-4) is the verb-notation ↔ command-routing
// binding. For every prose-function declaration in the restructured cores it extracts the
// migration target and asserts: a SHIPPED (backtick) verb is ACCEPTED by the compiled root
// (none of the three rejection markers), and a PROSE (guillemet) target that names a future
// `becomes` verb is REJECTED (one of the three markers). The binary's routing switch — not the
// contract prose — is the expected value: a contract that backticks a verb the binary does not
// route fails, and a guillemet whose named target now routes fails as should-have-flipped. The
// empty-walk guard keeps it from passing vacuously.
func TestProseFunctionNotationBindsToRouting(t *testing.T) {
	root := repoRootFromPkg(t)
	totalShipped, totalProse := 0, 0
	for _, rel := range proseFunctionCores {
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("read prose-function core %s: %v", rel, err)
		}
		bindings, err := extractVerbBindings(string(data))
		if err != nil {
			t.Fatalf("%s: %v", rel, err)
		}
		for _, b := range bindings {
			rejected := verbRejected(t, b.argv)
			switch b.notation {
			case "shipped":
				totalShipped++
				if rejected {
					t.Errorf("%s: backtick (shipped) verb %q is REJECTED by the binary routing — a backtick must name a routed command (line: %s)", rel, b.verb, b.line)
				}
			case "prose":
				totalProse++
				if !rejected {
					t.Errorf("%s: guillemet (prose) target %q is ACCEPTED by the binary routing — it should have flipped to a backtick (line: %s)", rel, b.verb, b.line)
				}
			}
		}
	}
	if totalShipped == 0 {
		t.Fatal("extracted zero shipped (backtick) verb bindings from the restructured cores — extraction bug; the binding check would pass vacuously")
	}
	if totalProse == 0 {
		t.Fatal("extracted zero prose (guillemet) targets from the restructured cores — extraction bug; the should-have-flipped half is never exercised")
	}
}

// TestProseFunctionRoutingGuardFailsOnViolation is the AC-1 RED control (mirroring
// boot_resident_closure_test.go's …GuardFailsOnDanglingTarget): it plants the two violation
// shapes against a fixture and proves the binding logic goes RED for each. (a) A backtick verb
// the binary does NOT route (`spacedock merge teleport`) — the shipped-must-route half reds.
// (b) A guillemet whose `→` target IS routed (`spacedock merge guard`) — the should-have-flipped
// half reds. It drives the same extractVerbBindings + verbRejected the real guard uses, so the
// control exercises the real code path rather than re-implementing it.
func TestProseFunctionRoutingGuardFailsOnViolation(t *testing.T) {
	fixture := "## «merge.phantom»(slug): a planted un-routed backtick\n" +
		"- → **shipped** (this sprint): `spacedock merge teleport` — invoke it directly.\n" +
		"## «merge.shouldflip»(slug): a planted routed guillemet target\n" +
		"- → **prose**, becomes `spacedock merge guard <slug>` (later) — hand-follow for now.\n"
	bindings, err := extractVerbBindings(fixture)
	if err != nil {
		t.Fatalf("control fixture failed to parse: %v", err)
	}
	if len(bindings) != 2 {
		t.Fatalf("control fixture extracted %d bindings, want 2 — the violation cases are not both exercised", len(bindings))
	}

	var sawPhantomRed, sawShouldFlipRed bool
	for _, b := range bindings {
		rejected := verbRejected(t, b.argv)
		switch {
		case b.notation == "shipped" && strings.Contains(b.verb, "teleport"):
			// shipped-must-route: an un-routed backtick is REJECTED, so the real guard's
			// shipped-branch would red here.
			if rejected {
				sawPhantomRed = true
			} else {
				t.Errorf("control: planted un-routed backtick %q was unexpectedly ACCEPTED — the shipped-must-route half cannot fail", b.verb)
			}
		case b.notation == "prose" && strings.Contains(b.verb, "merge guard"):
			// should-have-flipped: a routed guillemet target is ACCEPTED, so the real guard's
			// prose-branch would red here.
			if !rejected {
				sawShouldFlipRed = true
			} else {
				t.Errorf("control: planted routed guillemet target %q was unexpectedly REJECTED — the should-have-flipped half cannot fail", b.verb)
			}
		}
	}
	if !sawPhantomRed {
		t.Fatal("control: the un-routed-backtick violation was not exercised — the shipped-must-route guard cannot fail")
	}
	if !sawShouldFlipRed {
		t.Fatal("control: the routed-guillemet-target violation was not exercised — the should-have-flipped guard cannot fail")
	}
}

// TestProseFunctionRoutingOracleDiscriminates is the non-vacuity discriminator for the routing
// ORACLE: it proves verbRejected returns false for a representative SHIPPED verb and true for a
// representative UN-SHIPPED target at EACH of the three marker sites — so the oracle genuinely
// distinguishes shipped from un-shipped and the binding check cannot pass vacuously (e.g. by an
// oracle that always returns one value). Stubbing verbRejected to a constant reds this control.
func TestProseFunctionRoutingOracleDiscriminates(t *testing.T) {
	shipped := [][]string{{"merge", "guard"}, {"state", "commit"}, {"status", "--boot", "--json"}}
	for _, argv := range shipped {
		if verbRejected(t, argv) {
			t.Errorf("oracle: shipped verb %v wrongly reported REJECTED", argv)
		}
	}
	// One per marker site: root (`gate`), parent switch (`merge teleport`), dispatch handler.
	unshipped := [][]string{{"gate", "assemble-verdict"}, {"merge", "teleport"}, {"dispatch", "next-action"}}
	for _, argv := range unshipped {
		if !verbRejected(t, argv) {
			t.Errorf("oracle: un-shipped target %v wrongly reported ACCEPTED — a marker site is unhandled", argv)
		}
	}
}

// repoRootFromPkg returns the repository root, derived from this package's source dir
// (internal/cli), so the contract-file walk is independent of the test's working directory.
// go test runs with cwd = the package source dir, so the root is two levels up.
func repoRootFromPkg(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}
