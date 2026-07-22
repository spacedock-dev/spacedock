# Independent plan review: folder-form `state commit`

## Review question

Is vn's current ideation plan ready to approve for implementation as the durable-decisions sprint's persistence companion? Identify material correctness, concurrency, scope, or proof gaps; distinguish required plan corrections from deferrable polish.

## Why this joined the sprint

`spacedock state commit <slug>` resolves a folder-form entity to `<slug>/index.md` and gives Git that single-file pathspec. It therefore omits reports, Briefings, exact Results, associations, and other workflow-record artifacts stored beside the index.

The defect reproduced twice in live work:

1. The Roborev setup entity committed only its index while two canonical artifact files remained dirty.
2. The canonical-v1 recorder dogfood added `durable-gate-approval-pending-blockers/review/validation/briefing-v1/{briefing.json,gate-review.md}`. `state commit durable-gate-approval-pending-blockers` reported a clean no-op while both files remained untracked. They required manual exact-path commit `d8e4180c`; the recorder's later index mutation committed normally as `2c616b7e`.

That makes the sprint journey leaky: the recorder can create correct durable metadata, but the normal state command cannot durably commit its retained package.

## Proposed behavior

- Resolve one commit unit from a canonical top-level entity slug.
- Flat form commits exactly `<slug>.md`.
- Folder form commits exactly `<slug>/`, including tracked modifications, deletions, and new non-ignored reports/artifacts anywhere below it.
- Never stage the checkout wholesale. Dirty flat siblings, folder siblings, and top-level files remain untouched and unstaged.
- Preserve existing push, reject/rebase, conflict-HALT, no-force, local-only, JSON/text output, and clean no-op behavior.
- Treat concurrent changes anywhere within one folder-form entity as changes to the same entity. Disjoint entities may rebase and push; same-entity conflicts halt without discarding either writer.
- Reject path-bearing or noncanonical operands before filesystem resolution. A documented entity slug must not act as a path alias for an arbitrary nested artifact.

The intended implementation is a pathspec-boundary correction in the existing state-sync command, not a new artifact registry, manifest, transaction log, or multi-entity commit mode.

## Acceptance proof in the current plan

1. A real-Git fixture changes a folder entity's index, modifies a tracked report, and creates an untracked artifact. One command must push one commit containing exactly those paths and leave the target folder clean.
2. The same fixture dirties unrelated flat/folder siblings and a top-level file; none may be staged or committed.
3. An artifact-only change must produce a commit instead of a false no-op; an immediate rerun must produce the existing clean no-op.
4. Existing flat-form, output, no-origin, retry-on-reject, and conflict-HALT behavior remains green with an exact flat committed-path assertion.
5. Two-host real-Git fixtures prove disjoint-entity success and same-folder conflict HALT, including clean rebase abort and no force/discard.
6. CLI and real-Git cases reject separators, traversal, absolute paths, `.`/`..`, and the observed nested pseudo-slug workaround without changing HEAD, index, worktree, or remote.

Required repository gates remain `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

## Boundaries

Out of scope: committing multiple entities, `git add -A` over the state checkout, changing entity discovery/archive layout, changing merge-guard semantics, moving product deliverables into workflow state, or interpreting gate package contents.

vn owns only how one folder-form entity becomes one durable Git commit. It does not interpret `gates`, create Briefings, validate Results, or change recorder semantics.

## Plan-completeness questions for the reviewer

The current entity declares behavior and six falsifiable ACs, but does not yet give an expected file/LOC surface and tolerance, an explicit riskiest-path spike determination, a mechanism-to-value comparison, or a concrete user-facing documentation diff. Determine which omissions must be repaired before approval. Also inspect whether "canonical slug" should reuse an existing discovery/parser boundary rather than mint a second validation grammar.

## References

- [Current vn entity](../../index.md)
- [Current state-sync implementation](/Users/clkao/git/spacedock-research/spacedock-v1/internal/cli/state_sync.go)
- [Existing real-Git state-commit tests](/Users/clkao/git/spacedock-research/spacedock-v1/internal/cli/state_commit_test.go)
- [Durable-decisions Commander package](/Users/clkao/git/spacedock-research/spacedock-v1/docs/roadmap/durable-decisions/dispatch-sprint-execution.md)
