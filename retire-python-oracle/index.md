---
id: 02a92vbcn4x7453bhszprpje
title: Retire the Python oracle — remove the embedded vendor/status + VendorRunner, graduate the differential-parity suite
status: implementation
source: FO investigation (2026-06-01) — the sole open precondition gating zj (yaml-parser-migration), currently orphaned (0x scoped it out; sprint-notes deferred it). Post-bootstrap.
started: 2026-06-02T15:29:46Z
completed: 2026-06-02T19:48:24Z
verdict: PASSED
score: "0.25"
worktree: 
issue:
mod-block: 
pr: "#270"
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

## Stage Report: implementation

- DONE: Execute the 4-step ordered plan keeping the suite green AT EACH STEP (capture goldens; graduate status+dispatch; strip drivers; delete artifacts, oracle LAST).
  5 commits af709249→4d5fdb4f, each ending `go test ./...` green; oracle deleted only in the final commit (reversible until then).
- DONE: AC-1 proof — full `go test ./...` passes with python3 REMOVED from PATH (behavior, not grep).
  `env -i PATH=<git+gh+go, no python3> go test ./...` → all 12 packages ok, uncached; `git ls-files | grep -E 'vendor/status|commission/bin/(status|claude-team)'` = 0.
- DONE: AC-3 flip-test — mutate a native path → a graduated golden FAILS, then revert.
  Dispatch: perturb the build body `Stage:` emit → build-mods golden fails; revert → 112 green. Status: perturb the table header → many graduated goldens fail; revert → 324 green.
- DONE: AC-2 — no runOracle/indRunOracle/vendoredOracle/SPACEDOCK_ORACLE symbol remains; zz_independent builders survive (re-pointed to golden compares).
  `grep -rn 'runOracle|indRunOracle|vendoredOracle|SPACEDOCK_ORACLE|VendorRunner|go:embed' internal/` = empty; indSeq/indSDB32/indSlug fresh-fixture builders kept, re-pointed indDiff→indGolden.
- DONE: AC-4 — `go test ./skills/integration/` green; no doc names the retired commission/bin/status invocation.
  Integration 26 passed; skill_text_test.go drops the claude-team read; docs/dev/README.md drops the python3 .../status compat lines (spacedock status stays); no `python3 .../status` invocation in docs.
- DONE: Full suite green + gofmt/vet clean.
  `go test ./...` 701 passed (python3 present); gofmt -l empty; go vet clean.

### Summary
Graduated the ~60 oracle-coupled parity subtests (52 status + 7 dispatch files) to frozen goldens captured from the certified-parity native output, via a golden harness in each package (normRun/normPaths with `<ROOT>`/`<HOME>`; envelope + text goldens over the existing `-update`+`normalize`). Validated the live-oracle→frozen-golden mechanism end-to-end before converting (capture→compare→flip→revert). Then stripped all oracle-driver code (keeping the normalizers the goldens reuse) and deleted all three vendored Python artifacts + the `//go:embed`+VendorRunner seam + procgroup helper; NativeRunner is now the sole Runner.

Three deviations from the named touchpoints, all surfaced to team-lead: (1) `workflow_dir_test.go` tested the launcher's `dirname(__file__)` empty-scan landmine with no native analog (NativeRunner discovers the enclosing workflow) — deleted with the seam tests; native no-flag discovery stays covered by discover_walkup_test.go. (2) The whole-suite AC-1 run surfaced a 4th, un-inventoried python3 exec — `internal/claudeteam/pyjson_test.go` cross-checked EmitPythonJSON against live `python3 json.dumps`; graduated to frozen `ensure_ascii` literals (the last python3 exec anywhere in the tree). (3) `portability_test.go:51` lists `skills/commission/bin/status` as a forbidden-substring guard (NOT a file read, despite the spec wording) — kept, since the "path absent on a fresh install" premise holds even more strongly after deletion; raised the keep-vs-drop choice with team-lead. AC-1 is proven behaviorally (python3-free `go test ./...` green), not by grep. zj (yaml-parser-migration) is now unblocked.

## Stage Report: validation

- DONE: AC-1 (load-bearing): actually STRIP python3 from PATH and run full `go test ./...` — confirm green AND that NO test SKIPS for a missing interpreter. Don't trust the report; reproduce. Also chase the `python3 missing` error string at internal/cli/status_test.go:112/117 — confirm it is NOT a live python3 code path in production (NativeRunner only).
  Built a hermetic bin dir (go/git/gh/coreutils, zero python: `command -v python3/python/python3.13/python3.14` all empty); `env -i PATH=<nopython> go test ./... -count=1` → all 10 test packages ok, exit 0. Verbose run: 0 FAIL, only 3 SKIPs — all host-CLI gating (codex/claude absent), NONE python-interpreter-missing. No `exec.Command/LookPath(python)` anywhere in .go. status_test.go:112/117 is a `fakeRunner` + `errFake("...python3 missing")` literal (file comment: "no Python, no exec, no real script") proving CLI surfaces runner errors loudly — not an exec. NativeRunner is sole Runner impl; cli.go:44 wires only it.
- DONE: AC-2/AC-3: confirm no runOracle/indRunOracle/vendoredOracle/SPACEDOCK_ORACLE/VendorRunner symbol remains; vendored artifacts gone. Flip-test BOTH a native status path and a native dispatch path → graduated goldens FAIL, then revert.
  `grep -rnE 'runOracle|indRunOracle|vendoredOracle|SPACEDOCK_ORACLE|VendorRunner'` empty; no `go:embed`; `git ls-files | grep -E 'vendor/status|commission/bin/(status|claude-team)'` empty; no `__pycache__`/`.pyc` tracked; all 7 named paths absent on disk. zz_independent builders (indSeq/indSDB32/indSlug) survive, re-pointed indDiff→indGolden. FLIP #1 (status): perturb format.go:111 header STATUS→STATE → 35 graduated goldens FAIL (native STATE vs frozen STATUS); revert → 319 PASS. FLIP #2 (dispatch): perturb build.go:309 body "Stage:"→"Phase:" → 14 goldens FAIL (TestBuildParityCrossProduct et al, native Phase vs frozen Stage); revert → 112 PASS. Bonus FLIP (pyjson): break EscapeNonASCII DEL branch → del\x7fx golden FAIL; revert → green. Goldens are captured-oracle bytes, not self-comparisons.
- DONE: AC-4: `go test ./skills/integration/` green (portability + skill_text scans pass; no doc names the retired commission/bin/status invocation form). PASSED/REJECTED with the reproduced evidence.
  `env -i PATH=<nopython> go test ./skills/integration/ -v` → 26 PASS, 0 FAIL, 0 SKIP. skill_text_test.go drops `commission/bin/claude-team` from `rel` (verified diff). portability_test.go keeps `skills/commission/bin/status` as a forbidden-SUBSTRING guard (internalHelperPaths, not a file read) — deviation #3 is correct; deleting the file strengthens the "absent on fresh install" premise. docs/dev/README.md: both `python3 /path/to/.../status` invocation forms removed, `spacedock status` form remains (verified diff). `grep -rn 'commission/bin/status' docs/` finds ONE hit — bootstrap-roadmap.md:60 "no skill path references plugin-private `skills/commission/bin/status`" — a NEGATIVE invariant statement, not an invocation form; AC-4's clause forbids invocation forms (the `python3 .../status` command), which are all gone.

### Summary
PASSED. All four ACs validated behaviorally (not by grep/report re-read), with python3 absent from PATH throughout the load-bearing runs. AC-1: full `go test ./...` green with zero python3 reachable; the only 3 skips are host-CLI (codex/claude) gating, no interpreter-missing skip; the `python3 missing` string is a test-fixture error literal, not a live exec — NativeRunner is the only Runner. AC-2/AC-3: all oracle symbols/artifacts gone, zz_independent builders preserved, and BOTH a native status table path and a native dispatch body path flip-tested — perturbation makes the graduated goldens FAIL (35 status, 14 dispatch), revert restores green, proving the goldens are real captured-oracle oracles. AC-4: integration suite green, guard-test scan targets reconciled, README invocation forms retired. The pyjson graduation (deviation #2) kept full ensure_ascii coverage: all 10 frozen literals independently re-derived against real `python3 json.dumps(ensure_ascii=True)` match exactly, and a production flip (DEL branch) is caught. Deviations #1–3 all checked sound: workflow_dir_test.go coverage subsumed by discover_walkup_test.go; pyjson literals certified; portability guard correctly kept. Worktree clean after all reverts; gofmt and go vet clean. zj (yaml-parser-migration) is unblocked.

## Feedback Cycles

**Cycle 1 (FO + 4-lens detached audit, 2026-06-02) — validation PASSED but was ENV-MASKED; the audit REJECTS: the suite is red-by-construction on CI / any fresh clone. Routed to implementation.** The validation ensign ran in the mint-time worktree where the frozen golden matched, so it passed; the detached audit reproduced the FAILURE in a different-path worktree at the same commit (native emits `prz9…` at the mint path, `jksm…` at another). Three material findings:

- **M1+M2 — the sd-b32 `--next-id` parity check was graduated to a STATIC golden (`prz9n8rry3q0kgaxs2rgbmb1`) that bakes in the absolute checkout path.** The candidate is `sha256(… "workflow="+realpathOf(definitionDir) …)` (identity.go:165), so the frozen byte reproduces ONLY on the freeze-time machine. `TestNextIDParity` + `TestNativeNextIDMatchesOracle` PASS in the ensign worktree, FAIL in any other path (proven) → CI runs `go test ./...` with no path pinning, so it is RED-by-construction. The test's own surviving comment (nextid_boot_test.go:9-11) says --next-id must be asserted at format+determinism level, NOT a static golden, because the candidate is SHA-derived; step-2 froze it anyway and mislabeled it "Deterministic under the pinned env". **CAPTAIN STEER (2026-06-02): byte-compatibility with the old Python implementation is NOT the goal for --next-id — only LOGICAL SOUNDNESS.** So do NOT chase the old oracle byte. FIX: replace the byte-frozen --next-id golden with a LOGICAL assertion — (a) FORMAT (24-char, valid sd-b32 alphabet), (b) DETERMINISM (same seed/actor/path → same id; pure function of the digest material), and (c) a path-INDEPENDENT derivation vector: feed a FIXED literal workflow path (not realpathOf(fixtureDir)) so the candidate is reproducible on any machine, and assert the SHA-256 + 5-bit-extraction derivation against that fixed vector. This KEEPS the minting-regression detector M2 worries about (the 5-bit-mask 31→30 flip at identity.go:151 must still FAIL) WITHOUT byte-parity to Python. Do NOT resolve by `-update` (that re-pins the local path-dependent value).

- **M3 — the over-broad `sdB32Re→"<ID>"` line added to `indNormalize` (zz_independent_parity_test.go:78-79, ABSENT on origin/next) collapses ANY 24-char token on read-path goldens → partial tautology.** Mutation proof: emitting a wrong-but-valid id left `TestIndReadFlagsSDB32/resolve` PASSING. The fixture entity id `abcdefghjkmnpqrstvwxyz23` is a FIXED literal, so read-path goldens never needed it. FIX: scope the `<ID>` normalization to ONLY the minted `--next-id`/`NEXT_ID` output (the seeded next-id + boot NEXT_ID line), NOT the read-path goldens where the stored id is fixed fixture content; freeze the real fixture id there.

- **Polish (fold in):** scrub the now-stale comments asserting oracle-equality / "deterministic under the pinned env" (nextid_boot_test.go header, boot_probe_parity_test.go:81); optionally fix the stale `python3 missing` errFakeNoPython wording.

**Required proof (behavioral):** the full suite passes in TWO different-path worktrees (path-independence — the exact lever the audit used to catch this); the 5-bit-mask 31→30 mutation still FAILS the --next-id derivation check (regression detector intact); the M3 read-path mutation (a wrong-but-valid id) now FAILS. The residual-coupling, dispatch, claudeteam, and pyjson graduations were audited CLEAN and mutation-flipping — do NOT touch them.

## Stage Report: implementation (cycle 1)

- DONE: M1+M2 — replace the path-dependent byte-frozen sd-b32 --next-id goldens with logical assertions (format + determinism + a path-independent derivation vector).
  TestNextIDParity / TestNativeNextIDMatchesOracle → *FormatAndDeterminism (24-char valid sd-b32, two runs match); new TestSDB32CandidateDerivationVector drives sdB32Candidate over a FIXED literal workflow path (realpathOf returns it unchanged on any host) and asserts the SHA-256 + 5-bit minting against an independently-derived id (n1y6x7mw00etcc0v1c9mzwwr). Deleted nextid-parity-candidate.txt + native-nextid-candidate.txt; no -update re-pin. Commit 389df537.
- DONE: M2 regression detector intact — the 5-bit-mask 31→30 flip (identity.go:151) still FAILS.
  Flipped mask to 30 → TestSDB32CandidateDerivationVector FAILS (got m0y6w6mw…, want n1y6x7mw…); revert → green. Identical result at the alt path (proving the vector is path-independent).
- DONE: M3 — scope the <ID> masking to ONLY the minted next-id/NEXT_ID output, not read-path goldens.
  Removed the over-broad sdB32Re→<ID> from indNormalize; added maskBootNextID (anchored ^NEXT_ID:) used by both the ind goldens and the in-tree boot tests; re-captured the sdb32 read/boot goldens to carry the literal fixture id abcdefghjkmnpqrstvwxyz23; ind next-id-seeded dropped its golden for format+determinism. Mutation proof: inject a wrong-but-valid resolve id (zzzz…) → ind-sdb32-resolve + native-read-resolve + GoldenRead/resolve all FAIL; revert → green.
- DONE: Polish — scrub stale oracle-equality / "deterministic under the pinned env" comments + the python3-missing errFake wording.
  nextid_boot header + TestBootStructuralParity comment + boot_probe_parity_test.go:81 rewritten to "frozen golden"; errFakeNoPython→errFakeRunnerFailed ("runner failed", generic) in cli/status_test.go.
- DONE: REQUIRED PROOF — full `go test ./...` passes in TWO different-path worktrees.
  Mint path: 702 passed. Copied the worktree (rsync, excl .git) to /var/folders/…/spacedock-altpath-proof and ran `go test ./...` there → all 10 test packages ok. This is the exact path-independence lever the cycle-1 audit used to catch the red-by-construction failure; it now passes.
- DONE: AC-1 still green python3-free; dispatch/claudeteam/pyjson graduations untouched.
  `env -i PATH=<git+gh+go, no python3> go test ./...` → 10 packages ok, no failures. `git diff --stat HEAD~1 -- internal/dispatch/ internal/claudeteam/` empty. gofmt/vet clean.

### Summary
Fixed the cycle-1 REJECT: the sd-b32 `--next-id` parity check was graduated to a static golden that baked in the absolute checkout path (the candidate hashes `realpathOf(definitionDir)`), so it only reproduced on the mint-time machine — red-by-construction on CI. Per the captain steer (logical soundness, not Python byte-parity), replaced the byte-freeze with format + determinism + a path-independent SHA-256/5-bit derivation vector over a fixed literal path; the 31→30 mask flip still fails it, identically at any path. Scoped the over-broad `<ID>` masking to only the minted NEXT_ID line so read-path goldens freeze the real fixture id (a wrong-but-valid id now fails). Verified path-independence by running the full suite in two different-path worktrees — the exact lever the audit used — both green. The residual-coupling, dispatch, claudeteam, and pyjson graduations the audit cleared were left untouched.

## Stage Report: validation (cycle 1)

- DONE: PATH-INDEPENDENCE (the lever): full `go test ./...` green in the entity worktree AND at a DIFFERENT absolute path — BOTH must be green.
  Two checkouts of fix commit 389df537 at genuinely different abs paths: mint (`.worktrees/spacedock-ensign-retire-python-oracle`) and a `git worktree add --detach /tmp/sd-02-altpath`. Python-free (`env -i PATH=<go+git+gh, no python3>`), uncached: BOTH → 10/10 pkgs ok, exit 0, 699 PASS / 0 FAIL / 3 SKIP. Negative control: the PRE-fix commit 4d5fdb4f at `/tmp/sd-02-prefix-altpath` FAILS exactly TestNativeNextIDMatchesOracle + TestNextIDParity (path-dependent goldens) — proving the detector is real, not vacuous, and the cycle-0 miss is closed.
- DONE: The --next-id fix is a LOGICAL assertion now (format + determinism + path-INDEPENDENT derivation vector over a fixed literal path), not a byte-frozen golden; confirm the asserted value does NOT depend on the checkout path, AND the 31→30 mask-flip at identity.go:151 still FAILS the derivation check. M3: read-path goldens pin the real fixture id — a wrong-but-valid resolve id FAILS TestIndReadFlagsSDB32/resolve.
  Independently re-derived `wantID` from the documented algorithm (SHA-256 over the 7-line digest material, big-endian 5-bit windows, Crockford-32 alphabet, mask 31) in a from-scratch Python reference → `n1y6x7mw00etcc0v1c9mzwwr`, EXACT match to the test literal (not a self-comparison). Fixed literal path `/spacedock/derivation-vector/wf` is path-independent (realpathOf returns it unchanged). M2 flip: mask 31→30 → derivation vector FAILS with `got "m0y6w6mw00etcc0t0c8mywwr"` (matches my independent mask-30 derivation) at BOTH paths identically; revert → green. M3 flip: change fixture id `abcdefghjkmnpqrstvwxyz23`→`zzzzzzzzzzzzzzzzzzzzzz23` (both valid sd-b32) → TestIndReadFlagsSDB32/resolve FAILS (native `id=zzzz…` vs golden literal `id=abcd…`); revert → PASS. maskBootNextID is anchored to `^NEXT_ID:` only, so stored fixture ids freeze literally.
- DONE: AC-1 python3-free full suite still green; the audited-clean parts (dispatch/claudeteam/pyjson graduations) untouched; gofmt/vet clean.
  Both path-independence runs above were python3-free and green; `command -v python3` returned nothing in the hermetic bin. `git show 389df537 --name-only` touches NO file under internal/dispatch/ or internal/claudeteam/ (graduations untouched). No `exec.Command/LookPath(python)` anywhere; oracle symbols still absent. The status_test.go polish removes the last `python3` string (errFakeNoPython→errFakeRunnerFailed) — still a fake error literal, not an exec. gofmt -l empty; go vet clean (python-free).

### Summary
PASSED. The cycle-1 fix resolves the cycle-0 false positive. The load-bearing path-independence check is satisfied: the full python3-free suite is green at TWO genuinely different absolute paths (mint + /tmp/sd-02-altpath), 699/0/3 each, and a negative control (pre-fix commit at a third path) FAILS exactly the two formerly-path-dependent goldens — proving the detector works and the red-by-construction bug is gone. M1/M2: the --next-id check is now a logical assertion (format + determinism + a path-independent SHA-256/5-bit derivation vector); I independently re-derived the expected id from the algorithm spec (exact match, not a self-comparison) and confirmed the 31→30 mask-flip still FAILS the vector identically at both paths (regression detector intact). M3: the over-broad masking is replaced by a `^NEXT_ID:`-anchored mask, so a wrong-but-valid resolve id now FAILS the read-path golden. The dispatch/claudeteam/pyjson graduations the prior audit cleared are untouched by the fix commit; gofmt and vet clean; all flip-tests reverted, worktrees and temp checkouts cleaned up.

## Stage Report: implementation (cycle 2)

- DONE: `git rebase origin/next` — branch now sits cleanly atop current next (ahead 7, behind 0), KEEPING the cycle-1 logical --next-id fix + the M3 masking scoping intact.
  6 commits replayed; the cycle-1 derivation vector + format/determinism tests + the M3-scoped read-path goldens survived the replay unchanged.
- DONE: Resolve the parity_harness_test.go conflict by taking BOTH intents.
  #263 (7c17a61f) updated the stateCommitGuidance strip regex + push-reminder handling in that file; step-3 had DELETED those oracle-reconciliation helpers (rewriteOracleFetch/stateCommitGuidanceLine/stripStateCommitGuidance/assertParity) because the graduated goldens freeze native's full body (no oracle to strip). Resolution: kept the step-3 deletion (the helpers have NO caller post-graduation — confirmed by grep), and adopted #263's behavioral change via golden re-capture (below). No dead code left.
- DONE: README + integration secondary changes reconciled (both sides kept).
  docs/dev/README.md: my python3-invocation removals survived AND #266's validation-stage policy is present (no `python3` anywhere in the file; "Workflow State" + "Testing Resources" rows show the `spacedock status` form). skills/integration/skill_text_test.go: my claude-team read-drop survived (only the forbidden-substring guard remains) alongside #268's commission-decision-rule demotion.
- DONE: CRITICAL — re-verify the dispatch goldens against current next's build.go (#263 changed stateCommitGuidance).
  3 split-root cross-product goldens (split+folder+worktree, split+flat+worktree, split+folder+nonworktree) froze the pre-#263 guidance ending "after a short wait." and broke; re-captured from current native output — they now carry the push reminder ("Then push the state branch ... `pull --rebase` ... then re-push", state branch = spacedock-state/state-checkout), <ROOT>-normalized. Scanned all guidance-bearing goldens: every one carries the push reminder, none stale. `go test ./internal/dispatch/` green (112). Commit 9180b8a3.
- DONE: Re-run the FULL python3-free `go test ./...` ATOP CURRENT NEXT + re-confirm path-independence.
  python3-free run (hermetic PATH: git/gh/go + coreutils, zero python) → all 10 test packages ok. (A transient release-pkg failure was MY harness's too-narrow PATH missing `tr` — the release empty-body guard is `[ -z "$(tr -d '[:space:]' < notes)" ]`, not a python3 dep; adding coreutils made it green, confirming no new python coupling.) Path-independence re-proven: the cycle-1 derivation vector + a full second-path worktree run (copy under a different abs path) both green; the 31→30 mask mutation still FAILS the vector; revert restores green. Full suite green (722, python present); gofmt/vet clean.

### Summary
Integration reconciliation atop the advanced origin/next (which merged #263 statecommit push reminder, #266 README validation policy, #268 integration test reshape). The one real conflict — parity_harness_test.go, where #263 edited the oracle-strip helper that step-3 had deleted — resolved by keeping the deletion (no caller post-graduation) and adopting #263's behavior through a golden re-capture: the 3 split-root build goldens now freeze native's full body including the new push/pull-rebase reminder. README and integration changes merged cleanly (both sides kept). All proofs re-run atop current next: python3-free full suite green (10 packages), path-independence re-proven at a second path, the 31→30 minting-regression detector intact, dispatch suite green. The cycle-1 logical --next-id fix and M3 masking scoping survived the rebase unchanged.

## Stage Report: validation (cycle 2)

- DONE: PATH-INDEPENDENCE (the lever): full python3-free `go test ./...` green in the entity worktree AND at a SECOND absolute path; cycle-1 derivation vector + 31→30 mask detector hold; M3 read-path scoping intact.
  Two checkouts of cycle-2 tip 9180b8a3 (`.worktrees/spacedock-ensign-retire-python-oracle` + `git worktree add --detach /tmp/sd-02-c2-altpath`), python-free (`env -i PATH=<go+git+gh+coreutils, no python3>`), uncached: BOTH → 10/10 pkgs ok, exit 0, 719 PASS / 0 FAIL / 3 SKIP (the 3 SKIPs are host-CLI codex/claude gating, not python). M2: wantID `n1y6x7mw00etcc0v1c9mzwwr` unchanged; 31→30 mask flip at identity.go:151 → TestSDB32CandidateDerivationVector FAILS (got `m0y6w6mw00etcc0t0c8mywwr`); revert → green. M3: read-path goldens still pin the literal fixture id (maskBootNextID anchored to `^NEXT_ID:` only).
- DONE: THE RE-CAPTURED DISPATCH GOLDENS ARE REAL (not a frozen divergence): `go test ./internal/dispatch/` green; flip a native dispatch path → a re-captured golden FAILS, revert → green; goldens carry #263's push-reminder `<ROOT>`-normalized, NOT a baked absolute path.
  `go test ./internal/dispatch/` green (112). Flip build.go:473 push-reminder "Then push"→"Then publish" → EXACTLY the 3 re-captured split-root cross-product goldens FAIL (split+folder+worktree+team, split+flat+worktree+team, split+folder+nonworktree+team); revert → green. The goldens carry the #263 reminder `<ROOT>`-normalized (`git -C <ROOT>/state-checkout/state-checkout push origin spacedock-state/state-checkout; … pull --rebase … then re-push`), state branch resolved to a literal, paths `<ROOT>`-placeholdered. The only abs-looking strings (`/tmp/spacedock-dispatch/…`) are the fixed DISPATCH_FILE contract literal, identical on every host.
- DONE: The conflict resolution kept BOTH intents (oracle-strip + #263 push-reminder); README (#266) + skill_text (#268) merged cleanly; audited-clean claudeteam/pyjson graduations untouched; gofmt/vet clean.
  `stripStateCommitGuidance` helper gone, no dangling reference (kept step-3's deletion); parity_harness_test.go keeps normalizers (runNative/gitInit/writeFile), no oracle symbols; no git conflict markers anywhere. README has only `spacedock status` (no python3); skill_text_test.go `rel` no longer lists `commission/bin/claude-team`. claudeteam/ touched ONLY by ef2204ac (step-4 graduation, the cycle-0-audited one) — neither cycle-1 (5719ee6d) nor cycle-2 (9180b8a3) touched it; all 10 pyjson ensure_ascii literals independently re-derived against real `python3 json.dumps(ensure_ascii=True)` (exact), DEL-branch flip caught. No python exec anywhere; oracle symbols gone; vendored artifacts absent. gofmt -l empty; go vet clean (python-free).

### Summary
PASSED atop current origin/next (tip 09e1fdce; branch 0 behind / 7 ahead, tip 9180b8a3). Path-independence — the load-bearing lever — re-proven: the full python3-free suite is green at TWO genuinely different absolute paths (719/0/3 each). The cycle-1 logical --next-id fix survived the rebase: wantID unchanged, the 31→30 mask-flip still FAILS the derivation vector, and M3 read-path scoping holds. The 3 re-captured split-root dispatch goldens are REAL oracles — a one-word flip of the #263 push-reminder fails exactly those 3 cases (revert → green) — and they carry the reminder `<ROOT>`-normalized, not a baked path, confirming the conflict resolution merged both the oracle-strip graduation and #263's behavioral change. README/#266 and skill_text/#268 merged cleanly; the audited-clean claudeteam/pyjson graduations were untouched by the cycle-1/cycle-2 work (re-verified faithful + flip-caught anyway). One process note: my first python-free run surfaced a single FAIL (TestReleaseYAMLGuardRejectsEmptyBody) traced to a missing `tr` in my hermetic bin — the release empty-body guard uses `tr -d '[:space:]'`, a coreutil, not python; adding coreutils (still zero python) turned it green, confirming it is NOT a code/python dependency. gofmt and vet clean; all flip-tests reverted, mint worktree pristine, temp checkout removed.

## Stage Report: implementation (cycle 3)

- DONE: Fix the second path-dependent golden (same disease class as cycle-1) — worktree-overlay-resolve.txt froze the macOS-only `workflow=/private%ROOT%`.
  Root cause: TestWorktreeOverlayActiveReads normalized output with a bare strings.ReplaceAll(out, nativeRoot, "%ROOT%") using the AS-SPELLED t.TempDir root, but resolve.go emits the workflow= field via realpathOf(workflowDir) — on macOS that resolves through /private, stranding it before the placeholder. Replaced the bare replace with overlayRoot(), which strips the realpath'd spelling first then the as-spelled root (mirroring normalize()), and re-captured the golden → `workflow=%ROOT%` (no /private). Commit 3b82adaf.
- DONE: Grep ALL status goldens + normalizers to confirm worktree_overlay was the only as-spelled/bare-replace case.
  Re-confirmed: every other root normalizer (normalize at harness_test.go, indNormalize at zz_independent) is already realpath-first-then-as-spelled; worktree_overlay's overlayRoot now matches. No other golden carries a stranded /private (`grep -rn /private internal/status/testdata/golden/` → only the golden-dir README prose, not a frozen value).
- DONE: REQUIRED PROOF — Linux CI offline=success (the oracle for this macOS-vs-Linux bug; local macOS cannot expose it).
  Pushed the branch (tip 3b82adaf) and triggered runtime-live-e2e.yml via workflow_dispatch. Run 26839286383 (spacedock-dev/spacedock, event workflow_dispatch, head SHA 3b82adaf). The `offline` job (runs-on: ubuntu-latest, NO secret) CONCLUDED **success** — steps "Build" (go build ./...) and "Run offline test suite" (go test ./...) both green on Linux, exercising TestWorktreeOverlayActiveReads/resolve against the re-captured workflow=%ROOT% golden (which would have FAILED on the old /private%ROOT% value on Linux). Job URL: https://github.com/spacedock-dev/spacedock/actions/runs/26839286383/job/79141651953. The two claude-live matrix legs sit in `waiting` on their per-variant environment approval gates (by design — they spend ANTHROPIC_API_KEY and are not needed for this verification); only offline=success is load-bearing.
- DONE: Local checks (non-load-bearing for this bug, but clean): full go test ./... green (722), gofmt -l empty, go vet clean. Cycles 1+2 intact (the --next-id logical fix, M3 masking scoping, and the #263 push-reminder golden re-captures all unchanged).

### Summary
Fixed a second path-dependent golden of the same class as cycle-1, in a test cycle-1 never touched: worktree-overlay-resolve.txt froze a macOS-specific `/private%ROOT%` because TestWorktreeOverlayActiveReads used a bare as-spelled-root replace while the --resolve workflow= field is realpath'd (resolve.go), and macOS t.TempDir resolves through /private. Replaced the bare replace with a realpath-aware overlayRoot() (strip realpath'd spelling first, then as-spelled — the same pattern every other status normalizer uses) and re-captured the golden to workflow=%ROOT%. Confirmed by grep that worktree_overlay was the sole as-spelled case. Because the bug is macOS-vs-Linux (both local paths resolve through /private, so a local run — even two paths — cannot expose it), the load-bearing proof is the Linux CI offline job: run 26839286383's `offline` (ubuntu-latest, secret-free `go build` + `go test ./...`) concluded success, with the overlay-resolve test green on Linux.
