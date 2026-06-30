---
title: «fn» binding refinements — promote reuse-condition-4 / «gate.ac-cross-check» / «halt.rebase-conflict» «fn»s; fix →prose + spawn/shutdown arrows; consolidate Dispatch/Merge pointers
status: ideation
score: 0.35
group: cleanup
issue:
id: z4qnbcsv0ffty4njszg214wd
sprint: 0240-lean-contract
started: 2026-06-30T16:20:13Z
---

Post-v0.23.0 refinement from the 0230 «fn»-binding effectiveness audit (5-lens, verdict "mostly-effective-with-fixes"). The binding program is sound and ship-safe — every load-bearing guarantee is machine-pinned and the arrow trichotomy is the standout legibility win — so these are quality refinements, NOT correctness fixes, and were deliberately deferred past the v0.23.0 cut. (The audit's scariest finding, a claimed worker.spawn/shutdown over-collapse in the M2 trim, was VERIFIED a false alarm: worker.shutdown content is preserved and worker.spawn is a legitimate DRY relocation into the heading + the promoted «dispatch.build» «fn».)

## The six refinements

1. **Promote reuse-condition-4 to a `«reuse.model-match»` «fn»** (fo-dispatch-core.md). It is the heaviest deterministic comparator in the dispatch core (four branches: match / null-skip / captain-session-fallback-never-matches / mismatch-emit-diagnostic) yet rendered as flat prose while thinner per-host probes ARE «fn»s — the prose/binding boundary is drawn on the wrong side. Give it guard (null-skip) / effect (resolve + compare in the «worker-identity» canonical model space) / block (a captain-session fallback value never matches → forced one-time fresh re-stamp) / done-when (the verbatim `does not match next stage effective_model` diagnostic), and MOVE the existing verbatim-diagnostic catastrophe-backstop pin onto the new «fn» body (keep it discriminating + non-vacuous).

2. **Disambiguate `→ prose`** for `«dispatch.next-action»` (fo-dispatch-core.md). runtime-support.md:21 defines `→ prose` as "judgment-owned, no binary backing," but the next-action body is a deterministic 3-step scheduler awaiting a binary (descoped to roadmap 0222), not judgment. Either retag it (a "mechanism, binary pending" status distinct from `→ prose`) or add one line at the definition stating it is hand-followed-but-deterministic, so a future implementer neither reads it as non-codifiable nor probes for a not-yet-shipped `spacedock dispatch next-action`.

3. **Add a `→ runtime-binding` pointer line to `«worker.spawn»` and `«worker.shutdown»`** (fo-dispatch-core.md). These two highest-stakes capabilities carry NO `→` line in either tree — a cold FO could read the host-neutral body as "nothing to do here" with no in-file signal that it is runtime-bound or which adapter section binds it (the Claude adapter realizes them as prose headings `## Spawn Call (Agent)` with no guillemet match). Add a kind-only `→ runtime-binding: bound in the host adapter's ## Runtime implementation` line; requires loosening TestDispatchCoreDefinesWorkerLifecycleCapabilities to permit a kind-only arrow while still banning a per-host `**Host:**` token in the core body.

4. **Promote the AC coverage cross-check (+ end-value re-anchor) to a `«gate.ac-cross-check»` «fn»** (first-officer-shared-core.md, `## Completion and Gates`). The cross-check is flat prose while `«gate.assemble-verdict»` right beside it is a «fn» — the prose/binding boundary is on the wrong side (the #1 defect again). Give it guard (a value-measuring AC is paired to the mechanism-only AC) / effect (scan each `**AC-N**` for evidence; resolve the mechanism→value "serves" pairing) / block (a mechanism-only AC whose served value-AC regressed → REJECT, not pass) / done-when (every AC has evidence or is named). Folding it into guard/block — instead of today's inference-prose — closes two of bm's own validation edge-findings: the "only bites when a paired value-AC exists" precondition becomes the **guard**, and the "mechanism→value pairing left to FO inference" becomes the **effect**.

   **bm LANDED (#441, 2026-06-30) — unblocked.** This RESTRUCTURES bm's EDIT A (the re-anchor clause), now merged on `main`. Cross-ref: bm = `bmt9h66tg1s3eda1e1vxmzj` gate-on-end-value; its validator's adversarial-audit edge findings are the motivation for the guard/block structure here. Same net-neutral value-guardrail as AC-4 (= bm's AC-3).

5. **Promote the rebase-conflict halt to a `«halt.rebase-conflict»(paths)` «fn»** (first-officer-shared-core.md). The contract refers to "the rebase-conflict halt" BY NAME at 2 FO sites (shared-core ~162, ~177; + claude-fo-dispatch ~162) but DEFINES it 0× in the FO context — the only full definition lives in the *ensign* core, a different load tier. Define it once in the boot-resident core as a `«fn»`: block (`git rebase --abort`; surface the conflicting entity path(s) + peer commit; stop; never `--force`/`-X ours/theirs`/silent-discard) / → prose (no binary resolves a two-writer frontmatter conflict); the FO sites then collapse to `«halt.rebase-conflict»(paths)`. The one genuine cross-file fail-safe the whole-contract view reveals; net-neutral-to-negative, fits AC-4. (FO and ensign load separate cores — the `«fn»` can't be physically shared across that boundary; the win is one canonical FO definition instead of two inline restatements + a dangling "below" reference whose target sits in the ensign file.)

6. **Consolidate the `## Dispatch (deferred module)` + `## Merge and Cleanup (deferred module)` pointer prose into one `«fn»`-registry block** (first-officer-shared-core.md). The boot core spends ~140 + ~158 tok on two prose sections describing *what loads when / how each is realized* — the same map a ~14-line arrow-grouped `«fn»` registry (shipped / host-bound / prose × load tier, composition edges inline) carries more scannably, and it mirrors what `prose_function_routing_test.go` already extracts. Replace the two pointer sections with the registry block, preserving every load-point + guard they carry (the dispatch-reference load point, the merge terminal-boundary load point, the greet-guards). Net-neutral-to-negative in-core; the registry doubles as the navigation index for the sibling `entity-status` deferral task.

Plus (conditional): if/after M2 lands, sync docs/runtime-support.md §"Runtime binding-block shape" so it does not describe a 2-into-1 worker.spawn/dispatch-build binding the core now contradicts (M2 promotes «dispatch.build» to a first-class shipped host-neutral «fn»).

## Acceptance criteria

- **AC-1** — reuse-condition-4 is a `«reuse.model-match»` «fn» with guard/effect/block/done-when; the verbatim diagnostic pin now sits on its body and stays discriminating (a mutation that drops the diagnostic FAILS; a non-vacuity control holds).
- **AC-2** — `«dispatch.next-action»`'s arrow no longer reads as judgment-owned: it is retagged or annotated as deterministic-mechanism-binary-pending, consistent with runtime-support.md's trichotomy definition.
- **AC-3** — `«worker.spawn»` and `«worker.shutdown»` each carry a kind-only `→ runtime-binding` pointer; the contractlint guard permits the kind-only arrow while still banning a per-host token in the core body; build + contractlint green.
- **AC-4 (value guardrail)** — measured two ways, not asserted (see `## AC-4 measurement` below for the per-file budget table): (a) z4's OWN per-file delta — every FO contract file z4 edits is net-neutral-or-NEGATIVE vs its state at z4-implementation start (`origin/main`), the property z4 controls and the promotions deliver (prose→structure); and (b) the absolute ceiling — each file z4 restructures stays ≤ its v0.22.0 baseline, `wc -c` of the working file vs `git show v0.22.0:<path>`. (a) is z4's gate; (b) holds for `first-officer-shared-core.md` (the Wave-2 deferrals open ~1.7k tok of headroom) and `fo-dispatch-core.md` (+24 B headroom today — TIGHT; #1's «fn» promotion must net-trim enough flat prose to absorb #2's +1 annotation line and #3's 2 arrow lines). `claude-fo-dispatch.md` is ALREADY +508 B over v0.22.0 on `main` and z4 only collapses one rebase-halt reference in it (net ~0) — so its absolute (b) is NOT z4's charge; z4 owns only (a) ≤ 0 for it (do not worsen the pre-existing overage). FLAGGED for the gate.
- **AC-5** — runtime-support.md reflects the final binding shapes (no authority-doc-vs-core drift), including the dispatch.build sync if M2 has landed.
- **AC-6** — the AC coverage cross-check is a `«gate.ac-cross-check»` «fn» (guard/effect/block/done-when) that folds bm's end-value re-anchor as an explicit mechanism→value guard/block, not inference-prose. Because the cross-check is FO judgment (`→ prose`, no binary), a prose-grep for it would be the banned tautology; its real behavioral proof is the ALREADY-EXISTING live scenario `TestLiveReanchorGateRejectsMeansOnlyRegressed` (`internal/ensigncycle/ac2_reanchor_live_test.go`, build-tag `live`), whose fixture (`livescenario.AuthorACReanchorScenario` — AC-1 mechanism-only + AC-2 regressed end-value) IS "a fixture whose mechanism-only AC's served value-AC regressed gates REJECT," graded on an observed REJECT + re-anchor reasoning with the unmutated-body grading as its non-vacuity control. The restructure is behavior-PRESERVING, so validation RE-RUNS this scenario once against a real credential to confirm the guard/block wording reproduces the REJECT — a prose restructure can shift model behavior, so this is exercised, not assumed. (bm merged #441 2026-06-30 — unblocked.)
- **AC-7** — the rebase-conflict halt is a `«halt.rebase-conflict»(paths)` «fn» defined once in first-officer-shared-core.md; the FO restatement sites (~162, ~177) and the claude-fo-dispatch reference resolve to it by name; every guard survives (abort, surface conflicting paths, never force/`-X`/silent-discard); a grep confirms 1 definition + ≥3 by-name references + 0 inline restatements, and the file delta is ≤ baseline.
- **AC-8** — the `## Dispatch (deferred module)` + `## Merge and Cleanup (deferred module)` pointer prose is replaced by one `«fn»`-registry block; every load-point + guard the two sections carried is preserved (a check confirms the dispatch-reference load point, the merge terminal-boundary load point, and the greet-guards still resolve); first-officer-shared-core.md is net-NEGATIVE for this refinement (the one occupancy-positive item in z4, still inside AC-4's ≤-baseline ceiling).

## Contractlint guard-loosening design (AC-3)

The only binary-adjacent change z4 needs: loosen `TestDispatchCoreDefinesWorkerLifecycleCapabilities` (`internal/contractlint/capability_binding_test.go`) to PERMIT a kind-only `→ **runtime-binding**` arrow on the `«worker.spawn»`/`«worker.shutdown»` blocks while STILL banning a per-host `**Host:**` token in the core body.

Today's body bans *every* arrow (`strings.Contains(block, "\n- → ")`) and every `**Claude/Codex/Pi:**` token. Loosen to: every `- → ` line in the block MUST be the kind-only runtime-binding pointer; any other arrow reds; the per-host token ban is unchanged.

```go
for _, name := range []string{"worker.spawn", "worker.shutdown"} {
    block := fnBlock(t, data, name)
    for _, line := range strings.Split(block, "\n") {
        trimmed := strings.TrimSpace(line)
        if !strings.HasPrefix(trimmed, "- → ") {
            continue
        }
        if !strings.HasPrefix(trimmed, "- → **runtime-binding**") {
            t.Errorf("capability «%s» core block carries a non-runtime-binding arrow %q; the only arrow permitted host-neutral is a kind-only `→ **runtime-binding**` pointer", name, trimmed)
        }
    }
    for _, host := range capabilityHosts { // UNCHANGED — still bans a concrete per-host binding
        if strings.Contains(block, "**"+host+":**") {
            t.Errorf("capability «%s» core block contains concrete host binding for %s", name, host)
        }
    }
}
```

The exact arrow line to add to each of the two `## «worker.spawn»` / `## «worker.shutdown»` blocks in `fo-dispatch-core.md` (kind-only, no host named):

    - → **runtime-binding**: bound in the host adapter's `## Runtime implementation`

**No second-order break in `TestCapabilityBinding`** (same file): worker.spawn/shutdown enter its `defined` set via the `isRuntimeBoundLifecycleCapability` branch (independent of any arrow) and are SKIPPED in the per-host coverage loop; `hostsBoundByArrow` parses the new arrow, finds no host token, returns empty — so the set-equality and host-coverage halves are unaffected. **No break in `TestProseFunctionNotationBindsToRouting`** (cli): `migrationBareRe` matches only `shipped|prose`, so a `runtime-binding` arrow contributes no verb binding.

**Mechanism spike (DONE — riskiest path exercised first per proof policy):** a throwaway driver replicating the loosened detector graded 5 cases: kind-only runtime-binding arrow PASSES; no-arrow (current shape) PASSES; a per-host arrow REDS; a `→ **shipped**` arrow REDS; and — the over-loosening guard — a runtime-binding arrow that smuggles a `**Claude:**` token in its prose STILL REDS on the per-host ban. All 5 graded as expected. The string-logic loosening is proven sound before the gate; the only unverified residue is the live-scenario behavior preservation for AC-6 (re-run at validation).

## #6 registry-block design + sequencing (AC-8)

Replace the two prose pointer sections (`## Dispatch (deferred module)` ~587 B, `## Merge and Cleanup (deferred module)` ~666 B; ~1.25 KB combined) in `first-officer-shared-core.md` with ONE arrow-grouped `«fn»`-registry block grouped by realization-kind × load-tier, composition edges inline. It indexes FOUR deferred references: Dispatch, Merge, AND the new Status-Viewer (`84`) + Write-Scope (`k4`) references. Shape (≤14 lines, net-NEGATIVE vs the ~1.25 KB it replaces — feasible at ~60 B/line):

- one row per deferred module: `module name → realization kind (runtime-binding) → core-file path → load-point trigger → guard`.
- preserved load-points the closure tests bind to (all three MUST survive): the **dispatch-reference** load point (first worker dispatch) naming `references/fo-dispatch-core.md`; the **merge** terminal-boundary load point naming `references/fo-merge-core.md`; the **greet-guards** (a greet-and-stop boot reads NONE of them).

**Two hard constraints the registry must honor** (verified against the on-disk tests):
1. The literal filenames `fo-dispatch-core.md` and `fo-merge-core.md` MUST appear (as `references/<file>.md`) — `TestHostNeutralCoresResolveAndCarryCeremony` asserts `sharedBody.Contains(base)` for each, and `TestBootResidentDeferredLoadPointsResolve` `os.Stat`s every `references/X.md` the body names. Drop a filename → red.
2. The Status-Viewer/Write-Scope reference paths z4 names are exactly the new `references/*.md` files that `84` and `k4` create. `TestBootResidentDeferredLoadPointsResolve` `os.Stat`s them, so **naming `references/<status-viewer>.md` before `84` creates that file REDS the closure walk.**

**Sequencing dependency (THIS is why z4 is Wave 3 / last):** because the registry references the reference-files `84` (entity-status, defers Status Viewer) and `k4` (defer-write-scope-id-styles) create, z4's #6 implementation MUST land AFTER `84` + `k4` (Wave 2). z4 design names them generically ("the Status-Viewer reference `84` creates", "the Write-Scope reference `k4` creates"); the exact paths bind at z4 implementation once Wave 2's filenames are on disk. `84` + `k4` are ideating in parallel now; the registry doubles as the navigation index those deferral tasks point back to.

## AC-4 measurement (per-file budget vs v0.22.0)

Measured with `wc -c` of the working file vs `git show v0.22.0:<path>` (bytes; corroborate with `wc -l`). Today on `main`:

| FO contract file | v0.22.0 | current (main) | headroom | z4 touches | (b) absolute ≤-baseline achievable by z4? |
|---|---|---|---|---|---|
| `first-officer-shared-core.md` | 28586 | 29181 | −595 | #4 #5 #6 | YES — Wave-2 deferrals remove ~1.7k tok first; z4 adds #4/#5 (net-neutral-to-neg) + #6 (net-NEGATIVE) into the opened headroom |
| `fo-dispatch-core.md` | 17488 | 17464 | +24 | #1 #2 #3 | TIGHT — not helped by deferrals; #1's promotion must net-trim ≥ (#2 +1 line + #3 +2 arrow lines) to stay ≤ baseline |
| `claude-fo-dispatch.md` | 22489 | 22997 | −508 | #5 (1 ref collapse) | NO via z4 alone — already +508 over; z4 owns only (a) net-neutral-or-neg, NOT the pre-existing overage |
| `fo-merge-core.md` | 8059 | 7454 | +605 | none | n/a (untouched) |

**Finding flagged for the gate:** AC-4 as written ("every FO file ≤ v0.22.0 baseline") cannot be literally satisfied for `claude-fo-dispatch.md` by z4 — it is already over baseline from post-v0.22.0 work and z4 does not reduce it by 508 B. Recommended ruling: AC-4's per-file ABSOLUTE ceiling (b) scopes to the files z4 (or the Wave-2 deferrals) materially restructure; for `claude-fo-dispatch.md` z4's obligation is the DELTA guardrail (a) only — its edit nets ≤ 0. The pre-existing overage is another member's charge or out of 0240 scope.

## Test plan

Primary lane is OFFLINE — `go build ./...` + `go test ./internal/contractlint/ ./internal/cli/ ./internal/ensigncycle/` — for AC-1/2/3/5/7/8 (content/structure, contractlint-provable). AC-6 additionally re-runs ONE existing live scenario. Per AC:

- **AC-1 (`«reuse.model-match»` + moved diagnostic pin):** offline. The verbatim pin `does not match next stage effective_model` is guarded by `TestProseFunctionCatastropheClausesSurvive` (catastropheClause for `fo-dispatch-core.md`) + its discriminator `…BackstopDiscriminates`. **Mutation test (red-then-green):** after moving the pin onto the new «fn» body, delete the pin substring from that body → `go test -run TestProseFunctionCatastropheClausesSurvive ./internal/contractlint/` goes RED → restore → GREEN. Non-vacuity is the existing `…BackstopDiscriminates` control. (`«reuse.model-match»` is dotted, so invisible to `capability_binding_test.go`'s `[a-z][a-z-]+` extractor and to the routing test if its `→` carries no `spacedock` verb — recommend `→ prose` annotated deterministic per #2, deferring the per-host model space to `«worker-identity»`.)
- **AC-2 (`«dispatch.next-action»` arrow):** offline. Keep the `→ **prose**, becomes `spacedock dispatch next-action`` line verbatim (preserves `TestProseFunctionNotationBindsToRouting`'s should-have-flipped guard — the `becomes`-form already signals "deterministic, binary-pending"); ADD one annotation line at the definition disambiguating it from judgment-`→ prose`. Paired with the runtime-support.md:21 clause (AC-5). No new notation token (avoids perturbing the routing test / trichotomy). Proof: routing test stays green; the disambiguation is the annotation's presence (legible-for-human, not a behavioral claim).
- **AC-3 (worker.spawn/shutdown `→ runtime-binding` + guard loosening):** offline. The loosened `TestDispatchCoreDefinesWorkerLifecycleCapabilities` above + add a discriminator/RED control (mirroring the file's existing controls): plant a per-host arrow AND a runtime-binding-arrow-with-smuggled-host-token against a fixture block, assert both red — proving the loosening did not over-loosen (the spike already graded these). `go build ./...` + full `go test ./internal/contractlint/` green.
- **AC-4 (value guardrail):** the per-file `wc -c` vs `git show v0.22.0:<path>` measurement above, recorded in the implementation stage report (the byte delta is the AC, not a prose claim). z4-own-delta (a) ≤ 0 per touched file vs `origin/main`.
- **AC-5 (runtime-support.md):** offline. Concrete doc diff below; no test (it's an authority-doc sync). Cross-checked against the core's final arrow shapes by reading both — AC-5's proof is that the binding shapes named in runtime-support match the cores' final `→` lines.
- **AC-6 (`«gate.ac-cross-check»`):** offline structure (the «fn» heading parses with the four fields) + ONE live re-run of `TestLiveReanchorGateRejectsMeansOnlyRegressed` (`-tags live`, real credential) confirming the restructured guard/block reproduces the REJECT on the existing regressed-end-value fixture. The live scenario's unmutated-body + REJECT-reasoning grading is its built-in non-vacuity control.
- **AC-7 (`«halt.rebase-conflict»(paths)`):** offline. The verbatim pin `do not force-push or auto-resolve` (catastropheClause for `first-officer-shared-core.md`) moves onto the new «fn» body. **Mutation test (red-then-green):** delete that substring from the «fn» body → `TestProseFunctionCatastropheClausesSurvive` RED → restore → GREEN. Plus a structural grep asserting 1 `## «halt.rebase-conflict»` definition + ≥3 by-name `«halt.rebase-conflict»` references (the two shared-core sites `«state.ensure-ready»`~162 & `«state.commit»`~177, and `claude-fo-dispatch.md`~162) + 0 inline `git rebase --abort` restatements at the call sites. The «fn» body must phrase the abort step to fit both the FO-holds-the-conflict case (`«state.ensure-ready»` manual pull → FO runs `git rebase --abort`) and the already-aborted case (`«state.commit»` exit 3 → the binary aborted; FO surfaces + stops).
- **AC-8 (registry block):** offline. `TestHostNeutralCoresResolveAndCarryCeremony` + `TestBootResidentDeferredLoadPointsResolve` (already green) re-run after the block lands — they assert the dispatch/merge filenames resolve and every named `references/X.md` exists on disk (this is the closure guard that enforces the Wave-3 sequencing). Plus the `wc -c` delta confirming `first-officer-shared-core.md` is net-NEGATIVE for #6 alone (block < the two sections it replaces).

**Cost/complexity:** low. One offline dispatched worker under contractlint; ~2 mutation cycles (AC-1, AC-7 pins, ~minutes each); one `-tags live` re-run for AC-6 (the only real-credential spend, ~one Claude session). No new fixtures needed — AC-6 reuses the shipped `AuthorACReanchorScenario`; AC-1/AC-7 reuse the shipped catastrophe-backstop.

## AC-5 runtime-support.md doc diff (concrete)

§ "A `→` line states where the capability is realized" (runtime-support.md:21), the `→ prose` clause — disambiguate judgment-owned from deterministic-binary-pending so #2's `«dispatch.next-action»` reads correctly:

- BEFORE: `… `→ prose` marks an obligation that stays judgment-owned with no binary backing.`
- AFTER: `… `→ prose` marks an obligation with no binary backing: judgment-owned when it names no verb, or a deterministic mechanism still hand-followed when it names a `becomes` verb a binary will later ship (e.g. `«dispatch.next-action»` → `spacedock dispatch next-action`, descoped to roadmap 0222).`

Conditional (only if M2 has landed — promotes `«dispatch.build»` to a first-class shipped host-neutral «fn»): sync § "Runtime binding-block shape" so the `«worker.spawn»` bullet no longer implies a 2-into-1 worker.spawn/dispatch-build binding the core now contradicts. Gate the conditional on `grep`-confirming M2's state at implementation.

## Spike record

- **Spiked (riskiest path, done in ideation):** the AC-3 contractlint guard-loosening string-logic — 5 graded cases incl. the over-loosening guard (runtime-binding arrow + smuggled host token still reds). Proven sound; details under `## Contractlint guard-loosening design`.
- **No spike needed for the rest, on these proven mechanisms:** the moved-pin guards already SHIP and are green (`TestProseFunctionCatastropheClausesSurvive` + discriminator) — AC-1/AC-7 ride them; the AC-6 behavioral guard already SHIPS as a live scenario — #4 re-runs it (a behavior-preserving restructure of text it pins, not a new mechanism); the closure walk that enforces the AC-8 sequencing already SHIPS and is green. No parser round-trip, on-disk format, or new tool flag is introduced (AC-3's change is a test-logic loosening, not a new binary surface — confirmed in scope by the roadmap's "No new binary behavior beyond the contractlint guard-loosening z4 AC-3 requires").

## Stage Report: ideation

- DONE: The net-neutral value guardrail (AC-4) holds: every FO contract file stays ≤ its v0.22.0 baseline — each of the 6 refinements trades prose for structure, and the one occupancy-positive item (#6 Dispatch/Merge → registry, AC-8) stays net-NEGATIVE in-core. Measured vs baseline, not asserted.
  Built the per-file `wc -c` vs `git show v0.22.0:<path>` budget table (`## AC-4 measurement`); refined AC-4 to a two-part measure (z4-own-delta ≤0 + absolute ≤-baseline). FINDING: `claude-fo-dispatch.md` is already +508 B over v0.22.0 on `main` and z4 only collapses one ref in it — absolute (b) is not z4's charge; flagged for the gate.
- DONE: The contractlint guard-loosening (AC-3): permit a kind-only `→ runtime-binding` arrow on «worker.spawn»/«worker.shutdown» while STILL banning a per-host `**Host:**` token; design the test change and a red-then-green mutation test for each moved verbatim pin.
  Designed the exact loosened `TestDispatchCoreDefinesWorkerLifecycleCapabilities` (`## Contractlint guard-loosening design`) and SPIKED it offline — 5 graded cases incl. the over-loosening guard, all as expected. Mutation tests for both moved pins (AC-1 `does not match next stage effective_model`; AC-7 `do not force-push or auto-resolve`) ride the shipped `TestProseFunctionCatastropheClausesSurvive`: delete-pin→RED, restore→GREEN.
- DONE: The #6 registry block indexes the NOW-DEFERRED structure — design it to index Dispatch, Merge, AND the new Status-Viewer (`84`) + Write-Scope (`k4`) references; record that z4 must sequence AFTER `84`+`k4` at implementation (Wave 3).
  `## #6 registry-block design + sequencing`: ≤14-line arrow-grouped block; the two hard closure constraints (`fo-dispatch-core.md`/`fo-merge-core.md` filenames must appear; the 84/k4 ref paths are `os.Stat`'d) make the Wave-3 sequencing load-bearing — naming a Wave-2 ref before it exists REDS `TestBootResidentDeferredLoadPointsResolve`.
- DONE: complete the test plan + AC-5 doc diff.
  Expanded `## Test plan` to a per-AC plan; added the concrete runtime-support.md:21 before/after for AC-2/AC-5. Confirmed current baseline green: `go build ./...`, `go test ./internal/contractlint/` and `TestProseFunction*` in `./internal/cli/` all pass.

### Summary

Refined and completed the ideation design without rewriting it: added the exact contractlint guard-loosening (proven by an offline 5-case spike, incl. the over-loosening guard), a per-AC test plan with red-then-green mutation recipes for both moved pins, the #6 registry-block shape with its load-bearing Wave-3 sequencing constraint, and a per-file AC-4 byte-budget table. Two findings need a captain ruling at the gate: (1) AC-6's re-anchor restructure has an EXISTING live behavioral guard that must be re-run once — "no live lane" holds for the other ACs but not AC-6; (2) `claude-fo-dispatch.md` is already +508 B over its v0.22.0 baseline, so AC-4's absolute ≤-baseline ceiling is not z4's charge for that file (z4 owns only its net-neutral delta). The riskiest mechanism (AC-3 string-logic) was spiked first per the proof policy; everything else rides already-shipping green guards.
