# 0223 — pi-dispatch-contract (0.22.3)

> **Carved 2026-06-19 (captain + Shaping FO session on z-ai/glm-5.2).** First Pi-specific sprint. The Pi runtime support was previously delivered (7 archived PASSED tasks — `pi-runtime-support`, `pi-runtime-frontdoor`, `pi-stage-dispatch-uses-build-artifact`, `pi-stage-dispatch-fresh-context`, `pi-live-ci-runtime-scenarios`, `pi-intercom-runtime-capability-probe` spike, `pi-runtime-readiness-intercom-prereqs`); the archived spike `cq9kb7cdpp9y48tn8gwzmqzq` (PR #301 closed-as-spike) proved the host talkback chain. This sprint **wires** the proven host capability into the Spacedock dispatch contract and hardens the contract to runtime-neutral named capabilities — it does not implement host behavior.

## Theme

Pi dispatch contract correctness — make a Pi FO dispatch behave the way the dispatch core's back-channel model already intends, and harden the core to runtime-neutral named capabilities so each host adapter implements them by name rather than the core speaking in host-specific tools.

## Membership is the query, never a hard-coded list

```bash
spacedock status --workflow-dir docs/dev --where sprint=0223-pi-dispatch-contract
spacedock status --workflow-dir docs/dev --where sprint=0223-pi-dispatch-contract --where 'sprint-readiness != defer'
```

## Members

| # | member | id | status | layer |
|---|--------|----|--------|-------|
| 1 | `pi-ensign-skill-injection` | `k8tbnmcbyqc5kkhj0m9vewq4` | ideation | child-side skill discovery (Layer 1) |
| 2 | `pi-launcher-repo-resolution` | `2m1cgn22ygmwtxe43z2hx7xw` | ideation | launcher repo resolution, no cwd fallback (Layer 2) |
| 3 | `pi-dispatch-model-stamping` | `bdtx7bmhekpy1x12ab53d9k3` | ideation | null model → parent live model, not settings default |
| 4 | `pi-back-channel-dispatch` | `b23y61pgk93ph44pz506m2wy` | ideation | **capstone** — declare + wire the worker↔FO back-channel over intercom; harden `fo-dispatch-core.md` to runtime-neutral named capabilities |

Members 1–3 are independent, small, and unblock correct dispatch — parallelizable. Member 4 (capstone) depends on correct dispatch being *possible*, so it rides after/alongside; its friction 4 (`worker-identity-capture`) is where member 3's model-stamping lands structurally.

## The binding concept

The dispatch core (`skills/first-officer/references/fo-dispatch-core.md`) organizes its entire model around one declared adapter capability — the **Worker back-channel**. Today no host adapter declares it cleanly in runtime-neutral terms; the Pi adapter declares none at all and treats Pi as bare/one-shot, even though the proven pi-intercom substrate implements exactly the back-channel the core intends. The capstone hardens the core to named capabilities (`worker-back-channel`, `async-dispatch`, `inbound-message-service`, `worker-identity-capture`, `completion-signal`) and wires the Pi adapter's bindings. Members 1–3 are the prerequisite dispatch-correctness fixes the capstone's verification depends on (the right model, the loaded ensign contract, the explicitly-resolved repo).

## Sequencing

- **Parallel start:** members 1, 2, 3 — independent, all in ideation, no inter-member dependency.
- **After 1–3 land (or alongside once their designs are locked):** member 4 capstone. The capstone's `pi-live` drive (AC-2/AC-3/AC-4) requires a dispatched ensign that runs on the parent model (member 3) with the ensign contract loaded (member 1) from the explicitly-resolved repo (member 2). The capstone's `claude-live`/`codex-live` regression is independent of 1–3 and can run as soon as the core rewrite is draft.

## Definition of Done

A Pi FO dispatches an ensign that **(a)** runs on the parent FO's live model (not `settings.json` defaultModel), **(b)** loads the Spacedock ensign contract (not a bare `worker`), **(c)** resolves the Spacedock repo explicitly via install-recorded path or `--plugin-dir`/`SPACEDOCK_REPO_ROOT` (not cwd-luck), and **(d)** talks back to the FO over intercom mid-run — with the dispatch core describing all of this in runtime-neutral named capabilities each host adapter binds to concrete tools. Proven by:

1. `pi-ensign-skill-injection`: `subagents-doctor` lists ensign; a probe `subagent(... skill:["ensign"])` dispatch loads the ensign contract (no `skillsWarning`), child exhibits ensign-contract behavior (ideation probe = design in entity body, not product edits), registration is project-scoped and survives launch from a cwd other than the repo.
2. `pi-launcher-repo-resolution`: `spacedock pi` launched from a non-repo cwd resolves the correct skill paths; the repo path is install-recorded not cwd-derived; claude/codex parity is explicit (in-scope for all hosts or pi-only, decided and tested).
3. `pi-dispatch-model-stamping`: a dispatch with no stage-declared model runs on the parent's live model (run-meta `model` == FO session live model); the stamp is captured in worker-identity metadata; an explicit stage-declared model still wins.
4. `pi-back-channel-dispatch` (capstone): a live `pi-live` drive — seeded-ambiguity ensign → `contact_supervisor need_decision` → FO reply within 10 min → ensign resumes → complete; both completion-signal paths (subagent return AND inbound done-message) file-verified; `claude-live`/`codex-live` regression green (the runtime-neutral core rewrite touches every host adapter — the dogfood); structural contractlint for capability-name↔adapter-binding (independent values, not prose-grep).

## Out of scope (deferred to a follow-up sprint — depends on this one landing)

Frictions 7–9 from `pi-back-channel-dispatch`:
- 7. Concurrency serialization of `ask` (intercom allows one pending ask per session vs `concurrency: 3`).
- 8. Standing `comm-officer` on Pi via a long-lived intercom session.
- 9. Reuse-advance on Pi (sending the next-stage assignment to a kept-alive worker).

These depend on the capstone's frictions 1–6 landing first. Seed the follow-up sprint from the capstone's archived body when it lands.

## Commander cold-boot package — pi-specific quirks

The Commander drives this sprint on Pi (boots `spacedock:first-officer`, creates its own team). The cold-boot package is `dispatch-sprint-execution.md` (below). These pi-specific quirks are load-bearing — the Commander WILL hit them; record and plan for each, do not improvise.

### Q1 — Async dispatch is mandatory, not optional

**Symptom:** A foreground `subagent(...)` call blocks the FO; a worker's mid-run `contact_supervisor need_decision` (10-min intercom timeout) arrives while the FO is blocked and times out before the FO can reply. The FO only regains control when the foreground run times out (1800s default), by which time the worker is torn down and the reply cannot deliver ("Session not found").

**Evidence from shaping:** run `b929622e` timed out foreground at 1800s; run `0637e2ed`'s `need_decision` arrived after teardown, reply undeliverable.

**Quirk:** Every dispatch MUST be `subagent(... async: true)`. Poll with `subagent({action:"status", id})`, answer `need_decision` via `intercom({action:"reply", message})` within 10 min, `interrupt` if stalled. Foreground blocking is the wrong shape for a back-channel FO — it is literally friction 2 the capstone ships.

### Q2 — Explicit model on every dispatch (until member 3 lands)

**Symptom:** `spacedock dispatch build` emits `model: null` for stages with no declared model. pi-subagents resolves `null` to `settings.json`'s `defaultModel` (`~openai/gpt-mini-latest` in this env), NOT the parent session's live model. The ensign silently runs on the cheap tier regardless of the FO's model.

**Evidence from shaping:** run `0637e2ed` ran on `~openai/gpt-mini-latest` while the FO was on `z-ai/glm-5.2`.

**Quirk:** Until `pi-dispatch-model-stamping` lands, the Commander MUST pass `model: <parent live model>` explicitly on every `subagent(...)` call. Read the parent's live model from the pi status bar (bottom of the TUI). After member 3 lands, the adapter stamps it; this quirk is moot.

### Q3 — Skill injection is broken until member 1 lands

**Symptom:** `skill: ["ensign"]` on a `subagent(...)` call emits `Warning: skills not found: ensign`. pi-subagents uses its OWN discovery (`discoverAvailableSkills(cwd)` → `buildSkillPaths` in `node_modules/pi-subagents/src/agents/skills.ts`): `.pi/skills`, `.agents/skills`, `~/.pi/agent/skills`, package roots, settings. The Spacedock ensign skill at `skills/ensign/SKILL.md` is in NONE — the main pi agent finds it via the launcher's `--skill` flags, but pi-subagents children do NOT inherit those flags. A skill-less worker runs as a bare `worker` (implementation-biased prompt: "make the edits") and silently defaults to implementation behavior even in ideation (gate, read-only).

**Evidence from shaping:** run `0637e2ed`'s meta carries `"skillsWarning": "Skills not found: ensign"`; the worker edited contract docs during ideation (contained — reverted, main clean).

**Quirk:** Until `pi-ensign-skill-injection` lands, dispatched workers run as bare `worker`s. The Commander MUST compensate: put explicit stage-output discipline in the dispatch prompt (ideation = design in entity body, NOT product edits; verify stage reports carefully; expect workers to attempt edits in gate stages and review against the ensign contract manually). AFTER member 1 lands, `skill:["ensign"]` loads and this quirk is moot. **The Commander should verify member 1 has landed (run `subagents-doctor` — ensign must be listed) before relying on skill injection.**

### Q4 — State branch is `spacedock-state/dev`, NOT `main`

**Symptom:** Pulling `origin main` into the state checkout rebases the unrelated upstream main onto the state branch and conflicts ("README.md had different types on each side" — distinct types). The state checkout tracks `origin/spacedock-state/dev`.

**Evidence from shaping:** the initial boot pull-on-boot hit this conflict by operator error.

**Quirk:** Every state sync is `git -C docs/dev/.spacedock-state pull --rebase origin spacedock-state/dev` (push to the same). Never `main`. On rebase CONFLICT (two writers editing the SAME entity's frontmatter concurrently): HALT, `git rebase --abort`, surface the conflicting path + peer commit, stop — do NOT force-push or auto-resolve.

### Q5 — State commits are path-scoped, then push

**Quirk:** The state checkout is a single non-branched git index shared with peers. A bare `git add -A` / `git commit` sweeps up a peer's staged entity and cross-attributes. Every state commit: `git -C <state> add <entity_path> && git -C <state> commit -m "..." -- <entity_path>` then `git -C <state> push origin spacedock-state/dev`. Retry on `index.lock` contention after ~2s.

### Q6 — Back-channel mid-run service (once the capstone's event-loop step ships)

**Quirk:** When a worker sends `contact_supervisor` with `reason: need_decision` (blocking, 10-min timeout) the Commander replies via `intercom({action:"reply", message:"..."})` promptly. With `reason: progress_update` (non-blocking) the Commander reads and acknowledges, no reply required. With `reason: interview_request` (10-min timeout, multiple structured answers) reply with the provided JSON shape. Until the capstone's `inbound-message-service` event-loop step ships, the Commander services these inline as they arrive (async dispatch makes this possible — Q1).

### Q7 — Completion is the subagent return + file-verify

**Quirk:** When a worker returns, verify the stage report against the entity file: `spacedock status --read <ref> --json` → take the last `## Stage Report` heading's `offset`/`lines` → `Read(offset, limit)` that range. Never advance state based only on a cheerful worker summary. The stage report is the source of truth.

### Q8 — Live lanes required for merge (the path→lane mapping is the gate)

**Quirk:** The capstone's runtime-neutral core rewrite touches every host adapter (`skills/**/references/**`) — the merge gate requires `claude-live` AND `codex-live` AND `pi-live` green (the dogfood). Members 1–3 are pi-specific (`pi-first-officer-runtime.md`, `pi-ensign-runtime.md`, `internal/cli/pi.go`) — `pi-live` required. Per the dev-workflow proof policy, "the live lane is unrelated" is a claim the diff must substantiate, not a dispatcher judgment — and a red live lane is diagnosed by reading THIS run's failing test, not by inheriting a prior session's label.

### Q9 — Existing members carry in-flight state — preserve or redispatch, don't re-spike

**Quirk:**
- `pi-back-channel-dispatch` (member 4) has a stage report at commit `0a1e2787` with **LIVE SPIKE EVIDENCE** from run `0637e2ed` — the FO↔worker back-channel round-trip was proven live (foreground dispatch auto-detached for intercom, `need_decision` reached the FO, reply sent). PRESERVE this evidence; do not re-spike the host talkback chain (it's already PASSED, archived spike `cq9kb7cdpp9y48tn8gwzmqzq`). The capstone's job is the WIRE-UP and the core hardening, not re-proving the host.
- `pi-ensign-skill-injection` (member 1) had a first ideation run (`b929622e`) that TIMED OUT foreground with no commit — redispatch clean (async, per Q1).
- `pi-back-channel-dispatch`'s prior run (`0637e2ed`) also made uncommitted contract-doc edits as a DRAFT — main is clean, the draft is abandoned. Redispatch clean.

### Q10 — Sandbox mode

**Quirk:** Sandbox is enabled (safehouse). `gh`, `git push`, `spacedock install` work via the safehouse allow-list. If a worker hits a sandbox denial, surface it rather than treating it as a hard failure — it is a known environment quirk, not a contract violation.

## Sprint lifecycle checklist (owner-tagged — copy into this index to track)

**Shape — Shaping FO (this session)**
- [x] **Scope-lock** with the captain — members in/deferred ✓
- [x] **Carve** — stamp `sprint` / `sprint-readiness` on all four; write this `index.md` ✓
- [ ] **Ideate** each gated member — riskiest mechanism first; check existing ideation state (members 1, 4 have partial ideation; members 2, 3 fresh)
- [ ] **⚠️ Preflight staff review (sprint-wide)** — dispatch ONE independent reviewer (not the FO, not the ideation ensigns) over the sprint as a whole → `staff-review.md`
- [ ] **Present ideation gates** — per member; never self-approve
- [ ] **Package** — write `dispatch-sprint-execution.md` (the cold-boot Commander package with Q1–Q10 baked in)

**Drive — Commander (separate cold-booted session on pi)**
- [ ] Implementation → validation → done per member; detached adversarial audit at validation for every high-stakes surface (the shipped FO/ensign contract + host adapters + the `spacedock pi` front door)
- [ ] Merge each to `main` (PR-merge); state commits concurrency-safe (Q4, Q5)
- [ ] ⚠️ Pre-cut antipattern audit — independent reviewer over the assembled sprint before the tag fires
- [ ] Cut the release — `go test ./...` green from the root, then `docs/releasing.md`

**Close — Shaping FO**
- [ ] Seed the next sprint — fold the deferred frictions 7–9 + any pre-cut findings into a follow-up
