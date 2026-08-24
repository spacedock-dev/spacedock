---
id: 76582d5pd91e33smcekrzm4q
title: Marketplace README documents the retired channel model and omits the separate CLI install
status: backlog
source: "Captain CL fresh-install VM experience report, 2026-08-24: edge install id wrong, next-branch claim wrong, plugin-ships-no-binary surprise"
started:
completed:
verdict:
score:
worktree:
issue:
---

The spacedock-dev/marketplace README still documents the pre-decouple channel model and gives a fresh machine broken instructions. Verified 2026-08-24: (1) it says `claude plugin install spacedock-edge@spacedock` for edge, but the root manifest's entries are [spacedock, subspace, cargento] — no `spacedock-edge` entry exists, so the command fails; correct is `claude plugin marketplace add spacedock-dev/marketplace@edge` then `claude plugin install spacedock@spacedock-edge`. (2) It claims edge tracks the `next` branch (`source.ref: next`); the actual edge-branch manifest has `"ref": "main"` and no `next` branch exists on the product repo. (3) It never states the plugin ships skills/hooks only — the `spacedock` CLI binary is a separate install — which surprised a fresh-VM agent mid-run. Bonus in-repo staleness: `docs/site/get-started/install.md:50` comment says "Edge (tracks next)" (the commands there are correct).

## Problem

{Ideation fills this in. Seeded: the README is the first doc a fresh machine reads; every one of its edge instructions is stale relative to the channel-via-marketplace-name model shipped in marketplace-repo-and-pinned-channels (#352) and spacedock-marketplace-source-env (#424).}

## Proposed approach

{Ideation fills this in. Seeded: rewrite the README channel section to the current model (edge = `@edge`-branch marketplace named `spacedock-edge`, entry `spacedock`, tracking main), add a "the plugin ships no binary — install the CLI separately" note linking the install docs, and fix the stale "(tracks next)" label in docs/site/get-started/install.md.}

## Out of scope

install.sh channel/prerelease behavior (separate task: install-sh-edge-prerelease-parity). HTTPS/SSH plugin-install guidance (separate task: plugin-install-https-keyless-machines).

## Expected surface and tolerance

Estimate net LOC change: ~+30/-20, across 2 files (marketplace repo README, docs/site/get-started/install.md).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Seeded; ideation refines.

**AC-1 - The README's edge install commands, run verbatim on a host without the marketplace added, install the edge plugin.**
Verified by: a live `claude plugin` run of the exact README commands (or the runtime-live-ci host-install smoke shape); fails if the README names a nonexistent entry id.

**AC-2 - No claim in the marketplace README or docs/site says edge tracks `next`; branch claims match the live edge manifest ref.**
Verified by: reading the claims against the fetched edge manifest (`ref: main`); fails if any doc reasserts `next`.

**AC-3 - The README states the plugin ships no binary and points at the CLI install path.**
Verified by: the README containing the note with a working link; fails if a fresh reader has no in-README signal that the CLI is separate.

## Test plan

{Doc change: review against live manifests plus one live plugin-install smoke on a machine/profile without the marketplace pre-added. The marketplace repo change ships as a PR to spacedock-dev/marketplace, precedent fea2/#431 and gp/#352.}
