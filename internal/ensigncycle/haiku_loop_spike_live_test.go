//go:build live

// ABOUTME: Live drive of a bare `claude --model haiku -p` FO following a hand-authored
// ABOUTME: simplified loop (no FO contract) over a throwaway split-root fixture; grades durable state + stream.
package ensigncycle

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The 0205 gating spike harness. It launches a BARE `claude --model haiku -p`
// (NOT `spacedock claude` — no `--agent`, the FO contract is NOT loaded) with the
// simplified loop's prose-functions as the `-p` prompt, drives one throwaway
// split-root fixture entity through dispatch->gate->merge->terminal, and grades the
// durable on-disk end-state PLUS the captured tool-call stream — never the FO's
// narration.
//
// It reuses the proven internal/ensigncycle substrate VERBATIM: isolatedClaudeEnv
// (clean HOME + OAuth/API-key auth, t.Skip when neither), withBinaryOnPath (the
// hand-loop's prose-functions call `spacedock`, so the built binary stays on the
// child PATH even though the launched executable is bare `claude`), the
// streamWatcher quiet-budget liveness (drainToExit), extractClaudeFinalMessage, and
// cmdPoller. The ONLY new surface is the launch target (bare `claude`, not the
// front door) and the hand-loop prompt — exactly the two NEW mechanisms AC-1/AC-2
// spiked.
//
// Captain-ratified residency (locked at the ideation gate): the «gate» step routes
// the verdict to a PER-JUDGMENT bare Agent(model="opus", …) BLOCKING call, NOT a
// standing _mods/level-3 mod. The grader traces the gate verdict's provenance to
// that opus sub-call in the captured stream (AC-3).
//
// Scope: this file BUILDS the harness and proves it with ONE clean full drive. N>=3
// and the AC-5 must-build / routing-boundary classification + the AC-4 5-failure-mode
// triage are operated by the FO at validation, not here. driveHaikuLoopOnce is
// factored so a count loop is a trivial wrapper the validation lane adds.

// haikuLoopFixture is the throwaway split-root fixture the drive runs against: a
// definitionDir holding the README (which declares `state: .spacedock-state`) and a
// stateDir (the .spacedock-state checkout) holding the one entity. Split-root means
// the entity body + the path-scoped commit live in stateDir, NOT beside the README —
// the same contract the real workflow uses, so the hand-loop's «state.commit» body
// exercises the real path-scoped-commit shape.
type haikuLoopFixture struct {
	definitionDir string // holds README.md
	stateDir      string // .spacedock-state checkout, holds the entity
	entityPath    string // stateDir/widget.md
	entitySlug    string
}

const haikuLoopSlug = "widget"

// haikuBuiltMarker is the note line the impl worker appends to the entity body and
// the durable grade greps for — proof the dispatched worker actually did work, not
// just that the FO moved a frontmatter field. It is the fixture's AC-1.
const haikuBuiltMarker = "WIDGET-BUILT"

// writeHaikuLoopFixture stands up the throwaway split-root fixture and git-inits both
// roots (the definition root for the README, the state checkout as its own index for
// path-scoped entity commits). It returns the fixture handle. The entity starts at
// the dispatchable `implementation` stage; the workflow is implementation(initial) ->
// validation(gate) -> done(terminal), so the drive has a DISTINCT dispatch step
// (implementation), a DISTINCT gate step (validation, where the L3 verdict routes),
// and a DISTINCT terminalize+archive step (done).
func writeHaikuLoopFixture(t *testing.T, root string) haikuLoopFixture {
	t.Helper()
	stateDir := filepath.Join(root, ".spacedock-state")
	writeFile(t, filepath.Join(root, "README.md"), haikuLoopReadme())
	entityPath := filepath.Join(stateDir, haikuLoopSlug+".md")
	writeFile(t, entityPath, haikuLoopEntity())

	// The state checkout is its own git index (split-root): path-scoped entity
	// commits land here, separate from the definition root's README commit.
	gitInit(t, stateDir)
	gitInit(t, root)

	return haikuLoopFixture{
		definitionDir: root,
		stateDir:      stateDir,
		entityPath:    entityPath,
		entitySlug:    haikuLoopSlug,
	}
}

func haikuLoopReadme() string {
	return "---\n" +
		"commissioned-by: spacedock@1\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"state: .spacedock-state\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: implementation\n" +
		"      initial: true\n" +
		"    - name: validation\n" +
		"      gate: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Haiku Loop Spike Fixture\n\n" +
		"### implementation\n\nDo the trivial work: append the exact note line `" + haikuBuiltMarker + "` to the entity body, then append a `## Stage Report: implementation` section with one `- DONE:` item.\n\n- **Outputs:** the note recorded and a stage report.\n\n" +
		"### validation\n\nGate: a level-3 reviewer adjudicates the verdict against the entity's acceptance criteria.\n\n- **Outputs:** a PASSED/REJECTED verdict.\n\n" +
		"### done\n\nTerminal state.\n"
}

func haikuLoopEntity() string {
	return "---\n" +
		"id: " + haikuLoopSlug + "\n" +
		"title: Build the Widget\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Build the Widget\n\n" +
		"Trivial fixture entity at the dispatchable implementation stage.\n\n" +
		"## Acceptance criteria\n\n" +
		"- AC-1: the implementation worker's required note line (named in the implementation stage definition) appears in the entity body after implementation.\n"
}

// assertHaikuFixtureFingerprint is the running-research-spikes precondition: it
// asserts the fixture is in the EXPECTED start state — split-root backend, the entity
// at the dispatchable `implementation` stage, the planted AC present, the built
// marker NOT yet present — BEFORE the expensive drive. A green drive against a
// wrong-state fixture proves nothing, so this fails LOUD if the fixture drifted. It
// shells the real binary's `status --boot --json` so the precondition reads the same
// resolution the FO will.
func assertHaikuFixtureFingerprint(t *testing.T, binary string, fx haikuLoopFixture) {
	t.Helper()
	out, err := exec.Command(binary, "status", "--boot", "--json", "--workflow-dir", fx.definitionDir).CombinedOutput()
	if err != nil {
		t.Fatalf("fixture fingerprint: status --boot failed: %v\n%s", err, out)
	}
	var boot struct {
		StateBackend     string `json:"state_backend"`
		EntityDir        string `json:"entity_dir"`
		EntityDirPresent string `json:"entity_dir_present"`
		Dispatchable     []struct {
			Slug    string `json:"slug"`
			Current string `json:"current"`
			Next    string `json:"next"`
		} `json:"dispatchable"`
	}
	if err := json.Unmarshal(out, &boot); err != nil {
		t.Fatalf("fixture fingerprint: boot JSON unparseable: %v\n%s", err, out)
	}
	if boot.StateBackend != "split-root" {
		t.Fatalf("fixture fingerprint: state_backend=%q, want split-root\n%s", boot.StateBackend, out)
	}
	if boot.EntityDir != fx.stateDir {
		t.Fatalf("fixture fingerprint: entity_dir=%q, want the state checkout %q", boot.EntityDir, fx.stateDir)
	}
	if boot.EntityDirPresent != "true" {
		t.Fatalf("fixture fingerprint: entity_dir_present=%q, want true", boot.EntityDirPresent)
	}
	var found bool
	for _, d := range boot.Dispatchable {
		if d.Slug == fx.entitySlug {
			found = true
			if d.Current != "implementation" || d.Next != "validation" {
				t.Fatalf("fixture fingerprint: %s at %s->%s, want implementation->validation", d.Slug, d.Current, d.Next)
			}
		}
	}
	if !found {
		t.Fatalf("fixture fingerprint: entity %q not dispatchable; the drive would have nothing to drive\n%s", fx.entitySlug, out)
	}
	// The built marker must NOT be present yet — else the impl worker's effect is
	// indistinguishable from a pre-seeded body and failure-mode 3 (bare-dispatch) is
	// unobservable.
	if strings.Contains(readFile(t, fx.entityPath), haikuBuiltMarker) {
		t.Fatalf("fixture fingerprint: built marker %q already present at start; the impl effect would be unobservable", haikuBuiltMarker)
	}
}

// haikuLoopPrompt is the hand-authored SIMPLIFIED loop — the entity body's
// "The simplified loop (the hand-followed prose-functions)" section, NOT the real FO
// contract. The whole spike question is whether a Haiku FO follows THIS, so the full
// contract is deliberately absent (bare `claude`, no `--agent`). Each prose-function
// names the binary call its body wraps and the future verb it becomes, so the
// recorded must-build list maps one-to-one onto loop steps.
//
// The «gate» body carries the captain-ratified residency: route the verdict to a
// PER-JUDGMENT bare Agent(model="opus") BLOCKING call and pass it the ACTUAL stage
// report + acceptance criteria (the shell probe proved an opus L3 reviewer refuses to
// adjudicate without the real content). Haiku adjudicates none alone.
func haikuLoopPrompt(fx haikuLoopFixture) string {
	return strings.Join([]string{
		"You are a first officer driving ONE workflow entity to its terminal stage by FOLLOWING the simplified loop below EXACTLY. Do this now, in one continuous run, with no clarifying questions and no waiting for more input. You do NOT have a separate FO contract — the loop below is your COMPLETE specification. Follow each prose-function (named in «guillemets») in order; do not improvise your own flow.",
		"",
		"Workflow definition directory: " + fx.definitionDir,
		"This is a SPLIT-ROOT workflow: the entity body and your commits live in the state checkout it reports, NOT beside the README. Use the entity_dir from «boot» for all entity reads/writes/commits.",
		"This is a SINGLE-ENTITY run, so you are in BARE mode: do NOT create a team; dispatch the worker with a bare Agent(...) call and wait for it to return.",
		"",
		"THE SIMPLIFIED LOOP (follow by hand, in order):",
		"",
		"«boot»() -> run `spacedock status --boot --json --workflow-dir " + fx.definitionDir + "`. Read entity_dir and state_backend from it. (becomes: spacedock team_action --boot)",
		"",
		"«next»() -> run `spacedock status --next --workflow-dir " + fx.definitionDir + " --json`. It yields the ONE dispatchable entity and its current+next stage. (becomes: spacedock dispatch next-action)",
		"",
		"«dispatch»(slug, stage) -> (1) build the worker dispatch artifact by piping this JSON to `spacedock dispatch build --workflow-dir " + fx.definitionDir + "` on stdin: {\"schema_version\":2,\"entity_path\":\"<entity_dir>/<slug>.md\",\"workflow_dir\":\"" + fx.definitionDir + "\",\"stage\":\"<stage>\",\"checklist\":[\"TODO: append the exact note line " + haikuBuiltMarker + " to the entity body\",\"TODO: append a ## Stage Report: <stage> section\"],\"bare_mode\":true}. (2) Spawn the worker with a bare Agent(...) call whose prompt is EXACTLY the `prompt` field from the dispatch build output — do NOT hand-assemble your own worker prompt. (3) Wait for the Agent to return; its return IS the worker's completion — do NOT re-dispatch the same stage. (becomes: spacedock dispatch build, already a verb)",
		"",
		"«state.commit»(slug) -> set the entity's frontmatter stage with `spacedock status --set <slug> status=<new-stage> --workflow-dir " + fx.definitionDir + " --json`, then make a PATH-SCOPED git commit naming ONLY the entity file in the state checkout: `git -C <entity_dir> add <entity_dir>/<slug>.md && git -C <entity_dir> commit -m \"...\" -- <entity_dir>/<slug>.md`. NEVER `git add -A`. If the commit hits a rebase/non-fast-forward conflict, HALT and report it — do NOT invent a recovery. (becomes: spacedock state set + state commit)",
		"",
		"«gate»(slug, stage) -> you may NOT decide the verdict yourself. Read the entity's `## Stage Report` section(s) and its `## Acceptance criteria` section, then ROUTE the judgment to a level-3 reviewer: call the Agent tool with subagent_type=\"general-purpose\", model=\"opus\", description=\"L3 gate verdict\", and a prompt that includes the FULL stage-report text and the FULL acceptance-criteria text and asks it to reply with exactly one line `VERDICT: APPROVED <reason>` or `VERDICT: REJECTED <reason>`. WAIT for the opus Agent to return. Record the verdict it returns as the gate decision — the verdict MUST originate from the opus Agent sub-call, not from you. (becomes: spacedock gate route + a level-3 verb)",
		"",
		"«merge»(slug) -> only if the gate verdict was APPROVED: terminalize the entity with `spacedock status --set <slug> status=done verdict=passed completed --workflow-dir " + fx.definitionDir + " --json`, then archive it with `spacedock status --archive <slug> --workflow-dir " + fx.definitionDir + " --json`, then make a final path-scoped commit of the entity in the state checkout. (becomes: spacedock merge guard)",
		"",
		"DRIVE ORDER: «boot» -> «next» -> «dispatch»(slug, implementation) -> «state.commit»(slug) advancing to validation -> «gate»(slug, validation) -> «merge»(slug). Then stop. Your final message must state the entity's terminal status and the gate verdict you recorded from the level-3 reviewer.",
		"",
		antiShutdownOverride,
	}, "\n")
}

// haikuLoopGrade is the durable-state-plus-stream grade of one drive: the on-disk
// end-state of the throwaway fixture AND the parsed tool-call stream facts AC-3/AC-4
// read off. Every field is a durable on-disk fact or a stream-derived fact — none is
// the FO's narration.
type haikuLoopGrade struct {
	// Durable end-state (read from disk after the drive).
	entityLocated    bool
	entityArchived   bool
	statusDone       bool
	verdictSet       bool
	builtMarker      bool
	pathScopedCommit bool

	// Stream-derived facts (parsed from the captured tool-call stream).
	dispatchBuildCalled bool // a `spacedock dispatch build` Bash call appeared (failure-mode 3)
	opusAgentSpawned    bool // an Agent(model=opus) tool_use appeared (the L3 route)
	gateVerdictFromL3   bool // the opus Agent returned a VERDICT line (gate provenance, AC-3)
}

// driveHaikuLoopOnce runs ONE bare-`claude --model haiku -p` hand-loop drive against
// a freshly-stood-up throwaway split-root fixture and returns the durable+stream
// grade plus the artifact paths. It is the unit the N>=3 validation lane wraps in a
// count loop. The launch is the AC-1 NEW path: bare `claude`, no front door, the
// contract NOT loaded; stdin is /dev/null so the harness does not eat the 3s
// no-stdin warning stall the shell probe surfaced.
func driveHaikuLoopOnce(t *testing.T, binary string, env []string, artifactDir string) (haikuLoopGrade, string) {
	t.Helper()

	root := t.TempDir()
	fx := writeHaikuLoopFixture(t, root)

	// running-research-spikes precondition: assert the start state BEFORE spending
	// the drive. A green run against a drifted fixture proves nothing.
	assertHaikuFixtureFingerprint(t, binary, fx)

	prompt := haikuLoopPrompt(fx)

	// Bare `claude` — NOT `spacedock claude`. No --plugin-dir, no --agent, no
	// --skip-contract-check: the FO contract is NOT loaded, which is the whole point
	// of the spike (does Haiku follow the SIMPLIFIED loop). withBinaryOnPath already
	// put the built `spacedock` on the child PATH (the prose-functions shell it).
	cmd := exec.Command("claude",
		"--model", "haiku",
		"-p", prompt,
		"--output-format", "stream-json",
		"--verbose",
		"--dangerously-skip-permissions",
	)
	cmd.Dir = fx.definitionDir
	cmd.Env = env
	// stdin from /dev/null: the shell probe showed a bare `claude -p` with an open
	// stdin warns "no stdin data received in 3s, proceeding without it" and stalls;
	// /dev/null makes it proceed immediately.
	devnull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer devnull.Close()
	cmd.Stdin = devnull

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	if startErr := cmd.Start(); startErr != nil {
		t.Fatalf("bare `claude --model haiku -p` failed to start: %v", startErr)
	}
	poller := newCmdPoller(cmd, pw)
	t.Cleanup(poller.kill)
	watcher := newStreamWatcher(newPipeLineSource(pr), poller, discardStreamLine)

	// drainToExit runs to exit accumulating the full transcript, or kills on a
	// 3-minute no-progress stall (a multi-turn Haiku drive with an opus sub-call has
	// legitimately quiet stretches; the budget resets on any drained line).
	stream, stallErr := watcher.drainToExit(quietBudgetDispatchClose, "haiku loop drive")

	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	streamPath := filepath.Join(artifactDir, "haiku-loop-stream.jsonl")
	if writeErr := os.WriteFile(streamPath, []byte(stream), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	t.Logf("haiku loop stream artifact: %s (%d bytes)", streamPath, len(stream))

	if stallErr != nil {
		t.Fatalf("%v\nArtifacts: %s", stallErr, artifactDir)
	}

	// A non-error result event is required: a 401/is_error result is a LOUD launch
	// failure (stale credential), never graded as behavior.
	if _, extractErr := extractClaudeFinalMessage(stream); extractErr != nil {
		t.Fatalf("bare claude launch failed: %v\nStream tail:\n%s\nArtifacts: %s", extractErr, tail(stream, 4000), artifactDir)
	}

	grade := gradeHaikuLoop(t, fx, stream)
	return grade, artifactDir
}

// gradeHaikuLoop reads the durable on-disk end-state of the fixture and parses the
// captured stream into the AC-3/AC-4 facts. It never reads the FO's narration: the
// durable facts come from the entity file + git log, the stream facts from the parsed
// tool-call entries.
func gradeHaikuLoop(t *testing.T, fx haikuLoopFixture, stream string) haikuLoopGrade {
	t.Helper()
	var g haikuLoopGrade

	body, where, found := locateEntity(fx.stateDir, fx.entitySlug)
	g.entityLocated = found
	if found {
		g.entityArchived = strings.Contains(where, "_archive")
		g.statusDone = strings.Contains(body, "\nstatus: done\n") || strings.Contains(body, "\nstatus: done\r\n")
		g.verdictSet = strings.Contains(body, "verdict: passed")
		g.builtMarker = strings.Contains(body, haikuBuiltMarker)
	}
	g.pathScopedCommit = someCommitNamesOnly(t, fx.stateDir, fx.entitySlug)

	parseHaikuStreamFacts(stream, &g)
	return g
}

// parseHaikuStreamFacts walks the captured stream-json and sets the stream-derived
// grade fields: a `spacedock dispatch build` Bash call (failure-mode 3: did the
// worker prompt come from dispatch build, not hand-assembly), an Agent(model=opus)
// tool_use (the L3 route — gate did not self-decide), and an opus Agent tool_result
// carrying a VERDICT line (AC-3 gate-verdict provenance: the verdict ORIGINATED from
// the opus sub-call).
func parseHaikuStreamFacts(stream string, g *haikuLoopGrade) {
	// opusAgentToolUseIDs collects the tool_use ids of Agent(model=opus) calls so the
	// matching tool_result (the verdict text) can be attributed to the opus sub-call.
	opusAgentToolUseIDs := map[string]bool{}

	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var e streamEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}

		// Assistant tool_use blocks: dispatch-build Bash + opus Agent spawn.
		for _, b := range e.toolUseBlocks() {
			if b.Name == "Bash" && strings.Contains(b.Input.Command, "dispatch build") {
				g.dispatchBuildCalled = true
			}
			if b.Name == "Agent" && haikuIsOpusModel(line) {
				g.opusAgentSpawned = true
				if b.ID != "" {
					opusAgentToolUseIDs[b.ID] = true
				}
			}
		}

		// User tool_result blocks: an opus Agent's returned verdict line is the
		// gate-verdict provenance. Match the opus Agent's tool_use_id and look for a
		// VERDICT: line in its result text.
		for _, rb := range e.resultBlocks() {
			if opusAgentToolUseIDs[rb.ToolUseID] && strings.Contains(rb.flatText(), "VERDICT:") {
				g.gateVerdictFromL3 = true
			}
		}
	}
}

// haikuIsOpusModel reports whether a raw assistant stream line carries an Agent
// tool_use whose input.model is opus. streamToolInput does not model the Agent
// `model` field, so this re-parses the raw line for the nested input.model — kept
// local to the spike rather than widening the shared streamToolInput shape.
func haikuIsOpusModel(line string) bool {
	var probe struct {
		Message struct {
			Content []struct {
				Type  string `json:"type"`
				Name  string `json:"name"`
				Input struct {
					Model string `json:"model"`
				} `json:"input"`
			} `json:"content"`
		} `json:"message"`
	}
	if json.Unmarshal([]byte(line), &probe) != nil {
		return false
	}
	for _, b := range probe.Message.Content {
		if b.Type == "tool_use" && b.Name == "Agent" && strings.Contains(strings.ToLower(b.Input.Model), "opus") {
			return true
		}
	}
	return false
}

// TestLiveHaikuLoopSpike proves the harness with ONE clean full
// dispatch->gate->merge->terminal drive (the assignment's implementation deliverable).
// It asserts the durable end-state AND the two NEW-mechanism stream facts the banked
// AC-1/AC-2 did not cover: the gate->L3 route from inside a Haiku `-p` session and the
// merge ceremony. N>=3 (AC-5) and the 5-failure-mode triage (AC-4) are operated by the
// FO at validation, wrapping driveHaikuLoopOnce in a count loop.
func TestLiveHaikuLoopSpike(t *testing.T) {
	binary := spacedockBinary(t)
	env := isolatedClaudeEnv(t, os.Getenv("HOME")) // t.Skip when no credential
	env = withBinaryOnPath(env, binary)

	artifactDir := claudeLiveArtifactDir(t, "haiku-loop-spike")
	started := time.Now()
	grade, artifacts := driveHaikuLoopOnce(t, binary, env, artifactDir)
	t.Logf("haiku loop drive finished in %s; artifacts: %s", time.Since(started).Round(time.Second), artifacts)

	// Durable end-state (AC-3): entity terminal + archived + verdict, and the impl
	// worker actually did work (the built marker landed), with a path-scoped commit.
	if !grade.entityLocated {
		t.Fatalf("entity %q not found in place or under _archive/ after the drive; artifacts: %s", haikuLoopSlug, artifacts)
	}
	if !grade.builtMarker {
		t.Errorf("built marker %q absent from the entity body — the dispatched impl worker did not do the work", haikuBuiltMarker)
	}
	if !grade.statusDone {
		t.Errorf("entity did not reach terminal `status: done`")
	}
	if !grade.entityArchived {
		t.Errorf("entity was not archived to _archive/ by the «merge» ceremony")
	}
	if !grade.verdictSet {
		t.Errorf("entity carries no `verdict: passed` from terminalization")
	}
	if !grade.pathScopedCommit {
		t.Errorf("no path-scoped commit naming only the entity in the state checkout history")
	}

	// Stream provenance (AC-3 + failure-mode 3/4): the worker dispatch came from
	// `dispatch build` (not hand-assembled), and the gate verdict ROUTED to an
	// Agent(model=opus) and ORIGINATED there (Haiku did not self-approve).
	if !grade.dispatchBuildCalled {
		t.Errorf("no `spacedock dispatch build` call in the stream — the FO bare-dispatched (hand-assembled the worker prompt)")
	}
	if !grade.opusAgentSpawned {
		t.Errorf("no Agent(model=opus) call in the stream — the «gate» step did not route to the level-3 reviewer")
	}
	if !grade.gateVerdictFromL3 {
		t.Errorf("the gate verdict did not trace to the opus Agent sub-call's returned VERDICT line — provenance unproven (AC-3)")
	}
}
