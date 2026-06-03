// ABOUTME: Self-referential-AC classifier + require-external-proof README opt-in,
// ABOUTME: powering runSet's terminal guard and validateWorkflow's sub-check.
package status

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ACFlag is one classifier finding: the AC header line, the cleaned proof
// clause that was scanned, and the self-phrase that triggered the flag. Tests
// drive the classifier and assert against these fields directly.
type ACFlag struct {
	Header        string
	ProofClause   string
	MatchedPhrase string
}

// classifierCallCount is incremented every time ClassifyEntityACs runs. The
// AC-4 shared-classifier test exercises both integrating callers in one body
// and asserts both bumped this counter, proving the two surfaces call the
// same function.
var classifierCallCount int

// acHeaderRe matches a line opening an `**AC-N — …**` block. The classifier
// only emits findings for blocks that match this shape; free-form prose is
// never scanned.
var acHeaderRe = regexp.MustCompile(`^\*\*AC-`)

// proofMarkerRe locates the first proof-clause marker within an AC block:
// `Verified by`, `Oracle:`, `Proof:`, or `End state:` / `End state.`.
// Case-insensitive. The match is the boundary between the body and the proof
// clause taken to block end.
var proofMarkerRe = regexp.MustCompile(`(?i)(verified by|oracle:|proof:|end state[:.])`)

// quotedSpanRe matches every double-quoted span (`"..."`) and backtick-fenced
// span (`` `...` ``) so a quoted example of the antipattern is excluded from
// the self-phrase match.
var quotedSpanRe = regexp.MustCompile("\"[^\"]*\"|`[^`]*`")

// selfPhraseRes are the case-insensitive self-reference phrases the
// classifier flags. The order is stable so a test can inspect MatchedPhrase
// against the literal alternative.
var selfPhraseRes = []*regexp.Regexp{
	regexp.MustCompile(`(?i)this entity'?s`),
	regexp.MustCompile(`(?i)review of (this|the) (entity|decision|design)[^.]*section`),
	regexp.MustCompile(`(?i)the entity'?s own (prose|decision|section)`),
	regexp.MustCompile(`(?i)re-reading (this|the) (entity|task|body)`),
}

// externalTokenRe captures any token that proves the AC's proof cites
// something runnable / observable outside the entity body. Presence of ANY
// token clears the AC, so a CI run on the entity's PR no longer false-positives
// even when the prose says "this entity's own PR".
var externalTokenRe = regexp.MustCompile(`(?i)` + strings.Join([]string{
	`\btest\b`, `\.go\b`, `exit\s`, `exit-code`, `exit code`, `command`, `\bstatus\b`,
	`--\w+`, `fixture`, `golden`, `byte`, `on-disk`, `stdout`, `stderr`, `assert`, `parser`,
	`mutator`, `frontmatter`, `code path`, `command/parser`,
	`runs? the`, `running the`, `invok`, `drive the`, `driving the`,
	`\bCI\b`, `\bPR\b`, `live job`, `green`, `workflow file`,
}, "|"))

// stripFrontmatter returns the body of an entity file — the bytes after the
// closing `---` fence. A file with no opening fence returns its bytes
// unchanged; a file with only an opening fence returns the empty string.
func stripFrontmatter(data []byte) string {
	if !contentHasOpeningFence(data) {
		return string(data)
	}
	lines := splitLines(string(data))
	var body []string
	fences := 0
	first := true
	for _, raw := range lines {
		line := raw
		if first {
			line = strings.TrimPrefix(line, utf8BOM)
			first = false
		}
		if line == "---" {
			fences++
			continue
		}
		if fences >= 2 {
			body = append(body, line)
		}
	}
	return strings.Join(body, "\n")
}

// ClassifyEntityACs scans an entity body and flags every AC whose proof clause
// names the entity itself (and only the entity itself) as the verification. The
// five-step algorithm: extract `**AC-N` blocks → isolate the clause from the
// first proof marker onward → strip quoted spans → match a self-phrase →
// require external-token absence. The classifier is pure — no I/O — so tests
// drive it with literal strings.
func ClassifyEntityACs(body string) []ACFlag {
	classifierCallCount++

	var flags []ACFlag

	lines := splitLines(body)
	var current []string
	var currentHeader string

	flush := func() {
		if currentHeader == "" {
			return
		}
		block := strings.Join(current, "\n")
		clause := isolateProofClause(block)
		cleaned := quotedSpanRe.ReplaceAllString(clause, " ")
		matched := matchSelfPhrase(cleaned)
		if matched == "" {
			return
		}
		if externalTokenRe.MatchString(cleaned) {
			return
		}
		flags = append(flags, ACFlag{
			Header:        currentHeader,
			ProofClause:   strings.TrimSpace(clause),
			MatchedPhrase: matched,
		})
	}

	for _, line := range lines {
		if acHeaderRe.MatchString(line) {
			flush()
			currentHeader = line
			current = nil
			continue
		}
		if strings.HasPrefix(line, "## ") {
			flush()
			currentHeader = ""
			current = nil
			continue
		}
		if currentHeader != "" {
			current = append(current, line)
		}
	}
	flush()

	return flags
}

// isolateProofClause returns the slice of an AC block from the first proof
// marker to the block end. When no marker is present the entire block is
// returned (the absence-of-proof case is a separate failure mode the
// downstream FO cross-check catches).
func isolateProofClause(block string) string {
	loc := proofMarkerRe.FindStringIndex(block)
	if loc == nil {
		return block
	}
	return block[loc[0]:]
}

// matchSelfPhrase returns the literal substring of clause that matched a
// self-phrase regex, or "" when none matched. The returned phrase populates
// ACFlag.MatchedPhrase for test introspection.
func matchSelfPhrase(clause string) string {
	for _, re := range selfPhraseRes {
		if m := re.FindString(clause); m != "" {
			return m
		}
	}
	return ""
}

// externalProofPolicy is a workflow's declared external-proof opt-in, read
// from the README's top-level `require-external-proof:` key. Default OFF
// (byte-identical to absent) leaves non-dev workflows untouched.
type externalProofPolicy int

const (
	externalProofOff externalProofPolicy = iota
	externalProofOn
)

// resolveExternalProofPolicy reads the README's top-level
// `require-external-proof:` key and returns the declared policy. An absent or
// empty key, or an explicit `false`, defaults to externalProofOff —
// byte-identical to a workflow that never declared the key. An unknown value
// is rejected loudly rather than silently coerced, so a typo
// (`require-external-proof: tru`) fails fast instead of silently allowing the
// close.
func resolveExternalProofPolicy(definitionDir string) (externalProofPolicy, error) {
	value := strings.TrimSpace(ParseFrontmatter(filepath.Join(definitionDir, "README.md"))["require-external-proof"])
	switch value {
	case "", "false":
		return externalProofOff, nil
	case "true":
		return externalProofOn, nil
	default:
		return externalProofOff, fmt.Errorf("README require-external-proof: must be 'true' or 'false' (or absent for the default 'false'), not '%s'", value)
	}
}

// classifyEntityFile reads the entity file at path and returns the classifier
// findings against its body. A missing or unreadable file yields no flags
// (the FO surfaces a missing-entity defect elsewhere; the classifier is not
// the right layer to escalate I/O errors).
func classifyEntityFile(path string) []ACFlag {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return ClassifyEntityACs(stripFrontmatter(data))
}

// flaggedACLabels returns a `[AC-1,AC-3]` shaped string for an error message,
// pulled from the leading `**AC-N — ...` of each header.
func flaggedACLabels(flags []ACFlag) string {
	labels := make([]string, 0, len(flags))
	for _, f := range flags {
		labels = append(labels, acLabel(f.Header))
	}
	return strings.Join(labels, ",")
}

// acLabelRe extracts `AC-N` from a `**AC-N — …**` header line.
var acLabelRe = regexp.MustCompile(`AC-[0-9]+`)

func acLabel(header string) string {
	return acLabelRe.FindString(header)
}
