// ABOUTME: dispatch reconcile — computes a five-class drift report (lingering /
// ABOUTME: superseded / un-advanced PR / stale branch / local main drift) from disk.
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

// The five drift-class names. These are the canonical class vocabulary: the
// emitted drift[].class value, the --include accepted tokens, and the FO dispatch
// contract event-loop step-0 action mapping all name this identical set. The
// contractlint dual-extraction lint binds this var to the contract step-0 token
// so neither side can rename, add, or drop a class without the other.
const (
	classLingering      = "lingering"
	classSuperseded     = "superseded"
	classUnadvancedPR   = "un-advanced-pr"
	classStaleBranch    = "stale-branch"
	classLocalMainDrift = "local-main-drift"
)

// driftClasses is the canonical ordered class set. parseInclude validates against
// it, Reconcile gates each detector by it, and the AC-2 lint reads it as the
// single helper-side source of the emitted class vocabulary.
var driftClasses = []string{
	classLingering,
	classSuperseded,
	classUnadvancedPR,
	classStaleBranch,
	classLocalMainDrift,
}

// canonicalStages are the workflow stage names the decomposer falls back to when
// a workflow's own stage names cannot be derived from its README. The entity's
// design names these as the "five-canonical ensign stages" so a workflow with
// custom stages still decomposes cleanly without code changes — the per-workflow
// README scan augments this set.
var canonicalStages = []string{"backlog", "ideation", "implementation", "validation", "done"}

// reconcileOpts captures the resolved CLI flags for one reconcile run. roster
// is the team-roster loader the helper calls — in production this is
// claudeteam.LoadReconcileTeam, which owns the only ~/.claude/teams read; tests
// inject a stub so the host-neutrality scan stays clean. sessionID is the
// current lead session (from $CLAUDE_CODE_SESSION_ID); auto-discovery matches it
// against config.leadSessionId when teamName is empty. Tests inject a fixture
// UUID; the empty value drives the degrade-to-git-only path.
type reconcileOpts struct {
	workflowDir string
	teamName    string
	sessionID   string
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
type rosterLoader func(home, teamName, sessionID string) (claudeteam.ReconcileTeamState, error)

// ghRunner runs `gh pr view {N} --json state --jq .state` returning the state
// string (e.g. "MERGED", "OPEN"). Tests inject a stub; production passes
// ghRunnerExec.
type ghRunner func(prRef string) (string, error)

// GhRunner is the exported alias of ghRunner so callers outside this package (the
// `state sweep` verb) can inject a stub for tests and pass GhRunnerExec in
// production, sharing the one merged-detection probe.
type GhRunner = ghRunner

// GhRunnerExec is the exported production gh probe `state sweep` passes.
func GhRunnerExec(prRef string) (string, error) { return ghRunnerExec(prRef) }

// gitRunner runs `git -C dir args...` returning combined output and error.
// Tests inject a stub when they need a fixture-controlled answer; production
// passes gitRunnerExec.
type gitRunner func(dir string, args ...string) (string, error)

// ghRunnerExec calls `gh pr view PR --json state --jq .state` and returns the
// trimmed, uppercased state. This mirrors boot.go's PR_STATE probe invocation
// exactly (same flags, same leading-"#" trim, same --jq extraction) rather than
// parsing the `--json state` envelope as a second JSON layer — the two probes
// must agree on what "gh can answer" means, since they run against the same `gh`
// and the FO trusts both.
func ghRunnerExec(prRef string) (string, error) {
	pr := strings.TrimPrefix(prRef, "#")
	out, err := exec.Command("gh", "pr", "view", pr, "--json", "state", "--jq", ".state").Output()
	if err != nil {
		return "", err
	}
	return strings.ToUpper(strings.TrimSpace(string(out))), nil
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
		sessionID:   os.Getenv("CLAUDE_CODE_SESSION_ID"),
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
// classes; otherwise the value is a comma-separated subset of the descriptive
// class names (driftClasses). An unknown class is a usage error (exit 2). Returned
// map sets true for included classes only.
func parseInclude(raw string) (map[string]bool, error) {
	all := map[string]bool{}
	for _, c := range driftClasses {
		all[c] = true
	}
	if raw == "" {
		return all, nil
	}
	out := map[string]bool{}
	for _, c := range strings.Split(raw, ",") {
		c = strings.TrimSpace(c)
		if !all[c] {
			return nil, fmt.Errorf("unknown --include class: %q (expected %s)", c, strings.Join(driftClasses, ","))
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
	Trunk    string `json:"trunk,omitempty"`
	Owned    bool   `json:"owned,omitempty"`
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

	team, err := opts.roster(opts.home, opts.teamName, opts.sessionID)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}

	// An empty TeamName is the degrade-to-git-only sentinel: no session-scoped
	// team resolved (bare reconcile with no current-session match), so the
	// roster-derived classes A/B/C are suppressed — a non-session-matched config
	// could be a stale prior session's or a parallel live session's team, and
	// must never be trusted as our roster. Only the session-independent git/
	// filesystem classes D/E run. This is not an error: the sweep ran (exit 0).
	rosterTrusted := team.TeamName != ""
	if !rosterTrusted {
		fmt.Fprintln(stderr,
			"note: no session-scoped team resolved; reporting git/filesystem drift only (roster reconciliation needs a team identity — pass --team-name)")
	}

	stateRoot, err := splitRootStateCheckout(opts.workflowDir)
	if err != nil {
		fmt.Fprintf(stderr, "error: cannot resolve split-root state checkout: %v\n", err)
		return 1
	}
	if stateRoot == "" {
		stateRoot = opts.workflowDir
	}
	active := loadEntityFrontmatter(activeEntityDir(stateRoot))
	archived := loadEntityFrontmatter(archivedEntityDir(stateRoot))

	stageNames := readStageNames(opts.workflowDir)

	// The integration trunk is resolved ONCE from the workflow's declared
	// top-level trunk: key (defaulting to main) and passed into the git-hygiene
	// detectors. The detectors carry no trunk literal — trunk knowledge lives in
	// the one resolveTrunk source, not scattered inside roster-helper detectors.
	trunk := resolveTrunk(opts.workflowDir)

	var drift []driftItem

	ensigns := filterEnsigns(team.Members)
	if rosterTrusted && opts.include[classLingering] {
		drift = append(drift, classA(ensigns, stageNames, active, archived)...)
	}
	if rosterTrusted && opts.include[classSuperseded] {
		drift = append(drift, classB(ensigns, stageNames)...)
	}
	if rosterTrusted && opts.include[classUnadvancedPR] {
		drift = append(drift, classC(active, opts.gh)...)
	}
	if opts.include[classStaleBranch] {
		owned := ownedSlugs(rosterTrusted, ensigns, stageNames, active, archived)
		drift = append(drift, classD(active, repoRoot, trunk, owned, opts.git)...)
	}
	if opts.include[classLocalMainDrift] {
		drift = append(drift, classE(repoRoot, trunk, opts.git)...)
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
// frontmatter top-level key (no body parsing). id carries the entity's stored id
// so a capped worker name (whose slug component was replaced by an id-prefix at
// dispatch build) resolves back to its entity by id-prefix, not string-split.
type entityRecord struct {
	slug     string
	id       string
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
			id:       fm["id"],
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

// resolveSlugToken maps a decomposed name token to its entity slug. For an
// uncapped name the token IS the slug, so an exact match against the active or
// archived maps wins immediately (the common case). For a capped sd-b32 name the
// token is a fixed-length id-prefix; when no slug matches exactly, the token is
// resolved by HasPrefix against the active+archived id set — the same resolution
// `status --resolve prefix:` performs. A unique id-prefix match yields the real
// slug; an ambiguous (≥2) or absent match leaves the token unresolved (ok=false),
// the conservative outcome that avoids a false Class-A against a live sibling.
func resolveSlugToken(token string, active, archived map[string]entityRecord) (string, bool) {
	if _, ok := active[token]; ok {
		return token, true
	}
	if _, ok := archived[token]; ok {
		return token, true
	}
	// id-prefix resolution: the token must be a non-trivial prefix of exactly one
	// stored id. sdB32MinPrefix guards against a too-short token matching many ids.
	if len(token) < status.SDB32MinPrefix {
		return "", false
	}
	matched := ""
	count := 0
	for _, m := range []map[string]entityRecord{active, archived} {
		for slug, rec := range m {
			if rec.id != "" && strings.HasPrefix(rec.id, token) {
				matched = slug
				count++
			}
		}
	}
	if count != 1 {
		return "", false
	}
	return matched, true
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
		// Resolve the decomposed token to a real slug (exact-slug-first, then
		// sd-b32 id-prefix). An unresolved token leaves the member unclassified.
		slug, ok := resolveSlugToken(d.slug, active, archived)
		if !ok {
			continue
		}
		// Archived: definitive terminal.
		if _, ok := archived[slug]; ok {
			out = append(out, driftItem{
				Class:  classLingering,
				Name:   m.Name,
				Slug:   slug,
				Reason: "entity archived",
			})
			continue
		}
		if rec, ok := active[slug]; ok {
			if rec.status == "done" {
				out = append(out, driftItem{
					Class:  classLingering,
					Name:   m.Name,
					Slug:   slug,
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
				Class:  classSuperseded,
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
				Class:  classUnadvancedPR,
				Slug:   slug,
				PR:     rec.pr,
				Reason: fmt.Sprintf("PR merged but status=%s", rec.status),
			})
		}
	}
	return out
}

// sweptEntity is one merged-but-not-terminalized entity in the `state sweep`
// report — the state-repo's read-only view of the un-advanced-PR drift class.
type sweptEntity struct {
	Slug   string `json:"slug"`
	PR     string `json:"pr"`
	Reason string `json:"reason"`
}

// sweepResult is the `state sweep` JSON envelope. swept is the merged-but-not-
// terminalized set; an empty sweep encodes an empty array, never null.
type sweepResult struct {
	Command     string        `json:"command"`
	StateBranch string        `json:"state_branch,omitempty"`
	Swept       []sweptEntity `json:"swept"`
	Reason      string        `json:"reason"`
	Gh          string        `json:"gh,omitempty"`
	Next        string        `json:"next,omitempty"`
}

// Sweep is the read-only `state sweep` computation: the entities whose code PR has
// MERGED but whose state has not been terminalized. It reuses classC's un-advanced-
// PR detection over the split-root state checkout's active entities — NO second
// merged-detection path — and makes no commit, push, or mutation. The gh probe is
// injected so tests pin a deterministic merged-state; production passes
// GhRunnerExec. An inline (single-root) workflow sweeps its own dir. Exit 0 the
// sweep ran (set may be empty or populated); 1 setup failure (no README/state).
//
// classC itself stays best-effort — a gh error silently skips that entity, so the
// idle hook never blows up on a transient failure. Sweep wraps the injected probe
// to COUNT calls and errors around that unchanged behavior: when every PR-pending
// entity's probe errored, "0 entity(ies)" would be indistinguishable from a real
// empty sweep, so Sweep reports merge state UNKNOWN instead (D2).
func Sweep(workflowDir string, gh GhRunner, jsonOut bool, stdout, stderr io.Writer) int {
	if info, err := os.Stat(workflowDir); err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "spacedock state sweep: workflow directory not found: %s\n", workflowDir)
		return 1
	}
	stateRoot, err := splitRootStateCheckout(workflowDir)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock state sweep: cannot resolve state checkout: %v\n", err)
		return 1
	}
	if stateRoot == "" {
		stateRoot = workflowDir
	}
	active := loadEntityFrontmatter(activeEntityDir(stateRoot))

	probeTotal, probeErrs := 0, 0
	counting := func(prRef string) (string, error) {
		probeTotal++
		state, err := gh(prRef)
		if err != nil {
			probeErrs++
		}
		return state, err
	}
	drift := classC(active, counting)

	swept := make([]sweptEntity, 0, len(drift))
	for _, d := range drift {
		swept = append(swept, sweptEntity{Slug: d.Slug, PR: d.PR, Reason: d.Reason})
	}
	branch, _ := status.StateBranch(workflowDir)

	var reason, ghField, next string
	switch {
	case probeTotal > 0 && probeErrs == probeTotal:
		reason = "merge state UNKNOWN — gh unavailable; sweep skipped, not empty."
		ghField = "unavailable"
	case len(swept) > 0:
		reason = fmt.Sprintf("%d entity(ies) merged but not yet terminalized.", len(swept))
		next = sweepNextStep(workflowDir)
	default:
		reason = fmt.Sprintf("%d entity(ies) merged but not yet terminalized.", len(swept))
	}

	if jsonOut {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		enc.Encode(sweepResult{
			Command: "state sweep", StateBranch: branch, Swept: swept, Reason: reason, Gh: ghField, Next: next,
		})
		return 0
	}
	fmt.Fprintln(stdout, reason)
	if next != "" {
		fmt.Fprintln(stdout, next)
	}
	for _, s := range swept {
		fmt.Fprintf(stdout, "  %s (PR %s): %s\n", s.Slug, s.PR, s.Reason)
	}
	return 0
}

// sweepNextStep names the FO's next action for a non-empty sweep: the registered
// startup-hook mod file(s) to advance each entity per, from the same hookPoint ->
// mod-name scan the boot MODS-REPORT uses (status.ScanMods). It points at the mod
// FILE, never a procedure — the shipped mods/pr-merge.md advances an entity
// directly while a local per-workflow pr-merge mod can delegate to sentinel +
// merge guard, and the binary has no way to know which one applies without
// picking a side. No startup mod registered falls back to a generic _mods/
// pointer.
func sweepNextStep(workflowDir string) string {
	hooks := status.ScanMods(workflowDir)["startup"]
	if len(hooks) == 0 {
		return "next: advance each per the workflow's startup-hook advancement (_mods/)."
	}
	paths := make([]string, len(hooks))
	for i, h := range hooks {
		paths[i] = fmt.Sprintf("_mods/%s.md", h)
	}
	return fmt.Sprintf("next: advance each per the workflow's startup-hook advancement (%s).", strings.Join(paths, ", "))
}

// classD flags worktrees whose branch HEAD is behind origin/{trunk}. The remedy
// is ownership-gated: a `pull --rebase` is prescribed ONLY when the worktree's
// entity slug is in `owned` — the set of slugs the CURRENT trusted roster's
// ensigns resolve to (the same decompose signal classes A/B/C use). A worktree
// whose entity is not owned (a peer session's branch, or any worktree when the
// roster is untrusted) is still reported, but report-only with Owned:false and
// no rebase/pull verb — reconcile never mutates a branch the running session
// does not own. The behind count is `git rev-list --count HEAD..origin/{trunk}`
// run inside the worktree dir; the trunk is resolved once in Reconcile and
// passed in — classD carries no trunk literal.
func classD(active map[string]entityRecord, repoRoot, trunk string, owned map[string]bool, git gitRunner) []driftItem {
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
		out2, err := git(wt, "rev-list", "--count", "HEAD..origin/"+trunk)
		if err != nil {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(out2))
		if err != nil || n <= 0 {
			continue
		}
		item := driftItem{
			Class:    classStaleBranch,
			Slug:     slug,
			Worktree: rec.worktree,
			Behind:   n,
			Trunk:    trunk,
		}
		if owned[slug] {
			item.Owned = true
			item.Reason = fmt.Sprintf("branch behind origin/%s by %d; pull --rebase (owned)", trunk, n)
		} else {
			item.Reason = fmt.Sprintf("branch behind origin/%s by %d; peer-owned or un-owned — reporting only, not rebasing", trunk, n)
		}
		out = append(out, item)
	}
	return out
}

// ownedSlugs is the set of entity slugs the CURRENT trusted roster's ensigns
// resolve to — built from the trusted roster's ensign member names by the same
// decompose + resolveSlugToken path classes A/B/C use. When the roster is
// untrusted the set is empty, so every class-D worktree is report-only: without
// a trusted roster we cannot prove we own anything.
func ownedSlugs(rosterTrusted bool, ensigns []claudeteam.ReconcileMember, stages []string, active, archived map[string]entityRecord) map[string]bool {
	owned := map[string]bool{}
	if !rosterTrusted {
		return owned
	}
	for _, m := range ensigns {
		d := decompose(m.Name, stages)
		if !d.ok {
			continue
		}
		slug, ok := resolveSlugToken(d.slug, active, archived)
		if !ok {
			continue
		}
		owned[slug] = true
	}
	return owned
}

// classE classifies local main against origin/{trunk} by counting BOTH
// directions — ahead (`origin/{trunk}..main`, commits main has that origin
// lacks, i.e. unpushed) and behind (`main..origin/{trunk}`, commits origin has
// that main lacks). The remedy is direction-aware and never destructive:
//
//   - behind only (ahead 0, behind >0): an `ff-merge` of origin/{trunk} into
//     main, which git refuses (non-zero exit, no mutation) if main has diverged —
//     so even a race between detection and action cannot lose a commit.
//   - ahead/unpushed (ahead >0, behind 0): report only; the FO pushes on the
//     captain's word. The reason NEVER prescribes a reset.
//   - diverged (ahead >0, behind >0): report only; manual reconcile. No reset.
//
// The helper only detects and prescribes. The trunk is resolved once in
// Reconcile and passed in; the drift item carries it so the FO remedy reads
// {drift.trunk} from JSON.
func classE(repoRoot, trunk string, git gitRunner) []driftItem {
	aheadOut, err := git(repoRoot, "rev-list", "--count", "origin/"+trunk+"..main")
	if err != nil {
		return nil
	}
	ahead, err := strconv.Atoi(strings.TrimSpace(aheadOut))
	if err != nil {
		return nil
	}
	behindOut, err := git(repoRoot, "rev-list", "--count", "main..origin/"+trunk)
	if err != nil {
		return nil
	}
	behind, err := strconv.Atoi(strings.TrimSpace(behindOut))
	if err != nil {
		return nil
	}
	if ahead == 0 && behind == 0 {
		return nil
	}
	item := driftItem{Class: classLocalMainDrift, Ahead: ahead, Behind: behind, Trunk: trunk}
	switch {
	case ahead == 0: // behind only
		item.Reason = fmt.Sprintf("local main behind origin/%s by %d; ff-merge main<-origin/%s", trunk, behind, trunk)
	case behind == 0: // ahead / unpushed
		item.Reason = fmt.Sprintf("local main ahead of origin/%s by %d (unpushed); push when ready — no auto-rewrite", trunk, ahead)
	default: // diverged
		item.Reason = fmt.Sprintf("local main diverged from origin/%s (ahead %d, behind %d); manual reconcile — no auto-rewrite", trunk, ahead, behind)
	}
	return []driftItem{item}
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
