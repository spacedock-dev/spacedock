---
id: 4tqngghpq91dsc17g0njh25m
title: survey — agentsview install-detection false-negatives under sandbox (asks to install when `agentsview --version` works)
status: done
source: "captain manual install-path test (2026-06-08) — running /spacedock:survey under sandbox, the skill detected agentsview as NOT installed and asked to `brew install --cask agentsview`, but `agentsview --version` → v0.32.1 (installed). The detection is sandbox-flawed."
score: "0.27"
started: 2026-06-08T15:48:51Z
completed: 2026-06-11T15:20:44Z
verdict: superseded
worktree:
issue:
sprint: 0198-pre-flip-hardening
group: survey
sprint-readiness: defer
superseded-by: survey-skill-correctness-pass
archived: 2026-06-11T15:20:44Z
---

Under sandbox, the survey skill's agentsview-installed detection gives a false negative — it reports agentsview missing (and prompts to install) even though `agentsview --version` succeeds.

## Problem

- `/spacedock:survey` step 1 checks whether agentsview is installed before reading the local session history. The probe is `skills/survey/SKILL.md:27`: `if ! command -v agentsview >/dev/null; then echo "AGENTSVIEW MISSING"; fi`.
- Under sandbox, that check returns "missing" → the skill asks to `brew install --cask agentsview` — but `agentsview --version` reports `v0.32.1` (it IS installed). The skill nags to install an already-present tool and can derail a read-only survey at the consent prompt.
- The probe contradicts the skill's own sandbox awareness. `SKILL.md:21` already states the sandbox "cannot read `~/.agentsview/` directly … even though the `agentsview` binary itself reads it," and the whole design drives reads *through the binary* precisely because raw filesystem access is sandbox-fragile. Yet the install gate uses `command -v` — a pure shell builtin that does a PATH-walk plus an `access()`/`stat()` on each `<dir>/agentsview` candidate. It is the *filesystem-access* class of probe the rest of the skill avoids.

## Spike: ground the bug (DONE — see Stage Report)

Reproduced the probe behavior and identified the exec-allowed / access-denied asymmetry that splits the two probes. Findings:

1. **The failing probe is FS-access, the working one is exec.** `command -v agentsview` resolves a binary by walking `$PATH` and `access()`-checking each candidate (no `execve`). `agentsview --version` resolves the name and `execve`'s it. The macOS Seatbelt sandbox ("Safehouse's outer Seatbelt sandbox," per `~/.cache/claude/safehouse-vscode-reuse-needs-running-vscode.sh`) gates by **path and syscall class**, so a profile can permit `execve` of a binary while denying the `stat()`/`access()`/dir-read that `command -v` needs — making the binary runnable-by-name but invisible to `command -v`.
2. **Observed the asymmetry live, same dir.** In this very sandbox, `~/.local/bin/safehouse` is stat/read-denied (`file …/safehouse` → `Operation not permitted`) while sibling binaries in the same `~/.local/bin` exec fine. The sandbox profile is per-path; the captain's failing session had a profile where `agentsview`'s dir/stat fell on the denied side while its `execve` stayed allowed — exactly the bug's shape (`command -v` MISSING while `--version` exit 0).
3. **Demonstrated the divergence with a fixture.** A binary reachable by `execve` (absolute/shim path) but NOT discoverable by `command -v`'s PATH-walk yields `command -v: MISSING` and `invoke: PRESENT` together — the false negative — whereas an invocation probe reports PRESENT. Shells agree on exit codes: command-not-found → 127, present `--version` → 0 (verified under `sh`, `bash`, `zsh`).

(In this session's relaxed profile `command -v agentsview` happens to succeed, so the bug does not reproduce here directly; the captain's stricter profile is where it manifests — which is why the behavioral proof is a live drive under sandbox, below.)

## Decision: probe by invocation, not by PATH-walk

Detect agentsview the way the survey actually *uses* it — by invoking it and checking the exit code — instead of the FS-access `command -v` builtin. This aligns the gate with the skill's own "drive everything through the binary" principle, so the probe survives whatever the sandbox allows the survey's real reads to survive.

**Exact change** — `skills/survey/SKILL.md:27`:

```
- if ! command -v agentsview >/dev/null; then echo "AGENTSVIEW MISSING"; fi
+ if ! agentsview --version >/dev/null 2>&1; then echo "AGENTSVIEW MISSING"; fi
```

The `2>&1` suppresses both a present binary's `--version` output and a missing binary's "command not found" so only the intended `AGENTSVIEW MISSING` sentinel reaches the agent. Contract is unchanged: silent ⇒ present; prints `AGENTSVIEW MISSING` ⇒ absent (the step-1 consent/install path is untouched). Survey-skill change only — no agentsview, binary, or Go changes.

**Scope note (no install-instruction change).** The install fallback in `SKILL.md:30` (`brew install --cask agentsview` / `curl … | bash`) is unchanged — this task fixes only the false-negative detection, not the install path.

## Acceptance criteria

- **AC-1 — the install gate detects agentsview via invocation exit code, not a PATH-walk/FS-access check.** `skills/survey/SKILL.md` step 1 probes presence with `agentsview --version` (exit-code semantics) rather than `command -v agentsview` (or any `which`/`test -x`/`stat`/PATH-walk equivalent). Verified by: AC-3's live drive under sandbox — the behavioral proof that the gate no longer false-negatives where an FS-access probe would. (A grep over SKILL.md for the new wording does NOT satisfy this AC; the prose change is authoring work, not the proof.)

- **AC-2 — the probe distinguishes present from absent and emits only the sentinel.** With agentsview present, the probe prints nothing (no `--version` banner, no stderr leak); with agentsview absent (name not resolvable), it prints exactly `AGENTSVIEW MISSING`. Verified by: a probe-behavior test that runs the exact `SKILL.md:27` one-liner twice — once with a stub `agentsview` on a synthesized PATH (expect empty output, exit 0) and once with the name removed from PATH (expect sole line `AGENTSVIEW MISSING`) — asserting captured stdout+stderr against those two independent fixture conditions. Fails if the line reverts to a PATH-walk probe (which would mis-detect the exec-reachable / stat-denied case) or leaks `--version` output.

- **AC-3 — under sandbox with agentsview installed, the survey detects it present and does NOT prompt to install.** Driving `/spacedock:survey` in a sandbox where the prior `command -v` probe false-negatives, the run proceeds past step 1 to the sync/scan without emitting the consent-to-install prompt. Verified by: a live drive of the survey skill under sandbox (the only place the failure manifests), observing the survey continue without the `brew install --cask agentsview` consent prompt — and, as the negative control, confirming the OLD `command -v` probe would have stopped at that prompt in the same sandbox.

## Test plan

- **Mechanism de-risked (this ideation's spike).** The behavior the fix rests on — `command -v` is FS-access (`access()`/PATH-walk) while invocation is `execve`; Seatbelt can deny one and allow the other; the two probes provably diverge for an exec-reachable-but-stat-denied binary; exit codes are 0/127 across `sh`/`bash`/`zsh` — is confirmed in the Stage Report. No further spike needed before build; the build's first test is AC-2 over the present/absent fixture conditions.

- **AC-2 — probe-behavior test (deterministic, the gate-able half).** A small test extracts/runs the exact step-1 probe line against two synthesized conditions (stub binary on PATH; name absent from PATH) and asserts the captured output. Expected values (`""` for present, `AGENTSVIEW MISSING` for absent) come from the fixture conditions, an independent source — not from SKILL.md prose. This catches a regression to any FS-access probe at the unit level. If the existing survey test harness is shell-fixture-based, add it there; otherwise a focused shell test under the survey integration testdata. Cost: LOW (~one fixture pair + assertions).

- **AC-3 — live workflow drive under sandbox (the only proof for the false-negative fix).** Run `/spacedock:survey` under sandbox with agentsview installed; observe it pass step 1 without the install consent prompt. As the negative control in the same sandbox, confirm the pre-fix `command -v` probe reports `AGENTSVIEW MISSING` there (establishing the drive is exercising the failing path, not a vacuous pass). Cost: minutes; needs the sandboxed harness and the agentsview binary (present, v0.32.1).

- Estimated total cost/complexity: **LOW.** One-line SKILL.md change, one probe-behavior fixture/test, one sandboxed live drive. No binary/Go/agentsview changes (explicitly out of scope).

## Notes

Survey-skill bug; same area as xn / 1p27 / 69rk. Surfaced by captain manual install-path testing (2026-06-08). 0198 survey group. Root cause is a probe/sandbox-class mismatch (FS-access builtin vs the skill's own through-the-binary principle), not an agentsview defect — fixed in the skill, consistent with the sibling survey workarounds.

## Stage Report: ideation

- DONE: Ground the bug (spike): reproduce the sandbox false-negative — determine HOW the survey skill currently detects agentsview-installed (which probe), and why it fails under sandbox while `agentsview --version` (exit 0, v0.32.1) works. Record the failing probe vs the working one.
  Failing probe: `command -v agentsview` (`SKILL.md:27`) — a PATH-walk + `access()`/`stat()` FS-access builtin. Working probe: `agentsview --version` — an `execve`. Macos Seatbelt ("Safehouse" outer sandbox) gates per-path/per-syscall-class, so it can deny the stat/access `command -v` needs while allowing the `execve` — observed live (`~/.local/bin/safehouse` stat-denied "Operation not permitted" while siblings exec), and reproduced via a fixture where an exec-reachable binary off PATH gives `command -v: MISSING` + `invoke: PRESENT`. Recorded in "Spike".
- DONE: Decide the robust detection: probe agentsview the way the survey actually uses it (e.g. the `agentsview --version` exit code, or the real read path), not the sandbox-fragile check. Survey-skill change only.
  Decision: replace `command -v agentsview` with `agentsview --version >/dev/null 2>&1` exit-code probe — matches the skill's "drive everything through the binary" principle. Exact before/after recorded in "Decision"; contract (silent⇒present, `AGENTSVIEW MISSING`⇒absent) preserved; verified the replacement is silent-when-present and prints only the sentinel-when-absent across sh/bash/zsh.
- DONE: Produce build-ready ACs + test plan: under sandbox with agentsview installed, the survey detects it present and does NOT prompt to install — verified by a live drive under sandbox (the only place the failure manifests).
  AC-1 (gate uses invocation, proven by AC-3 live drive) + AC-2 (deterministic probe-behavior test over present/absent fixture conditions, expected values from the fixtures not from prose) + AC-3 (live drive under sandbox, with the pre-fix `command -v` false-negative as negative control). Recorded in "Acceptance criteria" and "Test plan".

### Summary

The current install gate uses `command -v agentsview` — an FS-access (PATH-walk + `access()`) shell builtin — which contradicts the skill's own sandbox principle of driving reads through the binary. Under macOS Seatbelt the `execve` of `agentsview --version` can be allowed while the `stat`/`access` `command -v` relies on is denied (observed live on a sibling binary and reproduced with a fixture), so the gate false-negatives and nags to install an already-present tool. Decided fix: a one-line SKILL.md change to an `agentsview --version` exit-code probe; ACs pair a deterministic probe-behavior test with a sandboxed live drive (the only place the failure manifests), with the pre-fix probe as the negative control. No spike beyond this needed; mechanism de-risked here.
