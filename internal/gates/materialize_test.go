package gates

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// controlBriefing is the exact s4-prepared indented Briefing retained at state
// blob de790a44 (1,675 bytes, 31 lines). Its two identities are deliberately
// unequal, so an implementation that collapses the domains cannot pass.
const (
	controlRawSha256 = "sha256:c3b6d4d5ac8c766dcc56e08b57a41e207147d1319c61f066160e4e7d4bacfb1b"
	controlJCSDigest = "sha256:0782c65c06c7ee9378226b3a7ef88d92939a54c05d916fe3690cc7d99804278f"
)

// TestControlBriefingSeparatesCanonicalAndRawDigestDomains is the focused
// positive control for AC-3. It fails if the canonicalizer is replaced by a raw
// file hash, if the raw pin is replaced by the canonical digest, or if the two
// domains are ever made equal for this real s4 artifact.
func TestControlBriefingSeparatesCanonicalAndRawDigestDomains(t *testing.T) {
	data, err := os.ReadFile("testdata/materialize-s4-room/gate-briefing.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 1675 {
		t.Fatalf("control Briefing is %d bytes, want the exact 1675-byte s4 artifact", len(data))
	}
	raw := RawDigest(data)
	jcs, err := CanonicalDigest(data)
	if err != nil {
		t.Fatal(err)
	}
	if raw != controlRawSha256 {
		t.Fatalf("raw pin %s, want %s", raw, controlRawSha256)
	}
	if jcs != controlJCSDigest {
		t.Fatalf("canonical digest %s, want %s", jcs, controlJCSDigest)
	}
	if raw == jcs {
		t.Fatal("control Briefing must keep the two digest domains observably unequal")
	}
}

// TestControlBriefingReindentKeepsCanonicalDigestAndMovesRawPin proves the
// domains are independent rather than two names for one value: reindenting the
// same semantic document preserves the canonical digest and changes the raw pin.
func TestControlBriefingReindentKeepsCanonicalDigestAndMovesRawPin(t *testing.T) {
	data, err := os.ReadFile("testdata/materialize-s4-room/gate-briefing.json")
	if err != nil {
		t.Fatal(err)
	}
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	reindented, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	jcs, err := CanonicalDigest(reindented)
	if err != nil {
		t.Fatal(err)
	}
	if jcs != controlJCSDigest {
		t.Fatalf("reindented canonical digest %s, want unchanged %s", jcs, controlJCSDigest)
	}
	if RawDigest(reindented) == controlRawSha256 {
		t.Fatal("reindented raw pin must differ from the original file bytes")
	}
}

// TestMaterializePublishesClosedManifestWithBothBriefingIdentities is the
// end-to-end Spacedock positive path over a real prepared room. It fails if the
// manifest omits or swaps either identity, if the ambiguous `digest` spelling
// returns, if payload bytes drift from their Git objects, or if the manifest is
// not the last published file.
func TestMaterializePublishesClosedManifestWithBothBriefingIdentities(t *testing.T) {
	for _, form := range []string{"folder", "flat"} {
		t.Run(form, func(t *testing.T) {
			fixture := materializeFixture(t, form)

			result, err := Materialize(MaterializeInput{Room: fixture.room, WorkflowDir: fixture.workflow})
			if err != nil {
				t.Fatal(err)
			}
			if result.Sources != 2 {
				t.Fatalf("sources=%d, want the room's two git-root sources", result.Sources)
			}
			if result.Actor != "person:captain" || result.Approver != "person:captain" {
				t.Fatalf("launch tuple actor/approver %q/%q, want the request-frozen captain authority", result.Actor, result.Approver)
			}
			if result.Briefing != realpath(t, fixture.briefingPath) {
				t.Fatalf("launch tuple briefing %s, want the located canonical Briefing %s", result.Briefing, fixture.briefingPath)
			}
			wantManifest := filepath.Join(fixture.room, "provider", "resolved-sources", "resolved-sources.json")
			if result.Manifest != wantManifest {
				t.Fatalf("manifest=%s, want the derived %s", result.Manifest, wantManifest)
			}

			body, err := os.ReadFile(result.Manifest)
			if err != nil {
				t.Fatal(err)
			}
			var manifest struct {
				Type     string `json:"type"`
				Version  string `json:"version"`
				Briefing struct {
					ID        string `json:"id"`
					JCSDigest string `json:"jcsDigest"`
					RawSha256 string `json:"rawSha256"`
				} `json:"briefing"`
				Items []resolvedItem `json:"items"`
			}
			if err := json.Unmarshal(body, &manifest); err != nil {
				t.Fatal(err)
			}
			if manifest.Type != "spacedock-resolved-sources" || manifest.Version != "1" {
				t.Fatalf("manifest is not a closed v1 resolved-source document: %s", body)
			}
			if strings.Contains(string(body), `"digest"`) {
				t.Fatalf("manifest must not carry the ambiguous digest spelling: %s", body)
			}
			briefingBytes, err := os.ReadFile(fixture.briefingPath)
			if err != nil {
				t.Fatal(err)
			}
			wantJCS, err := CanonicalDigest(briefingBytes)
			if err != nil {
				t.Fatal(err)
			}
			if manifest.Briefing.JCSDigest != wantJCS {
				t.Fatalf("jcsDigest=%s, want the canonical %s", manifest.Briefing.JCSDigest, wantJCS)
			}
			if manifest.Briefing.RawSha256 != RawDigest(briefingBytes) {
				t.Fatalf("rawSha256=%s, want the raw file pin %s", manifest.Briefing.RawSha256, RawDigest(briefingBytes))
			}
			if manifest.Briefing.JCSDigest == manifest.Briefing.RawSha256 {
				t.Fatal("the two published Briefing identities must stay in separate domains")
			}

			if len(manifest.Items) != 2 {
				t.Fatalf("manifest items=%d, want 2", len(manifest.Items))
			}
			if manifest.Items[0].Type != "Artifact" || manifest.Items[1].Type != "Reference" {
				t.Fatalf("items are not in canonical presentation order: %#v", manifest.Items)
			}
			for i, item := range manifest.Items {
				if !strings.HasPrefix(item.URI, "git-root://") || item.ID == "" || item.MediaType == "" || item.Rev == "" {
					t.Fatalf("item %d is missing a canonical tuple field: %#v", i, item)
				}
				payload, err := os.ReadFile(filepath.Join(filepath.Dir(result.Manifest), item.Path))
				if err != nil {
					t.Fatal(err)
				}
				// Each payload must be the exact bytes its canonical rev pins,
				// so a payload swap or truncation fails here.
				if RawDigest(payload) != item.Rev {
					t.Fatalf("payload %s does not match its canonical rev %s", item.Path, item.Rev)
				}
				info, err := os.Lstat(filepath.Join(filepath.Dir(result.Manifest), item.Path))
				if err != nil {
					t.Fatal(err)
				}
				if info.Mode().Perm() != 0o600 || !info.Mode().IsRegular() {
					t.Fatalf("payload %s is not a mode-0600 regular file (%v)", item.Path, info.Mode())
				}
			}
		})
	}
}

// TestMaterializeKeepsPreparedRoomAtTwoAuthoritativeFiles is AC-2's durable
// authority proof: materialization adds provider-owned evidence but never a
// third prepared file and never a selected-source copy in the room itself.
func TestMaterializeKeepsPreparedRoomAtTwoAuthoritativeFiles(t *testing.T) {
	fixture := materializeFixture(t, "folder")
	before := roomAuthorityEntries(t, fixture.room)
	if len(before) != 2 {
		t.Fatalf("prepared room entries %v, want exactly request.json and the canonical Briefing", before)
	}

	if _, err := Materialize(MaterializeInput{Room: fixture.room, WorkflowDir: fixture.workflow}); err != nil {
		t.Fatal(err)
	}

	after := roomAuthorityEntries(t, fixture.room)
	if len(after) != 2 || after[0] != before[0] || after[1] != before[1] {
		t.Fatalf("room authority changed from %v to %v", before, after)
	}
	briefingAfter, err := os.ReadFile(fixture.briefingPath)
	if err != nil {
		t.Fatal(err)
	}
	if RawDigest(briefingAfter) != fixture.briefingRaw {
		t.Fatal("materialization rewrote the canonical Briefing bytes")
	}
	if _, err := os.Lstat(filepath.Join(fixture.room, "association.json")); err == nil {
		t.Fatal("materialization must not write association.json")
	}
}

// TestMaterializeRefusesBeforeAllocatingOrWritingProviderState covers AC-5's
// refusal surface. Every case must fail without publishing a resolved-source
// child, and the provider root itself must be left to its Subspace owner.
func TestMaterializeRefusesBeforeAllocatingOrWritingProviderState(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(t *testing.T, f materializeRoom) (room string)
		wantErr string
	}{
		{
			name: "provider root absent",
			mutate: func(t *testing.T, f materializeRoom) string {
				if err := os.Remove(filepath.Join(f.room, "provider")); err != nil {
					t.Fatal(err)
				}
				return f.room
			},
			wantErr: "must be allocated by the provider integration",
		},
		{
			name: "provider root is world readable",
			mutate: func(t *testing.T, f materializeRoom) string {
				if err := os.Chmod(filepath.Join(f.room, "provider"), 0o755); err != nil {
					t.Fatal(err)
				}
				return f.room
			},
			wantErr: "must be mode 0700",
		},
		{
			name: "provider root is a symlink",
			mutate: func(t *testing.T, f materializeRoom) string {
				provider := filepath.Join(f.room, "provider")
				if err := os.Remove(provider); err != nil {
					t.Fatal(err)
				}
				elsewhere := filepath.Join(t.TempDir(), "provider")
				if err := os.Mkdir(elsewhere, 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(elsewhere, provider); err != nil {
					t.Fatal(err)
				}
				return f.room
			},
			wantErr: "non-symlink directory",
		},
		{
			name: "canonical Briefing semantically mutated",
			mutate: func(t *testing.T, f materializeRoom) string {
				body, err := os.ReadFile(f.briefingPath)
				if err != nil {
					t.Fatal(err)
				}
				// Change one semantic string. The canonical digest moves, so the
				// frozen request binding no longer holds.
				mutated := strings.Replace(string(body), `"question": "`, `"question": "tampered `, 1)
				if mutated == string(body) {
					t.Fatal("fixture question member not found")
				}
				if err := os.WriteFile(f.briefingPath, []byte(mutated), 0o644); err != nil {
					t.Fatal(err)
				}
				return f.room
			},
			wantErr: "do not match the frozen",
		},
		{
			name: "git object pruned before materialization",
			mutate: func(t *testing.T, f materializeRoom) string {
				// Remove the state repository's object store so the addressed
				// blob is no longer present locally. Resolution must fail closed
				// rather than fetch or fall back to worktree bytes.
				if err := os.RemoveAll(filepath.Join(f.stateRoot, ".git", "objects")); err != nil {
					t.Fatal(err)
				}
				return f.room
			},
			wantErr: "resolve",
		},
		{
			name: "room is not the bound attempt room",
			mutate: func(t *testing.T, f materializeRoom) string {
				sibling := filepath.Join(filepath.Dir(f.room), "briefing-2")
				if err := os.MkdirAll(filepath.Join(sibling, "provider"), 0o700); err != nil {
					t.Fatal(err)
				}
				for _, name := range []string{"request.json", filepath.Base(f.briefingPath)} {
					body, err := os.ReadFile(filepath.Join(f.room, name))
					if err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(sibling, name), body, 0o644); err != nil {
						t.Fatal(err)
					}
				}
				return sibling
			},
			wantErr: "not the room bound to attempt",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := materializeFixture(t, "folder")
			room := tc.mutate(t, fixture)

			_, err := Materialize(MaterializeInput{Room: room, WorkflowDir: fixture.workflow})
			if err == nil {
				t.Fatal("materialization succeeded, want a closed refusal")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error %q, want it to name %q", err, tc.wantErr)
			}
			if _, statErr := os.Lstat(filepath.Join(room, "provider", "resolved-sources")); statErr == nil {
				t.Fatal("a refused materialization left a resolved-source child behind")
			}
			entries, readErr := os.ReadDir(filepath.Join(room, "provider"))
			if readErr == nil {
				for _, entry := range entries {
					t.Fatalf("a refused materialization left provider residue %q", entry.Name())
				}
			}
		})
	}
}

// TestMaterializeResolvesArbitraryCanonicalBriefingLocator proves the room-only
// path never reconstructs a `briefing.json` basename: a nested, differently
// named canonical Briefing resolves purely from the frozen request locator.
func TestMaterializeResolvesArbitraryCanonicalBriefingLocator(t *testing.T) {
	fixture := materializeFixture(t, "folder")
	nested := filepath.Join(fixture.room, "canonical", "review-package.json")
	if err := os.MkdirAll(filepath.Dir(nested), 0o700); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(fixture.briefingPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nested, body, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(fixture.briefingPath); err != nil {
		t.Fatal(err)
	}
	rewriteRequestLocator(t, filepath.Join(fixture.room, "request.json"), "canonical/review-package.json")

	result, err := Materialize(MaterializeInput{Room: fixture.room, WorkflowDir: fixture.workflow})
	if err != nil {
		t.Fatal(err)
	}
	if result.Briefing != realpath(t, nested) {
		t.Fatalf("located Briefing %s, want the arbitrary locator target %s", result.Briefing, nested)
	}
	if result.Sources != 2 {
		t.Fatalf("sources=%d, want 2", result.Sources)
	}
}

// TestMaterializeRefusesRepublishingOverAnExistingChild keeps the publication
// boundary single-shot: an already-present resolved-source child is never
// silently replaced, because its bytes may already be in provider hands.
func TestMaterializeRefusesRepublishingOverAnExistingChild(t *testing.T) {
	fixture := materializeFixture(t, "folder")
	if _, err := Materialize(MaterializeInput{Room: fixture.room, WorkflowDir: fixture.workflow}); err != nil {
		t.Fatal(err)
	}
	_, err := Materialize(MaterializeInput{Room: fixture.room, WorkflowDir: fixture.workflow})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("second materialization error=%v, want an existing-child refusal", err)
	}
}

type materializeRoom struct {
	workflow     string
	stateRoot    string
	entity       string
	room         string
	briefingPath string
	briefingRaw  string
}

// materializeFixture builds a real prepared room through the production Prepare
// path, then allocates the provider root the way the Subspace fixed entry does.
// Using Prepare rather than a hand-written room keeps this test bound to the
// landed s4 contract instead of a private copy of it.
func materializeFixture(t *testing.T, form string) materializeRoom {
	t.Helper()
	// The Artifact lives in the main root and the Reference in the state root,
	// so the room spans both logical roots exactly as a real s4 room does and a
	// per-root object prune is observable.
	workflow, state, entity, artifact, reference := prepareFixture(t, form)

	result, err := Prepare(entity, PrepareInput{
		WorkflowDir: workflow,
		Question:    "Present both Git-root sources?",
		Artifact:    artifact,
		Summary:     "exact caller-authored summary",
		References:  []string{reference},
	})
	if err != nil {
		t.Fatal(err)
	}
	// The provider root is allocated by invocation-common before it calls the
	// materializer; Spacedock only ever writes inside it.
	if err := os.Mkdir(filepath.Join(result.Room, "provider"), 0o700); err != nil {
		t.Fatal(err)
	}
	briefingPath := filepath.Join(result.Room, "gate-briefing.json")
	briefingBytes, err := os.ReadFile(briefingPath)
	if err != nil {
		t.Fatal(err)
	}
	return materializeRoom{
		workflow:     workflow,
		stateRoot:    state,
		entity:       entity,
		room:         result.Room,
		briefingPath: briefingPath,
		briefingRaw:  RawDigest(briefingBytes),
	}
}

// roomAuthorityEntries lists the room's own files, ignoring the provider-owned
// child, so the two-file authority count is observable after presentation.
func roomAuthorityEntries(t *testing.T, room string) []string {
	t.Helper()
	entries, err := os.ReadDir(room)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		if entry.Name() == "provider" {
			continue
		}
		names = append(names, entry.Name())
	}
	return names
}

// realpath resolves the temporary-directory symlink macOS puts in front of
// /var so located-path comparisons stay exact rather than approximate.
func realpath(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func rewriteRequestLocator(t *testing.T, requestPath, locator string) {
	t.Helper()
	body, err := os.ReadFile(requestPath)
	if err != nil {
		t.Fatal(err)
	}
	var request map[string]any
	if err := json.Unmarshal(body, &request); err != nil {
		t.Fatal(err)
	}
	briefing, ok := request["briefing"].(map[string]any)
	if !ok {
		t.Fatal("request has no briefing member")
	}
	briefing["locator"] = locator
	updated, err := indentedJSON(request)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(requestPath, updated, 0o644); err != nil {
		t.Fatal(err)
	}
	// The request digest is frozen in the entity, so rewriting the locator also
	// requires re-freezing that binding. Rebind it here rather than weakening
	// the production check.
	rebindRequestDigest(t, requestPath)
}

func rebindRequestDigest(t *testing.T, requestPath string) {
	t.Helper()
	body, err := os.ReadFile(requestPath)
	if err != nil {
		t.Fatal(err)
	}
	digest, err := CanonicalDigest(body)
	if err != nil {
		t.Fatal(err)
	}
	room := filepath.Dir(requestPath)
	entity := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(room))), "index.md")
	if _, statErr := os.Lstat(entity); statErr != nil {
		entity = filepath.Dir(filepath.Dir(filepath.Dir(room))) + ".md"
	}
	entityBody, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	old := doc.Records[0].Attempts[0].Briefing.RequestDigest
	updated := strings.Replace(string(entityBody), old, digest, 1)
	if updated == string(entityBody) {
		t.Fatal("frozen request digest not found in entity")
	}
	if err := os.WriteFile(entity, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
}
