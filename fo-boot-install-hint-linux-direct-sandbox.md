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

The first-officer boot hits the binary version gate (Startup step 1) and, when the `spacedock` binary is absent from PATH (the install class), prints a Mac-only Homebrew install hint and stops, leaving the human to copy-paste a command and restart the session — a hint a Linux host cannot even run. This task improves that install journey along the axes the issue names: make the install hint OS-aware (include the documented Linux `curl|sh` path, not just Homebrew), offer to run the install directly and resume startup once the binary lands (turn hint-and-abort into one approved action, bounded to a single attempt), and detect sandboxed execution so a sandboxed install does not silently no-op (tell the human to run the install command themselves outside the sandbox, naming the exact command). Cycle 2 (captain binding annotations) adds a fourth piece: an OS line in `--version` output, useful both to the FO's hint logic and to issue reports. The upgrade journey (binary present but wrong minor → check latest release, offer upgrade-and-resume) is filed separately as `fo-boot-upgrade-hint-latest-release` and is out of scope here.

## Problem

Today the FO version gate (`first-officer-shared-core.md` Startup step 1) hits the **binary-absent** class with a Mac-only, hint-and-abort response:

1. **Binary absent** — the FO prints `brew install spacedock-dev/homebrew-tap/spacedock` and stops. A Linux VM with no `spacedock` on PATH sees a Homebrew command it cannot run (no `brew` on Linux). The documented Linux path (`curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh`, `docs/site/get-started/install.md` line 20) is never mentioned.
2. **No self-repair path** — the human must copy-paste the install command themselves and restart the session. The gate could offer to run the install once the operator approves, then re-check `--version` and resume startup.
3. **Sandboxed execution** — the FO doesn't check sandbox state at the version gate. A sandboxed session that ran the install inside the sandbox would install the binary somewhere the host cannot see — a silent no-op from the human's perspective. The detection mechanism (`safehouse.Inside()`) already exists and `--version` already prints `Sandbox: inside (agent-safehouse)`, but the version-gate prose ignores it.
4. **No OS in `--version`** — `--version` reports version, runtime, sandbox, and the frozen contract token, but not the host OS/arch, which would help both the FO's hint selection and user issue reports (captain annotation 1).

The **binary-present-but-wrong-minor** class stays on its current journey: ABORT with the mismatch message and run `${SPACEDOCK_BIN:-spacedock} doctor` for the remedy. The interactive upgrade journey (check latest release, offer upgrade-and-resume) is a separate task (`fo-boot-upgrade-hint-latest-release`); this task touches the doctor remedy only to make its install-command text OS-correct (a doctor on Linux must not lead with `brew`). See Design decision D-4.

## Proposed approach

Four coordinated changes across the binary and the FO skill prose, plus the explicit design decisions the staff review required.

### 1. OS line in `--version` output (binary; captain annotation 1)

- `internal/cli/cli.go printVersion` prints `OS: <runtime.GOOS>/<runtime.GOARCH>` (e.g. `OS: linux/amd64`, `OS: darwin/arm64`) as **line 2, in both output shapes** — the outside-runtime human case (currently a single line, pinned by `internal/cli/version_session_test.go:65-69` and retargeted by this task) and the in-session case. Line 1 stays `spacedock <version>` and nothing else (load-bearing; the FO gate parses only line 1, so a line-2 insertion is gate-safe).
- The FO prose/hint logic reads OS from the `^OS: ` line of `--version` output when the line is present (prefix-anchored parse), replacing the blanket `uname -s` prose detection from cycle 1.
- **Transition caveat (design decision D-1):** a binary old enough to fail the version gate predates the OS line, and a truly absent binary prints nothing at all — in both cases there is no OS line to read. For those two cases only, the FO prose falls back to `uname -s` (`Darwin` / `Linux` / other) and names the fallback explicitly in the skill text. Once the installed binary is new enough to pass the gate, the OS line governs all later hint/issue-report logic. This is the one retention of `uname -s` and is surfaced here for the gate rather than silently reintroduced.

### 2. OS-aware install hint

- **Skill prose** (`first-officer-shared-core.md` Startup step 1, binary-absent class): the hint is OS-conditional — on Linux it leads with the documented `curl|sh` path from `install.md` line 20; on macOS it keeps the Homebrew install lead. The source-build fallback stays on both. On a non-Darwin/Linux OS, hint-and-abort with the source-build note (D-5).
- **Binary** (`internal/contract/contract.go`): `tooOldBinaryRemedy(edgeCask, goos string)` gains an OS-conditional remedy lead. `goos` is threaded explicitly (not read from `runtime.GOOS` inside the function) so tests pin the OS deterministically; the caller (`compareNamed`) supplies `runtime.GOOS`. On Linux the remedy leads with the `curl|sh` command; on macOS the Homebrew lead is unchanged. This keeps doctor's remedy text runnable on the host it prints on — doctor is the wrong-version class's remedy surface and its Mac-only lead is the same defect class this task fixes (D-4). Precedence (N6): on `runtime.GOOS == "linux"` the `curl|sh` lead **supersedes** the edge-cask branch — `spacedock@next` is a brew-only token, meaningless on Linux, so edgeCask is ignored there.

### 3. Direct install-and-resume (install path only; one bounded attempt)

- The binary-absent class changes from hint-and-abort to offer-and-resume: the FO offers to run the OS-aware install command itself, the operator approves, the FO runs it, then re-checks `${SPACEDOCK_BIN:-spacedock} --version` and resumes startup if line 1 parses to a compatible version.
- **Convergence mechanism** (cycle-1 spike): `install.sh` prints `install.sh: installed spacedock <version> to <dir>/spacedock` to stderr and warns when the dir is not on PATH. After a `curl|sh` install, the FO sets `SPACEDOCK_BIN` to the installed path (parsed from stderr, or probed at `$HOME/.local/bin/spacedock`) and re-checks `--version`. The wrong-version/brew-upgrade convergence discussion moves to the upgrade task.
- **Loop bound (design decision D-3 — the staff review's B1 demand for a falsifiable fallback on a boot path):** one install attempt per session, enforced by a **binary-visible guardrail, not prose alone**: the FO checks `SPACEDOCK_INSTALL_ATTEMPTED` before offering to run an install; if set, it goes straight to the hint-and-abort fallback. The prose sets the marker when it runs the install. If the post-install `--version` re-check still fails, the FO falls back to hint-and-abort with the exact OS-aware command for the human to run manually — no second install attempt, no proceed-without-recheck. Rationale for choosing the marker over prose-only: the version gate runs on every FO boot and is the highest-blast-radius path; an env token gives the failure-fallback behavior fixture a deterministic, observable bit to assert on (spiked — see Spike record). The marker is session-scoped env; it does not persist across sessions and does not guard the next session (D-2 states what the next session does instead).

### 4. Sandbox detection

- At the version gate, the FO parses the `Sandbox:` line from `--version` output, **anchored on the `^Sandbox:` prefix** (N9 — not a loose `inside` substring), three-way: no such line → treat as not-inside but note the parse assumption; `Sandbox: not sandboxed …` → not inside; `Sandbox: inside (…)` → inside.
- **Coupling (N9, stated):** `printVersion` prints the `Sandbox:` line only when inside a detected runtime session. The FO version gate always runs in-session by construction, so the line is present in the gate's `--version` output; the design relies on that and names it.
- If inside a sandbox, the FO does NOT offer to run the install (it would no-op for the host). It tells the human plainly: "You're running inside a sandbox. Run this command yourself, outside the sandbox: `<exact OS-aware install command>`."
- **AC-3 test-type disposition (N8):** the sandbox-aware message stays skill prose (no new binary seam — the OS-aware command text already exists binary-side in the remedy, and a new `doctor --install-hint` arm would duplicate it). The test plan therefore uses a behavior fixture and the plan **names the weakness explicitly**: LLM prose is non-deterministic, so the fixture asserts semantic content (message carries the exact command token and the "outside the sandbox" instruction; install command is not executed) rather than byte-for-byte wording.

### Design decisions (explicit, per staff review)

- **D-1 (OS line + transition):** `OS: <goos>/<goarch>` as `--version` line 2 in both output shapes; FO prefers the line, falls back to `uname -s` only when no `--version` OS line exists (absent or pre-feature binary). Retargets the pinned one-line human output — an intended observable-semantics change, listed below.
- **D-2 (SPACEDOCK_BIN override × launcher invariant; durability — B2):** the post-install `SPACEDOCK_BIN` override is **in-process, launching-session scope** and happens **at the gate** — it is the gate's launcher resolution ("resolve ONE launcher at the version gate; use THAT launcher for every later helper call"), not a mid-session drift onto a bare `spacedock`, so it satisfies the shared-core invariant rather than colliding with it. If `SPACEDOCK_BIN` was set-but-stale (the cause of the gate failure), the override replaces it for this session and the fallback message says so. **Durability model: session-only.** The override is NOT persisted to a shell profile in this task. The **next session's** bootstrap re-runs the gate: if the install dir is on PATH, bare `spacedock` resolves and the session converges with no override; if not, the gate fails again and the fallback hint names the exact installed path and tells the human to add the dir to PATH (or launch with `SPACEDOCK_BIN=<path>`). **Surfaced captain-decision item:** whether spacedock should ever *persist* the override (write the install dir to the user's shell profile, or record it in config) is left to the captain/gate — the recommended default is no auto-profile-writes.
- **D-3 (loop-bound guardrail):** `SPACEDOCK_INSTALL_ATTEMPTED` env marker checked before any install attempt (spiked, see Spike record), plus a dedicated failure-fallback behavior fixture in the test plan. Prose-only loop control on a boot path judged insufficient.
- **D-4 (doctor relationship — N10):** the wrong-version class is unchanged: ABORT with the mismatch message and run `${SPACEDOCK_BIN:-spacedock} doctor`; the FO does **not** offer a direct upgrade in this task (that journey is `fo-boot-upgrade-hint-latest-release`). The direct install offer applies to the binary-absent class only, and the binary-absent class still does **not** run doctor (there is no binary to run). Doctor's changed Linux remedy output (curl|sh lead instead of `brew upgrade`) is an **intended observable semantic** of this task, not a side effect, and is listed in Expected surface.
- **D-5 (OS boundaries — N6, N7):** edgeCask ignored on Linux (curl|sh supersedes). Non-Darwin/Linux OS (from `--version` OS line or `uname -s`): hint-and-abort; `install.sh` already rejects non-darwin/linux, so the message names the source build and notes the unsupported OS.

## Out of scope

- **The upgrade journey** — checking the latest release, hinting the user that an upgrade exists, offering upgrade-and-resume, and the `brew upgrade` convergence spike all live in the separately filed `fo-boot-upgrade-hint-latest-release`. No upgrade-path content is designed or tested here beyond keeping doctor's remedy text OS-correct (D-4).
- Persisting `SPACEDOCK_BIN` (or the install dir) to a shell profile or config file — session-only override (D-2); any persistence is a captain decision surfaced by this design, not implemented here.
- Changing the `safehouse.Inside()` detection mechanism itself — it already works; this task uses its output, not changes it.
- Windows support — `install.sh` already rejects non-darwin/linux; D-5 names the source-build fallback for unrecognized OS values.
- Changing the brew formula, `install.sh`, or `docs/site/get-started/install.md` — the canonical commands stay as-is; the hint cites them, not drifts from them (and the no-drift check enforces this).
- Auto-detecting whether `brew` is installed on macOS — the hint leads with Homebrew regardless; the human falls back to `curl|sh` if they don't use Homebrew.
- The Homebrew **upgrade** command form (`brew upgrade spacedock`, `contract.go:236`) in the no-drift check — the upgrade task owns its own drift check; this task's no-drift check covers the two install commands (see Test plan).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. ACs are install-scoped (captain annotation 2) and re-anchored to independent baselines.

**AC-1 — The install hint is OS-aware: a Linux host receives a command it can execute (the `curl|sh` path from `install.md` line 20), not a Homebrew command that fails on Linux.**
Verified by: A Go unit test that drives `tooOldBinaryRemedy` with a pinned `goos` parameter: under `"linux"` it asserts the remedy contains `curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh` and contains NO `brew` lead (including no `spacedock@next` — the edge token is superseded on Linux, D-5); under `"darwin"` it asserts the remedy keeps the Homebrew lead and keeps the edge-cask distinction byte-for-byte as today. Plus a skill-prose content check that the binary-absent hint in `first-officer-shared-core.md` carries both the curl|sh and brew-install command forms. The independent baseline: the Homebrew command fails on Linux (no `brew`) — removing the Linux branch makes the test fail because a Linux host would see a `brew` command that cannot run.

**AC-2 — On the binary-absent class, the FO offers to run the install itself and resumes startup once a compatible binary lands, instead of hint-and-abort. The self-modifying boot converges (runs install, re-checks `--version`, proceeds if compatible) and, on failure, falls back to hint-and-abort after exactly one attempt — no loop, no no-op.**
Verified by TWO behavior fixtures driving the version-gate abort path with captive commands:
1. **Convergence fixture** (cycle-1 shape): a captive install script places a compatible `spacedock` binary at a known path; assert the FO runs the install (not just prints it), repoints `SPACEDOCK_BIN`, re-checks `--version`, and proceeds past the gate.
2. **Failure-fallback fixture** (new, answers B1/N2): a captive install script succeeds but produces an incompatible (or no) binary; assert the FO runs the install exactly once (the `SPACEDOCK_INSTALL_ATTEMPTED` marker is set and a second simulated gate entry does not re-run the install), re-checks `--version`, observes failure, and emits hint-and-abort naming the exact OS-aware command — and does NOT proceed past the gate.
The independent baseline: the current behavior requires the human to copy-paste a command and restart (1+ human interventions); the new behavior requires 0 on success and exactly 1 bounded attempt on failure. These fixtures fail if the FO reverts to print-and-exit, if it attempts install more than once (guardrail absent or ignored), or if it proceeds without re-checking `--version`.

**AC-3 — When the FO is running inside a sandbox, it detects the sandbox (from the `Sandbox: ... inside` line in `--version` output, prefix-anchored on `^Sandbox:`) and tells the human to run the exact install command themselves outside the sandbox, naming the exact command. It does NOT attempt the install inside the sandbox.**
Verified by: A behavior fixture that drives the version gate with a sandbox marker present (`--version` output includes `Sandbox: inside (agent-safehouse)`) and asserts: the human-facing message carries the exact OS-aware install command token and the "outside the sandbox" instruction, and the install command is not executed (no captive-install side effect observed). Named weakness (N8 disposition in Proposed approach §4): the message is skill prose, so the assertion is on semantic content, not byte-for-byte prose. The independent baseline: a sandboxed install lands where the host cannot see it (silent no-op) — the fixture fails if the FO runs the install inside the sandbox or omits the exact command / "outside the sandbox" instruction.

**AC-4 — `spacedock --version` reports the host OS and arch (`OS: <goos>/<goarch>`) in both output shapes (in-session and outside-runtime human), so issue reports carry the platform and the FO hint logic can read OS from `--version`.**
Verified by: A Go unit test on `printVersion` asserting the `OS: ` line appears as line 2 with the correct `runtime.GOOS`/`runtime.GOARCH` values in both the one-line-human case and the in-session case; the retargeted `internal/cli/version_session_test.go` pins both shapes. The independent baseline: today's `--version` output contains no OS token at all (verified — `spacedock %s`, `Runtime:`, `Sandbox:`, `contract 3` only), so dropping the OS line fails the test.

## Test plan

- **AC-1 (OS-aware hint) — Go unit test**: drive `tooOldBinaryRemedy(edgeCask, goos)` with pinned `goos` values. `"linux"`: contains the `curl|sh` command from `install.md` line 20 and no `brew` content (edge token superseded); `"darwin"`: keeps the Homebrew lead and the edge-cask byte-for-byte block. ~30-50 lines, in `internal/contract/version_message_test.go` or a new `internal/contract/remedy_os_test.go`.
- **Existing-test retarget (N3, mandatory)**: `internal/contract/version_message_test.go` is **modified**: `TestTooOldBinaryRemedyLeadsWithBrew` (currently asserts `brew upgrade spacedock` via `RunDoctor`, which will produce the curl|sh lead on the ubuntu CI leg) is retargeted to assert the host-appropriate remedy branch on `runtime.GOOS` (or rerouted to direct `tooOldBinaryRemedy` calls with pinned `goos`); `TestTooOldBinaryRemedyEdgeChannel` (byte-for-byte pinned block, direct calls) passes an explicit `goos = "darwin"` (or the new signature) so its pinned block is GOOS-stable. This is what keeps `.github/workflows/runtime-live-e2e.yml`'s ubuntu leg green.
- **AC-2 (direct install-and-resume) — two behavior fixtures**: (1) convergence fixture: captive install places a compatible binary at a known path; assert install ran once, `SPACEDOCK_BIN` repointed, `--version` re-checked, gate passed. (2) **failure-fallback fixture** (B1): captive install "succeeds" but yields an incompatible binary; assert exactly one install run (marker observed; a second simulated gate entry runs zero installs), `--version` re-check observed failing, hint-and-abort emitted with the exact command, gate not passed. Each ~40-60 lines; moderate cost (fixture harness simulating the FO gate flow).
- **AC-3 (sandbox detection) — behavior fixture** (explicitly NOT a unit test, per the N8 disposition in Proposed approach §4): drive the gate with `Sandbox: inside (agent-safehouse)` in `--version` output; assert the message carries the exact command token + "outside the sandbox" instruction, and the captive install was not executed. ~20-40 lines.
- **AC-4 (OS line) — Go unit test**: `printVersion` writes `OS: <goos>/<goarch>` as line 2 in the outside-runtime case and the in-session case; `internal/cli/version_session_test.go` retargeted (its `outside-every-runtime-is-one-line` case becomes a two-line case). ~15-25 lines total.
- **No-drift check (widened per N4, extraction defined per N5)**: a repo-relative test asserts (a) the `curl|sh` token in the shared-core hint equals `install.md`'s fenced command at line 20 **after stripping the markdown fence's 4-space indentation** (extraction: the single line inside the `=== "Binary (macOS / Linux)"` fence matching `^curl `, trimmed), and (b) the shared-core brew **install** hint refers to the same tap and formula as `install.md`'s Homebrew tab (`spacedock-dev/homebrew-tap` + `spacedock`; the two-line `brew tap` + `brew install` form and the one-line full-token form are checked for token equality, not formatting equality). The brew **upgrade** form in `contract.go` is excluded — the upgrade task owns it (Out of scope). ~10-15 lines.
- **Skill smoke test**: per repo convention, skill text changes to `first-officer-shared-core.md` get a smoke test assert for the new gate text (OS-aware hint shape, marker name, sandbox instruction). Follow the existing skill smoke-test fixture.

## Expected surface

Files:
1. `internal/cli/cli.go` — `printVersion` gains the `OS: <goos>/<goarch>` line 2 in both output shapes (move it above the outside-runtime early return). ~4-8 lines added.
2. `internal/cli/version_session_test.go` — **modified** (required): the pinned `outside-every-runtime-is-one-line` case becomes two lines; in-session cases gain the OS line; a GOOS/GOARCH assertion is added. ~10-20 lines changed.
3. `internal/contract/contract.go` — `tooOldBinaryRemedy` gains a `goos string` parameter and an OS-conditional lead: Linux leads with `curl|sh` (edgeCask superseded), macOS keeps Homebrew. `compareNamed` passes `runtime.GOOS`. ~15-25 lines added/changed.
4. `internal/contract/version_message_test.go` — **modified** (required, N3): `TestTooOldBinaryRemedyLeadsWithBrew` and `TestTooOldBinaryRemedyEdgeChannel` retargeted with pinned/host-threaded GOOS so the ubuntu CI leg stays green; AC-1 Linux/darwin assertions may land here. ~20-40 lines changed.
5. `internal/contract/doctor.go` — only if the `goos` thread ripples past `compareNamed`; preferred shape keeps `RunDoctor`'s signature unchanged (read `runtime.GOOS` at `compareNamed`). ~0-10 lines.
6. `internal/cli/init.go`, `internal/cli/frontdoor.go` — unchanged if `RunDoctor` keeps its signature. ~0 lines.
7. `skills/first-officer/references/first-officer-shared-core.md` — Startup step 1 rewrite of the binary-absent class: OS-aware install hint (prefer `--version` OS line, `uname -s` transition fallback), direct install-and-resume offer with the `SPACEDOCK_INSTALL_ATTEMPTED` guard and one-attempt bound, `SPACEDOCK_BIN` session-scoped convergence, sandbox branch (`^Sandbox:` anchored parse, no-install-inside message). ~40-60 lines (substantial rewrite of the abort classes).
8. Skill smoke test fixture (name per existing convention) — asserts the new gate text shape. ~15-30 lines.
9. No-drift test (repo-relative; location alongside the skill smoke fixture or `internal/contract`) — curl|sh + brew-install token checks against `install.md`. ~10-15 lines.

LOC tolerance: ±20% per file.

Observable semantics that may change:
- `spacedock --version` gains an `OS: <goos>/<goarch>` line 2 in **both** shapes — including the outside-runtime human output previously pinned to a single line (intended, D-1).
- `spacedock doctor`'s too-old-binary remedy leads with the `curl|sh` install command on Linux instead of `brew upgrade ...` (intended, D-4); edge-cask naming is unchanged on macOS and ignored on Linux.
- The FO version gate's binary-absent class offers to run the install and re-checks `--version` (abort→install-and-resume behavior), bounded by the `SPACEDOCK_INSTALL_ATTEMPTED` session marker (new env semantics: session-scoped, prose-set, prose-checked).
- The FO version gate detects sandbox state via the `^Sandbox:` line and reroutes to a human-run-outside-sandbox instruction (stderr/message content).
- The FO's session-scoped `SPACEDOCK_BIN` override after install (launcher resolution at the gate; not persisted, D-2).

## Spike record

**Cycle 2 spike: the one-attempt failure-fallback path (the B1/N2 gap — cycle 1 exercised only the happy convergence path)**

Spiked: YES — exercised with shell code mirroring the intended gate flow, a captive failing install, and the proposed `SPACEDOCK_INSTALL_ATTEMPTED` guardrail.

Setup: a captive `install.sh` that reports success (`install.sh: installed spacedock 0.20.0 to /tmp/spike-fallback/bin/spacedock` on stderr, mirroring the real installer's stderr contract) but drops an **incompatible** binary (`--version` prints `spacedock 0.20.0`, below the required minor). Gate flow: check marker → run install once → set marker → repoint `SPACEDOCK_BIN` to the installed path → re-check `--version`.

Findings:

1. **Failure fallback converges to hint-and-abort, not a loop**: pass 1 ran the install once, observed the re-check failure (`spacedock 0.20.0`), and emitted the fallback carrying the exact OS-aware command (exit 3, gate NOT passed). Recorded output: `FALLBACK: --version still incompatible (spacedock 0.20.0) -> hint-and-abort with: curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh`.
2. **The guardrail blocks a second attempt**: pass 2 simulated a same-session gate re-entry with the marker set and a sentinel file removed — the install did NOT re-run (sentinel absent), and the flow went straight to hint-and-abort. `GUARDRAIL-OK: install did not re-run`. This is the deterministic bit the failure-fallback behavior fixture asserts on (D-3 chosen over prose-only control).
3. **The in-session `SPACEDOCK_BIN` override mechanics work as the cycle-1 spike showed** (repoint → `--version` resolves against the installed path); D-2 states its invariant relationship and session-only durability. Not re-spiked beyond the pass-1 repoint, which the flow exercised.

**Cycle-1 spike results retained (proven mechanisms, not re-spiked):**

- `--version` already reports sandbox state: `printVersion` (`internal/cli/cli.go:837-854`) calls `safehouse.Inside(getenv)` and prints `Sandbox: inside (agent-safehouse)` — confirmed live (`spacedock --version` in this session shows the line). The FO parses the existing line; no new detection mechanism.
- `curl|sh` happy-path convergence: `install.sh` prints the installed path to stderr and warns when the dir is off PATH; setting `SPACEDOCK_BIN` to the installed path makes the `--version` re-check succeed (exit 0).
- The gate parse targets line 1 only, so inserting `OS:` at line 2 is parse-safe; the outside-runtime output currently ends at line 1 (pinned) and will be retargeted (AC-4).

**Deferred:** the `brew upgrade` convergence spike moves to `fo-boot-upgrade-hint-latest-release` (captain annotation 2 — upgrade journey out of this task). The cycle-1 record's brew-upgrade assertion is superseded by that task's spike obligation.

**No spike needed for the OS line**: the seam was read directly — `printVersion` assembles output from `displayVersion()` + stdlib `runtime` values with no I/O; inserting `fmt.Fprintf(w, "OS: %s/%s\n", runtime.GOOS, runtime.GOARCH)` after line 1 is a proven-trivial mechanism (same class as the existing `Sandbox:` line from the same stdlib).

## Stage Report: ideation

- DONE: At least one AC per end-value (OS-aware hint, direct install/upgrade-and-resume, sandbox detection) MEASURES the outcome against an independent baseline that can move the wrong way — not a mechanism-shipped assertion like "the hint updates to X"
  AC-1 measures a Linux-no-`spacedock`/no-`brew` baseline (fails if hint is Homebrew-only or omits the documented Linux path); AC-2 measures the prior hint-and-abort baseline (fails if FO reverts to print-and-exit, loops, or no-ops without re-checking `--version`); AC-3 measures a sandboxed install landing where the host can't see it (fails if FO runs install inside the sandbox or omits the exact command / "outside the sandbox" instruction). See "Acceptance criteria", commit 11ffc7d31.
- DONE: The riskiest mechanism is spiked first and the result recorded: a self-modifying boot that runs the install/upgrade then re-invokes `--version` must converge (not loop, not no-op); "no spike needed" only if backed by named proven mechanisms
  Spike record ("Spiked: YES") exercised boot convergence end-to-end with a locally-built binary + simulated install; convergence works iff FO sets `SPACEDOCK_BIN` to the installed path post-`curl|sh` and re-checks `--version` with a one-retry bound, while the brew path converges automatically — so it converges (not loop/no-op), not a "no spike needed" claim. See "Spike record", commit 11ffc7d31.
- DONE: Expected surface names the files + LOC tolerance AND the observable semantics it may change (boot abort→install-resume behavior, sandbox detection, stderr content), with the three end-values each served by a named AC
  "Expected surface" lists 5 files with per-file LOC estimates and "LOC tolerance: ±20% per file", plus four observable-semantics changes (doctor Linux output, boot abort→install-resume, sandbox detection / stderr content, binary-absent hint text); AC-1/AC-2/AC-3 each serve one end-value (OS-aware hint, install-resume, sandbox detection). See "Expected surface", commit 11ffc7d31.

### Summary

Ideation design is complete and committed (11ffc7d31): OS-aware hint, direct install/upgrade-and-resume with a one-attempt convergence bound, and sandbox detection, each backed by an end-value-measuring AC and a riskiest-mechanism spike. This run fixes the completion gap by appending the missing `## Stage Report: ideation` section to the entity file (the prior run placed the completion summary in the output message instead of the file) and committing it; no design section or frontmatter was modified.

### Feedback Cycles

- Cycle 1: REVISE — captain via subspace:r gate review (q0 worktree skill+binary dogfood, resolution:binding-1785472404632857000), following a fresh-context staff review (`.pi-subagents/artifacts/7177028a_reviewer_0_output.md`); surface: N/A — ideation redo round, no worktree/implementation branch exists to numstat against (the estimate-vs-actual diff grammar presumes an implementation candidate; none yet); AC narrowed + expanded — upgrade path removed from scope (captain: "let's focus on install here. i want a separate upgrade journey where we can check latest release and hint the user for upgrade. file that separately"); `--version` OS reporting pulled into scope (captain: "actually i think this is a good idea, also helpful when user reports issue" — annotating the Out-of-scope item that excluded an OS line). This is the captain-visible design-reset decision of record for the narrowed/expanded AC; correction round re-enters ideation (no `feedback-to` on this stage).

  Concrete revise asks (captain annotations, binding):
  1. Focus this task on INSTALL. The upgrade journey (check latest release, hint the user to upgrade, optionally upgrade-and-resume) is OUT of this task's scope — filed separately.
  2. Add an OS line to `--version` output (binary-side, `internal/cli`), replacing `uname -s` prose detection as the review's described mechanism; helpful for user issue reports.
  3. Fold in the material staff-review findings: (a) the one-attempt loop bound needs a falsifiable fallback path — happy-path-only fixture is insufficient; consider a binary guardrail; (b) the `SPACEDOCK_BIN` override must state its interaction with the launcher invariant and its durability model; (c) `internal/contract/version_message_test.go` GOOS-dependent assertions break on the ubuntu CI leg — retarget them and list the file in the expected surface; (d) spike the one-attempt-fallback path; the brew-upgrade spike moves to the upgrade task; (e) widen or rescope the no-drift check (brew command has three forms), soften byte-for-byte vs markdown indentation, name `edgeCask`×Linux precedence and non-Darwin/Linux fallback in the boundary notes.
## Stage Report: ideation (cycle 2)

- DONE: Design narrowed to install-only and expanded to include an OS line in --version output, per the captain's two binding annotations — both are applied verbatim with no leftover upgrade-path content anywhere in Problem/Approach/ACs/Out-of-scope
  Problem/Approach/ACs/Out-of-scope rewritten: wrong-version class is unchanged (ABORT + doctor, D-4), upgrade journey (latest-release check, upgrade offer, brew-upgrade convergence) moved to `fo-boot-upgrade-hint-latest-release`; OS line added as Proposed approach §1 + AC-4 (`internal/cli/cli.go printVersion`, `OS: <goos>/<goarch>` line 2, both output shapes).
- DONE: The three material staff-review gaps are closed in the revised design: (a) loop bound has a falsifiable one-attempt-failure-fallback path (fixture or binary guardrail, named choice); (b) SPACEDOCK_BIN override x launcher-invariant interaction and durability model (session-only vs persisted) are stated as explicit design decisions; (c) internal/contract/version_message_test.go is in the expected surface with GOOS-threaded retargeted tests
  (a) D-3 names the choice: `SPACEDOCK_INSTALL_ATTEMPTED` env guardrail checked before any install attempt + a dedicated failure-fallback behavior fixture (AC-2, fixture 2). (b) D-2: override is in-process launching-session scope at the gate (satisfies the resolve-one-launcher invariant), session-only durability, next-session behavior stated, profile-persistence surfaced as a captain-decision item with a no-persist recommendation. (c) Expected surface file 4 lists `version_message_test.go` modified with `TestTooOldBinaryRemedyLeadsWithBrew`/`TestTooOldBinaryRemedyEdgeChannel` retargeted (pinned/host-threaded GOOS); `version_session_test.go` likewise retargeted for the OS line.
- DONE: The one-attempt-fallback convergence path is spiked and recorded (brew-upgrade spike deferred to the separate upgrade task); ACs re-anchor to the narrowed scope with independent baselines, and expected surface lists the new --version OS-line file touch with observable-semantics delta
  Spike record cycle 2: captive failing install → exactly one attempt, re-check failure observed, hint-and-abort with exact command; marker-set re-entry ran zero installs (sentinel-verified). AC-2 now requires BOTH convergence and failure-fallback fixtures; each AC names its independent baseline; Expected surface files 1-2 carry the --version OS-line touch and the observable-semantics list names it.

### Staff-review dispositions

- B1 (material): closed via D-3 guardrail + failure-fallback fixture, spiked (see Spike record).
- B2 (material): closed via D-2 (invariant relationship, session-only durability, next-session behavior, captain-decision item surfaced).
- N1: superseded — brew-upgrade spike deferred to `fo-boot-upgrade-hint-latest-release` (captain annotation 2).
- N2: closed — failure-fallback path spiked; sandbox-conservatism accepted by design (message names exact command; no false-negative risk beyond telling the human to run it outside, which always works).
- N3 (material, CI): closed — `version_message_test.go` in Expected surface (modified) with GOOS-pinned retarget in the Test plan.
- N4: closed — no-drift check widened to the two install command forms (curl|sh + brew install token equality vs install.md); brew upgrade form excluded to the upgrade task.
- N5: closed — extraction defined (fence indentation stripped, `^curl ` line trimmed).
- N6: closed — D-5: curl|sh supersedes the edge-cask brew token on Linux; AC-1 asserts no `spacedock@next` on the Linux branch.
- N7: closed — D-5: non-Darwin/Linux falls back to hint-and-abort with the source-build note.
- N8: closed — AC-3 uses a behavior fixture; the prose non-determinism weakness is named explicitly in Proposed approach §4 and AC-3.
- N9: closed — parse anchored on `^Sandbox:` prefix; in-session coupling stated (Proposed approach §4).
- N10: closed — D-4: wrong-version class still ABORTs and runs doctor; doctor's Linux remedy change is an intended observable semantic; direct offer is install-class only.

### Summary

Cycle-2 ideation revises the design body in place: install-only scope (upgrade journey deferred), an `OS: <goos>/<goarch>` line 2 in `--version` both shapes (AC-4, retargeting the pinned one-line human output), a spiked one-attempt failure-fallback guarded by `SPACEDOCK_INSTALL_ATTEMPTED`, and explicit design decisions D-1..D-5 closing B1/B2/N3-N10. One surfaced captain-decision item: whether to ever persist the post-install launcher path (recommended default: no; session-only override, next session self-heals via PATH or the fallback hint naming the installed path).
