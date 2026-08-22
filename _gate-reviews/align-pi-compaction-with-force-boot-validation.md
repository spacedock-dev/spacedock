# FO Gate Review — `align-pi-compaction-with-force-boot` validation

## Question
Did the implementation pass validation: does the compaction boundary fire a boot read (not re-inject contract), all 5 ACs satisfied, no scope breach?

## Verdict: APPROVE → done (terminal, merge)

## Validation proof
FO ran the focused test (`go test ./internal/piruntime/... -run TestSpacedockPiExtension -v`) — both `TestSpacedockPiExtensionBootstrapBehavior` and `TestSpacedockPiExtensionChildExemption` PASS. The test exercises the real `.pi/extensions/spacedock.ts` through the Node harness with a mocked `pi.exec` returning a fixed boot JSON.

## ACs (all satisfied)
- AC-1 (value): compaction injects boot record (`"command":"boot"`), 0 bootstraps, no contract pointer — the compacted Pi FO receives state, not a contract re-injection. Baseline that re-injects without reading is replaced.
- AC-2 (mechanism): `pi.exec` called once with `["status","--boot","--identify","--json"]` — the hook fires the boot read.
- AC-3 (scope): `session_start` unchanged (still injects `FO_BOOTSTRAP_TEXT` + "Load the skill").
- AC-4 (child exemption): `PI_SUBAGENT_CHILD=1` → undefined on both paths.
- AC-5 (dedup): existing boot record → undefined.

## Scope
- `go vet ./internal/piruntime/...` clean. `git diff --check` clean. Go test file gofmt-clean.
- Observable-semantics: only compaction-boundary context content changes; no FO bootstrap content, gate/dispatch/state mechanics, command grammar, stored formats, or authority.
- Pre-existing failures isolated: 4 `internal/cli` tests fail on 753's tip too (env-driven: PI_CODING_AGENT + /tmp helper-binaries) — not caused by this change.

## Merge
Code on branch `spacedock-ensign/align-pi-compaction-with-force-boot` (commit `d24fbaada`, off 753 tip `185b53477`). Ready to open as a stacked PR on #753, add to stack #748, and trigger the pi-live lane (the consolidated stack signal).
