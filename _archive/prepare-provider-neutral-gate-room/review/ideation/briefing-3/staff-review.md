# Independent staff review — s4 ideation cycles 6–9

Verdict: **REVISE**.

## Bottom line

The core correction is sound and Captain-directed: address each selected source as
`git-root://<root>/<full-commit>/<repo-path>`, verify its raw SHA-256, keep the room to
the request plus canonical Briefing, carry one caller-authored primary-Artifact
`summary`, and call the result recorder-ready rather than presentation-ready.

The current package is not ready for an ideation approval. It consumes the cycle-5
approval for a different architecture, assumes a pre-rebase ownership layout that
landed `main` contradicts, adds an ordinary-ref policy that the address/reopen proof
does not require, and calls a stale 19-file/1,694-LOC estimate the smallest surface.

## Material 1 — the only approval is for the rejected frozen-copy design

The cycle-5 Briefing asks whether s4 should freeze every selected source inside the
retained room. Its gate review says every Artifact and Reference lives under
`room/sources/`, explicitly says that this adds no Git-address schema, and recommends
the 16-file frozen-source design. The cycle-5 staff review approves that exact direction
and calls `+1,090/-161` proportionate.

Cycles 6–9 materially supersede it:

- Cycle 6 removes every copied payload and introduces a new Git-root URI and production
  Git resolver.
- Cycle 7 changes the readiness claim and provider boundary and files a separate
  materialization owner.
- Cycle 8 adds the required `--summary` interface and canonical Artifact field.
- Cycle 9 expands compatibility and exact-display proof.

The entity itself records the governing instruction twice: do not consume the pending
approval; a corrected Briefing must supersede it (`index.md:863-876`). Nevertheless,
attempt 2 still resolves “Frozen room-owned sources…” and its application is now
`state: consumed` (`index.md:63-82`). The state log shows `53cf909c` consumed that
application, `19cfad70` created a worktree, and `f29d8ca6` reset status/worktree without
un-consuming the stale authorization. The landed binary reports the resulting mismatch
exactly: `status=ideation`, `gate-state=closed`,
`gate-application-state=consumed`, `gate-condition=consumed`,
`gate-eligible=false`.

Required correction: retain attempt 2 as history, bind a fresh attempt-3 Briefing over
the corrected design, and obtain a new ideation resolution. Do not reinterpret or
re-consume attempt 2.

## Material 2 — the design's baseline reset trigger has fired

The design is explicitly based on pre-xb-rebase 6y commits `60adfc1f`/`e9415a17`; it
says that changed lifecycle ownership, recorder commands, authority capture, request
fields, assertions, or surface require a return to ideation (`index.md:474-503`,
`532-543`). Current `main` is `4ff98d8c`, with landed 6y at `deac7f8a`.

The mismatch is behavioral, not just a SHA change. The design says
`present-gate` is rendering-only and unchanged while `fo-gate-lifecycle` owns the
selected-override halt (`index.md:474-482`, `505-515`, `647-675`). On current `main`,
`skills/present-gate/SKILL.md` owns override probing, chat fallback, room handoff,
provider retention, and room recording; `fo-gate-lifecycle` still tells the FO to
assemble and bind `ROOM/briefing.json`. The declared 19-file surface excludes
`present-gate`, so it cannot produce its own claimed owner split.

The smaller repair is to rebase the design on current `main` and choose one owner. If
s4 remains recorder-only, keep its shipped behavior to prepare/bind/validate and let
`git-root-review-v1-materialization` own the provider-channel halt and later handoff.
If s4 itself changes selected-override routing, include the actual `present-gate`
contract and re-scope the task honestly. The present hybrid is not implementable.

## Material 3 — ordinary-ref containment is an invented retention policy

The Captain-directed identity is supported by the spike: logical root, full immutable
commit, repository-relative path, and raw SHA recover exact bytes after both checkouts
and worktree paths move (`index.md:196-228`, `863-891`). Dirty, untracked, and
third-repository refusal also follows the recorded staff direction.

The additional requirement that `git for-each-ref --contains` name a local,
remote-tracking, or tag ref does not follow from that address. It rejects reflog-only
and detached commits even when the selected checkout is clean and
`<commit>:<path>` resolves exactly (`index.md:257-272`, `618-623`), while still creating
no retention ref or cache and admitting that deletion of the last ref can later make
the room fail.

An isolated Git probe demonstrated the distinction. A clean detached checkout at
commit `019f1078b203a6d4acbd29c580bed1bd6476b856` had no
`for-each-ref --contains` result, yet
`git cat-file blob <commit>:f.md` returned the exact selected bytes. The proposed rule
would reject a source the resolver can address and verify; conversely, observing a ref
at prepare time cannot guarantee that the ref remains later.

Keep the actual local-object contract: classify `main`/`state`, resolve the full
commit/path from the current root's object database, compare selected worktree bytes at
prepare, verify raw SHA, and fail closed without fetch or worktree fallback whenever a
later operation cannot reopen it. Remove the ordinary-ref/reflog policy and its
shallow/prune matrix unless the Captain explicitly chooses a stronger retention
contract and assigns the mechanism that enforces it.

## Material 4 — 19 files / 1,694 LOC is not a current smallest-surface estimate

The estimate is computed against the pre-rebase composition even though the reset
condition has occurred. The diff from `60adfc1f` to current `main` spans 52 files and
changes the exact lifecycle, recorder, application, contract, and live-test surfaces
s4 plans to edit.

The corrected close-out audit reinforces the need to rebaseline rather than stack more
machinery onto the provisional table:

- it finds 13 dead `application` leaves and calls for their pre-v1 removal;
- it proves the ephemeral association is derived and then re-verified against itself,
  with a one-file net-75-LOC simplification;
- it proves the lifecycle command-text mutant test is a tautological prose grep and
  routes it for deletion; and
- it requires the remaining sprint members, including s4, to receive a fresh detached
  audit and live-lane evidence on the actual tag candidate.

s4 currently budgets edits to `application.go`, `operation.go`, and
`recorded_gate_lifecycle_test.go` without composing those corrections, and about 345
new `internal/gitsource` lines include the unproven ref-retention policy. The 21-file /
2,118-LOC cap therefore measures neither current `main` nor the smallest corrected
implementation.

Recompute the file/LOC surface after resolving Materials 2–3 and the close-out cleanup
baseline. Preserve the proof that can fail:

- arbitrary Briefing locator through bind, room record, validation, and eligibility;
- production Git-root reopen after independent checkout/worktree movement, with raw-SHA
  mismatch and no-fetch controls;
- exact caller-authored primary summary, summary-free References, repeated-flag and
  invalid-UTF-8 refusal, plus request-less/advisory compatibility controls;
- recursive duplicate-member rejection with byte-clean mutation failure;
- two-file room/error atomicity; and
- existing host live journeys only for shipped lifecycle behavior actually changed.

Do not retain ref/reflog/prune lanes for an unchosen retention policy, selected-provider
fixtures for behavior owned by the sibling, or prose-token checks as behavioral proof.

## What survives review

- **Git-root locator: APPROVE direction.** Root + full commit + repository path plus raw
  SHA is the smallest Captain-directed durable identity; no object-format segment,
  remote URL, copied payload, or generic URI framework is needed.
- **Canonical summary: APPROVE direction.** One required `--summary` passed unchanged
  into the primary Artifact's existing Review v1 extra-field model is the smallest way
  to avoid both handcrafted JSON and invented prose. References remain summary-free.
- **Recorder/provider split: APPROVE direction.** The inspected Subspace working copy
  still opens Artifact URIs as filesystem paths, rejects `://` References, preserves
  Artifact extras, and reads only `Extra["label"]` for Artifact chrome. It neither
  resolves Git-root content nor displays `summary`; the sibling materialization and
  exact-display E2E are real release dependencies.
- **Proof spine: APPROVE after narrowing.** The moved-root production resolver,
  arbitrary-basename differential, exact-summary controls, duplicate-JSON adversary,
  byte-clean failure checks, repository gates, and final detached/live evidence are
  proportionate.

## Gate recommendation

REVISE now. Return with a current-main owner map, the minimal local-object policy, a
recomputed surface/proof table, and a fresh attempt-3 Briefing that explicitly
supersedes the cycle-5 frozen-copy approval.
