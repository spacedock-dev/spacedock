---
id: tdpnhct3kqk99e5fj447c1xm
title: Full mdschema conformance validator (status --validate enforces a subset)
status: backlog
source: captain-approved, surfaced by pt0 docs-site port (PR #343, 2026-06-13)
started:
completed:
verdict:
score:
worktree:
issue:
---

`docs/schema/entity.mdschema.yml` + `workflow-readme.mdschema.yml` are now the SSOT for the frontmatter contract (ported from the v0 branch into the v1 repo during the #343 docs pass). But `spacedock status --validate` (`internal/status/validate.go`) enforces only a SUBSET of that schema; nothing checks full conformance against the mdschema files.

## Problem

`status --validate` currently checks: entity-form (flat/folder) conflicts, stage-name regex, per-id-style id presence/uniqueness, and the opt-in external-proof policy. It does NOT check per-field types/patterns (e.g. `verdict` in {PASSED, REJECTED}, `mod-block` pattern, `score` numeric coercion) or the schema invariants end to end. So the mdschema files can drift from what the binary actually enforces, and malformed frontmatter the schema would reject can pass `--validate`.

## Proposed approach

Ideation to weigh:
- (a) extend Go `--validate` to full mdschema coverage (per-field types/patterns + invariants), or
- (b) a standalone conformance checker against `docs/dev`. The v0 branch had `scripts/validate_frontmatter_contract.py` + a test, NOT ported — Python does not belong in the Go repo without a CI hook.

## Out of scope

The #343 doc pass deliberately did schemas + light prose only (captain decision = option (a) territory, deferred). This task is the cleanup that closes the enforce-what-you-document gap.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - {End-state property.}**
Verified by: {test name / command output or exit code / file the change produces / resulting on-disk state — something outside this task body that a future reader can reproduce and that can fail.}

## Test plan

{What verifies the implementation, estimated cost/complexity, and whether fixture, CLI, or live workflow tests are needed. A conformance validator is naturally testable: feed known-bad frontmatter fixtures and assert non-zero exit + the right diagnostic.}
