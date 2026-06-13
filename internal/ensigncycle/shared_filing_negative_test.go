package ensigncycle

import "testing"

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
