// ABOUTME: The ONLY legal instruction-file reads in tests — structural checks a
// ABOUTME: machine can see (ref-closure, frontmatter-validity, structural-absence, dedup, no-machine-dependency).
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// This file is the single quarantined read path the boundary guard exempts. Every
// check here is STRUCTURAL — it catches a defect a machine can see without reading
// prose for meaning. No check asserts an instruction file contains its own prose
// (prose-grep), and none asserts the prose matches a code value as a stand-in for a
// behavior test (code-bound). See boundary_guard_test.go for the full policy.

// userSkills is the published user-skill surface: each owns a SKILL.md the host
// discovers. The test-only `integration` package is deliberately absent.
var userSkills = []string{
	"commission", "debrief", "refit", "survey", "ensign",
	"first-officer", "present-gate", "feedback-rejection-flow",
}

// skillsRoot is the shipped skill tree under test.
func skillsRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "skills")
}

// frontmatter returns the YAML frontmatter block (between the leading `---` and the
// next `---`) and whether the document opened with one.
func frontmatter(doc string) (string, bool) {
	if !strings.HasPrefix(doc, "---\n") {
		return "", false
	}
	rest := doc[len("---\n"):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return "", false
	}
	return rest[:end], true
}

// TestUserSkillsParseWithFrontmatter is a frontmatter-VALIDITY check: each user
// skill ships a SKILL.md whose YAML frontmatter block parses and declares a `name`
// and a `description` key. A malformed or keyless frontmatter block fails host
// discovery — a real structural defect, not a prose property.
func TestUserSkillsParseWithFrontmatter(t *testing.T) {
	root := skillsRoot(t)
	for _, skill := range userSkills {
		path := filepath.Join(root, skill, "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("user skill %q has no SKILL.md at %s: %v", skill, path, err)
			continue
		}
		fm, ok := frontmatter(string(data))
		if !ok {
			t.Errorf("%s/SKILL.md has no parseable YAML frontmatter block", skill)
			continue
		}
		if !frontmatterHasKey(fm, "name") {
			t.Errorf("%s/SKILL.md frontmatter declares no `name` key", skill)
		}
		if !frontmatterHasKey(fm, "description") {
			t.Errorf("%s/SKILL.md frontmatter declares no `description` key", skill)
		}
	}
}

// frontmatterHasKey reports whether a top-level `key:` line is declared in a parsed
// frontmatter block. Structural (a declared key), not a prose match on its value.
func frontmatterHasKey(fm, key string) bool {
	prefix := key + ":"
	for _, line := range strings.Split(fm, "\n") {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

// frontmatterField returns the trimmed scalar value of a top-level `key:` line in
// a flat frontmatter block.
func frontmatterField(fm, key string) string {
	prefix := key + ":"
	for _, line := range strings.Split(fm, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
		}
	}
	return ""
}

// discoverUserInvocableSkills scans the shipped skills tree the way the host does:
// every subdirectory with a SKILL.md whose frontmatter declares `user-invocable: true`
// is exposed as `/spacedock:<name>`.
func discoverUserInvocableSkills(t *testing.T) map[string]string {
	t.Helper()
	root := skillsRoot(t)
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read skills root %s: %v", root, err)
	}
	out := map[string]string{}
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "integration" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, e.Name(), "SKILL.md"))
		if err != nil {
			continue
		}
		fm, ok := frontmatter(string(data))
		if !ok || frontmatterField(fm, "user-invocable") != "true" {
			continue
		}
		name := frontmatterField(fm, "name")
		if name == "" {
			t.Errorf("user-invocable skill dir %q has no name field", e.Name())
			continue
		}
		out[name] = e.Name()
	}
	return out
}

// TestSurveyIsDiscoverableUserCommand is a structural frontmatter/discovery check
// kept inside the instruction-read quarantine. The behavior proof for survey's scan
// lives in skills/integration; this check only guards that the host can discover the
// `/spacedock:survey` command from the shipped skill tree.
func TestSurveyIsDiscoverableUserCommand(t *testing.T) {
	discovered := discoverUserInvocableSkills(t)
	dir, ok := discovered["survey"]
	if !ok {
		t.Fatalf("survey is not discoverable as /spacedock:survey; discovered user commands: %v", sortedUniqueKeys(discovered))
	}
	if dir != "survey" {
		t.Errorf("survey command resolves from dir %q, want skills/survey", dir)
	}
}

func sortedUniqueKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return sortedUnique(keys)
}

// referenceRe matches the two reference-include forms a SKILL.md uses: an
// `@references/foo.md` directive and a bare `references/foo.md` read path.
var referenceRe = regexp.MustCompile(`@?(references/[A-Za-z0-9_./-]+\.md)`)

// TestUserSkillReferenceClosureResolves is a ref-CLOSURE check: every
// `@references/...md` / `references/...md` path mentioned in a user SKILL.md
// resolves to a real file on disk under that skill's directory. A dangling
// reference (a ported skill pointing at a path that does not exist) is a real
// structural defect the host would hit at load time. Brace-placeholder template
// paths resolve against their concrete siblings.
func TestUserSkillReferenceClosureResolves(t *testing.T) {
	root := skillsRoot(t)
	for _, skill := range userSkills {
		skillDir := filepath.Join(root, skill)
		data, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
		if err != nil {
			t.Errorf("%s: %v", skill, err)
			continue
		}
		for _, m := range referenceRe.FindAllStringSubmatch(string(data), -1) {
			rel := m[1]
			if strings.Contains(rel, "{") {
				parent := filepath.Join(skillDir, filepath.Dir(rel))
				glob, _ := filepath.Glob(filepath.Join(parent, "*.md"))
				if len(glob) == 0 {
					t.Errorf("%s: templated reference %q has no concrete .md under %s", skill, rel, parent)
				}
				continue
			}
			if _, err := os.Stat(filepath.Join(skillDir, rel)); err != nil {
				t.Errorf("%s: dangling reference %q (resolved %s): %v", skill, rel, filepath.Join(skillDir, rel), err)
			}
		}
	}
}

// piRuntimeAdapters are the Pi runtime adapter references each host-aware skill
// must ship as a loadable file on disk.
var piRuntimeAdapters = []struct{ skill, ref string }{
	{skill: "first-officer", ref: "references/pi-first-officer-runtime.md"},
	{skill: "ensign", ref: "references/pi-ensign-runtime.md"},
}

// TestPiRuntimeAdaptersResolveOnDisk is a ref-CLOSURE check: each declared Pi
// runtime adapter resolves to a real file on disk. (The retired prose-grep half —
// asserting the SKILL.md advertises the ref string — is gone; that the file LOADS
// is the structural fact, proven by os.Stat.)
func TestPiRuntimeAdaptersResolveOnDisk(t *testing.T) {
	root := skillsRoot(t)
	for _, tc := range piRuntimeAdapters {
		path := filepath.Join(root, tc.skill, tc.ref)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("%s: Pi runtime adapter %s is not loadable on disk: %v", tc.skill, tc.ref, err)
		}
	}
}

// retiredPluginPrivatePaths are the plugin-private status paths the shipped
// instruction surface must NEVER name once it calls `spacedock status`. A
// re-introduced reference silently breaks the launcher contract on a fresh install.
var retiredPluginPrivatePaths = []string{
	"skills/commission/bin/status",
	"commission/bin/status",
	"{spacedock_plugin_dir}",
	".agents/plugins/marketplace.json",
}

// TestRetiredPluginPrivatePathsAbsent is a structural-ABSENCE check: no shipped
// instruction file (skills/**/*.md excluding the test-only integration dir, plus
// the canonical mods/) names a retired plugin-private status path. No positive
// behavioral seam can prove an absence; a re-introduced path is a real defect a
// fresh install would hit. (Consolidates the three prior near-duplicate
// plugin-private-path absence checks into one read path.)
func TestRetiredPluginPrivatePathsAbsent(t *testing.T) {
	files := shippedInstructionMarkdown(t)
	if len(files) == 0 {
		t.Fatal("walked zero shipped instruction files — scope bug; the absence check would pass vacuously")
	}
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		content := string(data)
		for _, p := range retiredPluginPrivatePaths {
			if strings.Contains(content, p) {
				t.Errorf("%s names retired plugin-private path %q — it does not exist on a fresh install", path, p)
			}
		}
	}
}

// TestNoUnexpectedModHookOrPRMergeIntroduced is a structural-ABSENCE check over
// the shipped surface: outside the pre-existing hook convention docs and the
// canonical pr-merge mod, instruction files must not introduce lifecycle-mod
// headings or PR-merge invocations. A new `## Hook:` mod or PR-merge command in a
// skill silently changes dispatch lifecycle, so the allowed files are explicit.
func TestNoUnexpectedModHookOrPRMergeIntroduced(t *testing.T) {
	allowedHookFiles := map[string]bool{
		filepath.Join("mods", "pr-merge.md"):                                                      true,
		filepath.Join("skills", "first-officer", "references", "claude-first-officer-runtime.md"): true,
		filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md"):    true,
		// The host-neutral extraction re-homed Mod-Block Enforcement (which names the
		// `## Hook: merge` mechanism surface) into the merge core; it legitimately
		// carries the `## Hook:` token.
		filepath.Join("skills", "first-officer", "references", "fo-merge-core.md"): true,
	}
	allowedPRMergeFiles := map[string]bool{
		filepath.Join("mods", "pr-merge.md"): true,
	}
	prMergeMarkers := []string{"gh pr merge", "git merge --no-ff", "git merge --ff-only main"}
	root := repoRoot(t)
	for _, path := range shippedInstructionMarkdown(t) {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			t.Errorf("rel %s: %v", path, err)
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		content := string(data)
		if strings.Contains(content, "## Hook:") && !allowedHookFiles[rel] {
			t.Errorf("%s introduces a lifecycle hook heading outside the allowed hook surfaces", rel)
		}
		for _, marker := range prMergeMarkers {
			if strings.Contains(content, marker) && !allowedPRMergeFiles[rel] {
				t.Errorf("%s introduces PR-merge invocation %q outside the canonical pr-merge mod", rel, marker)
			}
		}
	}
}

// The startup-gate single-source guard was RETIRED here per the package policy in
// doc_test.go:11 ("delete the read and report the owed test"). It walked skills/+agents/
// for hardcoded gate phrases (`"Contract version gate"`, `"per-class remedy"`,
// `"spacedock doctor"`) and required exactly one file to carry them, claiming to red
// when the guidance is mirrored so "the two surfaces drift". But the drift it guards
// is a MEANING mirrored, not a byte copy: agents/first-officer.md could restate the
// gate in its own words and the phrase-grep would stay green, so the check did not
// actually prove the single-source property it advertised — it proved only that those
// exact bytes appear once. That is the banned prose-grep (a literal-phrase substring
// as a proxy for "the gate guidance is not restated elsewhere"); narrowing it to
// verbatim-absence keeps it a prose-phrase-absence standing in for the meaning. No
// token here IS the gate guidance, so the genuine-structural path (cf. the literal
// claudeTeamDispatchTokens command tokens, where the token IS the thing) does not
// apply. OWED PROOF: the human gate-review of the FO contract is what verifies the
// startup gate is owned by first-officer-shared-core.md and delegated to (not mirrored
// by) agents/first-officer.md — agents/first-officer.md today carries only the
// delegating line "begin the Startup procedure from the shared core", which is the
// design that keeps the surfaces from drifting; that delegation is prose a machine
// cannot distinguish from a paraphrased mirror.

// homeRootedClaudeRe matches only HOME-rooted personal-config forms: a `~/.claude`
// tilde path, or `$HOME` / `os.UserHomeDir` joined with `.claude` on the same line.
// It does NOT match a project-relative `.claude/` path, which exists in any checkout
// and is portable.
var homeRootedClaudeRe = regexp.MustCompile(`~/\.claude|\$HOME[^\n]*\.claude|os\.UserHomeDir[^\n]*\.claude`)

// interpreterRe matches an interpreter-on-PATH dependency: a `python`/`python3`
// shell-out or a `commission/bin/...` plugin-private helper invocation.
var interpreterRe = regexp.MustCompile(`\bpython3?\b|commission/bin`)

// machineDependentPaths are plugin-private absolute paths that do not exist on a
// fresh install.
var machineDependentPaths = []string{
	"skills/commission/bin/status",
	"{spacedock_plugin_dir}",
	".agents/plugins/marketplace.json",
}

// isClaudeAdapter reports whether a shipped file is a Claude-host coupling surface
// (a claude-*-runtime.md adapter, a claude-fo-*.md FO module reference, or the
// survey skill), where a `~/.claude/teams` read is the legitimate quarantined
// coupling. ONLY the personal-config check excludes these; the interpreter /
// machine-path checks apply.
func isClaudeAdapter(path string) bool {
	base := filepath.Base(path)
	if strings.HasPrefix(base, "claude-") && strings.HasSuffix(base, "-runtime.md") {
		return true
	}
	// The Claude-host dispatch coupling (the `~/.claude/teams` and subagent-jsonl
	// reads) lives in the claude-fo-dispatch reference; it is the legitimate Claude
	// coupling surface, exempt from the HOME-rooted check.
	if strings.HasPrefix(base, "claude-fo-") && strings.HasSuffix(base, ".md") {
		return true
	}
	// The survey skill carries the same legitimate ~/.claude coupling (its
	// session-history scan reads the host's local agent state).
	return strings.Contains(path, filepath.Join("survey", "SKILL.md"))
}

// TestShippedSurfaceHasNoHiddenMachineDependency is a no-MACHINE-DEPENDENCY
// structural-absence check: the shipped instruction surface names none of three
// non-portable markers — a HOME-rooted personal-config path, an interpreter-on-PATH
// shell-out, or a plugin-private absolute path. A clean install must run for any
// user; a re-introduced marker is a real defect that breaks a fresh install. The
// empty-walk guard keeps it from passing vacuously.
func TestShippedSurfaceHasNoHiddenMachineDependency(t *testing.T) {
	files := shippedInstructionMarkdown(t)
	if len(files) == 0 {
		t.Fatal("walked zero shipped instruction files — scope bug; the portability check would pass vacuously")
	}
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		content := string(data)
		if !isClaudeAdapter(path) {
			if m := homeRootedClaudeRe.FindString(content); m != "" {
				t.Errorf("%s carries a HOME-rooted personal-config dependency %q — a clean install has no such file", path, m)
			}
		}
		if m := interpreterRe.FindString(content); m != "" {
			t.Errorf("%s carries an interpreter-on-PATH dependency %q — the dispatch path must not assume an installed interpreter/helper", path, m)
		}
		for _, p := range machineDependentPaths {
			if strings.Contains(content, p) {
				t.Errorf("%s bakes in plugin-private path %q — it does not exist on a fresh install", path, p)
			}
		}
	}
}

// TestPortabilityCheckDiscriminatesHostSpecific is the DISCRIMINATOR control for
// the no-machine-dependency check: it proves — against the real shipped surface —
// that the legitimately host-specific forms are present yet not flagged, so the
// absence check is not vacuous. The Claude adapter's `~/.claude/teams` read is
// present (the adapter exclusion is load-bearing), and the project-relative
// `.claude/` paths are present yet the HOME-rooted regex does not match them.
func TestPortabilityCheckDiscriminatesHostSpecific(t *testing.T) {
	files := shippedInstructionMarkdown(t)
	var sawAdapterHomeClaude, sawProjectRelativeClaude bool
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		content := string(data)
		if isClaudeAdapter(path) && strings.Contains(content, "~/.claude") {
			sawAdapterHomeClaude = true
		}
		if !isClaudeAdapter(path) {
			for _, line := range strings.Split(content, "\n") {
				if !strings.Contains(line, ".claude/") {
					continue
				}
				if strings.Contains(line, "~/.claude") || strings.Contains(line, "$HOME") {
					continue
				}
				sawProjectRelativeClaude = true
				if homeRootedClaudeRe.MatchString(line) {
					t.Errorf("%s: project-relative .claude line wrongly matched the HOME-rooted regex (false positive): %q", path, strings.TrimSpace(line))
				}
			}
		}
	}
	if !sawAdapterHomeClaude {
		t.Error("positive control missing: no Claude adapter carries a ~/.claude read — the adapter exclusion is no longer load-bearing")
	}
	if !sawProjectRelativeClaude {
		t.Error("positive control missing: no shipped file carries a project-relative .claude/ path — the discriminator has nothing to discriminate")
	}
}

// claudeTeamDispatchTokens are the LITERAL Claude-team command/flag tokens that
// must NOT appear as core dispatch imperatives in the host-neutral dispatch core.
// Their presence on a non-exempt line is itself the defect (structural-absence,
// same family as retiredPluginPrivatePaths) — not a prose property. The
// narrowly-paraphrasable phrase `team-mode` is deliberately NOT in this set: its
// fix is the prose move, not a token ban, and banning a paraphrasable word would
// be the prose-grep tautology the quarantine policy forbids.
var claudeTeamDispatchTokens = []string{
	"spawn-standing-all",
	"--team {team_name}",
}

// adapterAsSubjectExemptions are the PHRASE-LEVEL exemptions: a token line is
// allowed only when its sentence's subject is the adapter REALIZING the call. A
// bare same-line `runtime adapter` mention is deliberately NOT an exemption — the
// real current leak (line 9) names the adapter only as the forward destination
// ("forward each spawn spec … to the runtime adapter's spawn call") while still
// issuing the host command imperatively, so a bare-word exemption false-negatives
// it. The exemption must name the adapter as the actor.
var adapterAsSubjectExemptions = []string{
	"Claude adapter",
	"the adapter realizes",
	"adapter maps",
	"Claude realization",
}

// lineLeaksClaudeTeamDispatch reports whether a line carries a literal Claude-team
// dispatch token UNEXEMPTED by a phrase-level adapter-as-subject phrase. This is the
// scanner the host-neutral-core check and its discriminator control both drive.
func lineLeaksClaudeTeamDispatch(line string) bool {
	carriesToken := false
	for _, tok := range claudeTeamDispatchTokens {
		if strings.Contains(line, tok) {
			carriesToken = true
			break
		}
	}
	if !carriesToken {
		return false
	}
	for _, exempt := range adapterAsSubjectExemptions {
		if strings.Contains(line, exempt) {
			return false
		}
	}
	return true
}

// TestDispatchCoreHasNoClaudeTeamImperative is a structural-ABSENCE check: the
// host-neutral dispatch core (fo-dispatch-core.md) must carry no LITERAL Claude-team
// command/flag token (`spawn-standing-all`, `--team {team_name}`) as a core
// imperative. Codex and Pi load this core too, so a host command issued here reads
// to them like a universal requirement they cannot satisfy. The expected value
// comes from the rule (these are literal host command tokens that must not appear as
// core imperatives), not from the file's own prose, so a host-neutral paraphrase
// that drops the literal passes and an inverted/host-coupled imperative fails — same
// family as TestRetiredPluginPrivatePathsAbsent. The Claude realization legitimately
// lives in claude-fo-dispatch.md, whose lines name the adapter as the actor. The
// paired discriminator control below keeps this from passing vacuously.
func TestDispatchCoreHasNoClaudeTeamImperative(t *testing.T) {
	path := filepath.Join(skillsRoot(t), "first-officer", "references", "fo-dispatch-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read dispatch core %s: %v", path, err)
	}
	for i, line := range strings.Split(string(data), "\n") {
		if lineLeaksClaudeTeamDispatch(line) {
			t.Errorf("%s:%d carries a Claude-team dispatch imperative in the host-neutral core (Codex/Pi load this too); move the host realization to claude-fo-dispatch.md: %q", path, i+1, strings.TrimSpace(line))
		}
	}
}

// TestDispatchCoreClaudeTeamScannerDiscriminates is the DISCRIMINATOR control for
// the host-neutral-core check: it proves the scanner flags a genuinely host-coupled
// line and passes a host-neutral paraphrase, an adapter-as-subject realization, and
// a meaning-inverting paraphrase — so TestDispatchCoreHasNoClaudeTeamImperative can
// never pass vacuously (e.g. by a typo'd token never matching anything).
func TestDispatchCoreClaudeTeamScannerDiscriminates(t *testing.T) {
	// Host-coupled imperative — the real current-leak shape. MUST flag.
	mustFlag := "run `spacedock dispatch spawn-standing-all --workflow-dir {wd} --team {team_name}` and forward each spawn spec to the runtime adapter's spawn call"
	if !lineLeaksClaudeTeamDispatch(mustFlag) {
		t.Errorf("discriminator control: host-coupled imperative was NOT flagged (the scanner would pass vacuously): %q", mustFlag)
	}

	// A bare same-line `runtime adapter` forward-target mention is NOT a sufficient
	// exemption — the real leak carries exactly this and must still flag.
	bareAdapterWord := "forward each spawn spec to the runtime adapter's spawn call via `spawn-standing-all`"
	if !lineLeaksClaudeTeamDispatch(bareAdapterWord) {
		t.Errorf("discriminator control: a bare `runtime adapter` forward-target word wrongly exempted a token line (false negative on the real leak shape): %q", bareAdapterWord)
	}

	// Host-neutral paraphrase that drops the literal token. MUST pass.
	hostNeutral := "Before the first worker dispatch, inject standing teammates via the runtime adapter's standing-injection call."
	if lineLeaksClaudeTeamDispatch(hostNeutral) {
		t.Errorf("discriminator control: a host-neutral paraphrase was wrongly flagged: %q", hostNeutral)
	}

	// Adapter-as-subject realization naming the host command. MUST pass (this is the
	// shape that legitimately lives in claude-fo-dispatch.md).
	adapterSubject := "The Claude adapter realizes the standing-injection call: run `spacedock dispatch spawn-standing-all --team {team_name}` and forward each spawn spec to Agent()."
	if lineLeaksClaudeTeamDispatch(adapterSubject) {
		t.Errorf("discriminator control: an adapter-as-subject realization was wrongly flagged: %q", adapterSubject)
	}

	// Meaning-inverting paraphrase that still drops the literal. MUST pass — proving
	// the check is not a prose-grep that a paraphrase could defeat.
	inverted := "Standing teammates are NEVER injected before any dispatch on any host."
	if lineLeaksClaudeTeamDispatch(inverted) {
		t.Errorf("discriminator control: a meaning-inverting host-neutral paraphrase was wrongly flagged: %q", inverted)
	}
}

// shippedInstructionMarkdown returns every markdown file under skills/ (excluding
// the test-only integration dir) plus the canonical mods/ — the full shipped
// instruction surface the structural-absence and portability checks walk. This is
// the single instruction-surface walk for the quarantine package.
func shippedInstructionMarkdown(t *testing.T) []string {
	t.Helper()
	root := repoRoot(t)
	var out []string
	walk := func(base string) {
		filepath.WalkDir(base, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && d.Name() == "integration" {
				return filepath.SkipDir
			}
			if !d.IsDir() && strings.HasSuffix(p, ".md") {
				out = append(out, p)
			}
			return nil
		})
	}
	walk(filepath.Join(root, "skills"))
	walk(filepath.Join(root, "mods"))
	return out
}
