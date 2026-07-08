package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

// TestDevPreversionCommandPrintsComputedVersion — the subcommand prints exactly
// the post-release dev pre-version (X.(Y+1).0-pre1) to stdout, so release.yml's
// `DEV_VERSION="$(... dev-preversion "$RELEASE_VERSION")"` capture stamps `next`
// PAST the just-released stable version rather than re-stamping it.
func TestDevPreversionCommandPrintsComputedVersion(t *testing.T) {
	out, code := captureStdout(t, func() int { return devPreversion([]string{"0.24.0"}) })
	if code != 0 {
		t.Fatalf("dev-preversion exit = %d, want 0 on a bare stable semver", code)
	}
	if got := strings.TrimSpace(out); got != "0.25.0-pre1" {
		t.Fatalf("dev-preversion stdout = %q, want 0.25.0-pre1", got)
	}
}

// TestDevPreversionCommandRejectsMissingArg — with no version argument the
// subcommand exits with a usage error and prints nothing consumable.
func TestDevPreversionCommandRejectsMissingArg(t *testing.T) {
	if code := devPreversion(nil); code == 0 {
		t.Fatalf("dev-preversion exit = 0 with no argument; want non-zero")
	}
}

// TestDevPreversionCommandRejectsExtraArg — the subcommand takes exactly one
// version; a second argument is a usage error, so a release.yml miswiring fails
// loud rather than silently ignoring the extra token.
func TestDevPreversionCommandRejectsExtraArg(t *testing.T) {
	if code := devPreversion([]string{"0.24.0", "0.25.0"}); code == 0 {
		t.Fatalf("dev-preversion exit = 0 with two arguments; want non-zero")
	}
}

// TestDevPreversionCommandRejectsHyphenatedInput — a pre-release (hyphenated)
// input is malformed for this subcommand (it runs only on the stable branch),
// so it exits non-zero instead of emitting a doubly-suffixed version.
func TestDevPreversionCommandRejectsHyphenatedInput(t *testing.T) {
	if code := devPreversion([]string{"0.24.0-pre1"}); code == 0 {
		t.Fatalf("dev-preversion exit = 0 on a hyphenated input; want non-zero")
	}
}

// captureStdout runs fn with os.Stdout redirected to a pipe and returns what fn
// printed plus fn's return code, so a subcommand whose contract IS its stdout
// (release.yml captures it via $()) is verified on the real bytes it emits.
func captureStdout(t *testing.T, fn func() int) (string, int) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	code := fn()
	w.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read captured stdout: %v", err)
	}
	return string(data), code
}
