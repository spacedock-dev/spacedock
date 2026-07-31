---
id: z3j0tsbr6t3mqd39rhs8bbvq
title: "First-officer boot install journey: Linux-aware hint, direct install/upgrade offer, sandbox detection"
status: ideation
source: "GitHub issue spacedock-dev/spacedock#581 (nomen429, 2026-07-30): the install journey the first-officer skill hits at the version gate in first-officer-shared-core.md Startup step 1."
started: 2026-07-31T02:32:43Z
completed:
verdict:
score:
worktree:
issue: spacedock-dev/spacedock#581
gates:
    version: 1
    current:
        gate: gate:z3j0tsbr6t3mqd39rhs8bbvq:ideation
    records:
        - id: gate:z3j0tsbr6t3mqd39rhs8bbvq:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:z3j0tsbr6t3mqd39rhs8bbvq-backlog-1
              briefing:
                id: briefing:z3j0tsbr6t3mqd39rhs8bbvq:backlog:attempt-1:revision-1
                digest: sha256:36af7b28170688e41da0817d982617ceb0a1f940079357dcc366a1b51645650c
                digest-domain: canonical-bytes
                request-digest: sha256:50e8a9607112885a2d5a6617c28b2e10290f6e6cc307de5a4876e0681bf09cc0
                room-ref: ./fo-boot-install-hint-linux-direct-sandbox/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:z3j0tsbr6t3mqd39rhs8bbvq:backlog:1
                briefing: briefing:z3j0tsbr6t3mqd39rhs8bbvq:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-31T02:31:16.935682Z"
                decision: approve
                reason: Captain approves the seed entering ideation to design the OS-aware hint, direct install/upgrade-and-resume offer, and sandbox detection with end-value-measuring ACs and the convergence spike.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:z3j0tsbr6t3mqd39rhs8bbvq:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:z3j0tsbr6t3mqd39rhs8bbvq-ideation-1
              briefing:
                id: briefing:z3j0tsbr6t3mqd39rhs8bbvq:ideation:attempt-1:revision-1
                digest: sha256:e1ce7f0a4d254da36db5ad188231cd0d4f0f602f2e19075bc0426c0512180744
                digest-domain: canonical-bytes
                request-digest: sha256:33d4ee7592c0c772b14fe8d5b516fba514d4a3023f747ca4a2ae7fcc4d3cec36
                room-ref: ./fo-boot-install-hint-linux-direct-sandbox/review/ideation/briefing-1
              provider-evidence:
                result-digest: sha256:4163c7a352ea3437e8ee9304c55abcca128a3598efef66fca9af526b349d96d2
                presented-inventory-digest: sha256:4659586febb0b4f9d72d1c0e853ff709874c5d5c04e29d8111bbc8b6d999c8e8
              resolution:
                type: Resolution
                id: resolution:binding-1785472404632857000
                briefing: briefing:z3j0tsbr6t3mqd39rhs8bbvq:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-07-31T04:33:24.63286Z"
                decision: revise
                includes:
                    - annotation:captain-1785472336725531000
                    - annotation:captain-1785472400412596000
              application:
                action: feedback
                target-stage: ideation
                state: pending
---

The first-officer boot hits the binary version gate (Startup step 1) and, on either abort class — binary absent, or binary present but wrong minor — prints a Mac-only Homebrew install hint and stops, leaving the human to copy-paste a command and restart the session. This task improves that install journey along the three axes the issue names: make the hint OS-aware (include the documented Linux `curl|sh` path, not just Homebrew), offer to run the install/upgrade directly and resume startup once the binary lands (turn hint-and-abort into one approved action), and detect sandboxed execution so a sandboxed install does not silently no-op (tell the human to run the install command themselves outside the sandbox, naming the exact command).

## Problem

Today the FO version gate (`first-officer-shared-core.md` Startup step 1) hits two abort classes, both Mac-only and both hint-and-abort:

1. **Binary absent** — the FO prints `brew install spacedock-dev/homebrew-tap/spacedock` and stops. A Linux VM with no `spacedock` on PATH sees a Homebrew command it cannot run (no `brew` on Linux). The documented Linux path (`curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh`, `docs/site/get-started/install.md` line 20) is never mentioned.
2. **Binary present but wrong version** — the FO aborts and points to `spacedock doctor`. The doctor's `tooOldBinaryRemedy` in `internal/contract/contract.go` is also Mac-only (`brew upgrade spacedock`). No self-repair path: the human copy-pastes a command and restarts the session.
3. **Sandboxed execution** — the FO doesn't check sandbox state at the version gate. A sandboxed session that ran the install inside the sandbox would install the binary somewhere the host cannot see — a silent no-op from the human's perspective. The detection mechanism (`safehouse.Inside()`) already exists and `--version` already prints `Sandbox: inside (agent-safehouse)`, but the version-gate prose ignores it.

## Proposed approach

Three coordinated changes across the binary and the FO skill prose:

### 1. OS-aware hint

- **Binary** (`internal/contract/contract.go`): `tooOldBinaryRemedy` gains an OS-conditional branch. On Linux (`runtime.GOOS == "linux"`), the remedy leads with the documented `curl|sh` path from `install.md` line 20; on macOS, it keeps the Homebrew lead as today. The source-build fallback stays on both. OS is threaded through `compareNamed` → `tooOldBinaryRemedy` (or read from `runtime.GOOS` inside the function — the simplest path, since the binary always knows its own OS).
- **Skill prose** (`first-officer-shared-core.md`): the binary-absent hint adds the Linux `curl|sh` path alongside Homebrew. The FO detects OS via `uname -s` (a prose mechanism — `--version` does not report OS; the model runs `uname -s` in bash and branches on `Darwin` vs `Linux`).

### 2. Direct install/upgrade-and-resume

- The version gate changes from hint-and-abort to offer-and-resume: the FO offers to run the install/upgrade command itself, the operator approves, the FO runs it, then re-checks `${SPACEDOCK_BIN:-spacedock} --version` and resumes startup if compatible.
- **Convergence mechanism** (from spike): after `curl|sh` install, `install.sh` prints `install.sh: installed spacedock <version> to <dir>/spacedock` to stderr and warns if the dir is not on PATH. The FO sets `SPACEDOCK_BIN` to the installed path (parsed from stderr, or probed at `$HOME/.local/bin/spacedock`) and re-checks `--version`. After `brew upgrade`, the binary stays at the brew prefix (on PATH), so re-check resolves automatically.
- **Retry bound**: one install attempt. If `--version` still fails after install, fall back to the hint-and-abort with the exact command for the human to run manually. The skill prose states this bound explicitly to prevent a loop.

### 3. Sandbox detection

- At the version gate, the FO checks sandbox state by parsing the `Sandbox:` line from `--version` output (already present via `safehouse.Inside()` → `printVersion` in `internal/cli/cli.go`).
- If `Sandbox: inside (...)`, the FO does NOT offer to run the install (it would no-op in the sandbox). Instead, it tells the human plainly: "You're running inside a sandbox. Run this command yourself, outside the sandbox: `<exact OS-aware install command>`." The command is the same OS-aware command from end-value 1.

## Out of scope

- Adding an OS line to `--version` output — the FO detects OS via `uname -s` in prose; adding it to the binary is a separate change.
- Changing the `safehouse.Inside()` detection mechanism itself — it already works; this task uses its output, not changes it.
- Windows support — `install.sh` already rejects non-darwin/linux; the FO gate is not expected to handle Windows.
- Changing the brew formula, `install.sh`, or `docs/site/get-started/install.md` — the canonical commands stay as-is; the hint cites them, not drifts from them.
- Auto-detecting whether `brew` is installed on macOS — the hint leads with Homebrew regardless; the human falls back to `curl|sh` if they don't use Homebrew.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Seeded from the issue; ideation refines, re-anchors each to the end value it measures, and fills the `Verified by` with a concrete falsifiable test.

**AC-1 — The install hint is OS-aware: a Linux host receives a command it can execute (the `curl|sh` path from `install.md` line 20), not a Homebrew command that fails on Linux.**
Verified by: A Go unit test that drives `tooOldBinaryRemedy` (or `RunDoctor` with a too-old-binary manifest) under a Linux OS context (`runtime.GOOS == "linux"` or an injected OS parameter) and asserts the output contains `curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh` and does NOT contain `brew upgrade`; and under a macOS context asserts the output contains `brew upgrade` and the `curl|sh` command is absent or is the secondary fallback. The independent baseline: the Homebrew command fails on Linux (no `brew`) — removing the Linux branch makes the test fail because the Linux host would see `brew upgrade`, a command that cannot run on Linux.

**AC-2 — On either abort class (binary absent, or binary present but wrong minor), the FO offers to run the install/upgrade itself and resumes startup once the binary lands, instead of hint-and-abort. The self-modifying boot converges (runs install, re-checks `--version`, proceeds if compatible, falls back to hint after one attempt — no loop, no no-op).**
Verified by: A behavior fixture that drives the version-gate abort path with a captive install command (a script that places a compatible `spacedock` binary at a known path) and asserts: the FO runs the install command (not just prints it), re-checks `${SPACEDOCK_BIN:-spacedock} --version`, and proceeds past the gate (exit/return indicating resume, not abort). The independent baseline: the current behavior requires the human to copy-paste a command and restart (1+ human interventions); the new behavior requires 0. The test fails if the FO reverts to print-and-exit (does not run the install), if it loops (attempts install more than once without converging), or if it no-ops (proceeds without re-checking `--version`).

**AC-3 — When the FO is running inside a sandbox, it detects the sandbox (from the `Sandbox:` line in `--version` output) and tells the human to run the exact install command themselves outside the sandbox, naming the exact command. It does NOT attempt the install inside the sandbox.**
Verified by: A behavior fixture (or Go unit test for the sandbox-aware message builder) that drives the version gate with a sandbox marker present (the `--version` output includes `Sandbox: inside (agent-safehouse)`) and asserts: the human-facing message names the exact OS-aware install command and the "outside the sandbox" instruction, and the FO does NOT run the install command. The independent baseline: a sandboxed install would install the binary somewhere the host cannot see (silent no-op). The test fails if the FO attempts the install inside the sandbox, or if the message omits the exact command or the "outside the sandbox" instruction.

## Test plan

- **AC-1 (OS-aware hint) — Go unit test**: drive `tooOldBinaryRemedy` (or `RunDoctor` with a too-old-binary manifest) under Linux and macOS OS contexts. Assert the Linux output contains the `curl|sh` command from `install.md` line 20 and not `brew`; the macOS output contains `brew upgrade` and not `curl|sh` as the primary remedy. Fixture test; ~30-50 lines. Mirrors the existing `TestTooOldBinaryRemedyLeadsWithBrew` pattern in `internal/contract/version_message_test.go`, adding the OS axis.
- **AC-2 (direct install/upgrade-and-resume) — behavior fixture**: drive the version-gate abort path with a captive install command (a script that places a compatible binary at a known path). Assert the FO runs the install, re-checks `--version`, and proceeds past the gate (does not loop, does not no-op). The convergence mechanism is partly prose (the skill text instructs the model) and partly binary (the doctor remedy the FO cites). The behavior fixture simulates the FO's version-gate flow. ~40-60 lines. Cost: moderate (requires a fixture that simulates the install-and-recheck loop).
- **AC-3 (sandbox detection) — behavior fixture or Go unit test**: drive the version gate with a sandbox marker present (`--version` output includes `Sandbox: inside (agent-safehouse)`). Assert the human-facing message names the exact install command and the "outside the sandbox" instruction, and the FO does NOT run the install. If the sandbox-aware message is a binary function (e.g. in contract or cli), a Go unit test suffices; if it's purely prose, a behavior fixture exercises the skill text. ~20-40 lines.
- **No-drift check**: a test or assertion that the `curl|sh` command in the hint/doctor remedy matches `docs/site/get-started/install.md` line 20 byte-for-byte, so the hint and the docs cannot drift. ~5-10 lines.

## Expected surface

Files:
1. `internal/contract/contract.go` — `tooOldBinaryRemedy` gains an OS-conditional branch: Linux leads with `curl|sh`, macOS keeps Homebrew. Source-build fallback on both. Thread OS through `compareNamed` or read `runtime.GOOS` directly. ~15-25 lines added.
2. `internal/contract/doctor.go` — `ManifestVerdict`/`RunDoctor` may thread OS to `compareNamed`/`tooOldBinaryRemedy` if the OS parameter is passed (not needed if `runtime.GOOS` is read inside `tooOldBinaryRemedy`). ~0-10 lines.
3. `internal/cli/init.go` — callers of `RunDoctor` pass OS or rely on `runtime.GOOS` internally. ~0-5 lines.
4. `skills/first-officer/references/first-officer-shared-core.md` — Startup step 1 rewrite: OS-aware hints, direct install/upgrade offer with convergence mechanism (set `SPACEDOCK_BIN` from install stderr, re-check `--version`, one-attempt bound), sandbox detection branch (parse `Sandbox:` line, tell human to run outside sandbox). ~40-60 lines (substantial rewrite of the abort classes).
5. `internal/contract/version_message_test.go` or new `internal/contract/remedy_os_test.go` — OS-aware remedy tests (Linux `curl|sh`, macOS `brew`). ~30-50 lines.

LOC tolerance: ±20% per file.

Observable semantics that may change:
- `spacedock doctor` on Linux prints the `curl|sh` upgrade path, not Homebrew (command output).
- The FO version gate offers to run install/upgrade and re-checks `--version` (boot abort→install-resume behavior).
- The FO version gate detects sandbox and tells the human to run the command outside the sandbox (sandbox detection, stderr content).
- The binary-absent hint in the skill prose includes the Linux `curl|sh` path (instruction text).

## Spike record

**Riskiest mechanism: self-modifying boot convergence (install then re-check `--version`)**

Spiked: YES — exercised the convergence path end-to-end with a locally-built binary and a simulated install.

Findings:

1. **`--version` already reports sandbox state**: `printVersion` in `internal/cli/cli.go` calls `safehouse.Inside(getenv)` and prints `Sandbox: inside (agent-safehouse)` when inside a sandbox. The FO model can parse this line from `--version` output — no new sandbox detection mechanism needed. Confirmed by running `spacedock --version` inside this session (which is sandboxed).

2. **`curl|sh` convergence risk is REAL but bounded**: `install.sh` installs to `~/.local/bin/spacedock` (default `SPACEDOCK_INSTALL_DIR`). On a system where `~/.local/bin` is NOT on `$PATH`, after `curl|sh` install, bare `spacedock --version` fails with exit 127 (command not found) — convergence failure. However:
   - `install.sh` prints `install.sh: installed spacedock <version> to <dir>/spacedock` to stderr (line 178).
   - `install.sh` warns if the install dir is not on PATH (lines 180-181): `note: <dir> is not on PATH; add it to run 'spacedock' directly`.
   - The FO convergence mechanism: after install, set `SPACEDOCK_BIN=$HOME/.local/bin/spacedock` (parsed from stderr or probed directly) and re-check `${SPACEDOCK_BIN} --version`. Confirmed: setting `SPACEDOCK_BIN` to the installed path makes `--version` succeed (exit 0).

3. **`brew upgrade` convergence is safe**: `brew upgrade spacedock` installs to the brew prefix (on `$PATH`), so `spacedock --version` resolves after upgrade without any `SPACEDOCK_BIN` override. Convergence is automatic.

4. **No OS line in `--version`**: `--version` does not report the host OS. The FO model detects OS via `uname -s` (a prose mechanism — the model runs `uname -s` in the bash tool and branches on `Darwin` vs `Linux`). This is sufficient; adding OS to `--version` is out of scope.

5. **Retry bound**: the FO must attempt install exactly once. If `--version` still fails after install, fall back to hint-and-abort. The skill prose must state this bound explicitly to prevent a loop.

**Conclusion**: the convergence mechanism works IF the FO (a) detects OS via `uname -s` to choose the install command, (b) checks sandbox state via `--version`'s `Sandbox:` line before offering to run install, (c) after `curl|sh` install, sets `SPACEDOCK_BIN` to the installed path and re-checks `--version`, (d) bounds to one retry then falls back to hint-and-abort. The brew path converges automatically; the `curl|sh` path converges via the `SPACEDOCK_BIN` override.

## Stage Report: ideation

- DONE: At least one AC per end-value (OS-aware hint, direct install/upgrade-and-resume, sandbox detection) MEASURES the outcome against an independent baseline that can move the wrong way — not a mechanism-shipped assertion like "the hint updates to X"
  AC-1 measures a Linux-no-`spacedock`/no-`brew` baseline (fails if hint is Homebrew-only or omits the documented Linux path); AC-2 measures the prior hint-and-abort baseline (fails if FO reverts to print-and-exit, loops, or no-ops without re-checking `--version`); AC-3 measures a sandboxed install landing where the host can't see it (fails if FO runs install inside the sandbox or omits the exact command / "outside the sandbox" instruction). See "Acceptance criteria", commit 11ffc7d31.
- DONE: The riskiest mechanism is spiked first and the result recorded: a self-modifying boot that runs the install/upgrade then re-invokes `--version` must converge (not loop, not no-op); "no spike needed" only if backed by named proven mechanisms
  Spike record ("Spiked: YES") exercised boot convergence end-to-end with a locally-built binary + simulated install; convergence works iff FO sets `SPACEDOCK_BIN` to the installed path post-`curl|sh` and re-checks `--version` with a one-retry bound, while the brew path converges automatically — so it converges (not loop/no-op), not a "no spike needed" claim. See "Spike record", commit 11ffc7d31.
- DONE: Expected surface names the files + LOC tolerance AND the observable semantics it may change (boot abort→install-resume behavior, sandbox detection, stderr content), with the three end-values each served by a named AC
  "Expected surface" lists 5 files with per-file LOC estimates and "LOC tolerance: ±20% per file", plus four observable-semantics changes (doctor Linux output, boot abort→install-resume, sandbox detection / stderr content, binary-absent hint text); AC-1/AC-2/AC-3 each serve one end-value (OS-aware hint, install-resume, sandbox detection). See "Expected surface", commit 11ffc7d31.

### Summary

Ideation design is complete and committed (11ffc7d31): OS-aware hint, direct install/upgrade-and-resume with a one-attempt convergence bound, and sandbox detection, each backed by an end-value-measuring AC and a riskiest-mechanism spike. This run fixes the completion gap by appending the missing `## Stage Report: ideation` section to the entity file (the prior run placed the completion summary in the output message instead of the file) and committing it; no design section or frontmatter was modified.
