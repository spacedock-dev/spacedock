package status

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"gopkg.in/yaml.v3"
)

// ReadGateHistory writes the complete stored gate object as JSON.
// It uses the status read resolver and changes no files.
func ReadGateHistory(workflowDir, dir, ref string, stdout, stderr io.Writer) int {
	if workflowDir == "" {
		var code int
		workflowDir, code = ResolveWorkflowDir(dir, stderr)
		if code != 0 {
			return code
		}
	}
	roots, err := resolveRoots(workflowDir, dir)
	if err != nil {
		return errExit(stderr, err.Error())
	}
	path, code := resolveReadTarget(roots, ref, stderr)
	if code != 0 {
		return code
	}
	_, node, err := gates.Read(path)
	if err != nil && !errors.Is(err, gates.ErrNoGateRecord) {
		return errExit(stderr, err.Error())
	}
	var value any = map[string]any{"records": []any{}, "version": 1}
	if node != nil {
		retainGateTimestampText(node, map[*yaml.Node]bool{})
		if err := node.Decode(&value); err != nil {
			return errExit(stderr, "decode gate history: "+err.Error())
		}
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return errExit(stderr, "encode gate history: "+err.Error())
	}
	if _, err := stdout.Write(buffer.Bytes()); err != nil {
		return errExit(stderr, err.Error())
	}
	return 0
}

// YAML timestamps retain their stored text instead of a normalized time.Time value.
func retainGateTimestampText(node *yaml.Node, seen map[*yaml.Node]bool) {
	if node == nil || seen[node] {
		return
	}
	seen[node] = true
	if node.Kind == yaml.ScalarNode && node.Tag == "!!timestamp" {
		node.Tag = "!!str"
	}
	for _, child := range node.Content {
		retainGateTimestampText(child, seen)
	}
	retainGateTimestampText(node.Alias, seen)
}
