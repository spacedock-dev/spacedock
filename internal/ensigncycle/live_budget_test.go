// ABOUTME: AC-1 static guard: parses the live path's source AST and asserts no
// ABOUTME: timeout literal/const exceeds 60s and no monolithic deadline ctx remains in live_test.go.
package ensigncycle

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"testing"
	"time"
)

// AC-1 bans any individual MONOLITHIC timeout > 60s AND the deadline ctx. This
// guard PARSES the source AST and evaluates the real duration each `N *
// time.Unit` literal denotes — a relationship over real values, NOT a substring
// grep for the spelling "60". A literal written `120 * time.Second` would be
// caught (its real value is 120s); a literal written `1 * time.Minute` passes
// (60s). It also asserts `live_test.go` carries no `liveTimeout` identifier and
// no `context.WithTimeout` call — the monolithic deadline ctx this entity
// removes. The guard runs under DEFAULT build tags and reads live_test.go as
// source (the parser reads it regardless of its //go:build live tag), so it
// covers the live path from the offline suite with no model spend.
//
// EXEMPTION — no-progress quiet budgets. A `quietBudget*` const is NOT a
// monolithic timeout: it caps stream SILENCE, resetting on every drained line,
// so a stage that legitimately runs for minutes never trips as long as the
// stream keeps moving (the streamWatcher docstring's reconciliation of the
// "no individual timeout > 60s" directive with stages that take minutes). A
// dispatch-close step is one such stage — a single live ensign turn (boot,
// team-create, work, report) can stay quiet > 60s between lines — so
// quietBudgetDispatchClose is sanctioned to exceed the cap. The guard skips the
// initializer of any `quietBudget*` const; the >60s ban stays absolute on every
// other timeout literal.

// liveBudgetSources are the source files on the live path whose timeout literals
// must all be ≤60s: the streamWatcher (the per-step budget discipline), the live
// test that wires it, and the shared-scenario runners that ALSO drive the watcher.
// The shared runners were the unguarded gap that let the old per-scenario basket
// timeout exist; scanning them here brings them under the same ≤60s discipline, so
// they can never carry a >60s literal again.
var liveBudgetSources = []string{
	"streamwatch_test.go",
	"live_test.go",
	"claude_live_runner_test.go",
	"codex_live_runner_test.go",
	"codex_single_run_test.go",
}

// posSpan is the [start, end) source range of a quiet-budget const initializer
// the AST guard exempts from the >60s ban.
type posSpan struct {
	start token.Pos
	end   token.Pos
}

// quietBudgetAllowlist names the EXACT no-progress quiet-budget constants exempt
// from the >60s ban. An explicit allowlist (not a `quietBudget*` prefix) means a
// new constant must be deliberately added here to be exempted — naming something
// `quietBudgetSneakyTimeout` does NOT auto-exempt it. The guard bans a monolithic
// overall timeout because it MASKS a hang (the process sits dead under one big
// deadline); a no-progress quiet budget trips on stream SILENCE and resets on
// every drained line, so it CANNOT mask a hang — it catches one. The streamWatcher
// docstring already draws this line; this encodes it.
var quietBudgetAllowlist = map[string]bool{
	"quietBudgetDefault":       true,
	"quietBudgetDispatchClose": true,
}

// quietBudgetInitializerSpans collects the initializer-expression spans of every
// allowlisted quiet-budget const so the AST guard can skip them.
func quietBudgetInitializerSpans(f *ast.File) []posSpan {
	var spans []posSpan
	ast.Inspect(f, func(n ast.Node) bool {
		spec, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for i, name := range spec.Names {
			if i < len(spec.Values) && quietBudgetAllowlist[name.Name] {
				spans = append(spans, posSpan{spec.Values[i].Pos(), spec.Values[i].End()})
			}
		}
		return true
	})
	return spans
}

func TestNoTimeoutLiteralExceeds60s(t *testing.T) {
	for _, file := range liveBudgetSources {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		exempt := quietBudgetInitializerSpans(f)
		ast.Inspect(f, func(n ast.Node) bool {
			be, ok := n.(*ast.BinaryExpr)
			if !ok {
				return true
			}
			// Skip the initializer of a no-progress quiet-budget const — it is
			// sanctioned to exceed the 60s cap (see the EXEMPTION note above).
			for _, span := range exempt {
				if be.Pos() >= span.start && be.End() <= span.end {
					return false
				}
			}
			// Evaluate the WHOLE binary expression as a duration — folding ADD/SUB
			// composites, not just a bare MUL. `time.Minute + 30*time.Second`
			// evaluates to 90s here; a MUL-only fold would see 30s and miss the sum
			// (the audit-cycle-1 M3 hole). On a successful duration fold, check the
			// total and return false so we do NOT descend and re-flag the inner MUL
			// operands separately.
			d, isDuration := durationOf(be)
			if !isDuration {
				return true
			}
			if d > 60*time.Second {
				t.Errorf("%s: timeout literal %q evaluates to %s, exceeding the 60s cap (AC-1)",
					file, exprText(fset, be), d)
			}
			return false
		})
	}
}

// TestBudgetConstantsAreUnder60s pins the live path's budget constants to their
// REAL evaluated values at compile time — a direct value check that complements
// the source-AST guard above. The AST guard reasons about the SPELLING of every
// timeout literal in the file; this asserts the actual constants the live cycle
// uses are each ≤60s. Together they catch both an over-budget literal anywhere in
// the source AND an over-budget VALUE on the specific constants the watcher
// drives (e.g. a `quietBudgetDefault = time.Minute + 30*time.Second` = 90s — the
// audit-cycle-1 additive-composite hole — reds here directly regardless of how it
// is spelled).
func TestBudgetConstantsAreUnder60s(t *testing.T) {
	const budgetCap = 60 * time.Second
	for name, d := range map[string]time.Duration{
		"quietBudgetDefault": quietBudgetDefault,
		"exitBudgetDefault":  exitBudgetDefault,
	} {
		if d > budgetCap {
			t.Errorf("%s = %s exceeds the 60s cap (AC-1)", name, d)
		}
	}
}

// TestDurationOfFoldsScalarSpellings pins durationOf's scalar fold over the
// spellings the cycle-8 audit named — most importantly the `time.Duration(N)`
// integer conversion a basket could hide behind to dodge the BasicLit scalar
// match. Each foldable spelling must evaluate to its real duration (so the AST
// guard would flag it when >60s); the float-conversion form is correctly NOT
// statically foldable (the value-guard on the wired consts is its backstop).
func TestDurationOfFoldsScalarSpellings(t *testing.T) {
	foldable := map[string]time.Duration{
		"time.Second":                      time.Second,
		"120 * time.Second":                120 * time.Second,
		"time.Duration(120) * time.Second": 120 * time.Second, // the cycle-8 conversion evasion
		"time.Duration(2) * time.Minute":   2 * time.Minute,
		"time.Minute + 30*time.Second":     time.Minute + 30*time.Second,
		"2 * (40 * time.Second)":           80 * time.Second,
	}
	for src, want := range foldable {
		expr, err := parser.ParseExpr(src)
		if err != nil {
			t.Fatalf("parse %q: %v", src, err)
		}
		got, ok := durationOf(expr)
		if !ok {
			t.Errorf("durationOf(%q) did not fold; the AST guard would miss it", src)
			continue
		}
		if got != want {
			t.Errorf("durationOf(%q) = %s, want %s", src, got, want)
		}
	}

	// A float conversion is NOT statically foldable here — durationOf returns
	// false, by design. The wired-const value-guard (TestBudgetConstantsAreUnder60s)
	// is the backstop for any spelling the AST fold cannot resolve.
	floatExpr, err := parser.ParseExpr("time.Duration(1.5 * float64(time.Minute))")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := durationOf(floatExpr); ok {
		t.Error("durationOf unexpectedly folded a float conversion; that form is left to the const value-guard")
	}
}

// TestQuietBudgetExemptionIsNameGated pins that the AST guard's >60s exemption
// applies ONLY to constants in the EXPLICIT quietBudgetAllowlist — not to any
// const, and not to a `quietBudget`-prefixed name that was never allowlisted. A
// non-allowlisted >60s literal (a plain monolithic timeout, OR a sneaky
// `quietBudget`-prefixed one) is still collected as a non-exempt span, so the
// guard would flag it. This stops the carve-out from rotting into a blanket
// bypass: the no-progress-budget exemption must never shelter a monolithic
// timeout, and a new exemption must be a deliberate allowlist addition.
func TestQuietBudgetExemptionIsNameGated(t *testing.T) {
	const src = `package p
import "time"
const (
	quietBudgetDispatchClose = 3 * time.Minute
	someMonolithicTimeout    = 2 * time.Minute
	quietBudgetSneaky        = 4 * time.Minute
)
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "fake.go", src, 0)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	spans := quietBudgetInitializerSpans(f)
	if len(spans) != 1 {
		t.Fatalf("expected exactly 1 exempt span (only the allowlisted quietBudgetDispatchClose), got %d", len(spans))
	}
	// Walk the BinaryExprs and confirm ONLY the allowlisted 3m quiet budget is
	// inside an exempt span; the 2m plain monolithic literal AND the 4m
	// non-allowlisted `quietBudget`-prefixed literal are NOT (they would be flagged).
	var monolithicExempt, sneakyExempt, quietExempt bool
	ast.Inspect(f, func(n ast.Node) bool {
		be, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}
		d, isDur := durationOf(be)
		if !isDur {
			return true
		}
		within := false
		for _, span := range spans {
			if be.Pos() >= span.start && be.End() <= span.end {
				within = true
			}
		}
		switch d {
		case 3 * time.Minute:
			quietExempt = within
		case 2 * time.Minute:
			monolithicExempt = within
		case 4 * time.Minute:
			sneakyExempt = within
		}
		return false
	})
	if !quietExempt {
		t.Error("allowlisted quietBudgetDispatchClose literal was NOT exempt; the carve-out failed to cover it")
	}
	if monolithicExempt {
		t.Error("a plain monolithic literal was exempt; the carve-out leaked into a blanket bypass")
	}
	if sneakyExempt {
		t.Error("a non-allowlisted quietBudget-prefixed literal was exempt; the allowlist degraded to a prefix match")
	}
}

func TestLiveProcessPathsHaveNoMonolithicDeadlineCtx(t *testing.T) {
	for file, bannedIdent := range map[string]string{
		"live_test.go":             "liveTimeout",
		"codex_single_run_test.go": "codexScenarioTimeout",
	} {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.Ident:
				if node.Name == bannedIdent {
					t.Errorf("%s still references banned monolithic deadline %q", file, bannedIdent)
				}
			case *ast.SelectorExpr:
				if pkg, ok := node.X.(*ast.Ident); ok && pkg.Name == "context" && node.Sel.Name == "WithTimeout" {
					t.Errorf("%s still calls context.WithTimeout — fixed process deadlines are banned", file)
				}
			}
			return true
		})
	}
}

// durationOf evaluates an expression to its real time.Duration, FULLY
// COMPOSITIONALLY: every sub-expression recurses, so the COMPOSITE forms a
// timeout literal can take all fold —
//   - a bare `time.Unit` selector (`time.Second` = a 1× duration),
//   - a `time.Duration(N)` integer conversion (`time.Duration(120)` = 120ns,
//     so `time.Duration(120) * time.Second` folds to 120s — the cycle-8 hole a
//     BasicLit-only scalar match missed),
//   - a parenthesized duration (`(40 * time.Second)`),
//   - a `scalar * duration` (or `duration * scalar`) multiplication where the
//     duration operand is itself folded recursively, so `2 * (40 * time.Second)`
//     folds to 80s — the audit-cycle-2 hole a MUL branch that demanded a BARE
//     `time.Unit` operand missed,
//   - an ADD/SUB of two duration sub-expressions (`time.Minute + 30*time.Second`
//     = 90s — the audit-cycle-1 hole a MUL-only fold missed).
//
// Returns isDuration=false for anything that is not a STATICALLY-FOLDABLE duration
// expression, so the package's int arithmetic is skipped rather than mis-flagged.
// Forms it cannot statically fold (a float conversion, a const-ident or runtime
// scalar) yield false here; the wired live budget constants are the authoritative
// backstop — TestBudgetConstantsAreUnder60s value-checks their REAL compile-time
// values, catching any spelling the AST fold cannot.
func durationOf(e ast.Expr) (time.Duration, bool) {
	switch ex := e.(type) {
	case *ast.ParenExpr:
		return durationOf(ex.X)
	case *ast.SelectorExpr:
		// A bare `time.Second` etc. is a 1× duration.
		return timeUnitDuration(ex)
	case *ast.CallExpr:
		// A `time.Duration(N)` integer conversion is a duration of N nanoseconds.
		if n, ok := timeDurationConversion(ex); ok {
			return time.Duration(n), true
		}
		return 0, false
	case *ast.BinaryExpr:
		switch ex.Op {
		case token.MUL:
			// A duration MUL is `scalar * duration` in either order; the duration
			// operand is folded RECURSIVELY, so a parenthesized or nested duration
			// (e.g. `2 * (40 * time.Second)`) is not skipped.
			if n, ok := intScalar(ex.X); ok {
				if d, ok := durationOf(ex.Y); ok {
					return time.Duration(n) * d, true
				}
			}
			if n, ok := intScalar(ex.Y); ok {
				if d, ok := durationOf(ex.X); ok {
					return time.Duration(n) * d, true
				}
			}
			return 0, false
		case token.ADD, token.SUB:
			// A duration ADD/SUB requires BOTH operands to be durations; otherwise
			// it is not a duration expression (e.g. `len(x) + 1`).
			lhs, lok := durationOf(ex.X)
			rhs, rok := durationOf(ex.Y)
			if !lok || !rok {
				return 0, false
			}
			if ex.Op == token.ADD {
				return lhs + rhs, true
			}
			return lhs - rhs, true
		}
	}
	return 0, false
}

// timeDurationConversion folds a `time.Duration(N)` integer conversion to N. The
// scalar form a basket could hide behind to dodge the BasicLit scalar match:
// `time.Duration(120) * time.Second` evaluates to 120s here.
func timeDurationConversion(call *ast.CallExpr) (int64, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Duration" {
		return 0, false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "time" || len(call.Args) != 1 {
		return 0, false
	}
	return intScalar(call.Args[0])
}

// intScalar reports whether e is a statically-foldable integer scalar, returning
// its value. A parenthesized integer literal folds too (`(2) * time.Second`), as
// does a `time.Duration(N)` conversion used in scalar position
// (`time.Duration(120) * time.Second`).
func intScalar(e ast.Expr) (int64, bool) {
	if p, ok := e.(*ast.ParenExpr); ok {
		return intScalar(p.X)
	}
	if call, ok := e.(*ast.CallExpr); ok {
		return timeDurationConversion(call)
	}
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return 0, false
	}
	n, err := strconv.ParseInt(lit.Value, 0, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}

// timeUnitDuration maps a `time.Nanosecond|Microsecond|Millisecond|Second|
// Minute|Hour` selector to its time.Duration value.
func timeUnitDuration(e ast.Expr) (time.Duration, bool) {
	sel, ok := e.(*ast.SelectorExpr)
	if !ok {
		return 0, false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "time" {
		return 0, false
	}
	switch sel.Sel.Name {
	case "Nanosecond":
		return time.Nanosecond, true
	case "Microsecond":
		return time.Microsecond, true
	case "Millisecond":
		return time.Millisecond, true
	case "Second":
		return time.Second, true
	case "Minute":
		return time.Minute, true
	case "Hour":
		return time.Hour, true
	}
	return 0, false
}

// exprText renders a binary expression back to source for the failure message.
func exprText(fset *token.FileSet, be *ast.BinaryExpr) string {
	start := fset.Position(be.Pos())
	end := fset.Position(be.End())
	return start.String() + "-" + strconv.Itoa(end.Column)
}
