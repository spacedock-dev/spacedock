// ABOUTME: Binds the runtime capability set across the dispatch core and both host adapters,
// ABOUTME: binds Pi's emitted tokens to their Go source, and contains tool names to their sections.
package contractlint

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
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

// capabilitySetFromHeadings reads the dispatch core's capability «fn» DEFINITION
// headings — the third, host-independent surface the adapter blocks are held to.
func capabilitySetFromHeadings(t *testing.T) map[string]bool {
	t.Helper()
	raw, err := os.ReadFile(dispatchCorePath(t))
	if err != nil {
		t.Fatalf("read dispatch core: %v", err)
	}
	out := map[string]bool{}
	for _, m := range fnHeadingRe.FindAllStringSubmatch(string(raw), -1) {
		out[m[1]] = true
	}
	return out
}

// TestRuntimeCapabilitySetAgreesAcrossCoreAndAdapters binds the runtime capability
// vocabulary across three independently divergeable surfaces — the Codex FO
// runtime-binding bullets, the Pi FO runtime-binding bullets, and the dispatch
// core's capability «fn» definition headings — as the SAME set. It replaces a
// hardcoded expected-capability slice: the expectation is computed from the core
// doc, so no re-typed literal can drift from all three. It also holds the shared
// feedback-rejection skill to capabilities the core actually defines, which is
// what keeps that skill's `«token»` vocabulary from being decorative.
func TestRuntimeCapabilitySetAgreesAcrossCoreAndAdapters(t *testing.T) {
	core := capabilitySetFromHeadings(t)
	if len(core) == 0 {
		t.Fatal("no capability «fn» definition headings found in the dispatch core — extractor bug; the binding would pass vacuously")
	}

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
// `teams.go` counterpart, and its runtime meaning belongs to the gated live lane
// in internal/ensigncycle/pi_live_runner_test.go, not to a binding widened to
// swallow it.
func piEmittedRuntimeTokens(t *testing.T) []piEmittedRuntimeToken {
	t.Helper()
	wrapper := piruntime.SubagentStageDispatch("assignment", "implementation", "label")
	field, ok := reflect.TypeOf(wrapper).FieldByName("Context")
	if !ok {
		t.Fatal("piruntime.SubagentDispatch has no Context field — the binding would pass vacuously")
	}
	contextKey := strings.Split(field.Tag.Get("json"), ",")[0]
	return []piEmittedRuntimeToken{
		{piruntime.TeamsDelegateAction("worker", "assignment").Action, "piruntime.TeamsDelegateAction"},
		{piruntime.TeamsDirectMessageAction("worker", "message").Action, "piruntime.TeamsDirectMessageAction"},
		{piruntime.TeamsShutdownAction("worker", "reason").Action, "piruntime.TeamsShutdownAction"},
		{piruntime.TeamsDoneAction(true).Action, "piruntime.TeamsDoneAction"},
		{fmt.Sprintf("%s: %q", contextKey, wrapper.Context), "piruntime.SubagentStageDispatch"},
	}
}

// TestPiEmittedRuntimeTokensBindGoSource binds every Pi runtime token Spacedock
// emits to the Pi FO adapter's runtime-binding block. Renaming a `TeamsAction`
// action string or changing `SubagentStageDispatch`'s context value reds against
// an unchanged doc; dropping the token from the doc reds against unchanged Go. The
// wrapper's absent `acceptance` field is proven by
// piruntime.TestSubagentStageDispatchAddsOnlyPiTransportFields and the built
// artifact by internal/dispatch/build_pi_host_test.go; neither is restated here.
func TestPiEmittedRuntimeTokensBindGoSource(t *testing.T) {
	tokens := piEmittedRuntimeTokens(t)
	if len(tokens) == 0 {
		t.Fatal("no Pi emitted tokens — the binding would pass vacuously")
	}
	section := markdownSectionFromText(t, readRepoFile(t, piFORuntimeRel), runtimeImplHeading)
	for _, tok := range tokens {
		if tok.token == "" {
			t.Fatalf("%s emitted an empty token — the binding would pass vacuously", tok.source)
		}
		if !strings.Contains(section, tok.token) {
			t.Errorf("%s runtime-binding block does not carry %q, the token %s emits; the adapter must name every runtime token Spacedock actually sends", piFORuntimeRel, tok.token, tok.source)
		}
	}
}

// codexToolTokens are the Codex tool names that must stay inside the Codex FO
// adapter's probe and runtime-binding sections. piSubstrateNativeTokens are the Pi
// tokens Spacedock does not emit; with the emitted tokens they form the Pi
// containment set. Both are absence-guard vocabularies: they say nothing about
// what a tool MEANS, only where its name may appear.
var (
	codexToolTokens = []string{
		"spawn_agent(",
		"send_message(",
		"followup_task(",
		"wait_agent(",
		"list_agents(",
		"send_input",
		"SendMessage",
		"interrupt_agent",
	}
	piSubstrateNativeTokens = []string{
		"subagent(",
		"intercom(",
		"member_spawn",
		"cwd: <resolved repo root>",
		"acceptance",
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
// TestCodexSpawnSignatureBindsToolArgs, the Pi tokens by
// TestPiEmittedRuntimeTokensBindGoSource, and the substrate-native tools by the
// gated live lanes in internal/ensigncycle.
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
// Pi substrate-native token each RED in a span that must not carry them. An
// emptied vocabulary or a predicate that stopped comparing fails here rather than
// waving through an adapter that has leaked tool names into host-neutral prose.
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
