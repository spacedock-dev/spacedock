# Independent staff review — s4 ideation cycle 10

Verdict: **REVISE**.

## Material

### 1. The planning estimate is still an implementation gate

Cycle 10 correctly rebases the named-file table on `4ff98d8c`, but immediately turns
the estimate into a `+2 files and +20% changed LOC` tolerance and a hard
`28 files / 2,285 changed LOC` cap. Neither number follows from the value, authority,
safety, compatibility, or dependency contract. The wording can reject a semantically
correct helper/fixture split or reward an implementation that stays under the count
while crossing an authority boundary.

This directly conflicts with the Captain-approved instruction that file/LOC tables are
planning estimates and drift evidence only. It also leaves s4 inconsistent with rq
cycle 5, which already corrected the same issue after independent review.

**Required correction:** keep the named-file table and expected deltas as
non-authoritative planning estimates, but delete the tolerance, hard cap, and every
count-based authorization/rejection/reset implication. Before product edits, require
the worker to declare intended path/LOC scope and explain known deviations from the
table. Before implementation review, reconcile path-scoped `git diff --stat` and
`git diff --numstat`; reviewers judge the semantic explanation, never the count alone.
No test or gate may assert repository file totals, LOC totals, or percentage variance.
Return to ideation only for a semantic expansion such as another public argument,
caller-selected authority, a third authoritative prepared metadata file, copied
sources, remote acquisition, stored association, compatibility behavior, weaker
Briefing binding, or a changed provider/validator owner.

### 2. The current-main skill contract does not yet implement the declared owner map

The proposed `fo-gate-lifecycle` wording calls `gate prepare`, but current main's
mandatory `gate --help` preflight requires only `record`, `validate`, `eligibility`,
`consume`, and record flags. Cycle 10 does not update that preflight or the matching
command-reference paragraph. A pre-s4 binary at the pinned baseline therefore passes
the advertised capability gate and fails only when the First Officer attempts the new
verb.

The override path has the inverse mismatch. Current `present-gate` tells the agent to
run a provider availability/version probe before preparing a room and to fall back to
chat when it fails. Cycle 10 says s4 never probes a provider and assigns all
materialization, capability probing, allocation, invocation, cleanup, and evidence to
rq's room-only fixed entry. Its concrete skill diff changes only the room handoff line,
so the existing agent-run probe/fallback rule survives and contradicts both the owner
map and corrected rq cycle 4.

**Required correction:** make the First Officer's one fresh capability preflight
require `prepare` and the prepare flags it will use, in both `fo-gate-lifecycle` and the
command reference; a stale help surface must halt before source commits, preparation,
or state mutation. Replace the selected-override instructions in `present-gate` with
one post-prepare, post-bind-commit handoff of the exact
`/subspace:r gate <room>` form. The agent performs no provider probe and reconstructs
no coordinate; rq's fixed entry owns provider discovery/capability and its failure
semantics. Extend the existing lifecycle/contract fixtures so removing `prepare` from
help or adding an agent-run provider probe makes them fail.

## Polish (non-material)

Cycle 10 says one recursive token-stream reader covers the future resolved-source
manifest, while its surface and owner map correctly leave that manifest to rq/Subspace.
Limit s4's shared reader claim to request, located Briefing, Result, and presented
inventory, then state separately that rq's closed manifest parser rejects duplicate
members recursively before display. This is an ownership wording repair; corrected rq
already carries the downstream behavior.

## Direction that stands

The folder and flat forms now converge on the collision-free
`<root>/<slug>/review/<stage>/briefing-N/` room while preserving the existing
folder-wins discovery rule. Treating flat `<slug>.md` plus companion `<slug>/` as one
literal commit/archive/rollback unit is the necessary current-main durability
extension; the planned tests cover sibling isolation, tracked deletions, archive
movement, and rollback.

The two-file claim is correctly limited to the instant after preparation:
`request.json` plus its located canonical Briefing, with zero copied selected-source
payloads. Later provider-owned evidence under `provider/` does not weaken that
invariant and cannot replace either authoritative file.

The local Git-object contract is proportionate. Logical root, full commit,
repository-relative path, and raw SHA-256 locate and verify the exact local object
without an ordinary-ref or reflog policy. Detached resolvable objects remain valid;
missing objects fail closed without fetch, hydration, worktree fallback, or invented
retention machinery.

The arbitrary request locator, Briefing id/full canonical digest, exact caller-authored
primary Artifact summary, request-less/advisory compatibility controls, recursive
duplicate-member refusal, shared bind/record/validate/eligibility resolver, and
in-memory association all match the approved authority boundary. rq also correctly
receives only the room and independently binds its manifest to the unchanged canonical
Briefing id and full digest.

## Gate recommendation

Revise only the numeric-authority language and the current-main skill/preflight owner
contract above. No product implementation, provider run, compatibility layer, copied
source, retention policy, or additional schema is needed. After the corrected design
and a confirming independent re-review are committed, the First Officer may bind the
fresh attempt-3 Briefing; attempts 1 and 2 remain immutable history.
