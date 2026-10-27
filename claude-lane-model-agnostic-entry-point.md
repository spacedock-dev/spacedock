---
title: The Claude live lane selects a model by cadence, so a new model needs a workflow edit
status: validation
source: "Captain, 2026-10-07: a manual run on claude-haiku-5-5 was requested and cannot be dispatched. The captain wants opus and sonnet to be a parameter, not a mode, and the underlying entry point to be model-agnostic."
id: f77r1rj2g7sgj3hkrht7n37f
gates:
    version: 1
    records:
        - id: gate:f77r1rj2g7sgj3hkrht7n37f:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:f77r1rj2g7sgj3hkrht7n37f-backlog-1
              briefing:
                id: briefing:f77r1rj2g7sgj3hkrht7n37f:backlog:attempt-1:revision-1
                digest: sha256:431932767c9eefcead1a0540c5d3118b8f86f71c1f6fa6e6d552029da7f51eae
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:f77r1rj2g7sgj3hkrht7n37f:backlog:1
                briefing: briefing:f77r1rj2g7sgj3hkrht7n37f:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-10-07T18:14:48.37703Z"
                decision: approve
                reason: 'Captain, 2026-10-07: "file a task to reorg how we assign model to the claude-code harness. opus/sonnet should be a param, not a mode... dispatch."'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:f77r1rj2g7sgj3hkrht7n37f:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:f77r1rj2g7sgj3hkrht7n37f-ideation-1
              briefing:
                id: briefing:f77r1rj2g7sgj3hkrht7n37f:ideation:attempt-1:revision-1
                digest: sha256:199aef9b60fd507abc0b87a1989b82019785b7cb27c726c2fe2197debe3107aa
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:f77r1rj2g7sgj3hkrht7n37f:ideation:1
                briefing: briefing:f77r1rj2g7sgj3hkrht7n37f:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-10-07T23:14:48.895734Z"
                decision: approve
                reason: 'Captain, 2026-10-07: "this should be a small change. let us send this to implementation and run ci"'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:f77r1rj2g7sgj3hkrht7n37f:validation
          stage: validation
          attempts:
            - id: gate-attempt:f77r1rj2g7sgj3hkrht7n37f-validation-1
              briefing:
                id: briefing:f77r1rj2g7sgj3hkrht7n37f:validation:attempt-1:revision-1
                digest: sha256:c9e924c982b1ee9758423be9e5c94dbaf693ce71c5ed458b4894f3f25094126a
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:f77r1rj2g7sgj3hkrht7n37f:validation:1
                briefing: briefing:f77r1rj2g7sgj3hkrht7n37f:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-10-21T18:51:45.991607Z"
                decision: revise
                reason: 'Captain rejects validation (REJECTED verdict accepted) and approves the three new standing checks; route to implementation. Asks: add to internal/release/livemodels_test.go (1) an assertion that the Claude live job keeps needs: offline, (2) an assertion that opus is excluded from the routine (sonnet cadence) permitted set, (3) an assertion that the matrix binds to the selected registry output claude_selected, not claude_sonnet. Each must fail under its own audit mutation (show red/green). No broad schema/wiring suite, no new lane/job/CI step. Stay within the +150 net LOC ceiling; if it would exceed, stop for the FO rather than weaken. Report the measured diff.'
            - id: gate-attempt:f77r1rj2g7sgj3hkrht7n37f-validation-2
              briefing:
                id: briefing:f77r1rj2g7sgj3hkrht7n37f:validation:attempt-2:revision-1
                digest: sha256:fdba5117730b7aeb7cd60bf231e118101bae649679539c7ab9ef1427ec6f386e
                room-ref: '@review/validation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:f77r1rj2g7sgj3hkrht7n37f:validation:2
                briefing: briefing:f77r1rj2g7sgj3hkrht7n37f:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-10-22T05:30:42.195908Z"
                decision: revise
                reason: 'Captain rejects validation attempt-2: the task''s own ''No inherited XFAIL'' claim fails because dispatch_recovery_live_test.go:108 assigns liveXFail("claude-sonnet", 060xp69y61yhrww23g3wvwqy) unconditionally, so Haiku run 37703335512 reported break-glass selected-bare as a Sonnet XFAIL. Asks: (1) apply that binding only when the runtime target is claude-sonnet, reusing liveJourney''s target match rather than a second copy; (2) one focused check that a non-Sonnet target does not inherit the binding, shown red without the guard and green with it; (3) correct the ''No inherited XFAIL'' record and note Haiku break-glass selected-bare as unowned, routed to the Haiku follow-up. Captain raises net LOC ceiling to +170; stop and report rather than exceed.'
            - id: gate-attempt:f77r1rj2g7sgj3hkrht7n37f-validation-3
              briefing:
                id: briefing:f77r1rj2g7sgj3hkrht7n37f:validation:attempt-3:revision-1
                digest: sha256:83a26ad5433f66f7f24ae9eab5cbec17741c1787b68f38dc94cb5ff3cf38a259
                room-ref: '@review/validation/briefing-3'
              resolution:
                type: Resolution
                id: resolution:spacedock:f77r1rj2g7sgj3hkrht7n37f:validation:3
                briefing: briefing:f77r1rj2g7sgj3hkrht7n37f:validation:attempt-3:revision-1
                by: person:captain
                at: "2026-10-22T16:24:59.14622Z"
                decision: revise
                reason: 'Captain sets cycle-3 scope directly: XFAIL Haiku''s five observed failures now using the existing run evidence (run 37703335512), with a BLANK owner string (no new ids minted) so analysis can later decide whether they share one underlying cause or need separate issues. Asks: (1) liveXFail("claude-haiku", "") on keep-moving-posture, default-headless-gate-stop, withdrawn-gate-recovery, filing, and break-glass selected-bare; Sonnet bindings unchanged and claude-haiku never inherits a Sonnet owner; (2) update the focused break-glass check accordingly; (3) add claude-haiku rows (owner blank, pending analysis) to docs/runtime-live-ci-registry.md and sync its missing claude-sonnet jfdrr2dbdtt952s4cvvy74k8 rows; (4) confirm blank-owner bindings compile, select, grade XFAIL, and do not fail any existing check. Net LOC ceiling raised to +200; stop and report rather than exceed.'
            - id: gate-attempt:f77r1rj2g7sgj3hkrht7n37f-validation-4
              briefing:
                id: briefing:f77r1rj2g7sgj3hkrht7n37f:validation:attempt-4:revision-1
                digest: sha256:0d92a025f8be8ef9db4728803843fdfc412db4e3a5f644e8735a16090554b152
                room-ref: '@review/validation/briefing-4'
              resolution:
                type: Resolution
                id: resolution:spacedock:f77r1rj2g7sgj3hkrht7n37f:validation:4
                briefing: briefing:f77r1rj2g7sgj3hkrht7n37f:validation:attempt-4:revision-1
                by: person:captain
                at: "2026-10-24T15:21:09.004555Z"
                decision: approve
                reason: 'Captain approves candidate 861a6e13a at validation attempt-4: Haiku blank-owner XFAILs, per-target break-glass binding, gap-aware withdrawn-gate-recovery, registry doc sync; accepts the listed deferred risks and presentation without an independent cycle-3 reviewer.'
              application:
                target-stage: done
                state: pending
started: 2026-10-07T18:15:05Z
worktree: .worktrees/spacedock-ensign-claude-lane-model-agnostic-entry-point
review-round:
    id: round:f77r1rj2g7sgj3hkrht7n37f:implementation:3
    stage: implementation
    cycle: 3
    briefing:
        id: briefing:f77r1rj2g7sgj3hkrht7n37f:implementation:round-3
        digest: sha256:1827cf2687131ddf893656e243a138c194fef18f67f3cc9800dc69fbc2844030
        room-ref: '@review/implementation/round-3'
mod-block: merge:pr-merge
pr: pr-merge:853
---

## Problem

The Claude live lane has no model parameter. `claude_version` selects the CLI version, not the model.
The captain requested `claude-haiku-5-5`, but the dispatch surface cannot select it.
Opus and Sonnet must be parameter values, not execution modes. Host-specific secrets do not justify model-specific entry points.

## Baseline and exact evidence

All product quotes below refer to `origin/main` at `7558f2e099b661927b75fdde22e00950a53c3532` (baseline B).
The working checkout is stale at `c450d5d526fca30bbcd2c06cf2b40358e5a6ea8f`.
Its missing registry files and older matrix are not the design baseline. The supervisor confirmed B; no checkout update was performed.
Read B with `git show 7558f2e099b661927b75fdde22e00950a53c3532:<path>`.

B's `internal/release/live_models.txt` contains exactly:

```text
claude.sonnet=claude-sonnet-5-5
claude.opus=claude-opus-5-5
codex.exec=gpt-6-luna
pi.oauth=openai-codex/gpt-6-luna:max
pi.api-key=openai/gpt-6-luna:max
```

B's `.github/workflows/runtime-live-e2e.yml` assigns the Claude matrix with these exact lines:

```yaml
          - cadence: ${{ github.event_name == 'pull_request' && 'pull-request' || inputs.live_cadence }}
            model: ${{ (github.event_name == 'pull_request' || inputs.live_cadence == 'sonnet') && needs.offline.outputs.claude_sonnet || needs.offline.outputs.claude_opus }}
            effort: max
            environment: ${{ (github.event_name == 'pull_request' || inputs.live_cadence == 'sonnet') && 'CI-E2E' || 'CI-E2E-OPUS' }}
```

Thus `claude-sonnet-5-5, max, CI-E2E` is the evaluated row, not a literal YAML line.
The offline step reads the registry, changes dots and hyphens in keys to underscores, and writes `$GITHUB_OUTPUT`.
Explicit job outputs expose the two Claude IDs. The matrix selects one output by cadence.
`SPACEDOCK_LIVE_MODEL: ${{ matrix.model }}` passes the selected ID to the existing live test invocation.
The same matrix value labels the summary, artifact, metrics directory, and Claude config directory.
Go embeds the same file through `livemodels.go`; its current parser accepts every `key=value` record.

`internal/claudeteam/` at B handles model identity and context estimates, not this dispatch selection.
`TestContextLimitForModelBoundary` deliberately treats Haiku as the default context size. This task does not change that policy.

## Proposed approach

Add optional string input `claude_model`. Its values are registry suffixes, such as `sonnet`, `opus`, and `haiku`, not raw IDs.
The common entry point resolves `claude.<value>` without a model-family switch or a model choice list.
Blank input retains the existing cadence defaults. `claude_version` remains independent and unchanged.

The registry remains the only production location containing model IDs.
Add the Haiku ID there once. Add two policy records there, without repeating IDs:

```text
claude.haiku=claude-haiku-5-5
claude.allowed.sonnet=sonnet,haiku
claude.allowed.opus-pre-release=sonnet,opus,haiku
```

The file now contains two record types: model ID records and cadence-policy records.
One file keeps model registration and approval eligibility together. A new model needs registry edits, not workflow edits.
Every policy member must resolve to a `claude.<key>` record. Duplicate records, empty members, unknown members, and unsafe IDs fail closed.
`LiveModels` continues to expose model IDs only; the Go parser excludes `claude.allowed.` records.

The environment follows cadence only. The model parameter never changes the environment, required reviewers, secrets, job selection, or effort.
An explicit model requires permission in the selected cadence's set.
`sonnet` plus `opus` fails with `use live_cadence=opus-pre-release`; it neither spends routine approval nor silently upgrades approval.
This is a pre-approval gate in the existing secret-free offline job, not a GitHub-native dispatch rejection.
All live jobs retain `needs: offline`. No new job or CI step is needed.

| Event / cadence | Blank model | Explicit model | Claude environment | Other live jobs |
|---|---|---|---|---|
| Pull request | `sonnet` | Dispatch inputs ignored | `CI-E2E` | Codex and Pi, unchanged |
| Manual `sonnet` | `sonnet` | Any member of `claude.allowed.sonnet` | `CI-E2E` | Codex, unchanged |
| Manual `opus-pre-release` | `opus` | Any member of `claude.allowed.opus-pre-release` | `CI-E2E-OPUS` | None, unchanged |
| Manual `pi` | No Claude model | Nonempty `claude_model` fails | No Claude job | Pi, unchanged |

The existing `labeled` guards stay unchanged. A Pi dispatch cannot silently discard a requested Claude model.
Unknown keys and raw IDs fail rather than falling back to Sonnet. PR inputs cannot override the PR default.
Claude IDs must match `[A-Za-z0-9][A-Za-z0-9._-]*` because existing consumers use them in paths and shell text.
This is an artifact-safety constraint, not a family allowlist. New model families use the same resolver.
Existing Codex and Pi IDs, outputs, credentials, and jobs remain unchanged.

### Mechanisms and alternatives

- Registry-key parameter serves AC-1. Raw IDs are simpler inputs but violate the single-source rule; a choice list requires workflow edits.
- Registry policy serves AC-2 and AC-3. A workflow condition per model recreates the original maintenance problem.
- Generic resolution in the existing offline step serves AC-1 and AC-3. A new CLI command, script, job, or lane adds unnecessary surface.
- Safe-ID and policy-reference checks serve AC-1 and AC-3. Unchecked strings can break artifacts or silently bypass intended selection.
- Go policy filtering serves AC-2. Returning comma-separated policies as model IDs changes the existing map's meaning.
- The proposed test file serves AC-2 and AC-3. Static model-name searches cannot establish schema validity or approval wiring.
  This test file is an explicit captain gate decision. The product design does not depend on keeping that particular test file.

## Exact proposed changes

These are design text, not edits to product files. All before text is from B.

### `.github/workflows/runtime-live-e2e.yml`

Before: no `claude_model` input. After: insert after `live_cadence.options`:

```yaml
      claude_model:
        description: "Claude model key from live_models.txt, such as sonnet, opus, or haiku. Empty uses the cadence default."
        required: false
        type: string
        default: ""
```

Before: `description: "Select routine Sonnet, Opus pre-release, or optional Pi evidence."`
After: `description: "Select the approval cadence and runtime jobs. claude_model can override the default model within the cadence policy."`
The existing input name, default, options, and required flag remain unchanged.

Before offline outputs: the five existing named model outputs. After: retain all five and add:

```yaml
      claude_selected: ${{ steps.live_models.outputs.claude_selected }}
```

Before model matrix expression: the exact `model:` line quoted above. After:

```yaml
            model: ${{ needs.offline.outputs.claude_selected }}
```

Before the Claude job comment:

```text
  # Pull requests normalize to one Sonnet 5.5 leg. Manual dispatches choose the
  # routine Sonnet cadence or the separately approved Opus pre-release cadence.
```

After:

```text
  # Pull requests use the Sonnet default. Manual dispatches can select a model key.
  # The cadence controls approval; the registry controls permitted model keys.
```

Before `Resolve live lane models`: no step environment, and this exact `run` body:

```bash
set -euo pipefail
# The Go harness embeds internal/release/live_models.txt; read the same
# file here rather than rebuilding a command or repeating a literal.
while IFS='=' read -r key id; do
  if [ -n "$key" ]; then
    echo "${key//[.-]/_}=$id"
  fi
done < internal/release/live_models.txt >> "$GITHUB_OUTPUT"
```

After: keep its name and `id: live_models`, add this environment, and replace its `run` block with the shell below.

```yaml
        env:
          LIVE_EVENT: ${{ github.event_name }}
          LIVE_CADENCE: ${{ inputs.live_cadence }}
          CLAUDE_MODEL_KEY: ${{ inputs.claude_model }}
```

The shell is the proposed production resolver and the throwaway trial subject:

```bash
set -euo pipefail
awk -F= '
function fail(message) {
  print "::error::" message > "/dev/stderr"
  exit 1
}
{
  if (NF != 2 || $1 == "" || $2 == "" || ($1 in records))
    fail("Invalid or duplicate live model registry record at line " NR)
  records[$1] = $2
}
END {
  # An exit during record parsing also enters END. Emit nothing on that path.
  if (NR != length(records)) exit 1
  for (key in records) {
    if (key ~ /^claude\.allowed\./) {
      cadence = substr(key, 16)
      if (cadence != "sonnet" && cadence != "opus-pre-release")
        fail("Unknown Claude policy cadence: " cadence)
      count = split(records[key], members, ",")
      for (i = 1; i <= count; i++) {
        member = members[i]
        if (member !~ /^[a-z][a-z0-9-]*$/ || !(("claude." member) in records))
          fail("Unknown Claude policy member: " member)
        if ((cadence SUBSEP member) in allowed)
          fail("Duplicate Claude policy member: " member)
        allowed[cadence, member] = 1
      }
    } else if (key ~ /^claude\./) {
      if (substr(key, 8) !~ /^[a-z][a-z0-9-]*$/ || records[key] !~ /^[A-Za-z0-9][A-Za-z0-9._-]*$/)
        fail("Unsafe Claude registry key or model ID: " key)
    }
  }
  if (!("sonnet" SUBSEP "sonnet" in allowed) || !("opus-pre-release" SUBSEP "opus" in allowed))
    fail("Missing Claude cadence default or policy")
  event = ENVIRON["LIVE_EVENT"]
  cadence = ENVIRON["LIVE_CADENCE"]
  model = ENVIRON["CLAUDE_MODEL_KEY"]
  if (event == "pull_request") {
    cadence = "sonnet"
    model = "sonnet"
  } else if (event != "workflow_dispatch") {
    fail("Unsupported live event: " event)
  }
  if (cadence == "pi") {
    if (model != "") fail("claude_model requires live_cadence=sonnet or live_cadence=opus-pre-release")
  } else {
    if (cadence != "sonnet" && cadence != "opus-pre-release") fail("Unknown live cadence: " cadence)
    if (model == "") model = (cadence == "sonnet" ? "sonnet" : "opus")
    if (model !~ /^[a-z][a-z0-9-]*$/ || !(("claude." model) in records))
      fail("Unknown claude_model key; register it in internal/release/live_models.txt")
    if (!((cadence SUBSEP model) in allowed)) {
      if (("opus-pre-release" SUBSEP model) in allowed)
        fail("claude_model is not permitted; use live_cadence=opus-pre-release")
      fail("claude_model is not permitted by the selected cadence policy")
    }
    selected = records["claude." model]
  }
  # Keep existing host outputs. Policy records never become model outputs.
  for (key in records) {
    if (key ~ /^claude\.allowed\./) continue
    output = key
    gsub(/[.-]/, "_", output)
    print output "=" records[key]
  }
  if (selected != "") print "claude_selected=" selected
}' internal/release/live_models.txt >> "$GITHUB_OUTPUT"
```

Model IDs never come from an Actions expression inside the shell. Input strings enter through environment variables.
The resolver validates all policy references before emitting outputs, including records not selected by the current dispatch.
The existing offline step stays after Build; moving it earlier is not required for the pre-approval guarantee.
No other workflow changes are proposed. The cadence/environment expression, effort shim, live invocation, and artifact consumers stay intact.

### `internal/release/live_models.txt`

Before: the five exact records quoted above. After: those five records unchanged, followed by the three proposed records quoted above.
Register future models by adding a model record and adding its key to the permitted cadence sets. No workflow edit is required.
The literal IDs in this entity are evidence and proposed data, not additional runtime configuration sources.

### `internal/release/livemodels.go`

Before opening comments:

```go
// ABOUTME: The single source of truth for the live E2E lane model ids. The data
// ABOUTME: file live_models.txt holds one `key=id` line per lane; this package
// ABOUTME: embeds it for Go, and the live workflow reads the same file directly,
// ABOUTME: so no consumer repeats a model literal.
```

After:

```go
// ABOUTME: live_models.txt holds live model IDs and Claude cadence policies.
// ABOUTME: Go exposes the model records; the workflow also reads the policies.
// ABOUTME: No consumer repeats a model ID.
```

Before `LiveModels` comment:

```go
// LiveModels maps each dotted lane key to its pinned model id, parsed from
// live_models.txt: claude.sonnet, claude.opus, codex.exec, pi.oauth, and
// pi.api-key. A model change is one edit to that file.
```

After:

```go
// LiveModels maps dotted model keys to IDs from live_models.txt.
// Cadence-policy records are excluded. A model change needs one registry edit.
```

Before:

```go
		models[strings.TrimSpace(key)] = strings.TrimSpace(id)
```

After:

```go
		key = strings.TrimSpace(key)
		if strings.HasPrefix(key, "claude.allowed.") {
			continue
		}
		models[key] = strings.TrimSpace(id)
```

The existing embed and map API remain unchanged. CI owns registry validation; this task does not redesign the general Go parser.

### `internal/release/livemodels_test.go` (one new test file)

Before: no file at B. After: one test file in the existing `release` package, with these exact test names and definition comments:

```go
// TestParseLiveModelsExcludesCadencePolicy keeps policy lists out of model consumers.
// TestLiveModelRegistrySchema rejects duplicate records and unresolved policy members.
// TestClaudeModelWorkflowApprovalWiring binds model selection to the offline gate without changing approval routing.
```

The first test feeds synthetic model and policy records to `parseLiveModels` and compares the complete map to independent literal expectations.
The schema test checks the embedded raw registry, duplicate keys, safe Claude IDs, both default memberships, and every permitted member's resolution.
It also requires `haiku` in both policy sets and excludes `opus` from the routine set, without repeating production model IDs.
The structural test inspects the actual YAML job blocks, not documentation, and checks these independent invariants:
optional string input; empty default; resolver environment bindings; selected output wiring; unchanged cadence/environment expression; all live `needs: offline` dependencies.
It also checks the existing job count, unchanged `if` expressions, and unchanged `SPACEDOCK_LIVE_MODEL` binding.
This is a bounded schema/wiring check, not a shell replay harness or proof of GitHub scheduling.
Implementation must present exact assertions for review; these comments specify the test contracts, not placeholder implementations.

### `docs/releasing.md`

Before: no model-parameter instructions; step 4 includes:

```bash
gh workflow run "Runtime Live E2E" --ref main
```

After: keep that release command unchanged. Immediately after its command block, insert exactly:

````markdown
   **Select a Claude model.** Leave `claude_model` empty for the cadence default.
   Use a registry key to select another model:

   ```bash
   gh workflow run "Runtime Live E2E" --ref main -f live_cadence=sonnet -f claude_model=haiku
   gh workflow run "Runtime Live E2E" --ref main -f live_cadence=opus-pre-release -f claude_model=opus
   ```

   Model IDs and permitted cadence sets are in `internal/release/live_models.txt`.
   Add a model record and its permitted memberships there. Do not edit the workflow to register a model.
   The cadence still selects runtime jobs and the approval environment. A model override does not change either selection.
   The offline job rejects an unknown key or a forbidden combination before live approval.
   Use `live_cadence=opus-pre-release` for `claude_model=opus`. A Pi dispatch cannot accept a Claude model key.
   `claude_version` selects the Claude Code version, not its model. Maximum effort remains unchanged.
   An exploratory model run does not replace the required release cadence evidence.
````

## Expected surface and tolerance

Estimate net LOC change: +180, across 5 files. Estimated insertions: 197; deletions: 17.
Net tolerance: +120 through +240 (±60). File tolerance: exactly the five listed files; no additional product files without a design reset.
The five files are the workflow, registry, Go model reader, proposed Go test file, and release documentation listed above.
Both the workflow and registry change. The one new file is test-only; no lane, job, or CI step is added.
The entity body and stage report are state artifacts, outside the product LOC estimate.

Observable semantic changes:

1. CI dispatch gains optional `claude_model` with registry-key grammar and documented precedence. Existing inputs retain their meanings.
2. The registry gains model records and `claude.allowed.<cadence>=<comma-separated keys>` policy records. Go consumers see only model records.
3. The offline resolver rejects malformed policy, unsafe Claude IDs, unknown model keys, Pi/model combinations, and insufficient cadence permission.
4. Permitted overrides change the actual Claude model and existing model-labelled artifacts, summaries, config paths, and metrics paths.
5. Default model choice, cadence job selection, maximum effort, environment authority, approvals, credentials, release gates, and Codex/Pi behavior do not change.

## Acceptance criteria

**AC-1 (VALUE) - A registered named Claude model runs without a workflow edit.**
Verified by: a manual dispatch with `live_cadence=sonnet` and `claude_model=haiku` records the registered Haiku ID in host output and completes the existing Claude lane.
Record the run ID, exact SHA, CLI version, model observation, environment, conclusion, and artifact.
The baseline B cannot select that key through dispatch: zero selectable Haiku runs through this input; the candidate permits at least one successful run.
After parameterization, adding or changing a permitted model requires zero changed workflow lines, versus B's missing model input and explicit two-output selection.
Prove extension with a throwaway registry-only change using another synthetic key; keep the candidate workflow bytes identical and compare the resolved ID independently.
Falsifiers: remove parameter precedence, hardcode the two known keys, or change only the artifact label while the actual host runs Sonnet.
A resolver pass alone does not satisfy the live-run portion. Provider errors are unverified evidence, not a model success.

**AC-2 - Existing cadence defaults and approval authority are preserved.**
Verified by: the unchanged-event/default matrix in the approach section, schema/wiring tests, and GitHub job/environment observations on the candidate.
PR and blank routine dispatch select Sonnet; blank pre-release selects Opus; blank Pi starts no Claude job.
All four existing environments retain their reviewers and host secrets. The explicit model never selects an environment.
Falsifiers: map routine to the Opus environment, remove a live dependency, admit Opus to the routine policy, or alter a blank default.
The Go parser excludes policy records while preserving all existing model keys, as checked against an independent synthetic map.

**AC-3 - A forbidden model/cadence combination fails before live approval with an observable remedy.**
Verified by: dispatch `live_cadence=sonnet`, `claude_model=opus`; offline fails nonzero with `use live_cadence=opus-pre-release`.
GitHub records no started live job or pending live environment approval for that run. All live jobs remain dependent on offline success.
The resolver trial also rejects missing keys, unsafe values, duplicate records, and unresolved permitted members without emitting model outputs.
Falsifiers: silently select Sonnet, ignore policy, silently upgrade the environment, or return success without a selected model.

**AC-4 - Operator instructions identify the new input without duplicating model IDs or weakening release evidence.**
Verified by: execute the documented Haiku dispatch during AC-1, and inspect its submitted fields against the documented command.
The unchanged default release command remains usable; the docs distinguish model keys from CLI version pins and cadence permissions.
Falsifiers: document a raw ID as the key, use `claude_version` to select a model, or claim exploratory evidence replaces required cadences.

## Test plan

Primary owners at B: the offline registry-output step and the existing live invocation own model delivery.
No dispatch-input or policy validator exists at B. `internal/claudeteam` tests own context accounting, not model dispatch, and stay unchanged.
The proposed test file covers the new schema and wiring failure modes in the existing package, not a new framework.

1. Before implementation, retain the one-off resolver trial as first-test evidence. Exercise defaults, overrides, unknown keys, policy typos, unsafe IDs, and output atomicity.
   Cost: seconds, deterministic. Its independent expected IDs and exit codes catch ignored overrides and permissive policy checks.
2. Add the proposed Go tests before product edits. Parser tests catch leaked policy values; schema tests catch unresolved permissions.
   Structural checks catch changed dependencies, approval expressions, and input/output wiring. Cost: seconds; no workflow-shell replay harness.
3. Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` in the future implementation worktree.
   This ideation round does not format the stale product checkout or claim its tests validate B.
4. Run the documented Haiku dispatch and the forbidden Opus/routine dispatch on the exact candidate SHA.
   Use existing environments and approvals. Cost: the current Claude suite can take 90 minutes; a cheap model-resolution check comes first.
   Observe the actual model in host artifacts, not only the matrix summary. Keep the default routine and pre-release behavior evidence separate.
5. Compare workflow bytes before and after the synthetic registry-only extension; compare its resolved output to the trial's independent expected ID.
   Cost: seconds. A two-key allowlist fails even when the Haiku-specific example passes.
6. Run a detached adversarial audit before merge. CI and release machinery trigger the full audit under `docs/dev/README.md#proof-policy`.
   Attempt to bypass `needs: offline`, admit Opus to routine policy, or hardcode the selected model; confirm each relevant check fails.
   The audit belongs to validation, on a throwaway checkout, not this no-worktree ideation round.

## Risk evidence

The trial record follows below. No `no spike needed` claim is made: registry policy and parameter precedence are new behavior.

### Throwaway resolver trial

Executed the exact proposed shell above with `python3` driving `bash -c` and `/usr/bin/awk`.
The temporary directory used prefix `spike-claude-lane-model-agnostic-entry-point-` and was removed after the run.
The registry came from `git show B:internal/release/live_models.txt`, plus the three proposed records.
The shell was extracted from this entity, not reimplemented by the trial driver.
Its SHA-256 was `f49a56e827f88a1bd60d2622a25f9c578c94fddac36198187365b298b93cb408` (no trailing newline).
Every case used an empty output file; the driver asserted exit status, selected output, and absence of outputs on rejection.
It also checked unchanged Codex/Pi output values, absence of policy outputs, and absence of an injected file.

| Trial case | Observed result |
|---|---|
| PR with supplied Opus override | Sonnet ID; exit 0 |
| Blank routine / blank pre-release | Sonnet / Opus IDs; exit 0 |
| Blank Pi | No selected Claude output; exit 0 |
| Explicit routine Sonnet | Sonnet ID; exit 0 |
| Haiku on routine / pre-release | Haiku ID on both; exit 0 |
| Sonnet / Opus on pre-release | Requested IDs; exit 0 |
| Opus on routine | Exit 1; `use live_cadence=opus-pre-release`; zero outputs |
| Claude key on Pi | Exit 1; cadence remedy; zero outputs |
| Unknown key / raw ID / `$(touch injected)` key | Exit 1; unknown-key error; zero outputs |
| Unknown permitted member / empty member | Exit 1; policy-member error; zero outputs |
| Duplicate registry record / duplicate permitted member | Exit 1; duplicate error; zero outputs |
| Haiku ID changed to `../../injected` | Exit 1; unsafe-ID error; zero outputs |
| Routine policy without its Sonnet default | Exit 1; missing-default error; zero outputs |
| Unknown cadence | Exit 1; unknown-cadence error; zero outputs |
| Registry-only `claude.trial=fixture-future-model`, permitted on routine | `fixture-future-model`; exit 0; resolver bytes unchanged |

Result: **22 cases passed**. This proves generic key selection and fail-closed policy handling in the proposed shell, not hosted CI execution.
To reproduce, extract the shell block beginning `set -euo pipefail` followed by `awk -F=` into a temporary `resolver.sh`.
Place the proposed registry at `internal/release/live_models.txt` relative to that temporary directory. Then run, for example:

```bash
: > outputs
LIVE_EVENT=workflow_dispatch LIVE_CADENCE=sonnet CLAUDE_MODEL_KEY=haiku GITHUB_OUTPUT="$PWD/outputs" bash resolver.sh
cat outputs
: > outputs
LIVE_EVENT=workflow_dispatch LIVE_CADENCE=sonnet CLAUDE_MODEL_KEY=opus GITHUB_OUTPUT="$PWD/outputs" bash resolver.sh
# The second command exits 1, prints the cadence remedy, and leaves outputs empty.
```

Vary the event, cadence, key, and registry as specified in the table. Expected model IDs are the exact baseline/proposed values quoted above.
No fixture repository, product-file edit, workflow dispatch, paid host call, or worktree was used.
A shell trial cannot prove GitHub scheduling, environment approval state, actual model availability, or provider support for maximum effort.
Those observations remain required live evidence. The bounded read did not inspect the live harness outside the permitted paths.
The existing `SPACEDOCK_LIVE_MODEL` binding is visible, but its complete path to host argv remains unverified in this round.
The unchanged effort shim requests `max` for Haiku too. If the host rejects that combination, stop for a captain decision; do not silently lower effort.

An approval is attached to a cadence, so the caller controls which approval they request.
This is acceptable only with explicit refusal when the cadence does not cover the requested model.
Registry policy edits can expand permission. They remain reviewed CI-policy changes, not an alternate route around environment reviewers.

## Out of scope

No new lane, job, environment, approval model, release gate, or host credential mechanism.
No change to effort, context-budget inference, Claude transport internals, Codex selection, or Pi selection.
No raw-ID override, arbitrary environment input, dynamic approval promotion, repository-wide model cleanup, worktree, or PR in this round.
No product edits in this round. Only this entity body and its stage report are committed to `spacedock-state/dev`.

## Stage Report: ideation

- DONE: Quote the exact lines that assign the Claude model today, from the workflow matrix and from internal/release/live_models.txt.
  Baseline B is origin/main at 7558f2e099b661927b75fdde22e00950a53c3532; the body quotes its registry and matrix, not the stale checkout.
- DONE: Decide how a model parameter reaches the lane without breaking the single-source rule, and give exact before/after wording for each file you would change.
  The body specifies registry-key input, generic offline resolution, registry policy, Go filtering, proposed test contracts, and exact operator documentation.
- DONE: Keep the cadence mapping as a convenience over the parameter, so existing cadences still work.
  The event/cadence table preserves defaults, environments, host jobs, and effort; the 22-case trial catches ignored overrides and changed defaults.
- DONE: Declare the expected surface, net LOC with tolerance, and every semantic change, treating CI dispatch inputs as a user-visible surface.
  Five named files; +180 net LOC (+197/-17), ±60 net; input grammar, registry format, rejection behavior, and preserved authority are explicit.
- DONE: Exercise the riskiest unverified mechanism in a throwaway trial, or record an auditable "no spike needed" naming the proven mechanisms it relies on.
  Exact proposed resolver SHA-256 f49a56e827f88a1bd60d2622a25f9c578c94fddac36198187365b298b93cb408 passed 22 cases; hardcoded keys or permissive policy fail them.
- DONE: Write the acceptance with at least one value criterion that measures running a named model without a workflow edit, against a baseline that can move the wrong way.
  AC-1 requires actual Haiku host evidence and registry-only extension; AC-3 requires observable Opus/routine refusal before live approval.
- DONE: State whether the detached adversarial audit applies, since the CI and release machinery is a declared high-stakes surface.
  Full detached audit applies before merge; the test plan names approval-bypass, policy-expansion, and hardcoded-model adversarial edits.
- DONE: Do not add a lane, and do not change the environment approval model.
  No product files changed; the proposal reuses the offline step and preserves every existing live dependency and cadence/environment expression.
- DONE: Report in the canonical item form, one DONE/SKIPPED/FAILED line per checklist item with an evidence or rationale line, ending with a non-empty Summary.
  This report uses canonical items; only the entity body changes, and its original frontmatter remains intact.
- SKIPPED: Live dispatch, Haiku host execution, full Go suites, product formatting, and detached audit in this ideation round.
  No product implementation or worktree exists; shell proof does not verify hosted approvals, actual model delivery, or Haiku support for maximum effort.

### Summary

Designed a model-agnostic Claude parameter with one registry for IDs and permitted cadence sets, while retaining cadence-controlled approval.
The throwaway resolver passed 22 cases, including registry-only extension and explicit pre-approval refusal; live host and GitHub evidence remain unverified.
The five-file implementation proposal, acceptance criteria, documentation changes, and audit requirement are ready for independent ideation review.


## Implementation contract and evidence (supersedes the earlier surface estimate)

The implementation checkout was rebased to baseline B before work. Code commit
`d3df726fe4123ba2d9c053fb361b828ad4850d51` is on
`spacedock-ensign/claude-lane-model-agnostic-entry-point`, not main.
The captain's small-change direction and FO authorization re-declare the surface:
**at most +150 net LOC across the five named product files**, superseding the
ideation estimate of +180 with +120..+240 tolerance. Actual diff against B is
**+105 net (+124/-19), five files**:

| File | Added/deleted | Net | User-visible reason |
|---|---:|---:|---|
| `.github/workflows/runtime-live-e2e.yml` | 76/11 | +65 | Optional input, pre-approval refusal, resolved model delivery |
| `internal/release/live_models.txt` | 3/0 | +3 | Haiku registration and cadence eligibility |
| `internal/release/livemodels.go` | 10/8 | +2 | Keep policy lists out of existing model-ID consumers |
| `internal/release/livemodels_test.go` | 18/0 | +18 | Regression proof for the changed Go reader, not schema duplication |
| `docs/releasing.md` | 17/0 | +17 | Operator dispatch syntax and preserved release obligations |

Both workflow and registry change. No lane, job, environment, or CI step was added.
**Architecture: model-agnostic resolution with cadence-dependent eligibility.**
Cadence gates both model eligibility and the approval environment, and selects jobs;
it is NOT merely a convenience default over the model parameter. The older ideation
checklist wording was incomplete. Model overrides cannot promote approval or bypass it.

The independent review findings are addressed as authorized by the FO:
- Keep `claude.allowed.*`: it earns its place by allowing future registered keys to
  declare eligibility without workflow edits. Removing it needs a hardcoded model
  exception or another policy format. These are reviewed approval-policy records.
- Cut the proposed duplicate schema validator and broad new wiring test. The only
  new test is `TestParseLiveModelsExcludesCadencePolicy`: synthetic model records
  must survive intact while policy records do not. It failed before the Go change
  and passes after it; removing the filter fails it. No second schema implementation.

### Actual resolver and approval evidence

Executed the candidate's actual offline shell, extracted from `Resolve live lane models`,
not a reimplementation. SHA-256 (dedented, trailing newline included):
`29784c5835a2803e0ee7e8928585f7e41f35fb8df2047b9548a51ea498cb62b3`.
**25 cases passed**: PR precedence; routine/pre-release/Pi defaults; all five allowed
cadence/key pairs; forbidden Opus/routine with the exact cadence remedy; Pi override;
unknown/raw/injected keys; unknown/empty/duplicate policy members; duplicate/malformed
records; unsafe ID; missing default; unknown policy cadence, event, and cadence;
registry-only `claude.future=fixture-model` extension permitted on routine.
Success cases checked selected IDs and unchanged Codex/Pi outputs; every rejection
checked nonzero exit, diagnostic, and zero outputs. Injection created no file.
The extension resolved `fixture-model` with byte-identical workflow. Ignoring the
parameter or allowing routine Opus makes these cases fail. These cases are NOT the
host spike and do not establish GitHub approval semantics or model availability.

Reproduce by extracting the existing step's `run: |` body, preserving shell quoting,
and executing it with `LIVE_EVENT`, `LIVE_CADENCE`, `CLAUDE_MODEL_KEY`, and
`GITHUB_OUTPUT` set against a scratch copy of `internal/release/live_models.txt`.
The case list above names the independent mutations and expected outcomes. Local
throwaway driver/log: `.claude/model-probe/resolver-check.py` and `resolver-check.log`
(ignored worktree evidence, not required installation dependencies).
A parsed-YAML comparison against B confirmed the same five jobs, every step count,
job guards/dependencies/environments, effort, and unchanged Codex/Pi/job definitions.
`matrix.model` alone changes to the selected offline output. The existing host path
is `SPACEDOCK_LIVE_MODEL` -> `newClaudeLiveRunner.modelName` -> front-door `--model`;
the unchanged workflow shim supplies `--effort max`.

### Host spike: Haiku at requested maximum effort, older versus newer CLI

Prior FO evidence, recorded explicitly: Claude **2.1.283**, exit 0,
`canonicalModel=claude-haiku-5-5`, 77 thinking tokens / 81 output tokens, stderr
`[claude-code:unrecognized_model] {"model":"claude-haiku-5-5","query_source":"sdk"}`.
That old warning is a finding, not accepted as noise. Independently reproduced below.
The first unisolated re-run also exited 0 with that canonical model and warning,
but loaded local instructions and returned a permission-related response rather
than `ok`; it is not counted as the clean probe. Both clean comparison runs used:

```bash
"$CLI" -p --model claude-haiku-5-5 --effort max --output-format json \
  --setting-sources '' --tools '' --disable-slash-commands \
  --system-prompt 'Respond to the user directly. No tools are available.' \
  'Reply with the single word ok and nothing else.'
```

CWD: worktree-local `.claude/model-probe`; existing local login, no CI API-key proof.
New CLI installed only under that directory using
`npm install --prefix "$PWD/.claude/model-probe/newer" --no-audit --no-fund @anthropic-ai/claude-code@2.1.293`.
Local npm disables install scripts, so ran that package's `install.cjs` explicitly;
`newer/node_modules/.bin/claude --version` then reported 2.1.293. Operator CLI stayed 2.1.283.

| Observation | Existing CLI 2.1.283 | Isolated newer CLI 2.1.293 |
|---|---|---|
| Exit / result / is_error | 0 / `ok` / false | 0 / `ok` / false |
| modelUsage key and canonicalModel | `claude-haiku-5-5` | `claude-haiku-5-5` |
| Thinking / output tokens | 264 / 268 | 687 / 691 |
| stderr | Same `unrecognized_model` warning | Empty (0 bytes) |
| costBasis / contextWindow | unknown / 200000 | list / 1000000 |
| Session ID | `598306e2-b165-44cc-a59e-47fea0db0bc2` | `28621284-0fd1-47ec-af50-e4e3f6ca179c` |

Both accepted `--effort max` and thinking was active; the older-CLI hypothesis is
supported: upgrading to 2.1.293 removed the warning in this comparison. No effort
fallback or product CLI pin was added. Context-size policy remains out of scope.
Raw local outputs: `.claude/model-probe/{2.1.283,2.1.293}.isolated.{stdout.json,stderr.txt}`.
Stdout SHA-256 respectively:
`fa010389af0ef7b799804c720f6e91914bb18d79b0464bc91dce45179345794b`,
`20a62bfc15a99bd5ec6c8fc7867f9056beeba2516ed52d96e6d193025daff273`.
This is direct local host proof, NOT a successful hosted live-lane run or proof of
provider-internal effort accounting. AC-1's full hosted execution and independent
adversarial audit remain required; no acceptance criterion has been narrowed.


### Authorized harness-role correction and final surface (supersedes five-file declaration)

Hosted run `37702366584`, Claude job `113068993047`, artifact `11517724513`
(`runtime-live-e2e-claude-live-claude-haiku-5-5`) exposed this exact error:
`SPACEDOCK_LIVE_MODEL="claude-haiku-5-5", want sonnet, claude-sonnet-5-5, or claude-opus-5-5`.
The hosted Claude job **failed because of the harness role switch**, not because
Haiku was unavailable. CLI 2.1.293 was installed; candidate provenance confirmed
checkout/embedded revision `d3df726fe4123ba2d9c053fb361b828ad4850d51`, unmodified.
Seventeen common journeys failed at the model-role guard before their host launch;
AC-1's hosted half did not pass. The matrix label alone was insufficient proof.
Artifact: https://github.com/spacedock-dev/spacedock/actions/runs/37702366584/artifacts/11517724513

Finding's four evidence fields:
1. Released user/normal workflow: maintainer submits the documented Haiku dispatch.
2. Observable harm: `claudeLiveRole` rejects the registered model before common journeys.
3. `value-ac[AC-1]` requires a registered model to complete the existing Claude lane.
4. Trigger evidence: the CI error above and the two-ID switch in
   `internal/ensigncycle/shared_live_runner_test.go:94`.

Worker proposal: Material, task-owned behavior, surface reset required. The FO
separately authorized **fix**, a sixth file, and a **+150 net LOC ceiling across
six files**. This supersedes the five-file declaration above; no AC is narrowed.
Corrective code commit: `29525bea0` (same assigned branch, pushed; no main push).
Actual cumulative diff against B is **+126 net (+151/-25), six files**:

| File | Added/deleted | Net | Reason |
|---|---:|---:|---|
| `.github/workflows/runtime-live-e2e.yml` | 76/11 | +65 | Input and pre-approval resolver |
| `internal/release/live_models.txt` | 3/0 | +3 | Registration and eligibility policy |
| `internal/release/livemodels.go` | 10/8 | +2 | Exclude policies from model-ID consumers |
| `internal/release/livemodels_test.go` | 18/0 | +18 | Focused Go parser regression |
| `docs/releasing.md` | 17/0 | +17 | Operator surface |
| `internal/ensigncycle/shared_live_runner_test.go` | 27/6 | +21 | Generic point-of-use model roles and focused test |

`claudeLiveRole` now resolves IDs through `release.LiveModels`, deriving each role
as `claude-<registry suffix>`. The existing `sonnet` alias still returns
`claude-sonnet`; registered Sonnet/Opus retain `claude-sonnet`/`claude-opus`.
A future key requires no harness edit. **No inherited XFAIL:** gap selection is
exact equality on role strings, so `claude-haiku`/`claude-future` cannot match
`claude-sonnet`'s known-failure entries. No gap bindings were changed.

`TestClaudeLiveRole` uses an independent synthetic registry and checks the alias,
Sonnet, Opus, Haiku, future-key identity, and rejection of Codex/unknown IDs.
It failed on Haiku and the future key before the fix, then passed:
`go test -tags live ./internal/ensigncycle -run '^TestClaudeLiveRole$' -count=1`.
The test remains live-tagged and is excluded from the default offline suite.
No lane, job, CI step, effort change, or per-new-key branch was added.
The rerun is `37703335512` on corrected candidate `29525bea0`; it requires hosted
completion before claiming AC-1. The original run's Codex job was not modified or cancelled.


### Hosted approval-boundary proof

**AC-3 PROVED on GitHub:** run [37702369289](https://github.com/spacedock-dev/spacedock/actions/runs/37702369289)
on `d3df726fe4123ba2d9c053fb361b828ad4850d51`, inputs
`live_cadence=sonnet`, `claude_model=opus`, failed in the offline resolver at
2026-10-07T23:27:42Z with exit 1 and
`claude_model is not permitted; use live_cadence=opus-pre-release`.
Every live job was skipped. The pending-deployments endpoint returned `[]`.
No live environment approval was spent. This establishes the pre-approval guarantee
on the real platform, unlike the shell-only cases. The correction in `29525bea0`
changes only the live harness file; the negative-tested workflow/registry are byte-identical.

Both Haiku dispatches used the documented fields with `--ref` set to the assigned
branch rather than main. The observed matrix was `sonnet, claude-haiku-5-5, max, CI-E2E`.
Run 37702366584 concluded failure (Claude role-switch failure; **Codex success**).
No Codex changes or cancellation occurred. Corrected run 37703335512 passed offline;
live completion is still required. Environment approval was left to the FO/captain,
never self-issued by the worker using operator credentials.


### Final local verification

- `gofmt -w ./cmd ./internal`: exit 0. Its unrelated pre-existing two-line
  alignment change in `runtime_live_evidence_workflow_test.go` was restored to
  keep the six-file surface; the changed Go files are formatted.
- `go test ./...`: passed before the harness fix and on final candidate 29525bea0.
  One intervening retry hit the tool's 240-second bound; the next completed with
  exit 0 (ensigncycle 181.387s). It is not hidden as a passing invocation.
- `go test ./... -race`: passed before and after the harness fix (final exit 0).
- Both focused tests had observed red/green transitions; the parser test runs in
  the default suite, and the role test uses the explicit live-tagged command above.
- `git diff --check`: passed. The 25 actual-resolver cases and parsed-YAML
  invariants passed; the harness-only correction did not change their workflow bytes.


### Hosted rerun handoff at monitoring bound

Last observation **2026-10-07T23:44:29Z**, inside the 23:44:36Z bound:
run [37703335512](https://github.com/spacedock-dev/spacedock/actions/runs/37703335512)
on final SHA `29525bea0cd558b8d305fd1d4417669568b384d7` is **in_progress**,
with no conclusion. Offline job `113071744105` succeeded; Claude job
`113072132719` is in `Run live Claude E2E`; Codex job `113072132702` is in
`Run live Codex shared scenarios`; Pi is skipped. The FO approved the environments.
No corrected-run artifact is available yet. **AC-1's hosted half is in progress,
not passed**; final CLI/model observations must be recovered from that run's artifact.
The local monitoring process was stopped at the bound, not the GitHub run.
Resume without starting another run:

```bash
gh run watch 37703335512 --interval 30 --exit-status
gh run view 37703335512 --json status,conclusion,headSha,jobs
gh run view 37703335512 --log-failed
gh run download 37703335512 --name runtime-live-e2e-claude-live-claude-haiku-5-5 --dir /tmp/haiku-final-evidence
```

### Historical follow-up: corrected hosted Claude job FAILED; implementation blocked under prior acceptance

Historical assessment before the captain ruling below: its failure evidence is
retained unchanged. The ruling supersedes the acceptance/blocker interpretation,
not the failed run conclusion or any assertion error.

The completion guard's FAILED item is retained, not relabelled or removed. Both
attempts failed the hosted Haiku acceptance obligation; **AC-1 is not met**.
No replacement run was started and no code, XFAIL binding, or acceptance criterion
was changed in this follow-up. Worktree HEAD remains `29525bea0`.

| Attempt | Run | Claude job | Artifact | Outcome |
|---|---|---|---|---|
| Before role correction (`d3df726fe`) | `37702366584` | `113068993047` | `11517724513` | Claude failed; whole run concluded failure |
| Corrected candidate (`29525bea0`) | `37703335512` | `113072132719` | `11517939481` | Claude failed at 2026-10-07T23:47:08Z, exit 1; whole run concluded failure |

Final API observation at **2026-10-07T23:50:23Z**: corrected run 37703335512 is
`completed`, conclusion `failure`. Offline job 113071744105 and Codex job
113072132702 concluded success; Pi and journey-delta-comment were skipped.
The run is finished, not awaiting another result.

Both artifacts are named `runtime-live-e2e-claude-live-claude-haiku-5-5`:
- Initial: https://github.com/spacedock-dev/spacedock/actions/runs/37702366584/artifacts/11517724513
- Corrected: https://github.com/spacedock-dev/spacedock/actions/runs/37703335512/artifacts/11517939481

The initial harness error, verbatim, remains:

```text
SPACEDOCK_LIVE_MODEL="claude-haiku-5-5", want sonnet, claude-sonnet-5-5, or claude-opus-5-5
```

The corrected job reached the common journeys but failed four named journeys.
Its harness errors and findings are quoted verbatim (timestamps omitted):

```text
shared_live_runner_test.go:174: FAIL /keep-moving-posture owner= observed=[keep-moving-violation]
keep-moving-violation: durable keep-moving journeys = 2/3: map[approved-gate:active entity remains beside canonical archive]
shared_live_runner_test.go:190: FAIL /default-headless-gate-stop owner= observed=[gate-hold-violation gate-not-held implementation-worker-not-dispatched]
gate-hold-violation: gate hold crossed its committed no-authority boundary: no successful gate prepare recorded
gate-not-held: read prepared gate expectation: prepared fixture Briefing count = 0, want 1
implementation-worker-not-dispatched: implementation lifecycle incomplete: spawns=1 completed=121 validation=-1 report=<nil>
shared_live_runner_test.go:230: withdrawn recovery crossed its prepare/commit gate-stop boundary
shared_live_runner_test.go:261: FAIL /filing owner= observed=[filing-command-not-observed]
filing-command-not-observed: filing command log has no spacedock new wire-the-thing invocation
DONE 26 tests, 8 failures in 312.526s
Process completed with exit code 1.
```

The eight failures include the four journeys, three parent slots, and the suite.
The log also reports a wrong-root diagnostic for DefaultHeadlessGateStop; the full
log and stream artifacts remain available under the corrected run/job/artifact.
This is a failed hosted acceptance outcome, not a pass with caveats.

Observed execution configuration: `sonnet` cadence, `claude-haiku-5-5`, `max`,
`CI-E2E`; CLI version **2.1.293**. The max-effort shim step completed and the
archived filing stream's init/result identify `claude-haiku-5-5`, including
`canonicalModel=claude-haiku-5-5` and 533 thinking tokens / 1962 output tokens
(session `45a52cce-e6dc-4a81-b7a6-3f10c12b4d39`). These observations identify
what ran; they do not turn the failed journey assertions into hosted AC-1 success.
Candidate provenance matches checkout and embedded SHA
`29525bea0cd558b8d305fd1d4417669568b384d7`, with `vcs_modified=false`.
Corrected artifact ZIP SHA-256:
`336115cff79802d2e6d430ddd9647f72b598801361223af6e3a763d83f298c7d`.

### Captain ruling: dispatch acceptance re-scoped; compliance failure retained

The captain's ruling, relayed in this follow-up on **2026-10-08**, re-scopes the
acceptance owed by this task to: **a named model runs the lane without a workflow
edit**. This supersedes the earlier AC-1 requirement for a green/completed-success
journey suite and the prior interpretation that those journey failures block this
dispatch task. The report item changes to DONE because the captain changed the
acceptance, not because the failed run was reclassified as green.

Run `37703335512`, job `113072132719`, actually executed `claude-live` with
`claude-haiku-5-5` at `max` in `CI-E2E`, using CLI 2.1.293; the artifact's host
init and result identify that model. The harness role rejection from initial run
`37702366584` is gone. Between those runs, the correction changed only the harness
role resolver: workflow and registry bytes were unchanged. The earlier registry-only
future-key trial separately established extension without a workflow edit.
Thus the re-scoped dispatch acceptance is proved. **The corrected live run still
concluded failure: 26 tests / 8 failures. No green-suite claim is made.**

The first role-switch failure, its verbatim error, all corrected-run assertion
errors, and the older CLI's `[claude-code:unrecognized_model]` warning remain in
this entity. No code, model-compliance policy, XFAIL, or CLI pin is changed here;
no new lane run was started. Code remains `29525bea0` (+126 net across six files).

### Finding: Haiku model compliance — owned outside this dispatch task

Ownership/disposition per captain: the eight reported failures are a model-compliance
finding owned outside this task, not a dispatch-mechanism defect. The FO will
commission a separate classifier; this report does not assign a new classifier,
waive a failure, add XFAILs, or attempt a repair.

Four evidence fields:

1. **Released user and normal workflow:** a maintainer dispatches the existing
   Claude lane with `live_cadence=sonnet`, `claude_model=haiku`, maximum effort,
   and the normal `CI-E2E` environment approval. The host runs the requested model.
2. **Observable harm:** comparison live outcome **26 tests / 8 failures on
   claude-haiku-5-5**, exit 1. The failing journeys are
   `TestLiveCommonKeepMovingPosture`, `TestLiveCommonDefaultHeadlessGateStop`,
   `TestLiveCommonWithdrawnGateRecovery`, and `TestLiveCommonFiling`.
   The eight failures count those four journeys, three parent slots, and the suite;
   gate/commit boundaries, durable archive state, and filing-command observation failed.
3. **Authority:** `captain-ruling[2026-10-08]` retains these model-compliance
   failures for separate ownership/classification while accepting named-model
   execution as this task's dispatch value; it does not waive the observed
   gate-boundary or durable-workflow assertion failures.
4. **Trigger evidence:** initial run `37702366584` / job `113068993047` /
   artifact `11517724513` rejected the model in the role switch; corrected run
   `37703335512` / job `113072132719` / artifact `11517939481` executed Haiku
   and produced the eight failures. The complete verbatim harness assertions are
   retained in the historical outcome section above, including:

   ```text
   shared_live_runner_test.go:261: FAIL /filing owner= observed=[filing-command-not-observed]
   filing-command-not-observed: filing command log has no spacedock new wire-the-thing invocation
   ```

Corrected artifact path:
https://github.com/spacedock-dev/spacedock/actions/runs/37703335512/artifacts/11517939481
(`runtime-live-e2e-claude-live-claude-haiku-5-5`). Inside it, the per-journey
streams are under
`spacedock/spacedock/live-artifacts/claude/claude-haiku-5-5/claude-shared-scenarios/`,
including `filing/claude-stream.jsonl`; the job log records all assertion errors.
The initial artifact remains at
https://github.com/spacedock-dev/spacedock/actions/runs/37702366584/artifacts/11517724513.

## Stage Report: implementation

- DONE: Reach the goal with the smallest change that works. Remove machinery rather than add it. Drop permitted-set metadata or the proposed test file if it does not earn its place, and say why.
  Code 29525bea0 is +126 net across six authorized files; policy is necessary for generic eligibility, duplicate schema/wiring tests were cut, and the sixth file removes the discovered point-of-use allowlist.
- DONE: Keep the entry point model-agnostic: a parameter selects a registry key, and model ids stay in live_models.txt only.
  Actual resolver passed registry-only future-key extension; TestClaudeLiveRole passed a synthetic future ID while preserving Sonnet/Opus identities and rejecting non-Claude IDs.
- DONE: Keep the approval boundary: refuse an explicit model that the selected cadence does not cover, inside the offline job, so no environment approval is spent. The failure must name the remedy.
  AC-3 PROVED: hosted run 37702369289 failed offline with use live_cadence=opus-pre-release; all live jobs skipped and pending_deployments=[]; tested workflow is unchanged in final candidate.
- DONE: Describe the architecture accurately as model-agnostic resolution with cadence-dependent eligibility, and say so in the task body.
  The body and operator docs state that the model is a parameter while cadence gates eligibility and the approval environment, not merely defaults.
- DONE: Record the local probe: claude-haiku-5-5 ran with canonicalModel claude-haiku-5-5, thinking active, exit 0, using claude 2.1.283. Record the stderr warning [claude-code:unrecognized_model]. Test whether an older local CLI causes it, and state the result. Do not silently accept the warning.
  Reproduced on 2.1.283 (264 thinking tokens); isolated 2.1.293 accepted max, returned ok/exit 0 with 687 thinking tokens and empty stderr; commands, sessions, hashes, and limits are recorded above.
- DONE: Report the actual surface: files and net LOC with a re-declared tolerance, and note that both the workflow and the registry change.
  Authorized six-file declaration supersedes five files; +151/-25 = +126 net under +150 ceiling, with per-file reasons and both workflow/registry changes.
- DONE: Run gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race, and report each result.
  All passed; one intervening ordinary-suite invocation exceeded its 240s tool bound, then retry passed; focused parser and live-tagged role tests both have red/green evidence.
- DONE: Add no lane, no job, and no CI step.
  Parsed YAML comparison preserves all five jobs and step counts, cadence approval expressions, effort, dependencies, and Codex/Pi definitions; role fix changes no XFAIL bindings.
- DONE: Report in the canonical item form, one DONE/SKIPPED/FAILED line per checklist item with an evidence or rationale line, ending with a non-empty Summary.
  This report records code commits d3df726fe and 29525bea0, both failed run outcomes, and the captain-authorized acceptance reset; DONE applies only to the re-scoped named-model dispatch obligation.
- DONE: A named model runs the lane without a workflow edit.
  Captain re-scoped the acceptance: initial run 37702366584/job 113068993047/artifact 11517724513 failed with `SPACEDOCK_LIVE_MODEL="claude-haiku-5-5", want sonnet, claude-sonnet-5-5, or claude-opus-5-5`; corrected run 37703335512/job 113072132719/artifact 11517939481 executed claude-haiku-5-5 at max with no further workflow edit and no role rejection. Its 26 tests / 8 failures remain an outside-task model-compliance finding, not a green run.
- SKIPPED: Claim completed hosted AC-1 on the corrected candidate within this worker's bound.
  The initial bound ended at 23:44:29Z with run 37703335512 in progress; this follow-up records Claude job 113072132719 failing at 23:47:08Z. No green-suite claim or replacement run; the later captain ruling changes the acceptance, not this outcome.
- SKIPPED: Independent adversarial audit.
  FO owns and will orchestrate the required independent review; worker did not self-audit or self-approve environments.

### Summary

Implemented registry-key selection and cadence eligibility, including the authorized generic harness-role correction; code remains +126 net LOC across six files and both Go suites pass. The captain re-scoped acceptance to a named model running the lane without a workflow edit, proved by run 37703335512 at Haiku/max, so that one item is DONE. Both failed run outcomes, the initial role-switch rejection, the corrected run's 26 tests / 8 failures (four failing journeys), and the unrecognized_model warning remain recorded; model compliance is owned outside this task for separate classification, and no further run was started.

## Review-finding disposition

### Validation: three central claims lack a standing proof owner

**Recommendation: REJECTED. Defect kind: evidence defect. Release scope: material. Task ownership: this dispatch task.**
The candidate's observed behavior passed the one-off checks below. The defect is in its regression proof, not Haiku compliance or the resolution architecture.

Four evidence fields:
1. **Released user and normal workflow:** a maintainer selects a registered Claude key and its permitted approval cadence through the existing manual dispatch.
2. **Observable harm:** the shipped checks stay green when the offline dependency is removed, routine policy admits Opus, or the matrix ignores the selected key. They cannot defend the claimed dispatch and approval boundaries.
3. **Authority:** `contract[docs/dev/README.md#proof-policy]` requires the detached high-stakes audit to catch claim-breaking edits; these survivors leave AC-1, AC-2, and AC-3 without a standing proof owner for their central claims.
4. **Trigger evidence:** detached checkout of `29525bea0`; apply each mutation below separately, run the stated command, then restore the file. These are actual green runs, not inferred coverage gaps.

| Adversarial edit in the throwaway checkout | Command | Result / duration |
|---|---|---|
| First `needs: offline` (Claude job) becomes `needs: []` | `go test ./internal/release ./internal/contractlint -count=1` | **exit 0**, 19.653s; AC-3 dependency bypass survives |
| `claude.allowed.sonnet=sonnet,haiku` becomes `claude.allowed.sonnet=sonnet,opus,haiku` | Same command | **exit 0**, 14.530s; AC-2 routine permission expansion survives |
| Matrix `needs.offline.outputs.claude_selected` becomes `needs.offline.outputs.claude_sonnet` | Same command | **exit 0**, 14.301s; AC-1 parameter bypass survives |
| Reader's policy-filter condition becomes `false` | `go test ./internal/release -run '^TestParseLiveModelsExcludesCadencePolicy$' -count=1` | **exit 1**, 0.799s; independent expected map catches policy leakage |
| Generic role return becomes `return "claude-sonnet", nil` | `go test -tags live ./internal/ensigncycle -run '^TestClaudeLiveRole$' -count=1` | **exit 1**, 1.772s; Haiku/future/Opus identities fail |

**FO-authorized disposition:** record REJECTED; keep the candidate unchanged and route the three new standing assertions to the captain. No acceptance of this hole was authorized.
**Cause and accountability, as directed by the FO:** the FO authorized the earlier coverage cut on independent review advice. That review requested removal of duplicated schema checks while retaining parser and wiring coverage. The authorization removed the wiring coverage too. This was the FO's authorization error, not this validator's error.

**Three NEW STANDING CHECKS — captain approval required:** one assertion for Claude's `needs: offline`; one excluding `opus` from the routine permitted set; one binding the matrix to the selected registry output (`claude_selected`, not the fixed `claude_sonnet` output). Each must fail under its own mutation above. Keep them in `internal/release/livemodels_test.go`; do not restore a broad schema/wiring suite. No lane, job, or CI step is proposed.

### Measured surface

Measured diff: `git diff --numstat origin/main HEAD`, baseline `7558f2e0`, candidate `a9f7e8143` (implementation cycle 2; validation measured `29525bea0` at +126).

| File | Added/deleted | Net |
|---|---:|---:|
| `.github/workflows/runtime-live-e2e.yml` | 76/11 | +65 |
| `docs/releasing.md` | 17/0 | +17 |
| `internal/ensigncycle/shared_live_runner_test.go` | 27/6 | +21 |
| `internal/release/live_models.txt` | 3/0 | +3 |
| `internal/release/livemodels.go` | 10/8 | +2 |
| `internal/release/livemodels_test.go` | 42/0 | +42 |
| **Total** | **175/25** | **+150** |

The three standing assertions add +24 net, all in `internal/release/livemodels_test.go`. The same six files remain. The total is exactly at the +150 ceiling, not over it. The assertions were compacted to fit; none was weakened, and each still fails under its own audit mutation (see Stage Report: implementation (cycle 2)).

### Validation cycle 2: adjacent mutations outside the three approved checks

Reviewer observation; proposed class **Deferred risk** for 1-4 and **Polish** for 5. FO authorization pending; no candidate change requested.
Probe: throwaway clone of `a9f7e8143`, one edit at a time, `go test ./internal/release ./internal/contractlint -count=1`.

1. Resolver ignores the key (`selected = records["claude." model]` -> `records["claude.sonnet"]`): exit 0. New behavior with no standing test.
2. `always() &&` prefixed to the `claude-live` `if:` (live job runs after offline refusal): exit 0. Guard line unchanged from B.
3. Opus-cadence environment `'CI-E2E-OPUS'` -> `'CI-E2E'`: exit 0. Expression unchanged from B.
4. `claude.big=claude-opus-5-5` registered and added to `claude.allowed.sonnet`: focused test and resolver pass. Assertion 3 guards the key `opus`, not the Opus ID.
5. `needs: [offline]` (semantically equal) fails assertion 1. Stricter than needed; fails closed.

Four evidence fields (1-4): released workflow is a maintainer editing the resolver, guard, environment expression, or registry policy; harm is a regression that keeps CI green while removing parameter precedence or refusal-before-approval; authority `captain-ruling[2026-10-21]` limits standing proof to three named assertions with no broad schema/wiring suite; trigger is hypothetical (each needs a future edit), and current candidate behavior is verified by the actual-resolver spot-check and hosted runs 37702369289 / 37703335512.
Promote to material when a change touches the resolver step, the `claude-live` `if:`/environment expressions, or adds a second key resolving to an Opus ID; or when the captain extends the standing-proof scope.

## Stage Report: validation

- DONE: Confirm the six-file surface and net LOC against the re-declared tolerance, measured with git diff --numstat against origin/main.
  Candidate 29525bea0 versus origin/main 7558f2e0 is +151/-25 = +126 net, six files, within +150; the measured split and captain-pending prospective reset are above.
- DONE: Re-run gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race, and report each result with durations.
  All exit 0: gofmt 0.127s, ordinary suite 230.262s, race suite 252.952s; formatting ran on the detached copy to preserve candidate bytes and exposed only the pre-existing two-line alignment in runtime_live_evidence_workflow_test.go. Candidate suites ran unchanged.
- DONE: Confirm the model is a parameter and that model ids live only in live_models.txt.
  Actual offline shell passed 33 independent exit/output/atomicity cases; a registry-only future key resolved fixture-future-model with identical workflow bytes. Tracked-file ID search found no pinned live IDs outside live_models.txt.
- DONE: Confirm the cadence gates both eligibility and the environment, and that an out-of-cadence model is refused before any approval is spent. Cite the hosted negative run.
  Run 37702369289 at d3df726fe: routine/opus failed offline with use live_cadence=opus-pre-release; all live jobs skipped, approvals=[] and pending_deployments=[]; final workflow/registry are byte-identical to that run.
- DONE: Confirm the harness resolves the model role generically, that sonnet, claude-sonnet, and claude-opus are unchanged, and that a new key inherits no known-failure binding.
  TestClaudeLiveRole passed with -tags live (1.306s); synthetic Sonnet/Opus/Haiku/future IDs assert distinct roles, alias sonnet stays claude-sonnet, and exact gap.target equality cannot transfer Sonnet's unchanged XFAILs to new suffixes.
- DONE: Confirm the restated acceptance item is honest: both run identifiers, the harness error, and the eight journey failures remain in the record. Say plainly whether any failure evidence was softened or removed.
  Fresh GitHub logs confirm 37702366584's 17 role-guard errors and 37703335512's 26 tests/8 failures: four journeys, three parent slots, suite (plus a package fail event). No failure evidence was softened or removed; the captain changed acceptance, not either failed conclusion.
- DONE: Confirm the model-compliance finding carries its four evidence fields and an owner outside this task.
  The retained finding names normal dispatch, concrete harm, captain-ruling[2026-10-08], and both run/job/artifact triples; ownership remains outside this task. Its four named failed journeys match the downloaded JSONL failure events.
- DONE: Run the workflow's detached adversarial audit on a throwaway checkout, since the CI and release machinery is a declared high-stakes surface.
  Detached 29525bea0 checkout under /tmp caught policy leakage and role inheritance but missed three central-claim mutations; exact commands/results and the FO's REJECTED disposition are recorded above. No candidate mutation occurred.
- DONE: Reject any evidence that asserts an artifact's own text against a copy of that text.
  No self-copy assertion is accepted: fixture outputs and exit codes are independent expectations; hosted host-init/result events prove actual identity, not artifact labels; separate baseline bytes establish only unchanged wiring. Green suite counts do not prove the three uncovered claims.
- DONE: Report with a verdict of PASSED or REJECTED, in the canonical item form, ending with a non-empty Summary.
  Verdict: REJECTED for a material evidence defect, not for the captain-owned model-compliance outcome. The FO owns the earlier authorization error; three focused new standing checks require captain approval.
- DONE: AC-1 (VALUE), as superseded by the captain: a named model runs the lane without a workflow edit.
  Artifact 11517939481 (SHA-256 336115cff79802d2e6d430ddd9647f72b598801361223af6e3a763d83f298c7d) and job logs verify Haiku init/canonicalModel, CLI 2.1.293, max/CI-E2E, clean embedded/checkout SHA 29525bea0; initial-to-corrected change is harness-only. Hardcoding the matrix is the uncovered regression falsifier, not a current dispatch failure.
- DONE: AC-2 and AC-3 current behavior: existing cadence defaults, approval authority, and explicit pre-approval refusal.
  Parsed candidate/baseline comparison preserves five jobs, step counts, guards, dependencies, environments, max effort and Codex/Pi definitions; local PR/routine/pre-release/Pi defaults and forbidden Opus exercise actual resolver output. GitHub confirms four reviewer-protected environments and expected secret names; dependency/policy mutation survivors block regression-proof acceptance.
- DONE: AC-4: operator instructions select a key without duplicating IDs or weakening release evidence.
  Corrected offline job 113071744105 logs workflow_dispatch/sonnet/haiku, matching documented fields; host stream independently confirms the requested ID. Default release command is unchanged; docs separate CLI version, cadence permission, and exploratory versus release evidence. No new dispatch was made.
- DONE: Semantic adversarial pass over parameter, registry, output, host argv, role, approval and terminal state.
  The 33-case shell matrix adds whitespace/Unicode/newline inputs, CRLF, absent final newline and record reordering to defaults/overrides, malformed/duplicate/empty policy data, injection, and generic extension; rejects emit zero outputs. Registry work is linear over reviewed records; no new blocking/network hot path or supported scaling hazard was found.
- FAILED: Central dispatch and approval claims have a standing regression proof owner.
  Removing offline dependency, allowing routine Opus, or forcing the Sonnet matrix output all leave release/contractlint checks green. Existing parser/role tests have genuine independent oracles but do not defend these boundaries; no broad schema suite is requested.
- SKIPPED: New live lanes, environment approvals, repeated Sonnet classification, and claims not observable in this validation.
  No run/approval/classifier was started. Historical reviewer/secret-value snapshots and fresh hosted executions of every blank cadence were not independently verified; current metadata and local defaults are not substitutes for those observations. Compliance remains failed and outside-task, not waived.

### Summary

REJECTED: current named-model dispatch and hosted pre-approval refusal are independently verified, and both Go suites pass, but three claim-breaking edits survive the detached audit. The FO accepts responsibility for the earlier coverage cut and routes three focused new standing assertions to the captain; their prospective LOC increment is unmeasured. Candidate 29525bea0 stays unchanged, both failed Haiku runs and all eight failed test outcomes remain unsoftened, and temporary audit artifacts are removed after committing this report.

## Stage Report: implementation (cycle 2)

- DONE: Three standing assertions in internal/release/livemodels_test.go (needs: offline, opus excluded from routine set, matrix bound to claude_selected), each shown red under its own audit mutation and green on the candidate.
  Commit a9f7e8143 adds `TestClaudeLiveKeepsCadenceBoundaryAndSelectedModel`. It parses the workflow YAML and the embedded registry. Audit: throwaway detached worktree of a9f7e8143, command `go test ./internal/release ./internal/contractlint -count=1`, each mutation applied alone and then reverted with `git checkout`:

  | State | Exit | Duration | Failing assertion |
  |---|---:|---:|---|
  | Candidate, unmutated | 0 | 81.6s | none |
  | Claude job `needs: offline` becomes `needs: []` | 1 | 28.6s | only `livemodels_test.go:34` (claude-live needs = []) |
  | `claude.allowed.sonnet=sonnet,opus,haiku` | 1 | 39.6s | only `livemodels_test.go:40` (routine policy permits opus) |
  | Matrix model becomes `needs.offline.outputs.claude_sonnet` | 1 | 27.5s | only `livemodels_test.go:37` (matrix not bound to claude_selected) |
  | Restored after all three | 0 | 52.1s | none |

  What each assertion checks: (1) the `claude-live` job's `needs` is exactly `offline`. (2) The `claude.allowed.sonnet` record exists and its members exclude `opus`. (3) The matrix `model` names an offline output, and that output is bound to `steps.live_models.outputs.claude_selected`. Assertion 3 compares the matrix with a separate part of the workflow. It does not compare the matrix text with a copy of itself.
- DONE: Measured diff vs origin/main stays within +150 net LOC across the same six files, with no lane/job/CI step or behavior change.
  `git diff --numstat origin/main HEAD` (7558f2e0..a9f7e8143) is 175/25, +150 net, across the same six files. Only `livemodels_test.go` changed (+24). The total is exactly at the ceiling. To fit, the assertions were compacted, not weakened: the YAML parse folds into the first check, and the policy cut uses an if-init. The prospective reset under `## Review-finding disposition` is now this measured table. The workflow, registry, harness, and CI steps are unchanged.
- DONE: gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race all pass on the final commit.
  gofmt exit 0 (0.38s). `go test ./...` exit 0 (453s). `go test ./... -race` exit 0 (499s). All runs were on a9f7e8143. gofmt also rewrote two alignment lines in `internal/release/runtime_live_evidence_workflow_test.go`. That formatting drift already exists on origin/main (`git show origin/main:… | gofmt -l` lists the file), so I reverted it and kept it out of the six-file surface.

### Summary

One test function now holds the three captain-approved standing assertions. Each one fails only under its own audit mutation and passes on candidate a9f7e8143. Both Go suites pass. The surface is the same six files at exactly +150 net, which is the ceiling. Haiku model compliance and the pre-existing gofmt drift on main stay outside this task.

### Feedback Cycles

- Cycle 1: REJECTED — validation detached audit (three central dispatch/approval claims had no standing regression check); surface 6 files / +150 net LOC vs estimate 5 files / +180 net LOC (83%); AC narrowed: AC-1 re-scoped by captain-ruling[2026-10-08] to named-model execution without a workflow edit, six-file surface authorized earlier. Captain revised the validation gate (66bb17a16) and approved three standing assertions within +150; producer correction a9f7e8143 is complete, independent re-review pending.
- Cycle 2: REJECTED — captain at validation gate attempt-2 (validator had recommended PASSED); break-glass selected-bare XFAIL was bound to claude-sonnet unconditionally, so the task's "No inherited XFAIL" claim was false and Haiku inherited a Sonnet known failure; surface 7 files / +170 net LOC vs estimate 5 files / +180 net LOC (94%); AC narrowed: AC-1 re-scoped by captain-ruling[2026-10-08]; captain raised the ceiling to +170 and the seventh file is the one the fix names. Producer correction 03123dd4f is complete, independent re-review pending.
- Cycle 3: REJECTED — captain at validation gate attempt-3 (validator had recommended PASSED) set scope directly: XFAIL Haiku's five observed failures with blank owners pending analysis, make withdrawn-gate-recovery gap-aware, sync the registry doc; surface 10 files / +220 net LOC vs estimate 5 files / +180 net LOC (122%); AC narrowed: AC-1 re-scoped by captain-ruling[2026-10-08]; captain raised the ceiling to +200 then +220 and directed the added files. Producer correction 861a6e13a is complete; reviewer re-run held for the captain under the cycle-3 rule.

## Stage Report: validation (cycle 2)

- DONE: Independently re-run the three audit mutations (Claude job needs: offline -> [], routine allowed set admits opus, matrix reads claude_sonnet) on a throwaway checkout of a9f7e8143: each must fail only its own new assertion, and the unmutated candidate must pass; the assertions must check real values against an independent source, not text against a copy of itself.
  Fresh clone at a9f7e8143 under /tmp, `go test ./internal/release ./internal/contractlint -count=1`: unmutated exit 0 (28.9s). `needs: []` exit 1, only `livemodels_test.go:34`; `sonnet,opus,haiku` exit 1, only `:40`; matrix `claude_sonnet` exit 1, only `:37`. No other test failed in any run. Oracles: (1) parsed `needs` vs the constant `offline`; (2) matrix output name looked up in the offline job's separate `outputs` map, which must bind `claude_selected`; (3) embedded registry policy vs the ruled exclusion of `opus`. None compares text with a copy of itself. Clone removed.
- DONE: Confirm surface (6 files, +150 net vs origin/main, within the captain ceiling), gofmt, go test ./..., and go test ./... -race on a9f7e8143, and that workflow/registry/harness behavior is unchanged from 29525bea0.
  `git diff --numstat 7558f2e0 a9f7e8143`: 6 files, +175/-25 = +150, at the ceiling. `go test ./... -count=1` exit 0 (538s); `-race` exit 0 (467s); 21 packages ok, zero FAIL lines. `gofmt -l` lists only `runtime_live_evidence_workflow_test.go`, which origin/main also lists. `git diff --name-only 29525bea0 a9f7e8143` = `livemodels_test.go` only; workflow, registry, and reader are byte-identical to d3df726fe.
- DONE: Re-check AC-1..AC-4 evidence against the current task body and captain rulings, then recommend PASSED or REJECTED with material findings separated from deferred risks.
  AC-1 (re-scoped 2026-10-08): `gh run view 37703335512` shows head 29525bea0, job `claude-live (sonnet, claude-haiku-5-5, max, CI-E2E)`. Its 8 journey failures stay owned outside this task. The actual resolver from the candidate workflow selects claude-haiku-5-5 for sonnet+haiku. A registry-only key `zeta=zz-synthetic-9` resolves with workflow bytes unchanged. AC-2: PR+haiku still yields Sonnet; blank sonnet yields Sonnet and blank pre-release yields Opus; pi+haiku is refused; assertions 1-3 now hold the boundary. AC-3: run 37702369289 shows offline failure and all live jobs skipped; pending_deployments length is 0; the resolver prints the `use live_cadence=opus-pre-release` remedy with zero outputs. AC-4: the docs diff uses registry keys, keeps `claude_version` as the CLI pin, and keeps release evidence intact. No material finding. Four deferred risks and one polish item are recorded under `## Review-finding disposition`.

### Summary

PASSED recommendation. The cycle-1 evidence defect is closed: each of the three captain-approved assertions fails only under its own mutation and passes on a9f7e8143. Surface is 6 files at +150 net. Both Go suites pass, and product behavior is unchanged from 29525bea0. Adjacent survivors (resolver precedence, `always()` guard, Opus environment expression, Opus-ID alias key) fall outside the captain's three-check scope. They are recorded as deferred risks with promote conditions, not as blockers. Haiku model compliance stays a separate finding.

## Stage Report: implementation (cycle 3)

- DONE: The break-glass selected-bare XFAIL binds only on the claude-sonnet target via the same target match liveJourney uses, and a focused check shows claude-haiku inherits no binding (red without the guard, green with it).
  Commit 03123dd4f. The commit adds `selectLiveGap(gaps, target)` (exact `gap.target == target`) to `shared_live_runner_test.go`, and `liveJourney` now calls it. `TestLiveBreakGlassShimRecovery` resolves its target with `claudeLiveRole(runner.model())`, and selected-bare takes `breakGlassSelectedBareGap(target)`, which is `selectLiveGap` over the single Sonnet entry `060xp69y61yhrww23g3wvwqy`. `TestBreakGlassSelectedBareGapBindsOnlySonnet` asserts that claude-sonnet still gets that xfail and that claude-haiku gets a zero gap. Command: `go test -tags live -count=1 -run '^TestBreakGlassSelectedBareGapBindsOnlySonnet$' ./internal/ensigncycle`. Each mutation was applied alone and then restored:

  | State | Exit | Duration |
  |---|---:|---:|
  | Guard in place | 0 | 2.8s |
  | Helper reverted to unconditional `liveXFail("claude-sonnet", …)` (the pre-fix line) | 1 | 1.2s |
  | `gap.target == target` removed from `selectLiveGap` | 1 | 1.2s |
  | Restored | 0 | 0.4s |

  Both red runs fail at `dispatch_recovery_live_test.go:83`, showing claude-haiku={kind:xfail target:claude-sonnet owner:060xp69y61yhrww23g3wvwqy}. `TestClaudeLiveRole` (-tags live): exit 0, 0.4s.
- DONE: The record corrects the false "No inherited XFAIL" claim and routes Haiku's break-glass selected-bare failure to the Haiku follow-up; no other bindings, workflow, registry, or CI change.
  Correction: the claim under `### Authorized harness-role correction and final surface` ("**No inherited XFAIL:** gap selection is exact equality on role strings …") was false for break-glass selected-bare until 03123dd4f. The same applies to validation cycle 1's DONE item "a new key inherits no known-failure binding". `dispatch_recovery_live_test.go:108` assigned the Sonnet xfail without checking the target, so Haiku run 37703335512 (job 113072132719, artifact 11517939481) graded its selected-bare failure as a `claude-sonnet` XFAIL owned by 060xp69y61yhrww23g3wvwqy. Haiku fails five journeys, not four. Haiku's break-glass selected-bare failure is now unowned. It is routed to the separate Haiku model-compliance follow-up (`### Finding: Haiku model compliance — owned outside this dispatch task`). No claude-haiku binding was added. No other gap binding, workflow, registry, resolver, lane, job, or CI step changed. The commit touches only the two `internal/ensigncycle` test files.
- DONE: Measured diff vs origin/main is within +170 net; gofmt, go test ./..., go test ./... -race, and the live-tagged focused checks pass.
  Surface was measured with the merge-base diff `git diff --numstat origin/main...HEAD`: merge-base 7558f2e09, HEAD 03123dd4f, +204/-34 = +170 net across 7 files, which is exactly the captain's raised ceiling. origin/main has advanced to 230c77e43, so a two-dot diff would also count unrelated upstream changes. The seventh file, `internal/ensigncycle/dispatch_recovery_live_test.go` (+16/-1), is the file the captain's ask names. `shared_live_runner_test.go` moved from +27/-6 to +40/-14. `gofmt -w ./cmd ./internal`: exit 0. Its pre-existing rewrite of `internal/release/runtime_live_evidence_workflow_test.go` (also on origin/main) was reverted and kept out of the surface. `go test ./... -count=1`: exit 0, 182.7s, 21 packages ok. `go test ./... -race -count=1 -p 2`: exit 0, 218.8s, 21 packages ok, no DATA RACE. `-p 2` lowered peak temp-disk use on a host with about 2.4Gi free. An earlier `go test ./...` attempt at 553Mi free failed only on disk exhaustion (608 "no space left on device"/"os error 28" lines in 5 packages). It was rerun, not counted. No live run.

### Summary

Commit 03123dd4f binds the break-glass selected-bare known failure only on the claude-sonnet target through `selectLiveGap`, the exact target match `liveJourney` now shares. A focused check fails under either guard removal and passes with the guard. The record now says the earlier "No inherited XFAIL" claim was false for this journey until this fix. Haiku's selected-bare failure goes, unowned, to the Haiku compliance follow-up. The surface is +170 net against the merge-base, at the ceiling, and both Go suites pass.

## Stage Report: validation (cycle 3)

- DONE: On a throwaway checkout of 03123dd4f, independently confirm no Claude target other than claude-sonnet inherits any known-failure binding: enumerate every liveXFail/liveTODO binding and every direct gap assignment, and show the break-glass selected-bare guard and selectLiveGap target match each fail a focused check when removed.
  Detached worktree /tmp/val3-f77 (removed). An AST scan of every `internal/ensigncycle/*.go` found 14 bindings: 13 liveXFail, 0 liveTODO. Targets are only `claude-sonnet` (5) and `pi` (9). Direct gap assignments: `dispatch_recovery_live_test.go:123` (now `breakGlassSelectedBareGap(target)`) and the two pi-stub fixtures in `durable_journey_gap_live_test.go:43,55`. For each `claude.*` registry ID, `liveDriverForRuntime` gave targets claude-sonnet, claude-opus, and claude-haiku. Only claude-sonnet got any journey or break-glass binding. Focused test `go test -tags live -run '^TestBreakGlassSelectedBareGapBindsOnlySonnet$' ./internal/ensigncycle`: unmutated exit 0. With the helper returning the unconditional `liveXFail` again: exit 1 at `:83`, claude-haiku given the Sonnet xfail. With `gap.target == target` changed to `true` in `selectLiveGap`: exit 1, same assertion. Restored: exit 0. Haiku run 37703335512's `live-e2e-detail.jsonl` contains `XFAIL claude-sonnet/break-glass-shim-selected-bare owner=060xp69y61yhrww23g3wvwqy`, which confirms the defect that was fixed.
- DONE: Re-confirm the cycle-1 standing assertions still fail under their three audit mutations, surface (7 files, +170 net vs merge-base 7558f2e09, within the captain ceiling), gofmt, go test ./..., and go test ./... -race (use -p 2 if disk is tight), and that workflow/registry/resolver behavior is unchanged.
  `go test ./internal/release ./internal/contractlint -count=1`: unmutated exit 0. Claude `needs: []`: exit 1, only `livemodels_test.go:34`. `claude.allowed.sonnet=sonnet,opus,haiku`: exit 1, only `:40`. Matrix `claude_sonnet`: exit 1, only `:37`. `git diff --numstat 7558f2e09 03123dd4f`: 7 files, +204/-34 = +170, exactly at the ceiling. `gofmt -l` lists only `runtime_live_evidence_workflow_test.go`, which is unformatted at the merge-base too. `go test ./... -count=1 -p 2`: exit 0, 191s, 21 ok. `-race -p 2`: exit 0, 209s, 21 ok, no DATA RACE, no disk errors. `a9f7e8143..03123dd4f` touches only the two ensigncycle test files. Workflow, registry, reader, and docs are byte-identical to d3df726fe, the SHA of hosted AC-3 run 37702369289. The extracted actual resolver gives: PR+opus selects Sonnet; blank routine selects Sonnet; sonnet+haiku selects claude-haiku-5-5; blank pre-release selects Opus. sonnet+opus, pi+haiku, and a raw ID each exit 1 with their remedies and zero outputs.
- DONE: Re-check AC-1..AC-4 and the corrected record against captain rulings, then recommend PASSED or REJECTED with material findings separated from deferred risks.
  AC-1 (re-scoped 2026-10-08): run 37703335512 ran claude-haiku-5-5 with no workflow edit. Since then only harness test files changed. AC-2: defaults and approval authority are unchanged, and assertions 1-3 hold. AC-3: hosted refusal on byte-identical workflow and registry. AC-4: docs use registry keys and keep `claude_version` as the CLI pin. The captain's 2026-10-22 asks: (1) satisfied, because `liveJourney` and break-glass share `selectLiveGap`; (2) red/green shown above; (3) corrected in the cycle-3 implementation report. Material findings: none. Deferred risks and one polish item:
  1. Deferred risk: the call-site target is not pinned. Replacing `claudeLiveRole(runner.model())` with `claudeLiveRole("sonnet")` keeps the focused check green. Restoring the literal pre-fix line fails the `-tags live` build (`target` declared and not used), so the live lane would catch it. Trigger: a future edit to `TestLiveBreakGlassShimRecovery`. Promote to material if that function changes.
  2. Deferred risk: `TestBreakGlassSelectedBareGapBindsOnlySonnet` and `TestClaudeLiveRole` run in no CI job. The default suite lacks `-tags live`, and the Claude lane runs only `^TestLiveScheduled$`. They are manual regression checks, the same placement disclosed for `TestClaudeLiveRole` in cycle 1. Promote to material if the captain requires automatic execution.
  3. Deferred risk: at B, Opus pre-release inherited the Sonnet-owned selected-bare xfail unconditionally. It now has none, as the captain ruled, because owner 060xp is "Repair Sonnet keep-moving and selected-bare live flakes". An Opus selected-bare failure now fails the pre-release run. Promote to material if a pre-release run fails selected-bare.
  4. Deferred risk: `claudeLiveRole` ranges over a map, and the reader has no duplicate-ID check. Two `claude.*` keys that share an ID would give a nondeterministic role and gap set. This extends cycle-2 item 4. Promote to material on any duplicate-ID registry record.
  5. Polish: the original "**No inherited XFAIL**" sentence and the Haiku finding's "four failing journeys" are corrected only in the appended cycle-3 report. Recommend the FO add an inline pointer so the Haiku follow-up sees five journeys, including the unowned selected-bare one.
  The `liveJourney` refactor changes todo-skip semantics only for gap lists with duplicate targets, and the enumeration found none.

### Summary

PASSED recommendation. On a detached checkout of 03123dd4f, every known-failure binding was enumerated. Only claude-sonnet receives any binding; claude-opus and claude-haiku receive none. The focused check fails when either the break-glass guard or the `selectLiveGap` target match is removed. The three cycle-1 assertions still fail only under their own mutations. The surface is 7 files at +170, at the ceiling, and both Go suites pass. Workflow, registry, and resolver are byte-identical to the hosted-proven SHA. Four deferred risks and one record-polish item are listed above; none blocks.

## Stage Report: implementation (cycle 4)

- DONE: Haiku's five observed failures are bound as liveXFail("claude-haiku", "") with Sonnet/Pi/Opus bindings unchanged, and the focused check shows claude-haiku gets its own blank-owner xfail, never Sonnet's (red under the old mutation, green as fixed).
  Commit 43547cabb. The blank-owner binding is on keep-moving-posture, default-headless-gate-stop, withdrawn-gate-recovery (`withdrawnGateRecoveryGaps`), filing, and break-glass selected-bare (`breakGlassSelectedBareGap` = `selectLiveGap` over {sonnet 060xp69y61yhrww23g3wvwqy, haiku ""}). `TestBreakGlassSelectedBareGapBindsPerTarget` (`go test -tags live`) asserts sonnet={xfail,claude-sonnet,060x…}, haiku={xfail,claude-haiku,""}, opus=zero. Sonnet-for-everyone mutation: exit 1, and haiku and opus both get the Sonnet binding. Cycle-2 Sonnet-only list: exit 1, haiku gets zero. Fixed: exit 0.
- DONE: Blank-owner bindings are shown to compile, select, grade XFAIL, and pass every existing check; the registry doc gains the five claude-haiku rows and the four missing jfdrr2 Sonnet rows, verified against code.
  `go vet -tags live ./internal/ensigncycle` passes. The same focused test grades a failing semantic as `xfail` and a pass as `xpass` with `liveGradeFailsLane` false. With the gap dropped, it grades fail/pass and the test goes red. Only one check parses binding syntax: contractlint `parseLiveGap`, which would reject `""` (ownerID `^[a-z0-9]{24}$`; a throwaway probe returned "malformed liveXFail binding"). It is called only on synthetic literals, so no real binding is parsed and contractlint passes. Deferred risk: this becomes material if `parseLiveGap` is wired to the real journey bindings. `journeymetrics` Owner is omitempty and has no production consumer. Each jfdrr2 row's code was taken from its journey comment, and each Haiku code from the run 37703335512 log in this entity.
- DONE: Withdrawn-gate-recovery is gap-aware (captain disposition, ceiling raised to +220): its checks grade through `finishLiveScenario`; Haiku's blank-owner binding owns a boundary failure as XFAIL; unbound targets still FAIL.
  Commit 861a6e13a. In `runClaudeWithdrawnGateRecoveryScenario` (`claude_live_runner_test.go`), every former `t.Fatalf` is now a graded finding with its original message. The codes are `withdrawn-recovery-record` (entity read, attempt shape), `withdrawn-recovery-room` (withdrawn bytes rewritten, successor not a one-file room, request-digest set), `withdrawn-recovery-boundary` (prepare/commit gate-stop crossed, log included), and `withdrawn-recovery-status`. The runner reads `command.log` tolerantly: a missing log means the model never ran spacedock, so it grades as a boundary finding, not a harness fatal. With no binding, Sonnet, Opus, and Pi still FAIL through `finishLiveGrade`'s `t.Fatalf`, and the same messages appear under "Findings". New `TestWithdrawnRecoveryGradesThroughGap` (`-tags live`, stub driver, no model) checks two things. First, a boundary crossing grades `xfail` for claude-haiku and `fail` for claude-sonnet, with code `withdrawn-recovery-boundary`, and a clean log grades nil. Second, it drives the real runner end to end under Haiku's binding and gets `XFAIL claude-haiku/withdrawn-gate-recovery observed=[withdrawn-recovery-boundary withdrawn-recovery-record]`. Mutations, re-run on the final text: (A) the runner grades with `gradeLive(false, …)` and ignores the gap: exit 1. (B) the pre-fix `t.Fatalf` on the boundary: exit 1. (C) an unlabeled boundary error: exit 1. (D) the Haiku binding removed: exit 1. Fixed: exit 0. Registry row now reads "claude-haiku red `withdrawn-recovery-boundary`".
- DONE: Surface within +220 net vs merge-base; gofmt, go test ./..., go test ./... -race -p 2, and the live-tagged focused checks pass.
  `git diff --numstat origin/main...HEAD` (merge-base..861a6e13a): +275/-55 = +220 across 10 files, exactly at the ceiling. The first draft was +226; to fit, only comments were trimmed (two helper doc comments removed, test comments shortened). No logic or coverage was cut. `gofmt -l ./cmd ./internal` lists only `runtime_live_evidence_workflow_test.go`, which drifts at the merge-base too and is untouched. `-tags live`: TestWithdrawnRecoveryGradesThroughGap, TestBreakGlassSelectedBareGapBindsPerTarget, TestClaudeLiveRole (7 subtests), and TestDurableJourneyFailureConsultsGapBinding all pass. On 43547cabb, once the disk was freed (1.9Gi): `go test ./internal/release -count=1` ok (14.7s); `go test ./... -race -count=1 -p 2`: exit 0, 241s, 21 ok. On 861a6e13a (1.6Gi free): `go test ./... -count=1`: exit 0, 225s, 21 ok; `go test ./... -race -count=1 -p 2`: exit 0, 316s, 21 ok, no DATA RACE, no disk errors.

### Summary

Bound Haiku's five observed failures as blank-owner XFAIL. Replaced the Sonnet-only break-glass check with a per-target check that also proves blank-owner grading. Made the withdrawn-gate-recovery runner grade through its gap binding, with a focused check that is red under four mutations. Synced the registry doc. Net is +220, at the ceiling. Both Go suites pass. One deferred risk remains: contractlint `parseLiveGap` would reject a blank owner if it is ever wired to the real bindings. A runner mutation that bypasses `finishLiveScenario` without fataling is not covered by the stub-driver check; mutations (A) and (B) are the realistic forms, and both are covered.
