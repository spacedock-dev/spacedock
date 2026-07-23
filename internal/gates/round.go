package gates

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

type loadedRound struct {
	Briefing, Log  []byte
	Manifest       *briefingManifest
	Review         reviewLog
	Triage, Digest string
}
type roundLocation struct {
	stage, room, entityID string
	cycle                 int
	entity                []byte
	pointer               RoundPointer
}

func resolveRound(entityPath, spec string) (roundLocation, error) {
	stage, rawCycle, ok := strings.Cut(spec, "/")
	cycle, err := strconv.Atoi(rawCycle)
	if !ok || !roundStageRE.MatchString(stage) || err != nil || cycle < 1 || strconv.Itoa(cycle) != rawCycle {
		return roundLocation{}, fmt.Errorf("--round must be a normalized STAGE/positive-cycle")
	}
	result := roundLocation{stage: stage, cycle: cycle}
	if result.entity, err = os.ReadFile(entityPath); err != nil {
		return result, err
	}
	if result.pointer, err = readRoundPointerData(result.entity); err != nil {
		return result, err
	}
	root, _, _, _ := frontmatterNode(result.entity)
	id := mappingValue(root, "id")
	if id == nil || strings.TrimSpace(id.Value) == "" {
		return result, fmt.Errorf("entity has no identity for advisory round")
	}
	result.entityID = id.Value
	if result.pointer.ID != "" && result.pointer.ID != fmt.Sprintf("round:%s:%s:%d", id.Value, result.pointer.Stage, result.pointer.Cycle) {
		return result, fmt.Errorf("review-round identity does not match the entity")
	}
	result.room = filepath.Join(filepath.Dir(entityPath), "review", stage, fmt.Sprintf("round-%d", cycle))
	for _, parent := range []string{filepath.Dir(result.room), filepath.Dir(filepath.Dir(result.room))} {
		if info, statErr := os.Lstat(parent); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
			return result, fmt.Errorf("derived round room crosses symlink %s", parent)
		}
	}
	return result, nil
}
func loadValidateRound(entityDir, room, briefingPath, logPath string, want *Briefing) (loadedRound, error) {
	var result loadedRound
	var err error
	if result.Briefing, err = os.ReadFile(briefingPath); err != nil {
		return result, err
	}
	if result.Manifest, err = parseBriefingManifest(result.Briefing); err != nil {
		return result, err
	}
	result.Digest, err = CanonicalDigest(result.Briefing)
	if err != nil {
		return result, err
	}
	if want != nil && (want.ID != result.Manifest.ID || want.Digest != result.Digest || want.DigestDomain != "canonical-bytes") {
		return result, fmt.Errorf("review-round pointer does not bind the retained Briefing")
	}
	if err := verifyRoundArtifacts(entityDir, room, result.Manifest); err != nil {
		return result, err
	}
	if result.Log, err = os.ReadFile(logPath); err != nil {
		return result, err
	}
	if result.Review, err = parseReviewLog(result.Log, result.Manifest.ID); err != nil {
		return result, err
	}
	result.Triage, err = classifyCompletedRound(result.Review)
	if err != nil {
		return result, err
	}
	return result, nil
}
func recordRoundLocked(entityPath string, input RecordInput) error {
	return recordRoundLockedWith(entityPath, input, nil, atomicWrite)
}
func recordRoundLockedWith(entityPath string, input RecordInput, beforePublish func(string), replace func(string, []byte) error) error {
	location, err := resolveRound(entityPath, input.Round)
	if err != nil {
		return err
	}
	inputRound, err := loadValidateRound(filepath.Dir(entityPath), location.room, input.BriefingPath, input.LogPath, nil)
	if err != nil {
		return err
	}
	pointer := RoundPointer{ID: fmt.Sprintf("round:%s:%s:%d", location.entityID, location.stage, location.cycle), Stage: location.stage, Cycle: location.cycle,
		Briefing: Briefing{ID: inputRound.Manifest.ID, Digest: inputRound.Digest, DigestDomain: "canonical-bytes",
			RoomRef: fmt.Sprintf("./review/%s/round-%d", location.stage, location.cycle)}}
	if location.pointer.ID == pointer.ID {
		if _, statErr := os.Lstat(location.room); os.IsNotExist(statErr) {
			return fmt.Errorf("round identity already has a pointer without its immutable room")
		}
	}
	line := ""
	project := inputRound.Triage != "no-findings"
	if project {
		if input.FeedbackCyclePath == "" {
			return fmt.Errorf("triaged round requires --feedback-cycle")
		}
		if line, err = readFeedbackCycle(input.FeedbackCyclePath, location.cycle); err != nil {
			return err
		}
	} else if input.FeedbackCyclePath != "" {
		return fmt.Errorf("no-findings round does not accept --feedback-cycle")
	}
	nextRoom := roundRoomBytes{Exists: true, Briefing: inputRound.Briefing, Log: inputRound.Log}
	if beforePublish != nil {
		beforePublish(location.room)
	}
	return publishRound(location.room, nextRoom, func(replay bool) error {
		return mutateEntity(entityPath, entityExpectation{Bytes: location.entity}, func(entity []byte) ([]byte, error) {
			next, err := rebuildRoundEntity(entity, pointer, line, location.cycle, project)
			if err == nil && replay && !bytes.Equal(next, entity) {
				return nil, fmt.Errorf("immutable round replay does not match the entity pointer and projection")
			}
			return next, err
		}, replace)
	})
}
func ValidateRoundFile(entityPath, spec string) (RoundSummary, error) {
	location, err := resolveRound(entityPath, spec)
	if err != nil {
		return RoundSummary{}, err
	}
	if location.pointer.ID == "" || location.pointer.Stage != location.stage || location.pointer.Cycle != location.cycle {
		return RoundSummary{}, fmt.Errorf("entity current review-round pointer does not resolve %s", spec)
	}
	if _, err := readRoundRoom(location.room); err != nil {
		return RoundSummary{}, err
	}
	loaded, err := loadValidateRound(filepath.Dir(entityPath), location.room, filepath.Join(location.room, "briefing.json"), filepath.Join(location.room, "briefing.review.jsonl"), &location.pointer.Briefing)
	if err != nil {
		return RoundSummary{}, err
	}
	summary := RoundSummary{ID: location.pointer.ID, Stage: location.stage, Cycle: location.cycle, Briefing: loaded.Manifest.ID, Triage: loaded.Triage}
	for _, entry := range loaded.Review.Entries {
		item := RoundEntrySummary{Type: entry.Type, ID: entry.ID, Advisory: entry.Resolution != nil}
		if entry.Resolution != nil {
			item.Decision = entry.Resolution.Decision
		}
		summary.Entries = append(summary.Entries, item)
	}
	return summary, nil
}
func readRoundPointerData(data []byte) (RoundPointer, error) {
	root, _, _, err := frontmatterNode(data)
	if err != nil {
		return RoundPointer{}, err
	}
	node := mappingValue(root, "review-round")
	if node == nil {
		return RoundPointer{}, nil
	}
	var pointer RoundPointer
	briefingNode := mappingValue(node, "briefing")
	if len(node.Content) != 8 || briefingNode == nil || len(briefingNode.Content) != 8 {
		return RoundPointer{}, fmt.Errorf("entity has invalid review-round pointer")
	}
	err = node.Decode(&pointer)
	wantRoom := fmt.Sprintf("./review/%s/round-%d", pointer.Stage, pointer.Cycle)
	if err != nil || pointer.ID == "" || !roundStageRE.MatchString(pointer.Stage) || pointer.Cycle < 1 ||
		pointer.Briefing.ID == "" || !digestRE.MatchString(pointer.Briefing.Digest) ||
		pointer.Briefing.DigestDomain != "canonical-bytes" || pointer.Briefing.RoomRef != wantRoom {
		return RoundPointer{}, fmt.Errorf("entity has invalid review-round pointer")
	}
	return pointer, nil
}
func rebuildRoundEntity(original []byte, pointer RoundPointer, projection string, cycle int, project bool) ([]byte, error) {
	block, err := yaml.Marshal(struct {
		Round RoundPointer `yaml:"review-round"`
	}{pointer})
	if err != nil {
		return nil, err
	}
	out, err := replaceTopLevels(original, true, topLevelReplacement{key: "review-round", data: block})
	if err != nil {
		return nil, err
	}
	out, err = spliceFeedbackCycle(out, projection, cycle, project)
	if err != nil {
		return nil, err
	}
	got, err := readRoundPointerData(out)
	if err != nil || got != pointer {
		return nil, fmt.Errorf("validate rebuilt review-round pointer")
	}
	if _, _, gateErr := readData(out); gateErr != nil && !strings.Contains(gateErr.Error(), "no gates record") {
		return nil, fmt.Errorf("validate rebuilt gates: %w", gateErr)
	}
	return out, nil
}
func verifyRoundArtifacts(entityDir, room string, manifest *briefingManifest) error {
	realRoot, err := filepath.EvalSymlinks(entityDir)
	if err != nil {
		return err
	}
	for _, artifact := range manifest.Artifacts {
		parsed, err := url.Parse(artifact.URI)
		if err != nil {
			return fmt.Errorf("artifact %s URI: %w", artifact.ID, err)
		}
		if parsed.Scheme != "" {
			continue
		}
		path := filepath.Join(room, filepath.FromSlash(parsed.Path))
		realPath, err := filepath.EvalSymlinks(path)
		rel, relErr := filepath.Rel(realRoot, realPath)
		if err != nil || relErr != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return fmt.Errorf("artifact %s escapes or cannot resolve within entity", artifact.ID)
		}
		if realPath == filepath.Join(realRoot, "index.md") {
			return fmt.Errorf("artifact %s resolves to the mutable entity file", artifact.ID)
		}
		body, err := os.ReadFile(realPath)
		if err != nil || RawDigest(body) != artifact.Rev {
			return fmt.Errorf("artifact %s raw digest does not match Briefing revision", artifact.ID)
		}
	}
	return nil
}
func readFeedbackCycle(path string, cycle int) (string, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(body) {
		return "", fmt.Errorf("--feedback-cycle must be UTF-8")
	}
	line := strings.TrimSuffix(string(body), "\n")
	if strings.ContainsAny(line, "\r\n") || !strings.HasPrefix(line, fmt.Sprintf("- Cycle %d: ", cycle)) {
		return "", fmt.Errorf("--feedback-cycle must contain exactly one canonical Cycle %d line", cycle)
	}
	return line, nil
}
