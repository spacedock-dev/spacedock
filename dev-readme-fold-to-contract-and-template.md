---
title: Fold docs/dev/README.md into contract pointers + template regeneration — leave only the ~13% local residue
source: "captain (CL), 2026-07-21 — 'the goal is to fold as much of the customized docs/dev/README.md back into contract and template as possible; tell me what remains locally modified and why.' Analysis done (two FO workflows); captain: 'do that on a worktree after everything else lands and let me review.'"
status: backlog
sprint:
id: 6xp4szfqpdnse897e1ypmvyh
---

docs/dev/README.md (254 lines / 29,651 bytes) is the one commissioned workflow README that predates PR #388's "defer universal rules to the contract" restructuring, so it is a hand-drifted fork of doctrine that now lives in the FO/ensign contract and the commission templates. Fold it back so only genuinely-local content remains.

## The mechanism distinction (captain-clarified 2026-07-21 — this is the crux)

"Fold into contract" and "fold into template" are OPPOSITE operations because the two surfaces reach the README by different routes:

- **CONTRACT is loaded at runtime.** Every FO/ensign operating docs/dev has it in context. So the README can POINT at it ("proof discipline: see the FO contract") and genuinely SHRINK. Real byte reduction.
- **TEMPLATE is consumed once at commission, never read at runtime.** The commissioned README is a STANDALONE artifact; nobody loads development.md while running docs/dev. So the README CANNOT point at the template — the text must stay INLINE. "Folding into template" means making the template the SOURCE OF TRUTH and REGENERATING the README from it (refit), not deleting bytes. The win is ending hand-drift, not shrinking.

## Three-bucket classification (of 29,651 B)

| bucket | mechanism | README size effect | bytes | share |
|---|---|---|---|---|
| CONTRACT | pointer (runtime-loaded) | shrinks ~5.6 KB → a few hundred B | ~5,600 | ~19% |
| TEMPLATE | regenerate from template (text stays inline) | no shrink; text becomes template-owned, drift-proof | ~18,300 | ~62% |
| SPRINT doctrine | template-shaped but BLOCKED (no shipped home) | stays for now | ~1,840 | ~6% |
| LOCAL residue | hand-authored | stays | ~3,900 | ~13% |

Only the ~19% contract bucket removes bytes. The 62% is "how much is template-boilerplate that should stop being a hand-maintained fork," not deletable bytes.

## The local residue (what stays and why — the captain's primary question, answered)

- **Testing Resources table (236–247, ~1.5 KB)** — real spacedock artifacts: `go test ./...`, `gotestsum@v1.13.0` + install script, `spacedock status --workflow-dir docs/dev`, spec paths. Another repo differs.
- **Title + mission blurb (31–33, 256 B)** — the workflow's own identity; no template can author another workflow's name.
- **Registered mods (37, 196 B)** — this workflow's `pr-merge` + `comm-officer` choices.
- **Concrete examples in stage prose (~700 B)** — named spacedock entities, issue #262, contract_gate_test.go, the Go language.
- **setsid/PTY review examples (81, ~250 B)** — spacedock building a terminal launcher.
- **path→lane mapping (79, ~250 B)** — this repo's actual CI lanes.
- **Smaller local anchors** — `runtime-live-ci.md` pointer, `.spacedock-state` bootstrap note, `internal/contractlint` carve-out, sprint anchors incl. tracker `xp`.

Residue total ≈ 3,900 B ≈ 13%.

## Fold plan

**→ CONTRACT (defer to a pointer; already-there unless noted):** proof-policy preamble (L72–74 → first-officer-shared-core.md:167–174,199); "evidence must fail" (L77 → ensign-shared-core.md:85); "completion not a stopping point" (L123 → first-officer-shared-core.md:76); read-one-section (L182–184 → :201 + ensign:18); two-checkout commit separation (L254 → ensign-shared-core.md:34–38); the L79/80/103/155/248-clause-1 doctrine halves. SHOULD-MOVE: the fuller prose-grep articulation currently only at L76 → first-officer-shared-core.md:174.

**→ TEMPLATE (regenerate from / add to development.md; text stays inline):** File Naming (L49–51), Schema (L53–70; RE-ADD the drifted-out `pr`/`mod-block` fields), stage semantics, Workflow State (L168–181), Task Template (L190–230; drifted), Commit Discipline (L250–253; L252 "archive boundaries" is stale, template says "merge boundaries"). NEW template content (currently nowhere, 0 grep hits): doc-diff-in-ideation (L112), semantic adversarial pass (L140–146), two-axis finding classification (L147–153), surface-selection ladder (L248 clause 2), complex-ideation staff-review pattern (L115), path→lane-is-the-gate + one-run-stop-at-90-min (L79/81 patterns).

## Two caveats

1. **az's prose-grep ruling (L76)** looks irreducibly local (dated, verbatim, deliberately placed) but its WORDS are generic — correct home is first-officer-shared-core.md:174. Move the DOCTRINE as evergreen prose; preserve the date/provenance/az-attribution in the COMMIT MESSAGE + journal, NOT contract prose (the contract bans temporal wording). Do NOT leave a hedge copy in the README — that recreates the drift.
2. **Sprints section (L41–47, ~2,040 B)** is template-shaped (general pattern) but BLOCKED: `_proposals/sprint-roadmap-construct.md` assigns it to a future `spacedock:roadmap` skill + separate sprint template (neither exists; 0 grep hits) and forbids baking sprints into base development.md. Stays as a blocked template item, N=1, until the sprint template ships.

## Drift findings (arguably bugs a refit fixes)

README Schema omits `pr` + `mod-block` (template has them; this pr-merge workflow needs them); Task Template block drifted the same way; L252 "archive boundaries" is stale pre-pr-merge wording.

## Sequencing (binding order)

1. Template-fold FIRST (edit development.md / commission §2a). 2. Contract-fold as a SEPARATE reviewed change to ratcheted first-officer-shared-core.md (esp. the az ruling move). 3. THEN refit docs/dev so foldable blocks regenerate from the template — if refit runs first it overwrites hand-authored prose or re-injects stale copies. Keep contract-fold out of the refit so the refit only deletes duplication, not adjudicates doctrine.

## Execution note

Captain instruction 2026-07-21: do this on a WORKTREE, AFTER the rest of 0260 lands (bw, 2ae, close-out), and present for CAPTAIN REVIEW — do not auto-merge. It edits ratcheted contract files and the live governing README active turns read, so it wants its own ideation→validation gate. AC candidate: cumulative byte delta vs origin/main is negative for the contract-fold AND every deferred rule resolves to a verified contract line AND a refit of docs/dev reproduces the template-owned sections.
