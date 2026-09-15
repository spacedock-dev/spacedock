# Detached validation audit

Candidate: `5e04442a9c89ef644135721522dbf9bfcc8b7de0`. Detached throwaway checkout: assigned code worktree `.audit-validation`; no candidate edits, commits, pushes, CI, or live reruns.

## Exercised evidence

Cloned retained state.bundle files into `.retained/{before,after,initial-controls-review-required,initial-controls-separate-review-required,initial-controls-cycle-limit,identity-fixed-round-missing,round-final}`. The attached temporary Go test runs inside internal/ensigncycle with `go test ./internal/ensigncycle -run '^TestDetachedSameStage' -count=1 -v`.

Seven bundle regrades pass. Before retains one rejected attempt despite committed correction; after retains one old rejected plus one unresolved, unapplied, unwithdrawn fresh attempt. Both gate-selected corrected plan sources resolve through real Git object bytes and their revision digests. Required-review controls and missing-round stop at one attempt; cycle-limit stops at three. Required-round retains exactly one canonical four-entry room, validated by gates.ValidateRoundFile. All seven corrected plan/frozen-input pairs equal the independent literal `KEEP message A; DELETE message B\n`. Original briefing-1 and frozen-input bytes match their earliest Git commits in every bundle.

Selection matrix passes: canonical artifact and canonical Reference accepted; wrong source type, wrong revision digest, wrong exact identity, zero-width Unicode suffix, empty source and truncated JSON rejected. Worker-event matrix rejects empty, completion-before-spawn, different completion owner, duplicate completion and duplicate spawn. These tests cannot substitute for raw observed native events.

Canonical round audit command: `go test ./internal/gates -run '^TestRound(RecordNeutralReplayAndRefusalsAreByteClean|NoFindingsAndPreflightRefusals)$' -count=1 -v`. Pass: missing/incomplete/dangling/hold evidence refuses without side effects; complete Roborev round preserves all five entries and lifecycle bytes; neutral recorder does not invent projection. No changed hot-path resource policy or production binary path was introduced, so no scaling probe is indicated.

The first temporary selection positive failed because the audit set Main and State to the same checkout, correctly producing a main URI. The corrected fixture declares distinct main/state Git roots; all positives and negatives then pass. Both logs are retained; this was an audit setup error, not a candidate finding.

## Authority and scope

Skill applicability is the only shipped behavior changed; no CLI/state/recorder schema changed. Declared checkpoint authority, concrete captain correction authority in workflows lacking that checkpoint, independent-review identity, separate round/projection obligations, and cycle-3 stop remain explicit. Existing CLI success cannot establish compliance with those agent-owned obligations. Retained before/after states support the behavioral outcome, but a raw-native retention gap blocks the complete evidence claim (V1 below).

Latest captain approval covers five files/+236 net and local targeted proof, not release/runtime parity. Final normal/race logs each show only TestCodexResolveManifestAgainstInstalledHost failing because spacedock@spacedock is absent but spacedock-local resolves. Prior FO decline is preserved; neither suite is wholly green. No new timeout/race failure was observed in those final logs.

## V1: native lifecycle retention gap

The saved public codex-exec.jsonl has no native spawn/completion record; after contains only four empty wait collab results. TSV rows are a derived observation, not replayable source events. Producer independently confirms every isolated CODEX_HOME parent rollout was deleted by t.Cleanup(os.RemoveAll(dir)); original paths are empty. Raw call-ID pairing, returned task-name ownership and author-attributed Done completion therefore cannot be independently replayed. Durable reports corroborate work but cannot prove native identity.

This is an evidence defect, not an observed workflow outcome failure. Minimal repair belongs to existing nativeLifecycleStream in claude_live_runner_test.go: retain the already-read correlated public+native stream in result.artifactDir before cleanup, and test exact retention/replay. No new controller or lifecycle framework is needed. One retention line plus approximately 15–25 focused test lines exceeds the exact +236 cap and requires captain scope approval. Six candidate controls require raw replacement under registry policy; a causal baseline adds one run, and the conventional independent-review proof adds one if its attribution is to be replayable (eight bounded runs total). No candidate repair or rerun was authorized or performed by this validator.
