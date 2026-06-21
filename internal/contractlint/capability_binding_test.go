// ABOUTME: AC-1 — binds the capability «fn» layer of the host-neutral dispatch core
// ABOUTME: across two divergeable surfaces: «fn» DEFINITIONS (## «name») and body CALLS
// ABOUTME: («name»), plus legacy per-host → coverage where the core still owns it.
package contractlint

import (
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

// capabilityHosts: legacy core → lines still present before the runtime-binding-block
// migration MUST name all three. New lifecycle capabilities are host-neutral here and
// are bound from runtime adapter `## Runtime implementation` blocks instead.
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
)

func isRuntimeBoundLifecycleCapability(name string) bool {
	return name == "worker.spawn" || name == "worker.shutdown"
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
		if isRuntimeBoundLifecycleCapability(m[1]) || len(hostsBoundByArrow(fnBlock(t, data, m[1]))) > 0 {
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
		if isRuntimeBoundLifecycleCapability(name) {
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

func TestDispatchCoreDefinesWorkerLifecycleCapabilities(t *testing.T) {
	raw, err := os.ReadFile(dispatchCorePath(t))
	if err != nil {
		t.Fatalf("read dispatch core: %v", err)
	}
	data := string(raw)
	for _, name := range []string{"worker.spawn", "worker.shutdown"} {
		block := fnBlock(t, data, name)
		if strings.Contains(block, "\n- → ") {
			t.Errorf("capability «%s» must stay host-neutral in fo-dispatch-core.md; bind concrete host realization in runtime adapters instead", name)
		}
		for _, host := range capabilityHosts {
			if strings.Contains(block, "**"+host+":**") {
				t.Errorf("capability «%s» core block contains concrete host binding for %s", name, host)
			}
		}
	}
}
