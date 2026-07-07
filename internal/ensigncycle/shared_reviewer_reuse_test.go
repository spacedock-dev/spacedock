package ensigncycle

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// The host-specific reviewer-reuse assertions for the rejection-flow scenario.
// They are the producer signal that the FO reused the kept-alive validation
// reviewer for the cycle-2 re-review rather than fresh-dispatching (the #141
// keepalive contract the Go port dropped): Claude reuses via a SendMessage tool
// call, Codex multi_agent_v2 via `followup_task`, and legacy Codex fixtures via
// `send_input`. They live under the DEFAULT build tags
// (parsing only stdlib JSON) so the offline table tests exercise them without
// spending a model, alongside the //go:build live runners that feed them real
// transcripts.

// claudeAgentIDResult extracts the `agentId: <id>` a completed subagent's
// tool_result returns ("agentId: a94abe89c85f9f4cc (use SendMessage with to: …)").
// The teams-mode reuse handle for a kept-alive completed agent.
var claudeAgentIDResult = regexp.MustCompile(`agentId:\s*(a[0-9a-f]+)`)

// assertClaudeReviewerReuse scans the stream-json transcript for the FO reusing
// the kept-alive cycle-1 validation reviewer for the cycle-2 re-review — the
// durable producer signal it did NOT fresh-dispatch the cycle-2 validator (the
// #141 keepalive contract). It enforces BOTH halves of "genuine reuse", because
// either alone false-passes a fresh-dispatch (the cycle-8 audit hole):
//
//  1. NO fresh cycle-2 validation spawn. The rejection-flow drives ONE cycle-1
//     validation; a reuse run re-engages that same reviewer, so there is EXACTLY
//     ONE validation-stage Agent/Task spawn. A run that fresh-dispatches the
//     cycle-2 validator emits a SECOND validation spawn — the forbidden behavior —
//     so >1 validation spawn is an immediate FAIL regardless of any reuse-shaped
//     message. This is the discriminator the bare name/agentId match lacked: a
//     fresh cycle-2 reviewer messaged by its validation NAME looks identical to a
//     reused one EXCEPT for the extra spawn.
//  2. A message to the cycle-1 reviewer's handle. With exactly one validation
//     reviewer in play, a SendMessage to it — by the agentId its spawn returned
//     (teams-mode kept-alive resume) OR by its spawn NAME (`…-validation`, the
//     production shape the recorded testdata uses) — is the re-review. Narration
//     or a message to a non-validation recipient does not count.
func assertClaudeReviewerReuse(stream string) error {
	// First pass: count validation-stage Agent/Task spawns and collect the agentIds
	// those spawns returned. A reused reviewer drives one validation spawn; a fresh
	// cycle-2 dispatch drives a second.
	validationSpawnIDs := map[string]bool{} // tool_use_id of each validation-stage spawn
	validationSpawnCount := 0
	validationAgentIDs := map[string]bool{} // agentIds those spawns returned (reuse handles)
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Message *struct {
				Content []struct {
					Type      string `json:"type"`
					Name      string `json:"name"`
					ID        string `json:"id"`
					ToolUseID string `json:"tool_use_id"`
					Input     struct {
						Description string `json:"description"`
					} `json:"input"`
					Content json.RawMessage `json:"content"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if block.Type == "tool_use" && (block.Name == "Agent" || block.Name == "Task") &&
				strings.Contains(strings.ToLower(block.Input.Description), "validation") {
				validationSpawnCount++
				validationSpawnIDs[block.ID] = true
			}
			if block.Type == "tool_result" && validationSpawnIDs[block.ToolUseID] {
				if m := claudeAgentIDResult.FindStringSubmatch(string(block.Content)); m != nil {
					validationAgentIDs[m[1]] = true
				}
			}
		}
	}

	if validationSpawnCount == 0 {
		return fmt.Errorf("no validation-stage Agent/Task spawn found — the FO never created a cycle-1 reviewer to reuse")
	}
	if validationSpawnCount > 1 {
		return fmt.Errorf("the FO emitted %d validation-stage Agent/Task spawns — it FRESH-dispatched the cycle-2 validator instead of reusing the kept-alive cycle-1 reviewer (the #141 keepalive contract violation)", validationSpawnCount)
	}

	// Second pass: with exactly one validation reviewer, a SendMessage to it (by its
	// returned agentId, or by a name naming the validation stage) is the cycle-2
	// re-review reuse signal.
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Message *struct {
				Content []struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Input struct {
						To string `json:"to"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if block.Type != "tool_use" || block.Name != "SendMessage" {
				continue
			}
			to := block.Input.To
			if strings.Contains(strings.ToLower(to), "validation") || validationAgentIDs[to] {
				return nil
			}
		}
	}
	return fmt.Errorf("the FO spawned exactly one validation reviewer but sent it no reuse SendMessage (by name or agentId) for the cycle-2 re-review")
}

// assertClaudeSingleEntityRejectionFlow is the single-entity (`-p`) Claude
// producer-signal assertion for the rejection-flow scenario. The Claude runner
// launches `spacedock claude -- -p {prompt}` with a prompt naming one entity. The
// `-p` FO drives EITHER bare OR team mode (it opts into the background back-channel
// when SendMessage is exposed), so the contract admits two valid end-states; the
// invariant across both is that the cycle-2 re-review reaches a VALIDATION worker —
// fresh-dispatched or the kept-alive reviewer reused — and NEVER the implementation
// worker serving as its own validator. It enforces two checks:
//
//  1. The cycle-2 re-review reaches a validation worker. In bare mode that is two
//     distinct validation-stage Agent/Task spawns (a fresh cycle-1 reviewer AND a
//     fresh cycle-2 reviewer). In team mode it is one validation spawn plus a reuse
//     SendMessage to that kept-alive `…-validation` reviewer. A run with neither
//     shape either never re-reviewed or collapsed the cycle-2 re-review onto a
//     non-validation worker — both forbidden.
//  2. The cycle-2 re-review is NOT routed to an implementation worker. A SendMessage
//     to an `…-implementation` handle telling it to validate / re-review is the
//     impl-as-validator violation. A `shutdown_request` to that handle is exempt —
//     it is the contract-mandated supersede teardown of the prior cycle's worker
//     (the FO reaps the superseded fix worker by exactly this message), not a
//     re-review routing, so the body, not just the recipient, decides the verdict.
func assertClaudeSingleEntityRejectionFlow(stream string) error {
	validationSpawnCount := 0
	reviewerReuse := false
	for _, line := range strings.Split(stream, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var entry struct {
			Message *struct {
				Content []struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Input struct {
						Description string          `json:"description"`
						To          string          `json:"to"`
						Message     json.RawMessage `json:"message"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err != nil || entry.Message == nil {
			continue
		}
		for _, block := range entry.Message.Content {
			if block.Type != "tool_use" {
				continue
			}
			desc := strings.ToLower(block.Input.Description)
			if (block.Name == "Agent" || block.Name == "Task") && strings.Contains(desc, "validation") {
				validationSpawnCount++
			}
			if block.Name != "SendMessage" {
				continue
			}
			to := strings.ToLower(block.Input.To)
			// A SendMessage to the kept-alive validation reviewer is the team-mode
			// cycle-2 re-review reuse signal (the addressable-worker reuse path the
			// `-p` run legitimately takes when it drives team mode, not bare).
			if strings.Contains(to, "validation") {
				reviewerReuse = true
			}
			// The impl-as-validator violation: a SendMessage telling an
			// implementation-named worker to validate / re-review. A
			// shutdown_request to that handle is the contract-mandated supersede
			// teardown of the prior cycle's worker, NOT a re-review routing — exempt
			// it; only a non-teardown instruction to the fix worker is the violation.
			if strings.Contains(to, "implementation") && !isShutdownRequest(block.Input.Message) {
				return fmt.Errorf("the cycle-2 re-review was routed to an implementation worker (%q) — the fix agent and the reviewer must be SEPARATE dispatches; the impl worker must never serve as its own validator", block.Input.To)
			}
		}
	}
	// Two contract-valid end-states for the `-p` run, because it may drive bare OR
	// team mode: bare fresh-dispatches a distinct reviewer per cycle (>= 2 validation
	// spawns); team keeps the cycle-1 reviewer alive and reuses it for the cycle-2
	// re-review (one validation spawn + a reuse message to that reviewer).
	if validationSpawnCount >= 2 {
		return nil
	}
	if validationSpawnCount >= 1 && reviewerReuse {
		return nil
	}
	return fmt.Errorf("single-entity rejection-flow produced %d validation-stage spawns with reviewerReuse=%v, want either >= 2 fresh validation spawns (bare mode) or >= 1 spawn plus a reuse message to the kept-alive validation reviewer (team mode) — the cycle-2 re-review must reach a validation worker, fresh or reused", validationSpawnCount, reviewerReuse)
}

// isShutdownRequest reports whether a SendMessage `message` payload is a
// cooperative teardown (`{"type":"shutdown_request",...}`), sent either as a JSON
// object or as a JSON string carrying that type. The supersede-shutdown contract
// requires the FO to reap a superseded worker by this message, so a shutdown_request
// to an implementation handle is mandatory teardown, not a re-review routing.
func isShutdownRequest(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var obj struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil && obj.Type == "shutdown_request" {
		return true
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		var inner struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal([]byte(s), &inner); err == nil && inner.Type == "shutdown_request" {
			return true
		}
		return strings.Contains(s, "shutdown_request")
	}
	return false
}

// codexCollabItem is one `codex exec --json` stream item. Codex surfaces its
// multi-agent calls as `collab_tool_call` items (tool = spawn_agent / followup_task /
// send_input for legacy pre-v2 fixtures / wait / close_agent); the worker is addressed by opaque `receiver_thread_ids`,
// not by a name. The validation reviewer's thread is bound by the spawn_agent whose
// prompt dispatches the validation stage (its completed item carries the spawned
// thread id in receiver_thread_ids).
type codexCollabItem struct {
	Type string `json:"type"`
	Item struct {
		Type              string   `json:"type"`
		Tool              string   `json:"tool"`
		ReceiverThreadIDs []string `json:"receiver_thread_ids"`
		Prompt            string   `json:"prompt"`
	} `json:"item"`
}

// assertCodexReviewerReuse scans the `codex exec --json` transcript for the FO's
// rejection-flow reviewer routing. When the transcript proves a turn-starting
// `«addressable-worker»` route exists, the FO must reuse the kept-alive cycle-1
// validation reviewer for the cycle-2 re-review — the durable producer signal it
// did NOT fresh-dispatch the cycle-2 validator (the #141 keepalive contract).
// Current public Codex live surfaces may expose only spawn/wait; when the FO
// explicitly observes that no turn-starting reuse route is exposed, the contract is
// characterized as addressable-worker ABSENT and the cycle-2 reviewer is fresh.
// Codex exec does not surface enough current multi_agent_v2 metadata to prove a distinct reviewer process.
// When spawn/thread evidence is absent, the Codex lane
// grades assignment-separation plus durable end-state: this assertion requires
// distinct implementation and validation dispatch-build surfaces, and
// assertRejectionFlow grades the two-cycle entity body. The Claude runner remains
// the process-level separation oracle because its stream exposes Agent/SendMessage
// evidence.
// Codex multi_agent_v2 reuses via a `followup_task` collab_tool_call to the
// reviewer's thread; legacy pre-v2 fixtures use `send_input`. The thread is bound
// by the `spawn_agent` whose prompt dispatched the validation stage. In the PRESENT
// branch it enforces BOTH halves of "genuine reuse", because binding ANY
// validation-prompt spawn alone false-passes a fresh-dispatch (the cycle-8 M2 hole):
//
//  1. NO fresh cycle-2 validation spawn. The rejection-flow drives ONE cycle-1
//     validation spawn_agent; a reuse run re-engages that same thread. A run that
//     fresh-spawns the cycle-2 validator emits a SECOND validation spawn_agent — the
//     forbidden behavior — so >1 validation spawn_agent is an immediate FAIL, even
//     though a follow-up to the fresh thread otherwise looks like reuse.
//  2. A turn-triggering reuse call to the cycle-1 reviewer's bound thread. With
//     exactly one validation spawn, a follow-up to its thread is the re-review —
//     not feedback-to-implementation routing (which targets the impl thread), not
//     loose narration.
func assertCodexReviewerReuse(jsonl string) error {
	if codexAddressableWorkerAbsent(jsonl) {
		return assertCodexFreshValidationWhenAddressableAbsent(jsonl)
	}

	// First pass: count validation spawn_agents and bind their thread id(s). A
	// reused reviewer drives one validation spawn; a fresh cycle-2 dispatch drives a
	// second.
	validationThreads := map[string]bool{}
	validationSpawnCount := 0
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev codexCollabItem
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		it := ev.Item
		// spawn_agent surfaces as item.started (no threads yet) AND item.completed
		// (with the spawned receiver_thread_ids); count the completed one so the
		// thread binding and the count stay in lockstep (one spawn = one count).
		if it.Type != "collab_tool_call" || it.Tool != "spawn_agent" || ev.Type != "item.completed" {
			continue
		}
		if codexDispatchesValidation(it.Prompt) {
			validationSpawnCount++
			for _, tid := range it.ReceiverThreadIDs {
				validationThreads[tid] = true
			}
		}
	}

	if validationSpawnCount == 0 {
		if codexReviewerReuseViaAssignmentSurfaces(jsonl) {
			return nil
		}
		return fmt.Errorf("Codex reviewer assignment-separation evidence missing validation-stage dispatch assignment")
	}
	if validationSpawnCount > 1 {
		return fmt.Errorf("the FO emitted %d validation spawn_agents — it FRESH-dispatched the cycle-2 validator instead of reusing the kept-alive cycle-1 reviewer (the #141 keepalive contract violation)", validationSpawnCount)
	}

	// Second pass: with exactly one validation reviewer, a turn-triggering reuse
	// call to its bound thread is the cycle-2 re-review reuse signal.
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev codexCollabItem
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		it := ev.Item
		if it.Type != "collab_tool_call" || !codexReviewerReuseTool(it.Tool) {
			continue
		}
		for _, tid := range it.ReceiverThreadIDs {
			if validationThreads[tid] {
				return nil
			}
		}
	}
	return fmt.Errorf("the FO spawned exactly one validation reviewer but sent it no followup_task/send_input for the cycle-2 re-review")
}

func codexReviewerReuseViaAssignmentSurfaces(jsonl string) bool {
	initialValidationBuilds := 0
	validationAdvanceBuilds := 0
	implementationAdvances := 0
	waitCalls := 0
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev struct {
			Type string `json:"type"`
			Item struct {
				Type    string `json:"type"`
				Tool    string `json:"tool"`
				Command string `json:"command"`
			} `json:"item"`
		}
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		switch ev.Item.Type {
		case "collab_tool_call":
			if ev.Item.Tool == "wait" || ev.Item.Tool == "wait_agent" || ev.Item.Tool == "collab:wait" {
				waitCalls++
			}
		case "command_execution":
			if ev.Type != "item.completed" {
				continue
			}
			build, ok := codexDispatchBuildFromCommand(ev.Item.Command)
			if !ok {
				continue
			}
			switch build.stage {
			case "validation":
				if build.advance {
					validationAdvanceBuilds++
				} else {
					initialValidationBuilds++
				}
			}
			if build.stage == "implementation" {
				if build.advance {
					implementationAdvances++
				}
			}
		}
	}
	return initialValidationBuilds >= 1 &&
		validationAdvanceBuilds >= 1 &&
		implementationAdvances >= 1 &&
		waitCalls >= 2
}

type codexDispatchBuild struct {
	stage   string
	advance bool
}

func codexDispatchBuildFromCommand(command string) (codexDispatchBuild, bool) {
	lower := strings.ToLower(command)
	if !strings.Contains(lower, "dispatch build") {
		return codexDispatchBuild{}, false
	}
	build := codexDispatchBuild{
		stage:   codexDispatchStageFromFlags(lower),
		advance: strings.Contains(lower, "--advance"),
	}
	for _, obj := range jsonObjectsIn(command) {
		var payload struct {
			Stage   string `json:"stage"`
			Advance bool   `json:"advance"`
		}
		if err := json.Unmarshal([]byte(obj), &payload); err != nil {
			continue
		}
		if payload.Stage != "" {
			build.stage = strings.ToLower(payload.Stage)
		}
		if payload.Advance {
			build.advance = true
		}
	}
	if build.stage == "" {
		return codexDispatchBuild{}, false
	}
	return build, true
}

func codexDispatchStageFromFlags(command string) string {
	fields := strings.Fields(command)
	for i, field := range fields {
		if field == "--stage" && i+1 < len(fields) {
			return strings.Trim(fields[i+1], `'"`)
		}
		if strings.HasPrefix(field, "--stage=") {
			return strings.Trim(strings.TrimPrefix(field, "--stage="), `'"`)
		}
	}
	return ""
}

func jsonObjectsIn(s string) []string {
	var objects []string
	for start := 0; start < len(s); start++ {
		if s[start] != '{' {
			continue
		}
		depth := 0
		inString := false
		escaped := false
		for end := start; end < len(s); end++ {
			ch := s[end]
			if inString {
				if escaped {
					escaped = false
					continue
				}
				if ch == '\\' {
					escaped = true
					continue
				}
				if ch == '"' {
					inString = false
				}
				continue
			}
			switch ch {
			case '"':
				inString = true
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					objects = append(objects, s[start:end+1])
					start = end
					end = len(s)
				}
			}
		}
	}
	return objects
}

func codexReviewerReuseTool(tool string) bool {
	return tool == "followup_task" || tool == "send_input"
}

// codexAddressableWorkerAbsent reports whether the FO observed that the live
// Codex surface exposes no turn-starting reuse route for a completed worker, so a
// fresh cycle-2 reviewer is the contract-correct choice. The live FO words this
// observation freely ("no follow-up/send binding for worker reuse", "the cycle-1
// reviewer is not addressable on this host", "no followup_task or equivalent
// turn-starting reuse route"), so a fixed phrase list false-RED-flags a
// contract-correct fresh-dispatch. The detector instead reads the FO's narration
// (agent_message text only — never command output, which can echo the runtime
// doc) and treats a message that negates a reuse-route concept as the absence
// observation. A reuse tool actually appearing in the transcript still overrides
// this via assertCodexFreshValidationWhenAddressableAbsent.
func codexAddressableWorkerAbsent(jsonl string) bool {
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev struct {
			Item struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"item"`
		}
		if err := json.Unmarshal([]byte(line), &ev); err != nil || ev.Item.Type != "agent_message" {
			continue
		}
		if codexNarrationNegatesReuseRoute(ev.Item.Text) {
			return true
		}
	}
	return false
}

// codexNarrationNegatesReuseRoute reports whether a single FO narration message
// states that the reuse/follow-up route is absent. The negation must bind to the
// route, binding, tool, support, or addressability claim; affirmative reuse
// narration can contain unrelated negation ("not doing validation") without
// choosing the addressable-worker-ABSENT branch.
func codexNarrationNegatesReuseRoute(text string) bool {
	lower := strings.ToLower(text)
	for _, pattern := range codexReuseRouteAbsencePatterns {
		if pattern.MatchString(lower) {
			return true
		}
	}
	return false
}

var codexReuseRouteAbsencePatterns = []*regexp.Regexp{
	regexp.MustCompile(`\bno\s+(?:turn-starting\s+|completed-worker\s+|addressable[- ]worker\s+|addressable\s+)?(?:follow[- ]?up(?:_task)?|followup_task|send[-_ ]?(?:message|input)|reuse|addressable)(?:[/_a-z0-9 -]{0,80})?\b(?:binding|route|tool|call)(?:\s+exposed)?\b`),
	regexp.MustCompile(`\b(?:do|does|did)\s+not\s+have\s+(?:a\s+|an\s+)?(?:follow[- ]?up(?:_task)?|followup_task|addressable|reuse)(?:[/_a-z0-9 -]{0,80})?\b(?:binding|route|tool|call)\b`),
	regexp.MustCompile(`\b(?:do|does|did)\s+not\s+expose\s+(?:a\s+|an\s+)?(?:turn-starting\s+)?(?:addressable[- ]worker|addressable\s+reviewer|follow[- ]?up(?:_task)?|followup_task|reuse)(?:[/_a-z0-9 -]{0,80})?\b(?:binding|route|tool|call)\b`),
	regexp.MustCompile(`\breviewer\s+reuse\s+is\s+not\s+supported\b`),
	regexp.MustCompile(`\breuse\s+is\s+not\s+supported\b`),
	regexp.MustCompile(`\b(?:reviewer|cycle-1 reviewer|worker)\s+is\s+not\s+addressable\b`),
	regexp.MustCompile(`\bcannot\s+address\s+(?:the\s+)?(?:kept-alive\s+)?reviewer\b`),
}

func assertCodexFreshValidationWhenAddressableAbsent(jsonl string) error {
	validationSpawnCount := 0
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev codexCollabItem
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		it := ev.Item
		if it.Type != "collab_tool_call" {
			continue
		}
		if codexReviewerReuseTool(it.Tool) {
			return fmt.Errorf("Codex addressable-worker was characterized ABSENT, but transcript contains turn-starting reuse tool %q", it.Tool)
		}
		if it.Tool == "spawn_agent" && ev.Type == "item.completed" && codexDispatchesValidation(it.Prompt) {
			validationSpawnCount++
		}
	}
	if validationSpawnCount < 2 {
		return fmt.Errorf("Codex addressable-worker ABSENT rejection-flow produced %d validation spawn_agents, want >= 2 fresh validation reviewers for cycle 1 and cycle 2", validationSpawnCount)
	}
	return nil
}

// codexDispatchesValidation reports whether a spawn_agent prompt dispatched the
// validation stage — either the validation dispatch file
// (spacedock-ensign-{slug}-validation.md) or a prompt naming the validation stage.
func codexDispatchesValidation(prompt string) bool {
	lower := strings.ToLower(prompt)
	return strings.Contains(lower, "-validation.md") || strings.Contains(lower, "validation")
}
