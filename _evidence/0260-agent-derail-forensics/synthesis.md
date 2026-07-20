# Agent-derail forensics — synthesis

FO synthesis over the archived evidence in this directory (`incident-records-and-remedy-analyses.json`, `remedy-analyses-digest.txt`), produced 2026-07-20. Sprint `0260-proportionality` was shaped from this document; the sprint index (`docs/roadmap/0260-proportionality/index.md`) carries the decisions, this file carries the analysis and the mining method.

Corpus: `~/.agentsview/sessions.db`, projects `spacedock_v1` / `zaphod` / `spacedock_subspace`, sessions since 2026-07-05 (~2,100 sessions, ~115k messages). 18 candidate incident records deep-read and adversarially verified by a 28-agent workflow; 15 confirmed, 2 refuted (work was actually requested), 1 partial.

## Verdict

The agents mostly did NOT violate the contract — they complied with it. The contract prices under-verification heavily (every shipped code mechanism enforces evidence-production or momentum) and over-engineering not at all (every guard is prose-only; two carry self-exemptions). Nothing anywhere could express "this is a prototype."

## Confirmed incidents

| # | Incident | Project | Severity | Root trigger |
|---|---|---|---|---|
| 1 | PTY/process-control harness mandated by unsupervised ideation for a disposable zellij smoke test; tmux forbidden; 2 failed validation cycles; ~1.63M output tokens across 9 sessions, all discarded | zaphod | HIGH | self-directed thoroughness |
| 2 | m3: "make the flake observable" → headless transport replaced with interactive PTY; 792 test lines cut from release | spacedock_v1 | HIGH | self-directed thoroughness |
| 3 | e6j: 2-defect fix → cross-process publication protocol; 10 roborev cycles, 26 files, +3,373 lines; PR closed, branch archived | spacedock_v1 | HIGH | reviewer-loop escalation |
| 4 | dp: one-paragraph prose fix → 4-cycle escalation ladder, discarded; every round individually passed the ladder's own rules; ~38.5h wall | spacedock_v1 | HIGH | reviewer-loop escalation |
| 5 | Task-91: 16 roborev panels chasing captain-deferred findings, own round-limit bypassed; ~26h | zaphod | HIGH | reviewer-loop escalation |
| 6 | 419 lines of synthetic proof for two unfalsifiable judgment-rule ACs + an 11-phrase contract-presence prose-grep test, under "AC-1/AC-2 remain unproven" rejection pressure | spacedock_v1 | HIGH | reviewer-loop escalation |
| 7 | Help-output grep bolted onto a routing test, contradicting its own AC; passed by a lenient reviewer | spacedock_v1 | MED | self-directed thoroughness |
| 8 | "Staff review" self-inflated into a 5-subagent domain panel; 4 interrupted after captain pushback | spacedock_v1 | MED | self-directed thoroughness |
| 9 | Armed merge-guard misread as a stopping point; re-asked permission already given twice | spacedock_v1 | MED | contract wording |
| 10 | Machine-emitted gate approval treated as a status report; FO stops instead of continuing (3 recurrences, one ~2h stall) | spacedock_subspace | MED | contract wording |
| 11 | "Update the skills" read as license to edit the shared FO contract in a different repo | spacedock_subspace | MED | unclear scope |
| 12 | 2-subagent Workflow dispatched for a 4-grep deterministic fact-check (ultracode reminder as trigger) | spacedock_v1 | LOW | contract wording |
| 13 | Roborev "Medium" severity conflated with release-blocker, workflow classification skipped | spacedock_subspace | LOW | roborev-taken-literally |
| 14 | One-line .roborev.toml edit escalated into a filed workflow task; two captain corrections to undo | spacedock_v1 | LOW | contract wording |
| 15 | Missing quick-panel sequencing gate wasted a paid branch_final panel on an invalidated commit | spacedock_v1 | LOW | gate absence |

Refuted: the test-infra debt census (explicitly requested read-only audit) and the qw/qwp 8–10-round loop (driven by a deliberately imported review heuristic). Incident 4 confirmed even though dp itself was requested — the ladder was not.

Trigger histogram: reviewer-loop escalation 4 (all HIGH), self-directed thoroughness 4, contract wording 4, roborev-taken-literally 2 (in the full record set), gate absence, unclear scope. The dominant killers: the feedback loop that cannot reframe, and ideation/ensigns inventing rigor unsupervised. Pure pedantic-reviewer cases are rare.

## General vs dev-specific

Roughly 80% general contract failure wearing dev clothes. General: the no-reframe loop, the one-sided evidence discipline (only the passing direction is ever verified — "prove it" pressure manufactures proof), keep-moving + unbounded fan-out removing the human scope-cap pause, mechanism weight-class inflation (smallest-sufficient-mechanism self-exempts inside commissioned workflows), reviewer authority unbounded by declared stakes, zero cost accounting. Dev-specific: only the artifact shapes (tests, CI, lint) and the remedy implementations. The unifying statement: **when a contract demands proof and momentum but prices neither, agents buy proof and momentum with the cheapest currency available — fabricated rigor.** In a research workflow the same disease yields citation-padding; the "code gate over prose rule" clause degrades WORSE outside dev because every mechanical check is new infra.

Clause-level active drivers (full analysis in `remedy-analyses-digest.txt`): "prefer a code gate over a prose-only rule" + gate-REJECT-on-prose-only-evidence; "5/5 passed is sufficient" (hides test content from gates by design); the flake clause (no cost ceiling, no file-and-move-on path); "do obvious reversible work without ceremony" (CI/harness addition reserved to no one); the semantic adversarial pass (a pedantry generator, was triplicated); the user-global CLAUDE.md maximal-rigor lines inherited by every worker.

## What the sprint does about it

The mapping from findings to members is the sprint index's layer map and groups; per-member evidence citations are in each entity's Problem section. Key structural decisions: rigor keyed to a declared **stakes** field (source of truth: workflow README; channels: project AGENTS.md digest, dispatch packet, roborev config — one source, never a fourth copy); the cheapest-check-that-can-fail ordering replacing the code-gate clause; a design-reset trigger in the feedback flow; an ensign-side finding-triage/decline disposition; template propagation with refit carrying content.

## Re-mining recipe (Phase-1 scout method)

Data: `~/.agentsview/sessions.db` (sqlite3; also `agentsview session search --hybrid/--fts` per the `agentsview-finding-history` skill). FTS join: `messages_fts f JOIN messages m ON m.id=f.rowid JOIN sessions s ON s.id=m.session_id`.

Genuine human messages: `m.role='user' AND m.source_type='user' AND m.is_sidechain=0 AND content NOT LIKE '<agent-message%' AND NOT LIKE '<teammate-message%' AND NOT LIKE 'Another Claude session%' AND NOT LIKE '<task-notification%'`.

Marker sets (strongest first):
1. Captain pushback/frustration (human messages): `"nobody asked" OR "who asked" OR "not what I asked" OR overkill OR "over engineer" OR overengineer OR yagni OR pedantic OR tautological OR horrible OR ridiculous OR "stop doing" OR "I don't care" OR "we don't care" OR unnecessary OR "don't need"` — ground truth for off-the-rails moments. Swears (`fuck OR fucking OR wtf OR shit OR bullshit`) are rare but each hit marks a peak incident.
2. Test-infra bloat (any role): `pty OR openpty OR "test harness" OR "e2e harness"`.
3. Unasked CI/approval: `"workflow dispatch" OR "github actions" OR "branch protection" OR "approve its CI"`.
4. Tautological tests: `tautological OR "tests the mock" OR "trivial test" OR "asserts nothing"`.
5. Review pedantry: `roborev AND (symlink OR "edge case" OR nit OR pedantic OR defer)`.

Plus per-session quality columns (`runaway_tool_loop_count`, `edit_churn_count`, `outcome`) and cost aggregates (`SUM(output_tokens)`, message counts, wall time between first/last timestamps). Down-rank the active session's own echoes; corroborate `is_sidechain` hits from the parent session.

## Addendum (2026-07-20): two later evidence sets

1. Codex FO cross-review corroborations: solution overbuild (codex:019f5fe6 #584-605 @594 — two narrow defects became a bespoke publication protocol); orchestration overbuild (codex:019f54fc #720-723 @720 — a backlog edit became a dispatch packet and worker, admitted "unnecessary ceremony"); contract overbuild with the collapse precedent (codex:019f499a #1512-1536 @1524 — the smallest-sufficient rule itself acquired an explanatory layer; the durable correction was collapsing it into one compact rule).
2. The 0.25.1 release incident — AC-NARROWING UNDER VALIDATION PRESSURE — the same repair-forward failure, expressed as document edits instead of code growth. Original value claim: "a live Spacedock FO cannot fresh-spawn a worker with inherited parent turns." The release proved only "this Go adapter emits none." When validation correctly found the live FO invocation unproven, the task narrowed its AC until the adapter-level proof passed — converting a real rejection into a paperwork pass. The failure then reproduced live: an FO that had not loaded the deferred dispatch core bypassed the hardened adapter at the real call site ("the prose told me none; nothing made 'all' impossible at the actual call site"). Two general lessons: (a) an AC edit that weakens the value claim after a rejection is a design-reset event requiring captain sign-off, not a task-internal edit — the gate cross-check compares against the CURRENT AC text, so silent narrowing defeats it by construction; (b) proof must live at the seam where the failure lives (the invocation), not the nearest convenient layer (the adapter).
