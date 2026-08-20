package ensigncycle

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// Offline proof for the rejection-flow worker topology. Two things are proven here,
// and they are deliberately proven by different means, because one fixture cannot
// honestly do both:
//
//   - EXTRACTION is proven against REAL captured Codex rollout bytes
//     (testdata/rejection_topology/codex-reuse-chain.jsonl). That file is a verbatim
//     slice of a real multi-agent parent rollout — real handles, real key names, real
//     event ordering, with only the encrypted `message` blobs shortened and the
//     `agent_message` texts truncated, neither of which the extractor reads. It
//     carries the three shapes that break a hand-written parser: a `followup_task`
//     naming its worker under `target` rather than `task_name`, a REFUSED spawn whose
//     output is a plain error string instead of JSON, and an intermediate non-final
//     `agent_message` from a worker that is already running.
//   - GRADING is proven over route-level mutants of that extracted chain. Grading is
//     where the shape logic lives, and mutating routes states each non-conforming
//     shape as the sequence it IS, rather than fabricating rollout bytes for
//     topologies that were never observed on a real host.

func codexCapturedRoutes(t *testing.T) []rejectionRoute {
	t.Helper()
	return codexRejectionRoutes(readFile(t, filepath.Join("testdata", "rejection_topology", "codex-reuse-chain.jsonl")))
}

// TestRejectionTopologyExtractsCapturedCodexChain is the real-bytes half: the
// captured rollout must yield exactly the determined reuse chain, in order, with the
// handles correlated across the rooted/bare spelling difference. The sub-checks below
// each name the real-world shape that would break if the extractor regressed.
func TestRejectionTopologyExtractsCapturedCodexChain(t *testing.T) {
	routes := codexCapturedRoutes(t)
	want := []struct{ event, stage string }{
		{routeSpawn, "implementation"},
		{routeDone, "implementation"},
		{routeSpawn, "validation"},
		{routeDone, "validation"},
		{routeReuse, "implementation"},
		{routeDone, "implementation"},
		{routeReuse, "validation"},
		{routeDone, "validation"},
	}
	if len(routes) != len(want) {
		t.Fatalf("captured rollout produced %d routing events, want %d: %s",
			len(routes), len(want), rejectionTopologySummary(routes))
	}
	for i, expected := range want {
		if routes[i].event != expected.event || routes[i].stage != expected.stage {
			t.Fatalf("captured route %d = %s/%s, want %s/%s: %s",
				i, routes[i].event, routes[i].stage, expected.event, expected.stage, rejectionTopologySummary(routes))
		}
	}
	// The reuse rounds must correlate back to the handles their spawns opened. The
	// rollout spells one of those two sides rooted (`/root/NAME`) and the other bare,
	// so an extractor that compared raw strings would report zero reuse here.
	if routes[4].target != routes[0].target || routes[6].target != routes[2].target {
		t.Fatalf("reuse rounds did not correlate to their spawned handles: %#v", routes)
	}
	// The refused spawn ("agent path ... already exists") opened no worker, so it
	// must contribute no routing event; it is the 3rd spawn_agent in the fixture.
	if spawns := countRouteEvents(routes, routeSpawn); spawns != 2 {
		t.Errorf("captured chain has %d spawn events, want 2 — the refused spawn was counted as a dispatch", spawns)
	}
	// The implementation worker emits a non-final `agent_message` before its
	// FINAL_ANSWER; counting it would double this round's completion.
	if dones := countRouteEvents(routes, routeDone); dones != 4 {
		t.Errorf("captured chain has %d completion events, want 4 — an intermediate worker message was counted", dones)
	}
	if err := assertRejectionWorkerTopology(rejectionBranchReuse, routes); err != nil {
		t.Fatalf("captured conforming reuse chain graded red on its own branch: %v", err)
	}
}

// TestRejectionTopologyBranchesAreMutuallyExclusive is AC-3's second half: a run
// conforming on one branch must NOT also pass under the other branch's expectations,
// which is what stops the grading from silently accepting two topologies.
func TestRejectionTopologyBranchesAreMutuallyExclusive(t *testing.T) {
	reuse := codexCapturedRoutes(t)
	if err := assertRejectionWorkerTopology(rejectionBranchFresh, reuse); err == nil {
		t.Error("the captured reuse chain passed under the FRESH branch's expectations")
	}
	fresh := freshChainRoutes()
	if err := assertRejectionWorkerTopology(rejectionBranchFresh, fresh); err != nil {
		t.Fatalf("the fail-safe fresh chain graded red on its own branch: %v", err)
	}
	if err := assertRejectionWorkerTopology(rejectionBranchReuse, fresh); err == nil {
		t.Error("the fail-safe fresh chain passed under the REUSE branch's expectations")
	}
}

// freshChainRoutes is the Claude fail-safe shape: four fresh dispatches, the cycle-2
// pair re-opening workers under the same (slug, stage)-derived names the superseded
// ones used. Reusing the name is what makes a spawn COUNT the only structural proof
// that the cycle-2 worker is a new process, so the fresh branch cannot grade
// distinctness by comparing handles.
func freshChainRoutes() []rejectionRoute {
	const impl, val = "spacedock-ensign-rejection-task-implementation", "spacedock-ensign-rejection-task-validation"
	return []rejectionRoute{
		{index: 0, event: routeSpawn, stage: "implementation", target: impl},
		{index: 1, event: routeDone, stage: "implementation", target: impl},
		{index: 2, event: routeSpawn, stage: "validation", target: val},
		{index: 3, event: routeDone, stage: "validation", target: val},
		{index: 4, event: routeSpawn, stage: "implementation", target: impl},
		{index: 5, event: routeDone, stage: "implementation", target: impl},
		{index: 6, event: routeSpawn, stage: "validation", target: val},
		{index: 7, event: routeDone, stage: "validation", target: val},
	}
}

// TestRejectionTopologyGreensRetainedRepeatChains is AC-1's value half. Run
// 32270990171 graded 1/2: its Claude lane's fresh chain passed and its Codex lane's
// reuse chain reddened with "owes 8 routing events, the run produced 10" while order
// held, identity held, and every other durable oracle on that lane was green. Both
// retained chains grade GREEN here, which is the 2/2.
func TestRejectionTopologyGreensRetainedRepeatChains(t *testing.T) {
	captured := codexCapturedRoutes(t)
	// The repeat pair that run actually produced, at the position its retained digest
	// records it: rows 4/5 of
	// live-artifacts/codex/codex-shared-scenarios/rejection-flow/rejection-topology.tsv,
	// `reuse validation` then `done validation` on the handle the validation stage
	// opened. freshChainRoutes is that same run's Claude lane, handle names included.
	reviewer := rejectionRoute{event: routeReuse, stage: "validation", target: captured[2].target}
	for _, tc := range []struct {
		name   string
		branch rejectionBranch
		routes []rejectionRoute
	}{
		{"retained codex reuse chain, the repeat before the rework (10 events)", rejectionBranchReuse, insertRound(captured, 4, reviewer)},
		{"retained claude fresh chain (8 events)", rejectionBranchFresh, freshChainRoutes()},
		// The budget is a TOTAL, so the same one round spent after the rework greens
		// too. A bar that tolerated only the position already observed would be the
		// position-fragile bar this entity exists to remove.
		{"the same one round spent after the rework", rejectionBranchReuse, insertRound(captured, len(captured), reviewer)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := assertRejectionWorkerTopology(tc.branch, tc.routes); err != nil {
				t.Fatalf("conforming chain graded red: %v\n%s", err, rejectionTopologySummary(tc.routes))
			}
		})
	}
}

// insertRound splices a whole ROUND — a dispatch and the completion that closes it —
// into a flat chain at `at`. A repeat has to enter as a pair: a dispatch without its
// own completion is a different (unpairable) shape.
func insertRound(routes []rejectionRoute, at int, dispatch rejectionRoute) []rejectionRoute {
	out := append([]rejectionRoute(nil), routes[:at]...)
	out = append(out, dispatch, rejectionRoute{event: routeDone, stage: dispatch.stage, target: dispatch.target})
	return append(out, routes[at:]...)
}

// TestRejectionTopologyRedsNonConformingShapes is AC-2: every non-conforming shape in
// hand reds under the named code, and every red names the invariant it violated. A
// single case greening here is the widened-bar failure this entity was warned about —
// a permissive bar does not get investigated the way a false red does.
func TestRejectionTopologyRedsNonConformingShapes(t *testing.T) {
	captured := codexCapturedRoutes(t)
	reviewer := rejectionRoute{event: routeReuse, stage: "validation", target: captured[2].target}
	retained := insertRound(captured, 4, reviewer)
	for _, tc := range []struct {
		name   string
		branch rejectionBranch
		routes []rejectionRoute
		want   string
	}{
		{
			// The preserved loop greens: ONE worker advanced through every round, so
			// it reviewed its own fix. Both runs graded green before this layer
			// because the Codex leg asserted no topology at all.
			name:   "single-worker self-reviewing chain (preserved loop greens)",
			branch: rejectionBranchReuse,
			routes: selfReviewChainRoutes(),
			want:   "the candidate was never reviewed",
		},
		{
			name:   "the candidate ran and nothing else did",
			branch: rejectionBranchReuse,
			routes: captured[:2],
			want:   "no validation round at all",
		},
		{
			name:   "the review rejected and the rework was never routed",
			branch: rejectionBranchReuse,
			routes: captured[:4],
			want:   "never routed as a rework",
		},
		{
			// FO residual mode 2: the fix target is fresh-dispatched while the
			// followup route to the live producer is available.
			name:   "fresh rework while the followup route is live",
			branch: rejectionBranchReuse,
			routes: mutateRoute(captured, 4, func(r *rejectionRoute) { r.event = routeSpawn }),
			want:   "re-opened the implementation stage with a spawn",
		},
		{
			name:   "re-review runs before the rework",
			branch: rejectionBranchReuse,
			routes: swapRoutes(captured, 4, 6),
			want:   "position 4 is not closed by its own completion",
		},
		{
			// The routing reaches a validation-stage round, but the handle it reaches
			// is the fix producer — the registry's "independently checked" outcome
			// violated while every count still looks right.
			name:   "cycle-2 re-review routed to the fix producer",
			branch: rejectionBranchReuse,
			routes: reviewerCollapsedRoutes(t),
			want:   "reviewing its own output is not the independent check",
		},
		{
			// On the reuse branch a retargeted review usually trips handle
			// consistency first, so independence also has to be proven where handle
			// comparison is the ONLY check left.
			name:   "fresh-branch chain whose reviews run on the producer",
			branch: rejectionBranchFresh,
			routes: freshSelfReviewRoutes(),
			want:   "reviewing its own output is not the independent check",
		},
		{
			name:   "no re-review after the rework",
			branch: rejectionBranchReuse,
			routes: captured[:6],
			want:   "was never re-reviewed",
		},
		{
			name:   "two validators spawned up front instead of one per round",
			branch: rejectionBranchFresh,
			routes: mutateRoute(freshChainRoutes(), 1, func(r *rejectionRoute) { r.event, r.stage = routeSpawn, "validation" }),
			want:   "position 0 is not closed by its own completion",
		},
		{
			name:   "two repeated reviews in the same segment",
			branch: rejectionBranchReuse,
			routes: insertRound(retained, 4, reviewer),
			want:   "repeated 2 validation rounds (3 before the rework",
		},
		{
			// Split one per segment. The budget is a total, so this reds exactly as
			// two in one segment do.
			name:   "one repeated review on each side of the rework",
			branch: rejectionBranchReuse,
			routes: insertRound(retained, len(retained), reviewer),
			want:   "repeated 2 validation rounds (2 before the rework",
		},
		{
			// In a one-cycle journey a second rework means the first fix did not
			// hold, which is the flailing the bar exists to keep red.
			name:   "a second rework after the re-review",
			branch: rejectionBranchReuse,
			routes: insertRound(captured, len(captured), rejectionRoute{event: routeReuse, stage: "implementation", target: captured[0].target}),
			want:   "opens a second rework",
		},
		{
			// Widening the conforming set must not collapse the two branches: the
			// 10-event reuse chain must not ALSO pass as the fail-safe fresh shape.
			name:   "the retained reuse chain graded on the fresh branch",
			branch: rejectionBranchFresh,
			routes: retained,
			want:   "took the fail-safe FRESH branch",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := assertRejectionWorkerTopology(tc.branch, tc.routes)
			if err == nil {
				t.Fatalf("non-conforming shape graded green: %s", rejectionTopologySummary(tc.routes))
			}
			graded, ok := err.(*gradedErr)
			if !ok {
				t.Fatalf("non-conforming shape produced an infrastructure failure, not a graded red: %v", err)
			}
			if graded.code != "rejection-worker-topology" {
				t.Errorf("code = %q, want rejection-worker-topology", graded.code)
			}
			if !strings.Contains(graded.Error(), tc.want) {
				t.Errorf("detail = %q, want it to name %q", graded.Error(), tc.want)
			}
			// AC-2's other half, asserted on EVERY red rather than one of them: no
			// diagnostic path may fall back to an owed-event count. That message is
			// what made the retained live red unreadable — it named a number instead
			// of a violated invariant, so no reader could tell a benign extra round
			// from a broken journey without re-deriving the chain by hand.
			for _, banned := range []string{"routing events", "owes 8"} {
				if strings.Contains(graded.Error(), banned) {
					t.Errorf("detail = %q fell back to the bare-count vocabulary %q", graded.Error(), banned)
				}
			}
		})
	}
}

// freshSelfReviewRoutes is the fresh branch's collapsed shape: four fresh spawns whose
// validation rounds re-open the IMPLEMENTATION worker's name. Its routing is
// internally conforming — every round after a stage's first is a fresh spawn — so only
// the round-level independence check can catch it.
func freshSelfReviewRoutes() []rejectionRoute {
	routes := freshChainRoutes()
	producer := routes[0].target
	for i := range routes {
		routes[i].target = producer
	}
	return routes
}

// selfReviewChainRoutes is the shape the two preserved composed-tree greens actually
// drove. Their public streams carry no topology (only `wait` collab items), so the
// routing is taken from the FO's own `dispatch build` sequence in those streams: one
// initial `--stage implementation` dispatch followed by three `--advance` calls, i.e.
// one worker advanced through validation, the rework, and the re-review.
func selfReviewChainRoutes() []rejectionRoute {
	const only = "spacedock_ensign_rejection_task_implementation"
	routes := make([]rejectionRoute, 0, 8)
	for round, event := range []string{routeSpawn, routeReuse, routeReuse, routeReuse} {
		routes = append(routes,
			rejectionRoute{index: round * 2, event: event, stage: "implementation", target: only},
			rejectionRoute{index: round*2 + 1, event: routeDone, stage: "implementation", target: only})
	}
	return routes
}

// reviewerCollapsedRoutes keeps the conforming chain's every event and stage, and
// changes only WHICH handle the validation rounds reached: the fix producer's. Both
// events of each retargeted ROUND move together, because a dispatch to one worker
// closed by another worker's completion is a different (unpairable) defect and no
// extractor can emit it — the collapse this case is about is internally consistent,
// which is exactly why only a round-level identity check catches it.
func reviewerCollapsedRoutes(t *testing.T) []rejectionRoute {
	t.Helper()
	routes := codexCapturedRoutes(t)
	producer := routes[4].target
	retarget := func(r *rejectionRoute) { r.target = producer }
	return mutateRound(mutateRound(routes, 6, retarget), 2, retarget)
}

func mutateRoute(routes []rejectionRoute, at int, mutate func(*rejectionRoute)) []rejectionRoute {
	out := append([]rejectionRoute(nil), routes...)
	mutate(&out[at])
	return out
}

// mutateRound applies a mutation to a round's dispatch AND to the completion that
// closes it, the pair being the unit the grammar reads.
func mutateRound(routes []rejectionRoute, at int, mutate func(*rejectionRoute)) []rejectionRoute {
	return mutateRoute(mutateRoute(routes, at, mutate), at+1, mutate)
}

func swapRoutes(routes []rejectionRoute, a, b int) []rejectionRoute {
	out := append([]rejectionRoute(nil), routes...)
	out[a], out[b] = out[b], out[a]
	return out
}

func countRouteEvents(routes []rejectionRoute, event string) (n int) {
	for _, route := range routes {
		if route.event == event {
			n++
		}
	}
	return n
}

// TestRejectionTopologyProbeFailSafe is AC-4: the contract's own fail-safe must not
// grade red. The three streams below differ ONLY in the probe's result and the
// routing that follows it, so the branch key — not the shape this environment happens
// to produce — is what decides which chain is owed.
func TestRejectionTopologyProbeFailSafe(t *testing.T) {
	for _, tc := range []struct {
		name       string
		stream     string
		wantBranch rejectionBranch
		wantGreen  bool
	}{
		{
			name:       "probe unavailable then fresh dispatch is the conforming fail-safe",
			stream:     claudeRejectionStream(false, false),
			wantBranch: rejectionBranchFresh,
			wantGreen:  true,
		},
		{
			name:       "probe unavailable but the FO reused anyway",
			stream:     claudeRejectionStream(false, true),
			wantBranch: rejectionBranchFresh,
			wantGreen:  false,
		},
		{
			name:       "probe reported reuse_ok and the FO reused",
			stream:     claudeRejectionStream(true, true),
			wantBranch: rejectionBranchReuse,
			wantGreen:  true,
		},
		{
			name:       "probe reported reuse_ok but the FO fresh-dispatched anyway",
			stream:     claudeRejectionStream(true, false),
			wantBranch: rejectionBranchReuse,
			wantGreen:  false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			routes, branch := claudeRejectionRoutes(tc.stream)
			if branch != tc.wantBranch {
				t.Fatalf("branch key = %q, want %q", branch, tc.wantBranch)
			}
			err := assertRejectionWorkerTopology(branch, routes)
			if tc.wantGreen && err != nil {
				t.Fatalf("conforming run graded red: %v\n%s", err, rejectionTopologySummary(routes))
			}
			if !tc.wantGreen && err == nil {
				t.Fatalf("non-conforming run graded green: %s", rejectionTopologySummary(routes))
			}
		})
	}
}

// claudeRejectionStream builds a Claude stream-json transcript for the rejection
// journey in the dialect the live runner consumes: a named background `Agent` per
// dispatch, a `task_notification` completing it, a `dispatch context-budget` Bash
// probe whose tool_result decides the branch, and a `SendMessage` carrying either the
// reuse advance or the supersede teardown. probeOK controls only the probe's stdout;
// reuse controls only how the cycle-2 rounds are routed.
func claudeRejectionStream(probeOK, reuse bool) string {
	const impl, val = "spacedock-ensign-rejection-task-implementation", "spacedock-ensign-rejection-task-validation"
	spawn := func(id, name, stage string) string {
		return fmt.Sprintf(`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Agent","id":%q,"input":{"description":"Rejection Task: %s","name":%q}}]}}`, id, stage, name)
	}
	done := func(id string) string {
		return fmt.Sprintf(`{"type":"system","subtype":"task_notification","status":"completed","tool_use_id":%q}`, id)
	}
	probe := func(id, name string) string {
		result := `dispatch context-budget: no subagent jsonl found for '` + name + `'`
		if probeOK {
			result = `{"reuse_ok":true,"usage_pct":18.2}`
		}
		return fmt.Sprintf(`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Bash","id":%q,"input":{"command":"${SPACEDOCK_BIN:-spacedock} dispatch context-budget --name %s"}}]}}`, id, name) +
			"\n" + fmt.Sprintf(`{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":%q,"content":%q}]}}`, id, result)
	}
	advance := func(name string) string {
		return fmt.Sprintf(`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":%q,"message":"Advancing to next stage: re-run for cycle 2."}}]}}`, name)
	}
	teardown := func(name string) string {
		return fmt.Sprintf(`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"SendMessage","input":{"to":%q,"message":{"type":"shutdown_request","reason":"superseded"}}}]}}`, name)
	}
	lines := []string{
		spawn("toolu_i1", impl, "implementation"), done("toolu_i1"),
		spawn("toolu_v1", val, "validation"), done("toolu_v1"),
		probe("toolu_p1", impl),
	}
	if reuse {
		lines = append(lines, advance(impl), done("toolu_i1"), probe("toolu_p2", val), advance(val), done("toolu_v1"))
	} else {
		lines = append(lines,
			teardown(impl), spawn("toolu_i2", impl, "implementation"), done("toolu_i2"),
			probe("toolu_p2", val),
			teardown(val), spawn("toolu_v2", val, "validation"), done("toolu_v2"))
	}
	return strings.Join(lines, "\n")
}

// TestRejectionCycleLineAndGateShapes covers the two remaining AC-2 shapes that live
// in durable entity state rather than in the topology: an entity carrying a second
// Cycle line, and one whose gate was never prepared. Both were fixture-determined and
// ungraded before this layer — the Cycle line was never asserted at all, and the
// gate-prepared check was wired Codex-only, so FO residual mode 1 was invisible on
// Claude and Pi.
func TestRejectionCycleLineAndGateShapes(t *testing.T) {
	entityWith := func(t *testing.T, cycles string) string {
		t.Helper()
		path := filepath.Join(t.TempDir(), "rejection-task.md")
		writeFile(t, path, "---\nstatus: validation\n---\n# Rejection Task\n\n### Feedback Cycles\n\n"+cycles)
		return path
	}
	for _, tc := range []struct {
		name   string
		cycles string
		want   string
	}{
		{"the fixture's verbatim line is the conforming end state", rejectionFeedbackCycle, ""},
		{"a passing re-validation must add no second line", rejectionFeedbackCycle + "- Cycle 2: PASSED — validation reviewer\n", "want exactly 1"},
		{"a round the FO never recorded", "", "want exactly 1"},
		{"a Cycle line the FO reworded instead of copying", "- Cycle 1: REJECTED — reviewer said no\n", "not the fixture's verbatim text"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := assertRejectionCycleLine(entityWith(t, tc.cycles))
			if tc.want == "" {
				if err != nil {
					t.Fatalf("conforming Cycle section graded red: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("non-conforming Cycle section graded green")
			}
			graded, ok := err.(*gradedErr)
			if !ok || graded.code != "rejection-cycle-line" {
				t.Fatalf("want a graded rejection-cycle-line red, got %#v", err)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("detail = %q, want it to name %q", err.Error(), tc.want)
			}
		})
	}

	// Heading-level drift, an OBSERVED FO deviation rather than a hypothetical: a
	// tip-CI run wrote the byte-exact Cycle line under `## Feedback Cycles` (H2)
	// while the workflow declares the projection under `### Feedback Cycles` (H3),
	// and the grade honestly counted zero. Captain ruled STRICT: the declared
	// grammar is the contract, so this still reds and the fixture instruction now
	// names the level instead. Accepting the line by name would be a loosening.
	driftedHeading := filepath.Join(t.TempDir(), "rejection-task.md")
	writeFile(t, driftedHeading, "---\nstatus: validation\n---\n# Rejection Task\n\n## Feedback Cycles\n\n"+rejectionFeedbackCycle)
	err := assertRejectionCycleLine(driftedHeading)
	if err == nil {
		t.Fatal("a byte-exact Cycle line under an H2 heading graded green — the declared projection is H3")
	}
	if graded, ok := err.(*gradedErr); !ok || graded.code != "rejection-cycle-line" {
		t.Fatalf("heading-drift red = %#v, want a graded rejection-cycle-line", err)
	}
	// The grade stays strict, but the DIAGNOSTIC must say what happened. A bare
	// "holds 0 entries" reads as "the FO never recorded the round", which is exactly
	// the wrong thing to tell the next reader about an FO that wrote the line
	// byte-exact one heading level off.
	for _, phrase := range []string{"byte-exact", "the round WAS recorded", "## Feedback Cycles"} {
		if !strings.Contains(err.Error(), phrase) {
			t.Errorf("heading-drift diagnostic = %q, want it to name %q", err.Error(), phrase)
		}
	}
	if strings.Contains(err.Error(), "holds 0 `- Cycle N:` entries") {
		t.Errorf("heading-drift diagnostic fell back to the bare count, implying the round was never recorded: %q", err.Error())
	}

	// An entity the FO left without ever preparing the cycle-2 gate carries no gates
	// record, which is exactly FO residual mode 1's durable signature.
	unprepared := entityWith(t, rejectionFeedbackCycle)
	if err := assertRejectionGatePrepared(unprepared); err == nil {
		t.Fatal("an entity with no prepared gate graded green")
	}
}
