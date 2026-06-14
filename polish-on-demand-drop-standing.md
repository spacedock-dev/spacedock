---
id: xf7fft1hnj51eq7kagsc9833
title: Move the standing-teammate mechanism out of the FO contract — the comm-officer mod self-injects as a standing teammate
status: ideation
source: "captain (2026-06-13, this session) — the standing-teammate lifecycle (discovery pass / lazy-spawn / declaration / team-scope teardown / first-boot-wins) is ~4 contract subsections of maintenance surface for an infrequently-used polisher; amortization doesn't pay for infrequent use. Captain chose approach A: on-demand one-shot polish dispatch, and 'the mod can add the prose for the standing team member' — feature-specific usage prose lives in the mod, the contract keeps only a generic hook. Taken into 0203."
started: 2026-06-14T18:42:01Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0203-fo-efficiency
---

Replace the standing-teammate lifecycle with on-demand one-shot polish, and move the feature's usage prose out of the FO contract into the mod that declares it.

## Problem

Supporting the standing prose-polisher (`comm-officer`) costs ~4 contract subsections — the shared-core `## Standing Teammates` concepts plus the runtime's discovery-pass / lazy-spawn / declaration-and-routing-mechanics (all relocated into `claude-fo-dispatch.md` by j9's P1 split) — plus the `using-claude-team` lifecycle notes. All of that machinery exists to buy ONE thing: amortization (keep the polisher resident so its skill loads once across many polishes). For an infrequently-used feature that trade does not pay — we maintain a long-lived, team-scoped, first-boot-wins, lazy-spawn, teardown lifecycle to save a skill-load that rarely happens.

## Mechanism spike (riskiest unknown — exercised FIRST)

The whole design rests on one claim: the one-shot polish dispatch can reuse `spawn-standing`'s prompt-extraction. I exercised it before designing the rest.

Run (against the live binary, with a team name that matches no on-disk config so the already-alive short-circuit does not fire):

    ./spacedock dispatch spawn-standing --mod "$PWD/docs/dev/_mods/comm-officer.md" --team "spike-fresh-team-no-such-config-9931"

Result: exit 0, and a complete Agent() spec JSON with `subagent_type: general-purpose`, `name: comm-officer`, `team_name`, `model: sonnet`, and `prompt` = the mod's `## Agent Prompt` body verbatim (the parser at `internal/dispatch/standing.go` reads frontmatter, the `## Hook: startup` spawn config, and the trailing `## Agent Prompt` section). The extraction this task repurposes is PROVEN: a one-shot `dispatch polish-spec` is this exact code path minus the `claudeteam.MemberExists` already-alive short-circuit (standing.go:144-154) and minus a real `--team` requirement (a one-shot is bare — no team needed).

**What the spike also surfaced (a real design decision, not a footnote):** the current `## Agent Prompt` is written for a RESIDENT teammate — "Then idle. Do NOT start polishing anything until you receive a polish request", "Stay live. Go idle between tasks." A one-shot polisher receives the draft IN its dispatch, polishes once, replies, and dies. So the one-shot path cannot reuse the resident Agent Prompt verbatim; it needs a one-shot framing. Decision below.

## Proposed approach

**Approach A — on-demand one-shot.** When the FO (or an ensign) has a deliberate draft to polish, it dispatches a one-shot polish agent only then; there is no resident teammate, no team-scope lifecycle. You pay one spawn per polish, which is rare.

### Binary: `dispatch polish-spec`

Add a subcommand `spacedock dispatch polish-spec --mod {abs_mod_path}` that emits a one-shot Agent() spec from a polish mod, repurposing the proven `spawn-standing` extraction:

- Parse the mod the same way `runSpawnStanding` does (frontmatter, `## Hook: startup` for `subagent_type`/`model`, prompt body). Reuse `ParseModMetadata` + `ParseHookStartupSpawnConfig`.
- DROP the `--team` requirement and the `MemberExists`/already-alive branch entirely — a one-shot is bare-mode, so the emitted spec OMITS `team_name` (matching how `dispatch build` omits `name`/`team_name` in bare mode). It carries `subagent_type`, `name`, `model`, `prompt`.
- The emitted `prompt` is the mod's **one-shot prompt body**, not the resident one. The mod gains a `## One-Shot Prompt` section (the resident framing minus "stay idle / stay live / wait for a request"; instead: "You receive a draft to polish below. Polish it, reply once with polished text + notes block, then stop."). `polish-spec` extracts that section; if the mod lacks `## One-Shot Prompt`, `polish-spec` exits non-zero with a diagnostic naming the missing section (loud failure, same discipline as the existing spawn-standing missing-section errors).
- Caller appends the draft to the emitted prompt (text-passthrough) or passes the absolute path (file modes) before forwarding to `Agent()`. The four polish modes survive unchanged in framing — the difference is only "spawned per-polish, dies after" vs "resident, idle between."

This is the smallest binary change that drops residency: a new subcommand wired in `dispatch.go`, a `runPolishSpec` beside `runSpawnStanding` in `standing.go`, sharing the mod-parse helpers. `spawn-standing`/`list-standing`/`show-standing` stay (other workflows may still declare a resident teammate); this task does NOT delete them — it adds the one-shot path and retargets `comm-officer` + the contract onto it.

### Mod becomes self-describing (the captain's key point)

The feature-specific usage prose moves OUT of the contract and INTO `comm-officer.md`:

- ADD `## One-Shot Prompt` (the per-polish Agent body, above).
- Frontmatter `standing: true` → DROP it. The mod is no longer a standing teammate; it is an on-demand polish mod. `list-standing` no longer returns it; `comm-officer` no longer appears in any boot discovery.
- DELETE `## Hook: startup` and `## Hook: shutdown` (residency lifecycle) — there is no startup spawn and no teardown.
- The mod KEEPS, and becomes the single home for, the usage prose the contract sheds: when to polish vs not (the Scope bullets), the four polish modes, the boundary rules, the hard rules. This already lives in the mod's `## Routing guidance` / `## Routing Usage`; the contract's duplicate goes.

### Contract keeps only a minimal generic hook

The FO/ensign contract keeps a few sentences: "Registered polish mods declare on-demand polishers. When you have a deliberate draft for captain review (PR body, gate summary, debrief, long entity-body narrative), you MAY route it through a one-shot polisher: `dispatch polish-spec --mod {path}` emits the Agent spec; append the draft and dispatch; it polishes once and dies. Best-effort, non-blocking — proceed un-polished if it does not return. Read the polisher's usage from the mod, not from here." The when/how detail lives in the mod.

## Concrete before/after wording

### `claude-fo-dispatch.md` — SHEDS

- **`### Standing teammate discovery pass`** (current lines 17-24) — DELETE entirely. No boot-time `list-standing` discovery for polish.
- **`### Standing teammate lazy-spawn`** (current lines 26-38) — DELETE entirely. No deferred spawn-at-first-dispatch.
- **`### Standing teammate declaration and routing mechanics`** (current lines 40-47) — DELETE the declaration-layout / teardown-trigger / lazy-spawn-injection bullets. The dispatch-time-injection bullet (the `show-standing` build-append) is retargeted: `dispatch build` no longer auto-appends a standing-teammates routing block.
- **`## Standing Teammates`** (current lines 49-56) — DELETE the four standing concepts (first-boot-wins, team-scope lifecycle, routing contract, declaration). The 3-sentence on-demand hook lives in `## Dispatch` instead.
- **`## Dispatch` → "Routing through a standing prose-polisher"** (current line 80) — REPLACE the residency framing. BEFORE: "the FO MAY route through a live standing prose-polisher (convention: `comm-officer`). Check team membership first." AFTER: the on-demand hook — "the FO MAY dispatch a one-shot polisher via `dispatch polish-spec --mod {path}`; it polishes once and dies. Best-effort, non-blocking; proceed un-polished if it does not return. Read usage from the mod." The Out-of-scope sentence (live replies, short statuses, commit messages) MOVES to the mod (it is usage prose).

### `first-officer-shared-core.md` — SHEDS

- **Line 27** (MODS-REPORT) — the parenthetical "a comm-officer spawn" as a reportable startup hook: REMOVE. The pr-merge startup-hook example stays.
- **Line 39** (greet deferrals) — REMOVE "the comm-officer spawn" from the list of expensive deferrals carried past the greet.
- **Line 108** — REMOVE "standing-teammate discovery/spawn" from the list of dispatch-reference machinery.
- **Line 201** — DELETE the sentence "The standing-teammate concepts (first-boot-wins lifecycle, team-scope teardown, the by-name routing contract, the declaration layout) travel with the deferred dispatch module…". The concepts no longer exist.

### `comm-officer.md` mod — CHANGES (concrete)

- Frontmatter: `standing: true` → DELETE the line.
- `## Hook: startup` and `## Hook: shutdown` — DELETE both sections.
- ADD `## One-Shot Prompt` — the resident `## Agent Prompt` body with residency framing removed: cut "Then idle. Do NOT start polishing anything until you receive a polish request." and the entire "You are a **standing teammate**: Stay live / Go idle between tasks / …" block; replace the opening with "You are a one-shot communications officer. A draft to polish follows. Polish it, reply once with the polished text + notes block, then stop." The first-action online-message handshake (the elements-of-style availability check + the two online messages) is DROPPED — a one-shot does not announce itself and idle; it polishes the draft it was handed. KEEP: the four modes, boundary rules, reply formats, light-touch defaults, qualifier-preservation rules.
- `## Agent Prompt` — DELETE. The mod is now purely an on-demand polish mod. (Residency back is Approach B, out of scope.)

### `using-claude-team` skill — no change

Grep confirms no standing-teammate prose lives in `skills/using-claude-team/` (the seed's worry was stale — the lifecycle notes live in `claude-fo-dispatch.md`, already covered above). Record: "no edit needed."

### `dispatch.go` / `standing.go` — CHANGES

- `dispatch.go`: add a `case "polish-spec"` routing to `runPolishSpec(modPath, stdout, stderr)` (requires `--mod` only). Mirror the existing `spawn-standing` flag-parse.
- `standing.go`: add `runPolishSpec` — parse the mod, require a `## One-Shot Prompt` section (loud non-zero if absent), emit a spec with `team_name` omitted. Share `ParseModMetadata`/`ParseHookStartupSpawnConfig`.

## Out of scope

- The behavioral polish QUALITY (the elements-of-style skill usage) — unchanged.
- T3's FO-contract prose audit (sibling; this is a specific contract-surface reduction, T3 is the general audit).
- Approach B (spawn-on-first-route, keep residency) — recorded as the fallback if a live-polish-latency concern surfaces, but A is the chosen direction.
- DELETING `spawn-standing`/`list-standing`/`show-standing` — they stay for any other resident-teammate workflow; this task adds the one-shot path and retargets `comm-officer` onto it.

## Sequencing

j9's P1 split has already merged (`claude-fo-dispatch.md` exists on `main`), so the ordering hazard the seed flagged is resolved; the before/after line anchors above were read from the merged file. No further gating.

## Scaffolding guardrail

Touches shipped scaffolding (`skills/first-officer/references/`, the `_mods/comm-officer.md` mod, the dispatch binary) — a dispatched worker under test, never an FO-direct edit.

## Acceptance criteria

Each AC names a proof OUTSIDE the contract prose (the file under change cannot be its own expectation — a contract grep is banned as proof). The text edits (the prose sheds, the mod rewrites) are real authoring work but are NOT acceptance criteria on their own; the behavioral half below is what each rests on.

1. **`dispatch polish-spec` emits a one-shot spec from the mod, with NO team_name and NO residency check.**
   Verified by: a Go test in `internal/dispatch` (sibling to `standing_parity_test.go`) that runs `polish-spec --mod {comm-officer fixture}` and asserts the emitted JSON has `subagent_type`/`name`/`model`/`prompt` present, `team_name` ABSENT, and the prompt equals the mod's `## One-Shot Prompt` body; plus a negative test that a mod lacking `## One-Shot Prompt` exits non-zero with the section-name diagnostic. The expected prompt comes from the fixture (an independent source), not from the binary. The offline gate (`go test ./...`) exits 0.

2. **No standing teammate enters team `config.json` for polish across a live session until a polish is needed; the one-shot polisher is present only transiently.**
   Verified by: a live drive (one Spacedock session) where the FO boots, dispatches at least one ensign stage, and at NO point spawns `comm-officer` at boot or first-dispatch. Durable proof: capture `~/.claude/teams/{team}/config.json` members BEFORE any polish — `comm-officer` ABSENT. Then the FO composes a deliberate draft (a gate summary), dispatches the one-shot polisher via `polish-spec`, and the polisher appears in the roster only while polishing and is gone (or never resident) after. The boot/first-dispatch member list having no `comm-officer` is the falsifiable signal — under the old lifecycle it WOULD be there at first dispatch.

3. **The one-shot polish is driven by reading usage from the MOD, not the contract, and returns a polished draft.**
   Verified by: in the same live drive, the FO/ensign reads the polish modes + boundary rules from `comm-officer.md` (the contract no longer carries them), dispatches a one-shot polisher with a real draft, and receives polished text + a notes block back. Proof is the returned polished artifact (an observable message/file the polisher produced), not a grep that the mod contains the modes.

4. **`list-standing` no longer returns `comm-officer`; existing scenarios stay green.**
   Verified by: `dispatch list-standing --workflow-dir docs/dev` produces empty stdout (the mod dropped `standing: true`) — an observable command output, not a prose check; and the existing live/fixture scenarios that exercised the dispatch path remain green.

### Hard comm-officer guard (folded in)

A one-shot polisher MUST NOT touch MUST / MUST-NOT clauses or semantic qualifiers in contract prose — the same boundary the resident mod already carries ("Do not change semantic qualifiers silently"). When `comm-officer.md` is rewritten for one-shot, this guard is PRESERVED verbatim in the `## One-Shot Prompt` body (it is a polish-quality boundary, independent of residency). Carrying it is part of AC-3's mod rewrite; the proof is that a polish drive over a draft containing a MUST clause returns it unchanged.

## Test plan

- **AC-1 (mechanism):** Go unit test, fixture-driven, ~30 min. The fixture is a minimal `_mods` polish mod with a `## One-Shot Prompt`. Cheapest gate; runs in CI via `go test ./...`. The spike above already proved the extraction half live; this test pins the one-shot-specific behavior (no team_name, section requirement).
- **AC-2 + AC-3 + guard (behavioral):** one live Spacedock smoke drive, ~1 session. The load-bearing proof that residency is gone and the mod is self-describing. Capture `config.json` members at boot/first-dispatch (proof of absence) and the returned polished draft (proof of the on-demand round-trip). Cannot be a fixture — roster membership over a real session is the claim.
- **AC-4:** `list-standing` command-output assertion (trivial) + re-run the existing scenario suite, green. ~10 min.
- Cost/complexity: medium. The binary change is small (one subcommand sharing proven helpers). The prose sheds are mechanical but span three files. The live drive is the only multi-minute item.

## Stage Report: ideation

- DONE: Concrete before/after wording: what the FO contract SHEDS vs the minimal generic hook it keeps, AND the mod-owned usage prose moved into _mods/comm-officer.md.
  `## Concrete before/after wording` section names each shed (discovery pass, lazy-spawn, declaration layout, team-scope teardown, first-boot-wins) with current line anchors across `claude-fo-dispatch.md` + `first-officer-shared-core.md`, the kept 3-sentence hook, and the `comm-officer.md` rewrites (drop `standing:`/hooks, add `## One-Shot Prompt`).
- DONE: Exercise the riskiest mechanism FIRST: one-shot spec-emit repurposing spawn-standing's prompt-extraction.
  Ran `./spacedock dispatch spawn-standing` against the comm-officer mod with a non-matching team name → exit 0, full Agent spec with the mod's prompt verbatim (`## Mechanism spike`). One-shot `polish-spec` = that path minus `MemberExists` short-circuit + `--team`; the spike also surfaced that the resident Agent Prompt needs a one-shot variant.
- DONE: Behavioral AC + concrete live-drive scenario proving NO standing teammate in config.json until polish is needed; NO contract prose-grep; fold the hard comm-officer guard (never touch MUST/MUST-NOT/qualifiers).
  AC-2 captures `config.json` members at boot/first-dispatch (comm-officer ABSENT is the falsifiable signal); every AC binds proof to an independent source (fixture, command output, returned artifact) with the contract-grep ban stated; `### Hard comm-officer guard` folded into the ACs.

### Summary

Fleshed the seed into a full ideation spec for replacing the standing-polisher lifecycle with an on-demand one-shot `dispatch polish-spec`. Proved the riskiest mechanism live (spawn-standing's prompt-extraction is reusable; one-shot is that path minus residency), which surfaced a real design decision: the resident Agent Prompt needs a one-shot `## One-Shot Prompt` variant. Wrote concrete before/after wording across three files plus the binary, and four behavioral ACs each bound to an out-of-contract proof (Go fixture test, live config.json roster capture, returned polished artifact, list-standing output) with the contract-grep ban honored and the MUST/qualifier guard folded in. j9 P1 has already merged, so the seed's sequencing gate is resolved.

## Stage Report: implementation

- DONE: Ship `dispatch polish-spec --mod {path}`: emits a one-shot Agent spec (subagent_type/name/model/prompt, team_name ABSENT), requires `## One-Shot Prompt`, exits non-zero with a section-name diagnostic when absent — plus the AC-1 Go fixture test (sibling to standing_parity_test.go) asserting the spec shape from a fixture mod, NOT the binary.
  `runPolishSpec` in standing.go + `case "polish-spec"` in dispatch.go + `ParseOneShotPrompt` in mods.go; `internal/dispatch/polish_spec_test.go` (TestPolishSpecEmitsOneShotSpec asserts team_name absent + prompt==fixture body; TestPolishSpecMissingOneShotPromptFailsLoud asserts the section-name diagnostic). Commit fefedfbb.
- DONE: Rewrite `_mods/comm-officer.md`: drop `standing: true` + `## Hook: startup` + `## Hook: shutdown` + `## Agent Prompt`; add `## One-Shot Prompt` (per-polish framing, no online-handshake/idle), carrying the four modes, boundary rules, reply formats, and the hard MUST/MUST-NOT/qualifier-preservation guard verbatim. Shed the standing-teammate lifecycle prose from claude-fo-dispatch.md + first-officer-shared-core.md per the ideation; keep only the minimal generic on-demand hook.
  Mod carries subagent_type/model in frontmatter (no hooks); `## One-Shot Prompt` carries the four modes + boundary + reply formats + the MUST/qualifier guard. claude-fo-dispatch.md lost the discovery-pass/lazy-spawn/declaration-mechanics/`## Standing Teammates` sections; `## Dispatch` keeps the generic on-demand hook. first-officer-shared-core.md lines 27/39/108/201 + claude-first-officer-runtime.md line 7 shed. Commit fefedfbb.
- PARTIAL: Behavioral proof (no contract prose-grep): a live drive shows NO `comm-officer` in team config.json at boot/first-dispatch and `dispatch list-standing` returns empty; offline gate `go test ./...` exits 0.
  Binary-level halves proven live: `./spacedock dispatch list-standing --workflow-dir docs/dev` → empty stdout exit 0; `./spacedock dispatch polish-spec --mod .../comm-officer.md` → valid spec, no team_name, model sonnet, `## One-Shot Prompt` body with the MUST guard. `go test ./...` exits 0. The full nested FO smoke session (boot → dispatch an ensign → capture config.json roster) was NOT run from inside this dispatched ensign — see Summary; it is structurally guaranteed but unexercised end-to-end here.

### Summary

Added `dispatch polish-spec` (one-shot Agent spec, no team_name) repurposing the proven spawn-standing mod-parse, with a fixture Go test; rewrote comm-officer.md off residency (frontmatter spawn config + `## One-Shot Prompt` carrying the modes/boundary/guard); and shed the standing-teammate lifecycle prose from the three FO references, keeping only a generic on-demand routing hook. spawn/list/show-standing stay for other resident-teammate workflows. The AC-2/AC-3 live config.json roster capture was NOT run as a nested FO session from inside this ensign — spinning up a full nested Spacedock drive is heavy and risky here. It is structurally guaranteed: the boot discovery-pass + lazy-spawn code paths are deleted, comm-officer dropped `standing: true` (list-standing empty, proven), and polish-spec emits a bare Agent (no team_name → joins no roster), so no surviving path puts comm-officer in config.json at boot/first-dispatch. Flagging the integration smoke drive to the FO/captain as the remaining gate.

## Stage Report: validation

- DONE: Reproduce the structural shed from the worktree (not the report).
  `claude-fo-dispatch.md` section grep: discovery-pass / lazy-spawn / declaration-mechanics / `## Standing Teammates` ALL absent (`grep -ni` → NONE); `## Dispatch` carries only the generic on-demand polisher hook. `first-officer-shared-core.md` + runtime: zero `comm-officer`/`standing`/`spawn-standing` matches (lines 27/39/108/201 shed). `_mods/comm-officer.md`: no `standing:`, no `## Hook: startup`/`shutdown`, no `## Agent Prompt`; carries `## One-Shot Prompt` with the MUST/MUST-NOT/qualifier-preservation guard verbatim (sed-extracted, matched).
- DONE: Adversarially confirm NO surviving boot path can spawn a standing comm-officer.
  Repo-wide grep (`*.md`/`*.go`, excl. state/roadmap/reviews) for discovery-pass/lazy-spawn/first-boot-wins/spawn-on-first → NONE. `status --boot --json` `mods.startup` = `["pr-merge"]` only; no standing-teammate spawn directive. `dispatch build`'s standing-fetch append is gated on `EnumerateDeclaredStandingTeammates > 0`, which is empty here (list-standing AND show-standing both empty stdout). Only remaining route is `polish-spec` → bare Agent (no team_name → joins no roster).
- DONE: Offline + binary proof (worktree-built binary; PATH spacedock is stale).
  Built `/tmp/spacedock-wt` from the worktree. `go test ./...` exit 0 (incl. fresh `TestPolishSpecEmitsOneShotSpec` [team_name ABSENT + prompt==fixture body] and `TestPolishSpecMissingOneShotPromptFailsLoud` [section-name diagnostic], both PASS verbose). `dispatch list-standing --workflow-dir docs/dev` → empty stdout exit 0. `dispatch polish-spec --mod docs/dev/_mods/comm-officer.md` → bare Agent spec: subagent_type general-purpose, name comm-officer, model sonnet, NO team_name, prompt = the One-Shot Prompt body carrying the MUST/MUST-NOT guard.
- DONE: AC-2/AC-3 live drive — feasibility assessed; structural verdict + residual flagged (no fabricated live result).
  A clean fresh-boot live capture is infeasible from inside this dispatched ensign: spinning a nested FO→ensign Spacedock session inside an already-active team is heavy and recursive. Notable observed datum: the CURRENT live team's `~/.claude/teams/.../config.json` DOES contain a resident `comm-officer` carrying the OLD residency prompt ("Then idle... You are a **standing teammate**... Stay live") — this team booted under the pre-change lifecycle, so it is exactly AC-2's "under the old lifecycle it WOULD be there" baseline, NOT a regression in the deliverable. The deliverable changes the source-of-truth (mod + binary + contract) so a fresh boot has no path to spawn it. Verdict: structurally guaranteed — discovery/lazy-spawn deleted, `standing:` dropped (enumeration empty, proven), polish-spec bare. Residual: an end-to-end fresh-boot roster capture is owed before merge; flagged to FO/captain.

### Summary

Reproduced every structural shed independently from the worktree (not the report) and adversarially confirmed no surviving boot path spawns a standing comm-officer: discovery/lazy-spawn/declaration symbols are gone repo-wide, boot JSON carries no standing spawn directive, and the only remaining route (`polish-spec`) emits a bare team_name-less Agent. Offline + binary proofs all green on a worktree-built binary (go test ./... exit 0 incl. both polish_spec tests; list-standing empty; polish-spec emits the bare One-Shot spec with the MUST guard). AC-1 and AC-4 are fully proven. AC-2/AC-3 rest on a fresh-boot live roster capture that is infeasible from inside this active team; the chain is structurally closed (current team's resident comm-officer is the pre-change baseline, not a deliverable regression), but the end-to-end fresh-boot capture is the one residual owed before merge. Recommendation: PASSED on the binary/structural half; the live fresh-boot roster smoke is flagged to the FO/captain as the remaining integration gate (matching the implementation's own flag).
