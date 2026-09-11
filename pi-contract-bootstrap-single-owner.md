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
---

Pi sessions receive the FO contract through two channels with unclear ownership, and neither is reliable today. The frontdoor launch prompt (`Use $spacedock:first-officer for this whole Pi session.`, `internal/cli/pi.go:20`) is inert syntax pi cannot expand (`agent-session.js:_expandSkillCommand` only expands `/skill:`), and the extension's session-start contract bootstrap (`FO_BOOTSTRAP_TEXT`, `.pi/extensions/spacedock.ts`) names the skill with the same unexpandable reference plus a relative path that is ENOENT from any workflow cwd. The failure is not delivery — it is resolution: even with both injections landing, an FO outside the package root ends up hunting the filesystem, and a stale visible checkout wins. Nothing in either contract verifies the loaded skill's version against the binary.

## Scope (one entity, four coupled changes)

1. **Extension-owned bootstrap.** The frontdoor launch prompt drops the contract sentence entirely (task passing and the `containsResume` suppression semantics stay). The extension's bootstrap becomes the single owner: wording switches to the pi-resolvable `/skill:first-officer` trigger (no `$pkg:skill`, no relative path), and injection is gated on a frontdoor launch marker (e.g. `SPACEDOCK_BIN` passed through the wrap, or a dedicated `PI_SPACEDOCK_LAUNCH=1`) so plain `pi` sessions for unrelated work stop receiving it.
2. **Launch-time health check.** The frontdoor already computes package/extension resolvability (`spacedockPackageOK`, extension os.Stat guards). When the extension cannot be resolved at launch, warn loudly (or refuse) instead of silently launching a contract-less session — the failed-launch case the old prompt could not actually cover.
3. **FO contract version self-check.** Startup (shared core, at the binary gate) asserts the loaded skill's contract version against the binary's (`spacedock --version` already prints `contract N`); on mismatch the FO refuses to proceed with a loud diagnostic. This is the catch for any residual stale-contract path regardless of how it was loaded.
4. **Docs + vocabulary (folded, no doc-only entity).** Site reference (`docs/site/contributing/adding-a-runtime.md` or a new reference section) documents: the single-owner bootstrap design; the compaction contract (boot-record re-read per #738, never contract re-injection); that context-hook injections are request-time-only and never appear in session logs; and the duplicate-registration hazard (two registered packages both shipping `first-officer` — expansion is first-match). The extension header and marker text stop saying "commissions the parent session" — the system vocabulary reserves commissioning for workflows; this installs the FO contract.

## Acceptance criteria (sketch — ideation fleshes out, externally proven)

- A plain `pi` session in an unrelated cwd receives NO contract bootstrap (the model reports the marker absent from context).
- A frontdoor-launched `spacedock pi` session receives the bootstrap exactly once via the extension, with `/skill:first-officer` wording that resolves the installed package's skill from an arbitrary workflow cwd (proven by a live session actually loading the 0.28 contract from a non-package-root cwd).
- No frontdoor launch emits the `$spacedock:first-officer` sentence; a task given to `spacedock pi <task>` still reaches the session; resume launches still do not re-bootstrap.
- Frontdoor launch with an unresolvable extension/package produces a loud warning or refusal, observable in launch output.
- A session whose loaded skill contract version mismatches the binary aborts at startup with a diagnostic naming both versions (falsifiable by pointing a session at a stale checkout).
- Doctor/launch flags the duplicate-registration condition (two packages providing `first-officer`).
- The reference doc section exists and states all four documented facts above; no standalone docs entity was filed.

## Out of scope

- Dev-override (`--plugin-dir` / `SPACEDOCK_REPO_ROOT`) version guard: latent hazard, not the mechanism in play in the 2026-09-10 incident; candidate for its own entity if the captain wants it.
- Safehouse profile `add-dirs` grant hygiene (stale grants pointing at old checkouts): doctor-note material at most; visibility is not resolution.
- The `gate record` exit-0-with-frozen-error quirk observed once mid-incident: re-verify against the current binary before filing anything.
