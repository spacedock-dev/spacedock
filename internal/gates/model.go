// ABOUTME: Canonical v1 durable gate-resolution and one-use application model.
// ABOUTME: Validation keeps unknown or conflicting application state fail-closed.
package gates

import (
	"fmt"
	"regexp"
	"strings"
)

var digestRE = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

type Document struct {
	Version int          `yaml:"version" json:"version"`
	Current Selection    `yaml:"current" json:"current"`
	Records []GateRecord `yaml:"records" json:"records"`
}

type Selection struct {
	Gate string `yaml:"gate" json:"gate"`
}

type GateRecord struct {
	ID       string    `yaml:"id" json:"id"`
	Stage    string    `yaml:"stage" json:"stage"`
	Attempts []Attempt `yaml:"attempts" json:"attempts"`
}

type Attempt struct {
	ID          string       `yaml:"id" json:"id"`
	Briefing    Briefing     `yaml:"briefing" json:"briefing"`
	Resolution  *Resolution  `yaml:"resolution,omitempty" json:"resolution,omitempty"`
	Application *Application `yaml:"application,omitempty" json:"application,omitempty"`
}

type Application struct {
	Action        string         `yaml:"action" json:"action"`
	TargetStage   string         `yaml:"target-stage,omitempty" json:"target-stage,omitempty"`
	State         string         `yaml:"state" json:"state"`
	Blockers      *[]Blocker     `yaml:"blockers,omitempty" json:"blockers,omitempty"`
	ExecutionHold *ExecutionHold `yaml:"execution-hold,omitempty" json:"execution-hold,omitempty"`
	Feedback      *Feedback      `yaml:"feedback,omitempty" json:"feedback,omitempty"`
}

type Blocker struct {
	ID               string `yaml:"id,omitempty" json:"id,omitempty"`
	Kind             string `yaml:"kind,omitempty" json:"kind,omitempty"`
	Ref              string `yaml:"ref,omitempty" json:"ref,omitempty"`
	ExpectedRevision string `yaml:"expected-revision,omitempty" json:"expected-revision,omitempty"`
	ExpectedState    string `yaml:"expected-state,omitempty" json:"expected-state,omitempty"`
	State            string `yaml:"state,omitempty" json:"state,omitempty"`
}

type ExecutionHold struct {
	ID     string `yaml:"id,omitempty" json:"id,omitempty"`
	State  string `yaml:"state" json:"state"`
	By     string `yaml:"by,omitempty" json:"by,omitempty"`
	At     string `yaml:"at,omitempty" json:"at,omitempty"`
	Reason string `yaml:"reason,omitempty" json:"reason,omitempty"`
}

type Feedback struct {
	Cycle         int    `yaml:"cycle,omitempty" json:"cycle,omitempty"`
	FindingRef    string `yaml:"finding-ref,omitempty" json:"finding-ref,omitempty"`
	FindingDigest string `yaml:"finding-digest,omitempty" json:"finding-digest,omitempty"`
}

type Briefing struct {
	ID           string `yaml:"id" json:"id"`
	Digest       string `yaml:"digest" json:"digest"`
	DigestDomain string `yaml:"digest-domain" json:"digest-domain"`
	RoomRef      string `yaml:"room-ref" json:"room-ref"`
}

type Resolution struct {
	Type     string   `yaml:"type" json:"type"`
	ID       string   `yaml:"id" json:"id"`
	Briefing string   `yaml:"briefing" json:"briefing"`
	By       string   `yaml:"by" json:"by"`
	At       string   `yaml:"at" json:"at"`
	Decision string   `yaml:"decision" json:"decision"`
	Reason   string   `yaml:"reason,omitempty" json:"reason,omitempty"`
	Includes []string `yaml:"includes,omitempty" json:"includes,omitempty"`
	Adoption string   `yaml:"adoption-note,omitempty" json:"adoption-note,omitempty"`
}

type Summary struct {
	Gate             string
	Attempt          string
	State            string
	Briefing         string
	Resolution       string
	Decision         string
	Application      string
	ApplicationState string
	Condition        string
	Eligible         bool
	TargetStage      string
}

type Eligibility struct {
	Gate             string
	Attempt          string
	Action           string
	TargetStage      string
	ApplicationState string
	Condition        string
	Eligible         bool
}

type ConsumeResult struct {
	Eligibility
	Consumed bool
}

func Validate(doc *Document) error {
	if doc.Version != 1 {
		return fmt.Errorf("gates.version must be 1")
	}
	if doc.Current.Gate == "" {
		return fmt.Errorf("gates.current must name a gate")
	}
	gateIDs, attemptIDs, briefingIDs, resolutionIDs := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}
	selected := false
	for ri := range doc.Records {
		r := &doc.Records[ri]
		pendingApplications := 0
		if r.ID == "" || r.Stage == "" || gateIDs[r.ID] || len(r.Attempts) == 0 {
			return fmt.Errorf("record %d has missing or duplicate identity or no attempts", ri+1)
		}
		gateIDs[r.ID] = true
		if r.ID == doc.Current.Gate {
			selected = true
		}
		for ai := range r.Attempts {
			a := &r.Attempts[ai]
			if a.ID == "" || attemptIDs[a.ID] {
				return fmt.Errorf("gate %s attempt %d has missing or duplicate id", r.ID, ai+1)
			}
			attemptIDs[a.ID] = true
			if a.Briefing.ID == "" || briefingIDs[a.Briefing.ID] || !digestRE.MatchString(a.Briefing.Digest) || a.Briefing.RoomRef == "" {
				return fmt.Errorf("attempt %s has invalid or duplicate briefing binding", a.ID)
			}
			briefingIDs[a.Briefing.ID] = true
			if a.Briefing.DigestDomain != "canonical-bytes" && a.Briefing.DigestDomain != "raw-file-pin" {
				return fmt.Errorf("attempt %s has unknown digest-domain %q", a.ID, a.Briefing.DigestDomain)
			}
			if a.Resolution == nil {
				if a.Application != nil {
					return fmt.Errorf("open attempt %s cannot carry application data", a.ID)
				}
				continue
			}
			if err := validateResolution(a.Resolution, a.Briefing.ID); err != nil {
				return fmt.Errorf("attempt %s: %w", a.ID, err)
			}
			if a.Application != nil {
				if err := validateApplication(a.Application, a.Resolution.Decision); err != nil {
					return fmt.Errorf("attempt %s: %w", a.ID, err)
				}
				if a.Application.State == "pending" {
					pendingApplications++
				}
			}
			if resolutionIDs[a.Resolution.ID] {
				return fmt.Errorf("duplicate resolution id %s", a.Resolution.ID)
			}
			resolutionIDs[a.Resolution.ID] = true
		}
		if pendingApplications > 1 {
			return fmt.Errorf("gate %s carries more than one pending application", r.ID)
		}
	}
	if !selected {
		return fmt.Errorf("gates.current pointer does not resolve to one logical gate")
	}
	return nil
}

func validateApplication(a *Application, decision string) error {
	switch a.State {
	case "pending", "consumed", "superseded", "not-applicable":
	default:
		return fmt.Errorf("application state must be pending, consumed, superseded, or not-applicable")
	}
	switch decision {
	case "approve":
		if a.Action != "advance" || strings.TrimSpace(a.TargetStage) == "" || a.State == "not-applicable" {
			return fmt.Errorf("approve application must be advance with a target stage")
		}
	case "revise":
		if a.Action != "feedback" || strings.TrimSpace(a.TargetStage) == "" || a.State == "not-applicable" {
			return fmt.Errorf("revise application must be feedback with a target stage")
		}
	case "hold":
		if a.Action != "none" || a.TargetStage != "" || a.State != "not-applicable" {
			return fmt.Errorf("hold application must be none/not-applicable")
		}
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
		if r.ID != doc.Current.Gate || len(r.Attempts) == 0 {
			continue
		}
		a := &r.Attempts[len(r.Attempts)-1]
		s := Summary{Gate: r.ID, Attempt: a.ID, State: attemptState(a), Briefing: a.Briefing.ID}
		if a.Resolution != nil {
			s.Resolution, s.Decision = a.Resolution.ID, a.Resolution.Decision
		}
		if a.Application != nil {
			s.Application = a.Application.Action + "/" + a.Application.State
			s.ApplicationState = a.Application.State
			s.TargetStage = a.Application.TargetStage
		}
		return s
	}
	return Summary{}
}

func attemptState(a *Attempt) string {
	if a.Resolution != nil {
		return "closed"
	}
	return "open"
}
