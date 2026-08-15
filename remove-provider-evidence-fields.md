---
id: 7c4w88fnmnbtc0tgkrvx0vxj
title: Remove the provider-evidence gate fields
status: ideation
source: "Captain directive, 2026-08-14: value review found zero value; no writer, no retained bytes, no verifier"
started: 2026-08-15T02:55:32Z
completed:
verdict:
score:
worktree:
issue:
sprint: durable-decisions
gates:
    version: 1
    records:
        - id: gate:7c4w88fnmnbtc0tgkrvx0vxj:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:7c4w88fnmnbtc0tgkrvx0vxj-backlog-1
              briefing:
                id: briefing:7c4w88fnmnbtc0tgkrvx0vxj:backlog:attempt-1:revision-1
                digest: sha256:163027efa9408cfc6b01e548d89f063de4a6cbc6257f1022c08a4b536bc2b27f
                request-digest: sha256:6107a996b958133a0620aa02c9a5ee9071804f657c4ae61347c9c4b190ec1ccd
                room-ref: ./remove-provider-evidence-fields/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:7c4w88fnmnbtc0tgkrvx0vxj:backlog:1
                briefing: briefing:7c4w88fnmnbtc0tgkrvx0vxj:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:31.794778Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
---

Remove `ProviderEvidence` from the gate model: the struct, the `provider-evidence` field, the `Validate` branches that police it on open and withdrawn attempts, and their tests. Zero provider-closed attempts exist across 424 recorded attempts, so no stored record carries the field. Re-introduce the fields only together with a real provider integration: writer, retention, and verifier.

## Problem

**The seed's premise above is false. The ideation spike falsified it, and the correction changes the design.** The audit recorded "no writer, no retained bytes, no verifier" and "zero provider-closed attempts exist across 424 recorded attempts, so no stored record carries the field". Measured at HEAD (4d1912a69) over the whole `docs/dev/.spacedock-state` corpus — 924 markdown files, 135 carrying a `gates:` block, 437 attempts across 119 decodable documents — **six attempts in two archived entities do carry `provider-evidence`**, their retained rooms are intact on disk, and the verifier passes against them today. Only the "no writer" claim survives.

What is actually true at HEAD:

- **Writer: gone.** `gate record --room` and the `attempt.ProviderEvidence = &ProviderEvidence{...}` assignment were cut in `4ff999250` (2026-08-07, "gate: cut provider-backed closure from stable v1"). No code path writes the field. This is the one audit claim that holds.
- **Retention: present.** `_archive/fo-boot-install-hint-linux-direct-sandbox.md` (4 attempts) and `_archive/workflow-owned-review-finding-disposition/index.md` (2 attempts) carry the field, and every referenced room still holds its `provider/result.json` and `provider/presented-inventory.json`.
- **Verifier: present and green.** `validateRetainedAuthorityExcept` (`internal/gates/io.go:267-301`) checks both digests, compares the retained provider Resolution against the durable one, then derives and verifies the presentation association. Driven directly against both archived entities it returns no error over all six attempts.
- **Reachability: nil in practice.** The verifier only runs from live gate commands (`SummaryFileAt`, `application.go`, `delivery.go`, `prepare.go`, `operation.go`). Archived entities never reach any of them, so this code has never fired on the only records that would exercise it.

The removal is still right, but for a different reason than the audit gave. The declared contract already excludes this machinery: `docs/specs/gate-resolution-frontmatter-contract.md:266-268` lists "Provider-specific room-backed recording, `gate record --room`, Result or inventory ingestion, retained provider evidence, and provider package selection" under **Explicitly outside v1**. The model contradicts its own contract — roughly 190 lines of verification machinery whose writer was deliberately cut, kept alive only by six frozen archived records that no live command reads.

**The trap the audit's framing would have walked into.** `readDataDiagnostics` decodes with `KnownFields(true)` (`internal/gates/io.go:98`), and unknown keys are tolerated only under `application` mappings. Deleting `Attempt.ProviderEvidence` therefore makes those six records undecodable. Spiked by overlaying the stripped model and rebuilding the binary:

- `spacedock status --workflow-dir docs/dev --validate` goes from `VALID` at rc=0 to **rc=1** with two `Error: invalid gates: decode canonical gates v1: ... field provider-evidence not found in type gates.Attempt` lines naming both archived entities.
- Corpus read-ok drops 119 → 117; read-err rises 16 → 18.

That regression is permanent and no peer entity repairs it. `scope-validate-warnings-to-active-entities` filters warn-tier findings only, states that "the structural error above stays scope-inclusive", and explicitly rules that the read tolerance itself must be kept because "it is load-bearing for every gates read over the legacy corpus".

## Proposed approach

Retire the field behind the read tolerance that already exists, then delete everything the field was holding up.

**1. Drop the retired key before the strict decode.** In `filterApplicationMappings` (`internal/gates/io.go`), which already walks `gates.records[*].attempts[*]`, remove any `provider-evidence` key from the attempt mapping of the validation clone. Silent — the key is *retired*, not *unknown*, so it raises no `Warning` and adds no `status --validate` line. The source bytes are untouched: the filter runs on `cloneYAMLNode`'s copy, and writes still go through the original node.

**2. Delete the model surface.** `Attempt.ProviderEvidence` and `type ProviderEvidence` (`model.go:30,41-44`); the `Validate` open-attempt branch (`:267-269`); the provider half of the withdrawn-attempt branch (`:281-283`, which must keep its application check under the message `withdrawn attempt %s cannot carry application data`); and the provider-closed digest branch (`:289-295`).

**3. Delete the verification chain the field gated.** The `io.go:267-301` block (which also frees the `reflect` import), and then the helpers that become dead with it in `operation.go`: `providerResult`, `presentedInventory`, `resultAssociation`, `decodeProviderResult`, `decodePresentedInventory`, `verifyAssociation`, `deriveAssociation`, `verifyProviderResolution` — 154 lines. Go does not reject unused package-level declarations, so these will compile and silently survive unless deleted on purpose; the strip spike confirmed the tree builds with them left behind.

**Justifying the one new mechanism.** The attempt-level key drop (~11 lines) is the only thing added. It serves AC-2. Two simpler alternatives were considered and both fail:

- *Plain deletion, accept the breakage.* Falsified by measurement above: this repo's own dogfood workflow stops validating, permanently.
- *Hand-migrate the six archived records.* Ruled out by recorded precedent. Commit `6c45fd59c`'s ruling, cited by `scope-validate-warnings-to-active-entities`, is that "already-archived terminal entities must not be hand-migrated to silence a diagnostic about a value the tool itself wrote". The old writer wrote these values, so the ruling applies exactly.

**Documentation diff.** `docs/specs/gate-resolution-frontmatter-contract.md:87-91` documents the tolerance rule and must change with it.

Before:

```
and explicit attempt `state` encodings are rejected. A read tolerates unknown keys
only under each `records[*].attempts[*].application` mapping, reports them as warnings
on explicit `status --validate` or `gate validate`, ignores them for authority, and
never writes them. All other unknown or malformed fields fail closed. There is no
migration or compatibility rewrite.
```

After:

```
and explicit attempt `state` encodings are rejected. A read tolerates unknown keys
only under each `records[*].attempts[*].application` mapping, reports them as warnings
on explicit `status --validate` or `gate validate`, ignores them for authority, and
never writes them. A read also drops the retired
`records[*].attempts[*].provider-evidence` key silently: its writer was cut with
provider-backed closure, frozen archived records still carry it, and it is retired
rather than unknown, so it raises no warning and never reaches the model. All other
unknown or malformed fields fail closed. There is no migration or compatibility
rewrite.
```

## Out of scope

- The retained-authority checks for `request.json`. They have a live verifier. Keep them.
- `presentedItem`, `canonicalPresentationItems`, `boundBriefingManifest`, `validatePresentationGitSources`, and `RawDigest` — all reached by live callers (`prepare.go`, `round.go`, `gitsource`). Keep them.
- Removing `gate validate` (peer entity) and the nine `gate-*` status columns (peer entity).
- Deleting the retained `provider/` room files themselves. The bytes stay; only the durable pin to them goes.

## Coordination

- **`trim-dead-gate-model-surface`** edits `model.go` (`Summary.Condition`/`Eligible`, `:95-96`) and `io.go` (`ReadWithWarnings`, `Diagnostic`) — same two files, disjoint spans. Textual merge risk only, no semantic conflict. Both are at ideation; whichever lands second rebases.
- **`remove-gate-validate-subcommand`** keeps `SummaryFileDiagnosticsAt`, whose body this entity trims. No conflict. It also removes the `gate validate` phrase from the contract sentence this entity rewrites — same sentence, so the second to land must reconcile both edits.
- **`scope-validate-warnings-to-active-entities`** shares that sentence's subject and depends on the read tolerance this entity extends. Its "keep the read tolerance" ruling is the basis for step 1.

## Expected surface and tolerance

Estimate net LOC change: **about -190**, across **8 files**. Tolerance: ±60 lines, ±2 files.

Measured spans: `model.go` -16, `io.go` -36 then +13 for the filter (net -23), `operation.go` -154, `gates_test.go` about -30, plus term-level edits in `internal/cli/gate_test.go` and the two `internal/ensigncycle` tests, a new tolerance test (about +30), and the contract-spec diff.

Observable semantics changed, declared:

- **Stored format:** read tolerance widens by exactly one retired key at attempt level. Nothing else becomes tolerated; every other unknown key still fails closed.
- **Command grammar:** unchanged. No flag, subcommand, or output field is added or removed.
- **Authority:** unchanged for every live path. `provider-evidence` ceases to be an authority input, which it already was not, since no writer produces it.
- **Runtime behavior:** `status --validate` output over `docs/dev` must stay byte-identical to HEAD.

## Acceptance criteria

**AC-1 — The change removes far more lines than it adds.**
Verified by: cumulative line delta of the diff against `origin/main` is negative and at least 130 lines net. Baseline that can move the wrong way: the tolerance filter and its test are additions, so a bloated implementation fails this.

**AC-2 — Stored state stays exactly as readable as it is today.**
Verified by: a corpus decode pass over `docs/dev/.spacedock-state` reports read-ok=119 and read-err=16, matching the HEAD baseline measured in this ideation, with zero attempts decoding a provider-evidence value. A naive strip scores 117/18 and fails.

**AC-3 — This repo's own workflow still validates.**
Verified by: `spacedock status --workflow-dir docs/dev --validate` exits 0 and prints `VALID`, with stdout and stderr byte-identical to the HEAD baseline.

**AC-4 — No production code references the provider-evidence model.**
Verified by: `grep -rn "ProviderEvidence\|providerResult\|resultAssociation\|deriveAssociation\|verifyAssociation\|verifyProviderResolution\|presentedInventory" internal cmd skills` matches nothing outside the single retired-key string literal and its comment in `io.go`, and its test.

**AC-5 — The suite is no worse than the HEAD baseline.**
Verified by: `go test ./... -timeout 30m` and `go test ./... -race -timeout 30m` show no failure that is not already present at HEAD.

Baseline measured during this ideation at `4d1912a69`, so the implementer does not mistake pre-existing breakage for their own: `internal/gates` ok (151s) and `internal/status` ok (266s), but `internal/cli` and `internal/ensigncycle` both **panic at the default 600s `go test` timeout**, and `internal/cli` additionally carries one real pre-existing failure, `TestCodexResolveManifestAgainstInstalledHost`, which depends on an installed Codex host and so is environment-sensitive. Neither is caused by, nor in scope for, this entity. `-timeout 30m` is required for those two packages to finish at all.

## Test plan

Deletion, plus one new test and two existing tests that become the regression fence. Fixture and unit level only; no live workflow run is needed, since no runtime handoff changes.

1. **Retired-key tolerance (new, `internal/gates`).** Feed an entity whose attempt carries `provider-evidence` with both digests through `Read`. Assert the read succeeds, the decoded attempt holds no evidence state, and the file bytes are unchanged after a subsequent canonical gate write. Fails if the key drop is removed, moved after the strict decode, or applied to the write node instead of the validation clone.
2. **`TestPrototypeAndUnknownGateShapesFailClosed`, case `unknown attempt field` (existing, must keep passing).** It already plants a `note: prototype` key at attempt level and requires an unknown-field refusal plus byte-unchanged bytes. This is the fence against implementing step 1 as blanket attempt-level tolerance rather than one named key — the exact over-reach the widening invites. No edit needed; its continued pass is the assertion.
3. **`TestWithdrawalDefinesThirdValidatedFrozenAttemptState`, case `with application` (existing, must keep passing).** Step 2 collapses `if a.ProviderEvidence != nil || a.Application != nil` into an application-only check, and silently dropping the application arm during that edit is the most plausible way to break this change. The `with evidence` case in the same table is deleted with the field; `with application` must still refuse.

Corpus decode and `status --validate` comparisons (AC-2, AC-3) run as one-off evidence in the implementation report against the baselines recorded here, not as committed tests — they depend on this repo's state checkout, which no fresh clone reproduces.

**Spike record.** The riskiest unverified mechanism was whether stored state survives the model strip. It was exercised before design, by overlaying stripped `model.go`/`io.go`, rebuilding `cmd/spacedock`, and running both the corpus decode and `status --validate`. Result: the naive strip breaks two archived entities and fails validation; the tolerant variant reproduces the HEAD baseline byte-for-byte and was likewise built and measured, not reasoned about. Both binaries and both outputs were compared directly.

## Stage Report: ideation

- FAILED: One-off state decode pass proves zero stored records carry the fields before the model strip
  The pass ran and disproved the claim: 6 attempts across 2 archived entities carry `provider-evidence` (924 md files, 135 with gates, 437 attempts, 119 decodable), and their retained `provider/result.json` rooms are intact.
- DONE: Keep-boundary confirmed: retained-authority request.json verification stays
  `validateRetainedAuthorityExcept` keeps its request.json digest/binding block and `validatePresentationGitSources`; only the `ProviderEvidence != nil` arm at io.go:267-301 goes. `presentedItem`, `canonicalPresentationItems`, `boundBriefingManifest`, and `RawDigest` traced to live callers in prepare.go, round.go, and gitsource, and are all named in Out of scope.
- DONE: Value AC measures negative LOC and a green suite over live and archived state
  AC-1 measures net line delta against origin/main (>=130 net removed); AC-2 measures corpus decode read-ok=119/read-err=16 over live and archived state; AC-3 measures `status --validate` byte-identical to the HEAD baseline. "Green suite" was calibrated to "no worse than HEAD baseline" because HEAD is not green - see AC-5.

### Summary

The seed's premise was false and the spike caught it before implementation. The audit claimed no writer, no retained bytes, and no verifier; only the writer claim survives. Six archived attempts carry the field, their rooms are intact, and the verifier passes against them today, so a naive strip is a breaking change: rebuilding `cmd/spacedock` over a stripped model takes `spacedock status --workflow-dir docs/dev --validate` from `VALID`/rc=0 to rc=1 with two `field provider-evidence not found` errors, and drops corpus read-ok 119 to 117.

The design therefore retires the key behind the read tolerance that already exists rather than deleting it outright, which reproduces the HEAD baseline byte-for-byte. Hand-migrating the six archived records was rejected on recorded precedent (6c45fd59c, cited by `scope-validate-warnings-to-active-entities`). Both variants were built and measured, not reasoned about.

Scope grew from the seed's "~3 files" to 8 files and about -190 lines, because deleting the field orphans 154 lines of association-verification helpers in operation.go that Go will happily compile as dead code unless removed deliberately. Coordination overlaps with `trim-dead-gate-model-surface` (same two files, disjoint spans), `remove-gate-validate-subcommand`, and `scope-validate-warnings-to-active-entities` (shared contract sentence) are recorded in the body.
