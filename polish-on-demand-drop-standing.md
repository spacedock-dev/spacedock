---
id: xf7fft1hnj51eq7kagsc9833
title: Move the standing-teammate mechanism out of the FO contract — the comm-officer mod self-injects as a standing teammate
status: implementation
source: "captain (2026-06-13, this session) — the standing-teammate lifecycle (discovery pass / lazy-spawn / declaration / team-scope teardown / first-boot-wins) is ~4 contract subsections of maintenance surface for an infrequently-used polisher; amortization doesn't pay for infrequent use. Captain chose approach A: on-demand one-shot polish dispatch, and 'the mod can add the prose for the standing team member' — feature-specific usage prose lives in the mod, the contract keeps only a generic hook. Taken into 0203."
started: 2026-06-14T18:42:01Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-polish-on-demand-drop-standing
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

---

# CYCLE 3 — DIRECTION REVERSED (supersedes cycles 1-2)

> **Captain reversal (2026-06-14):** cycles 1-2 above (drop residency → on-demand one-shot `polish-spec`) are REJECTED. A one-shot polisher is pointless — a slash-command skill already covers ad-hoc polish. **KEEP the comm-officer RESIDENT.** The win is the SAME contract-bloat reduction, achieved by RELOCATING the standing-teammate injection mechanism out of the FO contract prose into the binary + the mod's self-declaration, NOT by dropping the feature. Everything below replaces the cycle-1/2 Proposed approach, before/after wording, ACs, and test plan. The cycle-1/2 text is retained above for the record, not as the plan.

## Corrected goal

Move the standing-teammate INJECTION MECHANISM out of the FO contract prose. Today the FO contract spells out HOW standing teammates are discovered, lazy-spawned, deduped, and torn down — ~4 subsections of machinery (`claude-fo-dispatch.md` lines 17-24 discovery-pass, 26-38 lazy-spawn, 40-47 declaration/routing mechanics, 49-56 `## Standing Teammates` concepts). The captain wants that HOW relocated into the binary (which already has the spawn helpers) plus the mod's self-declaration, so the contract carries at most a minimal generic trigger line. **Residency is preserved** — comm-officer keeps `standing: true` + spawn config + agent prompt and stays a long-lived team member.

This is "prefer a code gate over prose" applied to a lifecycle: the loop the FO currently executes by reading prose (enumerate declared standing mods → per-mod dedup against the live team config → emit the absent ones' Agent specs) becomes ONE binary call. The mechanism lives in the binary, not in contract sentences.

## Mechanism spike (cycle 3, riskiest unknown — exercised FIRST)

**The claim:** the binary can drive the FULL standing lifecycle generically (enumerate → first-boot-wins dedup → emit spawn specs) behind one helper call, so the contract's discovery-pass + lazy-spawn + declaration-mechanics + first-boot-wins prose can be DELETED.

**What I verified against `main` (HEAD 3779370f) + the live binary:**

1. **Enumeration already generic.** `EnumerateDeclaredStandingTeammates(workflowDir, teamName)` (`internal/dispatch/mods.go:95`) already scans `{wd}/_mods/*.md`, filters `standing: true`, resolves name from `## Hook: startup` (falling back to frontmatter), sorted. This IS the discovery pass — already in the binary. `./spacedock dispatch list-standing --workflow-dir docs/dev` → emits the comm-officer mod path, exit 0 (ran live).
2. **Dedup already generic.** `MemberExists(home, team, name)` (`internal/claudeteam/standing.go:71`) reads `~/.claude/teams/{team}/config.json` members — this IS the first-boot-wins predicate. `runSpawnStanding` already calls it and short-circuits to `{"status":"already-alive","name":...}` (`standing.go:144-154`).
3. **Spec emit already proven.** `./spacedock dispatch spawn-standing --mod {comm-officer} --team {fresh-team-with-no-config}` → exit 0, full Agent spec with `subagent_type`, `name`, `team_name`, `model: sonnet`, and `prompt` = the `## Agent Prompt` body verbatim (ran live this cycle; output captured). So per-mod spec emit works end-to-end today.

**Conclusion — the binary already has all three pieces; only the DRIVER that composes them is missing.** Today the FO composes them in prose (loop over `list-standing` output, call `spawn-standing` per path, branch on `already-alive`). The MINIMAL new capability is a thin subcommand — `dispatch spawn-standing-all --workflow-dir {wd} --team {team}` — that runs the loop in Go and emits a JSON ARRAY of the Agent specs for the not-already-alive declared standing mods (already-alive ones omitted). No new lifecycle logic; it calls `EnumerateDeclaredStandingTeammates`, then for each runs the existing `runSpawnStanding` body, collecting the non-`already-alive` specs. **Precedent already in tree:** `dispatch build` (`build.go:570`) ALREADY auto-injects the `show-standing` fetch line gated on `EnumerateDeclaredStandingTeammates(...) > 0` — the binary already drives one half of the standing machinery generically without FO prose. `spawn-standing-all` extends the same pattern to the spawn half.

This is throwaway-spike-confirmed: the riskiest claim ("the binary can carry the full lifecycle") is true with one small additive subcommand, not a rewrite. No on-disk format change, no parser round-trip risk — it reuses proven parse + member-probe code.

## Proposed approach (cycle 3)

**Approach C — relocate the mechanism, keep residency.**

### Binary: add `dispatch spawn-standing-all`

A new subcommand `spacedock dispatch spawn-standing-all --workflow-dir {wd} --team {team}` that drives the full inject loop in one call:

- Run `EnumerateDeclaredStandingTeammates(wd, team)`. For empty (bare mode / no `_mods` / no standing mods) → emit `[]` (empty JSON array), exit 0.
- For each declared standing mod, in sorted order, run the existing `runSpawnStanding` logic (parse, validate, `MemberExists` dedup). Collect the emitted Agent spec for each member that is NOT already-alive; SKIP already-alive members (they are the first-boot-wins case — no spec needed).
- Emit a JSON ARRAY of the collected specs (`[]spawnSpec`) to stdout, two-space-indented, matching the existing `emitSpawnJSON` byte conventions. Exit 0.
- Loud failure (exit 1, stderr naming the offending mod) on the same conditions `runSpawnStanding` already fails on (missing `## Agent Prompt`, missing/invalid `subagent_type`/`name`/`model`) — the validation is shared, not re-implemented.

The FO calls this ONCE before the first team-mode `Agent()` dispatch, forwards each spec in the returned array to `Agent()`, and is done. The enumerate/dedup/per-mod-loop that the contract currently spells out is now entirely inside the binary. `spawn-standing` / `list-standing` / `show-standing` STAY (other callers and the per-dispatch `show-standing` injection still use them); this is purely additive.

### Mod stays self-describing AND keeps residency (the captain's key point)

`comm-officer.md` KEEPS `standing: true`, its `## Hook: startup` spawn config (`subagent_type`/`name`/`model`), and its `## Agent Prompt` — residency fully preserved; the binary reads all of these. What MOVES into the mod is the *mechanism narrative* the contract sheds: the mod's `## Hook: startup` currently re-states "check the team config / first-boot-wins / spawn if absent" in PROSE (lines 12-27) — that narrative is now redundant because the binary's `spawn-standing-all` performs the check. The mod's `## Hook: startup` is trimmed to the spawn-config bullets the binary parses (the `- subagent_type:` / `- name:` / `- model:` block), dropping the human-readable check-the-config narrative. The mod's `## Routing guidance` / `## Routing Usage` (the when-to-polish scope, the four modes, boundary rules) STAY — they were always mod-owned and become the single home for usage prose the contract also sheds.

### Contract keeps only a minimal generic trigger line

`claude-fo-dispatch.md` sheds all four subsections (below) and keeps ONE line in `## Dispatch`: *"Before the first team-mode `Agent()` dispatch, inject declared standing teammates: run `spacedock dispatch spawn-standing-all --workflow-dir {wd} --team {team_name}` and forward each Agent spec in the returned JSON array to `Agent()`. Idempotent (already-alive members are omitted), so re-running is safe. Read each teammate's routing usage from its mod."* The HOW (enumerate, dedup, first-boot-wins, declaration layout) is gone from prose — it lives in `spawn-standing-all` + the mod.

## Concrete before/after wording (cycle 3)

### `claude-fo-dispatch.md` — SHEDS (re-confirmed against `main` HEAD 3779370f)

- **`### Standing teammate discovery pass`** (lines 17-24) — DELETE entirely. `spawn-standing-all` enumerates internally.
- **`### Standing teammate lazy-spawn`** (lines 26-38) — DELETE entirely. The per-mod loop + already-alive branch move into the binary. KEEP the one operational caveat that is NOT mechanism — the "polish round-trips can reach several minutes; treat routing as non-blocking" sentence — relocated to the routing hook at line 80 (it is a routing-discipline rule, not lifecycle machinery).
- **`### Standing teammate declaration and routing mechanics`** (lines 40-47) — DELETE the declaration-layout / routing-call / teardown-trigger / lazy-spawn-injection bullets. The declaration layout is now an implementation detail of `spawn-standing-all` (it parses the mod); the teardown trigger ("dies with the team at Claude Code teardown") is the only survivor and folds into the `## Standing Teammates` replacement as a single clause. The dispatch-time `show-standing` injection bullet STAYS true but is already enforced by `build.go` — no contract prose needed; DELETE the bullet (it described binary behavior the binary already guarantees).
- **`## Standing Teammates`** (lines 49-56) — REPLACE the four concept bullets (first-boot-wins, team-scope lifecycle, routing contract, declaration) with the single generic trigger line above (placed in `## Dispatch`). The section heading is removed; its one surviving fact — "standing teammates are team-scoped and die with the team at teardown" — becomes a clause in the trigger line.
- **`## Dispatch` → "Routing through a standing prose-polisher"** (line 80) — KEEP the residency framing (comm-officer IS live), but DROP the "Check team membership first" mechanism hint (membership is the binary's concern now) and the "Dispatched workers discover the same teammates through their build-time prompt" sentence (that is the `show-standing` build injection, now undocumented-because-automatic). Fold in the non-blocking-several-minutes caveat moved from line 38.
- **Line 129 sequencing rule** — KEEP. The "spawn-standing emits Agent specs; never in the same message as TeamCreate/TeamDelete" race rule still applies to `spawn-standing-all` (it also emits Agent specs forwarded to Agent dispatch). Update the wording to name `spawn-standing-all` as the spec-emitter the FO actually calls.

### `first-officer-shared-core.md` — SHEDS (re-confirmed)

- **Line 27** (MODS-REPORT) — KEEP the "comm-officer spawn" parenthetical as a reportable deferred startup hook (residency is preserved, so the greet still reports it as a pending deferral). No change — this was a cycle-1/2 shed that no longer applies.
- **Line 39** (greet deferrals) — KEEP "the comm-officer spawn" in the deferral list (still a real deferred spawn). No change.
- **Line 108** — KEEP "standing-teammate discovery/spawn" in the dispatch-machinery list, but the machinery now means "the `spawn-standing-all` call," not a prose procedure. Reword: "standing-teammate injection (`spawn-standing-all`)."
- **Line 201** — REPLACE "The standing-teammate concepts (first-boot-wins lifecycle, team-scope teardown, the by-name routing contract, the declaration layout) travel with the deferred dispatch module…" with a single sentence: "Standing-teammate injection is driven by `spacedock dispatch spawn-standing-all` at first dispatch; the concept is team-scoped (members die with the team)." The lifecycle CONCEPTS no longer need a prose home — the binary carries them.

### `claude-first-officer-runtime.md` — line 7

- REPLACE "standing-teammate discovery/lazy-spawn/declaration" in the dispatch-machinery list with "standing-teammate injection (`spawn-standing-all`)." One phrase, not a procedure.

### `comm-officer.md` mod — CHANGES (residency PRESERVED)

- Frontmatter `standing: true` — KEEP.
- `## Hook: startup` — TRIM to the spawn-config bullets the binary parses (`- subagent_type: general-purpose` / `- name: comm-officer` / `- model: sonnet`). DELETE the human-readable "check the team config members list / first-boot-wins / if present skip / if absent spawn" narrative (lines 14-27) — that mechanism now lives in `spawn-standing-all`. The mod declares WHAT to spawn; the binary decides WHEN/IF.
- `## Hook: shutdown` — KEEP (teardown is still real; the binary does not own teardown, Claude Code does).
- `## Agent Prompt` — KEEP verbatim (residency). The resident framing ("Then idle… You are a standing teammate… Stay live") is correct and stays.
- `## Routing guidance` / `## Routing Usage` — KEEP (mod-owned usage prose, now the single home).

### `using-claude-team` skill — no change

Grep confirms no standing-teammate mechanism prose lives in `skills/using-claude-team/` (re-verified against `main`). Record: "no edit needed."

## Acceptance criteria (cycle 3)

Each AC names a proof OUTSIDE the contract prose (a contract grep is BANNED as proof — the file under change cannot be its own oracle). The prose sheds are real authoring work but rest on the behavioral half below.

1. **`dispatch spawn-standing-all` drives the full inject loop: it emits the Agent specs for declared standing mods absent from the team config, omits already-alive members, and emits `[]` when none are declared.**
   Verified by: a Go test in `internal/dispatch` (sibling to `standing_parity_test.go`) that, against a fixture `_mods` dir with a `standing: true` polish mod and a fixture team config, asserts: (a) with the member ABSENT from config, `spawn-standing-all` emits a one-element JSON array whose spec has `subagent_type`/`name`/`team_name`/`model`/`prompt` and prompt == the fixture mod's `## Agent Prompt` body (expected value sourced from the fixture, NOT the binary); (b) with the member PRESENT in the fixture config, the array OMITS it (already-alive dedup); (c) with no `standing: true` mod, output is `[]`. The offline gate `go test ./...` exits 0.

2. **A live first-dispatch injects comm-officer into the team `config.json` roster as a standing teammate — driven by the mod's self-declaration + the binary — WITHOUT the FO contract carrying the discovery/lazy-spawn/first-boot-wins mechanism prose.**
   Verified by a REGISTERED live scenario (per the 58/lean-boot lesson — an ad-hoc live test proves nothing). The scenario boots an FO, reaches the first team-mode dispatch, and the FO calls `spawn-standing-all` (the ONLY standing-injection step it has). Durable proof: capture `~/.claude/teams/{team}/config.json` members AFTER first dispatch — `comm-officer` is PRESENT (residency works). The independent oracle is the roster file, not a contract grep. This AC REDS if residency breaks (comm-officer absent after first dispatch) — the inverse of the cycle-1/2 AC, because residency is now the desired state.

3. **The mechanism is shed from the FO contract — discovery-pass / lazy-spawn / declaration-mechanics / first-boot-wins prose is removed on-disk — while injection still works.**
   Verified by the SAME registered live scenario as AC-2: the contract files (`claude-fo-dispatch.md` + shared-core + runtime) no longer contain the four mechanism subsections (the shed is an on-disk fact the scenario's setup confirms by structural absence of the named sections), AND the live roster from AC-2 still shows comm-officer resident. The oracle is the live roster (independent), guarded against the contract regressing the mechanism back in: a structural check that the named section headings are ABSENT is the shed-half, but the LOAD-BEARING pass condition is AC-2's roster (injection works through the binary, not prose). Note: per the proof policy, the structural-absence half alone is not sufficient — it is bound to AC-2's live roster so the AC cannot pass on wording alone.

4. **`list-standing` STILL returns comm-officer (residency preserved); existing scenarios stay green.**
   Verified by: `dispatch list-standing --workflow-dir docs/dev` STILL emits the comm-officer mod path (the mod KEEPS `standing: true`) — an observable command output proving residency survived the refactor; and the existing live/fixture standing scenarios remain green. This is the falsifiable guard that the cycle-3 refactor did NOT accidentally drop residency (the cycle-1/2 mistake).

### Hard comm-officer guard (carried forward)

The resident polisher MUST NOT touch MUST / MUST-NOT clauses or semantic qualifiers in contract prose — the boundary the mod already carries ("Do not change semantic qualifiers silently"). Residency is preserved, so the `## Agent Prompt` keeps this guard verbatim with NO change. Proof: a polish drive over a draft containing a MUST clause returns it unchanged (part of AC-2's live scenario where the resident polisher is exercised).

## Test plan (cycle 3)

- **AC-1 (mechanism):** Go unit test, fixture-driven, ~40 min. Fixtures: a minimal `_mods` standing polish mod + two team-config fixtures (member-absent, member-present). Cheapest gate; runs in CI via `go test ./...`. The spike proved the three component functions live; this test pins the new driver's array-emit + dedup + empty-`[]` behavior.
- **AC-2 + AC-3 + guard (behavioral):** ONE registered live Spacedock scenario, ~1 session. Load-bearing proof that residency works through the binary and the mechanism is shed. Capture `config.json` members after first dispatch (comm-officer PRESENT = residency proof) and the structural absence of the four contract subsections. MUST be a registered scenario, not ad-hoc (lean-boot lesson). Cannot be a fixture — roster membership over a real boot→dispatch is the claim.
- **AC-4:** `list-standing` command-output assertion (trivial, comm-officer path present) + re-run the existing standing scenario suite, green. ~10 min.
- Cost/complexity: medium. The binary change is small and ADDITIVE (one subcommand composing three proven functions — no rewrite of `spawn-standing`/`list`/`show`). The prose sheds are mechanical across four files. The registered live scenario is the only multi-minute item and is the load-bearing proof.

## Doc-diff note

No user-facing CLI output, banner, or docs-site surface changes for the human captain: `spawn-standing-all` is an FO-internal dispatch subcommand (like `build`/`reconcile`), not a captain-typed command. The shed prose is internal contract/reference text. No `docs/` site diff needed; recorded here per the ideation doc-diff rule.

## Stage Report: ideation (cycle 3)

- DONE: Design HOW the comm-officer mod SELF-INJECTS as a standing teammate via the binary's spawn helpers — KEEP residency, drop the on-demand polish-spec direction. Exercise the riskiest mechanism FIRST (can the binary drive the full lifecycle behind ONE call); record what the binary already supports vs the minimal new capability.
  `## Mechanism spike (cycle 3)`: verified live against `main` HEAD 3779370f that `EnumerateDeclaredStandingTeammates` (discovery), `MemberExists` (first-boot-wins dedup), and `runSpawnStanding` (spec emit) ALL exist; only a thin composing driver is missing. Designed `dispatch spawn-standing-all` (one call → JSON array of absent-member specs); precedent is `build.go:570` already auto-injecting the show-standing line generically. Residency preserved (mod keeps `standing: true` + spawn config + `## Agent Prompt`).
- DONE: Concrete before/after wording — which mechanism sections move OUT of claude-fo-dispatch.md into the binary + mod self-declaration, and the minimal surviving contract footprint (~one line).
  `## Concrete before/after wording (cycle 3)`: claude-fo-dispatch.md sheds all four subsections (discovery-pass 17-24, lazy-spawn 26-38, declaration-mechanics 40-47, `## Standing Teammates` 49-56), keeping ONE trigger line in `## Dispatch`; shared-core line 201 + runtime line 7 reduced to a single `spawn-standing-all` phrase; mod's `## Hook: startup` trimmed to spawn-config bullets (mechanism narrative dropped) while `standing: true`/`## Agent Prompt`/`## Routing guidance` stay. Line anchors re-confirmed against origin/main.
- DONE: Behavioral AC, NO contract prose-grep — a REGISTERED live scenario or binary-level test proving comm-officer IS injected into config.json as a standing teammate via mod self-declaration + binary, WHILE the FO contract no longer carries the mechanism prose. Bind to an independent oracle that reds if residency breaks OR the mechanism leaks back.
  Four ACs: AC-1 binary-level Go fixture test (array-emit + dedup + empty `[]`, expected prompt from fixture); AC-2 registered live scenario (config.json roster shows comm-officer PRESENT after first dispatch — reds if residency breaks); AC-3 shed-half bound to AC-2's live roster (structural-absence alone insufficient per proof policy); AC-4 `list-standing` STILL returns comm-officer (residency-preservation guard against the cycle-1/2 mistake). Hard MUST/qualifier guard carried forward verbatim.

### Summary

Re-ideated cycle 3 from the captain's reversed goal: KEEP comm-officer resident, RELOCATE the standing-teammate injection mechanism out of the FO contract prose into the binary + the mod's self-declaration. The riskiest claim — the binary can drive the full lifecycle behind one call — was spike-confirmed live: enumeration, first-boot-wins dedup, and spec-emit all already exist in the binary (`EnumerateDeclaredStandingTeammates` + `MemberExists` + `runSpawnStanding`); only a thin composing driver `dispatch spawn-standing-all` is missing, extending the same generic-injection pattern `build.go` already uses for the show-standing line. The FO contract sheds all four mechanism subsections down to one trigger line; the mod keeps `standing: true` + spawn config + agent prompt (residency intact), shedding only its redundant first-boot-wins narrative. Four ACs bind to independent oracles (Go fixture array-emit test, registered live config.json roster capture showing comm-officer PRESENT, shed-bound-to-roster, `list-standing` residency guard) with the contract-grep ban honored and residency-breakage as the red condition. Line anchors re-confirmed against origin/main HEAD 3779370f; cycles 1-2 retained above for the record.

## Stage Report: implementation (cycle 3)

- DONE: Add the binary driver `dispatch spawn-standing-all` (one call → JSON array of the ABSENT declared-standing Agent specs) composing EnumerateDeclaredStandingTeammates + MemberExists + runSpawnStanding, mirroring build.go; plus its Go fixture test (array-emit shape, dedup skips already-present member, empty `[]` when none/all-present — expected prompt from the fixture).
  `runSpawnStandingAll` + shared `buildSpawnSpec` (extracted from runSpawnStanding so single/all share validation+dedup) in standing.go; `case "spawn-standing-all"` in dispatch.go. `internal/dispatch/spawn_standing_all_test.go`: 4 tests (absent-member array spec w/ prompt==fixture body, already-alive dedup→`[]`, no-standing→`[]`, broken-mod→exit 1 naming the mod). Commit 776a5f82.
- DONE: Shed the mechanism from the FO contract into binary + mod self-declaration — claude-fo-dispatch.md's four subsections down to ONE trigger line in `## Dispatch`; shared-core + claude-first-officer-runtime.md reduced to a single `spawn-standing-all` phrase; mod's `## Hook: startup` trimmed to spawn-config bullets while KEEPING `standing: true` + `## Agent Prompt` + routing guidance. Anchors re-confirmed against CURRENT worktree base (ad6b149e, post-#373) — NOT the ideation's pre-5e/3779370f anchors.
  claude-fo-dispatch.md lost discovery-pass/lazy-spawn/declaration-mechanics/`## Standing Teammates`; `## Dispatch` carries the generic `spawn-standing-all` trigger line + trimmed routing-polisher para; sequencing rule renamed to spawn-standing-all. shared-core lines 99/192 + runtime line 7 reduced to the single phrase; lines 18/30 (comm-officer spawn deferral) KEPT (residency). Mod `## Hook: startup` = 3 spawn-config bullets only. Commit 776a5f82.
- DONE: Residency + guard intact: `dispatch list-standing` STILL returns comm-officer (AC-4), the hard MUST/MUST-NOT/qualifier-preservation guard stays verbatim in the mod, `go test ./...` green. AC-2 registered live scenario noted as the validation-stage behavioral AC (gates at PR CI).
  Worktree binary (`/tmp/spacedock-wt`): `list-standing --workflow-dir docs/dev` → comm-officer path, exit 0; `spawn-standing-all --workflow-dir docs/dev --team {fresh}` → one-element array, full `## Agent Prompt` body incl. "Do not change semantic qualifiers silently". Guard verbatim at mod lines 48/49/135 + Agent Prompt. `go test ./...` exit 0, 15 pkgs ok (cli/integration disk-pressure failures cleared after cache clean; re-ran green).

### Summary

Cycle-3 (residency-preserving) implementation: added `dispatch spawn-standing-all` driving the full inject loop in one Go call (enumerate → MemberExists dedup → array of absent-member Agent specs, `[]` when none), extracting a shared `buildSpawnSpec` from `runSpawnStanding` so single- and all-spawn share validation. Shed the four standing-teammate mechanism subsections from claude-fo-dispatch.md to one `## Dispatch` trigger line, reduced shared-core/runtime to a single `spawn-standing-all` phrase, and trimmed the mod's `## Hook: startup` to spawn-config bullets — while KEEPING `standing: true` + `## Agent Prompt` + the MUST/qualifier guard (residency intact, list-standing still returns comm-officer). Re-confirmed every anchor against the CURRENT worktree base (ad6b149e, post-#373), not the ideation's stale pre-5e anchors. AC-1 (Go fixture array-emit/dedup/empty test) and AC-4 (list-standing residency) proven; `go test ./...` exit 0. AC-2 (live config.json roster showing comm-officer PRESENT after first dispatch) is authored as the validation-stage behavioral AC and gates at PR CI — not run as a nested FO session from inside this dispatched ensign.

## Stage Report: validation (cycle 3)

- DONE: Reproduce the binary driver from a worktree-built binary.
  Built `/tmp/spacedock-wt-val` from `./cmd/spacedock` (the PATH `go build .` produced a library archive — there is no root main package; the CLI is `cmd/spacedock`). `dispatch spawn-standing-all --workflow-dir docs/dev --team {fresh}` → one-element array, keys subagent_type/name/team_name/model/prompt ALL present (team_name PRESENT = residency, the cycle-3 inverse of cycles 1-2), model sonnet, prompt = the resident `## Agent Prompt` (markers "Then idle"/"standing teammate"/"Stay live" present) carrying "Do not change semantic qualifiers silently" verbatim. Live dedup: a throwaway HOME with comm-officer in `config.json` → `[]`. `[]`/exit-1 paths covered by AC-1 fixtures below.
- DONE: AC-1's 4 fixture tests green (expected prompt from fixture, not binary).
  `go test ./internal/dispatch -run SpawnStandingAll -count=1 -v` → 4/4 PASS. Read `spawn_standing_all_test.go`: prompt oracle is the synthetic `standingMod()` fixture body `"You are comm-officer.\n"`, not the real mod or binary output; `runNative` drives the real `Run(...)` entrypoint. Genuine behavioral test, not self-referential.
- DONE: AC-4 RESIDENCY GUARD — `dispatch list-standing --workflow-dir docs/dev` STILL returns comm-officer.
  Binary output: the comm-officer mod path, exit 0. `show-standing` still renders the routing block. Residency preserved.
- DONE: Structural shed reproduced from the worktree (not the report).
  `claude-fo-dispatch.md`: discovery-pass / lazy-spawn / declaration-mechanics / `## Standing Teammates` ALL absent (grep → NONE), replaced by ONE `spawn-standing-all` trigger line (L21) + the renamed sequencing rule (L90). shared-core (L99/L192) + claude-first-officer-runtime.md (L7) reduced to a single `spawn-standing-all` phrase. `_mods/comm-officer.md`: `standing: true` KEPT, `## Agent Prompt` KEPT, `## Hook: startup` trimmed to the 3 spawn-config bullets (first-boot-wins narrative removed — confirmed in the commit diff), MUST/MUST-NOT + qualifier-preservation guard verbatim (mod L48/L49/L135).
- DONE: Word-level contract-shed audit — NO non-standing obligation or qualifier dropped.
  Every removed `claude-fo-dispatch.md` line is standing-MECHANISM prose. The "Out of scope" usage prose (live replies, short statuses, commit messages) was RELOCATED into the mod (L32-33), not dropped; the non-blocking-several-minutes caveat moved into the routing-polisher para; the verbatim-forward discipline and the MUST-NOT-precede-TeamCreate sequencing obligation are preserved.
- DONE: Driver composition audit — correct, no behavior change vs prior lazy-spawn path.
  `buildSpawnSpec` was extracted from `runSpawnStanding`; BOTH single- and all-spawn call it, so validation + `MemberExists` first-boot-wins dedup are shared (no drift). `runSpawnStandingAll` enumerates via the same `sortedModPaths`+`meta.Standing` scan `EnumerateDeclaredStandingTeammates` uses, dedups per-mod, emits `[]spawnSpec`. (Deliberate, defensible divergence: a broken mod is loud exit-1 in the all-spawn path vs skipped in the enumerator — the all-spawn path is the one that actually spawns, so failing loud is safer and matches the checklist.)
- FAILED: Confirm the AC-2 live-roster scenario (comm-officer PRESENT in config.json after first dispatch) is REGISTERED in the CI live gate.
  It is NOT registered. The only `comm-officer` token in `internal/ensigncycle/` live runners is a path COMMENT (`claude_live_runner_test.go:284`) inside the shallow-boot scenario, whose own "AC-2" is the lean-boot no-TeamCreate-before-greet assertion — the OPPOSITE (it asserts NO team at boot). The 6 registered shared scenarios (gate-guardrail, rejection-flow, feedback-3-cycle-escalation, merge-hook-guardrail, filing, shallow-boot) and the `runtime-live-e2e.yml` `-run` selectors (TestLiveEnsignCycle, TestLiveClaudeSharedScenarios, TestLiveCodexSharedScenarios, the Pi pair) include NO scenario that boots an FO, reaches first dispatch, and asserts comm-officer is PRESENT in `~/.claude/teams/{team}/config.json`. No registered live gate REDs if residency breaks.
- DONE: Offline gate.
  `go test ./...` exit 0 (15 pkgs ok), incl. the 4 fresh spawn-standing-all tests.

### Feedback Cycles

- **validation cycle 3 — detached adversarial audit (separate Explore checkout + independent grep/binary reproduction):** The audit confirmed the binary driver, the structural shed, the word-level no-dropped-obligation check, and AC-1/AC-4 all PASS. The material catch: **AC-2 — the load-bearing residency proof — is authored-as-intent but UNREGISTERED in the CI live gate.** The implementation report (cycle 3) claims AC-2 "gates at PR CI"; it does not. This is exactly the lean-boot/58 lesson the AC text names ("an ad-hoc live test proves nothing") and the proof-policy bar: residency over a real boot→dispatch is the falsifiable claim, and nothing in `runtime-live-e2e.yml` exercises it. A unit fixture (`spawn_standing_all_test.go`) proves the binary's array/dedup behavior in isolation, NOT that a live FO boot + first dispatch lands comm-officer in `config.json`. Required to close: register a live scenario (in the shared-runtime set or a dedicated live test) that boots an FO, reaches the first team-mode dispatch, calls `spawn-standing-all`, and asserts comm-officer PRESENT in the team `config.json` roster — wired into a `runtime-live-e2e.yml` `-run` selector so it gates at PR CI. Until then the residency claim rests only on offline/structural proof.

### Summary

REJECTED. The binary half is solid and independently reproduced from a worktree-built binary: `spawn-standing-all` correctly composes the enumeration + `MemberExists` first-boot-wins dedup + the shared `buildSpawnSpec` (extracted from `runSpawnStanding`, so no behavior drift), emitting a one-element residency-preserving array (team_name PRESENT) for an absent member and `[]` on live dedup. AC-1 (4 genuine fixture tests, prompt oracle from the fixture not the binary) and AC-4 (list-standing still returns comm-officer) PASS. The structural shed is clean and the word-level audit found NO dropped non-standing obligation — the out-of-scope usage prose was relocated into the mod, not lost; the mod keeps `standing: true` + `## Agent Prompt` + the MUST/qualifier guard verbatim. The blocking defect is AC-2: it is the explicitly load-bearing residency proof, the checklist requires it REGISTERED in the CI live gate, and it is NOT — no scenario in `runtime-live-e2e.yml` boots→dispatches→asserts comm-officer present in `config.json`. The implementation report's "gates at PR CI" claim is false. Close by registering that live scenario into a `runtime-live-e2e.yml` `-run` selector; then re-validate.

### Feedback Cycles

- **Cycle 1 (validation REJECTED, 2026-06-14):** the binary half is solid (spawn-standing-all composition, AC-1 fixtures, AC-4 list-standing residency, structural shed, word-level audit all reproduced clean). Blocking defect: AC-2 — the load-bearing residency live scenario (FO boots → dispatches an ensign → comm-officer PRESENT in team config.json roster) is authored-as-intent but NOT registered/runnable in the CI live gate. runtime-live-e2e.yml has nothing that boots→dispatches→asserts the roster; the impl's "gates at PR CI" claim is false (same gap class as lean-boot cycle-1). Fix: author the actual live scenario test AND register it into runtime-live-e2e.yml's -run selector so CI invokes it; prove `go test -tags live -list` selects it. Routed back to implementation.
