---
id: 9x6xw292fsz1b4648x9hn40y
title: Make shipped contract content self-contained
status: validation
source: "Captain review of the 0.27 stack + audit-r2 (2026-08-15); captain directive: file, dispatch off stack tip, PR as stack layer"
started: 2026-08-15T18:35:12Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-make-shipped-contracts-self-contained
issue:
gates:
    version: 1
    records:
        - id: gate:9x6xw292fsz1b4648x9hn40y:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:9x6xw292fsz1b4648x9hn40y-backlog-1
              briefing:
                id: briefing:9x6xw292fsz1b4648x9hn40y:backlog:attempt-1:revision-1
                digest: sha256:de4553c953de929a52bfb362a7d19a5ac077470125db6f1fad8ab8263e978581
                request-digest: sha256:3d3d92fe0a1b05ad75500987f12bc7e7d65f3b1e395d76480f0691267f354cc0
                room-ref: ./make-shipped-contracts-self-contained/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:9x6xw292fsz1b4648x9hn40y:backlog:1
                briefing: briefing:9x6xw292fsz1b4648x9hn40y:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T18:15:19.816976Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: file, dispatch based off stack tip, PR on top of the stack'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:9x6xw292fsz1b4648x9hn40y:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:9x6xw292fsz1b4648x9hn40y-ideation-1
              briefing:
                id: briefing:9x6xw292fsz1b4648x9hn40y:ideation:attempt-1:revision-1
                digest: sha256:0886859da6795882e3972e6a871f8517108ada59ff8a4464701b17a172eb9daa
                request-digest: sha256:e1b1e98ed917f7291bc457b090654b7d6b902fc3cde47ce06354047bee7fd872
                room-ref: ./make-shipped-contracts-self-contained/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:9x6xw292fsz1b4648x9hn40y:ideation:1
                briefing: briefing:9x6xw292fsz1b4648x9hn40y:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T18:35:06.42448Z"
                decision: approve
                reason: 'Captain approved 2026-08-15 with a design correction: those references were likely never resolved by any FO, so prefer deleting the unresolvable fragment over expanding it - expand only where deletion loses operative meaning; PR review follows'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:9x6xw292fsz1b4648x9hn40y:validation
          stage: validation
          attempts:
            - id: gate-attempt:9x6xw292fsz1b4648x9hn40y-validation-1
              briefing:
                id: briefing:9x6xw292fsz1b4648x9hn40y:validation:attempt-1:revision-1
                digest: sha256:f3d474cff65cd468a3cd32463f891ad56dbfcd6b3c2580640582b9e0c98e7ca8
                request-digest: sha256:58662cbffcf0a5f82c1aa1f36325ceb53a272ae2ec9026a7610e2c37f2cc24ae
                room-ref: ./make-shipped-contracts-self-contained/review/validation/briefing-1
---

Shipped skills reference artifacts a user's machine does not have. Rewrite the seven audited instances so every shipped sentence resolves within the progressively-disclosed contract set. Base on stack layer 10 (retire-prose-grep-contract-tests); the deliverable becomes stack layer 11 - layer 10 removes the pins these rewrites would red.

The seven, with audited rewrites:
1. skills/first-officer/references/first-officer-shared-core.md:35 - "(driver binary descoped to roadmap 0222)" -> "no driver binary backs the drive yet".
2. skills/first-officer/references/fo-dispatch-core.md:213 - roadmap 0222 + runtime-support.md references -> "no driver binary backs it yet; hand-follow the deterministic skeleton above and do not probe for the unshipped command."
3. skills/first-officer/references/fo-install-gate.md:22 - internal/safehouse/state.go and insideRegistry -> "every sandbox the binary itself recognizes - the --version Sandbox: line names the same registry". CAUTION: TestVersionGateSandboxRegistry requires each env var name AND value in prose; keep the APP_SANDBOX_CONTAINER_ID = agent-safehouse row.
4. skills/first-officer/references/fo-write-core.md:11-12 - drop the docs/dev/README.md literal (covered by {workflow_dir}/README.md) and generalize docs/dev/_mods/** to {workflow_dir}/_mods/**; update internal/contractlint/fo_write_core_mutation_gate_test.go fixtures in step.
5. skills/commission/SKILL.md:678 - "the failure mode #201 addresses" -> "the failure mode this step exists to prevent".
6. skills/survey/SKILL.md:99-108 - delete tracker numbers (#318, #69, #321-#324, #317.2, 9h, za); keep the descriptive halves.
7. skills/commission/references/templates/development.md:129 - state the rule without the dated captain-ruling attribution (line 113's captain-ruling[YYYY-MM-DD] format spec stays).

Keep-verdicts (do not touch): illustrative example ids in commission SKILL.md:385-418 and development.md:64-65; references/... and spacedock:* loader-resolved cross-references; split-root .spacedock-state mentions.

## Problem

Seven sentences in shipped skill files point at artifacts that exist only in the spacedock development repo: two roadmap numbers, a Go source path, a GitHub tracker number, eleven survey tracker ids, a dated captain-ruling attribution, and this repo's own `docs/dev` workflow directory. A user who installs the plugin reads these sentences with nothing to resolve them against. The instruction still parses, but the reason behind it is a dead pointer.

One of the seven is worse than dead. The `fo-write-core` mutation-gate classifier — the machine-readable table that decides what the FO may write — carries the literal `docs/dev/README.md` alongside the generic `{workflow_dir}/README.md`, and blocks `docs/dev/_mods/**` rather than `{workflow_dir}/_mods/**`. The literals are this repo's own paths embedded in a contract every commissioned project loads.

## Proposed approach

Seven sentence rewrites in shipped prose, plus one test-fixture change that instance 4 forces.

### The seven rewrites (verified against stack tip `stack27/09-trim-version-output`, fdf008939)

1. `skills/first-officer/references/first-officer-shared-core.md:35`
   - before: ``…`«dispatch.next-action»()` skeleton (driver binary descoped to roadmap 0222).``
   - after: ``…`«dispatch.next-action»()` skeleton (the driver binary is deliberately unshipped, not missing).``

2. `skills/first-officer/references/fo-dispatch-core.md:213`
   - before: ``… — no driver binary backs it yet (descoped to roadmap 0222); the FO hand-follows the deterministic skeleton above and does not probe for the unshipped command (runtime-support.md's `→ prose` trichotomy).``
   - after: ``… — no driver binary backs it yet; hand-follow the deterministic skeleton above and do not probe for the unshipped command.``

3. `skills/first-officer/references/fo-install-gate.md:22`
   - before: ``…, extended to every row of the binary's `insideRegistry` (`internal/safehouse/state.go` is the source of truth; matching is on the VALUE, not mere presence).``
   - after: ``…, extended to every name+value pair the binary treats as a sandbox marker. Matching is on the VALUE, not mere presence: `APP_SANDBOX_CONTAINER_ID` is a generic macOS app-sandbox variable other containers also set, so presence alone would claim any of them.``

4. `skills/first-officer/references/fo-write-core.md:11-12` (classifier table)
   - before line 11: ``| allowed-process | `docs/dev/README.md`; `{workflow_dir}/README.md` | …``
   - after line 11: ``| allowed-process | `{workflow_dir}/README.md` | …``
   - before line 12: ``…; `fixtures/**`; `docs/dev/_mods/**` | …``
   - after line 12: ``…; `fixtures/**`; `{workflow_dir}/_mods/**` | …``

5. `skills/commission/SKILL.md:678`
   - before: `…is the failure mode #201 addresses — a commissioned FO that…`
   - after: `…is the failure mode this step exists to prevent — a commissioned FO that…`

6. `skills/survey/SKILL.md:99-109` — strip the leading tracker token and its em dash from all eleven `run_query` comments, keeping the descriptive half verbatim. Example: `run_query scoping            # #318 — sessions|blank_cwd|span over the cwd-prefix-scoped repo` becomes `run_query scoping            # sessions|blank_cwd|span over the cwd-prefix-scoped repo`. Covers `#318`, `#69`, `#321`, `#322`, `#323`, `#319`, `#317.2`, `#320`, `9h`, `#324`, `za`.

7. `skills/commission/references/templates/development.md:129`
   - before: `…is not prose-grep. Captain ruling (2026-07-20, verbatim): prose-greps are one-off validation evidence, never committed tests.`
   - after: `…is not prose-grep. Prose-greps are one-off validation evidence, never committed tests.`
   - Line 113's `captain-ruling[YYYY-MM-DD]` format specifier is untouched: it is a field-format placeholder, not a dated attribution.

### Two deviations from the audited rewrites, with reasons

**Instance 1.** The audited replacement was "no driver binary backs the drive yet". That duplicates the clause it follows — the sentence already opens with "no binary backs the drive". The replacement above carries the same information the roadmap number carried (the absence is deliberate, so do not go looking) without the restatement.

**Instance 3.** The audited replacement was "the `--version` `Sandbox:` line names the same registry". That does not work here: `fo-install-gate.md` line 3 states the file is "Loaded by the first officer ONLY when the version gate … lands in the **binary-absent class**". In that class there is no binary and therefore no `--version` output, so the replacement pointer is unresolvable exactly when the file is read. The wording above is self-contained instead: it states the matching rule and gives the reason value-matching is required, which is the information `internal/safehouse/state.go` actually held.

### The fixture change instance 4 forces

`internal/contractlint/fo_write_core_mutation_gate_test.go` parses the classifier table and asserts path classifications. Dropping the `docs/dev` literals reds it. The change is not a fixture edit but a fixture repair, because the naive edit is a tautology:

- `parseFOWriteClassifierTable` gains a `workflowDir` parameter and expands `{workflow_dir}` in each pattern. `TestFOWriteCoreMutationGateClassifiesTargets` runs the whole case table twice, against `docs/dev` and `ops/release`. A table that only worked for the repo it shipped from passes one and fails the other — which is precisely the property instance 4 introduces.
- A `workflowDir + "/_archive/shipped-task.md"` case is added; `{workflow_dir}/_archive/**` was in the table but never exercised.
- Every `blocked-product` expectation additionally requires an explicit pattern match via a new `foWriteClassMatches` helper. `classifyFOWriteTarget` defaults to `blocked-product`, so a bare `blocked-product` expectation is satisfied by a table that never mentions the path at all. This was not hypothetical: the first version of the fixture stayed green when the rewrite was reverted (probe recorded below).

### Audit grep (AC-1 evidence, one-off — nothing new is committed)

Run over `skills/` and `agents/`, `--include='*.md' --exclude-dir=integration` (`skills/integration/` is Go test code and testdata, never loaded by a user session):

    grep -rnE 'roadmap [0-9]{3,4}|docs/roadmap/[0-9]' --include='*.md' skills/ agents/ --exclude-dir=integration
    grep -rnE '#[0-9]{2,4}(\.[0-9]+)?\b' --include='*.md' skills/ agents/ --exclude-dir=integration
    grep -rnoE '\b(internal|cmd)/[a-z0-9_/-]+|[a-z0-9_]+\.go\b' --include='*.md' skills/ agents/ --exclude-dir=integration
    grep -rn 'docs/dev' --include='*.md' skills/ agents/ --exclude-dir=integration
    grep -rniE 'captain ruling|ruling \([0-9]{4}-' --include='*.md' skills/ agents/ --exclude-dir=integration
    grep -rnE 'runtime-support\.md|docs/(specs|site)/[a-z]' --include='*.md' skills/ agents/ --exclude-dir=integration

This grep is one-off validation evidence pasted into the report, never a committed test — the form the workflow's own prose-grep rule blesses, and the claim it settles ("no shipped file contains an unresolvable reference") is an existence fact, which is the case that rule says a grep establishes soundly.

Honest limitation: bare short ids such as `9h` and `za` (instance 6) cannot be matched by a regex with a tolerable false-positive rate. Instance 6 removes them by inspection; the grep does not police them.

### Keep-verdicts (complete recorded set)

Carried from the audit:
- Illustrative id examples: `skills/commission/SKILL.md:385-386` and `references/templates/development.md:64-65` (`#42` / `#57` as GitHub issue/PR field examples).
- `references/…` and `spacedock:*` cross-references — resolved by the skill loader.
- Split-root `.spacedock-state` mentions.

Added by this ideation — the grep surfaces them, so AC-1 cannot pass without them recorded:
- `skills/first-officer/references/fo-write-core.md:12` — the classifier's generic path-class globs (`cmd/**`, `internal/**`, `**/*_test.go`, `skills/**`, `agents/**`, `references/**`, `.github/**`, `docs/site/**`, `docs/specs/**`, `docs/roadmap/**`, `fixtures/**`). These are the contract's own path vocabulary, not references to files in this repo, and a reader in any project resolves them.
- `skills/commission/SKILL.md:286` — ``{state_branch}` = `spacedock-state/{workflow-dir-basename}`, e.g. `spacedock-state/dev` for `docs/dev``. An illustrative substitution that demonstrates the rule and resolves inside its own sentence.
- `skills/first-officer/references/first-officer-shared-core.md:10` and `references/claude-fo-dispatch.md:123` — `go build -o spacedock ./cmd/spacedock`. A build command run inside a clone of this repo, resolvable by the reader who clones it.

### Spike result

Not a paper design: all seven rewrites plus the fixture repair were applied to a detached worktree at fdf008939 and exercised.

- `go test ./internal/contractlint/ ./skills/integration/` green, plain and `-race`.
- Audit grep re-run after the rewrites returns only the keep-verdicts listed above.
- Three falsifying probes, each reverting one part of the change and confirming the suite reds:
  - re-hardcode `docs/dev/_mods/**` → `workflow_dir="ops/release": "ops/release/_mods/pr-merge.md" is blocked only by default-deny`.
  - re-hardcode `docs/dev/README.md` as the only `allowed-process` pattern → `classify "ops/release/README.md" = "blocked-product", want "allowed-process"`.
  - change `agent-safehouse` to another value in `fo-install-gate.md` → `TestVersionGateSandboxRegistry` reds, confirming instance 3 kept a live pin rather than a dead one.

`go test ./...` has one failure, `internal/cli TestCodexResolveManifestAgainstInstalledHost`, which reproduces identically on the pristine unmodified stack tip: the agent sandbox blocks reading `~/.codex/config.toml`, so the `codex` CLI cannot start (`Operation not permitted (os error 1)`). Environmental and unrelated to this change; implementation and validation should expect it on a sandboxed host and confirm it against a pristine tree rather than treating it as a regression.

### Coordination with layer 10

The entity's premise was that layer 10 (`retire-prose-grep-contract-tests`) must land first because it "removes the pins these rewrites would red". The spike does not bear that out. Run against the stack tip **without** layer 10 applied, the seven rewrites plus the fixture repair leave the whole suite green. The only test that red was `fo_write_core_mutation_gate_test.go`, which layer 10 explicitly **keeps** (it is listed there as gray, leaning KEEP), and this entity repairs it in step. The three prose-grep functions layer 10 deletes from `version_gate_smoke_test.go` (`:35`, `:77`, `:158`) pin `uname -s`, the sentinel tokens, and the launcher-invariant sentence — none of which any of the seven touches.

So the ordering is a merge-conflict convenience, not a correctness dependency: both layers edit `internal/contractlint/`, and layer 10 also touches `version_gate_smoke_test.go` in the same package. Recommend keeping the declared stack order (layer 11 on layer 10) to avoid rebase churn, but the gate should know the hard blocker does not exist — if layer 10 slips, this layer can land on layer 9 unchanged.

### Semantic changes declared

- **Authority (`fo-write-core` classifier):** the table stops naming `docs/dev` and addresses the workflow directory only through `{workflow_dir}`. No effective authority change for any project: in this repo `{workflow_dir}` *is* `docs/dev`; in any other project `README.md` was already covered by the placeholder and `_mods/**` already fell to default-deny. What changes is that the `_mods` block becomes an explicit rule rather than a fallthrough.
- **Command grammar, stored formats, runtime behavior:** unchanged. The other six rewrites are prose only.
- **Test-package internals:** `parseFOWriteClassifierTable` gains a parameter and `foWriteClassMatches` is added. Package-private test helpers, not a shipped surface.

## Out of scope

New content; anything beyond the seven plus their in-step test fixtures.

## Expected surface and tolerance

Measured from the spike diff at fdf008939, not estimated: **8 files, +70 / −44** (net +26).

    52  26  internal/contractlint/fo_write_core_mutation_gate_test.go
    11  11  skills/survey/SKILL.md
     2   2  skills/first-officer/references/fo-write-core.md
     1   1  skills/commission/SKILL.md
     1   1  skills/commission/references/templates/development.md
     1   1  skills/first-officer/references/first-officer-shared-core.md
     1   1  skills/first-officer/references/fo-dispatch-core.md
     1   1  skills/first-officer/references/fo-install-gate.md

Tolerance: ±1 file, ±20 LOC. The shipped prose is seven one-line replacements plus survey's eleven; it should not move. The whole tolerance exists for the test file, where review may want the two-workflow-dir loop shaped differently.

Note the net is **positive**, unlike the other cleanup layers in this program. This entity removes unresolvable references, not machinery; the fixture repair that instance 4 forces costs more lines than the seven prose edits save. A negative-delta target here would be met only by weakening the fixture, which is what the first spike version did.

## Acceptance criteria

**AC-1 (value) — No shipped skill or agent file references a roadmap number, tracker number, repo source path, dev-workflow path, or dated ruling that a reader outside this repo cannot resolve.**
The count of such references goes from 20 at fdf008939 (2 roadmap numbers, 1 tracker number, 11 survey tracker ids, 1 Go source path, 3 `docs/dev` literals, 1 unshipped doc reference, 1 dated ruling) to 0.
Verified by: the six-pattern audit grep in the Proposed approach, run over `skills/` and `agents/`, returning only the keep-verdicts recorded above and nothing else. Falsifying change: restoring any one of the seven sentences puts a non-keep line back in the output.

**AC-2 — The rewritten sentences carry the same instruction, and every contract check that was green before is green after.**
Verified by: `go test ./internal/contractlint/` green, with `TestVersionGateSandboxRegistry` and `TestFOWriteCoreMutationGateClassifiesTargets` specifically green; plus the before/after table in the report, which the gate reads to judge meaning preservation. Falsifying change: dropping `agent-safehouse` from the instance-3 rewrite reds `TestVersionGateSandboxRegistry` (probe recorded).

**AC-3 — The write-core fixture fails when the classifier table is re-hardcoded to a single repo's workflow directory.**
This is the criterion the naive fixture edit silently failed. Verified by: with the change applied, reverting `{workflow_dir}/_mods/**` to `docs/dev/_mods/**` reds `TestFOWriteCoreMutationGateClassifiesTargets` at the `ops/release` case, and reverting `{workflow_dir}/README.md` to `docs/dev/README.md` reds it at `ops/release/README.md`. Both probes recorded in the validation report, both reverted.

**AC-4 — The suite stays green, including the repaired write-core fixture.**
Verified by: `go test ./internal/contractlint/ ./skills/integration/` plain and `-race`, plus `go test ./...`.

## Test plan

No new committed test. Seven prose rewrites, one fixture repair in `internal/contractlint/fo_write_core_mutation_gate_test.go`, and the audit grep as one-off validation evidence pasted into the report.

The riskiest mechanism — whether the rewrites survive the pins that outlive layer 10, and whether the repaired fixture can actually fail — was spiked before this gate rather than deferred to implementation. Results are recorded under "Spike result" above: suite green plain and `-race`, three falsifying probes each confirmed to red. Implementation reproduces that diff on a branch off layer 10; validation re-runs the grep and the three probes independently.

## Stage Report: ideation

- DONE: Design confirms all seven rewrites against the stack tip with exact before/after text; the write-core fixture and sandbox-registry pin constraints verified, not assumed
  All seven applied to a detached worktree at fdf008939 and exercised: `go test ./internal/contractlint/ ./skills/integration/` green plain and `-race`. Both pin constraints verified by probe rather than reading: dropping `agent-safehouse` reds `TestVersionGateSandboxRegistry`; re-hardcoding `docs/dev/_mods/**` reds `TestFOWriteCoreMutationGateClassifiesTargets` at the `ops/release` case.
- DONE: Value AC: the audit pattern grep over shipped skills returns only recorded keep-verdicts
  Six-pattern grep defined in the body and re-run after the rewrites; output contains only keep-verdict lines. Three keep-verdicts the audit had not recorded were found and added (the classifier's generic path globs, commission SKILL.md:286's substitution example, and `go build ./cmd/spacedock` in two files) — without them AC-1 was unsatisfiable as written.
- DONE: Coordination: layer 11 on top of retire-prose-grep-contract-tests; must not land first
  Recorded, but the stated reason does not hold. Verified by running the full change against the stack tip with layer 10 NOT applied: the suite is green. Layer 10 keeps the only test this change reds. The ordering is merge-conflict convenience (both layers edit `internal/contractlint/`), not correctness.

### Summary

Two of the seven audited rewrites needed changing. Instance 1's proposed text duplicated the clause it followed; instance 3's proposed pointer to the `--version` `Sandbox:` line is unresolvable in the binary-absent class, which is the only class where `fo-install-gate.md` loads — both replaced with self-contained wording and the reasons recorded.

The instance-4 fixture update turned out to be a fixture repair. The naive version stayed green with the rewrite reverted, because `classifyFOWriteTarget` defaults to `blocked-product` and a bare `blocked-product` expectation is satisfied by a table that never mentions the path. The design expands `{workflow_dir}` over two workflow dirs and requires blocked-product rows to match an explicit pattern. That probe is why AC-3 exists as its own criterion.

The declared surface is +70/−44 across 8 files — measured from the spike diff, and positive, unlike the other layers in this cleanup program. The fixture repair costs more than the seven prose edits save, and buying a negative delta here means weakening the fixture.

## Stage Report: implementation

- DONE: Captain design correction, binding: these references were likely never resolved by any FO reading the contract - for each of the seven, DELETE the unresolvable fragment when the surrounding sentence already carries the instruction; use the spiked expansion wording ONLY where deletion loses operative meaning (the write-core table generalization and the sandbox env-var rows are the expected survivors); the +70/-44 estimate should shrink toward net-negative on the prose side
  Instances 1, 2, 5, 6, 7 are pure fragment deletions (roadmap 0222 twice, the runtime-support.md pointer, #201, eleven survey tracker ids, the dated ruling attribution); only the two predicted survivors carry spiked wording — instance 3 (sandbox env-var sentence, `agent-safehouse` pin kept) and instance 4 (`{workflow_dir}` table generalization). Prose side: 18/18 replaced lines, net −182 bytes of content.
- DONE: The fixture repair stands as designed - it caught the reverted-rewrite hole and is not prose
  `parseFOWriteClassifierTable` expands `{workflow_dir}` and the case table runs against both `docs/dev` and `ops/release`; blocked-product verdicts must match an explicit pattern (`foWriteClassMatches`), so default-deny cannot satisfy them. All three falsifying probes re-run against committed HEAD b826307c4: re-hardcoding `_mods/**` reds at "blocked only by default-deny"; re-hardcoding `README.md` reds at `classify "ops/release/README.md" = "blocked-product", want "allowed-process"`; breaking `agent-safehouse` reds `TestVersionGateSandboxRegistry`. Each probe reverted; suite green after.
- DONE: Value AC: audit-pattern grep returns only keep-verdicts; suite green plain and -race; deliverable is stack layer 11 - report the final delta honestly against the revised expectation
  Six-pattern grep over `skills/` and `agents/` returns exactly the recorded keep-verdicts (the `#42`/`#57` field examples, the classifier's `_test.go` glob, two `go build ./cmd/spacedock` lines, commission SKILL.md:286's substitution example) and nothing else; roadmap, ruling, and doc-reference patterns return empty. `go test ./...` plain and `-race`: green except `internal/cli TestCodexResolveManifestAgainstInstalledHost`, reproduced identically at pristine fdf008939 (sandbox blocks `~/.codex/config.toml`) — environmental, not a regression. Final delta: 8 files, +70/−49 (net +21) vs the audited +70/−44; test file 52/31; prose 18/18 lines, net −182 bytes. Deliverable committed as b826307c4 on `spacedock-ensign/make-shipped-contracts-self-contained` (pushed), based at fdf008939, which is also layer 10's current branch tip (layer 10 carries no commits yet; ideation already recorded the ordering as convenience-only).

### Summary

All seven rewrites applied under the deletion bias: five are pure fragment deletions, and only the two survivors the captain predicted needed expansion wording. The repaired write-core fixture proves the `{workflow_dir}` generalization can fail — all three probes red against the committed state with the designed messages. Total delta +21 (the fixture repair dominates); prose content is net-negative as the correction expected.

## Stage Report: validation

- DONE: Independently re-exercise the ACs against worktree commit b826307c4, never by reading the report: audit-pattern grep (roadmap numbers, tracker numbers, repo source paths, dev-workflow paths, dated ruling citations) over shipped skills/ and agents/ returns only the recorded keep-verdicts
  All six patterns re-run at b826307c4: roadmap, ruling, and unshipped-doc patterns return empty; the rest return exactly the recorded keeps (`#42`/`#57` field examples, `_test.go` classifier glob, two `cmd/spacedock` build lines, commission SKILL.md:286); bare `9h`/`za` confirmed gone from survey by direct grep. Before-state re-run at a pristine fdf008939 checkout shows every removed line present there and only keeps surviving.
- DONE: Verify the deletion bias held: five fragments deleted with surrounding sentences intact and carrying the same instruction; only the write-core generalization and sandbox rows carry expansion wording; prose delta net-negative measured separately from the fixture repair
  Diff fdf008939..b826307c4 read in context: instances 1, 2, 5, 6, 7 are pure fragment deletions, each host sentence intact and still carrying its instruction; new wording appears only in fo-install-gate.md (sandbox rows, `agent-safehouse` pin kept) and the fo-write-core table. Prose measured separately via `git cat-file -s`: seven prose files net −182 bytes, 18/18 lines, with fo-install-gate.md the only positive file (+99); fixture repair separate at 52/31 lines.
- DONE: The repaired fixture proves the generalization can fail: run the three probes, expect red with the designed messages; contractlint and skills/integration green plain and -race; verdict PASSED or REJECTED with per-AC citations
  Probes run in a throwaway detached checkout at b826307c4 (implementation worktree untouched), each red with the designed message, restored after: `workflow_dir="ops/release": "ops/release/_mods/pr-merge.md" is blocked only by default-deny`; `classify "ops/release/README.md" = "blocked-product", want "allowed-process"`; `TestVersionGateSandboxRegistry` red on `agent-safehouse` change. `go test ./internal/contractlint/ ./skills/integration/` green plain and `-race`. Verdict PASSED, per-AC below.

### Per-AC citations

- AC-1 PASSED — six-pattern grep at b826307c4 returns only recorded keep-verdicts; the pristine-base run shows the exact non-keep lines a restored sentence would put back. Counting nuance: the AC's before-count of 20 includes SKILL.md:286's `docs/dev`, which the recorded keep-verdict retains; references actually removed number 19.
- AC-2 PASSED — each rewritten sentence checked in diff context carries the same instruction; `TestVersionGateSandboxRegistry` and `TestFOWriteCoreMutationGateClassifiesTargets` green at b826307c4; probe 3 proves the sandbox pin still fails on divergence, and the test reads name+value rows from `internal/safehouse/state.go`, an independent source.
- AC-3 PASSED — both re-hardcode probes red at the `ops/release` cases with the designed messages; grep confirms no non-test code parses the classifier table, so the generalization has no hidden runtime reader.
- AC-4 PASSED — full `go test ./...` plain and `-race` at b826307c4 green except `internal/cli TestCodexResolveManifestAgainstInstalledHost`, which fails identically at pristine fdf008939 (`~/.codex/config.toml: Operation not permitted` under the sandbox) — environmental, not a regression.

### Findings

Material: none. Deferred risks: none (the codex failure is environmental and pre-existing, not a risk of this change). Polish, proposed disposition decline: 1) AC-1's "20 at fdf008939" arithmetic includes the SKILL.md:286 keep, so net removed is 19 — the operative verification (only keep-verdicts remain) holds; 2) the entity's Proposed approach still shows pre-correction expansion wording for instances 1, 2, 5 — the captain's deletion-bias correction and the actual deletions are recorded in the implementation report.

### Summary

Verdict: PASSED. Every AC re-exercised from scratch against b826307c4 — greps re-run both directions (base and tip), the deletion bias verified in the diff with the byte delta measured independently, and all three falsifying probes reproduced with the designed failure messages in a throwaway checkout. No material findings; two polish notes on report arithmetic and stale pre-correction wording in the entity body.
