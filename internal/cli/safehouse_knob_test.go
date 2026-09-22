// ABOUTME: AC-6 sandbox-knob parsing — space/equals/repeat forms produce the same
// ABOUTME: safehouse extra argv, the reported space-form bug leaks nothing, bad value errors.
package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/safehouse"
)

// TestSafehouseKnobFormsEquivalent pins AC-6: each value-taking knob accepts the
// space form, the equals form, and repeats — all producing the same safehouse
// `extra` argv. The equivalence is asserted through the full chain
// parseFrontDoorArgs → safehouse.TranslateFlags, which is what the launcher feeds.
func TestSafehouseKnobFormsEquivalent(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantExtra []string
	}{
		{"enable-equals", []string{"--safehouse-enable=docker"}, []string{"--enable=docker"}},
		{"enable-space", []string{"--safehouse-enable", "docker"}, []string{"--enable=docker"}},
		{"enable-comma-split", []string{"--safehouse-enable=ssh,docker"}, []string{"--enable=ssh", "--enable=docker"}},
		{"enable-repeat-accumulates", []string{"--safehouse-enable", "ssh", "--safehouse-enable", "docker"}, []string{"--enable=ssh", "--enable=docker"}},
		{"add-dirs-equals", []string{"--safehouse-add-dirs=/a"}, []string{"--add-dirs=/a"}},
		{"add-dirs-space", []string{"--safehouse-add-dirs", "/a"}, []string{"--add-dirs=/a"}},
		{"add-dirs-repeat", []string{"--safehouse-add-dirs", "/a", "--safehouse-add-dirs", "/b"}, []string{"--add-dirs=/a", "--add-dirs=/b"}},
		{"add-dirs-ro-space", []string{"--safehouse-add-dirs-ro", "/c"}, []string{"--add-dirs-ro=/c"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fd, err := parseFrontDoorArgs(tc.args)
			if err != nil {
				t.Fatalf("parseFrontDoorArgs(%v) err = %v", tc.args, err)
			}
			extra, err := safehouse.TranslateFlags(fd.safehouseFlags)
			if err != nil {
				t.Fatalf("TranslateFlags(%v) err = %v", fd.safehouseFlags, err)
			}
			if !equalArgv(extra, tc.wantExtra) {
				t.Fatalf("extra = %v, want %v", extra, tc.wantExtra)
			}
		})
	}
}

// TestSafehouseAddDirsSpaceFormNoLeak is the regression for the captain-reported
// bug: `--safehouse-add-dirs ~/a --safehouse-add-dirs ~/b` (space form) used to
// fail with `malformed flag` and leak the paths into host passthrough. In the end
// state cobra captures each value into two `--add-dirs=` entries in the safehouse
// extra slot, and the paths appear NOWHERE in the host argv (they are not host
// passthrough). Driven through runClaude to the recorded launch.
func TestSafehouseAddDirsSpaceFormNoLeak(t *testing.T) {
	dir := safehouseFixtureDir(t)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(),
		[]string{"--safehouse-add-dirs", "/home/a", "--safehouse-add-dirs", "/home/b"},
		dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--add-dirs=/home/a", "--add-dirs=/home/b", "--",
		"claude", "--dangerously-skip-permissions", "--agent", "spacedock:first-officer", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
	// The paths ride in the safehouse extra slot (before the inner `--`), never as
	// host passthrough adjacent to the inner claude argv.
	dash := -1
	for i, tok := range fake.launchedArg {
		if tok == "--" {
			dash = i
			break
		}
	}
	if dash < 0 {
		t.Fatalf("no safehouse `--` separator in argv: %v", fake.launchedArg)
	}
	for _, tok := range fake.launchedArg[dash:] {
		if tok == "/home/a" || tok == "/home/b" || strings.HasPrefix(tok, "--safehouse-add-dirs") {
			t.Fatalf("a knob path/token leaked past the safehouse `--` into host passthrough: %v", fake.launchedArg)
		}
	}
}

// TestSafehouseBadValueNamesKnob pins AC-6's clear-error end state: a genuinely
// bad knob value surfaces a knob-named error, not the internal `malformed flag`
// text. An unknown safehouse key is the representative bad value; the error names
// the knob (`--safehouse-bogus`) and Launch is never reached.
func TestSafehouseBadValueNamesKnob(t *testing.T) {
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--safehouse-bogus=x"}, dir, fake, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero for a bad knob value")
	}
	if fake.launchedArg != nil {
		t.Fatalf("Launch invoked on a bad knob value: %v", fake.launchedArg)
	}
	if !strings.Contains(stderr.String(), "--safehouse-bogus") {
		t.Fatalf("error does not name the knob --safehouse-bogus: %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "malformed flag") {
		t.Fatalf("error leaked the internal malformed-flag text: %q", stderr.String())
	}
}

// Count launches so failure cannot hide a second, unprofiled attempt.
type appendProfileHost struct {
	*fakeHost
	calls int
}

func (f *appendProfileHost) Launch(argv, env []string) (int, error) {
	f.calls++
	return f.fakeHost.Launch(argv, env)
}

func TestAppendProfileLiteralParsing(t *testing.T) {
	for _, parse := range []struct {
		name string
		fn   func([]string) (frontDoorArgs, error)
	}{
		{"shared", parseFrontDoorArgs}, {"pi", func(args []string) (frontDoorArgs, error) { fd, _, err := parsePiFrontDoorArgs(args); return fd, err }},
	} {
		for _, value := range []string{"relative.sb", "/absolute path/a,b:=c.sb", "~/$(echo data);*.sb", "", "--host-looking", "--"} {
			for _, args := range [][]string{{"--safehouse-append-profile=" + value}, {"--safehouse-append-profile", value}} {
				t.Run(parse.name+"/"+strings.Join(args, " "), func(t *testing.T) {
					fd, err := parse.fn(args)
					if err != nil {
						t.Fatal(err)
					}
					extra, err := safehouse.TranslateFlags(fd.safehouseFlags)
					if err != nil || !equalArgv(extra, []string{"--append-profile=" + value}) || len(fd.passthrough) != 0 {
						t.Fatalf("literal value lost: fd=%+v extra=%q err=%v", fd, extra, err)
					}
				})
			}
		}
		fd, err := parse.fn([]string{"--safehouse-append-profile=a.sb", "--safehouse-add-dirs-ro=/ro", "--safehouse-append-profile=b.sb", "--safehouse-enable=ssh", "--safehouse-add-dirs=/rw", "--safehouse-append-profile=b.sb"})
		if err != nil {
			t.Fatal(err)
		}
		extra, err := safehouse.TranslateFlags(fd.safehouseFlags)
		want := []string{"--enable=ssh", "--add-dirs=/rw", "--add-dirs-ro=/ro", "--append-profile=a.sb", "--append-profile=b.sb", "--append-profile=b.sb"}
		if err != nil || !equalArgv(extra, want) {
			t.Fatalf("%s grouped order: %q, %v", parse.name, extra, err)
		}
	}
}

func TestAppendProfileLaunchContract(t *testing.T) {
	withExecutablePath(t, executableFixture(t), nil)
	repo, pkg, home, dir := t.TempDir(), t.TempDir(), t.TempDir(), t.TempDir()
	writePiSkillFixtures(t, repo)
	writePiSubagentsFixtures(t, pkg)
	manifest := compatibleManifest(t)
	for _, host := range []string{"claude", "codex", "pi"} {
		t.Run(host, func(t *testing.T) {
			run := func(args []string, exit int, missing bool) (int, []string, []string, int) {
				var out, errout bytes.Buffer
				if host == "pi" {
					ops := piSafehouseReadyOps(repo, pkg)
					ops.launchCode = exit
					if missing {
						delete(ops.lookPath, "safehouse")
					}
					code := runPi(context.Background(), append([]string{"--plugin-dir", repo}, args...), dir, piTestEnv(pkg, home), ops, &out, &errout)
					return code, ops.launched, ops.launchedEnv, ops.launchCalls
				}
				ops := &appendProfileHost{fakeHost: &fakeHost{manifest: manifest, launchCode: exit}}
				look := lookFound
				if missing {
					look = func(string) (string, error) { return "", errors.New("not found") }
				}
				launch := runClaude
				if host == "codex" {
					launch = runCodex
				}
				code := launch(context.Background(), args, dir, ops, look, &out, &errout)
				return code, ops.launchedArg, ops.launchedEnv, ops.calls
			}
			tail := []string{"do task", "--", "--model", "example"}
			for _, inside := range []string{"", "agent-safehouse"} {
				t.Setenv("APP_SANDBOX_CONTAINER_ID", inside)
				code, baseline, baselineEnv, calls := run(append([]string{"--safehouse"}, tail...), 0, false)
				if code != 0 || calls != 1 {
					t.Fatalf("baseline exit=%d calls=%d", code, calls)
				}
				for _, exit := range []int{0, 23} {
					args := append([]string{"--safehouse-append-profile=first,a:=b.sb", "--safehouse-append-profile", "second.sb", "--safehouse-append-profile=second.sb", "--safehouse-append-profile="}, tail...)
					code, argv, env, calls := run(args, exit, false)
					var stripped, profiles []string
					before := true
					for _, arg := range argv {
						if arg == "--" {
							before = false
						}
						if before && strings.HasPrefix(arg, "--append-profile=") {
							profiles = append(profiles, arg)
						} else {
							stripped = append(stripped, arg)
						}
					}
					want := []string{"--append-profile=first,a:=b.sb", "--append-profile=second.sb", "--append-profile=second.sb", "--append-profile="}
					if code != exit || calls != 1 || !equalArgv(profiles, want) || !equalArgv(stripped, baseline) || !equalArgv(env, baselineEnv) {
						t.Fatalf("exit=%d calls=%d profiles=%q argv=%q; baseline=%q", code, calls, profiles, argv, baseline)
					}
				}
			}
			for _, missingBinary := range []bool{false, true} {
				arg := "--safehouse-append-profile"
				if missingBinary {
					arg += "=file.sb"
				}
				code, argv, _, calls := run([]string{arg}, 0, missingBinary)
				if code == 0 || calls != 0 || len(argv) != 0 {
					t.Fatalf("bad request launched: exit=%d calls=%d argv=%q", code, calls, argv)
				}
			}
			code, argv, _, calls := run([]string{"--", "--safehouse-append-profile=host.sb"}, 0, false)
			if code != 0 || calls != 1 || argv[0] != host || !strings.Contains(strings.Join(argv, " "), "--safehouse-append-profile=host.sb") {
				t.Fatalf("delimiter changed: exit=%d calls=%d argv=%q", code, calls, argv)
			}
		})
	}
}
