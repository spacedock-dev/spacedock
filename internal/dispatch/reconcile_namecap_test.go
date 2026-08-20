// ABOUTME: name-cap reconcile round-trip — a capped (id-prefix) worker name
// ABOUTME: resolves to the correct entity by id, no false Class-A against a sibling.
package dispatch

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// nameCapReconcileFixture seeds an active dir + _archive dir with two sd-b32
// entities whose slugs share a long common prefix but whose ids differ, emits the
// capped name for each via real dispatch build, and runs real Reconcile over a
// stub roster carrying those capped names. The expectation (which slug each
// capped name resolves to, and which is Class-A) is bound to the seeded on-disk
// id set — an independent source from the resolution under test.
type nameCapReconcileFixture struct {
	t            *testing.T
	home         string
	repoRoot     string
	workflowDir  string
	stateRoot    string
	teamName     string
	archivedName string // capped name of the archived long-slug entity
	activeName   string // capped name of the active long-slug sibling
}

// sdb32NameReadme is a split-root README declaring id-style: sd-b32 (so the build
// cap takes the id-first path) and a state checkout reconcile reads entities from.
const sdb32NameReadme = `---
entity-type: task
id-style: sd-b32
state: .spacedock-state
stages:
  defaults:
    worktree: false
    concurrency: 1
  states:
    - name: backlog
      initial: true
    - name: implementation
      worktree: true
    - name: validation
      worktree: true
      feedback-to: implementation
    - name: done
      terminal: true
---
# Fixture Workflow

### backlog

seed.

- **Outputs:** x.

### implementation

work.

- **Outputs:** y.

### validation

verify.

- **Outputs:** z.

### done

term.
`

// buildCappedName drives real dispatch build over a sd-b32 entity at the state
// checkout carrying the given slug + id, returning the capped name it emits.
func buildCappedName(t *testing.T, stateRoot, workflowDir, slug, id string) string {
	t.Helper()
	ep := filepath.Join(stateRoot, slug+".md")
	writeFile(t, ep, entityFMID(id, "Thing", "backlog"))
	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   workflowDir,
		"stage":          "backlog",
		"checklist":      []string{"- a"},
		"bare_mode":      false,
	}, nil)
	native := runNative(stdin, "build", "--workflow-dir", workflowDir)
	if native.exit != 0 {
		t.Fatalf("build %s exit=%d\nstderr:\n%s", slug, native.exit, native.stderr)
	}
	return nameFromStdout(t, native.stdout)
}

func newNameCapReconcileFixture(t *testing.T) *nameCapReconcileFixture {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	repoRoot := t.TempDir()
	teamName := "team-namecap-001"

	workflowDir := filepath.Join(repoRoot, "docs", "wf")
	writeFile(t, filepath.Join(workflowDir, "README.md"), sdb32NameReadme)
	stateRoot := filepath.Join(workflowDir, ".spacedock-state")
	gitInit(t, repoRoot)

	// Emit the capped name for the archived long-slug entity (idAlpha) and the
	// active long-slug sibling (idBravo). The build writes both entity files into
	// the active state root first; we then move the archived one under _archive.
	archivedName := buildCappedName(t, stateRoot, workflowDir, longSlug, idAlpha)
	activeName := buildCappedName(t, stateRoot, workflowDir, longSlugShare, idBravo)

	// Move the archived entity's file into the _archive convention dir; rewrite it
	// with status=done so the archived-branch of classA fires. The active sibling
	// stays in the active dir with a non-terminal status.
	archiveDir := filepath.Join(stateRoot, "_archive")
	writeFile(t, filepath.Join(archiveDir, longSlug+".md"), entityFMID(idAlpha, "Thing", "done"))
	if err := os.Remove(filepath.Join(stateRoot, longSlug+".md")); err != nil {
		t.Fatalf("remove archived source: %v", err)
	}
	// Active sibling: non-terminal status so it must NOT classify as Class-A.
	writeFile(t, filepath.Join(stateRoot, longSlugShare+".md"), entityFMID(idBravo, "Thing", "implementation"))

	// Team config: team-lead (exempt) plus both capped-name ensigns.
	cfgPath := filepath.Join(home, ".claude", "teams", teamName, "config.json")
	writeFile(t, cfgPath, teamConfigJSON(teamName, []claudeteam.ReconcileMember{
		{Name: "team-lead", AgentType: "team-lead", Model: "claude-opus-4-7"},
		{Name: archivedName, AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
		{Name: activeName, AgentType: "spacedock:ensign", Model: "claude-opus-4-7"},
	}))

	return &nameCapReconcileFixture{
		t:            t,
		home:         home,
		repoRoot:     repoRoot,
		workflowDir:  workflowDir,
		stateRoot:    stateRoot,
		teamName:     teamName,
		archivedName: archivedName,
		activeName:   activeName,
	}
}

func (f *nameCapReconcileFixture) run() reconcileResult {
	f.t.Helper()
	var stdout, stderr bytes.Buffer
	opts := reconcileOpts{
		workflowDir: f.workflowDir,
		teamName:    f.teamName,
		repoRoot:    f.repoRoot,
		include:     map[string]bool{classLingering: true, classSuperseded: true, classUnadvancedPR: true, classStaleBranch: true, classLocalMainDrift: true},
		home:        f.home,
		roster:      claudeteam.LoadReconcileTeam,
		gh:          func(string) (string, error) { return "", nil },
		git:         gitRunnerExec,
	}
	code := Reconcile(opts, &stdout, &stderr)
	if code != 0 {
		f.t.Fatalf("Reconcile exit=%d stderr=%s", code, stderr.String())
	}
	var result reconcileResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		f.t.Fatalf("decode result: %v\nstdout=%s", err, stdout.String())
	}
	return result
}

// TestReconcileCappedNameResolvesArchived — AC3: the capped name of an archived
// long-slug sd-b32 entity classifies Class-A against the CORRECT slug (resolved
// by id-prefix), and does NOT mis-resolve to the active sibling sharing the slug
// prefix. The expectation is bound to the seeded on-disk id set.
func TestReconcileCappedNameResolvesArchived(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newNameCapReconcileFixture(t)
	result := f.run()

	byClass := groupDriftByClass(result.Drift)
	as := byClass[classLingering]
	if len(as) != 1 {
		t.Fatalf("expected exactly 1 lingering entry (archived only); got %d: %s",
			len(as), formatDrift(result.Drift))
	}
	a := as[0]
	if a.Name != f.archivedName {
		t.Errorf("lingering name=%q, want archived capped name %q", a.Name, f.archivedName)
	}
	// Resolution must recover the archived entity's REAL slug from the id-prefix,
	// not the active sibling's slug.
	if a.Slug != longSlug {
		t.Errorf("lingering slug=%q, want %q (id-prefix mis-resolved?)", a.Slug, longSlug)
	}
	if a.Slug == longSlugShare {
		t.Errorf("lingering mis-resolved to the active sibling slug %q", longSlugShare)
	}
}

// TestReconcileCappedNameNoFalseClassA — AC4: the capped name of a STILL-ACTIVE
// long-slug sd-b32 entity produces NO Class-A entry (guards the false-shutdown
// hazard — a capped name must not be flagged as lingering against a live worker).
func TestReconcileCappedNameNoFalseClassA(t *testing.T) {
	if !hasGit(t) {
		t.Skip("git not available")
	}
	f := newNameCapReconcileFixture(t)
	result := f.run()

	for _, d := range result.Drift {
		if d.Class == classLingering && d.Name == f.activeName {
			t.Errorf("active sibling's capped name %q falsely classified lingering: %+v",
				f.activeName, d)
		}
	}
}

// Adversarial id sets pinning resolveSlugToken's two safety properties. The
// AC3/AC4 end-to-end fixture seeds ids that diverge at char 0, so it cannot
// exercise either property; these prefix-sharing ids do. All are 24-char ids
// from the sd-b32 alphabet (0123456789abcdefghjkmnpqrstvwxyz).
const (
	// idShareShort + idShareShortSibling share only a 6-char prefix (aaaaaa) and
	// diverge by char 10, so the full 10-char token discriminates them but a
	// 2-char comparison would not.
	idShareShort        = "aaaaaa0000bbbbccccddddee"
	idShareShortSibling = "aaaaaa1111bbbbccccddddee"
	// idShareLong + idShareLongSibling share a 10-char prefix (cccccccccc), so a
	// 10-char id-prefix token matches BOTH — the ambiguity case.
	idShareLong        = "cccccccccc0000000011112e"
	idShareLongSibling = "cccccccccc1111111122223e"
)

// TestResolveSlugTokenSafetyProperties pins the two safety properties of the
// id-prefix resolver directly against real code (no mocks). The AC3/AC4 fixture's
// ids diverge at char 0 and so leave both properties untested; an adversarial
// reviewer could weaken the ambiguity guard (count != 1 → count < 1) or the
// over-match comparison (full token → token[:2]) and stay green against that
// fixture. These cases turn RED under each weakening.
func TestResolveSlugTokenSafetyProperties(t *testing.T) {
	// Full-token discrimination (refutes the token[:SDB32MinPrefix] over-match
	// edit): two ids sharing only a 6-char prefix, both ACTIVE. The 10-char token
	// of one is a prefix of exactly one id, so it must resolve to that entity. A
	// 2-char-only comparison would match both ids → ambiguous → unclassified,
	// failing this assertion.
	t.Run("full-token discriminates a short-shared-prefix sibling", func(t *testing.T) {
		active := map[string]entityRecord{
			"alpha-slug": {slug: "alpha-slug", id: idShareShort},
			"beta-slug":  {slug: "beta-slug", id: idShareShortSibling},
		}
		archived := map[string]entityRecord{}
		token := idShareShort[:sdB32NameIDPrefixLen] // 10 chars: "aaaaaa0000"
		got, ok := resolveSlugToken(token, active, archived)
		if !ok {
			t.Fatalf("token %q should resolve uniquely; got ok=false (over-match would see both ids)", token)
		}
		if got != "alpha-slug" {
			t.Errorf("token %q resolved to %q, want alpha-slug (the only id with this 10-char prefix)", token, got)
		}
	})

	// Ambiguity guard (refutes the count < 1 edit): two ACTIVE ids share a ≥10-char
	// prefix, so the 10-char token is a prefix of BOTH. The resolver must leave the
	// token UNCLASSIFIED (ok=false) — resolving to either would be a false Class-A
	// against the live sibling, the exact #366 hazard. A count < 1 guard would
	// pass an arbitrary one through.
	t.Run("ambiguous 10-char prefix stays unclassified", func(t *testing.T) {
		active := map[string]entityRecord{
			"gamma-slug": {slug: "gamma-slug", id: idShareLong},
			"delta-slug": {slug: "delta-slug", id: idShareLongSibling},
		}
		archived := map[string]entityRecord{}
		token := idShareLong[:sdB32NameIDPrefixLen] // 10 chars: "cccccccccc" (prefix of both)
		got, ok := resolveSlugToken(token, active, archived)
		if ok {
			t.Errorf("ambiguous token %q resolved to %q, want unclassified (count==2 must not classify)", token, got)
		}
	})

	// Exact-slug-first still wins (the uncapped common case is untouched): a token
	// equal to a real slug resolves to that slug even when an unrelated id exists.
	t.Run("exact slug match wins over id-prefix", func(t *testing.T) {
		active := map[string]entityRecord{
			"alpha-slug": {slug: "alpha-slug", id: idShareShort},
		}
		got, ok := resolveSlugToken("alpha-slug", active, map[string]entityRecord{})
		if !ok || got != "alpha-slug" {
			t.Errorf("exact slug token resolved to (%q,%v), want (alpha-slug,true)", got, ok)
		}
	})
}
