// ABOUTME: Fence-safe canonical v1 gate reads and atomic gates-subtree writes.
// ABOUTME: Writes preserve every byte outside the binary-owned gates mapping.
package gates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func Read(path string) (*Document, *yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	return readData(data)
}

func readData(data []byte) (*Document, *yaml.Node, error) {
	root, _, _, err := frontmatterNode(data)
	if err != nil {
		return nil, nil, err
	}
	gatesNode := mappingValue(root, "gates")
	if gatesNode == nil {
		return nil, nil, fmt.Errorf("entity has no gates record")
	}
	encoded, err := yaml.Marshal(gatesNode)
	if err != nil {
		return nil, nil, fmt.Errorf("encode gates for validation: %w", err)
	}
	var doc Document
	decoder := yaml.NewDecoder(bytes.NewReader(encoded))
	decoder.KnownFields(true)
	if err := decoder.Decode(&doc); err != nil {
		return nil, nil, fmt.Errorf("decode canonical gates v1: %w", err)
	}
	if err := Validate(&doc); err != nil {
		return nil, nil, err
	}
	return &doc, gatesNode, nil
}

func SummaryFile(path string) (Summary, error) {
	doc, _, err := Read(path)
	if err != nil {
		return Summary{}, err
	}
	summary := CurrentSummary(doc)
	eligibility, err := EligibilityFile(path)
	if err != nil {
		return Summary{}, err
	}
	summary.Condition = eligibility.Condition
	summary.Eligible = eligibility.Eligible
	return summary, nil
}

func entityStatus(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	root, _, _, err := frontmatterNode(data)
	if err != nil {
		return "", err
	}
	status := mappingValue(root, "status")
	if status == nil || strings.TrimSpace(status.Value) == "" {
		return "", fmt.Errorf("entity has no current workflow status")
	}
	return status.Value, nil
}

func writeDocument(path string, expected *yaml.Node, doc *Document) error {
	if err := Validate(doc); err != nil {
		return err
	}
	original, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	root, fmStart, fmEnd, err := frontmatterNode(original)
	if err != nil {
		return err
	}
	current := mappingValue(root, "gates")
	if !sameYAMLNode(expected, current) {
		return fmt.Errorf("gates record changed during locked update")
	}
	block, err := yaml.Marshal(struct {
		Gates *Document `yaml:"gates"`
	}{Gates: doc})
	if err != nil {
		return err
	}
	start, end := fmEnd, fmEnd
	if current != nil {
		for i := 0; i+1 < len(root.Content); i += 2 {
			if root.Content[i].Value != "gates" {
				continue
			}
			start = fmStart + root.Content[i].Line
			if i+2 < len(root.Content) {
				end = fmStart + root.Content[i+2].Line
			}
			break
		}
	}
	replacement := block
	if bytes.Contains(original, []byte("\r\n")) {
		replacement = []byte(strings.ReplaceAll(string(replacement), "\n", "\r\n"))
	}
	startByte, endByte := lineOffset(original, start), lineOffset(original, end)
	out := make([]byte, 0, len(original)-(endByte-startByte)+len(replacement))
	out = append(out, original[:startByte]...)
	out = append(out, replacement...)
	out = append(out, original[endByte:]...)
	if _, _, err := readData(out); err != nil {
		return fmt.Errorf("validate rebuilt gates: %w", err)
	}
	return atomicWrite(path, out)
}

func sameYAMLNode(left, right *yaml.Node) bool {
	if left == nil || right == nil {
		return left == right
	}
	lb, lerr := yaml.Marshal(left)
	rb, rerr := yaml.Marshal(right)
	return lerr == nil && rerr == nil && bytes.Equal(lb, rb)
}

func normalizedLines(data []byte) []string {
	return strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
}

func lineOffset(data []byte, line int) int {
	if line <= 0 {
		return 0
	}
	seen := 0
	for i, b := range data {
		if b != '\n' {
			continue
		}
		seen++
		if seen == line {
			return i + 1
		}
	}
	return len(data)
}

func frontmatterNode(data []byte) (*yaml.Node, int, int, error) {
	lines := normalizedLines(data)
	fmStart, fmEnd := -1, -1
	for i, line := range lines {
		if strings.TrimSpace(line) != "---" {
			continue
		}
		if fmStart < 0 {
			fmStart = i
		} else {
			fmEnd = i
			break
		}
	}
	if fmStart < 0 || fmEnd < 0 {
		return nil, 0, 0, fmt.Errorf("entity has no complete frontmatter")
	}
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(strings.Join(lines[fmStart+1:fmEnd], "\n")+"\n"), &doc); err != nil {
		return nil, 0, 0, fmt.Errorf("parse frontmatter: %w", err)
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, 0, 0, fmt.Errorf("frontmatter must be a mapping")
	}
	return doc.Content[0], fmStart, fmEnd, nil
}

func mappingValue(m *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(m.Content); i += 2 {
		if m.Content[i].Value == key {
			return m.Content[i+1]
		}
	}
	return nil
}

func atomicWrite(path string, data []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".gates-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err := tmp.Chmod(info.Mode()); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}
