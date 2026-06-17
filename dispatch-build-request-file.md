---
title: dispatch build --request-file — single Markdown request file for FO dispatch input (schema shared with next-action)
status: backlog
source: captain proposal (2026-06-17) — reduce FO dispatch-input ceremony; symmetric to the output file-pointer handoff
score: 0.4
sprint: 0205-layered-fo
sprint-readiness: ready
id: vppqhf877vz8p9xykgdw98x6
---

The dispatch OUTPUT side already uses a file-pointer handoff: `dispatch build` writes the real assignment to /tmp/spacedock-dispatch/*.md, the worker gets a tiny "Read this and treat it as your assignment," and the stage prose is further deferred behind `### Fetch commands`. The INPUT side never got that treatment: the FO hand-creates 2-3 scratch files per dispatch (checklist-file, plus optional scope-notes-file / feedback-context-file). That is model-action + ceremony overhead (not the main token-bloat path), and it is the remaining asymmetry.

Proposal: add `--request-file FILE` to `spacedock dispatch build` — one Markdown request file the FO authors instead of 2-3 scratch files: YAML frontmatter scalars (entity_path, workflow_dir, stage, team_name, bare_mode, feedback_reflow) followed by `## Checklist`, `## Scope Notes`, `## Feedback Context` sections. Invoked as `spacedock dispatch build --request-file /tmp/spacedock-dispatch/request-foo.md`.

Why: one FO-created file instead of two or three; no JSON escaping for Markdown/backticks/shell vars; shell-quoting-safe; easy partial reads / section edits; the existing OUTPUT file-pointer handoff is unchanged; schema-v2 JSON stdin stays the programmatic path.

LOAD-BEARING CONSTRAINT (decided with the captain): the request-file schema is defined ONCE and SHARED with the `next-action` member. `--request-file` is the CONSUMER (FO-authored) now; `next-action` / `prepare` becomes the PRODUCER later — it generates the request file with the scalar fields prefilled and candidate checklist items, and the FO edits only the judgment-bearing sections (which signals belong in the checklist). Design the artifact as ONE contract so the two paths cannot drift. (next-action depends on 2y, which has MERGED — PR #385 — so next-action is now unblockable.)

Design points (ideation resolves):
- Additive: KEEP the `--checklist-file` / `--scope-notes-file` / `--feedback-context-file` flag form (scripts/tests) and the schema-v2 JSON stdin (API). Record a FOLLOW-UP to deprecate the flag/file form once `--request-file` + JSON cover FO and scripts (removal needs captain sign-off — do NOT remove in this task).
- `team_name` / `bare_mode` reflect LIVE team state, never inferred from the stage; the generator (next-action) fills them from live state, and a hand-authored file must set them correctly. Loud-failure (non-zero exit, named diagnostic) on missing/malformed scalars or sections, matching the existing dispatch-build validation idiom.
- Synergy with 6re (generalize `status --read`): the request file is markdown + frontmatter + named sections — the SAME shape 6re's `status --read` now parses. Coordinate so the request-file parse reuses that machinery rather than a parallel parser.

Acceptance criteria (ideation fleshes; behavior-first / oracle-based):
- `dispatch build --request-file F` produces the SAME dispatch output JSON as the equivalent `--checklist-file` (+ scope/feedback) invocation — a behavioral-equivalence test is the strongest proof the new path is faithful.
- The shared request-file schema is documented and is the artifact `next-action` will emit (recorded so the two paths cannot drift).
- Loud-failure on a malformed / missing-scalar request file; the existing flag/file and JSON-stdin forms are byte-unchanged.

Relationship: pairs with the `next-action` 0205 member (shared schema). Not 6re, but shares its markdown-with-frontmatter worldview. Sprint-tagged 0205-layered-fo as the dispatch-input seam of the layered-FO arc — untag if you'd rather it ride a separate release.
