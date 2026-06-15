# 0204 structured-reads - preflight staff review

**Verdict: gates-ready after minor folds.** Independent preflight over the pooled sprint, per-member plus cross-cutting. Refute-not-bless: a fresh reviewer per member tried to break each design, and the j7 spike was independently re-verified against the pinned harness. Shaping-FO session, 2026-06-15.

The drivable set is `6r`, `0q`, `48` (all `ideation`, `sprint-readiness: ready`) plus the backbone `e6a` (`status-section-reader`, already in `implementation`). `j7` (`status-set-staleness-echo`) dissolves to a roadmap decision below and drops from the queue. `ey` (`proof-policy-shipped-scaffolding`) is NOT a sprint member (`status: backlog`, no `sprint:` stamp) and composes with `48` from the backlog.

## Per-member readiness

### 6r - ci-log-read-summary - minor-gaps

The read-time-over-CI-emitted choice is sound and survived refutation: the savings is context tokens not disk, so a read-time parser captures it over every historical log with zero CI change, and reading-not-replacing buys "full log reachable / no info loss" for free. The `-v`-text grammar is the right and only target (`.github/workflows/runtime-live-e2e.yml:179,198` pipe `go test ... -v 2>&1 | tee`; no `-json`/`gotestsum` anywhere). Scope is well-fenced. AC-1/AC-3/AC-4 are oracle-based Go tests; AC-2 is genuinely behavioral (triage answer-set matches an external oracle, not a phrase match).

**Fold before the gate:**

- **AC-6 is the weak AC** and must rise to the `e6a` sibling bar. Its sole proof today is "the doc diff ... plus the existing skill-text/doc-contract guards passing over the edited files" - a prose-grep over the file under test, the exact pattern the project's own archived precedent rejects (`ban-readme-substring-assertions`, PR #378, PASSED: doc-contract guards over instruction prose are tautological and should be DROPPED, not retargeted). The entity concedes the gap ("only pins that the wording exists and is under-test"). `e6a`'s own AC6 was rewritten after a captain bounce-back to demand a LIVE drive of the helper being called at the real site. AC-6 must require a live drive where an FO/ensign fetches a transcript and runs the summarizer at the real triage site, the proof being that trace (helper invocation plus the short read replacing the whole-log read), never the instruction text. The recurring token savings has no behavioral oracle until this lands. *Owner: ideation-rework.*
- **AC-6's target site is partly invented.** Grep for the "read the CI transcript whole" wording across `skills/` and `docs/dev/` returns only the two entity bodies; no live skill or README carries it. The named `docs/dev/README.md` "Runtime Live CI" section is about proving host behavior via live runs, not reading a transcript for triage. AC-6 must own LOCATING or creating the real triage-read site, not assert one. *Owner: ideation-rework.*
- **AC-5 may be unsatisfiable as written.** The design says output line 2 is "source path + total line count"; a stdin read (`-`) has no path, so the path-vs-stdin byte-identical assertion fails unless the header is normalized. Flag for implementation. *Owner: captain-ratify-at-gate (note for the implementer).*
- Polish: reconcile the 143KB (cited-historical) vs 170KB/172KB (spike) figures; sequence the README edit to land after the known pre-existing `TestSharedScenarioDocsContract` red from the README-slim resolves, so this work does not inherit or re-create a prose-grep guard.

The doc/skill edit (`first-officer-shared-core.md`, the dev README) is shipped agent-facing scaffolding and is owed a detached adversarial audit at validation.

### 0q - status-source-not-default - minor-gaps

Sound and minimal: one line removed from `defaultStatusFields`, SOURCE re-enters through the existing `--fields`/`--all-fields` projection, no new flag (correct YAGNI). The spike was independently reproduced live against the binary (every claim, byte-for-byte). ACs are runtime header-token / JSON-key assertions over the live binary, no prose-grep, correctly mirroring the captain-approved `fields_dedupe_test.go` precedent. Direction A for the open fork (bare `--json` also drops `source`) is the right call, and the in-repo blast radius is effectively nil: the FO reads `--boot --json` (a different envelope that never carried `source`), and no in-repo consumer reads `entities[].source` from bare `status --json`. The captain decision rests purely on external-consumer policy, where A is cleaner than B's split human-vs-JSON default sets.

**Fold before the gate (two material gaps from the reviewer):**

- **Parity-handling under-scopes the oracle divergence.** The design pulls only the `default` case out of oracle-parity, but `archived`, `where-status`, and `all-fields` ALSO diverge from the frozen oracle under this change (archived/where lose the SOURCE column; all-fields moves SOURCE from fixed slot 6 to a sorted extra). An implementer following "pull the default case out" literally would regenerate those goldens IN PLACE while leaving them inside the parity assertions, silently re-pinning native to bytes the frozen oracle no longer produces - corrupting the "goldens ARE the oracle bytes" invariant `zz_independent_parity_test.go:19-24` asserts. The pull-out must extend to all four diverging cases. *Owner: ideation-rework.*
- **AC-3 / test-plan name only `seq-default.json`.** Under direction A, `seq-archived.json` (5 `source` keys) and `seq-where.json` (2) also drop `source` and are oracle-derived cases. AC-3 must cover the bare `--json` archived and where reads too, or explicitly scope to default and note the others regenerate under the same one-knob change. The implementer is currently pointed at a set ~7x too small (3 named text goldens vs ~28 that shift; 1 JSON golden vs 3). *Owner: ideation-rework.*
- Polish: add a column-ORDER assertion for `--all-fields` (SOURCE now sorts among extras, not slot 6) so the reorder is checked, not only frozen; reword the "Out of scope: any reorder" line to acknowledge the unavoidable `--all-fields` re-sort as an in-scope consequence.

The regression sweep is the safety net that catches the wider set, so the net is recoverable, but the test plan must be widened so the implementer is not surprised.

### 48 - commission-templates-defer-to-contract - minor-gaps

Sound, well-scoped prose-slim. The central read-cost finding was verified against the files: commission `Read`s the whole template at Phase 1 (where the cost lands) but auto-generates the commissioned README from `SKILL.md`'s inline structure at Phase 2, so per-stage template prose is reference-only and safe to slim. The one load-bearing coupling - `development.md:113` "Adopt them by copying the guidance into the validation stage" - is correctly identified as the opt-in block whose dev-specifics must survive as pointers, and the rewire is specced. AC-1 is genuinely behavioral (live commission + boot + tautological-AC rejection drive, graded on durable end-state, legitimate absence-grep on the generated README). AC-2 is a sanctioned single-source check modeled on `TestStartupGateGuidanceHasSingleSource` (binds divergeable contract-vs-template content, fires on drift), not a prose-grep tautology. The three slim moves faithfully port the verified dev-README slim commit `48edae4c`.

**Fold before the gate (residuals to name, not blockers):**

- **AC-2's discriminating power is only as strong as its unnamed marker set.** The entity and the boundary-guard doc both concede a meaning-inverting paraphrase keeps every grepped token, so AC-2 proves "these 3-5 phrasings live in one place," not "the doctrine does." Pin the marker set in the body, or explicitly defer it to the detached audit (which must attack it: a paraphrase that re-restates the doctrine while dodging the markers means the set is too narrow). *Owner: ideation-rework (pin) or captain-ratify-at-gate (defer to audit).*
- **State the development-only descope plainly at the gate.** AC-1 uses development as the single live representative; experiment/refinement are covered structurally (AC-2/AC-3) plus a no-spend commission dry-check, NOT a live drive. This is acceptable cost-scope, not a hidden gap (the governance mechanism - the FO loads the contract at boot - is identical across all three shapes, and only development has a paragraph-level removal). The condition: the dry-check MUST run Phase-1 to completion (or a faithful hand-trace) for experiment and refinement too, not just development, so their adoption machinery is confirmed intact. *Owner: captain-ratify-at-gate.*
- **AC-1's grader must reward the durable rejection end-state only**, never a specific rule-citation in the transcript. The current contract rejects the seeded scenario via the broad "prove by exercising" principle, not a verbatim independent-source rule (that verbatim rule is `ey`'s to add). A validator must not red a correct rejection that cited the wrong sentence. The body already says "never transcript phrasing"; make the no-rule-citation point explicit. *Owner: captain-ratify-at-gate.*

The templates are shipped scaffolding (a named high-stakes surface) and are owed a detached adversarial audit at validation. Implementation runs in a worktree under a dispatched worker, never an FO-direct edit.

## j7 - status-set-staleness-echo - TERMINALIZE as a roadmap decision

**Re-verification verdict: confirmed-with-caveats. Safe to terminalize.** An independent verifier rebuilt the binary and re-ran the falsification live on the pinned Claude Code 2.1.177 across four `claude -p --output-format stream-json --verbose` drives, real entities up to ~27k tokens: a `Read` followed by `status --set` of the same file does NOT re-emit the file. Post-mutation `cache_creation` stayed at 217-381 tokens, including the exact clean contract scenario (single full Read then immediate set), with `cache_read` climbing monotonically and no staleness reminder firing. Two orders of magnitude below a whole-file re-emit.

Critically, the verifier closed the spike's one methodological hole - it had never shown its oracle could detect the echo it claimed was absent. The verifier added the missing positive sensitivity control: forcing a genuine `Read` after the set produced `cache_creation=20505`, proving the harness's own `usage.cache_creation_input_tokens` oracle IS sensitive to a real whole-file write. So j7's negative is a true negative, not an instrument miss. The oracle is external (the inference layer's usage record, not prose), read at the right place, and the conclusion holds across two file sizes and multi-turn continuation. AC-2's narration is verified live and in code (`handlers.go:306`, `atomicWrite` at `mutate.go:244`/`os.Rename:260`).

Caveats recorded, none blocking: a single uninterrupted full Read of a >25k-token entity cannot occur on this harness (Read truncates at 25000 tokens and paginates, so the largest-file worst case the prompt asked about is already chunked); folder-form `{slug}/index.md` entities and non-`status` Bash mutations were not directly tested (the same-inode in-place rewrite arm is the closest proxy and showed no echo); a full live FO-boot cache topology differs from short headless sessions (the sealed-cache arm and the multi-turn drive bracket it). The result is correctly scoped to 2.1.177, which the entity pins in Out-of-scope - a future harness regression reopens it, the right hedge.

This is the DoD-sanctioned outcome: the DoD's second bullet names this sink by example ("e.g. the read-then-`set` echo if it turns out harness-inherent ... recorded as a roadmap decision rather than forced into code"). Nothing is shippable (`atomicWrite` is already in place and is not the lever; even a same-inode rewrite produced no echo), so the roadmap decision IS the satisfied outcome via AC-1's roadmap-decision branch.

**Contract reason-correction rides inside j7's roadmap-decision record** (per the captain's decision). The falsification contradicts the live contract, which asserts the echo as current behavior at TWO stale sites that must be corrected as a separate change against the live skill surface:

- `skills/first-officer/references/first-officer-shared-core.md:218` - "On Claude Code, a `Read` followed by a Bash mutation of the same file (including `status --set`) triggers the file-staleness safety net, echoing the file back as cache-write tokens."
- `skills/first-officer/references/claude-first-officer-runtime.md:39` - "The Claude Code runtime is where the Read-then-Bash-mutation staleness echo fires."

Both state the echo as a live fact; both are wrong on 2.1.177 for the Bash path. The Grep-over-Read PREFERENCE itself stays valid (Grep is cheaper regardless); only the stated staleness-echo REASON is stale. The reason-update is out of j7's ideation scope (it touches the live contract) and is flagged in the roadmap-decision record for a follow-up, naming both sites.

## Cross-entity coherence

- **48 <-> ey composition is coherent and clobber-safe.** Verified: `ey` is `status: backlog` with no `sprint:` stamp, so it is not a 0204 member; it composes from the backlog. The rule's home is the contract (`ey`'s scope), the templates defer (`48`'s scope), and `ey`'s captain-named `development.md` target (the live-scenario paragraph ~L135) is genuinely SUBSUMED - once the template defers, there is nothing in `development.md` for `ey` to edit. The defer's pointer target already exists in the current contract (`first-officer-shared-core.md` "Prefer a code gate over a prose-only rule"; `ensign-shared-core.md` "Prove by exercising, not by re-reading"), so AC-1 composes whether or not `ey` ships first. The single live hazard is concurrent editing of the `development.md` live-scenario prose; `48` names this as the failure-mode-to-avoid and routes to a read-ey-first reconciliation in the dispatch checklist. That enforcement lives in dispatch, outside the body, so it is a real-but-delegated risk the package must carry forward.
- **0q's A/B fork resolves to A** (bare `--json` drops `source` too), and the gate should be told the in-repo blast radius is nil so the decision rests purely on external-consumer policy. No documented JSON-key stability contract guards the bare-`--json` `source` key.
- **48's development-only live-test is acceptable cost-scope**, not a proof gap: the boot-path governance mechanism is identical across all three template shapes, so proving it once on the highest-risk shape (the only one with a paragraph-level removal) plus structural coverage on the other two is the right bill. Three live commissions would re-prove the same mechanism at 3x spend.
- **No cross-member file conflict.** `6r` touches a new native subcommand plus `first-officer-shared-core.md` / the dev README; `0q` touches `internal/status` Go plus goldens; `48` touches `skills/commission/references/templates/*.md`. The only shared surface is `first-officer-shared-core.md`: `6r`'s AC-6 edit (a triage-read instruction) and the j7 reason-correction follow-up (the :218 bullet) both land there. They are distinct lines, but the package must sequence them so neither clobbers the other, and so `6r`'s edit does not re-introduce a prose-grep guard into a contested file.

## DoD coverage

The `index.md` Definition of Done has three bullets. With j7 dissolved:

- **Bullet 1 (the `status` read helper backbone)** - covered by `e6a` (`status-section-reader`, already in implementation): FM + section-heading→offset map, Go tests over real helper output plus a live section-read exercise.
- **Bullet 2 (the other read/mutation-cost reductions)** - covered by `6r`, `0q`, `48`, each meeting its ACs, PLUS the explicit escape clause: "where a spike proves a sink is not tool-fixable (e.g. the read-then-`set` echo if it turns out harness-inherent), it is recorded as a roadmap decision rather than forced into code." j7 is the named example of exactly this branch.
- **Bullet 3 (cut after a clean pre-cut antipattern audit)** - a drive/close step, unaffected by membership.

**j7's removal opens no DoD hole.** The DoD anticipated this disposition in its own text and routed it to a roadmap decision rather than a code deliverable. The remaining set (`6r`, `0q`, `48`) plus the `e6a` backbone satisfies the DoD; j7's roadmap-decision record is the DoD-sanctioned satisfaction of its slice of bullet 2, not a gap. The only thing j7 leaves behind for a follow-up is the contract reason-correction, which is correctly out of the sprint's code scope and tracked inside the roadmap-decision record.

## Verdict

**Gates-ready after minor folds.** No material design defect blocks the gates; the read-side helper, the SOURCE-render trim, and the template defer are sound, and j7 terminalizes cleanly. Close these before the gates lock, owner-tagged:

1. **`6r` AC-6: rebuild to the `e6a` live-drive bar** - require an FO/ensign actually fetching a transcript and running the summarizer at a real triage site, proof being that trace, never the instruction text; and own LOCATING or creating that site (grep shows it does not exist today). *ideation-rework.*
2. **`0q` parity-divergence: extend the oracle-parity pull-out to all four diverging cases** (default, archived, where-status, all-fields) and widen the JSON test plan to `seq-archived.json` + `seq-where.json`. *ideation-rework.*
3. **`48` AC-2: pin the single-source marker set in the body, or explicitly hand it to the detached audit to attack.** *ideation-rework, or captain-ratify-at-gate if deferred to the audit.*
4. **`0q` fork decision (direction A) - captain ratifies**, told the in-repo blast radius is nil so the call rests on external-consumer policy. *captain-ratify-at-gate.*
5. **`48` development-only descope - captain ratifies**, with the condition that the no-spend dry-check runs Phase-1 completion for experiment and refinement too. *captain-ratify-at-gate.*
6. **`6r` AC-5 stdin header normalization** and the package's sequencing of the two `first-officer-shared-core.md` edits (`6r` AC-6 + the j7 :218 follow-up) after the known README-slim doc-contract red. *captain-ratify-at-gate (implementer notes in the package).*

Both shipped-scaffolding surfaces (`6r`'s doc/skill edit, `48`'s templates) are owed a detached adversarial audit at validation per the lifecycle checklist.

## Provenance

Per-member design reviews (independent, refute-not-bless): the reviewer findings folded above. j7 spike re-verification: an independent verifier built the binary and ran four live `claude -p` stream-json drives on the pinned 2.1.177, closing the spike's missing positive-control gap. Entity bodies under `docs/dev/.spacedock-state/`. Stale-site and composition facts verified against the live `skills/` surface this session.
