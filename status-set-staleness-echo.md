---
id: j7tbbka0fpd9twbqgr7bhrqd
title: Eliminate the Read-then-`status --set` staleness echo — mutating a just-read file re-emits the whole file as cache-write tokens
status: ideation
source: FO + 0.20.4 scope survey (2026-06-14, this session) — on Claude Code, a `Read` of an entity file followed by a `status --set` mutation of the same file trips the harness file-staleness safety net, re-emitting the whole file as cache-write tokens. A recurring FO/ensign tax (validate-then-set, dispatch-then-set). Cited in e6a's Problem section and the sprint friction notes. Distinct from e6a (which avoids the whole-file READ); this attacks the whole-file WRITE-echo on mutation. 0.20.4 read-cost theme.
started: 2026-06-15T05:19:11Z
completed:
verdict:
score: 0.33
worktree:
issue:
sprint: 0204-structured-reads
sprint-readiness: ready
---

A `Read` of an entity file followed by a `status --set` on the same file re-emits the whole file as cache-write tokens (the Claude Code file-staleness safety net). The FO contract works around it today (Grep-over-Read; trust `--set` stdout), but the echo still fires whenever a Read before a set is unavoidable. Kill the echo, or prove it is harness-inherent and document the avoidance.

## Problem

When an agent `Read`s a file and then mutates it via Bash (including `spacedock status --set`), Claude Code's file-staleness safety net re-emits the whole file as cache-write tokens on the next turn. For a 280-line entity that is a whole-file tax on every read-then-mutate. The FO contract (`first-officer-shared-core.md:218`) already names this and routes around it, but the workaround is avoidance, not a fix.

## Proposed approach

**The spike falsified the premise. On Claude Code 2.1.177 there is no whole-file echo to eliminate.** A controlled, fully-observed measurement of the exact `Read -> status --set -> next turn` pattern shows the harness keeps the original Read tool_result cached and does NOT re-emit the file. The post-mutation turn writes ~80-300 cache-creation tokens (the Bash result string + the reply), never the ~16k tokens the 280-line file would cost to re-emit. No `<system-reminder>` and no "modified since read" reminder fire on the Bash path.

### Spike (the riskiest unknown, exercised first)

Oracle = `cache_creation_input_tokens` from `claude -p --output-format stream-json --verbose` (the harness's own usage records — an external token-count oracle, never prose). A 280-line / ~16k-token entity, identical prompt across arms, only the middle Bash command varied:

| arm | middle Bash step | post-mutation turn `cache_creation` | staleness reminder |
|-----|------------------|-------------------------------------|--------------------|
| treatment — adjacent | Read, then `status --set` mutating the read file | **244** | none |
| control — no mutation | Read, then `cat file >/dev/null` | 167 | none |
| control — other file | Read, then `status --set` on a *different* file | 178 | none |
| treatment — sealed (3 intervening turns) | Read, seal cache, then `status --set` on the read file | **143** | none |
| treatment — in-place same-inode rewrite | Read, seal, then truncate+rewrite (keeps inode, bumps mtime) | **186** | none |

The treatment-minus-noop delta (244 − 167 = 77) is fully explained by the longer Bash result string (`score: 0.33 -> 0.50`), not a re-emit. `cache_read_input_tokens` climbs monotonically across the mutation (the file content keeps being served FROM cache), proving the read tool_result is never invalidated or re-written. The sealed-cache and same-inode variants are the strongest forms of the claimed problem (the file sits behind a sealed cache breakpoint; the mutation keeps the inode and bumps mtime) and still show no echo.

What the "modified since read" net actually guards: it BLOCKS the **Edit/Write** tools with a `<tool_use_error>` (the agent must re-Read first), observed in real transcripts. It does NOT manifest as a Bash-mutation echo. The two manifestations are distinct; only the Edit/Write block is real on this harness.

### Outcome: roadmap decision, not a code deliverable

Per the dev-workflow rule (a determination with nothing shippable belongs in the roadmap, not this queue), this task records a roadmap decision carrying the measured evidence:

1. The Read-then-`status --set` whole-file cache-write echo does not reproduce on Claude Code 2.1.177 (5 measured arms, all clean).
2. The FO contract's avoidance (`first-officer-shared-core.md:218` — "Grep over Read; trust `--set` stdout") therefore defends against a behavior this harness does not exhibit on the Bash path. The contract line is not wrong to prefer Grep (it is cheaper regardless), but its stated *reason* — the staleness echo — no longer holds for Bash mutations. The roadmap decision flags the line for a reason-update in a follow-up (out of scope here; touching the live contract is its own change).
3. No `status --set` code change ships: `atomicWrite` (write-temp + `os.Rename`, `mutate.go:244`) is already in place and is not the lever — even a same-inode in-place rewrite produced no echo, so the write strategy is not what the (absent) echo keys on.

The decision IS the satisfied outcome: AC-1's "or the task is recorded as a roadmap decision with the measured reason it cannot" branch, where the measured reason is "the echo does not fire on this harness," carried by the token table above.

## Out of scope

- e6a's read-side section helper (separate task; e6a avoids the whole-file READ, this task was about the WRITE-echo).
- Non-`status` mutations.
- Editing the live FO contract line (`first-officer-shared-core.md:218`) to correct its stated reason — the roadmap decision flags it; the actual contract edit is a separate change against the live skill surface, not an ideation deliverable.
- Behavior on Claude Code versions other than 2.1.177 and on Codex/Pi runtimes — the spike is Claude-Code-specific (the echo claim is Claude-Code-specific). A future harness regression that reintroduces the echo would reopen this, which is why the roadmap decision pins the measured version.

## Acceptance criteria

Each AC names a property of the finished outcome, not a stage action, and how it is verified.

**AC-1 — A read-then-`status --set` on the same file does not re-emit the whole file as cache-write tokens (or the task is recorded as a roadmap decision with the measured reason it cannot).**
Satisfied via the second branch: the controlled spike measured that a read-then-`status --set` already does NOT re-emit the file (244 cache-creation tokens post-mutation vs 167 for the no-op control; a ~16k re-emit would dwarf both). The recorded measurement (the token table in Proposed approach) + the roadmap decision IS the satisfied outcome. No code change can "fix" a behavior that does not occur.

**AC-2 — `status --set` still narrates its mutation (`field: old -> new`) so the FO needs no re-read.**
Already satisfied by shipped code: `handlers.go:306` emits `fmt.Fprintf(stdout, "%s: %s -> %s\n", field, oldValue, val)`, pinned by the AC-3 mutation-parity test (`internal/status/mutation_test.go`). Verified live during the spike — `status --set entity score=0.50` printed `score: 0.33 -> 0.50`. No change needed; this AC was never at risk under the no-code outcome.

## Test plan

The spike WAS the test, and its outcome is the no-code branch, so there is no implementation test to write — there is no code change to prove. The proof is the measurement, reproducible by a teammate on a fresh setup:

- **Reproduce the spike (the AC-1 proof).** Build the binary (`go build -o /tmp/sd/spacedock ./cmd/spacedock`), drop a ~280-line entity + a minimal `README.md` with a `## Stages` block, then run `claude -p` headless across the three+ arms with `--output-format stream-json --verbose --allowedTools Read Bash --permission-mode bypassPermissions`, scripting a deterministic Read -> Bash -> reply sequence. Parse `cache_creation_input_tokens` off each assistant turn's `usage`. Expected: post-mutation `cache_creation` stays in the low hundreds, never approaching the file's read size; no "modified since read" string in the stream. Cost: ~5 headless runs, a few minutes. No fixtures committed — the spike is a one-off reproduction, declared here so the result is auditable, not a standing test.
- **AC-2 is already covered** by `internal/status/mutation_test.go` (no new test).
- **Why no Go test, no live-workflow test:** both would assert against code that this task does not change. The dev-workflow rule routes a nothing-shipped determination to the roadmap; the deliverable is the recorded decision, whose proof is the token measurement above, not a passing build.

## Stage Report: ideation

- DONE: Spike the Read -> `status --set` -> next-turn echo FIRST and record a measured before/after token trace locating where it originates
  Controlled `claude -p` stream-json spike, 5 arms (adjacent / sealed-cache / in-place same-inode / no-op control / other-file control); oracle = `cache_creation_input_tokens`. Post-mutation turn = 143-244 tokens vs a ~16k whole-file re-emit; no staleness reminder fired. Trace + table in `## Proposed approach`. Origin: harness keeps the original Read tool_result cached and never re-emits on the Bash path — there is no echo to control.
- DONE: On the spike evidence resolve the bimodal fork — code fix OR recorded roadmap decision carrying the measured harness-inherence reason
  Resolved to a THIRD outcome: the premise does not reproduce on Claude Code 2.1.177, so neither a code fix (nothing to fix) nor the "harness-inherent echo" framing applies. Per the dev-workflow nothing-shipped rule, recorded as a roadmap decision; the measured reason is "the echo does not fire on this harness," carried by the token table. AC-1 satisfied via its roadmap-decision branch; AC-2 already shipped (`handlers.go:306`, `mutation_test.go`), verified live (`score: 0.33 -> 0.50`).

### Summary

The spike falsified the task's premise: a `Read` followed by `status --set` of the same file does NOT re-emit the file as cache-write tokens on Claude Code 2.1.177 (measured across 5 arms including the strongest sealed-cache and same-inode variants). The harness serves the original Read tool_result from cache and re-injects nothing on the Bash path; the "modified since read" net only blocks the Edit/Write tools, never echoes a Bash mutation. Outcome is a roadmap decision (nothing shippable — `atomicWrite` is already in place and is not the lever), flagging the FO contract line `first-officer-shared-core.md:218` for a reason-update in a follow-up. Strange things are afoot at the Circle K: this contradicts both the entity's stated problem and the live FO contract, so the FO/captain should weigh whether the contract's stated *reason* needs correcting (the Grep-over-Read preference itself is still fine — it is cheaper regardless).
