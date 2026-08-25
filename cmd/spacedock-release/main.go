// ABOUTME: Release-pipeline CLI: `stamp-version` writes the release version into
// ABOUTME: the plugin.json manifests for a release.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
	"github.com/spacedock-dev/spacedock/internal/release"
)

// main is the release-tooling entry point invoked by CI (not the user binary).
//
// Usage:
//
//	spacedock-release stamp-version <release-version> <manifest-or-prose> [<manifest-or-prose> ...]
//	spacedock-release dev-preversion <stable-version>
//	spacedock-release journey-delta <previous-ledger.json> --metrics-dir <dir> --pr <number>
//	spacedock-release e2e-gate <release-commit-sha>
//
// stamp-version rewrites each argument to the release version: a `.json`
// plugin.json gets its top-level `version` field rewritten (AC-4); a `.md`
// argument is the FO shared-core prose, whose single release-stamped "required
// binary minor" literal is rewritten to the release's major.minor (D5) —
// erroring unless the literal appears exactly once. All rewrite in place.
// dev-preversion prints the post-release dev pre-version
// (X.(Y+1).0-pre1) the stable-tag edge advance stamps onto `next`.
// edge-advance-decision prints `advance` or `skip` (exit 0 either way) deciding
// whether a tag advances the edge line: it computes the tag's target edge
// version and skips unless that is strictly greater than the highest known
// release version (fed from a manifest release.yml synthesizes via
// highest-known-edge-version), so the auto-pre0 cut no-ops on an old-line/patch
// tag. highest-known-edge-version prints the greatest version among its <tag>
// arguments (prereleases included), or nothing when none parses — the git-tag
// scan that replaced reading `next`'s manifest. edge-pre0-version prints the
// auto-cut edge prerelease version (X.(Y+1).0-pre0) the stable path tags on the
// greened release commit. journey-delta
// renders and posts the per-PR journey-cost delta against the previously
// published release ledger's latest-by-captured_at baseline per scenario/model,
// updating a single sticky PR comment found by its HTML marker. The AC-4
// release-ledger backfill (a captain-flagged one-time manual procedure, never a
// CI step) is a documented runbook, not a shipped subcommand: extraction reuses
// the exported ensigncycle.BuildShallowBootWindowRecord via a throwaway `go run`
// script against each archived stream, and the pre-upload safeguard is
// `jq -S .scenarios <original> | diff - <(jq -S .scenarios <rebuilt>)` —
// any diff output means STOP, do not upload. e2e-gate is the
// release-time precondition: it passes (exit 0) only when a conclusion:success
// Runtime Live E2E run exists for the commit, or when SPACEDOCK_E2E_GATE_WAIVER
// is set, and blocks the cut (exit 1) otherwise. manifest-tag-gate blocks the cut
// unless every tagged `.json` manifest's version equals the tag semver AND every
// tagged `.md` prose's stamped minor equals the tag's major.minor (the
// stamp-then-tag ordering). stable-regression-gate blocks the cut when the tag is
// older than the version in the plugin manifest that `refs/heads/stable` points
// at. Such a tag moves /releases/latest and the stable Homebrew cask DOWN.
// notes summarizes the commit log
// since the last tag into clean release notes and, on confirmation, cuts the
// annotated tag whose body carries them (CI extracts that body and feeds
// goreleaser via --release-notes).
func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "stamp-version":
		os.Exit(stampVersion(os.Args[2:]))
	case "dev-preversion":
		os.Exit(devPreversion(os.Args[2:]))
	case "edge-advance-decision":
		os.Exit(edgeAdvanceDecision(os.Args[2:]))
	case "highest-known-edge-version":
		os.Exit(highestKnownEdgeVersion(os.Args[2:]))
	case "edge-pre0-version":
		os.Exit(edgePre0Version(os.Args[2:]))
	case "journey-costs":
		os.Exit(journeyCosts(os.Args[2:]))
	case "journey-delta":
		os.Exit(journeyDelta(os.Args[2:], ghFindComment, ghPostComment))
	case "e2e-gate":
		os.Exit(runE2EGate(os.Args[2:], ghRunListForCommit))
	case "manifest-tag-gate":
		os.Exit(runManifestTagGate(os.Args[2:]))
	case "stable-regression-gate":
		os.Exit(runStableRegressionGate(os.Args[2:]))
	case "notes":
		os.Exit(notes(os.Args[2:]))
	default:
		fmt.Fprintf(os.Stderr, "spacedock-release: unknown command %q\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func journeyCosts(args []string) int {
	if len(args) < 5 {
		fmt.Fprintln(os.Stderr, "spacedock-release journey-costs: need <release-version> --metrics-dir <dir> --out <path>")
		return 2
	}
	version := strings.TrimPrefix(args[0], "v")
	var metricsDir, out string
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--metrics-dir":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "spacedock-release journey-costs: --metrics-dir needs a value")
				return 2
			}
			metricsDir = args[i]
		case "--out":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "spacedock-release journey-costs: --out needs a value")
				return 2
			}
			out = args[i]
		default:
			fmt.Fprintf(os.Stderr, "spacedock-release journey-costs: unknown argument %q\n", args[i])
			return 2
		}
	}
	if version == "" || metricsDir == "" || out == "" {
		fmt.Fprintln(os.Stderr, "spacedock-release journey-costs: need <release-version> --metrics-dir <dir> --out <path>")
		return 2
	}
	wantName := "journey-costs-v" + version + ".json"
	if filepath.Base(out) != wantName {
		fmt.Fprintf(os.Stderr, "spacedock-release journey-costs: output filename must be %s\n", wantName)
		return 1
	}
	records, err := journeymetrics.ReadRecordsDir(metricsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read journey metrics: %v\n", err)
		return 1
	}
	ledger, err := journeymetrics.AggregateLedger(version, records, time.Now().UTC())
	if err != nil {
		fmt.Fprintf(os.Stderr, "aggregate journey metrics: %v\n", err)
		return 1
	}
	data, err := journeymetrics.MarshalLedger(ledger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render journey ledger: %v\n", err)
		return 1
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "create output dir: %v\n", err)
		return 1
	}
	if err := os.WriteFile(out, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
		return 1
	}
	fmt.Printf("wrote %s (%d scenarios, %d observations)\n", out, ledger.Summary.ScenarioCount, ledger.Summary.ObservationCount)
	return 0
}

func stampVersion(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "spacedock-release stamp-version: need <release-version> <manifest> [<manifest> ...]")
		return 2
	}
	version, manifests := args[0], args[1:]
	for _, path := range manifests {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
			return 1
		}
		// D5: a `.md` argument is the FO shared-core prose — one release-stamped
		// "required binary minor" literal, rewritten by StampProseVersion. Every
		// other extension (the plugin.json manifests) keeps the existing JSON
		// full-version stamp. One invocation, one atomic multi-file rewrite.
		var out []byte
		if strings.HasSuffix(path, ".md") {
			out, err = release.StampProseVersion(data, version)
		} else {
			out, err = release.StampVersion(data, version)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "stamp %s: %v\n", path, err)
			return 1
		}
		if err := os.WriteFile(path, out, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", path, err)
			return 1
		}
		fmt.Printf("stamped %s version=%s\n", path, version)
	}
	return 0
}

// devPreversion prints the post-release dev pre-version for the `next` edge line
// (X.(Y+1).0-pre1) computed from a just-released stable version. release.yml's
// stable-tag path captures this stdout to stamp `next` PAST the released stable
// version. It takes exactly one <stable-version> arg and errors on a hyphenated
// or malformed input, since it runs only on the `!contains(github.ref, '-')`
// branch that already guarantees a bare X.Y.Z.
func devPreversion(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "spacedock-release dev-preversion: need exactly one <stable-version> (e.g. 0.24.0)")
		return 2
	}
	version, err := release.DevPreVersion(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "dev-preversion: %v\n", err)
		return 1
	}
	fmt.Println(version)
	return 0
}

// notes summarizes the commits since the last tag into clean release notes and,
// on the captain's confirmation, cuts an annotated tag whose body IS those
// notes. The tag is created locally only — pushing it (which fires the release
// build) stays a deliberate manual step the captain runs after review.
func notes(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "spacedock-release notes: need exactly one <release-version> (e.g. 0.19.2)")
		return 2
	}
	version := strings.TrimPrefix(args[0], "v")
	tag := "v" + version

	io := release.NotesIO{
		RawLog: commitLogSinceLastTag,
		Claude: runClaude,
	}
	proposed, err := release.GenerateNotes(version, io)
	if err != nil {
		fmt.Fprintf(os.Stderr, "generate notes: %v\n", err)
		return 1
	}

	fmt.Println("==========================================")
	fmt.Printf("PROPOSED RELEASE NOTES FOR %s\n", tag)
	fmt.Println("==========================================")
	fmt.Println(proposed)
	fmt.Println("==========================================")

	tagIO := release.TagIO{
		Confirm: confirmNotes,
		CutTag: func(body string) error {
			return cutAnnotatedTag(tag, version, body)
		},
	}
	if err := release.ConfirmAndTag(proposed, tagIO); err != nil {
		fmt.Fprintf(os.Stderr, "cut tag: %v\n", err)
		return 1
	}
	return 0
}

// commitLogSinceLastTag returns the `git log --oneline` range from the most
// recent tag to HEAD (or all of HEAD when no tag exists yet).
func commitLogSinceLastTag() (string, error) {
	prev, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	rangeSpec := "HEAD"
	if err == nil {
		if p := strings.TrimSpace(string(prev)); p != "" {
			rangeSpec = p + "..HEAD"
		}
	}
	out, err := exec.Command("git", "log", rangeSpec, "--oneline", "--no-decorate").Output()
	if err != nil {
		return "", fmt.Errorf("git log %s: %w", rangeSpec, err)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

// runClaude pipes the filtered log through `claude -p <prompt> --model opus
// --effort low`. A non-nil error (claude absent or failing) tells GenerateNotes
// to fall back to the filtered raw log.
func runClaude(prompt, input string) (string, error) {
	if _, err := exec.LookPath("claude"); err != nil {
		return "", err
	}
	cmd := exec.Command("claude", "-p", prompt, "--model", "opus", "--effort", "low")
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\n"), nil
}

// confirmNotes lets the captain review the proposed notes and choose to cut the
// tag (y), edit the notes in $EDITOR before tagging (e), or decline (anything
// else → no tag). The edited body is what gets tagged.
func confirmNotes(proposed string) (string, bool) {
	fmt.Print("Cut the annotated tag with these notes? [y/e=edit/N] ")
	r := bufio.NewReader(os.Stdin)
	line, _ := r.ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return proposed, true
	case "e", "edit":
		edited, err := editInEditor(proposed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "edit failed (%v); no tag cut\n", err)
			return proposed, false
		}
		return edited, true
	default:
		fmt.Fprintln(os.Stderr, "declined; no tag cut")
		return proposed, false
	}
}

// editInEditor opens the proposed notes in $EDITOR (falling back to vi) and
// returns the captain's edited text.
func editInEditor(proposed string) (string, error) {
	f, err := os.CreateTemp("", "spacedock-release-notes-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(proposed); err != nil {
		f.Close()
		return "", err
	}
	f.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, f.Name())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	edited, err := os.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(edited), "\n"), nil
}

// cutAnnotatedTag creates the annotated tag locally with the notes in the tag
// message BODY (under a `Release <version>` subject), so `git tag -l
// --format='%(contents:body)'` in CI round-trips exactly these notes into
// goreleaser's --release-notes. The tag is not pushed — that stays a manual
// step.
func cutAnnotatedTag(tag, version, body string) error {
	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("refusing to cut %s with empty notes body", tag)
	}
	if err := exec.Command("git", release.AnnotatedTagArgs(tag, version, body)...).Run(); err != nil {
		return err
	}
	fmt.Printf("created annotated tag %s (local only)\n", tag)
	fmt.Printf("push to trigger the release build: git push origin %s\n", tag)
	return nil
}

func usage() {
	fmt.Fprint(os.Stderr, `spacedock-release is the release-pipeline version tool.

Usage:
  spacedock-release stamp-version <release-version> <manifest-or-prose> [<manifest-or-prose> ...]
  spacedock-release dev-preversion <stable-version>
  spacedock-release edge-advance-decision <tag> <known-version-plugin.json>
  spacedock-release highest-known-edge-version [<tag> ...]
  spacedock-release edge-pre0-version <stable-version>
  spacedock-release journey-costs <release-version> --metrics-dir <dir> --out <path>
  spacedock-release journey-delta <previous-ledger.json> --metrics-dir <dir> --pr <number>
  spacedock-release e2e-gate <release-commit-sha>
  spacedock-release manifest-tag-gate <tag> <manifest-or-prose> [<manifest-or-prose> ...]
  spacedock-release stable-regression-gate <tag> <stable-plugin.json>
  spacedock-release notes <release-version>
`)
}
