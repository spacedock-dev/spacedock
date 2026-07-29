# Validation gate review: provider-neutral gate preparation

## Decision requested

Approve the validated s4 implementation at code commit `ee2cab74` for PR update, CI, and terminal merge handling.

## Capability delivered

`spacedock gate prepare` turns the First Officer's question, Markdown review Artifact, exact summary, and selected References into one open recorder-ready gate room. The binary derives the request, canonical Briefing, identifiers, digests, Git-root locators, room path, and binding. Immediately after preparation the room contains exactly two authoritative regular files and no copied source payloads.

The agent-facing journey is:

1. Commit the selected review sources.
2. Run `gate prepare`.
3. Commit the emitted room and binding.
4. Present the emitted room through chat or one selected room-only provider.
5. Record the direct decision with `gate record`.
6. Commit, consume, and commit the resulting stage transition.

The First Officer authors no JSON and reconstructs no room authority. Provider evidence may appear only beneath the prepared room after presentation. Recording recomputes the request, Briefing, Result, and inventory pins and writes no `association.json`.

## Reviewed snapshot and scope

- Code commit: `ee2cab74`
- Base: `origin/main` at `4738d583`
- Current diff: 40 files, `+4017/-1076`
- Recovery reference: `recovery/prepare-provider-neutral-gate-room-004368de`
- PR and CI: not yet updated after reconstruction

The correction-run expansion was removed. The retained later changes are task-owned integrity protections: Git path identity, prepared-room symlink containment, successor authority validation, exact provider Resolution equality, two-file room enforcement, arbitrary Briefing locators, dynamic Briefing identity in the real journey, and failed-prepare authority preservation.

## Validation evidence

Independent validation cycle 32 reproduced all five current acceptance criteria at `ee2cab74`.

- Folder and flat tasks prepare, commit, archive, and roll back without copied sources or sibling capture.
- Selected sources reopen from movable local Git objects with full commit, path, and raw SHA identity.
- Arbitrary canonical Briefing locators, exact Unicode summaries, recursive duplicate refusal, validation, eligibility, and consumption all pass.
- The real-binary journey observes prepare, bind commit, dynamic presentation, decision, close commit, consume, consumed commit, and successor authorization.
- Provider recording recomputes four retained pins, permits post-prepare provider evidence, and rejects tampering without an association artifact.
- Four adversarial mutations failed at their intended boundaries.
- `go test -count=1 ./...`, `go test -count=1 -race ./...`, formatting, diff checks, strict MkDocs, and the isolated real-binary replay passed.

## Findings

Material outcome defects: none.

Material evidence defects: none.

Polish findings: none.

Deferred risk: source and authority readers allocate whole files. The trigger is unusually large selected sources or authority JSON, outside the current small local review-artifact journey. Promote this risk if large binary References become supported or representative memory/latency probes regress.

## Recommendation

Approve. The supported v1 journey is independently demonstrated, the review-runaway machinery has been removed, and no material finding remains.
