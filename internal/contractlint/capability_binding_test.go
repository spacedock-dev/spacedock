// ABOUTME: AC-1 — binds the named-capability set in two independent sources
// ABOUTME: (fo-dispatch-core.md `## Named Capabilities` and each host adapter's
// ABOUTME: `## Capability implementations` subsection) as equal sets.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// dispatchCorePath is the host-neutral dispatch core that declares the named
// capabilities in its `## Named Capabilities` section.
func dispatchCorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(skillsRoot(t), "first-officer", "references", "fo-dispatch-core.md")
}

// adapterPaths returns the three FO runtime adapters that bind each capability to
// concrete tools in their `## Capability implementations` subsection.
func adapterPaths(t *testing.T) []string {
	t.Helper()
	base := filepath.Join(skillsRoot(t), "first-officer", "references")
	return []string{
		filepath.Join(base, "claude-first-officer-runtime.md"),
		filepath.Join(base, "codex-first-officer-runtime.md"),
		filepath.Join(base, "pi-first-officer-runtime.md"),
	}
}

// sectionHeadingRe matches a top-level Markdown heading line (`## …`). Used to
// bound a section's slice to the text between its heading and the next heading.
var sectionHeadingRe = regexp.MustCompile(`(?m)^## `)

// sectionBlock returns the slice of a Markdown file bounded by the given heading
// title (`## {title}`) and the next top-level heading. Fails if the heading is
// missing. The slice includes the heading line.
func sectionBlock(t *testing.T, path, title string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	heading := "## " + title
	loc := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(heading) + `\b`).FindIndex(data)
	if loc == nil {
		t.Fatalf("heading %q not found in %s", heading, path)
	}
	rest := data[loc[0]:]
	// Drop the matched heading line, then find the next `## ` heading after it.
	nextNL := strings.IndexByte(string(rest), '\n')
	search := rest
	if nextNL >= 0 {
		search = rest[nextNL+1:]
	}
	if end := sectionHeadingRe.FindIndex(search); end != nil {
		rest = append([]byte(nil), rest[:nextNL+1+end[0]]...)
	}
	return string(rest)
}

// capabilityBulletRe matches a Markdown bullet that binds a named capability: a
// line starting with `- “name“ ` where name is a lowercase hyphenated token.
// The backtick-delimited name is the single capability identifier both the core
// and each adapter use, so the extractor binds that identifier — not surrounding
// prose — to the capability set.
var capabilityBulletRe = regexp.MustCompile(`(?m)^- ` + "`" + `([a-z][a-z-]+)` + "`" + `(?:[ —:])`)

// extractCapabilities scans a section block for capability-binding bullets and
// returns the set of capability names. Fails (via the caller's empty-set guard)
// if it yields none — a broken extractor that returns [] on both sides would
// otherwise make the equality pass vacuously.
func extractCapabilities(t *testing.T, block, source string) []string {
	t.Helper()
	var out []string
	for _, m := range capabilityBulletRe.FindAllStringSubmatch(block, -1) {
		out = append(out, m[1])
	}
	if len(out) == 0 {
		t.Fatalf("%s capability extraction yielded zero capabilities — extractor bug; the binding would pass vacuously", source)
	}
	return out
}

// TestCapabilityBinding (AC-1) asserts the named-capability set the host-neutral
// dispatch core declares (`fo-dispatch-core.md` `## Named Capabilities`) and the
// bound-capability set EACH host adapter declares (`## Capability implementations`)
// are the SAME set. They are independent values that can diverge: a capability
// renamed, added, or dropped in the core OR in any single adapter reds the
// binding. It is a structural dual-extraction check (two delimited-token parses
// over independent files), NOT prose-grep: it never asserts a doc contains a
// given word; it compares extracted enumerations. The behavior that the adapters
// bind the right concrete tools is proven by the live lanes (AC-2/AC-6), not
// here — this test binds only the capability-name contract surface.
func TestCapabilityBinding(t *testing.T) {
	coreBlock := sectionBlock(t, dispatchCorePath(t), "Named Capabilities")
	core := extractCapabilities(t, coreBlock, "core `## Named Capabilities`")
	coreSet := toSet(core)

	for _, path := range adapterPaths(t) {
		rel := relPath(t, path)
		block := sectionBlock(t, path, "Capability implementations")
		adapter := extractCapabilities(t, block, rel+" `## Capability implementations`")
		adapterSet := toSet(adapter)
		if !setEqual(coreSet, adapterSet) {
			t.Errorf("capability set mismatch between the core and adapter %s:\n  core (fo-dispatch-core.md `## Named Capabilities`): %v\n  adapter (%s `## Capability implementations`): %v\nevery core-declared capability must be bound by each adapter; neither side may rename, add, or drop a capability without the other",
				rel, sortedSet(coreSet), rel, sortedSet(adapterSet))
		}
	}
}

// relPath returns a repo-relative path for readable error output.
func relPath(t *testing.T, path string) string {
	t.Helper()
	rel, err := filepath.Rel(repoRoot(t), path)
	if err != nil {
		return path
	}
	return rel
}
