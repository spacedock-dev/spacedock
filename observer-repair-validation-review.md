# Validation gate — repair-pi-worker-lifecycle-observation

## Outcome

Validated, with its live-lane clause pending. This is the strongest evidence in the stack.

## Candidate

Layer 2 of the stack. Its three commits credit `bg_wait` and the native `subagent-notify` completion in the
shared worker-lifecycle assert, commit a hash-pinned replay fixture with provenance, and remove the
auto-continue Pi binding because the assert now covers it.

## Evidence, with falsifiers

- Removing **both** new credits reds the replay: `validation lifecycle incomplete: spawns=1 completed=-1`.
- Removing **either** credit singly also reds it.
- Removing only the run-id correlation reds it with `a native completion notice for another run credited the
  spawned run`.
- The ordering control rejects, and removing the ordering clause reds it.
- The fixture and its `provenance.json` are committed; the provenance source hash matches the retained source.
- `go test ./internal/ensigncycle/...` passes. Claude and Codex grading branches are unchanged.

## Pending and caveats

- **The live lane clause is pending.** There is no green lane at this tip yet.
- The fixture's own hash is not asserted at runtime; pinning is the provenance source hash plus git.
- One stated mutation overstates the spawn guard: removing `spawns < 1` alone does not red the zero-spawn
  control, because run-id correlation independently gates it. The behaviour holds; the wording was loose.

## Question

Do you accept the validation of this layer? This approves the validation only, not any merge.
