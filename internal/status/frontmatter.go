// ABOUTME: Frontmatter reader — fence finder + EOF/CRLF normalization are
// ABOUTME: hand-rolled; the in-fence slice is parsed by gopkg.in/yaml.v3.
package status

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// utf8BOM is the leading byte-order mark stripped from a file's first line
// before the opening-fence check.
const utf8BOM = "\uFEFF"

// hasOpeningFence reports whether the file's first non-empty, non-BOM line is
// exactly `---`. Leading truly-empty lines (`\n` only) are skipped; a
// whitespace-only first content line disqualifies the file. A leading UTF-8 BOM
// on the first line is stripped before the check.
func hasOpeningFence(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return contentHasOpeningFence(data)
}

// contentHasOpeningFence is hasOpeningFence over in-memory bytes (for --new,
// which validates STDIN before any file exists).
func contentHasOpeningFence(data []byte) bool {
	first := true
	for _, raw := range splitLines(string(data)) {
		line := raw
		if first {
			line = strings.TrimPrefix(line, utf8BOM)
			first = false
		}
		if line == "" {
			continue
		}
		return line == "---"
	}
	return false
}

// normalizeNewlines translates CRLF and lone CR into LF (Python text-mode
// universal-newlines equivalent), so a `---\r` fence compares equal to `---`
// and a CRLF README's stages block parses. `\r\n` is collapsed first so a CRLF
// pair never yields two LFs.
func normalizeNewlines(s string) string {
	if !strings.ContainsRune(s, '\r') {
		return s
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}

// splitLines splits on '\n' the way Python's `for line in f` iteration does
// after rstrip('\n'): the trailing newline is removed from each line. A file
// ending in '\n' yields no extra empty trailing element. Newlines are
// normalized first.
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(normalizeNewlines(s), "\n")
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

// ParseFrontmatter extracts top-level key/value pairs between the first and
// second `---` fences. Returns an empty map when there is no opening fence.
// Implementation: the hand-rolled fence finder slices the in-fence bytes;
// gopkg.in/yaml.v3 parses them as a YAML mapping; only top-level scalar
// values are surfaced (nested mappings / sequences render as empty strings,
// preserving the prior "indented lines ignored" semantic). Last-key-wins is
// native to yaml.v3.
func ParseFrontmatter(path string) map[string]string {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]string{}
	}
	return parseFrontmatterContent(data)
}

// parseFrontmatterContent is ParseFrontmatter over in-memory bytes.
func parseFrontmatterContent(data []byte) map[string]string {
	fields := map[string]string{}
	if !contentHasOpeningFence(data) {
		return fields
	}
	slice := frontmatterSlice(data)
	if len(slice) == 0 {
		return fields
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(slice, &doc); err != nil {
		// A malformed YAML body is surfaced as no fields. The migration check
		// asserts every live entity parses; a parse failure here is the loud
		// signal we want for a true corruption (e.g., the retired
		// mismatched-quote quirk).
		return fields
	}
	if len(doc.Content) == 0 {
		return fields
	}
	mapping := doc.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return fields
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		k := mapping.Content[i]
		v := mapping.Content[i+1]
		if k.Kind != yaml.ScalarNode {
			continue
		}
		key := k.Value
		if v.Kind == yaml.ScalarNode {
			fields[key] = v.Value
		} else {
			// Nested mapping or sequence — match the prior "indented lines
			// ignored" semantic by surfacing the key with an empty value.
			fields[key] = ""
		}
	}
	return fields
}

// isSpaceByte reports whether b is one of the ASCII whitespace bytes Python's
// str.isspace treats as space for the leading-char check used by stages.go,
// mutate.go, and new.go to skip indented YAML lines in their stages-block /
// frontmatter-line passes.
func isSpaceByte(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	}
	return false
}

// frontmatterSlice returns the raw bytes between the first two `---` fences
// of data (the YAML body the reader feeds to yaml.v3). A missing closing
// fence yields the bytes from after the opening fence to EOF. BOM on the
// first line is stripped; CRLF/CR are universal-newline normalized. Returns
// nil when there is no opening fence.
func frontmatterSlice(data []byte) []byte {
	if !contentHasOpeningFence(data) {
		return nil
	}
	lines := splitLines(string(data))
	var body []string
	inFM := false
	first := true
	for _, raw := range lines {
		line := raw
		if first {
			line = strings.TrimPrefix(line, utf8BOM)
			first = false
		}
		if line == "---" {
			if inFM {
				break
			}
			inFM = true
			continue
		}
		if !inFM {
			continue
		}
		body = append(body, line)
	}
	if len(body) == 0 {
		return []byte{}
	}
	return []byte(strings.Join(body, "\n") + "\n")
}
