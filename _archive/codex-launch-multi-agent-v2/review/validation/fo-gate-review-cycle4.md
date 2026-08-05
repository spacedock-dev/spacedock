# Codex launcher guarantees multi-agent v2 — validation cycle 4

## Chosen direction

Approve the cycle-4 implementation and validation. The candidate binds `parent_thread_id` and `agent_path` atomically from one typed child session record, closing M5 without changing the approved launch surface.

## Evidence

- Four validation checklist items are DONE.
- AC-1 through AC-4 have independent behavioral evidence.
- The focused lifecycle/config checks, `go test ./...`, `go test ./... -race`, `gofmt -d ./cmd ./internal`, and `git diff --check` passed.
- The isolated-home live lifecycle passed in 49.99 seconds with the exact `[features.multi_agent_v2]` table: 16 threads, `agents` namespace, and visible spawn metadata.
- A detached audit rejected the spoofed typed-child identity at commit `155747c2fc`; no material or evidence finding remains.
- Candidate scope remains 8 files, +379/-26 against `origin/main`.

## Recommendation

Approve to advance this validation gate to `done`; terminal merge remains subject to the normal PR gate and Captain approval.
