---
title: "Pi live E2E runs as an associated, manually-approved PR check"
status: validation
source: "Captain directive 2026-09-11 at the s98 validation gate: 'approve the pi ci lane' + 'the pi live is simply skipped and we can't manually trigger that check to be really associated with the PR'."
id: 6v7sdg4pfzwcfnjstt80pvqj
gates:
    version: 1
    records:
        - id: gate:6v7sdg4pfzwcfnjstt80pvqj:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:6v7sdg4pfzwcfnjstt80pvqj-backlog-1
              briefing:
                id: briefing:6v7sdg4pfzwcfnjstt80pvqj:backlog:attempt-1:revision-1
                digest: sha256:6faef9094b173d98daaeba38e9c398603310446f1490ee0edf6340897a4eafd9
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:6v7sdg4pfzwcfnjstt80pvqj:backlog:1
                briefing: briefing:6v7sdg4pfzwcfnjstt80pvqj:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-11T06:01:25.950478Z"
                decision: approve
                reason: 'Captain directed in chat 2026-09-11 at the seed presentation: ''let''s get 6v done first so we can get s98 with proper tag and pr'' — approves the seed to enter ideation; the lane must land before s98''s PR delivery'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:6v7sdg4pfzwcfnjstt80pvqj:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:6v7sdg4pfzwcfnjstt80pvqj-ideation-1
              briefing:
                id: briefing:6v7sdg4pfzwcfnjstt80pvqj:ideation:attempt-1:revision-1
                digest: sha256:6eb6ed68eb0dd6b146d0d3dce21e01fb84465cd89621583b850a28fc355d6374
                room-ref: '@review/ideation/briefing-1'
              withdrawal:
                by: agent:first-officer
                at: "2026-09-11T06:11:08.934571Z"
                reason: 'Artifact normalized after prepare under captain direct-edit grant (''fix it'', 2026-09-11): AC heading renamed to the scanner-exact ''## Acceptance criteria''; the bound briefing digest would be stale at record time'
            - id: gate-attempt:6v7sdg4pfzwcfnjstt80pvqj-ideation-2
              briefing:
                id: briefing:6v7sdg4pfzwcfnjstt80pvqj:ideation:attempt-2:revision-1
                digest: sha256:c88af45c39e52fa0e3dd0e6215a5395a643b07259c1da9e9fc7ef75d8541b74e
                room-ref: '@review/ideation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:6v7sdg4pfzwcfnjstt80pvqj:ideation:2
                briefing: briefing:6v7sdg4pfzwcfnjstt80pvqj:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-09-11T06:35:42.856775Z"
                decision: approve
                reason: 'Captain approved in chat 2026-09-11 at the attempt-2 presentation: minimal label-gated diff with spillover guard, lean proof per the no-test-infra directive'
              application:
                target-stage: implementation
                state: consumed
started: 2026-09-11T06:01:40Z
worktree: .worktrees/spacedock-ensign-pi-live-pr-check-lane
---
## Problem

In `.github/workflows/runtime-live-e2e.yml`, `claude-live` and `codex-live` run on every `pull_request` (behind the `CI-E2E` / `CI-E2E-CODEX` required-reviewer environments), but `pi-live` is gated `if: github.event_name == 'workflow_dispatch' && inputs.live_cadence == 'pi'` — it never runs on a PR, and a manual `workflow_dispatch` run against a PR branch produces no check run on that PR (the workflow's own journey-delta-comment job notes "a workflow_dispatch run has no PR"). So pi live evidence is neither PR-associated nor manually triggerable per-PR; it only exists as an unassociated manual run or the release-time precondition. Today's baseline: **zero PR-associated `pi-live` check runs have ever existed** — that is the number this entity moves to ≥1.

## Proposed approach (the minimal workflow diff)

Mirrors the existing claude/codex job pattern. Three coordinated edits to `.github/workflows/runtime-live-e2e.yml`, nothing else:

**1. Widen the PR trigger with `labeled`** (the only functional addition; `synchronize` is already a default type but is listed explicitly because declaring `types:` replaces the default set):

```yaml
# before
  pull_request:
    branches: [main]
# after
  pull_request:
    branches: [main]
    types: [opened, synchronize, reopened, labeled]
```

**2. Label-gate `pi-live` onto PR events** (manual dispatch leg unchanged):

```yaml
# before (pi-live)
    if: ${{ github.event_name == 'workflow_dispatch' && inputs.live_cadence == 'pi' }}
# after
    if: ${{ (github.event_name == 'pull_request' && contains(github.event.pull_request.labels.*.name, 'live:pi')) || (github.event_name == 'workflow_dispatch' && inputs.live_cadence == 'pi') }}
```

**3. Guard the other PR-gated jobs against the widened trigger.** Adding `labeled` makes `github.event_name == 'pull_request'` true on label events, so `claude-live` (line 88), `codex-live` (line 317), and `journey-delta-comment` (line 842) each gain `&& github.event.action != 'labeled'` — without this, adding e.g. a `bug` label would re-run both paid live lanes and re-post the journey comment. This is not a new mechanism; it preserves today's semantics for the existing lanes under the widened trigger. `offline` stays unconditional: on a `labeled` event it is the gate `pi-live` `needs`, and it is cheap and secret-free.

**CI-E2E-PI environment:** already declared on `pi-live` and already in use by manual dispatches. The only wiring to verify at implementation is repo-settings: the `CI-E2E-PI` environment exists with a required reviewer (it does for the manual cadence today); if it is missing, that is a one-time repo-settings step to surface to the captain — not YAML.

**Observable semantics:** the `live:pi` label is the cost opt-in; the `CI-E2E-PI` required-reviewer approval is the human gate; the run appears in the PR's Checks tab associated with the PR head SHA, where it can be a required check where warranted. Labels persist, so later `synchronize` pushes re-run pi-live while labeled; `reopened` re-runs it too (consistent with claude/codex). Unlabeled PRs and non-pi labels run no pi-live. The manual `workflow_dispatch` cadence is byte-for-byte unchanged in behavior.

**Security posture unchanged:** the trigger stays `pull_request` (not `pull_request_target`), so fork PRs still receive no secrets; only triage+ can add labels, so the label introduces no new exfiltration surface — the environment approval remains the human gate.

**Rejected (stays rejected unless a blocker appears):** workflow_dispatch + out-of-band commit-status on the PR head SHA (bypasses the Checks-tab model, needs a `pr_number` input); scheduled nightly pi cadence on main (cheap, but never PR-associated). No `pr_number` inputs, no commit statuses, no scheduled lanes.

## Acceptance criteria

- **AC-1 (value, measured against today's baseline of zero):** a PR carrying the `live:pi` label exhibits a PR-associated `pi-live` check run in its Checks tab (head SHA) after `CI-E2E-PI` approval. Proof owner: the lane's first REAL run on 6v's own delivery PR (dogfood), not a synthetic test. Falsifying check: label present + approval granted but no pi-live check run on the PR ⇒ fail.
- **AC-2:** no PR event runs `pi-live` without the `live:pi` label (the cost opt-in holds). Proof owner: workflow-syntax inspection at the ideation/implementation gate, plus the dogfood PR observed unlabeled before labeling. Falsifying check: pi-live appears on an unlabeled PR ⇒ fail.
- **AC-3:** the existing manual `workflow_dispatch` `live_cadence=pi` cadence is unchanged. Proof owner: the next release-time manual dispatch (existing lane). Falsifying check: a `live_cadence=pi` dispatch skips pi-live ⇒ fail.
- **AC-4:** a non-pi label add re-runs neither claude-live, codex-live, nor the journey-delta comment, while unlabeled PRs still run offline + claude-live + codex-live as today. Proof owner: workflow-syntax inspection + observation on the dogfood PR (scratch label added before `live:pi`). Falsifying check: a `bug`-style label add re-runs a paid lane ⇒ fail.
- **AC-5:** the edited workflow remains valid Actions workflow syntax. Proof owner: actionlint at implementation. Falsifying check: actionlint error on the edited file ⇒ fail.

## Test plan (captain directive 2026-09-11: "no tests for test infra bs")

No unit tests of CI YAML. Proof is two-layered:

1. **Deterministic, zero-cost:** actionlint (or equivalent workflow parse) over the edited `.github/workflows/runtime-live-e2e.yml`. Distinct falsifying edit: introduce a YAML/expression syntax error → validation fails.
2. **Live, one-off manual validation (the dogfood, on 6v's own PR):** (a) confirm no pi-live check run before labeling (AC-2); (b) add a scratch non-pi label, confirm claude-live/codex-live do not re-run (AC-4); (c) add `live:pi`, approve `CI-E2E-PI`, confirm the pi-live check run appears on the PR and runs the real journeys (AC-1). Estimated cost: one real pi-live run (~the existing manual cadence cost); the deterministic step is trivial.

## Expected surface

- `.github/workflows/runtime-live-e2e.yml` — net +5 (≈ +9 insertions / −4 deletions), tolerance net ±10, 1 file.
- `docs/runtime-live-ci.md` — net +1 (≈ +3 / −2), 1 file (user-visible CI-lane doc; the only doc that describes the lanes).

Observable semantics declared: CI runtime behavior only — when pi-live runs (label-gated PR events; manual dispatch unchanged) and what gates it (`live:pi` label + `CI-E2E-PI` required-reviewer approval). No command grammar, stored-format, or authority changes.

### Doc diff (docs/runtime-live-ci.md, "Workflow:" section, line-171 bullet tail)

- Before: "…Pull requests still run only Sonnet and Codex; Pi is optional and is not a merge requirement."
- After: "…Pull requests run Sonnet and Codex; Pi is opt-in per PR — add the `live:pi` label and approve the `CI-E2E-PI` environment to attach a Pi live check to the PR. Pi is not a merge requirement."

## Spike

No spike needed: the design relies on already-proven mechanisms — GitHub's `pull_request` `labeled` trigger type and the `contains(github.event.pull_request.labels.*.name, …)` expression are documented Actions behavior, and environment-gated live jobs on pull_request events are the existing claude/codex pattern in this same workflow. The one behavior the dogfood must observe (AC-4's labeled-event spillover guard) is covered by the static inspection in the same pass plus the behavioral observation on the dogfood PR.

## Stage Report: ideation

- DONE: Fleshed-out body with the MINIMAL concrete workflow diff (labeled+synchronize trigger types, label-gated pi-live if-condition, CI-E2E-PI wiring verification)
  Body now carries three concrete before/after YAML edits to `.github/workflows/runtime-live-e2e.yml`; CI-E2E-PI is already declared on pi-live (workflow line 586), so wiring = repo-settings verification, not new YAML.
- DONE: Labeled-event spillover guard identified and folded into the minimal diff
  Adding `labeled` makes `github.event_name == 'pull_request'` true on label events; claude-live (l.88), codex-live (l.317), journey-delta-comment (l.842) each gain `&& github.event.action != 'labeled'`, offline stays unconditional (pi-live's `needs` gate). Not a new mechanism — preserves today's semantics under the widened trigger; without it a `bug` label re-runs both paid lanes.
- DONE: Entity-level ACs, each naming proof owner and falsifying check
  AC-1 value AC measured against today's baseline (zero PR-associated pi-live check runs ever) with proof = the lane's first REAL dogfood run on 6v's own PR; AC-2 unlabeled-PR opt-in; AC-3 manual dispatch unchanged; AC-4 non-pi-label no-re-run; AC-5 actionlint validity.
- DONE: No unit tests of CI YAML per captain directive 2026-09-11
  Test plan is two-layered: deterministic actionlint validation + one-off live dogfood run (unlabeled → scratch label → live:pi + approval) on the dogfood PR.
- DONE: Expected surface and observable semantics declared
  1 workflow file net +5 (≈+9/−4, tolerance ±10) + docs/runtime-live-ci.md net +1 (the only user-visible doc describing the lanes, with concrete before/after wording); semantics = CI runtime behavior only (when pi-live runs and what gates it).
- DONE: Mechanisms beyond the label gate rejected in the body
  pr_number inputs, out-of-band commit statuses, scheduled lanes all stay rejected; spike section records "no spike needed" with the proven mechanisms relied on (documented Actions labeled trigger + contains() expression; existing claude/codex PR+environment pattern in the same workflow).

### Summary

Fleshed out the 6v entity body into a gated ideation design: a three-edit diff to runtime-live-e2e.yml that makes pi-live an opt-in PR check via a `live:pi` label plus the existing CI-E2E-PI required-reviewer environment, with a necessary spillover guard on the other PR-gated jobs (the only addition beyond the seed's named edits, justified as semantics-preservation of the trigger widening, not a new mechanism). Proof stays lean per the captain directive: actionlint + the lane's first real dogfood run on 6v's own PR; AC-1 measures against today's zero-baseline. No code was changed — ideation output is the body itself, committed path-scoped to the state checkout.

## Stage Report: implementation

- DONE: Apply edit 1 — widen the pull_request trigger with labeled+synchronize
  `.github/workflows/runtime-live-e2e.yml:52` now reads `types: [opened, synchronize, reopened, labeled]` under the existing `pull_request: branches: [main]`.
- DONE: Apply edit 2 — label-gate pi-live onto PR events (manual dispatch leg unchanged)
  `pi-live` if-condition (line 584) is now `(pull_request && contains(pull_request.labels.*.name, 'live:pi')) || (workflow_dispatch && inputs.live_cadence == 'pi')`; the manual cadence disjunct is byte-identical to before.
- DONE: Apply edit 3 — spillover guard on claude-live, codex-live, journey-delta-comment
  claude-live (l.89), codex-live (l.318), journey-delta-comment (l.843) each gained `&& github.event.action != 'labeled'`; offline left unconditional as the gate pi-live `needs`.
- DONE: Parenthesization detail (flagged, not absorbed): the guard wraps each job's full disjunction in parens
  Actions `&&` binds tighter than `||`, so the literal "gain `&& …`" reading would guard only the last disjunct (e.g. leave a PR-label event re-running claude-live via the first disjunct). The approved guard semantics required the paren form; no other wording changed.
- DONE: Verify CI-E2E-PI environment wiring (repo-settings read, no new YAML)
  `gh api repos/spacedock-dev/spacedock/environments/CI-E2E-PI`: exists with a `required_reviewers` protection rule (reviewer clkao, prevent_self_review false); environment secrets list shows `CODEX_AUTH_JSON`. `environment: name: CI-E2E-PI` remains declared on pi-live (workflow l.587).
- DONE: Run actionlint over the edited workflow (AC-5)
  actionlint 1.7.7 on `.github/workflows/runtime-live-e2e.yml`: "Found total 0 errors"; exit 0 (shellcheck/pyflakes rules skipped — binaries not installed; only the two new `if:` expressions and the trigger types changed, no embedded scripts touched).
- DONE: Apply the doc diff to docs/runtime-live-ci.md exactly as approved
  Line-171 bullet tail replaced verbatim with the approved after-wording ("Pull requests run Sonnet and Codex; Pi is opt-in per PR — add the `live:pi` label and approve the `CI-E2E-PI` environment …"); net +1 line as declared.
- DONE: Commit deliverable to the worktree branch
  Commit 10ffbf7a8 on `spacedock-ensign/pi-live-pr-check-lane` — 2 files, +6/−5 (within the declared surface tolerance); no tests added per the captain directive.
- SKIPPED: AC-1/AC-2/AC-4 live dogfood observation (scratch label → live:pi → approval on 6v's own PR)
  Proof owner for those ACs is the delivery PR itself after merge/dispatch; implementation proves AC-2/AC-3/AC-5 statically (label-only pi-live condition; unchanged manual disjunct; actionlint 0 errors) and leaves the live observation to validation/dogfood.

### Summary

Implemented the approved three-edit workflow diff plus the one-line doc diff, committed as 10ffbf7a8 on the worktree branch. CI-E2E-PI verified in repo settings (required reviewer + CODEX_AUTH_JSON secret) so no new YAML or secrets were needed. Only deviation-shaped detail is the guard parenthesization, which is required for the approved semantics and is flagged above rather than absorbed. actionlint passes with 0 errors; live dogfood (AC-1) remains owned by the delivery PR.
