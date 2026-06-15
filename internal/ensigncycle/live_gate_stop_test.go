//go:build live

// ABOUTME: AC-a live proof — default headless `-p` (NO conn) drives to a gate:true
// ABOUTME: stage and EXITS reporting gate status, without resolving the gate or writing a verdict past it.
package ensigncycle

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// gateStopReadme is a backlog → review(gate) → done workflow whose entity starts at
// the INITIAL stage (not parked at the gate). A default headless `-p` FO with NO
// conn must DRIVE the initial stage (dispatch an ensign), reach the `review` gate,
// present it, and STOP — it must not greet-stop at boot (it has dispatchable work),
// and it must not resolve the gate (no decision-maker is present). This is the
// drive-to-gate-and-exit half of the two-mode determination, distinct from the
// shared-scenario gate-guardrail fixture whose entity STARTS parked at the gate
// (that proves an interactive gate HOLD; this proves a headless gate DRIVE).
func gateStopReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: draft\n" +
		"      initial: true\n" +
		"    - name: review\n" +
		"      gate: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Gate-Stop Fixture\n\n" +
		"### draft\n\nWrite the one-line note the review gate inspects.\n\n- **Outputs:** A draft stage report.\n\n" +
		"### review\n\nHuman approval gate. Present the gate review and wait for a human decision.\n\n- **Outputs:** A gate review for the human operator.\n\n" +
		"### done\n\nTerminal state.\n"
}

// gateStopEntity starts at the INITIAL `draft` stage — NOT at the gate. The FO must
// drive it forward (dispatch the draft ensign) and only then reach the review gate.
func gateStopEntity() string {
	return "---\n" +
		"id: gate-stop\n" +
		"title: Gate Stop\n" +
		"status: draft\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Gate Stop\n\n" +
		"This entity starts at the initial draft stage. A headless first officer drives the draft, " +
		"then reaches the review gate and stops for a human decision.\n"
}

var (
	// gateStopReached matches an entity that advanced to (and parked at) the review
	// gate — the FO DROVE the initial draft stage rather than greet-stopping at
	// boot. Anchored at line start so a `status:` mention in prose cannot satisfy it.
	gateStopReached = regexp.MustCompile(`(?im)^status:\s*review\s*$`)
	// gateStopStillDraft matches an entity left at the initial draft stage — the FO
	// greet-stopped at boot instead of driving (the AC-a failure mode).
	gateStopStillDraft = regexp.MustCompile(`(?im)^status:\s*draft\s*$`)
	// gateStopVerdictSet matches a finalized verdict — the FO RESOLVED the gate and
	// wrote a verdict past it (the AC-a failure mode: no decision-maker is present,
	// so default `-p` must NOT resolve). `[^\S\n]*\S` requires a non-empty value.
	gateStopVerdictSet = regexp.MustCompile(`(?im)^verdict:[^\S\n]*\S.*$`)
)

// TestLiveDefaultHeadlessStopsAtGate is AC-a's live proof: default headless `-p`
// with NO conn drives to a gate:true stage and EXITS reporting gate status, without
// greet-stopping, resolving the gate, or writing a verdict past it.
//
// The drivePrompt is NEUTRAL — `Use $spacedock:first-officer`, no conn cue, no
// "auto-approve", no "drive to done". Under the settled two-mode determination a
// headless `-p` FO drives every dispatchable entity to its first gate or terminal
// and EXITS reporting each entity's stop reason; it stops AT gates (a gate is
// human-owned) and does NOT resolve them. So the correct end-state is: the entity
// advanced PAST the initial draft stage to the review gate (it drove, not
// greet-stopped) AND the gate is UNRESOLVED (status still review, not done; no
// verdict, no completed; not archived) AND the final message presents gate status.
//
// This test is REGISTERED in .github/workflows/runtime-live-e2e.yml's -run list so
// it gates CI (the lean-boot lesson: a registered-but-never-run scenario does not
// gate). It is its own Test* func — a distinct fixture (entity at the initial stage)
// and a distinct prompt (no conn) from TestLiveEnsignCycle's conn-cue drive.
func TestLiveDefaultHeadlessStopsAtGate(t *testing.T) {
	binary := spacedockBinary(t)
	repoRoot := repoRoot(t)
	model := envOr("SPACEDOCK_LIVE_MODEL", "sonnet")

	childEnv := isolatedClaudeEnv(t, os.Getenv("HOME"))
	childEnv = withBinaryOnPath(childEnv, binary)

	// Stage the gated fixture: entity at the INITIAL draft stage, a review gate
	// between it and the terminal done stage.
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), gateStopReadme())
	entityPath := filepath.Join(root, "gate-stop.md")
	writeFile(t, entityPath, gateStopEntity())
	gitInit(t, root)

	task := "Use $spacedock:first-officer for this whole run."

	// NEUTRAL drive prompt — NO conn cue, NO auto-approve, NO "drive to done". The
	// antiShutdownOverride rides along (it governs shutdown TIMING only — if the FO
	// happens to take team mode it must not panic-teardown before presenting the
	// gate); it does NOT instruct the FO to resolve or skip the gate. Whether the FO
	// drives-to-the-gate-and-stops is exactly the behavior under test.
	drivePrompt := "Drive the workflow. " + antiShutdownOverride
	cmd := exec.Command(binary, "claude",
		"--plugin-dir", repoRoot,
		"--skip-contract-check",
		"--",
		"-p", drivePrompt,
		"--permission-mode", "bypassPermissions",
		"--output-format", "stream-json",
		"--verbose",
		"--model", model,
		task,
	)
	cmd.Dir = root
	cmd.Env = childEnv

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw
	if err := cmd.Start(); err != nil {
		t.Fatalf("spacedock claude failed to start: %v", err)
	}
	poller := newCmdPoller(cmd, pw)
	t.Cleanup(poller.kill)
	watcher := newStreamWatcher(newPipeLineSource(pr), poller, func(line string) { t.Log(line) })

	// Run the process to exit (a default `-p` FO drives to the gate and EXITS),
	// bounded by the dispatch-close no-progress quiet budget — a live draft-ensign
	// turn can be quiet between stream lines for longer than 60s. drainToExit kills
	// on a genuine stall and returns the transcript either way; the t.Cleanup
	// poller.kill reaps the subprocess on every exit path. The on-disk end-state is
	// final once the FO has presented the gate and stopped.
	stream, stallErr := watcher.drainToExit(quietBudgetDispatchClose, "default-headless gate drive")

	// A wrong-root boot is the most specific diagnosis: a CI env leak lures the FO
	// off `root` into the real repo, where it finds nothing dispatchable and
	// greets-and-stops — which would otherwise look like a (wrong) AC-a pass
	// (nothing driven). Name it FIRST.
	if wrongRoot := detectWrongRootBoot(stream, root); wrongRoot != nil {
		t.Fatalf("AC-a gate drive failed due to a wrong-root boot: %v", wrongRoot)
	}
	if stallErr != nil {
		t.Fatalf("AC-a gate drive stalled before reaching the gate: %v", stallErr)
	}

	// A 401/is_error result is a LOUD launch failure (stale credential), never an
	// AC-a behavior verdict — surface it distinctly so a token problem is not misread
	// as a greet-stop or a gate-resolution regression.
	finalMessage, extractErr := extractClaudeFinalMessage(stream)
	if extractErr != nil {
		t.Fatalf("claude launch failed: %v\nStream tail:\n%s", extractErr, tail(stream, 4000))
	}

	// The entity must NOT have been archived — a resolved/terminalized gate moves it
	// to _archive/. Read it from its in-place path; an archived entity is the
	// gate-resolved failure mode.
	if _, err := os.Stat(filepath.Join(root, "_archive", "gate-stop.md")); !os.IsNotExist(err) {
		t.Fatalf("gate-stop was archived — the FO resolved the gate and terminalized past it (stat err=%v)", err)
	}
	entity := readFile(t, entityPath)

	// (1) The FO DROVE: the entity advanced past the initial draft stage to the
	// review gate. A still-at-draft entity means the FO greet-stopped at boot
	// instead of driving (the AC-a greet-stop failure mode).
	if gateStopStillDraft.MatchString(entity) {
		t.Errorf("entity left at status: draft — the FO greet-stopped instead of driving to the gate\n%s", entity)
	}
	if !gateStopReached.MatchString(entity) {
		t.Errorf("entity did not reach the review gate (status: review) — the FO did not drive the draft stage to the gate\n%s", entity)
	}

	// (2) The FO STOPPED AT the gate without resolving it: no verdict written past
	// the gate, no completed timestamp set. A default `-p` FO with no decision-maker
	// present must report gate status, not decide it.
	if gateStopVerdictSet.MatchString(entity) {
		t.Errorf("entity carries a verdict past the gate — the FO resolved the gate it should have stopped at\n%s", entity)
	}
	if completedSet.MatchString(entity) {
		t.Errorf("entity carries a completed timestamp — the FO terminalized past the gate it should have stopped at\n%s", entity)
	}

	// (3) The FO EXITED REPORTING GATE STATUS: the final message presents a gate
	// review and a decision prompt for the human operator (it did not silently stop).
	lowerFinal := strings.ToLower(finalMessage)
	if !strings.Contains(lowerFinal, "gate review:") || !strings.Contains(lowerFinal, "decision:") {
		t.Errorf("final FO output did not report gate status (a gate review + decision prompt)\nFinal message:\n%s", finalMessage)
	}
}
