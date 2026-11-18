---
id: 3kn22a1d5v5x7rdnrcva0h5a
title: Commissioned seed tasks carry an Acceptance criteria section
status: implementation
source: launch probe 2026-11-17, release defect 4 (CL)
started: 2026-11-18T01:26:38Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-commission-seed-entities-carry-acceptance-criteria
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:3kn22a1d5v5x7rdnrcva0h5a:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3kn22a1d5v5x7rdnrcva0h5a-backlog-1
              briefing:
                id: briefing:3kn22a1d5v5x7rdnrcva0h5a:backlog:attempt-1:revision-1
                digest: sha256:82123b6b0dd21245098c471638edfb4809b48d3821078c332e427a040137d1a2
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:3kn22a1d5v5x7rdnrcva0h5a:backlog:1
                briefing: briefing:3kn22a1d5v5x7rdnrcva0h5a:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-11-18T01:26:07.250128Z"
                decision: approve
                reason: 'Captain in chat 2026-11-17: ''file defect 4. and check if any existing live tests cover the commission skill first. then dispatch.'' then ''override. you can push and continue dispatching.'' Accepts the seed and directs ideation.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:3kn22a1d5v5x7rdnrcva0h5a:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:3kn22a1d5v5x7rdnrcva0h5a-ideation-1
              briefing:
                id: briefing:3kn22a1d5v5x7rdnrcva0h5a:ideation:attempt-1:revision-1
                digest: sha256:82a3725c527228bdef52affba7ce12fa550f941b0479146ca279da220e4f3cdf
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:3kn22a1d5v5x7rdnrcva0h5a:ideation:1
                briefing: briefing:3kn22a1d5v5x7rdnrcva0h5a:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-11-18T02:23:44.468907Z"
                decision: approve
                reason: 'Captain approved in chat 2026-11-17: accepts landing c3895378d unchanged with TestCommissionedEntitiesPassACScan as the AC-1 proof.'
              application:
                target-stage: implementation
                state: consumed
---

A workflow made with `spacedock:commission` must give its seed tasks, and every task its templates produce, a `## Acceptance criteria` section, so the first officer's `status --read --ac-scan` works at the first gate.

## Problem

Commission writes the 2–3 seed tasks as frontmatter and a prose body only (`skills/commission/SKILL.md` ~505-521 on origin/main). The README task template it writes does carry `## Acceptance criteria` (~479), and so does the development template. The experiment and refinement templates (`skills/commission/references/templates/{experiment,refinement}.md`) carry no AC heading in their task snippet, so every later task of those workflows also lacks it.

The first officer runs `status --read <ref> --ac-scan` at every gate after the initial stage (`skills/fo-gate-lifecycle/SKILL.md` step 3; the initial-stage gate skips it because no stage report exists yet). It exits non-zero by design when the section is absent (`internal/status/gate_extract.go:366`, "no ## Acceptance criteria section in this file"). So the first scanned gate on a seed task fails in every freshly commissioned workflow: `ideation` for development, `smoke` for experiment, `review` for refinement and for a fallback workflow with a later gate. On experiment and refinement workflows it also fails for every task the first officer files later, because the Claude first officer files from the README's entity template (`skills/first-officer/references/claude-first-officer-runtime.md:35`), and those templates carry no AC section. A new user meets this on their first workflow.

## Proposed approach

Keep the prepared candidate `fix/commission-seed-ac` (`c3895378d`) as is. Implementation lands that commit unchanged on a branch from origin/main. The candidate is the smallest change that puts the scanned heading into every entity body commission writes:

1. The seed body in `skills/commission/SKILL.md` §2b gains an `## Acceptance criteria` block with one `**AC-1 — …**` stub. The stub uses the same parenthesised fill-in style as the existing body instruction, so the commissioning agent writes a real criterion from the captain's seed input instead of copying a brace placeholder.
2. The experiment template's `## Experiment Template` snippet and the refinement template's `## Artifact Template` snippet gain the same block the generic README template and the development template already carry. Each template's `### Inject entity-template snippet` line names the new section. The refinement line also states that variant snippets (outreach, integration, content-production, prd-authoring) keep the section, because those variants are prose lists with no fenced snippet to edit.
3. `internal/contractlint/commission_entity_ac_scan_test.go` extracts each fenced entity body commission instructs (the §2b seed, the SKILL's generic README entity template, and each template's single `## <Entity> Template` block), appends a stage report citing AC-1, and runs the real `status` native runner with `--read --stage implementation --ac-scan --json`. It asserts exit 0 and that the first scanned AC is `AC-1`.

Why keep it unchanged:

- The layer-assembly fallback needs no extra edit. A fallback workflow gets its seeds from §2b and its README entity template from the SKILL's generic `## {Entity_label} Template`, which already carries the section on origin/main. The test covers both bodies.
- The experiment snippet keeps `## Success criteria` next to the new `## Acceptance criteria`. They answer different questions: success criteria decide whether the hypothesis is accepted or rejected; acceptance criteria decide whether the experiment is finished. A rejected hypothesis is still a finished experiment, so merging the two would make the AC scan report a valid rejection as unmet.
- The AC-1 stub is kept, not only the heading. A heading with no AC makes `--ac-scan` exit 0 with an empty list, which passes the scan and gives the gate nothing to cite.

### Before/after wording

`skills/commission/SKILL.md` §2b seed body. Before, the fenced block ends:

```
(Body: write the description or thesis from the captain's seed input as plain prose. Do NOT carry brace-syntax placeholders into the body — rewrite any `{var}` phrasing as natural language or backticked words.)
```

After, the same line is followed by:

```
## Acceptance criteria

Each AC names a property of the finished {entity_label} (not a stage action) and how it is verified.

**AC-1 — (End-state property drawn from the captain's seed input.)**
Verified by: (test name / command output or exit code / file the change produces / resulting on-disk state — something outside this {entity_label} body that a future reader can reproduce and that can fail; name the concrete change that would make it fail.)
```

`skills/commission/references/templates/experiment.md` `## Experiment Template` snippet. Before: `## Success criteria` is followed directly by `## Smoke result`. After, between them:

```
## Acceptance criteria

Each AC names a property of the finished experiment (not a stage action) and how it is verified.

**AC-1 — {End-state property.}**
Verified by: {test name / command output or exit code / file the change produces / resulting on-disk state — something outside this experiment body that a future reader can reproduce and that can fail; name the concrete change that would make it fail.}
```

Adoption line before: `Use the hypothesis-result snippet (hypothesis / methodology / success criteria / smoke result / run result / holdout result / verdict) shown in the Experiment Template section above.`
After: `Use the hypothesis-result snippet (hypothesis / methodology / success criteria / acceptance criteria / smoke result / run result / holdout result / verdict) shown in the Experiment Template section above.`

`skills/commission/references/templates/refinement.md` `## Artifact Template` snippet. Before: the description line is followed directly by `## Draft`. After, between them, the same block with `artifact` in place of `experiment`.

Adoption line before: `Use the refinement snippet (draft / review notes / final) shown in the Artifact Template section above. Variants override this with their own snippet (see below).`
After: `Use the refinement snippet (acceptance criteria / draft / review notes / final) shown in the Artifact Template section above. Variants override this with their own snippet (see below); every variant snippet keeps the `## Acceptance criteria` section, because each gate's AC scan requires it.`

### Documentation

No doc diff. `docs/site` does not describe the sections of a generated entity body; `docs/site/concepts/workflows-and-entities.md:13` already says each record carries acceptance criteria, which this change makes true for commissioned workflows.

## Risk evidence

Backlog check, 2026-11-17: no live test drives the commission skill. The live fixtures in `internal/ensigncycle` hand-write a workflow README with `commissioned-by: spacedock@1`; the "commissioned" live scenario (`claude_live_runner_test.go:568`) uses such a prepared fixture. The only test that reads the commission skill is the static `internal/contractlint/boundary_guard_control_test.go`. The layer-assembly fallback path (no template) is not yet checked.

Ideation spike, 2026-11-17. All runs use a `spacedock` binary built from origin/main `647d55163` (`GOPROXY=off go build ./cmd/spacedock`, exit 0); the candidate changes no Go code.

1. Seed body, before and after. A throwaway script wrote a workflow README per template shape and a seed exactly as §2b instructs (frontmatter plus a prose body), appended a stage report for the first scanned gate citing AC-1, and ran `spacedock status --workflow-dir <wf> --read <wf>/first-seed.md --stage <gate> --ac-scan --json`.
   - Before (origin/main §2b body): exit 1, `Error: no ## Acceptance criteria section in this file`, for development (`ideation`), experiment (`smoke`), refinement (`review`) and fallback (`review`).
   - After (candidate §2b body with the AC-1 stub filled in): exit 0 for all four, `acs[0].id == "AC-1"`, `unevidenced: false`.
2. Candidate test on both trees (`git archive` exports, `GOPROXY=off go test ./internal/contractlint -run TestCommissionedEntitiesPassACScan -count=1`):
   - origin/main plus the test file: FAIL. Three subtests fail with ac-scan exit 1 and the missing-section error: the §2b seed, `experiment.md ## Experiment Template` and `refinement.md ## Artifact Template`. The generic README template and the development template pass.
   - Candidate `c3895378d`: PASS, all five subtests.
3. Full suite on the candidate worktree: `GOPROXY=off go test ./...` passes; the worktree stayed clean.
4. One-off live commission run, to check that a real agent copies the new block. In a fresh git repo: `claude -p --model sonnet --permission-mode bypassPermissions --plugin-dir <candidate export> "/spacedock:commission Batch mode, skip confirmation … experiment template … two seeds … stop after Phase 2"` exited 0, and the agent reported that `status --validate` printed `VALID`. The installed plugin (0.29.0-pre0) has no AC block in §2b or in the experiment snippet, so the AC sections in the output came from the candidate skill text.
   - Both generated seeds carry `## Acceptance criteria` with a filled AC-1 written from the seed input, for example "**AC-1 — The experiment carries a falsifiable verdict on the day-one welcome email effect.**", with no brace placeholder.
   - The generated README's `## Experiment Template` carries the section between `## Success criteria` and `## Smoke result`.
   - A real `status --read <copy> --stage smoke --ac-scan` on each generated seed, and on the generated README template body, each with a smoke stage report citing AC-1, exits 0 with `ac=AC-1 unevidenced=false citations=1`.

## Out of scope

- Making `--ac-scan` tolerate a missing section; its loud failure is the intended contract.
- Repairing workflows that users already commissioned.
- A new standing live lane for commission, unless ideation shows no cheaper check can fail.
- `skills/commission/references/decomposed-snippets.md` snippet lists. They describe the templates loosely (the hypothesis-result list already omits success criteria), and the fallback builds its README entity template from the SKILL's generic template, not from these lists.
- How first officers file later tasks. The Claude runtime reads the README's `## Task Template` before filing, but experiment and refinement READMEs name that section `## Experiment Template` and `## Artifact Template`. Whether filing finds those sections is a separate question for the first-officer contract.

## Expected surface and tolerance

Estimate net LOC change: +146 (148 insertions, 2 deletions), across 4 files: `internal/contractlint/commission_entity_ac_scan_test.go` (new, +125), `skills/commission/SKILL.md` (+7), `skills/commission/references/templates/experiment.md` (+8 net), `skills/commission/references/templates/refinement.md` (+8 net).

Tolerance: net +120 to +170, at most 5 files. Going over means the design changed and needs a new ideation gate.

Semantics this task may change: the content of files commission generates. Seed entity bodies gain an `## Acceptance criteria` section, and the experiment and refinement README entity templates gain the same section. It does not change command grammar, the `--ac-scan` contract, stored frontmatter, authority, or Go runtime behavior.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - The first gate on a seed task of a freshly commissioned workflow passes `--ac-scan`, for each template and for the layer-assembly fallback, and so does the first gate on a task built from each generated README entity template.**
Verified by: `TestCommissionedEntitiesPassACScan` in `internal/contractlint` builds each entity body commission instructs: the §2b seed (shared by every template and the fallback), the SKILL's generic README entity template (the fallback's), and each `references/templates/*.md` entity snippet. It runs the real `status --read --stage implementation --ac-scan --json` on each body with a stage report citing AC-1, and each run exits 0 with `AC-1` as the first scanned AC. Baseline: on origin/main `647d55163` the seed, experiment and refinement subtests exit 1 with "no ## Acceptance criteria section in this file". Removing the AC block from §2b in `skills/commission/SKILL.md`, or from either template snippet, makes the matching subtest fail.

## Test plan

Primary proof owner: `TestCommissionedEntitiesPassACScan` (new, `internal/contractlint/commission_entity_ac_scan_test.go`). No existing test feeds commission's generated bodies to the parser; the only current test that reads the commission skill is the static `boundary_guard_control_test.go`. The new test is deterministic, runs in well under a second, and has no machine dependencies.

Falsifying edits, each shown in the ideation spike or cheap to repeat:
1. Delete the AC block from the §2b seed body: the seed subtest fails with ac-scan exit 1 (shown: the origin/main tree).
2. Delete the AC block from the experiment or refinement snippet: that template's subtest fails (shown: the origin/main tree).
3. Decorate the heading (for example `## Acceptance criteria (draft)`): ac-scan finds no section and the subtest fails, because the scanner matches only the exact heading.
4. Keep the heading but drop the `**AC-1 — …**` stub: the test fails its `AC-1` assertion even though ac-scan exits 0. This guards against a section that passes the scan with nothing to cite.

Deterministic or live: the deterministic test is enough to keep this fixed. The test cannot prove that the commissioning agent copies the fenced block into the files it writes. That is the same literal-copy step that already carries the generic README template's AC section, and the one-off live run in Risk evidence exercised it once on the candidate skill text: both seeds and the README template carried the section, and real ac-scan exited 0 at the `smoke` gate. A standing live commission lane is not justified for this change; it stays out of scope.

Validation runs `go test ./...`, `go test ./... -race` and `gofmt -l ./cmd ./internal` on the implementation branch, and re-runs the named test against an origin/main export to confirm it still fails there.

### Feedback Cycles

## Stage Report: ideation

- DONE: The task body has Problem, Proposed approach, Acceptance criteria and Test plan. AC-1 measures the end value against an origin/main baseline that fails: the first scanned gate on a seed task and on each generated README entity template passes the real `status --read --ac-scan`. It covers every template and the layer-assembly fallback, which takes its bodies from §2b and the SKILL's generic template.
- DONE: Decision on the candidate: keep `fix/commission-seed-ac` (`c3895378d`) unchanged. The reasons are in Proposed approach: the fallback needs no extra edit, success criteria and acceptance criteria stay separate in the experiment template, and the AC-1 stub prevents a heading-only scan that passes with nothing to cite.
- DONE: The riskiest path was exercised, not asserted. Real `spacedock status --read --ac-scan`, built from origin/main, exits 1 on a §2b seed body before the change and 0 after, for the development, experiment, refinement and fallback gates. The candidate test FAILs three subtests on an origin/main export and PASSes all five on `c3895378d`. `go test ./...` passes on the candidate. A one-off live `claude -p` commission run with the candidate skill exited 0. It produced seeds and an Experiment Template that carry the section, and the real ac-scan on them exits 0 at `smoke`.
- DONE: The Test plan names the proof owner, `TestCommissionedEntitiesPassACScan`, and four falsifying edits. Edits 3 and 4 were checked: a decorated heading makes ac-scan exit 1, and a heading-only section makes it exit 0 with `acs: []`. The plan states that the deterministic test is enough; the one-off live run covered the literal-copy step once, and no standing live lane is proposed.
- DONE: Expected surface is net +146 (+148/-2) across 4 files, with a tolerance of +120 to +170 and at most 5 files. Semantics: only the content of commission-generated files changes. Before/after wording is recorded for the §2b seed, the experiment snippet and Adoption line, and the refinement snippet and Adoption line. No doc diff is needed, and the reason is recorded.

### Summary

The defect is confirmed. On origin/main, a commissioned seed body, and every experiment or refinement task built from its README template, lacks `## Acceptance criteria`, so the first `--ac-scan` gate exits 1. The prepared candidate fixes all of these shapes and its deterministic test fails on origin/main, so ideation recommends landing `c3895378d` unchanged. A one-off live commission run on the candidate skill showed a real agent writing a filled AC-1 into the seeds. Two neighbouring gaps are recorded as out of scope: the loose snippet lists in `decomposed-snippets.md`, and the Claude first officer's literal `## Task Template` lookup when it files later tasks.
