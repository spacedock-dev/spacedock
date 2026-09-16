# Retained validator: CI archive identity correction

Recommendation: **PASSED** for `c8c6d1125..6a9e4c839`. This focused recheck supplements the approved validation report without changing its bytes, acceptance criteria, or gate state. No new material findings or deferred risks were identified.

- DONE: Independently verify the fixture-only identity correction and trace-confirmed red/green archive evidence within approved scope.
  The exact diff changes only `skills/integration/launcher_smoke_test.go`, +5/-1: use existing `testgit.InitRepo` for the state checkout, then preserve the same add and seed commit. The helper persists repository-local user.name/user.email, whereas the old helper's `git -c` identity applied only to its seed command. Production behavior and authority are unchanged. Aggregate diff above `70165e9b5`: 22 files, +635/-44, +591 net, within approved 22-file/+800-net limits; `git diff --check` passed.
- DONE: Retain a durable focused recheck with exact final normal/race results and no duplicated broad/native runs or gate changes.
  The producer's [correction report](ci-archive-identity/report.md), exact red/green traces, and final logs were independently inspected. Every retained log was byte-compared against its `/tmp/durability-ci-identity-investigation` source. No tests or broad/native suites were rerun; no candidate, old report, frontmatter, gate, PR or remote was changed.

## Exact behavioral evidence

The producer ran the existing real-launcher `TestLauncherListSetArchive` with isolated global configuration containing `user.useConfigOnly=true`, system configuration disabled, and `GIT_AUTHOR_*`, `GIT_COMMITTER_*`, `GIT_CONFIG_*` and `EMAIL` overrides scrubbed before assigning the isolated config. The report retains the reproduction command. Clearing ambient identity alone was insufficient on this host because Git could auto-detect it; disabling auto-detection established the failing boundary.

I parsed both Trace2 files as JSON records and joined each archive command to its exit event by exact session ID. Each trace has exactly one `git commit -q -m "archive pilot-entity (retirement)" -- :(literal)pilot-entity :(literal)_archive/pilot-entity` start and one corresponding exit:

| Evidence | Exact process and result |
| --- | --- |
| [Red trace](ci-archive-identity/required-trace.jsonl) | SID `20260916T044737.499932Z-H3e8ac57e-P00015bed`; error `no email was given and auto-detection is disabled`; exit **128**. The [test log](ci-archive-identity/identity-required.log) reports archive exit 1 and fails. |
| [Green trace](ci-archive-identity/green-trace.jsonl) | SID `20260916T044947.426222Z-H3e8ac57e-P00009bd5`; same archive command exits **0** after repo-local identity setup. The [test log](ci-archive-identity/identity-required-green.log) passes in **1.375s**. The later missing-origin diagnostic is the expected local-only publication path. |

This establishes AC-2 fixture validity through the real retirement commit, not merely a seed-commit success or an error-string match. Removing persistent repository identity restores the exact red boundary. A further detached mutation probe would duplicate the existing trace-confirmed red/green comparison and was unnecessary for this one-file correction.

## Reused final required checks

- [Integration](ci-archive-identity/integration.log): `go test ./skills/integration -count=1`, passed **3.490s**.
- [Normal](ci-archive-identity/normal.log): `go test ./...`, exit **1**, solely `TestCodexResolveManifestAgainstInstalledHost` at `codex_resolve_test.go:44`; CLI package **218.292s**, integration passed **5.805s**, every other package passed.
- [Race](ci-archive-identity/race.log): `go test ./... -race`, exit **1**, solely the same assertion; CLI package **199.290s**, integration passed **13.086s**, every other package passed. No data-race report appears.
- Exact retained limitation: host discovery says `spacedock@spacedock` is absent, while resolution returns `/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`. FO declined this existing environment-dependent resolver correction as out of scope. These are not wholly green normal/race runs.

The normal and race log SHA-256 values independently read here are `ef68215780abd1646683e3bb1eaf46ce8b8ce519fbc797ffa11807960b9260a7` and `0564ab1e57910958c47be18bd6c955ab431f36ab4d8694b38454c2ce2afd5239`. Red and green trace SHA-256 values are `d0417b70ebed3d6639c450a39e01d16912ef011e63105b84d6cca1188816cbf2` and `02a6d33414d3b988e2fda0718bfbc76ea0e2b8281466540d13f7d4fde911c1af`.

### Summary

The authorized fixture correction fixes the demonstrated missing-identity failure without widening the approved durability scope. The existing isolated red/green execution and exact Git process exits support PASSED for the correction; the known resolver limitation remains explicit. Only this focused artifact was written and committed by the validator; no push was performed.
