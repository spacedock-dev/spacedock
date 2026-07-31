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
3. **Sandboxed execution** — the FO doesn't check sandbox state at the version gate. A sandboxed session that ran the install inside the sandbox would install the binary somewhere the host cannot see — a silent no-op from the human's perspective. The detection mechanism (`safehouse.Inside()`) already exists and `--version` already prints `Sandbox: inside (agent-safehouse)`, but the version-gate prose ignores it — and that line can never exist in the binary-absent class anyway (no binary, no `--version` output; NF1), so detection for the install offer must come from the session's env markers, which `safehouse.Inside` itself reads.
4. **No OS in `--version`** — `--version` reports version, runtime, sandbox, and the frozen contract token, but not the host OS/arch, which would help both the FO's hint selection and user issue reports (captain annotation 1).

The **binary-present-but-wrong-minor** class stays on its current journey: ABORT with the mismatch message and run `${SPACEDOCK_BIN:-spacedock} doctor` for the remedy. The interactive upgrade journey (check latest release, offer upgrade-and-resume) is a separate task (`fo-boot-upgrade-hint-latest-release`); this task touches the doctor remedy only to make its install-command text OS-correct (a doctor on Linux must not lead with `brew`). See Design decision D-4.

## Proposed approach

Four coordinated changes across the binary and the FO skill prose, plus the explicit design decisions the staff review required.

### 1. OS line in `--version` output (binary; captain annotation 1)

- `internal/cli/cli.go printVersion` prints `OS: <runtime.GOOS>/<runtime.GOARCH>` (e.g. `OS: linux/amd64`, `OS: darwin/arm64`) as **line 2, in both output shapes** — the outside-runtime human case (currently a single line, pinned by `internal/cli/version_session_test.go:64-67` and retargeted by this task) and the in-session case. Line 1 stays `spacedock <version>` and nothing else (load-bearing; the FO gate parses only line 1, so a line-2 insertion is gate-safe).
- The FO prose/hint logic reads OS from the `^OS: ` line of `--version` output when the line is present (prefix-anchored parse), replacing the blanket `uname -s` prose detection from cycle 1.
- **OS-source caveat (design decision D-1; NF4 restatement):** a binary old enough to fail the version gate predates the OS line, and a truly absent binary prints nothing at all — in both cases there is no OS line to read. The FO prose therefore obtains OS in the binary-absent class from `uname -s` (`Darwin` / `Linux` / other), and names this explicitly in the skill text. The `uname -s` path is **permanently load-bearing, not transitional and not circular** (NF4): the binary-absent class IS the install-offer class — the one class whose hint is OS-conditional — and it can never have an OS line; `uname -s` needs no binary at all. The `^OS: ` line governs only classes where a compatible binary exists (post-gate logic, issue-report helper text, and the D-1 transition for later gate-logic versions); the wrong-version class needs no FO-side OS detection at all (doctor computes `runtime.GOOS` itself).

### 2. OS-aware install hint

- **Skill prose** (`first-officer-shared-core.md` Startup step 1, binary-absent class): the hint is OS-conditional — on Linux it leads with the documented `curl|sh` path from `install.md` line 20; on macOS it keeps the Homebrew install lead. The source-build fallback stays on both. On a non-Darwin/Linux OS, hint-and-abort with the source-build note (D-5).
- **Binary** (`internal/contract/contract.go`): `tooOldBinaryRemedy(edgeCask, goos string)` gains an OS-conditional remedy lead. `goos` is threaded explicitly (not read from `runtime.GOOS` inside the function) so tests pin the OS deterministically; the caller (`compareNamed`) supplies `runtime.GOOS`. On Linux the remedy leads with the `curl|sh` command; on macOS the Homebrew lead is unchanged. This keeps doctor's remedy text runnable on the host it prints on — doctor is the wrong-version class's remedy surface and its Mac-only lead is the same defect class this task fixes (D-4). Precedence (N6): on `runtime.GOOS == "linux"` the `curl|sh` lead **supersedes** the edge-cask branch — `spacedock@next` is a brew-only token, meaningless on Linux, so edgeCask is ignored there.

### 3. Direct install-and-resume (install path only; one bounded attempt)

- The binary-absent class changes from hint-and-abort to offer-and-resume: the FO offers to run the OS-aware install command itself, the operator approves, the FO runs it, then re-checks `${SPACEDOCK_BIN:-spacedock} --version` and resumes startup if line 1 parses to a compatible version.
- **Convergence mechanism** (cycle-1 spike): `install.sh` prints `install.sh: installed spacedock <version> to <dir>/spacedock` to stderr and warns when the dir is not on PATH. After a `curl|sh` install, the FO sets `SPACEDOCK_BIN` to the installed path (parsed from stderr, or probed at `$HOME/.local/bin/spacedock`) and re-checks `--version`. The wrong-version/brew-upgrade convergence discussion moves to the upgrade task.
- **Loop bound (design decision D-3 — the staff review's B1 demand for a falsifiable fallback on a boot path; NF2 re-choice):** one install attempt per session, enforced by a **durable sentinel FILE, not prose alone and not a shell env variable**: the FO prose checks `test -f <sentinel>` before offering to run an install; if present, it goes straight to the hint-and-abort fallback. The prose creates the sentinel (via `write`/`touch`) BEFORE the install runs — create-before-run, so an attempt that dies mid-run still counts as the one attempt (creating it only after a successful run would destroy once-only-ness). Sentinel path: `${TMPDIR:-/tmp}/spacedock-install-attempted-<key>`, where `<key>` is the host runtime's session-identity env var where one is exposed — `CLAUDE_CODE_SESSION_ID` (claude) or `CODEX_THREAD_ID` (codex), the identity column of the runtimehost marker table. Pi exposes NO session-identity env var (`internal/runtimehost/runtimehost.go:23-24` — both pi rows carry an empty identity), so on pi `<key>` falls back to a short hash of the gate's working directory: the sentinel is then **project-scoped and tmp-durable** — a failed install in session N suppresses the install offer in session N+1 of the same project until tmp cleanup. That over-reach is admitted explicitly and is still a safe degradation mode: the sentinel-blocked path is the same hint-and-abort carrying the exact OS-aware command, so the human is never left without a working remedy. The sentinel never lands in the project tree or state checkout (where a stray commit could sweep it). Why a file and not the cycle-2 `SPACEDOCK_INSTALL_ATTEMPTED` env marker (NF2): env durability across tool calls is undocumented on EVERY supported runtime (nothing in `pi-first-officer-runtime.md` or `docs/runtime-support.md` documents cross-call env persistence — the repo's own `${SPACEDOCK_BIN:-spacedock}` inline-expansion launcher pattern exists precisely because env does NOT persist), so an env marker silently degrades to the prose-only control B1 rejected — and on fresh-shell runtimes is *less* reliable than prose discipline, because an LLM reliably remembers its own context while a fresh shell reliably forgets the export. The sentinel file is the very oracle the cycle-2 spike used for verification: `test -f` IS durable across tool calls. If the post-install `--version` re-check still fails, the FO falls back to hint-and-abort with the exact OS-aware command for the human to run manually — no second install attempt, no proceed-without-recheck. The fallback message (both at first failure and on any sentinel-blocked re-entry) MUST tell the human the sentinel exists, print its path, and tell them to remove it (`rm <sentinel-path>`) to re-enable the install offer. D-2 states what the next session does instead.

### 4. Sandbox detection (session-env layer; the `--version` line corroborates only — NF1 redesign)

- **Axis correction (NF1):** sandbox state is a property of the RUNNING SESSION, not of the binary — runtime-session detection and binary presence are different axes (the corrected N2/N9 coupling claim). The install-offer class (binary absent) has no binary, so `--version` fails with NO output and no `^Sandbox:` line can ever exist in that class. A `--version`-parse-based detection mechanism therefore cannot fire in exactly the class it protects: the three-way parse falls to "no such line → not-inside" and the FO would offer the install INSIDE the sandbox — the silent no-op AC-3 exists to prevent. Stated plainly: `--version`-based sandbox detection can never exist in the install-offer class.
- **Primary detection: session-env markers, checked in every gate class.** `safehouse.Inside` is itself pure env detection (`internal/safehouse/state.go:29-36`): a registry match on env var VALUE (`APP_SANDBOX_CONTAINER_ID=agent-safehouse` today; matching is on the value, not mere presence, because the variable is a generic macOS app-sandbox signal other containers also set). The FO prose performs the SAME check directly via Bash — `[ "${APP_SANDBOX_CONTAINER_ID:-}" = "agent-safehouse" ]`, extended to any further registry rows — so detection works identically with NO binary present (nothing to execute) and in binary-present classes. The `insideRegistry` table is the source of truth for WHICH markers to check; the skill prose names the current rows and the skill smoke test asserts the prose registry matches `insideRegistry`.
- **Corroboration only (binary-present classes):** when `--version` output exists, the `^Sandbox: ` line (prefix-anchored, not a loose `inside` substring — N9 parse hygiene) corroborates the env verdict. It is never the primary source and is never consulted in the binary-absent class, where it cannot exist. Disagreement (env marker set, output absent or contradicting) resolves to the env check and is named in the human-facing message.
- If inside a sandbox, the FO does NOT offer to run the install (it would no-op for the host). It tells the human plainly: "You're running inside a sandbox. Run this command yourself, outside the sandbox: `<exact OS-aware install command>`."
- **AC-3 test-type disposition (N8, updated per NF1):** the sandbox-aware message stays skill prose (no new binary seam — the OS-aware command text already exists binary-side in the remedy, and a new `doctor --install-hint` arm would duplicate it). The test plan therefore uses a behavior fixture and the plan **names the weakness explicitly**: LLM prose is non-deterministic, so the fixture asserts semantic content (message carries the exact command token and the "outside the sandbox" instruction; install command is not executed) rather than byte-for-byte wording.

### Design decisions (explicit, per staff review)

- **D-1 (OS line + OS-source, NF4 restated):** `OS: <goos>/<goarch>` as `--version` line 2 in both output shapes. Its value (corrected per NF4): user issue reports carry the platform; the skills-too-old class gets a machine-readable platform line; and later gate-logic versions can read OS from the line once a compatible binary exists (the D-1 transition). It does NOT serve this task's own hint selection — the binary-absent class has no binary (see the OS-source caveat in §1) and the wrong-version class lets doctor compute `runtime.GOOS` itself. The FO's `uname -s` path is permanently load-bearing for the binary-absent class, not a transitional or circular fallback. Retargets the pinned one-line human output — an intended observable-semantics change, listed below.
- **D-2 (SPACEDOCK_BIN override × launcher invariant; durability — B2):** the post-install `SPACEDOCK_BIN` override is **in-process, launching-session scope** and happens **at the gate** — it is the gate's launcher resolution ("resolve ONE launcher at the version gate; use THAT launcher for every later helper call"), not a mid-session drift onto a bare `spacedock`, so it satisfies the shared-core invariant rather than colliding with it. If `SPACEDOCK_BIN` was set-but-stale (the cause of the gate failure), the override replaces it for this session and the fallback message says so. **Durability model: session-only.** The override is NOT persisted to a shell profile in this task. The **next session's** bootstrap re-runs the gate: if the install dir is on PATH, bare `spacedock` resolves and the session converges with no override; if not, the gate fails again and the fallback hint names the exact installed path and tells the human to add the dir to PATH (or launch with `SPACEDOCK_BIN=<path>`). **Surfaced captain-decision item:** whether spacedock should ever *persist* the override (write the install dir to the user's shell profile, or record it in config) is left to the captain/gate — the recommended default is no auto-profile-writes.
- **D-3 (loop-bound guardrail — sentinel FILE; choice stated per NF2):** one install attempt per session, enforced by a durable sentinel file at `${TMPDIR:-/tmp}/spacedock-install-attempted-<key>`. **Key source:** `<key>` = the host runtime's session-identity env var where one is exposed — `CLAUDE_CODE_SESSION_ID` (claude) or `CODEX_THREAD_ID` (codex) per the runtimehost marker table (`internal/runtimehost/runtimehost.go`; NOT anything `--version` reports — in the binary-absent class `--version` produces nothing to read). Pi exposes no session-identity env var (`runtimehost.go:23-24`), so there `<key>` is a short hash of the gate's working directory: the sentinel is then **project-scoped and tmp-durable** — a failed install in session N suppresses the install offer in session N+1 of the same project until tmp cleanup. That over-reach is admitted explicitly (an earlier draft's "does not guard the next session" was wrong on pi) and is a safe degradation: the blocked path is hint-and-abort with the exact OS-aware command, and the fallback message MUST tell the human the sentinel exists, print its path, and direct them to `rm` it to re-enable the install offer. **Create-before-run:** the prose creates the sentinel BEFORE the install runs, so an attempt that dies mid-run still counts as the one attempt. The prose checks `test -f` before offering an install. **The choice:** a sentinel file, NOT the cycle-2 `SPACEDOCK_INSTALL_ATTEMPTED` env marker — env durability across tool calls is undocumented on every supported runtime, so the env marker would present as a "binary-visible guardrail" while actually being prose discipline in env clothing (on fresh-shell runtimes, *less* reliable than prose discipline: an LLM reliably remembers its own context; a fresh shell reliably forgets the export). **NF1/NF2 coherence note:** the env rejection here is consistent with §4's env-based sandbox detection because the two live on different layers — sandbox env markers (e.g. `APP_SANDBOX_CONTAINER_ID`) are injected by the LAUNCHER into the session process and inherited by every Bash tool-call shell, so an env-equality read in any later tool call still sees them; a session-internal `export` (the rejected NF2 env form) is set by one tool-call shell and dies with it. The sentinel file is the mechanism the cycle-2 spike itself used as its verification oracle — `test -f` is shell-visible and durable across tool calls on every runtime. The dedicated failure-fallback behavior fixture remains (it asserts on the sentinel), so B1's bar is met by a real guardrail plus the fixture, not by a fallback-fixture-only claim. Prose-only loop control on a boot path is still judged insufficient.
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
2. **Failure-fallback fixture** (answers B1/N2; NF2 sentinel): a captive install script succeeds but produces an incompatible (or no) binary; assert the FO runs the install exactly once (the sentinel file is created, and a second simulated gate entry observes it via `test -f` and does not re-run the install), re-checks `--version`, observes failure, and emits hint-and-abort naming the exact OS-aware command — and does NOT proceed past the gate.
The independent baseline: the current behavior requires the human to copy-paste a command and restart (1+ human interventions); the new behavior requires 0 on success and exactly 1 bounded attempt on failure. These fixtures fail if the FO reverts to print-and-exit, if it attempts install more than once (guardrail absent or ignored), or if it proceeds without re-checking `--version`. Named weakness (NF5 — the same caveat AC-3 names, applied symmetrically): both fixtures are prose-driven like AC-3's — LLM gate prose is non-deterministic — so they assert observable shell side effects and semantic content (install ran once, sentinel present, `--version` re-checked, fallback message carries the exact command token), not byte-for-byte prose.

**AC-3 — When the FO is running inside a sandbox, it detects the sandbox from the SESSION'S ENVIRONMENT MARKERS (the `safehouse.Inside` registry check — env var/value match such as `APP_SANDBOX_CONTAINER_ID=agent-safehouse`, `internal/safehouse/state.go:29-36` — performed by the FO prose via Bash, in every gate class including the binary-absent one) and tells the human to run the exact install command themselves outside the sandbox, naming the exact command. It does NOT attempt the install inside the sandbox.** In binary-present classes the `--version` `^Sandbox: ` line corroborates the env verdict; detection never depends on `--version` output, which cannot exist in the install-offer class (NF1).
Verified by: A behavior fixture that drives the version gate in the REAL install-offer state — a captive ABSENT-or-failing binary producing NO `--version` output (the state a truly absent binary produces, not a synthetic binary that prints a `Sandbox:` line while failing the line-1 gate parse — no real spacedock binary produces that) — plus the sandbox env marker set (`APP_SANDBOX_CONTAINER_ID=agent-safehouse`). Assert: detection fires from the env marker with NO `--version` output present; the human-facing message carries the exact OS-aware install command token and the "outside the sandbox" instruction; and the install command is not executed (no captive-install side effect observed). A second fixture arm drives the binary-present wrong-version class with the marker set and asserts the `^Sandbox: ` corroboration agrees with (never overrides) the env verdict. Named weakness (N8 disposition in Proposed approach §4): the message is skill prose, so the assertion is on semantic content, not byte-for-byte prose. The independent baseline: a sandboxed install lands where the host cannot see it (silent no-op) — the fixture fails if the FO runs the install inside the sandbox, omits the exact command / "outside the sandbox" instruction, or requires `--version` output to detect the sandbox (unreachable in this class).

**AC-4 — `spacedock --version` reports the host OS and arch (`OS: <goos>/<goarch>`) in both output shapes (in-session and outside-runtime human), so user issue reports carry the platform, the skills-too-old class gets a machine-readable platform line, and later gate-logic versions can read OS from `--version` once a compatible binary exists (the D-1 transition).** Its value is NOT this task's hint selection (NF4 correction): the binary-absent class has no binary and reads OS via the permanently load-bearing `uname -s` path, and the wrong-version class needs no FO-side OS detection (doctor computes `runtime.GOOS` itself).
Verified by: A Go unit test on `printVersion` asserting the `OS: ` line appears as line 2 with the correct `runtime.GOOS`/`runtime.GOARCH` values in both the one-line-human case and the in-session case; the retargeted `internal/cli/version_session_test.go` pins both shapes. The independent baseline: today's `--version` output contains no OS token at all (verified — `spacedock %s`, `Runtime:`, `Sandbox:`, `contract 3` only), so dropping the OS line fails the test.

## Test plan

- **AC-1 (OS-aware hint) — Go unit test**: drive `tooOldBinaryRemedy(edgeCask, goos)` with pinned `goos` values. `"linux"`: contains the `curl|sh` command from `install.md` line 20 and no `brew` content (edge token superseded); `"darwin"`: keeps the Homebrew lead and the edge-cask byte-for-byte block. ~30-50 lines, in `internal/contract/version_message_test.go` or a new `internal/contract/remedy_os_test.go`.
- **Existing-test retarget (N3, mandatory)**: `internal/contract/version_message_test.go` is **modified**: `TestTooOldBinaryRemedyLeadsWithBrew` (currently asserts `brew upgrade spacedock` via `RunDoctor`, which will produce the curl|sh lead on the ubuntu CI leg) is retargeted to assert the host-appropriate remedy branch on `runtime.GOOS` (or rerouted to direct `tooOldBinaryRemedy` calls with pinned `goos`); `TestTooOldBinaryRemedyEdgeChannel` (byte-for-byte pinned block, direct calls) passes an explicit `goos = "darwin"` (or the new signature) so its pinned block is GOOS-stable. This is what keeps `.github/workflows/runtime-live-e2e.yml`'s ubuntu leg green.
- **AC-2 (direct install-and-resume) — two behavior fixtures**: (1) convergence fixture: captive install places a compatible binary at a known path; assert install ran once, `SPACEDOCK_BIN` repointed, `--version` re-checked, gate passed. (2) **failure-fallback fixture** (B1; NF2 sentinel): captive install "succeeds" but yields an incompatible binary; assert exactly one install run (sentinel file observed via `test -f`; a second simulated gate entry runs zero installs), `--version` re-check observed failing, hint-and-abort emitted with the exact command, gate not passed. Each ~40-60 lines; moderate cost (fixture harness simulating the FO gate flow). Named weakness (NF5): these fixtures are as prose-driven as AC-3's — LLM prose is non-deterministic, so assertions are on observable shell side effects and semantic content, not byte-for-byte wording.
- **AC-3 (sandbox detection) — behavior fixture** (explicitly NOT a unit test, per the N8 disposition in Proposed approach §4; rewritten per NF1): drive the gate in the real install-offer state — a captive absent-or-failing binary producing NO `--version` output — with the sandbox env marker set (`APP_SANDBOX_CONTAINER_ID=agent-safehouse`); assert detection fires from the env marker alone, the message carries the exact command token + "outside the sandbox" instruction, and the captive install was not executed. Second arm: binary-present wrong-version class with the marker set; assert the `^Sandbox: ` corroboration agrees with the env verdict. ~25-45 lines.
- **AC-4 (OS line) — Go unit test**: `printVersion` writes `OS: <goos>/<goarch>` as line 2 in the outside-runtime case and the in-session case; `internal/cli/version_session_test.go` retargeted (its `outside-every-runtime-is-one-line` case becomes a two-line case). ~15-25 lines total.
- **No-drift check (widened per N4, extraction defined per N5)**: a repo-relative test asserts (a) the `curl|sh` token in the shared-core hint equals `install.md`'s fenced command at line 20 **after stripping the markdown fence's 4-space indentation** (extraction: the single line inside the `=== "Binary (macOS / Linux)"` fence matching `^curl `, trimmed), and (b) the shared-core brew **install** hint refers to the same tap and formula as `install.md`'s Homebrew tab (`spacedock-dev/homebrew-tap` + `spacedock`; the two-line `brew tap` + `brew install` form and the one-line full-token form are checked for token equality, not formatting equality). The brew **upgrade** form in `contract.go` is excluded — the upgrade task owns it (Out of scope). ~10-15 lines.
- **Skill smoke test**: per repo convention, skill text changes to `first-officer-shared-core.md` get a smoke test assert for the new gate text (OS-aware hint shape, sentinel path + `test -f` check, sandbox env-marker registry rows matching the `safehouse.Inside` registry as FULL name+value rows — each env name AND its expected value, since `state.go` matches on `wantValue` (`{env: "APP_SANDBOX_CONTAINER_ID", wantValue: "agent-safehouse"}`); a names-only assertion would stay green through a `wantValue` change while the prose check silently never fires. The test reads the `insideRegistry` table via source-grep of `internal/safehouse/state.go` — the source-grep form is chosen, so no new exported accessor/API is added — sandbox instruction, amended launcher-invariant sentence). Follow the existing skill smoke-test fixture.

## Expected surface

Files:
1. `internal/cli/cli.go` — `printVersion` gains the `OS: <goos>/<goarch>` line 2 in both output shapes (move it above the outside-runtime early return). `printVersion`'s doc comment (`cli.go:820-836`) and the inline comment carrying the "One line, nothing else" wording (`cli.go:842`) must be rewritten for the two-line outside shape (NF3, stated explicitly). ~6-10 lines added/changed.
2. `internal/cli/version_session_test.go` — **modified** (required): the pinned `outside-every-runtime-is-one-line` case becomes two lines; in-session cases gain the OS line; a GOOS/GOARCH assertion is added. The test file's organizing-rule comment (`:40-43`, "the outside case is ONE line — no Runtime line, no Sandbox line…") must be rewritten to the new two-line rule (NF3, stated explicitly). ~12-22 lines changed.
3. `internal/contract/contract.go` — `tooOldBinaryRemedy` gains a `goos string` parameter and an OS-conditional lead: Linux leads with `curl|sh` (edgeCask superseded), macOS keeps Homebrew. `compareNamed` passes `runtime.GOOS`. ~15-25 lines added/changed.
4. `internal/contract/version_message_test.go` — **modified** (required, N3): `TestTooOldBinaryRemedyLeadsWithBrew` and `TestTooOldBinaryRemedyEdgeChannel` retargeted with pinned/host-threaded GOOS so the ubuntu CI leg stays green; AC-1 Linux/darwin assertions may land here. ~20-40 lines changed.
5. `internal/contract/doctor.go` — only if the `goos` thread ripples past `compareNamed`; preferred shape keeps `RunDoctor`'s signature unchanged (read `runtime.GOOS` at `compareNamed`). ~0-10 lines.
6. `internal/cli/init.go`, `internal/cli/frontdoor.go` — unchanged if `RunDoctor` keeps its signature. ~0 lines.
7. `skills/first-officer/references/first-officer-shared-core.md` — Startup step 1 rewrite of the binary-absent class: OS-aware install hint (`uname -s` as the permanently load-bearing OS source in this class — NF4; prefer the `--version` `^OS: ` line only where a binary exists), direct install-and-resume offer with the **sentinel-file** guard (`test -f` on `${TMPDIR:-/tmp}/spacedock-install-attempted-<key>`) and one-attempt bound, `SPACEDOCK_BIN` session-scoped convergence, sandbox branch (**session-env marker check mirroring the `safehouse.Inside` registry** — NF1 — with the `--version` `^Sandbox:` line as corroboration only in binary-present classes; no-install-inside message). **B2 residual nit:** the rewrite also AMENDS the launcher-invariant sentence ("resolve ONE launcher at the version gate") to bless the gate-internal post-install resolution explicitly — a gate that fails the binary-absent class and then succeeds via the approved install performs its ONE launcher resolution at that point, still inside Startup step 1 before any helper call, and THAT resolution governs every later helper call. Without this amendment, D-2's claim that the override SATISFIES the invariant contradicts the prose this task ships. ~45-65 lines (substantial rewrite of the abort classes plus the invariant-sentence amendment).
8. Skill smoke test fixture (name per existing convention) — asserts the new gate text shape. ~15-30 lines.
9. No-drift test (repo-relative; location alongside the skill smoke fixture or `internal/contract`) — curl|sh + brew-install token checks against `install.md`. ~10-15 lines.
10. `docs/site/reference/command-reference.md` — **modified** (required, NF3): the `--version` section pins both output shapes verbatim (lines 5-16: "Outside any session it prints **one line**…" plus the four-line in-session block), and the doc ALSO pins the shapes at the line-3 intro paragraph (`(the binary version, and — inside an agent session — that session's runtime and sandbox state)`) and at line 26 ("the output is the single version line shown above"); the design changes both shapes, so all of these pins change and the doc and its diff belong in the surface per the stage-def doc-diff obligation. The concrete doc diff is recorded below, covering lines 3-26. ~8-14 lines changed (the line 3 and 26 additions are two sentences within that estimate).

### Documentation diff (NF3 — stage-def doc-diff obligation; against `docs/site/reference/command-reference.md` lines 3-26: the line-3 intro pin, the verbatim line 5-16 blocks, and the line-26 outside-shape pin)

Before (line 3, intro paragraph):

    The `spacedock` binary groups its subcommands into Launch, Setup, and Workflow, plus a top-level `spacedock --version` (the binary version, and — inside an agent session — that session's runtime and sandbox state). For the exact flags of any command, run `spacedock <command> --help`, the always-current source of truth; `spacedock` with no arguments prints the grouped help.

After (line 3, intro paragraph — both `--version` shapes gain the `OS:` line 2):

    The `spacedock` binary groups its subcommands into Launch, Setup, and Workflow, plus a top-level `spacedock --version` (the binary version and the host OS/arch, and — inside an agent session — that session's runtime and sandbox state). For the exact flags of any command, run `spacedock <command> --help`, the always-current source of truth; `spacedock` with no arguments prints the grouped help.

Before (lines 7-16, the `--version` section's verbatim shape pins):

    `spacedock --version` reports the binary version, and — when it is running inside an agent session — that session's runtime and sandbox state. Outside any session it prints one line:

        spacedock 0.26.0

    Inside a session it also names the runtime it detected, the marker that proved it, which session this is, and whether this process is running inside a sandbox:

        spacedock 0.26.0
        Runtime: claude (CLAUDECODE, session afd74765)
        Sandbox: inside (agent-safehouse)
        contract 3

After:

    `spacedock --version` reports the binary version and the host OS/arch, and — when it is running inside an agent session — that session's runtime and sandbox state. Outside any session it prints two lines:

        spacedock 0.26.0
        OS: darwin/arm64

    Inside a session it also names the host OS/arch, the runtime it detected, the marker that proved it, which session this is, and whether this process is running inside a sandbox:

        spacedock 0.26.0
        OS: darwin/arm64
        Runtime: claude (CLAUDECODE, session afd74765)
        Sandbox: inside (agent-safehouse)
        contract 3

Before (line 26, the outside-shape pin — the stale one-line pin must not survive):

    Being outside every runtime is a normal state, not a fault — it means a human at a terminal. There is no `Runtime:` line at all in that case, because there is no session to report: the output is the single version line shown above.

After (line 26 — no `Runtime:` line remains true, but the outside output is now two lines):

    Being outside every runtime is a normal state, not a fault — it means a human at a terminal. There is no `Runtime:` line at all in that case, because there is no session to report: the output is the two lines shown above (the version line plus the `OS:` line).

LOC tolerance: ±20% per file.

Observable semantics that may change:
- `spacedock --version` gains an `OS: <goos>/<goarch>` line 2 in **both** shapes — including the outside-runtime human output previously pinned to a single line (intended, D-1).
- `spacedock doctor`'s too-old-binary remedy leads with the `curl|sh` install command on Linux instead of `brew upgrade ...` (intended, D-4); edge-cask naming is unchanged on macOS and ignored on Linux.
- The FO version gate's binary-absent class offers to run the install and re-checks `--version` (abort→install-and-resume behavior), bounded by a durable sentinel file at `${TMPDIR:-/tmp}/spacedock-install-attempted-<key>` (new on-disk semantics: session/project-scoped tmp file, prose-created on install, prose-checked via `test -f`).
- The FO version gate detects sandbox state from session-env markers (the `safehouse.Inside` registry env checks via Bash — NF1), corroborated by the `^Sandbox:` line only when `--version` output exists, and reroutes to a human-run-outside-sandbox instruction (stderr/message content).
- The FO's session-scoped `SPACEDOCK_BIN` override after install (launcher resolution at the gate; not persisted, D-2; the shared-core invariant sentence is amended to bless this gate-internal resolution — B2 nit).
- **Undocumented coupling named (NF6):** the install-resume convergence parses `install.sh`'s stderr contract `install.sh: installed spacedock <version> to <dir>/spacedock` (`install.sh:178`) — a coupling on an out-of-scope-unchanged script, so drift there silently breaks convergence; the `$HOME/.local/bin/spacedock` probe fallback (`install.md:23`) mitigates, and the no-drift check does not cover this stderr format.

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

**Cycle-3 addendum (NF1/NF2 — corrections and proven-mechanism record):**

- **Cycle-2 record correction (NF2):** the cycle-2 spike's guardrail claim covered only intra-process shell logic (passes 1 and 2 ran in one shell lineage); it never exercised cross-tool-call env persistence in an actual FO runtime — which is undocumented on every supported runtime — so "spiked (sentinel-verified)" overclaimed the ENV marker's delivery mechanism. The sentinel file the spike used to VERIFY the flow is what the design now adopts as the guardrail itself: `test -f` is durable across tool calls and was exercised by the spike. Recorded here so the record of record stops standing behind a claim the evidence does not support.
- **No new spike required for the cycle-3 mechanisms; proven mechanisms named (per stage-def no-spike record):** (1) the sentinel-file guardrail is the cycle-2 spike's own verification oracle, exercised end-to-end there; (2) the sandbox env-marker check (`[ "${APP_SANDBOX_CONTAINER_ID:-}" = "agent-safehouse" ]`) is the identical registry read `safehouse.Inside` performs (`internal/safehouse/state.go:29-36`) — a shell env-equality test needs no binary, no I/O, and nothing exotic; (3) that `--version` cannot carry a `Sandbox:` line in the binary-absent class requires no spike — it follows from `printVersion` being unreachable when the binary does not exist (NF1's axis correction).

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

## Stage Report: ideation (cycle 3)

- DONE: Sandbox detection works in the binary-absent class: detection moves to session-env markers (safehouse.Inside-style) with the --version Sandox line as corroboration only, and AC-3's test is rewritten to a real absent-binary + env-marker state (reachable, falsifiable)
  Proposed approach §4 redesigned: primary detection is the `safehouse.Inside` registry env-var/value check (`APP_SANDBOX_CONTAINER_ID=agent-safehouse`, internal/safehouse/state.go:29-36) performed by FO prose via Bash in EVERY gate class; the `--version` `^Sandbox:` line is corroboration only in binary-present classes; §4 states that `--version`-based detection can never exist in the install-offer class; AC-3 rewritten — fixture drives a captive absent/failing binary (NO `--version` output) + env marker set, asserting detection fires, the message carries the exact command + "outside the sandbox" instruction, and install is not executed; N2/N9 coupling claim corrected to the two-axes statement.
- DONE: The loop-bound guardrail is a durable file sentinel (test -f) or the design honestly names prose discipline; no durable-environment-variable claim remains, and D-3 states the choice
  D-3 and Proposed approach §3 re-choose a sentinel FILE (`${TMPDIR:-/tmp}/spacedock-install-attempted-<key>`, prose-created, `test -f`-checked — the cycle-2 spike's own verification oracle, durable across tool calls); every `SPACEDOCK_INSTALL_ATTEMPTED` env-marker claim removed from §3/D-3/AC-2/test plan/Expected surface; Spike record cycle-3 addendum corrects the cycle-2 overclaim (intra-process shell logic is not cross-tool-call env durability); failure-fallback fixture retained, now asserting on the sentinel.
- DONE: command-reference.md is added to Expected surface with the --version doc-shape diff; AC-4's value statement is corrected to issue reports + skills-too-old class, and the invariant text the skill rewrite ships is amended to bless the post-install launcher resolution
  Expected surface file 10 adds `docs/site/reference/command-reference.md` with a recorded before/after doc diff for both `--version` shapes (stage-def doc-diff obligation); files 1-2 now name the `printVersion` doc comment (cli.go:820-836) and the version_session_test organizing-rule comment (:40-43) rewrites explicitly; AC-4's value is restated (issue reports + skills-too-old class + D-1 transition; NOT this task's hint selection) and §1/D-1 name `uname -s` permanently load-bearing, not circular/transitional (NF4); Expected surface file 7 amends the "resolve ONE launcher at the version gate" sentence to bless the gate-internal post-install resolution (B2 residual nit).

### Staff-review dispositions (cycle 3 — second fresh-context review, `.pi-subagents/artifacts/b3ede353_reviewer_0_output.md`)

- NF1 (blocker — sandbox detection unreachable in the class it protects): fixed. §4 redesign moves detection to the session-env layer (the same registry `safehouse.Inside` reads) in all classes; `--version` `^Sandbox:` demoted to corroboration-only; the design states `--version`-based detection can never exist in the install-offer class; AC-3 + test plan rewritten to the real absent-binary + env-marker state (the synthetic `Sandbox:`-printing-while-failing binary is explicitly excluded).
- NF2 (material, folds into B1 — env-marker durability undocumented): fixed. Guardrail re-chosen as a sentinel file via `test -f` (the durable mechanism the spike already validated as its oracle); D-3 states the choice and why env was rejected; no durable-env claim remains anywhere in the body; the failure-fallback fixture stays, asserting on the sentinel, so B1's bar is a real guardrail + fixture, not a fallback-fixture-only claim.
- NF3 (moderate — doc surface incomplete): fixed. `docs/site/reference/command-reference.md` is Expected surface file 10 with the concrete two-shape doc diff recorded in the body; the `printVersion` doc comment (cli.go:820-836) and the `version_session_test.go` organizing-rule comment (:40-43) rewrites are named explicitly in files 1-2.
- NF4 (note — OS-line value overstated): fixed by restating. AC-4's value is issue reports + the skills-too-old class + the D-1 transition; §1's OS-source caveat and D-1 state `uname -s` is permanently load-bearing for the binary-absent class, not circular and not transitional.
- NF5 (note — AC-2 lacks the nondeterminism caveat AC-3 names): fixed. AC-2 and its test-plan entry name the prose-non-determinism caveat; assertions are on observable shell side effects and semantic content.
- NF6 (note — third undocumented coupling): fixed. The observable-semantics list names the `install.sh:178` stderr-contract parsing (`install.sh: installed spacedock <version> to <dir>/spacedock`) as an undocumented coupling, with the `$HOME/.local/bin/spacedock` probe fallback (`install.md:23`) as mitigation.
- B2 residual nit: fixed. Expected surface file 7 requires the shared-core rewrite to amend the "resolve ONE launcher at the version gate" sentence to bless the gate-internal post-install resolution, so D-2's satisfaction claim matches the shipped prose.
- B1 (re-examined with NF2): stays closed-in-design via the sentinel-file guardrail + failure-fallback fixture; the env-marker mechanism is removed, not repaired.
- N2 (reopened-in-fact by NF1): re-closed. The disposition's coupling assumption is corrected along the two-axes statement (runtime-session detection ≠ binary presence); sandbox conservatism (human runs the command outside; always works) unchanged.
- N8: remains closed; the named prose-non-determinism weakness now applies symmetrically to AC-2 (NF5).
- N9: parse hygiene (`^Sandbox:` anchor) retained for the corroboration role; the coupling statement ("the gate always runs in-session, so the line is present") is corrected — in-session is necessary but binary presence is the axis the install-offer class fails.

### Summary

Cycle-3 ideation applies the FO's six authorized fixes against the second fresh-context staff review. The blocker (sandbox detection unreachable in the install-offer class) is resolved by moving detection to session-env markers (`safehouse.Inside`-style registry checks) with the `--version` line as corroboration only, and AC-3's fixture now drives a real absent-binary + env-marker state. The loop-bound guardrail becomes a durable sentinel file (the cycle-2 spike's own oracle); `command-reference.md` joins the expected surface with a recorded two-shape doc diff; AC-4's value statement, AC-2's caveat naming, the `install.sh:178` stderr coupling, and the launcher-invariant amendment are all corrected in place. This is not a new gate decision: the captain's two binding annotations (install-only scope; OS line in `--version`) still govern and remain applied.

## Stage Report: ideation (cycle 4)

- DONE: The recorded command-reference.md doc diff covers the shape pins at lines 3 and 26, not just 5-16 (no stale one-line pin survives the edited doc)
  Documentation diff section extended with Before/After blocks for the line-3 intro paragraph (now names host OS/arch) and the line-26 "the output is the single version line shown above" pin (now "the two lines shown above (the version line plus the `OS:` line)"); heading and Expected surface file 10 updated; ~8-14 LOC estimate retained (two sentences).
- DONE: D-3 is internally consistent on pi: the tmp-durable/project-scoped sentinel fallback is admitted, the fallback message names the sentinel path and its rm recovery, creation is specified as before-install-runs, and identity env vars are named (CLAUDE_CODE_SESSION_ID / CODEX_THREAD_ID)
  §3 loop-bound bullet and D-3 rewritten together: pi exposes no session-identity env var (verified `internal/runtimehost/runtimehost.go:23-24` — both pi rows carry an empty identity column), so `<key>` falls back to a gate-cwd hash; the session-N-failure suppresses session-N+1-offer over-reach is admitted as a safe degradation (hint-and-abort with the exact command); create-before-run stated; fallback message MUST print the sentinel path and the `rm` recovery.
- DONE: The NF1/NF2 env distinction (launcher-injected inherited markers vs session-internal exports) is stated; the registry smoke test asserts full name+value rows; the cli.go/version_session_test.go line citations are corrected
  D-3 coherence note added (launcher-injected `APP_SANDBOX_CONTAINER_ID` inherited by every tool-call shell vs session-internal exports dying with the shell); test-plan smoke test now asserts full `insideRegistry` rows (env name AND `wantValue`) via source-grep of `internal/safehouse/state.go` (no new exported API); citations corrected to `cli.go:820-836` / `cli.go:842` / `version_session_test.go:40-43` / `:64-67`, verified with `grep -n` against the current sources.

### Dispositions (cycle 4 — third fresh-context review, `.pi-subagents/artifacts/a1a7448b_reviewer_0_output.md`)

1. NF3 completion (doc-diff obligation): fixed — the recorded diff now covers lines 3 and 26 plus 5-16; no stale one-line pin survives.
2. D-3 sentinel scoping contradiction: fixed — (a) pi's project-scoped/tmp-durable fallback admitted, over-reach named with its safe degradation mode; (b) fallback message names the sentinel path + `rm` recovery; (c) create-before-run specified; (d) identity env vars named (`CLAUDE_CODE_SESSION_ID` / `CODEX_THREAD_ID`), the "`--version` reports" identity citation removed.
3. Env-distinction coherence note: fixed — stated in D-3 (launcher-injected inherited markers vs session-internal exports).
4. Registry smoke-test spec: fixed — full-row name+value assertion; source-grep form chosen (no new exported accessor).
5. Citation corrections: fixed — `cli.go:820-836`, `cli.go:842`, `version_session_test.go:40-43`, `:64-67`; all four verified against current source.

### Summary

Cycle-4 applies the FO's five verbatim fix-in-place edits from the third fresh-context staff review: the doc diff is extended to the line-3 and line-26 shape pins; D-3's sentinel scoping is made pi-consistent (tmp-durable/project-scoped fallback admitted, `rm` recovery in the fallback message, create-before-run, named identity env vars); the launcher-injected vs session-internal env distinction is stated; the registry smoke test is specified as full name+value rows via source-grep; and four stale line citations are corrected against the current sources. No redesign: all prior sections stand, and the captain's two binding annotations (install-only scope; OS line in `--version`) still govern.
