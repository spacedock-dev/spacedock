package ensigncycle

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"gopkg.in/yaml.v3"
)

var (
	validatingStatus = regexp.MustCompile(`(?im)^status:\s*validation\s*$`)
	reviewStatus     = regexp.MustCompile(`(?im)^status:\s*review\s*$`)
	completedSet     = regexp.MustCompile(`(?im)^completed:[^\S\n]*\S.*$`)
	verdictSetFM     = regexp.MustCompile(`(?im)^verdict:[^\S\n]*\S.*$`)
)

func assertGateHeld(before, after, review string) error {
	if before == after || !validatingStatus.MatchString(after) || completedSet.MatchString(after) || verdictSetFM.MatchString(after) {
		return fmt.Errorf("gated entity is not held at its open validation boundary")
	}
	doc, err := decodeGateDocument(after)
	if err != nil {
		return fmt.Errorf("canonical gate document: %w", err)
	}
	var selected *gates.GateRecord
	for i := range doc.Records {
		if doc.Records[i].ID == doc.Current.Gate {
			selected = &doc.Records[i]
			break
		}
	}
	if selected == nil || selected.ID != "gate:docs-dev:3k:validation" || selected.Stage != "validation" || len(selected.Attempts) == 0 {
		return fmt.Errorf("selected validation gate has no attempt")
	}
	attempt := selected.Attempts[len(selected.Attempts)-1]
	if attempt.ID != "gate-attempt:3k-validation-1" ||
		attempt.Briefing.ID != recordedGateBriefingID || attempt.Briefing.Digest != recordedGateDigest ||
		attempt.Resolution != nil || attempt.Application != nil {
		return fmt.Errorf("selected attempt is not open on the expected Briefing")
	}
	if err := assertConciseRecordedGateReview(review); err != nil {
		return fmt.Errorf("semantic gate review: %w", err)
	}
	return nil
}

func decodeGateDocument(entity string) (*gates.Document, error) {
	parts := bytes.SplitN([]byte(entity), []byte("---\n"), 3)
	if len(parts) != 3 || len(parts[0]) != 0 {
		return nil, fmt.Errorf("entity has no frontmatter")
	}
	var root yaml.Node
	if err := yaml.Unmarshal(parts[1], &root); err != nil || len(root.Content) == 0 {
		return nil, fmt.Errorf("decode frontmatter: %w", err)
	}
	mapping := root.Content[0]
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value != "gates" {
			continue
		}
		encoded, err := yaml.Marshal(mapping.Content[i+1])
		if err != nil {
			return nil, err
		}
		var doc gates.Document
		decoder := yaml.NewDecoder(bytes.NewReader(encoded))
		decoder.KnownFields(true)
		if err := decoder.Decode(&doc); err != nil {
			return nil, err
		}
		if err := gates.Validate(&doc); err != nil {
			return nil, err
		}
		return &doc, nil
	}
	return nil, fmt.Errorf("entity has no gates record")
}

type recordedGateCodexEvent struct {
	Type string `json:"type"`
	Item struct {
		Type             string `json:"type"`
		Text             string `json:"text"`
		Command          string `json:"command"`
		Status           string `json:"status"`
		AggregatedOutput string `json:"aggregated_output"`
		ExitCode         *int   `json:"exit_code"`
	} `json:"item"`
}
