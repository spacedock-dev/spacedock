---
title: «fn» binding refinements — promote reuse-condition-4 + «gate.ac-cross-check» (pending bm), disambiguate →prose, add worker.spawn/shutdown →runtime-binding pointer
status: backlog
score: 0.35
group: cleanup
issue:
id: z4qnbcsv0ffty4njszg214wd
---

Post-v0.23.0 refinement from the 0230 «fn»-binding effectiveness audit (5-lens, verdict "mostly-effective-with-fixes"). The binding program is sound and ship-safe — every load-bearing guarantee is machine-pinned and the arrow trichotomy is the standout legibility win — so these are quality refinements, NOT correctness fixes, and were deliberately deferred past the v0.23.0 cut. (The audit's scariest finding, a claimed worker.spawn/shutdown over-collapse in the M2 trim, was VERIFIED a false alarm: worker.shutdown content is preserved and worker.spawn is a legitimate DRY relocation into the heading + the promoted «dispatch.build» «fn».)

## The four refinements

1. **Promote reuse-condition-4 to a `«reuse.model-match»` «fn»** (fo-dispatch-core.md). It is the heaviest deterministic comparator in the dispatch core (four branches: match / null-skip / captain-session-fallback-never-matches / mismatch-emit-diagnostic) yet rendered as flat prose while thinner per-host probes ARE «fn»s — the prose/binding boundary is drawn on the wrong side. Give it guard (null-skip) / effect (resolve + compare in the «worker-identity» canonical model space) / block (a captain-session fallback value never matches → forced one-time fresh re-stamp) / done-when (the verbatim `does not match next stage effective_model` diagnostic), and MOVE the existing verbatim-diagnostic catastrophe-backstop pin onto the new «fn» body (keep it discriminating + non-vacuous).

2. **Disambiguate `→ prose`** for `«dispatch.next-action»` (fo-dispatch-core.md). runtime-support.md:21 defines `→ prose` as "judgment-owned, no binary backing," but the next-action body is a deterministic 3-step scheduler awaiting a binary (descoped to roadmap 0222), not judgment. Either retag it (a "mechanism, binary pending" status distinct from `→ prose`) or add one line at the definition stating it is hand-followed-but-deterministic, so a future implementer neither reads it as non-codifiable nor probes for a not-yet-shipped `spacedock dispatch next-action`.

3. **Add a `→ runtime-binding` pointer line to `«worker.spawn»` and `«worker.shutdown»`** (fo-dispatch-core.md). These two highest-stakes capabilities carry NO `→` line in either tree — a cold FO could read the host-neutral body as "nothing to do here" with no in-file signal that it is runtime-bound or which adapter section binds it (the Claude adapter realizes them as prose headings `## Spawn Call (Agent)` with no guillemet match). Add a kind-only `→ runtime-binding: bound in the host adapter's ## Runtime implementation` line; requires loosening TestDispatchCoreDefinesWorkerLifecycleCapabilities to permit a kind-only arrow while still banning a per-host `**Host:**` token in the core body.

4. **Promote the AC coverage cross-check (+ end-value re-anchor) to a `«gate.ac-cross-check»` «fn»** (first-officer-shared-core.md, `## Completion and Gates`). The cross-check is flat prose while `«gate.assemble-verdict»` right beside it is a «fn» — the prose/binding boundary is on the wrong side (the #1 defect again). Give it guard (a value-measuring AC is paired to the mechanism-only AC) / effect (scan each `**AC-N**` for evidence; resolve the mechanism→value "serves" pairing) / block (a mechanism-only AC whose served value-AC regressed → REJECT, not pass) / done-when (every AC has evidence or is named). Folding it into guard/block — instead of today's inference-prose — closes two of bm's own validation edge-findings: the "only bites when a paired value-AC exists" precondition becomes the **guard**, and the "mechanism→value pairing left to FO inference" becomes the **effect**.

   **PENDING bm.** This RESTRUCTURES bm's EDIT A (the re-anchor clause), which is live in bm's worktree pending re-validation + merge — it MUST sequence AFTER bm lands, or it conflicts. Cross-ref: bm = `bmt9h66tg1s3eda1e1vxmzj` gate-on-end-value; its validator's adversarial-audit edge findings are the motivation for the guard/block structure here. Same net-neutral value-guardrail as AC-4 (= bm's AC-3).

Plus (conditional): if/after M2 lands, sync docs/runtime-support.md §"Runtime binding-block shape" so it does not describe a 2-into-1 worker.spawn/dispatch-build binding the core now contradicts (M2 promotes «dispatch.build» to a first-class shipped host-neutral «fn»).

## Acceptance criteria

- **AC-1** — reuse-condition-4 is a `«reuse.model-match»` «fn» with guard/effect/block/done-when; the verbatim diagnostic pin now sits on its body and stays discriminating (a mutation that drops the diagnostic FAILS; a non-vacuity control holds).
- **AC-2** — `«dispatch.next-action»`'s arrow no longer reads as judgment-owned: it is retagged or annotated as deterministic-mechanism-binary-pending, consistent with runtime-support.md's trichotomy definition.
- **AC-3** — `«worker.spawn»` and `«worker.shutdown»` each carry a kind-only `→ runtime-binding` pointer; the contractlint guard permits the kind-only arrow while still banning a per-host token in the core body; build + contractlint green.
- **AC-4 (value guardrail)** — every FO file stays ≤ its v0.22.0 baseline (net-neutral or negative; the promotions trade prose for structure).
- **AC-5** — runtime-support.md reflects the final binding shapes (no authority-doc-vs-core drift), including the dispatch.build sync if M2 has landed.
- **AC-6** — the AC coverage cross-check is a `«gate.ac-cross-check»` «fn» (guard/effect/block/done-when) that folds bm's end-value re-anchor as an explicit mechanism→value guard/block, not inference-prose; a fixture whose mechanism-only AC's served value-AC regressed gates REJECT, with a non-vacuity control. **BLOCKED-ON:** bm (gate-on-end-value) merged first.

## Test plan

Offline `go test ./internal/contractlint/ ./internal/ensigncycle/` + `go build ./...`. The «fn» promotions are contract-file edits → dispatched worker under contractlint; no live lane required (no behavioral change — the bodies already say what the prose said). Mutation-test each moved pin (red-then-green).
