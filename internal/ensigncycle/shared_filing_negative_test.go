package ensigncycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Offline positive + negative cases for the `filing` scenario assertions. They
// build synthetic host streams — a stream that filed via `new` (passes) and the
// SPECIFIC manual-flow streams the assertion guards against (`--next-id` + a
// hand-write, must go red) — so a tautological assertion that only checked "a new
// command appeared" would stay green on the manual flow and these cases fail it.
// Offline (default tag): the assertions are pure functions over the transcript.

// claudeToolUse builds a stream-json assistant line carrying one tool_use block.
func claudeToolUse(name, inputJSON string) string {
	return `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"` + name + `","input":` + inputJSON + `}]}}`
}

// codexCommand builds a `codex exec --json` command_execution item line.
func codexCommand(command string) string {
	if !strings.Contains(command, "\n") {
		var decoded string
		if json.Unmarshal([]byte(`"`+command+`"`), &decoded) == nil {
			command = decoded
		}
	}
	exitCode := 0
	var event codexCommandItem
	event.Type = "item.completed"
	event.Item.Type = "command_execution"
	event.Item.Command = command
	event.Item.Status = "completed"
	event.Item.ExitCode = &exitCode
	encoded, _ := json.Marshal(event)
	return string(encoded)
}

func codexCommandResult(command, output, status string, exitCode *int) string {
	var event codexCommandItem
	event.Type = "item.completed"
	event.Item.Type = "command_execution"
	event.Item.Command = command
	event.Item.AggregatedOutput = output
	event.Item.Status = status
	event.Item.ExitCode = exitCode
	encoded, _ := json.Marshal(event)
	return string(encoded)
}

const codexPR679PublicCommand = `{"type":"item.completed","item":{"id":"item_9","type":"command_execution","command":"/bin/bash -lc \"printf '%s\\\\n' '---' 'title: Wire The Thing' 'status: backlog' '---' '' 'Wire the thing so it is connected and ready for future work.' | \\\"\"'${SPACEDOCK_BIN:-spacedock}\" new wire-the-thing'","aggregated_output":"created: /tmp/TestLiveCommonFiling3005718106/002/wire-the-thing.md id=001\n","exit_code":0,"status":"completed"}}`

// codexRun32105482382PublicCommand is the exact public item.completed event from
// run 32105482382 (codex-exec.jsonl line 14, artifact 9313785785): the atomic
// create piped from printf, then — in the SAME bash -lc item, on the next line —
// a read-only `status --read` verification. This is the shape that graded
// "0 atomic creates" before the newline terminator fix: the create verb is
// immediately followed by a literal `\n`, not a space/tab/quote/end-of-string.
const codexRun32105482382PublicCommand = `{"type":"item.completed","item":{"id":"item_6","type":"command_execution","command":"/bin/bash -lc \"printf '%s\\\\n' '---' 'title: Wire The Thing' 'status: backlog' '---' '' 'Wire the thing so it is connected and ready for follow-on work.' | \\\"\"'${SPACEDOCK_BIN:-spacedock}\" new wire-the-thing\n\"${SPACEDOCK_BIN:-spacedock}\" status --read wire-the-thing --json'","aggregated_output":"created: /tmp/TestLiveCommonFiling3063241947/002/wire-the-thing.md id=001\n{\"command\":\"read\",\"path\":\"/tmp/TestLiveCommonFiling3063241947/002/wire-the-thing.md\",\"total_lines\":\"7\",\"frontmatter\":{\"id\":\"001\",\"status\":\"backlog\",\"title\":\"Wire The Thing\"},\"headings\":[]}\n","exit_code":0,"status":"completed"}}`

func TestCodexPR679ExactPublicCommandTransaction(t *testing.T) {
	const root = "/tmp/TestLiveCommonFiling3005718106/002"
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("retained fixture path already exists: %s", root)
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Dir(root)) })
	entityPath := filepath.Join(root, filingSlug+".md")
	writeFile(t, entityPath, "---\nid: 001\ntitle: Wire The Thing\nstatus: backlog\n---\n\nWire the thing for follow-up work.\n")
	if _, err := assertCodexPublicFilingTransaction(codexPR679PublicCommand, entityPath, filingSlug); err != nil {
		t.Fatalf("exact retained PR #679 public command rejected: %v", err)
	}
}

func TestAssertClaudeFilingViaNew(t *testing.T) {
	slug := filingSlug

	// Positive: the FO filed via `spacedock new <slug>` piping a body on stdin.
	filed := claudeToolUse("Bash", `{"command":"spacedock new `+slug+` --workflow-dir . <<'EOF'\n# Wire The Thing\nbody\nEOF"}`)
	if err := assertClaudeFilingViaNew(filed, slug); err != nil {
		t.Fatalf("expected a `new`-filed stream to pass: %v", err)
	}

	// Positive: the `--new` flag alias also counts.
	filedAlias := claudeToolUse("Bash", `{"command":"spacedock status --new `+slug+` --workflow-dir ."}`)
	if err := assertClaudeFilingViaNew(filedAlias, slug); err != nil {
		t.Fatalf("expected the `--new` alias to count as atomic filing: %v", err)
	}

	// Positive: the contract-blessed var-capture launcher idiom — capture the
	// resolved launcher once (`B=${SPACEDOCK_BIN:-spacedock}`) then invoke it as
	// `$B new <slug>`. The `$B new` line carries NEITHER `spacedock` nor
	// `SPACEDOCK_BIN` literally, so a regex keyed only on those tokens misses it.
	// This is correct FO behavior (the launcher invariant), so it must pass.
	filedVarCapture := claudeToolUse("Bash", `{"command":"B=${SPACEDOCK_BIN:-spacedock}\n$B new `+slug+` --workflow-dir ."}`)
	if err := assertClaudeFilingViaNew(filedVarCapture, slug); err != nil {
		t.Fatalf("expected the `$B new` var-capture idiom to count as atomic filing: %v", err)
	}

	// Negative: a `$VAR new <slug>` call whose var was NEVER captured from
	// `${SPACEDOCK_BIN:-spacedock}` is not a launcher filing — the var-capture path
	// must stay tied to a real launcher resolution, not any `$X new`.
	bareVar := claudeToolUse("Bash", `{"command":"$EDITOR new `+slug+`.md"}`)
	if err := assertClaudeFilingViaNew(bareVar, slug); err == nil {
		t.Fatal("expected a `$VAR new` call with no launcher capture to fail")
	}

	// Negative: no atomic filing at all — the FO only previewed the id and never
	// committed to a create path. Must fail on the missing-`new` half.
	previewOnly := claudeToolUse("Bash", `{"command":"spacedock status --next-id --workflow-dir ."}`)
	if err := assertClaudeFilingViaNew(previewOnly, slug); err == nil {
		t.Fatal("expected a stream with no `new` command to fail")
	}

	// Negative: the manual pair — `--next-id` preview THEN a `Write` of the entity
	// file. This is the drift-prone flow `new` replaces; it must fail even though
	// the durable file would look identical.
	manualPair := claudeToolUse("Bash", `{"command":"spacedock status --next-id --workflow-dir ."}`) + "\n" +
		claudeToolUse("Write", `{"file_path":"001-`+slug+`.md","content":"---\nid: 001\n---\n"}`)
	if err := assertClaudeFilingViaNew(manualPair, slug); err == nil {
		t.Fatal("expected the manual `--next-id` + `Write` pair to fail even with no `new` command")
	}

	// Negative: BOTH `new` AND the manual pair appear — a run that filed atomically
	// but ALSO hand-wrote must still fail on the pair check, so the positive half
	// cannot mask a manual write.
	newPlusManual := filed + "\n" + manualPair
	if err := assertClaudeFilingViaNew(newPlusManual, slug); err == nil {
		t.Fatal("expected `new` plus the manual `--next-id` + `Write` pair to fail on the pair check")
	}

	// A `--next-id` alone alongside `new` (no entity Write) is fine — previewing the
	// candidate is not the manual flow without the hand-write that pairs with it.
	newWithPreview := filed + "\n" + claudeToolUse("Bash", `{"command":"spacedock status --next-id --workflow-dir ."}`)
	if err := assertClaudeFilingViaNew(newWithPreview, slug); err != nil {
		t.Fatalf("expected `new` plus a bare `--next-id` preview (no entity Write) to pass: %v", err)
	}
}

func TestAssertCodexFilingViaNew(t *testing.T) {
	slug := filingSlug

	artifactCommand := `/bin/bash -lc 'set -eu
sd_bin="${SPACEDOCK_BIN:-spacedock}"
printf '"'%s\\n' '---' 'title: Wire The Thing' 'status: backlog' '---' '' 'Wire the thing so it is connected and ready for follow-up work.' | \""'$sd_bin" new wire-the-thing --workflow-dir /tmp/TestLiveCommonFiling1757938389/002
"$sd_bin" status --boot --identify --json
"$sd_bin" status --read wire-the-thing --json'`
	if err := assertCodexFilingViaNew(codexCommand(artifactCommand), slug); err != nil {
		t.Fatalf("expected artifact 9078284956's successful atomic filing to count: %v", err)
	}
	failed := strings.Replace(codexCommand(artifactCommand), `"exit_code":0`, `"exit_code":1`, 1)
	if err := assertCodexFilingViaNew(failed, slug); err == nil {
		t.Fatal("expected a failed atomic filing command to fail")
	}
	redControls := map[string]string{
		"manual file creation": codexCommand("printf body > " + slug + ".md"),
		"preview plus write":   codexCommand("spacedock status --next-id; printf body > " + slug + ".md"),
		"wrong slug":           codexCommand("spacedock new other-slug --workflow-dir /tmp/wire-the-thing"),
		"narration only":       `{"type":"item.completed","item":{"type":"agent_message","text":"spacedock new wire-the-thing"}}`,
	}
	for name, stream := range redControls {
		t.Run(name, func(t *testing.T) {
			if assertCodexFilingViaNew(stream, slug) == nil {
				t.Fatalf("%s counted as atomic filing", name)
			}
		})
	}
	// Positive: the FO filed via a `spacedock new <slug>` command_execution.
	filed := codexCommand("spacedock new " + slug + " --workflow-dir .")
	if err := assertCodexFilingViaNew(filed, slug); err != nil {
		t.Fatalf("expected a `new`-filed Codex stream to pass: %v", err)
	}

	// Positive: the contract-blessed var-capture launcher idiom on Codex — the
	// capture and the `$B new <slug>` call land in one command_execution split by
	// a newline, so the call segment carries no literal `spacedock`/`SPACEDOCK_BIN`.
	// Correct FO behavior; must pass.
	filedVarCapture := codexCommand("B=${SPACEDOCK_BIN:-spacedock}\\n$B new " + slug + " --workflow-dir .")
	if err := assertCodexFilingViaNew(filedVarCapture, slug); err != nil {
		t.Fatalf("expected the `$B new` var-capture idiom to count as atomic filing on Codex: %v", err)
	}

	// Positive regression: PR #496's Codex command quoted both the captured
	// expansion and its invocation. It succeeded and created the entity atomically;
	// quote placement must not erase that durable action from the detector.
	filedQuotedCapture := codexCommand("launcher=\\\"${SPACEDOCK_BIN:-spacedock}\\\"\\nprintf '%s\\\\n' stub | \\\"$launcher\\\" new " + slug)
	if err := assertCodexFilingViaNew(filedQuotedCapture, slug); err != nil {
		t.Fatalf("expected PR #496's quoted launcher-capture filing to count as atomic: %v", err)
	}
	filedSingleQuotedCapture := codexCommand("launcher='${SPACEDOCK_BIN:-spacedock}'\\n\\\"$launcher\\\" new " + slug)
	if err := assertCodexFilingViaNew(filedSingleQuotedCapture, slug); err != nil {
		t.Fatalf("expected balanced single-quoted launcher capture to count as atomic: %v", err)
	}

	// Positive regression: PR #496's superseded f49b7826 Codex lane resolved the
	// fallback through command -v before invoking the same captured launcher.
	filedCommandVCapture := codexCommand("launcher=\\\"${SPACEDOCK_BIN:-$(command -v spacedock)}\\\"\\n\\\"$launcher\\\" new " + slug + " <<EOF")
	if err := assertCodexFilingViaNew(filedCommandVCapture, slug); err != nil {
		t.Fatalf("expected the exact command-v resolved launcher filing to count as atomic: %v", err)
	}

	malformedCaptures := map[string]string{
		"mismatched quotes":   `launcher=\"${SPACEDOCK_BIN:-spacedock}'\n\"$launcher\" new ` + slug,
		"leading-only quote":  `launcher=\"${SPACEDOCK_BIN:-spacedock}\n\"$launcher\" new ` + slug,
		"trailing-only quote": `launcher=${SPACEDOCK_BIN:-spacedock}\"\n\"$launcher\" new ` + slug,
	}
	for name, command := range malformedCaptures {
		t.Run(name, func(t *testing.T) {
			if err := assertCodexFilingViaNew(codexCommand(command), slug); err == nil {
				t.Fatalf("malformed launcher assignment counted as atomic filing: %s", command)
			}
		})
	}

	crossSegmentCalls := map[string]string{
		"different command after launcher status": `launcher=\"${SPACEDOCK_BIN:-spacedock}\"\n\"$launcher\" status; $EDITOR new ` + slug,
		"mismatched invocation quotes":            `launcher=\"${SPACEDOCK_BIN:-spacedock}\"\n\"$launcher' new ` + slug,
		"new token after launcher version":        `launcher=\"${SPACEDOCK_BIN:-spacedock}\"\n$launcher --version; touch new ` + slug,
	}
	for name, command := range crossSegmentCalls {
		t.Run(name, func(t *testing.T) {
			if err := assertCodexFilingViaNew(codexCommand(command), slug); err == nil {
				t.Fatalf("unrelated simple command counted as captured-launcher filing: %s", command)
			}
		})
	}

	// Negative: no atomic filing — must fail on the missing-`new` half.
	none := codexCommand("spacedock status --workflow-dir .")
	if err := assertCodexFilingViaNew(none, slug); err == nil {
		t.Fatal("expected a Codex stream with no `new` command to fail")
	}

	// Negative: the manual flow's id source — a `--next-id` command — appears. On
	// Codex (no Write tool) the `--next-id` command itself is the discriminator;
	// `new` needs none. Must fail even if `new` was also run.
	newPlusNextID := filed + "\n" + codexCommand("spacedock status --next-id --workflow-dir .")
	if err := assertCodexFilingViaNew(newPlusNextID, slug); err == nil {
		t.Fatal("expected a `--next-id` filing command on Codex to fail even alongside `new`")
	}

	// Negative: only the manual `--next-id` preview, no `new` — fails on both halves
	// (caught by the missing-`new` check first).
	nextIDOnly := codexCommand("spacedock status --next-id --workflow-dir .")
	if err := assertCodexFilingViaNew(nextIDOnly, slug); err == nil {
		t.Fatal("expected a `--next-id`-only Codex stream to fail")
	}
}

func TestCodexPublicFilingTransactionSevenRungMatrix(t *testing.T) {
	root := t.TempDir()
	entityPath := filepath.Join(root, filingSlug+".md")
	validEntity := "---\nid: 001\ntitle: Wire The Thing\nstatus: backlog\n---\n\nWire the thing for follow-up work.\n"
	receipt := "created: " + entityPath + " id=001\n"
	exit0, exit1 := 0, 1
	valid := codexCommandResult(`printf body | "${SPACEDOCK_BIN:-spacedock}" new wire-the-thing`, receipt, "completed", &exit0)
	unrelated := codexCommandResult("echo unrelated", receipt, "completed", &exit0)
	printedEvent := codexCommandResult("printf fake", valid+"\n"+receipt, "completed", &exit0)
	pr679Local := strings.ReplaceAll(codexPR679PublicCommand, "/tmp/TestLiveCommonFiling3005718106/002/wire-the-thing.md", entityPath)
	run32105482382Local := strings.ReplaceAll(codexRun32105482382PublicCommand, "/tmp/TestLiveCommonFiling3063241947/002/wire-the-thing.md", entityPath)
	run32105482382MalformedVerb := strings.Replace(run32105482382Local, `" new `, `" mew `, 1)

	tests := []struct {
		name, stream, entity string
		missing, wantErr     bool
	}{
		{"rung 1 direct alias", valid, validEntity, false, false},
		{"rung 1 PATH alias", codexCommandResult("spacedock new wire-the-thing", receipt, "completed", &exit0), validEntity, false, false},
		{"rung 1 bound alias", codexCommandResult("B=${SPACEDOCK_BIN:-spacedock}\n$B --new wire-the-thing", receipt, "completed", &exit0), validEntity, false, false},
		{"rung 2 exact PR679 public bytes", pr679Local, validEntity, false, false},
		{"rung 2 changed PR679 quote seam", strings.Replace(pr679Local, "SPACEDOCK_BIN:-spacedock", "OTHER_BIN:-other", 1), validEntity, false, true},
		{"rung 2 changed PR679 event type", strings.Replace(pr679Local, "item.completed", "item.started", 1), validEntity, false, true},
		{"rung 2 changed PR679 status", strings.Replace(pr679Local, `"status":"completed"`, `"status":"failed"`, 1), validEntity, false, true},
		{"rung 2 changed PR679 exit", strings.Replace(pr679Local, `"exit_code":0`, `"exit_code":1`, 1), validEntity, false, true},
		{"rung 2 exact run-32105482382 public bytes (create + newline + verify)", run32105482382Local, validEntity, false, false},
		{"rung 2 run-32105482382 malformed create verb", run32105482382MalformedVerb, validEntity, false, true},
		{"rung 3 unreachable native text and unrelated success", `{"type":"item.completed","item":{"type":"agent_message","text":"tools.exec_command({cmd:\"spacedock new wire-the-thing\"})"}}` + "\n" + unrelated, validEntity, false, true},
		{"rung 4 duplicate in one item", codexCommandResult("spacedock new wire-the-thing; spacedock --new wire-the-thing", receipt, "completed", &exit0), validEntity, false, true},
		{"rung 4 duplicate across items", valid + "\n" + valid, validEntity, false, true},
		{"rung 4 mixed aliases", valid + "\n" + codexCommandResult("spacedock --new wire-the-thing", receipt, "completed", &exit0), validEntity, false, true},
		{"rung 5 failed next-id", codexCommandResult("spacedock status --next-id", "", "failed", &exit1) + "\n" + valid, validEntity, false, true},
		{"rung 5 same-item next-id", codexCommandResult("spacedock status --next-id; spacedock new wire-the-thing", receipt, "completed", &exit0), validEntity, false, true},
		{"rung 6 failed create", codexCommandResult("spacedock new wire-the-thing", receipt, "failed", &exit1), validEntity, false, true},
		{"rung 6 wrong slug", codexCommandResult("spacedock new other-slug", receipt, "completed", &exit0), validEntity, false, true},
		{"rung 6 manual create", codexCommandResult("printf body > wire-the-thing.md", receipt, "completed", &exit0), validEntity, false, true},
		{"rung 6 started only", strings.Replace(valid, "item.completed", "item.started", 1), validEntity, false, true},
		{"rung 6 receipt from wrong item", codexCommandResult("spacedock new wire-the-thing", "", "completed", &exit0) + "\n" + unrelated, validEntity, false, true},
		{"rung 6 missing receipt", codexCommandResult("spacedock new wire-the-thing", "", "completed", &exit0), validEntity, false, true},
		{"rung 6 duplicate receipt", codexCommandResult("spacedock new wire-the-thing", receipt+receipt, "completed", &exit0), validEntity, false, true},
		{"rung 6 wrong receipt path", codexCommandResult("spacedock new wire-the-thing", "created: "+filepath.Join(root, "other.md")+" id=001\n", "completed", &exit0), validEntity, false, true},
		{"rung 6 wrong receipt id", codexCommandResult("spacedock new wire-the-thing", "created: "+entityPath+" id=002\n", "completed", &exit0), validEntity, false, true},
		{"rung 6 missing entity", valid, validEntity, true, true},
		{"rung 6 wrong entity id", valid, strings.Replace(validEntity, "id: 001", "id: 002", 1), false, true},
		{"rung 6 wrong entity body", valid, strings.Replace(validEntity, "Wire the thing for follow-up work.", "one\ntwo", 1), false, true},
		{"rung 7 printed event is nested output", printedEvent, validEntity, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.missing {
				_ = os.Remove(entityPath)
			} else {
				writeFile(t, entityPath, tt.entity)
			}
			_, err := assertCodexPublicFilingTransaction(tt.stream, entityPath, filingSlug)
			if (err != nil) != tt.wantErr {
				t.Fatalf("filing transaction error = %v, want error %v", err, tt.wantErr)
			}
		})
	}
}
