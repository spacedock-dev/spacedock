// ABOUTME: The instruction-file-read detector the boundary guard keys on — a
// ABOUTME: direct-read predicate (read sink + recognized instruction path, or a WalkDir .md collector).
package contractlint

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// readSinks are the call selectors that ingest a file's content given a path.
var readSinks = map[string]bool{
	"ReadFile":   true, // os.ReadFile
	"Open":       true, // os.Open
	"ReadAll":    true, // io.ReadAll
	"NewScanner": true, // bufio.NewScanner
	"NewReader":  true, // bufio.NewReader
}

// instructionPathSegments are the skill-tree / contract path segments that mark a
// path literal as an instruction file the model ingests (a skill, contract, agent,
// or runtime adapter). A path carrying one of these segments is an instruction
// path, recognized even before a `.md` suffix is appended.
var instructionPathSegments = map[string]bool{
	"skills":        true,
	"references":    true,
	"agents":        true,
	"first-officer": true,
	"ensign":        true,
	"commission":    true,
	"present-gate":  true,
	"SKILL.md":      true,
}

// directlyReadsInstructionFile reports whether fn DIRECTLY ingests a recognized
// instruction file's content: a read sink (ReadFile/Open/ReadAll/bufio) whose path
// arg subtree carries a recognized instruction literal/segment, or a WalkDir/Walk
// whose root arg carries a recognized instruction literal/segment and whose
// function body collects `.md` files. This is the direct-read predicate — no
// data-flow tracking. A read reached only through a param/local/field/method/
// closure flow, a transitive helper chain, a range-element flow, a cross-package
// reader, or an unrecognized surface is NOT detected here; that reader-shape axis
// is the detached-audit-backstopped boundary documented in the package guard doc.
func directlyReadsInstructionFile(fn *ast.FuncDecl) bool {
	found := false
	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if readSinks[sel.Sel.Name] {
			for _, arg := range call.Args {
				if exprCarriesInstructionPath(arg) {
					found = true
				}
			}
		}
		if (sel.Sel.Name == "WalkDir" || sel.Sel.Name == "Walk") &&
			len(call.Args) > 0 &&
			exprCarriesInstructionPath(call.Args[0]) &&
			fnFiltersInstructionMarkdown(fn) {
			found = true
		}
		return true
	})
	return found
}

// exprCarriesInstructionPath reports whether expr DIRECTLY carries an instruction-file
// path anywhere in its subtree: an instruction-path literal/segment — so the +/
// strings.Join/filepath.Join/fmt.Sprintf/string(...) path-build idioms (whose
// instruction operand is a subtree node) are covered when their operands are
// literals.
func exprCarriesInstructionPath(expr ast.Expr) bool {
	hit := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if x, ok := n.(*ast.BasicLit); ok {
			if x.Kind == token.STRING && isInstructionPathLiteral(strings.Trim(x.Value, "`\"")) {
				hit = true
			}
		}
		return true
	})
	return hit
}

// fnFiltersInstructionMarkdown reports whether fn's body filters paths by a `.md`
// suffix — the WalkDir-collector signal.
func fnFiltersInstructionMarkdown(fn *ast.FuncDecl) bool {
	hit := false
	ast.Inspect(fn, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if strings.HasSuffix(strings.Trim(lit.Value, "`\""), ".md") {
				hit = true
			}
		}
		return true
	})
	return hit
}

// isInstructionPathLiteral reports whether a string literal is (a fragment of) an
// instruction-file path: it carries a skill-tree/contract segment. A `.json`
// manifest path carries none and is not instruction.
func isInstructionPathLiteral(s string) bool {
	if strings.HasSuffix(s, ".json") {
		return false
	}
	for seg := range instructionPathSegments {
		if s == seg || strings.Contains(s, seg) {
			return true
		}
	}
	return false
}

// repoRoot returns the repository root, derived from this package's source dir
// (internal/contractlint), so the guard's filesystem walk is independent of the
// test's working directory.
func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	// go test runs with cwd = the package source dir (internal/contractlint).
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}
