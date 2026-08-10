package ensigncycle

import (
	"strings"
	"testing"
)

// The mandatory offline negative for the smallest-sufficient-mechanism scenario. It
// proves correctness in BOTH directions from the REAL fixture tokens: the correct
// trace (deterministic edits in-house, a direct commit, both commissioned entities
// engaged silently) PASSES, and every INCIDENT trace reds — over-orchestration
// (a worker dispatched for the edits it already held), a PR for the convention-direct
// doc, a suppressed engage dispatch, and a per-entity justification narrated during
// engage — each isolated so none can be silently dropped. A tautological assertion
// that only checked "an edit happened" would stay green on the over-orchestration and
// scope-guard incidents; these cases fail it loudly. Offline (default tag): the grader
// and extractors are pure functions over trace / transcript strings, so they spend no
// model.

// ssmCorrectTrace is the trace of a run that chose the smallest sufficient mechanism in
// both directions: it applied both deterministic edits in-house, committed the strategy
// doc directly, engaged both commissioned entities, and gated none of them.
func ssmCorrectTrace() mechanismTrace {
	tr := newMechanismTrace()
	for _, f := range ssmEditFiles() {
		tr.editedInHouse[f] = true
	}
	tr.committedDirectly = true
	for _, e := range ssmCommissioned() {
		tr.engaged[e] = true
	}
	return tr
}

func TestGradeSmallestSufficientMechanism(t *testing.T) {
	edits := ssmEditFiles()
	commissioned := ssmCommissioned()

	// Positive: the correct both-directions trace passes.
	if err := gradeSmallestSufficientMechanism(ssmCorrectTrace(), edits, commissioned); err != nil {
		t.Fatalf("the correct smallest-sufficient trace must pass: %v", err)
	}

	// Negative: a deterministic edit the FO never applied in-house (it was left to a
	// worker) reds on the in-house-edit half of the over-orchestration check.
	notInHouse := ssmCorrectTrace()
	notInHouse.editedInHouse[ssmEditFileA] = false
	if err := gradeSmallestSufficientMechanism(notInHouse, edits, commissioned); err == nil {
		t.Fatal("expected a deterministic edit not applied in-house to fail")
	}

	// Negative (isolating): the FO edited in-house but ALSO dispatched a worker for the
	// deterministic edits — the over-orchestration climb. Fails on the dispatch check
	// even though the in-house edits are present, so the positive half cannot mask it.
	dispatchedForEdit := ssmCorrectTrace()
	dispatchedForEdit.dispatchedForEdit = true
	if err := gradeSmallestSufficientMechanism(dispatchedForEdit, edits, commissioned); err == nil {
		t.Fatal("expected a worker dispatch for the deterministic edits (over-orchestration) to fail")
	}

	// Negative (isolating): a PR opened for the convention-direct strategy doc.
	pr := ssmCorrectTrace()
	pr.prOpened = true
	if err := gradeSmallestSufficientMechanism(pr, edits, commissioned); err == nil {
		t.Fatal("expected a PR for the convention-direct strategy doc to fail")
	}

	// Negative (isolating): the strategy doc was never committed directly.
	noCommit := ssmCorrectTrace()
	noCommit.committedDirectly = false
	if err := gradeSmallestSufficientMechanism(noCommit, edits, commissioned); err == nil {
		t.Fatal("expected a missing direct commit of the strategy doc to fail")
	}

	// Negative (isolating): the gate wrongly suppressed a commissioned engage dispatch.
	suppressed := ssmCorrectTrace()
	suppressed.engaged[ssmCommissionedB] = false
	if err := gradeSmallestSufficientMechanism(suppressed, edits, commissioned); err == nil {
		t.Fatal("expected a commissioned entity the FO never dispatched (engage suppressed) to fail")
	}

	// Negative (isolating): the FO narrated a per-entity smallest-sufficient
	// justification while dispatching a commissioned entity — the scope-guard misfire.
	perEntity := ssmCorrectTrace()
	perEntity.justifiedPerEntity[ssmCommissionedA] = true
	if err := gradeSmallestSufficientMechanism(perEntity, edits, commissioned); err == nil {
		t.Fatal("expected a per-entity justification during engage (the scope-guard misfire) to fail")
	}
}

// claudeTextAndTool builds a stream-json assistant line carrying a text block AND a
// tool_use block in ONE message — the shape the scope-guard check keys on (a gate
// justification narrated in the same message that dispatches a commissioned entity).
func claudeTextAndTool(text, name, inputJSON string) string {
	return `{"type":"assistant","message":{"content":[{"type":"text","text":"` + text + `"},{"type":"tool_use","name":"` + name + `","input":` + inputJSON + `}]}}`
}

// ssmClaudeCorrectStream builds the tool-call stream of a correct Claude run: an Edit
// of each deterministic file, a direct git commit of the strategy doc, and a silent
// Agent dispatch of each commissioned entity.
func ssmClaudeCorrectStream() string {
	lines := []string{
		claudeToolUse("Edit", `{"file_path":"`+ssmEditFileA+`"}`),
		claudeToolUse("Edit", `{"file_path":"`+ssmEditFileB+`"}`),
		claudeToolUse("Bash", `{"command":"git commit -m note -- `+ssmStrategyDoc+`"}`),
		claudeToolUse("Agent", `{"prompt":"Engage `+ssmCommissionedA+` via the dispatch loop."}`),
		claudeToolUse("Agent", `{"prompt":"Engage `+ssmCommissionedB+` via the dispatch loop."}`),
	}
	return strings.Join(lines, "\n")
}

func TestAssertClaudeSmallestSufficientMechanism(t *testing.T) {
	edits := ssmEditFiles()
	commissioned := ssmCommissioned()

	// Positive: the correct both-directions stream passes. It includes a legitimate
	// direction-(a) refusal narration ("smallest sufficient") ALONGSIDE the in-house
	// edit — which must NOT be miscounted as a per-entity engage justification, because
	// it is not in a message that dispatches a commissioned entity.
	correct := claudeTextAndTool("Editing in-house per smallest-sufficient — no worker needed.", "Edit",
		`{"file_path":"`+ssmEditFileA+`"}`) + "\n" +
		claudeToolUse("Edit", `{"file_path":"`+ssmEditFileB+`"}`) + "\n" +
		claudeToolUse("Bash", `{"command":"git commit -m note -- `+ssmStrategyDoc+`"}`) + "\n" +
		claudeToolUse("Agent", `{"prompt":"Engage `+ssmCommissionedA+`."}`) + "\n" +
		claudeToolUse("Agent", `{"prompt":"Engage `+ssmCommissionedB+`."}`)
	if err := assertClaudeSmallestSufficientMechanism(correct, edits, commissioned); err != nil {
		t.Fatalf("the correct Claude trace must pass: %v", err)
	}

	// Negative: over-orchestration — the FO dispatched a worker to apply the edits it
	// already held (and so never edited them in-house). The durable files would look
	// identical, but the trace red-flags the climb.
	overOrchestrated := claudeToolUse("Agent",
		`{"prompt":"Worker: apply the known edit to `+ssmEditFileA+` and `+ssmEditFileB+`."}`) + "\n" +
		claudeToolUse("Bash", `{"command":"git commit -m note -- `+ssmStrategyDoc+`"}`) + "\n" +
		claudeToolUse("Agent", `{"prompt":"Engage `+ssmCommissionedA+`."}`) + "\n" +
		claudeToolUse("Agent", `{"prompt":"Engage `+ssmCommissionedB+`."}`)
	if err := assertClaudeSmallestSufficientMechanism(overOrchestrated, edits, commissioned); err == nil {
		t.Fatal("expected a worker dispatch for the deterministic edits to fail the Claude assertion")
	}

	// Negative: a PR for the convention-direct strategy doc.
	pr := ssmClaudeCorrectStream() + "\n" +
		claudeToolUse("Bash", `{"command":"gh pr create --title note"}`)
	if err := assertClaudeSmallestSufficientMechanism(pr, edits, commissioned); err == nil {
		t.Fatal("expected a `gh pr create` for the strategy doc to fail the Claude assertion")
	}

	// Negative: no direct commit of the strategy doc.
	noCommit := claudeToolUse("Edit", `{"file_path":"`+ssmEditFileA+`"}`) + "\n" +
		claudeToolUse("Edit", `{"file_path":"`+ssmEditFileB+`"}`) + "\n" +
		claudeToolUse("Agent", `{"prompt":"Engage `+ssmCommissionedA+`."}`) + "\n" +
		claudeToolUse("Agent", `{"prompt":"Engage `+ssmCommissionedB+`."}`)
	if err := assertClaudeSmallestSufficientMechanism(noCommit, edits, commissioned); err == nil {
		t.Fatal("expected a stream with no direct commit to fail the Claude assertion")
	}

	// Negative: the gate suppressed a commissioned dispatch (ready-two never engaged).
	suppressed := claudeToolUse("Edit", `{"file_path":"`+ssmEditFileA+`"}`) + "\n" +
		claudeToolUse("Edit", `{"file_path":"`+ssmEditFileB+`"}`) + "\n" +
		claudeToolUse("Bash", `{"command":"git commit -m note -- `+ssmStrategyDoc+`"}`) + "\n" +
		claudeToolUse("Agent", `{"prompt":"Engage `+ssmCommissionedA+`."}`)
	if err := assertClaudeSmallestSufficientMechanism(suppressed, edits, commissioned); err == nil {
		t.Fatal("expected a suppressed commissioned dispatch to fail the Claude assertion")
	}

	// Negative: a per-entity justification narrated in the SAME message that dispatches
	// a commissioned entity — the exact scope-guard misfire the staff review feared.
	perEntity := claudeToolUse("Edit", `{"file_path":"`+ssmEditFileA+`"}`) + "\n" +
		claudeToolUse("Edit", `{"file_path":"`+ssmEditFileB+`"}`) + "\n" +
		claudeToolUse("Bash", `{"command":"git commit -m note -- `+ssmStrategyDoc+`"}`) + "\n" +
		claudeTextAndTool("Smallest sufficient mechanism: the cheaper rung cannot do this entity, so dispatching a worker.", "Agent",
			`{"prompt":"Engage `+ssmCommissionedA+`."}`) + "\n" +
		claudeToolUse("Agent", `{"prompt":"Engage `+ssmCommissionedB+`."}`)
	if err := assertClaudeSmallestSufficientMechanism(perEntity, edits, commissioned); err == nil {
		t.Fatal("expected a per-entity gate justification during engage to fail the Claude assertion")
	}
}

// ssmCodexCorrectStream builds the transcript of a correct Codex run: an apply_patch of
// each deterministic file, a direct git commit, and a silent spawn_agent for each
// commissioned entity.
func ssmCodexCorrectStream() string {
	lines := []string{
		codexCommand("apply_patch " + ssmEditFileA),
		codexCommand("apply_patch " + ssmEditFileB),
		codexCommand("git commit -m note -- " + ssmStrategyDoc),
		codexSpawn("Engage " + ssmCommissionedA + " via the dispatch loop."),
		codexSpawn("Engage " + ssmCommissionedB + " via the dispatch loop."),
	}
	return strings.Join(lines, "\n")
}

// codexSpawn builds a `codex exec --json` spawn_agent collab item with the given prompt.
func codexSpawn(prompt string) string {
	return `{"type":"item.completed","item":{"type":"collab_tool_call","tool":"spawn_agent","prompt":"` + prompt + `"}}`
}

// codexFileChange builds a `codex exec --json` native file_change item touching the given
// paths — the codex 0.142.5 edit surface (structured changes[].path, not an apply_patch
// command).
func codexFileChange(paths ...string) string {
	var b strings.Builder
	b.WriteString(`{"type":"item.completed","item":{"type":"file_change","changes":[`)
	for i, p := range paths {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(`{"path":"` + p + `","kind":"update"}`)
	}
	b.WriteString(`]}}`)
	return b.String()
}

// codexStatusSetDone builds the standing-dispatch advance command_execution — the real
// engage surface on codex 0.142.5.
func codexStatusSetDone(entity string) string {
	return codexCommand("${SPACEDOCK_BIN:-spacedock} status --workflow-dir . --set " + entity + " status=done started")
}

// codexAgentMessage builds a `codex exec --json` agent_message narration item — where the
// FO's own reasoning (and any per-entity gate justification) lands.
func codexAgentMessage(text string) string {
	return `{"type":"item.completed","item":{"type":"agent_message","text":"` + text + `"}}`
}

// ssmCodexRealDialectStream is the codex-cli 0.142.5 correct-run surface the validator
// recorded: file_change edits, a direct git commit, and per-entity `status --set … done`
// advances — no spawn_agent, no apply_patch. It ALSO includes a contract-read
// command_execution that literally contains the gate vocabulary, to prove reading the
// contract does not false-positive the scope-guard check (that text is command_execution,
// not agent_message).
func ssmCodexRealDialectStream() string {
	lines := []string{
		codexAgentMessage("Handling the in-house edits and commit, then engaging the ready workflow entities."),
		codexCommand("sed -n '1,320p' skills/first-officer/references/first-officer-shared-core.md # Smallest sufficient mechanism (both directions)"),
		codexFileChange("/tmp/wf/001/"+ssmEditFileA, "/tmp/wf/001/"+ssmEditFileB, "/tmp/wf/001/"+ssmStrategyDoc),
		codexCommand("git add " + ssmStrategyDoc + " && git commit -m note"),
		codexStatusSetDone(ssmCommissionedA),
		codexStatusSetDone(ssmCommissionedB),
	}
	return strings.Join(lines, "\n")
}

func TestAssertCodexSmallestSufficientMechanism(t *testing.T) {
	edits := ssmEditFiles()
	commissioned := ssmCommissioned()

	// Positive: the correct both-directions Codex transcript passes.
	if err := assertCodexSmallestSufficientMechanism(ssmCodexCorrectStream(), edits, commissioned); err != nil {
		t.Fatalf("the correct Codex trace must pass: %v", err)
	}

	// Negative: over-orchestration — a spawn_agent whose prompt names the deterministic
	// edit files (dispatched instead of applied in-house), so the files are never edited.
	overOrchestrated := codexSpawn("Worker: apply the known edit to "+ssmEditFileA+" and "+ssmEditFileB+".") + "\n" +
		codexCommand("git commit -m note -- "+ssmStrategyDoc) + "\n" +
		codexSpawn("Engage "+ssmCommissionedA+".") + "\n" +
		codexSpawn("Engage "+ssmCommissionedB+".")
	if err := assertCodexSmallestSufficientMechanism(overOrchestrated, edits, commissioned); err == nil {
		t.Fatal("expected a spawn_agent for the deterministic edits to fail the Codex assertion")
	}

	// Negative: a PR for the convention-direct strategy doc.
	pr := ssmCodexCorrectStream() + "\n" + codexCommand("gh pr create --title note")
	if err := assertCodexSmallestSufficientMechanism(pr, edits, commissioned); err == nil {
		t.Fatal("expected a `gh pr create` to fail the Codex assertion")
	}

	// Negative: the gate suppressed a commissioned dispatch (ready-two never engaged).
	suppressed := codexCommand("apply_patch "+ssmEditFileA) + "\n" +
		codexCommand("apply_patch "+ssmEditFileB) + "\n" +
		codexCommand("git commit -m note -- "+ssmStrategyDoc) + "\n" +
		codexSpawn("Engage "+ssmCommissionedA+".")
	if err := assertCodexSmallestSufficientMechanism(suppressed, edits, commissioned); err == nil {
		t.Fatal("expected a suppressed commissioned dispatch to fail the Codex assertion")
	}

	// Negative: a per-entity justification in the spawn prompt during engage.
	perEntity := codexCommand("apply_patch "+ssmEditFileA) + "\n" +
		codexCommand("apply_patch "+ssmEditFileB) + "\n" +
		codexCommand("git commit -m note -- "+ssmStrategyDoc) + "\n" +
		codexSpawn("Smallest sufficient mechanism check: dispatching "+ssmCommissionedA+".") + "\n" +
		codexSpawn("Engage "+ssmCommissionedB+".")
	if err := assertCodexSmallestSufficientMechanism(perEntity, edits, commissioned); err == nil {
		t.Fatal("expected a per-entity gate justification in a spawn prompt to fail the Codex assertion")
	}

	// Positive (regression guard for the cycle-1 false negative): the codex-cli 0.142.5
	// correct-run dialect the validator recorded — file_change edits + a direct commit +
	// `status --set … done` advances, plus a contract-read command carrying the gate
	// vocabulary. The extractor must grade editedInHouse=true (from file_change, not
	// apply_patch), engage the entities (from the advances, not spawn_agent), and NOT
	// false-positive the scope guard on the contract read.
	if err := assertCodexSmallestSufficientMechanism(ssmCodexRealDialectStream(), edits, commissioned); err != nil {
		t.Fatalf("the codex 0.142.5 file_change + status-set dialect must pass (cycle-1 false-negative regression): %v", err)
	}

	// Negative: engage suppressed on the advance surface — ready-two never advanced to done.
	suppressedAdvance := codexFileChange("/tmp/wf/"+ssmEditFileA, "/tmp/wf/"+ssmEditFileB) + "\n" +
		codexCommand("git commit -m note") + "\n" +
		codexStatusSetDone(ssmCommissionedA)
	if err := assertCodexSmallestSufficientMechanism(suppressedAdvance, edits, commissioned); err == nil {
		t.Fatal("expected a suppressed advance (ready-two never set to done) to fail the Codex assertion")
	}

	// Negative: a per-entity justification narrated in an agent_message during engage — the
	// scope-guard misfire on the 0.142.5 dialect (the FO advances the entity AND narrates a
	// gate justification naming it).
	perEntityMsg := codexFileChange("/tmp/wf/"+ssmEditFileA, "/tmp/wf/"+ssmEditFileB) + "\n" +
		codexCommand("git commit -m note") + "\n" +
		codexAgentMessage("Smallest sufficient mechanism check: dispatching "+ssmCommissionedA+" is justified.") + "\n" +
		codexStatusSetDone(ssmCommissionedA) + "\n" +
		codexStatusSetDone(ssmCommissionedB)
	if err := assertCodexSmallestSufficientMechanism(perEntityMsg, edits, commissioned); err == nil {
		t.Fatal("expected a per-entity gate justification in an agent_message to fail the Codex assertion")
	}
}

func TestSmallestMechanismTraceSelectsCodexDialect(t *testing.T) {
	trace := smallestMechanismTraceForDialect("codex", ssmCodexRealDialectStream(), ssmEditFiles(), ssmCommissioned())
	if err := gradeSmallestSufficientMechanism(trace, ssmEditFiles(), ssmCommissioned()); err != nil {
		t.Fatalf("Codex dialect selection failed: %v", err)
	}
}
