---
id: 2ae8r33r18g0w0g21559yc57
title: Dev template ships the rigor scar tissue and refit propagates it to commissioned workflows
status: done
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.6"
sprint: 0260-proportionality
group: template
started: 2026-07-20T06:40:51Z
gates:
    version: 1
    current:
        gate: gate:docs-dev:2ae:ideation
        attempt: gate-attempt:2ae-ideation-2
    records:
        - id: gate:docs-dev:2ae:ideation
          stage: ideation
          current-attempt: gate-attempt:2ae-ideation-2
          attempts:
            - id: gate-attempt:2ae-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:single-file:b567c1211ed3a2257a92f1725c2e93bc
                digest: sha256:4760be51f28b83d92b5671119ab26916113c4b470d44603377a2a88cc2800448
                room-ref: review/ideation/briefing-1/gate-summary.md
                note: Subspace advisory float (single-file, working-copy skill launcher); the artifact is the gate summary with the frozen entity snapshot appended
              resolution:
                type: Resolution
                id: resolution:actor-1784539060005018000
                briefing: briefing:single-file:b567c1211ed3a2257a92f1725c2e93bc
                by: person:reviewer
                at: 2026-07-20T09:17:40Z
                decision: approve
                reason: "on validation gate, present the refitted delta on the workflow readme for human review"
              application:
                action: advance
                target-stage: implementation
                state: superseded
              note: "The resolution reason is a binding captain instruction for the VALIDATION gate: its presentation must include the refit diff against the workflow README for human review. Carry into the Commander package. bw's Feedback Cycles format stays deferred (surfaced in the briefing, no annotation overriding it). Superseded by attempt 2 (captain-approved staff-review folds); the approval and the validation-gate instruction both stand."
            - id: gate-attempt:2ae-ideation-2
              sequence: 2
              previous-attempt: gate-attempt:2ae-ideation-1
              state: closed
              briefing:
                id: briefing:2ae-ideation-2-chat
                digest: sha256:1f71f711733bec3fe6d6d6a243c818767938cdc78388dab2cca5056ef32f3132
                note: "chat presentation; ADVISORY digest — it hashes the working file at recording time (body folds applied, this attempt's own record excluded), which no single committed tree reproduces because an entity cannot self-bind its gates record. For drift checking, diff the entity BODY against the state commit that introduced this attempt; do not re-hash the current file. Digest refreshed once in the same fold round after the closure pass caught a residual cross-reference (Piece 6's placement parenthetical)."
              resolution:
                type: Resolution
                id: resolution:captain-chat-2ae-ideation-2
                briefing: briefing:2ae-ideation-2-chat
                by: person:captain
                at: 2026-07-20T10:20:31Z
                decision: approve
                reason: "Staff-review folds, captain-approved in chat: Piece 4 restored word-for-word to 02av's approved block (the compressed sentences returned) and moved from the validation to the implementation stage-def mirroring 02av's dev placement; the commission-skill pattern sentence names implementation-with-review-rounds stages; AC-2 reworded to a live refit-skill dry-run (a dispatched agent drives Phase 3b; validator-performed regeneration proves nothing about the skill); Pieces 1-3 re-anchor against landed sibling text at implementation."
              application:
                action: advance
                target-stage: implementation
                state: consumed
              note: "FO applied the folds directly under the captain's edit-directly grant; fable delta findings 4-6 and the codex finding-4 wave-4 sharpening. The attempt-1 validation-gate instruction (refit README delta presented for human review) carries forward unchanged."
              carried-finding:
                from: "roborev branch review of sibling member az, 2026-07-20"
                routed-by: agent:first-officer
                record: "Declined on az as out-of-surface and ROUTED HERE because this entity owns the files. Finding, verified: skills/commission/SKILL.md and skills/commission/references/templates/development.md still present `grep` as first-class proof in their AC suggestion lists and omit the falsifying-change clause, so a normal `commission` run generates acceptance criteria that CONTRADICT the policy az landed in docs/dev/README.md (evidence must name the concrete change that would make it fail; a one-off grep is structural evidence only and cannot satisfy a behavioral AC). This is the template-gap half of the original anti-tautology scope, which the 0260 re-lock moved to the template group — so it is this entity's by design, not new scope. Fold it into the Verified-by fix already in this entity's declared surface rather than treating it as an addition. IMPORTANT: the reviewer also proposed adding a cross-file consistency test to keep the templates aligned; that was DECLINED OUTRIGHT and must not be implemented — a new committed check needs explicit captain approval and its own entity, and a test asserting two instruction files agree in wording is a prose-to-prose consistency check, the banned shape. Verify against az's LANDED text, not this note."
              coverage-gap:
                raised-by: agent:first-officer
                at: 2026-07-20T14:40:00Z
                question: "Captain asked whether az's docs/dev/README.md changes consolidate into the dev workflow TEMPLATE later. FO checked the Pieces against az's landed edits and against the template's current text; recording the answer here because two gaps are real and this entity is the only propagation path."
                covered-by-re-anchor: "Piece 1 re-anchors on docs/dev/README.md's AC-template Verified-by line, which az sharpened with 'name the concrete change that would make it fail' — so it picks that up automatically. Piece 2 carries docs/dev/README.md:76 verbatim, which now includes the captain's prose-grep ruling AND the honesty-of-evidence bounding clause — so it picks those up too. Both work ONLY if implementation honours the wave-4 re-anchor rule and quotes the LANDED line rather than this entity's ideation-time snapshot."
                gap-1: "az's Edit B standalone bullet — 'Evidence must be able to fail: each AC's cited evidence names the concrete change that would flip it; an author who cannot name what would make the evidence fail has not shown it can fail, and the criterion does not count' — is carried by NO Piece. Verified receiving surface: development.md has ZERO occurrences of a can-fail rule; its 'External-proof acceptance criteria' bullet (~line 93) requires evidence from outside the task body but never asks the author to name the falsifying change. Piece 1 fixes only the AC-template stub. So a workflow commissioned today inherits the weaker rule. DECIDE DELIBERATELY: port it into the external-proof bullet, or record a decline with grounds — do not let it fall through by omission."
                gap-2: "az's Edit D — the detached audit ALSO fires on AC provenance (an AC whose expected value derives from the same package's production functions or constants), scoped to that AC's adversarial-edit check — is carried by NO Piece. Verified receiving surface: development.md DOES already carry a detached-adversarial-audit bullet (~line 94), so there is somewhere for it to land; the gap is not structural. DECIDE DELIBERATELY: port, or decline on the grounds that the provenance trigger is a dev-repo-specific sharpening a fresh workflow does not need yet."
                do-not-propagate: "NOT everything az landed belongs in the template, and porting indiscriminately would be its own error. The required-CI-lane rule and its path-to-lane mapping are a dev-lane realization tied to this repo's specific lanes, and the validation stage-def's routine-change exemption qualification is likewise docs/dev-specific. Those stay put. The generic disciplines (evidence must be able to fail; the prose-grep honesty boundary; arguably the provenance trigger) are the propagation candidates. Judge each on whether a NEWLY COMMISSIONED workflow with no CI and no lanes would be served by it."
worktree: .worktrees/spacedock-ensign-template-rigor-propagation
pr: pr-merge:542
verdict: passed
completed: 2026-07-21T04:43:00Z
archived: 2026-07-21T04:43:00Z
---

Every rigor cap in the ecosystem is post-incident scar tissue stuck in the repo where the incident happened: zaphod's no-PR-machinery rule and offline/interactive AC split, spacedock-v1's mechanism-to-value trace and materiality taxonomy. The shipped commission template contains none of them, its proof-discipline menu can only raise rigor, its interview never asks rigor level, and its "Verified by: grep" example models the exact tautology the Proof policy bans. Refit refreshes scaffolding version, not accumulated content. Template gains: `## Stakes` scaffold + interview question, materiality taxonomy, AC split, small-change fast path, size-gated semantic adversarial pass, fixed Verified-by example; refit gains content propagation so the three commissioned workflows receive the delta. Grouped with `ey`.

## Problem

Every anti-over-engineering cap in the ecosystem is post-incident scar tissue that lives only where the incident happened. The forensic audit (`_evidence/0260-agent-derail-forensics/remedy-analyses-digest.txt`, GAPS block :110-118) names the two structural gaps this entity closes:

- **The shipped template ships almost no rigor caps.** "The commission interview and trait detection identify workflow SHAPE but never ask rigor LEVEL — so a freshly commissioned README contains zero anti-over-engineering content beyond the Out-of-scope template section" (digest:112). The dev template's proof-discipline menu "only RAISES rigor; there is no rigor-capping counterpart" (digest:113). Worse, the template actively models the banned pattern: the AC-template `Verified by: {grep / test name / file path / command}` example "lists grep FIRST as the example proof, actively modeling the exact prose-grep tautology the ecosystem's Proof policy bans" (digest:106, `development.md:146` + `SKILL.md:468`). `docs/dev/README.md` already fixed this to "something outside this task body … that can fail" (`README.md:224`); the template diverged and never caught up.
- **Refit refreshes scaffolding VERSION, not accumulated CONTENT.** "No propagation loop from evolved READMEs back into the template: spacedock-v1's mechanism-to-value trace, zaphod's interactive/offline AC split, and the materiality taxonomy are template-worthy post-incident inventions that exist only where the incident already happened. Refit updates scaffolding ve[rsion]…" (digest:114). The refit spike below confirms it: today's refit against a commissioned README changes only the `commissioned-by:` line.

The natural homes already exist and are scaffolded (digest:118): the template's `## Workflow-specific rules` section and the per-stage Outputs/Bad bullets. This entity fills them with the settled sibling wording — z7's cheapest-check ordering, 02av's finding-triage taxonomy, ht's fixed Verified-by example, the proof-policy port (absorbed `proof-policy-shipped-scaffolding`), zaphod's offline/interactive AC split, and a small-change fast path — and makes refit regenerate the template's full body so a commissioned README receives the content delta, not just a version bump. The three commissioned dev workflows (this repo's `docs/dev`, zaphod, spacedock-subspace) then start protected; zaphod/subspace receive the delta via refit under their own workflows (zaphod's README is EPERM-blocked on read — that permission must be resolved there; out of this sprint's scope).

## Scope reduction (0260 re-lock 2026-07-20)

Sheds two pieces: (1) the stakes-ontology scaffold — per the reduced stakes framing, the
template scaffolds a declared-posture section only for workflows that want one, carrying
existing-declaration reach, not a new concept; (2) the project AGENTS.md router scaffold +
maintenance — deferred to a followup after the packet channel proves out (the ingestion
canary showed Claude does not read AGENTS.md, so the router serves codex ad-hoc sessions
only for now). Absorbs the commission-template half trimmed from
anti-tautology-enforcement-and-template-gap. Core kept: existing scar-tissue propagation
(materiality taxonomy, AC split, small-change fast path, "Verified by: grep" fix) + refit
carries content.

## Merged scope (adopted cross-review re-lock, 2026-07-20)

Absorbs `proof-policy-shipped-scaffolding` — porting the proof policy (no tautological string-match over LLM-ingested files) into the shipped dev-shape scaffolding is one slice of this entity's propagation payload, not a separate delivery.

## Proposed approach

Two coordinated prose changes, no code:

1. **Fill the template's already-scaffolded homes** (`skills/commission/references/templates/development.md`, with the two generic pieces echoed into `experiment.md` and `skills/commission/SKILL.md` so the wording does not diverge again — divergence is the disease). Each piece is carried VERBATIM from its approved sibling; this entity re-drafts nothing.
2. **Make refit propagate content** (`skills/refit/SKILL.md`, Phase 3b) — regenerate the commissioned README from the current template's FULL body (the workflow-independent scar-tissue sections included), not just re-render extracted stage values, so the existing "show diff" surfaces the content delta the captain adopts. The mechanism is plain `diff`; the spike proves it end-to-end.

The pieces and their verbatim sources:

- **Verified-by example fix** — carry `docs/dev/README.md:224`'s fixed clause (ht's settled wording) into `development.md:146` and `SKILL.md:468`, replacing the grep-first tautology example.
- **Proof-policy port** — carry `docs/dev/README.md:76`'s "No prose-grep over instruction files" bullet (the independent-source test) into `development.md`'s `## Workflow-specific rules` (absorbed `proof-policy-shipped-scaffolding`; the FO/ensign contract halves are az's and z7's, out of scope here).
- **Cheapest-check rename** — z7's delegated downstream propagation (`falsifiability-ladder.md:231` explicitly delegates `development.md:88` and `experiment.md:120` to the template group): replace "prefer a code gate over a prose-only rule" with z7's verbatim rename.
- **Materiality taxonomy** — carry 02av's three-class taxonomy (`ensign-finding-triage-disposition.md:137-139`) as a validation stage-def bullet (02av's stage-def-bullet pattern), and propagate the pattern to non-dev workflows via `SKILL.md`'s stage-generation guidance.
- **Offline/interactive AC split** — zaphod's split (recovered in digest:70-72): ideation splits each AC into offline-verifiable vs interactive (human/live-drive) so a plan that would automate interactive validation is visible at the ideation gate, before a harness is built.
- **Small-change fast path** — the missing tier (digest:117): routine, low-blast-radius changes scale the validation checks to the diff, mirroring the detached-audit's existing "routine changes exempt" clause.
- **Declared-posture section (opt-in scaffold)** — the parked stakes member's residue (roadmap Goal; scope-reduction above): an OPT-IN `## Workflow-specific rules` note that a workflow wanting one may declare its posture (maturity, default test depth, infra-addition policy, review-finding priority) in this section — carrying existing-declaration reach, not a new concept, and scaffolded only for workflows that want it.

**One deferral recorded, not assumed:** bw's `### Feedback Cycles` entry FORMAT is not shipped as a new template section — the roadmap re-lock defers bw's "template and docs lines … until live drives show narration ignored." The materiality bullet references the `### Feedback Cycles` record BY NAME (02av's approved reference) without shipping bw's full format string. If the captain wants bw's format in the template now, that is a one-line addition surfaced for the gate.

## Documentation changes (concrete before/after — ideation proposes, implementation applies)

Each piece names the value AC it serves, the simplest alternative, and why that alternative is insufficient. The Before/After texts for Pieces 1-3 quote the sibling sources as of ideation; implementation re-anchors each against the LANDED post-z7/az/02av wording before applying — the landed text is the verbatim source, not these quotes (staff-review fold; the Commander package's wave 4 states the same rule).

### Piece 1 — Verified-by example fix (ht). Serves AC-1.

`skills/commission/references/templates/development.md:146` AND `skills/commission/SKILL.md:468`, identical line in both:

> - Before: `Verified by: {grep / test name / file path / command a future reader can reproduce.}`
> - After: `Verified by: {test name / command output or exit code / file the change produces / resulting on-disk state — something outside this task body that a future reader can reproduce and that can fail.}`

Alt: drop only the leading `grep` token. Insufficient: the fix is not de-listing grep, it is modeling a check that CAN FAIL (ht's point); the fixed clause names an independent source, matching the `docs/dev/README.md:224` wording verbatim so the template stops diverging from the workflow that already fixed it.

### Piece 2 — Proof-policy port (absorbed proof-policy-shipped-scaffolding). Serves AC-1.

`development.md`, add one bullet to `## Workflow-specific rules` (after the repo-mutation bullet, ~line 90), carried verbatim from `docs/dev/README.md:76`:

> `- **No prose-grep over instruction files.** A string, substring, or regex match over an instruction file the model reads (the FO/ensign contract, this README, a skill) never proves a behavioral claim. The matched text was written by the same implementer the check polices, so it asserts only that the file contains what we put in it. A valid paraphrase fails it and an inverted clause passes it. To settle a case, ask whether the expected value comes from outside the file under test; if it does not, the check is a tautology and is banned. A check that binds two independent values that can diverge, such as the plugin manifest's version sharing a major.minor with the binary's version, is legitimate and is not prose-grep.`

Alt: rely on the fixed Verified-by example (Piece 1) to imply the rule. Insufficient: the example shows one good shape; the rule states the general test (independent source) that a validator adjudicates by — the captain-named primary target for the port was this template file.

### Piece 3 — Cheapest-check rename (z7's delegated propagation). Serves AC-1.

`development.md:88` AND `experiment.md:120`, same phrase in both intros:

> - Before: `… already govern generic stage semantics and proof discipline — prefer a code gate over a prose-only rule, prove by exercising rather than re-reading, …`
> - After: `… already govern generic stage semantics and proof discipline — prefer the cheapest check that can fail (a new check or enforcement process is the last resort, not the default), prove by exercising rather than re-reading, …`

The parenthetical is z7's exact downstream-propagation wording (`falsifiability-ladder.md:231`). Alt: leave the phrase (z7 changes only the contract). Insufficient: z7's own doc-diff delegates these two template lines to the template group by name; leaving them keeps the retired instruction live in every newly commissioned workflow, which "read expansively … licenses building enforcement/CI infrastructure nobody asked for" (digest:109).

### Piece 4 — Materiality taxonomy (02av's stage-def-bullet pattern). Serves AC-1.

`development.md`, add to the `### `implementation`` stage-def (mirroring 02av's dev placement — in-stage reviewer rounds happen during implementation, which is where the dutiful-fix incidents lived) a bullet carrying 02av's three classes verbatim (`ensign-finding-triage-disposition.md:137-143`):

> `- When consuming a review round's findings, triage before fixing:`
> `  - **Material** — breaks a value AC, or a declared non-negotiable boundary (safety, security, data-integrity, compatibility) reachable through the supported workflow. Fix it.`
> `  - **Correct-but-disproportionate** (deferred risk or polish) — substantively right, but no value AC breaks and its trigger is outside the supported/promised workflow. Record a decline; do not fix it. The decline is your licensed disposition, not a dodge: name the finding, its class, and why it is not material (no value AC at risk; trigger outside the promise; the condition that would promote it to material).`
> `  - **Needs decision** — a genuine product or compatibility fork. Escalate to the first officer; do not resolve it privately.`
> `  Record the disposition — which findings were fixed as material, which were declined and why — in the entity's `### Feedback Cycles` record so the gate sees it. A finding you neither fix nor record is not triaged. **Narrowing an acceptance criterion to make a finding or rejection pass is not a licensed disposition.** Declining a disproportionate finding and narrowing the claim it targets are opposite moves under the same pressure: the first leaves the product unchanged and is yours to make; the second weakens the value the entity promised and is a design-reset event requiring the captain's sign-off, recorded so it is captain-visible — never a task-internal edit.`

(Wording above is word-for-word from 02av's approved block per AC-3; staff-review fold 2026-07-20 restored the sentences an earlier draft compressed and moved the bullet from the validation stage-def to implementation.)

And `skills/commission/SKILL.md` §2a stage-generation guidance (the `### {stage_name}` Outputs bullet block, ~line 426), add one sentence propagating the PATTERN to non-dev workflows:

> `- After the existing "Stage-output bullets become checklist items…" sentence, add: A stage whose rounds consume review findings (implementation with in-stage review rounds, validation, review, evaluation) may carry a triage-before-fixing bullet: classify each review finding as material (fix), correct-but-disproportionate (record a decline, do not fix), or needs-decision (escalate) before acting — so a reviewer's non-material finding does not force a dutiful fix.`

Alt: ship only the dev bullet. Insufficient: 02av's charter for the template group is to propagate the stage-def-bullet PATTERN to non-dev workflows, which is the `SKILL.md` guidance, not just the dev file.

### Piece 5 — Offline/interactive AC split (zaphod). Serves AC-1.

`development.md`, add to the `### `ideation`` stage-def (~line 72) an Outputs note:

> `- Split each acceptance criterion by how it is verified: **offline** (a test, command, or on-disk state a fresh agent reproduces) or **interactive** (requires a human or a live drive to judge). Declare the split at ideation. A plan that would build a harness to automate an interactive AC is visible here, at the gate, before the harness is built — interactive ACs are validated by a live drive or the captain, not by new automation.`

Alt: rely on the mechanism-to-value trace to catch harness-building. Insufficient: the trace fires at implementation dispatch; the AC split makes the interactive-automation intent "visible as a violation at design time, before any harness is built" (digest:72) — the cheapest interception point, and the direct scar tissue for the 7h PTY-harness incident.

### Piece 6 — Small-change fast path. Serves AC-1.

`development.md`, add to the `### `validation`` stage-def (after the existing sentence — Piece 4 now lands in the implementation stage-def):

> `- **Small-change fast path.** Scale the validation checks to the diff's blast radius. A routine, low-blast-radius change (a doc line, a one-line fix, a rename) does not need the full checklist or the detached adversarial audit — the same "routine changes exempt" carve-out the audit already grants, applied to validation as a whole. Match the rigor to the change; a trivial diff over-validated is its own waste.`

Alt: leave validation uniform. Insufficient: "Validation/review checklists lack a small-change fast path in all three READMEs" (digest:117); a uniform checklist makes "a compliant validator on a trivial task over-produce" (digest:107) — this is the proportionality thesis applied to the review stage.

### Piece 7 — Declared-posture section, opt-in scaffold (parked-stakes residue). Serves AC-1.

`development.md`, add to `## Workflow-specific rules` an opt-in note (not a mandatory section):

> `- **Declaring a posture (optional).** A workflow that wants a single findable answer to "how much engineering does this project want?" may declare it here: project maturity (prototype / product), default test depth, infra-addition policy (may a worker add a CI lane or lint unasked?), and review-finding priority. This is a place to write an existing posture down, not a new required concept — omit it unless the workflow benefits from a stated posture.`

Alt: a mandatory posture section, or a new interview question. Insufficient/over-reach: the stakes member is PARKED (roadmap Goal) — a new ontology or required field is exactly what the re-lock removed; the residue is an OPT-IN scaffold carrying existing-declaration reach, nothing more.

### Piece 8 — Refit content propagation. Serves AC-2.

`skills/refit/SKILL.md` Phase 3b (`### 3b. README (Show Diff)`, ~lines 96-110). Today step 1 reads "Generate what the current commission template would produce for this workflow, using the extracted values (mission, stages, schema, etc.)." — ambiguous enough that a refit re-renders only the workflow-specific stage values and never re-emits the template's workflow-independent sections, so scar-tissue content never reaches the diff. Sharpen step 1:

> - Before (step 1): `Generate what the current commission template would produce for this workflow, using the extracted values (mission, stages, schema, etc.).`
> - After (step 1): `Re-select the matching template for this workflow's shape (development / experiment / refinement) and generate its FULL body for this workflow, substituting the extracted values (mission, stages, schema). Regenerate the workflow-INDEPENDENT sections too — `## Workflow-specific rules`, the proof/triage bullets, and the `## {Entity} Template` / `Verified by` example — not only the stage prose. These carry the accumulated scaffolding content; a refit that re-renders only the stage values propagates the version stamp but not the content the template gained since this workflow was commissioned.`

And add one line to the Phase 3b diff-presentation prose so the captain sees content deltas as adoptable:

> - After the "I have NOT modified your README" line, add: `Additive sections the current template gained since your workflow was commissioned (new proof/triage bullets, a fixed example) appear as additions in this diff — those are the accumulated scaffolding content, safe to adopt; hunks that touch your customized stage prose are yours to judge.`

Alt: leave Phase 3b as-is (it already regenerates and diffs). Insufficient: the spike's control proves today's path surfaces only the `commissioned-by:` line — the regeneration must explicitly re-emit the workflow-independent sections, or content does not propagate. Alt: a new refit sub-command that injects sections. Over-reach: new machinery for what `diff` + a fuller regeneration already do; the captain ruling reserves new mechanisms to their own approved entity.

## Riskiest-mechanism spike (done first, per ideation policy)

**Claim under test (what invalidates the rest if false):** refit can propagate CONTENT (not just the version stamp) to an already-commissioned README, via the existing regenerate-and-diff path, on the smallest real case.

**Exercised end-to-end** (files in the ideation scratch; recreatable — a dev-shape README as commissioned from the CURRENT template genuinely lacks the scar tissue):

1. `site-workflow-old.md` — a minimal dev-shape README stamped `commissioned-by: spacedock@0.25.0`, carrying the tautological `Verified by: {grep …}` example, the "prefer a code gate" intro, and no triage bullet (i.e. a workflow commissioned from today's template).
2. Control (today's refit) = the same file with only the stamp bumped to `@0.26.0`. `diff -u old control` → **2 changed lines, the `commissioned-by:` pair only.** This reproduces the reported gap: refit refreshes version, not content.
3. Content propagation = the same README regenerated from the UPDATED template (Pieces 1-6 applied). `diff -u old regenerated` → **13 changed lines across 4 hunks**: the cheapest-check rename, the no-prose-grep bullet (additive), the materiality triage bullet + small-change fast path (additive), and the fixed Verified-by line — a reviewable, mostly-additive diff.

**Determination:** the mechanism holds. The proven mechanism it rides is POSIX `diff` over a full-body regeneration; no parser round-trip, runtime handoff, on-disk format, or tool-flag support is involved. The one required change is refit prose (Piece 8): the regeneration must re-emit the workflow-independent sections, or the diff collapses to the version-only control. The spike seeds AC-2's validation replay directly (the two diffs are the fixture: content-delta > 0, control == version-only).

## Acceptance criteria

Each AC names a property of the finished entity and how it is verified. Value ACs measure the end-value against a baseline that can move the wrong way.

**AC-1 (VALUE) — A workflow commissioned from the updated dev template ships the anti-over-engineering scar tissue that a workflow commissioned from `main` ships zero of, and models the banned tautology zero times.**
The generated README carries the three-class finding-triage taxonomy, the offline/interactive AC split, the small-change fast path, the no-prose-grep rule, and an AC-template `Verified by` example that names an independent source; it contains no grep-first proof example. The numbers that move the wrong way: on `main`, count of shipped scar-tissue pieces = 0 and count of banned grep-first proof examples ≥ 1; after, scar-tissue pieces = 6 and banned examples = 0.
Verified by: a live commission drive on a scratch mission ("a personal-site development workflow") under the branch vs `main`; the two generated READMEs compared and the counts pasted into the validation report (behavior observed from a real commission run, not a committed prose-grep). This is the DoD line "a scratch workflow commissioned from the updated template contains the materiality taxonomy and the fixed Verified-by example." Independent baseline: the `main` commission, which the digest predicts (and AC re-confirms) contains none of the six and models the tautology.

**AC-2 (VALUE) — Refit propagates CONTENT, not just the version stamp: a refit dry-run against a README commissioned from the pre-scar-tissue template surfaces the scar-tissue sections as an additive content diff, whereas today's refit surfaces only the `commissioned-by:` line.**
The number that moves the wrong way: content lines in the refit diff beyond the version stamp — 0 today, >0 after (the spike measured 2 control vs 13 propagated).
Verified by: a live refit dry-run at validation — a dispatched agent DRIVING `skills/refit/SKILL.md` Phase 3b (not the validator regenerating by hand) against the commissioned fixture README; assert the emitted diff contains the materiality-taxonomy and fixed-Verified-by hunks AND that the version-only control diff contains only the stamp line. The claim under test is the refit skill's own behavior after Piece 8 — validator-performed regeneration would reproduce the ideation spike and prove nothing about the skill (staff-review fold, per the live-drive proof rule). Independent baseline that can move wrong: the version-only control (today's behavior) surfacing 0 content lines. This is the DoD line "a refit dry-run against a commissioned README shows the content delta arriving." One-off drive recorded in the report, not a committed test.

**AC-3 (mechanism — the pieces land verbatim, no divergence, no minted terms, no new enforcement).** Each of the seven doc-diff pieces is applied to its named file; z7's cheapest-check clause, 02av's three-class taxonomy, ht's Verified-by wording, and `docs/dev/README.md:76`'s no-prose-grep bullet appear in the template WORD-FOR-WORD against their sibling source (the anti-divergence property this entity exists to protect); the shipped diff adds zero `*_test.go`, zero lint/CI/gate, and touches only `development.md`, `experiment.md`, `skills/commission/SKILL.md`, and `skills/refit/SKILL.md`.
Verified by: a one-off validation comparison of each shipped piece against its cited sibling line (output pasted into the report, not a committed grep — per the captain's prose-grep ruling); plus `git diff --name-only` showing the four-file set and no added test/gate. Paired with AC-1/AC-2 (this is the means; those measure the value). A divergent re-draft, a fifth touched file, or any new committed check fails it.

## Test plan

- **AC-1 value drive (live commission) → medium.** Two live `commission` drives (branch vs `main`) on one scratch mission; compare the generated READMEs and count scar-tissue pieces / banned examples. Commission is an LLM skill, so this is a live behavioral drive observed once at validation, not a committed test. Owned by validation per the live-drive proof rule.
- **AC-2 live refit dry-run → medium.** A dispatched agent drives the updated refit skill Phase 3b against the fixture README; assert content hunks present and the control is version-only. One live skill drive at validation; the spike's fixtures seed it.
- **AC-3 verbatim + git-state check → low.** A one-off read comparing each shipped piece to its sibling source, and `git diff --name-only` for the four-file set with no added test/gate. Output pasted, never committed.
- **No Go/product code, no new committed test, gate, or lint in this cut.** The deliverable is template + skill prose; the negative deliverable (no new enforcement) is the point and is a git-state check (AC-3). High-stakes note: the commission template and refit skill are shipped scaffolding (a high-stakes surface per the Proof policy), so the detached adversarial audit applies at validation before merge.

## Expected surface + tolerance (declared, per captain ruling)

- `skills/commission/references/templates/development.md`: ~7 edits (cheapest-check rename ~0 net lines; no-prose-grep bullet ~1 line; Verified-by fix ~0 net lines; materiality bullet ~6 lines; AC split ~1 line; small-change fast path ~2 lines; declared-posture opt-in ~2 lines) → ~12 net lines.
- `skills/commission/references/templates/experiment.md`: cheapest-check rename, ~0 net lines (1 line replaced).
- `skills/commission/SKILL.md`: Verified-by fix (~0 net) + stage-def-bullet-pattern sentence (~2 lines) → ~2 net lines.
- `skills/refit/SKILL.md`: Phase 3b regeneration sharpening + one diff-presentation line → ~4 net lines.
- **4 files, ~18 net lines of template/skill prose, 0 Go source, 0 product LOC, 0 new committed tests/gates/lints.** Fixtures: 1 dev-shape README fixture (+ its version-only control) for AC-2, seeded from the spike.
- **Tolerance: 2×** (≤ ~36 net prose lines). **Hard self-check:** any Go/product code, any new committed check/gate/lint, a fifth touched instruction file, a re-drafted (non-verbatim) sibling wording, or shipping bw's deferred `### Feedback Cycles` format without a captain yes trips a reconfirm — carrying settled wording verbatim into the already-scaffolded homes with zero new machinery IS the point.

## Stage Report: ideation

- DONE: Concrete before/after doc diffs for every propagated piece (materiality taxonomy, AC split, small-change fast path, Verified-by example fix, proof-policy port into shipped scaffolding, declared-posture section as opt-in scaffold) — each naming the value AC it serves, no "change X" hand-waving
  `## Documentation changes` carries Pieces 1-7 (plus the z7 cheapest-check rename z7 delegated to the template group, and Piece 8 refit), each with exact BEFORE/AFTER text against a named file:line, the value AC served, and the simplest-alternative-and-why-insufficient.
- DONE: Riskiest mechanism exercised first — refit propagating CONTENT (not just version) to a commissioned README, demonstrated end-to-end on the smallest real case
  `## Riskiest-mechanism spike`: a dev-shape README commissioned from today's template vs regenerated from the updated template; today's refit control diff = 2 lines (`commissioned-by:` only), content-propagation diff = 13 lines across 4 hunks. Mechanism = POSIX `diff` over a full-body regeneration; the one required change is refit prose (Piece 8). Seeds AC-2.
- DONE: Declared expected surface + tolerance in the task body, settled sibling wording carried verbatim (no re-drafting, no minted terminology)
  `## Expected surface + tolerance`: 4 files, ~18 net prose lines, 0 product LOC, 0 new tests/gates, 2× tolerance with a hard self-check. Verbatim sources cited per piece (z7:231, 02av:137-139, ht/`README.md:224`, proof-policy `README.md:76`); bare ordinals; no coined abstractions. AC-3 makes verbatim-ness and the no-new-enforcement negative a checkable property.

### Summary
Ideated the template-rigor propagation as two prose changes: fill the commission template's already-scaffolded homes with six settled scar-tissue pieces (carried verbatim from z7, 02av, ht, and the proof-policy/AC-split/fast-path sources, echoed into `experiment.md`/`SKILL.md` so wording stops diverging), and sharpen refit Phase 3b to regenerate the template's full body so a commissioned README receives the content delta, not just a version bump. The riskiest path — refit propagating content — was spiked end-to-end on the smallest real case: today's refit surfaces only the version stamp (2-line control) while regeneration from the updated template surfaces 13 content lines, via plain `diff`. Two VALUE ACs measure the end-value against baselines that move the wrong way (commissioned README ships 0 scar-tissue pieces + ≥1 tautology on `main`; refit surfaces 0 content lines today). Key decisions recorded for the gate: bw's `### Feedback Cycles` format is deferred per the re-lock (referenced by name, not shipped) and surfaced for a captain call; az's Edit D and the FO/ensign contract halves stay out of scope.

## Stage Report: implementation

- DONE: Pieces 1-3 are re-anchored against the LANDED sibling text, not the ideation-time quotes in your body: az and ht are MERGED to main, so read docs/dev/README.md as it now stands; z7 is approved and in PR #540 but NOT yet merged, so take its rename wording from its branch and say plainly in your report that this one source was pre-merge
  z7 MERGED as #540 mid-stage, so I re-anchored it a second time against merged main — `main:docs/dev/README.md:74` is byte-identical to the branch text, so no source in this cut is pre-merge. Carriage verified by exact string containment against merged main and 02av's body, not by eyeball: Piece 1 (233 chars), Piece 2 (1547), Piece 3 (252), gap-1 (250), gap-2 (323), and all five of 02av's Piece 4 fragments (173/411/123/212/456). Commit `54e26fd6`.
- DONE: The commission templates stop contradicting the policy az landed: `grep` is no longer offered as first-class proof in the AC suggestion lists, and the falsifying-change clause is present — this is the carried finding in your gates record, folded into the Verified-by fix rather than treated as new scope
  Both AC suggestion lists (`development.md`, `commission/SKILL.md`) now carry az's landed clause ending "…and that can fail; name the concrete change that would make it fail."; zero `{grep /` occurrences remain in either file. Folded into Piece 1, no new scope. The reviewer's proposed cross-file consistency test was NOT implemented, per the record's outright decline.
- DONE: A LIVE refit dry-run proves AC-2 — a DISPATCHED AGENT drives skills/refit/SKILL.md Phase 3b against the commissioned fixture and the emitted diff carries the content hunks, with the version-only control diff showing only the stamp line
  Four dispatched-agent drives, neutral prompts, treatment on the branch and control on a detached `main` checkout. Treatment: 29 changed lines carrying the materiality-taxonomy hunk, the fixed Verified-by hunk, the AC split, the fast path, and the whole `## Workflow-specific rules` section. Control: zero scar-tissue content lines in every run. The control was NOT literally stamp-only — see the finding below. Falsifying change: revert Piece 8's step-1 sharpening and the treatment diff collapses to the control's shape.
- DONE: Coverage-gap decisions, both PORTED, neither by omission
  gap-1 (evidence must be able to fail) landed as its own bullet mirroring az's README ordering; gap-2 (AC-provenance trigger) appended to the existing detached-audit bullet. Each drops one trailing sentence that cross-references docs/dev-only machinery ("the validation bar below"; the routine-change exemption the record marks do-not-propagate). Shipped words are word-for-word; nothing re-drafted. The do-not-propagate items (required-CI-lane rule, path-to-lane mapping) stayed put.

### Findings

**The control drive is not stamp-only, and that is a real result, not noise.** AC-2 predicted today's refit surfaces only the `commissioned-by:` line — the ideation spike measured a 2-line control. A live control drive does not reproduce that. Three control runs produced 107, 4, and 23 changed lines; the variance is entirely which template source the agent picks for the entity-template snippet, because `main`'s Phase 3b step 1 ("what the current commission template would produce") does not say. The load-bearing half of the claim holds in every run: zero scar-tissue content lines, and the `Verified by: {grep …}` tautology survives untouched. So the baseline moves the wrong way as AC-2 requires, but the hand-built spike control was cleaner than reality and the AC's "only the stamp line" wording overstates it. Validation should assert zero content lines, not a 2-line diff.

**`## Workflow-specific rules` is not a section commission emits.** Chasing the first control run surfaced this: `commission/SKILL.md` §2a lists the required README sections and that heading is not among them; SKILL.md:116 uses the template only for its `## Adoption` section. Zero real commissioned READMEs in this repo carry the heading — `docs/dev/README.md`'s `## Proof policy` is hand-added. Consequence for AC-1: of its five named items, four (triage taxonomy, AC split, fast path, Verified-by example) sit in sections §2a does emit; the no-prose-grep rule does not, and neither do Piece 7 and both gap ports. Refit propagates them anyway — the treatment drive emitted the whole section, because Piece 8's step 1 names it — but a FRESH commission would not. I considered adding the section to §2a (~3 lines, inside my declared file set) and DECLINED to: AC-2 is proven without it, the DoD's two named items are emitted without it, and making commission emit a new required section is a behavioral change the ideation gate approved on a premise it never tested. AC-1's live commission drive is validation-owned and is the right place to measure this. Surfaced for the gate as a decision, not silently patched. **SUPERSEDED BY CAPTAIN RULING (see cycle 2 below).** My reasoning above was sound on the question I was answering — AC-2 — and I preserve it unchanged. But the ruling was decided on REGRESSION grounds, not AC-2: the FO took my exact position ("land 2ae as-is, record AC-1 as failed with grounds, file the commission fix next-train") to the captain as one of three options, and the captain chose "fix the wiring in 2ae" instead. Provenance found after I asked shows why: the slot is not a new concept the gate approved untested — PR #388/`b9aad379` (2026-06-16) created it deliberately and §2a never listed it (`git log -S` → 0), so a freshly commissioned workflow has shipped none of the dev-shape rigor block for five weeks. The §2a edit is now applied (cycle 2). Decline and its supersession are both visible here by design.

**Piece 4 ships a rule this workflow's own parent does not carry.** 02av was parked for delivery-mechanism reasons (bytes on a ratcheted surface), not because its rule was wrong, and its approved block survives verbatim in its entity body — so AC-3's word-for-word carriage still has its source. But after this lands, a workflow commissioned from the template carries a finding-triage taxonomy that `docs/dev` itself does not, because 02av's delivery into the parent was parked. That asymmetry is real and resolves when the 3k advisory gate-resolution record lands. Recorded here so it is not discovered later as a surprise.

### Summary

Carried the sprint's settled wording verbatim into the four declared files: the fixed Verified-by example and the no-prose-grep rule from az's merged README, z7's cheapest-check clause into both templates, 02av's three-class triage taxonomy into the implementation stage-def plus a pattern sentence for non-dev workflows, zaphod's AC split, the small-change fast path, an opt-in declared-posture note, and both coverage-gap ports. Refit Phase 3b now regenerates the template's full body, which four dispatched-agent drives confirm propagates the content delta where today's refit propagates none. Surface: 4 instruction files, +24/-7 = 17 net prose lines against ~18 declared and a 36-line tolerance, plus the declared AC-2 fixture; zero Go, zero product code, zero new committed tests, gates, or lints, and bw's `### Feedback Cycles` format referenced by name only. `go test ./...` green. Three findings above are for the gate: the control drive is noisier than AC-2 predicted, `## Workflow-specific rules` never reaches a freshly commissioned workflow (declined to fix unilaterally), and Piece 4 now outruns its own parent workflow.

## Stage Report: implementation (cycle 2 — §2a slot wiring)

Captain ruled YES on the finding raised in cycle 1: add `## Workflow-specific rules` to commission §2a's required-section list. Landed as commit `a45b4ec0`, inside the declared four-file surface. Combined surface is now 21 net prose lines (28+/7−) across the four instruction files — inside the 36-line tolerance, no fifth file, still zero Go/tests/gates/lints.

- DONE: `## Workflow-specific rules` emitted from §2a, sourced from the selected template's own slot
  `commission/SKILL.md` §2a now lists the section between `## Scoring` and `## Workflow State`, instructing the generator to copy the selected template's slot verbatim. Falsifying change: delete that section-block and a freshly commissioned README ships none of the dev-shape rigor.
- DONE: provenance verified independently before it went in the commit message
  `git log -S 'Workflow-specific rules' -- skills/commission/SKILL.md` → 0 commits (§2a never named the slot, ever). The slot was born in `b9aad379` / PR #388 on 2026-06-16 ("defer universal rules to the FO/ensign contract"). So this is finishing wiring #388 left undone five weeks ago, not a new behavioral concept.
- DONE: AC-2 re-proven by two live refit drives against the corrected faithful §2a-shaped fixture, grounded in the emitted diff files (not the agents' summaries)
  Treatment (my branch): 35 changed lines; every scar-tissue marker is an ADDED line — the `## Workflow-specific rules` heading, no-prose-grep, evidence-must-fail, triage-before-fixing, small-change fast path, declaring-a-posture, the offline AC split, and the fixed Verified-by "…and that can fail; name the concrete change". Control (pre-fix main): 25 changed lines but ZERO scar-tissue content lines (4 added lines, all Task-Template scaffolding churn), and the `{grep / test name …}` tautology SURVIVES in the regenerated README. Baseline moves the wrong way exactly as AC-2 requires.

### Additional finding (the sprint's own thesis, in the sprint's own tooling)

`internal/contractlint/template_defer_test.go`'s `TestTemplatesCarryWorkflowSpecificRulesSlot` asserts each TEMPLATE file contains the `## Workflow-specific rules` heading. It has been green for five weeks — while the section it guards reached zero commissioned workflows, because commission never emitted it. The check guards the CONTAINER (the template has the slot) and never the DELIVERY (a commissioned README receives it). That is the exact failure mode this sprint exists to name: a green check whose passing says nothing about the value shipping. A delivery check — a test that a commissioned README carries the slot — was explicitly considered and DECLINED here to keep the "zero new committed tests/gates/lints" line intact; team-lead is filing it as next-train work. Recorded, not overlooked.

### Scope note

The §2a fix is prose whose behavioral proof is a live COMMISSION drive (branch vs `main`), which is AC-1's validation-owned verification — not re-run here by design. The refit treatment drive above demonstrates the same slot content is followable and lands correctly when an agent is told to emit it; it does not itself exercise commission's §2a. Cycle-1 finding-1 stands unchanged: a live refit control is not literally stamp-only (scaffolding churn varies by run), but it carries zero scar-tissue content lines in every run, which is the load-bearing half of AC-2's baseline. The Piece-4 asymmetry recorded in cycle 1 is likewise unchanged.

### Human-reviewable refit delta (for the validation gate's binding captain condition)

Reproducible in one dispatched drive: drive `skills/refit/SKILL.md` Phase 3b against `fixtures/refit-content-propagation/site-workflow/README.md` on this branch; the emitted README diff IS the delta a human reviews. Against a pre-fix `main` checkout of the same fixture, the control diff carries only stamp + scaffolding churn and no scar tissue. The fixture and its explainer are committed for exactly this.

### Feedback Cycles

FO-owned per feedback-rejection-flow. 2ae exercises the correction-round convention `bw` shipped this sprint on its own in-stage rounds; deviation is measured against the ~18-line ideation estimate, never the prior round.

- Cycle 1: changes-requested — FO branch review; surface 17 vs estimate 18 (94%); AC unchanged. The cycle-1 stage report recorded a reasoned decline to wire §2a; the captain ruling (fix the wiring, on regression grounds) superseded it. Bounced back for the wiring plus the report correction.
- Cycle 2: passed — FO ruling applied (§2a slot wiring, `a45b4ec0`); surface 21 vs estimate 18 (117%); AC narrowed: AC-2's control was specified as "the version-only control diff contains only the stamp line," but a live refit control carries variable Task-Template scaffolding churn (25 / 107 / 4 / 23 changed lines across runs) with ZERO scar-tissue content lines every time — narrow the control assertion to "zero scar-tissue content lines," which is the load-bearing invariant the 2-line figure only approximated. Cause of the variance is `main`'s refit step-1 template-source ambiguity (a finding about `main`, not this change). Within 2× tolerance, so the design-reset decision is RECONFIRM — no re-scope, park, or escalate.

FO note (recorded against this FO, not buried): the commission drive I ordered at implementation was declined by the ensign on the ground that AC-1's live commission drive is validation-owned and an implementer-run proof weakens its independence — the same principle the staff review used to keep AC-2's refit drive out of the validator's own hands. The ensign was right and I withdrew the instruction. AC-1's live commission drive (branch vs `main`, scar-tissue and tautology counts) is carried to the validation gate as a detached-validator requirement so it does not slip.

## Stage Report: implementation (cycle 3 — commission drive proves AC-1 here)

Per the captain ruling, ran the commission path so AC-1 stops being validation-deferred and becomes falsifiable in this stage. A dispatched agent drove `skills/commission/SKILL.md` (this branch) Phases 1-2 against a scratch mission ("a personal-site development workflow") into a throwaway dir, with a NON-LEADING prompt — it was not told what sections to emit, only to follow the skill. Trait detection selected `development`. I then grepped the generated README independently.

- DONE: DoD line "a scratch workflow commissioned from the updated template contains the materiality taxonomy and the fixed Verified-by example" — now runnable HERE, and it PASSES
  Three greps on the generated `README.md`:
  1. `## Workflow-specific rules` slot: PRESENT (line 131) — the §2a wiring fires. The full slot landed: cheapest-check clause, repo-mutation layer, no-prose-grep (verbatim), evidence-must-fail (verbatim), the opt-in proof-discipline block, and the declaring-a-posture note.
  2. Three-class materiality taxonomy: PRESENT (implementation stage Outputs) — all three classes named (material / correct-but-disproportionate / needs-decision) plus "record the disposition … in the `### Feedback Cycles` record".
  3. Fixed Verified-by: PRESENT (verbatim) and the `{grep / test name …}` tautology occurs ZERO times.
  Falsifying change: revert the §2a edit and the slot vanishes from the generated README (the pre-fix control refit drive and the pre-#388 state both show its absence); revert Piece 1 and the grep-first example returns.

### The honest nuance validation must know: commission adapts, it does not transcribe

The commission skill is an LLM generator; §2a instructs it to craft mission-specific content, and my slot instruction says to adapt mission-specific wording. So the generated README is an ADAPTATION, not a byte copy of the template. Concretely, in this run: the no-prose-grep rule, evidence-must-fail bullet, declaring-a-posture note, and the Verified-by example came through verbatim; but the triage taxonomy was compressed into one implementation-stage Outputs sentence (all three classes preserved), the AC split was paraphrased into ideation Outputs, the small-change fast path became "rigor matches the change's blast radius" in a validation Good bullet, and the AC-provenance audit trigger (gap-2) was paraphrased away entirely — replaced by a mission-specific "over-fires on ordinary content or styling edits" clause. This is expected behavior for an LLM commission step, not a defect: verbatim carriage is a property AC-3 pins on the TEMPLATE (`development.md`), which holds; the commissioned artifact is judged by CONTAINS (the DoD word), which passes. Validation should grep the commissioned README for the taxonomy's substance and the fixed example, not for byte-identity with the template — a byte-identity assertion on a generated artifact would fail for reasons that are correct LLM behavior. One consequence worth stating: a template bullet can be paraphrased out of any single commission (as gap-2's provenance trigger was here), so "the template carries X" and "every commission emits X verbatim" are different claims; only the first is guaranteed.

### Summary

Applied the captain's ruling (superseding my AC-2-grounded decline, which is preserved and marked superseded in the cycle-1 findings): §2a now emits `## Workflow-specific rules`, committed `a45b4ec0`, finishing wiring PR #388 left undone on 2026-06-16. A non-leading commission drive confirms the DoD line HERE — the generated README carries the slot, the three-class taxonomy, and the fixed Verified-by with zero grep-first tautology — closing the exact hole (nobody ran the commission path) that let a five-week regression survive. Surface holds at 21 net prose lines across the four instruction files (inside 36 tolerance, no fifth file, zero Go/tests/gates/lints). The one nuance for validation: commission adapts template prose to the mission, so judge the commissioned README by contains-substance, not byte-identity. The AC-2 narrowing is recorded by the FO in `### Feedback Cycles`; findings 2 and 3 (delivery-vs-container lint gap, Piece-4 asymmetry) stand as written.

### Roborev branch review (pre-validation, FO-triaged 2026-07-21)

`roborev review --branch --panel branch_final` (job 397, codex, roles correctness+product) returned **changes requested — three medium**. FO triage under the sprint posture: **all three DECLINED with grounds; 2ae not re-dispatched.**

1. **"Held feedback-disposition mechanism is being propagated" (`development.md:76`) — DECLINED.** The template references `### Feedback Cycles` exactly once, BY NAME, and ships zero instances of bw's deferred `- Cycle {N}: …` entry format (verified: `git diff main...HEAD` over the commission/refit surface finds the one-liner nowhere). `dispatch-sprint-execution.md:55` explicitly permits this: "bw's `### Feedback Cycles` entry FORMAT in the template: deferred by the re-lock; 2ae references the record by name only." The triage TAXONOMY shipping is DoD-required and captain-blessed — 2ae ships Piece 4 despite 02av's park because the DoD states the materiality-taxonomy requirement independently. The parent-vs-template asymmetry the reviewer names (a commissioned workflow carries a triage rule `docs/dev` lacks) is the Piece-4 asymmetry the stage report already records; the reviewer re-finding it independently corroborates recording it prominently. It resolves when the 3k advisory-resolution reframe lands. Not a scope violation; does not trip 2ae's hard self-check (bw's format is not shipped). The reviewer reasoned from wave-2/3 text without the DoD line or the 2ae-scope §3 carve-out — reasonable given that, settled with it.

2. **"Refit fixture does not verify full-body propagation" (`fixtures/refit-content-propagation/README.md:10`) — DECLINED.** AC-2's proof is a validation-owned LIVE refit Phase-3b drive with the version-only control as the negative case that moves the wrong way — that negative case EXISTS (control surfaces zero scar-tissue content lines). The fixture is INPUT for that live drive, not a self-contained committed test. The reviewer's requested fix — a committed smoke scenario comparing generated to hand-authored expected README — freezes LLM output, which is the "test that tests mocked behavior" anti-pattern this sprint exists to remove.

3. **"Commission's emission path lacks behavioral coverage" (`SKILL.md:436`) — DECLINED.** Same class. AC-1's proof is the validation-owned live commission drive (branch vs `main`). A committed "compare generated README to a hand-authored expected section" fixture is either frozen-LLM-output or a prose-to-prose comparison near the banned grep shape. The delivery check — a durable test that a commissioned README carries the slot — was explicitly considered and DECLINED this sprint, filed next-train. LLM-skill behavior is proven by live drive at validation, never a frozen committed fixture.

Findings 2 and 3 are the sprint's thesis inverted: an adversarial reviewer's honest instinct ("add durable behavioral coverage") reaching for exactly the committed prose-comparison / frozen-output check 0260 retires. The right coverage is the live drive, scheduled at the validation gate, with the human-reviewable refit delta as the binding captain condition.

## Stage Report: validation (retroactive gate record, 2026-07-21)

Recorded after the fact: 2ae was merged (#542) and archived under the captain's "just merge it" / "approve both" direction WITHOUT a recorded validation gate — the gate pointer still named ideation and the attempt-1 binding condition ("on the validation gate, present the refitted delta on the workflow readme for human review") had no recorded disposition. The 0260 pre-cut staff seats (codex `pre-cut-staff-review.md` blocker 2; fable `pre-cut-staff-review-fable.md`) flagged this. This record closes it evidence-only (the product is already merged and independently validated; no product correction needed).

**Independent validation evidence (fresh, detached — not the implementer's self-report):**
- **AC-1 (commission):** a fresh non-leading commission drive produced a README carrying the three-class materiality taxonomy, the fixed `Verified by` example verbatim, and ZERO grep-tautology occurrences (fable seat's detached drive; corroborated by the implementation-stage commission drive that motivated the §2a slot-wiring fix).
- **AC-2 (refit):** a fresh agent drove `refit` Phase 3b against the committed fixture `fixtures/refit-content-propagation/` and emitted the full content delta (the `## Workflow-specific rules` slot + no-prose-grep + evidence-must-fail + triage + fixed Verified-by); the version-only control moves the wrong way (stamp-only, zero scar-tissue content). AC-2's narrowing (control asserts zero scar-tissue content lines, not a literal 2-line diff) is recorded in the `### Feedback Cycles` entry.

**Captain-required human review (the binding condition):** the refit README delta (control 1-line stamp-only vs treatment +21/−2 carrying the full rigor payload) was presented to the captain 2026-07-21 and **APPROVED**.

**Verdict: PASSED.** Both AC drives pass on fresh detached agents; the captain reviewed and approved the refit delta; the Piece-4 asymmetry (a commissioned workflow carries the finding-triage rule docs/dev lacks until 3k lands) stays recorded, not papered over.
