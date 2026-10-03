package release

import (
	"fmt"
	"strings"
	"testing"
)

// Independent npm view <spec> name version dist.integrity --json snapshot,
// verified 2026-10-03. Do not derive this oracle from the workflow under test.
var piLivePublishedPins = []struct{ variable, spec, version, integrity, tarball string }{
	{"PI_CODING_AGENT", "@earendil-works/pi-coding-agent", "1.0.0", "sha512-/FtbxoSQU/mEv1QnichJjRjqteqaIaMWxmhB4G367+MwZfX7/DI5B9YAg5lqbN7nztFskBEtUSZ+FlmMBECtMw==", "pi_coding_agent_tgz"},
	{"PI_SUBAGENTS", "pi-subagents", "0.75.0", "sha512-RO4DiTJM6pnK8y9PnD7Y6TLeiX2c8Kh6QhmkecFo+ou5qtnz5qiM+vUKru48UmJG1nB9eJ6rALd1M04uBOR4XQ==", "pi_subagents_tgz"},
	{"PI_INTERCOM", "pi-intercom", "0.16.0", "sha512-ClGQuovPsz7r1iQwMRjEN+8wxywfrDrMILAkCSf/z19Nezzyfz77U8362Wb/VFNoZOwmF4PBsuGWCtB9AsEJMQ==", "pi_intercom_tgz"},
}

// This is a structural wiring guard, not a JS interpreter. Keep the complete
// executable heredoc bound to its invocation so comments or another job's copy
// cannot stand in for either checkpoint. Real-package execution is separate proof.
const piLiveManifestCheck = `const fs = require('fs');
const path = require('path');
const root = process.argv[2];
const manifest = JSON.parse(fs.readFileSync(path.join(root, 'package.json'), 'utf8'));
const extensions = manifest.pi?.extensions;
if (!Array.isArray(extensions) || extensions.length === 0) throw new Error('missing pi.extensions');
const bridge = manifest.exports?.['./intercom-bridge']?.default;
if (typeof bridge !== 'string' || !bridge) throw new Error('missing intercom-bridge runtime export');
for (const entry of [...extensions, bridge]) {
  if (typeof entry !== 'string' || !entry) throw new Error('invalid runtime entry');
  const target = path.resolve(root, entry);
  if (!fs.statSync(target).isFile()) throw new Error('not a runtime file: ' + target);
  console.log('verified runtime file ' + target);
}
NODE`

var piLiveCheckpoints = []struct{ name, root string }{
	{"Install Pi CLI and substrates", "$pi_npm_root/node_modules/pi-subagents"},
	{"Verify Pi current-checkout setup", "$PI_SUBAGENTS_PACKAGE_ROOT"},
}

func assertPiLivePinsAndSubstrateAssertions(workflow string) error {
	steps := map[string]workflowStep{}
	for _, job := range parseWorkflowJobs(workflow) {
		if job.name == "pi-live" {
			for _, step := range job.steps {
				if step.name == "" {
					continue
				}
				if _, duplicate := steps[step.name]; duplicate {
					return fmt.Errorf("duplicate pi-live step %q", step.name)
				}
				steps[step.name] = step
			}
		}
	}
	for _, checkpoint := range piLiveCheckpoints {
		step := steps[checkpoint.name]
		check := "node - \"" + checkpoint.root + "\" <<'NODE'\n" + piLiveManifestCheck + "\n"
		if step.ifCond != "" || !strings.Contains("\n"+step.run, "\n"+check) {
			return fmt.Errorf("%s lacks unconditional installed-manifest check", checkpoint.name)
		}
		for _, stale := range []string{"src/extension/index.ts", "src/intercom/intercom-bridge.ts"} {
			if strings.Contains(step.run, stale) {
				return fmt.Errorf("%s still asserts %s", checkpoint.name, stale)
			}
		}
	}
	commands := executableShellCommands(steps[piLiveCheckpoints[0].name].run)
	require := func(want string) error {
		for _, command := range commands {
			if command == want {
				return nil
			}
		}
		return fmt.Errorf("pi-live install missing executable command: %s", want)
	}
	for _, pin := range piLivePublishedPins {
		for _, command := range []string{
			fmt.Sprintf(`%s_VERSION="%s"`, pin.variable, pin.version),
			fmt.Sprintf(`%s_INTEGRITY="%s"`, pin.variable, pin.integrity),
			fmt.Sprintf(`%s="$(verified_pack "%s@$%s_VERSION" "$%s_INTEGRITY")"`, pin.tarball, pin.spec, pin.variable, pin.variable),
		} {
			if err := require(command); err != nil {
				return err
			}
		}
	}
	for _, command := range []string{
		`npm install -g "$pi_coding_agent_tgz" --ignore-scripts --no-audit --no-fund --omit=dev`,
		`npm install --prefix "$pi_npm_root" "$pi_subagents_tgz" "$pi_intercom_tgz" --ignore-scripts --no-audit --no-fund --omit=dev`,
	} {
		if err := require(command); err != nil {
			return err
		}
	}
	// The install step must register both substrate packages in the real agent
	// directory's settings.json so Pi's own discovery resolves both roots with
	// neither package-root override exported. Root symlinks alone do not register
	// resources, so the registration is load-bearing, not cosmetic.
	install := steps[piLiveCheckpoints[0].name]
	for _, command := range []string{
		`settings_path="$HOME/.pi/agent/settings.json"`,
		`node - "$settings_path" <<'NODE'`,
	} {
		if err := require(command); err != nil {
			return err
		}
	}
	for _, source := range []string{`'npm:pi-subagents'`, `'npm:pi-intercom'`} {
		if !strings.Contains(install.run, source) {
			return fmt.Errorf("pi-live install step does not register %s in the agent settings", source)
		}
	}
	for _, command := range commands {
		if !strings.Contains(command, "PI_SUBAGENTS_PACKAGE_ROOT") && !strings.Contains(command, "PI_INTERCOM_PACKAGE_ROOT") {
			continue
		}
		if strings.Contains(command, "GITHUB_ENV") || strings.Contains(command, "export") {
			return fmt.Errorf("pi-live install step must not export package-root overrides: %s", command)
		}
	}
	return nil
}

func TestPiLivePinsAndSubstrateAssertions(t *testing.T) {
	if err := assertPiLivePinsAndSubstrateAssertions(readWorkflowFiles(t)["runtime-live-e2e.yml"]); err != nil {
		t.Fatal(err)
	}
}

func TestPiLivePinsAndSubstrateAssertionsRejectMutations(t *testing.T) {
	workflow := readWorkflowFiles(t)["runtime-live-e2e.yml"]
	if err := assertPiLivePinsAndSubstrateAssertions(workflow); err != nil {
		t.Fatal(err)
	}
	mutations := []struct{ name, old, replacement string }{
		{"wrong-job", "  pi-live:", "  not-pi-live:"},
		{"agent-install-bypass", `npm install -g "$pi_coding_agent_tgz"`, `npm install -g @earendil-works/pi-coding-agent`},
		{"substrate-install-bypass", `"$pi_subagents_tgz" \`, `pi-subagents \`},
		{"intercom-install-bypass", `"$pi_intercom_tgz" \`, `pi-intercom \`},
		{"settings-registration-removed", `node - "$settings_path" <<'NODE'`, `# node - "$settings_path" <<'NODE'`},
		{"settings-path-not-agent-dir", `settings_path="$HOME/.pi/agent/settings.json"`, `settings_path="$RUNNER_TEMP/settings.json"`},
		{"settings-registration-subagents-missing", `for (const source of ['npm:pi-subagents', 'npm:pi-intercom']) {`, `for (const source of ['npm:pi-intercom']) {`},
		{"settings-registration-intercom-missing", `for (const source of ['npm:pi-subagents', 'npm:pi-intercom']) {`, `for (const source of ['npm:pi-subagents']) {`},
		{"re-export-subagents-root", `test -f "$pi_npm_root/node_modules/pi-intercom/skills/pi-intercom/SKILL.md"`, `test -f "$pi_npm_root/node_modules/pi-intercom/skills/pi-intercom/SKILL.md"` + "\n          echo \"PI_SUBAGENTS_PACKAGE_ROOT=$pi_npm_root/node_modules/pi-subagents\" >> \"$GITHUB_ENV\""},
		{"re-export-intercom-root", `test -f "$pi_npm_root/node_modules/pi-intercom/skills/pi-intercom/SKILL.md"`, `test -f "$pi_npm_root/node_modules/pi-intercom/skills/pi-intercom/SKILL.md"` + "\n          export PI_INTERCOM_PACKAGE_ROOT=\"$pi_npm_root/node_modules/pi-intercom\""},
	}
	oldHashes := []string{
		"sha512-FGRN+OHbWaefBPGaTggAdLjrIHW+s2PzLyglz/5dfLzb9of7uuXMXYC0fJIeZTw+shS32o2cuQ9jF7YSDuL/oQ==",
		"sha512-43FGBi82sbxEGhEfsaj0P7I/rRJfsw1UPuFSVSkjZw0Fy92lrBIB9d3gDTpwmK5hLj3ZzZUnvtlqiKQatySCTw==",
		"sha512-+QjKJRAEhrgQZj4+M9OW/8unRLvCzeCp0K66lZmbS5/me0fsClXRtANBgM6mY+EoX+Fcd+qBE8dTThp8+ND//g==",
	}
	for i, pin := range piLivePublishedPins {
		version := fmt.Sprintf(`%s_VERSION="%s"`, pin.variable, pin.version)
		hash := fmt.Sprintf(`%s_INTEGRITY="%s"`, pin.variable, pin.integrity)
		call := fmt.Sprintf(`%s="$(verified_pack "%s@$%s_VERSION" "$%s_INTEGRITY")"`, pin.tarball, pin.spec, pin.variable, pin.variable)
		for _, m := range []struct{ name, old, replacement string }{
			{pin.variable + "-old-version", version, fmt.Sprintf(`%s_VERSION="%s"`, pin.variable, []string{"0.85.1", "0.67.0", "0.13.0"}[i])},
			{pin.variable + "-corrupt-integrity", hash, pin.variable + `_INTEGRITY="sha512-corrupted"`},
			{pin.variable + "-old-integrity", hash, fmt.Sprintf(`%s_INTEGRITY="%s"`, pin.variable, oldHashes[i])},
			{pin.variable + "-comment-pin", version, "# " + version},
			{pin.variable + "-remove-pack", call, "# " + call},
			{pin.variable + "-bypass-pack", call, strings.Replace(call, "verified_pack", "npm pack", 1)},
		} {
			mutations = append(mutations, m)
		}
	}
	for _, checkpoint := range piLiveCheckpoints {
		invocation := "node - \"" + checkpoint.root + "\" <<'NODE'"
		block := invocation + "\n          " + strings.ReplaceAll(piLiveManifestCheck, "\n", "\n          ")
		for _, missing := range []string{"extensions.length === 0", "typeof bridge !== 'string' || !bridge", "fs.statSync(target).isFile()"} {
			mutations = append(mutations, struct{ name, old, replacement string }{checkpoint.name + missing, block, strings.Replace(block, missing, "false", 1)})
		}
		for _, replacement := range []string{"# " + invocation, `test -f "` + checkpoint.root + `/src/extension/index.ts"`, `test -f "` + checkpoint.root + `/src/intercom/intercom-bridge.ts"`} {
			mutations = append(mutations, struct{ name, old, replacement string }{checkpoint.name + replacement, invocation, replacement})
		}
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			changed := strings.Replace(workflow, mutation.old, mutation.replacement, 1)
			if changed == workflow {
				t.Fatal("mutation did not change workflow")
			}
			if err := assertPiLivePinsAndSubstrateAssertions(changed); err == nil {
				t.Fatal("mutation escaped guard")
			}
		})
	}
}
