# 0250 fo-behavioral-discipline — Commander dispatch (cold-boot execution package)

> **Self-contained cold-boot package.** A Commander session boots `spacedock:first-officer` on Claude, reads this file + the sprint `index.md`, and drives the **4 members** to the 0250 Definition of Done, then cuts **0.25.0**. The Shaping-FO phase is done (ideate → independent preflight staff review → fold the one blocking finding + two of three recommended trims → ideation gates approved); this is the handoff to the drive phase. Approved by the captain 2026-07-06 after the preflight staff review (NEEDS MINOR REWORK → the rework landed and was verified).

## 0. Boot

1. Resolve the launcher at the version gate: `SPACEDOCK_BIN` (set/executable) else `spacedock` on PATH. Same major.minor as this session (0.24); patch/prerelease skew is fine.
2. Workflow dir: `docs/dev` (split-root; entity state on `spacedock-state/dev`, checked out at `docs/dev/.spacedock-state`). `«state.ensure-ready»` (pull-on-boot) before any dispatch.
3. Members query: `${SPACEDOCK_BIN:-spacedock} status --workflow-dir docs/dev --where sprint=0250-fo-behavioral-discipline` → the 4 members (all `status: ideation`, gate-approved, implementation-ready):
   - `k7` / `fo-boot-engage-split` — `docs/dev/.spacedock-state/fo-boot-engage-split.md`
   - `z25` / `fo-self-evidence-bar` — `docs/dev/.spacedock-state/fo-self-evidence-bar.md`
   - `zm` / `fo-smallest-sufficient-mechanism` — `docs/dev/.spacedock-state/fo-smallest-sufficient-mechanism.md`
   - `vcm` / `fo-contract-keep-moving-posture` — `docs/dev/.spacedock-state/fo-contract-keep-moving-posture/index.md`

## 1. Gate status — do NOT re-present ideation gates

**Ideation gates are APPROVED** (captain, 2026-07-06) after the independent preflight staff review found one blocking gap and three non-blocking leanness trims:

- **Blocking (fixed, verified):** `zm`'s resident gate originally read literally over `k7`'s `«engage»` loop — a commissioned stage's already-declared dispatch would have needed a one-line justification per entity, defeating the "fire and walk away" verb the discipline exists to protect. `zm` cycle-2 rework scoped the gate to the FO's *discretionary* mechanism choice for ad-hoc tasks, explicitly exempting a commissioned stage's standing dispatch — landed in the resident blockquote itself (not just surrounding prose), plus a falsifiable AC-1 test-plan guard (engage segment must emit zero per-entity justifications). Verified by the Shaping FO by direct read, not taken on faith.
- **Trims taken:** `z25` folded its ~230-byte closing restatement into a `see present-gate` pointer (verified). `vcm` ran a wording-economy pass on S1–S4, recovering ~266 bytes (2,250→1,984) with every guard clause preserved in meaning (verified).
- **Trim deferred to implementation:** `k7`'s `scope:` bullet in the new `«engage»` block carries ~150 bytes of forward-compat rationale recoverable with tighter wording. No live ensign existed to dispatch this cheaply (k7's ideation ensign was from a prior session); take it while implementing k7's Startup edit — it's free there, not worth a dedicated dispatch.
- **Not applied, not blocking:** `z25`'s own Composition cross-check had proposed rewording its "green means ran and passed" bullet to explicitly name that a fresh validation report satisfies it. That rewording was never applied to z25's resident text. The staff review confirmed the composition holds anyway — `zm`'s resident blockquote ("re-doing a stage's verification is justified ONLY when its report shows the required check did not run green") carries the actual mechanism, and the two bullets sit adjacent in the same Working-Principles section. **Do not reorder or separate z25's and zm's bullets when composing the combined section** — the interlock is emergent from adjacency, not self-contained in either bullet alone. If you have the ideation budget to make z25 self-contained too, it's a strict improvement, but it is not required to close this gate.

Each member's body carries the gate-approved design (problem / approach / measured ACs / test plan / spike). The Commander does **NOT** re-present ideation gates; it drives `implementation → validation → done` per member.

## 2. Members, waves, and the strict-sequence constraint

Two waves — the code half and the contract-prose half have different collision profiles (per index.md's own Sequencing section):

- **Wave 0 — `k7`'s bootstrap-prompt code half, FIRST, parallel-safe.** `internal/cli/frontdoor.go`: drop `" Engage."` from `bootstrapPrompt` (line 25) and `codexBootstrapPrompt` (line 533); mirror both in the oracle constants `wantBootstrapPrompt` / `wantCodexBootstrapPrompt` (`internal/cli/safehouse_frontdoor_test.go` lines 18, 234) — the existing argv assertions already pin these as the appended last token. `pi.go`'s `piBootstrapPrompt` is untouched (no flourish there; pi parity is `7v`'s separate, not-yet-started job).
  - **`7v` coordination, checked 2026-07-06: `7v` is still `status: backlog`, not yet ideated or sprint-tagged.** Its current body hardcodes the *pre-trim* `codexBootstrapPrompt` text (including "Engage.") as pi's byte-identity target. Land this Wave 0 edit before `7v` starts, or `7v`'s implementer must re-derive its target from the live constant at implementation time rather than trusting `7v`'s own written quote — otherwise `7v` reintroduces the flourish on pi.
- **Wave 1 — the contract-prose cluster, STRICT SERIAL (never parallel — same 2-3 files, same-paragraph collision risk).** Order, per each member's own dependency reasoning:
  1. **`k7` contract-prose** — Startup step 3 (no forced workflow pick), Startup step 8 (no gate render at greet, `Use engage` hint), the deferred-load-points note, and the new `## «engage»(workflow)` function block (placed after Startup — a different section from Working Principles, so it's collision-light against 2-4 below). Take the `scope:` bullet trim here (see §1).
  2. **`z25`** — `skills/first-officer/references/first-officer-shared-core.md` Working Principles (S1, already trimmed), `skills/present-gate/SKILL.md` assembly rules (S2), `docs/dev/README.md` Proof-policy path→lane bullet, line 78 (S3, slims to reference S1). AC-2 ships a real code gate: a contractlint structural-absence test asserting no host-lane literal in S1.
  3. **`zm`** — one resident blockquote (now scope-clause-corrected) in Working Principles + a new lazy reference `references/fo-smallest-sufficient-mechanism.md` carrying the ladder, the named-busywork list, and the two-direction worked examples (AC-5's leanness split).
  4. **`vcm`** — S1–S3 in Working Principles, S4 in Clarification and Communication (now wording-trimmed). Lands last because it's the counterweight reconciled against z25's verify-bar and zm's mechanism-gate — do not resequence ahead of them.

## 3. Per-member build notes (already in the bodies; restated here as drive cues)

- **`k7`** — AC-1's value proof is a live interactive-boot drive on a *constructed heavy fixture* (≥2 discoverable workflows, ≥2 ready gates), asserting categorical signals (present-gate-at-greet count, pick-question vs. engage-hint text) captured at the controlled greet-stop — **not** a raw tool-call count (the ideation spike ruled that out on 91 real transcripts as nondeterministic and boundary-unrecoverable). AC-2 (`engage` runs the existing loop) and AC-3 (launcher flourish gone) share the same fixture; AC-3 is cheap/deterministic (`go test ./internal/cli/`), AC-1+AC-2 are the live-drive cost center.
- **`z25`** — AC-1 (the FO's own merge/triage decision holds the evidence bar) needs a new shared runtime scenario reconstructing the ezf/hf incident, per the four-step recipe in `docs/runtime-live-ci.md`, **plus the mandatory offline negative** in `shared_scenarios_negative_test.go` proving the assertion goes red on the incident end-state. AC-2 is the contractlint structural-absence test (cheap, deterministic). AC-4 is a structural dedup check between the shared-core bullet and the README's slimmed path→lane bullet.
- **`zm`** — AC-1 is a live FO drive over a constructed deterministic-local-task scenario, **plus the falsification/scope guard**: a commissioned `«engage»` segment in the same drive where the gate must stay silent through N dispatches (FAIL = any per-entity justification narrated). AC-2/3/4 are structural review at the coherence pass, not a prose-grep. AC-5 is the measured `wc -c` delta against §5 below.
- **`vcm`** — AC-1 is structural (the four strengthenings land, reconciled) verified by the cross-member review; AC-2/AC-3 are live-drive behavioral claims (post-approval advance, parallel dispatch, no-turn-end-on-async, correction-narrows-not-halts) on a scenario reconstructing the 0223 false-stop patterns as the baseline that moved the wrong way. No fresh spike needed — inherits k7's live-drive-observability proof from its spike.

## 4. Drive procedure (per member)

1. **Implementation** (worktree stage): the deliverable is **contract/skill prose + `internal/cli` Go code**, built by the dispatched ensign — do not hand-edit the shipped contract yourself. Commit in the worktree (Go) / the split-root state checkout path (entity body + stage report).
2. **Validation** (`fresh: true` — a fresh validator): MEASURE every AC's end-value per the live-drive test plans above; reproduce each "Verified by" / "Tested by" clause. No AC on this sprint is provable by prose-grep — the dev-workflow Proof policy explicitly bans it, and the entities were ideated with that in mind.
3. **Detached adversarial audit — REQUIRED before merging any of these four.** The shipped FO/ensign contract and scaffolding is one of docs/dev's four high-stakes surfaces, and every member touches it (`first-officer-shared-core.md` directly, or `present-gate/SKILL.md` / `docs/dev/README.md` which the shared core's own proof policy governs). Run on a throwaway checkout, not the impl worktree.
4. **Required CI lanes: all three host live lanes, every member, no exceptions.** All four touch the host-neutral shared core — per the dev-workflow Proof policy's path→lane mapping, that requires `claude-live` + `codex-live` + `pi-live` green before merge for each. A flake is re-run to green (serial, isolated), never skipped or left unapproved.
5. **Merge** to `main` (PR-merge), then advance to `done`. Wave 1's strict-serial constraint (§2) means these four merge one at a time in order — do not open overlapping worktrees on `first-officer-shared-core.md`.

## 5. Leanness gate — pinned baseline, not-to-exceed ceiling

Full detail in `index.md`'s "Leanness baseline" section (committed 2026-07-06). Restated for the drive:

- Baseline (`v0.24.0`, confirmed byte-identical to `origin/main` HEAD at pin time): `skills/first-officer/references/first-officer-shared-core.md` = 21,663 bytes; `skills/present-gate/SKILL.md` = 5,337 bytes.
- **Ceiling: combined resident additions to those two files MUST NOT exceed ≈5,600 bytes (+25.8%) — a not-to-exceed cap**, not a target. The three trims in §1 are margin already excluded from needing to hit it.
- Measurement command: `wc -c skills/first-officer/references/first-officer-shared-core.md skills/present-gate/SKILL.md` post-implementation; delta = post − the baseline above. Lazy-reference files (`references/fo-smallest-sufficient-mechanism.md`, etc.) are excluded by design (zm AC-5) — the whole point of the lazy split is that they don't count against boot-resident cost.
- Run this check before the final merge of Wave 1's last member (`vcm`), not just at release-cut time — catching an overshoot after all four are merged means unwinding a merged PR instead of catching it in `vcm`'s own validation.

## 6. Release cut — 0.25.0

After all 4 members are merged to `main`:
1. `go test ./...` green from the repo root.
2. **Pre-cut antipattern audit** (independent reviewer over the assembled sprint) before the tag fires; ship-blockers fixed pre-cut, non-blockers recorded for the next sprint.
3. Confirm the leanness ceiling held (§5) — if a merged member overshot, that's a ship-blocker for this audit, not a "next sprint" deferral.
4. Confirm the `7v` coordination (§2 Wave 0) — either `7v` hasn't started, or it re-derived its target from the live constant rather than its stale written quote.
5. Stamp the dev/plugin manifests → tag so the tagged commit's manifest matches its tag → publish → bump the Homebrew cask → advance the `stable` branch. Authoritative procedure: `docs/releasing.md`. *(Captain authorizes the cut.)*

## 7. Close (Shaping FO)

After the cut: fold the pre-cut audit's deferred findings into the next sprint's backlog. Two items surfaced during this sprint's shaping are explicitly **not** part of 0250 and were never filed — raise them with the captain separately, don't fold them in here: (1) four confirmed tautological tests found in `internal/ensigncycle` and `internal/status` during an unrelated audit this session (mutation-tested, not just asserted); (2) lessons from an external `testing-without-tautologies` skill worth partially adopting into docs/dev's Proof policy or the commission templates. Both are still awaiting a captain go-ahead to file.
