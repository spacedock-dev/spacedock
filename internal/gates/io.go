// ABOUTME: Fence-safe gate block reads and surgical, atomic gates-only writes.
// ABOUTME: Non-gates frontmatter and the Markdown body remain byte-identical.
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
	root, _, _, err := frontmatterNode(data)
	if err != nil {
		return nil, nil, err
	}
	gatesNode := mappingValue(root, "gates")
	if gatesNode == nil {
		return nil, nil, fmt.Errorf("entity has no gates record")
	}
	var doc Document
	if err := gatesNode.Decode(&doc); err != nil {
		return nil, nil, fmt.Errorf("decode gates: %w", err)
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
	return CurrentSummary(doc), nil
}

func write(path string, doc *Document) error {
	if err := Validate(doc); err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	root, fmStart, fmEnd, err := frontmatterNode(data)
	if err != nil {
		return err
	}
	var wrapper yaml.Node
	if err := wrapper.Encode(map[string]*Document{"gates": doc}); err != nil {
		return err
	}
	block, err := yaml.Marshal(&wrapper)
	if err != nil {
		return err
	}
	blockLines := strings.Split(strings.TrimSuffix(string(block), "\n"), "\n")
	crlf := bytes.Contains(data, []byte("\r\n"))
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	start, end := -1, -1 // absolute line indexes, start inclusive / end exclusive
	for i := 0; i+1 < len(root.Content); i += 2 {
		key := root.Content[i]
		if key.Value != "gates" {
			continue
		}
		start = fmStart + key.Line
		end = fmEnd
		if i+2 < len(root.Content) {
			end = fmStart + root.Content[i+2].Line
		}
		break
	}
	if start < 0 {
		start, end = fmEnd, fmEnd
	}
	rebuilt := make([]string, 0, len(lines)-(end-start)+len(blockLines))
	rebuilt = append(rebuilt, lines[:start]...)
	rebuilt = append(rebuilt, blockLines...)
	rebuilt = append(rebuilt, lines[end:]...)
	out := strings.Join(rebuilt, "\n")
	if crlf {
		out = strings.ReplaceAll(out, "\n", "\r\n")
	}
	return atomicWrite(path, []byte(out))
}

func frontmatterNode(data []byte) (*yaml.Node, int, int, error) {
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
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
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}
