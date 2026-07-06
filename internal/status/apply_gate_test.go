// ABOUTME: status --apply-gate turns an explicit gate verdict into the
// ABOUTME: workflow's durable stage transition without teaching callers --set.
package status

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplyGateApproveAdvancesFromGatedStage(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	out, errOut, code := runNative(t, root, env,
		"--workflow-dir", root,
		"--apply-gate",
		"--gate", "helm-gate-123",
		"--entity", "002-vendor-script",
		"--verdict", "approve",
	)

	if code != 0 {
		t.Fatalf("apply-gate approve exit=%d stderr=%q", code, errOut)
	}
	for _, want := range []string{
		"apply-gate slug=002-vendor-script",
		"gate=helm-gate-123",
		"verdict=approve",
		"status=ideation->implementation",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout missing %q:\n%s", want, out)
		}
	}
	fm := readFrontmatter(t, filepath.Join(root, "002-vendor-script.md"))
	for _, want := range []string{
		"status: implementation",
		"gate-id: helm-gate-123",
		"gate-verdict: approve",
	} {
		if !strings.Contains(fm, want) {
			t.Fatalf("frontmatter missing %q:\n%s", want, fm)
		}
	}
}

func TestApplyGateRejectRoutesToFeedbackTarget(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "seq-workflow", map[string]string{
		"README.md": `---
entity-type: task
entity-label: task
entity-label-plural: tasks
id-style: slug
stages:
  defaults:
    worktree: false
    concurrency: 1
  states:
    - name: implementation
      initial: true
    - name: validation
      feedback-to: implementation
      gate: true
    - name: done
      terminal: true
---
# Feedback Workflow
`,
		"needs-fix.md": "---\nid: needs-fix\ntitle: Needs fix\nstatus: validation\nscore: \"0.8\"\nsource: test\n---\n# Needs fix\n",
	})

	out, errOut, code := runNative(t, root, env,
		"--workflow-dir", root,
		"--apply-gate",
		"--gate", "review-gate",
		"--entity", "needs-fix",
		"--verdict", "reject",
	)

	if code != 0 {
		t.Fatalf("apply-gate reject exit=%d stderr=%q", code, errOut)
	}
	for _, want := range []string{
		"apply-gate slug=needs-fix",
		"verdict=reject",
		"status=validation->implementation",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout missing %q:\n%s", want, out)
		}
	}
	fm := readFrontmatter(t, filepath.Join(root, "needs-fix.md"))
	for _, want := range []string{
		"status: implementation",
		"gate-id: review-gate",
		"gate-verdict: reject",
	} {
		if !strings.Contains(fm, want) {
			t.Fatalf("frontmatter missing %q:\n%s", want, fm)
		}
	}
}

func TestApplyGateRejectTrimsFeedbackTarget(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "seq-workflow", map[string]string{
		"README.md": `---
entity-type: task
entity-label: task
entity-label-plural: tasks
id-style: slug
stages:
  states:
    - name: implementation
      initial: true
    - name: validation
      feedback-to: " implementation "
      gate: true
    - name: done
      terminal: true
---
# Feedback Workflow
`,
		"needs-fix.md": "---\nid: needs-fix\ntitle: Needs fix\nstatus: validation\nscore: \"0.8\"\nsource: test\n---\n# Needs fix\n",
	})

	out, errOut, code := runNative(t, root, env,
		"--workflow-dir", root,
		"--apply-gate",
		"--gate", "review-gate",
		"--entity", "needs-fix",
		"--verdict", "reject",
	)

	if code != 0 {
		t.Fatalf("apply-gate reject with padded feedback-to exit=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(out, "status=validation->implementation") {
		t.Fatalf("stdout should show trimmed feedback target:\n%s", out)
	}
	fm := readFrontmatter(t, filepath.Join(root, "needs-fix.md"))
	if !strings.Contains(fm, "status: implementation") || strings.Contains(fm, "status: ' implementation '") {
		t.Fatalf("frontmatter should use trimmed feedback target:\n%s", fm)
	}
}

func TestApplyGateRejectsConflictingActions(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	cases := []struct {
		name string
		args []string
	}{
		{
			name: "read",
			args: []string{"--read", "002-vendor-script"},
		},
		{
			name: "archive",
			args: []string{"--archive", "002-vendor-script"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"--workflow-dir", root}, tc.args...)
			args = append(args,
				"--apply-gate",
				"--gate", "conflict-gate",
				"--entity", "002-vendor-script",
				"--verdict", "approve",
			)
			out, errOut, code := runNative(t, root, env, args...)

			if code != 1 {
				t.Fatalf("conflicting apply-gate action exit=%d, want 1 stdout=%q stderr=%q", code, out, errOut)
			}
			if !strings.Contains(errOut, "--apply-gate cannot be combined") {
				t.Fatalf("stderr should explain apply-gate incompatibility, got %q", errOut)
			}
			fm := readFrontmatter(t, filepath.Join(root, "002-vendor-script.md"))
			if !strings.Contains(fm, "status: ideation") || strings.Contains(fm, "gate-id:") {
				t.Fatalf("conflicting action mutated entity:\n%s", fm)
			}
		})
	}
}

func TestApplyGateRejectsGateFlagsWithoutApplyGate(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	cases := []string{"--gate", "--entity", "--verdict"}
	for _, flag := range cases {
		flag := flag
		t.Run(flag, func(t *testing.T) {
			out, errOut, code := runNative(t, root, env,
				"--workflow-dir", root,
				flag, "partial-value",
			)

			if code != 1 {
				t.Fatalf("%s without --apply-gate exit=%d, want 1 stdout=%q stderr=%q", flag, code, out, errOut)
			}
			if !strings.Contains(errOut, flag+" requires --apply-gate") {
				t.Fatalf("stderr should explain partial apply-gate flag, got %q", errOut)
			}
			if out != "" {
				t.Fatalf("partial apply-gate flag should not fall through to read output, got %q", out)
			}
		})
	}
}

func TestApplyGateJSONUsesGateIDKey(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	out, errOut, code := runNative(t, root, env,
		"--workflow-dir", root,
		"--apply-gate",
		"--gate", "helm-gate-123",
		"--entity", "002-vendor-script",
		"--verdict", "approve",
		"--json",
	)

	if code != 0 {
		t.Fatalf("apply-gate --json exit=%d stderr=%q", code, errOut)
	}
	var payload map[string]string
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse apply-gate --json: %v\n%s", err, out)
	}
	if payload["gate_id"] != "helm-gate-123" {
		t.Fatalf("gate_id = %q, want helm-gate-123; payload=%v", payload["gate_id"], payload)
	}
	if _, ok := payload["gate"]; ok {
		t.Fatalf("apply-gate --json should not emit ambiguous gate key: %v", payload)
	}
}

func TestApplyGateRefusesNonGateStage(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	out, errOut, code := runNative(t, root, env,
		"--workflow-dir", root,
		"--apply-gate",
		"--gate", "bad-gate",
		"--entity", "003-wire-cli",
		"--verdict", "approve",
	)

	if code != 1 {
		t.Fatalf("apply-gate non-gate exit=%d, want 1 stdout=%q stderr=%q", code, out, errOut)
	}
	if !strings.Contains(errOut, "is not at a gate stage") {
		t.Fatalf("stderr should explain non-gate refusal, got %q", errOut)
	}
	fm := readFrontmatter(t, filepath.Join(root, "003-wire-cli", "index.md"))
	if !strings.Contains(fm, "status: implementation") || strings.Contains(fm, "gate-id:") {
		t.Fatalf("non-gate refusal mutated entity:\n%s", fm)
	}
}
