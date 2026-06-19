# 0223-pi-dispatch-contract — Commander cold-boot dispatch package

> **Self-contained cold-boot package.** A Commander session boots `spacedock:first-officer` on pi, reads this file + the sprint `index.md`, and drives the three members to the sprint DoD. The Shaping FO's work is done; this is the handoff to the drive phase. Approved by captain 2026-06-19 after two independent staff reviews (both gaps-closed, no material redesign).

## Boot recipe

1. **Boot the FO on pi.** Launch `pi` from the Spacedock repo (or, post member-1 landing, from any cwd after `spacedock install --host pi`). Load `spacedock:first-officer`. The FO's Startup procedure runs (contract gate → workflow discover → `docs/dev` → boot roll-up → split-root pull-on-rebase → greet).
2. **Read the sprint index.** `docs/roadmap/0223-pi-dispatch-contract/index.md` — goal, members, DoD, sequencing, Q1–Q13 quirks, both staff reviews. This package complements the index; it does not replace it. The index is the source of truth for membership/DoD; this package is the drive procedure.
3. **Confirm the state branch.** `git -C docs/dev/.spacedock-state pull --rebase origin spacedock-state/dev` (NOT `main` — see Q4). On conflict, HALT, `git rebase --abort`, surface, stop.
4. **Verify the three members are at `ideation` (gate-approved, implementation-ready).** `spacedock status --workflow-dir docs/dev --where 'sprint=0223-pi-dispatch-contract'`. All three should be `ideation` with `verdict` unset (gates approved, awaiting implementation dispatch).

## The three members (all gate-approved)

| # | slug | id | deliverable |
|---|------|----|-------------|
| 1 | `pi-install-managed-skill-placement` | `eqrcrxcyye56nfwm997bj33d` | Ship Spacedock as a pi package; `spacedock install --host pi` actually installs; both parent + child discover skills (D1–D5: root `package.json`, `.pi/extensions/spacedock.ts`, install rewrite, retire `--skill`/cwd-fallback/Stat-checks, add `spacedockPackageOK`) |
| 2 | `pi-dispatch-model-stamping` | `bdtx7bmhekpy1x12ab53d9k3` | Pi FO adapter stamps parent's live model (via `intercom list`) when `dispatch build` emits `model:null`; stage-declared models still win |
| 3 | `pi-back-channel-dispatch` (capstone) | `b23y61pgk93ph44pz506m2wy` | Harden `fo-dispatch-core.md` to 7 runtime-neutral named capabilities; wire Pi adapter bindings for frictions 1–6; ensign talkback via `contact_supervisor` |

## Sequencing

- **Parallel start: members 1 + 2.** Independent code surfaces (1: `pi.go` + `package.json` + `.pi/extensions/`; 2: `pi-first-officer-runtime.md` prose). No inter-member dependency. Dispatch both into implementation.
- **Member 3 (capstone) Deliverable A (core rewrite) can start in parallel** with 1+2 — it's prose-structural reorganization of `fo-dispatch-core.md`; the `claude-live`/`codex-live` regression (AC-6) can run as soon as the core rewrite is draft.
- **Member 3's pi-live drive (AC-2/AC-3/AC-4/AC-5) requires members 1 + 2 landed.** The drive needs a dispatched ensign that loads the ensign contract (member 1's install-managed discovery) running on the parent's live model (member 2's stamping) from a non-repo cwd (exercising member 1's no-cwd-dependency). Do NOT attempt the capstone's pi-live drive until Q12's preflight confirms 1+2 landed.
- **Merge order:** 1, 2, then 3. Each to `main` via PR-merge; state commits concurrency-safe (Q4, Q5).

## Drive procedure (per member)

For each member, the Commander runs the standard FO dispatch cycle:

1. **Advance `ideation → implementation`** (gate approved; implementation is a worktree stage). `spacedock status --workflow-dir docs/dev --set {slug} status=implementation worktree=.worktrees/spacedock-ensign-{slug} started`. Commit the state transition path-scoped, push (Q5).
2. **Create the worktree** on first dispatch to the worktree stage.
3. **Build the dispatch** via `spacedock dispatch build --workflow-dir docs/dev --entity-path docs/dev/.spacedock-state/{slug}/index.md --stage implementation --checklist-file {file} --host pi --bare-mode`. Write the checklist to a file first (one item per non-empty line, ≤3 items — the dispatch-core cap).
4. **Dispatch async** via `subagent(agent:"worker", async:true, context:"fresh", model:"z-ai/glm-5.2", cwd:<repo>, task:<build's file-pointer prompt>)`. See Q1 (async mandatory), Q2 (explicit model until member 2 lands), Q3 (skill injection broken until member 1 lands).
5. **Service mid-run `contact_supervisor` escalations** (Q6): `need_decision` → `intercom reply` within 10 min; `progress_update` → read, no reply; `interview_request` → reply with the JSON shape.
6. **On completion, verify the stage report** against the entity file (`status --read {ref} --json` → last `## Stage Report` → `Read(offset,limit)`). Never advance on a cheerful summary alone (Q7).
7. **Advance `implementation → validation`** (validation is `fresh: true` + `feedback-to: implementation` + `gate: true`). Fresh validator dispatch.
8. **At the validation gate**, present it (invoke `present-gate`); the captain decides. On captain approve → terminal ceremony (merge via PR-merge mod). On captain reject → `feedback-rejection-flow` back to implementation.
9. **Detached adversarial audit at validation** for every high-stakes surface (the shipped FO/ensign contract + host adapters + the `spacedock pi` front door) — per the dev-workflow proof policy.

## In-drive gates (captain-owned)

- **Ideation gates: APPROVED** (captain, 2026-06-19, after two staff reviews). The Commander does NOT re-present ideation gates.
- **Validation gates: captain decides.** The Commander presents each (via `present-gate`); the captain approves/rejects. Never self-approve.
- **Merge gate:** each member merges to `main` via the `pr-merge` mod (opens a code-branch PR at the merge boundary, tracks to merge). The mod-block enforcement runs at terminalization.
- **Pre-cut antipattern audit:** with all members merged to `main` and the tag NOT yet fired, dispatch an independent reviewer (staff-eng persona; not the Commander, not the implementers) over the assembled sprint to catch cross-cutting antipatterns + integration holes BEFORE the tag. Ship-blockers fixed before the cut; non-blockers recorded for the next sprint.
- **Release cut:** `go test ./...` green from the root, then `docs/releasing.md` (manifest bumps, the `vN.N.N` tag, what the tag push fires). Captain authorizes the cut.

## Q1–Q13 — Commander cold-boot quirks (load-bearing)

The Commander WILL hit these. They are not optional. The index carries the full text; this package summarizes the load-bearing ones. **Read Q1–Q13 in the index before driving.**

- **Q1 — Async dispatch is mandatory.** Every `subagent(...)` call is `async:true`. Foreground blocking cannot service a mid-run `need_decision` (10-min intercom timeout) — the FO only regains control on the foreground timeout, by which time the worker is gone (shaping-run evidence: run `b929622e` timed out foreground; run `0637e2ed`'s `need_decision` arrived after teardown). Poll with `subagent({action:"status",id})`, answer via `intercom({action:"reply"})`, `interrupt` if stalled.
- **Q2 — Explicit model on every dispatch UNTIL member 2 lands.** `dispatch build` emits `model:null` for stages with no declared model; pi-subagents resolves null to `settings.json` `defaultModel` (`~openai/gpt-mini-latest`), NOT the parent's live model. Pass `model:"z-ai/glm-5.2"` (or the current parent live model, read from the pi status bar) explicitly on every `subagent(...)` call. After member 2 lands, the adapter stamps it; this quirk is moot.
- **Q3 — Skill injection is broken UNTIL member 1 lands.** `skill:["ensign"]` emits "Skills not found: ensign" because pi-subagents children don't inherit the parent's `--skill` flags and the repo's `skills/` isn't in the child's `buildSkillPaths`. Dispatched workers run as bare `worker`s (implementation-biased). The Commander MUST compensate: put explicit stage-output discipline in the dispatch prompt (ideation = design in entity body, NOT product edits; verify stage reports carefully). AFTER member 1 lands, `skill:["ensign"]` loads via the package-root scan. **Verify member 1 landed (`subagents-doctor` — ensign as `user-package` source) before relying on skill injection.**
- **Q4 — State branch is `spacedock-state/dev`, NOT `main`.** Every state sync: `git -C docs/dev/.spacedock-state pull --rebase origin spacedock-state/dev` (push to same). Never `main` (pulling `main` rebases unrelated upstream onto the state branch and conflicts). On rebase CONFLICT (two writers editing the SAME entity's frontmatter): HALT, `git rebase --abort`, surface the conflicting path + peer commit, stop. Do NOT force-push or auto-resolve.
- **Q5 — State commits are path-scoped, then push.** The state checkout is a single non-branched git index shared with peers. A bare `git add -A`/`git commit` sweeps up a peer's staged entity. Every state commit: `git -C <state> add <entity_path> && git -C <state> commit -m "..." -- <entity_path>` then `git -C <state> push origin spacedock-state/dev`. Retry on `index.lock` contention after ~2s.
- **Q6 — Back-channel mid-run service.** `need_decision` (blocking, 10-min) → `intercom reply` promptly. `progress_update` (non-blocking) → read, acknowledge, no reply required. `interview_request` (10-min, multiple structured answers) → reply with the provided JSON shape. Until the capstone's `inbound-message-service` event-loop step ships, service these inline as they arrive (async dispatch makes this possible — Q1).
- **Q7 — Completion is the subagent return + file-verify.** When a worker returns, verify the stage report against the entity file: `status --read <ref> --json` → last `## Stage Report` heading's `offset`/`lines` → `Read(offset, limit)`. Never advance state based only on a cheerful worker summary. The stage report is the source of truth.
- **Q8 — Live lanes required for merge (path→lane mapping is the gate).** Member 1 (`internal/cli/pi.go` + `package.json` + `.pi/extensions/`) and member 2 (`pi-first-officer-runtime.md`) are pi-only → `pi-live` required. Member 3 capstone touches every host adapter (`skills/**/references/**`) → `claude-live` AND `codex-live` AND `pi-live` required (the dogfood). Per the dev-workflow proof policy, "the live lane is unrelated" is a claim the diff must substantiate, not a dispatcher judgment — and a red live lane is diagnosed by reading THIS run's failing test, not by inheriting a prior session's label.
- **Q9 — Existing members carry in-flight state — preserve, don't re-spike.** `pi-back-channel-dispatch` (member 3) has LIVE SPIKE EVIDENCE from run `0637e2ed` (the FO↔worker back-channel round-trip proven live) + archived spike `cq9kb7cdpp9y48tn8gwzmqzq` (host talkback chain PASSED). PRESERVE this evidence; do not re-spike the host talkback chain — the capstone's job is the WIRE-UP and the core hardening, not re-proving the host.
- **Q10 — Sandbox mode.** Sandbox is enabled (safehouse). `gh`, `git push`, `spacedock install` work via the safehouse allow-list. If a worker hits a sandbox denial, surface it rather than treating it as a hard failure — known environment quirk, not a contract violation.
- **Q11 — Capstone live-drive launch cwd (working-directory concern, NOT skill discovery).** The capstone's AC-2 `pi-live` drive MUST launch from a **non-repo cwd** — this exercises install-managed skill discovery (proving no cwd dependency): ensign loads via the package-root scan (`collectSettingsPackageSkillPaths`), NOT via a cwd-keyed symlink. The drive passes `cwd:<resolved repo root>` on the `subagent(...)` dispatch so the ensign's **working directory** is the repo (ensigns read entity files, run `go test`, commit to the repo) — a working-directory concern, NOT skill-discovery (skill discovery is cwd-independent post member 1, proven by `eq`'s spike). The repo path is sourced from member 1's install-registered package root (or `--plugin-dir`/`SPACEDOCK_REPO_ROOT` dev override).
- **Q12 — Preflight: confirm members 1 + 2 landed before the capstone's pi-live drive.** Before attempting the capstone's AC-2 live drive: (a) member 1 (`eq`) — `subagents-doctor` lists `ensign` as `user-package` source AND `settings.json` `packages` contains the Spacedock entry; (b) member 2 (`bdt`) — a null-model probe dispatch stamps the parent's live model (run-meta `model` == FO session live model, not `settings.json` defaultModel). A Commander who jumps to the capstone without 1+2 will hit the Q2/Q3 workarounds as blockers, not as pre-landing workarounds.
- **Q13 — Core/adapter null-model contradiction window.** Between member 2 landing and the capstone landing, the core (`fo-dispatch-core.md`) says "when null, OMIT the model argument entirely" while the Pi adapter says "stamp the parent's live model when null." This contradiction is **intentional and temporary** — member 2 documents it; the capstone generalizes it into a named `model-resolution` rule. A Commander reading both during member 2's verification will see the disagreement; it is not a bug, it is the planned transition window.

## Release-cut recipe

1. All three members merged to `main` via PR-merge.
2. **Pre-cut antipattern audit** — independent reviewer over the assembled sprint before the tag fires (not the Commander, not the implementers). Ship-blockers fixed before the cut; non-blockers recorded for the next sprint.
3. `go test ./...` green from the root.
4. Follow `docs/releasing.md` — manifest bumps, the `vN.N.N` tag, what the tag push fires. **Captain authorizes the cut.**
5. Cut from `main` (NOT `next` — `next` is the dev/edge line).

## Escalation (Commander → captain)

Escalate only on: a 3rd feedback cycle on any member, a budget blowout, an irrecoverable block, or a genuine scope fork. Otherwise the Commander drives to the DoD and presents validation gates + the pre-cut audit + the release cut for captain authorization.

## Close (Shaping FO, post-drive)

- Seed the next sprint — fold the deferred frictions 7–9 (concurrency serialization of `ask`; standing `comm-officer` on Pi; reuse-advance on Pi) + any pre-cut audit findings into a follow-up sprint.
- Light post-cut release verification (some release-machinery issues only manifest when the tag actually fires).
