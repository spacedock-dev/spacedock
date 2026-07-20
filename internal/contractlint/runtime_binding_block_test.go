// ABOUTME: Binds the runtime capability set across the dispatch core and both host adapters,
// ABOUTME: binds Pi's emitted tokens to their Go source, and contains tool names to their sections.
package contractlint

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/piruntime"
)

var (
	codexFORuntimeRel  = filepath.Join("skills", "first-officer", "references", "codex-first-officer-runtime.md")
	piFORuntimeRel     = filepath.Join("skills", "first-officer", "references", "pi-first-officer-runtime.md")
	feedbackFlowRel    = filepath.Join("skills", "feedback-rejection-flow", "SKILL.md")
	codexEnsignRel     = filepath.Join("skills", "ensign", "references", "codex-ensign-runtime.md")
	piEnsignRel        = filepath.Join("skills", "ensign", "references", "pi-ensign-runtime.md")
	codexIdleProbeRel  = filepath.Join("docs", "dev", "codex-idle-notification-probe.md")
	runtimeImplHeading = "## Runtime implementation"
)

// capabilityHeadings reads the dispatch core's capability «fn» DEFINITION headings
// — the third, host-independent surface the adapter blocks are held to. It returns
// a slice, not a set, so a duplicated declaration survives to be reported.
func capabilityHeadings(t *testing.T) []string {
	t.Helper()
	raw, err := os.ReadFile(dispatchCorePath(t))
	if err != nil {
		t.Fatalf("read dispatch core: %v", err)
	}
	var out []string
	for _, m := range fnHeadingRe.FindAllStringSubmatch(string(raw), -1) {
		out = append(out, m[1])
	}
	return out
}

// repeatedMembers returns the members appearing more than once. Set comparison
// collapses a duplicate, so a conflicting declaration passed until this ran first.
func repeatedMembers(xs []string) []string {
	counts, out := map[string]int{}, []string{}
	for _, x := range xs {
		counts[x]++
		if counts[x] == 2 {
			out = append(out, x)
		}
	}
	sort.Strings(out)
	return out
}

// TestRuntimeCapabilitySetAgreesAcrossCoreAndAdapters binds the runtime capability
// vocabulary across three independently divergeable surfaces — the Codex FO
// runtime-binding bullets, the Pi FO runtime-binding bullets, and the dispatch
// core's capability «fn» definition headings — as the SAME set, computed from the
// core doc rather than a re-typed slice. It also holds the shared
// feedback-rejection skill to capabilities the core actually defines.
//
// Multiplicity is checked before set equality. DECLARATION ORDER IS NOT ENFORCED:
// the retired check pinned each adapter to an ordered slice, but bullet order
// carries no runtime meaning — a doc author may reorder for readability without
// any capability changing. Duplication is different: a repeated or conflicting
// declaration is a real defect set equality cannot see, so it reds explicitly.
func TestRuntimeCapabilitySetAgreesAcrossCoreAndAdapters(t *testing.T) {
	headings := capabilityHeadings(t)
	if len(headings) == 0 {
		t.Fatal("no capability «fn» definition headings found in the dispatch core — extractor bug; the binding would pass vacuously")
	}
	if dup := repeatedMembers(headings); len(dup) > 0 {
		t.Fatalf("dispatch core declares capability heading(s) more than once: %v", dup)
	}
	core := toSet(headings)

	for _, rel := range []string{codexFORuntimeRel, piFORuntimeRel} {
		t.Run(rel, func(t *testing.T) {
			text := readRepoFile(t, rel)
			if count := strings.Count(text, runtimeImplHeading+"\n"); count != 1 {
				t.Fatalf("%s has %d `%s` sections, want exactly 1", rel, count, runtimeImplHeading)
			}
			bullets := bindingBulletCapabilities(markdownSectionFromText(t, text, runtimeImplHeading))
			if len(bullets) == 0 {
				t.Fatalf("%s runtime-binding block yielded no capability bullets — extractor bug; the binding would pass vacuously", rel)
			}
			if dup := repeatedMembers(bullets); len(dup) > 0 {
				t.Fatalf("%s runtime-binding block declares capability bullet(s) more than once: %v", rel, dup)
			}
			if adapter := toSet(bullets); !setEqual(adapter, core) {
				t.Fatalf("%s runtime capability set differs from the dispatch core:\n  adapter: %v\n  core:    %v\nneither surface may add, drop, or rename a capability without the other", rel, sortedSet(adapter), sortedSet(core))
			}
		})
	}

	refs := map[string]bool{}
	for _, m := range fnRefRe.FindAllStringSubmatch(readRepoFile(t, feedbackFlowRel), -1) {
		if !metaTokens[m[1]] {
			refs[m[1]] = true
		}
	}
	if len(refs) == 0 {
		t.Fatalf("%s references no capabilities — extractor bug; the subset check would pass vacuously", feedbackFlowRel)
	}
	for name := range refs {
		if !core[name] {
			t.Errorf("%s references capability «%s», which the dispatch core does not define; the shared skill may only speak in capabilities every host adapter binds", feedbackFlowRel, name)
		}
	}
}

// piEmittedRuntimeToken is one runtime token Spacedock's Go code emits for Pi,
// paired with the source that emits it.
type piEmittedRuntimeToken struct {
	token, source string
}

// piEmittedRuntimeTokens computes the Pi runtime tokens Spacedock emits by calling
// the emitters rather than re-typing their strings. Substrate-native tokens Pi
// owns but Spacedock never emits — `subagent(`, `intercom(`, `member_spawn` — are
// deliberately absent: `member_spawn` appears in the Pi FO adapter with no
// `teams.go` counterpart, and widening the binding to swallow it would be the
// make-it-pass move this entity exists to remove. It has no owner elsewhere
// either; see UNCOVERED RUNTIME TOKENS below.
func piEmittedRuntimeTokens(t *testing.T) []piEmittedRuntimeToken {
	t.Helper()
	wrapper := piruntime.SubagentStageDispatch("assignment", "implementation", "label")
	field, ok := reflect.TypeOf(wrapper).FieldByName("Context")
	if !ok {
		t.Fatal("piruntime.SubagentDispatch has no Context field — the binding would pass vacuously")
	}
	// An absent or empty json tag yields an empty key, which would compose the
	// meaningless token `: "fresh"` and match on substring. Refuse it.
	contextKey := strings.Split(field.Tag.Get("json"), ",")[0]
	if contextKey == "" {
		t.Fatal("piruntime.SubagentDispatch.Context has no json tag name — the binding would compose an empty key")
	}
	return []piEmittedRuntimeToken{
		{piruntime.TeamsDelegateAction("worker", "assignment").Action, "piruntime.TeamsDelegateAction"},
		{piruntime.TeamsDirectMessageAction("worker", "message").Action, "piruntime.TeamsDirectMessageAction"},
		{piruntime.TeamsShutdownAction("worker", "reason").Action, "piruntime.TeamsShutdownAction"},
		{piruntime.TeamsDoneAction(true).Action, "piruntime.TeamsDoneAction"},
		{fmt.Sprintf("%s: %q", contextKey, wrapper.Context), "piruntime.SubagentStageDispatch"},
	}
}

// codeSpanRe matches one `backticked` span. Comparing whole spans rather than raw
// substrings is load-bearing: ordinary prose ("we delegate the work") satisfied a
// `strings.Contains` check for `delegate`, binding nothing.
var codeSpanRe = regexp.MustCompile("`([^`]+)`")

func codeSpans(section string) map[string]bool {
	out := map[string]bool{}
	for _, m := range codeSpanRe.FindAllStringSubmatch(section, -1) {
		out[m[1]] = true
	}
	return out
}

// TestPiEmittedRuntimeTokensBindGoSource binds every Pi runtime token Spacedock
// emits to the Pi FO adapter's runtime-binding block; either side moving alone
// reds. The wrapper's absent `acceptance` field is proven by
// piruntime.TestSubagentStageDispatchAddsOnlyPiTransportFields and the built
// artifact by internal/dispatch/build_pi_host_test.go; neither is restated here.
func TestPiEmittedRuntimeTokensBindGoSource(t *testing.T) {
	tokens := piEmittedRuntimeTokens(t)
	if len(tokens) == 0 {
		t.Fatal("no Pi emitted tokens — the binding would pass vacuously")
	}
	spans := codeSpans(markdownSectionFromText(t, readRepoFile(t, piFORuntimeRel), runtimeImplHeading))
	if len(spans) == 0 {
		t.Fatalf("%s runtime-binding block yielded no code spans — extractor bug; the binding would pass vacuously", piFORuntimeRel)
	}
	for _, msg := range piTokenBindingViolations(spans, tokens) {
		t.Errorf("%s %s", piFORuntimeRel, msg)
	}
}

// piTokenBindingViolations requires each emitted token to appear as a whole code
// span in the adapter's binding block.
func piTokenBindingViolations(spans map[string]bool, tokens []piEmittedRuntimeToken) []string {
	var out []string
	for _, tok := range tokens {
		if tok.token == "" {
			out = append(out, fmt.Sprintf("%s emitted an empty token — the binding would pass vacuously", tok.source))
			continue
		}
		if !spans[tok.token] {
			out = append(out, fmt.Sprintf("runtime-binding block has no `%s` code span, the token %s emits; the adapter must name every runtime token Spacedock actually sends", tok.token, tok.source))
		}
	}
	return out
}

// codexToolTokens and piSubstrateNativeTokens are the tool names that must stay
// inside their adapter's binding sections — absence-guard vocabularies saying
// nothing about what a tool MEANS, only where its name may appear.
//
// UNCOVERED RUNTIME TOKENS (roborev job 328 finding 4b; captain ruling 2026-07-20).
// Tokens marked UNCOVERED have NO owner anywhere in this repo: no Go emitter, no
// build fixture, no live-lane assertion. Each was verified by grepping every
// non-contractlint Go file. Their prose checks are DELETED rather than retained: a
// phrase check for a token nothing exercises proves only that we wrote the word,
// which is the fabricated rigor this entity exists to remove. 0.26.0 ships these
// gaps knowingly — the entries still bound WHERE the names may appear, but nothing
// proves WHAT they do. Follow-up: record-uncovered-runtime-tokens.
var (
	codexToolTokens = []string{
		"spawn_agent(",    // bound by TestCodexSpawnSignatureBindsToolArgs
		"send_message(",   // UNCOVERED
		"followup_task(",  // ensigncycle/shared_reviewer_reuse{,_table}_test.go
		"wait_agent(",     // ensigncycle/codex_single_run_test.go, codex_dispatch_evidence_test.go
		"list_agents(",    // UNCOVERED
		"send_input",      // ensigncycle/shared_reviewer_reuse_table_test.go (reuse transcript fixtures)
		"SendMessage",     // Claude team tool; banned here, owned by the Claude lanes
		"interrupt_agent", // UNCOVERED
	}
	piSubstrateNativeTokens = []string{
		"subagent(",                 // UNCOVERED as an assertion — appears only as PROMPT text in the pi live lanes
		"intercom(",                 // UNCOVERED
		"member_spawn",              // UNCOVERED — absent from teams.go AND from all of internal/ensigncycle
		"cwd: <resolved repo root>", // UNCOVERED
		"acceptance",                // piruntime.TestSubagentStageDispatchAddsOnlyPiTransportFields (absence)
	}
)

func piContainmentTokens(t *testing.T) []string {
	out := append([]string{}, piSubstrateNativeTokens...)
	for _, tok := range piEmittedRuntimeTokens(t) {
		out = append(out, tok.token)
	}
	return out
}

func toolTokenContainmentViolations(span string, tokens []string) []string {
	var out []string
	for _, token := range tokens {
		if strings.Contains(span, token) {
			out = append(out, fmt.Sprintf("carries runtime tool token %q outside its binding section", token))
		}
	}
	return out
}

// TestRuntimeToolTokensStayInBindingSections is the section-scoping containment
// guard: host tool names live only where the adapter binds them, so prose
// elsewhere stays readable by a cold agent on any host. It asserts where a name
// may appear, never what the tool does — the Codex spawn call is bound by
// TestCodexSpawnSignatureBindsToolArgs and the Pi tokens by
// TestPiEmittedRuntimeTokensBindGoSource. Several substrate-native tools have no
// owner at all; the per-token annotations on the vocabularies below say which.
func TestRuntimeToolTokensStayInBindingSections(t *testing.T) {
	codexText := readRepoFile(t, codexFORuntimeRel)
	_, outsideProbe := extractMarkdownSection(t, codexText, "## Live Tool Surface Probe")
	_, outsideRuntime := extractMarkdownSection(t, outsideProbe, runtimeImplHeading)
	waitSection, outsideWait := extractMarkdownSection(t, outsideRuntime, "## Codex wait notes")

	for rel, span := range map[string]string{
		codexFORuntimeRel: outsideWait,
		codexEnsignRel:    readRepoFile(t, codexEnsignRel),
		feedbackFlowRel:   readRepoFile(t, feedbackFlowRel),
		codexIdleProbeRel: readRepoFile(t, codexIdleProbeRel),
	} {
		for _, msg := range toolTokenContainmentViolations(span, codexToolTokens) {
			t.Errorf("%s %s", rel, msg)
		}
	}

	// The wait notes are the one section allowed to name a concrete tool, and only
	// `wait_agent` — the tool the notes exist to describe.
	waitAllowed := []string{}
	for _, token := range codexToolTokens {
		if token != "wait_agent(" {
			waitAllowed = append(waitAllowed, token)
		}
	}
	for _, msg := range toolTokenContainmentViolations(waitSection, waitAllowed) {
		t.Errorf("%s Codex wait notes may name only wait_agent as a concrete runtime tool: %s", codexFORuntimeRel, msg)
	}

	piText := readRepoFile(t, piFORuntimeRel)
	outsidePi := piText
	for _, heading := range []string{runtimeImplHeading, "## Live Harness Isolation"} {
		_, outsidePi = extractMarkdownSection(t, outsidePi, heading)
	}
	for _, msg := range toolTokenContainmentViolations(outsidePi, piContainmentTokens(t)) {
		t.Errorf("%s %s", piFORuntimeRel, msg)
	}
}

// TestToolTokenContainmentGuardDiscriminates is the non-vacuity control:
// capability-only prose PASSES, while a Codex tool call, a Pi emitted token, and a
// Pi substrate-native token each RED in a span that must not carry them.
func TestToolTokenContainmentGuardDiscriminates(t *testing.T) {
	piTokens := piContainmentTokens(t)
	pass := []struct {
		why, span string
		tokens    []string
	}{
		{"capability-only Codex prose", "Route the follow-up through `«addressable-worker»`.\n", codexToolTokens},
		{"capability-only Pi prose", "Teardown maps to `«worker.shutdown»`.\n", piTokens},
	}
	for _, c := range pass {
		if v := toolTokenContainmentViolations(c.span, c.tokens); len(v) != 0 {
			t.Fatalf("control: the %s was wrongly flagged: %v", c.why, v)
		}
	}
	red := []struct {
		why, span string
		tokens    []string
	}{
		{"leaked Codex spawn call", "For every ready entity call spawn_agent(task_name,message).\n", codexToolTokens},
		{"leaked Claude team tool", "Reply with SendMessage to the first officer.\n", codexToolTokens},
		{"leaked Pi emitted action", "Follow-up steering maps to message_dm.\n", piTokens},
		{"leaked Pi substrate-native tool", "Call subagent(...) with a fresh context.\n", piTokens},
	}
	for _, c := range red {
		if v := toolTokenContainmentViolations(c.span, c.tokens); len(v) == 0 {
			t.Fatalf("control: the %s was not flagged — the guard stopped biting", c.why)
		}
	}
}

// TestRuntimeBindingGuardsDiscriminate is the non-vacuity control for the three
// bindings' extractors. Every RED case below is an input the earlier extractors
// ACCEPTED: unparseable signature text silently skipped, a duplicate declaration
// collapsed by set comparison, and a token matched as a substring of ordinary
// prose. A regression restoring any of those tolerances fails here.
func TestRuntimeBindingGuardsDiscriminate(t *testing.T) {
	emitted := map[string]string{"task_name": "w", "message": "m", "fork_turns": "none"}
	goodSig := "`spawn_agent(task_name,message,fork_turns=\"none\")`"
	if v := codexSpawnSignatureViolations(goodSig, emitted); len(v) != 0 {
		t.Fatalf("control: the well-formed signature was wrongly flagged: %v", v)
	}
	for _, c := range []struct{ why, text string }{
		{"empty argument entry", "`spawn_agent(task_name,,message,fork_turns=\"none\")`"},
		{"unparseable argument characters", "`spawn_agent(task_name!!,message,fork_turns=\"none\")`"},
		{"repeated argument", "`spawn_agent(task_name,task_name,message,fork_turns=\"none\")`"},
		{"default disagreeing with the emitter", "`spawn_agent(task_name,message,fork_turns=\"all\")`"},
	} {
		if v := codexSpawnSignatureViolations(c.text, emitted); len(v) == 0 {
			t.Fatalf("control: the %s was not flagged — the extractor skipped what it could not parse", c.why)
		}
	}

	if dup := repeatedMembers([]string{"worker.spawn", "context-budget"}); len(dup) != 0 {
		t.Fatalf("control: distinct members were reported as duplicates: %v", dup)
	}
	if dup := repeatedMembers([]string{"worker.spawn", "context-budget", "context-budget"}); len(dup) == 0 {
		t.Fatal("control: a duplicated capability declaration was not reported — set comparison alone would collapse it")
	}

	tokens := []piEmittedRuntimeToken{{"delegate", "piruntime.TeamsDelegateAction"}}
	if v := piTokenBindingViolations(codeSpans("map creation to `delegate` here"), tokens); len(v) != 0 {
		t.Fatalf("control: the `delegate` code span was wrongly flagged: %v", v)
	}
	for _, c := range []struct {
		why    string
		spans  map[string]bool
		tokens []piEmittedRuntimeToken
	}{
		{"token present only as ordinary prose", codeSpans("we delegate the work to a worker"), tokens},
		{"token present only inside a larger code span", codeSpans("call `delegate_v2(...)` instead"), tokens},
		{"empty token from a missing json tag", codeSpans("carries `: \"fresh\"` here"), []piEmittedRuntimeToken{{"", "piruntime.SubagentStageDispatch"}}},
	} {
		if v := piTokenBindingViolations(c.spans, c.tokens); len(v) == 0 {
			t.Fatalf("control: the %s was not flagged — substring matching would have bound it", c.why)
		}
	}
}

func readRepoFile(t *testing.T, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot(t), rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(data)
}

func markdownSectionFromText(t *testing.T, text, heading string) string {
	t.Helper()
	section, _ := extractMarkdownSection(t, text, heading)
	return section
}

func bindingBulletCapabilities(section string) []string {
	re := regexp.MustCompile(`(?m)^- ` + "`" + `«([^»]+)»` + "`" + `:`)
	out := []string{}
	for _, match := range re.FindAllStringSubmatch(section, -1) {
		out = append(out, match[1])
	}
	return out
}
