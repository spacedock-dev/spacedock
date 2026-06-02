---
id: 02a92vbcn4x7453bhszprpje
title: Retire the Python oracle — remove the embedded vendor/status + VendorRunner, graduate the differential-parity suite
status: ideation
source: FO investigation (2026-06-01) — the sole open precondition gating zj (yaml-parser-migration), currently orphaned (0x scoped it out; sprint-notes deferred it). Post-bootstrap.
started: 2026-06-02T15:29:46Z
completed:
verdict:
score: "0.25"
worktree:
issue:
---

Retire the vendored Python status oracle so the codebase no longer carries a ~100KB embedded
interpreter-dependent reference. RUNTIME python-freedom is already done (verified on this checkout:
the only production runner wired in `internal/cli/cli.go:44` is `NativeRunner`; `VendorRunner` and
the oracle are referenced ONLY from `_test.go` files). What remains is **test-time**: the oracle
(`internal/status/vendor/status`, 106164 bytes) is still `//go:embed`-ed (`vendor_runner.go:13`) and
run by `VendorRunner` (`internal/status/vendor_runner.go`) as the differential-parity reference —
native Go output is byte-asserted against the live oracle across the parity suites in
`internal/status` + `internal/dispatch`.

This is the **sole open precondition** gating `zj` (yaml-parser-migration): byte-parity with the
Python oracle is the only thing forcing the hand-rolled frontmatter parser; once the oracle is
retired, `zj` can adopt a YAML library with documented divergences.

## Touchpoint inventory (verified on this checkout, 2026-06-02)

The original spec named the status oracle only. Two grep sweeps (`grep -rln runOracle\|indRunOracle\|vendoredOracle internal/`, `grep -rln commission/bin internal/ skills/`) found **three** vendored Python artifacts, not one:

1. **`internal/status/vendor/status`** (106164 bytes) — the `//go:embed`-ed status oracle. Driven by `runOracle`/`indRunOracle`/`runLauncher` in the status parity suite.
2. **`skills/commission/bin/status`** (106164 bytes — **byte-identical** to #1; `diff -q` confirms). A second copy of the same script. No runtime invokes it (skills now say `spacedock status`); it survives only as a prose-neutrality scan target.
3. **`skills/commission/bin/claude-team`** (51361 bytes) — the **dispatch** oracle. Driven by `runOracle`/`vendoredOracle` in `internal/dispatch/*_parity_test.go`; `parity_harness_test.go:36` `t.Fatalf` if absent. Also the `context-budget`/`standing` reference.

Plus `skills/commission/bin/__pycache__/{status,claude-team}cpython-314.pyc` (compiled-bytecode debris that travels with #2 and #3).

### Test-side touchpoints to remove or rewrite

**Status seam (delete after the suite is graduated):**
- `internal/status/vendor_runner.go` — `VendorRunner`, `//go:embed vendor/status`, `materializeScript`.
- `internal/status/vendor_runner_test.go` — `VendorRunner` seam/interpreter/procgroup tests (these test the oracle plumbing itself; they die with it).
- `internal/status/oracle_resolution_test.go` — `TestOracleResolvesInTree` (asserts the in-tree oracle resolves; meaningless once there is no oracle).
- `internal/status/harness_test.go` — `oraclePath`, `runOracle`, `runLauncher`, the `-update` flag, `SPACEDOCK_ORACLE`. Keep `normalize`/`realpath`/`pinnedEnv`/`goldenPath`/`readGolden` (the graduated goldens still need them).
- `internal/status/zz_independent_parity_test.go` — `indOraclePath`, `indRunOracle`, the `indDiff` "native==oracle" half. Keep the fixture builders + `indRunNative` + the already-standalone `--new` tests.

**Status parity cases to graduate (52 `t.Run` subtests across the oracle-coupled files):** `archive_guard`, `set_stage_membership`, `nextid_boot`, `next_suppressed_by`, `golden_read` (already golden-backed — just drop its `-update` oracle path), `comment_roundtrip_parity`, `argv_passthrough`, `native_state_dir`, `native_read`, `native_mutation`, `native_guard`, `workflow_dir`, `native_validate`, `native_slug`, `native_discover`, `mutation`, `merge_policy_guard`, `enum_scope`, `worktree_overlay`, `stages_comment_parity`, `native_usage`, `native_eofnewline`, `boot_orphan_abs`, `symlink_profile`. Each currently compares `runNative` against a live `runOracle`/`runLauncher`; graduation re-points the `want` side at a frozen golden captured from the oracle at retirement time.

**Dispatch oracle (delete after graduated):**
- `internal/dispatch/parity_harness_test.go` — `runOracle`, `vendoredOracle`, `assertParity`'s oracle arg. Keep `runNative`, `rewriteOracleFetch`, `stripStateCommitGuidance`, `gitInit`, `writeFile` (the graduated goldens reuse the normalizers).
- The 7 dispatch `*_parity`/`build_*` files (`standing_parity`, `cycle2_parity`, `showstagedef_parity`, `build_hazards`, `build_errors`, `build_parity`, `contextbudget_parity`) — graduate each `runOracle` call to a frozen golden.

**Cross-references that break when `commission/bin/{status,claude-team}` are removed (must update in the same step):**
- `skills/integration/skill_text_test.go:34` — reads `commission/bin/claude-team`; drop it from `rel`.
- `skills/integration/portability_test.go:51` — reads `skills/commission/bin/status` as a scan target; drop it.
  (The contract/surface guard tests at `contract_status_path_test.go` / `skill_surface_test.go` only assert OTHER files do not reference these paths — they do not read the files, so they are unaffected by removal.)

**Stale doc (not load-bearing, clean up):**
- `docs/dev/README.md:128-132` documents `python3 .../commission/bin/status` as the "compatibility phase" command. Replace with the `spacedock status` form already shown below it.

## Verified mechanism (spike, throwaway)

The whole plan rests on ONE unproven mechanism: that a parity case which today diffs `native` against a **live** oracle run can be converted to diff `native` against a **frozen** golden while still catching divergence. I exercised the smallest end-to-end version before committing to ~60 conversions:

- On `internal/dispatch` (`build` backlog case): captured oracle output, normalized + froze it, asserted `runNative` matched the frozen bytes **with the oracle not re-run**, then injected a one-line corruption into the native body and confirmed the frozen golden **rejected** it. `PHASE2 OK` + `PHASE3 OK`, test passed. (Throwaway test deleted; the git log carries no spike file.)
- Surfaced design constraint: the dispatch build body embeds the absolute fixture root **3×** (`strings.Count(body, root)==3`) and it appears in stdout. So freezing dispatch goldens REQUIRES the `<ROOT>` placeholder normalization — exactly the `normalize()`/`indNormalize()`/`stripStateCommitGuidance` helpers the suites already carry. The freeze step normalizes the captured oracle output, and the graduated assertion normalizes the native output the same way before comparing. This is why the harness normalizers are KEPT (not deleted) above.

Baseline confirmed green before planning: `go test ./internal/status/ ./internal/dispatch/` → `ok` (status 14.9s, python3 3.14.4 on PATH).

## Disposition of the differential-parity suite (decided)

**Graduate, do not delete.** Rationale: the ~50+ oracle-coupled cases encode real coverage (EOF-newline identity, CRLF universal-newline read, exotic-score sort, realpath asymmetry, unknown-field preservation, archive-dest spelling, the dispatch cross-product). Deleting them drops that coverage; the spec's AC-3 (a flip-test must still fail) is exactly the anti-deletion guard. Graduation freezes each oracle output as a native-only golden so the case keeps asserting native behavior with the oracle gone.

**How native parsing stays trustworthy once the oracle is gone (the precondition zj depends on):** the goldens ARE the oracle's outputs captured at retirement time. A native regression still fails its golden; only an *intended* native change requires a golden refresh (a deliberate, reviewable diff — the same discipline as any golden suite). The trust does not come from re-running Python; it comes from the frozen capture being the certified-parity bytes. The `zz_independent` builders that build fresh fixtures stay (they prove native behavior on inputs the in-tree testdata never saw), re-pointed from `indDiff(native, oracle)` to `assertGolden(native, frozen)`.

## Ordered retirement plan (smallest reversible first, green at every step)

Each step ends with `go test ./...` green. Steps 1–2 add goldens WITHOUT removing the oracle (fully reversible); the oracle is deleted only at step 4, after the suite no longer reads it.

1. **Add a golden-capture mode and capture dispatch goldens.** Teach the dispatch parity files to write the normalized oracle output to `internal/dispatch/testdata/golden/<case>.txt` (via a `-update`-style flag, mirroring `golden_read`'s existing one) and to assert `runNative` against the frozen golden instead of `runOracle`. Run with capture once to populate goldens, commit goldens, then run the suite normally (oracle still present, but the assertion no longer touches it). Reversible: the oracle and `runOracle` still exist; revert the assertion swap to restore live-diff.
   *Proof:* `go test ./internal/dispatch/` green with `runOracle` calls removed from the assertion path; a flip-test (corrupt one native build branch) fails its golden.

2. **Graduate the status parity suite the same way.** For each of the 52 oracle-coupled status subtests, capture the normalized launcher/oracle output to `internal/status/testdata/golden/` and re-point `want` at the golden. `golden_read` is already golden-backed — only drop its `-update`→`runOracle` regeneration branch (replace with capture-from-native or freeze-as-is). The `zz_independent` builders keep their fixtures; swap `indDiff(...,oracle...)` for a golden compare.
   *Proof:* `go test ./internal/status/` green; flip-test (mutate a native code path) fails a graduated golden — this is AC-3.

3. **Remove the oracle-driver code, keeping the normalizers.** Delete `runOracle`/`indRunOracle`/`runLauncher`/`vendoredOracle`/`oraclePath`/`indOraclePath`/`assertParity`'s oracle side, the `SPACEDOCK_ORACLE` env, the `-update`→oracle regeneration paths, and `oracle_resolution_test.go`. Keep `normalize`/`indNormalize`/`stripStateCommitGuidance`/`pinnedEnv`/golden helpers. At this point NO test reads either oracle.
   *Proof:* `grep -rn 'runOracle\|indRunOracle\|vendoredOracle\|SPACEDOCK_ORACLE' internal/` returns nothing; `go test ./...` green.

4. **Delete the oracle artifacts and the seam.** Remove `internal/status/vendor_runner.go` (+ `vendor_runner_test.go`), `internal/status/vendor/status`, the `//go:embed`, `skills/commission/bin/status`, `skills/commission/bin/claude-team`, and the `__pycache__` debris. Update `skill_text_test.go` + `portability_test.go` (drop the removed scan targets) and `docs/dev/README.md` (drop the `python3 .../status` compat line). The `Runner` interface stays (it is the native runner's contract), but its only implementer is now `NativeRunner`.
   *Proof:* `go test ./...` green **with `python3` removed from PATH** (or `PATH` stripped of python) — proving AC-1's "requires no python3"; `git ls-files | grep -c 'vendor/status\|commission/bin/status\|commission/bin/claude-team'` is 0.

## Bootstrap-graduation precondition (recorded)

**Is every native command proven at parity now?** Yes — the differential-parity suite is currently green (baseline run above) across read (default/next/validate/resolve/short-id/boot/next-id/fields/all-fields/where), mutation (`--set` field/clear/bare-fill/insert, `--archive` flat/folder), `--new`, validation defects, usage errors, and the byte-traps (EOF-newline, CRLF, exotic-score, realpath asymmetry, unknown-field preservation), plus the dispatch build cross-product, standing, show-stage-def, cycle2, context-budget. The capture in steps 1–2 is therefore a faithful certification snapshot, not an aspirational one. This green-now state is the bootstrap-graduation precondition: the goldens inherit the oracle's certified parity.

## Boundary with zj (no scope bleed)

- **02 (this entity) removes the oracle and freezes parity to goldens** — parity-freeze risk only. It does NOT touch the frontmatter parser/mutator (`frontmatter.go`, `mutate.go`, `orderedmap.go`, `stages.go`); they stay byte-for-byte as-is.
- **zj swaps the hand-rolled parser for a YAML library** with documented divergences — library-divergence risk only. zj re-points the desired-behavior tests and accepts deliberate breaks from the frozen Python quirks.
- The two risks stay isolated per sprint-notes "Parity-with-Python is a migration scaffold": 02 is the NOW (match Python, certify, freeze); zj is the POST-PYTHON (standard compliance, deliberate divergences). The **byte-PRESERVATION contract** (unknown fields + order survive `--set`/`--archive`) is KEPT by both — 02 freezes it into the goldens, zj re-implements it via `yaml.Node`. zj is dispatched only AFTER 02 lands (it is `status: backlog`, gated).

## Preconditions (per zj's gating clause)
1. **Native parity certified** — the differential-parity suite is green; native status + dispatch are trusted. **Effectively MET.**
2. **`claude-runtime-segregation` (zs) landed** — removed the last python RUNTIME shell-out. **MET (#246, archived).**
3. **VendorRunner + the embedded vendor/status retired** — **this entity.** The long pole.

## Scope
- **Graduate the differential-parity suite to standalone assertions.** Replace the oracle-comparison
  calls (`runOracle`/`indRunOracle`/`vendoredOracle` in the `native_*` / `zz_independent` / dispatch
  `*_parity` files + the two shared harnesses) with frozen goldens or embedded expected literals, so
  no test needs `python3` or either oracle script. The in-tree templates already exist: the
  `zz_independent --new` desired-behavior tests and `golden_read`'s frozen oracle-capture goldens.
- **Delete the oracle + its seam:** `internal/status/vendor/status` (96568 bytes), the
  `//go:embed vendor/status` (`vendor_runner.go:14`), `VendorRunner` (`vendor_runner.go`) and its
  test instantiations; drop `golden_read`'s `-update`/`runOracle` regeneration path.
- **Retire the dispatch-side vendored test dep** `skills/commission/bin/claude-team` (a hard test dep —
  `parity_harness_test.go` `t.Fatalf` if absent) once dispatch parity is frozen.
- Leave the already-standalone intentional-divergence tests untouched (build-statecommit, json_boot, fields-dedupe).

## Acceptance criteria

**AC-1 — No vendored Python oracle remains, and the test suite runs without python3.** End state:
all three vendored artifacts (`internal/status/vendor/status`, `skills/commission/bin/status`,
`skills/commission/bin/claude-team`) and the `__pycache__` debris are gone; the `//go:embed
vendor/status` and `VendorRunner` are gone; only `NativeRunner` implements `Runner`.
Verified by: `git ls-files | grep -E 'vendor/status|commission/bin/(status|claude-team)'` is empty
AND `go test ./...` passes with `python3` removed from `PATH` (the build no longer embeds or execs a
Python script). The PATH-stripped run is the proof that exercises AC-1 rather than a grep standing in
for behavior.

**AC-2 — The differential-parity suite is graduated, not deleted.** End state: the 52 status
oracle-coupled subtests plus the 7 dispatch parity files assert `runNative` against frozen goldens
captured (and normalized) from the oracle at retirement time, and stay green with the oracle scripts
removed. The `zz_independent` fresh-fixture builders survive (re-pointed to golden compares), so the
coverage breadth is preserved, not just the in-tree fixtures.
Verified by: `go test ./internal/status/ ./internal/dispatch/` passes after step 4, with no
`runOracle`/`indRunOracle`/`vendoredOracle`/`SPACEDOCK_ORACLE` symbol remaining
(`grep -rn` returns nothing).

**AC-3 — Behavior preserved (flip-test).** End state: the graduated goldens still catch the
divergences the oracle-comparison caught — they ARE the captured oracle outputs. Verified by: a
flip-test in the change's own validation — mutate a native code path (e.g. perturb a status table
column or a dispatch body line) and observe a graduated golden assertion FAIL; the spike already
demonstrated this round-trip on the dispatch build case (PHASE3 OK).

**AC-4 — Cross-references and stale docs reconciled.** End state: removing the scanned scripts does
not break the integration guard suite, and no doc tells a reader to run the retired Python script.
Verified by: `go test ./skills/integration/` passes (the `skill_text_test.go`/`portability_test.go`
scan lists no longer name removed files) AND `grep -rn 'commission/bin/status' docs/` finds no
invocation form (the `spacedock status` form remains).

## Sequencing
oracle-retirement → (bootstrap graduation) → `zj` (YAML-library frontmatter migration). Do NOT fold
into `zj`: this freezes parity to goldens (parity-freeze risk); `zj` swaps the parser (library-
divergence risk) — keep the two risks isolated. Post-bootstrap; not on the current dev-workflow-
ergonomics sprint's critical path.

## Test plan

The deliverable is the graduated test suite itself plus the deletions, so the change IS the proof —
no separate test scaffolding is needed beyond the goldens. Costs are local-Go, no live workflow run.

- **Golden capture (steps 1–2):** add a capture flag mirroring `golden_read`'s existing `-update`;
  populate `internal/{status,dispatch}/testdata/golden/`. Cost: low; the normalizers already exist.
  Fixture-level (frozen goldens), CLI-driven through the existing `runNative`/`Run` harnesses.
- **Graduated-suite green (AC-2):** `go test ./internal/status/ ./internal/dispatch/` — ~15s status,
  fast dispatch. Go unit/golden tests, the right altitude for parser+command+emitter behavior.
- **Flip-test (AC-3):** one-off local mutation of a native path → a graduated golden fails; revert.
  Demonstrated in the ideation spike on the dispatch build case. Cost: minutes.
- **python3-free run (AC-1):** `PATH=$(dirs-without-python) go test ./...` — exercises that no test
  execs python. Cost: one run. This is the load-bearing AC-1 proof (behavior, not a grep).
- **Cross-ref/doc (AC-4):** `go test ./skills/integration/` + a `grep` over `docs/` for the retired
  invocation form. Cost: trivial.

No spike beyond the one already run: the only unverified mechanism was the live-oracle→frozen-golden
graduation round-trip, and it was exercised end-to-end (capture → native-vs-frozen → flip-test) on a
real dispatch case before this plan committed to ~60 conversions. Everything else composes
already-proven behavior (the suite is green today; the normalizers and golden helpers exist).

## Notes
The native code lives in `internal/status`, `internal/dispatch`, and `internal/claudeteam` (the
Claude seam holding the `~/.claude` reads). Verified present on THIS checkout (`main`, 2026-06-02):
`internal/claudeteam/{contextbudget,standing,claudeteam,pyjson}.go` all exist, and
`internal/cli/cli.go:44` wires only `NativeRunner` — the earlier FO investigation's "stale `main`
lacks `internal/claudeteam`" reading does not hold here. (Seam #244 + zs #246.)

Touchpoint count correction: the original spec named only `internal/status/vendor/status`. This
ideation found two more vendored Python artifacts — a byte-identical `skills/commission/bin/status`
copy and the dispatch oracle `skills/commission/bin/claude-team` — plus their `__pycache__` debris
and two integration-test scan-target references. All are folded into the touchpoint inventory and
AC-1/AC-4 above. Implementation runs under the `internal/status` serialized-lane rule.

## Stage Report: ideation

- DONE: The retirement plan enumerates EVERY Python-oracle touchpoint to remove and what replaces each, in an order that keeps the binary green + tests passing at each step (smallest reversible first).
  "Touchpoint inventory" + "Ordered retirement plan" sections; two grep sweeps found 3 vendored artifacts (the spec named 1), all listed with their test drivers and the 4-step green-at-each-step order (goldens added before oracle deleted).
- DONE: The differential-parity suite's disposition is decided with rationale (graduate vs delete) — and HOW native parsing stays trustworthy once the oracle is gone.
  "Disposition" section: graduate (freeze captured oracle bytes as native-only goldens); trust comes from the frozen certified-parity capture, not re-running Python. AC-3 flip-test is the anti-deletion guard.
- DONE: Records the bootstrap-graduation precondition (is every native command proven at parity now?) and the boundary with zj (no scope bleed).
  "Bootstrap-graduation precondition" section: yes — suite green now (baseline run, status 14.9s); "Boundary with zj" section: 02 freezes parity, zj swaps the parser; the two risks stay isolated, zj gated until 02 lands.

### Summary
Hardened the retirement into a 4-step ordered plan (add goldens → graduate status suite → strip oracle drivers → delete artifacts), each step ending green and the oracle deleted last (fully reversible until then). Ran the load-bearing spike: proved a live-oracle parity case can be graduated to a frozen-golden assertion that still fails on injected divergence (PHASE2/PHASE3 OK), and surfaced the `<ROOT>` normalization constraint (dispatch body embeds the abs root 3×) that dictates keeping the existing normalizers. Key correction to the spec: there are THREE vendored Python artifacts (a byte-identical `commission/bin/status` copy and the `claude-team` dispatch oracle were uninventoried), plus two integration-test scan references and a stale doc line — all now in the touchpoint list and AC-1/AC-4. Production already ships only `NativeRunner`, so this is purely test-time removal.
