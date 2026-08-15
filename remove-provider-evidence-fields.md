---
id: 7c4w88fnmnbtc0tgkrvx0vxj
title: Remove the provider-evidence gate fields
status: validation
source: "Captain directive, 2026-08-14: value review found zero value; no writer, no retained bytes, no verifier"
started: 2026-08-15T02:55:32Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-remove-provider-evidence-fields
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
        - id: gate:7c4w88fnmnbtc0tgkrvx0vxj:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:7c4w88fnmnbtc0tgkrvx0vxj-ideation-1
              briefing:
                id: briefing:7c4w88fnmnbtc0tgkrvx0vxj:ideation:attempt-1:revision-1
                digest: sha256:b496766d84d0bb0205e020157e7e84855669c88b3e411fac57b643f7ab638b90
                request-digest: sha256:371679a1811510ae4dcd770a6ac5696aefe63569641263ed7a19cba5bea2ff50
                room-ref: ./remove-provider-evidence-fields/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-15T03:38:03.431911Z"
                reason: 'Entity corrected post-prepare (b4bff335b): collision-contaminated baselines re-measured in isolation; re-preparing against corrected ACs'
            - id: gate-attempt:7c4w88fnmnbtc0tgkrvx0vxj-ideation-2
              briefing:
                id: briefing:7c4w88fnmnbtc0tgkrvx0vxj:ideation:attempt-2:revision-1
                digest: sha256:c7e23dbc596153c28aa115cd94a7975a73f68edd98babb10961599d599cd12b4
                request-digest: sha256:f1b927d79501223f17d84d48a1512d2b4d41de7777cad1f97e19f572245c1f36
                room-ref: ./remove-provider-evidence-fields/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:7c4w88fnmnbtc0tgkrvx0vxj:ideation:2
                briefing: briefing:7c4w88fnmnbtc0tgkrvx0vxj:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-15T03:55:39.311199Z"
                decision: approve
                reason: 'Captain ruling 2026-08-15 (approve all except x8): approved into implementation'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:7c4w88fnmnbtc0tgkrvx0vxj:validation
          stage: validation
          attempts:
            - id: gate-attempt:7c4w88fnmnbtc0tgkrvx0vxj-validation-1
              briefing:
                id: briefing:7c4w88fnmnbtc0tgkrvx0vxj:validation:attempt-1:revision-1
                digest: sha256:11b054ad5582a543e76eb10760a9697e0fb7c8720463dc985b3e1ffd3fa79814
                request-digest: sha256:89678880ad3050c7171f563ef9ea48f1dd3626d8ff0ca15652bb2e72ee2d860a
                room-ref: ./remove-provider-evidence-fields/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:7c4w88fnmnbtc0tgkrvx0vxj:validation:1
                briefing: briefing:7c4w88fnmnbtc0tgkrvx0vxj:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T16:08:32.933999Z"
                decision: approve
                reason: 'Captain batch approval 2026-08-15: validation PASSED; land via releng-27 train'
              application:
                target-stage: done
                state: pending
pr: "#702"
---

Remove `ProviderEvidence` from the gate model: the struct, the `provider-evidence` field, the `Validate` branches that police it on open and withdrawn attempts, and their tests. Zero provider-closed attempts exist across 424 recorded attempts, so no stored record carries the field. Re-introduce the fields only together with a real provider integration: writer, retention, and verifier.

## Problem

**The seed's premise above is false. The ideation spike falsified it, and the correction changes the design.** The audit recorded "no writer, no retained bytes, no verifier" and "zero provider-closed attempts exist across 424 recorded attempts, so no stored record carries the field". Measured over the whole `docs/dev/.spacedock-state` corpus at code `ef8f55c83` and state `86e7eebc2` — 925 markdown files, 135 carrying a `gates:` block, 443 attempts across 119 decodable documents — **six attempts in two archived entities do carry `provider-evidence`**, their retained rooms are intact on disk, and the verifier passes against them today. Only the "no writer" claim survives.

All measurements below were produced in an isolated, slug-named spike directory and re-run end to end after a shared-scratchpad path collision was reported; the corpus counts drift by a few files as peers commit, so both SHAs are pinned. `internal/gates`, `internal/status`, and `internal/cli` are byte-identical between `4d1912a69` and `ef8f55c83`.

What is actually true at HEAD:

- **Writer: gone.** `gate record --room` and the `attempt.ProviderEvidence = &ProviderEvidence{...}` assignment were cut in `4ff999250` (2026-08-07, "gate: cut provider-backed closure from stable v1"). No code path writes the field. This is the one audit claim that holds.
- **Retention: present.** `_archive/fo-boot-install-hint-linux-direct-sandbox.md` (4 attempts) and `_archive/workflow-owned-review-finding-disposition/index.md` (2 attempts) carry the field, and every referenced room still holds its `provider/result.json` and `provider/presented-inventory.json`.
- **Verifier: present and green.** `validateRetainedAuthorityExcept` (`internal/gates/io.go:267-301`) checks both digests, compares the retained provider Resolution against the durable one, then derives and verifies the presentation association. Driven directly against both archived entities it returns no error over all six attempts.
- **Reachability: nil in practice.** The verifier only runs from live gate commands (`SummaryFileAt`, `application.go`, `delivery.go`, `prepare.go`, `operation.go`). Archived entities never reach any of them, so this code has never fired on the only records that would exercise it.

The removal is still right, but for a different reason than the audit gave. The declared contract already excludes this machinery: `docs/specs/gate-resolution-frontmatter-contract.md:266-268` lists "Provider-specific room-backed recording, `gate record --room`, Result or inventory ingestion, retained provider evidence, and provider package selection" under **Explicitly outside v1**. The model contradicts its own contract — roughly 190 lines of verification machinery whose writer was deliberately cut, kept alive only by six frozen archived records that no live command reads.

**The trap the audit's framing would have walked into.** `readDataDiagnostics` decodes with `KnownFields(true)` (`internal/gates/io.go:98`), and unknown keys are tolerated only under `application` mappings. Deleting `Attempt.ProviderEvidence` therefore makes those six records undecodable. Spiked by overlaying the stripped model and rebuilding the binary:

- `spacedock status --workflow-dir docs/dev --validate` goes from `VALID` at rc=0 to **rc=1** with two `Error: invalid gates: decode canonical gates v1: ... field provider-evidence not found in type gates.Attempt` lines naming both archived entities.
- Corpus read-ok drops 119 → 117; read-err rises 16 → 18.
- Less obvious, and the reason a "does it still say VALID" check is not sufficient: the strip also **loses 14 warning lines**, 125 → 111. The two entities now fail the structural read before their per-attempt application-field diagnostics are emitted, so the naive strip silently reduces reporting on exactly the records it broke.

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
Verified by: a corpus decode pass over `docs/dev/.spacedock-state` reports the same read-ok and read-err counts as a HEAD binary run against the same state SHA in the same session, with zero attempts decoding a provider-evidence value. Measured at state `86e7eebc2`: HEAD and the tolerant build both score read-ok=119 / read-err=16; a naive strip scores 117/18 and fails. Re-measure the HEAD side rather than hard-coding 119/16 — the corpus grows as peers commit.

**AC-3 — This repo's own workflow still validates, with undiminished reporting.**
Verified by: `spacedock status --workflow-dir docs/dev --validate` exits 0, prints `VALID`, and is byte-identical on **both** stdout and stderr to a HEAD binary run against the same state SHA. At `86e7eebc2` that baseline is `VALID` at rc=0 with 125 stderr `Warning:` lines. Comparing stdout alone is not enough: the naive strip fails this at rc=1 with 2 `Error:` lines and only 111 warnings, and a variant that restored rc=0 while still dropping warning lines would pass a stdout-only check.

**AC-4 — No production code references the provider-evidence model.**
Verified by: `grep -rn "ProviderEvidence\|providerResult\|resultAssociation\|deriveAssociation\|verifyAssociation\|verifyProviderResolution\|presentedInventory" internal cmd skills` matches nothing outside the single retired-key string literal and its comment in `io.go`, and its test.

**AC-5 — The suite stays green, with one declared environment carve-out.**
Verified by: `go test ./... -timeout 30m` and `go test ./... -race -timeout 30m` pass, except for failures also present on an unmodified HEAD run in the same environment.

Baseline measured during this ideation at `ef8f55c83`, isolated, so the implementer does not mistake environment noise for their own breakage. `internal/gates` ok (68s), `internal/status` ok (128s), `internal/ensigncycle` ok (384s), `internal/cli` **one failure**: `TestCodexResolveManifestAgainstInstalledHost`, which shells out to the installed Codex host and failed here only because the sandbox denied reading `~/.codex/config.toml` (`Operation not permitted`). That is a property of this execution environment, not the repo; expect it to pass where Codex config is readable. Nothing else fails.

Two practical notes. `-timeout 30m` is needed: `internal/cli` and `internal/ensigncycle` each exceed the 600s `go test` default when the machine is shared with other running agents — an earlier contended run had both panic on timeout, while the same packages finished in 314s and 384s uncontended. And run the HEAD comparison in the same conditions as the change run, or machine load alone will move the result.

## Test plan

Deletion, plus one new test and two existing tests that become the regression fence. Fixture and unit level only; no live workflow run is needed, since no runtime handoff changes.

1. **Retired-key tolerance (new, `internal/gates`).** Feed an entity whose attempt carries `provider-evidence` with both digests through `Read`. Assert the read succeeds, the decoded attempt holds no evidence state, and the file bytes are unchanged after a subsequent canonical gate write. Fails if the key drop is removed, moved after the strict decode, or applied to the write node instead of the validation clone.
2. **`TestPrototypeAndUnknownGateShapesFailClosed`, case `unknown attempt field` (existing, must keep passing).** It already plants a `note: prototype` key at attempt level and requires an unknown-field refusal plus byte-unchanged bytes. This is the fence against implementing step 1 as blanket attempt-level tolerance rather than one named key — the exact over-reach the widening invites. No edit needed; its continued pass is the assertion.
3. **`TestWithdrawalDefinesThirdValidatedFrozenAttemptState`, case `with application` (existing, must keep passing).** Step 2 collapses `if a.ProviderEvidence != nil || a.Application != nil` into an application-only check, and silently dropping the application arm during that edit is the most plausible way to break this change. The `with evidence` case in the same table is deleted with the field; `with application` must still refuse.

Corpus decode and `status --validate` comparisons (AC-2, AC-3) run as one-off evidence in the implementation report against the baselines recorded here, not as committed tests — they depend on this repo's state checkout, which no fresh clone reproduces.

**Spike record.** The riskiest unverified mechanism was whether stored state survives the model strip. It was exercised before design, by overlaying stripped `model.go`/`io.go`, rebuilding `cmd/spacedock`, and running both the corpus decode and `status --validate`. Result: the naive strip breaks two archived entities and fails validation; the tolerant variant reproduces the HEAD baseline byte-for-byte and was likewise built and measured, not reasoned about. Three binaries (HEAD, strip, tolerant) were built and compared directly.

The spike used `go test -overlay` and `go build -overlay` throughout, so the shared working tree was never modified — relevant because peers build from the same checkout. It was then re-run end to end in a slug-named directory after the FO reported a shared-scratchpad collision. That re-run mattered: the first pass recorded a HEAD stderr baseline of zero lines, which the isolated re-run showed to be 125 warning lines, matching the independently documented count in `scope-validate-warnings-to-active-entities`. The conclusion held, but AC-3's baseline was wrong and is now corrected. Reproduce with the commands in this section against a freshly built HEAD binary; do not trust a stale one.

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

## Stage Report: ideation (cycle 2 - isolated re-verification)

Re-run of every cited measurement in a slug-named directory after the FO's shared-scratchpad collision advisory. The design conclusion is unchanged; two recorded baselines were wrong and are corrected in the body.

- DONE: One-off state decode pass proves zero stored records carry the fields before the model strip
  Re-run at code `ef8f55c83` / state `86e7eebc2`: 925 md files, 135 with gates, 443 attempts. HEAD decodes 6 provider-evidence attempts; tolerant build decodes 0 at identical read-ok=119 / read-err=16; naive strip scores 117/18. Still FAILED as originally worded - the premise it asked me to prove is false - but the pass itself is now collision-free.
- DONE: Keep-boundary confirmed: retained-authority request.json verification stays
  Unchanged; re-confirmed by the verifier run passing on both archived entities (4 and 2 attempts) with the request.json path intact.
- DONE: Value AC measures negative LOC and a green suite over live and archived state
  AC-1/AC-2/AC-3 re-baselined against freshly built HEAD, strip, and tolerant binaries. AC-5 rewritten: the suite is effectively green, not broken.

### Corrections

1. **AC-3 baseline was wrong.** The contaminated run recorded HEAD stderr as 0 lines; isolated re-run shows **125 `Warning:` lines**, which matches the count `scope-validate-warnings-to-active-entities` documents independently. AC-3 now requires byte-identity on both streams against a freshly built HEAD binary.
2. **AC-5 baseline was wrong.** The contaminated run showed `internal/cli` and `internal/ensigncycle` panicking at the 600s default. Uncontended, `ensigncycle` passes in 384s and the only failure is `TestCodexResolveManifestAgainstInstalledHost`, caused by this sandbox denying `~/.codex/config.toml`. The suite is green modulo that environment carve-out.
3. **New finding.** The naive strip also drops 14 warning lines (125 to 111), because the two broken entities fail the structural read before their per-attempt diagnostics emit. A stdout-only "still says VALID" check would miss this, so AC-3 compares stderr too.

### Summary

Everything load-bearing reproduced: 6 archived attempts carry the field, their rooms are intact, the verifier passes against them, the naive strip takes `status --validate` to rc=1, and the tolerant variant is byte-identical to HEAD on both streams. The collision did not change any conclusion, but it did corrupt two baselines that later stages would have compared against, and both were numbers I would have had no way to catch without re-running. Spike artifacts are under `scratchpad/spike-remove-provider-evidence-fields/`; the shared working tree was never modified, since all variants ran through `go build -overlay` / `go test -overlay`.

## Stage Report: implementation

- DONE: Retire-behind-tolerance exactly as gated: one named retired key, the prototype fail-closed test untouched green
  `retiredAttemptKey = "provider-evidence"` is the only key `filterApplicationMappings` drops, and it drops it from the validation clone only; `TestPrototypeAndUnknownGateShapesFailClosed` is unedited and green. Code commit `5b7bb8cfd` on `spacedock-ensign/remove-provider-evidence-fields`.
- DONE: Corpus decode 119/16 and status --validate byte-identical on both streams against the recorded baselines
  At state `296876f46`, re-checked unchanged across all four runs: HEAD and changed builds both score read-ok=119 / read-err=16 / attempts=452 / warnings=121; HEAD decodes 6 `provider-evidence` attempts, the change decodes 0, and the two archived files still carry their 6 keys on disk. `status --workflow-dir docs/dev --validate` is rc=0 on both with `cmp`-identical stdout and stderr (`VALID`; 125 `Warning:`, 0 `Error:`).
- DONE: Within -190 plus or minus 60 over 8 files; -timeout 30m; never chase the known environmental failure
  `git diff --numstat` against merge-base `4d1912a69`: 8 files, +108 / -253, net **-145**, inside both declared tolerances. `go test ./... -timeout 30m` and `go test ./... -race -timeout 30m` are green in every package on both runs except `internal/cli`'s `TestCodexResolveManifestAgainstInstalledHost`, the declared sandbox carve-out (`~/.codex/config.toml`: Operation not permitted); not chased. `-timeout 30m` was load-bearing under contention: `internal/cli` took 1019s / 798s and `internal/ensigncycle` 1186s / 871s, all past the 600s default.

### Tests and what falsifies them

- New `TestRetiredProviderEvidenceKeyReadsSilentlyWithoutWideningAttemptTolerance` (`internal/gates`) asserts that a stored `provider-evidence` attempt reads with zero warnings, never reaches the model, remains present on the returned compare-and-swap node, passes the frozen-attempt transition check, and that an unknown sibling attempt key still refuses. Two mutation runs prove it falsifiable: deleting the key-drop loop fails it with `field provider-evidence not found`; changing `cloneYAMLNode(gatesNode)` to `gatesNode` fails it on the write-node assertion.
- `TestPrototypeAndUnknownGateShapesFailClosed` is unedited. Its `unknown attempt field` case fails if the drop is ever widened from one named key to blanket attempt-level tolerance.
- `TestWithdrawalDefinesThirdValidatedFrozenAttemptState` keeps `with application`, which fails if the withdrawn branch loses its application arm while the provider arm is removed. The `with evidence` case was deleted with the field.

### Deviations and notes

1. **Test-plan wording on byte preservation.** The plan asked that entity bytes be unchanged "after a subsequent canonical gate write". The canonical writer re-marshals the gates block from the model, so it has never round-tripped any tolerated-but-unmodeled key and it also renormalizes indentation; the literal assertion was not achievable before this change either. The test asserts what the writer actually promises: the read leaves bytes untouched, the write is accepted rather than CAS-refused, and every byte outside the gates mapping survives. No live command writes archived entities, so this is unreachable in practice.
2. **One adjacency the design did not name.** `ValidateTransition` decodes the stored node without `KnownFields`, so the retired key is ignored symmetrically on the old-node and new-model sides. The frozen-attempt check still accepts an unchanged attempt and still refuses a mutated one. Proved by probe, then folded into the new test.
3. **Left alone.** `docs/roadmap/durable-decisions/staff-review-sprint-close.md` still names the deleted helpers. It is a historical review record rather than a contract, and is outside the declared surface.
4. **Coordination.** Branched from `4d1912a69`. The contract sentence rewritten here is shared with `remove-gate-validate-subcommand`; whichever lands second reconciles both edits.

### Summary

The retire-behind-tolerance design landed as approved and the recorded baselines reproduced exactly: 119/16 corpus decode and byte-identical `status --validate` on both streams against a freshly built HEAD binary at one verified-stable state SHA. Net -145 lines across 8 files, inside tolerance, with the 154 orphaned association helpers in `operation.go` carrying most of the deletion.

Two things beyond the plan were worth the time. The new test was mutation-checked in both directions the design warned about, so its pass is evidence rather than assertion. And `ValidateTransition` turned out to read the raw stored node, which would have been a live regression path had the retired key decoded asymmetrically; it does not, and that is now fenced. The only red in the suite is the Codex config sandbox denial the ideation already classified as environmental.

## Stage Report: validation

- DONE: Independently re-exercise the ACs against worktree commit 5b7bb8cfd, never by reading the report: corpus decode 119/16 with zero provider-evidence values decoded, status --validate byte-identical to the recorded HEAD baseline on BOTH stdout and stderr (125 warning lines), grep clean per AC-4, net delta at least -130
  All re-measured from scratch: frozen copy of docs/dev, freshly built HEAD (4d1912a69) and changed (5b7bb8cfd) binaries, own corpus-decode program; per-AC evidence below.
- DONE: The two guard tests are the fence: TestPrototypeAndUnknownGateShapesFailClosed unknown-attempt-field case must still refuse byte-clean (proves tolerance is one named key, not blanket), and the withdrawal with-application case must still refuse
  Both pass at 5b7bb8cfd (refusal + unchanged-bytes assertions confirmed in source), and both were falsified by overlay mutation: a blanket attempt-tolerance mutant fails unknown_attempt_field plus the new test's sibling-key case; deleting the withdrawn application arm fails with_application.
- DONE: Suite per AC-5 no-worse-than-baseline with -timeout 30m; reproduce any failure on clean 4d1912a69 before attributing; verdict PASSED or REJECTED with per-AC citations
  `go test ./... -timeout 30m` and `-race -timeout 30m` fail only TestCodexResolveManifestAgainstInstalledHost (`~/.codex/config.toml`: Operation not permitted), reproduced identically on a clean 4d1912a69 checkout. Verdict: **PASSED**.

### Per-AC evidence

- AC-1 PASS: `git diff --numstat 4d1912a69..5b7bb8cfd` = 8 files, +119/-253, net **-134** (floor is -130). The implementation report's "+108/-253, net -145" does not reproduce; the corrected figure still passes.
- AC-2 PASS: frozen-snapshot corpus decode (925 md files, 135 with gates): HEAD build read-ok=119 / read-err=16 / attempts=455 with 6 provider-evidence attempts decoded (the two `_archive` entities); changed build 119/16/455 with **0** decoded; error-file lists identical.
- AC-3 PASS: both binaries rc=0 printing `VALID` on the same snapshot; stdout AND stderr cmp-identical; HEAD stderr is exactly 125 `Warning:` lines, 0 `Error:` lines.
- AC-4 PASS: the AC grep over `internal cmd skills` matches only the new test's name in `gates_test.go`. A repo-wide sweep additionally finds only `docs/roadmap/durable-decisions/staff-review-sprint-close.md`, the declared historical review record.
- AC-5 PASS: both suites green except the environmental codex sandbox failure, attributed by clean-baseline reproduction, not by trusting the report.

### Adversarial pass

Four overlay mutants each fail exactly their fencing assertion: blanket attempt tolerance; deleted key-drop (`field provider-evidence not found`); filtering the write node instead of the validation clone (CAS-node assertion); dropped withdrawn-application arm. Level scoping probed empirically on the changed build: `provider-evidence` at record level and briefing level still refuses; only the attempt-level key is dropped. The contract-spec wording matches the approved ideation diff verbatim, and the keep-boundary held: request.json digest checks and `validatePresentationGitSources` are intact in `validateRetainedAuthorityExcept`.

### Findings

Material: none. Deferred risks: none new — the shared contract sentence with `remove-gate-validate-subcommand` stays a land-second reconciliation, already recorded in Coordination. Polish: the implementation report's AC-1 arithmetic (+108/-253, net -145) is not reproducible; actual is +119/-253, net -134. Evidence-report inaccuracy only; no AC fails under the corrected number and no candidate change is needed.

### Summary

Every AC was re-exercised from scratch rather than replayed: fresh HEAD and changed binaries, a frozen snapshot of docs/dev (the live state checkout moved twice during setup, confirming the pin was necessary), and an independent reflection-based corpus decoder. All five ACs hold, all four mutation directions fail their fences, tolerance is provably one named key at one level, and the only suite failure is environmental and present on clean 4d1912a69. Recommend PASSED with one polish-level reporting correction and no material or deferred findings.
