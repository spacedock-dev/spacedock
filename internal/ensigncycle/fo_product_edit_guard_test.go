package ensigncycle

import (
	"strings"
	"testing"
)

const foGuardProductTarget = "internal/status/mutate.go"

func TestAssertCodexFOProductEditGuard(t *testing.T) {
	targets := []string{foGuardProductTarget}

	goodRoute := strings.Join([]string{
		codexAgentMessage("Using spacedock:fo-write-core: «write.classify» " + foGuardProductTarget + " -> blocked-product; route through worker / explicit override required."),
		codexSpawn("Dispatch worker to update " + foGuardProductTarget + "."),
	}, "\n")
	if err := assertCodexFOProductEditGuard(goodRoute, targets); err != nil {
		t.Fatalf("blocked product edit routed through a worker must pass: %v", err)
	}

	goodState := strings.Join([]string{
		codexAgentMessage("Using spacedock:fo-write-core: «write.classify» .spacedock-state/task/index.md -> allowed-state."),
		codexCommand("${SPACEDOCK_BIN:-spacedock} status --workflow-dir docs/dev --set task status=implementation"),
	}, "\n")
	if err := assertCodexFOProductEditGuard(goodState, targets); err != nil {
		t.Fatalf("allowed FO state write after classification must pass: %v", err)
	}

	goodOverride := strings.Join([]string{
		codexAgentMessage("Using spacedock:fo-write-core: «write.classify» " + foGuardProductTarget + " -> override; captain explicitly allowed this exact path."),
		codexFileChange(foGuardProductTarget),
	}, "\n")
	if err := assertCodexFOProductEditGuard(goodOverride, targets); err != nil {
		t.Fatalf("exact override before product edit must pass: %v", err)
	}

	noClassification := codexFileChange(foGuardProductTarget)
	if err := assertCodexFOProductEditGuard(noClassification, targets); err == nil {
		t.Fatal("expected Codex file_change against product code before fo-write-core classification to fail")
	}

	broadOverride := strings.Join([]string{
		codexAgentMessage("Using spacedock:fo-write-core: «write.classify» internal/** -> override; fix the code directly."),
		codexFileChange(foGuardProductTarget),
	}, "\n")
	if err := assertCodexFOProductEditGuard(broadOverride, targets); err == nil {
		t.Fatal("expected broad override classification to fail for exact product target")
	}

	applyPatch := codexCommand("apply_patch " + foGuardProductTarget)
	if err := assertCodexFOProductEditGuard(applyPatch, targets); err == nil {
		t.Fatal("expected Codex apply_patch against product code before classification to fail")
	}
}

func TestAssertClaudeFOProductEditGuard(t *testing.T) {
	targets := []string{foGuardProductTarget}

	goodRoute := strings.Join([]string{
		claudeText("Using spacedock:fo-write-core: «write.classify» " + foGuardProductTarget + " -> blocked-product; route through worker / explicit override required."),
		claudeToolUse("Agent", `{"prompt":"Dispatch worker to update `+foGuardProductTarget+`."}`),
	}, "\n")
	if err := assertClaudeFOProductEditGuard(goodRoute, targets); err != nil {
		t.Fatalf("blocked product edit routed through a worker must pass: %v", err)
	}

	goodOverride := strings.Join([]string{
		claudeText("Using spacedock:fo-write-core: «write.classify» " + foGuardProductTarget + " -> override; captain explicitly allowed this exact path."),
		claudeToolUse("Edit", `{"file_path":"`+foGuardProductTarget+`"}`),
	}, "\n")
	if err := assertClaudeFOProductEditGuard(goodOverride, targets); err != nil {
		t.Fatalf("exact override before product edit must pass: %v", err)
	}

	noClassification := claudeToolUse("Edit", `{"file_path":"`+foGuardProductTarget+`"}`)
	if err := assertClaudeFOProductEditGuard(noClassification, targets); err == nil {
		t.Fatal("expected Claude Edit against product code before fo-write-core classification to fail")
	}

	blockedThenEdit := strings.Join([]string{
		claudeText("Using spacedock:fo-write-core: «write.classify» " + foGuardProductTarget + " -> blocked-product; route through worker / explicit override required."),
		claudeToolUse("Write", `{"file_path":"`+foGuardProductTarget+`","content":"package status"}`),
	}, "\n")
	if err := assertClaudeFOProductEditGuard(blockedThenEdit, targets); err == nil {
		t.Fatal("expected writing a target after blocked-product classification to fail")
	}
}

func claudeText(text string) string {
	return `{"type":"assistant","message":{"content":[{"type":"text","text":` + mustJSONString(text) + `}]}}`
}
