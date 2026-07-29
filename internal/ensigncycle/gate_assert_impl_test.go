package ensigncycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

type gateHeldExpectation struct {
	gateID, attemptID, briefingID, digest string
}

func recordedGateHeldExpectation(fixture recordedGateFixture) (gateHeldExpectation, error) {
	requests, err := filepath.Glob(filepath.Join(filepath.Dir(fixture.entity), "review", "validation", "briefing-*", "request.json"))
	if err != nil {
		return gateHeldExpectation{}, err
	}
	if len(requests) != 1 {
		return gateHeldExpectation{}, fmt.Errorf("prepared fixture request count = %d, want 1", len(requests))
	}
	body, err := os.ReadFile(requests[0])
	if err != nil {
		return gateHeldExpectation{}, err
	}
	var request struct {
		Gate     string `json:"gate"`
		Attempt  string `json:"attempt"`
		Briefing struct {
			ID     string `json:"id"`
			Digest string `json:"digest"`
		} `json:"briefing"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		return gateHeldExpectation{}, fmt.Errorf("decode prepared request: %w", err)
	}
	if request.Gate == "" || request.Attempt == "" || request.Briefing.ID == "" || request.Briefing.Digest == "" {
		return gateHeldExpectation{}, fmt.Errorf("prepared request is incomplete")
	}
	return gateHeldExpectation{
		gateID: request.Gate, attemptID: request.Attempt,
		briefingID: request.Briefing.ID, digest: request.Briefing.Digest,
	}, nil
}

func assertGateHeld(before, after string, expected gateHeldExpectation) error {
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
	if selected == nil || selected.ID != expected.gateID || selected.Stage != "validation" || len(selected.Attempts) == 0 {
		return fmt.Errorf("selected validation gate has no attempt")
	}
	attempt := selected.Attempts[len(selected.Attempts)-1]
	if attempt.ID != expected.attemptID ||
		attempt.Briefing.ID != expected.briefingID || attempt.Briefing.Digest != expected.digest ||
		attempt.Resolution != nil || attempt.Application != nil {
		return fmt.Errorf("selected attempt is not open on the expected Briefing")
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
