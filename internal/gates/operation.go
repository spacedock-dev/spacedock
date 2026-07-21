// ABOUTME: Recorder operations, pointer-CAS checks, digest binding, and result adoption.
// ABOUTME: These operations never model application state or invoke workflow effects.
package gates

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	jsoncanonicalizer "github.com/cyberphone/json-canonicalization/go/src/webpki.org/jsoncanonicalizer"
	"gopkg.in/yaml.v3"
)

type Operation struct {
	Operation  string   `yaml:"operation"`
	Expected   Expected `yaml:"expected"`
	GateID     string   `yaml:"gate-id"`
	Stage      string   `yaml:"stage,omitempty"`
	AttemptID  string   `yaml:"attempt-id"`
	Briefing   Briefing `yaml:"briefing,omitempty"`
	RawFilePin bool     `yaml:"raw-file-pin,omitempty"`
	Result     Result   `yaml:"result,omitempty"`
}

type Expected struct {
	Gate     string `yaml:"gate"`
	Attempt  string `yaml:"attempt"`
	Briefing string `yaml:"briefing"`
	Digest   string `yaml:"digest"`
}

type Result struct {
	BriefingDigest string  `yaml:"briefing-digest"`
	AuthorizedBy   string  `yaml:"authorized-by"`
	DelegatedFO    bool    `yaml:"delegated-fo,omitempty"`
	Entries        []Entry `yaml:"entries"`
}

type Entry struct {
	Type     string               `yaml:"type" json:"type"`
	ID       string               `yaml:"id" json:"id"`
	Briefing string               `yaml:"briefing" json:"briefing"`
	By       string               `yaml:"by,omitempty" json:"by,omitempty"`
	At       string               `yaml:"at,omitempty" json:"at,omitempty"`
	Decision string               `yaml:"decision,omitempty" json:"decision,omitempty"`
	Reason   string               `yaml:"reason,omitempty" json:"reason,omitempty"`
	Includes []string             `yaml:"includes,omitempty" json:"includes,omitempty"`
	Adoption string               `yaml:"adoption-note,omitempty" json:"adoption-note,omitempty"`
	Extra    map[string]yaml.Node `yaml:",inline" json:"-"`
}

func Record(entityPath, operationPath, briefingPath string) error {
	unlock, err := lockEntity(entityPath)
	if err != nil {
		return err
	}
	defer unlock()

	b, err := os.ReadFile(operationPath)
	if err != nil {
		return err
	}
	var op Operation
	if err := yaml.Unmarshal(b, &op); err != nil {
		return fmt.Errorf("parse operation: %w", err)
	}

	doc, oldNode, readErr := Read(entityPath)
	if readErr != nil {
		if op.Operation != "open" || !strings.Contains(readErr.Error(), "no gates record") {
			return readErr
		}
		doc = &Document{Version: 1}
	}
	if err := checkExpected(doc, op.Expected); err != nil {
		return err
	}

	switch op.Operation {
	case "open":
		if briefingPath == "" {
			return fmt.Errorf("open requires --briefing")
		}
		if findRecord(doc, op.GateID) != nil {
			return fmt.Errorf("gate %s already exists; use supersede after closure", op.GateID)
		}
		binding, err := bindBriefing(op.Briefing, briefingPath, op.RawFilePin)
		if err != nil {
			return err
		}
		doc.Records = append(doc.Records, GateRecord{ID: op.GateID, Stage: op.Stage, CurrentAttempt: op.AttemptID, Attempts: []Attempt{{ID: op.AttemptID, Sequence: 1, State: "open", Briefing: binding}}})
		doc.Current = Selection{Gate: op.GateID, Attempt: op.AttemptID}
	case "rebind":
		a, err := currentAttempt(doc, op.GateID, op.AttemptID)
		if err != nil {
			return err
		}
		if a.State != "open" {
			return fmt.Errorf("attempt %s is frozen closed", a.ID)
		}
		binding, err := bindBriefing(op.Briefing, briefingPath, op.RawFilePin)
		if err != nil {
			return err
		}
		a.Briefing = binding
	case "close":
		a, err := currentAttempt(doc, op.GateID, op.AttemptID)
		if err != nil {
			return err
		}
		if a.State != "open" {
			return fmt.Errorf("attempt %s is frozen closed", a.ID)
		}
		if !digestRE.MatchString(a.Briefing.Digest) {
			return fmt.Errorf("open legacy attempt %s has no verifiable digest; rebind it before closing", a.ID)
		}
		if op.Result.BriefingDigest != a.Briefing.Digest {
			return fmt.Errorf("result digest %s does not match current briefing digest %s", op.Result.BriefingDigest, a.Briefing.Digest)
		}
		if briefingPath != "" {
			data, err := os.ReadFile(briefingPath)
			if err != nil {
				return err
			}
			canonical, cerr := CanonicalDigest(data)
			raw := RawDigest(data)
			if a.Briefing.DigestDomain == "canonical-bytes" && (cerr != nil || canonical != a.Briefing.Digest) || a.Briefing.DigestDomain == "raw-file-pin" && raw != a.Briefing.Digest || a.Briefing.DigestDomain == "" && (cerr != nil || canonical != a.Briefing.Digest) && raw != a.Briefing.Digest {
				return fmt.Errorf("briefing bytes do not reproduce recorded digest")
			}
		}
		resolution, err := selectResolution(op.Result, op.Result.DelegatedFO)
		if err != nil {
			return err
		}
		resolution.Briefing = a.Briefing.ID // normalize only after digest validation
		a.State, a.Resolution = "closed", resolution
	case "supersede":
		r := findRecord(doc, op.GateID)
		if r == nil {
			return fmt.Errorf("unknown gate %s", op.GateID)
		}
		prev, err := currentAttempt(doc, op.GateID, r.CurrentAttempt)
		if err != nil {
			return err
		}
		if prev.State != "closed" {
			return fmt.Errorf("cannot supersede open attempt %s; rebind or close it", prev.ID)
		}
		binding, err := bindBriefing(op.Briefing, briefingPath, op.RawFilePin)
		if err != nil {
			return err
		}
		r.Attempts = append(r.Attempts, Attempt{ID: op.AttemptID, Sequence: prev.Sequence + 1, PreviousAttempt: prev.ID, State: "open", Briefing: binding})
		r.CurrentAttempt, doc.Current = op.AttemptID, Selection{Gate: r.ID, Attempt: op.AttemptID}
	default:
		return fmt.Errorf("operation must be open, rebind, close, or supersede")
	}
	if err := Validate(doc); err != nil {
		return err
	}
	if oldNode != nil {
		if err := ValidateTransition(oldNode, doc); err != nil {
			return err
		}
	}
	return write(entityPath, doc)
}

// lockEntity makes the pointer comparison and gates-only rename one
// process-atomic critical section. Contention fails closed; there is no retry
// loop or lease lifecycle for callers to coordinate.
func lockEntity(path string) (func(), error) {
	lockPath := path + ".gates.lock"
	f, err := os.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("concurrent gate writer holds %s", lockPath)
		}
		return nil, err
	}
	return func() {
		_ = f.Close()
		_ = os.Remove(lockPath)
	}, nil
}

func checkExpected(doc *Document, e Expected) error {
	actual := Expected{Gate: doc.Current.Gate, Attempt: doc.Current.Attempt}
	if a, _ := currentAttempt(doc, doc.Current.Gate, doc.Current.Attempt); a != nil {
		actual.Briefing, actual.Digest = a.Briefing.ID, a.Briefing.Digest
	}
	if actual != e {
		return fmt.Errorf("pointer conflict: expected gate=%q attempt=%q briefing=%q digest=%q; current gate=%q attempt=%q briefing=%q digest=%q", e.Gate, e.Attempt, e.Briefing, e.Digest, actual.Gate, actual.Attempt, actual.Briefing, actual.Digest)
	}
	return nil
}

func findRecord(doc *Document, id string) *GateRecord {
	for i := range doc.Records {
		if doc.Records[i].ID == id {
			return &doc.Records[i]
		}
	}
	return nil
}

func currentAttempt(doc *Document, gateID, attemptID string) (*Attempt, error) {
	r := findRecord(doc, gateID)
	if r == nil || r.CurrentAttempt != attemptID || doc.Current.Gate != gateID || doc.Current.Attempt != attemptID {
		return nil, fmt.Errorf("pointer conflict: %s/%s is not the current gate attempt", gateID, attemptID)
	}
	for i := range r.Attempts {
		if r.Attempts[i].ID == attemptID {
			return &r.Attempts[i], nil
		}
	}
	return nil, fmt.Errorf("pointer conflict: current attempt %s is missing", attemptID)
}

func bindBriefing(b Briefing, path string, raw bool) (Briefing, error) {
	if b.ID == "" || path == "" {
		return Briefing{}, fmt.Errorf("briefing id and --briefing file are required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Briefing{}, err
	}
	if raw {
		b.Digest, b.DigestDomain = RawDigest(data), "raw-file-pin"
	} else {
		b.Digest, err = CanonicalDigest(data)
		if err != nil {
			return Briefing{}, fmt.Errorf("canonicalize briefing: %w", err)
		}
		b.DigestDomain = "canonical-bytes"
	}
	return b, nil
}

func selectResolution(result Result, delegatedFO bool) (*Resolution, error) {
	seen := map[string]Entry{}
	for _, e := range result.Entries {
		if e.ID == "" || seen[e.ID].ID != "" {
			return nil, fmt.Errorf("provider log has missing or duplicate entry id")
		}
		if e.Type == "Resolution" && e.By == result.AuthorizedBy {
			// The provider entry is an envelope. Copy only the portable Resolution
			// fields so wrapper-owned or future envelope data cannot cross into the
			// durable decision record. Historical durable extras remain round-trippable
			// through Resolution.Extra, but the recorder never mints them from a result.
			r := &Resolution{Type: e.Type, ID: e.ID, Briefing: e.Briefing, By: e.By, At: e.At, Decision: e.Decision, Reason: e.Reason, Includes: e.Includes, Adoption: e.Adoption}
			for _, entry := range result.Entries {
				if entry.Briefing != r.Briefing {
					return nil, fmt.Errorf("provider log entries must belong to the same Briefing")
				}
			}
			for _, id := range r.Includes {
				included, ok := seen[id]
				if !ok || included.Briefing != r.Briefing {
					return nil, fmt.Errorf("includes must name an earlier entry from the same Briefing")
				}
			}
			if delegatedFO && r.Decision == "approve" && strings.TrimSpace(r.Reason) == "" {
				return nil, fmt.Errorf("delegated First Officer approval requires a nonblank reason")
			}
			if (r.Decision == "revise" || r.Decision == "hold") && strings.TrimSpace(r.Reason) == "" {
				validAnnotation := false
				for _, id := range r.Includes {
					included := seen[id]
					if included.Type == "Annotation" {
						validAnnotation = true
					}
				}
				if !validAnnotation {
					return nil, fmt.Errorf("%s resolution requires a nonblank reason or included earlier Annotation", r.Decision)
				}
			}
			if err := validateResolution(r, r.Briefing); err != nil {
				return nil, err
			}
			return r, nil
		}
		seen[e.ID] = e
	}
	return nil, fmt.Errorf("provider log has no Resolution by authorized actor %q", result.AuthorizedBy)
}

func ValidateTransition(oldNode *yaml.Node, next *Document) error {
	var old Document
	if err := oldNode.Decode(&old); err != nil {
		return err
	}
	for _, oldRecord := range old.Records {
		nr := findRecord(next, oldRecord.ID)
		if nr == nil {
			return fmt.Errorf("gate %s cannot be deleted", oldRecord.ID)
		}
		for _, oldAttempt := range oldRecord.Attempts {
			if oldAttempt.State != "closed" {
				continue
			}
			var found *Attempt
			for i := range nr.Attempts {
				if nr.Attempts[i].ID == oldAttempt.ID {
					found = &nr.Attempts[i]
				}
			}
			if found == nil || !nodesEqual(oldAttempt, *found) {
				return fmt.Errorf("frozen closed attempt %s cannot be deleted or mutated", oldAttempt.ID)
			}
		}
	}
	return nil
}

func nodesEqual(a, b Attempt) bool {
	ab, _ := yaml.Marshal(a)
	bb, _ := yaml.Marshal(b)
	return bytes.Equal(ab, bb)
}

func RawDigest(data []byte) string {
	s := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(s[:])
}

func CanonicalDigest(data []byte) (string, error) {
	canonical, err := jsoncanonicalizer.Transform(data)
	if err != nil {
		return "", err
	}
	return RawDigest(canonical), nil
}
