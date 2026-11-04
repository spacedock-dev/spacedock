---
title: Expose complete gate history as machine-readable JSON
status: ideation
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
                state: pending
started: 2026-11-04T07:47:25Z
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
