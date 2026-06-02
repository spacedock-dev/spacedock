// ABOUTME: Table test for parseFrontDoorArgs — the Option-2 grammar that takes the
// ABOUTME: task before --, host flags after --, and the safehouse/skip knobs anywhere.
package cli

import (
	"strings"
	"testing"
)

// TestParseFrontDoorArgs pins the Option-2 front-door grammar (AC-3 + AC-6). The
// task is the joined non-flag positionals BEFORE `--`; host value-taking flags
// ride AFTER `--` and forward verbatim as passthrough; the spacedock-owned flags
// (--skip-contract-check, --safehouse, the three repeatable --safehouse-* knobs)
// are consumed wherever they appear before `--` in BOTH space and equals form and
// are never forwarded.
func TestParseFrontDoorArgs(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		passthrough    []string
		task           string
		hasTask        bool
		forceSafehouse bool
		safehouseFlags []string
		skipCheck      bool
	}{
		{name: "bare"},
		{
			name:    "task-positional",
			args:    []string{"do the thing"},
			task:    "do the thing",
			hasTask: true,
		},
		{
			name:    "multi-word-task-joins",
			args:    []string{"do", "the", "thing"},
			task:    "do the thing",
			hasTask: true,
		},
		{
			name:        "task-then-fenced-host-flag",
			args:        []string{"do the thing", "--", "--plugin-dir", "/p"},
			passthrough: []string{"--plugin-dir", "/p"},
			task:        "do the thing",
			hasTask:     true,
		},
		{
			name:        "host-flags-after-fence-no-task",
			args:        []string{"--", "--model", "gpt-x"},
			passthrough: []string{"--model", "gpt-x"},
		},
		{
			name:      "skip-contract-check-consumed",
			args:      []string{"--skip-contract-check"},
			skipCheck: true,
		},
		{
			name:           "bare-safehouse-consumed",
			args:           []string{"--safehouse"},
			forceSafehouse: true,
		},
		{
			name:           "safehouse-knob-equals-form",
			args:           []string{"--safehouse-enable=docker"},
			safehouseFlags: []string{"enable=docker"},
		},
		{
			name:           "safehouse-knob-space-form",
			args:           []string{"--safehouse-enable", "docker"},
			safehouseFlags: []string{"enable=docker"},
		},
		{
			name:           "knob-and-task-and-fenced-host-flags",
			args:           []string{"--safehouse-enable=ssh", "task text", "--", "--model", "x"},
			safehouseFlags: []string{"enable=ssh"},
			passthrough:    []string{"--model", "x"},
			task:           "task text",
			hasTask:        true,
		},
		{
			name:        "plugin-dir-before-dash-space-form",
			args:        []string{"--plugin-dir", "/p"},
			passthrough: []string{"--plugin-dir", "/p"},
		},
		{
			name:        "plugin-dir-before-dash-equals-form",
			args:        []string{"--plugin-dir=/p"},
			passthrough: []string{"--plugin-dir", "/p"},
		},
		{
			name:        "plugin-dir-before-dash-repeated",
			args:        []string{"--plugin-dir", "/a", "--plugin-dir=/b"},
			passthrough: []string{"--plugin-dir", "/a", "--plugin-dir", "/b"},
		},
		{
			name:        "plugin-dir-before-dash-then-task",
			args:        []string{"--plugin-dir", "/p", "review the PRs"},
			passthrough: []string{"--plugin-dir", "/p"},
			task:        "review the PRs",
			hasTask:     true,
		},
		{
			name:        "plugin-dir-before-dash-with-skip-and-task",
			args:        []string{"--plugin-dir", "/p", "--skip-contract-check", "do it"},
			passthrough: []string{"--plugin-dir", "/p"},
			skipCheck:   true,
			task:        "do it",
			hasTask:     true,
		},
		{
			name:        "plugin-dir-before-and-after-dash-both-forward",
			args:        []string{"--plugin-dir", "/before", "--", "--plugin-dir", "/after"},
			passthrough: []string{"--plugin-dir", "/before", "--plugin-dir", "/after"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fd, err := parseFrontDoorArgs(tc.args)
			if err != nil {
				t.Fatalf("parseFrontDoorArgs(%v) err = %v, want nil", tc.args, err)
			}
			if !equalArgv(fd.passthrough, tc.passthrough) {
				t.Errorf("passthrough = %v, want %v", fd.passthrough, tc.passthrough)
			}
			if fd.task != tc.task || fd.hasTask != tc.hasTask {
				t.Errorf("task = (%q,%v), want (%q,%v)", fd.task, fd.hasTask, tc.task, tc.hasTask)
			}
			if fd.forceSafehouse != tc.forceSafehouse {
				t.Errorf("forceSafehouse = %v, want %v", fd.forceSafehouse, tc.forceSafehouse)
			}
			if !equalArgv(fd.safehouseFlags, tc.safehouseFlags) {
				t.Errorf("safehouseFlags = %v, want %v", fd.safehouseFlags, tc.safehouseFlags)
			}
			if fd.skipCheck != tc.skipCheck {
				t.Errorf("skipCheck = %v, want %v", fd.skipCheck, tc.skipCheck)
			}
		})
	}
}

// TestStrayHostFlagBeforeDashErrors pins the spike WINNER's loud-failure boundary:
// only --plugin-dir is promoted as a spacedock-parsed flag before `--`; every other
// host flag left before `--` is an UNKNOWN flag and produces a parse error (not a
// silent mis-split, the rejected generic-pre-pass failure mode). This is exactly why
// the live test must move its non-plugin-dir host flags after `--`.
func TestStrayHostFlagBeforeDashErrors(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "shorthand-p-before-dash", args: []string{"-p", "Drive."}},
		{name: "model-before-dash", args: []string{"--model", "sonnet"}},
		{name: "plugin-dir-then-stray-host-flag", args: []string{"--plugin-dir", "/p", "--model", "sonnet"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := parseFrontDoorArgs(tc.args); err == nil {
				t.Fatalf("parseFrontDoorArgs(%v) err = nil, want a parse error for the stray host flag", tc.args)
			}
		})
	}
}

// TestBareValuelessPluginDirErrors pins the load-bearing loud-error property that
// chose StringArray over ParseErrorsWhitelist.UnknownFlags: a bare before-`--`
// `--plugin-dir` with no following value is a LOUD parse error naming the flag
// (pflag knows its arity), NOT a silent drop. Under UnknownFlags the missing value
// would vanish without a trace — the worst outcome the spike rejected. The returned
// error means runClaude/runCodex return 1 before reaching Launch (no launch).
func TestBareValuelessPluginDirErrors(t *testing.T) {
	_, err := parseFrontDoorArgs([]string{"--plugin-dir"})
	if err == nil {
		t.Fatalf("parseFrontDoorArgs([--plugin-dir]) err = nil, want a loud arity error (not a silent drop)")
	}
	if !strings.Contains(err.Error(), "--plugin-dir") {
		t.Fatalf("error %q does not name --plugin-dir; the loud-error property must be specific", err.Error())
	}
}
