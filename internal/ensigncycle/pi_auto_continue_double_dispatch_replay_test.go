package ensigncycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPiAutoContinueReplayDoubleDispatch is the correction-round replay for the
// Pi auto-continue single-root leg from GitHub Actions run 32749385148 (artifact
// runtime-live-e2e-pi-live, sha bb814f024). The FO dispatched the validation
// worker twice: the first completed via subagent_wait, then the FO hit a
// dispatch-build entity_path error, retried with the correct path, and
// re-dispatched. The assert's spawn counter correctly counted 2 (a real
// double-dispatch), but the old `spawns != 1` check rejected a legitimate
// retry. The fix tolerates `spawns` in {1, 2} when the rest of the lifecycle
// completes.
//
// This test reconstructs the full stream the assert sees (stdout + stderr +
// root session, matching piSharedLiveDriver.run) and feeds it through
// assertWorkerLifecycle with stage="validation". It grades GREEN after the
// spawns-tolerance fix. t.Skip keeps it hermetic when the artifact is absent.
func TestPiAutoContinueReplayDoubleDispatch(t *testing.T) {
	artifactDir := os.Getenv("SPACEDOCK_PI_AC_ARTIFACT_DIR")
	if artifactDir == "" {
		artifactDir = "/tmp/pi-live-10-art/live-artifacts/pi/pi-common/auto-continue-after-implementation--auto-continue/single-root"
	}
	stdout, err := os.ReadFile(filepath.Join(artifactDir, "pi-stdout.txt"))
	if err != nil {
		t.Skipf("artifact absent (pi-stdout.txt): %v", err)
	}
	stderr, err := os.ReadFile(filepath.Join(artifactDir, "pi-stderr.txt"))
	if err != nil {
		t.Skipf("artifact absent (pi-stderr.txt): %v", err)
	}
	matches, err := filepath.Glob(filepath.Join(artifactDir, "sessions", "*.jsonl"))
	if err != nil || len(matches) != 1 {
		t.Skipf("artifact absent (root session): %v (%d matches)", err, len(matches))
	}
	rootSession, err := os.ReadFile(matches[0])
	if err != nil {
		t.Skipf("artifact absent (root session read): %v", err)
	}
	// Reconstruct exactly as piSharedLiveDriver.run() does:
	// stream = stdout + "\n" + (stderr + "\n" + rootSession)
	stream := string(stdout) + "\n" + string(stderr) + "\n" + string(rootSession)
	report := autoContinueGatedEndState(false)
	if err := assertWorkerLifecycle(stream, report, "validation", "gate prepare", artifactDir); err != nil {
		t.Fatalf("double-dispatch replay graded RED after the spawns-tolerance fix: %v", err)
	}
}

// This hermetic replay retains native parent events and separate child evidence
// from CI 35058669297. provenance.json records the source bytes and projection.
func capturedPiCompletion(t *testing.T, name string) (stream, artifacts, childPath string) {
	t.Helper()
	stream = readFile(t, filepath.Join("testdata", "pi_native_completion", name+"-parent.jsonl"))
	child := readFile(t, filepath.Join("testdata", "pi_native_completion", name+"-child.jsonl"))
	var result, notice piSessionRecord
	lines := strings.Split(stream, "\n")
	if err := json.Unmarshal([]byte(lines[2]), &result); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(lines[3]), &notice); err != nil {
		t.Fatal(err)
	}
	owner := result.Message.Details.Mission.OwnerSessionID
	locator := strings.Split(notice.Content, "Session file: ")[1]
	rel, err := filepath.Rel(strings.TrimSuffix(owner, ".jsonl"), locator)
	if err != nil {
		t.Fatal(err)
	}
	artifacts = t.TempDir()
	writeFile(t, filepath.Join(artifacts, "sessions", filepath.Base(owner)), stream)
	childPath = filepath.Join(artifacts, "sessions", strings.TrimSuffix(filepath.Base(owner), ".jsonl"), rel)
	writeFile(t, childPath, child)
	return
}

func TestPiNativeCompletionCapturedReplay(t *testing.T) {
	for _, name := range []string{"default", "auto"} {
		t.Run(name, func(t *testing.T) {
			stream, artifacts, _ := capturedPiCompletion(t, name)
			stage, boundary := "implementation", "status=validation"
			if name == "auto" {
				stage, boundary = "validation", "gate prepare"
			}
			report := "## Stage Report: " + stage + "\n\n- DONE: work\n"
			if err := assertWorkerLifecycle(stream, report, stage, boundary, artifacts); err != nil {
				t.Fatal(err)
			}
			if err := assertWorkerLifecycle(stream, report, stage, boundary); err == nil {
				t.Fatal("notice without child evidence passed")
			}
			if err := assertWorkerLifecycle(stream, "", stage, boundary, artifacts); err == nil {
				t.Fatal("missing report passed")
			}
			if name == "auto" {
				// Reconstruct the captured absent-gate end state with its actual report.
				// No downloaded Git repository exists; commit independently in a test repo.
				root, entity := stageAutoContinueEndState(t, false, false)
				var reportRead piSessionRecord
				if err := json.Unmarshal([]byte(strings.Split(stream, "\n")[4]), &reportRead); err != nil {
					t.Fatal(err)
				}
				body := strings.Replace(autoContinueEntity(), "status: implementation", "status: validation", 1) + "\n" + piTextContent(reportRead.Message.Content)
				writeFile(t, entity, body)
				gitCommitPathScoped(t, root, filepath.Base(entity), "reconstruct captured absent gate")
				if err := assertAutoContinueDispatchEvidence(t, stream, root, entity, artifacts); err == nil || !strings.Contains(err.Error(), "gate") {
					t.Fatalf("captured missing gate must still fail: %v", err)
				} else {
					t.Logf("captured downstream failure preserved: %v", err)
				}
				// Git, not the completion notice or child's commit prose, owns durability.
				git(t, root, "checkout", "--orphan", "uncommitted-report")
				git(t, root, "rm", "--cached", "-r", ".")
				writeFile(t, filepath.Join(root, "seed"), "seed")
				gitCommitPathScoped(t, root, "seed", "seed without report")
				if err := assertAutoContinueDispatchEvidence(t, stream, root, entity, artifacts); err == nil || !strings.Contains(err.Error(), "no durable commit") {
					t.Fatalf("uncommitted report passed: %v", err)
				}
			}
		})
	}
}

func TestPiNativeCompletionAttributionAndOrdering(t *testing.T) {
	const run = "c7f1d5e7-ef00-4235-a860-790fda66500a"
	for _, name := range []string{"run", "agent", "epoch", "task", "cwd", "missing-child", "symlink", "external", "sibling", "traversal", "missing-parent", "parent-identity", "failed-spawn", "wrong-result", "missing-run", "no-stop", "error-stop", "late-child", "early-child", "before-spawn", "after-boundary", "duplicate", "ambiguous-identity", "ambiguous-locator", "no-user", "conflicting-epoch", "conflicting-identities"} {
		t.Run(name, func(t *testing.T) {
			stream, artifacts, path := capturedPiCompletion(t, "default")
			child := readFile(t, path)
			lines := strings.Split(strings.TrimSpace(stream), "\n")
			switch name {
			case "run":
				child = strings.ReplaceAll(child, run, "another-run")
			case "agent":
				child = strings.ReplaceAll(child, "subagent-worker-", "subagent-other-")
			case "epoch":
				child = strings.ReplaceAll(child, run+"-1", run+"-2")
			case "task":
				child = strings.ReplaceAll(child, "recorded-gate-task-implementation.md", "another-task-implementation.md")
			case "cwd":
				child = strings.ReplaceAll(child, "335368489/004", "335368489/005")
			case "no-user":
				child = strings.ReplaceAll(child, `"role":"user"`, `"role":"system"`)
			case "missing-child":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				outside := filepath.Join(t.TempDir(), "session.jsonl")
				writeFile(t, outside, child)
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(outside, path); err != nil {
					t.Fatal(err)
				}
			case "external":
				lines[3] = strings.ReplaceAll(lines[3], "/home/runner/work/spacedock/spacedock/live-artifacts", "/outside/live-artifacts")
			case "sibling":
				lines[3] = strings.ReplaceAll(lines[3], "01a0a8ac-6c68-77ef-8d90-279835984ba4/", "another-parent/")
			case "traversal":
				lines[3] = strings.ReplaceAll(lines[3], "/run-0/", "/run-0/../run-0/")
			case "missing-parent":
				parents, _ := filepath.Glob(filepath.Join(artifacts, "sessions", "*.jsonl"))
				if err := os.Remove(parents[0]); err != nil {
					t.Fatal(err)
				}
			case "parent-identity":
				lines[0] = strings.ReplaceAll(lines[0], "279835984ba4", "279835984baa")
			case "failed-spawn":
				lines[2] = strings.ReplaceAll(lines[2], `"isError":false`, `"isError":true`)
			case "wrong-result":
				lines[2] = strings.ReplaceAll(lines[2], `"toolCallId":"`, `"toolCallId":"other-`)
			case "missing-run":
				lines[2] = strings.ReplaceAll(lines[2], `"runId":"`+run+`"`, `"runId":""`)
			case "no-stop":
				child = strings.ReplaceAll(child, `"stopReason":"stop"`, `"stopReason":"toolUse"`)
			case "error-stop":
				child = strings.ReplaceAll(child, `"stopReason":"stop"`, `"stopReason":"error"`)
			case "late-child":
				child = strings.ReplaceAll(child, "05:26:52.921Z", "05:28:52.921Z")
			case "early-child":
				child = strings.ReplaceAll(child, "05:25:51.375Z", "05:24:51.375Z")
			case "before-spawn":
				lines = []string{lines[0], lines[3], lines[1], lines[2], lines[4]}
			case "after-boundary":
				lines[3], lines[4] = lines[4], lines[3]
			case "duplicate":
				lines = append(lines, lines[3])
			case "conflicting-epoch", "conflicting-identities":
				other := filepath.Join(filepath.Dir(filepath.Dir(path)), "run-1", "session.jsonl")
				otherChild := strings.ReplaceAll(child, run+"-1", run+"-2")
				if name == "conflicting-identities" {
					otherChild = strings.Split(child, "\n")[1] + "\n" + otherChild
				}
				writeFile(t, other, otherChild)
				lines = append(lines, strings.ReplaceAll(lines[3], "/run-0/", "/run-1/"))
			case "ambiguous-identity":
				child += strings.Split(child, "\n")[1] + "\n"
			case "ambiguous-locator":
				lines[3] = strings.ReplaceAll(lines[3], "Session file: ", "Session file: /other\\nSession file: ")
			}
			if name != "missing-child" && name != "symlink" {
				writeFile(t, path, child)
			}
			stream = strings.Join(lines, "\n")
			if err := assertImplementationWorkerLifecycle(stream, "## Stage Report: implementation\n- DONE: work\n", artifacts); err == nil {
				t.Fatal("invalid completion passed")
			}
		})
	}
}

func capturedPiTipCompletion(t *testing.T, name string) (stream, artifacts, childPath string) {
	t.Helper()
	parent := readFile(t, filepath.Join("testdata", "pi_native_completion", name+"-parent.jsonl"))
	child := readFile(t, filepath.Join("testdata", "pi_native_completion", name+"-child.jsonl"))
	offset, resultLine, notifyLine := 20, 46, -1
	if name == "route" {
		offset, resultLine, notifyLine = 19, 94, 97
	}
	lines := strings.Split(parent, "\n")
	var result piSessionRecord
	if err := json.Unmarshal([]byte(lines[resultLine]), &result); err != nil {
		t.Fatal(err)
	}
	owner := result.Message.Details.Mission.OwnerSessionID
	var locator string
	if notifyLine >= 0 {
		var notice piSessionRecord
		if err := json.Unmarshal([]byte(lines[notifyLine]), &notice); err != nil {
			t.Fatal(err)
		}
		locator = strings.Split(notice.Content, "Session file: ")[1]
	} else {
		var record struct {
			Message struct {
				Details struct {
					Results []struct {
						SessionFile string `json:"sessionFile"`
					} `json:"results"`
				}
			}
		}
		if err := json.Unmarshal([]byte(lines[resultLine]), &record); err != nil {
			t.Fatal(err)
		}
		locator = record.Message.Details.Results[0].SessionFile
	}
	rel, err := filepath.Rel(strings.TrimSuffix(owner, ".jsonl"), locator)
	if err != nil {
		t.Fatal(err)
	}
	artifacts = t.TempDir()
	writeFile(t, filepath.Join(artifacts, "sessions", filepath.Base(owner)), parent)
	childPath = filepath.Join(artifacts, "sessions", strings.TrimSuffix(filepath.Base(owner), ".jsonl"), rel)
	writeFile(t, childPath, child)
	return strings.Repeat("\n", offset) + parent, artifacts, childPath
}

func TestPiNativeSynchronousCapturedReplay(t *testing.T) {
	stream, artifacts, _ := capturedPiTipCompletion(t, "sync")
	report := "## Stage Report: implementation\n- DONE: retained report\n"
	if err := assertImplementationWorkerLifecycle(stream, report, artifacts); err != nil {
		t.Fatal(err)
	}
	routes, _, routeErr := piRejectionRoutes(stream, artifacts)
	if routeErr != nil || len(routes) != 2 || routes[0].index != 65 || routes[1].index != 66 {
		t.Fatalf("sync route boundary: %v %v", routes, routeErr)
	}
	completed, err := piNativeCompletion(stream, "implementation", artifacts)
	if err != nil || completed != 66 {
		t.Fatalf("completion=%d want 66 before boundary 73: %v", completed, err)
	}
}

// The synchronous boundary must have explicit native success, never missing
// fields defaulting to exit zero or completion prose standing in for a child.
func TestPiNativeSynchronousNegativeControls(t *testing.T) {
	for _, name := range []string{"wrong-call", "wrong-run", "wrong-agent", "wrong-owner", "missing-results", "multiple-results", "missing-exit", "nonzero-exit", "missing-success", "error-result", "missing-output", "missing-locator", "missing-child", "task", "cwd", "epoch", "error-stop", "late-stop", "result-before-spawn", "result-after-advance", "duplicate-result"} {
		t.Run(name, func(t *testing.T) {
			stream, artifacts, path := capturedPiTipCompletion(t, "sync")
			rows := strings.Split(stream, "\n")
			var result piSessionRecord
			if err := json.Unmarshal([]byte(rows[66]), &result); err != nil {
				t.Fatal(err)
			}
			child := readFile(t, path)
			switch name {
			case "wrong-call":
				result.Message.ToolCallID = "other"
			case "wrong-run":
				result.Message.Details.RunID = "other"
			case "wrong-agent":
				result.Message.Details.Results[0].Agent = "other"
			case "wrong-owner":
				result.Message.Details.Mission.OwnerSessionID = "/outside/other.jsonl"
			case "missing-results":
				result.Message.Details.Results = nil
			case "multiple-results":
				result.Message.Details.Results = append(result.Message.Details.Results, result.Message.Details.Results[0])
			case "missing-exit":
				result.Message.Details.Results[0].ExitCode = nil
			case "nonzero-exit":
				n := 1
				result.Message.Details.Results[0].ExitCode = &n
			case "missing-success":
				result.Message.IsError = nil
			case "error-result":
				b := true
				result.Message.IsError = &b
			case "missing-output":
				result.Message.Details.Results[0].OutputState = ""
			case "missing-locator":
				result.Message.Details.Results[0].SessionFile = ""
			case "missing-child":
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			case "task":
				child = strings.ReplaceAll(child, "recorded-gate-task-implementation.md", "other-task-implementation.md")
			case "cwd":
				child = strings.ReplaceAll(child, "2958864088/004", "2958864088/005")
			case "epoch":
				child = strings.ReplaceAll(child, "6b3a7071-a2eb-48ba-8235-e1d0352bad43-1", "6b3a7071-a2eb-48ba-8235-e1d0352bad43-2")
			case "error-stop":
				child = strings.ReplaceAll(child, `"stopReason":"stop"`, `"stopReason":"error"`)
			case "late-stop":
				child = strings.ReplaceAll(child, "22:07:58.294Z", "22:08:58.294Z")
			}
			raw, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			rows[66] = string(raw)
			switch name {
			case "result-before-spawn":
				rows[64], rows[66] = rows[66], rows[64]
			case "result-after-advance":
				rows = append(rows, rows[66])
				rows[66] = ""
			case "duplicate-result":
				rows = append(rows, rows[66])
			}
			if name != "missing-child" {
				writeFile(t, path, child)
			}
			stream = strings.Join(rows, "\n")
			if err := assertImplementationWorkerLifecycle(stream, "## Stage Report: implementation\n- DONE: report\n", artifacts); err == nil {
				t.Fatal("invalid sync completion accepted")
			}
			routes, _, err := piRejectionRoutes(stream, artifacts)
			if err == nil && assertSameStageWorkers(routes, false) == nil {
				t.Fatal("invalid sync evidence accepted by route observer")
			}
		})
	}
}
