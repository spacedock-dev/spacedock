package journeymetrics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Ledger struct {
	SchemaVersion int                   `json:"schema_version"`
	Release       ReleaseMetadata       `json:"release"`
	GeneratedAt   string                `json:"generated_at"`
	Summary       LedgerSummary         `json:"summary"`
	Scenarios     []ScenarioLedgerEntry `json:"scenarios"`
}

type ReleaseMetadata struct {
	Version  string `json:"version"`
	Artifact string `json:"artifact"`
}

type LedgerSummary struct {
	ScenarioCount      int `json:"scenario_count"`
	ObservationCount   int `json:"observation_count"`
	MeasuredCount      int `json:"measured_count"`
	CharacterizedCount int `json:"characterized_count"`
}

type ScenarioLedgerEntry struct {
	ScenarioID   string   `json:"scenario_id"`
	Source       string   `json:"source,omitempty"`
	Observations []Record `json:"observations"`
}

func AggregateLedger(releaseVersion string, records []Record, generatedAt time.Time) (Ledger, error) {
	version := normalizeVersion(releaseVersion)
	if version == "" {
		return Ledger{}, fmt.Errorf("release version is required")
	}
	if len(records) == 0 {
		return Ledger{}, fmt.Errorf("no accepted journey metric records")
	}
	records = append([]Record(nil), records...)
	for i := range records {
		records[i] = normalizeRecord(records[i])
		if records[i].ScenarioID == "" {
			return Ledger{}, fmt.Errorf("record %d missing scenario_id", i+1)
		}
		if records[i].MetricsState == "" {
			return Ledger{}, fmt.Errorf("record %s missing metrics_state", records[i].ScenarioID)
		}
		if records[i].Outcome.Status == "" {
			return Ledger{}, fmt.Errorf("record %s missing outcome.status", records[i].ScenarioID)
		}
	}
	sort.Slice(records, func(i, j int) bool {
		return recordSortKey(records[i]) < recordSortKey(records[j])
	})
	byScenario := map[string][]Record{}
	scenarioSources := map[string]string{}
	for _, record := range records {
		byScenario[record.ScenarioID] = append(byScenario[record.ScenarioID], record)
		if scenarioSources[record.ScenarioID] == "" {
			scenarioSources[record.ScenarioID] = record.Source
		}
	}
	scenarios := make([]ScenarioLedgerEntry, 0, len(byScenario))
	for scenarioID, observations := range byScenario {
		scenarios = append(scenarios, ScenarioLedgerEntry{
			ScenarioID:   scenarioID,
			Source:       scenarioSources[scenarioID],
			Observations: observations,
		})
	}
	sort.Slice(scenarios, func(i, j int) bool {
		return scenarios[i].ScenarioID < scenarios[j].ScenarioID
	})
	summary := LedgerSummary{
		ScenarioCount:    len(scenarios),
		ObservationCount: len(records),
	}
	for _, record := range records {
		switch record.MetricsState {
		case StateMeasured:
			summary.MeasuredCount++
		case StateCharacterized:
			summary.CharacterizedCount++
		}
	}
	return Ledger{
		SchemaVersion: LedgerSchemaVersion,
		Release: ReleaseMetadata{
			Version:  version,
			Artifact: "journey-costs-v" + version + ".json",
		},
		GeneratedAt: generatedAt.UTC().Format(time.RFC3339),
		Summary:     summary,
		Scenarios:   scenarios,
	}, nil
}

func MarshalLedger(ledger Ledger) ([]byte, error) {
	data, err := json.MarshalIndent(ledger, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func ReadRecordsDir(dir string) ([]Record, error) {
	var records []Record
	if err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var record Record
		if err := json.Unmarshal(data, &record); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		record = normalizeRecord(record)
		if record.ScenarioID == "" || record.MetricsState == "" || record.Outcome.Status == "" {
			return nil
		}
		records = append(records, record)
		return nil
	}); err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("no accepted journey metric records in %s", dir)
	}
	return records, nil
}

func normalizeVersion(version string) string {
	return strings.TrimPrefix(strings.TrimSpace(version), "v")
}

func recordSortKey(record Record) string {
	return strings.Join([]string{
		record.ScenarioID,
		record.Runtime,
		record.Executor,
		record.Host,
		record.Mode,
		record.Model,
		string(record.MetricsState),
	}, "\x00")
}
