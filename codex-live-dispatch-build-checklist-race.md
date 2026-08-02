---
id: 0gm98sm87730hq8bapqysxew
title: Codex live agent raced dispatch build --checklist-file ahead of writing the checklist file (2 attempts, 1 success)
status: backlog
source: "Found on PR #600 (collapse-gate-approval-ceremony) codex-live CI, 2026-08-02, run 30754109029, job 91513297838: TestLiveCodexSharedScenarios/recorded-gate-lifecycle failed 'successor dispatch build attempts/successes = 2/1, want 1/1'. Pulled the runtime-live-e2e-codex-live artifact's codex-exec.jsonl directly: the live codex agent ran `spacedock dispatch build --workflow-dir . --entity-path .spacedock-state/recorded-gate-task/index.md --stage handoff --checklist-file /tmp/recorded-gate-task-handoff-checklist --bare-mode --stamp` before the checklist file existed (exit 1: 'failed to read checklist file ... no such file or directory'), then created the file and re-ran the identical command successfully (exit 0). Not caused by collapse-gate-approval-ceremony's changes -- --stamp is new from that entity, but the checklist-file read path and the dispatch.builds==1/successfulBuilds==1 assertion are both pre-existing and untouched. This is the first live run to get far enough (past two earlier, now-fixed, unrelated CI failures on the same PR) to actually exercise this exact ordering. Captain directed: treat as a candidate flake, file for codex diagnosis, do not block the merge on it -- re-run to confirm green."
started:
completed:
verdict:
score: 0.4
worktree:
issue:
---

The live Codex agent, following the ensign dispatch-build instructions, issued `dispatch build --checklist-file <path> --stamp` referencing a checklist file it had not yet written, got a benign exit-1 "no such file" error, then correctly wrote the file and retried -- succeeding on the second attempt. `TestRecordedGateLifecycleAC7ResumeMatrix`'s sibling assertion (`dispatch.builds != 1 || successfulBuilds != 1`) requires exactly one attempt and treats any retry as a failure, so this reads as a hard CI failure even though the agent self-corrected and the entity's actual outcome (successor dispatched, marker recorded, durable commit made) was unaffected.

## Open questions for diagnosis

- Is this a one-off live-model ordering slip (the agent read/composed the dispatch command before finishing the checklist-file write in the same turn), or does something in the current ensign/dispatch-build instructions (skill prose, checklist-file example ordering, `--stamp`'s own docs) invite writing the command before the file exists?
- Is `dispatch.builds/successfulBuilds == 1/1` the right bar for a *live* agent scenario, where a benign self-correcting retry is expected/normal agentic behavior, versus a *scripted* CLI-replay fixture where exactly one command is unambiguously correct? Consider whether the live assertion should tolerate one recovered retry while the scripted-replay tests keep the strict 1/1 bar.
- Reproduce on at least one more codex-live run (or a targeted local repro) before concluding it's non-deterministic rather than a systematic codex behavior under this exact prompt shape.

## Out of scope (for now)

Any code change to `dispatch build`'s checklist-file handling or to the assertion itself -- this entity is for diagnosis first; a fix (if warranted) is a follow-up decision once the pattern is understood, not assumed here.
