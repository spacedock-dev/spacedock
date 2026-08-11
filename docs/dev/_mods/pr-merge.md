---
name: pr-merge
description: Push branches and create/track GitHub PRs for workflow entities
version: 0.27.0
---

# PR Merge

Manages the PR lifecycle for workflow entities processed in worktree stages. Pushes branches, creates PRs, detects merged PRs, and advances entities accordingly.

The terminal merge ceremony keys off the pending terminal-target application that `spacedock gate consume` produces (route `approved-awaiting-merge`) — not off the entity's stage: `merge guard` discovers and arms this hook at delivery time, and a failed delivery that needs rework sends the entity back through its declared `feedback-to` via `merge guard --rework`.

## Hook: startup

Scan all entity files (in the workflow directory only, not `_archive/`) for entities with a non-empty `pr` field and either non-terminal status, terminal status with `mod-block: merge:pr-merge`, or terminal status carrying a valid merged sentinel (`pr-merge:` or `local-merge:`). A sentinel row already proves MERGED: bypass `gh`, run `spacedock merge guard {slug} --workflow-dir {dir} --verdict passed` directly, and stop processing that row. For every bare PR row, preserve its repository qualifier while parsing it: `owner/repo#N` is PR `N` in `owner/repo`, while `#N` or `N` is PR `N` in the entity's current code repository. Check qualified references with `gh pr view {N} --repo {owner}/{repo} --json state --jq '.state'`. For an unqualified reference, resolve the current code repository from the entity's `{worktree}` using `git -C {worktree} remote get-url origin` and `gh repo view {code-origin-url} --json nameWithOwner --jq '.nameWithOwner'`, then use that value in `gh pr view {N} --repo {code-owner}/{code-repo} --json state --jq '.state'`. Stop and report an unresolved worktree or repository instead of consulting the launch directory.

If `MERGED`, first record and commit the landed merge sentinel with `spacedock status --workflow-dir {dir} --set {slug} pr=pr-merge:{N}` then `spacedock state commit {slug} --workflow-dir {dir}`; next finalize through `spacedock merge guard {slug} --workflow-dir {dir} --verdict passed`. The sentinel is the restart-safe durable signal; the guard clears any in-flight `mod-block`, terminalizes, archives, and commits the archive move atomically. Clean up any worktree/branch and report each auto-advanced entity to the captain.

If `CLOSED` (closed without merge), report to the captain: "{entity title} has PR {pr number} which was closed without merging. How to proceed? Options: reopen the PR, create a new PR from the same branch, fall back to local merge, or send the entity back for rework." Wait for the captain's direction before taking action. On captain direction to send back: `spacedock merge guard {slug} --rework --workflow-dir {dir}` routes the entity through its declared `feedback-to` and clears `pr`/`mod-block` — do not edit frontmatter by hand. Retryable trouble (reopen, new PR, local merge) keeps the pending approval and the delivery retry; only `--rework` supersedes it.

If `OPEN`, no action needed — the PR is still in review.

If `gh` is not available, warn the captain and skip PR state checks.

## Hook: idle

Check PR-pending entities using the same logic as the startup hook: a terminal row carrying a valid merged sentinel (`pr-merge:` or `local-merge:`) bypasses `gh` and resumes `merge guard` directly, while bare PR rows (including terminal rows still carrying `mod-block: merge:pr-merge`) run `gh pr view` and advance through the same committed sentinel then guard path. This is the workflow's PR-pending scan: the generic event loop fires this idle hook and owns no PR scan of its own, so a workflow with no `pr-merge` mod never reaches for `gh` in its loop. Report any advanced entities to the captain.

## Hook: merge

Resolve the PR base once: `BASE=$(spacedock dispatch trunk --workflow-dir {dir})` — the workflow's configured integration trunk (default `main` when no `trunk:` key is set). `dispatch trunk` emits exactly a **bare branch name** (e.g. `main`), so `$( )` yields `$BASE` clean (command substitution strips the single trailing newline). Retain the workflow-supplied `{branch}` as the data value `BRANCH`. Always quote `"$BASE"` and `"$BRANCH"` at use sites; do not evaluate either as shell syntax.

**PR APPROVAL GUARDRAIL — Do NOT push or create a PR without explicit captain approval.** Before presenting the draft, construct the full PR body so the captain reviews the actual prose that will land on GitHub.

Record the exact code commit being submitted before constructing the draft: `CANDIDATE_SHA=$(git -C {worktree} rev-parse HEAD)`. Resolve every other input explicitly rather than from the launch directory:

- Compute the candidate's short SHA with `git -C {worktree} rev-parse --short "$CANDIDATE_SHA"`; if it exits non-zero, report the error and stop. Use this recorded candidate for every later candidate identity—never read ambient `HEAD` after approval.
- Resolve the code repository with `git -C {worktree} remote get-url origin`, then pass that URL as the repository argument to `gh repo view {code-origin-url} --json nameWithOwner --jq '.nameWithOwner'` and retain its result as `CODE_REPO`.
- Locate entity state through the launcher: run `spacedock status --workflow-dir {dir} --resolve {entity ref} --json` and consume its `path` field as `ENTITY_PATH`. Derive `STATE_ROOT` with `git -C "$(dirname "$ENTITY_PATH")" rev-parse --show-toplevel`, the immutable state commit with `git -C "$STATE_ROOT" rev-parse HEAD`, and the state-relative entity path with `git -C "$STATE_ROOT" ls-files --full-name -- "$ENTITY_PATH"`. Resolve the state repository by passing `git -C "$STATE_ROOT" remote get-url origin` as the repository argument to `gh repo view ... --json nameWithOwner --jq '.nameWithOwner'`. Stop if the entity is not tracked or any resolution fails.
- Compute the short entity-id slot with `spacedock status --workflow-dir {dir} --short-id {entity ref}` (shortest-unique-prefix for sd-b32 workflows, literal stored ID for sequential and slug, matching the status table's ID column).

Build the full PR body using the template below — motivation lead, `## What changed`, `## Evidence`, `---` separator, `[{short-id}](...)` audit link, and `Closes {issue}` line if frontmatter `issue` is set. Create a mode-0600 temporary file `PR_BODY_FILE`, then write those exact bytes directly to it without a shell-interpolated heredoc or command string. Present the body by reading that file. This is the file that will be passed to `gh pr create`; do not reconstruct or rewrite it after approval.

Then present the draft to the captain:

- **Title:** {entity title}
- **Branch:** $BRANCH -> $BASE
- **Candidate:** $CANDIDATE_SHA
- **Changes:** {N} file(s) changed across {N} commit(s), computed from `git -C {worktree}` and `$BASE`
- **Files:** {list of changed files from `git -C {worktree}` and `$BASE`}
- **Body:**

  ```
  {constructed body}
  ```

Wait for the captain's explicit approval before pushing. Do NOT infer approval from silence, acknowledgment of the summary, or the gate approval that preceded this step — only an explicit "push it", "go ahead", "yes", or equivalent counts.

**On approval:** First, push the trunk from the code worktree to ensure the remote is up to date: `git -C {worktree} push origin "$BASE"`. If that push fails (no remote, auth error), report to the captain and fall back to local merge.

Resolve the current integration tip as `BASE_SHA=$(git -C {worktree} rev-parse "$BASE")`, then exercise the approved commit against it without changing the candidate ref, index, or worktree: `git -C {worktree} merge-tree --write-tree "$BASE_SHA" "$CANDIDATE_SHA"`. Inspect its stdout, stderr, exit status, and repository context; the exit status is one signal, not the semantic verdict.

- If the evidence indicates a clean merge, continue without rebasing.
- If the evidence indicates an actual content conflict, stop PR and local-merge delivery, surface the conflict evidence, and preserve the pending delivery authority. Do not rebase, automatically resolve, or use force operations; leave reconciliation owner selection to the consuming workflow.
- If the command fails or the evidence is incomplete or ambiguous, mergeability is unknown: report the error, preserve the pending authority, and stop delivery; do not rebase or use local merge as a fallback.

For a clean result, push the approved commit with the exact-SHA refspec: `git -C {worktree} push origin "${CANDIDATE_SHA}:refs/heads/${BRANCH}"`. If that push fails, report to the captain and fall back to local merge.

Then invoke `gh pr create` against the resolved code repository with title, branch, base, and body file as separate arguments: `gh pr create --repo "$CODE_REPO" --base "$BASE" --head "$BRANCH" --title "$PR_TITLE" --body-file "$PR_BODY_FILE"`. Supply `CODE_REPO`, `BRANCH`, and `PR_TITLE` as data values from the resolutions and workflow context above; never interpolate entity text into executable shell syntax or use `eval`. The submitted body must be the exact reviewed bytes in `PR_BODY_FILE`. Remove the temporary file after success or failure. If `gh` is not available, warn the captain and fall back to local merge.

### PR body template

Lead with motivation + end-user value; audit metadata goes at the bottom. The goal is that a reviewer or future debugger sees the "why" first and the audit link last.

**Template structure (top to bottom):**

| Section | Required | Content |
|---|---|---|
| Motivation lead | **yes** | 1 sentence, ≤ 25 words, blending motivation and end-user value. No parentheticals. |
| `## What changed` | **yes** | Action-verb bullets, 3–5 total, each ≤ 15 words. One change per bullet. No rationale inside the bullet — if a change needs justification, it belongs in the entity body, not the PR. |
| `## Evidence` | **yes when verification ran** | Test suites with `N/N passed` format, 1–2 bullets. Do not include per-test-class breakdowns or enumerated suite lists — one pass ratio per suite, plus at most one line confirming live-probe verification. |
| `## Review guidance` | optional | 1 line pointing reviewer at the critical file or risky change — include only when a stage report explicitly flagged it |
| `---` separator + `[{entity-id}](/{state-owner}/{state-repo}/blob/{state-sha}/{state-relative-path})` | **yes** | Audit link, at the bottom |
| `Closes {issue}` | **yes when issue set** | Under the audit link, using the value exactly as it appears in frontmatter, e.g., `#48` or `owner/repo#48` |
| `Related: {siblings}` | optional | Under Closes, only when stage reports flagged follow-ups |

**Extraction rules (apply deterministically from the entity file):**

| PR body section | Source in entity file | Transformation |
|---|---|---|
| Motivation lead | Entity body paragraph(s) between closing `---` and the first `##` heading | Condense first paragraph to 1-2 sentences. Lead with impact or action verb — not "This PR" or "This task". Blend motivation + value. |
| What changed | Latest stage report whose declared stage outputs describe completed deliverable work | Use DONE items as one action-verb bullet per meaningful unit. Collapse sibling bullets that describe the same thing. Do NOT include "what we deliberately did NOT change" bullets — scope boundaries belong in the entity body, unless a later verification report flagged them as risk. Select this report by the stage's declared outputs and content, never by requiring a particular stage name. |
| Evidence | Latest stage report whose declared stage outputs independently verify the candidate against acceptance criteria | One bullet per suite with `N/N passed` format. Include any quantitative result the report explicitly called out (wallclock delta, size %, perf). If the workflow declares no independent verification-report role, fall back to self-test evidence in the deliverable-work report. Select by the declared role and content, never by requiring a particular stage name. |
| Review guidance | Explicit "focus on X" / "risk here" notes in either stage report | 1 line. **Omit if no such note exists.** |
| Audit link | Short entity id from `spacedock status --workflow-dir {dir} --short-id {entity ref}`; resolved state-repository owner/name, immutable state commit, and state-relative entity path from the explicit state resolution above | Format as `[{short-id}](/{state-owner}/{state-repo}/blob/{state-sha}/{state-relative-path})` |
| Closes | Entity frontmatter `issue` field (exactly as written) | Prefix `Closes ` |
| Related | Explicit "related entity" / "follow-up" mentions in stage reports | 1 line. **Omit if none.** |

Target total length: **60-120 words**.

**Key design decisions:**

1. **Lead with motivation + end-user value.** First content is a 1-2 sentence user-facing impact statement. The audit link moves to the bottom as audit metadata.
2. **Prescribed sections + extraction rules** — not a strict verbatim template, not free-form. The mod specifies headings and source subsections; the FO paraphrases rather than pasting.
3. **Evidence follows declared report roles.** Workflows without an independent verification report fall back to deliverable-report self-test evidence.
4. **Review guidance and Related are opt-in.** They appear only when stage reports explicitly flagged them, to prevent bloat.

Set the entity's `pr` field to the PR number (e.g., `#57`). Report the PR to the captain.

**On decline:** Do NOT automatically fall back to local merge. Ask the captain how to proceed — options include local merge or leaving the branch unmerged. Only act on the captain's explicit choice.

Do NOT archive yet. The entity stays at its current stage with `pr` set until the PR is merged. The FO handles advancement to the terminal stage and archival when it detects the merge through this portable startup or idle hook. A runtime may invoke the same hook check as a generic backstop; no host-specific reconcile class is required.
