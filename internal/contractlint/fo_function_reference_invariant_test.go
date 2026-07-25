// ABOUTME: Guards the first-officer contract against mutable numeric procedure addresses.
// ABOUTME: Pins named-function ownership, local procedure order, load topology, and prompt size.
package contractlint

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// foFunctionReferencePaths is the 13-file union the mutable-address lint scans. The
// BUDGET is not this union: no single FO ever loads all three host adapters, so the
// byte ratchet below is per-host over foSharedLoadPaths + foHostLoadPaths.
var foFunctionReferencePaths = []string{
	"skills/first-officer/SKILL.md",
	"skills/first-officer/references/first-officer-shared-core.md",
	"skills/first-officer/references/fo-dispatch-core.md",
	"skills/first-officer/references/fo-merge-core.md",
	"skills/first-officer/references/fo-write-core.md",
	"skills/first-officer/references/claude-first-officer-runtime.md",
	"skills/first-officer/references/claude-fo-dispatch.md",
	"skills/first-officer/references/codex-first-officer-runtime.md",
	"skills/first-officer/references/pi-first-officer-runtime.md",
	"skills/present-gate/SKILL.md",
	"skills/feedback-rejection-flow/SKILL.md",
	"skills/fo-gate-lifecycle/SKILL.md",
	"skills/fo-dispatch-recovery/SKILL.md",
}

// foSharedLoadPaths is the load every host's FO pulls: the skill entry, the
// boot-resident shared core, the three deferred cores, and the trigger skills
// reachable on every host (present-gate, feedback-rejection-flow).
var foSharedLoadPaths = []string{
	"skills/first-officer/SKILL.md",
	"skills/first-officer/references/first-officer-shared-core.md",
	"skills/first-officer/references/fo-dispatch-core.md",
	"skills/first-officer/references/fo-merge-core.md",
	"skills/first-officer/references/fo-write-core.md",
	"skills/present-gate/SKILL.md",
	"skills/feedback-rejection-flow/SKILL.md",
	"skills/fo-gate-lifecycle/SKILL.md",
}

// foHostLoadPaths adds each host's adapter file(s) and the trigger skills reachable
// only on that host (fo-dispatch-recovery is named only from the Claude dispatch
// module). load(H) = foSharedLoadPaths + this.
var foHostLoadPaths = map[string][]string{
	"claude": {
		"skills/first-officer/references/claude-first-officer-runtime.md",
		"skills/first-officer/references/claude-fo-dispatch.md",
		"skills/fo-dispatch-recovery/SKILL.md",
	},
	"codex": {
		"skills/first-officer/references/codex-first-officer-runtime.md",
	},
	"pi": {
		"skills/first-officer/references/pi-first-officer-runtime.md",
	},
}

var mutableProcedureAddress = regexp.MustCompile(`(?i)(?:\bsteps?[- ]\d+(?:\.\d+)?(?:\s*(?:-|–|to)\s*\d+(?:\.\d+)?)?|\breuse[- ]conditions?[- ]?\d+|\btiers?[- ]\d+|\btiers?\s+\d+(?:\s+and\s+\d+)?|\bentry-point principle\s+\d+|\b(?:signals?|items?)\s*\(\d+(?:\s*,\s*\d+)*(?:\s*,?\s*or\s*\d+)?\s+above\))`)

type foAddressMatch struct {
	path  string
	line  int
	value string
}

func mutableFOAddresses(path, body string) []foAddressMatch {
	var out []foAddressMatch
	for i, line := range strings.Split(body, "\n") {
		for _, value := range mutableProcedureAddress.FindAllString(line, -1) {
			out = append(out, foAddressMatch{path: path, line: i + 1, value: value})
		}
	}
	return out
}

// foHostLoadBytes measures one host's real FO session load: the shared load plus
// that host's adapter files and host-only trigger skills.
func foHostLoadBytes(t *testing.T, host string) int {
	t.Helper()
	total := 0
	for _, rel := range append(append([]string{}, foSharedLoadPaths...), foHostLoadPaths[host]...) {
		total += len([]byte(readRepoFile(t, filepath.FromSlash(rel))))
	}
	return total
}

func TestFOFunctionReferenceInvariant(t *testing.T) {
	var matches []foAddressMatch
	for _, rel := range foFunctionReferencePaths {
		body := readRepoFile(t, filepath.FromSlash(rel))
		matches = append(matches, mutableFOAddresses(rel, body)...)
	}
	if len(matches) == 0 {
		return
	}
	var lines []string
	for _, match := range matches {
		lines = append(lines, fmt.Sprintf("%s:%d: %s", match.path, match.line, match.value))
	}
	t.Fatalf("found %d mutable procedure addresses:\n%s", len(matches), strings.Join(lines, "\n"))
}

func TestFOFunctionReferenceClassifierDiscriminates(t *testing.T) {
	failures := []struct{ path, body string }{
		{"shared", "use Startup step 2"},
		{"claude", "after step-0"},
		{"codex", "reuse-condition-1 fails"},
		{"pi", "follow reuse condition 4"},
		{"gate", "read at Startup step 3"},
		{"legacy", "tiers 1 and 2 failed"},
		{"dispatch", "completion signals (1, 2, or 3 above)"},
		{"shared-principle", "entry-point principle 1"},
	}
	for _, tc := range failures {
		if got := mutableFOAddresses(tc.path, tc.body); len(got) == 0 {
			t.Errorf("classifier missed planted %s address %q", tc.path, tc.body)
		}
	}
	passes := []string{
		"1. local ordered item", "0.5. drain messages", "feedback cycle 3",
		"exit 3", "AC-2", "version 0.24.0", "PR #490", "up to 5 attempts",
	}
	for _, body := range passes {
		if got := mutableFOAddresses("pass-control", body); len(got) != 0 {
			t.Errorf("classifier rejected semantic/local number %q: %+v", body, got)
		}
	}
}

func TestFOInstructionComponentCaps(t *testing.T) {
	for rel, cap := range map[string]int{
		"skills/first-officer/references/first-officer-shared-core.md": 26754,
		"skills/fo-gate-lifecycle/SKILL.md":                            6600,
	} {
		if got := len([]byte(readRepoFile(t, filepath.FromSlash(rel)))); got > cap {
			t.Errorf("%s = %d bytes, component cap %d", rel, got, cap)
		}
	}
}

// TestFOHostLoadSetsCoverAddressLintUnion keeps the two scopes in sync: the union of
// the shared load and every host's load equals the 13-file set the mutable-address
// lint scans, so neither list can drop or gain a file without the other noticing.
func TestFOHostLoadSetsCoverAddressLintUnion(t *testing.T) {
	union := map[string]bool{}
	for _, rel := range foSharedLoadPaths {
		union[rel] = true
	}
	for _, paths := range foHostLoadPaths {
		for _, rel := range paths {
			union[rel] = true
		}
	}
	if !setEqual(union, toSet(foFunctionReferencePaths)) {
		t.Fatalf("per-host load sets and the address-lint union diverged:\n  loads:   %v\n  address: %v", sortedSet(union), sortedSet(toSet(foFunctionReferencePaths)))
	}
}

func TestFOFunctionReferenceCheckpointMetrics(t *testing.T) {
	addresses := 0
	for _, rel := range foFunctionReferencePaths {
		addresses += len(mutableFOAddresses(rel, readRepoFile(t, filepath.FromSlash(rel))))
	}
	t.Logf("FO_FUNCTION_METRICS addresses=%d claude_bytes=%d codex_bytes=%d pi_bytes=%d",
		addresses, foHostLoadBytes(t, "claude"), foHostLoadBytes(t, "codex"), foHostLoadBytes(t, "pi"))
}

func TestFirstOfficerReferenceTopology(t *testing.T) {
	root := skillsRoot(t)
	entry := readRepoFile(t, filepath.Join("skills", "first-officer", "SKILL.md"))
	var imports []string
	for _, line := range strings.Split(entry, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "@references/") {
			imports = append(imports, strings.TrimSpace(line))
		}
	}
	want := []string{
		"@references/first-officer-shared-core.md",
	}
	if strings.Join(imports, "\n") != strings.Join(want, "\n") {
		t.Fatalf("eager imports = %v, want lazy write/merge topology %v", imports, want)
	}
	for _, rel := range []string{
		"references/first-officer-shared-core.md",
		"references/fo-write-core.md",
		"references/fo-merge-core.md",
	} {
		path := filepath.Join(root, "first-officer", rel)
		if info, err := os.Stat(path); err != nil || info.Size() == 0 {
			t.Errorf("canonical first-officer body %s does not resolve non-empty: %v", rel, err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "fo-write-core")); !os.IsNotExist(err) {
		t.Errorf("redundant standalone fo-write-core entry surface remains: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "first-officer", "references", "fo-smallest-sufficient-mechanism.md")); !os.IsNotExist(err) {
		t.Errorf("duplicated smallest-sufficient eager body remains: %v", err)
	}
}

func TestFODeferredDispatchOwnerLoadsBeforeUse(t *testing.T) {
	shared := readRepoFile(t, filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md"))
	loads := foMarkdownSection(t, shared, "## Deferred load points")
	for _, want := range []string{
		"references/fo-dispatch-core.md",
		"before invoking `«dispatch.next-action»()`",
		"`«dispatch.build»` output is not a dispatch",
		"`«worker.spawn»`",
	} {
		if !strings.Contains(loads, want) {
			t.Errorf("deferred dispatch load point missing precondition %q", want)
		}
	}

	interaction := foMarkdownSection(t, shared, "## «interaction.boundary»(): route interactive and headless launch behavior")
	if !strings.Contains(interaction, "Read the deferred dispatch owner before the first dispatch") {
		t.Error("headless interaction can enter dispatch without loading its named owner")
	}

	engage := foMarkdownSection(t, shared, "## «engage»(workflow): converge one named workflow, then run its event loop to a stopping condition")
	for _, want := range []string{
		"including one made ready by gate approval",
		"observed `«worker.spawn»`",
		"before waiting for completion",
	} {
		if !strings.Contains(engage, want) {
			t.Errorf("engage dispatch boundary missing %q", want)
		}
	}

	principles := foMarkdownSection(t, shared, "## Working Principles")
	for _, want := range []string{
		"Commissioned workflow dispatch is mandatory",
		"in-house execution is not a lower rung",
		"every ready entity, including one advanced by gate approval",
	} {
		if !strings.Contains(principles, want) {
			t.Errorf("smallest-sufficient commissioned-dispatch boundary missing %q", want)
		}
	}
}

func TestFOGateLifecycleOwnsEveryEngagedEntry(t *testing.T) {
	shared := readRepoFile(t, filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md"))
	loads := foMarkdownSection(t, shared, "## Deferred load points")
	for _, want := range []string{
		`Skill(skill="spacedock:fo-gate-lifecycle")`,
		"headless with or without conn", "`engage`", "gated worker completion",
		"open/pending/revise/hold/stale/consumed resume",
		"interactive gated greet only names the gate and stops load-free",
	} {
		if !strings.Contains(loads, want) {
			t.Errorf("gate-entry funnel missing %q", want)
		}
	}
	interaction := foMarkdownSection(t, shared, "## «interaction.boundary»(): route interactive and headless launch behavior")
	if strings.Contains(interaction, `Skill(skill="spacedock:present-gate")`) ||
		!strings.Contains(interaction, "At each gate enter through `«gate.lifecycle»`") {
		t.Error("headless gate route bypasses the deferred lifecycle owner")
	}
	completion := foMarkdownSection(t, shared, "## Completion and Gates")
	if !strings.Contains(completion, `Skill(skill="spacedock:fo-gate-lifecycle")`) {
		t.Error("worker-completion gate route does not load the deferred lifecycle")
	}
	autoRoute := "reviewer recommends `REJECTED`"
	if strings.Count(completion, autoRoute) != 1 ||
		!strings.Contains(completion, "`«feedback.route»` before Captain presentation") {
		t.Error("reviewer-REJECTED correction is not uniquely owned before Captain presentation")
	}
	lifecycle, presenter := readRepoFile(t, filepath.Join("skills", "fo-gate-lifecycle", "SKILL.md")), readRepoFile(t, filepath.Join("skills", "present-gate", "SKILL.md"))
	for _, want := range []string{
		"absence or exposure", "halt before mutation", "refresh or a fresh build", "reject retired `--directive` exposure",
		"presentation completes only after", "exact bound Briefing id/digest",
		"precede decision record", "delegated conn does not waive",
		"explicit Captain grant in the active conversation",
		"including one issued later in that conversation",
		"leave the gate open", "FO-rendered decision as `agent:first-officer` with a nonblank reason, never `person:captain`",
	} {
		if !strings.Contains(strings.ToLower(lifecycle), strings.ToLower(want)) {
			t.Errorf("gate lifecycle missing fail-closed/presentation contract %q", want)
		}
	}
	for _, forbidden := range []string{"--directive-file", "EXACT_QUOTED_CAPTAIN_GRANT"} {
		if strings.Contains(lifecycle, forbidden) {
			t.Errorf("gate lifecycle retains caller-controlled provenance surface %q", forbidden)
		}
	}
	lifecycleOrder := []string{
		"**Capability preflight.**",
		"exactly one fresh `gate --help`",
		"**Prepare and bind.**",
		"gate prepare ENTITY",
		"Require the emitted `room`, `briefing`, `digest`, and `state=open` lines",
		"**Record and durably close.**",
		"gate record ENTITY --decision",
		"**Route fail-closed.**",
		"gate consume ENTITY",
	}
	orderedLifecycle := func(body string) bool {
		offset := 0
		for _, anchor := range lifecycleOrder {
			i := strings.Index(body[offset:], anchor)
			if i < 0 {
				return false
			}
			offset += i + len(anchor)
		}
		return true
	}
	if !orderedLifecycle(lifecycle) {
		t.Fatal("gate lifecycle does not preserve preflight -> prepare/bind -> close -> consume order")
	}
	for _, anchor := range lifecycleOrder {
		if orderedLifecycle(strings.ReplaceAll(lifecycle, anchor, "removed")) {
			t.Fatalf("ordered gate lifecycle mutation survived removal of %q", anchor)
		}
	}
	for _, want := range []string{
		"exactly one fresh `gate --help`", "`prepare`, `record`, `validate`, `eligibility`, `consume`",
		"`--question`, `--artifact`, `--summary`, `--reference`, `--workflow-dir`",
		"gate prepare ENTITY", "emitted clean absolute room", "flat Markdown-plus-companion room unit",
		`Skill(skill="spacedock:present-gate")`, "an override replaces only chat display",
	} {
		if !strings.Contains(lifecycle, want) {
			t.Errorf("gate lifecycle missing provider-neutral preparation contract %q", want)
		}
	}
	for _, want := range []string{
		"Wait for `gate prepare` to succeed", "/subspace:r gate <room>",
		"Reconstruct no authority", "Perform no provider probe or fallback selection",
	} {
		if !strings.Contains(presenter, want) {
			t.Errorf("present-gate missing room-only override contract %q", want)
		}
	}
	for _, want := range []string{
		"exactly one root-assistant message", "entity and stage",
		"exact bound Briefing id and digest", "one recommendation",
		"decision ask", "before the next decision-mutation tool call",
	} {
		if !strings.Contains(presenter, want) {
			t.Errorf("present-gate missing semantic review owner %q", want)
		}
	}
	if strings.Contains(lifecycle, "Gate review:") || strings.Contains(lifecycle, "Decision:") {
		t.Error("lifecycle duplicates presenter markers instead of waiting for semantic presentation")
	}
	mapping := "Map Captain calls before recording: `approve` maps to `approve` with an accepts-direction evidence reason; `redo with feedback` maps to `revise` with an accepts-direction reason; `reject` with `feedback-to` maps to `revise` with a rejects-direction reason; `reject` without `feedback-to` maps to `hold` with a pause reason; `hold` maps to `hold` with a pause reason; `not yet` maps to `hold` with a pause reason naming what remains. Routed redo/reject reasons include concrete asks and invoke `«feedback.route»` after the close commit; hold decisions commit and stop at the gate."
	gradeMapping := func(body string) bool {
		parts := strings.SplitN(body, "Map Captain calls before recording: ", 2)
		return len(parts) == 2 && strings.SplitN(parts[1], "\n\n", 2)[0] == strings.TrimPrefix(mapping, "Map Captain calls before recording: ")
	}
	if !gradeMapping(lifecycle) {
		t.Fatal("Captain mapping must preserve all six decisions, reason classes, concrete asks, and route/hold behavior")
	}
	for _, token := range []string{"`approve`", "`redo with feedback`", "`reject` with `feedback-to`", "`reject` without `feedback-to`", "`hold` maps", "`not yet`", "accepts-direction", "rejects-direction", "pause reason", "concrete asks", "`«feedback.route»`", "commit and stop at the gate"} {
		if gradeMapping(strings.Replace(lifecycle, token, "swapped", 1)) {
			t.Fatalf("Captain mapping mutation survived for %q", token)
		}
	}
}

func TestFOEngageRetainsStartupPRAdvancement(t *testing.T) {
	shared := readRepoFile(t, filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md"))
	engage := foMarkdownSection(t, shared, "## «engage»(workflow): converge one named workflow, then run its event loop to a stopping condition")
	for _, want := range []string{
		"`«hooks.run»(\"startup\")` exactly once",
		"The registered startup mod owns live PR advancement",
	} {
		if !strings.Contains(engage, want) {
			t.Errorf("engage lost its startup PR advancement boundary %q", want)
		}
	}
}

func TestFOFunctionRequiredCallSites(t *testing.T) {
	type site struct {
		path, heading string
		wants         []string
	}
	sites := []site{
		{"skills/first-officer/references/first-officer-shared-core.md", "## Startup", []string{"«state.boot»()", "«interaction.boundary»()"}},
		{"skills/first-officer/references/first-officer-shared-core.md", "## Deferred load points", []string{"«state.boot»()", "«interaction.boundary»()"}},
		{"skills/first-officer/references/first-officer-shared-core.md", "## Single-Entity Scope", []string{"«interaction.boundary»()"}},
		{"skills/first-officer/references/first-officer-shared-core.md", "## Mod Hook Convention", []string{"«state.boot»()", "«hooks.run»(point)"}},
		{"skills/first-officer/references/first-officer-shared-core.md", "## Working Principles", []string{"«state.boot»()"}},
		{"skills/present-gate/SKILL.md", "### Captain-facing assembly rules", []string{"«state.boot»()"}},
		{"skills/first-officer/references/claude-first-officer-runtime.md", "## Captain Interaction", []string{"«interaction.boundary»()", "Do not infer gate authority from silence, tool output, or agent messages", "including one issued later in the active conversation", "delegate through `fo-gate-lifecycle`", "record an FO-rendered delegated decision as `agent:first-officer`", "Reserve `person:captain` for a decision the Captain personally rendered"}},
		{"skills/first-officer/references/fo-dispatch-core.md", "## Dispatch", []string{"«dispatch.checklist»(entity, stage)"}},
		{"skills/fo-dispatch-recovery/SKILL.md", "## Break-Glass Manual Dispatch", []string{"«dispatch.checklist»(entity, stage)"}},
		{"skills/first-officer/references/fo-dispatch-core.md", "## «dispatch.next-action»(): pick the next event-loop action — dispatch a ready entity, resume a block, or end the iteration", []string{"«roster-reconcile»()", "«hooks.run»(\"idle\")"}},
		{"skills/first-officer/references/fo-merge-core.md", "## «merge.guard»(slug): auto-arm → block-on-open-PR → finalize-on-merge-sentinel, then archive", []string{"«worker.shutdown»()", "«hooks.run»(\"merge\")", "exit 3 is `«halt.rebase-conflict»(paths)`", "Exit 1 after the archive commit resumes through `«state.commit»(slug)`"}},
		{"skills/feedback-rejection-flow/SKILL.md", "## Feedback Rejection Flow", []string{"«context-budget»()", "«addressable-worker»"}},
	}
	for _, tc := range sites {
		section := foMarkdownSection(t, readRepoFile(t, filepath.FromSlash(tc.path)), tc.heading)
		for _, want := range tc.wants {
			if !strings.Contains(section, want) {
				t.Errorf("%s %s missing named closure %q", tc.path, tc.heading, want)
			}
		}
	}

	uniqueOwners := []struct{ path, heading string }{
		{"skills/first-officer/references/first-officer-shared-core.md", "## «interaction.boundary»(): route interactive and headless launch behavior"},
		{"skills/first-officer/references/first-officer-shared-core.md", "## «hooks.run»(point): run registered lifecycle hooks"},
		{"skills/first-officer/references/fo-dispatch-core.md", "## «dispatch.checklist»(entity, stage): assemble dispatch linchpins"},
	}
	for _, owner := range uniqueOwners {
		body := readRepoFile(t, filepath.FromSlash(owner.path))
		if got := strings.Count(body, owner.heading+"\n"); got != 1 {
			t.Errorf("%s owner heading %q count=%d, want 1", owner.path, owner.heading, got)
		}
	}
}

func TestFOFunctionNormalizationPreservationSuite(t *testing.T) {
	checks := []struct {
		name, path string
		headings   []string
		anchors    []string
	}{
		{"startup-interaction", "skills/first-officer/references/first-officer-shared-core.md", []string{"## «interaction.boundary»(): route interactive and headless launch behavior", "## Startup"}, []string{"Interactive", "Headless", "given the conn", "STOP"}},
		{"dispatch-checklist", "skills/first-officer/references/fo-dispatch-core.md", []string{"## «dispatch.checklist»(entity, stage): assemble dispatch linchpins", "## Dispatch"}, []string{"≤3", "Outputs:", "acceptance criteria", "not stage actions"}},
		{"reuse", "skills/first-officer/references/fo-dispatch-core.md", []string{"## Reuse and Fresh Dispatch"}, []string{"«context-budget»()", "«addressable-worker»", "fresh: true", "«reuse.model-match»"}},
		{"completion", "skills/first-officer/references/fo-dispatch-core.md", []string{"## «completion-signal»: the signals that trigger the completion-verify path"}, []string{"runtime-binding", "stage report"}},
		{"terminal-teardown", "skills/first-officer/references/fo-merge-core.md", []string{"## «merge.guard»(slug): auto-arm → block-on-open-PR → finalize-on-merge-sentinel, then archive"}, []string{"teardown", "best-effort", "drop them from session memory"}},
		{"hooks", "skills/first-officer/references/first-officer-shared-core.md", []string{"## «hooks.run»(point): run registered lifecycle hooks", "## Mod Hook Convention"}, []string{"startup", "idle", "merge", "alphabetically"}},
		{"write-permission", "skills/first-officer/references/fo-write-core.md", []string{"## Mutation Gate"}, []string{"classify", "blocked-product", "exact task and target path"}},
		{"approval-evidence-feedback", "skills/first-officer/references/first-officer-shared-core.md", []string{"## Completion and Gates"}, []string{"Stage Report", "AC coverage cross-check", "captain", "«feedback.route»"}},
		{"claude-binding", "skills/first-officer/references/claude-fo-dispatch.md", []string{"## Inter-Agent Communication"}, []string{"Agent", "SendMessage"}},
		{"pi-binding", "skills/first-officer/references/pi-first-officer-runtime.md", []string{"## Runtime implementation"}, []string{"subagent", "intercom", "member_shutdown"}},
	}
	for _, tc := range checks {
		t.Run(tc.name, func(t *testing.T) {
			section := markdownSectionOneOf(t, readRepoFile(t, filepath.FromSlash(tc.path)), tc.headings...)
			for _, anchor := range tc.anchors {
				if !strings.Contains(section, anchor) {
					t.Errorf("preservation section missing %q", anchor)
				}
			}
		})
	}
}

func markdownSectionOneOf(t *testing.T, body string, headings ...string) string {
	t.Helper()
	for _, heading := range headings {
		if strings.Contains(body, heading+"\n") {
			return foMarkdownSection(t, body, heading)
		}
	}
	t.Fatalf("none of the preservation headings found: %v", headings)
	return ""
}

func TestFOLocalOrderedProceduresPreserved(t *testing.T) {
	type procedure struct {
		path, heading string
		want          []string
		anchors       []string
	}
	cases := []procedure{
		{"skills/first-officer/references/first-officer-shared-core.md", "## Startup", []string{"1", "2", "3"}, []string{"Binary version gate", "«state.boot»", "«interaction.boundary»"}},
		{"skills/first-officer/references/fo-dispatch-core.md", "## Dispatch", sequence(1, 9), []string{"entity file", "«dispatch.checklist»", "conflicts", "dispatch_agent_id", "status --workflow-dir", "Commit", "worktree", "«dispatch.build»", "«completion-signal»"}},
		{"skills/first-officer/references/fo-dispatch-core.md", "## Reuse and Fresh Dispatch", sequence(0, 4), []string{"«context-budget»", "«addressable-worker»", "fresh: true", "worktree", "«reuse.model-match»"}},
		{"skills/first-officer/references/fo-dispatch-core.md", "## «dispatch.next-action»(): pick the next event-loop action — dispatch a ready entity, resume a block, or end the iteration", []string{"0.5", "1", "2", "3"}, []string{"«addressable-worker»", "mod-block", "status --next", "«hooks.run»", "«roster-reconcile»"}},
		{"skills/present-gate/SKILL.md", "## Presentation channels", sequence(1, 4), []string{"Pass only the emitted room after its bind commit", "/subspace:r gate <room>", "Reconstruct no authority", "Perform no provider probe or fallback selection", "gate record <entity> --room <room>", "do not fall back to chat after the selected handoff"}},
		{"skills/present-gate/SKILL.md", "### Captain-facing assembly rules", sequence(1, 11), []string{"Lede first", "Chosen direction", "Stage Report", "Reviewer findings", "Recommendation", "Bounce-back", "format-pedantry", "worktree", "Target length", "declared label", "verification state"}},
		{"skills/feedback-rejection-flow/SKILL.md", "## Feedback Rejection Flow", sequence(1, 7), []string{"feedback-to", "Feedback Cycles", "cycle 3", "«context-budget»", "«addressable-worker»", "reviewer", "gate flow"}},
	}
	for _, tc := range cases {
		section := foMarkdownSection(t, readRepoFile(t, filepath.FromSlash(tc.path)), tc.heading)
		if got := orderedMarkers(section); strings.Join(got, ",") != strings.Join(tc.want, ",") {
			t.Errorf("%s %s markers=%v, want %v", tc.path, tc.heading, got, tc.want)
		}
		for _, anchor := range tc.anchors {
			if !strings.Contains(section, anchor) {
				t.Errorf("%s %s missing ordered-procedure anchor %q", tc.path, tc.heading, anchor)
			}
		}
	}
}

var orderedMarker = regexp.MustCompile(`(?m)^(\d+(?:\.\d+)?)\.`)

func orderedMarkers(section string) []string {
	var out []string
	for _, match := range orderedMarker.FindAllStringSubmatch(section, -1) {
		out = append(out, match[1])
	}
	return out
}

func sequence(first, last int) []string {
	out := make([]string, 0, last-first+1)
	for n := first; n <= last; n++ {
		out = append(out, strconv.Itoa(n))
	}
	return out
}

func foMarkdownSection(t *testing.T, body, heading string) string {
	t.Helper()
	needle := heading + "\n"
	start := strings.Index(body, needle)
	if start < 0 || (start > 0 && body[start-1] != '\n') {
		t.Fatalf("missing line-anchored heading %q", heading)
	}
	level := len(heading) - len(strings.TrimLeft(heading, "#"))
	rest := body[start+len(needle):]
	end := len(rest)
	for pos, line := range strings.Split(rest, "\n") {
		_ = pos
		if !strings.HasPrefix(line, "#") {
			continue
		}
		lineLevel := len(line) - len(strings.TrimLeft(line, "#"))
		if lineLevel <= level && strings.HasPrefix(line[lineLevel:], " ") {
			end = strings.Index(rest, "\n"+line) + 1
			break
		}
	}
	return body[start : start+len(needle)+end]
}
