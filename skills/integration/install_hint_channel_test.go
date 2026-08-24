// ABOUTME: Channel-aware binary-absent install hint — extracts the shipped classifier
// ABOUTME: one-liner and the hinted commands from the FO contract and executes them.
package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// The shipped files this test binds together: the boot-resident FO core (which now keeps
// only the version pin and the deferral), the deferred install reference that owns the
// classifier and the commands, and the user-facing doc. Nothing here asserts on their
// wording — the tests extract runnable commands and observe what those commands do.
func sharedCorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "skills", "first-officer", "references", "first-officer-shared-core.md")
}

// installRefPath is the deferred FO install reference. The channel classifier and the
// per-OS install commands live here, not in the boot-resident core: the core defers the
// whole binary-absent arm to this file, sandbox or not.
func installRefPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "skills", "first-officer", "references", "fo-install.md")
}

func installDocPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "docs", "site", "get-started", "install.md")
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// classifierRe matches the shipped classifier in fo-install.md's Channel selection
// section: a one-liner opening by stripping the skills dir off the retained
// {first_officer_base}. Pinning that opening binds the classifier to its documented input;
// a reshape that reads the environment or a marketplace file instead would not match, so
// extraction is itself part of the guard.
var classifierRe = regexp.MustCompile(`(?m)^\s*(R="\$\{B%/skills/first-officer\}";.*)$`)

// extractClassifier returns the one runnable classifier shipped in the install reference.
// The tests EXECUTE it rather than a copy, so contract and test cannot drift.
func extractClassifier(t *testing.T) string {
	t.Helper()
	m := classifierRe.FindAllStringSubmatch(readFile(t, installRefPath(t)), -1)
	if len(m) != 1 {
		t.Fatalf("expected exactly one channel-classifier one-liner in fo-install.md matching %q, found %d",
			classifierRe.String(), len(m))
	}
	return m[0][1]
}

// runClassifier executes the extracted one-liner with B bound to base, returning the
// channel it prints. base is data, never shell text: it is passed through the
// environment so a fixture path can never be interpreted as shell syntax.
func runClassifier(t *testing.T, classifier, base string) string {
	t.Helper()
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh not on PATH; the classifier is a POSIX shell one-liner")
	}
	cmd := exec.Command(sh, "-c", `B="$SPACEDOCK_TEST_BASE"; `+classifier)
	cmd.Env = append(os.Environ(), "SPACEDOCK_TEST_BASE="+base)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run classifier with B=%q: %v\n%s", base, err, out)
	}
	return strings.TrimSpace(string(out))
}

// TestChannelClassifierTable is AC-2: the shipped classifier maps every observed real
// install-path shape to its true channel. Fixture paths are recorded from the real Claude
// plugin cache and installed_plugins.json, so expected values originate in the observed
// install layout, not in the file under test. Each `why` naming a single arm is a row that
// reds when that arm is removed — the two documented counterexamples.
func TestChannelClassifierTable(t *testing.T) {
	classifier := extractClassifier(t)

	const cache = "/Users/clkao/.claude/plugins/cache/"
	cases := []struct {
		base string
		want string
		why  string
	}{
		{cache + "spacedock-edge/spacedock/0.27.0-pre8/skills/first-officer", "edge",
			"edge marketplace, prerelease version — both signals agree"},
		{cache + "spacedock-edge/spacedock/0.23.0-pre/skills/first-officer", "edge",
			"edge marketplace, bare -pre suffix"},
		{cache + "spacedock-edge/spacedock/0.27.0/skills/first-officer", "edge",
			"edge install on an UNSUFFIXED version — only the marketplace arm catches it"},
		{cache + "spacedock/spacedock/0.27.0-pre7-dev/skills/first-officer", "edge",
			"next dev build under the STABLE marketplace name — only the suffix arm catches it"},
		{cache + "spacedock/spacedock/0.26.0/skills/first-officer", "stable",
			"stable marketplace, released version"},
		{cache + "spacedock/spacedock/0.19.1/skills/first-officer", "stable",
			"older stable release"},
		{"/Users/pre-release-tester/.claude/plugins/cache/spacedock/spacedock/0.26.0/skills/first-officer", "stable",
			"adversarial home containing -pre must NOT widen the match to edge"},
		{"/Users/clkao/git/spacedock-research/spacedock-v1/skills/first-officer", "local",
			"--plugin-dir source checkout: no version-shaped directory"},
		{"/Users/pre-release-tester/git/spacedock-v1/skills/first-officer", "local",
			"source checkout under the adversarial home is still local, not edge"},
	}

	for _, tc := range cases {
		got := runClassifier(t, classifier, tc.base)
		if got != tc.want {
			t.Errorf("classify(%q) = %q, want %q (%s)", tc.base, got, tc.want, tc.why)
		}
	}
}

// The install-command extractors, all anchored on the real install.sh URL so an unrelated
// command cannot satisfy them: the contract's backticked stable command, the doc's stable
// command, and the doc's edge command.
var (
	stableCurlRe = regexp.MustCompile("`(curl -fsSL " + installURLRe + ` \| sh)` + "`")
	docStableRe  = regexp.MustCompile(`(?m)^\s*(curl -fsSL ` + installURLRe + ` \| sh)\s*$`)
	// docEdgeRe matches the published edge command in EITHER env-var placement — prefixed
	// on curl or on sh — so that the shape question is settled by running the command
	// (TestEdgeHintDeliversChannelToScript) rather than by which regex happens to match.
	docEdgeRe = regexp.MustCompile(`(?m)^\s*(\S+=\S+ curl -fsSL ` + installURLRe + ` \| sh|curl -fsSL ` + installURLRe + ` \| \S+=\S+ sh)\s*$`)
	// curlFetchRe isolates just the network fetch inside a published command, so a test
	// can swap it for a local stub while leaving every other byte — crucially, anything
	// that decides what the script's environment looks like — exactly as published.
	curlFetchRe = regexp.MustCompile(`curl -fsSL ` + installURLRe)
)

const installURLRe = `https://raw\.githubusercontent\.com/spacedock-dev/spacedock/main/install\.sh`

// TestStableHintMatchesPublishedDoc is AC-3: the stable command the FO contract hints and
// the stable command the docs publish are two independently maintained strings that must
// stay byte-identical. Editing either one alone diverges the pair and fails here.
func TestStableHintMatchesPublishedDoc(t *testing.T) {
	contractM := stableCurlRe.FindStringSubmatch(readFile(t, installRefPath(t)))
	if contractM == nil {
		t.Fatalf("no Linux stable curl command found in fo-install.md matching %q", stableCurlRe.String())
	}
	docM := docStableRe.FindStringSubmatch(readFile(t, installDocPath(t)))
	if docM == nil {
		t.Fatalf("no stable curl command found in install.md matching %q", docStableRe.String())
	}
	if contractM[1] != docM[1] {
		t.Errorf("the FO contract's stable install command and the published doc's have drifted:\n contract: %q\n doc:      %q",
			contractM[1], docM[1])
	}
}

// TestEdgeHintDeliversChannelToScript is the AC-1 offline guard on the edge command's
// SHAPE. A shell variable prefix binds to the FIRST command of a pipeline, so writing
// `SPACEDOCK_CHANNEL=edge curl … | sh` leaves the variable UNSET inside the script and
// silently installs stable — the exact channel skew this whole gate exists to close.
//
// The test runs the PUBLISHED edge command with the network fetch swapped for a local
// script that reports what it actually received, so the oracle is the variable the script
// observes, never the wording of the command. Moving the assignment onto `curl` reds it.
func TestEdgeHintDeliversChannelToScript(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh not on PATH; the install hint is a POSIX shell pipeline")
	}
	docM := docEdgeRe.FindStringSubmatch(readFile(t, installDocPath(t)))
	if docM == nil {
		t.Fatalf("no edge (channel-carrying) curl command found in install.md matching %q", docEdgeRe.String())
	}
	published := docM[1]

	// Swap ONLY the network fetch for a local stub. Every other byte of the published
	// command is preserved — including any variable assignment and which side of the pipe
	// it sits on, which is precisely what determines whether the script sees the channel.
	stub := filepath.Join(t.TempDir(), "install.sh")
	if err := os.WriteFile(stub, []byte("echo \"channel=${SPACEDOCK_CHANNEL:-UNSET}\"\n"), 0o644); err != nil {
		t.Fatalf("write stub install.sh: %v", err)
	}
	runnable := curlFetchRe.ReplaceAllLiteralString(published, "cat "+stub)
	if runnable == published {
		t.Fatalf("could not substitute the fetch in the published edge command: %q", published)
	}

	cmd := exec.Command(sh, "-c", runnable)
	cmd.Env = []string{"PATH=" + os.Getenv("PATH")} // no ambient SPACEDOCK_CHANNEL
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run %q: %v\n%s", runnable, err, out)
	}
	if got := strings.TrimSpace(string(out)); got != "channel=edge" {
		t.Errorf("the published edge install command does not deliver the channel to the script: got %q, want %q\n"+
			"  command: %q\n"+
			"  a variable prefix binds to the first command of a pipeline, so it must sit on `sh`, not on `curl`",
			got, "channel=edge", published)
	}
}

// contractPin reads the binary minor these skills require, straight from the shared core,
// so the cask assertions below track the contract instead of a hardcoded copy.
func contractPin(t *testing.T) (major, minor int) {
	t.Helper()
	m := regexp.MustCompile(`require binary minor (\d+)\.(\d+)`).FindStringSubmatch(readFile(t, sharedCorePath(t)))
	if m == nil {
		t.Fatal("no `require binary minor X.Y` declaration found in the shared core")
	}
	major, _ = strconv.Atoi(m[1])
	minor, _ = strconv.Atoi(m[2])
	return major, minor
}

// tapInfo resolves the tap the contract hints, returning its on-disk path and the cask
// tokens it publishes. `brew tap-info` is used rather than `brew info --cask` on purpose:
// the latter refuses to LOAD a cask from a tap the machine has not `brew trust`ed, which
// would make this leg depend on a per-machine trust decision. Skips (never fails) when
// brew or the tap is unavailable, keeping the offline lane deterministic.
func tapInfo(t *testing.T) (tapPath string, caskTokens []string) {
	t.Helper()
	if testing.Short() {
		t.Skip("-short: cask resolution needs Homebrew and the tap")
	}
	brew, err := exec.LookPath("brew")
	if err != nil {
		t.Skip("brew not on PATH")
	}
	cmd := exec.Command(brew, "tap-info", "--json", "spacedock-dev/tap")
	cmd.Env = append(os.Environ(), "HOMEBREW_NO_AUTO_UPDATE=1", "HOMEBREW_NO_ENV_HINTS=1")
	out, err := cmd.Output()
	if err != nil {
		t.Skipf("cannot resolve tap spacedock-dev/tap (not tapped, or brew unavailable): %v", err)
	}
	var taps []struct {
		Installed  bool     `json:"installed"`
		Path       string   `json:"path"`
		CaskTokens []string `json:"cask_tokens"`
	}
	if err := json.Unmarshal(out, &taps); err != nil {
		t.Skipf("cannot parse brew tap-info output: %v", err)
	}
	if len(taps) == 0 || !taps[0].Installed {
		t.Skip("tap spacedock-dev/tap is not installed on this machine")
	}
	return taps[0].Path, taps[0].CaskTokens
}

// caskVersion returns the version pinned by a cask definition in the resolved tap.
func caskVersion(t *testing.T, tapPath, token string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(tapPath, "Casks", token+".rb"))
	if err != nil {
		t.Skipf("cannot read cask definition %s: %v", token, err)
	}
	m := regexp.MustCompile(`(?m)^\s*version\s+"([^"]+)"`).FindStringSubmatch(string(data))
	if m == nil {
		t.Errorf("cask %s has no parseable version stanza", token)
		return ""
	}
	return m[1]
}

// TestContractCasks covers the macOS half of AC-3 and all of AC-4, sharing one tap
// resolution. Network/tool-gated.
//
//   - names: both cask tokens the contract hints must be real casks the tap publishes.
//     The expected tokens are read out of the shipped contract, so renaming a cask (or
//     the contract hinting a cask that does not exist) fails here.
//   - edge-satisfies-pin: the macOS edge remedy must be a real, pin-satisfying cask
//     rather than the curl fallback. Observed today: edge 0.27.0-pre8 vs stable 0.26.0
//     against a contract pin of 0.27. If the edge cask ever stops tracking prereleases
//     its version falls back to the stable line and this reds — exactly when the macOS
//     hint would need to revert to the curl form.
func TestContractCasks(t *testing.T) {
	core := readFile(t, installRefPath(t))
	hinted := regexp.MustCompile("brew install spacedock-dev/tap/(\\S+?)`").FindAllStringSubmatch(core, -1)
	if len(hinted) != 2 {
		t.Fatalf("expected fo-install.md to hint exactly two casks (stable and edge), found %d", len(hinted))
	}
	tapPath, tokens := tapInfo(t)

	t.Run("names", func(t *testing.T) {
		published := make(map[string]bool, len(tokens))
		for _, tok := range tokens {
			published[tok] = true
		}
		for _, m := range hinted {
			want := "spacedock-dev/tap/" + m[1]
			if !published[want] {
				t.Errorf("fo-install.md hints `brew install spacedock-dev/tap/%s`, but the tap publishes no cask %q (has: %v)",
					m[1], want, tokens)
			}
		}
	})

	t.Run("edge-satisfies-pin", func(t *testing.T) {
		wantMajor, wantMinor := contractPin(t)
		v := caskVersion(t, tapPath, "spacedock@next")
		if v == "" {
			return
		}
		t.Logf("edge cask %s vs stable cask %s, contract pin %d.%d",
			v, caskVersion(t, tapPath, "spacedock"), wantMajor, wantMinor)
		m := regexp.MustCompile(`^(\d+)\.(\d+)`).FindStringSubmatch(v)
		if m == nil {
			t.Fatalf("edge cask version %q does not parse as major.minor", v)
		}
		major, _ := strconv.Atoi(m[1])
		minor, _ := strconv.Atoi(m[2])
		if major < wantMajor || (major == wantMajor && minor < wantMinor) {
			t.Errorf("edge cask spacedock@next is at %s, below the contract pin %d.%d — "+
				"the macOS edge hint would install a binary the FO gate then rejects",
				v, wantMajor, wantMinor)
		}
	})
}
