// ABOUTME: AC-1's git-scaffold-init guard -- every non-bare `git ... init` site
// ABOUTME: in internal/**/*_test.go must go through testgit.InitRepo; bare-repo inits are exempt.
package contractlint

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// gitScaffoldExemptPkg is the one directory allowed to hand-roll `git init`
// calls directly: it is testgit.InitRepo's own implementation and its
// falsifiability tests, which intentionally scaffold a repo WITHOUT InitRepo as
// the negative case. Every other non-bare `git init` site under internal/ must
// call testgit.InitRepo instead.
const gitScaffoldExemptPkg = "internal/testgit"

// TestNoHandRolledGitInitOutsideTestgit is AC-1's guard: it walks
// internal/**/*_test.go and fails, listing file:line, for any non-bare
// `git ... init` invocation that isn't routed through testgit.InitRepo.
// `git init --bare` sites are excluded -- a bare repo never authors a commit,
// so it is not in the identity-config bug class this entity closes. The
// offender count is the moving baseline AC-1 measures: it must go from ~34
// (before migration) to 0 (after), and a future hand-rolled fixture reds this
// guard instead of silently reopening the hole.
func TestNoHandRolledGitInitOutsideTestgit(t *testing.T) {
	offenders := sweepHandRolledGitInit(t, repoRoot(t))
	for _, o := range offenders {
		t.Errorf("%s hand-rolls a git init scaffold -- use testgit.InitRepo instead so identity is always persisted (see internal/testgit)", o)
	}
	if len(offenders) > 0 {
		t.Fatalf("git-scaffold guard: %d non-bare `git init` site(s) outside testgit.InitRepo; the count must be zero", len(offenders))
	}
}

// sweepHandRolledGitInit returns sorted "file:line" locations, outside
// gitScaffoldExemptPkg, of non-bare `git init` invocations in *_test.go files
// under internal/.
func sweepHandRolledGitInit(t *testing.T, repoRootDir string) []string {
	t.Helper()
	internalDir := filepath.Join(repoRootDir, "internal")
	var offenders []string
	fset := token.NewFileSet()
	err := filepath.WalkDir(internalDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(repoRootDir, path)
		if relErr != nil {
			return relErr
		}
		if d.IsDir() {
			if rel == gitScaffoldExemptPkg {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		f, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return parseErr
		}
		for _, pos := range handRolledGitInitPositions(f) {
			p := fset.Position(pos)
			offenders = append(offenders, rel+":"+strconv.Itoa(p.Line))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("sweep git-init scaffolds under %s: %v", internalDir, err)
	}
	return sortedUnique(offenders)
}

// handRolledGitInitPositions returns the positions of every non-bare git-init
// scaffold site in f: a direct call that runs `git init` (or `--bare` init,
// excluded below), and a `{"init", ...}` subcommand-table entry of the shape
// used by the file-local `for _, args := range [][]string{...}` helpers.
func handRolledGitInitPositions(f *ast.File) []token.Pos {
	var positions []token.Pos
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if pos, ok := directGitInitCallPosition(node); ok {
				positions = append(positions, pos)
			}
		case *ast.RangeStmt:
			positions = append(positions, gitInitSubcommandTablePositions(node)...)
		}
		return true
	})
	return positions
}

// gitInitSubcommandTablePositions handles the file-local helper shape
// `for _, args := range [][]string{ {"init", "-q"}, {"add", "-A"}, ... } { ... }`
// where a subcommand-argv table is ranged over and each entry fed to a git
// call. It returns the position of any "init"-headed entry, gated on the loop
// body actually invoking git -- otherwise an unrelated `[][]string` test table
// (e.g. a CLI-args-per-subtest table) would false-positive on a coincidental
// `{"init", ...}` entry that has nothing to do with git.
func gitInitSubcommandTablePositions(rng *ast.RangeStmt) []token.Pos {
	lit, ok := rng.X.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	arr, ok := lit.Type.(*ast.ArrayType)
	if !ok {
		return nil
	}
	if _, ok := arr.Elt.(*ast.ArrayType); !ok {
		return nil
	}
	if !rangeBodyInvokesGit(rng.Body) {
		return nil
	}
	var positions []token.Pos
	for _, elt := range lit.Elts {
		if pos, ok := gitInitSubcommandLiteralPosition(elt); ok {
			positions = append(positions, pos)
		}
	}
	return positions
}

// rangeBodyInvokesGit reports whether body contains a call recognized as a
// direct git invocation (see calleeArgsIfGitCall), regardless of its
// subcommand -- the signal that a ranged subcommand table feeds an actual git
// call rather than an unrelated string-table.
func rangeBodyInvokesGit(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if _, ok := calleeArgsIfGitCall(call); ok {
				found = true
			}
		}
		return true
	})
	return found
}

// calleeArgsIfGitCall reports whether call is a direct git invocation --
// either `exec.Command("git", ...)` or a local wrapper whose name contains
// "git" (e.g. `git`, `gitC`, `gitOutput`, `mustGit`) -- and if so returns the
// argument list carrying the git subcommand (with a leading "git" literal
// stripped for the exec.Command form).
func calleeArgsIfGitCall(call *ast.CallExpr) ([]ast.Expr, bool) {
	var calleeName string
	switch fn := call.Fun.(type) {
	case *ast.SelectorExpr:
		calleeName = fn.Sel.Name
	case *ast.Ident:
		calleeName = fn.Name
	default:
		return nil, false
	}
	args := call.Args
	if calleeName == "Command" && len(args) > 0 && basicLitString(args[0]) == "git" {
		return args[1:], true
	}
	if strings.Contains(strings.ToLower(calleeName), "git") {
		return args, true
	}
	return nil, false
}

// directGitInitCallPosition reports the position of the "init" argument in
// call, if call is a direct git invocation (see calleeArgsIfGitCall) carrying
// a literal "init" subcommand and no sibling "--bare".
func directGitInitCallPosition(call *ast.CallExpr) (token.Pos, bool) {
	args, isGitCall := calleeArgsIfGitCall(call)
	if !isGitCall {
		return token.NoPos, false
	}

	subcommand, subcommandPos, ok := gitSubcommandLiteral(args)
	if !ok || subcommand != "init" {
		return token.NoPos, false
	}
	for _, arg := range args {
		if basicLitString(arg) == "--bare" {
			return token.NoPos, false
		}
	}
	return subcommandPos, true
}

// gitSubcommandLiteral scans args in order for the git subcommand literal --
// the first string literal not consumed as part of a `-C <dir>` or
// `-c <key>=<value>` global-option pair (both of which take exactly one
// following value and may precede the subcommand). Non-literal args (e.g. a
// `*testing.T` or a path variable) are skipped positionally.
func gitSubcommandLiteral(args []ast.Expr) (string, token.Pos, bool) {
	for i := 0; i < len(args); i++ {
		lit, ok := args[i].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			continue
		}
		s := basicLitString(args[i])
		if s == "-C" || s == "-c" {
			i++ // skip this flag's value too
			continue
		}
		return s, lit.Pos(), true
	}
	return "", token.NoPos, false
}

// gitInitSubcommandLiteralPosition reports the position of an "init" element
// heading a `[]string{"init", ...}` composite literal (typed, or type-elided as
// a nested element of a `[][]string{...}` table) -- the subcommand-table shape
// the file-local range-based init helpers use. A sibling "--bare" excludes it.
func gitInitSubcommandLiteralPosition(elt ast.Expr) (token.Pos, bool) {
	lit, ok := elt.(*ast.CompositeLit)
	if !ok {
		return token.NoPos, false
	}
	if lit.Type != nil {
		if _, ok := lit.Type.(*ast.ArrayType); !ok {
			return token.NoPos, false
		}
	}
	if len(lit.Elts) == 0 {
		return token.NoPos, false
	}
	first, ok := lit.Elts[0].(*ast.BasicLit)
	if !ok || first.Kind != token.STRING || basicLitString(lit.Elts[0]) != "init" {
		return token.NoPos, false
	}
	for _, elt := range lit.Elts[1:] {
		if basicLitString(elt) == "--bare" {
			return token.NoPos, false
		}
	}
	return first.Pos(), true
}

// basicLitString returns expr's unquoted string value, or "" if expr is not a
// string BasicLit.
func basicLitString(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	return strings.Trim(lit.Value, "`\"")
}
