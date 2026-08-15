---
id: 71btbxdrken4kdmfsk0vptav
title: Trim stale prose in the dispatch core contract
status: implementation
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
started: 2026-08-15T02:55:45Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-trim-dispatch-core-stale-prose
issue:
gates:
    version: 1
    records:
        - id: gate:71btbxdrken4kdmfsk0vptav:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:71btbxdrken4kdmfsk0vptav-backlog-1
              briefing:
                id: briefing:71btbxdrken4kdmfsk0vptav:backlog:attempt-1:revision-1
                digest: sha256:665a1ebd3a9258ea6dad9a58d7a1494000d82e6a2ed6b2d0a5f000c540943df5
                request-digest: sha256:b7a73f6a9d9f568a7d9c13110daf4286d155d40a011052c0b330e3434fe76b82
                room-ref: ./trim-dispatch-core-stale-prose/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:71btbxdrken4kdmfsk0vptav:backlog:1
                briefing: briefing:71btbxdrken4kdmfsk0vptav:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:55.975611Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:71btbxdrken4kdmfsk0vptav:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:71btbxdrken4kdmfsk0vptav-ideation-1
              briefing:
                id: briefing:71btbxdrken4kdmfsk0vptav:ideation:attempt-1:revision-1
                digest: sha256:b7dff507bf6cd8f42fda606661e7b10ec4d857d6cc71f9a2ec7ea1897687b836
                request-digest: sha256:196941c9aecb1dc3a0424d313cda8f6d6871ad6ef300956c0d088a9691f92300
                room-ref: ./trim-dispatch-core-stale-prose/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:71btbxdrken4kdmfsk0vptav:ideation:1
                briefing: briefing:71btbxdrken4kdmfsk0vptav:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T03:55:23.282375Z"
                decision: approve
                reason: 'Captain ruling 2026-08-15 (approve all except x8): approved into implementation'
              application:
                target-stage: implementation
                state: consumed
---

Two prose repairs in skills/first-officer/references/fo-dispatch-core.md.

1. Line 195: the replay sentence references "s4", a private dev-workflow entity ID no installed reader can resolve. Rewrite to name the behavior: `spacedock gate prepare` is replay-idempotent. Keep the prohibition (no attempt counter, retry token, cache, or alternate authority).
2. Lines 197 and 199: the stale abbreviated --next envelope shape plus the sentence that exists only to retract it. Fold the canonical ready_gates shape into line 197 and delete line 199. The canonical rule at line 191 already governs.

## Problem

`skills/first-officer/references/fo-dispatch-core.md` is the FO's lazy-loaded dispatch contract. Three lines in its Event Loop section assert things that are either unresolvable for every reader or untrue of the shipped binary. All three were confirmed against HEAD 4d1912a69 by exercising, not by reading.

**1. Line 195 credits a dead referent.** It attributes the gate-prepare replay idempotency to "s4", a dev-workflow display ID. `spacedock status --workflow-dir docs/dev --resolve s4` returns `Error: unknown reference: s4` — the ID no longer resolves inside the workflow that minted it, and it never resolved for an installed reader. The behavior it credits is real and is owned by `gate prepare`: `Prepare` short-circuits through `preparedReplay` (`internal/gates/prepare.go:165` and `:338`), returning the prior room, briefing, and digest with `state=open` when the briefing ID, question, summary, artifact/reference ordering, and every selected source's logical revision all match. The rule the sentence carries — FO creates no attempt counter, retry token, cache, or alternate authority — is load-bearing and must survive untouched.

**2. Line 197 states two machine envelope shapes, and both are stale.** Exercised against the 0.27.0-pre4 binary:

- `--next` is documented as `{"command":"next","dispatchable":[…]}`, but `nextJSON` (`internal/status/json_commands.go:97-99`) emits `command`, `dispatchable`, **and** `ready_gates`.
- `status`/`--where` is documented as `{"command":"status","entities":[…]}`, but `statusJSON` (`internal/status/json_commands.go:64`) emits `command`, `entities`, **and** `pagination`.

The seed filing named only the `--next` half. The `status` half has the identical defect and no retraction sentence covering it, so it is the quieter of the two errors.

**3. Line 199 is a retraction, not a rule.** It exists only to tell the reader that line 197 is wrong ("supersedes the abbreviated shape in this historical sentence"). Line 191 already states the canonical `--next` shape correctly. Once line 197 is corrected, line 199 carries nothing.

## Proposed approach

Three prose edits in one file. No code, no test, no new mechanism.

**Edit 1 — line 195, swap the dead referent for the command that owns the behavior.**

```text
BEFORE: An exact prior-stage replay candidate may exercise s4's existing
        selection/replay idempotency during that one invocation; FO creates no
        attempt counter, retry token, cache, or alternate authority.

AFTER:  An exact prior-stage replay candidate may exercise `gate prepare`'s existing
        selection/replay idempotency during that one invocation; FO creates no
        attempt counter, retry token, cache, or alternate authority.
```

Only the possessive changes. "existing selection/replay idempotency" and the whole prohibition clause stay byte-identical, so meaning preservation is provable by inspection rather than argued. Bare `` `gate prepare` `` (not `spacedock gate prepare`) matches the section's local convention — line 193 already writes "exactly one existing `gate prepare`", and line 185 writes `gate consume` and `dispatch build --stamp` the same way.

**Edit 2 — line 197, state both envelope shapes as the binary emits them.**

```text
BEFORE: … Envelopes: `status`/`--where` → `{"command":"status","entities":[…]}`;
        `--next` → `{"command":"next","dispatchable":[…]}`. …

AFTER:  … Envelopes: `status`/`--where` → `{"command":"status","entities":[…],"pagination":{…}}`;
        `--next` → `{"command":"next","dispatchable":[…],"ready_gates":[…]}`. …
```

The rest of line 197 is unchanged.

**Scope note for the gate — the `pagination` half is an addition to the seed filing.** The seed directed only the `--next` fix. Correcting the `status` half in the same sentence is prose-only, costs 17 bytes, and is what makes "the envelope shape is stated correctly" true rather than half-true; leaving it invites the next audit to file the same line again. If the gate prefers strict seed scope, drop this clause: the file still lands at −129 B instead of −110 B and AC-2 narrows to the `--next` envelope alone. Recommendation: include it.

**Edit 3 — delete line 199 and its blank separator.** The retraction has nothing left to retract.

## Out of scope

Any behavior change. The canonical envelope rule and its Go emitters. Also explicitly out of scope, and carried forward as a separate finding below: the FO contract has no pagination handling at all.

## Finding for a separate entity (not fixed here)

`status --where … --json` paginates at `defaultPageLimit = 25` (`internal/status/format.go:91`), and no FO contract text anywhere in `skills/` reads `pagination` or `has_next` — the only mentions of `--page`/`--limit` are in `skills/fo-status-viewer/SKILL.md:46` for the human overview table. Event-loop step 2 (`status --where "mod-block !=" --json`, line 206) says "For each row, re-read the blocking mod" and would silently miss rows 26+ on a workflow with more than 25 mod-blocked entities. Naming `pagination` in the envelope shape does not fix this; the fix is a contract rule, which is a behavior change this entity must not make. Recommend filing separately.

## Expected surface and tolerance

One file: `skills/first-officer/references/fo-dispatch-core.md`. Net −2 lines / −110 bytes (measured on the applied spike, not estimated); tolerance ±1 line and ±40 bytes for wording adjustments the gate requests. No other file changes.

Declared semantic changes: none. Command grammar, stored formats, authority, and runtime behavior are all unchanged — this edits contract prose only, and every clause it touches either keeps its bytes or is brought into agreement with an emitter that itself does not move. No user-visible surface changes (no CLI output, command surface, startup banner, host integration, or docs-site content), so no doc diff is owed.

## Spike record

The riskiest unverified assumption was that these lines are load-bearing to a contractlint check or the prose-function routing binder, so that editing them turns something red. It was exercised, not reasoned about. All three edits were applied to a throwaway detached worktree at HEAD 4d1912a69 (`scratchpad/spike-trim-dispatch-core-stale-prose`, verified clean before the copy, one-file diff after):

- `go build ./...` exit 0.
- `go test ./internal/contractlint/` ok, 3.678s.
- `go test ./internal/cli/ -run ProseFunction` ok, 0.711s — the binder that parses this file for `«fn»` migration targets. It scopes extraction to `→ **shipped**`/`→ **prose**` lines; none of lines 195/197/199 carry one, which is why the edit cannot reach it.
- `go test ./internal/gates/ -run Replay` ok, 18.999s — `TestPrepareReplaySurvivesRequiredStateCommit`, `TestPrepareReplayAcceptsEntitySelectedAsArtifactOrReference`, and the two replay refusal tests, which are the existing proof that the behavior edit 1 names is real.

Also established during the spike: no vendored second copy of the file exists (`find` returns one shipped path plus other tasks' worktrees); `prose_manifest_minor_sync_test.go` pins `first-officer-shared-core.md`, not this file, so no manifest bump is owed; `TestFOInstructionComponentCaps` leaves this file uncapped, so there is no byte cap to satisfy and the value AC must be a measured delta.

One observation that is **not** evidence and is recorded so it is not mistaken for a result later: a first attempt at this spike used the obvious shared path `scratchpad/spike` and collided with a concurrent sibling ensign working in the same session scratchpad. `go test ./internal/cli/` in that checkout hit the 600s per-package cap, but the checkout also carried that sibling's in-flight `internal/gates`/`internal/status` deletions, so the timeout is unattributable to anything here. It was discarded and the spike redone in a uniquely named worktree. Full-package `internal/cli` timing is not this entity's risk.

## Acceptance criteria

**AC-1 (value; independent baseline that can move the wrong way) — the dispatch contract is smaller than it was.**
`skills/first-officer/references/fo-dispatch-core.md` is under 29210 bytes and under 223 lines, the `origin/main` baseline at filing (working tree confirmed byte-identical to `origin/main`), and the signed byte delta is NEGATIVE. This can genuinely move the wrong way: edits 1 and 2 together ADD 48 bytes, and only deleting the retraction sentence and its separator (−158 B) makes the total negative — a worker who "fixed the wording" and kept line 199 fails this AC. Verified by reporting post-edit `wc -c` and `wc -l` and the signed delta against 29210 / 223. Spike measurement: 29100 B / 221 lines, delta −110 B / −2 lines.

**AC-2 (value; exercised against the binary) — every envelope shape line 197 documents equals what the binary actually emits.**
Verified by running `spacedock status --workflow-dir docs/dev --next --json` and `spacedock status --workflow-dir docs/dev --where "mod-block !=" --json`, taking each response's top-level key set, and comparing it to the key set parsed out of the rewritten line 197. Both must match exactly: `command,dispatchable,ready_gates` and `command,entities,pagination`. The check discriminates — run against HEAD it reports `match=False` for both envelopes; against the proposed text, `match=True` for both. It fails if the prose drifts or if an emitter gains or loses a top-level key.

**AC-3 — the rule survives the referent swap, and the new referent resolves.**
`grep -rn '\bs4\b' skills/` returns nothing (baseline: exactly one hit, `fo-dispatch-core.md:195`); the clause "FO creates no attempt counter, retry token, cache, or alternate authority" is byte-identical to `origin/main`; and the command the sentence now names routes — `spacedock gate prepare --help` exits 0. The last clause is what makes this more than a deletion: it fails if the swap trades one unresolvable referent for another.

**AC-4 — nothing outside the three target lines moves.**
`git diff --stat origin/main` shows exactly one changed file, and `git diff origin/main -- skills/first-officer/references/fo-dispatch-core.md` touches only lines 195, 197, and 199 plus 199's blank separator. Line 191's canonical scheduler-envelope rule is byte-identical.

**AC-5 — the guarded suites stay green.**
`go build ./...` exit 0; `go test ./internal/contractlint/` green (reference closure, FO reference topology, component caps, and the initial-worker spawn-guard anchors that read this file); `go test ./internal/cli/ -run ProseFunction` green (the prose-function routing binder); full `go test ./...` green.

## Test plan

No new tests and no new mechanism. Every AC is checked with an existing suite or a one-off exercised comparison, which is the right cost for a prose repair — a durable check would be new machinery in a program whose whole purpose is removing machinery whose value chain is broken.

- AC-1: two `wc` calls against recorded baselines. Seconds.
- AC-2: the exercised key-set comparison above, run against the built binary during implementation and re-run at validation. Seconds. **Deliberately not made a permanent contractlint case.** The simplest alternative to a new lint is exactly this manual exercised check, and it is sufficient because the emitter side is already pinned by `internal/status/testdata/golden/seq-next.json` and `gate_readiness_needs_preparation_test.go`; the only unpinned side is prose that changes about once a year. A new lint would add a permanent test to guard a sentence, which is the cost this program exists to avoid.
- AC-3, AC-4: `grep`, `git diff`, and one `--help` exit code. Seconds.
- AC-5: existing Go suites, already run green on the spike; full `go test ./...` at implementation. Note that `internal/cli` is a slow package — allow a raised `-timeout` rather than reading the go default 600s cap as a failure.

No fixture, CLI-behavior, or live workflow tests are needed: nothing this entity changes is reachable by the binary at runtime.

## Stage Report: ideation

- DONE: Line 195 names gate prepare replay-idempotent and keeps the prohibition; the envelope shape folds into line 197 and line 199 goes
  Design records verbatim before/after for all three edits; applied and measured on a clean throwaway worktree at HEAD 4d1912a69 (one-file diff, 29210->29100 B, 223->221 lines). Line 195 swaps only `s4's` -> `` `gate prepare`'s ``, leaving "existing selection/replay idempotency" and the whole no-counter/token/cache/authority clause byte-identical.
- DONE: Contractlint reference-closure tests stay green
  With all three edits applied: `go test ./internal/contractlint/` ok 3.678s; `go build ./...` exit 0; `go test ./internal/cli/ -run ProseFunction` ok 0.711s. The routing binder scopes extraction to `→ **shipped**`/`→ **prose**` lines and none of 195/197/199 carries one, so a red there would mean the edit escaped its intended span.
- DONE: No-behavior-change declared: prose only, meaning preserved from canonical sentences
  Declared under Expected surface: no command grammar, stored format, authority, or runtime change; no user-visible surface, so no doc diff owed. Line 197's new text is not asserted to match prose, it is asserted to match what the binary emits: exercised key sets are `command,dispatchable,ready_gates` and `command,entities,pagination`, matching AC-2's parse of the proposed line (HEAD reports match=False on both, proposal match=True on both).
- DONE: Spike the riskiest unverified mechanism (stage-def requirement)
  Risk was "these lines are load-bearing to a lint or the routing binder." Exercised rather than reasoned: see Spike record. Also confirmed `spacedock status --resolve s4` = `Error: unknown reference: s4`, `spacedock gate prepare --help` exit 0, and `go test ./internal/gates/ -run Replay` ok 18.999s for the behavior line 195 now names.

### Summary

Confirmed all three defects against HEAD by exercising the binary, not by reading prose: `s4` is unresolvable even inside its own workflow, and line 197 misstates BOTH envelopes, not just `--next` (the `status` shape omits `pagination`). Scoped the `pagination` half as a flagged addition to the seed filing with a one-line fallback if the gate prefers strict seed scope. Rewrote the ACs so the value is measured rather than asserted: AC-1 is a signed byte delta against the origin/main baseline that goes the wrong way if a worker fixes the wording but keeps the retraction sentence, and AC-2 is an exercised key-set comparison that discriminates HEAD from the proposal. Two things carried out of scope deliberately: the FO contract reads no `pagination`/`has_next` anywhere, so event-loop step 2 silently truncates past 25 mod-blocked entities (recommended as a separate filing, since fixing it is a behavior change); and a first spike attempt collided with a sibling ensign on the shared `scratchpad/spike` path, so its 600s `internal/cli` timeout is unattributable and was discarded rather than reported as a result.
