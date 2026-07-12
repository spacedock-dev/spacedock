---
id: eqn4ecmdy9d5a0meeqxwjwfa
title: Make survey work inside Claude Cowork
status: ideation
source: captain request
started: 2026-07-12T06:27:41Z
completed:
verdict:
score:
worktree:
issue:
---

Make the existing `survey` skill recognize Claude Cowork, obtain the Spacedock binary only after an explicit network-access handoff, and analyze Cowork-native session history without treating local Claude Code files as authoritative.

## Problem

`survey` currently enters its `agentsview`/`~/.claude/projects` path before distinguishing its host. Claude Cowork instead exposes a session inventory and transcript reader through its tool surface, and may start without a `spacedock` executable. The current order therefore sends a Cowork user toward irrelevant local-history setup and can turn a network permission boundary into a misleading “missing binary/runtime” failure.

The Cowork path must be selected from positive capabilities, not guessed from missing files. It must ask before the first network request, install the checksum-verified release for the observed OS/architecture into a writable user location, and use only host-native, sanitized session evidence. Repo-only conclusions must remain visibly unverified.

## Design choice

**Recommend blending Cowork support into `survey`; do not add a `using-spacedock-with-cowork` skill.**

| Evidence | Dedicated Cowork skill | Cowork branch in `survey` |
| --- | --- | --- |
| Discoverability | Requires users asking “survey” to know and invoke a second skill, or relies on overlapping descriptions to win routing. | Existing “survey/catch me up/orient me” triggers remain the single front door. |
| Ownership | Would own installation, survey adaptation, commission, and filesystem workarounds at once. | Survey owns its own input selection, consent boundary, and evidence degradation. |
| Overlap | Must restate or redirect most of survey and can drift from its report/offer contract. | Shares one report contract; only the evidence adapter differs. |
| Maintenance | A second user-invocable skill needs routing, packaging, and cross-skill synchronization tests. | One positive runtime gate and Cowork fixture extend the current survey smoke suite. |

The supplied Cowork skill remains research evidence only. Its separate connected-folder/git/commission issue is not absorbed here.

## Proposed approach

Add a short `0. Select the evidence adapter` section before current step 1. Detection is an executable positive probe, including deferred tools:

1. Run `ToolSearch(query="select:mcp__session_info__list_sessions,mcp__session_info__read_transcript", max_results=2)` to load both schemas. An initially absent schema is not evidence of absence.
2. Only when both exact definitions resolve, call `mcp__session_info__list_sessions` once with its schema's empty/default read-only input. A successful response selects `cowork` and that response becomes the scan inventory; do not call it twice. A missing definition selects the unchanged local path and records only `definition-missing`. If both definitions resolve but inventory errors, stop as `cowork-probe-failed`; do not fall through to local history after Cowork capabilities have been positively identified.
3. Do not infer Cowork from a missing executable, missing `~/.claude/projects`, Linux, or a sandbox error.

This binding is proven live by a sanitized Cowork run named `TestLiveClaudeCoworkSurveyProbe` and replayed without Cowork by `TestSurveyCoworkEventReplay` over `skills/integration/testdata/survey/cowork-tool-events.jsonl`. The fixture records tool names, success/error classes, counts, ordering, and synthetic payload fields only; it contains no real payload values. The replay must show ToolSearch before either deferred tool, both definitions plus successful inventory selecting Cowork, partial definition results selecting local, and a resolved-schema/inventory error stopping without transcript or local-history calls. The current `survey_probe_test.go` continues to execute the unchanged local agentsview probe.

### Cowork mounted-project persistence model

The selected Cowork lane supplies a mounted **current project folder** as a host capability and that mount persists across Cowork session environments. Bind exactly that supplied root; do not enumerate arbitrary connected roots and do not derive the root from `pwd` alone. `pwd` may confirm that shell execution is bound under the supplied mount, but it is not discovery. If the host does not supply the current-project-folder capability or the shell binding cannot address it, stop with `current-project-folder-unavailable` and make no install/network/project mutation. Durable evidence records only `current-project-folder: available|unavailable` and `shell-binding: confirmed|failed`, never a path or project name.

The persistent Cowork execution identity is exactly `<current-project-folder>/.spacedock/bin/spacedock`. Do not use `$HOME/.local/bin/spacedock`, PATH, another connected root, or a project cache/pointer. Probe and invoke this exact mounted path directly. An absent path enters the consent/install flow. Any existing filesystem entry—including a dangling symlink—is existing state: direct invocation success means working/reusable; non-executable or failed invocation means broken and stops without network or overwrite. Repairing broken state requires a separate explicit user decision and is not inferred from rerunning survey.

Lifecycle phases are fixed: **run 0** receives the host current-project-folder capability, observes the exact mounted binary absent, records zero network events, presents the Settings/Full-network/relaunch instruction, and stops. **Run 1** is post-settings: install and checksum-verify in session temp/HOME, invoke that session copy, then create the persistent mounted path as described below and directly verify it. **Run 2** uses a fresh session HOME with the same mounted project folder, directly verifies/executes the mounted binary, and records zero network/download/install events. Evidence contains only capability/result classes and booleans.

On the Cowork branch:

1. Bind `{project_root}` from Cowork's host-supplied current-project-folder capability and set `COWORK_BIN="{project_root}/.spacedock/bin/spacedock"`. Never substitute `pwd` or PATH. Execute this probe block as written; `COWORK_INSTALL_ABSENT` enters step 2, `COWORK_INSTALL_BROKEN` stops, and silent success reuses the mounted binary:

   ```sh
   if [ -e "$COWORK_BIN" ] || [ -L "$COWORK_BIN" ]; then
     if ! "$COWORK_BIN" --version >/dev/null 2>&1; then echo "COWORK_INSTALL_BROKEN"; fi
   else
     echo "COWORK_INSTALL_ABSENT"
   fi
   ```

   The `-L` arm classifies a dangling symlink as broken rather than absent. Success reuses the exact mounted file with no network activity. Non-executable or failing existing content stops; do not search PATH, download, or overwrite it. Absence proceeds to consent even when another valid `spacedock` exists on PATH.
2. When the exact mounted path is absent, display the exact prompt below and stop. The user changes Cowork Settings to **Full network access**, relaunches, and runs survey again; nothing is downloaded before that relaunch/rerun. Before consent there must be no network request or project mutation.
3. On the post-relaunch rerun, allocate unique session paths (`INSTALLER_TMP=$(mktemp "${TMPDIR:-/tmp}/spacedock-install.XXXXXX")` and `SESSION_INSTALL_DIR=$(mktemp -d "${TMPDIR:-/tmp}/spacedock-bin.XXXXXX")`), run the checksum-verifying installer with `SPACEDOCK_INSTALL_DIR="$SESSION_INSTALL_DIR"`, and directly verify `VERIFIED_BIN="$SESSION_INSTALL_DIR/spacedock"` via `"$VERIFIED_BIN" --version`. All acquisition, target, download, checksum, extraction, or session-verification failures happen before project mutation and leave `.spacedock/bin/spacedock` absent. Internal request-domain and temp-cleanup ledgers remain as previously specified and stay out of user prose.
4. First creation on the mounted folder is no-overwrite and does not require unlink or atomic rename: create `.spacedock/bin/`; re-check `[ ! -e "$COWORK_BIN" ] && [ ! -L "$COWORK_BIN" ]`; reserve the absent final name with a shell noclobber create (`( set -C; : > "$COWORK_BIN" )`); then `cp "$VERIFIED_BIN" "$COWORK_BIN"`, `chmod 0755 "$COWORK_BIN"`, and directly invoke `"$COWORK_BIN" --version`. Do not create a sibling temp/pointer and do not claim atomic installation. Directory creation or reservation failure leaves no binary content guarantee beyond observed state. Copy/chmod/final-verification failure may leave a partial/broken final path; report `Cowork Spacedock install incomplete` and stop. Future runs classify it broken and refuse overwrite until a separate explicit repair decision. This ordering works without unlink; the live lane must prove first creation on its deletion-denied mount.
5. Reuse the successful inventory probe. Sort idle sessions by the inventory's most-recent timestamp descending (stable tie-break: opaque id, used only in memory), omit active sessions, and sample at most the newest 20 idle sessions. For each sampled session, call `read_transcript(max_wait_seconds: 0, limit: 15)` and request the newest page. If the returned schema/result exposes both a continuation cursor and `truncated=true`, fetch at most two more 15-message pages for that session (45 messages maximum); otherwise do not invent pagination. Stop paging on no cursor, `truncated=false`, repeat cursor, or the 45-message ceiling. Tool errors omit that session from analyzed counts and increment a generic read-failure count.
6. Analyze only the returned window. **Retained existing survey behavior:** render directly without a scratch preamble; fill every displayed value from evidence; never invent an empty signal; keep the existing `manual` / `exploration` / `knowledge-work` / `unlabeled` meanings and their mode-keyed commission offers; keep conservative open-frontier wording when artifact evidence is unavailable. **Cowork-specific adapter behavior:** the 20-session/15-message/three-page bounds, sampling disclosure, message-level decision/interruption extraction, `sample-unverified` label, privacy-generalized cluster names, and repo-only section omissions. These are host adaptations, not changes to the local agentsview path.

   For the Cowork adapter, a **decision** is an observed assistant question offering alternatives or requesting approval; it is resolved only when a later sampled user message answers it, otherwise `OPEN (sample-unverified)`. A **hand-steering interruption** is a sampled user turn that corrects, rejects, or redirects work after an agent proposal/output; initial requests and simple affirmative answers do not count, and one user turn counts at most once. Apply the existing mode meanings to observed messages: repeated execute/build-and-review loops -> `manual`; iterative rejection/revision of creative or design paths -> `exploration`; repeated intake/process/commentary -> `knowledge-work`; insufficient or mixed evidence -> `unlabeled`. Cluster labels must be generic activity categories, never copied titles, names, or identifying attributes.
7. Project the report explicitly. Keep `WHAT THIS GIVES YOU`, sampled `BY THE NUMBERS`, `HOW YOU WORK`, `THREADS TO PULL`, `RECENT DECISIONS`, `WORKSTREAMS`, `WHAT THIS CAN'T SEE`, and the mode-keyed commission offer. Omit `BACKLOG`, `WORK BY AREA`, `CODEX`, scaffold comparison, dispatch/subagent facts, no-follow-up counts, shipped claims, and git/PR counts because Cowork session evidence cannot calculate them. Every open thread is labeled `sample-unverified`; no section may imply a repo cross-check occurred. The title is `SpaceDock survey — Cowork session sample`, followed by `sampled {analyzed} of {idle_total} idle sessions (newest first; up to 45 messages/session; {read_failures} read failures; active sessions excluded)`. If inventory timestamps are unavailable, omit a day-span rather than estimating one.
8. Fill every displayed value from the bounded scan, disclose sampling/truncation/read failures, show `none observed in this sample` for an applicable empty signal, and emit no literal `{slot}`. Preserve the end-of-report commission offer, scoped to the sampled evidence.

### Exact survey instruction change

Before (current opening route):

> Run the four steps in order: **check agentsview → scan → recognize scaffold → report and offer**.

After:

> Run `ToolSearch(query="select:mcp__session_info__list_sessions,mcp__session_info__read_transcript", max_results=2)` so lazy-loaded schemas cannot false-negative. If both exact tools resolve, call `mcp__session_info__list_sessions` once. A successful call selects **bind current project folder → check exact mounted Spacedock execution → ask before network → bounded Cowork scan → sampled Cowork report and offer**; reuse its inventory and never enter agentsview or repo probes. Missing Cowork definitions select the existing local path. Resolved Cowork tools plus an inventory error or missing current-project-folder capability stop; neither case falls into local history.

Insert this exact permission prompt when the Cowork branch cannot invoke Spacedock:

> Spacedock is not installed in this Cowork project. In Cowork Settings, choose **Full network access**, then relaunch this session and run survey again. Nothing will be downloaded before you relaunch.

### Privacy boundary

Session identifiers are in-memory tool handles only. Never echo or persist identifiers, session titles, raw tool payloads, transcript text, proper nouns, account details, credentials, tokens, machine-specific paths, or paraphrases/redactions that preserve a private person's, organization’s, project's, or event's identifying attributes. User-visible workstream and decision labels must be generic activity categories; durable evidence is limited to tool names, result classes, ordering, counts, sampling metadata, generic categories, and synthetic fixture values. Privacy tests inject raw and transformed canaries (case-folded, punctuation-stripped, partial-token, title-derived, and paraphrased-attribute forms) and reject every variant in output and durable files. Do not write raw Cowork responses to disk or pass any session material to the installer or a network call.

### Read-only boundary

The existing local agentsview survey remains strictly read-only with respect to the project and session sources. Cowork session-history reads are also read-only. The consented Cowork bootstrap is the narrow project mutation: after the settings relaunch it may create exactly `<current-project-folder>/.spacedock/bin/spacedock` (plus the `.spacedock/bin/` directories when absent). No transcript/session material enters that file. Aside from this disclosed bootstrap path, connected/project files remain unchanged. The report calls it a project-local Cowork helper installation.

### Documentation change proposed for implementation

In `docs/site/get-started/survey.md`, replace the local-only opening paragraph:

> When you use the `spacedock:survey` skill, it looks at your existing agent conversation logs on local disk (through agentsview, an open source session-history tool). It is read-only; if agentsview is missing, it asks before installing it.

with these two paragraphs:

> When you use the `spacedock:survey` skill, it reads the session history available to the current host. Local coding-agent sessions use agentsview; the local scan is read-only, and if agentsview is missing, Survey asks before installing it.
>
> In Claude Cowork, Survey reads a bounded session sample. On first use it stores the Spacedock helper at `.spacedock/bin/spacedock` in the current project folder so later Cowork sessions can reuse it; other connected files stay unchanged.
>
> If Spacedock is unavailable, set Cowork to **Full network access** in Settings, relaunch, and run survey again.

No install-page or commission documentation changes belong to this task.

## Out of scope

- General Claude Cowork workflow commissioning or git support.
- A new `using-spacedock-with-cowork` skill or general Cowork runtime adapter.
- Changes to connected-folder deletion semantics.
- Persisting or publishing Cowork transcripts, session identifiers, account data, filesystem paths, credentials, or other personal/private information.
- Automatic network-policy changes or downloads without the user's explicit enablement.
- GitHub issue creation.

## Acceptance criteria

**AC-1 - A successful, lazy-load-safe Cowork capability probe selects exactly one Cowork evidence path; missing definitions retain existing local behavior, while a resolved inventory error stops.**
Verified by: `TestLiveClaudeCoworkSurveyProbe` produces a sanitized event recording and `TestSurveyCoworkEventReplay` replays `cowork-tool-events.jsonl`; ToolSearch precedes deferred calls, both exact definitions plus one successful inventory call select Cowork, partial definitions select local, a resolved-schema/inventory error stops, and every non-local case records zero agentsview/repo-probe operations.

**AC-2 - Cowork binds only the host-supplied current project folder and probes exactly `.spacedock/bin/spacedock`; missing capability stops, working state reuses with zero network, absent state prompts, and all broken states refuse overwrite.**
Verified by: a capability/event-ledger fixture plus live events. `current-project-folder: unavailable` causes no path/network call. Available+absent emits the prompt with 0 requests/project changes. Working invokes only the mounted exact path with 0 requests. Non-executable, failing, and dangling-symlink entries stop broken. A valid later-PATH binary is never invoked in any case. Evidence contains only capability/result classes.

**AC-3 - After the settings relaunch, the verified binary is first-created at the exact mounted project path without overwrite/unlink, and a fresh-HOME second session reuses it with zero network.**
Verified by: run 1 installs/verifies in session space, reserves the absent mounted filename with noclobber, copies/chmods, then directly verifies it. Failures before reservation leave the project byte-identical; reservation/copy/chmod/final-verify failures may leave a broken path and must stop without retry. Run 2 uses a fresh HOME plus the same mounted project and directly executes the mounted binary with 0 request/install events. A deletion-denied fixture proves mkdir, noclobber reservation, copy, chmod, and direct execution without unlink. Byte-level comparison permits only creation of `.spacedock/bin/spacedock` and its parent directories.

**AC-4 - Cowork survey renders a complete, honestly sampled orientation with correct filled values and no repository-only claims.**
Verified by: a deterministic replay inventory with four idle sessions and one active session. Its newest-first transcript pages produce exactly `4 of 4` analyzed, `2` hand-steering interruptions, `1` resolved decision, `1` `OPEN (sample-unverified)`, four workstreams (`manual`, `exploration`, `knowledge-work`, `unlabeled`), `0` read failures, and the `up to 45 messages/session; active sessions excluded` disclosure. The output contains no literal `{slot}` and omits BACKLOG, WORK BY AREA, CODEX, scaffold, dispatch, no-follow-up, shipped, and git/PR sections. Mutant cases cover 21 idle sessions (20 sampled), truncation with one/two continuation pages, repeated cursor, missing pagination, and one read failure.

**AC-5 - Cowork output and durable artifacts contain only aggregate/synthetic evidence, not raw or transform-preserving private session material.**
Verified by: canary identifiers, titles, transcript sentences, proper nouns, account/token/path values, and identifying attributes injected into the in-memory harness; recursive assertions reject exact, case-folded, punctuation-stripped, partial-token, title-derived, and paraphrased-attribute canaries from output and files, while expected aggregate counts, sampling disclosure, and generic activity categories remain. The harness also asserts no raw tool response is written.

**AC-6 - User documentation preserves local read-only/ask-before-install behavior and separately discloses Cowork's project-local `.spacedock/bin/spacedock` mutation and short network relaunch instruction.**
Verified by: render and inspect the Survey page against the exact proposed replacement; runtime replay and live smoke remain the behavioral proof rather than a prose substring test.

## Test plan

Implementation starts with the fixture harness so routing, consent, privacy, and rendered-report assertions fail against the current skill before its instructions change.

- **Cowork event replay (medium):** add `TestSurveyCoworkEventReplay` and the committed sanitized `skills/integration/testdata/survey/cowork-tool-events.jsonl`. The runner consumes tool events and synthetic responses, enforcing the exact ToolSearch -> inventory -> bounded transcript sequence and rendering the Cowork projection. Fixture/mutant expectations are AC-1/AC-4's exact values. The existing `TestSurveyInstallProbe` remains the executable local-path control; no generic “skill smoke” is assumed.
- **Capability/path routing (medium):** replay host current-project-folder available/unavailable classes without storing roots. Execute exact-path probes for absent, working, non-executable, failing, and dangling-symlink state plus later-PATH shadows. Assert no arbitrary-root enumeration, no `pwd`-only discovery, and no PATH binary invocation.
- **Installer/first-create CLI tests (medium):** keep the consent/request/architecture/checksum/temp-cleanup seams. Install to session space first, then exercise mkdir -> absent recheck -> noclobber reservation -> copy -> chmod -> direct final verification. Mutate every boundary: pre-verification failure leaves project unchanged; race/existing path refuses overwrite; copy/chmod/final-verify failure leaves a classified broken path and no retry. Run on a deletion-denied mount fixture and assert no unlink/rename call. Snapshot comparison allows only `.spacedock/`, `.spacedock/bin/`, and `.spacedock/bin/spacedock` creation.
- **Fresh-session persistence (medium):** run 0 proves absent mounted binary -> 0 requests -> Settings/Full-network/relaunch stop. Run 1 uses HOME A and the mounted project to install/verify persistent identity. Run 2 uses empty HOME B and the same mount; assert direct mounted-path invocation, 0 network/download/install events, and no additional project bytes. Record only capability/result booleans/classes.
- **Report/privacy replay (medium):** render the exact four-session baseline and pagination/failure mutants from synthetic in-memory payloads. Assert filled values, no placeholders, sampling disclosure, omission matrix, open-thread label, mode mapping, and transformed-canary absence across stdout and the test temp tree.
- **Documentation render (low):** build the docs site and visually/readably inspect the changed Survey page.
- **Live Cowork smoke (required):** run 0 records current-project-folder capability available, exact mounted binary absent, and 0 requests before the relaunch stop. Run 1 (HOME A) verifies in session space, first-creates/directly verifies `.spacedock/bin/spacedock`, and records project delta class `expected-helper-only`. Run 2 (fresh HOME B, same mount) directly executes the mounted binary with 0 network/download/install events. A parallel deletion-denied lane proves first create without unlink; missing-capability and broken-state lanes prove stops. Durable evidence contains only generic capability/result classes and booleans—no root, project name, session identifier, or payload.

## Risk spike

The captain confirms the selected lane's host-supplied current project folder is mounted and persistent, superseding the no-cache/session-only design. The riskiest mechanism is first creation on a deletion-denied mount without overwrite or unlink. Implementation begins with the deletion-denied fixture and then the live HOME-A/HOME-B same-mount proof. No atomicity is claimed: failures after noclobber reservation may leave a broken final path, and recovery requires a separate explicit repair decision.

## Stage Report: ideation

- DONE: Compare a dedicated using-spacedock-with-cowork skill with integrating Cowork support into survey; recommend one smallest coherent design using discoverability, ownership, overlap, and maintenance evidence.
  Recommended one capability-gated branch in `survey`; the comparison table records why a second user-invocable skill would overlap and drift.
- DONE: Update the task with a concrete behavior-first design, including exact skill-routing or before/after instruction text, privacy boundaries, and the user-visible network-access prompt.
  The proposed route, exact prompt, privacy contract, acceptance criteria, and Survey documentation replacement are specified above.
- DONE: Define reproducible proof for Cowork detection, no download before permission, architecture-aware binary setup, and Cowork-native session evidence; spike the riskiest unverified mechanism or record why no spike is needed.
  Fixture, installer, privacy-canary, docs-render, and mandatory live-Cowork proofs are defined; existing executable installer seams and supplied live research make a throwaway spike unnecessary.

### Summary

Ideation selects integration into the existing survey skill as the smallest coherent design and explicitly excludes the separate commission/filesystem issue. The design detects Cowork from two positive session capabilities, stops before all network activity for user consent, delegates architecture/checksum handling to the existing installer, and requires a privacy-bounded live Cowork smoke before the runtime claim can pass.

## Stage Report: ideation

- DONE: Bind Cowork detection to a concrete, executable positive capability probe that covers lazy-loaded tools, and name a real fixture or live harness.
  The design now binds ToolSearch plus one successful read-only inventory call, with named live and sanitized event-replay tests.
- DONE: Specify a concrete installer acquisition/locating route and existing-install probe; inventory every network request/domain and prove zero network calls before affirmative consent.
  The exact PATH probe, raw installer URL, four-host allowlist, consent ledger, and no-redownload/broken-install behavior are specified.
- DONE: Remove or precisely narrow the unsupported atomic-install claim.
  Claims now stop at the real pre-final-write boundary for a fresh absent path; existing paths are never overwritten and no atomic replacement claim remains.
- DONE: Define exact Cowork survey report semantics, deterministic expectations, honest sampling disclosure, omission rules, and expanded privacy protection.
  The task fixes sorting, sample/page ceilings, decision/interruption/mode mappings, projection, expected fixture values, transformed canaries, and no-raw-payload persistence.

### Summary

This repair keeps Cowork inside `survey` but replaces every abstract runtime assumption with an executable ToolSearch/inventory binding and named live-to-replay proof. It also fully discloses and gates installer networking, narrows failure guarantees to shipped behavior, and specifies a bounded, privacy-preserving Cowork report whose values and omissions can be checked deterministically.

## Stage Report: ideation — SUPERSEDED (historical session-local path design)

- DONE: Replace the PATH-based install identity probe with explicit execution-path semantics.
  The design invokes and verifies only `"$HOME/.local/bin/spacedock"`; absent, working, and broken exact-path states cannot be shadowed by PATH.
- DONE: Add two executable test mutants for later-PATH binaries.
  AC-2 and the installer test plan require absent-exact-path/valid-PATH consent and broken-exact-path/valid-PATH stop cases with invocation and network ledgers.
- DONE: Correct AC-1's heading to distinguish missing definitions from resolved inventory errors.
  The heading now says missing definitions retain local behavior while resolved inventory errors stop, matching the body and verifier.
- DONE: Replace the fixed temporary installer filename with a unique `mktemp` path and require honest cleanup behavior.
  The acquisition contract uses a unique path, traps cleanup, proves normal deletion, and reports without overclaiming when deletion is blocked.

### Summary

SUPERSEDED HISTORICAL SUMMARY: this cycle chose an exact session-local HOME execution path rather than PATH resolution. The captain's later correction keeps the exact/no-shadow semantics but moves the persistent identity to `<current-project-folder>/.spacedock/bin/spacedock`.

## Stage Report: ideation — SUPERSEDED (historical cache design)

- DONE: Resolve binary persistence across Cowork session relaunches.
  SUPERSEDED: this cycle chose a connected-folder binary/manifest cache; the next historical cycle removed it, and the current operative design instead uses the captain-specified exact project-local `.spacedock/bin/spacedock` file without a cache manifest/pointer.
- DONE: Clarify existing survey behavior versus Cowork-specific sampling/report behavior.
  The analysis step now labels retained direct-report, evidence, mode, offer, and conservative-frontier rules separately from Cowork bounds, extraction, sampling labels, privacy generalization, and omissions.
- DONE: Shorten the proposed Cowork documentation into a new paragraph.
  The concrete replacement is now one local paragraph plus one short Cowork paragraph.
- DONE: Make Full network access plus session relaunch the only user-facing network instruction.
  The exact prompt and docs name Settings, Full network access, relaunch, and rerun; domains remain internal to the request ledger and the live proof is explicitly two-phase.

### Summary

SUPERSEDED HISTORICAL SUMMARY: this cycle proposed restoring from a connected-folder cache. That manifest/pointer cache does not ship; the operative design uses the single exact project-local `.spacedock/bin/spacedock` identity. Sampling and user-network wording remain historical inputs only where retained by the operative body.

## Stage Report: ideation — SUPERSEDED (historical no-cache/session-only design)

- DONE: Narrow the persistence premise and bind it to lane-specific live evidence.
  Official evidence is now limited to temporary remote sandboxes/multiple modes; run 0 plus runs 1-3 record only mode class/booleans and empirically characterize HOME reset for the selected lane.
- DONE: Define real connected-folder discovery and deterministic selection, or remove the unsupported mechanism.
  No root-enumeration/working-directory contract was found, so the design explicitly performs no discovery, selection, writability probe, cache read, or connected-folder write and uses honest per-environment reinstall.
- SKIPPED: Define cache schema and operations precisely.
  The cache is removed because safe discovery and replacement cannot be guaranteed; immutable-entry schema, candidate selection, noexec copy, and partial-write handling are therefore not shipped or claimed.
- DONE: Resolve the read-only contract.
  Local agentsview and both evidence scans remain project/session-source read-only; Cowork's disclosed exception writes only temporary/session-HOME helper files after the settings relaunch.
- DONE: Strengthen live persistence proof.
  Run 0 proves the zero-network stop; run 1 installs/verifies before creating the sentinel; runs 2/3 inspect both booleans before install, prove reinstall rather than restore on reset, and execute only the exact HOME path.

### Summary

SUPERSEDED HISTORICAL SUMMARY: this cycle removed all mounted-folder persistence and chose exact-HOME reuse/reinstall. The captain's later correction replaces that design with the host-supplied current project folder and exact `.spacedock/bin/spacedock` persistence. Its no-arbitrary-root and privacy constraints remain retained where stated in the operative body.

## Stage Report: ideation — SUPERSEDED (historical session-only consistency design)

- DONE: Preserve local agentsview's read-only and ask-before-install behavior in proposed docs.
  The local paragraph now states both guarantees; Cowork's session-local helper exception remains in its separate short paragraph and connected files stay unchanged.
- DONE: Make the live lifecycle phases unambiguous.
  Run 0 is the missing-binary zero-network stop; run 1 installs/verifies before sentinel creation; runs 2/3 inspect both booleans before install and classify reset versus persistence exactly.
- DONE: Mark historical cache claims as superseded without erasing workflow history.
  At that cycle the cache-choosing report was superseded by no-cache; this entire session-only cycle is now itself SUPERSEDED by the captain's exact project-local binary correction.
- DONE: Rename the session-local binary identity.
  Operative wording now calls `$HOME/.local/bin/spacedock` the exact Cowork execution identity, not a durable install identity.

### Summary

SUPERSEDED HISTORICAL SUMMARY: this cycle aligned the session-only/no-cache design. The captain's later correction replaces its HOME execution identity and run-0-to-3 lifecycle with persistent `<current-project-folder>/.spacedock/bin/spacedock` and a fresh-HOME run-2 reuse proof. Local read-only, consent, and historical-labeling fixes remain retained.

## Stage Report: ideation

- DONE: Bind persistence to Cowork's mounted current project folder.
  The operative design consumes only the host-supplied current-project-folder capability, records generic availability/binding classes, and never enumerates roots or logs a path/name.
- DONE: Make `.spacedock/bin/spacedock` the exact persistent execution identity.
  Absent enters consent; working executes directly with zero network; non-executable, failing, dangling-symlink, and PATH-shadow cases stop or ignore the shadow exactly as specified.
- DONE: Define deletion-denied first creation and honest failure boundaries.
  Run 1 verifies in session space, reserves the absent final name with noclobber, copies/chmods/verifies without unlink or rename, and treats any partial final state as broken pending explicit repair.
- DONE: Prove fresh-HOME reuse and constrain project mutation.
  Run 2 uses a fresh HOME plus the same mount, executes the mounted binary with zero network/install events, and byte-level evidence permits only `.spacedock/bin/spacedock` and absent parent-directory creation.
- DONE: Update read-only, docs, consent, privacy, and historical consistency.
  Local/history reads stay read-only; project-local helper storage is disclosed separately from short network guidance; no sensitive identifiers enter evidence; both immediately preceding session-only reports are marked SUPERSEDED.

### Summary

The captain's project-local correction is now the operative persistence design: Cowork installs once at the exact mounted `.spacedock/bin/spacedock` path and later fresh-HOME sessions reuse it directly without network. First creation requires no unlink or rename and makes no atomicity claim; broken or partial state is never overwritten without a separate repair decision.
