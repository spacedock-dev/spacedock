---
id: 7yd3mbsy2am5qggc17sxvz2v
title: docs/releasing.md step 3 (manual pre-stamp) is stale vs the tag-the-gated-commit practice
status: backlog
source: "0202 Commander drive (2026-06-13). releasing.md 'Cutting a Stable Release' step 3 says to stamp+commit the version before tagging, but v0.20.1 and v0.20.2 both tagged the gated commit directly and let release.yml stamp post-tag (a pre-stamp creates an ungated commit the exact-SHA e2e-gate blocks on)."
group: cleanup
---

`docs/releasing.md`'s "Cutting a Stable Release" procedure (step 3) documents a manual pre-stamp commit before the annotated tag. Actual practice (v0.20.1, v0.20.2) tags the gated commit directly; release.yml stamps the plugin manifests post-tag.

## Problem

The `e2e-gate` job binds the cut to a green Runtime Live E2E run for the EXACT tagged commit SHA. A manual pre-stamp (releasing.md step 3) creates a NEW commit with no green run, which the gate would block. So both recent cuts tagged the already-gated `main` HEAD directly and relied on release.yml's post-tag stamp + the moving `stable` ref. The written procedure and the real procedure diverge — a fresh cutter following step 3 literally would stamp, create an ungated commit, tag it, and be blocked by the e2e-gate.

## Proposed approach (ideation to firm)

Reconcile `docs/releasing.md` step 3 with the tag-the-gated-commit practice: tag `main` HEAD (the gated commit) directly with the changelog; document that release.yml stamps the manifests and advances `stable` post-tag; drop or correct the manual pre-stamp + push-to-main step. Cross-check against `docs/site/contributing/releasing.md` (the site mirror).

## Acceptance criteria (sketch)

**AC-1 (sketch) — releasing.md describes tagging the gated commit directly, with post-tag stamping by release.yml.**
Verified by: the reconciled doc; if any release-machinery test asserts the flow, it pins the gated-commit-tag invariant. (Doc-accuracy change; the behavioral truth is release.yml's existing e2e-gate + stamp steps.)

## Notes
Pure doc reconciliation; the machinery is already correct. Surfaced by the 0202 cut.
