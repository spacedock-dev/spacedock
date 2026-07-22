package ensigncycle

import (
	"encoding/json"
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
	return `{"type":"item.completed","item":{"type":"command_execution","command":"` + command + `"}}`
}

func codexRawCommand(command string) string {
	row := codexCommandItem{Type: "item.completed"}
	row.Item.Type = "command_execution"
	row.Item.Command = command
	encoded, err := json.Marshal(row)
	if err != nil {
		panic(err)
	}
	return string(encoded)
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

	// Positive regression: the archived PR #551 Codex lane records the successful
	// quoted launcher behind a pipeline using bash -lc's display quote-splice.
	// Keep the complete command here so the detector is pinned to the observed
	// transport shape that exited zero and created wire-the-thing as id 001.
	filedDisplayQuotedCapture := codexRawCommand(`/bin/bash -lc 'launcher="${SPACEDOCK_BIN:-spacedock}"
printf '"'%s\\n' '---' 'id:' 'title: Wire The Thing' 'status: backlog' '---' 'Connect the required components so the thing is fully wired.' | \""'$launcher" new wire-the-thing
"$launcher" status --read wire-the-thing --json'`)
	if err := assertCodexFilingViaNew(filedDisplayQuotedCapture, slug); err != nil {
		t.Fatalf("expected PR #551's exact /bin/bash -lc display-quoted filing to count as atomic: %v", err)
	}

	// Negative regression: the same display-quoted text inside a heredoc is
	// narration, not an executed simple command. The validator caught the matcher
	// losing separator provenance and treating this line as a pipeline invocation.
	heredocNarration := codexRawCommand(`/bin/bash -lc 'launcher="${SPACEDOCK_BIN:-spacedock}"
cat <<'EOF'
\""'$launcher" new wire-the-thing
EOF
'`)
	if err := assertCodexFilingViaNew(heredocNarration, slug); err == nil {
		t.Fatal("expected display-quoted launcher text inside heredoc narration to stay negative")
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
		"wrapped different variable":              `/bin/bash -lc 'launcher="${SPACEDOCK_BIN:-spacedock}"; printf body | \""'$EDITOR" new ` + slug + `'`,
		"wrapped narrated invocation":             `/bin/bash -lc 'launcher="${SPACEDOCK_BIN:-spacedock}"; echo \""'$launcher" new ` + slug + `'`,
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
