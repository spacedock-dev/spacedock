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

On the Cowork branch:

1. The durable Cowork install identity is exactly `$HOME/.local/bin/spacedock`, never whichever `spacedock` resolves later on `PATH`. Execute this probe block as written; `COWORK_INSTALL_ABSENT` enters step 2, `COWORK_INSTALL_BROKEN` stops, and silent success reuses the binary:

   ```sh
   COWORK_BIN="$HOME/.local/bin/spacedock"
   if [ -e "$COWORK_BIN" ] || [ -L "$COWORK_BIN" ]; then
     if ! "$COWORK_BIN" --version >/dev/null 2>&1; then echo "COWORK_INSTALL_BROKEN"; fi
   else
     echo "COWORK_INSTALL_ABSENT"
   fi
   ```

   The `-L` arm classifies a dangling symlink as broken rather than absent. Success reuses the exact file with no network activity. Non-executable or failing existing content stops; do not search PATH, download, or overwrite it. Absence proceeds to consent even when another valid `spacedock` exists on PATH.
2. When the exact path is absent, display the exact prompt below and stop. Before affirmative consent there must be no `WebFetch`, `WebSearch`, `curl`, GitHub API call, installer fetch, release fetch, or network-capability test.
3. After affirmation, allocate a unique installer path with `INSTALLER_TMP=$(mktemp "${TMPDIR:-/tmp}/spacedock-install.XXXXXX")`; abort without networking if allocation fails. Register `trap 'rm -f "$INSTALLER_TMP" 2>/dev/null || true' EXIT HUP INT TERM`, acquire with `curl -fsSL -o "$INSTALLER_TMP" https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh`, then run `SPACEDOCK_INSTALL_DIR="$HOME/.local/bin" sh "$INSTALLER_TMP"`. Attempt `rm -f "$INSTALLER_TMP"` on the success path and clear the trap when it succeeds. On a deletion-blocked mount, report only that the temporary installer could not be removed (do not echo its machine-specific path); do not claim cleanup succeeded. The complete allowlist disclosed before consent is: `raw.githubusercontent.com` for installer acquisition, `api.github.com` for latest-release resolution, `github.com` for the tarball and `checksums.txt`, and the redirect host `release-assets.githubusercontent.com`. Do not send credentials and do not try alternate mirrors. The installer owns `uname -s`/`uname -m` normalization (`Linux` plus `x86_64|amd64` -> `linux_amd64`; `arm64|aarch64` -> `linux_arm64`), release selection, and checksum verification. Verify the installed identity directly with `"$HOME/.local/bin/spacedock" --version`; only afterward may `PATH="$HOME/.local/bin:$PATH"` be exported for command convenience.
4. Claim only the installer's shipped failure boundary: acquisition, unsupported target, download, missing checksum, checksum mismatch, and archive-extraction failures occur before it writes the final install path, so a **fresh, absent** `$HOME/.local/bin/spacedock` remains absent. After checksum/extraction it writes the final path directly; the task makes no atomic replacement or preservation claim. Because this flow refuses to overwrite an existing path, overwrite-failure preservation is out of scope.
5. Reuse the successful inventory probe. Sort idle sessions by the inventory's most-recent timestamp descending (stable tie-break: opaque id, used only in memory), omit active sessions, and sample at most the newest 20 idle sessions. For each sampled session, call `read_transcript(max_wait_seconds: 0, limit: 15)` and request the newest page. If the returned schema/result exposes both a continuation cursor and `truncated=true`, fetch at most two more 15-message pages for that session (45 messages maximum); otherwise do not invent pagination. Stop paging on no cursor, `truncated=false`, repeat cursor, or the 45-message ceiling. Tool errors omit that session from analyzed counts and increment a generic read-failure count.
6. Analyze only the returned window. A **decision** is an observed assistant question offering alternatives or requesting approval; it is resolved only when a later sampled user message answers it, otherwise `OPEN (sample-unverified)`. A **hand-steering interruption** is a sampled user turn that corrects, rejects, or redirects work after an agent proposal/output; initial requests and simple affirmative answers do not count, and one user turn counts at most once. Mode mapping retains survey's existing classes using observed evidence: repeated execute/build-and-review loops -> `manual`; iterative rejection/revision of creative or design paths -> `exploration`; repeated intake/process/commentary -> `knowledge-work`; insufficient or mixed evidence -> `unlabeled`. Cluster labels must be generic activity categories, never copied titles, names, or identifying attributes.
7. Project the report explicitly. Keep `WHAT THIS GIVES YOU`, sampled `BY THE NUMBERS`, `HOW YOU WORK`, `THREADS TO PULL`, `RECENT DECISIONS`, `WORKSTREAMS`, `WHAT THIS CAN'T SEE`, and the mode-keyed commission offer. Omit `BACKLOG`, `WORK BY AREA`, `CODEX`, scaffold comparison, dispatch/subagent facts, no-follow-up counts, shipped claims, and git/PR counts because Cowork session evidence cannot calculate them. Every open thread is labeled `sample-unverified`; no section may imply a repo cross-check occurred. The title is `SpaceDock survey — Cowork session sample`, followed by `sampled {analyzed} of {idle_total} idle sessions (newest first; up to 45 messages/session; {read_failures} read failures; active sessions excluded)`. If inventory timestamps are unavailable, omit a day-span rather than estimating one.
8. Fill every displayed value from the bounded scan, disclose sampling/truncation/read failures, show `none observed in this sample` for an applicable empty signal, and emit no literal `{slot}`. Preserve the end-of-report commission offer, scoped to the sampled evidence.

### Exact survey instruction change

Before (current opening route):

> Run the four steps in order: **check agentsview → scan → recognize scaffold → report and offer**.

After:

> Run `ToolSearch(query="select:mcp__session_info__list_sessions,mcp__session_info__read_transcript", max_results=2)` so lazy-loaded schemas cannot false-negative. If both exact tools resolve, call `mcp__session_info__list_sessions` once. A successful call selects **check durable Spacedock install → ask before network → bounded Cowork scan → sampled Cowork report and offer**; reuse its inventory and never enter agentsview or repo probes. A missing definition selects the existing local path unchanged. If both definitions resolve but inventory fails, report `Cowork session inventory unavailable` and stop; do not misroute into local history. Missing files or binaries alone never identify Cowork.

Insert this exact permission prompt when the Cowork branch cannot invoke Spacedock:

> Spacedock is not available in `~/.local/bin`. To continue, I need network access to download the installer from `raw.githubusercontent.com`, resolve the latest release through `api.github.com`, and fetch its checksum and architecture-matched Linux asset through `github.com` / `release-assets.githubusercontent.com`. The verified binary will be written to `~/.local/bin`; no credentials or session data will be sent. Please enable network access for this Cowork session and tell me when it is enabled. I will not probe the network, call those hosts, or download anything until you confirm.

### Privacy boundary

Session identifiers are in-memory tool handles only. Never echo or persist identifiers, session titles, raw tool payloads, transcript text, proper nouns, account details, credentials, tokens, machine-specific paths, or paraphrases/redactions that preserve a private person's, organization’s, project's, or event's identifying attributes. User-visible workstream and decision labels must be generic activity categories; durable evidence is limited to tool names, result classes, ordering, counts, sampling metadata, generic categories, and synthetic fixture values. Privacy tests inject raw and transformed canaries (case-folded, punctuation-stripped, partial-token, title-derived, and paraphrased-attribute forms) and reject every variant in output and durable files. Do not write raw Cowork responses to disk or pass any session material to the installer or a network call.

### Documentation change proposed for implementation

In `docs/site/get-started/survey.md`, replace the local-only opening paragraph:

> When you use the `spacedock:survey` skill, it looks at your existing agent conversation logs on local disk (through agentsview, an open source session-history tool). It is read-only; if agentsview is missing, it asks before installing it.

with:

> When you use the `spacedock:survey` skill, it reads the session history available to the current host. Local coding-agent sessions use agentsview; Claude Cowork uses a bounded, newest-first sample from Cowork's session inventory and transcript tools. In Cowork, Survey checks exactly `~/.local/bin/spacedock` and never substitutes another PATH copy: a working file is reused, a broken file stops for repair, and an absent file leads to the consent prompt. Before downloading a uniquely named temporary installer, it discloses `raw.githubusercontent.com`, `api.github.com`, `github.com`, and `release-assets.githubusercontent.com`, then waits for your confirmation; it attempts to remove the temporary installer and reports when the mount blocks deletion. The report discloses its session/message limits and read failures, omits claims that require repository evidence, labels open threads sample-unverified, and keeps raw titles, transcripts, identifiers, and identifying attributes out of output and durable files.

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

**AC-2 - A durable working `~/.local/bin/spacedock` is reused without network activity; an absent install produces the specified complete host disclosure and zero network attempts before affirmative permission.**
Verified by: an event-ledger test with a fake `$HOME` and transport. A working exact-path stub yields one direct version invocation and 0 requests; absence or denial yields 0 requests and no file; affirmation permits requests only to `raw.githubusercontent.com`, `api.github.com`, `github.com`, and `release-assets.githubusercontent.com`, in the installer flow. Executable mutants prove (a) exact path absent plus a valid later-PATH binary still enters consent/install without invoking the PATH binary, and (b) exact path broken/non-executable plus a valid later-PATH binary stops broken with neither PATH-binary invocation nor network/overwrite.

**AC-3 - After permission, a fresh Cowork installation uses the release asset matching the observed Linux architecture, verifies its checksum, writes the configured user path, and successfully invokes the resulting binary.**
Verified by: the real `install.sh` inspection seam for `x86_64/amd64` and `arm64/aarch64`, fixture tarballs/checksums, and a fake `$HOME`. Acquisition/target/download/checksum/extraction failures assert an initially absent final binary remains absent; success asserts `"$HOME/.local/bin/spacedock" --version`. Unique-`mktemp` acquisition tests assert distinct concurrent paths, cleanup on a normal temp filesystem, and an honest generic cleanup warning with the file retained on a deletion-denied fixture. No atomic replacement or existing-binary preservation assertion is made.

**AC-4 - Cowork survey renders a complete, honestly sampled orientation with correct filled values and no repository-only claims.**
Verified by: a deterministic replay inventory with four idle sessions and one active session. Its newest-first transcript pages produce exactly `4 of 4` analyzed, `2` hand-steering interruptions, `1` resolved decision, `1` `OPEN (sample-unverified)`, four workstreams (`manual`, `exploration`, `knowledge-work`, `unlabeled`), `0` read failures, and the `up to 45 messages/session; active sessions excluded` disclosure. The output contains no literal `{slot}` and omits BACKLOG, WORK BY AREA, CODEX, scaffold, dispatch, no-follow-up, shipped, and git/PR sections. Mutant cases cover 21 idle sessions (20 sampled), truncation with one/two continuation pages, repeated cursor, missing pagination, and one read failure.

**AC-5 - Cowork output and durable artifacts contain only aggregate/synthetic evidence, not raw or transform-preserving private session material.**
Verified by: canary identifiers, titles, transcript sentences, proper nouns, account/token/path values, and identifying attributes injected into the in-memory harness; recursive assertions reject exact, case-folded, punctuation-stripped, partial-token, title-derived, and paraphrased-attribute canaries from output and files, while expected aggregate counts, sampling disclosure, and generic activity categories remain. The harness also asserts no raw tool response is written.

**AC-6 - User documentation describes capability-gated Cowork evidence, durable-install reuse, all pre-consent network hosts, bounded sampling, privacy, and omission of unsupported repo conclusions.**
Verified by: render and inspect the Survey page against the exact proposed replacement; runtime replay and live smoke remain the behavioral proof rather than a prose substring test.

## Test plan

Implementation starts with the fixture harness so routing, consent, privacy, and rendered-report assertions fail against the current skill before its instructions change.

- **Cowork event replay (medium):** add `TestSurveyCoworkEventReplay` and the committed sanitized `skills/integration/testdata/survey/cowork-tool-events.jsonl`. The runner consumes tool events and synthetic responses, enforcing the exact ToolSearch -> inventory -> bounded transcript sequence and rendering the Cowork projection. Fixture/mutant expectations are AC-1/AC-4's exact values. The existing `TestSurveyInstallProbe` remains the executable local-path control; no generic “skill smoke” is assumed.
- **Installer/consent CLI tests (medium, existing installer seams):** wrap the exact-path durable-install probe and post-consent acquisition commands with a fake `$HOME`, stub `curl` request ledger, stub `uname`, and fixture assets/checksums. Execute the real `install.sh`; assert all four allowed domains, zero pre-consent requests, direct exact-path post-install verification, architecture mapping, fresh-install success, and only the documented pre-final-write fail-closed cases. Add two shadow-PATH mutants: exact path absent + valid PATH binary must consent/install without invoking the shadow; exact path broken/non-executable + valid PATH binary must stop broken without invoking the shadow, network, or overwrite. Run two acquisitions to prove unique `mktemp` names; assert removal where supported and a generic, path-free warning plus retained file where deletion is denied.
- **Report/privacy replay (medium):** render the exact four-session baseline and pagination/failure mutants from synthetic in-memory payloads. Assert filled values, no placeholders, sampling disclosure, omission matrix, open-thread label, mode mapping, and transformed-canary absence across stdout and the test temp tree.
- **Documentation render (low):** build the docs site and visually/readably inspect the changed Survey page.
- **Live Cowork smoke (required, medium/manual host harness):** run `TestLiveClaudeCoworkSurveyProbe` in a fresh Cowork session with a synthetic/non-sensitive idle session. It executes ToolSearch, the inventory probe, consent stop, installer acquisition/install, bounded transcript read, and report render; then exports only the sanitized event shape used by replay. Durable evidence is limited to exit/status classes, exact tool/domain names, event ordering, counts, selected `linux_{arm64|amd64}` token, binary version, sampling metadata, generic categories, and privacy-canary absence.

## Risk spike

No separate throwaway product spike is needed before implementation. Installer OS/architecture/checksum behavior already has executable seams, while lazy tool loading is established in the current Claude runtime contract through `ToolSearch`; the repaired design binds those mechanisms rather than inventing an environment flag. The supplied sanitized research records the Cowork tool names/read shape, but does not prove availability here. Therefore `TestLiveClaudeCoworkSurveyProbe` is an explicit first implementation gate and the source of the replay fixture; replay cannot substitute for that live result.

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

- DONE: Replace the PATH-based durable-install identity probe with explicit path semantics.
  The design invokes and verifies only `"$HOME/.local/bin/spacedock"`; absent, working, and broken exact-path states cannot be shadowed by PATH.
- DONE: Add two executable test mutants for later-PATH binaries.
  AC-2 and the installer test plan require absent-exact-path/valid-PATH consent and broken-exact-path/valid-PATH stop cases with invocation and network ledgers.
- DONE: Correct AC-1's heading to distinguish missing definitions from resolved inventory errors.
  The heading now says missing definitions retain local behavior while resolved inventory errors stop, matching the body and verifier.
- DONE: Replace the fixed temporary installer filename with a unique `mktemp` path and require honest cleanup behavior.
  The acquisition contract uses a unique path, traps cleanup, proves normal deletion, and reports without overclaiming when deletion is blocked.

### Summary

This surgical repair makes the Cowork binary identity the exact durable path rather than PATH resolution and adds executable shadow-PATH mutants for both absent and broken states. Installer acquisition is uniquely named and cleanup is verified or honestly reported, while the previously resolved capability, report, and privacy contracts remain intact.
