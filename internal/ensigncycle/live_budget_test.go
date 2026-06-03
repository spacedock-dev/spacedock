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

// AC-1 bans any individual timeout > 60s AND the monolithic deadline ctx. This
// guard PARSES the source AST and evaluates the real duration each `N *
// time.Unit` literal denotes — a relationship over real values, NOT a substring
// grep for the spelling "60". A literal written `120 * time.Second` would be
// caught (its real value is 120s); a literal written `1 * time.Minute` passes
// (60s). It also asserts `live_test.go` carries no `liveTimeout` identifier and
// no `context.WithTimeout` call — the monolithic deadline ctx this entity
// removes. The guard runs under DEFAULT build tags and reads live_test.go as
// source (the parser reads it regardless of its //go:build live tag), so it
// covers the live path from the offline suite with no model spend.

// liveBudgetSources are the source files on the live path whose timeout literals
// must all be ≤60s: the streamWatcher (the per-step budget discipline) and the
// live test that wires it.
var liveBudgetSources = []string{"streamwatch.go", "live_test.go"}

func TestNoTimeoutLiteralExceeds60s(t *testing.T) {
	for _, file := range liveBudgetSources {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			be, ok := n.(*ast.BinaryExpr)
			if !ok {
				return true
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
// drives (e.g. a `holdConfirmDefault = time.Minute + 30*time.Second` = 90s — the
// audit-cycle-1 M3 hole — reds here directly regardless of how it is spelled).
func TestBudgetConstantsAreUnder60s(t *testing.T) {
	const budgetCap = 60 * time.Second
	for name, d := range map[string]time.Duration{
		"quietBudgetDefault": quietBudgetDefault,
		"exitBudgetDefault":  exitBudgetDefault,
		"holdConfirmDefault": holdConfirmDefault,
	} {
		if d > budgetCap {
			t.Errorf("%s = %s exceeds the 60s cap (AC-1)", name, d)
		}
	}
}

func TestLiveTestHasNoMonolithicDeadlineCtx(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "live_test.go", nil, 0)
	if err != nil {
		t.Fatalf("parse live_test.go: %v", err)
	}
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.Ident:
			if node.Name == "liveTimeout" {
				t.Errorf("live_test.go still references the banned monolithic `liveTimeout` const (AC-1)")
			}
		case *ast.SelectorExpr:
			// context.WithTimeout — the banned monolithic deadline ctx.
			if pkg, ok := node.X.(*ast.Ident); ok && pkg.Name == "context" && node.Sel.Name == "WithTimeout" {
				t.Errorf("live_test.go still calls context.WithTimeout — the monolithic deadline ctx is banned (AC-1)")
			}
		}
		return true
	})
}

// durationOf evaluates an expression to its real time.Duration, FULLY
// COMPOSITIONALLY: every sub-expression recurses, so the COMPOSITE forms a
// timeout literal can take all fold —
//   - a bare `time.Unit` selector (`time.Second` = a 1× duration),
//   - a parenthesized duration (`(40 * time.Second)`),
//   - a `scalar * duration` (or `duration * scalar`) multiplication where the
//     duration operand is itself folded recursively, so `2 * (40 * time.Second)`
//     folds to 80s — the audit-cycle-2 hole a MUL branch that demanded a BARE
//     `time.Unit` operand missed,
//   - an ADD/SUB of two duration sub-expressions (`time.Minute + 30*time.Second`
//     = 90s — the audit-cycle-1 hole a MUL-only fold missed).
//
// Returns isDuration=false for anything that is not a pure duration expression,
// so the package's int arithmetic is skipped rather than mis-flagged.
func durationOf(e ast.Expr) (time.Duration, bool) {
	switch ex := e.(type) {
	case *ast.ParenExpr:
		return durationOf(ex.X)
	case *ast.SelectorExpr:
		// A bare `time.Second` etc. is a 1× duration.
		return timeUnitDuration(ex)
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

// intScalar reports whether e is an integer literal, returning its value. A
// parenthesized integer literal folds too, so `(2) * time.Second` is handled.
func intScalar(e ast.Expr) (int64, bool) {
	if p, ok := e.(*ast.ParenExpr); ok {
		return intScalar(p.X)
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
