# Validation evidence: advisory review round recorder

Candidate: `1ae990f509b62240fc16a3069239959e8a9fa8dc`

Verdict: PASSED

Checklist: 12 done, 0 skipped, 0 failed.

## Capability exercised

The validated operation records one already-completed reviewer round. It writes the exact reviewed Briefing and ordered advisory Resolution graph, records every material fix or explicit decline, projects one Feedback Cycles entry, and leaves gate/application/lifecycle state unchanged.

## Independent falsifiers

- The successful-write oracle matches all 36 expected entity lines. Corrupting unrelated frontmatter and body bytes makes the focused test fail.
- The shared rejection-flow grader requires an observed branch-built recorder invocation, canonical room, pointer, advisory graph, projection, retained inputs, and unchanged lifecycle sentinels. Setting `invoked=false` fails.
- Malformed log/digest, occupied room, lock, CAS, replacement failure, immutable divergence, exact replay, and whole-tree rollback controls remain byte-clean.
- Removing the distinction between no findings and an all-declines triage Resolution fails the public-command fixture.

## Executed checks

- Focused recorder and rejection-flow suites: PASS.
- Codex `rejection-flow` live subtest with pinned branch-built `SPACEDOCK_BIN`: PASS in 360.93 seconds.
- `go test ./...`: PASS.
- `go test ./... -race`: PASS.
- Fresh affected-package race run: PASS.
- Live-tag compilation, `gofmt`, and `git diff --check`: PASS.

Claude's attempted live lane stopped before first-officer work on HTTP 401. It is classified as an external credential condition, not green evidence and not a product failure. Release use of that adapter still requires its existing credentialed CI lane.

## Scope

Commit `1ae990f5` changes seven test/testdata paths and zero production paths. The underlying recorder remains 795 additions and 156 deletions, or 639 net production LOC. No new package, schema, gate application, or workflow-state mechanism was added.

## Deferred risks

Exact-log digests, historical fixed-evidence reads, power-loss fsync, CRLF handling, and loose-heading hardening remain outside the approved promise. Each retains its recorded promotion trigger; none is material to the five accepted ACs.
