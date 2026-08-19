// ABOUTME: --boot section printer — MODS, ID_STYLE/NEXT_ID, ORPHANS, PR_STATE,
// ABOUTME: DISPATCHABLE, TEAM_STATE, matching print_boot and its probes.
package status

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
	"github.com/spacedock-dev/spacedock/internal/safehouse"
)

// teamStateNeutralHint is the boot TEAM_STATE present:false hint on a host with
// no team-state probe wired (Codex/bare). It carries no Claude-specific advice —
// the Claude present:false hint (claudeteam.PresentFalseHint) is supplied by the
// seam only when a probe is present.
const teamStateNeutralHint = "no active team runtime detected"

// orphan describes a worktree-backed entity for the ORPHANS section.
type orphan struct {
	id, slug, worktree, dirExists, branchExists string
}

// scanOrphans returns entities with a non-empty worktree field plus dir/branch
// existence info from `git worktree list --porcelain`. Matches scan_orphans.
func scanOrphans(entities []*entity, gitRoot string) []orphan {
	worktreePaths := map[string]bool{}
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = gitRoot
	out, err := cmd.Output()
	if err == nil {
		var current string
		for _, line := range strings.Split(string(out), "\n") {
			switch {
			case strings.HasPrefix(line, "worktree "):
				current = strings.TrimSpace(line[len("worktree "):])
			case strings.HasPrefix(line, "branch ") && current != "":
				worktreePaths[realpathOf(current)] = true
			case line == "":
				current = ""
			}
		}
	}

	var orphans []orphan
	for _, e := range entities {
		wt := e.fields["worktree"]
		if wt == "" {
			continue
		}
		// os.path.join(git_root, wt): an absolute wt discards git_root, so an
		// absolute worktree path is probed as-is. filepath.Join would instead
		// graft it under git_root and miss the existing dir. PyJoin matches the
		// oracle's os.path.join absolute-component-reset semantics.
		dirPath := PyJoin(gitRoot, wt)
		dirExists := "no"
		if st, err := os.Stat(dirPath); err == nil && st.IsDir() {
			dirExists = "yes"
		}
		branchExists := "no"
		if worktreePaths[realpathOf(dirPath)] {
			branchExists = "yes"
		}
		orphans = append(orphans, orphan{
			id: e.fields["id"], slug: e.fields["slug"], worktree: wt,
			dirExists: dirExists, branchExists: branchExists,
		})
	}
	return orphans
}

// prResult describes a PR-bearing entity's resolved state.
type prResult struct {
	id, slug, pr, state string
}

// checkPRStates returns (status, results) for entities with a non-empty pr and
// non-terminal status. Matches check_pr_states. status is "none",
// "gh not available", "ok", or (identify mode) "local". In identify mode the boot
// is a side-effect-free local read: it renders the stored `pr:` field as a local
// mirror (state "local", not-gh-checked) and makes NO `gh pr view` call — the live
// OPEN/MERGED/CLOSED state is filled in at «engage»'s convergence, not the greet.
func checkPRStates(entities []*entity, stages []Stage, e env, identify bool) (string, []prResult) {
	stageByName := map[string]Stage{}
	for _, s := range stages {
		stageByName[s.Name] = s
	}
	var prEntities []*entity
	for _, ent := range entities {
		if ent.fields["pr"] == "" {
			continue
		}
		if st, ok := stageByName[ent.fields["status"]]; ok && st.terminal {
			continue
		}
		prEntities = append(prEntities, ent)
	}
	if len(prEntities) == 0 {
		return "none", nil
	}

	if identify {
		results := make([]prResult, 0, len(prEntities))
		for _, ent := range prEntities {
			results = append(results, prResult{id: ent.fields["id"], slug: ent.fields["slug"], pr: ent.fields["pr"], state: "local"})
		}
		return "local", results
	}

	ghPath := lookupExecutable("gh", e.get("PATH"))
	if ghPath == "" {
		return "gh not available", nil
	}

	var results []prResult
	for _, ent := range prEntities {
		pr := strings.TrimPrefix(ent.fields["pr"], "#")
		state := "ERROR"
		out, err := exec.Command("gh", "pr", "view", pr, "--json", "state", "--jq", ".state").Output()
		if err == nil {
			state = strings.ToUpper(strings.TrimSpace(string(out)))
		}
		results = append(results, prResult{id: ent.fields["id"], slug: ent.fields["slug"], pr: ent.fields["pr"], state: state})
	}
	return "ok", results
}

// lookupExecutable returns the first executable named `name` on pathStr, or "".
// Matches the inline PATH scan in check_pr_states.
func lookupExecutable(name, pathStr string) string {
	for _, d := range filepath.SplitList(pathStr) {
		candidate := filepath.Join(d, name)
		if st, err := os.Stat(candidate); err == nil && st.Mode().IsRegular() && st.Mode()&0o111 != 0 {
			return candidate
		}
	}
	return ""
}

// bootData holds the gathered boot-section material the text and JSON renderers
// both consume, so the two output forms read from one source of truth.
type bootData struct {
	hooks        map[string][]string
	idStyle      string
	nextID       string
	orphans      []orphan
	prStatus     string
	prResults    []prResult
	dispatchable []dispatchable
	readyGates   []*entity
	teamPresent  bool
	teamHint     string
	// State backend: split-root when entityDir diverges from definitionDir (a
	// non-empty README state: field), else single-root. The dirs are the absolute
	// I/O spellings, so the FO reads the state checkout from boot with no inference.
	stateBackend     string
	definitionDir    string
	entityDir        string
	entityDirPresent bool
	// stateRemote is the state checkout's remote-sync availability, populated only
	// under split-root: "origin" when the checkout has a named origin remote,
	// "none" when it does not (local-only — state is not remotely synced). Empty
	// for single-root, where remote sync does not apply.
	stateRemote string
	// Sandbox posture: the three-way safehouse state (enabled / available, not
	// enabled / unavailable) computed from a .safehouse profile at the repo root and
	// whether the safehouse binary resolves on PATH, so the operator sees the
	// execution-isolation posture before dispatching work.
	sandbox string
	// Identify mode (the FO's opt-in local-identify boot). When set, the record
	// folds the workflow discovery result and the stage taxonomy into the same
	// envelope, PR_STATE is the local `pr:` mirror (checkPRStates skips gh), and the
	// whole boot is a side-effect-free local read. discovery/stages/ready_gates
	// are appended AFTER the existing key set so every prior key's order is preserved.
	identify  bool
	discovery []string
	stages    []Stage
}

// gatherBoot runs every boot probe once and returns the result. NEXT_ID is
// minted here (timestamp-dependent for sd-b32); on a minting error it returns
// the error after the caller has emitted the stderr diagnostic. transcriptProbe
// is the boot guard's receipt-write seam (nil on a non-Claude host); it carries
// no ~/.claude read itself — see claudeteam.TranscriptProbe.
func gatherBoot(probe claudeteam.TeamStateProbe, transcriptProbe claudeteam.TranscriptProbe, entities []*entity, stages []Stage, definitionDir, entityDir, gitRoot, idStyle string, e env, stderr io.Writer, identify bool) (*bootData, error) {
	d := &bootData{idStyle: idStyle, hooks: scanMods(definitionDir), identify: identify}

	if idStyle == "slug" {
		d.nextID = "n/a (id-style: slug)"
	} else {
		id, err := computeNextID(definitionDir, entityDir, idStyle, "", "", e, stderr)
		if err != nil {
			fmt.Fprintf(stderr, "Error: %s\n", err)
			return nil, err
		}
		d.nextID = id
	}

	d.orphans = scanOrphans(entities, gitRoot)
	d.prStatus, d.prResults = checkPRStates(entities, stages, e, identify)
	d.dispatchable = computeDispatchable(entities, stages)
	// Identify mode folds the two hand-issued pre-greet reads — workflow discovery
	// and the stage taxonomy — into this one record. Both are local reads (a
	// filesystem walk; the already-parsed stages), so the boot stays side-effect-free.
	if identify {
		d.discovery = discoverWorkflows(gitRoot)
		d.stages = stages
		materializeGateReadiness(entities, stages)
		d.readyGates = computeReadyGates(entities, stages)
	}
	// TEAM_STATE comes from the host-supplied probe. HOME resolution stays generic
	// here; only the ~/.claude read moves into the Claude seam. The hint for both
	// the present and absent cases is resolved here so the renderers carry no
	// host-specific string: with a probe wired, present:false uses the Claude seam's
	// PresentFalseHint; a nil probe (a non-Claude host) yields a host-neutral line.
	if probe != nil {
		home := e.get("HOME")
		if home == "" {
			home = os.Getenv("HOME")
		}
		d.teamPresent, d.teamHint, _ = probe(home, time.Now())
		if !d.teamPresent {
			d.teamHint = claudeteam.PresentFalseHint
		}
	} else {
		d.teamHint = teamStateNeutralHint
	}

	d.definitionDir = definitionDir
	d.entityDir = entityDir
	if entityDir != definitionDir {
		d.stateBackend = "split-root"
		// Remote-sync availability is read only under split-root: the FO uses it to
		// know whether the state checkout can push/pull origin or is local-only.
		if stateHasOrigin(entityDir) {
			d.stateRemote = "origin"
		} else {
			d.stateRemote = "none"
		}
	} else {
		d.stateBackend = "single-root"
	}
	if info, err := os.Stat(entityDir); err == nil && info.IsDir() {
		d.entityDirPresent = true
	}

	// SANDBOX: boot reports THIS session's posture, so it renders SessionState —
	// is this process sandboxed? — not whether a launch from here would wrap. That
	// distinction is why the .safehouse profile is no longer read here: a profile
	// is a launch fact, and boot is not a launch. The sandbox is detected from the
	// request env via the shared registry, and the safehouse binary is resolved
	// against the request PATH via the existing executable scan — no exec, no live
	// host CLI.
	//
	// This field lands in the First Officer's durable boot record, which is why it
	// mattered most that it inverted: inside the sandbox, safehouse is off PATH
	// precisely BECAUSE the wrap already happened, so the old availability-only
	// render wrote "unavailable (safehouse not on PATH)" into machine-read evidence
	// captured from inside a sandbox.
	available := lookupExecutable("safehouse", e.get("PATH")) != ""
	insideName, inside := safehouse.Inside(e.get)
	d.sandbox = safehouse.SessionState(insideName, inside, available)

	// Boot's own side effect (see docs/runtime-support.md, "Boot guard at the
	// compaction boundary"): a session boot receipt, host scratch only. This is
	// the mechanism the shipped `«state.boot»()` gains — workflow/entity state
	// stays read-only, unaffected by this call either way.
	writeBootReceipt(e, gitRoot, transcriptProbe, stderr, time.Now())
	return d, nil
}

// printBoot writes all boot sections in order. Matches print_boot.
func printBoot(probe claudeteam.TeamStateProbe, transcriptProbe claudeteam.TranscriptProbe, w io.Writer, entities []*entity, stages []Stage, definitionDir, entityDir, gitRoot, idStyle string, e env, stderr io.Writer, identify bool) error {
	d, err := gatherBoot(probe, transcriptProbe, entities, stages, definitionDir, entityDir, gitRoot, idStyle, e, stderr, identify)
	if err != nil {
		return err
	}

	// MODS
	if len(d.hooks) == 0 {
		fmt.Fprintln(w, "MODS: none")
	} else {
		fmt.Fprintln(w, "MODS")
		points := make([]string, 0, len(d.hooks))
		for p := range d.hooks {
			points = append(points, p)
		}
		sort.Strings(points)
		for _, point := range points {
			mods := append([]string(nil), d.hooks[point]...)
			sort.Strings(mods)
			fmt.Fprintf(w, "%s: %s\n", point, strings.Join(mods, ", "))
		}
	}

	// ID_STYLE / NEXT_ID
	fmt.Fprintf(w, "ID_STYLE: %s\n", d.idStyle)
	fmt.Fprintf(w, "NEXT_ID: %s\n", d.nextID)
	if d.idStyle == "sd-b32" {
		fmt.Fprintf(w, "MIN_PREFIX: %d\n", sdB32MinPrefix)
	}

	// ORPHANS
	if len(d.orphans) == 0 {
		fmt.Fprintln(w, "ORPHANS: none")
	} else {
		fmt.Fprintln(w, "ORPHANS")
		row := func(a, b, c, d, e string) string {
			return padRight(a, 6) + " " + padRight(b, 30) + " " + padRight(c, 43) + " " + padRight(d, 11) + " " + e
		}
		fmt.Fprintln(w, row("ID", "SLUG", "WORKTREE", "DIR_EXISTS", "BRANCH_EXISTS"))
		for _, o := range d.orphans {
			fmt.Fprintln(w, row(o.id, o.slug, o.worktree, o.dirExists, o.branchExists))
		}
	}

	// PR_STATE. Identify mode renders the local `pr:` mirror (state "local") under a
	// labeled banner so the reader knows the live gh state is filled in at «engage».
	switch d.prStatus {
	case "none":
		fmt.Fprintln(w, "PR_STATE: none")
	case "gh not available":
		fmt.Fprintln(w, "PR_STATE: gh not available")
	case "local":
		fmt.Fprintln(w, "PR_STATE (local view — not gh-checked)")
		row := func(a, b, c, d string) string {
			return padRight(a, 6) + " " + padRight(b, 30) + " " + padRight(c, 8) + " " + d
		}
		fmt.Fprintln(w, row("ID", "SLUG", "PR", "STATE"))
		for _, r := range d.prResults {
			fmt.Fprintln(w, row(r.id, r.slug, r.pr, r.state))
		}
	default:
		fmt.Fprintln(w, "PR_STATE")
		row := func(a, b, c, d string) string {
			return padRight(a, 6) + " " + padRight(b, 30) + " " + padRight(c, 8) + " " + d
		}
		fmt.Fprintln(w, row("ID", "SLUG", "PR", "STATE"))
		for _, r := range d.prResults {
			fmt.Fprintln(w, row(r.id, r.slug, r.pr, r.state))
		}
	}

	// DISPATCHABLE
	fmt.Fprintln(w, "DISPATCHABLE")
	printNextTable(w, entities, stages, nil, false)

	// TEAM_STATE. The present/absent hints are both resolved in gatherBoot so this
	// renderer carries no host-specific string.
	fmt.Fprintln(w, "TEAM_STATE")
	if d.teamPresent {
		fmt.Fprintln(w, "present: true")
	} else {
		fmt.Fprintln(w, "present: false")
	}
	fmt.Fprintf(w, "hint: %s\n", d.teamHint)

	// STATE_BACKEND. The remote clause is appended only under split-root: origin
	// when the state checkout can push/pull, else a local-only marker so the FO
	// sees state is not remotely synced. Single-root omits the clause entirely.
	remoteClause := ""
	switch d.stateRemote {
	case "origin":
		remoteClause = ", remote: origin"
	case "none":
		remoteClause = ", remote: none — state not remotely synced"
	}
	fmt.Fprintf(w, "STATE_BACKEND: %s (entity_dir: %s, present: %t%s)\n",
		d.stateBackend, d.entityDir, d.entityDirPresent, remoteClause)

	// SANDBOX: appended last so every prior section's order is preserved.
	fmt.Fprintf(w, "SANDBOX: %s\n", d.sandbox)

	// Identify mode folds discovery + the stage taxonomy in, appended after every
	// existing section so their order is untouched.
	if d.identify {
		fmt.Fprintln(w, "DISCOVERY")
		for _, wf := range d.discovery {
			fmt.Fprintln(w, wf)
		}
		fmt.Fprintln(w, "STAGES")
		for _, s := range d.stages {
			fmt.Fprintln(w, s.Name)
		}
	}
	return nil
}
