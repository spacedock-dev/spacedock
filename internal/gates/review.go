package gates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type Annotation struct {
	Type     string   `json:"type"`
	ID       string   `json:"id"`
	Briefing string   `json:"briefing"`
	By       string   `json:"by,omitempty"`
	At       string   `json:"at,omitempty"`
	Target   string   `json:"target,omitempty"`
	Kind     string   `json:"kind,omitempty"`
	Body     string   `json:"body,omitempty"`
	Includes []string `json:"includes,omitempty"`
}
type reviewEntry struct {
	Annotation
	Decision   string      `json:"decision,omitempty"`
	Reason     string      `json:"reason,omitempty"`
	Resolution *Resolution `json:"-"`
}
type reviewLog struct {
	Entries  []reviewEntry
	Reviewer *Resolution
}

func validateAnnotation(a Annotation, briefingID string, prior map[string]reviewEntry) error {
	if _, err := time.Parse(time.RFC3339Nano, a.At); a.Type != "Annotation" || a.ID == "" ||
		a.Briefing != briefingID || a.By == "" || err != nil {
		return fmt.Errorf("Annotation identity, attribution, or same Briefing binding is invalid")
	}
	for _, ref := range a.Includes {
		if _, ok := prior[ref]; !ok {
			return fmt.Errorf("Annotation %s includes non-earlier identity %s", a.ID, ref)
		}
	}
	return nil
}
func parseReviewLog(data []byte, briefingID string) (reviewLog, error) {
	var result reviewLog
	prior := map[string]reviewEntry{}
	if len(data) == 0 || data[len(data)-1] != '\n' || !utf8.Valid(data) {
		return result, fmt.Errorf("review log must be non-empty UTF-8 JSONL ending in a complete line")
	}
	for i, line := range bytes.Split(data[:len(data)-1], []byte{'\n'}) {
		var entry reviewEntry
		if len(bytes.TrimSpace(line)) == 0 || json.Unmarshal(line, &entry) != nil {
			return result, fmt.Errorf("review log entry %d is not JSON", i+1)
		}
		if _, duplicate := prior[entry.ID]; duplicate {
			return result, fmt.Errorf("review log entry %d has duplicate identity", i+1)
		}
		switch entry.Type {
		case "Annotation":
			if err := validateAnnotation(entry.Annotation, briefingID, prior); err != nil {
				return result, err
			}
		case "Resolution":
			entry.Resolution = &Resolution{Type: entry.Type, ID: entry.ID, Briefing: entry.Briefing, By: entry.By, At: entry.At,
				Decision: entry.Decision, Reason: entry.Reason, Includes: entry.Includes}
			if err := validateResolution(entry.Resolution, briefingID); err != nil {
				return result, fmt.Errorf("Resolution %s: %w", entry.ID, err)
			}
			if _, err := time.Parse(time.RFC3339Nano, entry.At); err != nil {
				return result, fmt.Errorf("Resolution %s has invalid attribution time", entry.ID)
			}
			hasAnnotation := false
			for _, ref := range entry.Includes {
				previous, ok := prior[ref]
				if !ok {
					return result, fmt.Errorf("Resolution %s includes non-earlier identity %s", entry.ID, ref)
				}
				hasAnnotation = hasAnnotation || previous.Resolution == nil
			}
			if entry.Decision != "approve" && strings.TrimSpace(entry.Reason) == "" && !hasAnnotation {
				return result, fmt.Errorf("Resolution %s lacks included Annotation rationale", entry.ID)
			}
			if result.Reviewer == nil {
				result.Reviewer = entry.Resolution
			}
		default:
			return result, fmt.Errorf("review log entry %d has unknown type", i+1)
		}
		result.Entries, prior[entry.ID] = append(result.Entries, entry), entry
	}
	if result.Reviewer == nil {
		return result, fmt.Errorf("review log has no Resolution")
	}
	return result, nil
}
