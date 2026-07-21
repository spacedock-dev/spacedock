// ABOUTME: AC-1 — binds the capability «fn» layer of the host-neutral dispatch core
// ABOUTME: across two divergeable surfaces: «fn» DEFINITIONS (## «name») and body CALLS
// ABOUTME: («name»), holding every runtime-bound block to the host-neutral arrow policy.
package contractlint

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func dispatchCorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(skillsRoot(t), "first-officer", "references", "fo-dispatch-core.md")
}

// capabilityHosts: the runtime hosts the adapters bind. A legacy core → line that
// names hosts MUST name all three; runtime-bound capabilities are host-neutral in the
// core and are bound from runtime adapter `## Runtime implementation` blocks instead.
var capabilityHosts = []string{"Claude", "Codex", "Pi"}

// fnHeadingRe: a capability «fn» DEFINITION heading `## «name»:` (the colon distinguishes
// a definition from an inline reference). fnRefRe: any capability `«name»` reference. Most
// dotted prose-functions («state.commit», «dispatch.next-action») are excluded; the
// first-class worker lifecycle capabilities are included explicitly.
var (
	fnNameRe    = `(?:worker\.(?:spawn|shutdown)|[a-z][a-z-]+)`
	fnHeadingRe = regexp.MustCompile(`(?m)^## «(` + fnNameRe + `)»:`)
	fnRefRe     = regexp.MustCompile(`«(` + fnNameRe + `)»`)
	arrowRe     = regexp.MustCompile(`(?m)^- → (.*)$`)
	metaTokens  = map[string]bool{"fn": true}
	// hostWordRe matches a host name (Claude/Codex/Pi) as a whole word, case-insensitive, so a
	// PROSE host binding on a runtime-binding arrow line ("on Claude, call Agent() directly") reds
	// even though it carries no `**Host:**` bold token. Derived from capabilityHosts so a new host
	// is covered without a second edit; word boundaries keep it from firing on "capital"/"spinning".
	hostWordRe = regexp.MustCompile(`(?i)\b(` + strings.Join(capabilityHosts, "|") + `)\b`)
)

// runtimeBoundCapabilities are the capabilities whose per-host binding is delegated to the
// runtime adapters' `## Runtime implementation` blocks via a kind-only `→ **runtime-binding**`
// arrow in the core; they carry no per-host `→` coverage in fo-dispatch-core.md itself.
var runtimeBoundCapabilities = []string{
	"worker.spawn",
	"worker.shutdown",
	"addressable-worker",
	"async-dispatch",
	"worker-identity",
	"completion-signal",
	"context-budget",
	"roster-reconcile",
}

func isRuntimeBoundCapability(name string) bool {
	for _, c := range runtimeBoundCapabilities {
		if c == name {
			return true
		}
	}
	return false
}

// fnBlock returns a «fn»'s definition body (heading to next `## `), including its → line.
func fnBlock(t *testing.T, data, name string) string {
	t.Helper()
	loc := regexp.MustCompile(`(?m)^## «` + regexp.QuoteMeta(name) + `»:`).FindStringIndex(data)
	if loc == nil {
		t.Fatalf("capability «%s» definition heading not found", name)
	}
	rest := data[loc[0]:]
	nl := strings.IndexByte(rest, '\n')
	search := rest
	if nl >= 0 {
		search = rest[nl+1:]
	}
	if end := regexp.MustCompile(`(?m)^## `).FindStringIndex(search); end != nil {
		rest = rest[:nl+1+end[0]]
	}
	return rest
}

// hostsBoundByArrow returns the capabilityHosts named (`**Host:**`) on a block's → line.
// Empty for a «fn» with no host-naming arrow (a non-capability prose «fn»).
func hostsBoundByArrow(block string) map[string]bool {
	out := map[string]bool{}
	m := arrowRe.FindStringSubmatch(block)
	if m == nil {
		return out
	}
	for _, host := range capabilityHosts {
		if strings.Contains(m[1], "**"+host+":**") {
			out[host] = true
		}
	}
	return out
}

// TestCapabilityBinding (AC-1) binds the capability «fn» layer across two surfaces of
// fo-dispatch-core.md that can diverge — the «fn» DEFINITIONS (`## «name»:` headings) and
// the body CALLS (`«name»` references beyond the heading) — and asserts they are the SAME
// set: a capability defined-but-never-called (dead definition) or called-but-never-defined
// (dangling call) reds. It then asserts every legacy core-bound capability «fn»'s → line
// binds all three hosts. It is a structural
// multi-extraction check (heading parse, reference count, per-host segment parse over one
// file), NOT prose-grep — it compares enumerations and host coverage, never asserts the doc
// contains a word. The bound tools' behavior is proven by the live lanes (AC-2/AC-6).
func TestCapabilityBinding(t *testing.T) {
	raw, err := os.ReadFile(dispatchCorePath(t))
	if err != nil {
		t.Fatalf("read dispatch core: %v", err)
	}
	data := string(raw)

	// DEFINED: `## «name»:` headings whose body has a per-host → line, plus lifecycle
	// capabilities whose concrete bindings live in runtime adapters (excludes the scheduling
	// «dispatch.next-action», whose → **prose** line names no host).
	defined := map[string]bool{}
	for _, m := range fnHeadingRe.FindAllStringSubmatch(data, -1) {
		if isRuntimeBoundCapability(m[1]) || len(hostsBoundByArrow(fnBlock(t, data, m[1]))) > 0 {
			defined[m[1]] = true
		}
	}
	if len(defined) == 0 {
		t.Fatal("no capability «fn» definitions (## «name»: + per-host → line) found — extractor bug; would pass vacuously")
	}

	// CALLED: an undotted, non-meta `«name»` referenced MORE than its single heading line.
	// Exactly-once = heading-only (dead definition); no heading = dangling call. Both
	// diverge from DEFINED and red the set-equality.
	refCount, headingCount := map[string]int{}, map[string]int{}
	for _, m := range fnRefRe.FindAllStringSubmatch(data, -1) {
		if !metaTokens[m[1]] {
			refCount[m[1]]++
		}
	}
	for _, m := range fnHeadingRe.FindAllStringSubmatch(data, -1) {
		headingCount[m[1]]++
	}
	called := map[string]bool{}
	for name, n := range refCount {
		if n-headingCount[name] >= 1 {
			called[name] = true
		}
	}
	if len(called) == 0 {
		t.Fatal("no capability «fn» body calls found — extractor bug; would pass vacuously")
	}

	if !setEqual(defined, called) {
		t.Errorf("capability «fn» set mismatch between DEFINITIONS and body CALLS in fo-dispatch-core.md:\n  defined (## «name»: + → line): %v\n  called  (body «name» refs):    %v\nevery capability «fn» must be both defined and called; neither side may rename, add, or drop one without the other",
			sortedSet(defined), sortedSet(called))
	}

	// Per-host coverage: every legacy core-bound capability «fn»'s → line binds all three hosts.
	for name := range defined {
		if isRuntimeBoundCapability(name) {
			continue
		}
		bound := hostsBoundByArrow(fnBlock(t, data, name))
		for _, host := range capabilityHosts {
			if !bound[host] {
				t.Errorf("capability «%s» → line does not bind host %q (missing `**%s:**`); every capability «fn» must carry a Claude/Codex/Pi realization on its → line", name, host, host)
			}
		}
	}
}

// runtimeBoundArrowViolations applies the host-neutral arrow policy to one runtime-bound
// capability block: every `- → ` arrow MUST be the kind-only `→ **runtime-binding**` pointer (any
// other arrow — a per-host `→ **Claude:**`, a `→ **shipped**` — reds), the runtime-binding arrow
// MUST name no host in ANY form (neither a `**Host:**` bold token NOR a prose host name like "on
// Claude, call Agent()"), and the block MUST carry no `**Host:**` token in its body either. The
// kind-only runtime-binding arrow is permitted so a cold FO has an in-file signal that the
// capability is host-bound and which adapter section binds it, without naming a host in the core.
// It returns one message per violation; the real guard and its discriminator both drive it, so a
// regression that re-bans the kind-only arrow, admits a prose host name, or stops banning a
// per-host token reds the control.
func runtimeBoundArrowViolations(block, name string) []string {
	var out []string
	for _, line := range strings.Split(block, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- → ") {
			continue
		}
		if !strings.HasPrefix(trimmed, "- → **runtime-binding**") {
			out = append(out, fmt.Sprintf("capability «%s» core block carries a non-runtime-binding arrow %q; the only arrow permitted host-neutral is a kind-only `→ **runtime-binding**` pointer", name, trimmed))
			continue
		}
		// The kind-only runtime-binding arrow must not name a host in ANY form — a `**Host:**` bold
		// token OR a bare prose host word ("on Claude, ..."). A host word anywhere on the arrow
		// line reds, keeping the core body host-neutral (the bold-token-only check missed prose).
		if m := hostWordRe.FindString(trimmed); m != "" {
			out = append(out, fmt.Sprintf("capability «%s» runtime-binding arrow names host %q in %q; the kind-only pointer must not name a host in any form", name, m, trimmed))
		}
	}
	for _, host := range capabilityHosts {
		if strings.Contains(block, "**"+host+":**") {
			out = append(out, fmt.Sprintf("capability «%s» core block contains concrete host binding for %s", name, host))
		}
	}
	return out
}

func TestDispatchCoreDefinesRuntimeBoundCapabilities(t *testing.T) {
	raw, err := os.ReadFile(dispatchCorePath(t))
	if err != nil {
		t.Fatalf("read dispatch core: %v", err)
	}
	data := string(raw)
	for _, name := range runtimeBoundCapabilities {
		block := fnBlock(t, data, name)
		for _, msg := range runtimeBoundArrowViolations(block, name) {
			t.Error(msg)
		}
	}
}

// TestDispatchCoreRuntimeBoundArrowGuardDiscriminates is the non-vacuity control for the
// runtime-bound arrow policy. It drives the same runtimeBoundArrowViolations the real guard uses
// against planted blocks: the legitimate kind-only runtime-binding arrow PASSES and a no-arrow
// block PASSES, while a per-host arrow, a `→ **shipped**` arrow, the `**Claude:**`-token smuggle,
// and — the audit-driven addition — a runtime-binding arrow that names a host in PROSE ("on
// Claude/Codex/Pi, ...") each RED. The prose-host cases prove the strengthened guard bites the
// hole the bold-token-only check left open; a loosening that over-permits (admits a prose host
// name or drops the per-host ban) or over-restricts (re-bans the kind-only arrow) fails here.
func TestDispatchCoreRuntimeBoundArrowGuardDiscriminates(t *testing.T) {
	pass := []struct {
		why, block string
	}{
		{"kind-only runtime-binding arrow", "body\n\n- → **runtime-binding**: bound in the host adapter's `## Runtime implementation`\n"},
		{"no arrow at all", "body only, no arrow\n"},
	}
	for _, c := range pass {
		if v := runtimeBoundArrowViolations(c.block, "worker.spawn"); len(v) != 0 {
			t.Fatalf("control: the %s block was wrongly flagged: %v", c.why, v)
		}
	}
	red := []struct {
		why, block string
	}{
		{"per-host arrow", "body\n- → **Claude:** Agent()\n"},
		{"`→ **shipped**` arrow", "body\n- → **shipped**: `spacedock spawn`\n"},
		{"runtime-binding arrow smuggling a `**Claude:**` token", "body\n- → **runtime-binding**: bound in **Claude:** Agent()\n"},
		{"runtime-binding arrow naming a host in prose (Claude)", "body\n- → **runtime-binding**: on Claude, call Agent() directly\n"},
		{"runtime-binding arrow naming a host in prose (Codex)", "body\n- → **runtime-binding**: on Codex, use spawn_agent\n"},
		{"runtime-binding arrow naming a host in prose (Pi)", "body\n- → **runtime-binding**: on Pi, use subagent()\n"},
	}
	for _, c := range red {
		if v := runtimeBoundArrowViolations(c.block, "worker.spawn"); len(v) == 0 {
			t.Fatalf("control: the %s was not flagged — the loosening lost a guard", c.why)
		}
	}
}
