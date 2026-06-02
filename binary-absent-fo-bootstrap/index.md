---
id: 5vb6mh9kewyh0p68r93mf6m1
title: Harden binary-absent FO startup — install hint in the contract-gate abort, don't route to the missing binary's doctor
status: backlog
source: "captain (2026-06-02) — live install-path test (plugin present, binary unlinked): the FO aborts cleanly but its install guidance is model-improvised, and the contract routes the binary-absent case to `spacedock doctor`, which is itself the missing binary"
started:
completed:
verdict:
score: "0.28"
worktree:
issue:
---

When a Spacedock plugin is present but the `spacedock` binary is absent from PATH, the FO contract-version gate (`first-officer-shared-core.md` §Startup step 1) runs `spacedock --version`, gets "command not found", and aborts startup — correctly, no silent failure. But the captain's live `claude -p` test surfaced two weaknesses in what happens next:

- **Binary-absent is routed to `spacedock doctor`.** The gate's abort instruction says "run `spacedock doctor` for the per-class remedy." For the *binary-absent* class that is a dead end — `doctor` is the same missing binary. A strong model worked this out and didn't loop, but a weaker model could loop on `doctor` or just print "command not found." `doctor` is the right remedy only for the *binary-present / version-mismatch* class (where the binary can run).
- **Install guidance is model-improvised, not contract-guaranteed.** `first-officer-shared-core.md` carries zero install guidance. The good behavior the captain saw (concrete `brew` + `go build ./cmd/spacedock` steps) was improvised by the model, not guaranteed by the contract. The contract should carry the install hint so even a weak model emits a correct, runnable remedy.

This is the FO-contract sibling of the `init→install` rename drift + unknown-subcommand-silent-exit code fixes folded into `1x` (same live install-path test). It is scaffolding: it edits the shipped FO contract that governs every dispatched session.

## Design direction (for ideation to flesh out)

- Split the contract-gate abort into two classes: **binary-absent** (`spacedock --version` not found / non-executable) → emit a literal install hint inline (the `brew` cask line + the `go build -o ./spacedock ./cmd/spacedock` source fallback), and do **not** route to `doctor`; **binary-present-but-out-of-range** (version token parsed, outside `>=1,<2`) → keep routing to `spacedock doctor` for the per-class remedy (the binary can run).
- Reconcile with `agents/first-officer.md` if it mirrors the startup gate text (keep the two in sync — the shared-core is the authority).

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — binary-absent emits a runnable install hint, not a doctor bounce.** The contract-gate abort step, for the binary-absent class, contains a literal install command (the `brew` cask install + the `go build -o ./spacedock ./cmd/spacedock` source fallback) and does not instruct the FO to run `spacedock doctor`.

**AC-2 — doctor stays the remedy for the binary-present version-mismatch class.** The version-out-of-range branch still routes to `spacedock doctor`; only the binary-absent branch is changed.

**AC-3 — shared-core and the first-officer agent prompt agree.** Any startup-gate text mirrored in `agents/first-officer.md` matches the shared-core change; no drift between the two.

## Test gates / proof discipline

The deliverable here is the FO's own loaded instructions — the binary is *absent*, so no `spacedock` command can run to emit the hint; the guarantee must live in the contract text the FO already has. This is the legitimate doc-as-deliverable case (the contract prose *is* the product for an agent-read contract), distinct from the prose-only-AC antipattern the workflow forbids (an AC claiming a *behavioral* guarantee that a code gate should enforce). Proof is grep-checkable presence of the install hint + the absence of a doctor-route in the binary-absent branch, plus a read-through that the two-class split is unambiguous. Ideation must record this reasoning so the gate does not bounce AC-1 as a bare wording-is-present criterion.
