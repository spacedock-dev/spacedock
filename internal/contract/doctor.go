// ABOUTME: spacedock doctor — read a plugin manifest's display version, compare
// ABOUTME: it against the binary's, and report one of five verdicts with an exit code.
package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// errNoManifest is returned by readManifest when the manifest file does not
// exist — distinguishing the no-plugin-found report from a malformed field.
var errNoManifest = errors.New("manifest not found")

// readManifest reads a plugin manifest JSON and returns its display version. A
// missing file yields errNoManifest. requires-contract, if present, is not read —
// it is a frozen cross-era sentinel for integer-era binaries only (D4). The
// display version is the user-facing semver shown in the verdict message.
func readManifest(manifestPath string) (version string, err error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errNoManifest
		}
		return "", err
	}
	var m struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return "", fmt.Errorf("parse manifest %s: %w", manifestPath, err)
	}
	return m.Version, nil
}

// ManifestVersion reads a plugin manifest JSON and returns its display version,
// the same source the doctor verdict reads. A missing manifest file yields
// errNoManifest; an unparseable manifest yields the parse error. Callers that only
// need the version to show it (e.g. --version's per-runtime block) use this rather
// than the full verdict.
func ManifestVersion(manifestPath string) (string, error) {
	return readManifest(manifestPath)
}

// ManifestVerdict resolves the compatibility verdict for the manifest at
// manifestPath against binaryVersion, for the named host. A missing manifest file
// yields NoPluginFound; an unparseable manifest JSON yields a MalformedVersion
// Result naming the parse error. The plugin's display version (from the manifest)
// and the binary's display version (threaded in by the cli caller) are woven into
// the user-facing message. src is the detected install source that tailors the
// too-old-binary remedy (the zero value reproduces the generic block). The front
// door inspects the verdict directly (a non-empty path to a missing file is
// NoPluginFound, NOT compatible); RunDoctor maps the same verdict to an exit code
// and stream.
func ManifestVerdict(manifestPath, host, binaryVersion string, src InstallSource) Result {
	pluginVersion, err := readManifest(manifestPath)
	if errors.Is(err, errNoManifest) {
		return Result{Verdict: NoPluginFound, Message: noPluginMessage(host)}
	}
	if err != nil {
		return Result{Verdict: MalformedVersion, Message: fmt.Sprintf("error: %s", err)}
	}
	return compareNamed(host, manifestPath, pluginVersion, binaryVersion, src)
}

// RunDoctor reports the compatibility verdict for the manifest at manifestPath
// against binaryVersion, for the named host. binaryVersion is the binary's
// display version threaded in for the user-facing message; src is the detected
// install source that tailors the too-old-binary remedy. A compatible verdict
// and a no-plugin-found report exit 0 (the report is non-fatal-by-default); every
// mismatch (too-old-binary, too-old-plugin, malformed-version) exits 1 with the
// pinned remedy on stderr.
func RunDoctor(manifestPath, host, binaryVersion string, src InstallSource, stdout, stderr io.Writer) int {
	res := ManifestVerdict(manifestPath, host, binaryVersion, src)
	switch res.Verdict {
	case Compatible:
		fmt.Fprintln(stdout, res.Message)
		return 0
	case NoPluginFound:
		fmt.Fprintln(stdout, res.Message)
		return 0
	default:
		fmt.Fprintln(stderr, res.Message)
		return 1
	}
}
