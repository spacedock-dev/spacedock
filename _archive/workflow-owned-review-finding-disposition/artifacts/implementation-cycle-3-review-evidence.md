# Implementation cycle 3 review evidence

## Candidate and scheduling

- Review kind: captain-requested end-of-implementation in-stage review, not validation.
- Initial exact candidate: `2ea9ac712e23aff3bc9efa4ee4088a5aac216ccf`.
- Command: `roborev review --panel branch_final --sha 2ea9ac712e23aff3bc9efa4ee4088a5aac216ccf`.
- Enqueue output: `Enqueued job 321 (panel: branch_final, 2 reviewers) for 2ea9ac7`.
- Roborev version: `v0.62.0`.

## Job 321 exact result identity

- Panel run: `0948a2a3-152f-438a-a182-7b1bf469a06e`.
- Patch: `94576d39f652c8184aec45d9f117011fa0223188`.
- Correctness member: job `319`, Codex `gpt-5.6-sol`, verdict `F`.
- Product member: job `320`, Claude Code `opus`, verdict `F`.
- Synthesis: job `321`, Codex, verdict `F`, output:

> Changes requested — two medium-severity issues leave the fixture false-green and commissioned workflows incomplete or undispatchable.

The synthesis retained two Medium findings:

1. `docs/specs/testdata/finding-triage-materiality.tsv`: the uncited-boundary row did not independently detect accepting bare `value-ac:`.
2. `skills/commission/references/templates/development.md` and `skills/commission/SKILL.md`: normal commissioned-development generation did not atomically carry the stage selectors and complete policy section.

The product member also retained three Low observations: the lost combined uncited/Material case, no short-row fixture, and no permanent CI/Go invocation.

## Finding checkpoint and authorizations

- F1: Material, task-owned, `value-ac[AC-4]`; FO-authorized fix. Changed evidence showed one row could not independently falsify both `valid_boundary` deletion and accepting a bare citation, so the finding re-entered consultation. The FO then authorized two rows on the existing TSV surface.
- F2: Material, initially outside the approved surface, `value-ac[AC-1]` and `value-ac[AC-2]`; held and routed for captain decision. Binding reset Resolution `resolution:spacedock:rhx820qrkn6vxpday10nch36:ideation:2` authorized `skills/commission/SKILL.md` as file 11 and a retained one-off generation proof.
- F3: declined as Polish; the independent promised guards remain covered.
- F4: declined as Deferred risk; promote if supported/generated fixture input can be short or a short-row control becomes required.
- F5: declined as Deferred risk outside the approved standing-infrastructure surface; promote if unattended CI becomes the promised AC-4 enforcement path.

No candidate mutation, commit, or reviewer rerun preceded the corresponding FO authorization. No Roborev rerun occurred while F2 was held.

## Authorized corrections

- `574cfb5542784883e3b44b1acf7d52989adb5749` — `Separate uncited boundary fixture controls`.
  - Normal oracle: all 19 rows matched expected.
  - Delete only `valid_boundary` from the final check: exactly `FAIL red-uncited-boundary-validity: got accept, want reject`.
  - Replace bare `value-ac:` with accepted `value-ac[AC-1]:`: exactly `FAIL red-uncited-boundary: got accept, want reject`.
- `a0391fdbdb957d703ee4f836a5dfd5e21cc70153` — `Carry review policy through commissioning`.
  - The development Adoption pre-fill carries both implementation/validation selectors.
  - Commission generation preserves every selected-template selector and copies its complete matching named section as one unit.

## Retained one-off commission proof

Throwaway evidence root: `/tmp/spacedock-rhx820qrkn-commission-proof.UlbG51`.

The README was generated from the changed development template, and the candidate launcher was built from `./cmd/spacedock`.

- Implementation:
  - `show-stage-def` exit `0`, exactly one `## Review-finding disposition`.
  - Stage output: 29 lines, SHA-256 `89f3711f5127a379661402f6811f1c1415816842d6e96966f640c911f63c3358`.
  - `dispatch build` exit `0`, SHA-256 `d913e804ed61718e315548d9163c47ef1b16c577188e4cf07e67b8ec0664dc38`.
- Validation:
  - `show-stage-def` exit `0`, exactly one `## Review-finding disposition`.
  - Stage output: 29 lines, SHA-256 `1215fce7601be3e95d7a7ef176ee84c9eabd5ad0c61dedaa4a4fc59b9e5f1e79`.
  - `dispatch build` exit `0`, SHA-256 `dff95da3af54d6321bce7864ff8b783a9dcb32888a4529dedee946dbee2546c3`.
- No-selector negative, both stages:
  - `show-stage-def` exit `0`, policy headings `0`, stage output 5 lines.
  - `dispatch build` exit `0`; the fetched definition omits the policy.
- Selector-without-section negative, both stages:
  - `show-stage-def` exit `1`; `dispatch build` exit `1`.
  - Exact diagnostic suffix: `selector "Review-finding disposition" matches 0 headings`.

## Clean rerun

- Exact candidate: `a0391fdbdb957d703ee4f836a5dfd5e21cc70153`.
- Command: `roborev review --panel branch_final --sha a0391fdbdb957d703ee4f836a5dfd5e21cc70153`.
- Enqueue output: `Enqueued job 331 (panel: branch_final, 2 reviewers) for a0391fd`.
- Panel run: `6d7a2e58-dc9e-4d29-877b-d8730d98508f`.
- Patch: `662d6d0f1d6ba1de774a27dc654df2b3b6477354`.
- Correctness member: job `329`, Codex `gpt-5.6-sol`, verdict `P`, output `No issues found.`
- Product member: job `330`, Claude Code `opus`, verdict `P`, output `No issues found.`
- Synthesis: job `331`, Codex, verdict `P`, exact output `No issues found.`

The product member retained four non-blocking observations:

- N1: no permanent selector/section guard — declined as Deferred risk because current generation is proven and a permanent guard is excluded; promote on an observed mismatch or approved standing-guard scope.
- N2: generation checklist omits explicit pair verification — declined as Polish/Deferred risk; promote on observed omission or a promised generation-time check.
- N3: copied calibration uses `main` — declined as Deferred risk because non-main integration branches are outside the current development-template promise; promote when supported or observed.
- N4: possible refit propagation gap — declined as Deferred risk after the FO checked the shipped refit contract; promote if a supported refit proposes the section without either selector, or selectors without the section.

## Canonical round boundary

Roborev v0.62.0 exposes `show --json --job` records and `export reviews` JSON documents. It does not emit a canonical Spacedock Briefing plus reviewer JSONL pair. The supported recorder invocation requires:

`spacedock gate record <entity> --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl`

No such paths were emitted for jobs 319/320/321 or 329/330/331. Therefore no `gate record --round` package was constructed or recorded; doing so would fabricate unsupported inputs.
