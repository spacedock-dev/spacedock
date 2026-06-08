# Sprint 0198 — pre-cut antipattern audit

**Verdict: SHIP-CLEAR — 0 blockers.**

Independent pre-cut antipattern audit of the assembled `next` (sprint
0198-pre-flip-hardening) BEFORE the v0.19.8 tag. Four merged members
(kb #327, qa #328, z9 #329, vh #330). Read-only over code; tests run. Audit
performed at HEAD `d4b89b3a` in the detached checkout
`.worktrees/audit-0198`. The audit found no test weakness, no proof-policy
violation, no broken integration, and no half-done work the tag would freeze.
`go test ./...` from the repo root is GREEN. Nothing must change before the cut.

## Ship-blockers

None.

## How each of the five dimensions was checked

### 1. Shipped test weakness / proof-policy violations — CLEAR

Every merged test in the four members was read against the validation section
and the "Instruction-file read quarantine" in `docs/dev/README.md`. None is a
tautology, a prose-grep over an instruction file, or a self-referential
assertion of the implementer's own text.

- **qa `internal/contract/version_message_test.go`** — behavior fixtures over
  `RunDoctor`, NOT prose-grep. The oracle (`binaryVersionForTest = "0.19.4"`,
  `version_message_test.go:17`) lives OUTSIDE the message under test: it is
  passed INTO `RunDoctor` and the test asserts the rendered stderr/stdout
  CONTAINS it (`:57`, `:84`) — proving the version reached the message. The
  plugin version comes from the fixture manifest (independent source, `:54`,
  `:81`). The negative assertions (`contractTokenPattern` must NOT match,
  `:60`; `hasHalfOpenRange` must be false, `:63`) test rendered runtime output,
  not file contents. Independently confirmed at runtime: `spacedock doctor
  --host claude --plugin-manifest …/too-old-binary.json` renders `Spacedock
  version mismatch: binary 0.19.0, plugin 0.13.0` with no `contract N` token
  and no `>=N,<M` range — the jargon is genuinely gone, not just absent from a
  string literal.
- **qa AC-3 (the admitted PR#262 legitimate-text exception)** — VERIFIED it did
  NOT add a banned prose-grep test. The deliverable is a one-line prose append
  to `skills/first-officer/references/first-officer-shared-core.md` ("Once
  spacedock is on PATH, launch with `spacedock claude` …"). Grep confirms NO
  `.go` test asserts that payoff phrase and NO test outside `contractlint` reads
  `first-officer-shared-core.md` as part of this sprint (the one ensigncycle
  file that reads it, `teardown_marker_consistency_test.go`, is pre-existing and
  untouched by any of the four commits — see Record R3). The qa entity frames
  AC-3 exactly as the PR#262 text exception, and its validation ran a detached
  adversarial audit (3 breaking edits, each reddens the suite — non-vacuous).
- **z9 `install_tolerance_codex_test.go` / `install_behavior_codex_test.go`** —
  strong behavioral tests. The tolerance tests drive the real `execHost.Install`
  control flow against a per-PATH stub codex whose echoed argv
  (`stub:<args>:exit=<code>`) is the independent source of truth; they assert
  order, the tolerance asymmetry (both cleanup steps tolerated, both pin steps
  fail-fast), error-argv wrapping, and `--ref` omission on empty branch
  (`install_tolerance_codex_test.go:165`). `TestCodexPluginInstallIsHostNative`
  drives the REAL `codex` CLI against an isolated `CODEX_HOME` and observes
  on-disk state (`codex plugin list --json` `installed:true` + cache manifest
  exists) — oracle is the host CLI's own JSON, not the source. It RAN on this
  box (0.23s, not skipped).
- **vh `survey_probe_test.go`** — EXECUTES the shipped step-1 one-liner extracted
  from `SKILL.md` under synthesized PATH conditions; the oracle is the two
  fixture conditions (`""` for present, `AGENTSVIEW MISSING` for absent), not a
  grep over the file. Reading the file to extract the runnable line is
  "executing the artifact," and the regex that rejects FS-access forms is part of
  the guard, not the assertion. vh's claim that it added NO prose-grep holds.
- **vh `survey_queries_test.go` / `survey_sync_codex_test.go`** — EXECUTE the
  labeled SQL from `references/queries.sql` against an independent fixture DB
  (query-smoke) and against a DB produced by a real `agentsview sync`
  (end-to-end). Expected values come from the fixture rows / synced rows, which
  diverge from skill prose. Both are behavioral, non-vacuous.
- **kb `migration_check_test.go`** — `TestMigrationCheckPrunesStateTree`
  (`:186`) is a hermetic positive prune-proof driving `filepath.Walk` through
  the production `isMigrationCheckPrunedDir` predicate; it asserts the
  `.spacedock-state` subtree is never visited AND the real entity IS visited
  (`:241`), so the prune does not over-reach. The `checked == 0` Fatal guard
  (`:157`) stays non-vacuous.

**Contractlint quarantine boundary holds:** none of the four commits touched any
`internal/contractlint/*` file or added any instruction-file read outside that
package. `go test ./internal/contractlint/` is green.

### 2. Front-door integration (qa × z9) — CLEAN

qa (#328) and z9 (#329) both edit `frontdoor.go` and `host_exec.go` and merged
sequentially. The second merge did not revert or break the first:

- qa's version-bearing gate (`gateHost`, `frontdoor.go:116`) is intact and is
  called by BOTH front doors after both merges: `gateHost(ops, "claude", …)`
  (`frontdoor.go:174`) and `gateHost(ops, "codex", …)` (`frontdoor.go:322`). It
  threads the binary `Version` into `contract.ManifestVerdict(…, Version)`
  (`frontdoor.go:130`).
- z9's `runCodex` switch (`frontdoor.go:321-342`) coexists with qa's gate —
  `NoPluginFound → ops.Install("codex", marketplaceSource, devBranch) → launch`.
- z9's now-false-comment fixes LANDED and describe the NEW behavior:
  - `host_exec.go:304` — was "only supported for claude" → now `"programmatic
    install is supported for claude and codex, not %q"`; `codexInstallArgvSequence`
    (`:264-286`) documents the codex arm.
  - `host_exec.go:29-34` — the stale "0.136.0 has no --json" comment is gone;
    now describes the codex `ResolveManifest` schema difference (no install
    path).
  - `frontdoor.go:317-320` — was "codex has nothing to auto-install" → now
    describes the codex auto-install mirroring `runClaude`.
- z9 did NOT touch the claude `installArgvSequence` body (the #329 diff shows no
  change to it), so qa's claude path is byte-identical post-merge.
- z9 correctly inverted the old `TestCodexFrontDoorNoPluginFailsFastWithoutInstalling`
  into `TestCodexFrontDoorNoPluginAutoInstalls` (`frontdoor_test.go:436`) +
  `TestCodexFrontDoorNoPluginNoInstallRefuses` (`:470`).

Channel-tracking constraint satisfied: `runCodex` installs via the shared
`devBranch` var (`frontdoor.go:49`), not a hardcoded `next`; the test asserts the
recorded seam `{"codex", marketplaceSource, devBranch}` (`frontdoor_test.go:455`)
and a separate test proves `--ref` is omitted on empty branch.

### 3. vh consolidation — CONFIRMED

- The three absorbed entities each carry `superseded-by:
  survey-skill-correctness-pass` and `sprint-readiness: defer` in frontmatter
  (`survey-codex-cwd-workaround` status `ideation`,
  `survey-scaffold-state-the-fact` status `backlog`,
  `survey-agentsview-detect-under-sandbox` status `ideation`). They are not
  terminalized, but they are the documented consolidation disposition (deferred
  out of the sprint, work folded into vh), not half-done state. The
  `survey-skill-correctness-pass` entity is `done`/`PASSED` and archived
  (`_archive/survey-skill-correctness-pass.md`).
- The SHIPPED survey skill carries all four deliverables:
  - **D1 git-root-basename model** — `SKILL.md:64-81` rewritten to the
    inverse-collision framing; `queries.sql` `:18-31` rationale inverted;
    `folded_keys` DROPPED from the `scoping` SELECT (`queries.sql:46-49` returns
    only `sessions|blank_cwd|span`); cwd-prefix-union kept (`queries.sql:53`).
  - **D2 codex-presence query + hint** — `queries.sql:55-70` (`codex-presence`,
    bound to `:repo_project`), SKILL.md step-2 derivation (`:91`) + step-4 hint
    line (`:154`).
  - **D3 SCAFFOLD state-the-fact** — `SKILL.md:125-129`, `:157`; states
    "recovered from behavior, not files" as a fact, with NO
    recovered/installed/active taxonomy LABELS.
  - **D4 sandbox probe swap** — `SKILL.md:29` is `agentsview --version
    >/dev/null 2>&1`, not `command -v`.

### 4. `go test ./...` green from the repo root (DoD#2) — GREEN

```
$ go test ./...
Go test: 1172 passed in 16 packages
```

(Run from the repo root with `.spacedock-state` and `.worktrees` present — the
condition kb's prune fix targets.) Re-ran the four members' packages explicitly:
`internal/status` ok, `internal/contract` ok, `internal/cli` ok,
`internal/contractlint` ok, `skills/integration` ok.

### 5. Release readiness — CLEAN

- No `TODO` / `FIXME` / `XXX` / `HACK` in any file changed by the four members.
- Every `t.Skip` in the merged tests is a legitimate environment-absent guard
  (codex / agentsview / git / sqlite3 / bash not on PATH, or a Windows-shell
  incompatibility for the `/bin/sh` stub). None hides unproven behavior or
  unfinished work. On this box codex, claude, agentsview, sqlite3, git, and bash
  are all present, so the host-native and sync e2e tests RAN rather than skipped.
- kb / qa / z9 / vh entities are all `status: done`, `verdict: PASSED`, archived,
  with PR refs #327 / #328 / #329 / #330.
- DoD#4 (z9 high-stakes detached adversarial audit) is satisfied on the record:
  the z9 validation ran a detached audit at `0b714fac`, refuted all five
  mandated probes (each reddens the suite then reverts), no Material/Polish
  findings, channel-tracking clean.

## Record for next sprint

These are non-blocking observations seeded for 0.19.9 — none affects the cut.

- **R1 — `TestMigrationCheckPrunesStateTree` re-implements the walk loop rather
  than invoking the production walk.** `migration_check_test.go:213-227` inlines
  a `filepath.Walk` callback instead of calling the same walk closure the
  production check uses (`:65-129`). It DOES exercise the real shared predicate
  (`isMigrationCheckPrunedDir`), which is the load-bearing prune logic, so this
  is not a hole — but the assembled walk body in
  `TestMigrationCheckFixturesParseConsistently` (the `info.IsDir()` →
  `SkipDir` composition) is duplicated, not shared. If that composition ever
  diverges from the production path, the hermetic test would not catch it. Low
  priority: extract the walk-step into a shared helper both call, so the prune
  composition is tested once.

- **R2 — the three superseded survey entities (69/1p/4t) are deferred, not
  archived.** They carry `superseded-by` + `sprint-readiness: defer` but remain
  at top-level in the state checkout with `status: ideation`/`backlog` (not
  terminalized). This is the intended consolidation disposition and is correct
  for a pre-flip cut, but they will linger in `--next`-style queries unless a
  later sweep terminalizes superseded entities. Consider a convention: a
  `superseded-by` entity routes to a terminal `superseded` disposition so it
  drops out of the ready frontier. Not a 0198 concern.

- **R3 — pre-existing instruction-file read outside contractlint (not this
  sprint's).** `internal/ensigncycle/teardown_marker_consistency_test.go:32-36`
  reads `skills/first-officer/references/first-officer-shared-core.md`
  (`TestGradeMarkerMatchesContract`). None of the four 0198 commits touched it,
  so it is out of scope for this audit, but it sits outside the
  `internal/contractlint` quarantine the README defines for instruction-file
  reads. Worth a separate look in a future hardening pass to confirm it is a
  structural/marker check (allowed in spirit) rather than a prose-grep, and to
  decide whether the boundary guard should cover ensigncycle. Flagged, not
  filed against 0198.

- **R4 — `RunDoctor` exits 0 on `NoPluginFound`.** `doctor.go:83-85` maps
  no-plugin-found to exit 0 (a non-fatal report), which is the documented
  contract and is correctly compensated for by `gateHost` inspecting the VERDICT
  rather than the doctor exit code (`frontdoor.go:104-114` comment). Noted only
  so a future reader does not mistake the exit-0 for "compatible" — the gate
  already handles it. No action needed; informational.
