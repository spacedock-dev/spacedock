# Haiku loop premise-spike — per-run drive streams

Durable tool-call streams from the live N=3 Haiku mechanical-loop drive
(`TestLiveHaikuLoopSpikeN`, ran 2026-06-19). Each `run-N-stream.jsonl` is the
stream-json transcript of one bare `claude --model haiku -p` FO driving the
throwaway split-root fixture through `boot → next → dispatch → advance(→integration)
→ terminalize(→done)+archive`.

These are EVIDENCE, not the grade. The README proof policy grades on durable
on-disk state (terminal frontmatter + committed transitions), which the harness
reads off the fixture git log; the fixture itself lives in `t.TempDir()` and is
cleaned up after the run. These streams are the surviving per-run artifact the
deviation-class table cites.

Re-grade from a stream: parse the `tool_use` Bash/Agent blocks for the verb
sequence and the absence of any `Agent(model=opus|sonnet)` call.

Per-run summary:
- run-1: single clean drive; one bare worker Agent (no subagent_type).
- run-2: first worker Agent used `subagent_type: spacedock:ensign` which does not
  exist in the bare (no-plugin) session and FAILED ("Agent type 'spacedock:ensign'
  not found"); the FO recovered by re-dispatching a plain `claude` worker, which
  did the work. Durable end-state still held (marker landed, terminal reached, no
  shortcut). One trailing `status --boot` after archive confirmed `dispatchable:[]`.
- run-3: single clean drive; one bare worker Agent (no subagent_type).

All three runs graded every durable field true, including
`integrationTransitionCommitted` (no implementation→done shortcut) and
`noStrongerModelAgent` (no opus/sonnet Agent touched the loop).
