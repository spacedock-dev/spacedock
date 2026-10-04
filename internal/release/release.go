// ABOUTME: Release-pipeline version steps — the single stamp target list and the
// ABOUTME: top-level plugin.json version stamp (AC-4).
package release

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

// StampTargets is the ONE authoritative list of files a release stamps: the
// per-host plugin descriptors plus the first-officer shared-core prose. The
// `stamp-version` and `manifest-tag-gate` subcommands default to it when given
// no explicit targets, and `stamp-paths` prints it, so the list has exactly one
// authority instead of being restated per ritual step.
func StampTargets() []string {
	return []string{
		".claude-plugin/plugin.json",
		".codex-plugin/plugin.json",
		".pi/plugin.json",
		"skills/first-officer/references/first-officer-shared-core.md",
	}
}

// replaceTopLevelVersion rewrites the value of the TOP-LEVEL `version` member of
// blob to value, preserving all surrounding bytes and formatting. It is
// JSON-structure aware: a `version` nested inside another object — which can
// appear BEFORE the top-level one — is skipped, because only the first object's
// members are inspected. Errors when the top-level version is absent or is not a
// JSON string, so a malformed manifest fails loud instead of writing unchanged
// bytes.
func replaceTopLevelVersion(blob []byte, value string) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(blob))
	tok, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return nil, fmt.Errorf("parse manifest: top level is not a JSON object")
	}
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("parse manifest: %w", err)
		}
		key, _ := keyTok.(string)
		if key != "version" {
			// Skip this member's value (object, array, or scalar) so the scan
			// stays at the top-level object's keys.
			var skip json.RawMessage
			if err := dec.Decode(&skip); err != nil {
				return nil, fmt.Errorf("parse manifest: %w", err)
			}
			continue
		}
		// The decoder sits just after the `"version"` key string; walk past
		// the colon and whitespace to the opening quote of the value, so only
		// the value's bytes are replaced.
		i := dec.InputOffset()
		for i < int64(len(blob)) && isJSONSpace(blob[i]) {
			i++
		}
		if i >= int64(len(blob)) || blob[i] != ':' {
			return nil, fmt.Errorf("parse manifest: malformed version member")
		}
		i++
		for i < int64(len(blob)) && isJSONSpace(blob[i]) {
			i++
		}
		if i >= int64(len(blob)) || blob[i] != '"' {
			return nil, fmt.Errorf("parse manifest: top-level version is not a string")
		}
		valTok, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("parse manifest: %w", err)
		}
		if _, ok := valTok.(string); !ok {
			return nil, fmt.Errorf("parse manifest: top-level version is not a string")
		}
		valueEnd := dec.InputOffset() // just past the value's closing quote
		out := make([]byte, 0, len(blob)+len(value))
		out = append(out, blob[:i+1]...) // through the value's opening quote
		out = append(out, value...)
		out = append(out, blob[valueEnd-1:]...) // from the value's closing quote
		return out, nil
	}
	return nil, fmt.Errorf("parse manifest: no top-level version member")
}

// isJSONSpace reports whether b is insignificant whitespace BETWEEN JSON
// tokens (the bytes json.Decoder itself skips).
func isJSONSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// stableVersionRe matches a bare stable semver (no `v` prefix, no `-pre`
// suffix). DevPreVersion uses it to reject anything that is not X.Y.Z.
var stableVersionRe = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

// DevPreVersion computes the post-release dev pre-version for the `next` edge
// line from a just-released stable version: X.Y.Z -> X.(Y+1).0-pre1. Input must
// be a bare stable semver (no hyphen) — release.yml only calls this on the
// `!contains(github.ref, '-')` branch, which already guarantees that shape.
func DevPreVersion(stableVersion string) (string, error) {
	m := stableVersionRe.FindStringSubmatch(stableVersion)
	if m == nil {
		return "", fmt.Errorf("stable version %q is not X.Y.Z", stableVersion)
	}
	minor, _ := strconv.Atoi(m[2])
	return fmt.Sprintf("%s.%d.0-pre1", m[1], minor+1), nil
}

// ManifestVersion reads the top-level `version` field of a plugin manifest
// (plugin.json / .codex-plugin/plugin.json). It returns an empty string with no
// error when the manifest parses but carries no top-level `version` (a stamp that
// never ran), so the manifest/tag gate can block on a missing version rather than
// erroring. A manifest that does not parse is an error.
func ManifestVersion(manifest []byte) (string, error) {
	var top struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(manifest, &top); err != nil {
		return "", fmt.Errorf("parse manifest: %w", err)
	}
	return top.Version, nil
}

// StampVersion rewrites the top-level `version` field of a plugin manifest
// (plugin.json / .codex-plugin/plugin.json / .pi/plugin.json) to version,
// preserving the rest of the file's formatting. When the manifest has no
// top-level `version` key (e.g. a marketplace.json, whose version lives on the
// nested plugin entry), the input is returned unchanged — the stamp is a named
// plugin descriptor operation and must not move the marketplace entry's calendar
// key. The release CLI turns that no-op into a loud failure for a named manifest.
func StampVersion(manifest []byte, version string) ([]byte, error) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(manifest, &top); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	if _, ok := top["version"]; !ok {
		// No top-level version (marketplace.json shape): nothing to stamp.
		return manifest, nil
	}
	return replaceTopLevelVersion(manifest, version)
}
