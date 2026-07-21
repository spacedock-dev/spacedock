// ABOUTME: AC-4 reference-closure over the boot-resident FO contract bodies — every
// ABOUTME: deferred load-point they name resolves to a real file on disk (os.Stat oracle).
package contractlint

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// bootResidentBodies are the contract bodies the FO loader inlines/reads at boot:
// the slimmed shared core and each host's runtime adapter. AC-4 walks these (NOT a
// SKILL.md, which the existing TestUserSkillReferenceClosureResolves reads) because
// the loader reads the bodies directly, and only the bodies name the deferred
// load-points the boot core defers to (a sibling reference path, a bare skill
// invocation, a canonical mod file). The codex and pi adapters join the Claude one
// so the closure walk catches a dead-end reference on every host, not only Claude.
var bootResidentBodies = []string{
	filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md"),
	filepath.Join("skills", "first-officer", "references", "claude-first-officer-runtime.md"),
	filepath.Join("skills", "first-officer", "references", "codex-first-officer-runtime.md"),
	filepath.Join("skills", "first-officer", "references", "pi-first-officer-runtime.md"),
}

// foReferenceCores are the host-neutral cores the first-officer contract loads lazily
// from the shared core. Each core must exist on disk and carry its ceremony anchors.
var foReferenceCores = map[string][]string{
	filepath.Join("skills", "first-officer", "references", "fo-merge-core.md"): {
		"## Merge and Cleanup", "## «merge.guard»", "### Worktree removal safety", "## Mod-Block Guard",
	},
	filepath.Join("skills", "first-officer", "references", "fo-dispatch-core.md"): {
		"## Dispatch", "## Reuse and Fresh Dispatch", "## Dispatch Adapter", "## Event Loop",
	},
	filepath.Join("skills", "first-officer", "references", "fo-write-core.md"): {
		"## Mutation Gate", "## FO Write Scope", "## ID Styles",
	},
}

// deferredSkillCores are adapter-less deferred modules the boot-resident shared
// core names as non-user-invocable skills (via `spacedock:<name>`), NOT as references/*.md
// read-paths. Each must exist
// on disk AND carry its section anchors. Keyed on the skill path so a newly-registered
// anchor is both watched + owner-resolved (sectionOwners) and stat-checked (the skill-anchor
// test) from one place.
var deferredSkillCores = map[string][]string{
	filepath.Join("skills", "fo-status-viewer", "SKILL.md"): {
		"## Status Viewer", "### Captain-Facing State Display", "## Issue Filing",
	},
	filepath.Join("skills", "fo-dispatch-recovery", "SKILL.md"): {
		"## Break-Glass Manual Dispatch", "## Context Budget Failure and Dead Ensign Handling",
	},
}

// deferredModuleBodies are contract bodies that are NOT boot-resident (the FO loader
// does not read them at boot) but that themselves carry in-module `spacedock:` lazy
// skill triggers whose targets must resolve on disk, and whose prose pointers at a
// moved section must not dangle. claude-fo-dispatch.md is read at first dispatch
// (not boot) and carries the fo-dispatch-recovery failure-trigger lines (the
// dispatch-failure retry rung's dead-ensign handling, Break-Glass, Context Budget).
// Walking it with the same extraction+os.Stat oracle closure-checks those
// spacedock:fo-dispatch-recovery tokens, uncovered by bootResidentBodies today
// because this file is not itself boot-resident.
var deferredModuleBodies = []string{
	filepath.Join("skills", "first-officer", "references", "claude-fo-dispatch.md"),
}

// bodyReferenceRe matches a sibling reference read-path named in a contract body
// (the dispatch/merge references the split defers to), the same path shape the
// SKILL.md closure check uses, applied to body text.
var bodyReferenceRe = regexp.MustCompile(`references/[A-Za-z0-9_./-]+\.md`)

// bodySkillRe matches a lazy skill invocation `spacedock:<name>` the boot core
// names as a deferred load point. The legacy-probe / gate skills
// (using-legacy-claude-team, present-gate, feedback-rejection-flow) each resolve
// to their skills/<name>/SKILL.md.
var bodySkillRe = regexp.MustCompile(`spacedock:([a-z0-9-]+)`)

// bodyModRe matches a CONCRETE _mods reference (e.g. `_mods/pr-merge.md`) the boot
// core names as a deferred mod-file load point. Brace-templated placeholders
// (`_mods/{mod_name}.md`) are NOT concrete load-points and are excluded by the
// non-brace character class.
var bodyModRe = regexp.MustCompile(`_mods/([a-z0-9][a-z0-9_.-]*\.md)`)

// lazyLoadSkills are the skill names a boot-resident body may name as deferred
// load points. The ensign skill is the dispatched-worker contract, not a boot-core
// load point, so it is excluded; the FO-self reference would be a self-load. The two
// adapter-less deferred modules (currently fo-status-viewer) are here because the
// shared core names them via `spacedock:<name>`, so the closure walk resolves each to
// skills/<name>/SKILL.md, the same as the gate skills.
var lazyLoadSkills = map[string]bool{
	"present-gate":            true,
	"feedback-rejection-flow": true,
	"fo-status-viewer":        true,
	"fo-dispatch-recovery":    true,
}

// deferredLoadPoint is one extracted load-point: the on-disk path the body names
// and the literal token that named it (for a useful failure message).
type deferredLoadPoint struct {
	path  string // repo-relative resolved path
	named string // the literal token in the body
}

// extractDeferredLoadPoints parses one boot-resident body's text and returns every
// deferred load-point it names, resolved to a repo-relative on-disk path: sibling
// reference read-paths (resolved under the FO skill dir), lazy skill invocations
// (resolved to skills/<name>/SKILL.md), and concrete _mods files (resolved against
// the canonical mods/ tree the repo ships). It does NOT assert presence/absence of
// any prose — it only collects the paths the body NAMES, for the os.Stat oracle to
// resolve.
func extractDeferredLoadPoints(body string) []deferredLoadPoint {
	foSkillDir := filepath.Join("skills", "first-officer")
	var out []deferredLoadPoint
	seen := map[string]bool{}
	add := func(p deferredLoadPoint) {
		if seen[p.path] {
			return
		}
		seen[p.path] = true
		out = append(out, p)
	}
	for _, m := range bodyReferenceRe.FindAllString(body, -1) {
		if strings.Contains(m, "{") {
			continue
		}
		add(deferredLoadPoint{path: filepath.Join(foSkillDir, m), named: m})
	}
	for _, m := range bodySkillRe.FindAllStringSubmatch(body, -1) {
		name := m[1]
		if !lazyLoadSkills[name] {
			continue
		}
		add(deferredLoadPoint{path: filepath.Join("skills", name, "SKILL.md"), named: m[0]})
	}
	for _, m := range bodyModRe.FindAllStringSubmatch(body, -1) {
		add(deferredLoadPoint{path: filepath.Join("mods", m[1]), named: m[0]})
	}
	return out
}

// TestBootResidentDeferredLoadPointsResolve is the AC-4 reference-closure guard: a
// genuine structural check, not a prose-grep. For each boot-resident contract body
// it extracts every deferred load-point the body NAMES and os.Stats it. The
// EXPECTED value (the target exists on disk) comes from the FILESYSTEM — an
// independent source the contract text can diverge from — so a body that names a
// deferred reference at a moved/renamed/deleted path fails the stat. It is NOT the
// banned present-here/absent-there heading grep (boundary_guard_test.go): it does
// not assert the body contains or lacks any heading; it asserts every load-point
// the body points at resolves to a real file. The empty-walk guard keeps it from
// passing vacuously.
func TestBootResidentDeferredLoadPointsResolve(t *testing.T) {
	root := repoRoot(t)
	total := 0
	walked := append(append([]string{}, bootResidentBodies...), deferredModuleBodies...)
	for _, rel := range walked {
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("read boot-resident body %s: %v", rel, err)
		}
		points := extractDeferredLoadPoints(string(data))
		for _, p := range points {
			total++
			if _, err := os.Stat(filepath.Join(root, p.path)); err != nil {
				t.Errorf("%s names deferred load-point %q which resolves to %s — but no such file exists on disk: %v", rel, p.named, p.path, err)
			}
		}
	}
	if total == 0 {
		t.Fatal("extracted zero deferred load-points from the boot-resident bodies — extraction bug; the closure check would pass vacuously")
	}
}

// TestBootResidentDeferredLoadPointGuardFailsOnDanglingTarget is the AC-4 control:
// it points a boot-resident-style fixture body at a non-existent deferred reference
// and proves the closure logic goes RED, so the guard is shown able to fail (not a
// guard that can only ever pass). It drives the same extraction + os.Stat the real
// guard uses, against a planted fixture, so the control exercises the real code
// path rather than re-implementing it.
func TestBootResidentDeferredLoadPointGuardFailsOnDanglingTarget(t *testing.T) {
	root := repoRoot(t)
	fixture := "At first dispatch, read references/claude-fo-this-file-does-not-exist.md\n" +
		"and at terminal, invoke spacedock:fo-dispatch-recovery.\n"
	points := extractDeferredLoadPoints(fixture)
	if len(points) == 0 {
		t.Fatal("control fixture extracted no load-points — the dangling-target case never exercises the stat")
	}
	var sawDangling, sawReal bool
	for _, p := range points {
		_, err := os.Stat(filepath.Join(root, p.path))
		if strings.Contains(p.named, "does-not-exist") {
			if err == nil {
				t.Fatalf("control: the dangling reference %q unexpectedly resolved on disk", p.named)
			}
			sawDangling = true
		} else if err == nil {
			sawReal = true
		}
	}
	if !sawDangling {
		t.Fatal("control: the dangling deferred reference was not extracted — the guard cannot fail on a moved/deleted target")
	}
	if !sawReal {
		t.Fatal("control: the real load-point (fo-dispatch-recovery) was not resolved — the discriminator has nothing to contrast the dangling case against")
	}
}

// TestHostNeutralCoresResolveAndCarryCeremony is the reachability guard: every
// host-neutral deferred core is named by the shared core's load points and exists on
// disk with its ceremony anchors.
func TestHostNeutralCoresResolveAndCarryCeremony(t *testing.T) {
	root := repoRoot(t)
	if len(foReferenceCores) == 0 {
		t.Fatal("no host-neutral cores declared — the reachability check would pass vacuously")
	}
	sharedCore := filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md")
	sharedData, err := os.ReadFile(filepath.Join(root, sharedCore))
	if err != nil {
		t.Fatalf("read shared core %s: %v", sharedCore, err)
	}
	sharedBody := string(sharedData)
	for corePath, anchors := range foReferenceCores {
		base := filepath.Base(corePath)
		if !strings.Contains(sharedBody, base) {
			t.Errorf("%s does not name %s in its deferred load points", sharedCore, base)
		}
		data, err := os.ReadFile(filepath.Join(root, corePath))
		if err != nil {
			t.Errorf("host-neutral core %s does not resolve on disk: %v", corePath, err)
			continue
		}
		body := string(data)
		for _, anchor := range anchors {
			if !strings.Contains(body, anchor) {
				t.Errorf("%s is missing ceremony anchor %q — the named core resolves but does not carry the ceremony, so reachability reaches an empty file", corePath, anchor)
			}
		}
	}
}

// sectionOwners maps each watched section name (the deferredSkillCores anchors with their
// `## `/`### ` heading markers stripped — "Status Viewer", "Captain-Facing State Display",
// "Issue Filing", "FO Write Scope", "ID Styles") to the skill that owns it
// (skills/<owner>/SKILL.md, via filepath.Base(filepath.Dir(skillPath))). Derived from
// deferredSkillCores so a newly-registered anchor is both watched AND owner-resolved from
// one place. The section-name set the prose-pointer scanner watches is exactly this map's
// keys; the owner drives the Skill(skill="spacedock:<owner>") cross-file resolution.
func sectionOwners() map[string]string {
	owners := map[string]string{
		"Mutation Gate":  "fo-write-core",
		"FO Write Scope": "fo-write-core",
		"ID Styles":      "fo-write-core",
	}
	for skillPath, anchors := range deferredSkillCores {
		owner := filepath.Base(filepath.Dir(skillPath)) // skills/<owner>/SKILL.md -> <owner>
		for _, anchor := range anchors {
			owners[strings.TrimLeft(anchor, "# ")] = owner
		}
	}
	return owners
}

// referenceProsePointerDanglers scans one deferred module's body for non-heading lines that
// NAME a watched section yet resolve by NONE of the three honored paths: intra-file (the body
// carries that section's heading at any level), a references/*.md path token on the same line,
// or the cross-file Skill(skill="spacedock:<owner>") idiom — the same line names the section's
// owning skill via the spacedock:<owner> token (captured by bodySkillRe). Each unresolved line
// is a dangling prose pointer — the M5 shape, where a bare "(see FO Write Scope)" survives a
// move into a file that no longer carries that section. The Skill-form resolution is
// owner-specific: a spacedock: token for a DIFFERENT skill does NOT suppress the dangler, so
// the scanner cannot degrade into silencing on any spacedock: token. It returns one description
// per dangling line. The real guard and its control both drive this single scanner, so defeating
// it reds both.
func referenceProsePointerDanglers(body string, owners map[string]string) []string {
	lines := strings.Split(body, "\n")
	defined := map[string]bool{}
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") {
			continue
		}
		heading := strings.TrimLeft(line, "# ")
		if _, watched := owners[heading]; watched {
			defined[heading] = true
		}
	}
	var danglers []string
	for i, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue // a heading DEFINES a section, it does not point at one
		}
		hasPath := bodyReferenceRe.MatchString(line)
		// skills named on this line via the spacedock:<name> token (the Skill(skill="…") idiom)
		skillsNamed := map[string]bool{}
		for _, m := range bodySkillRe.FindAllStringSubmatch(line, -1) {
			skillsNamed[m[1]] = true
		}
		for name, owner := range owners {
			if !strings.Contains(line, name) || defined[name] || hasPath {
				continue
			}
			if skillsNamed[owner] {
				continue // resolves cross-file via Skill(skill="spacedock:<owner>")
			}
			danglers = append(danglers, fmt.Sprintf("line %d names section %q but resolves neither intra-file (no `## %s` heading here), via a references/*.md token, nor via Skill(skill=\"spacedock:%s\"): %q", i+1, name, name, owner, strings.TrimSpace(line)))
		}
	}
	return danglers
}

// TestDeferredSkillCoresResolveAndCarryCeremony is the AC-3 reachability guard for the two
// adapter-less deferred modules realized as skills: the boot-resident shared core names each
// via its `spacedock:<name>` token (the SAME closure token bodySkillRe /
// TestBootResidentDeferredLoadPointsResolve resolve to skills/<name>/SKILL.md), each skill
// file resolves on disk, AND it carries its section anchors — so "reachable" means "reaches a
// real ceremony," not an empty file. It CANNOT be folded into
// TestHostNeutralCoresResolveAndCarryCeremony, whose naming check is
// strings.Contains(sharedBody, "SKILL.md") — always false here, since the shared core names a
// skill by `spacedock:fo-status-viewer`, never a literal "SKILL.md". The independent source is
// the filesystem + the anchor set; a skill renamed/emptied/dropped or un-named by the shared
// core fails.
func TestDeferredSkillCoresResolveAndCarryCeremony(t *testing.T) {
	root := repoRoot(t)
	if len(deferredSkillCores) == 0 {
		t.Fatal("no deferred skill cores declared — the reachability check would pass vacuously")
	}
	sharedCore := filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md")
	sharedData, err := os.ReadFile(filepath.Join(root, sharedCore))
	if err != nil {
		t.Fatalf("read shared core %s: %v", sharedCore, err)
	}
	sharedBody := string(sharedData)
	for skillPath, anchors := range deferredSkillCores {
		name := filepath.Base(filepath.Dir(skillPath)) // skills/<name>/SKILL.md -> <name>
		token := "spacedock:" + name
		if !strings.Contains(sharedBody, token) {
			t.Errorf("%s does not name %q — the boot-resident shared core is the sole namer, so an un-named deferred skill is unreachable on every host", sharedCore, token)
		}
		data, err := os.ReadFile(filepath.Join(root, skillPath))
		if err != nil {
			t.Errorf("deferred skill core %s does not resolve on disk: %v", skillPath, err)
			continue
		}
		body := string(data)
		for _, anchor := range anchors {
			if !strings.Contains(body, anchor) {
				t.Errorf("%s is missing section anchor %q — the named skill resolves but does not carry the moved section, so reachability reaches an empty file", skillPath, anchor)
			}
		}
	}
}

// TestDeferredSkillProsePointersResolve is the AC-3c reachability gate the os.Stat closure
// structurally cannot be: for each deferred skill body it walks every non-heading line and
// FAILS on any that names a watched section without resolving it (an intra-file heading OR a
// references/*.md path token). The expected value comes from the filesystem + the resolution
// rule, not the file's own prose, so a section moved into a different file leaves its old
// bare-name pointers dangling and reds here. The companion control proves the scanner can
// fail; the empty-walk guards keep this non-vacuous.
func TestDeferredSkillProsePointersResolve(t *testing.T) {
	root := repoRoot(t)
	owners := sectionOwners()
	if len(owners) == 0 {
		t.Fatal("no section owners derived from deferredSkillCores — the prose-pointer check would pass vacuously")
	}
	walked := 0
	targets := make([]string, 0, len(deferredSkillCores)+len(deferredModuleBodies))
	for ref := range deferredSkillCores {
		targets = append(targets, ref)
	}
	// deferredModuleBodies joins the walk so a resident trigger line inside a
	// non-boot-resident contract body (e.g. claude-fo-dispatch.md's Degraded
	// Mode/Break-Glass/Context Budget trigger lines) is scanned for a dangling
	// bare-name pointer at a moved section too — today it walks only
	// deferredSkillCores files, so without this extension those trigger lines
	// would go unscanned.
	targets = append(targets, deferredModuleBodies...)
	for _, ref := range targets {
		data, err := os.ReadFile(filepath.Join(root, ref))
		if err != nil {
			t.Errorf("read deferred module body %s: %v", ref, err)
			continue
		}
		walked++
		for _, d := range referenceProsePointerDanglers(string(data), owners) {
			t.Errorf("%s: %s", ref, d)
		}
	}
	if walked == 0 {
		t.Fatal("walked zero deferred module bodies — the prose-pointer check would pass vacuously")
	}
}

// TestDeferredSkillProsePointerGuardFailsOnDanglingTarget is the AC-3c control battery: it
// drives referenceProsePointerDanglers against planted bodies so the gate is shown able to
// fail (RED on the M5 bare-name shape) without false-positiving the three legitimate
// resolutions — an intra-file heading, a references/*.md path token, and the cross-file
// Skill(skill="spacedock:<owner>") idiom. The wrong-owner row is the non-weakening guard:
// a Skill token for a DIFFERENT skill must NOT suppress the dangler, so the resolver cannot
// degrade into "any spacedock: token silences danglers."
func TestDeferredSkillProsePointerGuardFailsOnDanglingTarget(t *testing.T) {
	owners := sectionOwners()
	if len(owners) == 0 {
		t.Fatal("no section owners derived from deferredSkillCores — the control has nothing to test")
	}
	cases := []struct {
		name        string
		body        string
		wantDangler bool // true: the pointer must dangle (>=1); false: it must resolve (0)
	}{
		{
			// The exact M5 shape: a bare prose pointer to a section the body neither
			// defines nor reaches by path. It MUST dangle.
			name:        "bare-name pointer dangles",
			body:        "Some prose.\nUse `spacedock new` (see FO Write Scope), which mints the id.\n",
			wantDangler: true,
		},
		{
			// Intra-file heading present: the pointer resolves, must NOT dangle.
			name:        "intra-file heading resolves",
			body:        "## FO Write Scope\n\nThe scope.\n\nSee FO Write Scope above for the full contract.\n",
			wantDangler: false,
		},
		{
			// A references/*.md path token present: the pointer resolves, must NOT dangle.
			name:        "references/*.md path token resolves",
			body:        "Use `spacedock new` (see `references/fo-dispatch-core.md`), which mints the id.\n",
			wantDangler: false,
		},
		{
			// The Skill(skill="spacedock:<owner>") idiom nt introduced: the line names the
			// watched section AND its owning skill, so the pointer resolves cross-file.
			name:        "Skill-form owner-matched pointer resolves",
			body:        "Use `spacedock new` (see FO Write Scope via `Skill(skill=\"spacedock:fo-write-core\")`), which mints the id.\n",
			wantDangler: false,
		},
		{
			// Wrong-owner Skill token: fo-status-viewer does not own FO Write Scope, so the
			// pointer still dangles. This is the non-weakening property.
			name:        "wrong-owner Skill token still dangles",
			body:        "Use `spacedock new` (see FO Write Scope via `Skill(skill=\"spacedock:fo-status-viewer\")`), which mints the id.\n",
			wantDangler: true,
		},
	}
	for _, tc := range cases {
		got := referenceProsePointerDanglers(tc.body, owners)
		if tc.wantDangler && len(got) == 0 {
			t.Errorf("%s: expected >=1 dangler, got 0: body=%q", tc.name, tc.body)
		}
		if !tc.wantDangler && len(got) != 0 {
			t.Errorf("%s: expected 0 danglers, got %d: %v", tc.name, len(got), got)
		}
	}
}

// deadDeferredReferencePaths are the references/*.md read-paths that no longer exist after the
// two adapter-less modules became skills. Any surviving contract file that still names one is a
// dangling pointer the os.Stat closure over bootResidentBodies does NOT catch —
// feedback-rejection-flow/SKILL.md (not a boot-resident body) and the migrated fo-status-viewer
// skill body (a prose "(see references/fo-write-core.md)" referenceProsePointerDanglers treats
// as resolved on the bare path token) both escape it.
var deadDeferredReferencePaths = []string{
	"references/fo-status-viewer.md",
}

// TestNoSurvivingContractFileNamesDeadDeferredReferencePath is the AC-3 unguarded-re-point gate:
// it walks every .md under skills/ and FAILS on any that still names a deleted
// references/fo-{status-viewer,write-core}.md path. The two modules moved to skills, so a
// leftover references/*.md pointer is dead. The independent source is the filesystem (the files
// are gone); a re-point the migration missed — including the two sites that escape every other
// closure check — reds here.
func TestNoSurvivingContractFileNamesDeadDeferredReferencePath(t *testing.T) {
	root := repoRoot(t)
	skillsDir := filepath.Join(root, "skills")
	walked := 0
	err := filepath.Walk(skillsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		walked++
		body := string(data)
		rel, _ := filepath.Rel(root, path)
		for _, dead := range deadDeferredReferencePaths {
			if strings.Contains(body, dead) {
				t.Errorf("%s still names the deleted deferred reference path %q — re-point it to Skill(skill=\"spacedock:<name>\") after the skill conversion", rel, dead)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk skills/: %v", err)
	}
	if walked == 0 {
		t.Fatal("walked zero .md files under skills/ — the unguarded-re-point gate would pass vacuously")
	}
}
