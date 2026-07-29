// ABOUTME: Binds the Codex spawn signature in the FO adapter to the Go emitter's argument shape.
package contractlint

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/dispatch"
)

var (
	codexSpawnCallRe = regexp.MustCompile(`spawn_agent\(([^)]*)\)`)
	codexSpawnArgRe  = regexp.MustCompile(`^([a-z_]+)(?:="([^"]*)")?$`)
)

type codexSpawnArg struct {
	name, value string
	hasDefault  bool
}

func codexSpawnArgs(argList string) ([]codexSpawnArg, error) {
	var out []codexSpawnArg
	for _, entry := range strings.Split(argList, ",") {
		trimmed := strings.TrimSpace(entry)
		m := codexSpawnArgRe.FindStringSubmatch(trimmed)
		if m == nil {
			return nil, fmt.Errorf("unparseable argument %q", trimmed)
		}
		out = append(out, codexSpawnArg{name: m[1], value: m[2], hasDefault: strings.Contains(trimmed, "=")})
	}
	return out, nil
}

func codexSpawnSignatureViolations(text string, toolArgs map[string]string) []string {
	calls := codexSpawnCallRe.FindAllStringSubmatch(text, -1)
	if len(calls) == 0 {
		return []string{"no `spawn_agent(...)` signature found"}
	}
	wanted := map[string]bool{}
	for name := range toolArgs {
		wanted[name] = true
	}
	var out []string
	for _, call := range calls {
		args, err := codexSpawnArgs(call[1])
		if err != nil {
			out = append(out, fmt.Sprintf("signature %q is malformed: %v", call[0], err))
			continue
		}
		names := map[string]bool{}
		for _, arg := range args {
			if names[arg.name] {
				out = append(out, fmt.Sprintf("signature %q declares argument %q twice", call[0], arg.name))
			}
			names[arg.name] = true
			if emitted, ok := toolArgs[arg.name]; arg.hasDefault && ok && emitted != arg.value {
				out = append(out, fmt.Sprintf("signature %q spells default %s=%q but the Go emitter produces %q", call[0], arg.name, arg.value, emitted))
			}
		}
		if !setEqual(names, wanted) {
			out = append(out, fmt.Sprintf("signature %q arg set %v does not match the Go emitter's ToolArgs keys %v", call[0], sortedSet(names), sortedSet(wanted)))
		}
	}
	return out
}

func TestCodexSpawnSignatureBindsToolArgs(t *testing.T) {
	emitted := dispatch.CodexMultiAgentV2Spawn{TaskName: "spacedock_worker", Message: "assignment"}.ToolArgs()
	if len(emitted) == 0 {
		t.Fatal("ToolArgs() returned no args")
	}
	for _, msg := range codexSpawnSignatureViolations(readRepoFile(t, codexFORuntimeRel), emitted) {
		t.Errorf("%s: %s", codexFORuntimeRel, msg)
	}
}

func extractMarkdownSection(t *testing.T, text, heading string) (section string, remainder string) {
	t.Helper()
	start := strings.Index(text, heading)
	if start == -1 {
		t.Fatalf("missing heading %q", heading)
	}
	afterStart := start + len(heading)
	nextRel := strings.Index(text[afterStart:], "\n## ")
	if nextRel == -1 {
		return text[start:], text[:start]
	}
	end := afterStart + nextRel
	return text[start:end], text[:start] + text[end:]
}
