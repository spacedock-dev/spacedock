// ABOUTME: The pi-live Runtime Live E2E lane: the pinned Pi family, its readiness
// ABOUTME: floors, and the install/guard/manifest logic the release binary invokes.
package pilive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// The family pins resolved from the npm `latest` dist-tag on 2026-10-04, and the
// readiness floors (pi-coding-agent's is also the launcher ready-gate floor).
// This is the ONLY place these numbers may be written; update here when the
// family moves.
const (
	PiCodingAgentSpec      = "@earendil-works/pi-coding-agent"
	PiCodingAgentVersion   = "1.0.2"
	PiCodingAgentIntegrity = "sha512-3ZdIghMSELMGV3sKi5iASOb1Jwb696fLjmNu0aezaqDxTLLWWoRpqBYkGxJ1CgAMCbtfqXWEF0lcRrlVXmiEGQ=="

	PiSubagentsSpec      = "pi-subagents"
	PiSubagentsVersion   = "0.75.0"
	PiSubagentsIntegrity = "sha512-RO4DiTJM6pnK8y9PnD7Y6TLeiX2c8Kh6QhmkecFo+ou5qtnz5qiM+vUKru48UmJG1nB9eJ6rALd1M04uBOR4XQ=="

	PiIntercomSpec      = "pi-intercom"
	PiIntercomVersion   = "0.16.0"
	PiIntercomIntegrity = "sha512-ClGQuovPsz7r1iQwMRjEN+8wxywfrDrMILAkCSf/z19Nezzyfz77U8362Wb/VFNoZOwmF4PBsuGWCtB9AsEJMQ=="

	NodeEngineFloor    = "22.19.0"
	PiCodingAgentFloor = "0.83.0"
	PiSubagentsFloor   = "0.53.0"
)

// substrateNpmSources are the substrate registrations the isolated Pi home's
// settings.json must carry so Pi's own discovery loads them.
var substrateNpmSources = []string{"npm:pi-subagents", "npm:pi-intercom"}

type pinnedPackage struct{ spec, version, integrity string }

var packages = []pinnedPackage{
	{PiCodingAgentSpec, PiCodingAgentVersion, PiCodingAgentIntegrity},
	{PiSubagentsSpec, PiSubagentsVersion, PiSubagentsIntegrity},
	{PiIntercomSpec, PiIntercomVersion, PiIntercomIntegrity},
}

// Runner runs an external command and returns stdout.
type Runner interface {
	Output(name string, args ...string) ([]byte, error)
}

// ExecRunner is the production Runner; a failing command's error carries stderr,
// so npm warnings cannot pollute the captured stdout.
type ExecRunner struct{}

func (ExecRunner) Output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

// Command dispatches the pi-live lane commands the release binary exposes and
// returns the process exit code.
func Command(name string, args []string) int {
	switch name {
	case "pins":
		for _, p := range packages {
			fmt.Printf("%s %s %s\n", p.spec, p.version, p.integrity)
		}
		fmt.Printf("floors node=%s pi-coding-agent=%s pi-subagents=%s\n", NodeEngineFloor, PiCodingAgentFloor, PiSubagentsFloor)
	case "print-install":
		fmt.Printf("npm install -g %s@%s\n", PiCodingAgentSpec, PiCodingAgentVersion)
		fmt.Printf("npm install --prefix \"$HOME/.pi/agent/npm\" %s@%s %s@%s\n", PiSubagentsSpec, PiSubagentsVersion, PiIntercomSpec, PiIntercomVersion)
	case "verify-manifest":
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "spacedock-release verify-manifest: need exactly one <package-root>")
			return 2
		}
		resolved, err := verifyManifest(args[0])
		if err != nil {
			return fail(name, err)
		}
		for _, path := range resolved {
			fmt.Println("verified runtime file " + path)
		}
	case "install", "guard":
		dir, err := agentDir()
		if err == nil {
			if name == "install" {
				err = install(ExecRunner{}, dir)
			} else {
				err = guard(ExecRunner{}, dir)
			}
		}
		if err != nil {
			return fail(name, err)
		}
	default:
		fmt.Fprintf(os.Stderr, "spacedock-release: unknown pi-live command %q\n", name)
		return 2
	}
	return 0
}

func agentDir() (string, error) {
	if d := os.Getenv("PI_CODING_AGENT_DIR"); d != "" {
		return d, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".pi", "agent"), nil
}

func fail(name string, err error) int {
	fmt.Fprintf(os.Stderr, "spacedock-release %s: %v\n", name, err)
	return 1
}

// install performs the "Install Pi CLI and substrates" step: verify each pin's
// integrity, install the CLI globally and both substrates into the agent npm
// prefix, verify every installed name+version, assert the substrate runtime
// manifest and packaged skills, and register both npm substrates in the agent
// settings. It never exports either package-root override.
func install(r Runner, agentDir string) error {
	packDir := filepath.Join(os.TempDir(), "pi-live-npm-packs")
	if err := os.MkdirAll(packDir, 0o755); err != nil {
		return err
	}
	tarballs := map[string]string{}
	for _, p := range packages {
		tarball, err := pack(r, p, packDir)
		if err != nil {
			return err
		}
		tarballs[p.spec] = tarball
	}
	if out, err := r.Output("npm", "install", "-g", tarballs[PiCodingAgentSpec], "--ignore-scripts", "--no-audit", "--no-fund", "--omit=dev"); err != nil {
		return fmt.Errorf("install %s: %w: %s", PiCodingAgentSpec, err, strings.TrimSpace(string(out)))
	}
	globalRoot, err := npmGlobalRoot(r)
	if err != nil {
		return err
	}
	if err := verifyInstalled(filepath.Join(globalRoot, PiCodingAgentSpec), PiCodingAgentSpec, PiCodingAgentVersion); err != nil {
		return err
	}
	piNpmRoot := filepath.Join(agentDir, "npm")
	if err := os.MkdirAll(piNpmRoot, 0o755); err != nil {
		return err
	}
	if out, err := r.Output("npm", "install", "--prefix", piNpmRoot, tarballs[PiSubagentsSpec], tarballs[PiIntercomSpec], "--ignore-scripts", "--no-audit", "--no-fund", "--omit=dev"); err != nil {
		return fmt.Errorf("install substrates: %w: %s", err, strings.TrimSpace(string(out)))
	}
	subagentsRoot := filepath.Join(piNpmRoot, "node_modules", PiSubagentsSpec)
	intercomRoot := filepath.Join(piNpmRoot, "node_modules", PiIntercomSpec)
	if err := verifyInstalled(subagentsRoot, PiSubagentsSpec, PiSubagentsVersion); err != nil {
		return err
	}
	if err := verifyInstalled(intercomRoot, PiIntercomSpec, PiIntercomVersion); err != nil {
		return err
	}
	if _, err := verifyManifest(subagentsRoot); err != nil {
		return err
	}
	for _, skill := range []string{
		filepath.Join(subagentsRoot, "skills", "pi-subagents", "SKILL.md"),
		filepath.Join(intercomRoot, "skills", "pi-intercom", "SKILL.md"),
	} {
		if info, err := os.Stat(skill); err != nil || !info.Mode().IsRegular() {
			return fmt.Errorf("missing packaged skill file %s", skill)
		}
	}
	if _, err := mergeSettings(filepath.Join(agentDir, "settings.json"), substrateNpmSources); err != nil {
		return err
	}
	fmt.Printf("installed %s@%s, %s@%s, %s@%s; registered %s\n",
		PiCodingAgentSpec, PiCodingAgentVersion, PiSubagentsSpec, PiSubagentsVersion,
		PiIntercomSpec, PiIntercomVersion, strings.Join(substrateNpmSources, ", "))
	return nil
}

// pack runs `npm pack` for p and returns the tarball path after re-deriving its
// published integrity and matching the pin — the lane's independent registry
// oracle. A mismatch aborts before any install.
func pack(r Runner, p pinnedPackage, dir string) (string, error) {
	out, err := r.Output("npm", "pack", p.spec+"@"+p.version, "--pack-destination", dir, "--ignore-scripts", "--json")
	if err != nil {
		return "", fmt.Errorf("npm pack %s@%s: %w: %s", p.spec, p.version, err, strings.TrimSpace(string(out)))
	}
	var entries []struct {
		Filename  string `json:"filename"`
		Integrity string `json:"integrity"`
	}
	if err := json.Unmarshal(out, &entries); err != nil {
		return "", fmt.Errorf("parse npm pack --json output: %w", err)
	}
	if len(entries) != 1 || entries[0].Filename == "" || entries[0].Integrity == "" {
		return "", fmt.Errorf("npm pack --json returned no usable entry: %s", strings.TrimSpace(string(out)))
	}
	if entries[0].Integrity != p.integrity {
		return "", fmt.Errorf("%s@%s: integrity mismatch: expected %s got %s", p.spec, p.version, p.integrity, entries[0].Integrity)
	}
	return filepath.Join(dir, entries[0].Filename), nil
}

func npmGlobalRoot(r Runner) (string, error) {
	out, err := r.Output("npm", "root", "-g")
	if err != nil {
		return "", fmt.Errorf("npm root -g: %w: %s", err, strings.TrimSpace(string(out)))
	}
	if root := strings.TrimSpace(string(out)); root != "" {
		return root, nil
	}
	return "", fmt.Errorf("npm root -g returned no path")
}

// packageIdentity reads name and version from an installed package.json.
func packageIdentity(root string) (name, version string, err error) {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return "", "", err
	}
	var m struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return "", "", fmt.Errorf("parse %s/package.json: %w", root, err)
	}
	return m.Name, m.Version, nil
}

func verifyInstalled(root, wantName, wantVersion string) error {
	name, version, err := packageIdentity(root)
	if err != nil {
		return err
	}
	if name != wantName {
		return fmt.Errorf("installed package name %q at %s, want %q", name, root, wantName)
	}
	if version != wantVersion {
		return fmt.Errorf("installed %s version %q, want %q", wantName, version, wantVersion)
	}
	return nil
}

// verifyManifest asserts an installed substrate declares a nonempty pi.extensions
// array and the exports["./intercom-bridge"].default runtime string, with every
// declared extension and the bridge resolving to a regular file. It returns the
// resolved paths.
func verifyManifest(root string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return nil, err
	}
	var m struct {
		Pi struct {
			Extensions []string `json:"extensions"`
		} `json:"pi"`
		Exports map[string]json.RawMessage `json:"exports"`
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s/package.json: %w", root, err)
	}
	if len(m.Pi.Extensions) == 0 {
		return nil, fmt.Errorf("%s/package.json declares no pi.extensions", root)
	}
	var bridge struct {
		Default string `json:"default"`
	}
	if raw, ok := m.Exports["./intercom-bridge"]; !ok || json.Unmarshal(raw, &bridge) != nil || bridge.Default == "" {
		return nil, fmt.Errorf("%s/package.json declares no exports[\"./intercom-bridge\"].default runtime export", root)
	}
	var resolved []string
	for _, entry := range append(append([]string{}, m.Pi.Extensions...), bridge.Default) {
		target := filepath.Join(root, entry)
		info, err := os.Stat(target)
		if err != nil {
			return nil, fmt.Errorf("runtime entry %q does not resolve to a file: %w", entry, err)
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("runtime entry %q resolves to %s, not a regular file", entry, target)
		}
		resolved = append(resolved, target)
	}
	return resolved, nil
}

// mergeSettings create-or-merges sources into the packages array at path,
// preserving unrelated keys and existing entries. A missing file starts empty;
// malformed JSON is an error. It is idempotent.
func mergeSettings(path string, sources []string) ([]string, error) {
	settings := map[string]json.RawMessage{}
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) > 0 {
		if err := json.Unmarshal(data, &settings); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
	}
	var pkgs []string
	if raw, ok := settings["packages"]; ok {
		var entries []json.RawMessage
		if json.Unmarshal(raw, &entries) == nil {
			for _, entry := range entries {
				var s string
				if json.Unmarshal(entry, &s) == nil {
					pkgs = append(pkgs, s)
				}
			}
		}
	}
	for _, source := range sources {
		seen := false
		for _, existing := range pkgs {
			seen = seen || existing == source
		}
		if !seen {
			pkgs = append(pkgs, source)
		}
	}
	encoded, err := json.Marshal(pkgs)
	if err != nil {
		return nil, err
	}
	settings["packages"] = encoded
	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, append(out, '\n'), 0o644); err != nil {
		return nil, err
	}
	return pkgs, nil
}

// guard enforces the Node engine floor, the pi-coding-agent ready-gate floor, and
// the pi-subagents pairing floor, then requires @earendil-works/pi-ai/compat to
// resolve AND load.
func guard(r Runner, agentDir string) error {
	nodeOut, err := r.Output("node", "--version")
	if err != nil {
		return fmt.Errorf("node --version: %w", err)
	}
	nodeVersion := strings.TrimSpace(string(nodeOut))
	if !versionAtLeast(nodeVersion, NodeEngineFloor) {
		return fmt.Errorf("Node %s does not satisfy the pi-coding-agent >=%s engine requirement", nodeVersion, NodeEngineFloor)
	}
	globalRoot, err := npmGlobalRoot(r)
	if err != nil {
		return err
	}
	agentRoot := filepath.Join(globalRoot, PiCodingAgentSpec)
	agentName, agentVersion, err := packageIdentity(agentRoot)
	if err != nil {
		return err
	}
	if !versionAtLeast(agentVersion, PiCodingAgentFloor) {
		return fmt.Errorf("%s %s is below the spacedock ready-gate floor %s", agentName, agentVersion, PiCodingAgentFloor)
	}
	subName, subVersion, err := packageIdentity(filepath.Join(agentDir, "npm", "node_modules", PiSubagentsSpec))
	if err != nil {
		return err
	}
	if !versionAtLeast(subVersion, PiSubagentsFloor) {
		return fmt.Errorf("%s %s is below the last verified pairing floor %s", subName, subVersion, PiSubagentsFloor)
	}
	_, compatPath, err := compatExportPath(agentRoot, globalRoot)
	if err != nil {
		return err
	}
	if err := compatExportLoads(r, compatPath); err != nil {
		return err
	}
	fmt.Printf("verified %s %s pairs with @earendil-works/pi-ai/compat at %s\n", subName, subVersion, compatPath)
	return nil
}

// compatExportPath resolves @earendil-works/pi-ai's "./compat" import export,
// checking the nested copy under pi-coding-agent first (npm usually does not
// hoist it to the global root), then the top-level global root.
func compatExportPath(agentRoot, globalNpmRoot string) (root, compatPath string, err error) {
	candidates := []string{
		filepath.Join(agentRoot, "node_modules", "@earendil-works", "pi-ai"),
		filepath.Join(globalNpmRoot, "@earendil-works", "pi-ai"),
	}
	for _, candidate := range candidates {
		data, readErr := os.ReadFile(filepath.Join(candidate, "package.json"))
		if readErr != nil {
			continue
		}
		var m struct {
			Exports map[string]struct {
				Import string `json:"import"`
			} `json:"exports"`
		}
		if err := json.Unmarshal(data, &m); err != nil {
			return "", "", fmt.Errorf("parse %s/package.json: %w", candidate, err)
		}
		if compat, ok := m.Exports["./compat"]; ok && compat.Import != "" {
			path := filepath.Join(candidate, compat.Import)
			if _, statErr := os.Stat(path); statErr != nil {
				return "", "", fmt.Errorf("@earendil-works/pi-ai/compat resolves to missing file %s", path)
			}
			return candidate, path, nil
		}
	}
	return "", "", fmt.Errorf("cannot locate @earendil-works/pi-ai with a ./compat import export (checked %s)", strings.Join(candidates, ", "))
}

// compatExportLoads evaluates the resolved compat module in Node, catching a
// module that resolves but fails to load. ESM load verification is inherently a
// Node action — the one place the lane still needs a JS runtime.
func compatExportLoads(r Runner, compatPath string) error {
	fileURL := (&url.URL{Scheme: "file", Path: compatPath}).String()
	const script = `import(process.argv[1]).catch((error) => { console.error(error.message); process.exit(1); })`
	out, err := r.Output("node", "-e", script, fileURL)
	if err != nil {
		return fmt.Errorf("@earendil-works/pi-ai/compat exists but failed to load: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// versionAtLeast reports whether version parses at or above floor; a component
// that is not a number counts as 0 (the output is a plain semver triple).
func versionAtLeast(version, floor string) bool {
	v, f := parseVersion(version), parseVersion(floor)
	for i := 0; i < 3; i++ {
		if v[i] != f[i] {
			return v[i] > f[i]
		}
	}
	return true
}

func parseVersion(version string) [3]int {
	core := strings.TrimPrefix(strings.TrimSpace(version), "v")
	if i := strings.IndexByte(core, '-'); i >= 0 {
		core = core[:i]
	}
	parts := strings.Split(core, ".")
	var out [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		n, _ := strconv.Atoi(strings.TrimSpace(parts[i]))
		out[i] = n
	}
	return out
}
