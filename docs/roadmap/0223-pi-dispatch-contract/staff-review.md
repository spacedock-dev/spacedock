# 0223 pi-dispatch-contract — sprint-wide preflight readiness review

Independent staff review of sprint `0223-pi-dispatch-contract` as a whole.
Scope: sprint-level coherence, cross-member composition, blast-radius, and
Commander cold-boot readiness — NOT per-task AC quality or mechanism choices
(the ideation gate owns those). Per-task gaps are surfaced only where they
threaten the sprint's deliverable.

Evidence base: the sprint `index.md`; the four member ideations
(`pi-ensign-skill-injection`, `pi-launcher-repo-resolution`,
`pi-dispatch-model-stamping`, `pi-back-channel-dispatch`); the dispatch core
(`skills/first-officer/references/fo-dispatch-core.md`); the Pi FO adapter
(`skills/first-officer/references/pi-first-officer-runtime.md`);
`internal/cli/pi.go`; and a live read of the installed
`pi-subagents` substrate (`~/.pi/agent/npm/node_modules/pi-subagents`) to
verify a cross-member composition claim.

## Verdict

**Gaps to close — not yet cold-boot drivable.**

The sprint is individually sound: each member's ideation is behavior-bound,
spike-backed where risky, and maps cleanly to a DoD bullet. The DoD is fully
owned and the sequencing is structurally correct. The blocker is a single
**three-way cross-member composition seam** — members 1, 2, and 4 each hold
one third of an assumption about the dispatched child's `cwd`, and no ideation
owns the wiring that makes the three compose. Under the sprint's own headline
scenario (member 2's "launch from a non-repo cwd"), member 1's cwd-keyed
symlink stops being discovered by the child, which silently re-breaks DoD
bullet (b) — and the capstone's AC-2 live drive does not pin the launch cwd,
so it would pass from a repo cwd and hide the gap. This is closeable (one
wiring instruction in the capstone + one AC pin + one quirk), not a material
redesign. Two minor gaps round out the list.

## DoD coverage

Every DoD bullet has an owning in-scope member. No unowned bullet.

| DoD bullet | Owner | Evidence |
|---|---|---|
| (a) runs on the parent FO's live model (not settings default) | member 3 | `pi-dispatch-model-stamping` AC-1; run-meta `model` == intercom-listed live model |
| (b) loads the Spacedock ensign contract (not bare `worker`) | member 1 | `pi-ensign-skill-injection` AC-2; no `skillsWarning`, ensign-contract behavior |
| (c) resolves the repo explicitly via install-record / `--plugin-dir` / `SPACEDOCK_REPO_ROOT` | member 2 | `pi-launcher-repo-resolution` AC-1/AC-2; install record ahead of cwd |
| (d) talks back to the FO over intercom mid-run | member 4 | `pi-back-channel-dispatch` AC-2; live `need_decision` → reply → resume |
| core describes all of this in runtime-neutral named capabilities | member 4 | capstone Deliverable A; AC-1 contractlint |

The four numbered "Proven by" items in the DoD map 1:1 to the four members.
Coverage is complete on paper. The risk is not coverage but **composition**
(see below): bullets (b) and (c) interact on the child-`cwd` axis, and neither
owning ideation resolves the interaction.

## Sequencing

The planned order is structurally sound, with one integration blind spot.

- **Parallel start of 1, 2, 3** — correct for *implementation*. Each is
  independent at the code/artifact level (member 1 = a repo symlink; member 2
  = `pi.go` resolution; member 3 = a prose instruction in
  `pi-first-officer-runtime.md`). No hard implementation dependency.
  The index's "no inter-member dependency" claim is true for implementation
  but **misleading for integration**: members 1 and 2 carry a composition
  assumption that only surfaces when the capstone's live drive assembles all
  three (see Cross-member composition).

- **Capstone Deliverable A (core rewrite) can start in parallel.** It is
  prose-structural (reorganizing existing contract text into named
  capabilities). It does NOT depend on 1–3. The index correctly states the
  `claude-live`/`codex-live` regression (AC-6) "can run as soon as the core
  rewrite is draft." Confirmed sound.

- **Capstone AC-2/AC-3/AC-4/AC-5 (the `pi-live` drive) genuinely requires
  members 1 AND 2 AND 3 landed first** — the index states this explicitly and
  correctly: the drive needs a dispatched ensign that runs on the parent model
  (3) with the ensign contract loaded (1) from the explicitly-resolved repo
  (2). The capstone's own `worker-identity-capture` binding cites member 3 as
  the stamped-model source; its `async-dispatch` binding presupposes a
  dispatchable ensign. This dependency is real and correctly recorded.

- **Hidden dependency the index misses:** the capstone's live drive also
  depends on a *fourth* wiring decision — the child's `cwd` — that no member
  currently owns (see Cross-member composition). This is the sequencing blind
  spot: the index treats 1–3 as the complete prerequisite set for the
  capstone's live drive, but the live drive will fail under member 2's
  headline scenario unless the capstone also wires `cwd: <repo>` into the
  dispatch call.

## Cross-member composition

This is where the sprint's coherence breaks. Three findings, one central.

### Central gap — the child-`cwd` seam between members 1, 2, and 4

Member 1's mechanism (`.pi/skills/ensign -> ../../skills/ensign`) is
**cwd-keyed by pi-subagents' own admission**. Member 1's ideation states
explicitly:

> `{cwd}/.pi/skills` is cwd-keyed, so it works only when the child's cwd is
> the repo. This is acceptable BECAUSE … the FO dispatches with cwd = project
> root, so the child inherits the repo as its cwd, and the committed
> `.pi/skills/ensign` symlink is found.

That "the child inherits the repo as its cwd" assumption is the load-bearing
precondition for member 1's entire fix. It holds **only** when the parent pi
session is launched from the repo cwd. Member 2's headline end value is the
opposite:

> `spacedock pi` launched from a cwd that is NOT the Spacedock repo resolves
> the correct skill paths after install.

Under member 2's scenario the parent pi session's cwd is the non-repo cwd.
The dispatched child's cwd defaults to the parent's cwd (`step.cwd ?? ctx.cwd`
in `subagent-runner.ts:721`, confirmed live). So under member 2's scenario,
**member 1's symlink is not in the child's discovery path** — the child scans
`{non-repo-cwd}/.pi/skills/`, which does not exist — and `skill:["ensign"]`
silently falls back to bare `worker`. DoD bullet (b) re-breaks.

I verified the seam is **closeable, not structural**. The `subagent(...)` tool
exposes a top-level `cwd` parameter to the FO caller
(`pi-subagents/src/extension/schemas.ts:69`, `cwd: Type.Optional(Type.String())`),
and the runner resolves `step.cwd ?? ctx.cwd`. So the FO *can* pass
`cwd: <resolved-repo>` on the dispatch call, forcing the child's cwd to the
repo regardless of the parent's launch cwd — which makes member 1's symlink
discoverable under member 2's non-repo launch. **But no ideation currently
owns this wiring:**

- Member 1 assumes the child cwd is already the repo (no instruction to set
  it).
- Member 2 fixes the *parent's* repo resolution and records the path, but
  says nothing about forwarding that path as the child's `cwd`.
- The capstone's `async-dispatch` binding says "Dispatch via
  `subagent(... async: true)`" and its `worker-identity-capture` captures the
  repo, but neither binding instructs the FO to pass `cwd: <repo>` on the
  spawn call.

The capstone — which owns the dispatch wiring (friction 2) — is the natural
owner. Its `async-dispatch` (or `worker-identity-capture`) binding must add:
"pass `cwd: <resolved repo root>` on every `subagent(...)` call, sourced from
the same resolution member 2 records, so the child inherits the repo as its
cwd and member 1's `.pi/skills/ensign` project-declared skill is discovered
even under a non-repo-cwd launch." Until that wiring lands, the capstone's
own AC-2 live drive is under-specified: it does not pin the launch cwd, so it
would pass from a repo cwd (exercising member 1, hiding member 2) and fail
from a non-repo cwd (exercising member 2, breaking member 1) — either way the
three-way composition is unproven.

This is the one gap that blocks cold-boot drivability.

### Secondary seam — member 3 / capstone ownership of the Pi model-space declaration

Member 3 identifies a real reuse-condition-4 hazard: Pi model strings
(`z-ai/glm-5.2`, `~openai/gpt-mini-latest`) are all outside the Claude-centric
enum (`sonnet`/`opus`/`haiku`), so the core's "captain-session fallback …
forces fresh dispatch" clause would defeat reuse on every Pi dispatch. Member
3 recommends shipping the Pi canonical-model-space declaration **in member 3**
but explicitly defers confirmation: "the capstone's ideation should confirm
this composition."

The capstone's `worker-identity-capture` schema includes "stamped model (sibling:
`pi-dispatch-model-stamping`)" and its core-rewrite note says "Condition-4's
'the host's canonical enum' stays adapter-declared" — but the capstone
**never explicitly claims or disclaims** the Pi model-space declaration. So
neither ideation firmly owns it. Risk: it ships twice (drift), or not at all
(reuse silently broken on Pi). This is a soft ownership gap, not a blocker —
close it by having the capstone's ideation explicitly confirm member 3 ships
the Pi canonical-model-space declaration (or claim it in the capstone).

### Note — core/adapter null-model contradiction window

Member 3 ships a Pi-adapter override ("stamp the parent's live model when
`dispatch build` emits `model: null`") that directly contradicts the core's
current instruction (`fo-dispatch-core.md`: "when null, OMIT the model
argument entirely"). Member 3 records this as a temporary tension the capstone
generalizes into a named `model-resolution` rule. This is sound *as a plan*,
but it means there is a window between member 3 landing and the capstone
landing where the core and the Pi adapter disagree on null handling. A
Commander reading both files during member 3's drive will see the
contradiction. Member 3's body documents it; the index's Q1–Q10 do not (see
Commander cold-boot readiness).

## Blast-radius

The path→lane mapping is correctly scoped. No gap.

- **Members 1–3 are pi-only surfaces** → `pi-live` required. Confirmed:
  - Member 1: `.pi/skills/ensign` symlink — a `.pi/` convention artifact
    read only by pi-subagents; does not touch claude/codex discovery.
  - Member 2: `internal/cli/pi.go` only (AC-3 explicitly confines the change
    and keeps `frontdoor.go`/`init.go` claude/codex paths unchanged).
  - Member 3: `skills/first-officer/references/pi-first-officer-runtime.md`
    only. The override is Pi-local; the core's "OMIT on null" still governs
    claude/codex, so no claude-live/codex-live impact from member 3 alone.
- **Capstone touches every host adapter** → `claude-live` AND `codex-live`
  AND `pi-live` required. AC-1's contractlint binds capability→tool across
  *all* adapters, so the capstone edits
  `claude-first-officer-runtime.md`, `codex-first-officer-runtime.md`, and
  the ensign runtime adapters too. AC-6 (claude-live + codex-live
  regression) is correctly scoped as the dogfood gate. Q8 states this
  correctly.

The "live lane is unrelated" proof-policy burden (Q8) is respected: each
member's lane claim is tied to a concrete diff surface, not a dispatcher
judgment. No path→lane mapping gap.

One observation (not a gap): member 1 ships a repo symlink with **no Go test**
— its proof is live (`subagents-doctor` + a probe dispatch). The sprint-wide
`go test ./...` gate is still met by members 2, 3 (Go unit tests) and 4
(structural contractlint). The symlink's existence *could* be Go-asserted
(`os.Readlink`), but its absence does not threaten the sprint deliverable
because the live probe is the behavior gate. Noting only; the ideation gate
owns AC mechanism.

## Scope

- **First-officer skill under `.pi/skills/`?** Correctly out of scope. The FO
  loads `first-officer` via the launcher's parent-side `--skill` flag (member
  2's domain); dispatched children are ensigns (workers), never FOs, and do
  not load `first-officer`. Member 1's follow-up note ("consider whether
  `first-officer` … should also be symlinked") is a genuine non-sprint
  follow-up, not a sprint-level gap. No action.
- **Over-scoped members?** None. Each member's blast radius matches its DoD
  bullet. The capstone's two-deliverable shape (A: core rewrite, B: Pi wiring)
  is appropriately bounded; frictions 7–9 are cleanly deferred with sharp
  seams in the named-capability declarations.
- **Missing member?** The child-`cwd` wiring (Cross-member composition) is
  the one missing piece. It does not require a new member — it is a binding
  instruction the capstone should absorb into its existing `async-dispatch`
  / `worker-identity-capture` wiring. No new member needed; one gap to close
  inside member 4.

## Commander cold-boot readiness

Q1–Q10 are strong and mostly complete. Q1 (async), Q2 (explicit model until
member 3), and Q3 (skill injection broken until member 1) correctly capture
the three in-flight dispatch frictions a Commander will hit, with shaping-run
evidence cited for each. Q4–Q10 (state branch, path-scoped commits,
back-channel service, completion+file-verify, live lanes, in-flight state,
sandbox) cover the operational surface. Two gaps:

1. **Missing quirk — the capstone live-drive launch cwd (the central gap's
   cold-boot face).** A Commander driving the capstone's AC-2 must decide the
   launch cwd. Q1–Q10 give no guidance. If the Commander launches from a
   repo cwd, member 2's fix is unexercised and the composition seam stays
   hidden; if from a non-repo cwd (the point of member 2), member 1's symlink
   is not discovered and AC-2 fails — with nothing in the quirks explaining
   why. The cold-boot package must either (a) bake the `cwd: <repo>` wiring
   into the capstone's dispatch instruction (preferred — closes the gap at
   the source), or (b) add a quirk directing the Commander to launch the
   capstone's live drive from a non-repo cwd AND verify ensign discovery, so
   the seam is exercised rather than hidden.

2. **Missing quirk — confirm members 1, 2, 3 have landed before the
   capstone's `pi-live` drive.** Q3 tells the Commander to verify member 1
   landed (`subagents-doctor`) before relying on skill injection. There is no
   equivalent pre-flight for the capstone's live drive: the sequencing
   section says the drive requires 1–3, but the load-bearing quirks list
   (Q1–Q10, "the Commander WILL hit them") does not include "before
   attempting the capstone's AC-2, confirm member 1 (ensign listed), member 2
   (install record present, `doctor --host pi` shows install-recorded source),
   and member 3 (a null-model dispatch stamps the parent live model) have all
   landed." A Commander who jumps to the capstone without 1–3 will hit the
   very frictions Q2/Q3 describe as workarounds, but the quirks frame those
   as pre-member-1/3 workarounds, not as "the capstone is blocked until its
   siblings land." Minor, but it costs a cold-boot session a failed drive to
   discover.

3. **Missing quirk — the core/adapter null-model contradiction window.** After
   member 3 lands and before the capstone lands, the core
   (`fo-dispatch-core.md`) says "OMIT on null" while the Pi adapter says
   "stamp on null." Member 3's body documents this; Q1–Q10 do not. A
   Commander reading both during member 3's verification will see the
   contradiction with no quirk explaining it is intentional and temporary.
   Minor.

Q1–Q10 are otherwise complete and correct for the frictions they cover.

## Gaps to close

1. **[Blocker — owner: member 4 (capstone)] Wire `cwd: <resolved repo root>`
   into the dispatch call.** The capstone's `async-dispatch` (or
   `worker-identity-capture`) binding must instruct the FO to pass
   `cwd: <repo>` on every `subagent(...)` call, sourced from member 2's
   install-recorded/explicitly-resolved repo path. This is the wiring that
   makes member 1's cwd-keyed `.pi/skills/ensign` symlink discoverable under
   member 2's non-repo-cwd launch. Verified closeable: the `subagent(...)`
   tool exposes a top-level `cwd` parameter
   (`pi-subagents/src/extension/schemas.ts:69`) and the runner resolves
   `step.cwd ?? ctx.cwd`. **Pin the capstone's AC-2 live drive to launch from
   a non-repo cwd** so the three-way composition (1 + 2 + 4) is exercised,
   not hidden behind a repo-cwd launch.

2. **[Gap — owner: member 4 (capstone), confirm against member 3] Claim the
   Pi canonical-model-space declaration.** The capstone's ideation must
   explicitly confirm that member 3 ships the reuse-condition-4 Pi
   canonical-model-space declaration (or claim it in the capstone's
   `worker-identity-capture` binding). Today neither ideation firmly owns it;
   member 3 defers confirmation to the capstone, and the capstone does not
   respond. Risk: double-ship drift or silent reuse breakage on Pi.

3. **[Gap — owner: Shaping FO, in the cold-boot package] Add the missing
   quirks.** (a) Direct the Commander to launch the capstone's live drive
   from a non-repo cwd and verify ensign discovery (or bake the `cwd: <repo>`
   wiring so the cwd choice no longer hides the seam — preferred, and
   supersedes this quirk). (b) Add a pre-flight: "confirm members 1, 2, 3
   have landed before the capstone's `pi-live` drive" (subagents-doctor lists
   ensign; `doctor --host pi` shows install-recorded source; a null-model
   dispatch stamps the parent live model). (c) Add a quirk noting the
   core/adapter null-model contradiction is intentional and temporary between
   member 3 landing and the capstone landing.

Once gap 1 is closed (the capstone wires `cwd: <repo>` and pins its AC-2
launch cwd) and gap 2 is confirmed, the sprint is cold-boot drivable. Gap 3
is cold-boot polish that should land with the `dispatch-sprint-execution.md`
package. No material redesign needed.
