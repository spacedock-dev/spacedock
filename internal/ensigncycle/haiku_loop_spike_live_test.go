//go:build live

// ABOUTME: Live drive of a bare `claude --model haiku -p` FO following a hand-authored
// ABOUTME: MECHANICAL-ONLY loop (no gate, no L3, no FO contract) over a throwaway split-root fixture; grades durable state + stream.
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

// The 0221 premise-spike harness. It launches a BARE `claude --model haiku -p`
// (NOT `spacedock claude` — no `--agent`, the FO contract is NOT loaded) with a
// STRICTLY MECHANICAL hand-loop as the `-p` prompt, drives one throwaway split-root
// fixture entity through dispatch->advance->advance->terminalize+archive, and grades
// the durable on-disk end-state PLUS the captured tool-call stream — never the FO's
// narration.
//
// This is the gate-stripped isolation of the w4 haiku-loop spike. The w4 loop was
// `«boot»→«next»→«dispatch»→«state.commit»→«gate»(route-to-opus)→«merge»`; this loop
// is `«boot»→«next»→«dispatch»→«advance»(->integration)→«advance»(->done)+archive`.
// The `«gate»` step, its opus `Agent` sub-call, and its verdict-PROVENANCE grade are
// REMOVED. "Review the report → advance" is a MECHANICAL FO step (read the worker's
// appended report, then `status --set` to the next stage), not a verdict. Haiku is
// never asked to make a judgment call, so there is nothing to escalate, no tier to
// self-identify, no level-3 vocabulary anywhere in the fixture or the loop prompt.
//
// It reuses the proven internal/ensigncycle substrate VERBATIM: isolatedClaudeEnv
// (clean HOME + OAuth/API-key auth, t.Skip when neither), withBinaryOnPath (the
// hand-loop's prose-functions call `spacedock`, so the built binary stays on the
// child PATH even though the launched executable is bare `claude`), the
// streamWatcher quiet-budget liveness (drainToExit), extractClaudeFinalMessage, and
// cmdPoller. The two NEW mechanisms this spike relies on are both PASSED in w4: the
// bare-`claude -p` launch (w4 AC-1) and the loop-substitution hold (w4 AC-2). This
// spike REMOVES surface (the opus L3 route) rather than adding it, so it introduces
// nothing unverified.

// haikuLoopFixture is the throwaway split-root fixture the drive runs against: a
// definitionDir holding the README (which declares `state: .spacedock-state`) and a
// stateDir (the .spacedock-state checkout) holding the one entity. Split-root means
// the entity body + the path-scoped commit live in stateDir, NOT beside the README —
// the same contract the real workflow uses, so the hand-loop's «advance» body
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
// integration -> done(terminal), so the drive has a DISTINCT dispatch step
// (implementation), a DISTINCT mechanical advance (->integration), and a DISTINCT
// terminalize+archive step (->done). A one-stage implementation->done fixture would
// collapse advance and terminalize and lose the integration-transition observation.
//
//spacedock:live-fixture id=haiku-loop/experimental
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
		"    - name: integration\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Haiku Loop Spike Fixture\n\n" +
		"### implementation\n\nDo the trivial work: append the exact note line `" + haikuBuiltMarker + "` to the entity body, then append a `## Stage Report: implementation` section with one `- DONE:` item.\n\n- **Outputs:** the note recorded and a stage report.\n\n" +
		"### integration\n\nA plain non-terminal stage. The entity is advanced here mechanically after the worker returns.\n\n- **Outputs:** the entity advanced.\n\n" +
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
// at the dispatchable `implementation` stage with next `integration`, the planted AC
// present, the built marker NOT yet present — BEFORE the expensive drive. A green
// drive against a wrong-state fixture proves nothing, so this fails LOUD if the
// fixture drifted. It shells the real binary's `status --boot --json` so the
// precondition reads the same resolution the FO will.
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
			if d.Current != "implementation" || d.Next != "integration" {
				t.Fatalf("fixture fingerprint: %s at %s->%s, want implementation->integration", d.Slug, d.Current, d.Next)
			}
		}
	}
	if !found {
		t.Fatalf("fixture fingerprint: entity %q not dispatchable; the drive would have nothing to drive\n%s", fx.entitySlug, out)
	}
	// The built marker must NOT be present yet — else the impl worker's effect is
	// indistinguishable from a pre-seeded body and the bare-dispatch class is
	// unobservable.
	if strings.Contains(readFile(t, fx.entityPath), haikuBuiltMarker) {
		t.Fatalf("fixture fingerprint: built marker %q already present at start; the impl effect would be unobservable", haikuBuiltMarker)
	}
}

// haikuLoopPrompt is the hand-authored MECHANICAL-ONLY loop — NOT the real FO
// contract. The whole spike question is whether a Haiku FO follows THIS, so the full
// contract is deliberately absent (bare `claude`, no `--agent`). Each prose-function
// names the binary call its body wraps. The gate/L3/verdict-adjudication/tier
// vocabulary is ENTIRELY absent: there is no judgment step, only mechanical verbs in
// order. The terminal `verdict: passed` token is a FIXED, STRUCTURALLY-REQUIRED write
// (the binary's finalize guard forces it), never an adjudication.
func haikuLoopPrompt(fx haikuLoopFixture) string {
	return strings.Join([]string{
		"You are a first officer driving ONE workflow entity to its terminal stage by FOLLOWING the simplified loop below EXACTLY. Do this now, in one continuous run, with no clarifying questions and no waiting for more input. You do NOT have a separate FO contract — the loop below is your COMPLETE specification. Follow each prose-function (named in «guillemets») in order; do not improvise your own flow.",
		"",
		"Workflow definition directory: " + fx.definitionDir,
		"This is a SPLIT-ROOT workflow: the entity body and your commits live in the state checkout it reports, NOT beside the README. Use the entity_dir from «boot» for all entity reads/writes/commits.",
		"This is a SINGLE-ENTITY run, so you are in BARE mode: do NOT create a team; dispatch the worker with a bare Agent(...) call and wait for it to return. Do NOT spawn any other Agent — there is no reviewer, no judge, no second model. Every step below is a mechanical binary call YOU make.",
		"",
		"THE SIMPLIFIED LOOP (follow by hand, in order):",
		"",
		"«boot»() -> run `spacedock status --boot --json --workflow-dir " + fx.definitionDir + "`. Read entity_dir and state_backend from it.",
		"",
		"«next»() -> run `spacedock status --next --workflow-dir " + fx.definitionDir + " --json`. It yields the ONE dispatchable entity and its current+next stage.",
		"",
		"«dispatch»(slug, stage) -> (1) build the worker dispatch artifact by piping this JSON to `spacedock dispatch build --workflow-dir " + fx.definitionDir + "` on stdin: {\"schema_version\":2,\"entity_path\":\"<entity_dir>/<slug>.md\",\"workflow_dir\":\"" + fx.definitionDir + "\",\"stage\":\"<stage>\",\"checklist\":[\"TODO: append the exact note line " + haikuBuiltMarker + " to the entity body\",\"TODO: append a ## Stage Report: <stage> section\"],\"bare_mode\":true}. (2) Spawn the worker with a bare Agent(...) call whose prompt is EXACTLY the `prompt` field from the dispatch build output — do NOT hand-assemble your own worker prompt. (3) Wait for the Agent to return; its return IS the worker's completion — do NOT re-dispatch the same stage.",
		"",
		"«advance»(slug, new-stage) -> set the entity's frontmatter stage with `spacedock status --set <slug> status=<new-stage> --workflow-dir " + fx.definitionDir + " --json`, then make a PATH-SCOPED git commit naming ONLY the entity file in the state checkout: `git -C <entity_dir> add <entity_dir>/<slug>.md && git -C <entity_dir> commit -m \"...\" -- <entity_dir>/<slug>.md`. NEVER `git add -A`. If the commit hits a rebase/non-fast-forward conflict, HALT and report it — do NOT invent a recovery.",
		"",
		"«terminalize»(slug) -> terminalize the entity with `spacedock status --set <slug> status=done verdict=passed completed --workflow-dir " + fx.definitionDir + " --json` (the `verdict=passed` token is REQUIRED by the binary's finalize guard — write it exactly, it is not a judgment to make), then archive it with `spacedock status --archive <slug> --workflow-dir " + fx.definitionDir + " --json`, then make a final path-scoped commit of the entity in the state checkout (same path-scoped form as «advance»).",
		"",
		"DRIVE ORDER: «boot» -> «next» -> «dispatch»(slug, implementation) -> «advance»(slug, integration) -> «terminalize»(slug). Then stop. Your final message must state the entity's terminal status.",
		"",
		antiShutdownOverride,
	}, "\n")
}

// haikuLoopGrade is the durable-state-plus-stream grade of one drive: the on-disk
// end-state of the throwaway fixture AND the parsed tool-call stream facts. Every
// field is a durable on-disk fact or a stream-derived fact — none is the FO's
// narration.
type haikuLoopGrade struct {
	// Durable end-state (read from disk after the drive).
	entityLocated  bool
	entityArchived bool
	statusDone     bool
	completedSet   bool
	verdictSet     bool
	builtMarker    bool
	// pathScopedCommit: at least one commit naming ONLY the entity file.
	pathScopedCommit bool
	// integrationTransitionCommitted (the BLOCKER fix): the git log carries a commit
	// whose entity blob held `status: integration` AND there were >=2 path-scoped
	// commits — so a Haiku FO that jumped implementation->done (skipping the advance)
	// FAILS this field even though every other durable field passes.
	integrationTransitionCommitted bool

	// Stream-derived facts (parsed from the captured tool-call stream).
	dispatchBuildCalled  bool // a `spacedock dispatch build` Bash call appeared
	noStrongerModelAgent bool // NO Agent(model=opus|sonnet) tool_use appeared anywhere
}

// driveHaikuLoopOnce runs ONE bare-`claude --model haiku -p` mechanical-loop drive
// against a freshly-stood-up throwaway split-root fixture and returns the
// durable+stream grade plus the artifact paths. It is the unit the N>=3 lane wraps in
// a count loop. The launch is the bare-`claude` path: no front door, the contract NOT
// loaded; stdin is /dev/null so the harness does not eat the 3s no-stdin warning stall.
func driveHaikuLoopOnce(t *testing.T, binary string, env []string, artifactDir string) (haikuLoopGrade, string) {
	t.Helper()

	root := t.TempDir()
	fx := writeHaikuLoopFixture(t, root)

	// running-research-spikes precondition: assert the start state BEFORE spending
	// the drive. A green run against a drifted fixture proves nothing.
	assertHaikuFixtureFingerprint(t, binary, fx)

	prompt := haikuLoopPrompt(fx)

	// Bare `claude` — NOT `spacedock claude`. No --plugin-dir, no --agent, no
	// --skip-compat-check: the FO contract is NOT loaded, which is the whole point
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
	// 3-minute no-progress stall (a multi-turn Haiku drive has legitimately quiet
	// stretches; the budget resets on any drained line).
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
// captured stream into the grade facts. It never reads the FO's narration: the
// durable facts come from the entity file + git log, the stream facts from the parsed
// tool-call entries.
func gradeHaikuLoop(t *testing.T, fx haikuLoopFixture, stream string) haikuLoopGrade {
	t.Helper()
	var g haikuLoopGrade

	body, where, found := locateEntity(fx.stateDir, fx.entitySlug)
	g.entityLocated = found
	if found {
		g.entityArchived = strings.Contains(where, "_archive")
		g.statusDone = frontmatterField.MatchString(body)
		g.completedSet = completedSet.MatchString(body)
		g.verdictSet = strings.Contains(body, "verdict: passed")
		g.builtMarker = strings.Contains(body, haikuBuiltMarker)
	}
	g.pathScopedCommit = someCommitNamesOnly(t, fx.stateDir, fx.entitySlug)
	g.integrationTransitionCommitted = integrationTransitionCommitted(t, fx.stateDir, fx.entitySlug)

	parseHaikuStreamFacts(stream, &g)
	return g
}

// parseHaikuStreamFacts walks the captured stream-json and sets the stream-derived
// grade fields: a `spacedock dispatch build` Bash call (did the worker prompt come
// from dispatch build, not hand-assembly) and the ABSENCE of any Agent(model=opus|
// sonnet) tool_use (no stronger model touched the mechanical loop — the only agents
// that ran were the Haiku FO and its bare-mode worker). noStrongerModelAgent starts
// true and is cleared if a stronger-model Agent call is seen.
func parseHaikuStreamFacts(stream string, g *haikuLoopGrade) {
	g.noStrongerModelAgent = true

	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var e streamEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}

		for _, b := range e.toolUseBlocks() {
			if b.Name == "Bash" && strings.Contains(b.Input.Command, "dispatch build") {
				g.dispatchBuildCalled = true
			}
			if b.Name == "Agent" && haikuIsStrongerModel(line) {
				g.noStrongerModelAgent = false
			}
		}
	}
}

// haikuIsStrongerModel reports whether a raw assistant stream line carries an Agent
// tool_use whose input.model is a stronger model (opus or sonnet). streamToolInput
// does not model the Agent `model` field, so this re-parses the raw line for the
// nested input.model — kept local to the spike rather than widening the shared
// streamToolInput shape. The mechanical loop must run with NO stronger-model Agent,
// so any opus/sonnet Agent call clears noStrongerModelAgent.
func haikuIsStrongerModel(line string) bool {
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
		if b.Type == "tool_use" && b.Name == "Agent" {
			m := strings.ToLower(b.Input.Model)
			if strings.Contains(m, "opus") || strings.Contains(m, "sonnet") {
				return true
			}
		}
	}
	return false
}

// TestLiveHaikuLoopSpike proves the harness with ONE clean full
// dispatch->advance->advance->terminalize+archive drive (the assignment's
// implementation deliverable). It asserts the durable end-state INCLUDING the
// committed `integration` transition (so a shortcut implementation->done fails) and
// the stream facts (dispatch-build called, no stronger-model Agent). N>=3 is operated
// by TestLiveHaikuLoopSpikeN, wrapping driveHaikuLoopOnce in a count loop.
func TestLiveHaikuLoopSpike(t *testing.T) {
	binary := spacedockBinary(t)
	env := isolatedClaudeEnv(t, os.Getenv("HOME")) // t.Skip when no credential
	env = withBinaryOnPath(env, binary)

	artifactDir := claudeLiveArtifactDir(t, "haiku-loop-spike")
	started := time.Now()
	grade, artifacts := driveHaikuLoopOnce(t, binary, env, artifactDir)
	t.Logf("haiku loop drive finished in %s; artifacts: %s", time.Since(started).Round(time.Second), artifacts)

	// Durable end-state: entity terminal + archived + completed + verdict, the impl
	// worker actually did work (the built marker landed), a path-scoped commit, AND
	// the committed integration transition (the full loop, not a shortcut).
	if !grade.entityLocated {
		t.Fatalf("entity %q not found in place or under _archive/ after the drive; artifacts: %s", haikuLoopSlug, artifacts)
	}
	if !grade.builtMarker {
		t.Errorf("built marker %q absent from the entity body — the dispatched impl worker did not do the work", haikuBuiltMarker)
	}
	if !grade.statusDone {
		t.Errorf("entity did not reach terminal `status: done`")
	}
	if !grade.completedSet {
		t.Errorf("entity carries no `completed:` timestamp from terminalization")
	}
	if !grade.entityArchived {
		t.Errorf("entity was not archived to _archive/ by the «terminalize» ceremony")
	}
	if !grade.verdictSet {
		t.Errorf("entity carries no `verdict: passed` from terminalization")
	}
	if !grade.pathScopedCommit {
		t.Errorf("no path-scoped commit naming only the entity in the state checkout history")
	}
	if !grade.integrationTransitionCommitted {
		t.Errorf("no committed `integration` transition — the FO shortcut implementation->done (skipped the advance)")
	}

	// Stream facts: the worker dispatch came from `dispatch build` (not hand-assembled),
	// and NO stronger-model Agent ran (the mechanical loop is Haiku-only).
	if !grade.dispatchBuildCalled {
		t.Errorf("no `spacedock dispatch build` call in the stream — the FO bare-dispatched (hand-assembled the worker prompt)")
	}
	if !grade.noStrongerModelAgent {
		t.Errorf("an Agent(model=opus|sonnet) call appeared in the stream — a stronger model touched the mechanical loop")
	}
}
