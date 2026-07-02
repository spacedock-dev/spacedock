# Pi intercom runtime capability probe — ideation

Implemented the ideation update for `docs/dev/.spacedock-state/pi-intercom-runtime-capability-probe/index.md` and committed only that entity path in the state checkout.

## Findings

- Clarified that `subagents-doctor` bridge-active is setup evidence, not sufficient proof of child-to-supervisor `contact_supervisor` talkback.
- Proposed a reusable runtime capability probe pattern: recipe, durable evidence JSON, focused integration tests, docs/probe invariant tests, and live/manual smoke path.
- Specialized the first target to Pi intercom supervisor talkback with exact progress/decision/approval/marker expectations.
- Recorded that the requested Pi runtime docs were missing from the main checkout and available only in `.worktrees/spacedock-ensign-pi-runtime-support/`; inspected those worktree copies and noted the limitation in the entity.
- Did not run a live Pi intercom smoke during ideation.

## Commit

- `63129f25 ideation: pi intercom runtime capability probe`

## Validation

- `git -C docs/dev/.spacedock-state status --short` returned clean after commit.
- `git -C docs/dev/.spacedock-state show --stat --oneline --name-only HEAD -- pi-intercom-runtime-capability-probe/index.md` showed only `pi-intercom-runtime-capability-probe/index.md` in the commit.
- `git -C docs/dev/.spacedock-state status --short -- non-interactive-teardown-exit/index.md` returned no output.
