// ABOUTME: standing-subcommand three-channel parity — list-standing/show-standing,
// ABOUTME: native vs vendored Python over fixtures.
package dispatch

import (
	"path/filepath"
	"strings"
	"testing"
)

// standingMod returns a standing-teammate mod body with the given frontmatter and
// optional ## Routing Usage / ## Agent Prompt sections. A mod's spawn name comes
// from its ## Hook: startup section.
func standingMod(name, model, description, routingUsage string) string {
	fm := "---\nstanding: true\nname: " + name + "\n"
	if description != "" {
		fm += "description: " + description + "\n"
	}
	fm += "---\n"
	body := "## Hook: startup\n" +
		"- subagent_type: general-purpose\n" +
		"- name: " + name + "\n" +
		"- model: " + model + "\n"
	if routingUsage != "" {
		body += "## Routing Usage\n" + routingUsage + "\n"
	}
	body += "## Agent Prompt\nYou are " + name + ".\n"
	return fm + body
}

// writeMods materializes a _mods dir under workflowDir with the given mod files
// (filename -> content) and returns the mods dir path.
func writeMods(t *testing.T, workflowDir string, mods map[string]string) string {
	t.Helper()
	modsDir := filepath.Join(workflowDir, "_mods")
	for name, content := range mods {
		writeFile(t, filepath.Join(modsDir, name), content)
	}
	return modsDir
}

// TestListStandingParity drives list-standing over a _mods fixture mixing standing
// and non-standing mods: the native and oracle emit the same sorted absolute paths.
func TestListStandingParity(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"comm-officer.md":    standingMod("comm-officer", "sonnet", "prose polisher", ""),
		"science-officer.md": standingMod("science-officer", "opus", "researcher", ""),
		"not-standing.md":    "---\nstanding: false\nname: nope\n---\nbody\n",
	})

	native := runNative("", "list-standing", "--workflow-dir", wd)
	assertGolden(t, "list-standing", goldenEnvelope{res: normRun(native, wd, home)})
}

// TestListStandingParityNoMods drives the degenerate path: a workflow with no
// _mods dir emits empty stdout, exit 0, on both sides.
func TestListStandingParityNoMods(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	native := runNative("", "list-standing", "--workflow-dir", wd)
	assertGolden(t, "list-standing-no-mods", goldenEnvelope{res: normRun(native, wd, home)})
}

// TestShowStandingParity drives show-standing over a _mods fixture: one mod with a
// ## Routing Usage section (extracted body) and one without (fallback one-liner),
// so both render branches are byte-compared native vs oracle.
func TestShowStandingParity(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"comm-officer.md": standingMod("comm-officer", "sonnet", "prose polisher",
			"Send a draft; reply is the polished text.\n\nKeep it tight."),
		"science-officer.md": standingMod("science-officer", "opus", "researcher", ""),
	})

	native := runNative("", "show-standing", "--workflow-dir", wd)
	assertGolden(t, "show-standing", goldenEnvelope{res: normRun(native, wd, home)})
}

// TestShowStandingParityEmpty drives the degenerate empty case (no standing mods):
// empty stdout, exit 0, on both sides.
func TestShowStandingParityEmpty(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"not-standing.md": "---\nstanding: false\nname: nope\n---\nbody\n",
	})
	native := runNative("", "show-standing", "--workflow-dir", wd)
	assertGolden(t, "show-standing-empty", goldenEnvelope{res: normRun(native, wd, home)})
}

// assertEmDashEscaped asserts stdout carries the literal \u2014 escape and no raw
// UTF-8 em-dash, so the ensure_ascii fix is demonstrably what makes the byte
// parity hold (a guard against the harness comparing raw == raw if reverted).
func assertEmDashEscaped(t *testing.T, stdout string) {
	t.Helper()
	if !strings.Contains(stdout, "\\u2014") {
		t.Errorf("stdout missing the \\u2014 ensure_ascii escape (em-dash emitted raw?):\n%s", stdout)
	}
	if strings.ContainsRune(stdout, '—') {
		t.Errorf("stdout contains a raw em-dash (ensure_ascii escaping not applied):\n%s", stdout)
	}
}

// TestBuildModsParity guards that a workflow declaring a standing teammate does
// NOT get an auto-injected show-standing fetch line from `dispatch build` — that
// auto-injection only ever fired on the retired legacy team_name path (merged and
// bare dispatches always omitted it, per build.go's documented behavior). The
// standing-teammate flow stays reachable directly via show-standing /
// spawn-standing-all, covered by the tests above and in spawn_standing_all_*_test.go.
func TestBuildModsParity(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	writeMods(t, root, map[string]string{
		"comm-officer.md": standingMod("comm-officer", "sonnet", "prose polisher", ""),
	})
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "backlog", ""))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "backlog",
		"checklist":      []string{"- a", "- b"},
		"bare_mode":      false,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	nativeBody := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	env := goldenEnvelope{res: normRun(native, root, home), body: normPaths(nativeBody, root, home)}
	assertGolden(t, "build-mods", env)
}
