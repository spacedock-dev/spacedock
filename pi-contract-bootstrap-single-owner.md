---
title: "Single-owner FO contract bootstrap on Pi: extension-owned injection, version self-check, resolvable skill trigger"
status: ideation
source: "Root-cause follow-up to the 2026-09-10 stale-skill incident in the email-triage-282 FO session: the FO fell back to a stale ~/git/spacedock checkout (release/0.12.1, v0.12.1-3-g8396a6de) after BOTH contract pointers failed to resolve from a workflow cwd ($spacedock:first-officer is not pi-expandable; skills/first-officer/SKILL.md is relative and ENOENT outside the package root). Findings in /tmp/spacedock-fo-stale-skill-findings.md, validated 2026-09-11 with corrections: absence of the bootstrap message in a session log is NOT evidence (context-hook injection is request-time-only, never persisted), and no dev override was in play on that launch. Captain consolidated the derived fixes into one entity."
sprint:
id: s98gb2f779fbz41gn54ja9c3
gates:
    version: 1
    records:
        - id: gate:s98gb2f779fbz41gn54ja9c3:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:s98gb2f779fbz41gn54ja9c3-backlog-1
              briefing:
                id: briefing:s98gb2f779fbz41gn54ja9c3:backlog:attempt-1:revision-1
                digest: sha256:441c6a5c5a88c3812bfe6ec6cb9500c47ce27ec3d3662a00e608762d41ce53e6
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:s98gb2f779fbz41gn54ja9c3:backlog:1
                briefing: briefing:s98gb2f779fbz41gn54ja9c3:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-11T03:45:54.378084Z"
                decision: approve
                reason: 'Captain approved in chat 2026-09-11: seed outcome clear, scope bounded, proof named'
              application:
                target-stage: ideation
                state: consumed
started: 2026-09-11T03:46:14Z
---
Pi sessions receive the FO contract through two channels with unclear ownership, and neither is reliable today. The frontdoor launch prompt (`Use $spacedock:first-officer for this whole Pi session.`, `internal/cli/pi.go:20`) is inert syntax pi cannot expand (`agent-session.js:_expandSkillCommand` only expands `/skill:`), and the extension's session-start contract bootstrap (`FO_BOOTSTRAP_TEXT`, `.pi/extensions/spacedock.ts`) names the skill with the same unexpandable reference plus a relative path that is ENOENT from any workflow cwd. The failure is not delivery — it is resolution: even with both injections landing, an FO outside the package root ends up hunting the filesystem, and a stale visible checkout wins. Nothing in either contract verifies the loaded skill's version against the binary.

## Ideation spike: riskiest unverified mechanisms, exercised first (2026-09-11)

Ran from `/tmp/spikedock-spike` (a non-package-root cwd), with `PI_SUBAGENT_CHILD` unset so the launch is a genuine fresh frontdoor session, using the exact frontdoor argv shape (`spacedock pi "<task>" -- -p --no-session` → inner argv `pi -p --no-session "Use $spacedock:first-officer … <task>"`):

1. **Extension context-hook injection under the frontdoor argv shape: WORKS.** The session model quoted the verbatim `[SPACEDOCK-FO-BOOTSTRAP-v1]` message from its context. Injection is request-time-only (consistent with the incident findings): the session log carries no bootstrap message, and the marker is absent from the system prompt.
2. **Plain `pi` baseline: INJECTION LEAKS TODAY.** The same probe without the frontdoor (`pi -p --no-session` from `/tmp/spikedock-spike`) also received the bootstrap — the extension injects into every fresh non-subagent session, not just frontdoor launches. This is the baseline AC-1 measures against.
3. **`/skill:first-officer` resolution from a non-package-root cwd: RESOLVES, BUT TO THE WRONG COPY.** `pi -p --no-tools "/skill:first-officer …"` expanded the skill (session log holds the `<skill name="first-officer" location="…">` block), but the location was `/Users/clkao/git/spacedock-research/spacedock-v1/skills/first-officer/SKILL.md` — the dev-link checkout, NOT the installed package (`~/.pi/agent/git/github.com/spacedock-dev/spacedock`). Cause: this machine's `~/.pi/agent/settings.json` registers spacedock TWICE (`../../git/spacedock-research/spacedock-v1` and `git:github.com/spacedock-dev/spacedock`); pi's skill scan is first-match by skill name, so the dev-link shadows the git package. The duplicate-registration hazard (change 4) is live on this machine, and the version self-check (change 3) is load-bearing, not belt-and-suspenders.
4. **The model cannot self-invoke `/skill:`; its resolvable path is `<available_skills>`.** `_expandSkillCommand` (pi `dist/core/agent-session.js`) expands only user-input text that starts with `/skill:` (prompt/steer/followUp paths). The FO model loads the contract by reading the SKILL.md at the location listed in the system prompt's `<available_skills>` section — which a deterministic probe extension (hooking `before_agent_start`, dumping `event.systemPrompt`) confirmed IS present, listing `first-officer` with an absolute location. **Method finding: model self-report is unusable for this surface** — asked twice to quote `<available_skills>`, probe models answered ABSENT while the deterministic dump showed the section present. All AC proofs below use probe-extension dumps or session-log artifacts, never model self-report.

Design consequences recorded from the spike: (a) the bootstrap wording must be actionable by the model — name the `<available_skills>` location to read; `/skill:first-officer` remains the human-invocable form; (b) duplicate registration must be detected and flagged (doctor + launch warning), because "resolves" is not "resolves the installed package's skill"; (c) proofs must be deterministic artifacts (probe dumps, session logs), not model answers.

## Scope (one entity, four coupled changes)

1. **Extension-owned bootstrap.** The frontdoor launch prompt drops the contract sentence entirely (task passing and the `containsResume` suppression semantics stay: a fresh launch still appends the operator task, a resume still appends no launch prompt). The extension's bootstrap becomes the single owner. Injection is gated on a frontdoor launch marker: `runPi` sets `PI_SPACEDOCK_LAUNCH=1` in the child env (via `launchEnv`, forwarded through the safehouse wrap with a pi-specific `--env-pass` addition), and the extension injects the FO bootstrap only when the marker is present (the existing `PI_SUBAGENT_CHILD` subagent exemption stays). Wording switches to the resolvable trigger per spike finding 4: "read the `first-officer` skill's SKILL.md at the location listed for it in your available skills" — no `$pkg:skill` syntax, no relative path; `/skill:first-officer` remains the human-invocable form in interactive input. The boot-record injection at compaction carries the same marker gate (a resumed frontdoor session still carries the env, so post-compaction boot reads keep working; a plain session that compacts is not an FO session and gets no boot record).
2. **Launch-time health check.** `runPi` already refuses when the runtime is not ready; the ready gate is extended so the failure it refuses on is the one that matters: it additionally requires (a) the installed package's `.pi/extensions/spacedock.ts` to exist (today only the pi-subagents extension is Stat-ed) and (b) `firstOfficerDiscoverable` in the package status (today only ensign is checked). Duplicate registration (>1 spacedock entry in settings.json) is detected in `piSpacedockPackageStatus`: launch prints a loud warning naming both package roots and which one wins (first match — non-fatal, since resolution stays deterministic); `doctor` flags it as a first-class report line with the same remedy (`pi remove` the stale entry).
3. **FO contract version self-check.** The FO shared core's binary gate (step 1 of `skills/first-officer/references/first-officer-shared-core.md`) additionally asserts the loaded contract's declared level against the binary: the shared core declares its contract level in the gate step text (`contract level 3`), the gate parses the `contract N` token from the `--version` output it already requires, and on mismatch (different N, or token absent) aborts with a loud diagnostic naming both the skill's declared level and the binary's token, remedy `spacedock doctor` / update the plugin. Coordination constraint: the binary's `contract 3` token is pinned by `internal/cli/version_session_test.go` (frozen token, retirement condition documented there); the skill-declared level starts at 3 and moves only in lockstep with a binary token change — the frozen-token pin is the coordination point, recorded here.
4. **Docs + vocabulary (folded, no doc-only entity).** `docs/site/contributing/adding-a-runtime.md` documents: the single-owner bootstrap design; the compaction contract (boot-record re-read per #738, never contract re-injection); that context-hook injections are request-time-only and never appear in session logs (so absence in a log is not evidence — the corrected incident finding); and the duplicate-registration hazard (two registered packages both shipping `first-officer` — skill scan/expansion is first-match). The extension header comment stops saying "commissions the parent session" — the system vocabulary reserves commissioning for workflows; this installs the FO contract. The stale "Skill install and load paths" section (it still describes `--skill` flags for the spacedock skills, retired when resources_discover landed) is corrected in the same diff.

## Proposed approach — how the four changes resolve into one design

One ownership transfer, one gate, two safety nets. Ownership: the frontdoor stops pretending to deliver the contract (its sentence is unexpandable syntax) and keeps only what argv can do — pass the task, suppress on resume. The extension, which already owns request-time injection and can act at compaction boundaries, becomes the single contract delivery path, gated to fire only for frontdoor launches. Gate: `PI_SPACEDOCK_LAUNCH=1` is set by the frontdoor on every launch (wrap and non-wrap) and read by the extension; it is deliberately NOT `SPACEDOCK_BIN`, which `fo-install.md` explicitly tells users to set session-scoped — a plain session carrying that env must still get no bootstrap. Safety nets for the residual stale-contract paths the gate cannot close: the launch health check catches an unresolvable package/extension before a contract-less session starts, and the contract-level self-check catches a stale skill that DID load, regardless of how. Alternatives considered and rejected: keeping the launch prompt and fixing its syntax (pi expansion is user-input-only — the launch prompt is not user input, so no argv-embeddable syntax can expand; spike finding 4); gating on `SPACEDOCK_BIN` (false-positive per above); a separate doc-only entity for change 4 (captain folded it here; the doc surface is one file).

## Acceptance criteria (entity-level, each with a test plan)

- **AC-1 (value, measured against an independent baseline).** A plain `pi` session (no `spacedock pi`, no `PI_SPACEDOCK_LAUNCH`) started in an unrelated cwd receives no FO bootstrap: neither the marker nor any first-officer contract instruction appears in its context. Baseline measured in the spike: today the identical session receives the injection (spike result 2). Proof: live behavior probe — `pi -p` with a probe extension scanning the injected context (deterministic artifact), plus the model-quote probe as the end-to-end confirmation; owner: the spike harness promoted to a scratch probe under `tmp/` (a one-off exercise, not a standing check), plus a deterministic bun unit test on the extension module. Falsifying edit: remove the `PI_SPACEDOCK_LAUNCH` gate in `spacedock.ts` → the probe flips to marker-present and the unit test fails. Cost: low (unit), one-off live validation at implementation.
- **AC-2 (value).** A frontdoor-launched `spacedock pi` session from an arbitrary non-package-root cwd can load the FO contract: the session's system prompt lists `first-officer` under `<available_skills>` with an absolute location inside a registered spacedock package root, and a `/skill:first-officer` user input expands the skill body. Baseline: the 2026-09-10 incident session, where neither pointer resolved from a workflow cwd. Proof: probe-extension system-prompt dump + session-log grep for the expanded `<skill name="first-officer"` block (both deterministic); owner: the same scratch probe harness. Falsifying edit: break the extension's `resources_discover` skillPaths (or unregister the package) → the listing and the expansion disappear. Cost: low; one-off live validation.
- **AC-3 (mechanism).** No `spacedock pi` fresh launch appends a `$spacedock:first-officer` sentence to the inner argv; a fenced task still lands in the launch prompt; resume launches (`--resume`, `--resume=<id>`, `-r`, `--continue`, `-c`) append no launch prompt. Proof: `internal/cli/pi_frontdoor_test.go` argv assertions (existing owner of the frontdoor argv shape). Falsifying edit: restore the `piBootstrapPrompt` append → the argv test fails. Cost: low, deterministic.
- **AC-4 (mechanism).** Extension injection is gated: the FO bootstrap injects only when `PI_SPACEDOCK_LAUNCH=1` is set and the session is not a pi-subagents child; the compaction boot-record injection carries the same gate; the bootstrap wording names the `<available_skills>` location and contains no `$spacedock:` syntax and no relative SKILL.md path. Proof: bun unit test importing `.pi/extensions/spacedock.ts` with a fake `pi` object (fire `session_start`, invoke the `context` handler with and without the env; assert injected-message presence/absence and wording). Falsifying edit: drop the env check or re-add `$spacedock:first-officer` to the text → the unit test fails. Cost: medium (new small bun test file; bun is available as the TS runner). Deterministic.
- **AC-5 (mechanism).** Launch refuses (existing not-ready path, doctor report attached) when the installed package's extension file is missing or `first-officer` is not discoverable, and warns loudly (non-fatal, naming both roots and the first-match winner) when spacedock is registered twice. Proof: Go unit tests with fake `piRuntimeOps` (Stat misses; double-entry settings.json fixture). Falsifying edit: revert the ready-gate extension/first-officer additions or the duplicate warning → the tests fail. Cost: low, deterministic. Owner: `internal/cli/pi_frontdoor_test.go` / `pi_launch_test.go`.
- **AC-6 (value, version-skew abort).** A session whose loaded FO contract declares a contract level that mismatches the binary's `contract N` token aborts at the binary gate with a diagnostic naming both versions. Baseline: a stale-checkout skill proceeding silently (the incident's residual path). Proof: one-off live falsification — run a session against a checkout whose shared-core declared level differs from the binary token (or with the declaration stripped) and record the abort diagnostic; plus the declaration literal pinned so any future bump without a lockstep binary change trips the gate. Falsifying edit: change the declared level in the shared core → the live run aborts. Owner: FO shared-core prose gate, proven by an exercised session per the proof policy (prose rules need a real check, not a presence-grep). Cost: medium (one-off live run against a deliberately stale fixture checkout).
- **AC-7 (mechanism).** `spacedock doctor --host pi` reports the first-officer skill check line and flags duplicate registration with a remedy naming `pi remove`. Proof: Go unit test on `printPiDoctorReport` with a double-entry `piPackageStatus` fixture. Falsifying edit: remove the warning/report lines → the test fails. Cost: low, deterministic.
- **AC-8 (mechanism, doc).** The reference doc carries the four documented facts with the concrete wording in this body's doc diffs, and the stale "Skill install and load paths" text is corrected. Proof: the doc diff in this body is the reviewable artifact; implementation applies it (modulo merge drift); checked at the ideation gate and re-checked at implementation review. Counts only paired with AC-1/AC-2/AC-6, which measure the value the docs describe. Falsifying edit: removing any of the four facts → gate review fails.

## Expected surface

Net LOC change: **+230 across 6 files** (insertions +281 / deletions −52), tolerance ±60 net lines / ±2 files:

| File | +ins | −del | Content |
|---|---|---|---|
| `internal/cli/pi.go` (893) | +55 | −12 | drop `piBootstrapPrompt` + prompt-append reshape; `PI_SPACEDOCK_LAUNCH=1` env + wrap `--env-pass`; ready-gate adds installed-ext Stat + `firstOfficerDiscoverable`; duplicate-registration detection in `piSpacedockPackageStatus` |
| `.pi/extensions/spacedock.ts` (~130) | +10 | −8 | marker gate on both injections; bootstrap wording; header-comment vocabulary |
| `skills/first-officer/references/first-officer-shared-core.md` | +10 | −2 | contract-level declaration + gate step |
| `docs/site/contributing/adding-a-runtime.md` | +55 | −20 | corrected load-paths section + new bootstrap subsection |
| `internal/cli/pi_frontdoor_test.go` + `pi_launch_test.go` | +130 | −8 | argv, gating env, duplicate-warning, doctor-report tests |
| new bun test for the extension module | +21 | −2 | extension injection gate unit test |

Observable semantics this task may change (the boundary the implementation must not cross silently): the frontdoor inner argv shape (a fresh launch no longer ends with the contract sentence; it ends with the operator task or nothing); the launched child env (new `PI_SPACEDOCK_LAUNCH=1` on every `spacedock pi` launch, wrap and non-wrap); the extension injection surface (bootstrap and boot-record injections now fire only for frontdoor-launched sessions); the FO startup step (the binary gate gains a contract-level assertion that can abort); doctor/launch output (new first-officer + installed-extension checks that refuse, and a duplicate-registration warning that does not). No stored formats, command grammar, or authority semantics change.

## Doc diffs (concrete, for the ideation gate)

**Diff 1 — `docs/site/contributing/adding-a-runtime.md`, "Skill install and load paths" intro (stale text corrected).**

Before:
> For Pi, `spacedock pi` launches the proven front door by loading local resources explicitly:

Before (path block):
> ```text
> <spacedock checkout>/skills/first-officer
> <spacedock checkout>/skills/ensign
> ~/.pi/agent/npm/node_modules/pi-subagents/skills/pi-subagents
> ~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
> ```

After:
> For Pi, `spacedock pi` launches the front door through the registered Spacedock package: the package's `.pi/extensions/spacedock.ts` discovers the Spacedock skills (`first-officer`, `ensign`, …) from the package's own `skills/` directory via pi's `resources_discover`, and installs the FO contract through Pi's context hook, gated on the `PI_SPACEDOCK_LAUNCH=1` marker the frontdoor sets on every launch:

After (path block):
> ```text
> <spacedock package root>/.pi/extensions/spacedock.ts   (registered package or --plugin-dir dev override)
> <spacedock package root>/skills/first-officer
> <spacedock package root>/skills/ensign
> ~/.pi/agent/npm/node_modules/pi-subagents/skills/pi-subagents
> ~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
> ```

**Diff 2 — new subsection `### First-officer contract bootstrap (Pi)` inserted after "Exact Pi parent prompt".**

> The FO contract has one owner on Pi: the Spacedock extension. The frontdoor launch prompt carries no contract sentence — it passes the operator task and suppresses the launch prompt on resume; the extension injects the contract through Pi's context hook, gated on the `PI_SPACEDOCK_LAUNCH=1` marker `spacedock pi` sets on every launch, so a plain `pi` session for unrelated work receives no bootstrap. The bootstrap names the resolvable trigger: the `first-officer` entry of the session's `<available_skills>` listing, whose location pi resolves from the registered package root regardless of cwd (`/skill:first-officer` remains the human-invocable form in interactive input; the model loads the skill by reading the listed location).
>
> Three facts about this surface that session evidence cannot show directly:
>
> - **Compaction contract.** At a compaction boundary the FO re-reads durable state via the boot-record injection (PR #738); the contract is never re-injected.
> - **Request-time only.** Context-hook injections are never persisted to session logs. Absence of the bootstrap message in a session log is NOT evidence that it was absent at request time.
> - **Duplicate registration.** When two registered packages both provide `first-officer` (e.g. a dev-link checkout and the installed git package), pi's skill scan is first-match: the first `settings.json` entry wins. `spacedock doctor --host pi` and launch output flag the condition; the remedy is `pi remove` of the stale entry.

**Diff 3 — `skills/first-officer/references/first-officer-shared-core.md`, binary gate step.**

Before:
> 1. **Binary gate.** … `${SPACEDOCK_BIN:-spacedock} --version` line 1 must be `spacedock <version>`. These skills require binary minor 0.28.

After (appended to the same numbered step):
> … The same output carries the contract token: this contract is **contract level 3**, and the `--version` output must contain `contract 3`. A different level, or a missing token → ABORT with a diagnostic naming the skill's declared level and the binary's token; remedy: `spacedock doctor` and update the plugin.

**Diff 4 — `.pi/extensions/spacedock.ts` header comment.**

Before:
> // It also commissions the parent session as Spacedock first officer through

After:
> // It also installs the FO contract in the parent session through

## Out of scope

- Dev-override (`--plugin-dir` / `SPACEDOCK_REPO_ROOT`) version guard: latent hazard, not the mechanism in play in the 2026-09-10 incident; candidate for its own entity if the captain wants it.
- Safehouse profile `add-dirs` grant hygiene (stale grants pointing at old checkouts): doctor-note material at most; visibility is not resolution.
- The `gate record` exit-0-with-frozen-error quirk observed once mid-incident: re-verify against the current binary before filing anything.

## Stage Report: ideation

- DONE: A fleshed-out task body whose proposed approach resolves the four coupled changes into one coherent design (extension-owned gated bootstrap, frontdoor health check, FO version self-check, folded docs)
  Body now carries "Proposed approach — how the four changes resolve into one design": one ownership transfer (frontdoor keeps argv-only duties), one gate (`PI_SPACEDOCK_LAUNCH=1`, deliberately not `SPACEDOCK_BIN`), two safety nets (ready-gate refusal, contract-level abort), with rejected alternatives named.
- DONE: The riskiest unverified mechanism exercised first and its result recorded in the body
  Spike section records 4 results: context hook injects under the frontdoor's exact argv shape from /tmp/spikedock-spike (model quoted the verbatim bootstrap); plain `pi` sessions receive the injection today (baseline); `/skill:first-officer` expands from a non-package-root cwd but resolves the dev-link checkout over the installed git package (duplicate registration, first-match wins); `/skill:` expansion is user-input-only so the model's resolvable path is the `<available_skills>` listing (deterministic probe dump).
- DONE: Entity-level acceptance criteria, each with a test plan naming its proof owner and distinct falsifying check
  AC-1..AC-8 each name proof owner, falsifying edit, cost, and deterministic/live classification; AC-1 measures against the measured plain-pi baseline and AC-6 against the stale-checkout abort baseline.
- DONE: At least one AC measuring the end value against an independent baseline
  AC-1 (plain `pi` session receives no bootstrap; baseline measured today: it does) and AC-2 (contract loadable from arbitrary cwd; baseline: the 2026-09-10 incident session could not) and AC-6 (version-mismatched session aborts; baseline: stale skill proceeds silently).
- DONE: Expected surface declared as net LOC change with insertions/deletions separate across N files with tolerance, observable semantics, and concrete doc diffs
  +230 net (+281/−52) across 6 files, tolerance ±60 net lines / ±2 files; five observable-semantics changes declared; four concrete before/after doc diffs recorded in the body.

### Summary

Fleshed out the entity body for the ideation gate and ran the required spike first. The spike materially changed the design twice: (1) `/skill:first-officer` cannot be the model's trigger (pi expands it only on user input), so the bootstrap wording points at the `<available_skills>` location instead, keeping `/skill:first-officer` as the human-invocable form; (2) on this machine spacedock is registered twice and the dev-link checkout shadows the installed git package, so the duplicate-registration flag and the contract-level self-check are load-bearing. Second method finding recorded in the body: model self-report about its own system prompt was wrong twice (claimed `<available_skills>` ABSENT when the deterministic probe dump showed it present), so every AC proof is specified as a deterministic artifact (probe-extension dump or session log), never a model answer. Spike scratch lives in /tmp/spikedock-spike (probe extensions, session logs) and is intentionally not committed.
