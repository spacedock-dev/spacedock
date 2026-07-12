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

### Cowork execution and reinstall model

Narrow evidence boundary: Anthropic documents isolated **remote** Cowork sandboxes as temporary and documents both remote and local execution modes ([Cowork safety](https://support.claude.com/en/articles/13364135-use-claude-cowork-safely), [architecture](https://support.claude.com/en/articles/14479288-claude-cowork-architecture-overview)). It does not establish that every mode resets `$HOME`, or that a connected folder is Cowork's only durable storage. For the selected/probed lane, a connected folder is only a candidate shell-readable, user-controlled cross-session boundary; persistence must be established by live behavior, not generalized from the documentation.

The available Cowork session-info capabilities enumerate sessions and transcripts, not connected-folder roots. The supplied live research does not identify a root-enumeration tool or a stable shell working-directory convention. `pwd` is insufficient: it cannot prove whether the directory is scratch or a mounted root and cannot enumerate multiple roots. A create-only writability probe would leave an artifact on a deletion-denied mount; a create/delete probe is not non-destructive there. Therefore survey does **not** discover, select, probe, or write a connected folder, and it does not create `.spacedock/cache/`. This avoids invented root selection, unsafe partial metadata, noexec assumptions, and replacement on mounts that cannot unlink.

Choose honest per-environment installation: `$HOME/.local/bin/spacedock` is the exact binary for the current Cowork execution environment. When that environment resets and the path disappears, survey repeats the settings/relaunch/install flow. The live probe records execution mode as only `remote`, `local`, or `unknown` when the host exposes it; never infer mode and never record a path. Lifecycle phases are fixed: **run 0** observes the missing exact binary, records zero network events, presents the Settings/Full-network/relaunch instruction, and stops. **Run 1** is the post-settings environment: install and verify the exact binary first, then create/record a non-sensitive random sentinel in that same HOME. **Runs 2 and 3** inspect the run-1 sentinel and exact path before any install action. Classify reset only when both values are absent in both later runs; each reset run reinstalls. Classify persistence when the exact path survives and direct invocation succeeds with zero network. Record only mode class and the two presence booleans; do not claim a cross-mode rule.

On the Cowork branch:

1. The exact Cowork execution identity is `$HOME/.local/bin/spacedock`, never whichever `spacedock` resolves later on `PATH`. Execute this probe block as written; `COWORK_INSTALL_ABSENT` enters step 2, `COWORK_INSTALL_BROKEN` stops, and silent success reuses the binary:

   ```sh
   COWORK_BIN="$HOME/.local/bin/spacedock"
   if [ -e "$COWORK_BIN" ] || [ -L "$COWORK_BIN" ]; then
     if ! "$COWORK_BIN" --version >/dev/null 2>&1; then echo "COWORK_INSTALL_BROKEN"; fi
   else
     echo "COWORK_INSTALL_ABSENT"
   fi
   ```

   The `-L` arm classifies a dangling symlink as broken rather than absent. Success reuses the exact file with no network activity. Non-executable or failing existing content stops; do not search PATH, download, or overwrite it. Absence proceeds to consent even when another valid `spacedock` exists on PATH.
2. When the exact path is absent, display the exact prompt below and stop. The user changes Cowork Settings to **Full network access**, relaunches the Cowork session, and runs survey again; nothing is downloaded before that relaunch/rerun. Before consent there must be no `WebFetch`, `WebSearch`, `curl`, GitHub API call, installer fetch, release fetch, or network-capability test. Do not search for a cache or touch connected folders.
3. The post-relaunch survey rerun is the affirmative continuation. Allocate a unique installer path with `INSTALLER_TMP=$(mktemp "${TMPDIR:-/tmp}/spacedock-install.XXXXXX")`; abort without networking if allocation fails. Register `trap 'rm -f "$INSTALLER_TMP" 2>/dev/null || true' EXIT HUP INT TERM`, acquire with `curl -fsSL -o "$INSTALLER_TMP" https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh`, then run `SPACEDOCK_INSTALL_DIR="$HOME/.local/bin" sh "$INSTALLER_TMP"`. Attempt `rm -f "$INSTALLER_TMP"` on the success path and clear the trap when it succeeds. On a deletion-blocked temp filesystem, report only that the temporary installer could not be removed; do not echo its path or claim cleanup succeeded. Internally, the transport allowlist and request ledger cover `raw.githubusercontent.com`, `api.github.com`, `github.com`, and `release-assets.githubusercontent.com`; these details do not belong in the user prompt or docs. Do not send credentials and do not try alternate mirrors. Verify the installed identity directly with `"$HOME/.local/bin/spacedock" --version`; only afterward may `PATH="$HOME/.local/bin:$PATH"` be exported for command convenience. The branch writes only VM/session temporary files and `$HOME/.local/bin/spacedock`, never project or connected-folder content.
4. Claim only the installer's shipped failure boundary: acquisition, unsupported target, download, missing checksum, checksum mismatch, and archive-extraction failures occur before it writes the final install path, so a **fresh, absent** `$HOME/.local/bin/spacedock` remains absent. After checksum/extraction it writes the final path directly; the task makes no atomic replacement or preservation claim. Because this flow refuses to overwrite an existing path, overwrite-failure preservation is out of scope.
5. Reuse the successful inventory probe. Sort idle sessions by the inventory's most-recent timestamp descending (stable tie-break: opaque id, used only in memory), omit active sessions, and sample at most the newest 20 idle sessions. For each sampled session, call `read_transcript(max_wait_seconds: 0, limit: 15)` and request the newest page. If the returned schema/result exposes both a continuation cursor and `truncated=true`, fetch at most two more 15-message pages for that session (45 messages maximum); otherwise do not invent pagination. Stop paging on no cursor, `truncated=false`, repeat cursor, or the 45-message ceiling. Tool errors omit that session from analyzed counts and increment a generic read-failure count.
6. Analyze only the returned window. **Retained existing survey behavior:** render directly without a scratch preamble; fill every displayed value from evidence; never invent an empty signal; keep the existing `manual` / `exploration` / `knowledge-work` / `unlabeled` meanings and their mode-keyed commission offers; keep conservative open-frontier wording when artifact evidence is unavailable. **Cowork-specific adapter behavior:** the 20-session/15-message/three-page bounds, sampling disclosure, message-level decision/interruption extraction, `sample-unverified` label, privacy-generalized cluster names, and repo-only section omissions. These are host adaptations, not changes to the local agentsview path.

   For the Cowork adapter, a **decision** is an observed assistant question offering alternatives or requesting approval; it is resolved only when a later sampled user message answers it, otherwise `OPEN (sample-unverified)`. A **hand-steering interruption** is a sampled user turn that corrects, rejects, or redirects work after an agent proposal/output; initial requests and simple affirmative answers do not count, and one user turn counts at most once. Apply the existing mode meanings to observed messages: repeated execute/build-and-review loops -> `manual`; iterative rejection/revision of creative or design paths -> `exploration`; repeated intake/process/commentary -> `knowledge-work`; insufficient or mixed evidence -> `unlabeled`. Cluster labels must be generic activity categories, never copied titles, names, or identifying attributes.
7. Project the report explicitly. Keep `WHAT THIS GIVES YOU`, sampled `BY THE NUMBERS`, `HOW YOU WORK`, `THREADS TO PULL`, `RECENT DECISIONS`, `WORKSTREAMS`, `WHAT THIS CAN'T SEE`, and the mode-keyed commission offer. Omit `BACKLOG`, `WORK BY AREA`, `CODEX`, scaffold comparison, dispatch/subagent facts, no-follow-up counts, shipped claims, and git/PR counts because Cowork session evidence cannot calculate them. Every open thread is labeled `sample-unverified`; no section may imply a repo cross-check occurred. The title is `SpaceDock survey — Cowork session sample`, followed by `sampled {analyzed} of {idle_total} idle sessions (newest first; up to 45 messages/session; {read_failures} read failures; active sessions excluded)`. If inventory timestamps are unavailable, omit a day-span rather than estimating one.
8. Fill every displayed value from the bounded scan, disclose sampling/truncation/read failures, show `none observed in this sample` for an applicable empty signal, and emit no literal `{slot}`. Preserve the end-of-report commission offer, scoped to the sampled evidence.

### Exact survey instruction change

Before (current opening route):

> Run the four steps in order: **check agentsview → scan → recognize scaffold → report and offer**.

After:

> Run `ToolSearch(query="select:mcp__session_info__list_sessions,mcp__session_info__read_transcript", max_results=2)` so lazy-loaded schemas cannot false-negative. If both exact tools resolve, call `mcp__session_info__list_sessions` once. A successful call selects **check exact session Spacedock execution → ask before network → bounded Cowork scan → sampled Cowork report and offer**; reuse its inventory and never enter agentsview or repo probes. A missing definition selects the existing local path unchanged. If both definitions resolve but inventory fails, report `Cowork session inventory unavailable` and stop; do not misroute into local history. Missing files or binaries alone never identify Cowork.

Insert this exact permission prompt when the Cowork branch cannot invoke Spacedock:

> Spacedock is not available in this Cowork environment. In Cowork Settings, choose **Full network access**, then relaunch this session and run survey again. Nothing will be downloaded before you relaunch.

### Privacy boundary

Session identifiers are in-memory tool handles only. Never echo or persist identifiers, session titles, raw tool payloads, transcript text, proper nouns, account details, credentials, tokens, machine-specific paths, or paraphrases/redactions that preserve a private person's, organization’s, project's, or event's identifying attributes. User-visible workstream and decision labels must be generic activity categories; durable evidence is limited to tool names, result classes, ordering, counts, sampling metadata, generic categories, and synthetic fixture values. Privacy tests inject raw and transformed canaries (case-folded, punctuation-stripped, partial-token, title-derived, and paraphrased-attribute forms) and reject every variant in output and durable files. Do not write raw Cowork responses to disk or pass any session material to the installer or a network call.

### Read-only boundary

The existing local agentsview survey remains strictly read-only with respect to the project and session sources. The Cowork evidence scan is also read-only. Its consented bootstrap is the narrow exception: after the settings relaunch it may create a temporary installer and `$HOME/.local/bin/spacedock` inside that Cowork execution environment. It never writes project/connected-folder content and creates no persistent cache. The report must call this a session-local helper installation, not a project change or durable install.

### Documentation change proposed for implementation

In `docs/site/get-started/survey.md`, replace the local-only opening paragraph:

> When you use the `spacedock:survey` skill, it looks at your existing agent conversation logs on local disk (through agentsview, an open source session-history tool). It is read-only; if agentsview is missing, it asks before installing it.

with these two paragraphs:

> When you use the `spacedock:survey` skill, it reads the session history available to the current host. Local coding-agent sessions use agentsview; the local scan is read-only, and if agentsview is missing, Survey asks before installing it.
>
> In Claude Cowork, Survey reads a bounded session sample and does not change connected files. If Spacedock is unavailable, set Cowork to **Full network access** in Settings, relaunch, and run survey again; Survey may reinstall its session-local helper when Cowork resets the environment.

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

**AC-2 - A working exact session binary resumes Cowork survey without network activity; an absent binary produces the Full-network/settings/relaunch instruction and no cache or connected-folder mutation.**
Verified by: run 0 of an event-ledger test with a fake `$HOME`, read-only connected-folder fixture, and transport. Exact-path absence emits the prompt and stops with 0 requests, 0 cache/discovery events, and no connected-folder changes. A persistent working exact-path mutant yields one direct invocation and 0 requests; shadow-PATH mutants retain their absent/broken guarantees.

**AC-3 - After the settings relaunch, a fresh Cowork installation verifies the architecture-matched release, invokes only the exact session path, and is honestly repeated when the selected lane resets HOME.**
Verified by: the real installer seams plus runs 1-3 after run 0's stop. Run 1 installs and directly verifies the exact path, then creates/records a random HOME sentinel. Before any install action, runs 2 and 3 record only mode class, `home-sentinel-present`, and `exact-binary-present`. Reset requires both booleans false in both runs and each run reinstalls; persistence requires direct exact-path success with 0 network. Unique-`mktemp`, cleanup, failure-boundary, and shadow-PATH tests remain. A read-only connected-folder snapshot is byte-identical throughout.

**AC-4 - Cowork survey renders a complete, honestly sampled orientation with correct filled values and no repository-only claims.**
Verified by: a deterministic replay inventory with four idle sessions and one active session. Its newest-first transcript pages produce exactly `4 of 4` analyzed, `2` hand-steering interruptions, `1` resolved decision, `1` `OPEN (sample-unverified)`, four workstreams (`manual`, `exploration`, `knowledge-work`, `unlabeled`), `0` read failures, and the `up to 45 messages/session; active sessions excluded` disclosure. The output contains no literal `{slot}` and omits BACKLOG, WORK BY AREA, CODEX, scaffold, dispatch, no-follow-up, shipped, and git/PR sections. Mutant cases cover 21 idle sessions (20 sampled), truncation with one/two continuation pages, repeated cursor, missing pagination, and one read failure.

**AC-5 - Cowork output and durable artifacts contain only aggregate/synthetic evidence, not raw or transform-preserving private session material.**
Verified by: canary identifiers, titles, transcript sentences, proper nouns, account/token/path values, and identifying attributes injected into the in-memory harness; recursive assertions reject exact, case-folded, punctuation-stripped, partial-token, title-derived, and paraphrased-attribute canaries from output and files, while expected aggregate counts, sampling disclosure, and generic activity categories remain. The harness also asserts no raw tool response is written.

**AC-6 - User documentation preserves local survey's read-only promise and briefly discloses Cowork's session-local helper reinstall without implying connected-folder writes or durable caching.**
Verified by: render and inspect the Survey page against the exact proposed replacement; runtime replay and live smoke remain the behavioral proof rather than a prose substring test.

## Test plan

Implementation starts with the fixture harness so routing, consent, privacy, and rendered-report assertions fail against the current skill before its instructions change.

- **Cowork event replay (medium):** add `TestSurveyCoworkEventReplay` and the committed sanitized `skills/integration/testdata/survey/cowork-tool-events.jsonl`. The runner consumes tool events and synthetic responses, enforcing the exact ToolSearch -> inventory -> bounded transcript sequence and rendering the Cowork projection. Fixture/mutant expectations are AC-1/AC-4's exact values. The existing `TestSurveyInstallProbe` remains the executable local-path control; no generic “skill smoke” is assumed.
- **Installer/consent CLI tests (medium, existing installer seams):** wrap the exact-path probe and post-relaunch acquisition commands with a fake `$HOME`, stub `curl` request ledger, stub `uname`, and fixture assets/checksums. Execute the real `install.sh`; assert the short settings/relaunch prompt and zero pre-relaunch requests, the internal four-domain allowlist after rerun, direct exact-path verification, architecture mapping, fresh-install success, and documented pre-final-write failures. Retain both shadow-PATH mutants and unique-`mktemp`/cleanup cases. Assert zero connected-folder discovery/write events.
- **Environment-reset characterization (medium):** run 0 proves missing exact binary -> 0 requests -> Settings/Full-network/relaunch stop. Run 1, after settings, installs and directly verifies the exact path before creating the random HOME sentinel. Runs 2/3 inspect sentinel and exact binary before any install event. Record only mode class (`remote|local|unknown`) and the two booleans. Reset requires both false in both runs and leads to reinstall; persistence requires direct exact-path reuse with 0 network. A synthetic connected-folder tree remains byte-identical and no cache operations occur.
- **Report/privacy replay (medium):** render the exact four-session baseline and pagination/failure mutants from synthetic in-memory payloads. Assert filled values, no placeholders, sampling disclosure, omission matrix, open-thread label, mode mapping, and transformed-canary absence across stdout and the test temp tree.
- **Documentation render (low):** build the docs site and visually/readably inspect the changed Survey page.
- **Live Cowork smoke (required, medium/manual host harness):** run 0 executes ToolSearch/inventory, records mode class when exposed, observes exact-binary absence, and stops at the settings/relaunch prompt with 0 requests. Run 1 is post-settings: install and verify the exact HOME path first, then create/record the random sentinel, then scan. Runs 2/3 record the two presence booleans before any install event. Reset requires both false twice and proves reinstall; persistence directly reuses the exact path with 0 network. Evidence also records `connected-folder-discovery: unsupported/not-used`, `connected-folder-handle: none`, `cache-read-event: none`, `zero-network-restore: not-applicable`, and `copy-from-cache: none`. No path, folder name, payload, or identifying handle is retained, and no command/file event targets connected-folder content.

## Risk spike

The focused review invalidates the cache premise: neither official guidance nor the supplied schema provides a general shell root-discovery/selection contract. No cache spike is warranted because the mechanism is removed. The remaining runtime uncertainty is lane-specific HOME persistence, so run 0 plus the run-1-to-3 sentinel characterization is the first live implementation gate. Its result determines reuse versus honest reinstall without changing the connected-folder read-only boundary; replay cannot substitute for that live lifecycle result.

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

## Stage Report: ideation

- DONE: Replace the PATH-based install identity probe with explicit execution-path semantics.
  The design invokes and verifies only `"$HOME/.local/bin/spacedock"`; absent, working, and broken exact-path states cannot be shadowed by PATH.
- DONE: Add two executable test mutants for later-PATH binaries.
  AC-2 and the installer test plan require absent-exact-path/valid-PATH consent and broken-exact-path/valid-PATH stop cases with invocation and network ledgers.
- DONE: Correct AC-1's heading to distinguish missing definitions from resolved inventory errors.
  The heading now says missing definitions retain local behavior while resolved inventory errors stop, matching the body and verifier.
- DONE: Replace the fixed temporary installer filename with a unique `mktemp` path and require honest cleanup behavior.
  The acquisition contract uses a unique path, traps cleanup, proves normal deletion, and reports without overclaiming when deletion is blocked.

### Summary

This surgical repair makes the Cowork binary identity the exact session-local execution path rather than PATH resolution and adds executable shadow-PATH mutants for both absent and broken states. Installer acquisition is uniquely named and cleanup is verified or honestly reported, while the previously resolved capability, report, and privacy contracts remain intact.

## Stage Report: ideation — SUPERSEDED (historical cache design)

- DONE: Resolve binary persistence across Cowork session relaunches.
  SUPERSEDED: this cycle chose a connected-folder binary/manifest cache; the later no-cache report and operative body remove that mechanism because discovery and safe updates were not proven.
- DONE: Clarify existing survey behavior versus Cowork-specific sampling/report behavior.
  The analysis step now labels retained direct-report, evidence, mode, offer, and conservative-frontier rules separately from Cowork bounds, extraction, sampling labels, privacy generalization, and omissions.
- DONE: Shorten the proposed Cowork documentation into a new paragraph.
  The concrete replacement is now one local paragraph plus one short Cowork paragraph.
- DONE: Make Full network access plus session relaunch the only user-facing network instruction.
  The exact prompt and docs name Settings, Full network access, relaunch, and rerun; domains remain internal to the request ledger and the live proof is explicitly two-phase.

### Summary

SUPERSEDED HISTORICAL SUMMARY: this cycle proposed restoring from a connected-folder cache. That proposal does not ship; the later no-cache design uses exact session-local reuse or honest reinstall and never reads or writes a connected-folder cache. The sampling and user-network wording from this cycle remain historical inputs only where retained by the operative body.

## Stage Report: ideation

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

This repair removes the under-specified connected-folder cache instead of inventing discovery and safe-update semantics that Cowork does not expose. The operative design now characterizes each live lane through run 0 plus runs 1-3, reuses a surviving exact HOME binary or honestly reinstalls it after reset, and preserves connected folders as read-only while keeping routing, report, consent, docs brevity, and privacy behavior intact.

## Stage Report: ideation

- DONE: Preserve local agentsview's read-only and ask-before-install behavior in proposed docs.
  The local paragraph now states both guarantees; Cowork's session-local helper exception remains in its separate short paragraph and connected files stay unchanged.
- DONE: Make the live lifecycle phases unambiguous.
  Run 0 is the missing-binary zero-network stop; run 1 installs/verifies before sentinel creation; runs 2/3 inspect both booleans before install and classify reset versus persistence exactly.
- DONE: Mark historical cache claims as superseded without erasing workflow history.
  The cache-choosing report heading, evidence, and summary are explicitly SUPERSEDED, and the later no-cache report/body remain operative.
- DONE: Rename the session-local binary identity.
  Operative wording now calls `$HOME/.local/bin/spacedock` the exact Cowork execution identity, not a durable install identity.

### Summary

This consistency repair preserves the no-cache design while making documentation, lifecycle evidence, and identity terminology agree. Historical cache evidence remains visible but unmistakably superseded; the operative proof now has a precise run-0/run-1/run-2/run-3 sequence using only mode class and presence booleans.
