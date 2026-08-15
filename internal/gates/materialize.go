// ABOUTME: Room-only materialization of Git-root gate sources for provider presentation.
// ABOUTME: Publishes a closed manifest last and never copies payloads into durable room state.
package gates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/gitsource"
)

const (
	resolvedSourcesType    = "spacedock-resolved-sources"
	resolvedSourcesVersion = "1"
	providerChildName      = "provider"
	resolvedSourcesChild   = "resolved-sources"
	resolvedManifestName   = "resolved-sources.json"
	payloadChildName       = "payload"
	candidateChildPrefix   = ".resolved-sources-candidate-"
)

var briefingRoomRE = regexp.MustCompile(`^briefing-[1-9][0-9]*$`)

// MaterializeInput carries only the one room operand. WorkflowDir is derived by
// fixed code from that room before this call; it is never a caller coordinate on
// the public provider grammar.
type MaterializeInput struct {
	Room        string
	WorkflowDir string
}

// MaterializeResult is the derived private launch tuple. The fixed provider
// branch consumes it directly; it is not an agent copy/paste interface.
type MaterializeResult struct {
	Manifest     string
	Sources      int
	Briefing     string
	ProviderRoot string
	Actor        string
	Approver     string
}

// resolvedSources is the closed v1 manifest. Unknown, missing, duplicate, or
// alternate-spelling members fail closed on the reading side; the ambiguous
// `digest` spelling is deliberately absent rather than retained as an alias.
type resolvedSources struct {
	Type     string             `json:"type"`
	Version  string             `json:"version"`
	Briefing resolvedBriefingID `json:"briefing"`
	Items    []resolvedItem     `json:"items"`
}

// resolvedBriefingID carries the two independent Briefing identities. jcsDigest
// is the existing canonical recorder authority over the RFC 8785 serialization;
// rawSha256 is a provider-handoff pin over the exact located file bytes. They
// are computed in separate domains and are never compared with each other.
type resolvedBriefingID struct {
	ID        string `json:"id"`
	JCSDigest string `json:"jcsDigest"`
	RawSha256 string `json:"rawSha256"`
}

type resolvedItem struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	URI       string `json:"uri"`
	MediaType string `json:"mediaType"`
	Rev       string `json:"rev"`
	Path      string `json:"path"`
}

// Materialize resolves every Git-root source addressed by a prepared room's
// canonical Briefing into provider-owned payloads, then publishes the closed
// manifest last. It requires the provider root to be already allocated by its
// Subspace owner and neither creates, removes, selects, nor relocates it.
func Materialize(input MaterializeInput) (MaterializeResult, error) {
	room, err := canonicalRoom(input.Room)
	if err != nil {
		return MaterializeResult{}, err
	}
	entityPath, err := entityForRoom(room)
	if err != nil {
		return MaterializeResult{}, err
	}
	providerRoot, err := allocatedProviderRoot(room)
	if err != nil {
		return MaterializeResult{}, err
	}

	requestBytes, err := os.ReadFile(filepath.Join(room, "request.json"))
	if err != nil {
		return MaterializeResult{}, fmt.Errorf("read gate room request: %w", err)
	}
	request, err := decodeGateRoomRequest(requestBytes)
	if err != nil {
		return MaterializeResult{}, err
	}
	if request.Type != "spacedock-gate-presentation-request" || request.Version != resolvedSourcesVersion {
		return MaterializeResult{}, fmt.Errorf("gate room request is not a complete presentation request v1")
	}
	if request.Actor == "" || request.Approver == "" {
		return MaterializeResult{}, fmt.Errorf("gate room request does not freeze actor and approver")
	}
	requestDigest, err := CanonicalDigest(requestBytes)
	if err != nil {
		return MaterializeResult{}, fmt.Errorf("canonicalize gate room request: %w", err)
	}

	// The located Briefing is read at the request's exact clean relative
	// locator. No basename is appended or required, so an arbitrary canonical
	// filename resolves without reconstruction.
	briefingPath, err := resolveBriefingLocator(room, request.Briefing.Locator)
	if err != nil {
		return MaterializeResult{}, err
	}
	briefingInfo, err := os.Lstat(briefingPath)
	if err != nil || !briefingInfo.Mode().IsRegular() {
		return MaterializeResult{}, fmt.Errorf("canonical Briefing must be a regular non-symlink file")
	}
	briefingBytes, err := os.ReadFile(briefingPath)
	if err != nil {
		return MaterializeResult{}, fmt.Errorf("read canonical Briefing: %w", err)
	}
	manifest, err := parseBriefingManifest(briefingBytes)
	if err != nil {
		return MaterializeResult{}, err
	}
	jcsDigest, err := CanonicalDigest(briefingBytes)
	if err != nil {
		return MaterializeResult{}, fmt.Errorf("canonicalize canonical Briefing: %w", err)
	}
	if manifest.ID != request.Briefing.ID || jcsDigest != request.Briefing.Digest {
		return MaterializeResult{}, fmt.Errorf("canonical Briefing identity and canonical digest do not match the frozen request binding")
	}
	if err := requireFrozenAttemptBinding(entityPath, room, request, requestDigest, jcsDigest); err != nil {
		return MaterializeResult{}, err
	}
	// The raw pin is computed only after the canonical authority above holds. It
	// is a separate domain over the exact located file bytes and never stands in
	// for the canonical digest.
	rawSha256 := RawDigest(briefingBytes)

	items, err := gitRootPresentationItems(manifest)
	if err != nil {
		return MaterializeResult{}, err
	}
	roots := gitsource.Roots{Main: input.WorkflowDir, State: filepath.Dir(entityPath)}
	payloads := make([][]byte, 0, len(items))
	for i := range items {
		body, err := gitsource.Resolve(roots, items[i].URI, items[i].Rev)
		if err != nil {
			return MaterializeResult{}, fmt.Errorf("resolve %s %s: %w", strings.ToLower(items[i].Type), items[i].ID, err)
		}
		payloads = append(payloads, body)
	}

	published, err := publishResolvedSources(providerRoot, resolvedSources{
		Type:    resolvedSourcesType,
		Version: resolvedSourcesVersion,
		Briefing: resolvedBriefingID{
			ID:        manifest.ID,
			JCSDigest: jcsDigest,
			RawSha256: rawSha256,
		},
		Items: items,
	}, payloads)
	if err != nil {
		return MaterializeResult{}, err
	}
	return MaterializeResult{
		Manifest:     published,
		Sources:      len(items),
		Briefing:     briefingPath,
		ProviderRoot: providerRoot,
		Actor:        request.Actor,
		Approver:     request.Approver,
	}, nil
}

func canonicalRoom(room string) (string, error) {
	if strings.TrimSpace(room) == "" {
		return "", fmt.Errorf("gate materialize requires --room")
	}
	abs, err := filepath.Abs(room)
	if err != nil {
		return "", fmt.Errorf("resolve gate room: %w", err)
	}
	abs = filepath.Clean(abs)
	info, err := os.Lstat(abs)
	if err != nil {
		return "", fmt.Errorf("resolve gate room: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return "", fmt.Errorf("gate room must be an existing non-symlink directory")
	}
	return abs, nil
}

// entityForRoom inverts the s4 room placement. Both entity forms share the room
// shape <entity-root>/<slug>/review/<stage>/briefing-N, so the folder and flat
// candidates are collision-free; zero or two candidates are refused rather than
// resolved by searching arbitrary parents.
func entityForRoom(room string) (string, error) {
	if !briefingRoomRE.MatchString(filepath.Base(room)) {
		return "", fmt.Errorf("gate room is not a canonical prepared room")
	}
	reviewDir := filepath.Dir(filepath.Dir(room))
	if filepath.Base(reviewDir) != "review" {
		return "", fmt.Errorf("gate room is not a canonical prepared room")
	}
	home := filepath.Dir(reviewDir)
	slug := filepath.Base(home)
	if slug == "." || slug == string(filepath.Separator) || slug == "" {
		return "", fmt.Errorf("gate room is not a canonical prepared room")
	}
	var found []string
	for _, candidate := range []string{filepath.Join(home, "index.md"), home + ".md"} {
		info, err := os.Lstat(candidate)
		if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			continue
		}
		found = append(found, candidate)
	}
	if len(found) != 1 {
		return "", fmt.Errorf("gate room resolves %d entity candidates; exactly one folder or flat entity is required", len(found))
	}
	return found[0], nil
}

// allocatedProviderRoot requires the Subspace-owned provider root to be present
// and private. Spacedock writes only inside it and never allocates it, so a
// failure here happens before any Git read or payload write.
func allocatedProviderRoot(room string) (string, error) {
	providerRoot := filepath.Join(room, providerChildName)
	info, err := os.Lstat(providerRoot)
	if err != nil {
		return "", fmt.Errorf("provider root %s must be allocated by the provider integration before materialization", providerRoot)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return "", fmt.Errorf("provider root must be a non-symlink directory")
	}
	if perm := info.Mode().Perm(); perm != 0o700 {
		return "", fmt.Errorf("provider root must be mode 0700, found %#o", perm)
	}
	return providerRoot, nil
}

// requireFrozenAttemptBinding requires the room's request to match the entity's
// current recorded attempt exactly. A wrong room, stale attempt, or drifted
// binding fails before Git reads and before any provider-visible write.
func requireFrozenAttemptBinding(entityPath, room string, request *gateRoomRequest, requestDigest, jcsDigest string) error {
	doc, _, err := Read(entityPath)
	if err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("entity has no recorded gate authority for this room")
	}
	for ri := range doc.Records {
		record := &doc.Records[ri]
		if record.ID != request.Gate {
			continue
		}
		for ai := range record.Attempts {
			attempt := &record.Attempts[ai]
			if attempt.ID != request.Attempt {
				continue
			}
			binding := attempt.Briefing
			if binding.RequestDigest == "" {
				return fmt.Errorf("attempt %s is not a request-backed prepared room", attempt.ID)
			}
			if binding.RequestDigest != requestDigest {
				return fmt.Errorf("gate room request does not match the frozen request digest")
			}
			if binding.ID != request.Briefing.ID || binding.Digest != jcsDigest {
				return fmt.Errorf("canonical Briefing identity and canonical digest do not match the frozen attempt binding")
			}
			bound := filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(binding.RoomRef))
			boundReal, boundErr := filepath.EvalSymlinks(bound)
			roomReal, roomErr := filepath.EvalSymlinks(room)
			if boundErr != nil || roomErr != nil || boundReal != roomReal {
				return fmt.Errorf("--room is not the room bound to attempt %s", attempt.ID)
			}
			if attempt.Resolution != nil {
				return fmt.Errorf("attempt %s is already closed", attempt.ID)
			}
			return nil
		}
	}
	return fmt.Errorf("entity has no open attempt %s for gate %s", request.Attempt, request.Gate)
}

// gitRootPresentationItems walks Artifacts in order and then recursively reached
// References in canonical context order, keeping only Git-root sources. Payload
// paths are assigned in that same canonical presentation order and are
// deliberately unrelated to the repository path.
func gitRootPresentationItems(manifest *briefingManifest) ([]resolvedItem, error) {
	items := make([]resolvedItem, 0, len(manifest.Artifacts))
	seen := map[string]bool{}
	add := func(kind, id, uri, mediaType, rev string) error {
		if seen[id] {
			return fmt.Errorf("canonical Briefing has a duplicate source id %s", id)
		}
		seen[id] = true
		if !strings.HasPrefix(uri, "git-root://") {
			return nil
		}
		if mediaType == "" {
			return fmt.Errorf("canonical Briefing source %s has no mediaType", id)
		}
		items = append(items, resolvedItem{Type: kind, ID: id, URI: uri, MediaType: mediaType, Rev: rev})
		return nil
	}
	for _, artifact := range manifest.Artifacts {
		if err := add("Artifact", artifact.ID, artifact.URI, artifact.MediaType, artifact.Rev); err != nil {
			return nil, err
		}
	}
	references, err := gitRootReferences(manifest.Context)
	if err != nil {
		return nil, err
	}
	for _, reference := range references {
		if err := add("Reference", reference.ID, reference.URI, reference.MediaType, reference.Rev); err != nil {
			return nil, err
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("canonical Briefing addresses no git-root source")
	}
	for i := range items {
		items[i].Path = fmt.Sprintf("%s/%04d", payloadChildName, i+1)
	}
	return items, nil
}

type referenceNode struct {
	Type      string            `json:"type"`
	ID        string            `json:"id"`
	URI       string            `json:"uri"`
	Rev       string            `json:"rev"`
	MediaType string            `json:"mediaType"`
	Summary   json.RawMessage   `json:"summary"`
	Children  []json.RawMessage `json:"children"`
}

func gitRootReferences(nodes []json.RawMessage) ([]referenceNode, error) {
	var references []referenceNode
	for _, raw := range nodes {
		var node referenceNode
		if err := json.Unmarshal(raw, &node); err != nil {
			return nil, fmt.Errorf("parse Briefing context: %w", err)
		}
		if node.Type == "Reference" {
			if node.ID == "" || node.URI == "" || !digestRE.MatchString(node.Rev) {
				return nil, fmt.Errorf("canonical Briefing has an incomplete Reference binding")
			}
			if node.Summary != nil {
				return nil, fmt.Errorf("canonical Briefing References must not carry summaries")
			}
			references = append(references, node)
		}
		nested, err := gitRootReferences(node.Children)
		if err != nil {
			return nil, err
		}
		references = append(references, nested...)
	}
	return references, nil
}

// publishResolvedSources stages every payload and the manifest inside a private
// candidate child, then makes the exact resolved-sources child visible with one
// rename. The manifest is written last, so no reader can observe a manifest that
// outruns its payloads, and a failure before publication removes the candidate
// while leaving the provider root to its owner.
func publishResolvedSources(providerRoot string, manifest resolvedSources, payloads [][]byte) (string, error) {
	target := filepath.Join(providerRoot, resolvedSourcesChild)
	if _, err := os.Lstat(target); err == nil {
		return "", fmt.Errorf("resolved-source child already exists at %s", target)
	}
	candidate, err := os.MkdirTemp(providerRoot, candidateChildPrefix+"*")
	if err != nil {
		return "", fmt.Errorf("stage resolved sources: %w", err)
	}
	published := false
	defer func() {
		if !published {
			os.RemoveAll(candidate)
		}
	}()
	if err := os.Chmod(candidate, 0o700); err != nil {
		return "", fmt.Errorf("stage resolved sources: %w", err)
	}
	if err := os.Mkdir(filepath.Join(candidate, payloadChildName), 0o700); err != nil {
		return "", fmt.Errorf("stage resolved sources: %w", err)
	}
	for i, item := range manifest.Items {
		payloadPath, err := containedPayloadPath(candidate, item.Path)
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(payloadPath, payloads[i], 0o600); err != nil {
			return "", fmt.Errorf("write resolved payload: %w", err)
		}
	}
	body, err := indentedJSON(manifest)
	if err != nil {
		return "", fmt.Errorf("encode resolved-source manifest: %w", err)
	}
	if err := os.WriteFile(filepath.Join(candidate, resolvedManifestName), body, 0o600); err != nil {
		return "", fmt.Errorf("write resolved-source manifest: %w", err)
	}
	if err := os.Rename(candidate, target); err != nil {
		return "", fmt.Errorf("publish resolved sources: %w", err)
	}
	published = true
	return filepath.Join(target, resolvedManifestName), nil
}

func containedPayloadPath(candidate, manifestRelative string) (string, error) {
	if manifestRelative == "" || strings.HasPrefix(manifestRelative, "/") ||
		strings.Contains(manifestRelative, "\\") || filepath.ToSlash(filepath.Clean(manifestRelative)) != manifestRelative {
		return "", fmt.Errorf("resolved payload path %q is not clean and manifest-relative", manifestRelative)
	}
	full := filepath.Join(candidate, filepath.FromSlash(manifestRelative))
	rel, err := filepath.Rel(candidate, full)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("resolved payload path %q escapes the provider candidate", manifestRelative)
	}
	return full, nil
}
