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
)

type loadedRound struct {
	Briefing, Log []byte
	Manifest      *briefingManifest
	Review        reviewLog
	Digest        string
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
		return result, fmt.Errorf("entity has no identity for correction round")
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
	if want != nil && (want.ID != result.Manifest.ID || want.Digest != result.Digest) {
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
	return result, nil
}
func recordRoundLockedWith(entityPath string, input RecordInput, beforePublish func(string), replace func(string, []byte) error) error {
	location, err := resolveRound(entityPath, input.Round)
	if err != nil {
		return err
	}
	if _, err := applicationForDecision(entityPath, input.WorkflowDir, location.stage, "revise"); err != nil {
		return err
	}
	inputRound, err := loadValidateRound(filepath.Dir(entityPath), location.room, input.BriefingPath, input.LogPath, nil)
	if err != nil {
		return err
	}
	if !roundLogClosed(inputRound.Review) {
		return fmt.Errorf("round %s is not closed: %s. Nothing was recorded; the room is immutable, so an open log would truncate the record permanently. Route the correction, append that round's entries to the review log (its disposition Annotations and a closing Resolution), then record the round",
			input.Round, roundOpenTail(inputRound.Review))
	}
	pointer := RoundPointer{ID: fmt.Sprintf("round:%s:%s:%d", location.entityID, location.stage, location.cycle), Stage: location.stage, Cycle: location.cycle,
		Briefing: Briefing{ID: inputRound.Manifest.ID, Digest: inputRound.Digest,
			RoomRef: fmt.Sprintf("./review/%s/round-%d", location.stage, location.cycle)}}
	if _, statErr := os.Lstat(location.room); location.pointer.ID == pointer.ID && os.IsNotExist(statErr) {
		return fmt.Errorf("round identity already has a pointer without its immutable room")
	}
	if beforePublish != nil {
		beforePublish(location.room)
	}
	return publishRound(location.room, roundRoomBytes{Exists: true, Briefing: inputRound.Briefing, Log: inputRound.Log}, func(replay bool) error {
		return mutateEntity(entityPath, entityExpectation{Bytes: location.entity}, func(entity []byte) ([]byte, error) {
			next, err := rebuildRoundEntity(entity, pointer)
			if err == nil && replay && !bytes.Equal(next, entity) {
				return nil, fmt.Errorf("immutable round replay does not match the entity pointer")
			}
			return next, err
		}, replace)
	})
}

// roundLogClosed reports whether a review log has been ANSWERED, which is what
// makes its round admissible to the immutable room. The verdict is the log's first
// Resolution (parseReviewLog's Reviewer), entries are ordered, and validateResolution's
// dichotomy splits verdicts into `approve`, which demands nothing further, and
// `revise`/`hold`, which demand a response. So the log is closed when its LAST entry
// is a Resolution that either sits past the verdict — that Resolution IS the
// response — or is an `approve` verdict; open when it ends at an Annotation (a
// dangling finding nobody closed) or exactly at a non-approve verdict (a demanded
// response never logged). No entry count, actor, or `includes` graph enters this.
//
// Record-time only: ValidateRoundFile shares loadValidateRound but not this check,
// so rooms recorded before the precondition stay readable and their truncation is
// graded at read time instead.
func roundLogClosed(review reviewLog) bool {
	if len(review.Entries) == 0 {
		return false
	}
	last := review.Entries[len(review.Entries)-1].Resolution
	return last != nil && (last != review.Reviewer || last.Decision == "approve")
}

// roundOpenTail names which open shape the log ended in, so the refusal tells the
// reader what to append rather than only that something is missing.
func roundOpenTail(review reviewLog) string {
	last := review.Entries[len(review.Entries)-1]
	if last.Resolution == nil {
		return fmt.Sprintf("it ends at Annotation %s with no closing Resolution", last.ID)
	}
	return fmt.Sprintf("it ends at the reviewer's %s Resolution %s", last.Resolution.Decision, last.ID)
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
	summary := RoundSummary{ID: location.pointer.ID, Stage: location.stage, Cycle: location.cycle, Briefing: loaded.Manifest.ID}
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
	if len(node.Content) != 8 || briefingNode == nil || len(briefingNode.Content) != 6 {
		return RoundPointer{}, fmt.Errorf("entity has invalid review-round pointer")
	}
	err = node.Decode(&pointer)
	wantRoom := fmt.Sprintf("./review/%s/round-%d", pointer.Stage, pointer.Cycle)
	if err != nil || pointer.ID == "" || !roundStageRE.MatchString(pointer.Stage) || pointer.Cycle < 1 ||
		pointer.Briefing.ID == "" || !digestRE.MatchString(pointer.Briefing.Digest) || pointer.Briefing.RoomRef != wantRoom {
		return RoundPointer{}, fmt.Errorf("entity has invalid review-round pointer")
	}
	return pointer, nil
}
func rebuildRoundEntity(original []byte, pointer RoundPointer) ([]byte, error) {
	block, err := yaml.Marshal(struct {
		Round RoundPointer `yaml:"review-round"`
	}{pointer})
	if err != nil {
		return nil, err
	}
	out, err := replaceTopLevels(original, topLevelReplacement{key: "review-round", data: block, insert: true})
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
