// ABOUTME: dispatch reconcile — computes a five-class drift report (lingering /
// ABOUTME: superseded / un-advanced PR / stale branch / stale local main) from disk.
package dispatch

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
	"github.com/spacedock-dev/spacedock/internal/status"
)

// canonicalStages are the workflow stage names the decomposer falls back to when
// a workflow's own stage names cannot be derived from its README. The entity's
// design names these as the "five-canonical ensign stages" so a workflow with
// custom stages still decomposes cleanly without code changes — the per-workflow
// README scan augments this set.
var canonicalStages = []string{"backlog", "ideation", "implementation", "validation", "done"}

// reconcileOpts captures the resolved CLI flags for one reconcile run. roster
// is the team-roster loader the helper calls — in production this is
// claudeteam.LoadReconcileTeam, which owns the only ~/.claude/teams read; tests
// inject a stub so the host-neutrality scan stays clean.
type reconcileOpts struct {
	workflowDir string
	teamName    string
	repoRoot    string
	include     map[string]bool
	home        string
	roster      rosterLoader
	gh          ghRunner
	git         gitRunner
	now         func() int64
}

// rosterLoader is the team-roster injection point. claudeteam.LoadReconcileTeam
// satisfies it; tests pass an in-memory stub.
type rosterLoader func(home, teamName string) (claudeteam.ReconcileTeamState, error)

// ghRunner runs `gh pr view {N} --json state` returning the state string (e.g.
// "MERGED", "OPEN"). Tests inject a stub; production passes ghRunnerExec.
type ghRunner func(prRef string) (string, error)

// gitRunner runs `git -C dir args...` returning combined output and error.
// Tests inject a stub when they need a fixture-controlled answer; production
// passes gitRunnerExec.
type gitRunner func(dir string, args ...string) (string, error)

// ghRunnerExec calls `gh pr view PR --json state` and extracts the state field.
func ghRunnerExec(prRef string) (string, error) {
	cmd := exec.Command("gh", "pr", "view", prRef, "--json", "state")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	var obj struct {
		State string `json:"state"`
	}
	if err := json.Unmarshal(out, &obj); err != nil {
		return "", err
	}
	return obj.State, nil
}

// gitRunnerExec runs git with the given args under -C dir and returns trimmed
// combined output. The error preserves the exit status so callers can decide.
func gitRunnerExec(dir string, args ...string) (string, error) {
	full := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", full...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// runReconcile parses argv for the `reconcile` subcommand and routes to Reconcile.
// Exit codes match the entity contract: 0 sweep ran (drift may be empty or
// populated); 1 setup failure (team config missing, workflow dir invalid, git
// not in repo); 2 usage error.
func runReconcile(args []string, stdout, stderr io.Writer) int {
	flags := parseFlags(args, map[string]bool{
		"--workflow-dir": true, "--team-name": true,
		"--repo-root": true, "--include": true,
	})
	workflowDir, ok := flags["--workflow-dir"]
	if !ok {
		fmt.Fprintln(stderr, "error: dispatch reconcile requires --workflow-dir")
		return 2
	}
	include, err := parseInclude(flags["--include"])
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 2
	}
	opts := reconcileOpts{
		workflowDir: workflowDir,
		teamName:    flags["--team-name"],
		repoRoot:    flags["--repo-root"],
		include:     include,
		home:        os.Getenv("HOME"),
		roster:      claudeteam.LoadReconcileTeam,
		gh:          ghRunnerExec,
		git:         gitRunnerExec,
	}
	return Reconcile(opts, stdout, stderr)
}

// parseInclude resolves the --include flag value. An empty string means all five
// classes; otherwise the value is a comma-separated subset of {A,B,C,D,E}. An
// unknown class is a usage error (exit 2). Returned map sets true for included
// classes only.
func parseInclude(raw string) (map[string]bool, error) {
	all := map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true}
	if raw == "" {
		return all, nil
	}
	out := map[string]bool{}
	for _, c := range strings.Split(raw, ",") {
		c = strings.TrimSpace(c)
		if !all[c] {
			return nil, fmt.Errorf("unknown --include class: %q (expected A,B,C,D,E)", c)
		}
		out[c] = true
	}
	return out, nil
}

// driftItem is one entry in the drift[] array. Fields are populated per class —
// omitempty keeps each class's JSON object minimal.
type driftItem struct {
	Class    string `json:"class"`
	Name     string `json:"name,omitempty"`
	Slug     string `json:"slug,omitempty"`
	Stage    string `json:"stage,omitempty"`
	Worktree string `json:"worktree,omitempty"`
	PR       string `json:"pr,omitempty"`
	Behind   int    `json:"behind,omitempty"`
	Ahead    int    `json:"ahead,omitempty"`
	Reason   string `json:"reason"`
}

// reconcileResult is the {"command":"reconcile","team_name","drift":[]} envelope.
type reconcileResult struct {
	Command  string      `json:"command"`
	TeamName string      `json:"team_name"`
	Drift    []driftItem `json:"drift"`
}

// Reconcile runs one drift sweep and emits the JSON envelope. It does NO writes:
// every state read (team config, entity frontmatter, git refs, gh) is read-only,
// and every action (shutdown, reset, PR-advance) is the FO's responsibility per
// the action table. Exit 0 sweep ran; 1 setup failure; 2 caller-side usage
// (validated upstream in runReconcile).
func Reconcile(opts reconcileOpts, stdout, stderr io.Writer) int {
	if info, err := os.Stat(opts.workflowDir); err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "error: workflow directory not found: %s\n", opts.workflowDir)
		return 1
	}
	repoRoot := opts.repoRoot
	if repoRoot == "" {
		repoRoot = status.FindGitRoot(opts.workflowDir)
		if repoRoot == "" {
			fmt.Fprintf(stderr, "error: not in a git repo (cannot resolve --repo-root from %s)\n",
				opts.workflowDir)
			return 1
		}
	}

	team, err := opts.roster(opts.home, opts.teamName)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}

	stateRoot := splitRootStateCheckout(opts.workflowDir)
	if stateRoot == "" {
		stateRoot = opts.workflowDir
	}
	active := loadEntityFrontmatter(activeEntityDir(stateRoot))
	archived := loadEntityFrontmatter(archivedEntityDir(stateRoot))

	stageNames := readStageNames(opts.workflowDir)

	var drift []driftItem

	ensigns := filterEnsigns(team.Members)
	if opts.include["A"] {
		drift = append(drift, classA(ensigns, stageNames, active, archived)...)
	}
	if opts.include["B"] {
		drift = append(drift, classB(ensigns, stageNames)...)
	}
	if opts.include["C"] {
		drift = append(drift, classC(active, opts.gh)...)
	}
	if opts.include["D"] {
		drift = append(drift, classD(active, repoRoot, opts.git)...)
	}
	if opts.include["E"] {
		drift = append(drift, classE(repoRoot, opts.git)...)
	}

	sortDrift(drift)

	result := reconcileResult{
		Command:  "reconcile",
		TeamName: team.TeamName,
		Drift:    drift,
	}
	if result.Drift == nil {
		result.Drift = []driftItem{}
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintf(stderr, "error: encoding result: %s\n", err)
		return 1
	}
	return 0
}

// filterEnsigns returns only members whose agentType is spacedock:ensign. The
// team-lead and any standing teammates (comm-officer, science-officer, etc.)
// are exempt and never get classified.
func filterEnsigns(members []claudeteam.ReconcileMember) []claudeteam.ReconcileMember {
	var out []claudeteam.ReconcileMember
	for _, m := range members {
		if m.AgentType == "spacedock:ensign" {
			out = append(out, m)
		}
	}
	return out
}

// entityRecord is the minimal frontmatter the sweep reads — every field is a
// frontmatter top-level key (no body parsing).
type entityRecord struct {
	slug     string
	status   string
	worktree string
	pr       string
}

// loadEntityFrontmatter scans dir for `{slug}/index.md` or `{slug}.md` entries
// and returns a slug→record map. Missing dir → empty map (callers iterate).
func loadEntityFrontmatter(dir string) map[string]entityRecord {
	out := map[string]entityRecord{}
	if dir == "" {
		return out
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return out
	}
	for _, e := range entries {
		name := e.Name()
		// skip dot-prefixed dirs and the _archive convention (caller handles archived).
		if strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
			continue
		}
		var entityPath, slug string
		if e.IsDir() {
			slug = name
			entityPath = filepath.Join(dir, name, "index.md")
		} else if strings.HasSuffix(name, ".md") {
			slug = strings.TrimSuffix(name, ".md")
			entityPath = filepath.Join(dir, name)
		} else {
			continue
		}
		if _, err := os.Stat(entityPath); err != nil {
			continue
		}
		fm := status.ParseFrontmatter(entityPath)
		out[slug] = entityRecord{
			slug:     slug,
			status:   fm["status"],
			worktree: fm["worktree"],
			pr:       fm["pr"],
		}
	}
	return out
}

// activeEntityDir is the directory holding non-archived entity files in the
// state checkout (or workflow dir for non-split-root).
func activeEntityDir(stateRoot string) string { return stateRoot }

// archivedEntityDir is the _archive convention dir under the state root.
func archivedEntityDir(stateRoot string) string { return filepath.Join(stateRoot, "_archive") }

// readStageNames returns the workflow's declared stage names plus the canonical
// fallback set. README parse failures fall back to the canonical set alone —
// the decomposer still works on any workflow using standard stage names.
func readStageNames(workflowDir string) []string {
	seen := map[string]bool{}
	stages, _ := status.ParseStagesWithDefaults(filepath.Join(workflowDir, "README.md"))
	for _, s := range stages {
		if s.Name != "" {
			seen[s.Name] = true
		}
	}
	for _, n := range canonicalStages {
		seen[n] = true
	}
	out := make([]string, 0, len(seen))
	for n := range seen {
		out = append(out, n)
	}
	// Sort longest-first so a stage like "implementation" is tried before "imp" —
	// strip-trailing-stage must be greedy on the longest match. Lexicographic
	// tiebreaker keeps the order deterministic.
	sort.Slice(out, func(i, j int) bool {
		if len(out[i]) != len(out[j]) {
			return len(out[i]) > len(out[j])
		}
		return out[i] < out[j]
	})
	return out
}

// decomposeResult is what decompose returns: the entity slug, the stage suffix
// (or "" when unknown), and the cycle suffix (e.g. "cycle2", or "" for cycle 1).
type decomposeResult struct {
	slug  string
	stage string
	cycle string
	ok    bool
}

// decompose parses an ensign member name into its (slug, stage, cycle) parts.
// The input shape is `{workerKey}-{slug}-{stage}` where workerKey is
// "spacedock-ensign" by convention; an optional trailing `-cycleN` or `-N` rides
// after the stage. stageNames is the workflow's stage set (canonicalStages
// suffices when the workflow uses standard stages).
//
// Order of operations:
//  1. strip the worker key prefix; ok=false if absent.
//  2. peel an optional trailing -cycleN or -N from the tail.
//  3. peel the longest matching stage from the new tail.
//  4. what remains is the slug.
//
// An empty slug, missing worker prefix, or unrecognized stage all yield ok=false.
func decompose(name string, stageNames []string) decomposeResult {
	const workerKey = "spacedock-ensign-"
	if !strings.HasPrefix(name, workerKey) {
		return decomposeResult{ok: false}
	}
	rest := strings.TrimPrefix(name, workerKey)
	if rest == "" {
		return decomposeResult{ok: false}
	}

	cycle := ""
	if i := strings.LastIndex(rest, "-"); i >= 0 {
		tail := rest[i+1:]
		if strings.HasPrefix(tail, "cycle") {
			n := strings.TrimPrefix(tail, "cycle")
			if n != "" && isAllDigits(n) {
				cycle = tail
				rest = rest[:i]
			}
		} else if isAllDigits(tail) && tail != "" {
			cycle = tail
			rest = rest[:i]
		}
	}

	// Strip the longest matching stage from the tail. stageNames is sorted
	// longest-first so the first match wins.
	stage := ""
	for _, s := range stageNames {
		suffix := "-" + s
		if strings.HasSuffix(rest, suffix) {
			stage = s
			rest = strings.TrimSuffix(rest, suffix)
			break
		}
	}
	if stage == "" {
		return decomposeResult{ok: false}
	}
	if rest == "" {
		return decomposeResult{ok: false}
	}
	return decomposeResult{slug: rest, stage: stage, cycle: cycle, ok: true}
}

// isAllDigits reports whether s is non-empty and contains only ASCII digits.
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// classA flags ensigns whose entity is archived OR whose entity status is
// terminal ("done" or empty next-stage equivalent). Lingering = the entity is
// past terminal but the agent is still in the team roster.
func classA(ensigns []claudeteam.ReconcileMember, stages []string, active, archived map[string]entityRecord) []driftItem {
	var out []driftItem
	for _, m := range ensigns {
		d := decompose(m.Name, stages)
		if !d.ok {
			continue
		}
		// Archived: definitive terminal.
		if _, ok := archived[d.slug]; ok {
			out = append(out, driftItem{
				Class:  "A",
				Name:   m.Name,
				Slug:   d.slug,
				Reason: "entity archived",
			})
			continue
		}
		if rec, ok := active[d.slug]; ok {
			if rec.status == "done" {
				out = append(out, driftItem{
					Class:  "A",
					Name:   m.Name,
					Slug:   d.slug,
					Reason: "entity status=done",
				})
			}
		}
	}
	return out
}

// classB flags supersede losers — when multiple ensigns share the same
// (slug, stage) cohort, the highest cycle wins (or unsuffixed = cycle1). Losers
// emit one driftItem each pointing at the winner.
func classB(ensigns []claudeteam.ReconcileMember, stages []string) []driftItem {
	cohorts := map[string][]decomposeResult{}
	namesByDecomp := map[string]string{}
	for _, m := range ensigns {
		d := decompose(m.Name, stages)
		if !d.ok {
			continue
		}
		key := d.slug + "\x00" + d.stage
		cohorts[key] = append(cohorts[key], d)
		namesByDecomp[key+"\x00"+d.cycle] = m.Name
	}
	keys := make([]string, 0, len(cohorts))
	for k := range cohorts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var out []driftItem
	for _, k := range keys {
		members := cohorts[k]
		if len(members) < 2 {
			continue
		}
		winner := pickCycleWinner(members)
		winnerName := namesByDecomp[k+"\x00"+winner.cycle]
		for _, c := range members {
			if c.cycle == winner.cycle {
				continue
			}
			loserName := namesByDecomp[k+"\x00"+c.cycle]
			out = append(out, driftItem{
				Class:  "B",
				Name:   loserName,
				Slug:   c.slug,
				Stage:  c.stage,
				Reason: fmt.Sprintf("superseded by %s", winnerName),
			})
		}
	}
	return out
}

// pickCycleWinner returns the cohort member with the highest cycle number;
// unsuffixed (cycle == "") is treated as cycle 1.
func pickCycleWinner(members []decomposeResult) decomposeResult {
	winner := members[0]
	winnerN := cycleNumber(winner.cycle)
	for _, c := range members[1:] {
		n := cycleNumber(c.cycle)
		if n > winnerN {
			winner = c
			winnerN = n
		}
	}
	return winner
}

// cycleNumber returns the cycle integer for a decomposed cycle suffix. "" is
// cycle 1 (unsuffixed); "cycleN" or bare digits are the explicit number.
func cycleNumber(cycle string) int {
	if cycle == "" {
		return 1
	}
	s := strings.TrimPrefix(cycle, "cycle")
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return 1
}

// classC flags entities whose pr field is set and whose PR is MERGED upstream
// but whose own status is not yet terminal — the FO advancing reactively
// missed the merge.
func classC(active map[string]entityRecord, gh ghRunner) []driftItem {
	var out []driftItem
	slugs := sortedKeys(active)
	for _, slug := range slugs {
		rec := active[slug]
		if rec.pr == "" {
			continue
		}
		// A local-merge sentinel (pr=local-merge:{sha}) is not a remote PR; skip.
		if strings.HasPrefix(rec.pr, "local-merge:") {
			continue
		}
		state, err := gh(rec.pr)
		if err != nil {
			// gh unavailable or rate-limited — skip silently. The sweep is best-
			// effort; a transient gh failure should not blow up the FO's idle hook.
			continue
		}
		if state == "MERGED" && rec.status != "done" {
			out = append(out, driftItem{
				Class:  "C",
				Slug:   slug,
				PR:     rec.pr,
				Reason: fmt.Sprintf("PR merged but status=%s", rec.status),
			})
		}
	}
	return out
}

// classD flags worktrees whose branch HEAD is behind origin/next (need a
// rebase). The behind count is `git rev-list --count HEAD..origin/next` run
// inside the worktree dir.
func classD(active map[string]entityRecord, repoRoot string, git gitRunner) []driftItem {
	var out []driftItem
	slugs := sortedKeys(active)
	for _, slug := range slugs {
		rec := active[slug]
		if rec.worktree == "" {
			continue
		}
		wt := rec.worktree
		if !filepath.IsAbs(wt) {
			wt = filepath.Join(repoRoot, wt)
		}
		if info, err := os.Stat(wt); err != nil || !info.IsDir() {
			continue
		}
		out2, err := git(wt, "rev-list", "--count", "HEAD..origin/next")
		if err != nil {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(out2))
		if err != nil || n <= 0 {
			continue
		}
		out = append(out, driftItem{
			Class:    "D",
			Slug:     slug,
			Worktree: rec.worktree,
			Behind:   n,
			Reason:   fmt.Sprintf("branch behind origin/next by %d", n),
		})
	}
	return out
}

// classE flags a local main that carries commits not on origin/next. The FO's
// action is `git fetch && git reset --hard origin/next` plus a rebuild of the
// binary. The helper only detects.
func classE(repoRoot string, git gitRunner) []driftItem {
	out, err := git(repoRoot, "rev-list", "--count", "origin/next..main")
	if err != nil {
		return nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil || n <= 0 {
		return nil
	}
	return []driftItem{{
		Class:  "E",
		Ahead:  n,
		Reason: fmt.Sprintf("local main carries %d commits not on origin/next; reset main->origin/next", n),
	}}
}

// sortedKeys returns the entity-record map's keys in sorted order so the drift
// array is deterministic across runs.
func sortedKeys(m map[string]entityRecord) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// sortDrift sorts drift entries by (class, slug, name) so the JSON output is
// stable for tests and golden-comparable across runs.
func sortDrift(drift []driftItem) {
	sort.Slice(drift, func(i, j int) bool {
		a, b := drift[i], drift[j]
		if a.Class != b.Class {
			return a.Class < b.Class
		}
		if a.Slug != b.Slug {
			return a.Slug < b.Slug
		}
		return a.Name < b.Name
	})
}
