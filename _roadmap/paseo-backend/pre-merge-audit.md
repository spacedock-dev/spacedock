# Paseo backend pre-merge audit

Date: 2026-11-23. Scope: sprint `paseo-backend`, release train `0.29.x`.

## Executive Verdict

**SHIP-BLOCKERS FOUND — 2 ship-blockers, both owned by `x9m` / PR #895.**

The production change is small and justified. It adds an explicit transport choice without a second worker controller or an unused production simulator. The recovery contract contains conflicting decisions. Its live proof also supplies the adoption procedure through test prompts and releases the worker before observing adoption.

The prose needs consolidation, especially the duplicated recovery authority. The code needs a stronger adoption observation, not more recovery machinery. No independent ship-blocker was found in `936` / PR #893.

This audit covers commit `01d66ed160eb6d44b60c7c34c58534094d063ae6`. All code citations refer to that tree. The code checkout is read-only. Section 3 separates static counterexamples from executed checks.

## Assembled Stack Summary

| Boundary | Content | Relation to sprint DoD |
|---|---|---|
| `origin/main..e1b971b32`, `936` / PR #893 | Explicit `--host paseo`, host-neutral assignment, deferred FO binding, ensign adapter, isolated live harness, optional CI | Supplies fresh dispatch, installed skill discovery, worktree adoption, Git completion checks, and worker archival. |
| `e1b971b32..01d66ed16`, `x9m` / PR #895 | Durable worker record, restart lookup, recovery matrix, write authority, status fixture, adoption smoke | Supplies restart adoption and single-retry recovery. The two blockers concern this contract and its proof. |

The audit baseline `origin/main` resolves to `439abf789d6cdddd846bed8753a42ea67e18a68f`. The merge base is `37dd5c876d152aed273a5a755650e78f39e97c0d`.

The requested two-dot comparison contains **23 files, 1,493 insertions, and 29 deletions**. Three files reflect review-consent work present only on main. They are not sprint-authored changes. The actual stack from the merge base contains **20 files, 1,483 insertions, and 17 deletions**. `936` contributes 17 files with 1,117 insertions and 14 deletions. `x9m` contributes five files with 377 insertions and 14 deletions. Two files overlap.

A normal merge preserves the main-only consent changes. A tree replacement or an indiscriminate two-dot patch can revert them. This audit does not classify those apparent reversions as sprint scope creep.

The sprint DoD requires a real local Paseo FO and worker, committed code and split-root reports, notifications-off observation, archival, and restart adoption. Remote workers, leases, worker reuse, and child-to-parent messaging remain outside scope. The optional Paseo CI lane does not replace the required existing host lanes.

## Section 1 — Complete Prose Change Inventory & Necessity Audit

Each row covers a changed section or an independently meaningful clause. **Necessary** means the text serves the declared DoD. **Tighten** means the obligation remains useful but its wording or authority needs correction. **Cuttable/Redundant** identifies text that adds no necessary obligation.

### Skill entry points and ensign binding

| Location and owner | What changes and why it exists | Verdict | Exact excess or gap |
|---|---|---|---|
| `skills/first-officer/SKILL.md:38`, `936` | Selects Paseo through `PASEO_AGENT_ID` or explicit choice. Loads the transport beside the shared dispatch core. | Tighten | Necessary selection. “At first dispatch” needs the same resume trigger as `fo-dispatch-core.md:231`. State that Paseo owns all worker lifecycle bindings, including wait and roster. See Section 4. |
| `skills/ensign/SKILL.md:13`, `936` | Selects the Paseo adapter through the environment or assignment reference. | Necessary | The prompt fallback addresses F8. When the inner host marker also exists, clarify first-match precedence. No additional adapter is needed. |
| `skills/ensign/references/paseo-ensign-runtime.md:1-5`, `936` | Names the adapter and assigns shared discipline to the shared core. | Necessary | A short ownership statement prevents another full ensign contract. |
| `skills/ensign/references/paseo-ensign-runtime.md:7`, `936` | Loads shared instructions, reads the pointer, and runs fetch commands. | Cuttable/Redundant | Repeats `SKILL.md:8,13,18` and `build.go:882`. It also tells the already loaded adapter to load itself. Keep only a genuinely Paseo-specific obligation. |
| `skills/ensign/references/paseo-ensign-runtime.md:8`, `936` | Adds noninteractive SSH failure handling and durable local completion evidence. | Tighten | BatchMode is necessary operational guidance. The opening code/state commit sentence repeats shared discipline and generated assignment guidance. Split the long push-failure sentence. |
| `skills/ensign/references/paseo-ensign-runtime.md:9`, `936` | Reports a blocker through Git and ends the turn because ordinary ensigns lack daemon messaging access. | Necessary | This is the supported clarification route under F8. The FAILED report gives the FO evidence for intervention. |
| `skills/ensign/references/paseo-ensign-runtime.md:10`, `936` | Defines final-turn completion text and forbids native mailbox delivery. | Necessary | The generated completion block repeats this for the assignment recipient. That narrow boundary repetition is justified. Git remains authoritative. |

### Paseo FO reference: purpose and common invariants

All locations in the next three tables refer to `skills/first-officer/references/paseo-fo-dispatch.md`.

| Section / lines and owner | What changes and why it exists | Verdict | Exact excess or gap |
|---|---|---|---|
| Title and introduction, 1-3, `936` | Defines the deferred, host-neutral transport. | Tighten | “Five verbs” describes five groups, not five CLI verbs. Project lookup and workspace listing add commands. Use “operations.” |
| JSON exit contract, 7, `936` | Distinguishes successful stdout from error stderr. | Necessary | Essential for CLI control. No extra parser layer is warranted. |
| Stdout-only parsing, 8, `936` | Handles successful `run` stderr noise. | Tighten | The rule is essential. The full failure example and repeated rule at line 65 can reduce to one rule plus a short call-site reminder. |
| Parent identity and daemon target, 9, `936` | Preserves caller labels and the session daemon. | Necessary | Prevents an accidental switch to the default daemon. F3 remains an upstream environment limitation. |
| Least privilege, 10, `936` | Denies ordinary workers the FO daemon grant. | Necessary | Required because the available grant exposes more than parent-only messaging. |
| Provider/model/thinking, 11, `936` | Requires explicit provider/model and separate thinking effort. | Tighten | `model=null` makes this necessary. The concrete Antigravity model example is unnecessary future drift. Generic provider/model syntax suffices. |
| Worker label, 12-13, `936` | Keeps the roster key invariant across retries. | Necessary | Separates roster lookup from attempt identity and avoids hidden retry workers. |
| Epoch label, 14, `936` | Defines cycle and attempt identity. | Necessary | A deterministic epoch supports lost-receipt lookup. The “v1 fresh” wording can become “first attempt.” |

### Paseo FO reference: pre-spawn and CLI operations

| Section / lines and owner | What changes and why it exists | Verdict | Exact excess or gap |
|---|---|---|---|
| Pre-spawn item 1, 18, both PRs | Builds/stamps the assignment. `x9m` captures the state revision and excludes stamping on retry. | Tighten | Necessary inputs, but line 162 repeats the full sequence without the retry exception. Define it once and refer to it. |
| Pre-spawn item 2, 19-22, `936` | Resolves the physical code directory from structured status. | Necessary | Avoids parsing prompt prose and preserves split-root worktree ownership. |
| Pre-spawn item 3, 23, `936` | Adopts the existing worktree instead of asking Paseo to create one. | Necessary | Prevents occupied-branch errors and mismatched worktrees. The same ownership rule need not recur in the runtime summary. |
| Pre-spawn item 4, 24, `936` | Finds a matching project, or lets Paseo register a new checkout. | Necessary | Supports an installed consumer repository without a preexisting Paseo project ID. |
| Pre-spawn item 5, 25, `936` | Stops on inaccessible sibling state. | Tighten | Correct failure boundary. F7 also blocks common `.git` metadata and workflow reads. Naming only state writes understates the known diagnostic scope. |
| Verb 1, 29-46, `936` | Lists and reuses workspaces before creation. Documents required fields and failures. | Necessary | Resolves preflight M3. The `background` and display-name fields are not consumed here and can be omitted from the field inventory. |
| Verb 2 command, 48-60, both PRs | Starts a labeled asynchronous worker with explicit provider, mode, and thinking. `x9m` uses the computed epoch. | Necessary | This is the transport call. Bracketed optional arguments make it a command template, not a directly runnable shell block. |
| Verb 2 receipt/error handling, 63-66, both PRs | Captures the handle. `x9m` persists the running record before waiting. | Tighten | Persisting the receipt is essential. The exhaustive symbolic error list adds little beyond reconcile-on-failure and diagnostic preservation. |
| Verb 3, 68-81, `936` | Polls with bounded waits and branches on JSON status despite exit 0. | Tighten | Timeout and idle rules are necessary. Line 80 says to hold on error, conflicting with `x9m` retry rows 5-6 at 150-151. Ship-blocker 1. |
| Verb 4, 83-93, `936` | Separates exact inspect data from lossy list data. | Tighten | `Cwd`, `CreatedAt`, and `Archived` are essential. `LastUsage` details, `UpdatedAt`, and display fields are unnecessary to the FO algorithm. Keep only consumed fields. |
| Verb 5, 95-105, `936` | Archives verified or explicitly superseded workers and protects parents. | Necessary | The cascade explanation addresses a real destructive boundary. Exact failure handling belongs in the recovery owner, not another independent decision list. |

### Paseo FO reference: durable state and recovery

| Section / lines and owner | What changes and why it exists | Verdict | Exact excess or gap |
|---|---|---|---|
| Worker record ownership and update boundary, 107-109, `x9m` | Defines one FO-owned H2 section and path-scoped writes at launch/adoption and terminal outcomes. | Tighten | Necessary write boundary. Parser function names and “no intermediate verified” preserve design history rather than reader needs. H2 and update timing suffice. |
| Worker record commit/sync and gate timing, 109, `x9m` | Commits the record before waiting and before gate preparation. Handles index contention and noninteractive SSH. | Tighten | Gate timing is essential. Inline Git mechanics duplicate the shared state-commit mechanism. A blind second commit also needs an unchanged-record no-op case. |
| Worker record YAML, 111-123, `x9m` | Persists identity, stage, epoch, directory, baseline, requested/returned runtime, and terminal state. | Necessary | Ten required top-level keys plus conditional reason are proportionate. No fixed YAML ordering or new production parser is necessary. |
| Restart trigger and identity derivation, 125-129, `x9m` | Recovers expected identity from existing task state and retry history. | Necessary | Deterministic labels make lost-receipt recovery possible. Clarify selection of the matching stage/cycle entry instead of introducing a new registry. |
| Restart step 1, 130, `x9m` | Preserves holds and resumes an already archived valid completion. | Necessary | Handles a restart between archival and stage advancement. Row 10 needs a narrower scope so it cannot override this recovery. |
| Restart step 2, 131, `x9m` | Inspects the recorded worker instead of spawning another. | Necessary | This is the principal restart path. |
| Restart step 3.1-3.2, 132-134, `x9m` | Searches epoch labels, then worker labels, including archived records only as fallback. | Tighten | Necessary receipt recovery. Fallback currently keys on raw empty results rather than the locally filtered result. Section 4 records the collision risk. |
| Restart step 3.3, 135, `x9m` | Filters candidates through exact inspect fields and physical directory. | Necessary | Resolves preflight M2. A worker-label fallback also needs explicit treatment of a different epoch. |
| Restart step 3.4, 136, `x9m` | Adopts one match, holds multiple active matches, or spawns after no match. | Tighten | Necessary decisions. Multiple closed candidates and recovery of a missing retry baseline remain underspecified. Avoid expanding this into an unused Go simulator. |
| Completion predicate, 140, `x9m` | Requires new committed report evidence, clean paths, expected branch, and current identity. | Tighten | Essential Git gate. A nonempty entity revision range alone also sees FO record commits. Explicitly compare the report change and required deliverable, not only entity activity. |
| Retry ledger, 142, `x9m` | Persists one retry independently from feedback cycles. | Necessary | Required for restart-safe retry limits. The legacy cycle-1 fallback is harmless but has no demonstrated Paseo migration need. |
| Matrix row 1, 146, `x9m` | Treats timeout as progress without spending a retry. | Necessary | F9 and slow workers make this distinction necessary. |
| Matrix row 2, 147, `x9m` | Holds permission-blocked work without killing the worker. | Necessary | Requires precedence over the generic running-worker row. |
| Matrix row 3, 148, `x9m` | Adopts a running worker and waits. | Necessary | Central sprint value. An unchanged record need not create another commit. |
| Matrix row 4, 149, `x9m` | Archives and advances a completed attempt. | Tighten | It accepts `AGENT_NOT_FOUND` as input but only tolerates `AGENT_ALREADY_ARCHIVED` on archive. See ship-blocker 1. |
| Matrix row 5, 150, `x9m` | Archives a failed first attempt, records one retry, and reuses the worktree. | Necessary | The bounded recovery serves the task. It conflicts with line 80 and must become the sole decision owner. |
| Matrix row 6, 151, `x9m` | Holds after the retry fails. | Tighten | Necessary retry cap. Explicitly tolerate a missing attempt-2 agent while recording the hold. |
| Matrix row 7, 152, `x9m` | Adopts a worker after a lost launch response. | Necessary | Prevents retries from creating duplicate workers. |
| Matrix row 8, 153, `x9m` | Retries a confirmed failed launch once. | Necessary | Distinct from a lost receipt because no worker exists. |
| Matrix row 9, 154, `x9m` | Holds after the second failed launch. | Necessary | Makes the retry cap durable. |
| Matrix row 10, 155, `x9m` | Ignores stale or duplicate completion evidence. | Tighten | “Repeated wait” is too broad. A later wait can observe new completion. Scope the no-op to an already processed terminal event. Ship-blocker 1. |
| Matrix row 11, 156, `x9m` | Holds ambiguous active workers without killing unowned agents. | Necessary | The conservative decision is smaller and safer than automatic reconciliation. |
| Matrix columns and closing paragraph, 144-158, `x9m` | Restates retries, archival, unrelated-work behavior, and parent protection. | Cuttable/Redundant | Every `Unrelated` cell says `Unaffected`. State that once. Parent protection repeats line 105 and line 167. Most retry/archive cells repeat the action cell. |

### Paseo FO reference: capability bindings

| Binding / line and owner | What it adds and why | Verdict | Exact excess or gap |
|---|---|---|---|
| `worker.spawn`, 162, both PRs | Maps the shared capability to the Paseo launch sequence and durable record. | Tighten | Repeats pre-spawn, Verb 1, Verb 2, and record mechanics. Refer to those owners and retain the retry exception. |
| `addressable-worker`, 163, `936` | Declares fresh workers without ordinary reuse or messaging. | Necessary | Prevents inherited native-host reuse. |
| `async-dispatch`, 164, `936` | Maps asynchronous launch and polling. | Necessary | A compact binding is appropriate. |
| `worker-identity`, 165, both PRs | Points at the durable worker identity. | Tighten | Repeats all ten YAML keys from 112-121. A reference to the record is sufficient. |
| `completion-signal`, 166, both PRs | Connects idle wakeup to Git evidence. | Necessary | Keep the predicate reference. Repeating its internals can drift. |
| `worker.shutdown`, 167, both PRs | Connects verified archival, record update, and gate timing. | Tighten | The binding is necessary. Parent and gate rules already have owners at 105 and 109. |
| `context-budget`, 168, `936` | Declares the capability absent for fresh workers. | Necessary | Prevents unsupported context-budget probing. |
| `roster-reconcile`, 169, both PRs | Binds shared reconciliation to the restart algorithm. | Necessary | Referencing the algorithm is sufficient. The literal “11-row” count is unnecessary maintenance coupling. |

### Write authority, operator documentation, and generated CLI prose

| Location and owner | What changes and why it exists | Verdict | Exact excess or gap |
|---|---|---|---|
| `skills/first-officer/references/fo-write-core.md:10`, `x9m` | Adds retry and worker-record sections to the allowed-state classifier. | Necessary | This is the mechanical authority boundary. Product files remain blocked. |
| Same file, 34, `x9m` | Defines the FO-owned operational sections and correct split-root commit location. | Necessary | Aligns the human scope statement with the classifier. It also corrects the earlier unconditional worktree-copy wording for feedback cycles. |
| Same file, 40, `x9m` | Adds the same narrow exceptions to the body-edit prohibition. | Necessary | Without this edit, the scope statement contradicts the classifier. No broader entity-content authority is added. |
| `docs/runtime-live-ci.md:103-109`, both PRs | Documents local dispatch smoke, then expands the command to adoption smoke. | Tighten | Necessary operator entry point. The shared 15-minute package timeout can be shorter than two serial 12-minute per-test caps. |
| Same file, 111, `x9m` | Describes stop, restart adoption, advancement, and archival. | Necessary | The claim needs the stronger observation in ship-blocker 2. |
| Same file, 113, both PRs | Documents binary, worker provider, FO provider, and thinking overrides. | Necessary | Supports local probes without hardcoded model defaults. |
| Same file, 115, `936` | Describes auth, environment, and optional CI dispatch lane. | Necessary | Accurately says CI runs dispatch smoke only. It does not claim automatic adoption coverage. |
| `internal/dispatch/dispatch.go:265,277,294`, `936` | Adds Paseo to general usage, build usage, and host help. | Necessary | Keeps the CLI surface discoverable and consistent. |
| Same file, 300, `936` | Documents the Paseo bare-mode refusal. | Necessary | Matches the early runtime guard. |
| `internal/dispatch/build.go:197,210`, `936` | Adds Paseo to invalid-host and ambiguous-host diagnostics. | Necessary | Gives an actionable explicit selection without changing runtime detection. |
| Same file, 245-247, `936` | Explains and reports named-worker requirements. | Necessary | Labels are required for restart lookup. The error can use two sentences instead of a semicolon. |
| Same file, 794-796, `936` | Adds host-neutral first-action instructions and code/state scope. | Tighten | The first sentence is oversized and repeats the outer pointer. Shared commit discipline also appears in the generated assignment body. Keep the adapter load and assignment boundary. |
| Same file, 840-844, `936` | Adds the final-turn message and Git completion explanation. | Necessary | Binds the assignment to Paseo completion instead of a native mailbox. |
| Same file, 882, `936` | Adds the host-neutral installed-skill pointer prompt. | Necessary | This is the worker-facing discovery path. It must remain intact through the FO transport. |
| `.github/workflows/paseo-live-e2e.yml:1,8,13,18,70-74,104,112`, `936` | Names the optional lane, manual controls, auth failure, run step, and artifact step. | Necessary | “Pin” describes a manual npm version selection, not a new version-pinning subsystem. |
| `internal/ensigncycle/paseo_live_smoke_test.go:40-47`, `936` | Supplies fixture work and explicit dispatch instructions to the live FO. | Tighten | Fixture paths, receipt, and provider are necessary. The prompt also teaches the transport sequence. During installed-binding proof work, reduce this duplication. |
| `internal/ensigncycle/paseo_adoption_live_test.go:41,56-57,78-79`, `x9m` | Supplies worker synchronization and teaches the full record/adoption procedure. | Tighten | The fixture barrier is justified. The FO prompts duplicate the shipped behavior under test. Ship-blocker 2 covers the adoption portion. |

The inherited model-ignore diagnostic now applies to Paseo without a new string definition. `internal/dispatch/build_paseo_host_test.go:139-143` checks `model=null` and the host-specific diagnostic. This is necessary disclosure of a model choice owned by Paseo.

### Main-only comparison artifacts

| Two-dot location | Apparent change | Necessity disposition |
|---|---|---|
| `skills/review/SKILL.md:20` | Product-name/private-link wording appears to revert to the public host wording. | Cuttable/Redundant as sprint work. No sprint commit changes this file. Preserve main through normal merge. |
| `skills/review/SKILL.md:112` | Longer access/comment/expiry consent appears to become shorter wording. | Cuttable/Redundant as sprint work. Same ancestry explanation. |
| `docs/specs/extension-install-contract.md:316-320` | Consent specification appears shortened. | Cuttable/Redundant as sprint work. Same ancestry explanation. |
| `skills/integration/review_choice_live_test.go:367-372,414-419,728-735,761-768` | Four consent assertions appear to revert from product/access wording to host wording. | Cuttable/Redundant as sprint work. These test differences accompany the main-only consent commits. |

### Per-section analysis

The CLI-specific details earn their place when they prevent concrete runtime errors. These include separate streams, physical paths, explicit models, stable labels, and parent archival protection. Preflight findings M2-M5 and M7 explain these details. Removing them to meet an arbitrary line budget weakens the transport.

The record and retry ledger also earn their place. They provide restart identity and a durable retry limit without changing the status schema. The H2 boundary avoids a new parser and preserves report scans.

The excess is repeated procedure. Pre-spawn, verb descriptions, recovery rows, and capability bindings each state parts of the same lifecycle. The assembled reference is **20,316 bytes**. The shared boot core remains **23,497 / 23,500 bytes**, unchanged. Preserving that cap is correct, but moving repetition into an uncapped reference does not remove its reading cost.

The smallest correction keeps one owner for each decision. CLI sections describe inputs and observations. The recovery section decides what those observations mean. Capability bindings point to those owners. This structure also removes the contradiction in ship-blocker 1.

## Section 2 — Complete Code Behavior Change Inventory & Necessity Audit

**Necessary** means the mechanism directly supplies or proves the DoD. **Tighten** means the mechanism remains justified but needs a narrower implementation or stronger evidence. **Unjustified/Cuttable** identifies machinery without a distinct value claim.

| Component / location and owner | Major behavior and value | Verdict | Complexity, correctness, and proof assessment |
|---|---|---|---|
| `internal/dispatch/build.go:195-219`, `936` | Accepts explicit Paseo while preserving inner-host auto-detection. | Necessary | Small extension to the existing validator. Adding `PASEO_AGENT_ID` to runtimehost can create false ambiguity for nested hosts. |
| Same file, 245-247 | Rejects Paseo bare mode before stamping or artifact creation. | Necessary | Small guard protects named roster attribution. No separate builder is needed. |
| Same file, 600-611, 837-845 | Emits Paseo completion guidance for named dispatch. | Necessary | Uses the existing assignment assembly. No transport API enters production Go. |
| Same file, 619-637 | Prefixes pointer filenames with a sanitized Paseo agent ID. | Necessary | Reuses the existing token sanitizer. Avoids concurrent FO file collisions when the marker exists. Marker stripping remains F3. |
| Same file, 790-796, 880-882 | Emits host-neutral skill loading and pointer text. | Necessary | Uses the established envelope and split-root paths. Wording needs tightening, not another abstraction. |
| Same file, 682-710 and existing model resolution | Paseo inherits schema v2, named envelope, null model, and omission of native-only fields. | Necessary | The branch restrictions already express this shape. There is no extra provider/model registry. |
| `internal/dispatch/dispatch.go:260-304` | Exposes the new option and refusal in help. | Necessary | No new command router or subcommand is introduced. |
| `internal/dispatch/build_paseo_host_test.go:19-88` | Exercises real build output, split-root path, null model, omitted fields, prompt, and incompatible-token absence. | Necessary | The oracle is emitted CLI output, not an instruction-file read. Exact wording assertions are brittle but serve output-shape claims, not live behavior proof. |
| Same file, 90-129 | Proves stamp-created worktree and bare refusal without an artifact. | Necessary | Uses existing real Git fixtures. It can also assert frontmatter preservation on refusal without creating another harness. |
| Same file, 131-170 | Proves ignored model disclosure, filename separation, and explicit-host precedence. | Necessary | Distinct failure modes justify these checks. Fixed global pointer filenames inherit the existing cross-process collision risk. |
| `build_host_test.go:38`, `build_json_ergonomics_test.go:161,189` | Updates existing error expectations. | Necessary | Maintains existing proof owners instead of adding duplicate tests. |
| `paseo_live_controls_test.go:18-27` | Selects the model from release SSOT and disables unrelated daemon services. | Necessary | A short value table is enough. There is no fake ACP service. |
| Same file, 29-57 | Scrubs outer runtime identity, sets isolated home, forces BatchMode, and constructs launch arguments. | Necessary | Explicit Codex full-access mode supports unattended fixture writes. This is a test-runtime choice, not production sandbox policy. |
| Same file, 59-75 | Runs Paseo with separate stdout/stderr and decodes stdout JSON. | Necessary | Directly addresses successful stderr noise. Missing per-command cancellation is a follow-up risk. |
| Same file, 77-94 | Selects run, local skip, or required-CI failure from CLI/auth availability. | Necessary | Reuses existing Codex auth policy. It does not silently skip missing CI auth. |
| Same file, 96-175 | Exercises config, scrubbed environment, provider arguments, stream separation, and preconditions. | Tighten | The stream regression and scrub checks are useful. The default-model equality mirrors the helper, and the config table repeats implementation. Keep claims limited to wiring. |
| `paseo_live_helpers_test.go:22-44,86-95` | Shares a daemon fixture with artifact name and timeout controls. | Tighten | Reuse is necessary for two smoke tests. One functional option for one string field is more API than needed. A name parameter suffices. This is polish. |
| Same file, 45-85,96-131 | Resolves prerequisites, creates canonical isolated state, starts a real daemon, polls readiness, and stops it in cleanup. | Necessary | Uses supported CLI lifecycle rather than a second supervisor. The close-then-bind port race and unbounded CLI process remain reliability limits. |
| Same file, 99-115 | Checks stopped state and captures one daemon log on failure. | Necessary | Makes cleanup observable. Hard process timeout can bypass Go cleanup. Raw log capture needs caution if upstream adds sensitive fields. |
| Same file, 133-150 | Provides the JSON wrapper and workflow-free non-Codex skill staging. | Necessary | Avoids real-workflow discovery. Non-Codex staging has weaker host isolation than the Codex lane and needs separate live evidence. |
| Same file, 152-170 | Creates a real split-root repository, state branch, task, and worktree-producing stage. | Necessary | This is the smallest useful fixture for both code and state writes. No mock worker runtime is introduced. |
| Same file, 172-202 | Creates an isolated Codex home, seeds supported auth, and installs a local marketplace from current skills. | Necessary | Reuses existing helpers and proves consumer installation. It avoids raw checkout paths in agent prompts. |
| Same file, 205-220 | Chooses a loopback ephemeral port outside the common defaults. | Tighten | Necessary isolation. Sixteen attempts and special default-port exclusions are minor defensive excess, not a new subsystem. |
| `paseo_live_smoke_test.go:20-56` | Launches a real installed-skill FO with a receipt assignment. | Necessary | Covers the actual transport. The explicit test prompt weakens the claim that selection and procedure originate only from installed skills. |
| Same file, 58-91 | Checks worktree branch, committed receipt, stage-report shape, entity-only final commit, and clean trees. | Tighten | Strong happy-path outcome checks. After `x9m`, the final entity commit can be an FO archive record, so it no longer proves the worker report commit was path-scoped. |
| Same file, 93-117 | Checks one labeled worker, epoch, exact parent/directory, archived timestamp, and live FO. | Necessary | Independent daemon observations enforce lifecycle outcomes. It does not count rogue workers missing the label. Adoption smoke adds a broader roster check. |
| Same file, 120-142 | Bounds outer polling and distinguishes idle, timeout, and errors. | Necessary | Reads actual wait status. The subprocess itself remains unbounded. |
| `paseo_adoption_live_test.go:22-35,201-221`, `x9m` | Parses fixture worker YAML and counts retry lines. | Necessary | Test-local parsing avoids a dead production schema implementation. No reflection-heavy registry is necessary. |
| Same file, 37-75 | Holds the worker on a file barrier, observes a committed running record, stops FO-1, and checks the worker survives. | Necessary | Supported `paseo stop` preserves parent identity. This is better evidence than timing-only interruption. Cleanup archives FO-1 even on early failure. See Section 4. |
| Same file, 78-93 | Launches FO-2 and releases the worker after generic FO activity. | Tighten | Ship-blocker 2: activity is not adoption. Prompts supply the procedure that the installed binding is meant to own. |
| Same file, 95-118 | Checks one worker across labeled/all rosters, archival, FO survival, and server identity. | Necessary | The broad roster closes the renamed-retry hole. It still trusts `role=first-officer` for exclusions instead of the two known FO IDs. |
| Same file, 120-150 | Follows entity archival and checks all durable record fields and zero retries. | Necessary | Binds expected values to independent requests and daemon snapshots. This closes the prior metadata assertion hole. |
| Same file, 151-178 | Checks task branch, committed receipt, a post-baseline entity-only report diff, and clean trees. | Necessary | This closes the earlier stale-report and multi-file-commit proof holes. It does not identify the actor that authored those commits. |
| `internal/status/paseo_worker_compat_test.go:14-87`, `x9m` | Exercises four record states in two H2 positions through `--validate`, `--read`, `--checklist`, and `--ac-scan`. | Necessary | No production status change is required. The eight cases target parser boundaries and held-reason contamination. Retry-ledger placement after a report is not covered. |
| `internal/contractlint/boot_resident_closure_test.go:18-23`, `936` | Adds the deferred reference to existing path-closure checks. | Necessary | One list entry extends the existing owner. It proves paths resolve, not that an FO loads them at the correct time. |
| `internal/contractlint/structural_checks_test.go:193-206`, `936` | Checks both adapter files exist. | Tighten | The ensign file is not covered by the FO closure walker. The FO existence check overlaps its mandatory read there. During simplification, retain only distinct coverage. |
| `internal/contractlint/paseo_live_lane_test.go:12-90`, `936` | Parses workflow structure and checks triggers, permissions, protected environment, required auth, smoke wiring, and model source. | Tighten | Useful structural checks. Shell substring presence and absence of the current model literal do not prove exact dataflow or exclude a different hardcoded model. |
| `.github/workflows/paseo-live-e2e.yml:4-56`, `936` | Adds manual/PR triggers, a read-only token, offline tests, and model output. | Necessary | Uses existing CI tools. The PR path filter omits several shared helpers consumed by the lane. Record this as follow-up. |
| Same file, 58-103 | Adds protected auth, runtime installation, binary build, and version diagnostics. | Necessary | Inputs enter shell through environment variables, not executable interpolation. Secrets remain absent from the offline job. |
| Same file, 104-120 | Runs required-auth dispatch smoke in an optional job and uploads diagnostics. | Necessary | Matches the locked optional-lane scope. Does not schedule adoption smoke. No adoption CI coverage is claimed here. |
| `skills/integration/review_choice_live_test.go`, two-dot-only | Shows four changed review-consent assertions. | Unjustified/Cuttable as sprint work | These are main-only ancestry differences, not a new test mechanism in either PR. |

### Per-component analysis

**Dispatch builder.** The production implementation is approximately forty changed lines in an existing package. It reuses envelope, name, path, stamp, and model machinery. A dedicated transport service or provider hierarchy adds no value here. Accepting `--advance` remains a generic builder capability. The Paseo binding declares worker reuse absent.

**Live harness.** Real daemon lifecycle, isolated credentials, staged skills, and real Git fixtures are justified by the DoD. The harness does not emulate ACP or own production worker recovery. Its largest weaknesses concern what the assertions prove, not excessive architecture.

**Status compatibility.** The H2 record works with the existing parser. The tests drive public status surfaces against fixture data. These checks are materially different from a test that searches an instruction file for required prose.

**Contractlint and CI.** Structural path and YAML checks belong in contractlint. They do not prove model behavior or adoption. The new lane correctly consumes `codex.exec` through the release source in Go and CI. Its substring lint remains a limited structural guard.

## Section 3 — Ship-Blockers

### 1. Recovery has conflicting decisions and incomplete terminal error handling

**Owner: `x9m` / PR #895.** `936` introduced the original wait guidance. `x9m` adds the conflicting recovery authority and owns integration of the assembled contract.

**Exact locations:** `skills/first-officer/references/paseo-fo-dispatch.md:80`, `:149-151`, and `:155`. Restart completion recovery is at `:130`.

A failed first worker with no valid report reaches two incompatible instructions. Verb 3 says inspect and hold. Recovery row 5 says archive and retry once. An FO that obeys the earlier rule never performs the promised recovery. An FO that obeys the later rule ignores an unconditional hold instruction. The document gives no precedence rule.

Row 10 adds another conflict. It treats “repeated wait” as a strict no-op. A second wait can report idle after a first timeout. Rows 1 and 4 require continued polling and completion handling for that sequence. Duplicate event delivery must not mean every repeated observation.

Terminal handling also contains a direct failed-command path. Row 4 accepts `AGENT_NOT_FOUND` with valid Git evidence, then requires an archive command that returns `AGENT_NOT_FOUND` again. It only names `AGENT_ALREADY_ARCHIVED` as successful. Row 6 similarly needs an explicit missing-agent no-op while recording an exhausted hold. Row 5 already demonstrates the required tolerance.

**Why this blocks:** `x9m` promises deterministic single-retry recovery and restart-safe completion. These paths can prematurely hold, suppress new completion, or stop after valid completion evidence. More tests cannot make contradictory instructions deterministic.

**Minimal fix:** Make the recovery section the sole decision owner. Change Verb 3 to inspect and route through that section. Restrict row 10 to already processed terminal evidence for the same attempt. After valid Git completion, treat an absent agent as already shut down. Apply the same rule during an exhausted hold. Keep permission precedence explicit.

**Proof to close:** Exercise an error/no-report first attempt, timeout followed by valid idle completion, and valid completion with an absent agent. Use independent durable state and retry counts as outcomes. A small live exercise or recorded one-off proof is sufficient. Do not add an unused Go recovery simulator or prose-grep tests.

**Evidence limit:** This is a static contradiction and command-path counterexample. The audit did not run a new live recovery scenario or modify the shipped contract.

### 2. Adoption smoke can pass without proving the installed running-worker adoption path

**Owner: `x9m` / PR #895.** The shared helpers from `936` need no architectural change to correct this proof.

**Exact locations:** `internal/ensigncycle/paseo_adoption_live_test.go:56-57`, `:78-93`, and `:140-178`. The shipped procedure under test is `skills/first-officer/references/paseo-fo-dispatch.md:125-169`.

FO-1 receives the worker-record schema and persistence procedure directly in its prompt. FO-2 receives explicit instructions to read that record, adopt without spawning, wait, archive, update the record, and advance. A missing or unloaded adoption binding can therefore be replaced by the test prompt. The final assertions do not distinguish those instruction sources.

The barrier does not establish running-worker adoption either. Line 85 accepts general `UpdatedAt` activity, idle, or error. After the 30-second loop expires, execution also continues without a successful activity condition. Line 89 proves only that W1 remains running immediately before release. It does not prove that FO-2 inspected or waited on W1.

A concrete false-green schedule is possible. FO-2 begins an unrelated startup action. The harness releases W1. W1 commits its report before FO-2 reads the worker record. FO-2 then processes an already completed worker using its explicit prompt instructions. One worker, zero retries, correct metadata, valid commits, and archival all still satisfy the final assertions.

**Why this blocks:** The sprint DoD specifically requires adoption of a running worker across an FO restart through the installed binding. The current test proves survival across restart and eventual completion. It does not isolate the installed adoption behavior or establish adoption before completion.

**Minimal fix:** Keep fixture facts, provider selection, and the worker barrier in the test. Reduce the FO prompts to task intent and a normal resume request. Remove the duplicate record/adoption implementation from those prompts. First observe an FO-2 adoption action, such as its actual inspect/wait call, while W1 is running. Then release W1. Use a narrow test observation of the real command stream or an existing runtime event. If that observation does not occur, fail the test. No new production record field, recovery controller, or lease is needed.

**Proof to close:** A run without the installed adoption binding must fail. A run where FO-2 only performs startup activity must not release W1. The corrected live smoke must pass with one worker, zero retries, durable report/code commits, and archival.

**Evidence limit:** The false-green schedule follows the test control flow. The audit did not execute a destructive mutation or claim an observed live failure. Existing successful smoke runs do not distinguish this schedule from the intended one.

## Section 4 — Non-Blockers / Next-Sprint Seeds

### Cross-cutting integration and necessity

1. **Make the transport override complete and early.** `first-officer/SKILL.md:38` says “first dispatch,” while shared reconciliation runs before dispatch at `fo-dispatch-core.md:231`. The loaded Codex adapter also contains unconditional native wait instructions at `codex-first-officer-runtime.md:30-38`. Clarify that selecting Paseo replaces worker lifecycle bindings and loads before reconciliation on resume. The existing “alongside” wording supports that interpretation, so this audit treats the ambiguity as tightening rather than another proven failure.

2. **Strengthen candidate scoping without another registry.** `paseo-fo-dispatch.md:133-136` falls back on raw roster emptiness before local directory filtering. An unrelated same-label worker can suppress the fallback. The worker-label fallback also lacks an explicit epoch check. `started` is not an attempt timestamp. `internal/dispatch/stamp.go:128` only fills it when empty. Filter candidates first, preserve observed epoch identity, and define an ambiguous closed-roster result. Multiple FOs on one workflow remain outside scope.

3. **Make completion freshness explicit.** `paseo-fo-dispatch.md:140` requires a new report but uses entity revision activity as its concrete Git check. FO record commits also change that entity. Require a report change after the baseline and the assigned deliverable evidence. `paseo_adoption_live_test.go:159-174` already demonstrates the stronger report-diff check. Do not promote clean Git status alone into proof of a new code deliverable.

4. **Keep schema and commit mechanics small.** The ten-field record is sufficient. `server_id` records the current local daemon, not remote-worker support. The `2s` index-lock retry at `paseo-fo-dispatch.md:109` is an arbitrary local recovery detail. Prefer the existing state commit owner and unchanged-record handling. Neither a fixed YAML key order nor a production worker-state package is justified.

5. **Consolidate prose before changing caps.** The clearest cuts are the repeated spawn chain, ten-key identity list, self-loading ensign bullet, redundant matrix columns, and repeated cascade guard. New prose also contains long compound instructions, especially `paseo-fo-dispatch.md:18,109,129,136,140` and `build.go:794`. Apply the declared simple-English rules to touched prose. These issues do not justify holding the release beyond the two correctness/proof blockers.

6. **Preserve quarantine and interpret lint honestly.** Changed instruction-file reads remain in `internal/contractlint`. The new checks concern existence, closure, YAML structure, and literal model absence. Runtime behavior still belongs to the live tests. The unchanged shared core has only three bytes of cap headroom. Do not solve that by adding another prose-to-code consistency test.

### Live harness, CI security, and model ownership

1. **Security controls are present.** The workflow uses `pull_request`, never `pull_request_target`, and grants only `contents: read` at `.github/workflows/paseo-live-e2e.yml:22-36`. Secrets belong to the protected live job at `:62-66`. The offline job receives none. Manual versions and provider values enter through environment variables at `:67,84-89`. This avoids shell execution through direct input interpolation. Environment approval remains a trust boundary before running PR code.

2. **No direct credential artifact copy was found.** Codex auth uses the existing isolated-home helpers at `paseo_live_helpers_test.go:172-201`. Their normal CI home location stays outside the uploaded repository artifact subtree. `codex_live_runner_test.go:184-186` removes that home in cleanup. `paseo_live_helpers_test.go:107-113` copies raw daemon logs on failure, and `paseo_live_smoke_test.go:135` logs recent activity. An upstream logging change or an agent printing a secret can still expose sensitive output. This is a residual diagnostic risk, not evidence of a current leak.

3. **Current model defaults use SSOT.** `paseo_live_controls_test.go:18-19` reads `release.LiveModels["codex.exec"]`. CI reads the same file at `.github/workflows/paseo-live-e2e.yml:48-54,109`. The lint at `paseo_live_lane_test.go:78` only bans the current model literal. A different literal can evade it. The run-step substring checks also do not prove exact shell dataflow. Strengthen those structural checks only for a distinct failure oracle.

4. **Bound individual CLI calls.** `runPaseoJSON` uses `exec.Command` at `paseo_live_controls_test.go:62`. A wedged CLI can outlast the surrounding polling deadline. The package timeout then prevents ordinary Go cleanup. A context timeout around each real CLI call is sufficient. Avoid a second process supervisor. The free-port close at `paseo_live_helpers_test.go:213` also leaves a small bind race.

5. **Use a package timeout that covers both tests.** The local combined command at `docs/runtime-live-ci.md:108` allows 15 minutes total. Each helper allows 12 minutes, and adoption has separate phases. Increase the documented budget or run the tests separately. This concerns slow-run reliability, not the short successful runs already recorded.

6. **Expand coverage through existing owners.** The optional workflow runs only dispatch smoke at `.github/workflows/paseo-live-e2e.yml:111`. Adoption therefore has no automatic lane in this stack. Its path filter also omits shared helpers such as `codex_liveenv.go`, `live_test.go`, and parts of runtime installation. If recurring cost is acceptable, add adoption and shared-path coverage in a follow-up. Required Claude/Codex/Pi live lanes still apply under `docs/dev/README.md:85`.

7. **Tighten the dispatch smoke after adding FO record commits.** `paseo_live_smoke_test.go:86` checks only the final state commit. That can now be the FO archival record rather than the worker report. Reuse the stronger report-history assertion from adoption smoke. Neither smoke independently attributes every write to the worker. Current evidence proves outcomes more strongly than authorship.

8. **Use known FO IDs for roster exclusion.** `paseo_adoption_live_test.go:105-108` excludes every agent labeled `role=first-officer`. Comparing against `fo1.AgentID` and `fo2.AgentID` avoids hiding an accidental extra agent with that label. This is a small assertion improvement.

9. **Clarify failure cleanup and non-Codex isolation.** FO-1 cleanup at `paseo_adoption_live_test.go:63` runs even if W1 remains active after a failure. The whole fixture daemon then stops, so cleanup is not successful-path cascade proof. Non-Codex staging at `paseo_live_helpers_test.go:142-149` uses `.agents/skills` but leaves the operator home in the environment. The default Codex lane does not prove every non-Codex skill discovery or sandbox configuration.

### Dogfood friction disposition

| Existing friction | Disposition for these PRs | Reason and next action |
|---|---|---|
| F1/F2: Paseo executable build and sandbox path | Non-blocker. Resolved operating setup | The ledger records a working CLI. Keep executable readiness outside the transport implementation. |
| F3: stripped `PASEO_AGENT_ID` | Next-sprint/upstream seed | Explicit backend selection and assignment fallback exist. Parent attribution and filename separation still need the environment fix. Do not add this marker to inner-host detection. |
| F4: Antigravity version gate omits contract token | Next-sprint Spacedock seed. Known support limitation | `internal/cli/cli.go:913-917` returns before the token. The shared gate requires it at `first-officer-shared-core.md:9`. Fix version reporting independently. This does not invalidate the declared default Codex lane, but unrestricted “any FO” support is not yet proved. |
| F5: installed Antigravity skill drift | Next-sprint install/freshness seed | The test stages current skills. The real installation still needs a reliable freshness path. |
| F6: bootstrap through Codex host mode | Resolved by `936` | Explicit `--host paseo` and a neutral pointer remove the bootstrap workaround. |
| F7: worktree sandbox denies common Git/state/workflow paths | Next-sprint upstream seed. Deployment prerequisite | Temp fixtures cannot reproduce the non-temp boundary. The ledger records a scoped `.safehouse` workaround. Repair upstream path grants for real worktrees. Do not grant ordinary ensigns daemon control. |
| F8: ordinary ensign lacks daemon/auth/identity access | Expected boundary plus F3 follow-up | Run live daemon tests from Commander or CI. Keep the assignment-reference adapter fallback and durable blocker report. |
| F9: idle before tests and commits finish | Completion gate works. Next-sprint adapter/UX seed | Git evidence prevents false completion. The current fresh-only recovery will replace an incomplete idle worker once. Same-worker nudging needs an explicit lifecycle decision, not an inferred exception. |
| F10: provider quota interruption | Positive evidence. No new mechanism | Host-neutral assignment allowed provider failover. It does not prove every recovery branch. |
| F11: stacked base and worktree preparation | Next-sprint dispatch ergonomics seed | The existing stamp starts from trunk. Stacked base selection remains general dispatch work, not a Paseo transport defect. |
| Workspace/tab lifecycle | Next-sprint Paseo integration seed | Find-or-create avoids duplicate workspaces. Worker archive does not express workspace ownership or safe removal. Consider `--background` for temporary audit workspaces and archive owned unused workspaces after all children stop. Never archive a shared workspace merely because one task completes. |

The F4/F7 classification follows the locked local-lane scope and the explicit outside-sprint boundaries in the sprint index. It does not claim those environments work without the recorded prerequisites. The release notes must not turn a successful default Codex test into a universal sandbox or Antigravity claim.

### Preflight review closure

| Existing preflight finding | Assembled-stack disposition |
|---|---|
| M1: unnecessary simulator and oversized record | Closed. No production worker simulator exists. The record has ten required keys and a conditional reason. |
| M2: lossy roster fields | Closed for field selection. Inspect supplies exact directory, creation time, and archival state. Candidate scoping remains a follow-up. |
| M3: duplicate workspaces | Closed in the binding through find-or-create. Safe workspace removal remains separate lifecycle work. |
| M4: mixed JSON streams | Closed. The real CLI helper separates stdout/stderr, and the offline regression exercises successful stderr noise. |
| M5: parent cascade archival | Closed on the intended restart path through `paseo stop` and the live parent check. Failure cleanup has a narrower limitation. |
| M6: installed-skill discovery | Marketplace staging and named loading address the raw checkout-path problem. Prompt-supplied procedure still weakens behavioral attribution, especially adoption. |
| M7: label and epoch contract | Closed for invariant worker labels and defined epochs. Adoption additionally counts unlabelled non-FO agents. |
| M8: shared helper, model source, and CI controls | Closed in the implementation. The structural lint is useful but does not prove every security or dataflow property. |

## Proof & Hygiene Spot-Check

The checks below ran against `01d66ed160eb6d44b60c7c34c58534094d063ae6`.

| Check | Result |
|---|---|
| `go test -race -count=1 ./internal/dispatch ./internal/status ./internal/contractlint ./internal/release` | PASS. Dispatch 47.423s, status 59.803s, contractlint 7.135s, release 60.379s. |
| `go test -race -count=1 ./internal/ensigncycle -run '^TestPaseo'` | PASS, 1.662s. These are offline controls, not live smoke execution. |
| `go test -tags live -c -o /tmp/audit-live.test ./internal/ensigncycle` | PASS. Temporary binary removed after successful compilation. |
| `go vet -tags live ./internal/ensigncycle` | PASS. |
| `git diff --check origin/main..HEAD` | PASS. |
| `gofmt -l ./cmd ./internal` | Only unchanged `internal/release/runtime_live_evidence_workflow_test.go` appears. `origin/main` has the same formatting difference. Changed Go files are formatted. |
| `go test ./...` | PASS across all packages, including CLI (215.435s) and ensigncycle (245.772s). |
| `go test ./... -race` | PASS across all packages, including CLI (232.448s) and ensigncycle (273.677s). |
| `git status --porcelain` in the audit code worktree | Empty at start, after targeted checks, and after both full suites before the report commit. |

The explicit read-only boundary overrides the general `gofmt -w` instruction. The audit used `gofmt -l` and did not rewrite code. No new live daemon/model run occurred in this audit. Live compilation and offline tests do not establish live adoption behavior. Prior smoke results remain supporting sprint evidence, subject to the proof gap in Section 3.

The audit inspected all 20 sprint-authored files and all three additional two-dot comparison artifacts. It did not inspect GitHub environment configuration or attest required lane status on the eventual merge SHA. Those remain merge checks for the Commander.

## Audit Task Ledger

No journal tool is available. This report retains the investigation facts, findings, limits, and work status.

- Complete: Read sprint goal, preflight review, friction ledger, runtime support guidance, and declared proof policy.
- Complete: Inspect both commit ranges and identify main-only ancestry differences.
- Complete: Inventory every changed prose section and major behavior component.
- Complete: Assess preflight closure, recovery seams, operational frictions, CI security, model ownership, and proof limits.
- Complete: Run the requested targeted checks without code edits.
- Complete: Finish both broad suites, final code hygiene, and the report for path-scoped publication.

Only `_roadmap/paseo-backend/pre-merge-audit.md` is authorized for the state commit. No entity file or frontmatter is part of this audit deliverable.
