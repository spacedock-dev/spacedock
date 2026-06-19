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
| 1 | `pi-install-managed-skill-placement` | `eqrcrxcyye56nfwm997bj33d` | ideation | **merged** — ship Spacedock as a pi package; `spacedock install --host pi` actually installs; both parent + child discover skills (supersedes archived `k8t` + `2m1`) |
| 2 | `pi-dispatch-model-stamping` | `bdtx7bmhekpy1x12ab53d9k3` | ideation | null model → parent live model, not settings default |
| 3 | `pi-back-channel-dispatch` | `b23y61pgk93ph44pz506m2wy` | ideation | **capstone** — declare + wire the worker↔FO back-channel over intercom; harden `fo-dispatch-core.md` to runtime-neutral named capabilities |

> **Re-carved 2026-06-19 (captain):** members `pi-ensign-skill-injection` (`k8t`) and `pi-launcher-repo-resolution` (`2m1`) ARCHIVED REJECTED — both picked clone-bound workarounds (repo symlink; cwd-fallback record) for the fact that `spacedock install --host pi` is check-only. Superseded by the merged `pi-install-managed-skill-placement` (`eq`), which ships Spacedock as a pi package (`package.json` `pi.skills` + `.pi/extensions/spacedock.ts`) so `pi install git:github.com/spacedock-dev/spacedock` makes both parent (extension `resources_discover`) and child (pi-subagents `collectSettingsPackageSkillPaths`) discover skills with no clone/cwd/symlink. Spike PASSED. The capstone's staff-review gap-1 (`cwd:<repo>` for skill discovery) is re-checked and reframed as a working-directory concern (the child-cwd seam is gone). Second staff review pending.

Members 1 (install) and 2 (model-stamping) are independent and parallelizable. Member 3 (capstone) depends on 1+2 landed for its `pi-live` drive (needs install-managed skill discovery + parent-model stamping); its core-rewrite Deliverable A can start in parallel.

## Staff preflight review #1 (run `3adf00ee`, 2026-06-19)

Independent reviewer (glm-5.2, not the FO, not the ideation ensigns) refuted the sprint as a whole. Full review at `staff-review.md`. **Verdict: Gaps to close — not yet cold-boot drivable.** Found one blocker (child-cwd seam) + two minor gaps. **All three gaps are now closed by the re-carve:**

1. ~~[Blocker] Child-cwd seam~~ — CLOSED by the re-carve. The gap assumed members 1 (`k8t` symlink) + 2 (`2m1` cwd-fallback) would compose with the capstone via `cwd:<repo>` wiring. The merged `pi-install-managed-skill-placement` removes the seam at its root (package-root scan, not cwd-keyed). The capstone's `cwd:<repo>` is re-checked (run `af9a980a`) and reframed as a working-directory concern.
2. ~~[Gap] Pi canonical-model-space declaration ownership~~ — CLOSED (capstone references member 3's declaration, doesn't re-declare; folded via run `73dadc9f`).
3. ~~[Gap] Three missing cold-boot quirks~~ — CLOSED (Q11–Q13 added to this index).

## Staff preflight review #2 (run `efff49c9`, 2026-06-19)

Independent reviewer (glm-5.2) re-refuted the re-carved 3-member sprint. Full review at `staff-review-2.md`. **Verdict: Gaps to close — not yet cold-boot drivable.** Review #1's blocker (child-cwd seam) is GENUINELY CLOSED by the re-carve (verified end-to-end); merge is clean; blast-radius correct; no material redesign. Seven follow-through gaps from the re-carve not being propagated through the docs:

1. **[Gap — owner: `eq`] `repoRoot` source post-install undefined** — `eq`'s D3 writes no repo-path record; AC-3 references an "install-record resolution" its design doesn't produce; the doctor (`pi.go:293-294,324-325`) still `Stat`s `cfg.firstOfficer/ensign` from `repoRoot`, so under a non-repo launch it reports broken skills that actually work via the package-root scan. **Closing via async fold-in re-dispatch (run in flight).**
2. **[Gap — owner: Shaping FO] Q11 contradictory + partly false** — claims ensign doesn't load without `cwd:<repo>`, disproven by `eq`'s spike and retracted by the capstone's re-check. **Closed below (Q11 rewritten).**
3. **[Gap — owner: Shaping FO] Q12 stale numbering + references `2m1`'s install-record mechanism.** **Closed below (Q12 rewritten).**
4. **[Gap — owner: Shaping FO] `## Sequencing` section stale (old 4-member layout).** **Closed below (Sequencing rewritten).**
5. **[Gap — owner: Shaping FO] DoD "Proven by" items 1–2 name archived `k8t`/`2m1`; bullet (c) says "install-recorded path".** **Closed below (DoD refreshed).**
6. **[Gap — owner: `b2` capstone] Fold-in AC-2 sub-bullet still reads "via member 1's `.pi/skills/ensign` symlink" (archived member, absent mechanism).** **Closing via async fold-in re-dispatch (run in flight).**
7. **[Cosmetic — owner: Shaping FO] Q3, Q9, Q13 stale member names/numbers.** **Closed below (Q3/Q9/Q13 refreshed).**

Once gaps 1 + 6 land (async) and 2–5, 7 are fixed inline (below), the sprint is cold-boot drivable per the review's exit criterion.

## The binding concept

The dispatch core (`skills/first-officer/references/fo-dispatch-core.md`) organizes its entire model around one declared adapter capability — the **Worker back-channel**. Today no host adapter declares it cleanly in runtime-neutral terms; the Pi adapter declares none at all and treats Pi as bare/one-shot, even though the proven pi-intercom substrate implements exactly the back-channel the core intends. The capstone hardens the core to named capabilities (`worker-back-channel`, `async-dispatch`, `inbound-message-service`, `worker-identity-capture`, `completion-signal`) and wires the Pi adapter's bindings. Members 1–3 are the prerequisite dispatch-correctness fixes the capstone's verification depends on (the right model, the loaded ensign contract, the explicitly-resolved repo).

## Sequencing

- **Parallel start:** members 1 (`eq`) + 2 (`bdt`) — independent code surfaces (`eq`: `pi.go` + `package.json` + `.pi/extensions/`; `bdt`: `pi-first-officer-runtime.md` prose). No inter-member dependency.
- **Member 3 (`b2` capstone) `pi-live` drive requires 1 + 2 landed** — AC-2 needs ensign discoverable via install-managed placement (`eq`) + parent model stamped (`bdt`). The capstone's gap-1 re-check updates the dependency from the archived `k8t`+`2m1` to `eq`.
- **Capstone Deliverable A (core rewrite) starts in parallel** — prose-structural reorganization of `fo-dispatch-core.md`; does not depend on 1+2. The `claude-live`/`codex-live` regression (AC-6) can run as soon as the core rewrite is draft.

## Definition of Done

A Pi FO dispatches an ensign that **(a)** runs on the parent FO's live model (not `settings.json` defaultModel), **(b)** loads the Spacedock ensign contract (not a bare `worker`), **(c)** resolves the Spacedock repo explicitly via the install-registered pi package (not cwd-luck), and **(d)** talks back to the FO over intercom mid-run — with the dispatch core describing all of this in runtime-neutral named capabilities each host adapter binds to concrete tools. Proven by:

1. `pi-install-managed-skill-placement` (`eq`): `spacedock install --host pi` actually installs (registers in `settings.json` `packages`); `subagents-doctor` lists ensign as `user-package` source; a probe `subagent(... skill:["ensign"])` dispatch from a NON-REPO cwd loads the ensign contract (no `skillsWarning`), child exhibits ensign-contract behavior; the `--skill` flags + cwd fallback are retired; the dev override (`pi install ./local/path` / `--plugin-dir`) discovers from a local checkout.
2. `pi-dispatch-model-stamping` (`bdt`): a dispatch with no stage-declared model runs on the parent's live model (run-meta `model` == FO session live model); the stamp is captured in worker-identity metadata; an explicit stage-declared model still wins.
3. `pi-back-channel-dispatch` (`b2`, capstone): a live `pi-live` drive — seeded-ambiguity ensign → `contact_supervisor need_decision` → FO reply within 10 min → ensign resumes → complete; both completion-signal paths (subagent return AND inbound done-message) file-verified; `claude-live`/`codex-live` regression green (the runtime-neutral core rewrite touches every host adapter — the dogfood); structural contractlint for capability-name↔adapter-binding (independent values, not prose-grep).

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

### Q2 — Explicit model on every dispatch (until member 2 lands)

**Symptom:** `spacedock dispatch build` emits `model: null` for stages with no declared model. pi-subagents resolves `null` to `settings.json`'s `defaultModel` (`~openai/gpt-mini-latest` in this env), NOT the parent session's live model. The ensign silently runs on the cheap tier regardless of the FO's model.

**Evidence from shaping:** run `0637e2ed` ran on `~openai/gpt-mini-latest` while the FO was on `z-ai/glm-5.2`.

**Quirk:** Until `pi-dispatch-model-stamping` (member 2, `bdt`) lands, the Commander MUST pass `model: <parent live model>` explicitly on every `subagent(...)` call. Read the parent's live model from the pi status bar (bottom of the TUI). After member 2 lands, the adapter stamps it; this quirk is moot.

### Q3 — Skill injection is broken until member 1 lands

**Symptom:** `skill: ["ensign"]` on a `subagent(...)` call emits `Warning: skills not found: ensign`. pi-subagents uses its OWN discovery (`discoverAvailableSkills(cwd)` → `buildSkillPaths` in `node_modules/pi-subagents/src/agents/skills.ts`): `.pi/skills`, `.agents/skills`, `~/.pi/agent/skills`, package roots, settings. The Spacedock ensign skill at `skills/ensign/SKILL.md` is in NONE of these pre-install — the main pi agent finds it via the launcher's `--skill` flags, but pi-subagents children do NOT inherit those flags. A skill-less worker runs as a bare `worker` (implementation-biased prompt: "make the edits") and silently defaults to implementation behavior even in ideation (gate, read-only).

**Evidence from shaping:** run `0637e2ed`'s meta carries `"skillsWarning": "Skills not found: ensign"`; the worker edited contract docs during ideation (contained — reverted, main clean).

**Quirk:** Until `pi-install-managed-skill-placement` (member 1, `eq`) lands, dispatched workers run as bare `worker`s. The Commander MUST compensate: put explicit stage-output discipline in the dispatch prompt (ideation = design in entity body, NOT product edits; verify stage reports carefully; expect workers to attempt edits in gate stages and review against the ensign contract manually). AFTER member 1 lands, `skill:["ensign"]` loads via the package-root scan and this quirk is moot. **The Commander should verify member 1 has landed (run `subagents-doctor` — ensign must be listed as `user-package` source) before relying on skill injection.**

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

**Quirk:** The capstone's runtime-neutral core rewrite touches every host adapter (`skills/**/references/**`) — the merge gate requires `claude-live` AND `codex-live` AND `pi-live` green (the dogfood). Members 1 (`eq`) + 2 (`bdt`) are pi-specific (`internal/cli/pi.go` + `package.json` + `.pi/extensions/`, resp. `pi-first-officer-runtime.md`) — `pi-live` required. Per the dev-workflow proof policy, "the live lane is unrelated" is a claim the diff must substantiate, not a dispatcher judgment — and a red live lane is diagnosed by reading THIS run's failing test, not by inheriting a prior session's label.

### Q9 — Existing members carry in-flight state — preserve or redispatch, don't re-spike

**Quirk:**
- `pi-back-channel-dispatch` (member 3, capstone) has a stage report at commit `0a1e2787` with **LIVE SPIKE EVIDENCE** from run `0637e2ed` — the FO↔worker back-channel round-trip was proven live (foreground dispatch auto-detached for intercom, `need_decision` reached the FO, reply sent). PRESERVE this evidence; do not re-spike the host talkback chain (it's already PASSED, archived spike `cq9kb7cdpp9y48tn8gwzmqzq`). The capstone's job is the WIRE-UP and the core hardening, not re-proving the host.
- `pi-back-channel-dispatch`'s prior run (`0637e2ed`) also made uncommitted contract-doc edits as a DRAFT — main is clean, the draft is abandoned. Redispatch clean.
- Members 1 (`eq`) + 2 (`bdt`) were ideated clean this session (async, glm-5.2, no foreground timeouts). No in-flight state to preserve; implement from the filed ideation.

### Q10 — Sandbox mode

**Quirk:** Sandbox is enabled (safehouse). `gh`, `git push`, `spacedock install` work via the safehouse allow-list. If a worker hits a sandbox denial, surface it rather than treating it as a hard failure — it is a known environment quirk, not a contract violation.

### Q11 — Capstone live-drive launch cwd (re-carved: working-directory concern, NOT skill discovery)

**Quirk:** The capstone's AC-2 `pi-live` drive MUST launch from a **non-repo cwd** — this exercises install-managed skill discovery (proving no cwd dependency): ensign loads via the package-root scan (`collectSettingsPackageSkillPaths` reading `settings.json` `packages` → `package.json` `pi.skills`), NOT via a cwd-keyed symlink. The drive passes `cwd: <resolved repo root>` on the `subagent(...)` dispatch call so the ensign's **working directory** is the repo (ensigns read entity files, run `go test`, commit to the repo) — this is a working-directory concern, NOT a skill-discovery concern (skill discovery is cwd-independent post-install, proven by `eq`'s spike). The repo path is sourced from `eq`'s install-registered package root. The capstone's gap-1 re-check section carries the authoritative stance; the prior fold-in's "required for ensign discovery" claim is SUPERSEDED.

### Q12 — Preflight: confirm members 1 + 2 landed before the capstone's pi-live drive (staff-review gap 3b)

**Quirk:** Before attempting the capstone's AC-2 live drive, confirm both siblings have landed: (a) member 1 (`eq`) — `subagents-doctor` lists `ensign` as `user-package` source AND `settings.json` `packages` contains the Spacedock entry; (b) member 2 (`bdt`) — a null-model probe dispatch stamps the parent's live model (run-meta `model` == FO session live model, not `settings.json` defaultModel). A Commander who jumps to the capstone without 1+2 will hit the very frictions Q2/Q3 describe as workarounds, but framed as pre-member-1/2 workarounds, not as "the capstone is blocked until its siblings land."

### Q13 — Core/adapter null-model contradiction window (staff-review gap 3c)

**Quirk:** Between member 2 (`bdt`) landing and the capstone (member 3) landing, the core (`fo-dispatch-core.md`) says "when null, OMIT the model argument entirely" while the Pi adapter (`pi-first-officer-runtime.md`) says "stamp the parent's live model when null." This contradiction is **intentional and temporary** — member 2 documents it; the capstone generalizes it into a named `model-resolution` rule. A Commander reading both during member 2's verification will see the disagreement; it is not a bug, it is the planned transition window.

## Sprint lifecycle checklist (owner-tagged — copy into this index to track)

**Shape — Shaping FO (this session)**
- [x] **Scope-lock** with the captain — members in/deferred ✓
- [x] **Carve** — stamp `sprint` / `sprint-readiness` on all members; write this `index.md` ✓ (re-carved to 3 members after captain review)
- [x] **Ideate** each gated member — riskiest mechanism first; all three ideations complete with spiked riskiest mechanisms ✓
- [x] **⚠️ Preflight staff review #1** (run `3adf00ee`) — verdict Gaps to close; one blocker (child-cwd seam) + two minor gaps. Folded into capstone + Q11–Q13 added. ✓
- [x] **Re-carve** (captain 2026-06-19) — archived `k8t`+`2m1` (clone-bound workarounds), filed merged `eq` (install-managed package placement, spike PASSED); capstone gap-1 re-checked (cwd:<repo> reframed as working-directory concern). ✓
- [x] **⚠️ Preflight staff review #2** (run `efff49c9`) — verdict Gaps to close; blocker GENUINELY CLOSED by re-carve; seven follow-through gaps (1 spec + 6 doc). ✓
- [ ] **Close staff-review #2 gaps** — gap 1 (`eq` repoRoot source) + gap 6 (`b2` AC-2 sub-bullet reconciliation) async; gaps 2–5, 7 fixed inline in this index ✓ (this commit)
- [ ] **Present ideation gates** — per member; never self-approve (pending gap 1 + 6 fold-ins landing)
- [ ] **Package** — write `dispatch-sprint-execution.md` (cold-boot Commander package with Q1–Q13 baked in)

**Drive — Commander (separate cold-booted session on pi)**
- [ ] Implementation → validation → done per member; detached adversarial audit at validation for every high-stakes surface (the shipped FO/ensign contract + host adapters + the `spacedock pi` front door)
- [ ] Merge each to `main` (PR-merge); state commits concurrency-safe (Q4, Q5)
- [ ] ⚠️ Pre-cut antipattern audit — independent reviewer over the assembled sprint before the tag fires
- [ ] Cut the release — `go test ./...` green from the root, then `docs/releasing.md`

**Close — Shaping FO**
- [ ] Seed the next sprint — fold the deferred frictions 7–9 + any pre-cut findings into a follow-up
