// ABOUTME: Containment lint — the five shared/host-neutral contract files name no
// ABOUTME: runtime host (word or model name) and no host-specific tool token; adapters own those.
package contractlint

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// hostNeutralSharedFiles are the shared contract files every host loads. A host name
// or host tool token here reads to the other two hosts as a universal claim their
// runtime cannot satisfy — the class both field incidents (standing-teammate
// idempotency, bare-mode streaming fan-out) belong to. Host-specific coverage lives
// in the runtime adapters' binding sections instead.
var hostNeutralSharedFiles = []string{
	"skills/first-officer/references/first-officer-shared-core.md",
	"skills/first-officer/references/fo-dispatch-core.md",
	"skills/first-officer/references/fo-write-core.md",
	"skills/first-officer/references/fo-merge-core.md",
	"skills/ensign/references/ensign-shared-core.md",
}

// hostNeutralWordRe matches a host name or a host model name as a whole word,
// case-insensitive — the same shape as hostWordRe (capability_binding_test.go),
// widened with the Claude model enum, and derived from capabilityHosts so a new
// host is covered without a second edit. Word boundaries keep `pi` from firing on
// lookalikes (`api`, `spike`, `pickup`) — the discriminator control proves it.
var hostNeutralWordRe = regexp.MustCompile(`(?i)\b(` +
	strings.Join(append(append([]string{}, capabilityHosts...), "Opus", "Sonnet", "Haiku", "Fable"), "|") + `)\b`)

// claudeHostToolTokens are the Claude-host tool tokens banned from the shared files;
// the Codex and Pi vocabularies ride codexToolTokens / piContainmentTokens
// (runtime_binding_block_test.go), so one token list per host exists in one place.
//
// Excluded from the token set, with grounds: generic harness verbs every host
// exposes (`Bash`, `Read`, `Grep`, `Write`, `Edit`, `Skill`) and spacedock's own
// cross-host dispatch-prompt convention (`Skill(skill="spacedock:ensign")`, emitted
// by internal/dispatch/build.go for every host).
var claudeHostToolTokens = []string{
	"Agent(",
	"SendMessage",
	"TeamCreate",
	"TeamDelete",
	"task_notification",
	"run_in_background",
	"subagent_type",
	"ToolSearch",
	"BashOutput",
}

// hostNeutralToolTokens is the union of all three hosts' tool-token vocabularies,
// deduplicated (codexToolTokens also carries the Claude SendMessage ban). The bare
// Pi token `acceptance` is excluded here with grounds: in the shared files it is
// core workflow vocabulary ("acceptance criteria"), not the `subagent(...
// acceptance: ...)` transport field; the field token stays contained to the Pi
// adapter's binding sections by TestRuntimeToolTokensStayInBindingSections.
func hostNeutralToolTokens(t *testing.T) []string {
	t.Helper()
	seen := map[string]bool{"acceptance": true}
	var out []string
	for _, tok := range append(append(append([]string{}, claudeHostToolTokens...), codexToolTokens...), piContainmentTokens(t)...) {
		if !seen[tok] {
			seen[tok] = true
			out = append(out, tok)
		}
	}
	return out
}

// sharedCoreHostViolations is the single scanner the real check and its discriminator
// control both drive: per line, a host/model word (word-bounded, case-insensitive)
// or a host tool token is a violation. Existence-fact containment — it asserts where
// a name may appear, never what a host or tool DOES.
func sharedCoreHostViolations(body string, tokens []string) []string {
	var out []string
	for i, line := range strings.Split(body, "\n") {
		if m := hostNeutralWordRe.FindString(line); m != "" {
			out = append(out, fmt.Sprintf("line %d names host %q outside a runtime adapter: %q", i+1, m, strings.TrimSpace(line)))
		}
		for range toolTokenContainmentViolations(line, tokens) {
			out = append(out, fmt.Sprintf("line %d carries a host tool token outside a runtime adapter: %q", i+1, strings.TrimSpace(line)))
		}
	}
	return out
}

// TestSharedContractFilesStayHostNeutral (AC-1) is the containment invariant: the
// five shared files carry no host name and no host tool token. Green requires the
// legacy per-host `→` coverage to have RELOCATED into the adapters — a suppression
// that deleted meaning instead would strand each host's binding, which the adapter
// set-equality and preservation checks red on independently.
func TestSharedContractFilesStayHostNeutral(t *testing.T) {
	tokens := hostNeutralToolTokens(t)
	if len(tokens) == 0 {
		t.Fatal("no host tool tokens enumerated — the containment check would pass vacuously")
	}
	for _, rel := range hostNeutralSharedFiles {
		body := readRepoFile(t, filepath.FromSlash(rel))
		if len(body) == 0 {
			t.Fatalf("%s read empty — the containment check would pass vacuously", rel)
		}
		for _, msg := range sharedCoreHostViolations(body, tokens) {
			t.Errorf("%s %s", rel, msg)
		}
	}
}

// TestSharedContractHostNeutralGuardDiscriminates is the non-vacuity control: it
// drives the same sharedCoreHostViolations the real check uses. Capability-speak and
// `pi`-lookalike words PASS; a planted host word (each host, plus a model name and a
// prose-case variant) and a planted tool token from each host's vocabulary RED.
func TestSharedContractHostNeutralGuardDiscriminates(t *testing.T) {
	tokens := hostNeutralToolTokens(t)
	pass := []struct{ why, body string }{
		{"capability-speak prose", "Await the worker result per `«async-dispatch»`; completion is recognized via `«completion-signal»`.\n"},
		{"pi-lookalike words", "the api surface, a spike, and a pickup of the worktree\n"},
		{"acceptance-criteria domain language", "entity-level acceptance criteria this stage naturally advances\n"},
		{"kind-only runtime-binding arrow", "- → **runtime-binding**: bound in the host adapter's `## Runtime implementation`\n"},
	}
	for _, c := range pass {
		if v := sharedCoreHostViolations(c.body, tokens); len(v) != 0 {
			t.Fatalf("control: the %s was wrongly flagged: %v", c.why, v)
		}
	}
	red := []struct{ why, body string }{
		{"planted Claude host word", "on Claude, call the spawn tool directly\n"},
		{"planted Codex host word", "Codex blocks until the worker returns\n"},
		{"planted Pi host word", "poll the run id on Pi\n"},
		{"planted lower-case host word", "the claude enum owns the model space\n"},
		{"planted model name", "prefer opus for gate stages\n"},
		{"planted Claude tool token", "spawn via Agent(name=…) and wait\n"},
		{"planted Claude polling token", "wait via BashOutput polling\n"},
		{"planted Codex tool token", "call spawn_agent(task_name,message) for every ready entity\n"},
		{"planted Pi substrate token", "initial creation maps to member_spawn\n"},
	}
	for _, c := range red {
		if v := sharedCoreHostViolations(c.body, tokens); len(v) == 0 {
			t.Fatalf("control: the %s was not flagged — the containment guard stopped biting", c.why)
		}
	}
}
