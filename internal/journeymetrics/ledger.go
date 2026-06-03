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
	SchemaVersion int             `json:"schema_version"`
	Release       ReleaseMetadata `json:"release"`
	GeneratedAt   string          `json:"generated_at"`
	Summary       LedgerSummary   `json:"summary"`
	Journeys      []Record        `json:"journeys"`
}

type ReleaseMetadata struct {
	Version  string `json:"version"`
	Artifact string `json:"artifact"`
}

type LedgerSummary struct {
	JourneyCount       int `json:"journey_count"`
	MeasuredCount      int `json:"measured_count"`
	CharacterizedCount int `json:"characterized_count"`
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
		if records[i].SchemaVersion == 0 {
			records[i].SchemaVersion = RecordSchemaVersion
		}
		if records[i].JourneyID == "" {
			return Ledger{}, fmt.Errorf("record %d missing journey_id", i+1)
		}
		if records[i].MetricsState == "" {
			return Ledger{}, fmt.Errorf("record %s missing metrics_state", records[i].JourneyID)
		}
		if records[i].Outcome.Status == "" {
			return Ledger{}, fmt.Errorf("record %s missing outcome.status", records[i].JourneyID)
		}
		records[i].Tokens = records[i].Tokens.withTotal()
		records[i].ModelUsage = normalizeModelUsage(records[i].ModelUsage)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].JourneyID < records[j].JourneyID
	})
	summary := LedgerSummary{JourneyCount: len(records)}
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
		Journeys:    records,
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
		if record.JourneyID == "" || record.MetricsState == "" || record.Outcome.Status == "" {
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
