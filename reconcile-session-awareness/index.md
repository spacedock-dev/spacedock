---
id: 0nadpgpzer0jhrvcxeg52az2
title: reconcile auto-discovery (no --team-name) is not session-aware — newest-mtime team glob picks stale prior-session or parallel-session configs, poisoning roster-derived drift classes
status: validation
source: session-11 FO (2026-06-03) — observed bare `dispatch reconcile` resolve a two-day-old prior-session team config and report archived-entity agents as Class A; captain flagged it ("did reconcile not consider repeated or parallel sessions?")
score: "0.19"
worktree: .worktrees/spacedock-ensign-reconcile-session-awareness
started: 2026-06-03T07:09:59Z
completed:
verdict:
issue:
mod-block: merge:pr-merge
pr: "#280"
---

`spacedock dispatch reconcile` run without `--team-name` resolves the team by globbing `~/.claude/teams/*/config.json` and taking the **most-recently-modified** config containing a `spacedock:ensign` member (`internal/claudeteam/reconcile.go` `LoadReconcileTeam`, lines 43-87). The code comment calls this "a stable proxy for 'the live team in this session'" — but it is a newest-mtime heuristic, not a session-identity check, and it misfires on both repeated and parallel sessions.

## Problem

The roster-derived drift classes (A lingering, B superseded, C un-advanced-PR) are computed from the resolved team's member roster cross-referenced with entity state. When auto-discovery resolves the **wrong** team, those classes are poisoned:

- **Repeated (sequential) sessions.** Claude Code leaves `config.json` on disk after a session ends (registry-desync, anthropics/claude-code#36806). A session that does no dispatch never runs `TeamCreate`, so no fresh team dir exists — newest-mtime resolves to a **stale prior-session team**. Observed in session 11: bare reconcile picked `spacedock-v1-dev-20260601-0729-...` (two days old) and reported its archived-entity agents (`status-enumeration-and-validation-{implementation,validation}`) as Class A "lingering." The prescribed Class A action is `SendMessage shutdown_request` — meaningless against a dead session's agents.
- **Parallel (concurrent) sessions.** Two live sessions (e.g. a Claude FO + a codex peer FO, or two Claude FOs) each own a team dir. Newest-mtime **races** between them and can resolve to the *other* session's live team. Class A/B teardown then fires `shutdown_request` at a **parallel session's legitimately-live agents** — active cross-session interference, not just a false positive.

**Precision:** Class D (stale branch) and E (stale local main) are computed from git + filesystem state (worktree-behind, local-main-behind), independent of the team roster — they are session-safe and were accurate in the observed run. Only the **roster-derived classes (A/B/C)** are vulnerable.

**Unused affordance:** the config carries `leadSessionId`, loaded into `ReconcileTeamState`, with a comment saying it is "useful for narrowed discovery when --team-name is omitted" — but discovery never references it; it sorts by mtime only. The contract's event-loop step 0 also writes `[--team-name {team_name}]` as **optional**, which invites the unsafe bare invocation.

## Proposed approach

Make roster-derived drift classes provably session-scoped. Ideation picks among / combines:

1. **Degrade-to-git-only without a team identity.** When `--team-name` is omitted (or no session-scoped team resolves), emit only the git/filesystem classes (D/E) and suppress the roster classes (A/B/C), with a one-line note that roster reconciliation needs a team identity. Safe by construction: no team → no roster claims.
2. **Narrow auto-discovery by `leadSessionId`.** Match the discovered config's `leadSessionId` against the current session, using the affordance the struct already loads, instead of mtime. (Requires plumbing the current session id into the helper — confirm the env/source first; this is the riskiest unknown to exercise in ideation.)
3. **Contract update.** Reflect the safe usage in the FO event-loop step 0: roster reconciliation requires `--team-name` (the FO's own `TeamCreate` name); bare `reconcile` is git-only. Pair the prose with the code gate above so the guarantee is enforced, not just documented.

## Ideation decision

**Chosen fix: combine all three — session-id narrowing replaces mtime, degrade-to-git-only is the safe floor, contract reflects it.** Not an either/or. The newest-mtime sort is the bug and is removed entirely. The resolved loader behaves as:

- `--team-name {name}` given → load that config directly (unchanged explicit path).
- `--team-name` omitted → narrow `~/.claude/teams/*/config.json` to the config whose `leadSessionId` **equals the current session id**, among configs carrying a `spacedock:ensign` member. Exactly-one match → use it (this is the affordance the struct already loads but discovery never referenced).
- No session-matched team (zero matches, or no current session id available) → return a sentinel "no session-scoped team" result, NOT an error. The reconcile assembly then **suppresses roster classes A/B/C** and emits only the git/filesystem classes D/E, with a one-line stderr note that roster reconciliation needs a team identity. Safe by construction: no trusted roster → no roster claims, so a stale prior-session or parallel-session config can never be mistaken for "the live team."

Multiple session-id matches are not expected (one team per lead session) — on >1 match the loader degrades to git-only rather than guessing, same as zero matches. This keeps the "never trust an unverifiable roster" invariant absolute.

### Spike result (riskiest unknown — RESOLVED, no further spike needed)

The riskiest unknown was whether a current-session-id source exists to make approach #2 viable. **Exercised end-to-end during ideation:** `$CLAUDE_CODE_SESSION_ID` is present in the environment and carries the **team-lead's** session id even when read from a dispatched sub-agent's shell. Verified against the live team config: `$CLAUDE_CODE_SESSION_ID` (`aec06dd4-…86011`) is byte-identical to this session's `config.json` `leadSessionId`. So `leadSessionId == os.Getenv("CLAUDE_CODE_SESSION_ID")` is a sound session-identity check, and the FO (which runs reconcile) shares the lead session id that gets written into the config. The match is the proven mechanism the implementation's first test seeds.

### Plumbing the session id (hermetic-test constraint)

The current session id MUST be an injectable parameter, not read via `os.Getenv` inside the loader — otherwise the regression tests (which build hermetic `~/.claude/teams` fixtures with fixture UUIDs) cannot exercise the match. Signature change: `LoadReconcileTeam(home, teamName, sessionID string)`. `runReconcile` passes `os.Getenv("CLAUDE_CODE_SESSION_ID")`; tests pass a fixture UUID. The empty-`sessionID` case is the degrade-to-git-only path (handles real environments where the env var is absent). The `reconcileOpts` struct gains a `sessionID` field threaded through to the `rosterLoader` call, mirroring the existing `home` injection seam.

### Distinguishing "no team" from "git-only by design"

The loader returns a sentinel (e.g. `ReconcileTeamState{}` with `TeamName: ""` plus a typed marker, or a distinct `(state, nil)` where the assembly checks `team.TeamName == ""`) so `Reconcile` can tell "degrade to git-only" apart from a real setup failure. A real setup failure (e.g. `--team-name` given but config missing) stays exit-1 as today. The git-only degrade is exit-0 (sweep ran; roster classes intentionally absent) with the stderr note — it is not an error condition.

## Out of scope

- Changing the drift-class definitions themselves (A/B/C/D/E semantics stay).
- Deleting stale on-disk team dirs (forbidden by the runtime adapter's NEVER-delete constraint; this is about not *trusting* them, not removing them).
- The git-derived classes D/E — already session-safe.

## Acceptance criteria

Proof is Go tests over the reconcile loader with multiple team-config fixtures (the existing `internal/dispatch/reconcile_test.go` already builds hermetic `~/.claude/teams` fixtures via a tmp `home`; `reconcile_e_test.go` covers Class E). Tests inject `home` and `sessionID` and call the real `claudeteam.LoadReconcileTeam` — not a stub — so they exercise the actual loader. ACs are entity-level end-state properties of the finished fix.

**AC-1 — A stale or foreign (non-current-session) team config never poisons Class A/B/C.**
Verified by: a Go test seeding TWO `~/.claude/teams/*/config.json` fixtures — (a) a stale/foreign team whose `leadSessionId` differs from the injected current session id and whose roster + entity state would trip Class A, and (b) NO config matching the current session id — then running reconcile bare (`--team-name` empty, `sessionID` = the current-session UUID) and asserting the drift array contains **zero** Class A/B/C entries. This is the regression that locks both the repeated-session (stale prior team) and parallel-session (another live session's team) cases: in both, the only on-disk team is one whose `leadSessionId` ≠ ours, and the fix must refuse to derive a roster from it. The pre-fix newest-mtime loader fails this (it returns the foreign team and emits Class A); the fixed loader passes (degrade-to-git-only).

**AC-2 — Session-matched discovery produces correct roster classes without `--team-name`.**
Verified by: a Go test seeding a config whose `leadSessionId` EQUALS the injected current session id (plus a decoy stale config with a different `leadSessionId`), running reconcile bare, and asserting Class A/B/C are computed against the session-matched roster (e.g. an archived-entity ensign yields exactly one Class A entry for that member). This proves discovery now follows session identity, not mtime — the decoy is newer-mtime to make the distinction load-bearing (mtime would pick the decoy; session-id picks the match).

**AC-3 — Explicit `--team-name` still produces correct roster classes (no regression to the explicit path).**
Verified by: a Go test passing the team name directly and asserting Class A/B/C are computed against that roster regardless of session id (the explicit path does not consult `leadSessionId`).

**AC-4 — Git/filesystem classes (D, E) are still emitted when no team identity resolves.**
Verified by: a Go test running reconcile bare with NO session-matched team over a worktree behind `origin/next`, asserting the output still reports Class D (and Class E for a stale local main). Confirms the degrade-to-git-only path emits the session-independent classes — the sweep ran (exit 0), it just suppressed the roster classes.

**AC-5 — The git-only degrade is exit-0 with a roster-suppressed note, distinct from a setup failure.**
Verified by: a Go test asserting that bare reconcile with no session-matched team exits 0 (sweep ran), emits the one-line stderr note that roster reconciliation needs a team identity, and that an explicit `--team-name` pointing at a missing config still exits 1 (real setup failure — the degrade must not mask genuine errors).

**AC-6 — The FO event-loop step-0 contract reflects the safe usage.**
Verified by: a presence-check / oracle test over `skills/first-officer/references/claude-first-officer-runtime.md` step-0 prose confirming it states roster reconciliation (Class A/B/C) requires a team identity — explicit `--team-name {team_name}` or current-session match — and that bare reconcile without one is git-only (Class D/E). Proof at the text's own level, paired with the AC-1/AC-4 code gates that enforce it. The step-0 wording change makes `--team-name {team_name}` no longer read as a free-floating "optional" flag whose omission silently falls back to an unsafe heuristic.

## Test plan

Pinned regression-test shape (all in `internal/dispatch/reconcile_test.go`, extending the existing `reconcileFixture`):

- The fixture gains the ability to write **multiple** team configs into the tmp `home`, each with a chosen `leadSessionId` and roster, and the `reconcileOpts` it builds gains an injectable `sessionID`. This is the one structural change to the fixture; everything else reuses its existing tmp git repo + worktrees + stubbed `gh`.
- **AC-1** (foreign-team-no-poison): two configs, neither matching `sessionID`, one with a Class-A-tripping roster → assert no A/B/C in drift. Write the foreign config with a NEWER mtime than any decoy so the test also proves mtime is no longer consulted.
- **AC-2** (session-matched discovery): a matching config (+ newer-mtime decoy) → assert the expected A/B/C entries against the matched roster.
- **AC-3** (explicit `--team-name` regression): pass `teamName` directly → assert unchanged A/B/C.
- **AC-4** (git-only degrade still emits D/E): no session match, worktree behind `origin/next` (reuse `wtStale`) + stale local main → assert D and E present, A/B/C absent.
- **AC-5** (exit codes + note): assert exit 0 + stderr note for the degrade; assert exit 1 for `--team-name` at a missing path (reuse the existing setup-failure assertion).
- **AC-6**: instruction-text oracle over the step-0 prose. Cost: trivial. Same pattern as existing prose-presence checks (`hostneutrality`/`integration` style).

Cost: low — all hermetic Go unit tests over `LoadReconcileTeam` + the `Reconcile` assembly; no live workflow run needed. The signature change (`LoadReconcileTeam` gains `sessionID`, `reconcileOpts` gains `sessionID`) touches `runReconcile` and the three existing `reconcileOpts` test call sites — mechanical.

- **Spike: RESOLVED in ideation, no further spike needed.** The current-session-id source (`$CLAUDE_CODE_SESSION_ID` == config `leadSessionId`) was exercised end-to-end above; the implementation's first failing test (AC-1) seeds directly from it. The only remaining behavior to prove is the loader's match/degrade logic, which the hermetic tests cover.
- **High-stakes surface** (FO-event-loop teardown machinery that can `shutdown_request` a parallel session) → detached adversarial audit required before merge, per the entity's validation note. The audit's specific charge: confirm the degrade-to-git-only path cannot, under any fixture (missing env var, multiple matches, malformed config, empty roster), emit a roster class A/B/C derived from a non-session-matched team.

## Notes

- Surfaced session 11 while finishing the v0.19.4 release recovery; the FO ran bare reconcile with no live team and got a stale-team Class A report. Captain's question ("did reconcile not consider repeated or parallel sessions?") is the title.
- 0.19.5, test-improvement themed: the deliverable is the fix plus the regression tests that lock the session-scoping guarantee.
- The reconcile helper is named as a Phase 1 candidate in `docs/dev/_proposals/binary-simplification-roadmap.md` (FO-event-loop hot path) — this correctness fix should land before any consolidation of that surface.

## Stage Report: ideation

- DONE: Choose the fix path: degrade-to-git-only without a team identity vs leadSessionId-narrowed discovery; if the latter, exercise the current-session-id source end-to-end first (riskiest unknown).
  Chosen: combine both — leadSessionId-narrowing REPLACES newest-mtime discovery, with degrade-to-git-only as the safe floor and a step-0 contract update. Spike resolved end-to-end: `$CLAUDE_CODE_SESSION_ID` is byte-identical to the live team's `config.json` `leadSessionId` even from a dispatched sub-agent (verified `aec06dd4-…86011` == config). Recorded in the "Ideation decision" + "Spike result" sections.
- DONE: Pin the regression-test shape: a stale/foreign team config must not poison Class A/B/C, while D/E still emit without a team identity.
  Pinned AC-1..AC-6 with concrete Go-test fixtures in `internal/dispatch/reconcile_test.go` (multi-config tmp `home`, injectable `sessionID`, newer-mtime decoys to prove mtime is no longer consulted). AC-1 = foreign-team-no-poison; AC-4 = D/E still emit on degrade. See "Acceptance criteria" + "Test plan".

### Summary

Resolved both ideation checklist items. The fix removes the newest-mtime team glob entirely and narrows roster discovery by matching `config.leadSessionId` against the injected current session id (`$CLAUDE_CODE_SESSION_ID`), degrading to git-only (Class D/E) when no session-scoped team resolves — safe by construction against both repeated (stale prior-session) and parallel (concurrent-session) configs. The riskiest unknown (does a current-session-id source exist) was exercised end-to-end during ideation and confirmed: the FO's `$CLAUDE_CODE_SESSION_ID` equals the value written as `leadSessionId`, so no implementation-time spike remains. Key design constraint surfaced for implementation: the session id must be an injectable parameter (`LoadReconcileTeam(home, teamName, sessionID)`, threaded through `reconcileOpts`), not read via `os.Getenv` inside the loader, so the hermetic fixtures can drive the match. Six ACs pinned, each with a hermetic Go-test oracle; detached adversarial audit flagged for the teardown-machinery surface before merge.

## Stage Report: implementation

- DONE: Replace LoadReconcileTeam's newest-mtime team glob (internal/claudeteam/reconcile.go) with leadSessionId-matching against $CLAUDE_CODE_SESSION_ID; degrade to git-only (Class D/E) when no session-scoped team resolves.
  Mtime glob removed; signature is now `LoadReconcileTeam(home, teamName, sessionID)`. Empty/zero-match/multi-match → degrade sentinel (empty TeamName, no error). `runReconcile` passes `os.Getenv("CLAUDE_CODE_SESSION_ID")`. Commit bf382636.
- DONE: Land the regression tests in internal/dispatch/reconcile_test.go: AC-1 foreign-team-never-poisons-A/B/C (red on pre-fix), AC-3 D/E still emit on degrade — injectable sessionID + multi-config fixtures with newer-mtime decoys.
  Added reconcile_session_test.go (dispatch) covering AC-1..AC-5 over the real LoadReconcileTeam + reconcile_session_test.go (claudeteam) unit table. Proven red against a restored mtime loader: AC-1/AC-2/AC-4/AC-5 failed for the documented reasons, AC-3 (explicit path) stayed green. fixture extended via teamConfigJSONWithSession + writeTeamConfigFile (mtime-stamped decoys).
- DONE: Update the FO contract event-loop step-0 prose (AC-4): roster reconciliation needs a team identity; bare reconcile is git-only.
  claude-first-officer-runtime.md step-0 now passes `--team-name {team_name}` (no longer bracketed-optional) and states roster classes A/B/C need a team identity (explicit or current-session match) while bare reconcile is git-only D/E. Locked by skills/integration/reconcile_session_contract_test.go (region-scoped oracle over the step-0 list item), proven red before the prose edit.

### Summary

Removed the newest-mtime team glob entirely and narrowed bare-reconcile discovery to the config whose `leadSessionId` equals the injected current session id; zero/multiple/empty-session cases degrade to a git-only sentinel (empty TeamName, exit 0, one-line stderr note) so a stale prior-session or parallel-session config can never be trusted as the live roster. The session id is an injectable `LoadReconcileTeam`/`reconcileOpts` parameter (production wires `$CLAUDE_CODE_SESSION_ID`), so the new hermetic fixtures drive the match with newer-mtime decoys that prove mtime is no longer consulted. Note for validation/audit: AC-1/AC-4 suppress all three roster classes A/B/C on degrade per the entity spec; Class C is read from entity frontmatter (pr+gh) rather than the roster, but the spec's explicit instruction is to suppress A/B/C together, which the implementation follows. Full suite green (846 tests), `go vet` clean, `gofmt` clean. The high-stakes teardown surface still wants the detached adversarial audit the entity flagged before merge.

## Stage Report: validation

- DONE: Reproduce evidence for every AC (AC-1 foreign-team-no-poison, AC-2 session-matched roster, AC-3 D/E-on-degrade, AC-4 contract prose, AC-5): run the new reconcile_session tests and confirm the red-on-pre-fix (restored mtime loader) claim is reproducible.
  All 6 ACs verified against real `LoadReconcileTeam` (not a stub). On the fix: `claudeteam` session table 6/6, `dispatch` AC-1..5 5/5, `integration` AC-6 2/2 — all green. Red-on-pre-fix REPRODUCED: temporarily restored the newest-mtime loader (signature-compatible, ignored `sessionID`) and the 4 degrade/session-match unit tests + AC-1/AC-2/AC-4/AC-5 dispatch tests all went red for the documented reasons (mtime resolved the foreign/decoy team, emitted Class A+C, no degrade note); AC-3 (explicit path) stayed green as designed. Loader reverted via `git checkout`; worktree clean.
- DONE: Run go test ./... and report results; verify no roster class A/B/C can derive from a non-current-session team config.
  `go test ./...` green: 10/10 packages ok (846 tests), `go vet ./...` clean, `gofmt -l` clean. Suppression-is-total confirmed: the red-on-pre-fix AC-1 run showed the mtime loader leaked BOTH Class A and Class C from a foreign team; on the fix all of A/B/C are gated by `rosterTrusted` (reconcile.go:209-217). A throwaway spot-check (not committed) proved Class C suppression is load-bearing — the frontmatter-derived `pr-merged` (gh MERGED) entity is hidden on degrade and surfaces only on the explicit `--team-name` path.

### Summary

PASSED. The fix removes the newest-mtime team glob and resolves bare-reconcile roster discovery solely by `leadSessionId == sessionID`; every other case (zero/multi/empty-session match, foreign config) degrades to a git-only sentinel that suppresses A/B/C entirely, so no roster class can derive from a non-current-session config. All six ACs have reproducible external evidence over the real loader/assembly and the contract file — none self-referential. The red-on-pre-fix claim was independently reproduced by restoring the mtime loader (8 tests red for the documented reasons, the explicit-path tests correctly green). AC-6's prose oracle is genuinely load-bearing: an adversarial revert to the bracketed-optional framing turns it red on every clause. Full suite, vet, and gofmt are clean. This is a high-stakes teardown surface — the detached adversarial audit the entity flagged is still required before merge; validation passing is necessary but not sufficient on its own.

## Feedback Cycles

### Cycle 1 — detached adversarial audit (2026-06-03) — MATERIAL

Validation PASSED, but the detached audit found a test-strength hole the suite misses:

- **MATERIAL — the exactness of the `leadSessionId` match is not pinned.** The shipped loader (`internal/claudeteam/reconcile.go:70`) uses exact `!=` and is correct, but loosening it to `strings.HasPrefix(sessionID, lead)` or `strings.Contains(...)` breaks the core claim (a foreign/stale config resolves) while **all 20 targeted tests stay GREEN**. The realistic exploit: a foreign config with an **empty `leadSessionId`** — `HasPrefix(x, "")` / `Contains(x, "")` are always true in Go, so an empty-lead foreign team matches ANY current session. The audit reproduced the leak end-to-end (foreign `team-foreign-empty-lead` with an archived-entity ensign emitted Class A + Class C from bare reconcile — the exact session-11 stale-team bug) while `TestReconcileForeignTeamNeverPoisonsRoster` (AC-1) stayed green. The entity's own audit charge demanded the malformed/empty-config fixture; it is missing.

**Fix:** add a regression fixture that pins exactness — a foreign team whose `leadSessionId` is (a) empty and (b) a strict prefix/substring of the current session id must NOT resolve → must degrade to git-only (no roster A/B/C). The fixture must turn RED under a `HasPrefix`/`Contains` comparison and GREEN under exact `!=`. Then re-run the targeted suite (must stay green on the shipped exact code) + `go test ./...`.

(Polish, optional: pin the assembly-level `rosterTrusted` gate at `internal/dispatch/reconcile.go:191` directly via a stub `rosterLoader` returning a non-empty-roster degrade sentinel — currently it's only exercised through the loader.)

## Stage Report: implementation (cycle 2)

Addresses the Cycle 1 detached-audit MATERIAL finding: the exactness of the `leadSessionId` match was not pinned, so loosening the shipped `!=` to `strings.HasPrefix`/`strings.Contains` (notably a foreign config with an EMPTY `leadSessionId`) would leak roster classes A/B/C while all targeted tests stayed green.

- DONE: Add a regression fixture pinning exactness — a foreign team whose `leadSessionId` is empty / a strict prefix / a strict substring of the current session id must NOT resolve → degrade to git-only.
  `TestLoadReconcileTeamSessionMatchIsExact` (internal/claudeteam/reconcile_session_test.go) — table of 3 foreign near-miss configs. Verified RED under both audit-named mutants: `!strings.HasPrefix(session, lead)` fails the empty-lead + strict-prefix sub-cases (1 leaked member each); `!strings.Contains(session, lead)` fails all 3. GREEN under the shipped exact `!=`. Mutants applied/reverted via file backup; loader unchanged.
- DONE: (audit polish) Pin the assembly-level `rosterTrusted` gate directly via a stub rosterLoader.
  `TestReconcileGateSuppressesEvenWithPopulatedRoster` (internal/dispatch/reconcile_session_test.go) — stub returns the empty-`TeamName` sentinel carrying a Class-A-tripping ensign roster; asserts A/B/C suppressed while D/E still emit. Verified RED when `rosterTrusted &&` is dropped from the Class A emit (reconcile.go:209).
- DONE: Re-run targeted suite + `go test ./...`.
  Targeted (claudeteam + dispatch + integration) 185/185; full suite 851/851 (10/10 packages); `go vet ./...` clean; `gofmt -l` clean. Commit e1d5b785, test-only (the shipped loader/gate were already correct).

### Summary

Test-only cycle closing the mutation-coverage gap the audit found. The shipped exact `leadSessionId != sessionID` comparison and the `rosterTrusted` assembly gate are both now pinned by mutation-proven regression tests: the empty-`leadSessionId` exploit (the session-11 stale-team leak) and the prefix/substring near-misses all degrade to git-only under exact match and leak under any loosening. No production change. Full suite, vet, and gofmt clean.

## Stage Report: validation (cycle 2)

- DONE: Independently re-run the cycle-1 audit's adversarial edits in a throwaway worktree (loosen the leadSessionId comparison to strings.HasPrefix and to strings.Contains; a foreign config with empty leadSessionId) against TestLoadReconcileTeamSessionMatchIsExact + TestReconcileGateSuppressesEvenWithPopulatedRoster — confirm they now go RED under the mutants and GREEN on the shipped exact `!=`.
  Reproduced all three mutants on a detached throwaway worktree (`/tmp/...`, never the impl worktree), reverting each before the next. Mutant 1 `!strings.HasPrefix(session, lead)` → exactness test RED on empty-lead + strict-prefix sub-cases (foreign team resolved, 1 member leaked), GREEN on strict-substring. Mutant 2 `!strings.Contains(session, lead)` → RED on ALL 3 sub-cases (empty-lead the realistic exploit: `Contains(x,"")` always true). Mutant 3 drop `rosterTrusted &&` from the Class A emit (reconcile.go:209) → gate test RED ("leaked class A despite empty TeamName sentinel"). Shipped exact `!=` / `rosterTrusted` → both tests GREEN (baseline + after each revert). Cross-check: under Mutant 1 the cycle-1 AC-1 `TestReconcileForeignTeamNeverPoisonsRoster` STAYS green (its foreign decoy `decoy-session-zzz` is not a prefix near-miss) — confirming the new `TestLoadReconcileTeamSessionMatchIsExact` is exactly the test that closes the hole AC-1 alone missed.
- DONE: Confirm `go test ./...` green and that no cycle-1 AC regressed (the cycle-2 fix is test-only; the shipped loader/gate are unchanged).
  `go test ./...` green: 10/10 packages ok, 402 `--- PASS` total; `go vet ./...` clean. Cycle-1 ACs all still pass: dispatch AC-1..AC-5 5/5, claudeteam session table 7/7 (6 cycle-1 + the new exactness test), AC-6 contract oracle 2/2 (`TestReconcileStep0RequiresTeamIdentityForRoster` + `...DropsOptionalTeamNameFraming`, both genuine region-scoped oracles over the step-0 list item, not free-floating greps). The cycle-2 commit e1d5b785 is test-only — `git diff bf382636 HEAD -- internal/claudeteam/reconcile.go internal/dispatch/reconcile.go` is empty (shipped loader at reconcile.go:70 exact `!=`, gate at reconcile.go:191 `team.TeamName != ""` byte-identical to cycle-1's PASSED validation).

### Summary

PASSED. The cycle-1 detached-audit MATERIAL finding (leadSessionId exactness unpinned) is closed. I independently reproduced all three adversarial mutants against the new regression tests on a throwaway detached worktree: `HasPrefix` and `Contains` loosenings both turn `TestLoadReconcileTeamSessionMatchIsExact` RED (each mutant caught by the prefix vs substring sub-case the audit required), and dropping the `rosterTrusted` guard turns `TestReconcileGateSuppressesEvenWithPopulatedRoster` RED — all GREEN on the shipped exact code. The fix is test-only: the shipped loader/gate are byte-identical to cycle-1's PASSED validation, so no cycle-1 AC could regress, and all six ACs plus the new exactness test pass under `go test ./...` (402 passes, 10/10 packages, vet clean). Since this re-validation reproduced the audit's mutants directly, it IS the re-audit — no separate detached-audit dispatch is needed. Out-of-scope note: `gofmt -l` flags three files this branch does NOT touch (`internal/status/mutate.go` et al.) — they are pre-existing and equally dirty on origin/next under the now-installed go1.26.1 gofmt (consecutive-blank-line rule); pure toolchain noise, not a regression from this fix. The branch's own files are gofmt-clean.
