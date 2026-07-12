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

Add a short `0. Select the evidence adapter` section before current step 1. Detection is a positive capability match: when the current tool inventory contains both `mcp__session_info__list_sessions` and `mcp__session_info__read_transcript`, select `cowork`; otherwise preserve the current local `agentsview` path byte-for-byte. Do not infer Cowork from a missing executable, missing `~/.claude/projects`, Linux, or a sandbox error.

On the Cowork branch:

1. Invoke `spacedock --version` to test the executable, without a PATH filesystem probe.
2. If unavailable, display the exact prompt below and stop. Do not call `curl`, a release API, or any other network tool before the user affirms that network access is enabled.
3. After affirmation, run the repository's stable, checksum-verifying installer with `SPACEDOCK_INSTALL_DIR="$HOME/.local/bin"`. The installer is the single owner of `uname -s`/`uname -m` normalization (`Linux` plus `x86_64|amd64` -> `linux_amd64`; `arm64|aarch64` -> `linux_arm64`), release resolution, checksum verification, and atomic failure. Export that directory for the current session and verify `spacedock --version`. Unsupported OS/architecture or denied/blocked access leaves no executable installed and reports the exact next action; it does not retry alternate URLs or request credentials.
4. Call the Cowork session inventory, then batch bounded transcript reads for idle sessions. Use titles/state only for inventory and transcript content only in-memory to cluster workstreams, decisions, interruptions, and modes. Do not call `agentsview`, read `~/.claude/projects`, run repo scaffold probes, attribute Codex by workdir, or infer work-by-file.
5. Render the existing survey report. Mark shipped/open, git/PR, scaffold, and file-area conclusions `unverified — Cowork session evidence has no repository cross-check`. Preserve the existing end-of-report commission offer.

### Exact survey instruction change

Before (current opening route):

> Run the four steps in order: **check agentsview → scan → recognize scaffold → report and offer**.

After:

> Run the evidence-adapter check first. If both Cowork session tools (`mcp__session_info__list_sessions` and `mcp__session_info__read_transcript`) are available, run **check Spacedock permission → Cowork session scan → Cowork report and offer** and do not enter the agentsview or repo-probe steps. Otherwise run the existing **check agentsview → scan → recognize scaffold → report and offer** path unchanged. Missing files or binaries alone never identify Cowork.

Insert this exact permission prompt when the Cowork branch cannot invoke Spacedock:

> Spacedock is not available in this Cowork session. To continue, I need to download the checksum-verified Linux release from `github.com/spacedock-dev/spacedock` into `~/.local/bin`; the installer will choose arm64 or amd64 from this sandbox's architecture. Please enable network access for this Cowork session and tell me when it is enabled. I will not make any network request or download anything until you confirm.

### Privacy boundary

Session identifiers are handles for tool calls only. Never echo or persist them. Do not copy transcript text, account details, credentials, tokens, or machine-specific paths into diagnostics, fixtures, reports, task artifacts, or logs. Derived counts, generic synthetic labels, redacted capability results, and aggregate workstream summaries are permitted. Do not pass transcript text to the Spacedock installer or any network call.

### Documentation change proposed for implementation

In `docs/site/get-started/survey.md`, replace the local-only opening paragraph:

> When you use the `spacedock:survey` skill, it looks at your existing agent conversation logs on local disk (through agentsview, an open source session-history tool). It is read-only; if agentsview is missing, it asks before installing it.

with:

> When you use the `spacedock:survey` skill, it reads the session history available to the current host. Local coding-agent sessions use agentsview; Claude Cowork uses Cowork's session inventory and transcript tools. Survey asks before installing a missing helper or making the network request needed to obtain Spacedock, and it labels repository conclusions unverified when Cowork provides no repository evidence.

No install-page or commission documentation changes belong to this task.

## Out of scope

- General Claude Cowork workflow commissioning or git support.
- A new `using-spacedock-with-cowork` skill or general Cowork runtime adapter.
- Changes to connected-folder deletion semantics.
- Persisting or publishing Cowork transcripts, session identifiers, account data, filesystem paths, credentials, or other personal/private information.
- Automatic network-policy changes or downloads without the user's explicit enablement.
- GitHub issue creation.

## Acceptance criteria

**AC-1 - A Cowork-capable tool inventory selects exactly one Cowork evidence path, while the non-Cowork baseline retains existing local behavior.**
Verified by: a table-driven routing fixture where both Cowork tools yield `cowork`, either tool alone yields `local`, and the Cowork case records zero `agentsview`, `~/.claude/projects`, and repo-probe operations versus at least one local-history operation in the unchanged local control.

**AC-2 - A missing Spacedock binary produces the specified network-access prompt and zero network/download attempts before affirmative permission.**
Verified by: a fake network transport and event ledger; before affirmation its request count and install-directory file count are both 0, denial keeps both at 0, and affirmation permits exactly the installer requests.

**AC-3 - After permission, Cowork installation uses the release asset matching the observed Linux architecture, verifies its checksum, installs only to the configured writable user directory, and successfully invokes the resulting binary.**
Verified by: the real `install.sh` inspection seam for `x86_64/amd64` and `arm64/aarch64`, fixture tarballs/checksums for install success and mismatch failure, and a fake `$HOME` asserting `~/.local/bin/spacedock --version`; unsupported architecture and blocked network leave the install directory empty.

**AC-4 - Cowork survey derives its orientation from host-native session evidence and does not overstate conclusions that require repository evidence.**
Verified by: a synthetic Cowork tool harness with two idle sessions and one active session; bounded reads of the two idle sessions change the rendered workstream/decision counts from the zero-session baseline, the active session is not read, and repo-only fields carry the required unverified label.

**AC-5 - Cowork survey output and durable test artifacts contain no raw private session material.**
Verified by: synthetic canary identifiers, transcript sentences, token strings, account values, and absolute paths injected into the harness; recursive assertions over captured output and fixture-generated files find none of the canaries while aggregate counts and generic workstream labels remain.

**AC-6 - User documentation describes host-native Cowork evidence, consent-before-network behavior, and the repository-evidence limitation.**
Verified by: render the documentation site and inspect the Survey page for the proposed behavior; the integration harness remains the behavioral proof rather than a prose substring test.

## Test plan

Implementation starts with the fixture harness so routing, consent, privacy, and rendered-report assertions fail against the current skill before its instructions change.

- **Fixture/skill smoke (medium):** add a Cowork tool-inventory and invocation-ledger harness beside the current survey integration tests. Exercise the real instruction artifact through the supported skill smoke mechanism, not a grep for phrases. Cover complete/partial/absent capabilities, consent yes/no, bounded idle-session reads, repo-probe absence, report degradation, and privacy canaries.
- **Installer CLI tests (low, mostly existing seams):** reuse the real `install.sh` with stubbed `uname`, fixture assets/checksums, fake network transport, and fake `$HOME`. Assert architecture mapping, before-consent request count, verified install state, and fail-closed cases. Do not duplicate installer selection logic in survey.
- **Documentation render (low):** build the docs site and visually/readably inspect the changed Survey page.
- **Live Cowork smoke (required for the runtime claim, bounded):** in a sanitized Cowork session, confirm the two positive tool capabilities, begin with no usable Spacedock binary, observe that no download occurs before the exact prompt and approval, enable network access, install/verify the architecture-correct binary, list sessions, read one synthetic or explicitly non-sensitive idle session, and render a report. Durable evidence is limited to process exit, redacted capability booleans, selected `linux_{arm64|amd64}` token, binary version, invocation counts, and absence of privacy canaries.

## Risk spike

No new throwaway product spike is needed before implementation. The two load-bearing mechanisms already have executable proof surfaces: `install.sh` owns and tests OS/architecture selection, checksum verification, writable-directory installation, and `--version`; the supplied sanitized Cowork research artifact records a live observation of the exact session inventory/transcript tools and their idle-session read shape. The remaining uncertainty is host availability and permission UX, so it cannot be honestly closed from this non-Cowork worker; the required live Cowork smoke above is the first implementation validation gate, and fixture tests must not substitute for it.

## Stage Report: ideation

- DONE: Compare a dedicated using-spacedock-with-cowork skill with integrating Cowork support into survey; recommend one smallest coherent design using discoverability, ownership, overlap, and maintenance evidence.
  Recommended one capability-gated branch in `survey`; the comparison table records why a second user-invocable skill would overlap and drift.
- DONE: Update the task with a concrete behavior-first design, including exact skill-routing or before/after instruction text, privacy boundaries, and the user-visible network-access prompt.
  The proposed route, exact prompt, privacy contract, acceptance criteria, and Survey documentation replacement are specified above.
- DONE: Define reproducible proof for Cowork detection, no download before permission, architecture-aware binary setup, and Cowork-native session evidence; spike the riskiest unverified mechanism or record why no spike is needed.
  Fixture, installer, privacy-canary, docs-render, and mandatory live-Cowork proofs are defined; existing executable installer seams and supplied live research make a throwaway spike unnecessary.

### Summary

Ideation selects integration into the existing survey skill as the smallest coherent design and explicitly excludes the separate commission/filesystem issue. The design detects Cowork from two positive session capabilities, stops before all network activity for user consent, delegates architecture/checksum handling to the existing installer, and requires a privacy-bounded live Cowork smoke before the runtime claim can pass.
