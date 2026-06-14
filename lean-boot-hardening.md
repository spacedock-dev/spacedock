---
id: 58q4bynqqxd3dzjpyntz8m8w
title: Lean boot hardening — FO must report-and-stop on zero `--discover`, not broad-search the filesystem
status: ideation
source: "captain (2026-06-14) — an FO instance overstepped Startup step 3: after `spacedock status --discover` returned zero (exit 0, no output), it ran a broad find/grep filesystem sweep to hunt a workflow instead of reporting no-workflow-found and stopping. Contract + lean-boot violation."
started: 2026-06-14T19:16:23Z
completed:
verdict:
score: "0.30"
worktree:
issue:
sprint: 0203-fo-efficiency
---

Keep FO boot lean: when `spacedock status --discover` returns zero workflows, the Startup discovery step must report no workflow found and STOP — never fall back to a broad `find`/`grep` filesystem sweep to hunt one down. Harden the discipline so the zero-`--discover` path is provably report-and-stop, not an expensive search.

## Problem

- Startup step 3 is explicit: `status --discover` → one path → use it; zero → report no workflow found; multiple → present the list. The zero branch is terminal.
- Observed (captain, 2026-06-14): an FO, after `--discover` returned zero, ran a broad filesystem sweep to locate a workflow — violating the contract's zero-branch AND the lean-boot ethos (cf. j9 shallow-boot: boot is cheap, it does not sweep the filesystem).
- A broad filesystem search at boot is both a discipline violation (the zero-branch is report-and-stop) and a cost/latency regression — the opposite of lean boot.

## Out of scope

{Ideation fills — e.g. whether to scope strictly to the `--discover`-zero path or also guard other broad-search-at-boot temptations.}

## Acceptance criteria

{Ideation fills. PROOF POLICY: a contract clause that says "do not broad-search" is a prose rule with a wording-present ceiling — it does NOT prove the FO obeys it (a paraphrase fails the grep, an inverted clause passes it). The real proof is behavioral (a live drive observing the FO report-and-stop on a zero-`--discover` boot, taking NO broad filesystem search) or a code-level gate — never a string/regex match over the contract asserting the clause exists.}

## Test plan

{Ideation fills.}
