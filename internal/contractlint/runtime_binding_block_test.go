// ABOUTME: Binds Codex and Pi adapter tokens to independently declared Go source values.
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
	runtimeImplHeading = "## Runtime implementation"
)

type piEmittedRuntimeToken struct {
	token, source string
}

func piEmittedRuntimeTokens(t *testing.T) []piEmittedRuntimeToken {
	t.Helper()
	wrapper := piruntime.SubagentStageDispatch("assignment", "implementation", "label")
	field, ok := reflect.TypeOf(wrapper).FieldByName("Context")
	if !ok {
		t.Fatal("piruntime.SubagentDispatch has no Context field")
	}
	contextKey := strings.Split(field.Tag.Get("json"), ",")[0]
	if contextKey == "" {
		t.Fatal("piruntime.SubagentDispatch.Context has no json tag name")
	}
	return []piEmittedRuntimeToken{
		{piruntime.TeamsDelegateAction("worker", "assignment").Action, "piruntime.TeamsDelegateAction"},
		{piruntime.TeamsDirectMessageAction("worker", "message").Action, "piruntime.TeamsDirectMessageAction"},
		{piruntime.TeamsShutdownAction("worker", "reason").Action, "piruntime.TeamsShutdownAction"},
		{piruntime.TeamsDoneAction(true).Action, "piruntime.TeamsDoneAction"},
		{fmt.Sprintf("%s: %q", contextKey, wrapper.Context), "piruntime.SubagentStageDispatch"},
	}
}

var codeSpanRe = regexp.MustCompile("`([^`]+)`")

func codeSpans(section string) map[string]bool {
	out := map[string]bool{}
	for _, m := range codeSpanRe.FindAllStringSubmatch(section, -1) {
		out[m[1]] = true
	}
	return out
}

func TestPiEmittedRuntimeTokensBindGoSource(t *testing.T) {
	tokens := piEmittedRuntimeTokens(t)
	if len(tokens) == 0 {
		t.Fatal("no Pi emitted tokens")
	}
	spans := codeSpans(markdownSectionFromText(t, readRepoFile(t, piFORuntimeRel), runtimeImplHeading))
	if len(spans) == 0 {
		t.Fatalf("%s runtime-binding block yielded no code spans", piFORuntimeRel)
	}
	for _, msg := range piTokenBindingViolations(spans, tokens) {
		t.Errorf("%s %s", piFORuntimeRel, msg)
	}
}

func piTokenBindingViolations(spans map[string]bool, tokens []piEmittedRuntimeToken) []string {
	var out []string
	for _, tok := range tokens {
		if tok.token == "" {
			out = append(out, fmt.Sprintf("%s emitted an empty token", tok.source))
			continue
		}
		if !spans[tok.token] {
			out = append(out, fmt.Sprintf("runtime-binding block has no `%s` code span, the token %s emits", tok.token, tok.source))
		}
	}
	return out
}

func TestRuntimeBindingParsersDiscriminate(t *testing.T) {
	emitted := map[string]string{"task_name": "w", "message": "m", "fork_turns": "none"}
	goodSig := "`spawn_agent(task_name,message,fork_turns=\"none\")`"
	if v := codexSpawnSignatureViolations(goodSig, emitted); len(v) != 0 {
		t.Fatalf("well-formed signature was wrongly flagged: %v", v)
	}
	for _, c := range []struct{ why, text string }{
		{"empty argument entry", "`spawn_agent(task_name,,message,fork_turns=\"none\")`"},
		{"unparseable argument characters", "`spawn_agent(task_name!!,message,fork_turns=\"none\")`"},
		{"repeated argument", "`spawn_agent(task_name,task_name,message,fork_turns=\"none\")`"},
		{"default disagreeing with the emitter", "`spawn_agent(task_name,message,fork_turns=\"all\")`"},
	} {
		if v := codexSpawnSignatureViolations(c.text, emitted); len(v) == 0 {
			t.Fatalf("%s was not flagged", c.why)
		}
	}

	tokens := []piEmittedRuntimeToken{{"delegate", "piruntime.TeamsDelegateAction"}}
	if v := piTokenBindingViolations(codeSpans("map creation to `delegate` here"), tokens); len(v) != 0 {
		t.Fatalf("`delegate` code span was wrongly flagged: %v", v)
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
			t.Fatalf("%s was not flagged", c.why)
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
