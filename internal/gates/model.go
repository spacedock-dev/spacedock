// ABOUTME: Canonical v1 durable gate-resolution and one-use application model.
// ABOUTME: Canonical validation keeps unknown or conflicting application state fail-closed.
package gates

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var digestRE = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
var roundStageRE = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

type Document struct {
	Version int          `yaml:"version" json:"version"`
	Records []GateRecord `yaml:"records" json:"records"`
}

type GateRecord struct {
	ID       string    `yaml:"id" json:"id"`
	Stage    string    `yaml:"stage" json:"stage"`
	Attempts []Attempt `yaml:"attempts" json:"attempts"`
}

type Attempt struct {
	ID          string       `yaml:"id" json:"id"`
	Briefing    Briefing     `yaml:"briefing" json:"briefing"`
	Withdrawal  *Withdrawal  `yaml:"withdrawal,omitempty" json:"withdrawal,omitempty"`
	Resolution  *Resolution  `yaml:"resolution,omitempty" json:"resolution,omitempty"`
	Application *Application `yaml:"application,omitempty" json:"application,omitempty"`
}

type Withdrawal struct {
	By     string `yaml:"by" json:"by"`
	At     string `yaml:"at" json:"at"`
	Reason string `yaml:"reason" json:"reason"`
}

type Application struct {
	TargetStage string `yaml:"target-stage" json:"target-stage"`
	State       string `yaml:"state" json:"state"`
}

type Briefing struct {
	ID            string `yaml:"id" json:"id"`
	Digest        string `yaml:"digest" json:"digest"`
	RequestDigest string `yaml:"request-digest,omitempty" json:"request-digest,omitempty"`
	RoomRef       string `yaml:"room-ref" json:"room-ref"`
}

type RoundPointer struct {
	ID, Stage string
	Cycle     int
	Briefing  Briefing
}

type RoundEntrySummary struct {
	Type, ID, Decision string
	Advisory           bool
}

type RoundSummary struct {
	ID, Stage, Briefing string
	Cycle               int
	Entries             []RoundEntrySummary
}

type Resolution struct {
	Type     string   `yaml:"type" json:"type"`
	ID       string   `yaml:"id" json:"id"`
	Briefing string   `yaml:"briefing" json:"briefing"`
	By       string   `yaml:"by" json:"by"`
	At       string   `yaml:"at" json:"at"`
	Decision string   `yaml:"decision" json:"decision"`
	Reason   string   `yaml:"reason,omitempty" json:"reason,omitempty"`
	Conn     *Conn    `yaml:"conn,omitempty" json:"conn,omitempty"`
	Includes []string `yaml:"includes,omitempty" json:"includes,omitempty"`
}

// Conn cites the delegated authority a First Officer resolution acted under:
// the grant quoted verbatim from the conversation and where it was given. It
// attributes an FO-rendered decision back to its authority; it never
// authenticates the grant itself (the binary checks form and disjointness
// only) and it confers no authorization on its own — auto-continue's boundary
// negative proves a citation on a no-conn journey still reds. Present only on
// agent:first-officer resolutions; a person:captain decision cites no grant,
// so the two record shapes stay disjoint.
type Conn struct {
	Quote  string `yaml:"quote" json:"quote"`
	Source string `yaml:"source" json:"source"`
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

// RouteApprovedAwaitingMerge is the readiness vocabulary for an approval
// whose target stage is terminal: consume spends nothing and writes no
// status; the terminal merge ceremony (merge guard) is the sole terminal
// consumer and spends the still-pending approval with delivery proof.
// CurrentStageReadiness projects it; CLI consume reporting reuses it.
const RouteApprovedAwaitingMerge = "approved-awaiting-merge"

// RouteNeedsPreparation is the scheduler-only readiness value for a gated
// current stage whose report is mechanically complete but has no current
// attempt authority. It is a candidate for First Officer semantic review, not
// an approval or a permission to mutate the entity.
const RouteNeedsPreparation = "needs-preparation"

type ConsumeResult struct {
	Eligibility
	Consumed bool
	// Wrote reports whether THIS call wrote a real mutation (a fresh advance or
	// a fresh pending->superseded supersede) — distinct from ApplicationState,
	// which EvaluateEligibility always copies from the attempt's current
	// application state, including on a pure-refusal read against an
	// already-superseded or already-consumed application. Callers deciding
	// whether to sync/commit must branch on Wrote, never on
	// ApplicationState == "superseded" or "consumed" alone.
	Wrote bool
}

// ReadinessStage is the minimum workflow taxonomy needed to reduce a selected
// gate attempt into a local scheduling state.
type ReadinessStage struct {
	Name     string
	Gate     bool
	Terminal bool
}

// CurrentStageReadiness projects validated durable gate state without reading
// a retained room or entity body. Unknown and inconsistent states fail closed.
func CurrentStageReadiness(doc *Document, status string, stages []ReadinessStage) string {
	stageByName := make(map[string]ReadinessStage, len(stages))
	for _, stage := range stages {
		stageByName[stage.Name] = stage
	}
	current, ok := stageByName[status]
	if !ok || !current.Gate || current.Terminal {
		return ""
	}
	if doc == nil {
		return "validating"
	}
	record, err := recordForStage(doc, status)
	if err != nil {
		if strings.Contains(err.Error(), "no logical gate") {
			return "validating"
		}
		return "invalid"
	}
	if len(record.Attempts) == 0 {
		return "validating"
	}
	attempt := &record.Attempts[len(record.Attempts)-1]
	if attempt.Briefing.ID == "" || attempt.Briefing.Digest == "" || attempt.Briefing.RoomRef == "" {
		return "invalid"
	}
	switch attemptState(attempt) {
	case "open":
		return "awaiting-captain"
	case "withdrawn":
		return "withdrawn-awaiting-prepare"
	case "closed":
	default:
		return "invalid"
	}
	app := attempt.Application
	if app == nil {
		switch attempt.Resolution.Decision {
		case "hold":
			return "not-applicable"
		case "revise":
			return "feedback-pending"
		default:
			return "invalid"
		}
	}
	switch app.State {
	case "consumed", "superseded":
		return app.State
	case "pending":
	default:
		return "invalid"
	}
	if attempt.Resolution.Decision != "approve" ||
		strings.TrimSpace(app.TargetStage) == "" || app.TargetStage == status {
		return "invalid"
	}
	target, ok := stageByName[app.TargetStage]
	if !ok {
		return "invalid"
	}
	if target.Terminal {
		return RouteApprovedAwaitingMerge
	}
	return "approved-awaiting-advance"
}

// CurrentStageReadinessWithReport extends CurrentStageReadiness with the
// status-owned mechanical promotion proof. The caller owns what proves
// promotion: a complete committed stage report at an ordinary gated stage, or
// the committed clean seed itself at an initial one, which had no prior stage
// to write a report. A satisfied proof promotes only the no-current-authority
// shape; every selected, malformed, or otherwise classified authority keeps
// the canonical result above. Keeping this wrapper beside the reducer lets
// status and any future machine projection share the exact ordering and
// fail-closed behavior without adding durable gate state.
func CurrentStageReadinessWithReport(doc *Document, status string, stages []ReadinessStage, promotionProven bool) string {
	readiness := CurrentStageReadiness(doc, status, stages)
	if !promotionProven || readiness != "validating" {
		return readiness
	}
	if doc == nil {
		return RouteNeedsPreparation
	}
	if len(doc.Records) == 0 {
		return "invalid"
	}
	if _, err := recordForStage(doc, status); err != nil && strings.Contains(err.Error(), "no logical gate") {
		return RouteNeedsPreparation
	}
	return readiness
}

func Validate(doc *Document) error {
	if doc.Version != 1 {
		return fmt.Errorf("gates.version must be 1")
	}
	gateIDs, stages, attemptIDs, briefingIDs, resolutionIDs := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}
	for ri := range doc.Records {
		r := &doc.Records[ri]
		pendingApplications := 0
		if r.ID == "" || r.Stage == "" || gateIDs[r.ID] || len(r.Attempts) == 0 {
			return fmt.Errorf("record %d has missing or duplicate identity or no attempts", ri+1)
		}
		gateIDs[r.ID] = true
		if stages[r.Stage] {
			return fmt.Errorf("multiple logical gates claim workflow stage %s", r.Stage)
		}
		stages[r.Stage] = true
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
			if a.Briefing.RequestDigest != "" && !digestRE.MatchString(a.Briefing.RequestDigest) {
				return fmt.Errorf("attempt %s has invalid request-digest", a.ID)
			}
			switch attemptState(a) {
			case "open":
				if a.Application != nil {
					return fmt.Errorf("open attempt %s cannot carry application data", a.ID)
				}
				continue
			case "withdrawn":
				if a.Briefing.RequestDigest == "" {
					return fmt.Errorf("withdrawn attempt %s must retain a request-digest", a.ID)
				}
				if err := validateWithdrawal(a.Withdrawal); err != nil {
					return fmt.Errorf("attempt %s: %w", a.ID, err)
				}
				if a.Application != nil {
					return fmt.Errorf("withdrawn attempt %s cannot carry application data", a.ID)
				}
				continue
			case "closed":
			default:
				return fmt.Errorf("attempt %s has conflicting withdrawal and resolution state", a.ID)
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
	return nil
}

func validateWithdrawal(withdrawal *Withdrawal) error {
	if withdrawal == nil || withdrawal.By != "agent:first-officer" {
		return fmt.Errorf("withdrawal attribution must be agent:first-officer")
	}
	if parsed, err := time.Parse(time.RFC3339Nano, withdrawal.At); err != nil || parsed.Location() != time.UTC {
		return fmt.Errorf("withdrawal timestamp must be RFC3339Nano UTC")
	}
	if strings.TrimSpace(withdrawal.Reason) == "" {
		return fmt.Errorf("withdrawal reason must be nonblank")
	}
	return nil
}

func validateApplication(a *Application, decision string) error {
	switch a.State {
	case "pending", "consumed", "superseded":
	default:
		return fmt.Errorf("application state must be pending, consumed, or superseded")
	}
	if decision != "approve" {
		return fmt.Errorf("only approve resolutions may carry an application")
	}
	if strings.TrimSpace(a.TargetStage) == "" {
		return fmt.Errorf("approve application must have a target stage")
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
	if r.Conn != nil {
		// Write-side requires the citation only on new agent:first-officer chat
		// closes (recordChatLocked); this read-side check is the disjointness
		// half — a conn citation is only ever valid attribution evidence on an
		// FO-rendered decision. A hand-forged captain-plus-citation record fails
		// here, and every durable read (including the graders) inherits the
		// refusal. Historical FO resolutions written before this field existed
		// carry no conn: block at all and are unaffected — this only rejects a
		// conn: block that is PRESENT on the wrong actor.
		if r.By != "agent:first-officer" {
			return fmt.Errorf("conn citation is only valid on an agent:first-officer resolution, not %s", r.By)
		}
		if strings.TrimSpace(r.Conn.Quote) == "" || strings.TrimSpace(r.Conn.Source) == "" {
			return fmt.Errorf("conn citation requires a nonblank quote and source")
		}
	}
	return nil
}

func CurrentSummary(doc *Document, stage ...string) Summary {
	var r *GateRecord
	if len(stage) > 0 {
		r, _ = recordForStage(doc, stage[0])
	} else if len(doc.Records) == 1 {
		r = &doc.Records[0]
	}
	if r != nil && len(r.Attempts) > 0 {
		a := &r.Attempts[len(r.Attempts)-1]
		s := Summary{Gate: r.ID, Attempt: a.ID, State: attemptState(a), Briefing: a.Briefing.ID}
		if a.Resolution != nil {
			s.Resolution, s.Decision = a.Resolution.ID, a.Resolution.Decision
		}
		if a.Application != nil {
			s.Application = "advance/" + a.Application.State
			s.ApplicationState = a.Application.State
			s.TargetStage = a.Application.TargetStage
		}
		return s
	}
	return Summary{}
}

func attemptState(a *Attempt) string {
	switch {
	case a.Withdrawal == nil && a.Resolution == nil:
		return "open"
	case a.Withdrawal != nil && a.Resolution == nil:
		return "withdrawn"
	case a.Withdrawal == nil && a.Resolution != nil:
		return "closed"
	default:
		return "invalid"
	}
}
