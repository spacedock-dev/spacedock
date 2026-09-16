# Independent validation audit

Candidate: `1b3b70508`; base: `4ce49f1ea`. Detached checkout: `/tmp/semantic-validation-audit`.

Reproduce: detach a throwaway worktree at the candidate, copy `validation_audit_test.go.txt` into its `internal/dispatch/validation_audit_test.go`, and run `go test ./internal/dispatch -run '^TestValidationAudit|^TestSemantic' -count=1 -v`.

`checks.log` records the clean run. `negative-controls.log` records two independently applied, then restored mutations to `names.go`: prepend `spacedock-ensign-` in the short-name return; accept `len(matches) >= 1` in `uniqueWorker`. Both exit 1. No candidate changes.

Native source: `/tmp/semantic-native-evidence-v2/codex-shared-scenarios/semantic-names-codex/native-lifecycle.jsonl`. Independently parsed response_item records show exactly one agents.spawn_agent call (`call_o8vcAbCEqQHkNNYLPWP9jiBy`) with task_name `ci_duration_hints_ideation`, matched output `/root/ci_duration_hints_ideation`, then one agents.followup_task (`call_kB5YXvoKtlXbtDRozJYHjEu4`) to that exact handle. An additional send_message targets the same handle; it does not spawn or advance another worker. Message bodies are encrypted in the native rollout; this review does not assert independent byte-for-byte prompt verification. Parent command_execution records contain two instruction reads, no entity writes. Retained entity.md contains the unchanged seed and distinct marker reports; state-log.txt records only ci-duration-hints.md in commits 5c9057416c344afc9eb904d01ea5554d7342303a and 21354d633b08a59979d36dce5c891e9e062f5ffd. The owned live test establishes clean state and unchanged seed; validation reuses its 109.67-second result.

Broad results reused: `/tmp/semantic-full-normal.log` and `/tmp/semantic-full-race.log`. Both fail only TestCodexResolveManifestAgainstInstalledHost at internal/cli/codex_resolve_test.go:44, under the existing FO-authorized DECLINE. All other package results pass; no race diagnostic. This is not a wholly green suite or CI claim.
