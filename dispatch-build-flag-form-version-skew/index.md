---
id: jhazq8v9prphgs0bk5pn0xfw
title: Dispatch-build flag/file runtime docs must not outpace accepted binaries
status: ideation
source: "FO dogfood (2026-06-06) - plugin 0.19.5 Codex runtime instructs `dispatch build --checklist-file`, but installed spacedock 0.19.4 still accepts only stdin JSON while satisfying contract 1."
score: "0.27"
worktree: ""
issue:
sprint: 0221-layered-fo
group: binary-ux
sprint-readiness: defer
started: 2026-06-08T15:29:12Z
---

The Codex first-officer runtime now instructs the FO to build dispatch prompts
with the flag/file form:

```bash
spacedock dispatch build --workflow-dir DIR --entity-path FILE --stage STAGE --checklist-file FILE
```

That works in current source, but the installed `spacedock 0.19.4 (contract 1)`
still treats the request as stdin JSON and fails with `invalid JSON on stdin`.
The startup contract gate accepts both binaries because they both advertise
`contract 1`, so the skill can put an FO into a command form the accepted binary
does not support.

This is a runtime/source skew problem, not a Codex host-dispatch leak. Current
source with `--host codex` emits a Codex-shaped `Read /tmp/...` prompt; the old
binary simply does not support the flag/file input mode.

## Spike findings (riskiest path, exercised first)

The riskiest unknown: does a real older accepted binary reject the flag/file
dispatch-build form *while still advertising contract 1* — and, if so, does any
seam let the gate distinguish the two surfaces? Both questions were exercised
end-to-end against a binary built from the `v0.19.4` tag and a binary built from
current source, over an identical minimal workflow fixture.

**The skew is real and the contract token cannot see it.** Both binaries print
`(contract 1)` from `--version`. The same complete flag/file command —
`dispatch build --workflow-dir DIR --entity-path FILE --stage STAGE
--checklist-file FILE --team-name T` — produces:

- current source: exit 0, valid dispatch JSON envelope on stdout;
- `v0.19.4` binary: exit 1, `error: invalid JSON on stdin: unexpected end of
  JSON input`.

The old binary ignores the flags entirely, falls through to its stdin-JSON-only
`runBuild`, and fails on empty stdin — exactly the reported symptom. The startup
contract gate (`internal/contract`) compares the binary's `contract N` token
against the plugin manifest's `requires-contract` range; with both binaries at
`contract 1` and every shipped manifest declaring `>=1,<2`, the gate classifies
both as `compatible`. The gate is blind to this skew because the surface changed
without a contract bump.

**Root cause, confirmed in history.** The flag/file form landed in `cfa3b671`
("dispatch build flag input and host resolution", 2026-06-03), *after* the
`v0.19.4` tag (2026-06-02). `CONTRACT_VERSION` (`internal/contract/contract.go`)
has never changed from `1` (empty `git log -p` history on that constant). A new
observable command surface — now the MANDATORY FO dispatch form per
`skills/first-officer/references/claude-first-officer-runtime.md` — shipped
without the contract bump the constant's own doc comment requires ("Bump it only
when a change to the binary alters the observable surface the FO/ensign contracts
call").

**Detection seam probed.** `dispatch build --print-schema` is a clean
feature-detection seam: current source exits 0 emitting a JSON schema; `v0.19.4`
exits 2 with `error: dispatch build requires --workflow-dir` (it has no
`--print-schema` branch). `dispatch build --help` also splits by exit code
(current 0 / old 2), but its help text never mentions the flag form, making it a
weak schema oracle. The `contract N` token is the natural discriminator and the
one signal the gate already consumes — it just was never moved.

## Decision: bump the contract, let the existing gate reject the old binary

Of the three candidate mechanisms (raise the compatibility gate; feature-detect
and fall back to stdin JSON; clear upgrade-abort before dispatch), the spike
points at **raising the gate**, which delivers the upgrade-abort outcome through
machinery that already exists and is already tested:

- Bump `CONTRACT_VERSION` from `1` to `2` — the flag/file form is a genuine
  change to the observable surface the FO contract calls, which is precisely the
  documented trigger for a bump.
- Bump the shipped manifest ranges (`.claude-plugin/plugin.json`,
  `.codex-plugin/plugin.json`) from `>=1,<2` to `>=2,<3`, plus the vendored
  fixture `internal/contract/testdata/plugin.json`.

With contract 2, the old `contract 1` binary trips the existing `too-old-binary`
verdict at startup (`TestStartupGateAbortsBeforeDiscover` already proves an
out-of-range contract aborts before any discover/dispatch fires), surfacing the
pinned remedy that names the rebuild/upgrade command. No dispatch is attempted
against a binary that cannot honor the form.

The feature-detection-fallback mechanism is rejected: it would perpetuate two
command surfaces indefinitely, and its only natural home is FO-runtime prose
("probe `--print-schema`, else use stdin JSON") — a check over a file the model
reads, which the ideation gate bans as proof. The flag/file form is already
mandatory; the correct direction is to require the binary that supports it, not
to keep supporting the binary that does not.

## Acceptance criteria

**AC-1 - The binary's contract version reflects the dispatch-build flag/file
surface.** `CONTRACT_VERSION` is `2`, and the `dispatch build` flag/file form
(`--entity-path`/`--stage`/`--checklist-file`) is supported at that version.
*Verified by:* `TestBuildSchemaAndValidateOnly`/`print-schema` and the flag-form
build tests pass at `CONTRACT_VERSION == 2` — the value comes from the constant,
the support from running the command; independent sources.

**AC-2 - The startup gate rejects a contract-1 binary before any dispatch.** A
binary advertising `contract 1` against the shipped `>=2,<3` range yields the
`too-old-binary` verdict with the pinned rebuild remedy, and no discover/dispatch
runs. *Verified by:* `TestStartupGateAbortsBeforeDiscover` extended with a
`stubContract: "1"`, `embeddedRange: ">=2,<3"` case asserting abort +
`too-old-binary` + zero discover calls — the verdict comes from `Compare`, the
range from the manifest, neither from FO prose.

**AC-3 - The shipped manifests and binary contract cannot drift apart.** Every
shipped and vendored `requires-contract` range brackets `CONTRACT_VERSION`.
*Verified by:* the existing `TestVendoredFixtureBracketsContractVersion`,
`skills/integration/plugin_manifest_test.go`, and
`skills/integration/marketplace_manifest_test.go` — each parses a real manifest
file (independent of the binary) and asserts it brackets the `CONTRACT_VERSION`
constant; they fail until every range is bumped in lockstep with the constant.

(The prior AC-3 — "Codex host dispatch remains Codex-shaped" — is dropped. The
019x audit confirmed it is already covered verbatim by
`internal/dispatch/build_codex_host_test.go:45,53`, which asserts the emitted
prompt is the `Read … treat its content as your assignment` form, contains no
`Skill(skill=`, and that the dispatch body omits `Skill(skill="spacedock:ensign")`
and `SendMessage(to="team-lead"`. Re-stating it here adds no checkable change.)

## Test plan

- **Gate-rejection regression (AC-2), the load-bearing test.** Add the
  `contract 1` vs `>=2,<3` case to `TestStartupGateAbortsBeforeDiscover`
  (`internal/contract/gate_test.go`) — a real stub binary prints `contract 1`,
  the gate aborts with `too-old-binary`, discover is never invoked. This is the
  one test that, before the contract bump, lets the skew through; it binds the
  minimum supported command surface to a startup abort. Fixture cost: low
  (reuses the existing version-stub harness). Go unit/behavior test.
- **Contract/manifest lockstep (AC-1, AC-3).** No new test needed — bumping
  `CONTRACT_VERSION` to `2` is caught by the three existing bracketing tests
  until the three manifest ranges are bumped to `>=2,<3`. Run
  `go test ./internal/contract/... ./skills/integration/...` to confirm they go
  red on the bump, then green on the manifest edits. Zero new fixtures.
- **Full sweep (validation).** `go test ./...` from the repo root catches any
  other site that hard-codes `>=1,<2` or `contract 1` (e.g. doctor/frontdoor
  fixtures). A live workflow smoke test is not required — runtime behavior is the
  startup abort, fully exercised by the gate test above.

## Stage test gates

- Implementation: bump `CONTRACT_VERSION` to `2` and the three manifest ranges
  to `>=2,<3`; add the `contract 1` rejection case to the gate test; drop the
  old AC-3 Codex-shape duplication. The gate-rejection test is the riskiest path
  and is written first.
- Validation: run `go test ./internal/contract/... ./skills/integration/...`
  focused, then `go test ./...`.

## Stage Report: ideation

- DONE: SPIKE FIRST (riskiest, gates the design): obtain a real older accepted binary (build from the v0.19.4 tag, or brew) and confirm it rejects the flag/file dispatch-build form while still advertising contract 1 — and probe the detection seam (does `--print-schema` / `--help` / a version compare distinguish the supported command surface?).
  Built v0.19.4 + current binaries; identical complete flag/file command → old exits 1 `invalid JSON on stdin`, current exits 0 valid JSON; both print `(contract 1)`. `--print-schema` is a clean seam (current exit 0 schema / old exit 2); `--help` is weak (no flag mention); contract token is identical so the gate is blind. Recorded in "Spike findings".
- DONE: Decide the mechanism (from the spike): pick among raising the binary/plugin compatibility gate, feature-detection/fallback to stdin JSON, or a clear upgrade-abort before dispatch.
  Chose raising the gate (bump `CONTRACT_VERSION` 1→2 + manifest ranges to `>=2,<3`), which delivers the upgrade-abort via the already-tested startup gate; fallback rejected (perpetuates two surfaces; only home is FO prose, a banned check). See "Decision".
- DONE: Produce build-ready ACs + test plan: a version-compatibility test that fails if shipped runtime text requires a flag form the minimum accepted binary rejects. Drop or rewrite AC-3 (the 019x audit found "Codex prompt stays Read /tmp/..." already covered by build_codex_host_test.go:45,53 — tautological as written).
  ACs reframed around the contract bump + existing manifest-bracketing/gate tests (independent-source, not prose grep); load-bearing test is the `contract 1` vs `>=2,<3` rejection case in `TestStartupGateAbortsBeforeDiscover`. Old AC-3 dropped with rationale citing build_codex_host_test.go:45,53.

### Summary

The spike confirmed the version skew end-to-end: a v0.19.4-built binary and current source both advertise `contract 1`, yet the same flag/file `dispatch build` command fails on the old binary (`invalid JSON on stdin`) and succeeds on current — so the startup gate, which compares only the contract token, cannot distinguish them. Root cause confirmed in history: the flag/file form landed after v0.19.4 (`cfa3b671`) without the `CONTRACT_VERSION` bump its own doc comment mandates. The chosen mechanism bumps `CONTRACT_VERSION` 1→2 and the manifest ranges to `>=2,<3`, making the existing (already-tested) startup gate reject the old contract-1 binary with the pinned `too-old-binary` remedy before any dispatch; the load-bearing new test is a `contract 1` rejection case in the gate test, and the existing manifest-bracketing tests enforce the lockstep bump.
