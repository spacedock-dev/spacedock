// ABOUTME: Behavioural tests for the pi-live stamp logic and lane orchestration.
// ABOUTME: Every case names a real failure (bad manifest, bad integrity, bad settings), not a copy drift.
package pilive

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeRunner stubs command output by name+args. Tests install closures that
// assert the exact command they expect, so a dropped or reordered invocation
// surfaces as an error rather than passing silently.
type fakeRunner struct {
	run func(name string, args ...string) ([]byte, error)
}

func (f fakeRunner) Output(name string, args ...string) ([]byte, error) {
	if f.run == nil {
		return nil, fmt.Errorf("unexpected command: %s %s", name, strings.Join(args, " "))
	}
	return f.run(name, args...)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runtimeManifestJSON(name, version string, pi, exports map[string]any) string {
	doc := map[string]any{"name": name, "version": version}
	if pi != nil {
		doc["pi"] = pi
	}
	if exports != nil {
		doc["exports"] = exports
	}
	b, _ := json.Marshal(doc)
	return string(b)
}

func TestVerifyRuntimeManifestResolvesEveryDeclaredEntry(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "index.js"), "// ext\n")
	writeFile(t, filepath.Join(root, "second.js"), "// ext2\n")
	writeFile(t, filepath.Join(root, "src", "api", "intercom-bridge.js"), "// bridge\n")
	writeFile(t, filepath.Join(root, "package.json"), runtimeManifestJSON("pi-subagents", "0.0.0",
		map[string]any{"extensions": []string{"./index.js", "./second.js"}},
		exportsWithBridge("./src/api/intercom-bridge.js")))

	resolved, err := VerifyRuntimeManifest(root)
	if err != nil {
		t.Fatalf("valid manifest rejected: %v", err)
	}
	if len(resolved) != 3 {
		t.Fatalf("resolved %d entries, want 3 (two extensions + bridge): %v", len(resolved), resolved)
	}
}

func exportsWithBridge(bridge string) map[string]any {
	return map[string]any{"./intercom-bridge": map[string]any{"default": bridge, "types": bridge + ".d.ts"}}
}

func TestVerifyRuntimeManifestRejectsBadDeclarations(t *testing.T) {
	bridge := "./src/api/intercom-bridge.js"
	cases := []struct {
		name       string
		manifest   string
		files      []string
		wantSubstr string
	}{
		{"missing pi.extensions", runtimeManifestJSON("pi-subagents", "1", nil, exportsWithBridge(bridge)), []string{bridge}, "pi.extensions"},
		{"empty pi.extensions", runtimeManifestJSON("pi-subagents", "1", map[string]any{"extensions": []string{}}, exportsWithBridge(bridge)), []string{bridge}, "pi.extensions"},
		{"missing bridge export", runtimeManifestJSON("pi-subagents", "1", map[string]any{"extensions": []string{"./index.js"}}, nil), []string{"./index.js"}, "intercom-bridge"},
		{"bridge is a bare string", runtimeManifestJSON("pi-subagents", "1", map[string]any{"extensions": []string{"./index.js"}}, map[string]any{"./intercom-bridge": bridge}), []string{"./index.js", bridge}, "intercom-bridge"},
		{"types-only bridge", runtimeManifestJSON("pi-subagents", "1", map[string]any{"extensions": []string{"./index.js"}}, map[string]any{"./intercom-bridge": map[string]any{"types": "./bridge.d.ts"}}), []string{"./index.js"}, "intercom-bridge"},
		{"extension target missing", runtimeManifestJSON("pi-subagents", "1", map[string]any{"extensions": []string{"./absent.js"}}, exportsWithBridge(bridge)), []string{bridge}, "absent.js"},
		{"bridge target missing", runtimeManifestJSON("pi-subagents", "1", map[string]any{"extensions": []string{"./index.js"}}, exportsWithBridge(bridge)), []string{"./index.js"}, "intercom-bridge"},
		{"second extension missing", runtimeManifestJSON("pi-subagents", "1", map[string]any{"extensions": []string{"./index.js", "./second.js"}}, exportsWithBridge(bridge)), []string{"./index.js", bridge}, "second.js"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for _, f := range tc.files {
				writeFile(t, filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(f, "./"))), "x")
			}
			writeFile(t, filepath.Join(root, "package.json"), tc.manifest)
			if _, err := VerifyRuntimeManifest(root); err == nil || !strings.Contains(err.Error(), tc.wantSubstr) {
				t.Fatalf("error = %v, want containing %q", err, tc.wantSubstr)
			}
		})
	}
}

func TestVerifyRuntimeManifestRejectsDirectoryEntry(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "index.js"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "src", "api", "intercom-bridge.js"), "x")
	writeFile(t, filepath.Join(root, "package.json"), runtimeManifestJSON("pi-subagents", "1",
		map[string]any{"extensions": []string{"./index.js"}}, exportsWithBridge("./src/api/intercom-bridge.js")))
	if _, err := VerifyRuntimeManifest(root); err == nil {
		t.Fatal("a directory extension target must not pass the regular-file check")
	}
}

func TestMergeSettingsCreatesMergesAndIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	pkgs, err := MergeSettings(path, SubstrateNpmSources)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(pkgs, ",") != "npm:pi-subagents,npm:pi-intercom" {
		t.Fatalf("fresh packages = %v", pkgs)
	}

	writeFile(t, path, `{"theme":"dark","packages":["npm:other",123]}`)
	pkgs, err = MergeSettings(path, SubstrateNpmSources)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(pkgs, ",") != "npm:other,npm:pi-subagents,npm:pi-intercom" {
		t.Fatalf("merged packages = %v", pkgs)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatal(err)
	}
	if settings["theme"] != "dark" {
		t.Fatalf("unrelated key dropped: %v", settings)
	}

	first, _ := os.ReadFile(path)
	if _, err := MergeSettings(path, SubstrateNpmSources); err != nil {
		t.Fatal(err)
	}
	second, _ := os.ReadFile(path)
	if string(first) != string(second) {
		t.Fatalf("MergeSettings is not idempotent:\n%s\n---\n%s", first, second)
	}
}

func TestMergeSettingsRejectsMalformedJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	writeFile(t, path, "{not json")
	if _, err := MergeSettings(path, SubstrateNpmSources); err == nil {
		t.Fatal("malformed settings JSON must be an error, not silently overwritten")
	}
}

func TestParsePackMetadataAcceptsSingleEntry(t *testing.T) {
	meta, err := ParsePackMetadata([]byte(`[{"filename":"pi-subagents-0.75.0.tgz","integrity":"sha512-abc"}]`))
	if err != nil {
		t.Fatal(err)
	}
	if meta.Filename != "pi-subagents-0.75.0.tgz" || meta.Integrity != "sha512-abc" {
		t.Fatalf("meta = %+v", meta)
	}
}

func TestParsePackMetadataRejectsBadOutput(t *testing.T) {
	for _, tc := range []struct{ name, data string }{
		{"malformed", "{"},
		{"empty array", "[]"},
		{"two entries", `[{"filename":"a","integrity":"x"},{"filename":"b","integrity":"y"}]`},
		{"missing integrity", `[{"filename":"a"}]`},
		{"missing filename", `[{"integrity":"x"}]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParsePackMetadata([]byte(tc.data)); err == nil {
				t.Fatalf("ParsePackMetadata(%q) accepted bad output", tc.data)
			}
		})
	}
}

func TestPackTarballVerifiesIntegrityAgainstThePin(t *testing.T) {
	packDir := t.TempDir()
	pin := Packages[1] // pi-subagents
	var gotSpec string
	r := fakeRunner{run: func(name string, args ...string) ([]byte, error) {
		if name != "npm" || args[0] != "pack" {
			return nil, fmt.Errorf("unexpected command %s %v", name, args)
		}
		gotSpec = args[1]
		return []byte(fmt.Sprintf(`[{"filename":"pi-subagents-%s.tgz","integrity":%q}]`, pin.Version, "sha512-corrupted")), nil
	}}
	if _, err := PackTarball(r, pin, packDir); err == nil {
		t.Fatal("a corrupted npm-pack integrity must abort the pack: an implementation that drops VerifyIntegrity would return nil here")
	}

	r = fakeRunner{run: func(name string, args ...string) ([]byte, error) {
		return []byte(fmt.Sprintf(`[{"filename":"pi-subagents-%s.tgz","integrity":%q}]`, pin.Version, pin.Integrity)), nil
	}}
	path, err := PackTarball(r, pin, packDir)
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(packDir, "pi-subagents-"+pin.Version+".tgz") {
		t.Fatalf("path = %s", path)
	}
	if gotSpec != pin.Spec+"@"+pin.Version {
		t.Fatalf("packed spec = %q", gotSpec)
	}
}

func TestVersionAtLeastComparesFloors(t *testing.T) {
	for _, tc := range []struct {
		version, floor string
		want           bool
	}{
		{PiCodingAgentFloor, PiCodingAgentFloor, true},
		{"0.82.9", PiCodingAgentFloor, false},
		{"0.83.1", PiCodingAgentFloor, true},
		{"1.0.2", PiCodingAgentFloor, true},
		{"0.0.0", PiSubagentsFloor, false},
		{PiSubagentsVersion, PiSubagentsFloor, true},
	} {
		if got := VersionAtLeast(tc.version, tc.floor); got != tc.want {
			t.Errorf("VersionAtLeast(%q, %q) = %v, want %v", tc.version, tc.floor, got, tc.want)
		}
	}
}

func TestNodeEngineAtLeastEnforcesFloor(t *testing.T) {
	for _, tc := range []struct {
		version string
		want    bool
	}{
		{"22.18.0", false},
		{NodeEngineFloor, true},
		{"24.13.1", true},
	} {
		if got := NodeEngineAtLeast(tc.version); got != tc.want {
			t.Errorf("NodeEngineAtLeast(%q) = %v, want %v", tc.version, got, tc.want)
		}
	}
}

func TestCompatExportPathPrefersNestedThenGlobal(t *testing.T) {
	global := t.TempDir()
	writeFile(t, filepath.Join(global, "@earendil-works", "pi-ai", "package.json"),
		`{"version":"9.9.9","exports":{"./compat":{"import":"./compat.js"}}}`)
	writeFile(t, filepath.Join(global, "@earendil-works", "pi-ai", "compat.js"), "x")
	agentRoot := filepath.Join(t.TempDir(), "pi-coding-agent") // no nested pi-ai
	root, path, err := CompatExportPath(agentRoot, global)
	if err != nil {
		t.Fatal(err)
	}
	if root != filepath.Join(global, "@earendil-works", "pi-ai") || path != filepath.Join(root, "compat.js") {
		t.Fatalf("root=%s path=%s", root, path)
	}
}

func TestCompatExportPathRejectsMissingExportOrFile(t *testing.T) {
	global := t.TempDir()
	agentRoot := filepath.Join(t.TempDir(), "pi-coding-agent")
	nested := filepath.Join(agentRoot, "node_modules", "@earendil-works", "pi-ai")
	writeFile(t, filepath.Join(nested, "package.json"), `{"version":"1.0.0","exports":{"./other":{"import":"./other.js"}}}`)
	writeFile(t, filepath.Join(nested, "other.js"), "x")
	if _, _, err := CompatExportPath(agentRoot, global); err == nil {
		t.Fatal("a pi-ai without a ./compat import export must fail resolution")
	}

	writeFile(t, filepath.Join(nested, "package.json"), `{"version":"1.0.0","exports":{"./compat":{"import":"./compat.js"}}}`)
	if _, _, err := CompatExportPath(agentRoot, global); err == nil {
		t.Fatal("a ./compat export pointing at a missing file must fail resolution")
	}
}

func TestCompatExportLoadsPropagatesNodeFailure(t *testing.T) {
	r := fakeRunner{run: func(name string, args ...string) ([]byte, error) {
		if name != "node" || args[0] != "-e" {
			return nil, fmt.Errorf("unexpected command %s %v", name, args)
		}
		return []byte("boom"), fmt.Errorf("exit status 1")
	}}
	if err := CompatExportLoads(r, "/tmp/compat.js"); err == nil {
		t.Fatal("a non-zero node import must fail the load: an implementation that ignores node's exit would return nil")
	}
}

func TestGuardRejectsBelowFloorSubstrate(t *testing.T) {
	agentRoot, subagentsRoot := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(agentRoot, "package.json"), `{"name":"@earendil-works/pi-coding-agent","version":"`+PiCodingAgentVersion+`"}`)
	writeFile(t, filepath.Join(subagentsRoot, "package.json"), `{"name":"pi-subagents","version":"0.52.0"}`)
	err := Guard(fakeRunner{}, GuardOptions{
		AgentRoot:     agentRoot,
		SubagentsRoot: subagentsRoot,
		GlobalNpmRoot: t.TempDir(),
		NodeVersion:   "24.13.1",
	})
	if err == nil || !strings.Contains(err.Error(), PiSubagentsFloor) {
		t.Fatalf("err = %v, want the pi-subagents floor %s", err, PiSubagentsFloor)
	}
}

func TestGuardRejectsBelowFloorNode(t *testing.T) {
	err := Guard(fakeRunner{}, GuardOptions{NodeVersion: "22.18.0"})
	if err == nil || !strings.Contains(err.Error(), NodeEngineFloor) {
		t.Fatalf("err = %v, want the Node engine floor %s", err, NodeEngineFloor)
	}
}

func TestInstallVerifiesAndRegistersThePinnedFamily(t *testing.T) {
	agentDir := t.TempDir()
	globalRoot := t.TempDir()
	packDir := t.TempDir()
	opts := InstallOptions{
		PackDir:      packDir,
		PiNpmRoot:    filepath.Join(agentDir, "npm"),
		SettingsPath: filepath.Join(agentDir, "settings.json"),
	}
	writeFile(t, filepath.Join(globalRoot, PiCodingAgentSpec, "package.json"),
		`{"name":"`+PiCodingAgentSpec+`","version":"`+PiCodingAgentVersion+`"}`)
	subRoot := filepath.Join(opts.PiNpmRoot, "node_modules", PiSubagentsSpec)
	writeFile(t, filepath.Join(subRoot, "package.json"), runtimeManifestJSON(PiSubagentsSpec, PiSubagentsVersion,
		map[string]any{"extensions": []string{"./index.js"}}, exportsWithBridge("./src/api/intercom-bridge.js")))
	writeFile(t, filepath.Join(subRoot, "index.js"), "x")
	writeFile(t, filepath.Join(subRoot, "src", "api", "intercom-bridge.js"), "x")
	writeFile(t, filepath.Join(subRoot, "skills", "pi-subagents", "SKILL.md"), "x")
	icRoot := filepath.Join(opts.PiNpmRoot, "node_modules", PiIntercomSpec)
	writeFile(t, filepath.Join(icRoot, "package.json"), `{"name":"pi-intercom","version":"`+PiIntercomVersion+`"}`)
	writeFile(t, filepath.Join(icRoot, "skills", "pi-intercom", "SKILL.md"), "x")

	r := fakeRunner{run: func(name string, args ...string) ([]byte, error) {
		switch {
		case name == "npm" && args[0] == "pack":
			at := strings.LastIndex(args[1], "@")
			spec := args[1][:at]
			for _, pkg := range Packages {
				if pkg.Spec == spec {
					return []byte(fmt.Sprintf(`[{"filename":"%s.tgz","integrity":%q}]`, strings.NewReplacer("/", "_", "@", "").Replace(spec), pkg.Integrity)), nil
				}
			}
			return nil, fmt.Errorf("unexpected pack spec %q", args[1])
		case name == "npm" && args[0] == "root":
			return []byte(globalRoot + "\n"), nil
		case name == "npm" && args[0] == "install":
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected command %s %v", name, args)
	}}
	if err := Install(r, opts); err != nil {
		t.Fatalf("Install: %v", err)
	}
	settings, err := os.ReadFile(opts.SettingsPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, source := range SubstrateNpmSources {
		if !strings.Contains(string(settings), source) {
			t.Fatalf("settings %s missing %s", settings, source)
		}
	}
}
