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
			if !ok || be.Op != token.MUL {
				return true
			}
			d, isDuration := durationOfMul(be)
			if !isDuration {
				return true
			}
			if d > 60*time.Second {
				t.Errorf("%s: timeout literal %q evaluates to %s, exceeding the 60s cap (AC-1)",
					file, exprText(fset, be), d)
			}
			return true
		})
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

// durationOfMul evaluates a `N * time.Unit` (or `time.Unit * N`) binary
// expression to its real time.Duration. Returns isDuration=false when the
// expression is not an int-times-time.Unit shape, so non-duration multiplications
// (and the package's int arithmetic) are skipped rather than mis-flagged.
func durationOfMul(be *ast.BinaryExpr) (time.Duration, bool) {
	if n, unit, ok := intAndUnit(be.X, be.Y); ok {
		return time.Duration(n) * unit, true
	}
	if n, unit, ok := intAndUnit(be.Y, be.X); ok {
		return time.Duration(n) * unit, true
	}
	return 0, false
}

// intAndUnit reports whether `a` is an integer literal and `b` is a time.Unit
// selector, returning the parsed int and the unit's duration.
func intAndUnit(a, b ast.Expr) (int64, time.Duration, bool) {
	lit, ok := a.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return 0, 0, false
	}
	n, err := strconv.ParseInt(lit.Value, 0, 64)
	if err != nil {
		return 0, 0, false
	}
	unit, ok := timeUnitDuration(b)
	if !ok {
		return 0, 0, false
	}
	return n, unit, true
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
