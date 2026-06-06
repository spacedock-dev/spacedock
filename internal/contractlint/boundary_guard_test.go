// ABOUTME: The instruction-file-read boundary guard — instruction/prompt files are
// ABOUTME: read in tests ONLY here in the quarantine package; this sweep bans them everywhere else.
package contractlint

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// READING INSTRUCTION/PROMPT FILES IN A TEST IS BANNED BY DEFAULT.
//
// An instruction file is a markdown surface the model ingests — a skill body
// (skills/**/SKILL.md), a contract/runtime reference (skills/**/references/*.md),
// or an agent definition (agents/*.md). The ONLY legal place a test reads one is
// THIS quarantine package, and only for a STRUCTURAL check — a defect a machine
// can see without reading prose for meaning:
//
//   - a @-reference resolves to a real file on disk (ref-closure),
//   - YAML frontmatter parses and declares its required keys (frontmatter-validity),
//   - a retired path is ABSENT from the shipped surface (structural-absence),
//   - a fact appears in exactly one source / a count holds (dedup),
//   - a file clears a size floor (line-floor).
//
// What is BANNED, here and everywhere:
//
//   - PROSE-GREP — asserting an instruction file contains (or lacks) its own prose.
//     The file is its own source of truth, so the assertion is a tautology: it
//     cannot fail for a real reason, and a meaning-inverting paraphrase keeps every
//     grepped token. "The skill says X" is never proof the system DOES X.
//   - CODE-BOUND-AS-BEHAVIOR-SUBSTITUTE — asserting an instruction file's prose
//     matches a code value (a const, a router subcommand, a seam name). That is a
//     consistency lint, not a behavior test. If the behavior matters there must be a
//     BEHAVIOR test that RUNS it; the prose⟷code check never proved the behavior.
//
// Behavior is proven by RUNNING it — a live scenario, an offline command-level
// drive, a code-side invariant over real parsed source — never by reading prose.
// The reader-shape axis (an undeclared read hiding via an undiscovered read shape)
// is backstopped by the detached adversarial audit at every high-stakes gate, not
// by enumerating read shapes here.
//
// This guard is the polarity flip of the retired AC-3 marker sweep: that sweep let
// a test read prose if it declared a marker (a permission slip); this guard bans
// the read outright outside this package. There is no marker. There is no taint
// flow. A read outside the quarantine is a failure, full stop.

// quarantinePkg is the only package directory where an instruction-file read is
// legal. Paths are repo-relative (the guard walks from the repo root).
const quarantinePkg = "internal/contractlint"

// TestNoInstructionReadsOutsideQuarantine is the boundary guard, re-runnable
// offline. It parses every *_test.go under the repo and FAILS if any test FILE
// outside the quarantine package contains a function that reads an instruction
// file's content. The quarantine package itself is exempt — it is the single legal
// read path for the structural checks. The count of out-of-quarantine reader files
// must be zero.
func TestNoInstructionReadsOutsideQuarantine(t *testing.T) {
	offenders := sweepInstructionReadsOutsideQuarantine(t, repoRoot(t))
	for _, o := range offenders {
		t.Errorf("%s reads an instruction file's content in a test — instruction/prompt files are read ONLY in %s for structural checks; behavior is proven by RUNNING it, never by reading prose (see the package doc)", o, quarantinePkg)
	}
	if len(offenders) > 0 {
		t.Fatalf("boundary guard: %d file(s) read instruction content outside %s; the count must be zero", len(offenders), quarantinePkg)
	}
}

// sweepInstructionReadsOutsideQuarantine returns the repo-relative paths of
// *_test.go files OUTSIDE the quarantine package that contain a function reading
// an instruction file's content. Exported logic so the guard's mutation control
// can drive it against a planted fixture directory.
func sweepInstructionReadsOutsideQuarantine(t *testing.T, repoRootDir string) []string {
	t.Helper()
	var offenders []string
	err := filepath.WalkDir(repoRootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(repoRootDir, path)
		if relErr != nil {
			return relErr
		}
		if d.IsDir() {
			switch rel {
			case ".git", ".worktrees", "docs/dev/.spacedock-state", "vendor", quarantinePkg:
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		f, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return parseErr
		}
		if fileReadsInstructionContent(f) {
			offenders = append(offenders, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("sweep instruction reads under %s: %v", repoRootDir, err)
	}
	return sortedUnique(offenders)
}

// instructionReadingTestFiles returns the base names of *_test.go files in dir
// that contain at least one function reading an instruction file's content.
func instructionReadingTestFiles(t *testing.T, dir string) []string {
	t.Helper()
	fset := token.NewFileSet()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read package dir %s: %v", dir, err)
	}
	var files []string
	for _, ent := range entries {
		name := ent.Name()
		if ent.IsDir() || !strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		if fileReadsInstructionContent(f) {
			files = append(files, name)
		}
	}
	return sortedUnique(files)
}

// fileReadsInstructionContent reports whether any function declared in f reads an
// instruction file's content.
func fileReadsInstructionContent(f *ast.File) bool {
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if directlyReadsInstructionFile(fn) {
			return true
		}
	}
	return false
}

func sortedUnique(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1] > out[j]; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}
