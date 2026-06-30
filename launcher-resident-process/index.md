---
id: vva363ty3afdgvxcacrfx6mc
title: Keep the spacedock launcher resident — spawn the host instead of exec'ing away
status: validation
source: captain (2026-06-10) — zellij restart legibility today, sidecar launches later
started: 2026-06-29T22:28:07Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-launcher-resident-process
issue:
mod-block: merge:pr-merge
pr: "#442"
---

`spacedock claude` currently replaces itself with the host process via `syscall.Exec` (`internal/cli/host_exec.go:275`), so the launcher vanishes from the process tree the moment the host starts. The captain wants the launcher to stay resident as the parent process for two reasons:

1. **Future sidecars** — a resident launcher can spawn companion processes (sidecar services, watchers) alongside the host later; an exec'd-away launcher cannot.
2. **Zellij legibility now** — with exec, zellij's restart view shows the raw underlying command (`Waiting to run: bash …/safehouse --trust-workdir-config -- claude --dangerously-skip-permissions --agent spacedock:first-officer --plugin-dir … --resume <session-id>`) instead of a legible `spacedock claude …` line, making session restarts hard to read.

The change: spawn the host as a child (wait, forward signals, propagate exit code, pass through the TTY) instead of `syscall.Exec`.

## Problem

`execHost.Launch` (`internal/cli/host_exec.go:275-281`) ends every `spacedock claude/codex/pi` launch with `syscall.Exec`, which replaces the launcher's process image with the host. The instant the host starts, the launcher is gone from the process tree — there is no spacedock process between the shell and the host. Two consequences matter:

1. **No resident parent for future companion processes.** A launcher that has exec'd away cannot spawn or supervise sidecar processes (watchers, services) alongside the host. Establishing a resident parent now is the prerequisite for that later capability.
2. **Illegible session managers today.** Because spacedock is no longer in the tree, zellij's restart view shows the raw inner command (`Waiting to run: bash …/safehouse --trust-workdir-config -- claude --dangerously-skip-permissions --agent spacedock:first-officer --plugin-dir … --resume <session-id>`) instead of the legible `spacedock claude …` the user typed. Session restarts are hard to read.

## Proposed approach

Replace `syscall.Exec` in `execHost.Launch` with a **resident-parent, shared-foreground-group spawn-wait** model (validated by the spike below):

1. `cmd := exec.Command(bin, argv[1:]...)`, inherit the real terminal: `cmd.Stdin/Stdout/Stderr = os.Stdin/Stdout/Stderr`, `cmd.Env = env`.
2. **Do not set `SysProcAttr.Setpgid`.** The host stays in the launcher's process group, which the shell already made the terminal's foreground group. The host therefore owns the controlling TTY (no SIGTTIN/SIGTTOU on read/write) and the kernel delivers terminal-generated signals to the host directly.
3. **Signal handling via `signal.Notify` (NOT `signal.Ignore`).** This is load-bearing, not stylistic — see spike finding S3. Minimal correct set:
   - **Absorb** `SIGINT`, `SIGQUIT` in the launcher (so it stays resident; the host already received its own copy from the foreground group — do **not** re-forward, which would double-deliver).
   - **Forward** `SIGTERM`, `SIGHUP` to the host (these reach only the launcher when sent to its pid — e.g. a `kill` or a zellij pane close — so relaying them tears the pair down together).
   - **Leave `SIGTSTP`/`SIGCONT` at default disposition** so Ctrl-Z suspends the whole foreground group and the shell's `fg`/`bg` resumes it (correct shared-group job control). Interactive raw-mode hosts (claude/codex/pi TUIs) never generate SIGTSTP anyway — under raw mode ISIG is off and Ctrl-Z arrives as a keystroke.
   - **Leave `SIGWINCH` at default (ignore).** The launcher survives it; the host gets its own copy from the foreground group. No launcher action needed.
4. `cmd.Wait()`, then propagate the host's exit status. Portable core: `(*os.ProcessState).ExitCode()`. On unix, refine signal-death to the shell `128+signum` convention via the `WaitStatus` (`go test`-verifiable).

**Interface change.** `hostOps.Launch(argv, env) error` becomes `Launch(argv, env) (int, error)`: the `int` is the host's propagated exit code; the `error` is reserved for *launch* failure (host binary not found, fork failure), not a non-zero host exit. The three call sites change from `if err := ops.Launch(...); err != nil { return 1 }; return 0` to `code, err := ops.Launch(...); if err != nil { return 1 }; return code` — `runClaude` (`frontdoor.go:377`), `runCodex` (`frontdoor.go:559`), `runPi` (`pi.go:320`). The test fake (`frontdoor_test.go:35`), `execPiRuntimeOps.Launch` (`pi.go:60`), and the `piRuntimeOps` interface (`pi.go:32`) update signatures. (Today `syscall.Exec` never returns on success, so the existing `return 0` was effectively unreachable on success; the new model makes it the host's real code.)

**Host scope: all three hosts.** claude/codex/pi all route through the single `execHost.Launch`, so the change applies uniformly with no per-host branch. The resident-parent motivations apply to whichever host the user launches.

**Cross-platform scope: unix (macOS/Linux); keep Windows compiling, no behavior claim.** The portable core (`exec.Command`/`Wait`/`ProcessState.ExitCode`) compiles everywhere; the unix-only signal model (`SIGWINCH`/`Getpgrp`/`WaitStatus.Signaled`) lives in a `//go:build unix` helper with a `//go:build !unix` no-op so the package still builds under `GOOS=windows`. **Correction to the seed premise:** the seed says "Windows already emulates spawn-wait under exec" — this is false. `syscall.Exec` on Windows returns `EWINDOWS` immediately (verified against Go 1.26 `src/syscall/exec_windows.go:453`), so `spacedock <host>` is *non-functional* on Windows today. This task does not change that contract: it scopes the resident-parent model to unix and makes no validated Windows claim. (The portable core happening to spawn-wait on Windows is an unclaimed, untested side effect; an explicit Windows lane is a separate task.)

## Out of scope

- **Actual sidecar/companion-process features.** This task only establishes the resident-parent process model that makes them possible later.
- **Changing Windows launch behavior.** Non-functional today (EWINDOWS); remains unsupported. Only compile-safety is in scope.
- **Job-control polish beyond transparent passthrough.** The launcher leaves SIGTSTP/SIGCONT at default (correct shared-group suspend/resume); it does not implement its own fg/bg/Ctrl-Z bookkeeping.

## Acceptance criteria

**AC-1 — While the host runs, the launcher remains resident as the host's parent (host PPID == launcher PID), where today's exec model has no resident launcher at all.**
This is the value-measuring AC against a baseline that moves the wrong way: the measured end-value is "is the host a child of a live launcher". New model: `host.ppid == launcher.pid` and `host.pid != launcher.pid` (TRUE). Baseline (`syscall.Exec`): the host *is* the launcher — `host.pid == launcher.pid`, `host.ppid ==` the shell — so the metric is FALSE.
Verified by: a Go behavior test that launches a stub-host fixture in `wait` mode, reads the PPID the stub reports (written to a temp file), and asserts it equals the launcher subprocess's PID; plus the contrasting assertion that an `exec`-model launcher yields `host.pid == launcher.pid`. (Spike: AC-1 pair, both observed.)

**AC-2 — The launcher's exit code equals the host's exit code.**
End-state: for host exits of 0, a nonzero code N, and signal-death, the launcher exits with the same status (signal-death → `128+signum`). Baseline-movable: a spawn-wait that ignores `Wait`'s error would exit 0 regardless (wrong way).
Verified by: a Go behavior test running the stub host with several exit codes and asserting the launcher's exit code matches each. (Spike: codes 0/7/42 propagated; SIGINT-death → 130; SIGTERM-death → 143.)

**AC-3 — Terminal signals reach the host and interactive TTY behavior is preserved, with the launcher surviving to propagate.**
End-state (all observable): (a) the host's stdin is a real TTY (`isatty` true) and it sees the real window size; (b) a terminal Ctrl-C delivers SIGINT to the host (it is interrupted/handles it), the launcher does **not** die first, and the host's resulting exit code propagates; (c) a terminal resize delivers SIGWINCH to the host; (d) an external `SIGTERM` to the launcher pid is forwarded to the host and tears the pair down.
Verified by: a PTY-driven Go behavior test (allocate a controlling terminal, e.g. `github.com/creack/pty` as a test-only dep) that injects `\x03`/resize/`SIGTERM` and asserts the stub-host log + launcher exit code. (Spike: all four observed, 11/11.)

## Test plan

- **AC-1, AC-2 — no PTY needed.** Pure-stdlib Go subprocess behavior tests: build a stub-host fixture (`go test` helper-process pattern or a small fixture binary), run the launch path as a subprocess, assert reported PPID and propagated exit code. Cheap (sub-second), no new deps.
- **AC-3 — PTY-backed.** Needs a controlling terminal to exercise terminal-generated signal delivery and isatty; add `github.com/creack/pty` as a **test-only** dependency (widely used, MIT). The throwaway spike harness (`drive.py`, pty.fork) maps directly to this Go test — it is the implementation's first test. Moderate complexity (one helper to allocate the pty + drive bytes/signals).
- **Build-tag check.** A CI `GOOS=windows go vet ./internal/cli` (or build) guards that the package still compiles with the unix signal file tagged out.
- **Manual zellij observation** of the legible `spacedock claude …` restart line is supporting evidence, not the proof.
- Estimated cost: small. AC-1/AC-2 are quick stdlib tests; AC-3 is one PTY helper. The risky mechanism is already spiked, so implementation is mostly translating the spike into Go tests + the interface change.

## Spike result (riskiest path — exercised first)

The design rests on `syscall.Exec → exec.Command+Wait` with signal forwarding and interactive-TTY passthrough. Spiked on **darwin/arm64, Go 1.26.1**, driving real Go launcher/stub binaries under a **real controlling terminal** (`pty.fork`), with evidence routed through a logfile (not PTY-stdout scraping). Throwaway artifacts in the session scratchpad (`spike/launcher`, `spike/stubhost`, `spike/execlauncher`, `spike/drive.py`).

- **S1 — Shared-foreground-group model works (11/11 checks).** AC-1 resident parent (host.ppid==launcher.pid) vs the exec baseline (host.pid==forked.pid, no resident parent); AC-2 exit-code propagation (0/7/42, plus 130 on SIGINT-death, 143 on SIGTERM-death); AC-3 isatty true + real winsize, Ctrl-C reaches the host while the launcher survives and propagates 130, SIGWINCH reaches the host **without** the launcher forwarding it (kernel delivers to the foreground group), external SIGTERM-to-launcher forwarded to the host.
- **S2 — Windows seed premise is false.** `syscall.Exec` returns `EWINDOWS` on Windows (Go `src/syscall/exec_windows.go:453`); the launch path is non-functional there today. Drove the unix/Windows scope decision above.
- **S3 — `signal.Notify`, not `signal.Ignore`, is load-bearing.** `signal.Ignore` sets kernel `SIG_IGN`, which is **inherited across execve**, so a host that relies on the default SIGINT disposition inherits an ignored SIGINT and **cannot be Ctrl-C'd**. Demonstrated faithfully: a "dumb" stub (no own SIGINT handler) survives Ctrl-C under a `signal.Ignore` launcher (trap), but is killed (exit 130) under a `signal.Notify` launcher. Hosts that install their own SIGINT handler (claude/codex/pi TUIs) re-arm SIGINT and are immune either way, but `signal.Notify` is the choice that is correct for both kinds of host.

Docs-confidence only (reasoned, not separately spiked): Ctrl-Z/job-control suspend-resume (standard POSIX shared-foreground-group behavior with SIGTSTP left at default; not generated by the raw-mode TUI hosts) and the optional re-raise-for-WIFSIGNALED fidelity refinement over the `128+signum` exit convention.

## Documentation change (doc diff)

The only published doc describing the launch process model is `docs/runtime-support.md` (section "Launcher binary propagation through wrappers"), which currently describes the *exec* model. Implementation applies:

**Before** (`docs/runtime-support.md`, the section's single paragraph):

> `spacedock claude` and `spacedock codex` attach `SPACEDOCK_BIN` to the process they exec, including the outer `safehouse -- ...` process when safehouse wrapping is active. Spacedock does not modify safehouse internals or assume a private passthrough mechanism; if a wrapper or runtime strips `SPACEDOCK_BIN` before the agent session observes it, the skill contract's `${SPACEDOCK_BIN:-spacedock}` convention degrades to the existing `$PATH` lookup.

**After** (one phrase corrected + one paragraph added):

> `spacedock claude` and `spacedock codex` attach `SPACEDOCK_BIN` to the host process they spawn, including the outer `safehouse -- ...` process when safehouse wrapping is active. Spacedock does not modify safehouse internals or assume a private passthrough mechanism; if a wrapper or runtime strips `SPACEDOCK_BIN` before the agent session observes it, the skill contract's `${SPACEDOCK_BIN:-spacedock}` convention degrades to the existing `$PATH` lookup.
>
> The launcher stays resident as the host's parent: it spawns the host as a child, inherits the terminal, forwards externally-targeted signals (`SIGTERM`/`SIGHUP`) while letting terminal signals (Ctrl-C, resize) reach the host through the shared foreground process group, and exits with the host's exit code — rather than replacing itself with the host. This keeps the `spacedock <host> …` command legible in process listings and session managers (for example zellij's restart view) and lets the launcher supervise companion processes alongside the session in future. (Unix launch lane; `spacedock <host>` is not a supported launch path on Windows.)

## Stage Report: ideation

- DONE: Fill Problem / Proposed approach / Out-of-scope; firm the three seed ACs into measurable end-states with test methods.
  Problem/approach/scope written above; AC-1 firmed as the value-measuring AC with an exec baseline that moves the wrong way; AC-2/AC-3 firmed with observable end-states + Go-test "Verified by".
- DONE: Spike the riskiest mechanism first (syscall.Exec → exec.Command+Wait + signal forwarding + interactive-TTY passthrough); record it.
  Real Go binaries under a real `pty.fork` controlling terminal, 11/11 checks (S1). Artifacts: scratchpad `spike/` (`launcher`, `stubhost`, `execlauncher`, `launcher_ignore`, `dumbstub`, `drive.py`).
- DONE: Decide cross-host and cross-platform scope explicitly.
  All three hosts change (single `execHost.Launch`). Unix-scoped; Windows kept compiling, no behavior claim. Forensic correction: seed's "Windows emulates spawn-wait under exec" is false — `syscall.Exec` → `EWINDOWS` (S2).
- DONE: Produce the concrete doc diff for the user-visible change.
  Before/after for `docs/runtime-support.md` "Launcher binary propagation through wrappers" recorded above (one phrase corrected + one resident-parent paragraph).

### Summary

Replaced the placeholder body with a firmed design: a resident-parent, shared-foreground-group spawn-wait replacement for `syscall.Exec`, with `Launch` returning `(int, error)` for exit-code propagation. The riskiest path was spiked first against a real controlling terminal (11/11), which also surfaced two findings the seed missed: the Windows premise is false (`syscall.Exec` → `EWINDOWS`, launch already non-functional there), and `signal.Notify` (not `signal.Ignore`) is load-bearing because `SIG_IGN` is inherited across exec and would silently break Ctrl-C for default-disposition hosts. Scope: all three hosts, unix-only, Windows kept compiling. Doc diff targets the one published doc that describes the launch process model.

## Stage Report: implementation

- DONE: Implement the resident-parent spawn-wait — replace syscall.Exec in execHost.Launch with exec.Command+Wait, NO Setpgid, signal.Notify (absorb SIGINT/SIGQUIT, forward SIGTERM/SIGHUP, leave SIGTSTP/SIGCONT/SIGWINCH at default), interface error→(int,error) at ALL implementors/call sites INCLUDING fakePiRuntimeOps.Launch (M1); unix //go:build helper + !unix no-op.
  Commit ca47dfe7: host_exec.go Launch rewritten; signal model + 128+signum exit mapping in host_launch_unix.go (//go:build unix) with host_launch_other.go (//go:build !unix) no-op; hostOps/piRuntimeOps + execHost/execPiRuntimeOps/fakeHost/fakePiRuntimeOps + runClaude/runCodex/runPi all updated.
- DONE: Write the AC tests — AC-1 (host.ppid==launcher.pid) via a stub host PLUS the exec-launcher baseline fixture (P5) contrasting host.pid==launcher.pid; AC-2 (0, nonzero N, signal-death→128+signum); AC-3 PTY-backed via github.com/creack/pty as a TEST-ONLY dep.
  host_launch_test.go (re-exec helper-process harness: launcher/execlauncher/stub roles) covers AC-1+AC-2; host_launch_pty_test.go covers AC-3. AC-2 signal-death kills the stub with no handler (genuine WIFSIGNALED → 143/130/137), exercising the launcher's signal branch. creack/pty absent from `go list -deps ./cmd/...`.
- DONE: Apply the doc diff to docs/runtime-support.md AND fix the stale artifacts (in-code "execs/replaces" comments, the syscall.Exec line cite, a by-analogy note); confirm GOOS=windows build/vet passes.
  runtime-support.md resident-parent paragraph + phrase fix; frontdoor.go:38 + host_exec.go Launch comment + ABOUTME updated; two internal/ensigncycle stream-watch comments de-stale'd; entity cite host_exec.go:227→:275; SIGHUP-forward/SIGQUIT-absorb by-analogy + harmless group-SIGHUP double-fire noted in host_launch_unix.go. GOOS=windows build + vet both exit 0.

### Summary

Translated the gate-approved design into Go: execHost.Launch is now a resident-parent spawn-wait (no Setpgid, shared foreground group), returning the host's exit code via a `(int, error)` signature threaded through both interfaces, all four implementors, and three call sites. The riskiest paths are proven by a re-exec helper-process harness — AC-1 resident parent vs the syscall.Exec baseline, AC-2 exit-code propagation (incl. real signal death), AC-3 PTY-backed terminal-signal delivery — all green under the full suite, with the unix signal model build-tagged so GOOS=windows still compiles. One race found and fixed during validation: the stub must register signal handlers before announcing STARTED, or an injected terminal signal can beat the listener.

## Stage Report: validation

- DONE: Reproduce each AC by RUNNING the tests (host_launch_test.go + host_launch_pty_test.go); confirm GOOS=windows build + vet.
  All green on the worktree, Go 1.26.1 darwin/arm64. AC-1 TestLaunchResidentParent 2/2 (host.ppid==launcher.pid + the exec baseline host.pid==launcher.pid contrast). AC-2 TestLaunchExitCodePropagation 6/6 (0/7/42 + genuine WIFSIGNALED death SIGTERM→143/SIGINT→130/SIGKILL→137). AC-3 TestLaunchPTYTerminalSignals 3/3 (isatty+24x80 winsize+Ctrl-C→130, resize→SIGWINCH, external SIGTERM→143). `GOOS=windows go build`+`vet ./internal/cli` exit 0; native `vet` exit 0; full internal/cli suite ok (63.9s).
- DONE: Verify the M1 completeness fix + creack/pty test-only.
  `(int,error) Launch` threaded through BOTH interfaces (hostOps frontdoor.go:40, piRuntimeOps pi.go:32), ALL FOUR implementors (execHost host_exec.go:281, fakeHost frontdoor_test.go:36, execPiRuntimeOps pi.go:60, fakePiRuntimeOps pi_frontdoor_test.go:42), and the three call sites all `code, err := ops.Launch(...); if err != nil { return 1 }; return code` (runClaude frontdoor.go:379, runCodex :562, runPi pi.go:320). `internal/cli` builds. creack/pty is TEST-ONLY: 0 hits in `go list -deps ./cmd/...`, present in `.TestImports` but absent from `.Imports` (only host_launch_pty_test.go imports it).
- DONE: Detached adversarial audit of the signal matrix on a THROWAWAY checkout (detached HEAD @ ca47dfe7, never the impl worktree), with mutation-confirmed teeth.
  Three probes attacking the absorb/forward matrix, all green on shipped code, each proven to catch its bug by injecting the defect: (A) SIGQUIT absorb — Ctrl-\ reaches the host through the foreground group while the launcher absorbs its copy, survives, propagates 131; mutation (drop SIGQUIT from Notify) → launcher dies to its own SIGQUIT (exit 2), CAUGHT. (B) SIGINT no-double-delivery — one Ctrl-C yields exactly 1 SIGINT at the host; mutation (forward SIGINT) → host counts 2, CAUGHT. (C) SIGHUP forward — pid-targeted SIGHUP relayed to the host → 129; mutation (absorb instead of forward) → host never gets it, CAUGHT. Handlers-before-STARTED race fix is sound: the launcher calls forwardHostSignals (signal.Notify) BEFORE cmd.Start (host_exec.go:294), so no terminal signal hits the launcher's default disposition; post-Wait forwards are safe via the `cmd.Process != nil` guard + Go's ErrProcessDone (no PID-recycle). The audit found NO code defect — clean.

### Summary

PASSED. All three ACs reproduce by running the real tests (AC-1 2/2, AC-2 6/6, AC-3 3/3), the Windows compile guard holds (build+vet exit 0), the full internal/cli suite is green, and the M1 `(int,error)` thread-through is complete across both interfaces, all four implementors, and three call sites with creack/pty confirmed test-only. The detached adversarial audit closed the design's two never-spiked "by analogy" claims (SIGQUIT-absorb, SIGHUP-forward) and the load-bearing no-double-delivery property with three mutation-confirmed probes — the signal model is correct under attack. Two non-blocking observations (do not gate): (1) a SIGTERM/SIGHUP landing in the nanosecond window between signal.Notify and cmd.Start is swallowed (cmd.Process==nil) — inherent to spawn-wait, no worse than the old exec model which would kill the launcher there, self-correcting on the next signal; (2) the SHIPPED suite's AC-3 Ctrl-C test cannot detect SIGINT double-delivery (it exits on the first signal) and the two by-analogy properties are unasserted — the code is correct (proven by the audit probes), so this is a cheap test-hardening opportunity, not a defect. Recommend the FO/captain optionally fold the three counting/by-analogy probes (a counting-stub mode + three test funcs) into the shipped suite. Recommendation: PASSED.

### Feedback Cycles

- Cycle 1 (CI-surfaced — offline lane on the Linux runner; validation had a platform gap, passing only on the macOS box). FAILURE: `TestLaunchPTYTerminalSignals/terminal_resize_delivers_SIGWINCH_to_the_host` timed out (5s) — "host never saw SIGWINCH after resize" — yet the captured Linux stub log contained the SIGWINCH line. ROOT CAUSE (not a timing race, despite the timeout symptom): the test matched the host log against `syscall.Signal.String()` output, whose SIGWINCH rendering diverges by platform — "window size changes" (darwin/BSD) vs "window changed" (Linux, confirmed against the Go 1.26.1 zerrors tables). The signal WAS delivered; the matcher was platform-specific, so the Linux match never fired. FIX (test-only, commit b15df9bd): the stub logs fixed per-signal markers (`SIGNAL SIGINT/SIGTERM/SIGHUP/SIGWINCH`) instead of `.String()`; the SIGWINCH case re-reads and logs the new window size so the assertion pins the resize-triggered signal (`SIGNAL SIGWINCH winsize=40x120`), not an incidental one; all three PTY subtests hardened consistently (SIGINT/SIGTERM matchers switched to the stable tokens too); `waitForLog` budget 5s→10s (under the stub's 30s self-exit) for loaded-CI headroom. The launcher signal model is untouched (mutation-audited). Verified on macOS (PTY ×10 green, vet/gofmt clean); Linux re-verified via PR #442 CI.
