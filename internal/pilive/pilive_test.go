// ABOUTME: Behavioural tests for the pi-live lane: every case names a real
// ABOUTME: failure (bad manifest, bad integrity, bad settings), not a copy drift.
package pilive

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeRunner stubs command output by name+args; an unexpected command is an error
// so a dropped or reordered invocation surfaces rather than passing silently.
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

func manifestJSON(name, version string, pi, exports any) string {
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

func bridgeExport(path string) map[string]any {
	return map[string]any{"./intercom-bridge": map[string]any{"default": path, "types": path + ".d.ts"}}
}

func TestVerifyManifestResolvesEveryDeclaredEntry(t *testing.T) {
	root := t.TempDir()
	for _, f := range []string{"index.js", "second.js", "src/api/intercom-bridge.js"} {
		writeFile(t, filepath.Join(root, filepath.FromSlash(f)), "x")
	}
	writeFile(t, filepath.Join(root, "package.json"),
		manifestJSON("pi-subagents", "0.0.0", map[string]any{"extensions": []string{"./index.js", "./second.js"}}, bridgeExport("./src/api/intercom-bridge.js")))
	resolved, err := verifyManifest(root)
	if err != nil {
		t.Fatalf("valid manifest rejected: %v", err)
	}
	if len(resolved) != 3 {
		t.Fatalf("resolved %d entries, want two extensions + bridge: %v", len(resolved), resolved)
	}
}

func TestVerifyManifestRejectsBadDeclarations(t *testing.T) {
	bridge := "./src/api/intercom-bridge.js"
	cases := []struct {
		name, manifest, wantSubstr string
		files, dirs                []string
	}{
		{"missing pi.extensions", manifestJSON("pi-subagents", "0.0.0", nil, bridgeExport(bridge)), "pi.extensions", []string{bridge}, nil},
		{"empty pi.extensions", manifestJSON("pi-subagents", "0.0.0", map[string]any{"extensions": []string{}}, bridgeExport(bridge)), "pi.extensions", []string{bridge}, nil},
		{"missing bridge export", manifestJSON("pi-subagents", "0.0.0", map[string]any{"extensions": []string{"./index.js"}}, nil), "intercom-bridge", []string{"./index.js"}, nil},
		{"bridge is a bare string", manifestJSON("pi-subagents", "0.0.0", map[string]any{"extensions": []string{"./index.js"}}, map[string]any{"./intercom-bridge": bridge}), "intercom-bridge", []string{"./index.js", bridge}, nil},
		{"types-only bridge", manifestJSON("pi-subagents", "0.0.0", map[string]any{"extensions": []string{"./index.js"}}, map[string]any{"./intercom-bridge": map[string]any{"types": "./bridge.d.ts"}}), "intercom-bridge", []string{"./index.js"}, nil},
		{"extension target missing", manifestJSON("pi-subagents", "0.0.0", map[string]any{"extensions": []string{"./absent.js"}}, bridgeExport(bridge)), "absent.js", []string{bridge}, nil},
		{"extension target is a directory", manifestJSON("pi-subagents", "0.0.0", map[string]any{"extensions": []string{"./index.js"}}, bridgeExport(bridge)), "index.js", []string{bridge}, []string{"./index.js"}},
		{"bridge target missing", manifestJSON("pi-subagents", "0.0.0", map[string]any{"extensions": []string{"./index.js"}}, bridgeExport(bridge)), "intercom-bridge", []string{"./index.js"}, nil},
		{"second extension missing", manifestJSON("pi-subagents", "0.0.0", map[string]any{"extensions": []string{"./index.js", "./second.js"}}, bridgeExport(bridge)), "second.js", []string{"./index.js", bridge}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for _, f := range tc.files {
				writeFile(t, filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(f, "./"))), "x")
			}
			for _, d := range tc.dirs {
				if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(d, "./"))), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			writeFile(t, filepath.Join(root, "package.json"), tc.manifest)
			if _, err := verifyManifest(root); err == nil || !strings.Contains(err.Error(), tc.wantSubstr) {
				t.Fatalf("error = %v, want containing %q", err, tc.wantSubstr)
			}
		})
	}
}

func TestMergeSettingsCreatesMergesAndIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	pkgs, err := mergeSettings(path, substrateNpmSources)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(pkgs, ",") != "npm:pi-subagents,npm:pi-intercom" {
		t.Fatalf("fresh packages = %v", pkgs)
	}
	writeFile(t, path, `{"theme":"dark","packages":["npm:other",123]}`)
	pkgs, err = mergeSettings(path, substrateNpmSources)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(pkgs, ",") != "npm:other,npm:pi-subagents,npm:pi-intercom" {
		t.Fatalf("merged packages = %v", pkgs)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), `"theme": "dark"`) {
		t.Fatalf("unrelated key dropped: %s", data)
	}
	if _, err := mergeSettings(path, substrateNpmSources); err != nil {
		t.Fatal(err)
	}
	again, _ := os.ReadFile(path)
	if string(data) != string(again) {
		t.Fatalf("mergeSettings is not idempotent:\n%s\n---\n%s", data, again)
	}
}

func TestMergeSettingsRejectsMalformedJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	writeFile(t, path, "{not json")
	if _, err := mergeSettings(path, substrateNpmSources); err == nil {
		t.Fatal("malformed settings JSON must be an error, not silently overwritten")
	}
}

func TestPackVerifiesIntegrityAgainstThePin(t *testing.T) {
	pin := packages[1]
	for _, tc := range []struct {
		name, output string
		wantErr      bool
	}{
		{"corrupt integrity", fmt.Sprintf(`[{"filename":"a.tgz","integrity":"sha512-corrupted"}]`), true},
		{"malformed json", "{", true},
		{"empty array", "[]", true},
		{"missing integrity", `[{"filename":"a.tgz"}]`, true},
		{"matching integrity", fmt.Sprintf(`[{"filename":"pi-subagents-%s.tgz","integrity":%q}]`, pin.version, pin.integrity), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := fakeRunner{run: func(name string, args ...string) ([]byte, error) {
				if name == "npm" && args[0] == "pack" {
					return []byte(tc.output), nil
				}
				return nil, fmt.Errorf("unexpected command %s %v", name, args)
			}}
			dir := t.TempDir()
			path, err := pack(r, pin, dir)
			if tc.wantErr != (err != nil) {
				t.Fatalf("err = %v, wantErr = %v", err, tc.wantErr)
			}
			if want := filepath.Join(dir, "pi-subagents-"+pin.version+".tgz"); !tc.wantErr && path != want {
				t.Fatalf("path = %s, want %s", path, want)
			}
		})
	}
}

func TestVersionAtLeastComparesFloors(t *testing.T) {
	for _, tc := range []struct {
		version, floor string
		want           bool
	}{
		{PiCodingAgentFloor, PiCodingAgentFloor, true}, {"0.82.9", PiCodingAgentFloor, false},
		{"0.83.1", PiCodingAgentFloor, true}, {"1.0.2", PiCodingAgentFloor, true},
		{"0.0.0", PiSubagentsFloor, false}, {PiSubagentsVersion, PiSubagentsFloor, true},
		{"22.18.0", NodeEngineFloor, false}, {NodeEngineFloor, NodeEngineFloor, true}, {"24.13.1", NodeEngineFloor, true},
	} {
		if got := versionAtLeast(tc.version, tc.floor); got != tc.want {
			t.Errorf("versionAtLeast(%q, %q) = %v, want %v", tc.version, tc.floor, got, tc.want)
		}
	}
}

func TestGuardRejectsBelowFloorNodeAndSubstrate(t *testing.T) {
	agentDir := t.TempDir()
	nodeBelow := fakeRunner{run: func(name string, args ...string) ([]byte, error) { return []byte("22.18.0\n"), nil }}
	if err := guard(nodeBelow, agentDir); err == nil || !strings.Contains(err.Error(), NodeEngineFloor) {
		t.Fatalf("err = %v, want the Node engine floor %s", err, NodeEngineFloor)
	}

	subRoot := filepath.Join(agentDir, "npm", "node_modules", PiSubagentsSpec)
	writeFile(t, filepath.Join(subRoot, "package.json"), `{"name":"pi-subagents","version":"0.52.0"}`)
	globalRoot := filepath.Join(agentDir, "global")
	writeFile(t, filepath.Join(globalRoot, PiCodingAgentSpec, "package.json"),
		`{"name":"@earendil-works/pi-coding-agent","version":"`+PiCodingAgentVersion+`"}`)
	r := fakeRunner{run: func(name string, args ...string) ([]byte, error) {
		if name == "node" {
			return []byte("v24.13.1\n"), nil
		}
		return []byte(globalRoot + "\n"), nil
	}}
	if err := guard(r, agentDir); err == nil || !strings.Contains(err.Error(), PiSubagentsFloor) {
		t.Fatalf("err = %v, want the pi-subagents floor %s", err, PiSubagentsFloor)
	}
}

func writeCompatPiAi(t *testing.T, root, export string) string {
	t.Helper()
	writeFile(t, filepath.Join(root, "package.json"), `{"version":"1.0.0","exports":{"./compat":{"import":"`+export+`"}}}`)
	writeFile(t, filepath.Join(root, export), "x")
	return root
}

func TestCompatExportPathPrefersNestedThenGlobalAndRejectsMissing(t *testing.T) {
	global := t.TempDir()
	globalPi := writeCompatPiAi(t, filepath.Join(global, "@earendil-works", "pi-ai"), "compat.js")
	agentRoot := filepath.Join(t.TempDir(), "pi-coding-agent")
	nested := writeCompatPiAi(t, filepath.Join(agentRoot, "node_modules", "@earendil-works", "pi-ai"), "compat.js")
	if root, _, err := compatExportPath(agentRoot, global); err != nil || root != nested {
		t.Fatalf("nested pi-ai must win: root=%s err=%v", root, err)
	}
	if root, _, err := compatExportPath(filepath.Join(t.TempDir(), "pi-coding-agent"), global); err != nil || root != globalPi {
		t.Fatalf("global fallback root = %s err=%v, want %s", root, err, globalPi)
	}
	// A ./compat export pointing at a missing file must fail resolution.
	broken := filepath.Join(t.TempDir(), "pi-coding-agent")
	writeFile(t, filepath.Join(broken, "node_modules", "@earendil-works", "pi-ai", "package.json"),
		`{"version":"1.0.0","exports":{"./compat":{"import":"./compat.js"}}}`)
	if _, _, err := compatExportPath(broken, t.TempDir()); err == nil {
		t.Fatal("a ./compat export pointing at a missing file must fail resolution")
	}
}

func TestCompatExportLoadsPropagatesNodeFailure(t *testing.T) {
	r := fakeRunner{run: func(name string, args ...string) ([]byte, error) {
		if name == "node" && args[0] == "-e" {
			return []byte("boom"), fmt.Errorf("exit status 1")
		}
		return nil, fmt.Errorf("unexpected command %s %v", name, args)
	}}
	if err := compatExportLoads(r, "/tmp/compat.js"); err == nil {
		t.Fatal("a non-zero node import must fail the load")
	}
}

func TestInstallVerifiesAndRegistersThePinnedFamily(t *testing.T) {
	agentDir := t.TempDir()
	globalRoot := t.TempDir()
	writeFile(t, filepath.Join(globalRoot, PiCodingAgentSpec, "package.json"),
		`{"name":"`+PiCodingAgentSpec+`","version":"`+PiCodingAgentVersion+`"}`)
	subRoot := filepath.Join(agentDir, "npm", "node_modules", PiSubagentsSpec)
	writeFile(t, filepath.Join(subRoot, "package.json"), manifestJSON(PiSubagentsSpec, PiSubagentsVersion,
		map[string]any{"extensions": []string{"./index.js"}}, bridgeExport("./src/api/intercom-bridge.js")))
	writeFile(t, filepath.Join(subRoot, "index.js"), "x")
	writeFile(t, filepath.Join(subRoot, "src/api/intercom-bridge.js"), "x")
	writeFile(t, filepath.Join(subRoot, "skills/pi-subagents/SKILL.md"), "x")
	icRoot := filepath.Join(agentDir, "npm", "node_modules", PiIntercomSpec)
	writeFile(t, filepath.Join(icRoot, "package.json"), `{"name":"pi-intercom","version":"`+PiIntercomVersion+`"}`)
	writeFile(t, filepath.Join(icRoot, "skills/pi-intercom/SKILL.md"), "x")

	r := fakeRunner{run: func(name string, args ...string) ([]byte, error) {
		switch {
		case name == "npm" && args[0] == "pack":
			spec := args[1][:strings.LastIndex(args[1], "@")]
			for _, p := range packages {
				if p.spec == spec {
					return []byte(fmt.Sprintf(`[{"filename":"%s.tgz","integrity":%q}]`, strings.NewReplacer("/", "_", "@", "").Replace(spec), p.integrity)), nil
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
	if err := install(r, agentDir); err != nil {
		t.Fatalf("install: %v", err)
	}
	settings, err := os.ReadFile(filepath.Join(agentDir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, source := range substrateNpmSources {
		if !strings.Contains(string(settings), source) {
			t.Fatalf("settings %s missing %s", settings, source)
		}
	}
}
