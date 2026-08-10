---
commissioned-by: spacedock@template
entity-type: experiment
entity-label: experiment
entity-label-plural: experiments
id-style: sequential
stages:
  defaults:
    worktree: false
    concurrency: 2
  states:
    - name: hypothesis
      initial: true
      gate: true
    - name: smoke
      parked: true
      gate: true
    - name: run
      parked: true
    - name: analysis
      gate: true
    - name: holdout
      parked: true
      fresh: true
    - name: accepted
      terminal: true
    - name: rejected
      terminal: true
  transitions:
    - from: smoke
      to: run
      label: smoke evidence supports the hypothesis; proceed to the main run
    - from: smoke
      to: hypothesis
      label: smoke surfaces a flawed hypothesis; revise and re-smoke
    - from: smoke
      to: rejected
      label: smoke evidence is decisive against the hypothesis; reject without running
    - from: analysis
      to: holdout
      label: analysis supports the hypothesis; promote to out-of-sample holdout
    - from: analysis
      to: rejected
      label: analysis fails the hypothesis; reject
    - from: holdout
      to: accepted
      label: holdout confirms the run result; accept
    - from: holdout
      to: rejected
      label: holdout disconfirms the run result; reject
---

# Experiment Workflow Template

Each experiment is a hypothesis tested through tiers of evidence: a cheap smoke pre-flight, the main run, an analysis gate, an out-of-sample holdout, then accepted or rejected. This is the refinement shape specialized with parked tiers and a structured entity body: the parked-stages layer is active on `smoke`, `run`, and `holdout`, and the `silence-watcher` idle-mod is offered for any parked stage with timeout/nudge semantics.

Use this template when the captain's mission is to learn whether something works: hypothesis-test-learn loops, A/B tests, model evals, intervention trials, multi-tier evidence-driven promotion. Industry-of-art for the multi-tier evidence-driven promotion shape is **stage-gate** (Cooper, *Winning at New Products*); we use `experiment` as the captain-facing name for first-contact recognition and surface the lineage here for advanced captains.

## File Naming

Each experiment lives as either:

- a flat markdown file `{slug}.md` (default), or
- a folder `{slug}/` containing `index.md` plus per-tier evidence files (run logs, analysis notebooks, holdout records). The folder form is the common case for experiments because evidence files accumulate across tiers.

Slugs are lowercase, hyphens, no spaces. Example: `landing-headline-variant-b.md` or `landing-headline-variant-b/index.md`.

## Schema

Every experiment file has YAML frontmatter. Fields are documented below; see **Experiment Template** for a copy-paste starter.

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier, format determined by id-style in README frontmatter |
| `title` | string | Human-readable experiment name |
| `status` | enum | One of: hypothesis, smoke, run, analysis, holdout, accepted, rejected |
| `source` | string | Where this experiment came from (research lead, prior result, captain hunch) |
| `started` | ISO 8601 | When active work began |
| `completed` | ISO 8601 | When the experiment reached a terminal stage |
| `verdict` | enum | PASSED (= accepted) or REJECTED — set when entering a terminal stage |
| `score` | number | Priority score, 0.0–1.0 (optional) |
| `worktree` | string | Worktree path while a dispatched agent is active, empty otherwise. Once set on first dispatch into a `worktree: true` stage, it stays set across all non-terminal advancements (stickiness) and clears at terminal merge. |
| `issue` | string | GitHub issue reference (optional cross-reference) |
| `pr` | string | GitHub PR reference (set when an experiment has a paired code change) |

## Stages

### `hypothesis`

The experiment is being defined behind a captain-curated gate: a falsifiable hypothesis, the reproducible methodology that tests it, the success criteria fixed before any evidence is gathered, and the smoke pre-flight design.

- **Gate content:** Show the falsifiable hypothesis, method, fixed success criteria, and smoke design needed to decide whether testing should start.

### `smoke`

A genuinely cheap pre-flight (minutes, not days) that catches broken instrumentation or no-signal-at-all before the full run. Parked while it executes. Gate-approval to `run`, or rejection back to `hypothesis` to revise the methodology.

- **Gate content:** Show the pre-flight result, instrumentation health, observed signal, and any blocker to the full run.

### `run`

The main execution, parked while evidence accumulates. Run the same methodology that was smoke-validated and capture the evidence — raw data, sample sizes, intermediate metrics — in a form analysis can re-use.

### `analysis`

The run evidence is interpreted against the hypothesis's pre-fixed success criteria, behind a gate that distinguishes signal from noise without moving the goalposts. Gate-approval to `holdout` (verify out-of-sample), or rejection straight to `rejected`.

- **Gate content:** Show actual results against fixed success criteria, uncertainty and limitations, and the recommendation to verify on holdout or reject.

### `holdout`

Out-of-sample verification: the accepted hypothesis is re-tested on data or a setting the run did not see, parked while it executes. A `fresh` agent runs it independently of whoever ran the analysis. Holdout supports → `accepted`; holdout fails → `rejected`.

### `accepted`

Terminal state for hypotheses the holdout supported. `completed` set, `verdict: PASSED`, archived.

### `rejected`

Terminal state for hypotheses that failed at any tier (`smoke`, `analysis`, or `holdout`). `completed` set, `verdict: REJECTED`, archived — the body retains why, so future captains avoid re-running it.

## Workflow-specific rules

The FO/ensign operating contract already governs generic stage semantics and proof discipline: prefer the cheapest check that can fail — a shipped guard's run, an existing mechanical check, a one-off falsifiable exercise recorded in the report, then the captain's judgment — with new standing enforcement as the last resort rather than the default; prove by exercising rather than re-reading; fix success criteria before gathering evidence. Experiments inherit these rules from the contract their FO loads at boot; the rules below add only the experiment-shape specifics.

- **Falsifiable hypothesis, criteria fixed first.** The hypothesis must be falsifiable and its accept/reject success criteria must be written down before any evidence is gathered, so the verdict cannot be rationalized after seeing results.
- **Tier-gating.** Evidence promotes through `smoke → run → analysis → holdout` and a tier never skips: smoke fails fast before the run is burned, analysis applies the criteria as written, and a decisive smoke or analysis can reject without continuing.
- **Holdout out-of-sample independence.** The holdout must be genuinely out-of-sample (different data, population, or time window) and run by a `fresh` agent independent of the analyst — a holdout that overlaps the run sample or is interpreted by the analyst defeats its purpose.
- **Parked tiers + silence-watcher.** The parked-stages layer marks `smoke`, `run`, and `holdout` as normal-to-sit-in rather than stalled; the `silence-watcher` mod is offered for any parked stage with timeout/nudge semantics.

## Workflow State

View the workflow overview:

```bash
spacedock status --workflow-dir {dir}
```

Output columns: ID, SLUG, STATUS, TITLE, SCORE, SOURCE.

Find dispatchable experiments ready for their next stage:

```bash
spacedock status --workflow-dir {dir} --next
```

## Experiment Template

```yaml
---
id:
title: Experiment name here
status: hypothesis
source:
started:
completed:
verdict:
score:
worktree:
issue:
pr:
---

## Hypothesis

{The falsifiable claim being tested.}

## Methodology

{How the hypothesis will be tested. Reproducible.}

## Success criteria

{Fixed before evidence is gathered. What evidence accepts the hypothesis vs rejects.}

## Smoke result

{Filled in at the smoke tier — pass/fail and any instrumentation notes.}

## Run result

{Filled in at the run tier — raw evidence, sample sizes, intermediate metrics.}

## Holdout result

{Filled in at the holdout tier — out-of-sample evidence and verdict.}

## Verdict

{Final accepted or rejected, with the reasoning that drove the terminal decision.}
```

## Commit Discipline

- Commit status changes at every tier transition (smoke pass, run start, analysis verdict, holdout verdict)
- Commit experiment body updates at each tier — evidence accumulates and should be auditable

## Adoption

### Pre-fill stages

```yaml
- name: hypothesis
  initial: true
  gate: true
- name: smoke
  parked: true
  gate: true
- name: run
  parked: true
- name: analysis
  gate: true
- name: holdout
  parked: true
  fresh: true
- name: accepted
  terminal: true
- name: rejected
  terminal: true
```

### Apply layers

- **parked-stages**: fires on `smoke`, `run`, and `holdout`. These tiers wait on evidence accumulation or independent re-run rather than active in-the-loop work; the parked flag tells the FO that an entity sitting in these stages is normal, not stalled.

### Offer mods

<!-- TODO: silence-watcher mod not yet shipped — wire up when the mod lands. -->

- **silence-watcher** (idle-mod, hypothetical until the mod ships): offer when the captain indicates timeout/nudge semantics on any parked stage (e.g., "if smoke is still parked after 3 days, ping me", "if holdout has not produced a result in a week, advance to rejected"). Surface the offer in Phase 1 with this framing:

  > Parked stages are normal — an experiment can sit in `smoke` or `run` for a while because that is how evidence accumulates. The **silence-watcher** mod handles the case where parked drifts into stalled: it watches for timeout or nudge thresholds you set per parked stage, and either pings you for triage or advances the entity per the rule you defined. Without the mod, parked entities sit forever and the captain has to remember to check on them.

  Skip the offer if no parked stage has timeout semantics — silence-watcher only earns its keep when "stuck for too long" is a real concern.

### Inject entity-template snippet

Use the hypothesis-result snippet (hypothesis / methodology / success criteria / smoke result / run result / holdout result / verdict) shown in the Experiment Template section above.

### Surface variants

None. The hypothesis-test-learn shape and the multi-tier evidence-driven promotion shape are the same shape structurally — they share the smoke/run/analysis/holdout tier sequence — and both are covered by the default stage list. Captains who want a simpler 3-tier shape (skip the holdout) can drop it; that is an edit to this template, not a separate variant.

### Confirmation prose

Surface this in Phase 1 once the template is selected:

> I'll set this up as an **experiment** workflow: each {entity_label} moves through `hypothesis → smoke → run → analysis → holdout → accepted | rejected`. The `smoke` tier is a cheap pre-flight that catches broken setups before you spend the real run; the `holdout` tier is an out-of-sample check that runs with a fresh agent so it cannot be biased by whoever ran the analysis.
>
> Parked-stages layer fires on `smoke`, `run`, and `holdout` — those tiers are designed to sit waiting on evidence, not to be actively worked. If any of those parked stages has a "ping me after N days" or "auto-advance after timeout" rule in mind, I will offer the **silence-watcher** mod when we get to mod offers.
>
> Industry term-of-art for this shape is **stage-gate** (Cooper, *Winning at New Products*) — surfaced here so the lineage is discoverable. The captain-facing name stays `experiment` for first-contact recognition.
