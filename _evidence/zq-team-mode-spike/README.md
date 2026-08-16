# Team-mode invocation spike — rejection-flow (zqb683j8jth0tyr2eme231e2, ideation)

Two live runs of `TestLiveCommonRejectionFlow` on the composed tree (571017df3 = PR #718 + #719)
with the ideation-proposed team-mode prompt, one per runtime, 2026-08-16. Run in a throwaway
detached worktree; nothing was written to the stack branch.

What each file is:

- `codex-topology.txt` — every `function_call` in the parent rollout, in stream order, with the
  correlated `task_name` handles and each worker's `Done:` message. This is the evidence that
  `followup_task` carried both the rework and the re-review against the two spawned handles.
- `claude-topology.txt` — every `Agent`, `SendMessage`, `task_notification`, and load-bearing
  `Bash` call, in stream order. Shows four fresh dispatches, each preceded by a
  `dispatch context-budget` probe that exited 1 and a `shutdown_request` to the superseded worker.
- `*-entity-headings.txt` — the durable heading shape each run left. Codex wrote
  `## Stage Report: implementation (cycle 2)`; Claude wrote an exact duplicate heading. The
  fixture's rework sentence, not model variance, is the difference.
- `*-verdict.txt` — the graded outcome. Codex red on `rejection-gate-not-prepared` only (FO
  residual mode 1, recurring under team mode); Claude green.

Reproduce: check out 571017df3 into a scratch worktree, apply the invocation change from the
entity's `## Proposed approach` item 1, persist the Codex native rollout next to the public stream,
and run the journey once per runtime with the CI codex shim on PATH. The raw streams (135 KB public,
444 KB rollout per run) were not committed; implementation should capture its own from the AC-1 loop
and land the distilled fixtures under `internal/ensigncycle/testdata/`.
