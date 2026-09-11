---
title: "Single-owner FO contract bootstrap on Pi: extension-owned injection, version self-check, resolvable skill trigger"
status: validation
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
        - id: gate:s98gb2f779fbz41gn54ja9c3:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:s98gb2f779fbz41gn54ja9c3-ideation-1
              briefing:
                id: briefing:s98gb2f779fbz41gn54ja9c3:ideation:attempt-1:revision-1
                digest: sha256:80c82ce4ad3b5e541b8f23136ec5cf2f0dc23d38f59a6a9f37a18a03ab36607a
                room-ref: '@review/ideation/briefing-1'
              withdrawal:
                by: agent:first-officer
                at: "2026-09-11T04:11:02.44722Z"
                reason: Artifact normalized after prepare (heading renamed to the scanner-exact '## Acceptance criteria' under captain direct-edit grant); bound briefing digest would be stale at record time
            - id: gate-attempt:s98gb2f779fbz41gn54ja9c3-ideation-2
              briefing:
                id: briefing:s98gb2f779fbz41gn54ja9c3:ideation:attempt-2:revision-1
                digest: sha256:1d472e1deb2c5d90362ff1b8a188cd7df9011b49496dd5a09aa654806a56b3f6
                room-ref: '@review/ideation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:s98gb2f779fbz41gn54ja9c3:ideation:2
                briefing: briefing:s98gb2f779fbz41gn54ja9c3:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-09-11T04:25:39.248881Z"
                decision: revise
                reason: 'Captain-directed revision after gate probe, four concrete asks: (1) add AC coverage for the --plugin-dir dev-override arm — PI_SPACEDOCK_LAUNCH present in child env and the gated extension firing, live-probed once against a current checkout; (2) state the duplicate-extension load case deterministically (dev-override --extension plus installed-package extension can both register; dedupe coverage or ready-gate/doctor flag, not emergent behavior); (3) declare the pi-behavior dependency list (context-hook API session_start/session_compact/context, <available_skills> listing with absolute locations, /skill: user-input-only expansion) and pin a pi >= 0.83 version floor enforced in ready-gate/doctor reading ''pi --version'' at binary level — pi''s org moved to @earendil-works (old-org installs flagged, no package-path checks); (4) note expandPromptTemplates on pi.sendUserMessage() as the known upgrade path for programmatic skill expansion. Context: pi 0.83.0-0.85.1 changelog shows no changes to the three load-bearing behaviors; version floor protects the next drift.'
            - id: gate-attempt:s98gb2f779fbz41gn54ja9c3-ideation-3
              briefing:
                id: briefing:s98gb2f779fbz41gn54ja9c3:ideation:attempt-3:revision-1
                digest: sha256:23bed9815f6592ea8814d30782ac990e33f4a0e76856d68ad7151cf4087db3a7
                room-ref: '@review/ideation/briefing-3'
              resolution:
                type: Resolution
                id: resolution:spacedock:s98gb2f779fbz41gn54ja9c3:ideation:3
                briefing: briefing:s98gb2f779fbz41gn54ja9c3:ideation:attempt-3:revision-1
                by: person:captain
                at: "2026-09-11T04:57:47.171269Z"
                decision: approve
                reason: 'Captain approved in chat 2026-09-11: revision-2 fold is complete and verifiable — all four revise asks landed in ACs/approach/surface/doc diffs (commit 04ce892d4), stage report 3 done/0 skipped/0 failed, surface +270 net within approved tolerance, spike evidence stands'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:s98gb2f779fbz41gn54ja9c3:validation
          stage: validation
          attempts:
            - id: gate-attempt:s98gb2f779fbz41gn54ja9c3-validation-1
              briefing:
                id: briefing:s98gb2f779fbz41gn54ja9c3:validation:attempt-1:revision-1
                digest: sha256:b525f19fe5ed805b7687912881799f795d86f5b8b91d0b02948668c362a6d455
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:s98gb2f779fbz41gn54ja9c3:validation:1
                briefing: briefing:s98gb2f779fbz41gn54ja9c3:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-11T05:54:00.274763Z"
                decision: approve
                reason: 'Captain approved in chat 2026-09-11 after surface-drift and contract-comment review: PASSED verdict accepted with fresh falsifiable evidence for all eight ACs; drift adjudicated as revise-cycle test matrix + contract comments (captain-visible note)'
              application:
                target-stage: done
                state: pending
started: 2026-09-11T03:46:14Z
worktree: .worktrees/spacedock-ensign-pi-contract-bootstrap-single-owner
mod-block: merge:pr-merge
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

1. **Extension-owned bootstrap.** The frontdoor launch prompt drops the contract sentence entirely (task passing and the `containsResume` suppression semantics stay: a fresh launch still appends the operator task, a resume still appends no launch prompt). The extension's bootstrap becomes the single owner. Injection is gated on a frontdoor launch marker: `runPi` sets `PI_SPACEDOCK_LAUNCH=1` in the child env (via `launchEnv`, forwarded through the safehouse wrap with a pi-specific `--env-pass` addition), and the extension injects the FO bootstrap only when the marker is present (the existing `PI_SUBAGENT_CHILD` subagent exemption stays). Wording switches to the resolvable trigger per spike finding 4: "read the `first-officer` skill's SKILL.md at the location listed for it in your available skills" — no `$pkg:skill` syntax, no relative path; `/skill:first-officer` remains the human-invocable form in interactive input. The boot-record injection at compaction carries the same marker gate (a resumed frontdoor session still carries the env, so post-compaction boot reads keep working; a plain session that compacts is not an FO session and gets no boot record). The gate covers the dev-override arm identically: `PI_SPACEDOCK_LAUNCH=1` is set on every launch shape including `repoRoot` set (`SPACEDOCK_REPO_ROOT` / `--plugin-dir` install source), with one live probe of the gated extension under that shape against a CURRENT checkout; a stale checkout under the dev override is intentionally not made to work — a stale contract that still loads is AC-6's abort territory, not a supported path.
2. **Launch-time health check.** `runPi` already refuses when the runtime is not ready; the ready gate is extended so the failure it refuses on is the one that matters: it additionally requires (a) the installed package's `.pi/extensions/spacedock.ts` to exist (today only the pi-subagents extension is Stat-ed), (b) `firstOfficerDiscoverable` in the package status (today only ensign is checked), and (c) the pi binary's `--version` output to parse at ≥ 0.83.0 — the declared floor for the load-bearing pi behaviors (see the dependency declaration under Proposed approach). The floor is read at the BINARY level only (`pi --version` via a new `piRuntimeOps` version hook so tests can fake it); no package-path checks. A sub-floor install fails the same not-ready path; in `doctor` it is a first-class line whose diagnostic names the org move (pi now publishes as `@earendil-works`, 0.83.0–0.85.1; the old `@mariozechner` name is stale at 0.73.1 — a sub-floor binary IS the old-org install signal, no package paths consulted). Duplicate registration (>1 spacedock entry in settings.json) is detected in `piSpacedockPackageStatus` — route-aware: the count is packages × routes, where each package contributes TWO skill-registration routes (the manifest's `pi.skills` scan and the extension's `resources_discover`, which re-registers the same `skills/` directory its own manifest already declares — live-observed as doubled loser entries in pi's skill-conflict report), so up to 4 registrations per skill name on this machine's double-registered setup. Launch prints a loud warning naming all package roots and which one wins (first match — non-fatal, since resolution stays deterministic), and additionally warns on the dev-override double-extension load (`repoRoot` set AND a spacedock package registered: the checkout's spacedock.ts loads via `--extension` while pi's package discovery loads the installed copy — a static, code-visible consequence of the dev-override argv construction in `internal/cli/pi.go`); `doctor` flags both as first-class report lines with the same remedy (`pi remove` the stale entry). The double-injection question (two context hooks on one session event) is NOT left to emergent behavior: the extension's `hasStructuralBootstrap` dedupe provably covers only re-fires within ONE hook instance (the same `context` handler seeing an event whose messages already carry the marker); whether pi chains one hook's modified messages into the next hook's input is unverified, so the ready-gate/doctor double-extension flag above is the deterministic answer — and the cleanest remedy removes the redundancy outright: `resources_discover` returns no `skillPaths` when the package manifest's `pi.skills` already declares the same skills directory (unit-tested), eliminating the intra-package double registration entirely.
3. **FO contract version self-check.** The FO shared core's binary gate (step 1 of `skills/first-officer/references/first-officer-shared-core.md`) additionally asserts the loaded contract's declared level against the binary: the shared core declares its contract level in the gate step text (`contract level 3`), the gate parses the `contract N` token from the `--version` output it already requires, and on mismatch (different N, or token absent) aborts with a loud diagnostic naming both the skill's declared level and the binary's token, remedy `spacedock doctor` / update the plugin. Coordination constraint: the binary's `contract 3` token is pinned by `internal/cli/version_session_test.go` (frozen token, retirement condition documented there); the skill-declared level starts at 3 and moves only in lockstep with a binary token change — the frozen-token pin is the coordination point, recorded here.
4. **Docs + vocabulary (folded, no doc-only entity).** `docs/site/contributing/adding-a-runtime.md` documents: the single-owner bootstrap design; the compaction contract (boot-record re-read per #738, never contract re-injection); that context-hook injections are request-time-only and never appear in session logs (so absence in a log is not evidence — the corrected incident finding); and the duplicate-registration hazard (two registered packages both shipping `first-officer` — skill scan/expansion is first-match). The extension header comment stops saying "commissions the parent session" — the system vocabulary reserves commissioning for workflows; this installs the FO contract. The stale "Skill install and load paths" section (it still describes `--skill` flags for the spacedock skills, retired when resources_discover landed) is corrected in the same diff.

## Proposed approach — how the four changes resolve into one design

One ownership transfer, one gate, two safety nets. Ownership: the frontdoor stops pretending to deliver the contract (its sentence is unexpandable syntax) and keeps only what argv can do — pass the task, suppress on resume. The extension, which already owns request-time injection and can act at compaction boundaries, becomes the single contract delivery path, gated to fire only for frontdoor launches. Gate: `PI_SPACEDOCK_LAUNCH=1` is set by the frontdoor on every launch (wrap and non-wrap) and read by the extension; it is deliberately NOT `SPACEDOCK_BIN`, which `fo-install.md` explicitly tells users to set session-scoped — a plain session carrying that env must still get no bootstrap. Safety nets for the residual stale-contract paths the gate cannot close: the launch health check catches an unresolvable package/extension before a contract-less session starts, and the contract-level self-check catches a stale skill that DID load, regardless of how. Alternatives considered and rejected: keeping the launch prompt and fixing its syntax (pi expansion is user-input-only — the launch prompt is not user input, so no argv-embeddable syntax can expand; spike finding 4); gating on `SPACEDOCK_BIN` (false-positive per above); a separate doc-only entity for change 4 (captain folded it here; the doc surface is one file).

**Load-bearing pi behaviors (declared dependency, enforced as a floor).** The mechanism rests on three pi behaviors: (1) the context-hook API — `session_start` / `session_compact` events and a `context` hook returning modified messages (the request-time delivery path for the bootstrap and boot record); (2) the `<available_skills>` system-prompt listing carrying an absolute per-skill location (the model's resolvable trigger, spike finding 4); (3) `/skill:` expansion being user-input-only (`agent-session.js:_expandSkillCommand` — why the launch prompt cannot carry the contract, and why `/skill:first-officer` remains the human-invocable form only). Survey basis: the 0.83.0–0.85.1 changelog touches none of the three behaviors; injection-ordering and package-glob mechanics are adjacent churn the floor protects against. These are pinned as a `pi >= 0.83` floor enforced in the ready-gate and doctor from the binary's `pi --version` (verified 2026-09-11: the binary prints a bare semver — `0.85.1` on this machine — so the parse is trivial and path-free).

**Duplicate-extension load under dev override, stated deterministically.** Under the dev override the checkout's spacedock.ts (`--extension`, added by the argv construction whenever `repoRoot` is set) AND the installed package's spacedock.ts (settings.json package discovery) can both register; each package also registers its skills through two routes (manifest `pi.skills` + the extension's `resources_discover` on the same directory). The design does NOT rely on dedupe covering the two-hook case (see change 2 for the mechanism boundary): the launch warning and doctor flag fire on the static, code-visible condition (`repoRoot` set + spacedock package registered), and the `resources_discover` manifest-skip removes the intra-package doubling outright.

## Acceptance criteria

- **AC-1 (value, measured against an independent baseline).** A plain `pi` session (no `spacedock pi`, no `PI_SPACEDOCK_LAUNCH`) started in an unrelated cwd receives no FO bootstrap: neither the marker nor any first-officer contract instruction appears in its context. Baseline measured in the spike: today the identical session receives the injection (spike result 2). Proof: live behavior probe — `pi -p` with a probe extension scanning the injected context (deterministic artifact), plus the model-quote probe as the end-to-end confirmation; owner: the spike harness promoted to a scratch probe under `tmp/` (a one-off exercise, not a standing check), plus a deterministic bun unit test on the extension module. Falsifying edit: remove the `PI_SPACEDOCK_LAUNCH` gate in `spacedock.ts` → the probe flips to marker-present and the unit test fails. Cost: low (unit), one-off live validation at implementation.
- **AC-2 (value).** A frontdoor-launched `spacedock pi` session from an arbitrary non-package-root cwd can load the FO contract: the session's system prompt lists `first-officer` under `<available_skills>` with an absolute location inside a registered spacedock package root, and a `/skill:first-officer` user input expands the skill body. Baseline: the 2026-09-10 incident session, where neither pointer resolved from a workflow cwd. Proof: probe-extension system-prompt dump + session-log grep for the expanded `<skill name="first-officer"` block (both deterministic); owner: the same scratch probe harness. Falsifying edit: break the extension's `resources_discover` skillPaths (or unregister the package) → the listing and the expansion disappear. Cost: low; one-off live validation.
- **AC-3 (mechanism).** No `spacedock pi` fresh launch appends a `$spacedock:first-officer` sentence to the inner argv; a fenced task still lands in the launch prompt; resume launches (`--resume`, `--resume=<id>`, `-r`, `--continue`, `-c`) append no launch prompt. Proof: `internal/cli/pi_frontdoor_test.go` argv assertions (existing owner of the frontdoor argv shape). Falsifying edit: restore the `piBootstrapPrompt` append → the argv test fails. **Dev-override arm:** one Go unit test asserts `PI_SPACEDOCK_LAUNCH=1` is present in the launched child env when `repoRoot` is set (the `--plugin-dir` / `SPACEDOCK_REPO_ROOT` path) — same gate, same env, wrap and non-wrap. Cost: low, deterministic.
- **AC-4 (mechanism).** Extension injection is gated: the FO bootstrap injects only when `PI_SPACEDOCK_LAUNCH=1` is set and the session is not a pi-subagents child; the compaction boot-record injection carries the same gate; the bootstrap wording names the `<available_skills>` location and contains no `$spacedock:` syntax and no relative SKILL.md path. Proof: bun unit test importing `.pi/extensions/spacedock.ts` with a fake `pi` object (fire `session_start`, invoke the `context` handler with and without the env; assert injected-message presence/absence and wording). Falsifying edit: drop the env check or re-add `$spacedock:first-officer` to the text → the unit test fails. **Dev-override arm:** one live probe of the gated extension loaded via `--extension <checkout>/.pi/extensions/spacedock.ts` (the dev-override argv shape, `repoRoot` set) pointing at a CURRENT checkout confirms the marker gate holds outside package registration; a stale checkout under the dev override is intentionally not made to work (AC-6 aborts it). The same bun test also asserts the `resources_discover` manifest-skip: when the package manifest's `pi.skills` already declares the same skills directory, `resources_discover` returns no `skillPaths` (the intra-package double-registration remedy). Cost: medium (new small bun test file; bun is available as the TS runner) plus one one-off live probe. Deterministic except the one-off probe.
- **AC-5 (mechanism).** Launch refuses (existing not-ready path, doctor report attached) when the installed package's extension file is missing or `first-officer` is not discoverable, and warns loudly (non-fatal, naming all roots and the first-match winner, route-aware across packages × routes) when spacedock is registered twice, and on the dev-override double-extension load (`repoRoot` set AND a spacedock package registered — both extension copies load). It also refuses when the pi binary's `--version` parses below the 0.83.0 floor (binary-level read; no package paths). Proof: Go unit tests with fake `piRuntimeOps` (Stat misses; double-entry settings.json fixture; version below/at/above floor). Falsifying edit: revert the ready-gate extension/first-officer additions, the duplicate or double-extension warning, or the floor check → the tests fail. Cost: low, deterministic. Owner: `internal/cli/pi_frontdoor_test.go` / `pi_launch_test.go`.
- **AC-6 (value, version-skew abort).** A session whose loaded FO contract declares a contract level that mismatches the binary's `contract N` token aborts at the binary gate with a diagnostic naming both versions. Baseline: a stale-checkout skill proceeding silently (the incident's residual path). Proof: one-off live falsification — run a session against a checkout whose shared-core declared level differs from the binary token (or with the declaration stripped) and record the abort diagnostic; plus the declaration literal pinned so any future bump without a lockstep binary change trips the gate. Falsifying edit: change the declared level in the shared core → the live run aborts. Owner: FO shared-core prose gate, proven by an exercised session per the proof policy (prose rules need a real check, not a presence-grep). Cost: medium (one-off live run against a deliberately stale fixture checkout).
- **AC-7 (mechanism).** `spacedock doctor --host pi` reports the first-officer skill check line, reports the pi version against the 0.83.0 floor (a sub-floor install is flagged as the stale `@mariozechner`-org signal with an upgrade diagnostic naming `@earendil-works`), and flags duplicate registration with a route-aware count (packages × routes) and a remedy naming `pi remove`. Proof: Go unit test on `printPiDoctorReport` with double-entry `piPackageStatus` and sub-floor-version fixtures. Falsifying edit: remove the warning/report lines → the test fails. Cost: low, deterministic.
- **AC-8 (mechanism, doc).** The reference doc carries the five documented facts with the concrete wording in this body's doc diffs (compaction contract, request-time-only persistence, duplicate registration, dev-override double-extension load, `pi >= 0.83` floor), and the stale "Skill install and load paths" text is corrected. Proof: the doc diff in this body is the reviewable artifact; implementation applies it (modulo merge drift); checked at the ideation gate and re-checked at implementation review. Counts only paired with AC-1/AC-2/AC-6, which measure the value the docs describe. Falsifying edit: removing any of the five facts → gate review fails.

### Feedback Cycles

- Cycle 1 (2026-09-11, captain-directed revise at ideation gate attempt-2, no reviewer run): fold four additions into approach/ACs/surface/doc diffs — (1) `--plugin-dir` dev-override arm coverage (env + gated-extension assertions plus one live probe against a current checkout; stale checkout stays AC-6 territory); (2) duplicate-extension load under dev override stated deterministically (dedupe mechanism + test, or ready-gate/doctor flag); (3) declared pi-behavior dependency list (context-hook API, `<available_skills>` listing with absolute locations, `/skill:` user-input-only expansion) with a `pi >= 0.83` floor enforced in ready-gate/doctor via `pi --version` at binary level — org moved to `@earendil-works`, old-org installs flagged; (4) one-line note naming `expandPromptTemplates` on `pi.sendUserMessage()` as the known upgrade path for programmatic skill expansion. Route-dimension addendum (live-observed via pi's skill-conflict report): each package registers skills through two routes (manifest `pi.skills` scan + the extension's `resources_discover` re-registering the same directory its own manifest declares), so duplicates are packages × routes (up to 4 per skill name) — the doctor count must be route-aware, and the cleanest remedy is `resources_discover` skipping manifest-covered paths. Correction package assembled (`/tmp/s98-revise-context.md`); dispatch withheld on captain instruction. **Fold applied 2026-09-11 (revision 2):** all four additions are now folded into this body — scope item 1 (dev-override arm), scope item 2 (route-aware duplicate count, double-extension flag, `resources_discover` manifest-skip remedy, `pi >= 0.83` binary-level floor with old-org diagnostic), Proposed approach (declared pi-behavior dependency + deterministic duplicate answer), AC-3/AC-4/AC-5/AC-7/AC-8, the expected-surface table, and doc diffs 1–2. Verification evidence for the new claims is recorded in the revision-2 stage report below.

## Expected surface

Net LOC change: **+270 across 6 files** (insertions +327 / deletions −57), tolerance ±60 net lines / ±2 files:

| File | +ins | −del | Content |
|---|---|---|---|
| `internal/cli/pi.go` (893) | +70 | −14 | drop `piBootstrapPrompt` + prompt-append reshape; `PI_SPACEDOCK_LAUNCH=1` env + wrap `--env-pass`; ready-gate adds installed-ext Stat + `firstOfficerDiscoverable` + `pi --version` floor hook; route-aware duplicate detection + dev-override double-extension warning in `piSpacedockPackageStatus` |
| `.pi/extensions/spacedock.ts` (~130) | +16 | −9 | marker gate on both injections; bootstrap wording; `resources_discover` manifest-skip; header-comment vocabulary |
| `skills/first-officer/references/first-officer-shared-core.md` | +10 | −2 | contract-level declaration + gate step |
| `docs/site/contributing/adding-a-runtime.md` | +60 | −20 | corrected load-paths section + new bootstrap subsection (duplicate-extension load, two registration routes, version floor) |
| `internal/cli/pi_frontdoor_test.go` + `pi_launch_test.go` | +145 | −9 | argv, gating env (incl. dev-override arm), duplicate/double-extension warnings, doctor report, version-floor tests |
| new bun test for the extension module | +26 | −3 | extension injection gate + `resources_discover` manifest-skip unit tests |

Observable semantics this task may change (the boundary the implementation must not cross silently): the frontdoor inner argv shape (a fresh launch no longer ends with the contract sentence; it ends with the operator task or nothing); the launched child env (new `PI_SPACEDOCK_LAUNCH=1` on every `spacedock pi` launch, wrap and non-wrap); the extension injection surface (bootstrap and boot-record injections now fire only for frontdoor-launched sessions); the FO startup step (the binary gate gains a contract-level assertion that can abort); doctor/launch output (new first-officer + installed-extension checks that refuse, a `pi >= 0.83` floor that refuses sub-floor binaries with an old-org diagnostic, and duplicate-registration / dev-override double-extension warnings that do not). No stored formats, command grammar, or authority semantics change.

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
> - **Duplicate registration.** When two registered packages both provide `first-officer` (e.g. a dev-link checkout and the installed git package), pi's skill scan is first-match: the first `settings.json` entry wins. Each package also registers its skills through two routes (the manifest's `pi.skills` scan and the extension's `resources_discover`, which skips paths the manifest already declares). `spacedock doctor --host pi` and launch output flag the condition with a route-aware count; the remedy is `pi remove` of the stale entry.
> - **Dev-override double extension.** Under the dev override (`SPACEDOCK_REPO_ROOT` / `--plugin-dir` install source), the checkout's `spacedock.ts` loads via `--extension` while the installed package's `spacedock.ts` loads via package discovery; `spacedock doctor --host pi` and launch output flag the double load.
> - **Version floor.** The mechanism requires `pi >= 0.83` (context-hook injection; `<available_skills>` with absolute per-skill locations; `/skill:` user-input-only expansion). The ready gate and doctor read the binary's `pi --version`; a sub-floor install is refused at launch and flagged by `doctor` (pi now publishes as `@earendil-works`; the old `@mariozechner` line is stale).

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

**Upgrade-path note (not scope).** `expandPromptTemplates` on `pi.sendUserMessage()` (new in the 0.83–0.85 range) is the known future path for programmatic skill expansion: the extension could deliver the skill body itself instead of pointing at the `<available_skills>` listing. Recorded here so the follow-up has a name; no code in this entity touches it.

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

## Stage Report: ideation — revision 2 (captain-directed fold, 2026-09-11)

- DONE: Dev-override arm coverage folded
  Scope item 1 gains the dev-override sentence (gate set with `repoRoot` set; stale checkout intentionally not supported, AC-6 territory); AC-3 gains the child-env unit assertion (`PI_SPACEDOCK_LAUNCH=1` present with `repoRoot` set, wrap and non-wrap); AC-4 gains the one live probe of the gated extension under `--extension <checkout>/…` against a CURRENT checkout.
- DONE: Duplicate-extension load stated deterministically at packages × routes granularity
  Scope item 2 now: route-aware duplicate count (each package contributes two routes — manifest `pi.skills` scan + extension `resources_discover` on the same directory, live-observed as doubled loser entries); double-extension load under dev override flagged at launch and in doctor on the static, code-visible condition (`repoRoot` set AND a spacedock package registered), because `hasStructuralBootstrap` provably covers only re-fires within one hook instance and pi's hook-chaining of modified messages between two hook instances is unverified; the cleanest remedy adopted — `resources_discover` skips paths the package manifest already declares, with a bun-test assertion (AC-4).
- DONE: pi-behavior dependency declaration + `pi >= 0.83` floor folded
  Approach gains the declared dependency list (context-hook API: `session_start` / `session_compact` / `context` returning modified messages; `<available_skills>` with absolute per-skill locations; `/skill:` user-input-only expansion) with the 0.83.0–0.85.1 changelog survey as basis. Floor enforced in ready-gate (AC-5) and doctor (AC-7) parsed from the binary's `pi --version` — never package paths; sub-floor installs flagged as the stale `@mariozechner`-org signal with an `@earendil-works` upgrade diagnostic. Doc diff 2 gains the version-floor bullet.
- DONE: Upgrade-path note recorded (one line, not scope)
  `expandPromptTemplates` on `pi.sendUserMessage()` named in the doc-diff section as the known future path for programmatic skill expansion.
- DONE: Expected surface re-estimated and doc diffs extended
  Surface now +270 net (+327/−57) across the same 6 files (floor hook, route-aware count, manifest-skip, and their tests added to the rows); observable-semantics paragraph gains the floor refusal + old-org diagnostic; doc diff 1 wording unchanged, doc diff 2 gains the double-extension and floor bullets, AC-8 now counts five documented facts.
- VERIFIED (new-claim evidence, no spike re-run needed): `pi --version` prints a bare semver (`0.85.1` on this machine) — the floor parse is trivial and binary-level; the route doubling is confirmed from source (package.json `pi.skills: ["./skills"]` vs `spacedock.ts` `resources_discover` returning the same `skills/` dir); the dev-override double-extension load is a static consequence of pi.go's dev-override argv block (adds `--extension <checkout>/spacedock.ts` whenever `repoRoot` is set, independent of settings.json registration). The dev-override live probe and the floor tests are specified as implementation-stage AC proofs; the existing spike results stand unchanged.

### Summary (revision 2)

Folded the captain's four revise asks into the entity body without redesign: dev-override arm coverage (env assertion + one live probe, stale checkout explicitly out), deterministic duplicate-extension answer at packages × routes granularity (route-aware doctor count, static double-extension flag, `resources_discover` manifest-skip as the cleanest remedy), the declared pi-behavior dependency list with a `pi >= 0.83` binary-level floor and old-org diagnostic, and the `expandPromptTemplates` upgrade-path note. Surface re-estimated to +270 net across 6 files; doc diffs 1–2 extended; no new spike needed — every new claim was verifiable from source or a one-line binary check.

## Stage Report: implementation

- DONE: The four coupled changes are implemented per the body's design with each AC's named proof owner built
  Commit b6d1dc073 on spacedock-ensign/pi-contract-bootstrap-single-owner (worktree). Frontdoor argv-only + `PI_SPACEDOCK_LAUNCH=1` env/wrap `--env-pass` + ready-gate additions (extension Stat, first-officer discoverable, `pi --version` 0.83.0 floor) + route-aware duplicate count + dev-override double-extension warning in `internal/cli/pi.go`; marker gate + resolvable wording + `resources_discover` manifest-skip in `.pi/extensions/spacedock.ts`; contract-token gate in `skills/first-officer/references/first-officer-shared-core.md:9`; doc diffs 1–2 applied as written in `docs/site/contributing/adding-a-runtime.md`.
- DONE: Proofs land as the ACs name them — deterministic artifacts, not model self-report
  AC-1/AC-4: bun tests `.pi/extensions/spacedock.test.ts` (no-marker → no injection; marker → injection with resolvable wording, no `$spacedock:`, no relative path; boot-record same gate; `resources_discover` manifest-skip) — `bun test ./.pi/extensions/spacedock.test.ts` 4/4 pass. AC-3/AC-5/AC-7: Go tests in `internal/cli/pi_frontdoor_test.go` (TestPiFrontDoorLaunchesWithNativeResourcePaths task-only prompt, TestRunPi_SetsLaunchMarkerEnvOnEveryLaunchShape dev-override arm, TestRunPi_ReadyGateRequiresSpacedockExtension/RequiresFirstOfficerSkill/RefusesSubFloorPiVersion, TestRunPi_WarnsOnDuplicateSpacedockRegistration/WarnsOnDevOverrideDoubleExtension, TestPiDoctorReportsFirstOfficerVersionDuplicates, TestPiSpacedockPackageStatus_DoubleRegistration, TestPiVersionAtLeast) and `pi_launch_test.go` (TestRunPi_WrapCarriesLaunchMarker). `go test ./internal/cli/` green except two pre-existing failures (see deferred-risk note below).
- DONE: AC-1 live probe (baseline + gated), AC-2 live probe, AC-4 dev-override live probe — deterministic artifacts under /tmp/s98-probe/out/ (isolated PI_CODING_AGENT_DIR home, probe extension dumping systemPrompt + context messages)
  AC-1: `a1-baseline-old-context.txt` — the INSTALLED (ungated) extension injects `[SPACEDOCK-FO-BOOTSTRAP-v1]` into a plain `pi` session (baseline leak reproduced); `a1-gated-nomarker-context.txt` — the worktree gated extension injects NOTHING without the marker. AC-2: `a2b-system-prompt.txt` — `<available_skills>` lists `first-officer` with absolute location `<worktree>/skills/first-officer/SKILL.md` from cwd /tmp/s98-probe/cwd, and `a2-context.txt` carries the expanded `<skill name="first-officer" location="…worktrees/…/skills/first-officer/SKILL.md">` block for `/skill:first-officer` user input; each skill listed exactly once (manifest-skip holding live). AC-4 dev-override arm: `a4-marker-context.txt` — the gated extension loaded via `--extension` (the dev-override argv shape) injects the NEW wording verbatim only when `PI_SPACEDOCK_LAUNCH=1` is set. End-to-end: freshly built worktree binary `spacedock pi … -- -p --no-session --no-tools` from /tmp cwd launched green through the new ready gate (fd-stdout.txt: session replied OK).
- DONE: AC-6 stale-fixture abort run — one-off live falsification
  Fixture `/tmp/s98-probe/stale-fixture` (shared-core declared level edited to contract level 2) against the worktree binary printing `contract 3`: exercised session output (`ac6-stdout.txt`) ends `GATE-OUTCOME: ABORT — contract-level mismatch … "Contract level mismatch: this skill declares contract level 2, but the binary … reports contract 3. Run spacedock doctor and update the plugin."`; positive control with the CURRENT shared core ends `GATE-OUTCOME: PASS … contract 3, matching this contract's declared level 3`.
- DONE: Deliverable committed to the worktree branch with this stage report naming the proof artifacts
  Commit b6d1dc073; entity report appended here; state committed path-scoped and pushed.
- DONE: (deferred risk, not this entity's promise) two pre-existing internal/cli test failures observed during the run — defect recorded with promote-to-material condition
  Both failures are real defects outside this entity's promises, classified as deferred risk, not Material: TestVersionAmbiguousMarkersExitZero fails only with this subagent session's `PI_CODING_AGENT` env set (passes with `env -u PI_CODING_AGENT`) — an env-dependent trigger, not any change here; TestCodexResolveManifestAgainstInstalledHost fails identically on the clean HEAD (verified via `git stash`) — codex manifest resolution, untouched by this entity's diff. Promote-to-material condition: if validation attributes either failure to this entity's diff, or the codex manifest failure blocks validation's suite verdict, re-disposition it as Material under the review-finding policy.
- SURFACE DRIFT (flagged, not absorbed): net +783 (+876/−93) across 8 files vs declared +270 net (+327/−57) across 6 (tolerance ±60 net / ±2 files)
  Files 8 vs 6 is within ±2 (the 7th, `internal/cli/safehouse_frontdoor_test.go` +3/−3, is a compile-only edit — TestBootstrapPromptsDropEngageFlourish referenced the removed `piBootstrapPrompt` const; pi dropped from that flourish table, no semantic change). The LOC drift is tests + mechanism-contract comments: pi test files +493/−36 (vs +145/−9 declared) for the AC-3/AC-5/AC-7 proof matrix (8 test funcs incl. subtests), new bun test 125 lines (vs +26 declared), pi.go +187/−43 (vs +70/−14) mostly contract comments. Doc row came in UNDER estimate (+16/−3 vs +60/−20). No undeclared semantic change: the one observable delta beyond the declared set is the dev-override ready-gate tightening (a dev-override checkout without `.pi/extensions/spacedock.ts` now REFUSES instead of silently launching skill-less — AC-5a's refusal, test renamed accordingly).

### Summary

Implemented the single-owner FO contract bootstrap end to end and produced every AC's named proof artifact. The live probes closed the loop the spike opened: the old ungated extension still leaks into plain `pi` sessions (baseline re-measured), the gated worktree extension does not without the marker and does with it, the registered package's skills list resolvably from a foreign cwd with the manifest-skip holding (one registration per skill), and a stale declared contract level aborts at the binary gate naming both versions while the current one passes. Surface drift past the declared tolerance is flagged above with its composition; the only undeclared observable semantic (dev-override without extension now refuses) is AC-5a's own refusal semantics, not a silent widening.

## Stage Report: validation

- DONE: Every AC checked against its declared proof, re-run from this run's evidence — not accepted from the implementation report
  Bun extension tests re-run fresh: `bun test ./.pi/extensions/spacedock.test.ts` 4/4 pass (no-marker→no injection; marker→resolvable wording without `$spacedock:`/relative path; boot-record same gate; `resources_discover` manifest-skip). Falsification spot-check: flipping the gate to `return true` flips the suite to 2 fail/2 pass; revert restores 4 pass — the gate test is not tautological. Entity Go test set re-run green: TestPiFrontDoor* (task-only argv), TestRunPi_SetsLaunchMarkerEnvOnEveryLaunchShape (installed + dev_override), TestRunPi_ReadyGateRequiresSpacedockExtension/RequiresFirstOfficerSkill/RefusesSubFloorPiVersion (0.82.9 refuse / 0.83.0 / 0.85.1 / 1.0.0 / dev / garbage), TestPiVersionAtLeast, TestRunPi_WarnsOnDuplicateSpacedockRegistration, TestRunPi_WarnsOnDevOverrideDoubleExtension, TestPiDoctorReportsFirstOfficerVersionDuplicates (3 subtests), TestPiSpacedockPackageStatus_DoubleRegistration, TestRunPi_WrapCarriesLaunchMarker — all PASS.
- DONE: Live probes re-run fresh under /tmp/s98-val/out/ (isolated home, foreign cwd /tmp/s98-val/cwd, deterministic probe-extension dumps, never model self-report)
  AC-1: `val-a1-baseline-context.txt` — installed UNGATED extension still injects `[SPACEDOCK-FO-BOOTSTRAP-v1]` into a plain `pi` session (baseline leak reproduced); `val-a1-nomarker-context.txt` — worktree gated extension injects NOTHING without the marker. AC-2: `val-a2-list-system-prompt.txt` — `<available_skills>` lists `first-officer` exactly once with absolute worktree location, from a foreign cwd; `val-a2-expand-context.txt` — `/skill:first-officer` USER INPUT expands to `<skill name="first-officer" location="…worktree…/skills/first-officer/SKILL.md">`. AC-4 dev-override arm: `val-a4-marker-context.txt` — extension loaded via `-e` (the `--extension` argv shape) injects the NEW wording only with `PI_SPACEDOCK_LAUNCH=1`; `val-a4-nomarker-context.txt` — nothing without it.
- DONE: AC-6 re-run live — stale-fixture ABORT plus PASS control, fresh fixture at /tmp/s98-val/stale-fixture (declared level edited to contract level 2) against the freshly built worktree binary (contract 3)
  `val-ac6-stdout.txt` ends `GATE-OUTCOME: ABORT — … reports contract 3 while the first-officer skill declares contract level 2`; control with the CURRENT shared core ends `GATE-OUTCOME: PASS — … contract 3, matching the skill's contract level 3`.
- DONE: Full `go test ./internal/cli/` read and both flagged failures triaged against the deferred-risk note
  Suite with `PI_CODING_AGENT` unset: exactly ONE failure, TestCodexResolveManifestAgainstInstalledHost — verified failing IDENTICALLY on pre-entity HEAD af70297dd (detached worktree at /tmp, since removed): pre-existing codex-manifest defect, untouched by this diff, stays deferred risk. TestVersionAmbiguousMarkersExitZero: PASSES with `PI_CODING_AGENT` unset, fails only with the subagent env set (also verified failing on pre-entity HEAD with the env set) — env-dependent trigger, pre-existing, stays deferred risk. Neither is attributable to this entity's diff; neither promotes to Material.
- DONE: Observable semantics verified end-to-end on real launches from a foreign cwd
  Freshly built worktree binary: `spacedock pi "Reply OK" -- -p --no-session` from /tmp/s98-val/cwd → exit 0, session replies OK (`val-fd-stdout.txt`), through the new ready gate (extension Stat + first-officer discoverable + `pi --version` 0.85.1 ≥ 0.83.0 floor). Against the REAL home the same launch prints the live duplicate-registration WARNING naming both roots, the first-match winner, the route-aware ceiling, and the `pi remove` remedy (`val-fd-real-stderr.txt`). Sub-floor live probe (fake `pi` printing 0.73.1 first on PATH): launch REFUSED at the ready gate, `MISSING pi version: 0.73.1 (floor 0.83.0)` with the @earendil-works/old-org remedy, exit 1 (`val-subfloor-stdout.txt`). `spacedock doctor --host pi` live: first-officer line, `OK pi version: 0.85.1 (floor 0.83.0)`, route-aware duplicate WARN + remedy. Doc diffs 1–2 match shipped behavior verbatim (wording identical to the entity body; stale load-paths text corrected; the remaining `--skill` lines live in the separate live-smoke recipe, intentionally out of scope).
- DONE: Surface drift adjudicated — +783 net (+876/−93) across 8 files vs declared +270 (+327/−57) across 6, tolerance ±60 net / ±2 files
  Files 8 vs 6 within ±2 (7th = compile-only safehouse test edit, 8th = the bun test the table counted as one row). LOC past tolerance; numstat composition: `internal/cli/pi_frontdoor_test.go` +450/−36 (declared +145/−9) and bun test 125 lines (declared +26) — the AC-3/AC-5/AC-7 proof matrix the revise cycle demanded, plus mechanism-contract comments; `pi.go` +187/−43 (declared +70/−14), mostly comments; doc row came in UNDER (+16/−3 vs +60/−20). No undeclared semantic surface: the one observable delta beyond the declared set is the dev-override ready-gate tightening (a dev-override checkout without `.pi/extensions/spacedock.ts` now REFUSES) — that IS AC-5a's refusal semantics, not a widening. Recommendation: captain-visible surface note, not a design reset — the excess mass is test/comment, every AC's named proof is present and falsifiable.
- DONE: Verdict — PASSED recommended; no Material findings; deferred risks recorded with promote conditions
  All eight ACs PASSED with fresh evidence (AC-8 via verbatim doc application paired with AC-1/2/6 live value). Deferred risks: (1) two pre-existing internal/cli failures above; (2) the INSTALLED package's extension is still the pre-fix ungated copy, so plain `pi` sessions on this machine keep receiving the bootstrap until the package updates — deployment lag, not a diff defect; (3) pi 0.85.1 lists NO skills at all under `--no-tools` (observed during validation: `available_skills` absent with the flag, present without) — no supported launch shape passes `--no-tools`; promote if one ever does. Polish (non-blocking): doctor/launch "up to 4 registrations" states the pre-manifest-skip ceiling; with the skip active the actual is packages × 1; and the doc's "Three facts" lead-in introduces five bullets (the entity's own approved wording).

### Summary

Re-ran every AC's proof from fresh evidence: deterministic suites (bun 4/4 including a falsification spot-check; the full entity Go test set green), live probes under /tmp/s98-val (AC-1 baseline leak + gated silence, AC-2 listing/expansion from a foreign cwd, AC-4 dev-override arm both polarities, AC-6 stale-fixture ABORT + PASS control), and real launches (frontdoor green through the new gate, live duplicate warning, live sub-floor refusal, live doctor report). Both pre-existing test failures were re-triaged against pre-entity HEAD and stay deferred. Surface drift past tolerance is real but composed entirely of the revise-cycle's test matrix and contract comments; verdict PASSED with no Material findings.
