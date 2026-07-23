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
	return CurrentSummary(doc), nil
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

func readRoundPointerData(data []byte) (RoundPointer, bool, error) {
	root, _, _, err := frontmatterNode(data)
	if err != nil {
		return RoundPointer{}, false, err
	}
	node := mappingValue(root, "review-round")
	if node == nil {
		return RoundPointer{}, false, nil
	}
	encoded, err := yaml.Marshal(node)
	if err != nil {
		return RoundPointer{}, false, err
	}
	var pointer RoundPointer
	decoder := yaml.NewDecoder(bytes.NewReader(encoded))
	decoder.KnownFields(true)
	if err := decoder.Decode(&pointer); err != nil {
		return RoundPointer{}, false, fmt.Errorf("decode review-round pointer: %w", err)
	}
	wantRoom := fmt.Sprintf("./review/%s/round-%d", pointer.Stage, pointer.Cycle)
	if pointer.ID == "" || !roundStageRE.MatchString(pointer.Stage) || pointer.Cycle < 1 ||
		pointer.Briefing.ID == "" || !digestRE.MatchString(pointer.Briefing.Digest) ||
		pointer.Briefing.DigestDomain != "canonical-bytes" || pointer.Briefing.RoomRef != wantRoom {
		return RoundPointer{}, false, fmt.Errorf("entity has invalid review-round pointer")
	}
	return pointer, true, nil
}

func rebuildRoundEntity(original []byte, pointer RoundPointer, projection string, complete bool) ([]byte, error) {
	root, fmStart, fmEnd, err := frontmatterNode(original)
	if err != nil {
		return nil, err
	}
	current := mappingValue(root, "review-round")
	block, err := yaml.Marshal(struct {
		Round RoundPointer `yaml:"review-round"`
	}{Round: pointer})
	if err != nil {
		return nil, err
	}
	start, end := fmEnd, fmEnd
	if current != nil {
		for i := 0; i+1 < len(root.Content); i += 2 {
			if root.Content[i].Value != "review-round" {
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
	out := make([]byte, 0, len(original)-(endByte-startByte)+len(replacement)+len(projection)+32)
	out = append(out, original[:startByte]...)
	out = append(out, replacement...)
	out = append(out, original[endByte:]...)
	out, err = projectFeedbackCycle(out, projection, complete)
	if err != nil {
		return nil, err
	}
	got, ok, err := readRoundPointerData(out)
	if err != nil || !ok || !sameRoundPointer(got, pointer) {
		if err == nil {
			err = fmt.Errorf("rebuilt entity review-round pointer does not validate")
		}
		return nil, err
	}
	if _, _, gateErr := readData(out); gateErr != nil && !strings.Contains(gateErr.Error(), "no gates record") {
		return nil, fmt.Errorf("validate rebuilt gates: %w", gateErr)
	}
	return out, nil
}

type markdownHeading struct {
	level int
	text  string
	start int
}

func projectFeedbackCycle(data []byte, line string, complete bool) ([]byte, error) {
	cyclePrefix := line[:strings.Index(line, ":")+1]
	headings := markdownBodyHeadings(data)
	var feedback []markdownHeading
	for _, heading := range headings {
		if heading.level == 3 && heading.text == "Feedback Cycles" {
			feedback = append(feedback, heading)
		}
	}
	if len(feedback) > 1 {
		return nil, fmt.Errorf("entity has more than one Feedback Cycles heading")
	}
	bodyStart := bodyByteOffset(data)
	body := data[bodyStart:]
	exactCount := countStandaloneLine(body, line)
	cycleCount := countLinePrefix(body, cyclePrefix)
	if !complete {
		if cycleCount != 0 {
			return nil, fmt.Errorf("incomplete round already has a Feedback Cycles projection")
		}
		return data, nil
	}
	if exactCount == 1 && cycleCount == 1 {
		return data, nil
	}
	if exactCount != 0 || cycleCount != 0 {
		return nil, fmt.Errorf("Feedback Cycles projection changed for %s", cyclePrefix)
	}
	if len(feedback) == 0 {
		separator := "\n\n"
		if len(data) == 0 || bytes.HasSuffix(data, []byte("\n\n")) {
			separator = ""
		} else if bytes.HasSuffix(data, []byte("\n")) {
			separator = "\n"
		}
		return append(data, []byte(separator+"### Feedback Cycles\n\n"+line+"\n")...), nil
	}
	insert := len(data)
	for _, heading := range headings {
		if heading.start > feedback[0].start && heading.level <= 3 {
			insert = heading.start
			break
		}
	}
	prefix := "\n"
	if insert >= 2 && bytes.Equal(data[insert-2:insert], []byte("\n\n")) {
		prefix = ""
	}
	addition := []byte(prefix + line + "\n\n")
	out := make([]byte, 0, len(data)+len(addition))
	out = append(out, data[:insert]...)
	out = append(out, addition...)
	out = append(out, data[insert:]...)
	return out, nil
}

func markdownBodyHeadings(data []byte) []markdownHeading {
	start := bodyByteOffset(data)
	var headings []markdownHeading
	inFence := false
	for offset := start; offset < len(data); {
		end := bytes.IndexByte(data[offset:], '\n')
		if end < 0 {
			end = len(data) - offset
		}
		line := strings.TrimSuffix(string(data[offset:offset+end]), "\r")
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
		} else if !inFence {
			level := 0
			for level < len(line) && line[level] == '#' {
				level++
			}
			if level >= 1 && level <= 6 && level < len(line) && line[level] == ' ' {
				headings = append(headings, markdownHeading{level: level, text: strings.TrimSpace(line[level+1:]), start: offset})
			}
		}
		offset += end
		if offset < len(data) {
			offset++
		}
	}
	return headings
}

func bodyByteOffset(data []byte) int {
	_, _, fmEnd, err := frontmatterNode(data)
	if err != nil {
		return 0
	}
	return lineOffset(data, fmEnd+1)
}

func countStandaloneLine(data []byte, want string) int {
	count := 0
	for _, line := range normalizedLines(data) {
		if line == want {
			count++
		}
	}
	return count
}

func countLinePrefix(data []byte, prefix string) int {
	count := 0
	for _, line := range normalizedLines(data) {
		if strings.HasPrefix(line, prefix) {
			count++
		}
	}
	return count
}

func commitRound(entityPath string, expected []byte, room string, roomExists bool, briefing, log, rebuilt []byte, entityWriter func(string, []byte) error) error {
	current, err := os.ReadFile(entityPath)
	if err != nil {
		return err
	}
	if !bytes.Equal(current, expected) {
		return fmt.Errorf("entity changed during locked round update")
	}
	if roomExists {
		logPath := filepath.Join(room, "briefing.review.jsonl")
		oldLog, err := os.ReadFile(logPath)
		if err != nil {
			return err
		}
		if !bytes.Equal(oldLog, log) {
			if err := atomicWrite(logPath, log); err != nil {
				return err
			}
		}
		if err := entityWriter(entityPath, rebuilt); err != nil {
			if restoreErr := atomicWrite(logPath, oldLog); restoreErr != nil {
				return fmt.Errorf("write entity: %v; restore round log: %w", err, restoreErr)
			}
			return err
		}
		return nil
	}

	parent := filepath.Dir(room)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	cleanupParents := func() {
		_ = os.Remove(parent)
		_ = os.Remove(filepath.Dir(parent))
	}
	tmp, err := os.MkdirTemp(parent, ".round-*")
	if err != nil {
		cleanupParents()
		return err
	}
	defer os.RemoveAll(tmp)
	if err := writeSyncedFile(filepath.Join(tmp, "briefing.json"), briefing); err != nil {
		cleanupParents()
		return err
	}
	if err := writeSyncedFile(filepath.Join(tmp, "briefing.review.jsonl"), log); err != nil {
		cleanupParents()
		return err
	}
	if err := os.Rename(tmp, room); err != nil {
		cleanupParents()
		return err
	}
	if err := entityWriter(entityPath, rebuilt); err != nil {
		_ = os.RemoveAll(room)
		cleanupParents()
		return err
	}
	return nil
}

func writeSyncedFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}
