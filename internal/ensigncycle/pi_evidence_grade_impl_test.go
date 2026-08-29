package ensigncycle

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const piFrontDoorProofID = "pi-front-door-child-durable-boot-contract"

type piBootContractEvidence struct {
	Agent                 string
	Skills                []string
	DispatchFileForwarded bool
	ReadCallCount         int
	EnsignSkillReadRank   int
	FirstOfficerReads     int
	Transcript            string
}

type piSessionEvidence struct {
	Model      string   `json:"model"`
	DurationMS int64    `json:"duration_ms"`
	CostUSD    *float64 `json:"cost_usd,omitempty"`
}

type piEvidenceClaims struct {
	FrontDoor     bool `json:"front_door"`
	ChildDispatch bool `json:"child_dispatch"`
	DurableOutput bool `json:"durable_output"`
	BootContract  bool `json:"boot_contract"`
}

type piSpawnGrade struct {
	Agent      string   `json:"agent"`
	Skills     []string `json:"skills"`
	Transcript string   `json:"transcript"`
	ReadCalls  int      `json:"read_calls"`
	EnsignRank int      `json:"ensign_skill_read_rank"`
	FOReads    int      `json:"first_officer_reads"`
}

type piFrontDoorEvidenceGrade struct {
	ProofID string            `json:"proof_id"`
	Verdict string            `json:"verdict"`
	Claims  piEvidenceClaims  `json:"claims"`
	Root    piSessionEvidence `json:"root_session"`
	Child   piSessionEvidence `json:"child_session"`
	Spawn   piSpawnGrade      `json:"spawn_and_boot_contract"`
}

func buildPiFrontDoorEvidenceGrade(rootSessionPath, childSessionPath string, durableOutput bool, boot piBootContractEvidence) (piFrontDoorEvidenceGrade, error) {
	root, err := readPiSessionEvidence(rootSessionPath)
	if err != nil {
		return piFrontDoorEvidenceGrade{}, fmt.Errorf("grade Pi front door: %w", err)
	}
	child, err := readPiSessionEvidence(childSessionPath)
	if err != nil {
		return piFrontDoorEvidenceGrade{}, fmt.Errorf("grade Pi child dispatch: %w", err)
	}
	if !durableOutput {
		return piFrontDoorEvidenceGrade{}, fmt.Errorf("grade Pi durable output: entity mutation and commit were not observed")
	}
	if boot.Agent == "" || len(boot.Skills) == 0 || !boot.DispatchFileForwarded {
		return piFrontDoorEvidenceGrade{}, fmt.Errorf("grade Pi child dispatch: build artifact agent, skill, or dispatch pointer was not forwarded")
	}
	if boot.EnsignSkillReadRank < 1 || boot.EnsignSkillReadRank > 5 {
		return piFrontDoorEvidenceGrade{}, fmt.Errorf("grade Pi boot contract: ensign skill read rank %d is outside the first five reads", boot.EnsignSkillReadRank)
	}
	if boot.FirstOfficerReads != 0 {
		return piFrontDoorEvidenceGrade{}, fmt.Errorf("grade Pi boot contract: child made %d first-officer reads", boot.FirstOfficerReads)
	}

	return piFrontDoorEvidenceGrade{
		ProofID: piFrontDoorProofID,
		Verdict: "pass",
		Claims: piEvidenceClaims{
			FrontDoor:     true,
			ChildDispatch: true,
			DurableOutput: true,
			BootContract:  true,
		},
		Root:  root,
		Child: child,
		Spawn: piSpawnGrade{
			Agent:      boot.Agent,
			Skills:     append([]string(nil), boot.Skills...),
			Transcript: boot.Transcript,
			ReadCalls:  boot.ReadCallCount,
			EnsignRank: boot.EnsignSkillReadRank,
			FOReads:    boot.FirstOfficerReads,
		},
	}, nil
}

func readPiSessionEvidence(path string) (piSessionEvidence, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return piSessionEvidence{}, fmt.Errorf("read session %s: %w", path, err)
	}
	var first, last time.Time
	model := ""
	cost := 0.0
	hasCost := false
	for lineNo, line := range strings.Split(string(body), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var event struct {
			Type      string `json:"type"`
			Timestamp string `json:"timestamp"`
			Provider  string `json:"provider"`
			ModelID   string `json:"modelId"`
			Message   struct {
				Role  string `json:"role"`
				Model string `json:"model"`
				Usage struct {
					Cost struct {
						Total *float64 `json:"total"`
					} `json:"cost"`
				} `json:"usage"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return piSessionEvidence{}, fmt.Errorf("session %s line %d is not JSON: %w", path, lineNo+1, err)
		}
		if event.Timestamp != "" {
			at, err := time.Parse(time.RFC3339Nano, event.Timestamp)
			if err != nil {
				return piSessionEvidence{}, fmt.Errorf("session %s line %d has invalid timestamp: %w", path, lineNo+1, err)
			}
			if first.IsZero() || at.Before(first) {
				first = at
			}
			if last.IsZero() || at.After(last) {
				last = at
			}
		}
		if event.Type == "model_change" && event.ModelID != "" {
			model = event.ModelID
			if event.Provider != "" {
				model = event.Provider + "/" + event.ModelID
			}
		} else if model == "" && event.Message.Model != "" {
			model = event.Message.Model
		}
		if event.Message.Role == "assistant" && event.Message.Usage.Cost.Total != nil {
			cost += *event.Message.Usage.Cost.Total
			hasCost = true
		}
	}
	if first.IsZero() || last.IsZero() {
		return piSessionEvidence{}, fmt.Errorf("session %s has no timestamped records", path)
	}
	if model == "" {
		return piSessionEvidence{}, fmt.Errorf("session %s has no model record", path)
	}
	evidence := piSessionEvidence{Model: model, DurationMS: last.Sub(first).Milliseconds()}
	if hasCost {
		evidence.CostUSD = &cost
	}
	return evidence, nil
}

func TestPiFrontDoorEvidenceGradeIncludesClaimsModelsDurationsAndCosts(t *testing.T) {
	grade, err := buildPiFrontDoorEvidenceGrade(
		filepath.Join("testdata", "pi_front_door_grade", "root-session.jsonl"),
		filepath.Join("testdata", "pi_front_door_grade", "child-session.jsonl"),
		true,
		validPiBootContractEvidence(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if grade.ProofID != piFrontDoorProofID || grade.Verdict != "pass" {
		t.Fatalf("grade identity = %q/%q, want %q/pass", grade.ProofID, grade.Verdict, piFrontDoorProofID)
	}
	if !grade.Claims.FrontDoor || !grade.Claims.ChildDispatch || !grade.Claims.DurableOutput || !grade.Claims.BootContract {
		t.Fatalf("grade does not own all four Pi claims: %+v", grade.Claims)
	}
	if grade.Root.Model != "openai/gpt-5.4" || grade.Child.Model != "openai/gpt-5.4" {
		t.Fatalf("graded models = root %q child %q", grade.Root.Model, grade.Child.Model)
	}
	if grade.Root.DurationMS != 104808 || grade.Child.DurationMS != 98500 {
		t.Fatalf("graded durations = root %dms child %dms", grade.Root.DurationMS, grade.Child.DurationMS)
	}
	if grade.Root.CostUSD == nil || math.Abs(*grade.Root.CostUSD-0.15) > 0.000001 {
		t.Fatalf("root cost = %v, want 0.15", grade.Root.CostUSD)
	}
	if grade.Child.CostUSD == nil || math.Abs(*grade.Child.CostUSD-0.127493) > 0.000001 {
		t.Fatalf("child cost = %v, want 0.127493", grade.Child.CostUSD)
	}
	if grade.Spawn.Agent != "worker" || len(grade.Spawn.Skills) != 1 || grade.Spawn.Skills[0] != "ensign" {
		t.Fatalf("spawn grade = %+v", grade.Spawn)
	}
}

func TestPiFrontDoorEvidenceGradeRejectsMissingClaimGraders(t *testing.T) {
	root := filepath.Join("testdata", "pi_front_door_grade", "root-session.jsonl")
	child := filepath.Join("testdata", "pi_front_door_grade", "child-session.jsonl")
	for name, mutate := range map[string]func(*bool, *piBootContractEvidence){
		"durable output":  func(durable *bool, _ *piBootContractEvidence) { *durable = false },
		"child dispatch":  func(_ *bool, boot *piBootContractEvidence) { boot.DispatchFileForwarded = false },
		"boot skill read": func(_ *bool, boot *piBootContractEvidence) { boot.EnsignSkillReadRank = 0 },
		"boot read order": func(_ *bool, boot *piBootContractEvidence) { boot.EnsignSkillReadRank = 6 },
		"boot isolation":  func(_ *bool, boot *piBootContractEvidence) { boot.FirstOfficerReads = 1 },
	} {
		t.Run(name, func(t *testing.T) {
			durable := true
			boot := validPiBootContractEvidence()
			mutate(&durable, &boot)
			if _, err := buildPiFrontDoorEvidenceGrade(root, child, durable, boot); err == nil {
				t.Fatal("Pi evidence grade accepted a removed or falsified grader")
			}
		})
	}
}

func validPiBootContractEvidence() piBootContractEvidence {
	return piBootContractEvidence{
		Agent:                 "worker",
		Skills:                []string{"ensign"},
		DispatchFileForwarded: true,
		ReadCallCount:         2,
		EnsignSkillReadRank:   1,
		FirstOfficerReads:     0,
		Transcript:            "/run/child-transcript.jsonl",
	}
}
