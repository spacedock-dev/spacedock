---
title: Expose complete gate history as machine-readable JSON
status: validation
source: "Captain authorized issue 814 on 2026-11-03: Proceed with the 3. You have the conn to approve gate and open pr"
issue: spacedock-dev/spacedock#814
sprint: paseo-integration
id: 8v4gxdyfenqxz3t9wpc2c474
gates:
    version: 1
    records:
        - id: gate:8v4gxdyfenqxz3t9wpc2c474:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:8v4gxdyfenqxz3t9wpc2c474-backlog-1
              briefing:
                id: briefing:8v4gxdyfenqxz3t9wpc2c474:backlog:attempt-1:revision-1
                digest: sha256:1eff4abee71d034d2bb6e6bb157e4b2ecc422e69e87742e9773adb551d07a908
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:8v4gxdyfenqxz3t9wpc2c474:backlog:1
                briefing: briefing:8v4gxdyfenqxz3t9wpc2c474:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-04T07:47:02.951714Z"
                decision: approve
                reason: The committed seed addresses issue 814 with complete stored-tree fidelity, read-only behavior, output stability, and focused proof requirements.
                conn:
                    quote: Proceed with the 3. You have the conn to approve gate and open pr
                    source: Captain message in this Codex conversation on 2026-11-03
              application:
                target-stage: ideation
                state: consumed
        - id: gate:8v4gxdyfenqxz3t9wpc2c474:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:8v4gxdyfenqxz3t9wpc2c474-ideation-1
              briefing:
                id: briefing:8v4gxdyfenqxz3t9wpc2c474:ideation:attempt-1:revision-1
                digest: sha256:1452421613a39eda31b11e8d881ab43e3bd41fd482189b2aba94a3b14a9cb7f4
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:8v4gxdyfenqxz3t9wpc2c474:ideation:1
                briefing: briefing:8v4gxdyfenqxz3t9wpc2c474:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-04T08:18:29.366102Z"
                decision: approve
                reason: Three design checklist items are DONE. Independent fixtures prove complete stored-tree fidelity, exact timestamps and reasons, read-only behavior, and stable status output. The design adds only a focused reader and help/docs corrections within its declared surface. Full implementation checks remain owed.
                conn:
                    quote: Proceed with the 3. You have the conn to approve gate and open pr
                    source: Captain message in this Codex conversation on 2026-11-03
              application:
                target-stage: implementation
                state: consumed
        - id: gate:8v4gxdyfenqxz3t9wpc2c474:validation
          stage: validation
          attempts:
            - id: gate-attempt:8v4gxdyfenqxz3t9wpc2c474-validation-1
              briefing:
                id: briefing:8v4gxdyfenqxz3t9wpc2c474:validation:attempt-1:revision-1
                digest: sha256:c501d93474abc9d6b2ef537a27a42a0f1874982f695c981cf5cf5b07a1bdfc0c
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:8v4gxdyfenqxz3t9wpc2c474:validation:1
                briefing: briefing:8v4gxdyfenqxz3t9wpc2c474:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-04T10:35:00.757185Z"
                decision: approve
                reason: The independent validation report's per-criterion evidence section maps every AC to external command, tree, filesystem, and help proof. Three checklist items are DONE. Exact-candidate normal and race aggregates each passed 21 test-bearing packages; independent checks and the detached audit passed, with five claim-breaking edits detected. The nine-file, net +403 candidate is within the approved surface. Approve PR delivery while terminal authority remains pending until delivery.
                conn:
                    quote: Proceed with the 3. You have the conn to approve gate and open pr
                    source: Captain message in this Codex conversation on 2026-11-03
              application:
                target-stage: done
                state: pending
started: 2026-11-04T07:47:25Z
worktree: .worktrees/spacedock-ensign-gate-history-json-reader
mod-block: merge:pr-merge
---
Paseo can show gate history from the binary-owned schema without parsing Markdown itself.

## Problem

Issue [814](https://github.com/spacedock-dev/spacedock/issues/814) needs a timeline that distinguishes spent approvals from the live approval.
`status --read --json` currently returns `frontmatter.gates` as `""`.
The consumer loses each attempt, briefing binding, resolution, reason, and application state.
Top-level help also advertises `gate validate`, which the dispatcher rejects.

The source baseline is `origin/main` commit `15e0baea39a6d47a802a01f45cdcfd8063c6edfd`.
The shared root checkout was not advanced or cleaned.
`internal/status/frontmatter.go` deliberately converts nested mappings and sequences to empty strings.
`section_read.go` and `json_commands.go` carry that scalar map into the existing all-strings JSON contract.

## Proposed approach

Add one focused reader:

```bash
spacedock gate read <ref-or-path> --json [--workflow-dir DIR]
spacedock gate read gate-history-json-reader --json --workflow-dir docs/dev
spacedock gate read archive:finished-task --json --workflow-dir docs/dev
```

The command returns the complete stored `gates` object directly, without an envelope or a timeline summary.
For example, an empty history returns `{"records":[],"version":1}` plus a newline.
The `version` value remains a JSON number.
Records, attempts, and scalar arrays retain their stored order.
JSON object keys use deterministic lexical order.
Decoded reason strings retain their exact Unicode, whitespace, quotes, tabs, and newlines.
The command does not escape `<`, `>`, or `&`.

`internal/status/gate_history.go` owns the reader and its JSON output.
It reuses `resolveRoots` and `resolveReadTarget`, which already support paths, active references, archive references, and split-root workflows.
It calls `gates.Read` and serializes the original returned node, rather than the typed `Document`.
This preserves optional empty fields, tolerated application data, and retired `provider-evidence` data that still exist in storage.
Those stored compatibility fields gain no authority from this reader.

Before decoding the private node into JSON values, the reader retains timestamp scalar text as strings.
The existing YAML decoder then handles the tree, and the standard JSON encoder writes one buffered document.
An entity without `gates` returns the same empty-history object as `gates: {version: 1, records: []}`.
Invalid gate data keeps the existing gate-read refusal and produces no partial JSON.
A stored value that cannot become JSON, such as an object with non-string keys, also fails without partial output.

`internal/cli/gate_read.go` owns the small reader grammar and calls the status reader.
It requires exactly one reference and one `--json`, and accepts one optional `--workflow-dir`.
It rejects mutation flags before resolving the target.
Success exits 0, data or resolution failures exit 1, and usage failures exit 2.
The existing gate command routes `read` to this helper before its mutation parser.

The help row replaces unsupported `validate` with `read`.
`gate validate` remains unsupported. Existing `status --validate` remains the validation command.

### Mechanisms and alternatives

| Mechanism | Value AC | Simplest alternative | Reason the alternative does not meet the value |
|---|---|---|---|
| Focused native JSON reader | AC-1, AC-3 | Change `frontmatter.gates` from a string to an object | That changes the existing status all-strings contract. A focused reader adds one command without changing that output. |
| Original gate node | AC-1, AC-2 | Marshal the existing typed `Document` | The spike lost explicit `reason: ''`, `includes: []`, application compatibility data, and retired provider evidence. |
| Retain timestamp scalar text | AC-2 | Decode the original YAML node directly into `any` | The spike rewrote `2026-09-16T00:00:00.1200Z` as `2026-09-16T00:00:00.12Z`. |
| Buffered standard JSON encoding | AC-1, AC-2 | Write JSON while traversing nodes | Buffering prevents partial stdout on an encoding failure. The standard encoder already provides deterministic object keys. |

No dependency, stored schema, event model, mutation path, authority rule, or standing enforcement process changes.

## Risk evidence

The debugging skill preceded the spike. The temporary checkout used the exact source baseline stated above.
A block-YAML fixture and an independent JSON fixture contain two records, five attempts, and 66 scalar leaves.
The validation record contains three approvals with states `superseded`, `superseded`, and `pending`, in that order.
The other record contains a withdrawn attempt and an open attempt.
The fixtures also contain request digests, connection citations, Unicode reasons, empty optional values, scalar arrays, and stored compatibility fields.

Observed results:

- The baseline real `status --read --json` command exited 0 and returned `frontmatter.gates: ""`.
- Baseline `gate read` and `gate validate` each exited 2 with the existing unknown-subcommand diagnostic.
- The typed-Document shortcut failed the complete-tree comparison because it omitted stored fields.
- Plain raw-node decoding failed that comparison because it normalized the unquoted timestamp.
- The temporary focused CLI passed the complete independent tree for slug, absolute path, and archive references.
- Two identical reads emitted identical bytes, with one final newline and readable Unicode reason text.
- Both absent and explicit empty histories returned the documented empty object.
- Invalid gate version exited 1 with empty stdout. Mutation flags and missing `--json` exited 2.
- Entity, workflow, Git index, Git HEAD, and all fixture-repository file hashes remained identical after the reader exercise.
- Existing status-read output remained identical before and after the temporary reader command.

The temporary CLI and fixtures live only in `/tmp/spacedock-gate-history-ideation-8v4gxdyf`.
The runnable probe below contains its independent baseline and the risky serialization mechanism, so that directory is not a required artifact.
The probe first reproduces both failed shortcuts, then exercises the retained-node approach.

## Out of scope

Changing gate decision authority, application semantics, stored schema, the Paseo plugin, or unrelated status projections.
No gate record is filtered, summarized, migrated, consumed, or rewritten by the new reader.

## Expected surface and tolerance

Estimate net LOC change: +434, across 9 files.
Planned insertions: 444. Planned deletions: 10.
Tolerance: net +314 to +554, across 7 to 11 files.
The state entity and temporary spike do not count toward the product estimate.

| Expected file | Insertions | Deletions | Responsibility |
|---|---:|---:|---|
| `internal/cli/cli.go` | 10 | 2 | Route `read`, update command usage and the unknown-command list. |
| `internal/cli/gate_read.go` | 65 | 0 | Parse the read-only grammar. |
| `internal/cli/help.go` | 2 | 2 | Replace the stale top-level help row. |
| `internal/status/gate_history.go` | 80 | 0 | Reuse read resolution, retain stored values, and emit JSON. |
| `internal/cli/gate_test.go` | 145 | 4 | Exercise history, empty cases, refusals, resolution, and exact gate help. |
| `internal/cli/cli_test.go` | 2 | 2 | Update the existing exact top-level help expectation. |
| `internal/status/testdata/gate-history/complete.md` | 87 | 0 | Independent stored history with distinct application states. |
| `internal/status/testdata/gate-history/complete.json` | 41 | 0 | Independent expected tree. |
| `docs/site/reference/command-reference.md` | 12 | 0 | Document the command and its output contract. |

Observable semantic changes are limited to the new `gate read` grammar, its native JSON output, and corrected gate help.
Stored formats, authority, mutable status behavior, runtime launch behavior, and existing status output remain unchanged.

## Acceptance criteria

**AC-1**: A machine-readable command exposes complete stored gate records, attempts, briefing data, resolutions, reasons, and application states.
Verified by: exact independent JSON fixtures through the real command. Dropping a nested field or confusing pending and superseded must fail.

**AC-2**: The reader preserves order and exact scalar values, including captain reasons, and handles empty gate history.
Verified by: command fixtures with repeated attempts, Unicode, and empty history. Reordering or rewriting a value must fail.

**AC-3**: Reading history changes no entity or workflow bytes, and existing status output remains stable outside the explicitly selected JSON surface.
Verified by: before/after file and Git checks plus existing output fixtures.

**AC-4**: Gate help names supported commands and documents the chosen reader with usable examples.
Verified by: command behavior and help fixtures, plus review of the documentation diff.

## Test plan

Add the focused failing command cases before the implementation.
Reuse the existing proof owners rather than creating another test suite.

| AC and primary proof owner | Observable proof | Falsifying edit | Cost and mode |
|---|---|---|---|
| AC-1 and AC-2: `internal/cli/gate_test.go` | Run the actual CLI against the independent Markdown and JSON fixtures. Compare the complete decoded tree. Require two records, five attempts, and the three distinct validation application positions. | Emit only the current attempt, omit an empty or compatibility field, change pending to superseded, or normalize the timestamp. | One fixture-driven Go table, deterministic, seconds. |
| AC-2: the same command table | Repeat the command, compare bytes, and compare absent and explicit empty histories with the fixed empty object. | Reorder records or attempts, trim a reason, change a scalar type, emit `null`, or add nondeterministic bytes. | Extend the primary table, deterministic. |
| AC-3: the same command fixture | Snapshot entity/workflow bytes, fixture tree, Git HEAD, index, and status before and after successful reads and refusals. Exercise slug, path, archive, and split-root resolution. | Create a lock file, call state sync, rewrite frontmatter, or resolve the definition root instead of the state root. | Existing fixture Git helpers, deterministic. |
| AC-3: `TestGoldenRead`, `TestJSONReadGolden`, and `section_read_test.go` | Run existing status output fixtures unchanged. | Change the status string-valued frontmatter contract or any existing status output. | Existing deterministic tests. |
| AC-4: `TestTopLevelHelpFormsAreIdentical` and `TestGatePrepareCLIPrintsExactRoomBindingAndCurrentV1HelpSurface` | Update the existing exact help fixtures. Run the documented reader and the unsupported `gate validate` case. | Restore stale `validate` help or advertise a reader spelling that the parser rejects. | Existing help owners plus command cases, deterministic. |

The command table also covers malformed gate YAML, invalid canonical bindings, scalar application data, and values that JSON cannot encode.
Data refusals must leave stdout empty and exit 1. Mutation flags, missing arguments, and missing `--json` must exit 2.
The expected tree comes from literal fixture data, never production constructors, field lists, or constants.

Required completion commands are `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.
Validation must apply the workflow audit trigger to the actual diff.
The planned command-routing change warrants a detached audit of the new command and its refusal/read-only proofs.
A claim-breaking edit must fail its primary owner, such as typed serialization that drops stored compatibility fields.
No host adapter, dispatch path, live lane, or shipped skill changes in this design.
Independent review, if required, belongs to the First Officer. The worker did not commission a reviewer.

## Product documentation diff

The existing command reference owns this command. No new documentation page is needed.
The following additions use the current command table and a short example after that table:

```diff
--- a/docs/site/reference/command-reference.md
+++ b/docs/site/reference/command-reference.md
@@ Workflow command table
+| `spacedock gate read <ref-or-path> --json [--workflow-dir DIR]` | Read complete stored gate history as JSON. References and paths use the same resolution as `status --read`, including archived entities. The command changes no files. |
@@ After the workflow command table
+`gate read --json` returns the stored `gates` object, including every attempt and its application state.
+Records and attempts retain their stored order. Scalar values retain their stored values.
+An entity without gate history returns `{"records":[],"version":1}`.
+Invalid gate data exits nonzero without JSON output.
+
+```bash
+spacedock gate read gate-history-json-reader --json --workflow-dir docs/dev
+spacedock gate read archive:finished-task --json --workflow-dir docs/dev
+```
```

The exact top-level help replacement is:

```diff
-  gate        prepare | withdraw | record | validate | consume
-                                      Record, inspect, or consume durable gate resolutions
+  gate        prepare | withdraw | record | read | consume
+                                      Read or update durable gate history
```

The gate-specific help adds this usage line and reader description:

```text
       spacedock gate read <ref-or-path> --json [--workflow-dir DIR]

Read complete gate history as one JSON document. The reader changes no files.
```

The command `Use`, short description, and unknown-subcommand diagnostic also list `read` among the supported commands.

## Runnable risk probe

Requirements: Git, Go, and the repository checkout. The source baseline is available as the commit stated earlier.
Copy the following Go block to `internal/status/ideation_gate_history_risk_test.go` in an isolated checkout of that commit.
Run `go test ./internal/status -run '^TestGateHistoryIdeationRisk$' -count=1 -v` in that checkout.
The probe is temporary evidence. Implementation replaces it with the existing CLI proof owner's fixture cases.

```go
package status

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"gopkg.in/yaml.v3"
)

const gateHistoryRiskFixture = `{"version":1,"records":[{"id":"gate:task:validation","stage":"validation","attempts":[{"id":"attempt:validation:1","briefing":{"id":"briefing:validation:1","digest":"sha256:1111111111111111111111111111111111111111111111111111111111111111","request-digest":"sha256:2222222222222222222222222222222222222222222222222222222222222222","room-ref":"@review/validation/briefing-1"},"provider-evidence":{"result-digest":"sha256:4444444444444444444444444444444444444444444444444444444444444444"},"resolution":{"type":"Resolution","id":"resolution:validation:1","briefing":"briefing:validation:1","by":"person:captain","at":"2026-09-16T00:00:00.1200Z","decision":"approve","reason":"","includes":[]},"application":{"target-stage":"done","state":"superseded","legacy":{"live":false,"counts":[1,2],"note":"  old  "}}},{"id":"attempt:validation:2","briefing":{"id":"briefing:validation:2","digest":"sha256:2222222222222222222222222222222222222222222222222222222222222222","room-ref":"@review/validation/briefing-2"},"resolution":{"type":"Resolution","id":"resolution:validation:2","briefing":"briefing:validation:2","by":"agent:first-officer","at":"2026-09-17T00:00:00Z","decision":"approve","reason":"  Résumé — 保留 <reason> & \"quote\"\nline\tend  ","conn":{"quote":"Proceed with the 3.","source":"Captain message"}},"application":{"target-stage":"done","state":"superseded"}},{"id":"attempt:validation:3","briefing":{"id":"briefing:validation:3","digest":"sha256:3333333333333333333333333333333333333333333333333333333333333333","room-ref":"@review/validation/briefing-3"},"resolution":{"type":"Resolution","id":"resolution:validation:3","briefing":"briefing:validation:3","by":"person:captain","at":"2026-09-18T00:00:00Z","decision":"approve","reason":"Final approval","includes":["annotation:first","annotation:last"]},"application":{"target-stage":"done","state":"pending"}}]},{"id":"gate:task:ideation","stage":"ideation","attempts":[{"id":"attempt:ideation:1","briefing":{"id":"briefing:ideation:1","digest":"sha256:4444444444444444444444444444444444444444444444444444444444444444","request-digest":"sha256:5555555555555555555555555555555555555555555555555555555555555555","room-ref":"@review/ideation/briefing-1"},"withdrawal":{"by":"agent:first-officer","at":"2026-09-15T00:00:00.123400Z","reason":"  旧資料  "}},{"id":"attempt:ideation:2","briefing":{"id":"briefing:ideation:2","digest":"sha256:5555555555555555555555555555555555555555555555555555555555555555","room-ref":"@review/ideation/briefing-2"}}]}]}`

func TestGateHistoryIdeationRisk(t *testing.T) {
	root := t.TempDir()
	readme := "---\nentity-type: task\nid-style: slug\nstages:\n  - name: validation\n    gate: true\n  - name: done\n    terminal: true\n---\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "task.md")
	stored := strings.Replace(gateHistoryRiskFixture, `"at":"2026-09-16T00:00:00.1200Z"`, `"at": 2026-09-16T00:00:00.1200Z`, 1)
	source := "---\nid: task\ntitle: Risk fixture\nstatus: validation\ngates: " + stored + "\n---\n## Body\n"
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	var want any
	if err := json.Unmarshal([]byte(gateHistoryRiskFixture), &want); err != nil {
		t.Fatal(err)
	}
	out, stderr, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--read", "task", "--json")
	if code != 0 || stderr != "" {
		t.Fatalf("baseline exit=%d stderr=%q", code, stderr)
	}
	var baseline struct {
		Frontmatter map[string]string `json:"frontmatter"`
	}
	if err := json.Unmarshal([]byte(out), &baseline); err != nil {
		t.Fatal(err)
	}
	if baseline.Frontmatter["gates"] != "" {
		t.Fatal("baseline no longer reproduces the loss")
	}
	doc, raw, err := gates.Read(path)
	if err != nil {
		t.Fatal(err)
	}
	typed, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var typedTree any
	if err := json.Unmarshal(typed, &typedTree); err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(typedTree, want) {
		t.Fatal("typed shortcut unexpectedly retained all stored fields")
	}
	var plain any
	if err := raw.Decode(&plain); err != nil {
		t.Fatal(err)
	}
	plainBytes, err := json.Marshal(plain)
	if err != nil {
		t.Fatal(err)
	}
	var plainTree any
	if err := json.Unmarshal(plainBytes, &plainTree); err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(plainTree, want) {
		t.Fatal("timestamp normalization risk no longer reproduces")
	}
	first, err := gateHistoryRiskJSON(path)
	if err != nil {
		t.Fatal(err)
	}
	second, err := gateHistoryRiskJSON(path)
	if err != nil {
		t.Fatal(err)
	}
	var got any
	if err := json.Unmarshal(first, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatal("the complete stored tree differs from the independent fixture")
	}
	if !bytes.Equal(first, second) || bytes.Count(first, []byte("\n")) != 1 || !bytes.HasSuffix(first, []byte("\n")) {
		t.Fatal("JSON bytes are unstable")
	}
	if !bytes.Contains(first, []byte("<reason> & ")) {
		t.Fatal("JSON escapes readable reason text")
	}
	after, err := os.ReadFile(path)
	if err != nil || string(after) != source {
		t.Fatal("the reader changed source bytes")
	}
	for _, seed := range []string{"gates: {version: 1, records: []}\n", ""} {
		if err := os.WriteFile(path, []byte("---\ntitle: Empty\n"+seed+"---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		empty, err := gateHistoryRiskJSON(path)
		if err != nil || string(empty) != "{\"records\":[],\"version\":1}\n" {
			t.Fatalf("empty=%q err=%v", empty, err)
		}
	}
	t.Log("baseline loses the tree; raw-node JSON retains 2 records, 5 attempts, all values, and empty history")
}

func gateHistoryRiskJSON(path string) ([]byte, error) {
	_, node, err := gates.Read(path)
	if errors.Is(err, gates.ErrNoGateRecord) {
		return []byte("{\"records\":[],\"version\":1}\n"), nil
	}
	if err != nil {
		return nil, err
	}
	var preserve func(*yaml.Node)
	seen := map[*yaml.Node]bool{}
	preserve = func(n *yaml.Node) {
		if n == nil || seen[n] {
			return
		}
		seen[n] = true
		if n.Kind == yaml.ScalarNode && n.Tag == "!!timestamp" {
			n.Tag = "!!str"
		}
		for _, child := range n.Content {
			preserve(child)
		}
		preserve(n.Alias)
	}
	preserve(node)
	var tree any
	if err := node.Decode(&tree); err != nil {
		return nil, err
	}
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(tree); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
```

### Feedback Cycles

## Stage Report: ideation

- DONE: Demonstrate the nested gate-history loss against independent fixtures and identify the smallest read-only API that returns the complete stored tree.
  AC-1/AC-2: the baseline loses the nested tree. The focused CLI retained two records, five attempts, and 66 scalar leaves. Typed serialization and timestamp normalization each failed the independent oracle. AC-3: file and Git snapshots remained identical after reads and refusals.
- DONE: Produce a concrete design, unchanged value ACs with falsifiers, counted file/net-LOC estimate and tolerance, semantic changes, proof owners, and product-documentation diff.
  AC-1 through AC-4 remain byte-identical. The design declares +434 net LOC (444 insertions, 10 deletions), nine files, net tolerance +314 to +554, and seven to eleven files.
- DONE: Commit a complete ideation Stage Report with the three checklist items and runnable risk evidence; send its exact path and commit to the First Officer.
  This entity contains the complete design, documentation diff, and runnable Go risk probe. The path-scoped state commit and exact report path are sent to the First Officer.

### Summary

The selected command is `gate read <ref-or-path> --json`, which returns the complete stored gate object through the existing read resolver.
The temporary CLI preserved stored compatibility fields, exact reasons, timestamp text, attempt order, and absent or empty history without file or Git changes.
The implementation adds no stored format or authority change and replaces stale `gate validate` help with `read`.

### Proof and observed limits

- The temporary command probe passed slug, absolute-path, archive, split-root, empty-history, refusal, repeat-byte, and file/Git snapshot cases. A status-read before/after comparison also passed. Calling sync or rewriting one stored reason fails those oracles.
- `TestGateHistoryIdeationRisk` passed in 0.584 seconds. It compares the complete tree with literal fixture data. Dropping an empty field, changing a pending state, or shortening the timestamp fails it.
- The existing help owners passed in 1.752 seconds with the stable fnm PATH. Restoring unsupported help text or changing a documented spelling fails their exact output expectations.
- The focused status golden/read owners and risk probe passed in 2.333 seconds with the same PATH. Changing existing status output fails the independent golden fixtures.
- `go test ./...` exited 1 in the scratch checkout. CLI and ensigncycle reached the package ten-minute timeout. The observed active tests were `TestClaudeInstallLeavesOnlySelectedChannelEnabled/sibling-present` and `TestRecordedGateCommittedBeforeDispatchResolutionPaths`. Timeout causes were not diagnosed.
- The same run failed two piruntime harnesses. They launched `/opt/homebrew/Cellar/node/25.9.0_2/bin/node`, whose loader could not find `/opt/homebrew/opt/llhttp/lib/libllhttp.9.3.dylib`.
- `go test ./... -race` was started, then stopped with exit 143 after the unrelated broad-check failures. It is incomplete, not a passing race result.
- `gofmt -w ./cmd ./internal` completed in the scratch checkout. The shared product checkout was not modified.
- The focused reruns used `/Users/clkao/.local/share/fnm/node-versions/v24.13.1/installation/bin` as a per-command PATH prefix. No global Node configuration changed.

### Blockers

No ideation blocker remains. Implementation and validation must run the required full checks and report any remaining baseline limits.
The First Officer owns any independent design review and the gate decision. No worker approved a gate or opened a PR.

## Stage Report: implementation

- DONE: Deliver the approved behavior and documentation with focused red-first evidence for every value acceptance criterion.
  The committed CLI fixtures prove AC-1 through AC-3. The exact help owners, removed-verb owner, and command reference prove AC-4.
- DONE: Commit the scoped candidate after formatting and the complete Go and race suites pass; report denominators and unchanged approved semantics.
  Product 2aeab0a03e5b54d9644304ca7325e913d4f7697b follows formatting and successful complete normal/race commands. Each final command passed all 21 test-bearing packages.
- DONE: Commit and send a complete implementation Stage Report with candidate identity, three checklist outcomes, AC evidence, surface actuals, and any finding proposals.
  This report records all three outcomes, the candidate, AC proofs, surface, and authorized dispositions. Its path-scoped state commit is sent to the First Officer.

### Summary

`gate read <ref-or-path> --json` now returns the complete validated stored gate object. The reader retains timestamp text, scalar values, array order, empty fields, and compatibility data without changing stored authority.

### Candidate and surface

- Product: `2aeab0a03e5b54d9644304ca7325e913d4f7697b`, branch `gate-history-json-reader`, approved source baseline `15e0baea39a6d47a802a01f45cdcfd8063c6edfd`. The product worktree is clean.
- Actual surface: nine files, +411/-8, net +403. Estimate: nine files, net +434. Actual net is 7.1% smaller and within +314 to +554 and seven to eleven files.
- Approved semantics remain unchanged: stored schema, decision authority, application semantics, runtime launch behavior, and existing status output. Observable changes are the new reader and corrected gate help.

### Evidence and checks

- Red first: `go test ./internal/cli -run '^(TestGateReadCLI|TestTopLevelHelpFormsAreIdentical|TestGatePrepareCLIPrintsExactRoomBindingAndCurrentV1HelpSurface)$' -count=1 -v` exited 1 before implementation. All 23 reader cases and both help owners failed on the missing command or stale help (4.943s).
- AC-1/AC-2: `TestGateReadCLI` passes 23/23 cases against independent Markdown/JSON literals: two records, five attempts, and 66 scalar leaves. Typed serialization, timestamp normalization, field loss, reason edits, or array reordering fail the complete-tree oracle. Repeat reads also require identical bytes, readable Unicode, unescaped HTML characters, and one final newline.
- AC-3: the same owner snapshots five complete-history forms, two empty forms, and 16 refusals. It compares every fixture file and directory, Git HEAD, index, status, and before/after status-read output. A writer, sync, or incorrect split-root resolver fails these checks.
- AC-3 compatibility: nine existing status owners pass unchanged, including `TestGoldenRead`, `TestJSONReadGolden`, and section-read cases. Changing string-valued frontmatter or existing output fails these golden/read assertions.
- AC-4: both exact help owners and all three removed-verb cases pass. Restoring `validate` help or omitting `read` from diagnostics fails them. The documented slug, path, archive, discovery, and relative-workflow forms run through the real CLI.
- Final focused command: `go test ./internal/cli ./internal/status -run '^(TestGateReadCLI|TestTopLevelHelpFormsAreIdentical|TestGatePrepareCLIPrintsExactRoomBindingAndCurrentV1HelpSurface|TestRemovedGateVerbsAreAbsentAndSideEffectFree|TestGoldenRead|TestJSONReadGolden|TestReadFrontmatterParity|TestReadHeadingMapLocatesBytes|TestReadFencedHeadingsSkipped|TestFindSectionSpansUsesFenceSafeHeadingOwnership|TestReadEntityRefResolvesLikeResolve|TestReadPathFormResolvesPlainFile|TestReadFlagIncompatibility)$' -count=1 -v`: 13 top-level tests and 47 named subcases passed across two packages.
- Final Go environment: `PATH=/Users/clkao/.local/share/fnm/node-versions/v24.13.1/installation/bin:$PATH`, `GOFLAGS=-p=1`, `GOCACHE=/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/paseo-fo-64ye7hx6/gocache`. Full suites ran sequentially; no global install or shared-cache change occurred.
- `go test ./...` exited 0: 21/21 test-bearing packages passed fresh, plus two packages without tests. No package timed out.
- `go test ./... -race -timeout 30m` exited 0: 21/21 test-bearing packages passed (17 cached, four fresh), plus two without tests. No race finding appeared. The First Officer authorized the invocation budget after a measured 520.323-second normal package; no assertion, AC, or committed timeout changed.
- `go test ./internal/ensigncycle -race -timeout 30m` passed 1/1 packages in 519.922s before the successful aggregate rerun. It retained every assertion and used the same cache/PATH/GOFLAGS.
- `gofmt -w ./cmd ./internal` exited 0. All six changed Go files also passed `gofmt -l`; `git diff --check` passed. The authorized exclusion of unrelated baseline formatting appears below.
- `uv run --with-requirements docs/requirements.txt mkdocs build --strict --site-dir /tmp/spacedock-gate-history-json-reader-site` exited 0. The pinned requirements declare the documentation dependencies.
- Complete logs remain under `/tmp/spacedock-gate-history-json-reader-`: `focused.log`, `focused-final.log`, `normal.log`, `normal-final.log`, `race.log`, `race-ensigncycle-retry.log`, and `race-final.log`. Failed commands remain failed in this record.

### Findings and authorized dispositions

- Discovery: user/workflow: a developer runs the no-flag reader. Harm: the new fixture could not prove valid discovery. Authority: `value-ac[AC-3]` requires valid resolver evidence. Trigger: the first post-code run passed 22/23 reader cases and reported no commissioned workflow. Worker proposal: Polish, task-owned, fix. The First Officer classified an owned Material evidence defect and authorized only local `commissioned-by` setup. The final 23/23 reader cases pass; shared fixtures and discovery semantics remain unchanged.
- Diagnostic: user/workflow: a developer runs removed gate verbs. Harm: a stale exact-output oracle rejected the approved supported-command list. Authority: `value-ac[AC-4]` requires supported command names. Trigger: three removed-verb cases expected `prepare|withdraw|record|consume`; actual diagnostics also named `read`. Worker/FO disposition: Material evidence defect, task-owned, fix only that expected list. Rejection and side-effect assertions remain; all three cases pass.
- Missing cache archives: user/workflow: a developer runs the complete normal suite. Harm: six late packages could not build. Authority: `value-ac[AC-3]` requires complete compatibility evidence. Trigger: Go reported missing shared-cache archives for `testing`, `os`, and other imports. That first normal command exited 1: 14/21 packages passed, seven failed, plus two without tests. Worker/FO disposition: Material evidence/environment defect, environment-owned, fix with the session cache. Removal cause remains unknown; the complete normal rerun passed.
- Disk capacity: user/workflow: a developer runs the complete race suite. Harm: Go could not retain the workflow-cycle test log. Authority: `value-ac[AC-3]` requires a successful complete compatibility command. Trigger: assertions printed PASS, then `no space left on device` failed the package at 533.658s. That aggregate exited 1: 20/21 packages passed, one failed, plus two without tests. Worker proposal: Material environment defect, environment-owned, hold. The First Officer authorized an isolated failed-package retry, then the identical complete aggregate. Both passed; no files were deleted and no assertion changed.
- Formatting: user/workflow: released workflows remain unchanged. Harm: none. Authority: none: unrelated whitespace cannot establish a value failure. Trigger: required global gofmt aligned two pre-existing fields in `internal/release/runtime_live_evidence_workflow_test.go`. Worker/FO disposition: Polish, outside task ownership, decline. Only those two incidental edits were restored; that file is absent from the candidate.
- Concurrent suites: user/workflow: a developer runs the required checks. Harm: no shared-path interference was observed. Authority: `value-ac[AC-3]` protects the known shared dispatch proof boundary. Trigger: two external workflow-cycle test processes overlapped a green normal package. Worker/FO disposition: Deferred risk, environment-owned, decline product changes and continue. No other session was stopped. An observed shared-path failure promotes this risk for new triage.

### Remaining limits

No unresolved product finding remains. Cache-removal cause and low disk headroom remain environment limits. Independent validation and the detached command audit remain with the First Officer. No product branch was pushed, and no PR, merge, release, frontmatter, or gate-state change occurred.

## Stage Report: validation

- DONE: Independently verify each acceptance criterion against the exact committed candidate and report explicit PASS/FAIL with reproducible external evidence.
  PASS. All four ACs pass the real CLI, independent literal trees, unchanged status fixtures, and before/after filesystem and Git observations.
- DONE: Attack the highest-risk semantics and perform any workflow-required detached adversarial audit in a throwaway checkout; report concrete falsifiers and outcomes.
  PASS. The detached baseline passed 23/23 cases; existing proof owners killed 5/5 claim-breaking edits. The external adjacent matrix passed 14/14 cases.
- DONE: Commit and send a complete validation Stage Report with the three checklist outcomes, tests and denominators, exact candidate SHA, authorized findings, and PASSED/REJECTED recommendation.
  PASS. This path-scoped report contains every required result. Its state commit and publication result accompany the report sent to the First Officer.

### Summary

Recommendation: **PASSED** for product `2aeab0a03e5b54d9644304ca7325e913d4f7697b` on `gate-history-json-reader`. The exact candidate remains clean and unchanged. Validation found no new outcome defect or evidence defect.

### Acceptance criteria evidence

- PASS AC-1: `TestGateReadCLI` compares independent complete Markdown/JSON fixtures through the actual CLI: two records, five attempts, and 66 scalar leaves. Typed serialization that loses empty and compatibility fields failed the detached `slug` owner.
- PASS AC-2: Complete-tree and repeat-byte assertions preserve stored order, exact Unicode reasons, empty fields, timestamps, and both empty-history forms. Removing timestamp retention failed `slug`; the external matrix also preserved reordered records/attempts, consumed approvals, NUL/Unicode/EOF, and exact signed/unsigned 64-bit integers.
- PASS AC-3: All 23 command cases retain fixture files/directories, HEAD, index, status, and the existing status-read output. A write of `<entity>.audit-lock` failed both `slug` and `invalid_version`; the external matrix independently compared complete file hashes and Git/status-read bytes for all 14 cases. Existing status golden and section-read owners also pass unchanged.
- PASS AC-4: Both exact help owners and all three removed-verb cases pass. The real command exercises documented slug/archive forms and rejects unsupported `gate validate`. The documentation diff matches the approved command, empty result, refusal, and read-only contracts; changed prose meets the workflow's sentence and modal rules.

### Checks and detached audit

- Independent focused command: `go test ./internal/cli ./internal/status -run '^(TestGateReadCLI|TestTopLevelHelpFormsAreIdentical|TestGatePrepareCLIPrintsExactRoomBindingAndCurrentV1HelpSurface|TestRemovedGateVerbsAreAbsentAndSideEffectFree|TestGoldenRead|TestJSONReadGolden|TestReadFrontmatterParity|TestReadHeadingMapLocatesBytes|TestReadFencedHeadingsSkipped|TestFindSectionSpansUsesFenceSafeHeadingOwnership|TestReadEntityRefResolvesLikeResolve|TestReadPathFormResolvesPlainFile|TestReadFlagIncompatibility)$' -count=1 -json`: exit 0, 13/13 top-level tests, 47/47 named subcases, 2/2 packages.
- Reused exact-candidate completion evidence: `go test ./...` and `go test ./... -race -timeout 30m` each passed 21/21 test-bearing packages, plus two without tests. Independent log inspection found no failed package. Full suites were not duplicated.
- Reused the producer's required `gofmt -w ./cmd ./internal` run and its authorized unrelated-formatting disposition. Independent `gofmt -l` on the entrypoint and all six changed Go files printed nothing; `git diff --check` passed. The producer's strict MkDocs build is green.
- Detached checkout: `/tmp/spacedock-gate-history-json-reader-validation-audit-20261104-2aeab0a`, detached at the exact product SHA. `go test ./internal/cli -run '^TestGateReadCLI$' -count=1 -json` passed 23/23 subcases before mutation. The checkout is clean after all mutations were restored.
- Falsifier 1: serialize the typed `Document` instead of the original node. `TestGateReadCLI/slug` failed its complete-tree comparison; explicit empty fields and stored compatibility data cannot disappear silently.
- Falsifier 2: omit `retainGateTimestampText`. `TestGateReadCLI/slug` failed its complete-tree comparison because timestamp spelling changed.
- Falsifier 3: write an adjacent lock file before reading. `TestGateReadCLI/slug` and `/invalid_version` failed their filesystem/Git snapshots, including the refusal path.
- Falsifier 4: bypass canonical `gates.Validate`. `/invalid_version` and `/invalid_binding` failed when the CLI returned exit 0 and JSON instead of exit 1 and empty stdout.
- Falsifier 5: write `{` before returning a JSON-encoding refusal. `/non-string_key` and `/non-JSON_number` failed their empty-stdout assertions. All five edits failed tests, not compilation.
- External matrix: 14/14 cases passed, comprising 11 successful history variants and three refusals. Sizes 65,535, 65,536, 65,537, and 1,048,577 reason bytes were preserved; the largest read took 0.118 seconds without truncation. Duplicate record identity and positive infinity exited 1 with empty stdout.
- Alias hypothesis: an anchor inside `gates` preserves the exact aliased timestamp. An anchor outside `gates` is already refused by canonical reading: both `gate read` and existing `status --validate` report `gates.version must be 1`. It does not establish a regression or a newly accepted format restriction.

### Scope, findings, and retained evidence

- Actual product surface versus approved baseline `15e0baea39a6d47a802a01f45cdcfd8063c6edfd`: nine files, +411/-8, net +403. This is 7.1% below the net +434 estimate and within +314 to +554 and seven to eleven files. No AC, stored authority, schema, mutation behavior, or runtime behavior changed.
- No new finding requires disposition. The implementation's recorded discovery, diagnostic, cache, disk, formatting, and concurrent-suite dispositions remain authorized as recorded. No material finding remains. No new deferred risk was established.
- All commands used the assigned session `GOCACHE`, `GOFLAGS=-p=1`, and verified fnm Node PATH. Last observed free space was 315 MiB; no write failed and no shared cache or peer file was cleaned. No host adapter, dispatch, skill, or live-lane trigger occurs in this diff.
- Focused evidence: `/tmp/spacedock-gate-history-json-reader-validation-focused.jsonl`. Detached mutation commands, exact patches, logs, and `summary.json` remain in `/tmp/spacedock-gate-history-json-reader-validation-audit-evidence`.
- Reproduce adjacent evidence with `python3 /tmp/spacedock-gate-history-json-reader-validation-audit-evidence/adjacent.py <exact-candidate-checkout> <candidate-binary> <new-fixture-directory>`; the script uses repository literal fixtures and Python's standard library. `adjacent-summary.json` retains exact output hashes, timings, and canonical alias-refusal evidence.
