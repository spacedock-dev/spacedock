// ABOUTME: normalized build-field typed accessors, path joins, and the team/split-root
// ABOUTME: probes build needs, matching the oracle's helper semantics.
package dispatch

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// isJSONNull reports whether a raw JSON value is the literal null.
func isJSONNull(v json.RawMessage) bool {
	return string(bytes.TrimSpace(v)) == "null"
}

// jsonString decodes a raw JSON string value, returning "" for a non-string.
func jsonString(v json.RawMessage) string {
	var s string
	if json.Unmarshal(v, &s) == nil {
		return s
	}
	return ""
}

// optString returns the string value of an optional field, "" when absent, null,
// or non-string — matching inp.get(field) used as a truthy string in the oracle.
func optString(fields map[string]json.RawMessage, key string) string {
	v, ok := fields[key]
	if !ok || isJSONNull(v) {
		return ""
	}
	return jsonString(v)
}

// optBool returns the bool value of an optional field, false when absent, null,
// or non-bool — matching inp.get(field, False).
func optBool(fields map[string]json.RawMessage, key string) bool {
	v, ok := fields[key]
	if !ok || isJSONNull(v) {
		return false
	}
	var b bool
	if json.Unmarshal(v, &b) == nil {
		return b
	}
	return false
}

// jsonStringList decodes a raw JSON array of strings. ok is false when the value
// is not a JSON array (so the caller collapses both empty and non-list to the
// "checklist must not be empty" message, as the oracle does).
func jsonStringList(v json.RawMessage) ([]string, bool) {
	var list []string
	if json.Unmarshal(v, &list) != nil {
		return nil, false
	}
	return list, true
}

// isFile reports whether path is an existing regular file (os.path.isfile).
func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// splitRootStateCheckout returns the absolute state-checkout dir for a split-root
// workflow, or "" when the workflow is single-root. Under split root the README
// declares a state: checkout relative to the README/definition dir, so the
// resolved state checkout is workflowDir/<state> — NOT workflowDir itself (which
// is the definition dir where the state checkout is git-excluded). Returns "" when
// the README is unreadable or carries no non-empty state: field.
func splitRootStateCheckout(workflowDir string) string {
	readmePath := filepath.Join(workflowDir, "README.md")
	fm := status.ParseFrontmatter(readmePath)
	mode, relPath, err := status.ClassifyState(fm["state"])
	if err != nil || mode != status.StateSplitRoot {
		return ""
	}
	return filepath.Join(workflowDir, relPath)
}

// pyRelpath returns path relative to base the way os.path.relpath does for the
// absolute clean paths build passes (entity_path under git_root). filepath.Rel
// computes the same relative path for these inputs.
func pyRelpath(path, base string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return filepath.Base(path)
	}
	return rel
}
