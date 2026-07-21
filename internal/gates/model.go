// ABOUTME: Durable gate-resolution record model and invariant validation.
// ABOUTME: The application subtree is deliberately opaque and preserved verbatim.
package gates

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var digestRE = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

type Document struct {
	Version int                  `yaml:"version" json:"version"`
	Current Selection            `yaml:"current" json:"current"`
	Records []GateRecord         `yaml:"records" json:"records"`
	Extra   map[string]yaml.Node `yaml:",inline" json:"-"`
}

type Selection struct {
	Gate    string               `yaml:"gate" json:"gate"`
	Attempt string               `yaml:"attempt" json:"attempt"`
	Extra   map[string]yaml.Node `yaml:",inline" json:"-"`
}

type GateRecord struct {
	ID             string               `yaml:"id" json:"id"`
	Stage          string               `yaml:"stage" json:"stage"`
	CurrentAttempt string               `yaml:"current-attempt" json:"current-attempt"`
	Attempts       []Attempt            `yaml:"attempts" json:"attempts"`
	Note           string               `yaml:"note,omitempty" json:"note,omitempty"`
	Extra          map[string]yaml.Node `yaml:",inline" json:"-"`
}

type Attempt struct {
	ID              string               `yaml:"id" json:"id"`
	Sequence        int                  `yaml:"sequence" json:"sequence"`
	PreviousAttempt string               `yaml:"previous-attempt,omitempty" json:"previous-attempt,omitempty"`
	State           string               `yaml:"state" json:"state"`
	Briefing        Briefing             `yaml:"briefing" json:"briefing"`
	Resolution      *Resolution          `yaml:"resolution,omitempty" json:"resolution,omitempty"`
	Application     any                  `yaml:"application,omitempty" json:"-"`
	Note            string               `yaml:"note,omitempty" json:"note,omitempty"`
	Extra           map[string]yaml.Node `yaml:",inline" json:"-"`
}

type Briefing struct {
	ID           string               `yaml:"id" json:"id"`
	Digest       string               `yaml:"digest" json:"digest"`
	DigestDomain string               `yaml:"digest-domain,omitempty" json:"digest-domain,omitempty"`
	RoomRef      string               `yaml:"room-ref,omitempty" json:"room-ref,omitempty"`
	Note         string               `yaml:"note,omitempty" json:"note,omitempty"`
	Extra        map[string]yaml.Node `yaml:",inline" json:"-"`
}

type Resolution struct {
	Type     string               `yaml:"type" json:"type"`
	ID       string               `yaml:"id" json:"id"`
	Briefing string               `yaml:"briefing" json:"briefing"`
	By       string               `yaml:"by" json:"by"`
	At       string               `yaml:"at" json:"at"`
	Decision string               `yaml:"decision" json:"decision"`
	Reason   string               `yaml:"reason,omitempty" json:"reason,omitempty"`
	Includes []string             `yaml:"includes,omitempty" json:"includes,omitempty"`
	Adoption string               `yaml:"adoption-note,omitempty" json:"adoption-note,omitempty"`
	Extra    map[string]yaml.Node `yaml:",inline" json:"-"`
}

type Summary struct {
	Gate       string
	Attempt    string
	State      string
	Briefing   string
	Resolution string
	Decision   string
}

func Validate(doc *Document) error {
	if doc.Version != 1 {
		return fmt.Errorf("gates.version must be 1")
	}
	if doc.Current.Gate == "" || doc.Current.Attempt == "" {
		return fmt.Errorf("gates.current must name a gate and attempt")
	}
	gateIDs, attemptIDs, briefingIDs, resolutionIDs := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}
	var selected *Attempt
	for ri := range doc.Records {
		r := &doc.Records[ri]
		if r.ID == "" || r.Stage == "" || r.CurrentAttempt == "" || gateIDs[r.ID] {
			return fmt.Errorf("record %d has missing or duplicate identity/current pointer", ri+1)
		}
		gateIDs[r.ID] = true
		currentFound := false
		for ai := range r.Attempts {
			a := &r.Attempts[ai]
			if a.ID == "" || attemptIDs[a.ID] || a.Sequence < 1 || ai > 0 && a.Sequence != r.Attempts[ai-1].Sequence+1 {
				return fmt.Errorf("gate %s attempt %d has missing/duplicate id or non-contiguous sequence", r.ID, ai+1)
			}
			attemptIDs[a.ID] = true
			if ai > 0 && a.PreviousAttempt != r.Attempts[ai-1].ID {
				return fmt.Errorf("attempt %s previous-attempt does not name sequence %d", a.ID, ai)
			}
			if a.Briefing.ID == "" || briefingIDs[a.Briefing.ID] || a.Briefing.Digest != "" && !digestRE.MatchString(a.Briefing.Digest) {
				return fmt.Errorf("attempt %s has invalid or duplicate briefing binding", a.ID)
			}
			briefingIDs[a.Briefing.ID] = true
			if a.Briefing.DigestDomain != "" && a.Briefing.DigestDomain != "canonical-bytes" && a.Briefing.DigestDomain != "raw-file-pin" {
				return fmt.Errorf("attempt %s has unknown digest-domain %q", a.ID, a.Briefing.DigestDomain)
			}
			switch a.State {
			case "open":
				if a.Resolution != nil {
					return fmt.Errorf("open attempt %s cannot carry a resolution", a.ID)
				}
			case "closed":
				if a.Resolution == nil {
					return fmt.Errorf("closed attempt %s requires a resolution", a.ID)
				}
				if err := validateResolution(a.Resolution, a.Briefing.ID); err != nil {
					return fmt.Errorf("attempt %s: %w", a.ID, err)
				}
				if resolutionIDs[a.Resolution.ID] {
					return fmt.Errorf("duplicate resolution id %s", a.Resolution.ID)
				}
				resolutionIDs[a.Resolution.ID] = true
			default:
				return fmt.Errorf("attempt %s state must be open or closed", a.ID)
			}
			if a.ID == r.CurrentAttempt {
				currentFound = true
			}
			if r.ID == doc.Current.Gate && a.ID == doc.Current.Attempt {
				selected = a
			}
		}
		if !currentFound {
			return fmt.Errorf("gate %s current-attempt pointer does not resolve", r.ID)
		}
		if r.ID == doc.Current.Gate && r.CurrentAttempt != doc.Current.Attempt {
			return fmt.Errorf("current pointer conflict: gate %s selects %s but gates.current selects %s", r.ID, r.CurrentAttempt, doc.Current.Attempt)
		}
	}
	if selected == nil {
		return fmt.Errorf("gates.current pointer does not resolve to one attempt")
	}
	return nil
}

func validateResolution(r *Resolution, briefingID string) error {
	if r.Type != "Resolution" || r.ID == "" || r.Briefing != briefingID || r.By == "" || r.At == "" {
		return fmt.Errorf("resolution identity, attribution, or briefing binding is invalid")
	}
	switch r.Decision {
	case "approve":
	case "revise", "hold":
		if strings.TrimSpace(r.Reason) == "" && len(r.Includes) == 0 {
			return fmt.Errorf("%s resolution requires a reason or included earlier Annotation", r.Decision)
		}
	default:
		return fmt.Errorf("resolution decision must be approve, revise, or hold")
	}
	return nil
}

func CurrentSummary(doc *Document) Summary {
	for i := range doc.Records {
		r := &doc.Records[i]
		if r.ID != doc.Current.Gate {
			continue
		}
		for j := range r.Attempts {
			a := &r.Attempts[j]
			if a.ID != doc.Current.Attempt {
				continue
			}
			s := Summary{Gate: r.ID, Attempt: a.ID, State: a.State, Briefing: a.Briefing.ID}
			if a.Resolution != nil {
				s.Resolution, s.Decision = a.Resolution.ID, a.Resolution.Decision
			}
			return s
		}
	}
	return Summary{}
}
