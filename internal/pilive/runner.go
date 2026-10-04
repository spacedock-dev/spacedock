// ABOUTME: Orchestrates the pi-live install and compatibility guard through Go,
// ABOUTME: replacing the lane's inline Node heredocs and JS one-liners.
package pilive

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Runner runs an external command and returns its stdout. The lane supplies an
// exec-based runner; tests supply a fake.
type Runner interface {
	Output(name string, args ...string) ([]byte, error)
}

// ExecRunner is the production Runner. It returns stdout only; a failing
// command's error carries the captured stderr, so an npm warning emitted while
// resolving a path (e.g. `npm root -g`) cannot pollute the result.
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

// PackTarball runs `npm pack` for pkg and returns the downloaded tarball path
// after verifying its published integrity against the pin — the lane's
// independent registry oracle. A mismatch aborts before any install.
func PackTarball(r Runner, pkg Package, packDir string) (string, error) {
	out, err := r.Output("npm", "pack", pkg.Spec+"@"+pkg.Version, "--pack-destination", packDir, "--ignore-scripts", "--json")
	if err != nil {
		return "", fmt.Errorf("npm pack %s@%s: %w: %s", pkg.Spec, pkg.Version, err, strings.TrimSpace(string(out)))
	}
	meta, err := ParsePackMetadata(out)
	if err != nil {
		return "", err
	}
	if err := VerifyIntegrity(meta.Integrity, pkg.Integrity); err != nil {
		return "", fmt.Errorf("%s@%s: %w", pkg.Spec, pkg.Version, err)
	}
	return filepath.Join(packDir, meta.Filename), nil
}

// CompatExportLoads evaluates the resolved compat module in Node, so a module
// that resolves but fails to load is caught. This is the one place the lane
// still needs a JS runtime: verifying an ESM load is inherently a Node action.
func CompatExportLoads(r Runner, compatPath string) error {
	fileURL := (&url.URL{Scheme: "file", Path: compatPath}).String()
	const script = `import(process.argv[1]).catch((error) => { console.error(error.message); process.exit(1); })`
	out, err := r.Output("node", "-e", script, fileURL)
	if err != nil {
		return fmt.Errorf("@earendil-works/pi-ai/compat exists but failed to load: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// InstallOptions names the directories the lane install writes into.
type InstallOptions struct {
	PackDir      string // where npm pack downloads tarballs
	PiNpmRoot    string // <agentDir>/npm, the substrate install prefix
	SettingsPath string // <agentDir>/settings.json
}

// Install performs the whole "Install Pi CLI and substrates" step: verify each
// pin's integrity, install pi-coding-agent globally and both substrates into the
// agent npm prefix, verify every installed name+version, assert the substrate
// runtime manifests, check the packaged skill files, and register both npm
// substrates in the agent settings so Pi's own discovery loads them. It never
// exports either package-root override.
func Install(r Runner, opts InstallOptions) error {
	if err := os.MkdirAll(opts.PackDir, 0o755); err != nil {
		return err
	}
	tarballs := map[string]string{}
	for _, pkg := range Packages {
		tarball, err := PackTarball(r, pkg, opts.PackDir)
		if err != nil {
			return err
		}
		tarballs[pkg.Spec] = tarball
	}

	if out, err := r.Output("npm", "install", "-g", tarballs[PiCodingAgentSpec], "--ignore-scripts", "--no-audit", "--no-fund", "--omit=dev"); err != nil {
		return fmt.Errorf("install %s: %w: %s", PiCodingAgentSpec, err, strings.TrimSpace(string(out)))
	}
	globalRoot, err := npmGlobalRoot(r)
	if err != nil {
		return err
	}
	if err := VerifyInstalled(filepath.Join(globalRoot, PiCodingAgentSpec), PiCodingAgentSpec, PiCodingAgentVersion); err != nil {
		return err
	}

	if err := os.MkdirAll(opts.PiNpmRoot, 0o755); err != nil {
		return err
	}
	if out, err := r.Output("npm", "install", "--prefix", opts.PiNpmRoot,
		tarballs[PiSubagentsSpec], tarballs[PiIntercomSpec],
		"--ignore-scripts", "--no-audit", "--no-fund", "--omit=dev"); err != nil {
		return fmt.Errorf("install substrates: %w: %s", err, strings.TrimSpace(string(out)))
	}
	subagentsRoot := filepath.Join(opts.PiNpmRoot, "node_modules", PiSubagentsSpec)
	intercomRoot := filepath.Join(opts.PiNpmRoot, "node_modules", PiIntercomSpec)
	if err := VerifyInstalled(subagentsRoot, PiSubagentsSpec, PiSubagentsVersion); err != nil {
		return err
	}
	if err := VerifyInstalled(intercomRoot, PiIntercomSpec, PiIntercomVersion); err != nil {
		return err
	}
	if _, err := VerifyRuntimeManifest(subagentsRoot); err != nil {
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
	if _, err := MergeSettings(opts.SettingsPath, SubstrateNpmSources); err != nil {
		return err
	}
	fmt.Printf("installed %s@%s, %s@%s, %s@%s; registered %s\n",
		PiCodingAgentSpec, PiCodingAgentVersion, PiSubagentsSpec, PiSubagentsVersion,
		PiIntercomSpec, PiIntercomVersion, strings.Join(SubstrateNpmSources, ", "))
	return nil
}

func npmGlobalRoot(r Runner) (string, error) {
	out, err := r.Output("npm", "root", "-g")
	if err != nil {
		return "", fmt.Errorf("npm root -g: %w: %s", err, strings.TrimSpace(string(out)))
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("npm root -g returned no path")
	}
	return root, nil
}

// GuardOptions names the installed roots and running Node version the
// compatibility guard reads.
type GuardOptions struct {
	AgentRoot     string
	SubagentsRoot string
	GlobalNpmRoot string
	NodeVersion   string
}

// Guard is the lane's compatibility guard: it enforces the Node engine floor,
// the pi-coding-agent ready-gate floor, and the pi-subagents pairing floor, then
// requires @earendil-works/pi-ai/compat to resolve AND load.
func Guard(r Runner, opts GuardOptions) error {
	if !NodeEngineAtLeast(opts.NodeVersion) {
		return fmt.Errorf("Node %s does not satisfy the pi-coding-agent >=%s engine requirement", opts.NodeVersion, NodeEngineFloor)
	}
	agentName, agentVersion, err := ReadInstalledIdentity(opts.AgentRoot)
	if err != nil {
		return err
	}
	if !VersionAtLeast(agentVersion, PiCodingAgentFloor) {
		return fmt.Errorf("%s %s is below the spacedock ready-gate floor %s", agentName, agentVersion, PiCodingAgentFloor)
	}
	subName, subVersion, err := ReadInstalledIdentity(opts.SubagentsRoot)
	if err != nil {
		return err
	}
	if !VersionAtLeast(subVersion, PiSubagentsFloor) {
		return fmt.Errorf("%s %s is below the last verified pairing floor %s", subName, subVersion, PiSubagentsFloor)
	}
	_, compatPath, err := CompatExportPath(opts.AgentRoot, opts.GlobalNpmRoot)
	if err != nil {
		return err
	}
	if err := CompatExportLoads(r, compatPath); err != nil {
		return err
	}
	fmt.Printf("verified %s %s pairs with @earendil-works/pi-ai/compat at %s\n", subName, subVersion, compatPath)
	return nil
}
