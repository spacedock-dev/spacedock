package ensigncycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

// canonicalPreparedBriefingID is the v1 Briefing identity grammar. Group 1 is
// the entity identity, group 2 the stage, group 3 the attempt ordinal.
var canonicalPreparedBriefingID = regexp.MustCompile(`^briefing:(.+):([^:]+):attempt-([1-9][0-9]*):revision-[1-9][0-9]*$`)

// recordedGateHeldExpectation reads what the gate must hold from the published
// room, not from the entity the assertion checks. The room is one file, so the
// canonical Briefing carries the identity and the digest.
func recordedGateHeldExpectation(fixture recordedGateFixture) (gateHeldExpectation, error) {
	briefings, err := filepath.Glob(filepath.Join(filepath.Dir(fixture.entity), "review", "validation", "briefing-*", "index.json"))
	if err != nil {
		return gateHeldExpectation{}, err
	}
	if len(briefings) != 1 {
		return gateHeldExpectation{}, fmt.Errorf("prepared fixture Briefing count = %d, want 1", len(briefings))
	}
	body, err := os.ReadFile(briefings[0])
	if err != nil {
		return gateHeldExpectation{}, err
	}
	var briefing struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &briefing); err != nil {
		return gateHeldExpectation{}, fmt.Errorf("decode prepared Briefing: %w", err)
	}
	parts := canonicalPreparedBriefingID.FindStringSubmatch(briefing.ID)
	if parts == nil {
		return gateHeldExpectation{}, fmt.Errorf("prepared Briefing id %q is not canonical", briefing.ID)
	}
	digest, err := gates.CanonicalDigest(body)
	if err != nil {
		return gateHeldExpectation{}, fmt.Errorf("canonicalize prepared Briefing: %w", err)
	}
	entity, stage, ordinal := parts[1], parts[2], parts[3]
	if cut := strings.LastIndex(entity, ":"); cut >= 0 {
		entity = entity[cut+1:]
	}
	return gateHeldExpectation{
		gateID:     "gate:" + parts[1] + ":" + stage,
		attemptID:  "gate-attempt:" + entity + "-" + stage + "-" + ordinal,
		briefingID: briefing.ID,
		digest:     digest,
	}, nil
}

func semanticGateHeldExpectation(fixture recordedGateFixture) (gateHeldExpectation, error) {
	expected, err := recordedGateHeldExpectation(fixture)
	if err != nil {
		return gateHeldExpectation{}, &gradedErr{code: "gate-not-held", msg: fmt.Sprintf("read prepared gate expectation: %v", err)}
	}
	return expected, nil
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
		if doc.Records[i].Stage == "validation" {
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
