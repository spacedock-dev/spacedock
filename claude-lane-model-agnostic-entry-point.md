---
title: The Claude live lane selects a model by cadence, so a new model needs a workflow edit
status: ideation
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
started: 2026-10-07T18:15:05Z
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
