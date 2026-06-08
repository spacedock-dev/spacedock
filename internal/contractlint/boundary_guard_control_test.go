// ABOUTME: Mutation control for the boundary guard — a planted out-of-quarantine
// ABOUTME: instruction read must RED the detector; a non-instruction read must not.
package contractlint

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBoundaryGuardDetectsAPlantedInstructionRead is the guard's mutation control:
// the guard is the boundary oracle, so it must be demonstrated to RED on the exact
// shape it bans (a test that reads an instruction file) and stay GREEN on a read
// that is not an instruction file. Without this, the guard could silently degrade
// to a no-op (e.g. the detector predicate inverted) and pass vacuously. It parses
// synthetic source through the same detector the real sweep uses.
func TestBoundaryGuardDetectsAPlantedInstructionRead(t *testing.T) {
	dir := t.TempDir()

	// A planted out-of-quarantine instruction read: reads a skill body and inspects
	// it. The retired marker model would have let this pass with a declaration; the
	// guard bans it regardless.
	instructionRead := `package fixture
func TestReadsInstruction(t *T) {
	data, _ := os.ReadFile("../../skills/first-officer/references/first-officer-shared-core.md")
	if strings.Contains(string(data), "x") { _ = data }
}
`
	// A non-instruction read: reads a JSON manifest (the binary parses it) and a
	// generated state entity. Neither is an instruction surface, so the guard must
	// not flag it.
	nonInstructionRead := `package fixture
func TestReadsManifest(t *T) {
	a, _ := os.ReadFile("../../.agents/plugins/marketplace.json")
	b, _ := os.ReadFile(filepath.Join(stateDir, "entity", "index.md"))
	_ = a
	_ = b
}
`
	nonInstructionWalk := `package fixture
func TestWalksFixtureMarkdown(t *T) {
	filepath.WalkDir("testdata", func(path string, d DirEntry, err error) error {
		if strings.HasSuffix(path, ".md") {
			_ = path
		}
		return nil
	})
}
`
	writeFixture(t, filepath.Join(dir, "instruction_read_test.go"), instructionRead)
	if got := instructionReadingTestFiles(t, dir); !contains(got, "instruction_read_test.go") {
		t.Fatalf("guard failed to flag a planted instruction-file read; flagged=%v", got)
	}

	writeFixture(t, filepath.Join(dir, "manifest_read_test.go"), nonInstructionRead)
	writeFixture(t, filepath.Join(dir, "fixture_walk_test.go"), nonInstructionWalk)
	got := instructionReadingTestFiles(t, dir)
	if contains(got, "manifest_read_test.go") {
		t.Fatalf("guard wrongly flagged a non-instruction read (json manifest / generated state entity); flagged=%v", got)
	}
	if contains(got, "fixture_walk_test.go") {
		t.Fatalf("guard wrongly flagged a non-instruction markdown fixture walk; flagged=%v", got)
	}
	if !contains(got, "instruction_read_test.go") {
		t.Fatalf("adding a non-instruction fixture must not stop the guard flagging the instruction one; flagged=%v", got)
	}
}

// TestBoundaryGuardSweepFlagsAPlantedDir drives the directory sweep itself (not just
// the per-file detector) against a planted repo layout: a policed dir outside the
// quarantine that reads an instruction file must appear as an offender, and the
// quarantine dir's own instruction read must NOT.
func TestBoundaryGuardSweepFlagsAPlantedDir(t *testing.T) {
	root := t.TempDir()

	// A policed dir (not the quarantine) with an instruction read -> offender.
	policedDir := filepath.Join(root, "internal", "hostneutrality")
	if err := os.MkdirAll(policedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, filepath.Join(policedDir, "leak_test.go"), `package fixture
func TestLeak(t *T) {
	data, _ := os.ReadFile("../../skills/ensign/references/ensign-shared-core.md")
	_ = data
}
`)

	// The quarantine dir with an instruction read -> NOT an offender (the legal path).
	quarantineDir := filepath.Join(root, quarantinePkg)
	if err := os.MkdirAll(quarantineDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, filepath.Join(quarantineDir, "legal_read_test.go"), `package fixture
func TestLegalRead(t *T) {
	data, _ := os.ReadFile("../../skills/commission/SKILL.md")
	_ = data
}
`)
	// skills/integration must exist for the sweep's ReadDir; leave it empty (no reads).
	if err := os.MkdirAll(filepath.Join(root, "skills", "integration"), 0o755); err != nil {
		t.Fatal(err)
	}

	offenders := sweepInstructionReadsOutsideQuarantine(t, root)
	if !contains(offenders, filepath.Join("internal", "hostneutrality", "leak_test.go")) {
		t.Fatalf("sweep failed to flag the out-of-quarantine instruction read; offenders=%v", offenders)
	}
	otherDir := filepath.Join(root, "cmd", "spacedock")
	if err := os.MkdirAll(otherDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, filepath.Join(otherDir, "leak_test.go"), `package fixture
func TestLeak(t *T) {
	data, _ := os.ReadFile("../../skills/first-officer/SKILL.md")
	_ = data
}
`)
	offenders = sweepInstructionReadsOutsideQuarantine(t, root)
	if !contains(offenders, filepath.Join("cmd", "spacedock", "leak_test.go")) {
		t.Fatalf("sweep failed to flag an instruction read outside the historical packages; offenders=%v", offenders)
	}
	for _, o := range offenders {
		if filepath.Dir(o) == quarantinePkg {
			t.Fatalf("sweep wrongly flagged the quarantine package's own read %q; the quarantine is the legal read path", o)
		}
	}
}

// TestBoundaryGuardSweepSkipsAgentWorktrees pins that the sweep prunes the
// untracked agent-team worktree trees — both the Spacedock `.worktrees` and the
// Claude-Code `.claude` scratch. Those checkouts hold copies of instruction
// files that are not the repo's shipped surface, so an instruction read living
// inside one must NOT be reported as an offender. A real out-of-quarantine
// offender alongside them proves the sweep is still live, not skipping wholesale.
func TestBoundaryGuardSweepSkipsAgentWorktrees(t *testing.T) {
	root := t.TempDir()

	instructionRead := `package fixture
func TestReadsInstruction(t *T) {
	data, _ := os.ReadFile("skills/integration/SKILL.md")
	_ = data
}
`
	// A real offender outside the quarantine -> must be flagged.
	policedDir := filepath.Join(root, "internal", "hostneutrality")
	if err := os.MkdirAll(policedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, filepath.Join(policedDir, "leak_test.go"), instructionRead)

	// The same instruction read planted inside each agent-worktree scratch tree
	// -> must NOT be flagged.
	scratchDirs := []string{
		filepath.Join(root, ".worktrees", "agent-spacedock", "skills", "integration"),
		filepath.Join(root, ".claude", "worktrees", "agent-claude", "skills", "integration"),
	}
	for _, dir := range scratchDirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		writeFixture(t, filepath.Join(dir, "skill_surface_test.go"), instructionRead)
	}

	offenders := sweepInstructionReadsOutsideQuarantine(t, root)
	if !contains(offenders, filepath.Join("internal", "hostneutrality", "leak_test.go")) {
		t.Fatalf("sweep failed to flag the real out-of-quarantine read; offenders=%v", offenders)
	}
	for _, o := range offenders {
		if strings.HasPrefix(o, ".worktrees"+string(filepath.Separator)) ||
			strings.HasPrefix(o, ".claude"+string(filepath.Separator)) {
			t.Fatalf("sweep flagged an agent-worktree scratch read %q; those trees are untracked agent checkouts, not the shipped surface", o)
		}
	}
}

func contains(in []string, want string) bool {
	for _, s := range in {
		if s == want {
			return true
		}
	}
	return false
}

func writeFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ensure go/parser + token stay referenced for the fixture-shape sanity below.
var _ = parser.ParseFile
var _ = token.NewFileSet
