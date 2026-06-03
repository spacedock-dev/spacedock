// ABOUTME: Release-pipeline CLI: `stamp-version` writes the release version into
// ABOUTME: the plugin.json manifests; `bump-calendar` advances the marketplace key.
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
//	spacedock-release stamp-version <release-version> <plugin.json> [<plugin.json> ...]
//	spacedock-release bump-calendar <marketplace.json>
//
// stamp-version rewrites each manifest's top-level `version` to the release
// version (AC-4). bump-calendar advances the marketplace plugin entry's calendar
// key to today's `0.0.YYYYMMDDNN` (AC-2d). Both rewrite in place. notes
// summarizes the commit log since the last tag into clean release notes and, on
// confirmation, cuts the annotated tag whose body carries them (CI extracts that
// body and feeds goreleaser via --release-notes).
func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "stamp-version":
		os.Exit(stampVersion(os.Args[2:]))
	case "bump-calendar":
		os.Exit(bumpCalendar(os.Args[2:]))
	case "journey-costs":
		os.Exit(journeyCosts(os.Args[2:]))
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
	fmt.Printf("wrote %s (%d journeys)\n", out, len(records))
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
		out, err := release.StampVersion(data, version)
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

func bumpCalendar(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "spacedock-release bump-calendar: need exactly one <marketplace.json>")
		return 2
	}
	path := args[0]
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
		return 1
	}
	out, err := release.BumpCalendarVersion(data, time.Now())
	if err != nil {
		fmt.Fprintf(os.Stderr, "bump %s: %v\n", path, err)
		return 1
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", path, err)
		return 1
	}
	fmt.Printf("bumped %s\n", path)
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
  spacedock-release stamp-version <release-version> <plugin.json> [<plugin.json> ...]
  spacedock-release bump-calendar <marketplace.json>
  spacedock-release journey-costs <release-version> --metrics-dir <dir> --out <path>
  spacedock-release notes <release-version>
`)
}
