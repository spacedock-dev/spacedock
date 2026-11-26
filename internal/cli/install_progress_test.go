package cli

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestInstallVerbosePreservesOutputAndNamesSkippedWork(t *testing.T) {
	t.Setenv("PATH", "/usr/bin:/bin")
	stubClaudeInventory(t, "[]")
	var ordinaryOut, ordinaryErr, stdout, stderr bytes.Buffer
	args := []string{"install", "--check", "--host", "claude"}
	ordinaryCode := Run(args, &ordinaryOut, &ordinaryErr)
	code := Run(append(args, "--verbose"), &stdout, &stderr)

	if code != ordinaryCode || stdout.String() != ordinaryOut.String() {
		t.Fatalf("verbose install changed its result: (%d, %q) vs (%d, %q)", code, stdout.String(), ordinaryCode, ordinaryOut.String())
	}
	for _, line := range []string{
		"Skipped: Install Claude Code plugins (check only)",
		"Start: Resolve Claude Code plugin manifest",
		"Done: Resolve Claude Code plugin manifest",
		"Total:",
		"(exit 0)",
	} {
		if !strings.Contains(stderr.String(), line) {
			t.Errorf("stderr did not contain %q:\n%s", line, &stderr)
		}
	}
}

func TestInstallVerboseUsageDoesNotStartTiming(t *testing.T) {
	for _, args := range [][]string{
		{"install", "--verbose", "--host"},
		{"install", "--host", "--verbose"},
		{"install", "--verbose", "--plugin-dir"},
		{"install", "--plugin-dir", "--verbose"},
		{"install", "--verbose", "--unknown"},
		{"doctor", "--verbose"},
	} {
		var stdout, stderr bytes.Buffer
		if code := Run(args, &stdout, &stderr); code != 2 || strings.Contains(stderr.String(), "[+") {
			t.Errorf("%v: wanted untimed exit 2, got (%d, %q)", args, code, &stderr)
		}
	}
}

func sameInstallSteps(got, want []installStep) bool {
	return slices.EqualFunc(got, want, func(a, b installStep) bool {
		a.operation, b.operation = "", ""
		return reflect.DeepEqual(a, b)
	})
}

func TestInstallVerboseUsesCommandClock(t *testing.T) {
	t.Setenv("PATH", "/usr/bin:/bin")
	stubClaudeInventory(t, "[]")
	now := time.Unix(0, 0)
	clock := func() time.Time {
		now = now.Add(125 * time.Millisecond)
		return now
	}
	var stdout, stderr bytes.Buffer
	args := []string{"--verbose", "--host", "claude", "--check"}
	if code := runInstallCommand(context.Background(), args, os.Environ(), t.TempDir(), &stdout, &stderr, clock); code != 0 {
		t.Fatalf("install exited %d: %s", code, &stderr)
	}
	want := "[+0.125s] Skipped: Install Claude Code plugins (check only)\n" +
		"[+0.250s] Start: Resolve Claude Code plugin manifest\n" +
		"[+0.375s] Done: Resolve Claude Code plugin manifest (0.125s)\n"
	if !strings.HasPrefix(stderr.String(), want) {
		t.Errorf("trace did not use the command clock for offsets and durations:\n%s", &stderr)
	}
}
