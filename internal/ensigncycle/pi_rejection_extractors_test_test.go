package ensigncycle

import (
	"fmt"
	"strings"
	"testing"
)

// piToolCallLine builds a Pi session JSONL assistant record carrying one toolCall block.
func piToolCallLine(id, name string, args string) string {
	return fmt.Sprintf(`{"type":"message","message":{"role":"assistant","content":[{"type":"toolCall","id":%q,"name":%q,"arguments":%s}]}}`,
		id, name, args)
}

// piBashToolCallLine builds a Pi session JSONL assistant record carrying a bash toolCall.
func piBashToolCallLine(id, command string) string {
	return piToolCallLine(id, "bash", fmt.Sprintf(`{"command":%q}`, command))
}

// piSpawnToolCallLine builds a Pi session JSONL assistant record carrying a subagent
// spawn (async dispatch with a task pointing at the dispatch file).
func piSpawnToolCallLine(id, handle string) string {
	task := fmt.Sprintf("Read /tmp/spacedock-dispatch/%s.md and treat its content as your assignment.", handle)
	return piToolCallLine(id, "subagent", fmt.Sprintf(`{"agent":"worker","async":true,"task":%q}`, task))
}

// piStatusToolCallLine builds a Pi session JSONL assistant record carrying a
// subagent status poll.
func piStatusToolCallLine(id, runID string) string {
	return piToolCallLine(id, "subagent", fmt.Sprintf(`{"action":"status","id":%q}`, runID))
}

// piToolResultLine builds a Pi session JSONL toolResult record.
func piToolResultLine(toolCallID, toolName, text string, isError bool) string {
	return fmt.Sprintf(`{"type":"message","message":{"role":"toolResult","toolCallId":%q,"toolName":%q,"content":[{"type":"text","text":%q}],"isError":%t}}`,
		toolCallID, toolName, text, isError)
}

const piRoundCommand = `${SPACEDOCK_BIN:-spacedock} gate record rejection-task --workflow-dir . --round validation/1 --briefing rejection-task/inputs/briefing.json --log rejection-task/inputs/briefing.review.jsonl`
const piRoundResult = "round=round:rejection-task:validation:1 stage=validation cycle=1 briefing=briefing:rejection-task:validation:round-1 entries=4"

const (
	piImplHandle = "spacedock-ensign-rejection-task-implementation"
	piValHandle  = "spacedock-ensign-rejection-task-validation"
)

// TestPiRecordedRejectionRound pins the Pi-dialect round extractor: it finds a
// successful `gate record --round validation/1` invocation and its correlated
// toolResult, and reports true. Falsifiable: a missing call, a failed result, or a
// result text for a different round all report false.
func TestPiRecordedRejectionRound(t *testing.T) {
	session := strings.Join([]string{
		piBashToolCallLine("call-1", piRoundCommand),
		piToolResultLine("call-1", "bash", piRoundResult, false),
	}, "\n")
	if !piRecordedRejectionRound(session) {
		t.Fatal("Pi extractor missed correlated round invocation with successful result")
	}

	// Falsifier: a missing call.
	if piRecordedRejectionRound(piBashToolCallLine("call-1", piRoundCommand)) {
		t.Fatal("Pi extractor accepted a round invocation with no correlated result")
	}

	// Falsifier: a failed result.
	failed := strings.Join([]string{
		piBashToolCallLine("call-1", piRoundCommand),
		piToolResultLine("call-1", "bash", "Error: exit 1", true),
	}, "\n")
	if piRecordedRejectionRound(failed) {
		t.Fatal("Pi extractor accepted a failed round invocation")
	}

	// Falsifier: a result text for a different round.
	wrongRound := strings.Join([]string{
		piBashToolCallLine("call-1", piRoundCommand),
		piToolResultLine("call-1", "bash", strings.Replace(piRoundResult, "cycle=1", "cycle=2", 1), false),
	}, "\n")
	if piRecordedRejectionRound(wrongRound) {
		t.Fatal("Pi extractor accepted a success line for a different round")
	}
}

// TestPiRejectionRoundPublications pins the Pi-dialect publication counter: exactly one
// publication at validation/1. Falsifiable: a republication at validation/2 and a
// failed call both under-report correctly.
func TestPiRejectionRoundPublications(t *testing.T) {
	session := strings.Join([]string{
		piBashToolCallLine("call-1", piRoundCommand),
		piToolResultLine("call-1", "bash", piRoundResult, false),
	}, "\n")
	if err := assertSingleRejectionRoundPublication(piRejectionRoundPublications(session)); err != nil {
		t.Fatalf("Pi publication counter rejected the single publication: %v", err)
	}

	// Falsifier: a republication at validation/2.
	republished := strings.Join([]string{
		piBashToolCallLine("call-1", piRoundCommand),
		piToolResultLine("call-1", "bash", piRoundResult, false),
		piBashToolCallLine("call-2", strings.Replace(piRoundCommand, "validation/1", "validation/2", 1)),
		piToolResultLine("call-2", "bash", "round=validation/2", false),
	}, "\n")
	if err := assertSingleRejectionRoundPublication(piRejectionRoundPublications(republished)); err == nil {
		t.Fatal("Pi publication counter accepted a republication at validation/2")
	}

	// Falsifier: a failed call does not count as a publication.
	failedCall := strings.Join([]string{
		piBashToolCallLine("call-1", piRoundCommand),
		piToolResultLine("call-1", "bash", "Error: exit 1", true),
	}, "\n")
	if pubs := piRejectionRoundPublications(failedCall); len(pubs) != 0 {
		t.Fatalf("Pi publication counter counted a failed call as a publication: %v", pubs)
	}
}

// TestPiRejectionRoutesFreshChain pins the Pi-dialect topology extractor: four fresh
// subagent spawns (implementation, validation, implementation, validation) each
// completed via a status poll, matching the fail-safe fresh chain the Pi branch owes.
// Falsifiable: a chain missing the final completion, and a chain graded under the wrong
// branch, both red.
func TestPiRejectionRoutesFreshChain(t *testing.T) {
	session := piFreshChainSession(t)
	routes, branch := piRejectionRoutes(session)
	if branch != rejectionBranchFresh {
		t.Fatalf("Pi branch = %q, want %q", branch, rejectionBranchFresh)
	}
	want := []struct{ event, stage string }{
		{routeSpawn, "implementation"},
		{routeDone, "implementation"},
		{routeSpawn, "validation"},
		{routeDone, "validation"},
		{routeSpawn, "implementation"},
		{routeDone, "implementation"},
		{routeSpawn, "validation"},
		{routeDone, "validation"},
	}
	if len(routes) != len(want) {
		t.Fatalf("Pi routes produced %d events, want %d: %s",
			len(routes), len(want), rejectionTopologySummary(routes))
	}
	for i, expected := range want {
		if routes[i].event != expected.event || routes[i].stage != expected.stage {
			t.Fatalf("Pi route %d = %s/%s, want %s/%s: %s",
				i, routes[i].event, routes[i].stage, expected.event, expected.stage, rejectionTopologySummary(routes))
		}
	}
	// The fresh chain must grade green on its own branch.
	if err := assertRejectionWorkerTopology(branch, routes); err != nil {
		t.Fatalf("conforming Pi fresh chain graded red: %v\n%s", err, rejectionTopologySummary(routes))
	}
	// The fresh chain must NOT pass under the reuse branch.
	if err := assertRejectionWorkerTopology(rejectionBranchReuse, routes); err == nil {
		t.Fatal("the Pi fresh chain passed under the REUSE branch's expectations")
	}

	// Falsifier: a chain missing the final completion reds.
	truncated := strings.Join(strings.Split(session, "\n")[:len(strings.Split(session, "\n"))-2], "\n")
	truncRoutes, _ := piRejectionRoutes(truncated)
	if err := assertRejectionWorkerTopology(rejectionBranchFresh, truncRoutes); err == nil {
		t.Fatal("a truncated chain missing the final completion graded green")
	}
}

// piFreshChainSession builds a Pi session JSONL for the four-dispatch fresh rejection
// chain: implementation → validation (reject) → implementation (rework) → validation
// (re-review, pass). Each dispatch is a fresh `subagent` spawn; each completion is a
// status poll whose result shows "worker completed, exit 0, accept".
func piFreshChainSession(t *testing.T) string {
	t.Helper()
	mkRound := func(spawnID, runID, statusID, handle string) []string {
		return []string{
			piSpawnToolCallLine(spawnID, handle),
			piToolResultLine(spawnID, "subagent", "Async workflow ["+runID+"]", false),
			piStatusToolCallLine(statusID, runID),
			piToolResultLine(statusID, "subagent", "Run: "+runID+"\n1. worker completed, exit 0, accept", false),
		}
	}
	var lines []string
	lines = append(lines, mkRound("spawn-1", "run-impl-1", "status-1", piImplHandle)...)
	lines = append(lines, mkRound("spawn-2", "run-val-1", "status-2", piValHandle)...)
	lines = append(lines, mkRound("spawn-3", "run-impl-2", "status-3", piImplHandle)...)
	lines = append(lines, mkRound("spawn-4", "run-val-2", "status-4", piValHandle)...)
	return strings.Join(lines, "\n")
}

// TestPiRejectionRoutesHandleCorrelation confirms the Pi extractor derives the correct
// worker handle from the dispatch file path in the spawn task, so the independence
// check (validation worker ≠ implementation worker) holds on the derived handles.
func TestPiRejectionRoutesHandleCorrelation(t *testing.T) {
	session := piFreshChainSession(t)
	routes, _ := piRejectionRoutes(session)
	if routes[0].target != piImplHandle {
		t.Fatalf("route 0 target = %q, want %q", routes[0].target, piImplHandle)
	}
	if routes[2].target != piValHandle {
		t.Fatalf("route 2 target = %q, want %q", routes[2].target, piValHandle)
	}
	// The cycle-2 implementation and validation handles must match their cycle-1
	// counterparts (fresh dispatch under the same (slug, stage)-derived name).
	if routes[4].target != routes[0].target {
		t.Fatalf("cycle-2 implementation handle %q != cycle-1 %q", routes[4].target, routes[0].target)
	}
	if routes[6].target != routes[2].target {
		t.Fatalf("cycle-2 validation handle %q != cycle-1 %q", routes[6].target, routes[2].target)
	}
}

// TestPiRejectionRoutesRunIDCorrelation confirms the run-id correlation path: the
// spawn result's "Async workflow [run-id]" is matched to the status poll's
// `arguments.id`, so completion is attributed to the correct worker even when status
// polls for different workers interleave.
func TestPiRejectionRoutesRunIDCorrelation(t *testing.T) {
	// Two spawns with interleaved status polls — the run-id correlation must still
	// attribute each completion to the correct worker.
	lines := []string{
		piSpawnToolCallLine("spawn-1", piImplHandle),
		piToolResultLine("spawn-1", "subagent", "Async workflow [run-a]", false),
		piSpawnToolCallLine("spawn-2", piValHandle),
		piToolResultLine("spawn-2", "subagent", "Async workflow [run-b]", false),
		// Status polls in reverse order — run-b (validation) completes first.
		piStatusToolCallLine("status-b", "run-b"),
		piToolResultLine("status-b", "subagent", "Run: run-b\n1. worker completed, exit 0, accept", false),
		piStatusToolCallLine("status-a", "run-a"),
		piToolResultLine("status-a", "subagent", "Run: run-a\n1. worker completed, exit 0, accept", false),
	}
	session := strings.Join(lines, "\n")
	routes, _ := piRejectionRoutes(session)
	if len(routes) != 4 {
		t.Fatalf("interleaved session produced %d routes, want 4: %s", len(routes), rejectionTopologySummary(routes))
	}
	// spawn implementation, spawn validation, done validation, done implementation
	if routes[2].stage != "validation" || routes[2].event != routeDone {
		t.Fatalf("route 2 = %s/%s, want done/validation: %s", routes[2].event, routes[2].stage, rejectionTopologySummary(routes))
	}
	if routes[3].stage != "implementation" || routes[3].event != routeDone {
		t.Fatalf("route 3 = %s/%s, want done/implementation: %s", routes[3].event, routes[3].stage, rejectionTopologySummary(routes))
	}
}
