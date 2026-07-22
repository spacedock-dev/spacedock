// ABOUTME: Fence-safe gate block reads and surgical, atomic gates-only writes.
// ABOUTME: Non-gates frontmatter and the Markdown body remain byte-identical.
package gates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

type lineEdit struct {
	start    int
	end      int
	lines    []string
	oldValue string
	newValue string
}

type gateLocation struct {
	gates       *yaml.Node
	current     *yaml.Node
	records     *yaml.Node
	record      *yaml.Node
	attempts    *yaml.Node
	recordEnd   int
	attemptsEnd int
	recordsEnd  int
}

// writeNewGate adds the minimal v1 projection. It intentionally omits the
// legacy attempt pointer, sequence, lineage, and state fields.
func writeNewGate(path string, doc *Document, record GateRecord) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	root, _, fmEnd, err := frontmatterNode(data)
	if err != nil {
		return err
	}
	if mappingValue(root, "gates") == nil {
		minimal := &Document{Version: 1, Current: Selection{Gate: record.ID}, Records: []GateRecord{record}}
		block, err := yaml.Marshal(struct {
			Gates *Document `yaml:"gates"`
		}{Gates: minimal})
		if err != nil {
			return err
		}
		return applyLineEdits(path, data, []lineEdit{{start: fmEnd, end: fmEnd, lines: splitYAML(block)}})
	}

	loc, err := locateGate(data, "")
	if err != nil {
		return err
	}
	currentGate := mappingValue(loc.current, "gate")
	if currentGate == nil {
		return fmt.Errorf("gates.current has no gate selection")
	}
	item, err := sequenceItemLines(record, sequenceIndent(data, loc.records))
	if err != nil {
		return err
	}
	edits := []lineEdit{
		scalarEdit(data, currentGate, record.ID),
		{start: loc.recordsEnd, end: loc.recordsEnd, lines: item},
	}
	if node := mappingValue(loc.current, "attempt"); node != nil {
		edits = append(edits, scalarEdit(data, node, record.Attempts[0].ID))
	}
	if err := Validate(doc); err != nil {
		return err
	}
	return applyLineEdits(path, data, edits)
}

func writeBriefingRebind(path, gateID, attemptID string, binding Briefing) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, attempt, attemptEnd, err := locateAttempt(data, gateID, attemptID)
	if err != nil {
		return err
	}
	start, end, _, err := mappingPairRange(attempt, "briefing", locLineStart(data), attemptEnd)
	if err != nil {
		return err
	}
	fragment, err := yaml.Marshal(struct {
		Briefing Briefing `yaml:"briefing"`
	}{Briefing: binding})
	if err != nil {
		return err
	}
	return applyLineEdits(path, data, []lineEdit{{start: start, end: end, lines: indentLines(splitYAML(fragment), lineIndent(data, start))}})
}

// writeBriefingSuccessor makes only existing selection edits and one attempt
// append. Frozen closures and opaque legacy fields remain byte-identical.
func writeBriefingSuccessor(path, gateID string, next Attempt) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	loc, err := locateGate(data, gateID)
	if err != nil {
		return err
	}
	if len(loc.attempts.Content) == 0 {
		return fmt.Errorf("gate %s has no attempt history", gateID)
	}
	item, err := sequenceItemLines(next, sequenceIndent(data, loc.attempts))
	if err != nil {
		return err
	}
	edits := []lineEdit{
		scalarEdit(data, mappingValue(loc.current, "gate"), gateID),
		{start: loc.attemptsEnd, end: loc.attemptsEnd, lines: item},
	}
	if node := mappingValue(loc.current, "attempt"); node != nil {
		edits = append(edits, scalarEdit(data, node, next.ID))
	}
	if node := mappingValue(loc.record, "current-attempt"); node != nil {
		edits = append(edits, scalarEdit(data, node, next.ID))
	}
	return applyLineEdits(path, data, edits)
}

func writeResolution(path, gateID, attemptID string, resolution *Resolution) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	loc, attempt, attemptEnd, err := locateAttempt(data, gateID, attemptID)
	if err != nil {
		return err
	}
	if mappingValue(attempt, "resolution") != nil {
		return fmt.Errorf("attempt %s is frozen closed", attemptID)
	}
	briefingStart, insertAt, _, err := mappingPairRange(attempt, "briefing", locLineStart(data), attemptEnd)
	if err != nil {
		return err
	}
	fragment, err := yaml.Marshal(struct {
		Resolution *Resolution `yaml:"resolution"`
	}{Resolution: resolution})
	if err != nil {
		return err
	}
	edits := []lineEdit{{start: insertAt, end: insertAt, lines: indentLines(splitYAML(fragment), lineIndent(data, briefingStart))}}
	if node := mappingValue(attempt, "state"); node != nil {
		edits = append(edits, scalarEdit(data, node, "closed"))
	}
	if node := mappingValue(loc.current, "gate"); node != nil {
		edits = append(edits, scalarEdit(data, node, gateID))
	}
	if node := mappingValue(loc.current, "attempt"); node != nil {
		edits = append(edits, scalarEdit(data, node, attemptID))
	}
	return applyLineEdits(path, data, edits)
}

func locateGate(data []byte, gateID string) (gateLocation, error) {
	root, fmStart, fmEnd, err := frontmatterNode(data)
	if err != nil {
		return gateLocation{}, err
	}
	_, gatesEnd, gatesNode, err := mappingPairRange(root, "gates", fmStart, fmEnd)
	if err != nil {
		return gateLocation{}, err
	}
	current := mappingValue(gatesNode, "current")
	recordsStart, recordsEnd, records := 0, 0, (*yaml.Node)(nil)
	recordsStart, recordsEnd, records, err = mappingPairRange(gatesNode, "records", fmStart, gatesEnd)
	_ = recordsStart
	if err != nil || current == nil || records.Kind != yaml.SequenceNode {
		return gateLocation{}, fmt.Errorf("gates record has no editable current/records projection")
	}
	loc := gateLocation{gates: gatesNode, current: current, records: records, recordsEnd: recordsEnd}
	if gateID == "" {
		return loc, nil
	}
	for i, record := range records.Content {
		id := mappingValue(record, "id")
		if id == nil || id.Value != gateID {
			continue
		}
		loc.record = record
		loc.recordEnd = recordsEnd
		if i+1 < len(records.Content) {
			loc.recordEnd = fmStart + records.Content[i+1].Line
		}
		_, loc.attemptsEnd, loc.attempts, err = mappingPairRange(record, "attempts", fmStart, loc.recordEnd)
		if err != nil || loc.attempts.Kind != yaml.SequenceNode {
			return gateLocation{}, fmt.Errorf("gate %s has no attempts", gateID)
		}
		return loc, nil
	}
	return gateLocation{}, fmt.Errorf("unknown gate %s", gateID)
}

func locateAttempt(data []byte, gateID, attemptID string) (gateLocation, *yaml.Node, int, error) {
	loc, err := locateGate(data, gateID)
	if err != nil {
		return gateLocation{}, nil, 0, err
	}
	fmStart := locLineStart(data)
	for i, attempt := range loc.attempts.Content {
		id := mappingValue(attempt, "id")
		if id == nil || id.Value != attemptID {
			continue
		}
		end := loc.attemptsEnd
		if i+1 < len(loc.attempts.Content) {
			end = fmStart + loc.attempts.Content[i+1].Line
		}
		return loc, attempt, end, nil
	}
	return gateLocation{}, nil, 0, fmt.Errorf("gate %s attempt %s is missing", gateID, attemptID)
}

func mappingPairRange(parent *yaml.Node, key string, fmStart, containerEnd int) (int, int, *yaml.Node, error) {
	for i := 0; i+1 < len(parent.Content); i += 2 {
		if parent.Content[i].Value != key {
			continue
		}
		start, end := fmStart+parent.Content[i].Line, containerEnd
		if i+2 < len(parent.Content) {
			end = fmStart + parent.Content[i+2].Line
		}
		return start, end, parent.Content[i+1], nil
	}
	return 0, 0, nil, fmt.Errorf("mapping has no %s field", key)
}

func locLineStart(data []byte) int {
	_, start, _, _ := frontmatterNode(data)
	return start
}

func scalarEdit(data []byte, node *yaml.Node, value string) lineEdit {
	if node == nil {
		return lineEdit{start: -1}
	}
	line := locLineStart(data) + node.Line
	lines := normalizedLines(data)
	if line < 0 || line >= len(lines) || !strings.Contains(lines[line], node.Value) {
		return lineEdit{start: -1}
	}
	return lineEdit{start: line, end: line + 1, oldValue: node.Value, newValue: value}
}

func applyLineEdits(path string, original []byte, edits []lineEdit) error {
	lines := normalizedLines(original)
	for _, edit := range edits {
		if edit.start < 0 || edit.end < edit.start || edit.end > len(lines) {
			return fmt.Errorf("cannot locate surgical gate edit")
		}
	}
	sort.SliceStable(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	for _, edit := range edits {
		if edit.oldValue != "" {
			if !strings.Contains(lines[edit.start], edit.oldValue) {
				return fmt.Errorf("cannot locate scalar %q for surgical gate edit", edit.oldValue)
			}
			lines[edit.start] = strings.Replace(lines[edit.start], edit.oldValue, edit.newValue, 1)
			continue
		}
		rebuilt := make([]string, 0, len(lines)-(edit.end-edit.start)+len(edit.lines))
		rebuilt = append(rebuilt, lines[:edit.start]...)
		rebuilt = append(rebuilt, edit.lines...)
		rebuilt = append(rebuilt, lines[edit.end:]...)
		lines = rebuilt
	}
	out := strings.Join(lines, "\n")
	if bytes.Contains(original, []byte("\r\n")) {
		out = strings.ReplaceAll(out, "\n", "\r\n")
	}
	return atomicWrite(path, []byte(out))
}

func normalizedLines(data []byte) []string {
	return strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
}

func lineIndent(data []byte, line int) string {
	lines := normalizedLines(data)
	if line < 0 || line >= len(lines) {
		return ""
	}
	return lines[line][:len(lines[line])-len(strings.TrimLeft(lines[line], " \t"))]
}

func sequenceIndent(data []byte, sequence *yaml.Node) string {
	if len(sequence.Content) == 0 {
		return "    "
	}
	return lineIndent(data, locLineStart(data)+sequence.Content[0].Line)
}

func sequenceItemLines(value any, indent string) ([]string, error) {
	body, err := yaml.Marshal(value)
	if err != nil {
		return nil, err
	}
	lines := splitYAML(body)
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty YAML sequence item")
	}
	out := []string{indent + "- " + lines[0]}
	for _, line := range lines[1:] {
		out = append(out, indent+"  "+line)
	}
	return out, nil
}

func splitYAML(body []byte) []string {
	return strings.Split(strings.TrimSuffix(string(body), "\n"), "\n")
}

func indentLines(lines []string, indent string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = indent + line
	}
	return out
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
