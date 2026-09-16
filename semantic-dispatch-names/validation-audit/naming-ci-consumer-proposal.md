# Consolidated naming consumer proposal

FO classified all three findings Material and naming-owned. Disposition: HOLD pending captain approval of the consolidated scope. The open two-file binding review covers only its pinned proposal. No correction below is authorized yet. Candidate remains unchanged at `2a849180848c6d9d3742ef59d995330be95ea74f`; approved report, gate, and frontmatter remain unchanged.

## Exact scope

The current layer against `c7568815b` has 46 files and +428 net lines. All seven paths below are new to that manifest; no existing manifest file needs another edit. The complete proposed layer is 53 files, +428 net lines: 5 production Go files, 18 test files, 5 documentation/skill files, and 25 goldens. The 450-net ceiling is preserved; the file ceiling requires captain amendment.

| New path | Proposed net delta | Minimal correction |
|---|---:|---|
| `internal/dispatch/codex_v2_adapter.go` | -19 | Remove unsupported Slug/Stage fields and reverse-name helpers; preserve Name, tool arguments and collision refusal. |
| `internal/ensigncycle/conflict_owner_handoff_live_test.go` | +12 | Read canonical entity/status from the stamped entity; add deterministic stamped-owner tuple regression. |
| `internal/ensigncycle/pi_rejection_extractors_test.go` | 0 | Remove mandatory old prefix from dispatch-path matcher; preserve whole stem, path and stage constraints; update comment. |
| `internal/ensigncycle/pi_rejection_extractors_test_test.go` | +7 | Run existing complete fresh-chain proof with legacy and semantic names; assert exact handle and preserve negative controls. |
| `internal/dispatch/build_codex_host_test.go` | 0 | Correct the independently specified fresh artifact filename. |
| `internal/dispatch/build_input_mode_test.go` | 0 | Correct the independently specified fresh artifact filename. |
| `internal/dispatch/self_contained_assignment_test.go` | 0 | Correct the independently specified fresh artifact filename. |

The seven corrections total zero net lines. The prior two-file proposal includes only the first two rows, totaling -7 net lines.

## Finding 1: adapter identity inference

1. Released user and normal workflow: Codex first officer hands an owned conflict back after normal stamped dispatch.
2. Observable harm: the initial owner tuple loses its entity and aborts handoff. The adapter can also misrepresent a hyphenated stage or shortened slug.
3. Authority: value-ac[AC-3] lifecycle ownership must retain canonical entity and stage while preserving worker identity.
4. Trigger evidence: CI run `35058663190`, job `104674351749`; a temporary Go overlay exercised the existing stamp fixture, actual `dispatch build --stamp --host codex`, and `CodexMultiAgentV2SpawnInput`. It exited 1 in 0.666s with `{Name:conflict-owner-implementation Slug: Stage:implementation}`, while `status.EntitySlug` returned `conflict-owner`.

The adapter requires more than four dash components and discards two presumed prefix tokens. Names are not reversibly encoded identities. Only the conflict-owner fixture consumes Slug/Stage; remove those misleading fields and source the tuple from authoritative entity state. Keep every existing tuple and provenance assertion.

Minimal proof: deterministic `TestConflictOwnerStampedIdentity` invokes real stamping and adapter parsing, with literal expected entity, stage, worker name, and branch. Existing adapter isolation and sanitized-collision tests remain intact. Live-tag compilation checks integration.

## Finding 2: Pi dispatch-path decoding

1. Released user and normal workflow: Pi first officer executes the supported fresh rejection/rework/re-review chain.
2. Observable harm: canonical semantic dispatch prompts yield no observed spawn routes, so a conforming lifecycle cannot pass its grader.
3. Authority: value-ac[AC-3] new workers must retain correct lifecycle identity.
4. Trigger evidence: actual `dispatch build --stamp --host pi` emitted `Read /tmp/spacedock-dispatch/thing-implementation.md and treat its content as your assignment.` The actual Pi extractor returned no routes in a temporary overlay, exiting 1 in 0.248s. The complete chain replay returned zero of eight expected events with semantic names, while the legacy form passed. Log: `/tmp/semantic-pi-prefix-audit.log`.

Remove only the required prefix in the existing regex. Run the existing full fresh-chain proof for both namespaces, retain the exact filename stem as handle, and preserve the missing-completion and wrong-reuse-branch controls. This adds no generic name-to-entity inference.

## Finding 3: refusal-proof artifact blindness

1. Released user and normal workflow: maintainers verify that invalid dispatch input, unsupported Codex bare dispatch, and missing launcher resolution leave no artifact mutation.
2. Observable harm: three refusal assertions inspect obsolete artifact paths and therefore pass despite mutation at the current canonical path.
3. Authority: value-ac[AC-2] refusal must preserve artifacts, with a falsifiable non-mutation proof.
4. Trigger evidence: temporary overlays injected writes to `f02codexbare-implementation.md`, `thing-backlog.md`, and `unique-no-launcher-implementation.md` during their respective refusal checks. All three existing tests still passed, collectively in 0.856s. Log: `/tmp/semantic-stale-guards-audit.log`. Each injection restored previous bytes or removed its temporary artifact through cleanup.

Correct the three literal expected filenames. The same injected-write controls must then fail. No new standing test owner is needed; retain all other existing refusal assertions.

## Audit boundary and limits

The bounded repository audit inspected direct worker-name and dispatch-filename consumers, fixed prefixes, name splitting, path matching, and regex decoding. No further equivalent active consumer was found. The old `decompose` helper is used only by legacy tests. The production resolver enumerates both generations. Parity setup registers intentional historical worktrees. The merged-member prefix assertion checks a literal historical fixture. Other inspected lifecycle readers use prefix-independent handles, stage tokens or path suffixes.

The Pi completion-notification failures described in `/tmp/spacedock-tip-ci/pi-completion-diagnosis.md` are separate; this proposal does not claim them as naming-owned or fixed. Reproductions used deterministic overlays and existing fixtures, not live models. No candidate edits, candidate commits, push, CI, broad suite, native run, or reviewer rerun occurred. Proposed deltas are counted from prepared text outside the candidate; no green corrected-candidate claim is made.
