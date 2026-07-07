---
id: 81hn8vs2fv9wv34wm942r4zj
title: Codex --plugin-dir prevents stale sibling Spacedock skill providers
status: ideation
source: captain request 2026-07-07 after local --plugin-dir session loaded a cached first-officer path
started: 2026-07-07T12:49:53Z
completed:
verdict:
score: 0.6
worktree:
issue:
---

A Codex session launched for local Spacedock development was expected to use `--plugin-dir .`, but the session skill registry surfaced a cached `spacedock:first-officer` path first. Live inspection showed both stable `spacedock@spacedock` and edge `spacedock@spacedock-edge` installed/enabled in Codex, while `spacedock codex --plugin-dir .` only installs the selected channel.

## Problem

Codex has a shared skill namespace. When both Spacedock channels are enabled, `$spacedock:first-officer` may resolve through the wrong or stale provider, undermining the local checkout guarantee that `--plugin-dir .` is supposed to give developers.

## Proposed approach

Ideation should root-cause the actual Codex plugin resolution behavior before selecting a fix. Candidate fixes include removing/disabling the sibling Spacedock channel during local `--plugin-dir` installs, making the launcher emit a stronger diagnostic when duplicate Spacedock providers are enabled, or changing the dev launch path to make the project-local skills authoritative without relying on provider ordering.

## Out of scope

Do not change stable/edge release semantics unless the diagnosis proves channel coexistence is inherently unsafe for Codex. Do not rely on transcript wording alone as proof.

## Acceptance criteria

**AC-1 - A Codex `spacedock codex --plugin-dir <checkout>` launch cannot load stale Spacedock skills from a sibling channel.**
Verified by: fixture-backed or live Codex evidence showing the effective `spacedock:first-officer` skill path resolves to the checkout or to the local-marketplace cache created from that checkout, not an older stable/edge cache.

**AC-2 - Duplicate Spacedock provider state is handled deliberately.**
Verified by: tests over the install/resolve sequence that cover stable-only, edge-only, and both-installed states, with the chosen behavior documented by failing fixtures.

**AC-3 - Developer diagnostics identify the authoritative plugin source.**
Verified by: launcher or install output that names the local checkout/local marketplace source and catches a stale or duplicate-provider condition where applicable.

## Test plan

Start with focused unit tests around `codexInstallArgvSequence`, `WriteCodexLocalMarketplace`, and `execHost.ResolveManifest("codex")`, then add a live or hermetic Codex CLI smoke only if needed to prove actual provider resolution.
