# Sprint 0200-flip — preflight staff review

**Method:** independent adversarial panel (not the FO, not the ideation ensigns) — one reviewer per entity (`k6d`, `pj`) + a cross-coherence reviewer over the pair + `nzb`, then a skeptical synthesis. Lenses: design soundness, proof-gaps, tautological/prose-grep ACs, over-engineering, blast-radius. Each claim verified against the real code/CI, not the entities' self-reports.

**Verdict: NOT-READY** — five Material defects must be folded before the ideation gates lock. None is a bogus-release risk (the e2e-gate is fail-safe on publish), but each is a concrete design / proof / blast-radius defect with a named fix. Once folded, six polish items can ride implementation.

## Material findings (fold before the gate)

**M1 — `k6d` AC-b/M2: phantom two-channel stamp + tautological guard.**
AC-b specifies the release.yml version-stamp should "derive its branch from the release channel (stable→`main`, edge→`next`)", but `release.yml` triggers ONLY on `push: tags: ['v*']` (release.yml:7-10) — a tag push carries no channel signal, and there is no edge v*-tag cut scheme anywhere. So post-flip exactly one thing reaches the stamp step: a stable `v*` tag → it should stamp `main`. The "edge→next" arm is phantom. The coupled proof-defect: with no channel input, AC-b's guard test can only parse the value back from the release.yml the implementer authored — the banned tautological shape.
*Fix:* reframe AC-b/M2 to single-target (`next→main`); delete the phantom edge arm; make the guard non-tautological by cross-checking three INDEPENDENT surfaces that must agree on `main` — `release.yml` stamp target vs `.goreleaser.yaml` devBranch (stable) vs `.claude-plugin/marketplace.json` `source.ref`. Reflect the same wording back to pj checklist:128. → **routed to k6d ensign**

**M2 — `nzb` (cross): adding `needs: e2e-gate` breaks 14 journey-workflow guard tests; AC-4 "suite stays green" is empirically false.**
`internal/release/journey_workflow_test.go` anchors five functions on the byte-literal `"  goreleaser:\n    runs-on: macos-latest"` (:187/:231/:238/:274) and `TestReleaseWorkflowJobGraphMatchesGitHubActions:424` asserts goreleaser has NO `needs`. Inserting `needs: e2e-gate` makes the anchor vanish → the reviewer patched a scratch copy and `go test ./internal/release/` went green → **0 passed / 14 failed**. An implementer following nzb AC-4 literally ships a RED package believing the AC passed.
*Fix:* broaden nzb AC-4 from "the separation suite stays green" to "`go test ./internal/release/` green over the WHOLE package after the change", and add to nzb's blast-radius an explicit implementation co-edit: update the four `goreleaser:`/`runs-on` literal anchors + the no-`needs` assertion to the new header carrying `needs: e2e-gate`. Belongs to nzb (the entity adding the edge). → **routed to a fresh nzb ensign**

**M3 — `pj`: the freeze window must start at step 3, not step 4 (`workflow_dispatch` binds headSha to the branch tip).**
`runtime-live-e2e.yml`'s `workflow_dispatch` has no per-SHA input (inputs are claude_version/codex_version/effort) — it executes on `next`'s tip AT DISPATCH TIME, not the recorded `$PREPARED`. Dev-stays-on-`next` is explicit (AC-8) and `next` has provably moved since ideation (d13f355e → d6b27eb9). A push landing between step 3 (record `$PREPARED`) and step 4 (dispatch) burns real API spend across 4 environments on a run whose headSha ≠ `$PREPARED`. Not a bogus-release risk (the gate fail-safe-blocks at the tag), but a wasted-spend / cut-day-surprise hazard.
*Fix:* extend the tag-the-green-tip invariant to "`next` frozen from step 3 through step 7"; add a precondition (branch lock or captain-announced dev quiesce) and a re-verify (green run's headSha == `$PREPARED`) before step 5. Runbook discipline, no mechanism change. → **routed to pj ensign**

**M4 — cross: stable calendar-bump-on-`main` has two non-identical owner stories + an "implementation call" punt.**
k6d says stable bump is NOT k6d's — it "rides docs/releasing.md:38 release-prep"; pj step 2b folds the bump as a commit onto the prepared `next` line; pj AC-7 says "an implementation call." No entity owns an ongoing automated post-flip stable calendar-bump (next-publish.yml bumps + pushes `next` only; the release.yml stamp step leaves the calendar key untouched). An AC left "an implementation call" is not firmed.
*Fix:* pin AC-7 to ONE mechanism. → **routed to pj ensign (resolved jointly with M5)**

**M5 — cross: transient cross-ref exposure — `$PREPARED` on `next` carries `ref:main` while `main` is still the OLD tip.**
Step 2 folds the marketplace ref flip (`next→main`) + calendar-key bump onto the prepared line while it lives on `next`, held there through steps 4-6. In that window `next` carries `marketplace.json{ref:main, bumped key}` but `origin/main` is still the OLD pre-flip tip (8c069d95). A `@next` consumer running `plugin update` (the bumped key is exactly what fires it) reads `ref:main` and resolves the plugin payload from the OLD main — lacking the current plugin surface.
*Fix (resolves M4 too):* hold the calendar-key bump OUT of the pre-flip line and apply it on `main` only AFTER the flip (step 6+); pin AC-7 to that single path. Alternatively record the window as accepted with rationale. → **routed to pj ensign**

## Polish (may ride implementation)

- **P1 (cross):** reflect M1's single-target reframe into pj checklist:128 so k6d (Commander) and pj (FO) keep one source of truth on the stamp target.
- **P2 (k6d):** label AC-c's `next-publish.yml`-untouched check as a non-regression assertion, not functional proof of edge flow; cross-ref that post-flip edge behavior is verified under pj's calendar AC. Optionally drop the unexercised "edge keeps flowing" assertion from k6d's body.
- **P3 (pj):** AC-1(b)'s strict-semver `release.yml` guard silently requires an unowned release.yml edit. Decide: (i) the non-`v*` archive ref alone suffices → drop the guard from AC-1's mandatory Verified-by (optional follow-up), or (ii) name pj as the owner of the new guard step.
- **P4 (pj):** AC-4 doc-vs-machinery is a manual human read with no failing-on-violation artifact. Optional hardening: a small `internal/release` assertion that the post-flip stamp target + `.goreleaser.yaml` devBranch are `main`, anchoring the doc check to a machine-checked fact.
- **P5 (cross):** M2 version-stamp ownership is cleanly k6d (verified, not double-owned). Normalize the line citation to `release.yml:161,172` (the code) across both entities — k6d/pj cite `:156/:159` (the step name / a comment).
- **P6 (cross):** `journey-ledger`'s `--branch next` filter (release.yml:56) is correct for the 0.20.0 cut (its e2e runs on next at `$PREPARED`) but is orphaned for ongoing post-flip stable cuts; non-fatal (best-effort, exit 0 on no-run). Record as fix-later; no pre-gate owner needed.

## Resolution

All five Material findings folded and FO-verified by reading the revised sections (not just the ensign reports):

- **M1 (k6d, `94724ea4`)** — AC-b/M2 reframed to single-target `next→main`; the phantom edge arm deleted; the guard replaced with a **tri-surface agreement check** (release.yml stamp ↔ `.goreleaser.yaml` stable devBranch ↔ marketplace.json ref) across three independently-authored files — non-tautological, with honest pre-flip (binary-side pair green) / at-flip (full `==main`) sequencing. P2 (AC-c relabeled non-regression) + P5 (citations `:161/:172`) folded.
- **M3 (pj, `fc9cd169`)** — the tag-the-green-tip invariant now freezes `next` from step 3 (record `$PREPARED`) through step 7 (tag); runbook adds the freeze + a post-step-4 `headSha == $PREPARED` re-verify (re-record+re-dispatch on mismatch). Closes the workflow_dispatch-binds-to-tip hazard.
- **M4 + M5 (pj, `fc9cd169`)** — one fix: the calendar-bump is held OUT of the pre-flip fold and lands on `main` post-flip (runbook step 8); AC-7 pinned to that single mechanism. Removes the transient stale-payload window (`next` serving `ref:main` while `main` is the old tip) and the dual-owner punt. P1/P3/P4/P5 folded (single-target wording, strict-semver guard demoted to optional follow-up, AC-4 machine-anchor option, citations).
- **M2 (nzb, `fd9960c1`)** — AC-4 reframed to "`go test ./internal/release/` green over the WHOLE package"; a new "Implementation co-edit" section names the four byte-literal `goreleaser:`/`runs-on` anchors (:187/:231/:238/:274) + the no-`needs` assertion (:424) as required co-edits owned by nzb. The ensign re-reproduced the 14-failure break empirically and restored clean.

**Status: READY** — all Material findings closed; the six polish items are folded or recorded. Ideation gates presented to the captain.
