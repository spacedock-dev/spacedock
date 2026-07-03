---
title: "Live codex-live rejection-flow flake — addressable-worker characterized ABSENT but transcript has reuse send_input"
status: done
source: "0240/0qt PR #446 CI (2026-06-30): codex-live failed TestLiveCodexSharedScenarios/rejection-flow twice at codex_live_runner_test.go:40 — 'Codex addressable-worker was characterized ABSENT, but transcript contains turn-starting reuse tool send_input.' The scenario SUCCEEDS (both feedback cycles run); only the worker-characterization assertion fails. Intermittent (main green on the same base e3f85ec3). The handoff-named 'codex rejection-flow' flake; forced an admin-override merge of 0qt."
id: tt4sh23s0j8rxe74tx2r14xv
verdict: rejected
completed: 2026-07-03T02:16:57Z
archived: 2026-07-03T02:16:57Z
---
`TestLiveCodexSharedScenarios/rejection-flow` (`internal/ensigncycle`, `codex_live_runner_test.go:40`) characterizes the Codex addressable-worker capability from the run transcript. It intermittently asserts the worker was characterized ABSENT while the transcript contains a turn-starting reuse tool (`send_input`) — a characterization mismatch — even though the scenario itself completes correctly (cycle-1 REJECTED, rework applies the `shared-rejection-fix: applied` marker, cycle-2 PASSED, both workers closed).

It recurred on both PR #446 runs, but `origin/main` was green on the same base (`e3f85ec3`), so it is intermittent (depends on whether codex emits `send_input` that run), not deterministic — the handoff's named "codex rejection-flow" flake.

Impact: it intermittently blocks the codex live lane (it forced the 0qt admin-override merge). Likely fix: reconcile the characterization assertion with the actual reuse-tool behavior (the worker IS reusable when the transcript shows `send_input`), or correct the scenario's characterization expectation. Surfaced during the 0240 sprint's 0qt merge; filed at captain request.
