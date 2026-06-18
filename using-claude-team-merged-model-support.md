---
id: 9243ej9j3a434c3k2ba3g06j
title: Bridge using-claude-team to .178+ merged team semantics (with a legacy-deprecation path)
status: ideation
source: "FO + captain live spike (2026-06-18) on Claude Code 2.1.181: validated the subagent-as-team-member machinery after anthropics/claude-code#68721 removed TeamCreate/TeamDelete (followup comment 4741152246, by bjcoombs, confirms the team->subagent merge is intentional). using-claude-team is built on the now-removed TeamCreate primitive and, on .178+, silently degrades to bare mode (surrendering concurrency) because it gates all dispatch on a TeamCreate that no longer exists."
started: 2026-06-18T14:55:53Z
completed:
verdict:
score:
worktree:
issue: "anthropics/claude-code#68721"
---

Make `skills/using-claude-team` work on Claude Code .178+ (the merged auto-team model, where `TeamCreate`/`TeamDelete` are gone and team membership comes from spawning named background subagents), while remaining compatible with legacy hosts that still expose `TeamCreate` — and mark a clear, triggered path to deprecate the legacy branch once the supported floor is .178+.

**Deliverable layers (true scope — established by `## Scope finding` below):** this is NOT prose-only. It is (a) skill/contract PROSE — `skills/using-claude-team/SKILL.md` + `skills/first-officer/references/claude-fo-dispatch.md` — PLUS (b) a Go change in `internal/dispatch/build.go`: a merged-mode dispatch-emission path (emit `name` present + `team_name` absent, past the Rule-8 `team_name`-required guard at `build.go:442`), with Go unit tests. The precise nuance: `internal/claudeteam/reconcile.go` itself needs **no change** (auto-discovery already keys on `leadSessionId`/`$CLAUDE_CODE_SESSION_ID`); the Go change is in the SEPARATE `internal/dispatch/build.go` file.

## Problem

`using-claude-team` makes `TeamCreate` the mandatory first call and **blocks all `Agent` dispatch until it resolves**. On .178+, `ToolSearch(select:TeamCreate)` returns no match, so the skill falls into **bare mode**: sequential, no concurrent dispatch, no direct ensign chat. Nothing breaks (all workflow functionality is preserved), but spacedock **needlessly surrenders concurrency** by gating on a primitive that no longer exists, while the primitives it actually needs (named background `Agent`, `SendMessage`, Tasks, on-disk team registry) are all present and working. The skill, the `claude-fo-dispatch.md` reconcile step, and the Terminal-Teardown / Degraded-Mode machinery are all written around `TeamCreate`/`TeamDelete` and the `#36806` registry-desync model, which no longer describe the runtime.

## Proposed approach (seed — ideation to flesh out)

Two-mode bridge with a deprecation trigger:

- **Detect mode**, do not assume it. Merged (.178+) vs legacy is observable: `ToolSearch(select:TeamCreate)` empty => merged; present => legacy. Ideation pins the exact detection contract and where it lives (generic skill vs `claude-fo-dispatch.md`).
- **Merged path (.178+):** no `TeamCreate`. Team membership = `Agent(name=…, run_in_background=true)`. Lead<->teammate via `SendMessage(to=name)` / teammate->lead via `SendMessage(to="main")`. Teardown = per-name `SendMessage(shutdown_request)` -> `shutdown_approved` -> `teammate_terminated` (NO `TeamDelete`, no active-member race, no `TERMINAL_TEARDOWN_BOUNDED` apparatus). FO tracks its own ensign roster (it already does — `TaskList` is not used).
- **Legacy path:** keep the existing `TeamCreate`/`TeamDelete` flow intact, clearly fenced and labelled deprecated, with a removal trigger (e.g. "remove when the min supported Claude Code no longer exposes `TeamCreate`, and no live lane still drives the legacy branch").
- **Strip the `claude-team`/NAME_PATTERN leak** from the generic skill (line ~19): the FO no longer names the team, and the introspection helper keys on `leadSessionId`, never the team name, so the compat note is doubly obsolete.
- **Flip `claude-fo-dispatch.md`'s reconcile step** from "pass your `TeamCreate` `{team_name}`" to the already-documented `leadSessionId` **auto-discovery** mode (drop `--team-name`; no `--session-id` plumbing is needed — OQ-5 confirms `reconcile.go` reads `$CLAUDE_CODE_SESSION_ID`). The reconcile Go helper (`internal/claudeteam/reconcile.go`) needs **no change**. NOTE (ideation correction, see `## Scope finding`): the seed's broader "no Go change" reading is INCOMPLETE — a SEPARATE Go file, `internal/dispatch/build.go`, DOES need a merged-mode emission path (Rule 8 at `build.go:442` blocks the `name`-present + `team_name`-absent merged dispatch shape). So the deliverable is prose PLUS that `build.go` change; only `reconcile.go` is change-free.

## Out of scope

- The live pty/team-mode harness work — owned by `m4` (`live-team-mode-terminal-harness`, m40mphxan8phr3t3tp03gk89). This task is the **skill + FO-contract prose** layer; m4 is the **live test harness** layer. They must stay reconciled but are separate deliverables.
- The team-mode verdict-omission question (reeppr990pyzzaejmbnyrvt7).
- Any change to the Go `internal/claudeteam` helper (`reconcile.go`) — the spike showed it is already robust to the merged model (auto-discovery keys on `leadSessionId`). This is the ONLY Go file out of scope: `internal/dispatch/build.go` IS in scope (the merged-mode emission path — see `## Scope finding`), so "no Go change" applies to `internal/claudeteam` specifically, not to the deliverable as a whole.
- **tmux pane-backed teammate reap** — descoped per captain (2026-06-18). The merged model spawns teammates `in-process` (validated), where lead `shutdown_request` reap works; pane-backed reap is not a path this task targets.

## Spike findings — validated live 2026-06-18 on Claude Code 2.1.181

Probed with real `Agent` + `SendMessage` (NOT the Workflow harness — a different code path). Four background subagents spawned (`scout`, `engineer`, `supervisor`, + a blocked `helper`); all reaped cleanly. Backend was `in-process`, model `claude-opus-4-8[1m]`.

1. **TeamCreate/TeamDelete are gone** on 2.1.181: `ToolSearch(select:TeamCreate,TeamDelete)` => "No matching deferred tools found". Intentional per #68721 followup.
2. **Named background subagent works with zero TeamCreate.** `Agent(name=…, run_in_background=true)` spawns an addressable teammate.
3. **subagent -> lead channel is OPEN.** A background subagent's `SendMessage(to="main")` was delivered to the root verbatim. (Old model: a raw subagent could only reach the lead via its final return value.)
4. **Standing member works.** A persistent teammate retained state across waves (recalled a secret number 42 -> 84 and an incrementing counter across separate inbound messages), was addressable each round, and replied to the lead each round.
5. **Reap works via lead-initiated shutdown_request.** `SendMessage(shutdown_request)` -> teammate emits `shutdown_approved` -> `teammate_terminated`. This CONTRADICTS the #68721 comment's claim that "the lead's shutdown_request isn't honored" — at least for `in-process` backend on 2.1.181.
6. **Self-exit is NOT available — and this is NOT a regression.** A teammate can only send `shutdown_response` as a REPLY to an inbound `shutdown_request`; it cannot originate termination. Old `TeamCreate` members were identical. The comment's "can't send shutdown_response" was likely observed because no `shutdown_request` was ever sent (no `request_id` to echo).
7. **Roster is FLAT by construction.** A background teammate CANNOT spawn background grandchildren: `Agent(name=…, run_in_background=true)` from a teammate => "Teammates cannot spawn other teammates — the team roster is flat"; nameless `run_in_background=true` => "In-process teammates cannot spawn background agents. Use run_in_background=false for synchronous subagents." So the "orphan-grandchild" hazard cannot occur; recursive *background* depth is bounded to 1. Synchronous (blocking) subagents at depth are still allowed — so an ensign running a fan-out skill/Workflow still works.
8. **On-disk registry is intact + accurate.** The auto-team writes `~/.claude/teams/session-<id>/config.json`, named by **session id** (not the FO scheme), with the SAME schema (`name`/`leadAgentId`/`leadSessionId`/`members[]`/`inboxes/`). `members[]` is a **live pruned roster** (terminated members removed; `inboxes/*.json` linger). Each member records `agentType` = the spawn `subagent_type` (probe recorded `general-purpose`).
9. **The `claude-team` helper survives unchanged.** `internal/claudeteam/reconcile.go` auto-discovers by globbing all team dirs and matching `leadSessionId == current session` AND a member with `agentType == "spacedock:ensign"`; `--team-name` is only an optional override. So the session-id rename does NOT break discovery — the helper never used the team name. Only the CALLER must switch to auto-discovery mode (no `TeamCreate` name to pass).

## Capability matrix — Claude Code 2.1.181 (merged auto-team model)

Tested context: Claude Code **2.1.181**. The standing-member / reap / registry rows were observed via the interactive `cli` entrypoint (`in-process` teammate backend, model `claude-opus-4-8[1m]`, `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`, no tmux). The three rows that were ❓ UNTESTED in that interactive session — headless `-p` residency, flag-free operation, and `agentType` for a non-`general-purpose` `subagent_type` — were SUBSEQUENTLY resolved during ideation (see `## Open-question resolutions`): residency + flag-free via a nested `claude -p` stream-json spike (so those two are now **headless-validated**, not interactive-only), and the `agentType` stamping via on-disk registry inspection. Each ✅/❌ was observed live; no row remains ❓.

| Capability | 2.1.181 | Evidence / exact signal |
|---|---|---|
| `TeamCreate` / `TeamDelete` tools | ❌ absent (by design) | `ToolSearch(select:TeamCreate,TeamDelete)` → "No matching deferred tools found" |
| Spawn named background teammate without `TeamCreate` | ✅ | `Agent(name=…, run_in_background=true)` returns a live addressable teammate |
| Teammate → lead message | ✅ | `SendMessage(to="main")` delivered to the root verbatim |
| Lead → teammate message | ✅ | `SendMessage(to=<name>)` reaches the teammate inbox |
| Standing teammate persists + retains state across waves | ✅ | recalled a secret value + an incrementing counter across separate inbound messages |
| Lead-initiated reap via `shutdown_request` | ✅ (in-process) | `shutdown_request` → teammate `shutdown_approved` → `teammate_terminated` |
| Teammate self-originated exit | ❌ (NOT a regression) | `shutdown_response` is valid only as a REPLY to an inbound `shutdown_request`; legacy `TeamCreate` members were identical |
| Background teammate spawns a background **grandchild** | ❌ blocked (flat roster) | "Teammates cannot spawn other teammates — the team roster is flat"; nameless bg → "In-process teammates cannot spawn background agents. Use run_in_background=false for synchronous subagents." |
| Teammate spawns a **synchronous** subagent | ✅ allowed | per the harness's own error guidance (`run_in_background=false`) |
| `TaskList` / `TaskStop` reach background teammates | ❌ | `TaskStop <id>` → "No task found"; `TaskList` → empty |
| On-disk team registry | ✅ present | `~/.claude/teams/session-<id>/config.json`; same schema as the `TeamCreate` era; `members[]` is a live, pruned roster (terminated members removed) |
| Headless **`-p`** team mode (residency past `end_turn`) | ✅ | nested `claude -p` stream-json probe — spawn→deliver→reap within one resident `-p` run (both flag arms) |
| **Flag-free** (no `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`) | ✅ | OQ-2 probe — identical behavior flag-unset vs `=1`; flag is a no-op for the merged channel |
| `agentType` stamped for a non-`general-purpose` `subagent_type` (e.g. `spacedock:ensign`) | ✅ | OQ-3 — real on-disk `config.json`s stamp `agentType:"spacedock:ensign"`; FO-side fresh-dispatch confirm still flagged |

### What changed, 2.1.177 → 2.1.181

| | 2.1.177 (legacy `TeamCreate`) | 2.1.181 (merged) |
|---|---|---|
| Team tools | `TeamCreate`/`TeamDelete` present | absent (by design) |
| Establish membership | `TeamCreate`, then named members | named background `Agent`, no `TeamCreate` |
| Lead reap | `TeamDelete` (whole team) + per-member shutdown | per-member `shutdown_request` only (no bulk teardown) |
| Enable flag | `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` needed | flag-free |
| Team dir | caller-named | auto, `session-<id>` |
| Headless `-p` team mode | ❌ (the reason the m4 pty harness exists) | ✅ (merged channel stays resident) |

The 2.1.181 column is this session's live probes; the 2.1.177 column is the documented legacy `TeamCreate` model + m4's recorded pinned-2.1.177 pty spike (NOT re-tested this session).

Unchanged across both: teammate self-exit is reply-only; the on-disk `config.json` schema is the same shape.

## Open questions (for ideation)

- **OQ-1 — DESCOPED (captain, 2026-06-18): tmux pane-backed reap is out of scope.** The merged model spawns teammates `in-process`, where lead `shutdown_request` reap is validated (see capability matrix). Whether the #68721 comment's "shutdown not honored" was pane-specific is no longer a question this task chases.
- **OQ-2: does `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` still gate the named-background-subagent + `SendMessage(to="main")` channel on .178+, or is it now flag-free?** During the spike the flag was PRESENT (`=1` in the session env), so flag-free operation is UNTESTED. m4's CI sets the flag. The skill's bare-mode hint tells captains to set it — possibly stale.
- **OQ-3: does a real `subagent_type="spacedock:ensign"` dispatch record `agentType:"spacedock:ensign"` in the auto-team config**, so `reconcile.go`'s `hasEnsign` gate fires? Needs a one-shot real-ensign smoke test (probe used `general-purpose`, which would be correctly excluded).
- **OQ-4: exact mode-detection contract.** Is `ToolSearch(select:TeamCreate)` presence reliable across 2.1.177 <-> .178 <-> .181, and does it belong in the generic skill or the spacedock FO layer?
- **OQ-5: is `--session-id` plumbed** through `spacedock claude`/the launcher so the reconcile sweep can run in `leadSessionId` auto-discovery mode without a `TeamCreate` name?
- **OQ-6: reconcile with m4.** m4's transcript-observation work and this skill rewrite both touch team-mode semantics; ideation must define the seam so they do not collide (skill/contract here, live harness there).

## Provisional acceptance criteria (ideation finalizes; proof = behavior, never prose-grep)

Per this workflow's proof policy a skill change is PASSED only when a LIVE drive observed the behavior — a string match over the skill text proves nothing.

- **AC-1 (provisional):** On a .178+ host, a live FO boot dispatches a teammate concurrently WITHOUT calling `TeamCreate` (no bare-mode fallback). Verified by: a live drive (or the m4 pty harness) observing concurrent dispatch + the auto-team `config.json` on disk.
- **AC-2 (provisional):** On a legacy host that still exposes `TeamCreate`, the legacy path still runs. Verified by: the legacy live lane (or a fixture that stubs `TeamCreate` present) still green.
- **AC-3 (provisional):** Terminal teardown on .178+ reaps every tracked ensign via per-name `shutdown_request` (no `TeamDelete`), leaving `members[]` empty. Verified by: live drive observing `teammate_terminated` for each + empty roster.
- **AC-4 (provisional):** The reconcile sweep resolves the auto-team via `leadSessionId` auto-discovery with no `TeamCreate` name passed. Verified by: `spacedock dispatch reconcile` (no `--team-name`) resolving the live session's roster.
- **AC-5 (provisional):** The legacy branch is fenced + labelled deprecated with an explicit removal trigger. Verified by: the deprecation trigger references an external, checkable condition (supported-version floor / live-lane coverage), not just prose.

## Related

- `m4` — live-team-mode-terminal-harness (m40mphxan8phr3t3tp03gk89, PR #390): live-harness layer of the same regression; OQ-1/OQ-2/OQ-3 are partly answerable there.
- `team-mode-verdict-omission` (reeppr990pyzzaejmbnyrvt7).
- anthropics/claude-code#68721 + comment 4741152246 (intentional team->subagent merge).
- Spike transcript: this FO session (2026-06-18), `~/.claude/teams/session-a15e8296/`.

## Open-question resolutions (ideation spike, 2026-06-18, claude 2.1.181)

The live OQs were exercised FIRST, before any wording, via two nested `claude -p` Bash probes (`--output-format stream-json`). Nested-session env markers (`CLAUDECODE`, `CLAUDE_CODE_CHILD_SESSION`, `CLAUDE_CODE_SESSION_ID`, `CLAUDE_CODE_AGENT`, `CLAUDE_CODE_ENTRYPOINT`, `CLAUDE_CODE_EXECPATH`) were `env -u`-scrubbed so the child behaved like a top-level session (the m4 transcript-suppression lesson, commit `2b4beaa0`). No auth/quota blocker hit — both probes ran clean on the machine credential.

- **OQ-2 — RESOLVED: flag-free works; the bare-mode flag hint is STALE.** Ran the identical probe twice — arm A with `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` **unset**, arm B with it `=1`. Probe shape: lead loads `SendMessage`, spawns `Agent(name="pinger", run_in_background=true)`, the teammate runs `SendMessage(to="main", "PING_FROM_TEAMMATE")`, the lead must report receipt. **Both arms behaved identically:** `system init` carried `SendMessage` in `tools[]` and `TeamCreate`/`TeamDelete` absent; the `Agent` spawn returned "Async agent launched successfully"; the teammate's `SendMessage(to="main")` reached the lead (which printed `RESULT: RECEIVED_TEAMMATE_MESSAGE`); the teammate was reaped (`task_notification status: completed`). So on 2.1.181 the merged named-background-subagent + `SendMessage(to="main")` channel is flag-free — the flag is a no-op for it. The skill's bare-mode UX hint ("set `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` and restart") is obsolete on the merged floor and must go. (Version nuance, recorded: m4 found 2.1.179 *interactive* still needed the flag to scaffold the auto-team while native `TeamCreate` was already gone; on 2.1.181 *headless* the merged channel needs no flag. The flag gated the now-absent native-tool path, not the merged channel. Detection keys on tool presence, not the flag — see OQ-4 — so this version drift does not destabilize the contract.)
- **Headless `-p` residency — RESOLVED (folded into OQ-2): a `claude -p` stays resident through spawn→deliver→reap.** In both arms the `-p` session did not race to `end_turn` and exit before the teammate's message arrived — it received the inbound `SendMessage`, printed the receipt, and observed the teammate's terminal `task_notification` all within one `-p` invocation. This CONTRADICTS m4's headless-`-p`-goes-bare premise — but m4 tested the OLD native-`TeamCreate` path (gone on .178+). The MERGED named-background-subagent path is what stays resident headless. (This does not re-open m4's go-bare decision for the legacy regime; it scopes the bare-fallback to legacy hosts only — see the two-mode bridge below.)
- **OQ-3 — RESOLVED via on-disk registry inspection: a real `spacedock:ensign` dispatch stamps `agentType:"spacedock:ensign"`.** The flat roster blocks me from dispatching a named ensign from inside my own context, so I read the live registry instead: dozens of on-disk `~/.claude/teams/*/config.json` from real spacedock workflow sessions (002, dataagentbench, buggy-add-task, email-triage) carry members spawned with `subagent_type="spacedock:ensign"` stamped `"agentType":"spacedock:ensign"`, exactly the discriminator `reconcile.go`'s `hasEnsign` gate matches (`m.AgentType == "spacedock:ensign"`). The registry records the spawn's `subagent_type` faithfully. **FO-side confirmation still recommended** (the flat-roster rule blocked a fresh same-session dispatch): the FO should, on the first merged-mode live drive, dispatch one real ensign and confirm the live session's `config.json` member carries `agentType:"spacedock:ensign"` and that `spacedock dispatch reconcile` (no `--team-name`) resolves the roster. FLAGGED for the FO.
- **OQ-4 — RESOLVED: detection contract is `ToolSearch(query="select:TeamCreate")` empty ⟺ merged; lives in the generic skill.** Both 2.1.181 arms show `TeamCreate`/`TeamDelete` absent from `tools[]` and `ToolSearch(select:TeamCreate)` is the documented hop already in the skill. Detection is a host-capability fact, not a spacedock-policy fact, so it belongs in the generic `using-claude-team` skill (the same file that owns the deferred-tool ToolSearch hop), with the FO layer consuming the resolved mode. The probe is the existing pre-`TeamCreate` ToolSearch the skill already prescribes — no new mechanism. Reliability across 2.1.177↔.178↔.181: m4 independently verified `select:TeamCreate` returns "No matching deferred tools found" on 2.1.179 via the real hop (not an init eyeball), and present on 2.1.177; this probe confirms absent on 2.1.181. The discriminator is stable across the band.
- **OQ-5 — RESOLVED: no `--session-id` plumbing is needed; auto-discovery already runs off `$CLAUDE_CODE_SESSION_ID`.** Code read: `runClaude` (`internal/cli/frontdoor.go:268-353`) assembles `claude --agent spacedock:first-officer [--permission-mode auto] {passthrough} {prompt}` and does NOT inject a `--session-id` (the `--session-id` entry at `frontdoor.go:560` is only a `valueTakingHostFlags` parser hint so a passthrough `--session-id` token's successor is not mis-flagged as a stray positional). Claude assigns the FO's own session id at launch and sets it in the FO session env as `$CLAUDE_CODE_SESSION_ID` (verified present in this live FO session). `spacedock dispatch reconcile` reads exactly that (`internal/dispatch/reconcile.go:131` `sessionID: os.Getenv("CLAUDE_CODE_SESSION_ID")`) and passes it to `claudeteam.LoadReconcileTeam`, which matches `leadSessionId == sessionID` among ensign-bearing configs. So the `leadSessionId` auto-discovery path is already wired end-to-end with no `TeamCreate` name and no launcher change. The ONLY change is at the caller layer: `claude-fo-dispatch.md`'s reconcile step must stop telling the FO to pass `--team-name {team_name}` (a name it no longer has on the merged floor) and let the bare-reconcile session-match run.
- **OQ-6 — RESOLVED: seam with m4 is the host-regime split.** m4 (live pty harness) owns the **legacy-`TeamCreate` regime**: its CI live-e2e is pinned to claude **2.1.177** (`runtime-live-e2e.yml:111` `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"`, with `DISABLE_AUTOUPDATER: "1"`) precisely because 2.1.177 is the last release exposing native team tools — that pin exists ONLY to keep m4's `TeamCreate`-based AC-3/AC-4 green. THIS task owns the **merged-model regime** (current/unpinned claude, no `TeamCreate`). The skill must serve BOTH because CI runs the pinned legacy host while real users run the merged host. Non-overlap: m4 changes NO skill/contract prose (it is a test-harness layer); this task changes NO live-harness Go (it is the skill/contract layer). The one shared surface — the deprecation trigger — is owned here and *references* m4's pin (below), so the two stay reconciled through one externally-checkable condition rather than colliding.

## Scope finding — the Go change the seed missed (`internal/dispatch/build.go`, NOT `internal/claudeteam`)

The seed says "the Go `internal/claudeteam` helper needs NO change." That is TRUE and verified for `internal/claudeteam/reconcile.go` (the spike + this code read confirm auto-discovery already keys on `leadSessionId`). But a SEPARATE Go file gates the merged path and DOES need a change: **`internal/dispatch/build.go` Rule 8** (`build.go:442`):

    if !bareMode && host == "claude" && teamName == "" {
        return buildError(stderr, 1, "team mode requires team_name")
    }

`spacedock dispatch build` emits `out.Name` only when `!bareMode` (`build.go:648-649`) and `out.TeamName` only when `teamName != ""` (`build.go:650-651`). With Rule 8 in place, on `host=claude` exactly two states are reachable: **bare** (`name` absent, `team_name` absent → sequential, the concurrency we are trying to recover) or **legacy team** (`name` present, `team_name` present → requires `TeamCreate`). There is **no path that emits `name` present + `team_name` absent** — which is precisely the merged dispatch shape (`Agent(name=…, run_in_background=true)`, no `TeamCreate`, no `team_name`). So the merged path requires a `build.go` change: a third "merged team mode" that emits a `name` (for addressability + reuse-advance `SendMessage`) with NO `team_name`. Adjacent `teamName != ""` gates also assume team_name-present implies team mode and must be re-read for the merged case: the standing-teammate fetch line (`build.go:570`), the completion-signal block (`build.go:583`), and the dispatch-file naming (`build.go:606`, which keys the collision-free dispatch-file path on `team_name` for uniqueness — the merged path needs a different disambiguator, e.g. the session id, since there is no unique team name).

This is recorded as the load-bearing implementation-scope correction: the deliverable is NOT "skill prose only." It is **skill/contract prose PLUS a `build.go` merged-mode emission path** (with its own Go unit tests). The reconcile helper genuinely needs no change; the dispatch-build helper does. This is exactly the unverified-mechanism the proof policy demands ideation surface before the gate.

## Two-mode bridge — concrete before/after wording

Detection (OQ-4) selects the path once, at the first team-mode tool call, in `using-claude-team`. Merged is the default forward path; legacy is fenced + deprecated.

### Change 1 — `skills/using-claude-team/SKILL.md` frontmatter description (line 3)

BEFORE:
> Invoke when orchestrating any Claude Code agent team (TeamCreate/Agent/SendMessage), independent of any specific workflow.

AFTER:
> Invoke when orchestrating any Claude Code agent team (named background `Agent` + `SendMessage`; legacy `TeamCreate` on pre-2.1.178 hosts), independent of any specific workflow.

### Change 2 — `SKILL.md` Deferred Team Tools (line 13)

BEFORE:
> The Claude Code team tools (`TeamCreate`, `TeamDelete`, `SendMessage`, and related team-registry tools) are deferred — their schemas are not loaded at session start, so calling one directly fails until its schema is fetched. Before the first call to any team tool, run `ToolSearch(query="select:{ToolName}", max_results=1)` to fetch its schema (e.g. `ToolSearch(query="select:TeamCreate", max_results=1)` before the first `TeamCreate`, `ToolSearch(query="select:SendMessage", max_results=1)` before the first `SendMessage`). Once a tool's schema appears in the ToolSearch result, it is callable exactly like a normal tool. `Agent` is not deferred. Address an agent by its declared `name` via `SendMessage`; your plain text output is NOT visible to other agents.

AFTER:
> The Claude Code team tools (`SendMessage`, and on legacy hosts `TeamCreate`/`TeamDelete`) are deferred — their schemas are not loaded at session start, so calling one directly fails until its schema is fetched. Before the first call to any team tool, run `ToolSearch(query="select:{ToolName}", max_results=1)` to fetch its schema (e.g. `ToolSearch(query="select:SendMessage", max_results=1)` before the first `SendMessage`). This same probe is the **mode discriminator**: `ToolSearch(query="select:TeamCreate", max_results=1)` returning a `TeamCreate` definition means a **legacy** host (run the fenced `TeamCreate` path below); returning no match means a **merged** host (.178+: team membership comes from named background `Agent`, no `TeamCreate`). Once a tool's schema appears in the ToolSearch result, it is callable exactly like a normal tool. `Agent` is not deferred. Address an agent by its declared `name` via `SendMessage`; your plain text output is NOT visible to other agents.

### Change 3 — `SKILL.md` Team Creation, replace the `TeamCreate`-first mandate (lines 15-34)

The whole `## Team Creation` block (the `TeamCreate`-MUST-be-first mandate, the recovery procedure, the failure-recovery ladder, "Block all Agent dispatch until team setup resolves") is rewritten as a mode split. The merged path becomes the default; the legacy `TeamCreate` machinery is moved verbatim under a fenced, deprecated subsection.

BEFORE (line 19, the mandate):
> 1. **Probe for TeamCreate and run it first.** `TeamCreate` MUST be the first team-mode tool call in every session, before ANY `Agent` or `SendMessage` invocation. Run `ToolSearch(query="select:TeamCreate", max_results=1)`. If the result contains a TeamCreate definition, derive `{project_name}` … then run `TeamCreate(…)`. The timestamp token must be lowercase and hyphen-separated — no uppercase, no colons — to stay compatible with Claude Code's NAME_PATTERN and the `claude-team` helper. …

AFTER (new `## Team Setup` opening):
> At the first team-mode tool call, run the mode discriminator `ToolSearch(query="select:TeamCreate", max_results=1)`.
>
> **Merged host (no `TeamCreate` match — .178+, the default forward path):** there is NO team-creation step. Team membership is established by dispatching named background teammates: `Agent(name=…, run_in_background=true)`. The on-disk auto-team `~/.claude/teams/session-<id>/config.json` is written by Claude Code automatically, keyed by session id, recording each member's `agentType`. Do NOT call `TeamCreate`. Do NOT block `Agent` dispatch on any team-setup step — dispatch directly. Lead→teammate is `SendMessage(to=<name>)`; teammate→lead is `SendMessage(to="main")`.
>
> **Legacy host (`TeamCreate` match — pre-2.1.178, DEPRECATED):** follow the fenced `### Legacy TeamCreate path (deprecated)` subsection below.

The `NAME_PATTERN` / `claude-team` compat sentence is DELETED from the merged path entirely (the FO no longer names the team; the introspection helper keys on `leadSessionId`, never the team name — the compat note is doubly obsolete, per the seed). It survives only inside the fenced legacy subsection where `TeamCreate` still takes a name.

### Change 4 — `SKILL.md` bare-mode flag hint (line 22), delete the stale UX hint

BEFORE:
> 2. If ToolSearch returns no match, enter **bare mode**: dispatch is sequential … When reporting bare mode at startup, append this one-line UX hint verbatim … `Tip: set CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 and restart this session to enable team mode for concurrent dispatch and direct ensign chat (Shift+Up/Shift+Down).` Do not repeat the hint after startup.

AFTER:
> On a merged host, `ToolSearch(select:TeamCreate)` returning no match is NOT a bare-mode trigger — it is the normal merged path (Change 3). **Bare mode** is now entered only when `Agent`/`SendMessage` themselves are unavailable (a genuinely degraded host) or by explicit operator command. The `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` UX hint is REMOVED — OQ-2 proved the merged named-background-`Agent` + `SendMessage(to="main")` channel works flag-free on 2.1.181, so the hint points at a flag that no longer gates concurrency.

(The "TeamCreate-no-match ⟹ bare mode" wiring is the exact bug this task fixes: it is what makes spacedock surrender concurrency on .178+. The fenced legacy subsection keeps the real bare-mode fallback for a host where `Agent` itself is absent.)

### Change 5 — `SKILL.md` Terminal Team Teardown (lines 98-109), mode-split the teardown

BEFORE: the whole `## Terminal Team Teardown` block attempts `TeamDelete` with the bounded `TERMINAL_TEARDOWN_BOUNDED` settle-and-cap apparatus.

AFTER: split by mode.
> **Merged host:** there is NO `TeamDelete`. Tear down per-roster-member with `SendMessage({"type":"shutdown_request"})` → the member emits `shutdown_response`/`shutdown_approved` → `teammate_terminated`; the auto-team `members[]` prunes the terminated member (spike finding #8: the live roster is pruned). There is no `active member(s)` race and no `TERMINAL_TEARDOWN_BOUNDED` apparatus on this path. The FO tracks its own ensign roster (it already does).
> **Legacy host:** the existing `TeamDelete` bounded-teardown block (settle, attempt-cap, `TERMINAL_TEARDOWN_BOUNDED` marker) is preserved verbatim under the fenced legacy subsection — it remains the contract for a host that still has `TeamDelete`.

### Change 6 — fence the legacy machinery under one labelled, deprecated subsection with a removal trigger

All `TeamCreate`/`TeamDelete`-specific prose (the `TeamCreate`-first naming + NAME_PATTERN, the TeamCreate recovery procedure, the failure-recovery ladder, the `TeamDelete` bounded teardown, the diagnostic on-disk-teams probe, the `#36806` registry-desync model) is moved UNCHANGED under a single subsection:

> ### Legacy TeamCreate path (DEPRECATED — pre-2.1.178 hosts only)
>
> This path runs only when the mode discriminator finds a `TeamCreate` tool (Claude Code < 2.1.178). It is retained for hosts that still expose the native team registry. **Removal trigger (externally checkable):** remove this entire subsection — and the `TeamCreate`/`TeamDelete` branches in `internal/dispatch/build.go` and the FO contract — when `SPACEDOCK_PINNED_CLAUDE_VERSION` in `.github/workflows/runtime-live-e2e.yml` no longer pins a team-tools-capable version (i.e. moves to ≥ 2.1.178 or the pin is dropped) AND no live lane drives the legacy branch. That pin is the single live consumer of this path today (it exists, per its own comment, only because 2.1.177 is the last release exposing native team tools); when it stops pinning ≤ 2.1.177, the legacy branch has no live coverage and is dead code. [Then the contents follow, moved verbatim.]

### Change 7 — `skills/first-officer/references/claude-fo-dispatch.md` reconcile step (line 96), flip to auto-discovery

BEFORE:
> 0. **Reconcile sweep.** Run `spacedock dispatch reconcile --workflow-dir {workflow_dir} --team-name {team_name}` … Pass your own `TeamCreate` `{team_name}` — the roster-derived classes (lingering/superseded/un-advanced-pr) require a team identity … The team identity comes from either the explicit `--team-name {team_name}` or a current-session match …

AFTER:
> 0. **Reconcile sweep.** Run `spacedock dispatch reconcile --workflow-dir {workflow_dir}` (no `--team-name` on the merged host — you have no `TeamCreate` name to pass). The team identity resolves by **`leadSessionId` auto-discovery**: the helper narrows to the auto-team `config.json` whose `leadSessionId` equals this session's `$CLAUDE_CODE_SESSION_ID` (set by Claude Code at launch, read by the reconcile command — no launcher plumbing needed). The roster-derived classes (lingering/superseded/un-advanced-pr) are emitted against that session-matched roster. On a **legacy** host that still has a `TeamCreate` name, pass `--team-name {team_name}` as before (the explicit override path is unchanged). Bare reconcile with no resolvable session team stays git-only (the existing degrade).

### Change 8 — `claude-fo-dispatch.md` Team Creation references (lines 5-21), mode-split the dispatch contract

The `## Team Creation` and `## Spawn Call (Agent)` sequencing rules that assume a `TeamCreate`-first ordering are rewritten so the merged path dispatches named background teammates with no `TeamCreate` and no `team_name`, consuming `spacedock dispatch build`'s new merged-mode output (`name` present, `team_name` absent — see the `build.go` scope finding). The `TeamCreate`/`spawn-standing-all`-requires-team_name sequencing rule (line 21) is fenced to the legacy path. Standing-teammate injection (line 55, `spawn-standing-all --team {team_name}`) needs the merged-mode analog: on the merged host the standing teammate is injected by a named background `Agent` dispatch with no `team_name`, scoped to the session auto-team rather than a named team. (Concrete `build.go`-output-shape wording for this is finalized in implementation against the new merged-mode emission; ideation pins the contract: merged dispatch = `name` yes, `team_name` no, `TeamCreate` never.)

## Acceptance criteria (finalized — entity-level, behavior-anchored; proof = live drive, never prose-grep)

Per this workflow's proof policy a skill/contract change is PASSED only when a LIVE drive observes the behavior. Each AC names its fixture coverage (offline, cheap, runs in `go test`) AND its live coverage (the behavior oracle). A string match over the skill text is explicitly NOT proof for any AC.

- **AC-1 — Merged-host concurrent dispatch with no `TeamCreate`.** On a merged host (current claude, `ToolSearch(select:TeamCreate)` empty), a live FO boot dispatches an ensign as a named background teammate WITHOUT calling `TeamCreate` and WITHOUT falling into sequential bare mode.
  - *Fixture:* a `build.go` Go unit test asserting merged-mode `spacedock dispatch build` (host=claude, not bare, merged) emits `name` present + `team_name` absent + no "team mode requires team_name" error (the new Rule-8 path).
  - *Live:* a live merged-host drive (the m4 pty harness on an unpinned/current claude, or a manual FO drive) observes the FO issuing `Agent(name=…, run_in_background=true)` with no `TeamCreate` in the stream, and the auto-team `config.json` on disk gaining the ensign member with `agentType:"spacedock:ensign"`. Spike-proven mechanism: named background `Agent` + `SendMessage` channel works flag-free (OQ-2 probe).
- **AC-2 — Legacy host still runs the `TeamCreate` path.** On a legacy host that exposes `TeamCreate` (the CI pin, claude 2.1.177), the fenced legacy path runs unchanged: `TeamCreate` → dispatch → `TeamDelete` bounded teardown.
  - *Fixture:* a `build.go` unit test asserting legacy-mode build (team_name present) still emits `name` + `team_name` as today; the offline teardown-grade fixtures (`teardown_grade_watcher_test.go`) are untouched.
  - *Live:* m4's existing pinned-2.1.177 CI live lane (AC-3/AC-4) staying green is the legacy-path live proof — this task must not break it. (Reconciled via OQ-6 seam.)
- **AC-3 — Merged-host teardown reaps via `shutdown_request`, no `TeamDelete`.** On a merged host, terminal teardown reaps every tracked ensign with per-name `SendMessage(shutdown_request)` and the auto-team `members[]` prunes them — no `TeamDelete`, no `TERMINAL_TEARDOWN_BOUNDED` apparatus.
  - *Fixture:* none new (no Go behavior change for teardown — it is FO-contract prose); the discriminator is exercised by AC-1's build test.
  - *Live:* a live merged-host drive observes `teammate_terminated` (or `task_notification status: completed`) for each reaped ensign and the on-disk `members[]` shrinking. Spike-proven: lead-initiated `shutdown_request` reap works in-process (spike finding #5); the auto-team roster is live-pruned (finding #8); the OQ-2 probe observed a clean reap.
- **AC-4 — Reconcile resolves via `leadSessionId` auto-discovery with no `--team-name`.** `spacedock dispatch reconcile` (no `--team-name`) resolves the live merged session's roster by matching `$CLAUDE_CODE_SESSION_ID`, and the FO contract instructs the no-`--team-name` form on merged hosts.
  - *Fixture:* the EXISTING `internal/dispatch/reconcile_session_test.go` (AC-1..AC-5: foreign-team-never-poisons, session-matched discovery, explicit-`--team-name` override, degrade-to-git-only) already proves the helper behavior offline — no new helper test needed; this AC adds a contract-prose change, proven live.
  - *Live:* a live merged drive runs the FO's reconcile step (no `--team-name`) and it resolves the live session's ensign roster (non-empty `drift`/roster classes against the real auto-team), confirming the env-var session identity flows end-to-end. Spike-proven: `$CLAUDE_CODE_SESSION_ID` is set in the live FO session; `reconcile.go:131` reads it; the helper matches `leadSessionId`.
- **AC-5 — Legacy branch is fenced + deprecated with an externally-checkable removal trigger.** The legacy `TeamCreate`/`TeamDelete` machinery lives under one labelled deprecated subsection whose removal trigger references an external, checkable condition — not prose.
  - *Fixture:* a cheap structural check (a Go test or a CI guard) asserting the removal trigger names the real `SPACEDOCK_PINNED_CLAUDE_VERSION` pin in `.github/workflows/runtime-live-e2e.yml`, and that the pin currently still pins ≤ 2.1.177 (so the legacy branch still has a live consumer). This is the one AC whose proof is structural rather than behavioral, because it is a deprecation *contract*, not a runtime behavior — but it is anchored to a real grep-able file+value, not free prose.
  - *Live:* none — the trigger is a documentation+CI contract; its "liveness" is that AC-2's live lane (the 2.1.177 pin) is the very consumer the trigger names. When that pin moves to ≥ 2.1.178, AC-2's live lane goes dark and the trigger fires.

## Test plan (cost/complexity, fixture vs live, what verifies each)

- **Offline (cheap, runs in `go test ./...`):** the `build.go` merged-mode emission is the bulk of the testable mechanism — new Go unit tests for the third dispatch mode (merged: `name` yes, `team_name` no, no Rule-8 error; standing-teammate + completion-signal + dispatch-file-naming gates behave for the no-`team_name` merged case). The reconcile auto-discovery is ALREADY covered by `reconcile_session_test.go` (reuse, do not duplicate). AC-5's trigger gets a cheap structural guard. Est: ~half a day of Go test work; behavior-preserving for the legacy path (legacy build tests stay green).
- **Live (multi-minute, credentialed, CI-gated):** AC-1/AC-3/AC-4 are the merged-host behavior claims. They need a live drive on a MERGED (current/unpinned) claude — which the CI live-e2e currently does NOT run (it pins 2.1.177 for m4's legacy tests). So this task's live coverage requires either (a) a SECOND live-e2e lane on unpinned/current claude that drives the merged path (recommended; it also becomes the regression sensor for the merged floor), or (b) a manual FO live drive recorded in the validation stage. The two-lane split (legacy pinned 2.1.177 for m4; merged unpinned for this task) is the live realization of the OQ-6 seam. AC-2's live proof is m4's existing pinned lane staying green.
- **Riskiest-path-first ordering for implementation:** build the `build.go` merged-mode emission and its unit tests FIRST (it is the mechanism everything else rests on, and the seed mis-scoped it as "no Go change"), then the skill/contract prose, then the live merged-lane drive. A merged-mode dispatch that cannot be emitted by `build.go` would invalidate the whole skill rewrite — so it goes first.

## Docs impact

User-visible surface check: this task changes FO/skill *agent-facing* contract prose, not the docs-site or CLI output. The one user-observable behavior change — on a merged host the FO no longer prints the stale `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` bare-mode tip and no longer silently drops to sequential bare mode — is internal agent behavior, not a documented CLI surface, banner, or docs-site page. No `spacedock --help`, startup banner, or docs-site page describes the team-mode flag hint. So no doc diff is required at the ideation gate. (If a future docs-site page documents team mode, it would inherit the merged-default framing; none exists today — confirmed: the hint lives only in `skills/using-claude-team/SKILL.md`.)

## Stage Report: ideation

- DONE: Remaining open questions each resolved by a recorded spike result OR an auditable "no spike needed" naming the proven mechanism — the live ones (headless `-p` team-mode residency; flag-free operation without `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS`) exercised first via a nested `claude -p` Bash probe, with any auth/quota blocker recorded honestly rather than guessed.
  All OQs resolved in `## Open-question resolutions`. OQ-2 (flag-free) + headless `-p` residency: two nested `claude -p` `--output-format stream-json` probes (flag unset vs `=1`, nested-session markers `env -u`-scrubbed) — both spawned a named bg teammate, received its `SendMessage(to="main")`, reaped it, all within one resident `-p` run; flag is a no-op. No auth blocker hit. OQ-3: on-disk registry inspection (flat-roster blocks self-dispatch) — dozens of real `config.json` stamp `agentType:"spacedock:ensign"`; FO-side fresh-dispatch confirm FLAGGED. OQ-4: `ToolSearch(select:TeamCreate)` empty⟺merged, lives in the generic skill. OQ-5: code read — no `--session-id` plumbing needed (`reconcile.go:131` reads `$CLAUDE_CODE_SESSION_ID`). OQ-6: host-regime seam (m4 owns pinned-2.1.177 legacy; this task owns merged).
- DONE: Two-mode bridge specified as concrete before/after `using-claude-team` wording: the merged `.178+` path, the fenced+deprecated legacy `TeamCreate` path with an externally-checkable removal trigger, plus the `claude-fo-dispatch.md` reconcile-step flip to `leadSessionId` auto-discovery.
  `## Two-mode bridge — concrete before/after wording`: 8 changes with verbatim before/after — SKILL.md description/deferred-tools/team-creation/bare-mode-hint/teardown + fenced legacy subsection with the removal trigger anchored to `SPACEDOCK_PINNED_CLAUDE_VERSION` in `runtime-live-e2e.yml`; claude-fo-dispatch.md reconcile-step flip to no-`--team-name` auto-discovery + the dispatch-contract mode split.
- DONE: Acceptance criteria are entity-level and behavior-anchored (proven by a live drive, never a prose-grep over the skill text), each with a test plan naming fixture vs live coverage.
  `## Acceptance criteria` AC-1..AC-5 each name fixture (offline `go test`) + live (drive oracle) coverage; AC-5 is the one structural-not-behavioral case (a deprecation contract) anchored to a real grep-able file+value. `## Test plan` orders riskiest-path-first (build the `build.go` merged-mode emission + tests before the prose).

### Summary

Resolved every open question with recorded evidence: the two live OQs (flag-free operation, headless `-p` residency) via paired nested `claude -p` stream-json probes that proved the merged named-background-`Agent` + `SendMessage(to="main")` channel works flag-free and stays resident on 2.1.181 — the bare-mode flag hint is stale; OQ-3/4/5/6 via on-disk registry inspection plus launcher/reconcile code reads. The load-bearing correction to the seed: it under-scoped the Go change. `internal/claudeteam/reconcile.go` genuinely needs none (auto-discovery already keys on `leadSessionId`/`$CLAUDE_CODE_SESSION_ID`), but `internal/dispatch/build.go` Rule 8 (`build.go:442`) blocks the merged dispatch shape (`name` present + `team_name` absent), so the deliverable is skill/contract prose PLUS a `build.go` merged-mode emission path with Go tests — recorded as the riskiest-path-first implementation item. Specified the two-mode bridge as 8 concrete before/after edits with the legacy `TeamCreate` machinery fenced under one deprecated subsection whose removal trigger references the real `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"` pin (the legacy branch's only live consumer). Finalized 5 behavior-anchored ACs, each with fixture + live coverage; AC-4's helper behavior is already covered by `reconcile_session_test.go` (reuse, not duplicate). The OQ-6 seam with m4 is a host-regime split: m4 keeps the pinned-2.1.177 legacy live lane, this task adds a merged (unpinned) live lane.

### Addendum: consistency reconciliation (post-report)

Two internal-consistency fixes after the report was first written, reconciling the top-of-doc framing and the capability matrix with the evidence gathered in `## Open-question resolutions` (no new claims, no scope change):
1. Capability matrix: flipped the three formerly-❓ rows (headless `-p` residency, flag-free, `agentType` stamping) to ✅ with the OQ-2/OQ-3 evidence, and rewrote the "Tested context" sentence so it no longer asserts those rows are untested — they were subsequently resolved (residency + flag-free are now headless-validated via the nested `claude -p` spike). No row remains ❓.
2. Scope framing: added a top-of-doc **Deliverable layers** note and reconciled the lede/`## Proposed approach`/`## Out of scope` with `## Scope finding` — the deliverable is skill/contract prose PLUS the `internal/dispatch/build.go` merged-mode emission path; only `internal/claudeteam/reconcile.go` is change-free.
