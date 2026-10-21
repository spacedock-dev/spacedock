package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type stagedInstall struct {
	t        *testing.T
	first    string
	second   string
	dataHome string
	env      []string
	project  string
	host     *fakeHost
	lookPath func(string) (string, error)
}

func writeProgram(t *testing.T, path, output string) {
	t.Helper()
	mustMkdir(t, filepath.Dir(path))
	mustWrite(t, path, "#!/bin/sh\necho '"+output+"'\n")
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

// fakeExtensions is a bin extension with a check hook and a tree extension with an
// update note, so the tests drive the mechanism without naming any product.
func fakeExtensions() []extension {
	return []extension{
		{name: "alpha", program: "spacedock-alpha", check: alphaChoiceCheck},
		{name: "beta", program: "spacedock-beta", tree: true,
			afterUpdate: []string{"A running beta keeps its old version until you restart it:", "  spacedock beta --restart"}},
	}
}

func alphaChoiceCheck(in versionInputs) ([]pieceRow, []string) {
	path := filepath.Join(in.dir, "alpha.choice")
	data, err := os.ReadFile(path)
	switch {
	case os.IsNotExist(err):
		return []pieceRow{{"alpha choice", "not chosen", ""}}, nil
	case err != nil || strings.TrimSpace(string(data)) == "bad":
		return []pieceRow{{"alpha choice", "unreadable", path}}, []string{"UNREADABLE: the alpha choice for this directory."}
	}
	return []pieceRow{{"alpha choice", strings.TrimSpace(string(data)), path}}, nil
}

func stageInstall(t *testing.T) *stagedInstall {
	return stageInstallWith(t, fakeExtensions())
}

func stageInstallWith(t *testing.T, table []extension) *stagedInstall {
	t.Helper()
	savedBranch, savedTable := devBranch, extensions
	devBranch, extensions = "next", table
	t.Cleanup(func() { devBranch, extensions = savedBranch, savedTable })

	root := t.TempDir()
	s := &stagedInstall{t: t, first: filepath.Join(root, "first"), second: filepath.Join(root, "second"),
		dataHome: filepath.Join(root, "data"), project: filepath.Join(root, "project")}
	mustMkdir(t, s.project)
	s.env = []string{
		"PATH=" + s.first + string(os.PathListSeparator) + s.second + ":/usr/bin:/bin",
		"HOME=" + root, "XDG_STATE_HOME=" + filepath.Join(root, "state"), "XDG_DATA_HOME=" + s.dataHome,
	}
	writeProgram(t, filepath.Join(s.first, "spacedock"), "spacedock 0.28.0-pre6")
	s.lookPath = hostsOnPath("claude")
	s.host = &fakeHost{inventory: []pluginInventoryEntry{
		{ID: "spacedock@spacedock-edge", Version: "0.28.0-pre6", Installed: true, Enabled: true},
	}}
	for _, ext := range table {
		version := "0.16.0-beta.2"
		if ext.tree {
			version = "1.0.1"
			s.stageTreeProgram(ext, version)
		} else {
			writeProgram(t, filepath.Join(s.first, ext.program), version)
		}
		s.host.inventory = append(s.host.inventory,
			pluginInventoryEntry{ID: ext.name + "@spacedock", Version: version, Installed: true, Enabled: true})
	}
	return s
}

func (s *stagedInstall) stageTreeProgram(ext extension, version string) {
	s.t.Helper()
	program := filepath.Join(s.dataHome, ext.name, version, ext.name, ext.name)
	writeProgram(s.t, program, ext.name+" has no version flag")
	link := filepath.Join(s.first, ext.program)
	os.Remove(link)
	if err := os.Symlink(program, link); err != nil {
		s.t.Fatal(err)
	}
}

func (s *stagedInstall) report() versionReport {
	return collectVersions(versionInputs{env: s.env, dir: s.project, host: s.host, lookPath: s.lookPath,
		selfPath: filepath.Join(s.first, "spacedock"), selfVersion: "0.28.0-pre6"})
}

func TestVersionReportNamesNothingForAConsistentInstall(t *testing.T) {
	s := stageInstall(t)
	report := s.report()
	if len(report.problems) != 0 {
		t.Fatalf("a consistent install reported problems: %q", report.problems)
	}
	want := map[string]string{
		"spacedock": "0.28.0-pre6", "spacedock plugin (Claude Code)": "0.28.0-pre6",
		"alpha program": "0.16.0-beta.2", "alpha plugin (Claude Code)": "0.16.0-beta.2",
		"beta program": "1.0.1", "beta plugin (Claude Code)": "1.0.1", "alpha choice": "not chosen",
	}
	for label, version := range want {
		if got := report.version(label); got != version {
			t.Errorf("%s = %q, want %q", label, got, version)
		}
	}
	var out bytes.Buffer
	printVersionReport(&out, report)
	for _, word := range []string{"MIXED", "MISSING", "UNREADABLE"} {
		if strings.Contains(out.String(), word) {
			t.Errorf("clean report prints %s:\n%s", word, out.String())
		}
	}
}

func TestVersionReportNamesEachMixedInstallProblem(t *testing.T) {
	breakIt := map[string]struct {
		apply func(s *stagedInstall)
		want  string
	}{
		"program newer than its plugin": {
			func(s *stagedInstall) { writeProgram(t, filepath.Join(s.first, "spacedock-alpha"), "0.16.0-beta.3") },
			"MIXED: the alpha program 0.16.0-beta.3 does not match its plugin 0.16.0-beta.2 on Claude Code.",
		},
		"two spacedock copies": {
			func(s *stagedInstall) { writeProgram(t, filepath.Join(s.second, "spacedock"), "spacedock 0.27.2") },
			"MIXED: 2 copies of spacedock on PATH.",
		},
		"extension plugin missing": {
			func(s *stagedInstall) { s.host.inventory = s.host.inventory[:2] },
			"MISSING: the beta plugin is not installed on Claude Code from spacedock.",
		},
		"extension program missing": {
			func(s *stagedInstall) { os.Remove(filepath.Join(s.first, "spacedock-beta")) },
			"MISSING: the beta program is not on PATH.",
		},
		"an extension's own check": {
			func(s *stagedInstall) { mustWrite(t, filepath.Join(s.project, "alpha.choice"), "bad") },
			"UNREADABLE: the alpha choice for this directory.",
		},
	}
	for name, tc := range breakIt {
		t.Run(name, func(t *testing.T) {
			s := stageInstall(t)
			tc.apply(s)
			problems := s.report().problems
			if len(problems) != 1 || !strings.HasPrefix(problems[0], tc.want) {
				t.Fatalf("problems = %q, want exactly one starting %q", problems, tc.want)
			}
		})
	}
}

func (s *stagedInstall) withCodex(plugins []pluginInventoryEntry) *twoHosts {
	s.t.Helper()
	host := &twoHosts{fakeHost: s.host, codexPlugins: plugins}
	s.host = host.fakeHost
	s.lookPath = hostsOnPath("claude", "codex")
	return host
}

func (s *stagedInstall) reportWith(host hostOps) versionReport {
	return collectVersions(versionInputs{env: s.env, dir: s.project, host: host, lookPath: s.lookPath,
		selfPath: filepath.Join(s.first, "spacedock"), selfVersion: "0.28.0-pre6"})
}

func TestVersionReportListsEachHostsPluginsLabelledByHost(t *testing.T) {
	s := stageInstall(t)
	host := s.withCodex(codexPlugins("1.0.1"))
	report := s.reportWith(host)
	if len(report.problems) != 0 {
		t.Fatalf("a consistent install on both hosts reported problems: %q", report.problems)
	}
	for _, label := range []string{"spacedock plugin", "alpha plugin", "beta plugin"} {
		for _, title := range []string{"Claude Code", "Codex"} {
			if report.version(label+" ("+title+")") == "" {
				t.Errorf("no row %q (%s)", label, title)
			}
		}
	}
	var out bytes.Buffer
	printVersionReport(&out, report)
	if !strings.Contains(out.String(), "beta plugin (Codex)") || !strings.Contains(out.String(), "alpha plugin (Claude Code)") {
		t.Errorf("printed list lacks the host labels:\n%s", out.String())
	}
}

func TestVersionReportChecksEachHostOnItsOwn(t *testing.T) {
	cases := map[string]struct {
		apply func(*twoHosts)
		want  string
	}{
		"MIXED on Codex only": {
			func(h *twoHosts) { h.codexPlugins[1].Version = "0.16.0-beta.1" },
			"MIXED: the alpha program 0.16.0-beta.2 does not match its plugin 0.16.0-beta.1 on Codex.",
		},
		"MISSING on Codex only": {
			func(h *twoHosts) { h.codexPlugins = h.codexPlugins[:2] },
			"MISSING: the beta plugin is not installed on Codex from spacedock.",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			s := stageInstall(t)
			host := s.withCodex(codexPlugins("1.0.1"))
			tc.apply(host)
			problems := s.reportWith(host).problems
			if len(problems) != 1 || !strings.HasPrefix(problems[0], tc.want) {
				t.Fatalf("problems = %q, want exactly one starting %q", problems, tc.want)
			}
		})
	}
}

func TestVersionReportNamesAPluginBehindOnceWhenBothHostsAreBehind(t *testing.T) {
	s := stageInstall(t)
	host := s.withCodex(codexPlugins("1.0.1"))
	writeProgram(t, filepath.Join(s.first, "spacedock-alpha"), "0.16.0-beta.3")
	var out bytes.Buffer
	printVersionReport(&out, s.reportWith(host))
	if !strings.Contains(out.String(), "The alpha plugin waits on its listing") || strings.Contains(out.String(), "alpha and alpha") {
		t.Errorf("remedy does not name alpha exactly once:\n%s", out.String())
	}
	if strings.Count(out.String(), "MIXED: the alpha program") != 2 {
		t.Errorf("want one MIXED line per host:\n%s", out.String())
	}
}

func TestVersionReportGivesNoRowsAndNoProblemToAHostThatIsAbsentOrHoldsNoSpacedockPlugin(t *testing.T) {
	cases := map[string]struct {
		onPath  []string
		codex   []pluginInventoryEntry
		claude  func(*stagedInstall)
		absent  string
		present string
	}{
		"codex absent":                       {[]string{"claude"}, codexPlugins("9.9.9"), nil, "(Codex)", "(Claude Code)"},
		"codex without the Spacedock plugin": {[]string{"claude", "codex"}, codexPlugins("9.9.9")[1:], nil, "(Codex)", "(Claude Code)"},
		"claude without the Spacedock plugin": {[]string{"claude", "codex"}, codexPlugins("1.0.1"),
			func(s *stagedInstall) { s.host.inventory = s.host.inventory[1:] }, "(Claude Code)", "(Codex)"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			s := stageInstall(t)
			if tc.claude != nil {
				tc.claude(s)
			}
			host := s.withCodex(tc.codex)
			s.lookPath = hostsOnPath(tc.onPath...)
			report := s.reportWith(host)
			var out bytes.Buffer
			printVersionReport(&out, report)
			if strings.Contains(out.String(), tc.absent) || !strings.Contains(out.String(), tc.present) {
				t.Errorf("want rows for %s and none for %s:\n%s", tc.present, tc.absent, out.String())
			}
			if len(report.problems) != 0 {
				t.Errorf("problems = %q, want none", report.problems)
			}
		})
	}
}

func TestVersionsAreShownForClaudeAndCodexHostsOnly(t *testing.T) {
	cases := map[string]bool{
		"": true, "--host claude": true, "--host codex": true,
		"--host pi": false, "--plugin-manifest x": false, "--plugin-dir x": false,
	}
	for args, want := range cases {
		if got := versionsWanted(strings.Fields(args)); got != want {
			t.Errorf("versionsWanted(%q) = %v, want %v", args, got, want)
		}
	}
}

func TestVersionReportCopiesNamePathVersionAndTheOneThatRuns(t *testing.T) {
	s := stageInstall(t)
	writeProgram(t, filepath.Join(s.second, "spacedock"), "spacedock 0.27.2")
	var out bytes.Buffer
	printVersionReport(&out, s.report())
	for _, want := range []string{
		filepath.Join(s.first, "spacedock") + "    0.28.0-pre6  (runs)",
		filepath.Join(s.second, "spacedock") + "    0.27.2",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("report lacks %q:\n%s", want, out.String())
		}
	}
}

func TestVersionReportRemedyDependsOnWhichSideIsNewer(t *testing.T) {
	cases := map[string]struct {
		alphaVersion string
		want, absent string
	}{
		"program newer: the plugin waits on the marketplace": {"0.16.0-beta.3",
			"The alpha plugin waits on its listing in the spacedock marketplace, which each extension's release keeps current. Run `spacedock update` after the listing shows the new version.",
			"to bring every program and plugin"},
		"program older: update fixes it": {"0.16.0-beta.1",
			"Run `spacedock update` to bring every program and plugin to the same release.",
			"waits on"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			s := stageInstall(t)
			writeProgram(t, filepath.Join(s.first, "spacedock-alpha"), tc.alphaVersion)
			var out bytes.Buffer
			printVersionReport(&out, s.report())
			if !strings.Contains(out.String(), tc.want) {
				t.Errorf("report lacks %q:\n%s", tc.want, out.String())
			}
			if strings.Contains(out.String(), tc.absent) {
				t.Errorf("report contains %q:\n%s", tc.absent, out.String())
			}
		})
	}
}

func TestVersionReportExpectsNoExtensionOnStable(t *testing.T) {
	s := stageInstall(t)
	devBranch = "main"
	os.Remove(filepath.Join(s.first, "spacedock-alpha"))
	os.Remove(filepath.Join(s.first, "spacedock-beta"))
	s.host.inventory = s.host.inventory[:1]

	report := s.report()
	if len(report.problems) != 0 {
		t.Errorf("stable named a problem for extensions it does not install: %q", report.problems)
	}
	if got := report.version("alpha program"); got != "" {
		t.Errorf("stable lists the alpha program (%q)", got)
	}
}

func TestExtensionTableAgreesWithInstallScript(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`(?s)\nEXTRAS='([^']*)'\n`).FindStringSubmatch(string(src))
	if m == nil {
		t.Fatal("install.sh declares no EXTRAS table")
	}
	rows := map[string][]string{}
	if table := strings.TrimSpace(m[1]); table != "" {
		for _, line := range strings.Split(table, "\n") {
			fields := strings.Split(line, "|")
			if len(fields) != 6 {
				t.Fatalf("EXTRAS row %q has %d fields, want 6", line, len(fields))
			}
			rows[fields[0]] = fields
		}
	}
	if len(rows) != len(extensions) {
		t.Errorf("install.sh lists %d extensions and internal/cli lists %d", len(rows), len(extensions))
	}
	for _, ext := range extensions {
		fields, ok := rows[ext.name]
		if !ok {
			t.Errorf("%s is in internal/cli but not in install.sh", ext.name)
			continue
		}
		if (fields[5] == "tree") != ext.tree {
			t.Errorf("%s: install.sh kind %q, internal/cli tree=%v", ext.name, fields[5], ext.tree)
		}
		linked := false
		for _, link := range strings.Fields(fields[4]) {
			linked = linked || link[strings.Index(link, ":")+1:] == ext.program
		}
		if !linked {
			t.Errorf("%s: install.sh links %q, internal/cli looks for %q", ext.name, fields[4], ext.program)
		}
	}
}
