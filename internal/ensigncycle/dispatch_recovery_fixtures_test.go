package ensigncycle

import "fmt"

// Fixtures and prompts for the two fo-dispatch-recovery live scenarios
// (degraded-bare, break-glass-shim). Default-tagged (no //go:build live) so the
// pure string builders are reusable without a model; the fixture WRITERS that touch
// disk live alongside their //go:build live runner (dispatch_recovery_live_test.go),
// matching the shallow-boot precedent (shared_fixtures_test.go vs
// shallow_boot_fixture_live_test.go).

func dispatchRecoveryReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: implementation\n" +
		"      initial: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Dispatch Recovery Fixture\n\n" +
		"### implementation\n\nAppend a one-line marker to `widget-task.md` proving the worker ran.\n\n- **Outputs:** An implementation stage report.\n\n" +
		"### done\n\nTerminal state.\n"
}

// dispatchRecoveryMarker is what the dispatched ensign is asked to append, so a
// successful (non-observable-relevant) run still leaves durable evidence a human
// can eyeball alongside the stream-json oracle.
const dispatchRecoveryMarker = "DISPATCH-RECOVERY-WORKER-RAN"

const dispatchRecoveryStageDefinition = "Append a one-line marker to `widget-task.md` proving the worker ran.\n\n- **Outputs:** An implementation stage report."

func dispatchRecoveryEntity() string {
	return "---\n" +
		"id: widget-task\n" +
		"title: Widget Task\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Widget Task\n\n" +
		"A single dispatchable entity at implementation. Dispatch one worker to append the marker `" + dispatchRecoveryMarker + "` to this file, then stop — do not advance to done.\n"
}

// bareReachablePrompt carries the `/spacedock bare` instruction in the INITIAL `-p`
// prompt. Post-retirement, `/spacedock bare` is a plain captain instruction to
// dispatch bare from that point (not a mode transition), so every Agent() dispatch
// in this run must be bare-shaped from the first dispatch onward, with NO retired
// Degraded Mode captain report and NO recovery-skill load.
//
// The trigger phrase must NOT be the first four characters of the `-p` argument:
// a live probe proved `claude -p "/spacedock bare\n..."` gets intercepted by the
// CLI's own slash-command dispatcher (final message: "Unknown command:
// /spacedock") before the FO contract ever sees the text. Prefixing a sentence
// that quotes the captain's command keeps it plain prompt content instead.
func bareReachablePrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"The captain has issued the command: /spacedock bare",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"Dispatch one worker to process the entity `widget-task` at its implementation stage per the workflow README, then stop once that worker completes. Do not advance the entity to done.",
		"Your final response must confirm the worker's implementation stage report landed.",
	)
}

// breakGlassShimPrompt carries NO special trigger phrase — the break-glass path is
// triggered by the PATH-shimmed `spacedock dispatch build` failing for real, not by
// prompt text. This is an ordinary dispatch-and-wait prompt; the shim is what
// forces the FO onto the break-glass path.
func breakGlassShimPrompt() string {
	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
		"Use $spacedock:first-officer for this whole run.",
		"Workflow directory: .",
		"Dispatch one worker to process the next dispatchable entity at its implementation stage per the workflow README, then stop once that worker completes. Do not advance the entity to done.",
		"Your final response must confirm the worker's implementation stage report landed.",
		"If the dispatch helper fails, follow your contract's break-glass recovery path rather than giving up.",
	)
}

// breakGlassShimTeamPrompt changes only the dispatch-mode selection. The helper
// still fails at the same seam and the worker receives the same entity/stage work.
func breakGlassShimTeamPrompt() string {
	return "You MUST run in team mode for this run: dispatch every worker as a named background teammate. " +
		breakGlassShimPrompt()
}
