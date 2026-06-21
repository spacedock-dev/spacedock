---
title: Convert ensign runtime adapters to capability binding blocks (t0g ensign-side)
status: ideation
source: "Captain (2026-06-21): the same verbosity pattern t0g fixed for FO adapters exists in all three ensign runtime files — ~70-80% shared-core duplication (Agent Surface, worktree rules, split-root, frontmatter, path-scoped commit, file-pointer, feedback routing), ~20-30% genuinely host-specific (clarification tool name, completion signal mechanism, captain communication, shutdown response protocol). The shared ensign core already carries the discipline; the per-host ensign files should become compact binding blocks or be absorbed into the FO adapter's → lines."
score:
started: 2026-06-21T06:05:14Z
completed:
verdict:
worktree:
issue:
sprint: 0230-stable-finalization
sprint-readiness:
id: x1khmz0e80fyhe7vnjg8w59y
---

# Ensign runtime binding block cleanup

## End value

The three ensign runtime adapters (`claude-ensign-runtime.md`, `codex-ensign-runtime.md`, `pi-ensign-runtime.md`) are compact binding blocks carrying only the genuinely host-specific ensign content — clarification tool, completion signal mechanism, captain communication, shutdown response — not shared-core duplication. The shared ensign core (`ensign-shared-core.md`) owns the discipline (worktree, split-root, frontmatter, path-scoped commit, proof, assignment reading); the per-host files own only what differs by host. Same restructure `t0g` applied to the FO adapters, now for the ensign side.

## Problem — same pattern across all three hosts

Each ensign runtime file is ~70-80% duplication of `ensign-shared-core.md`:

### Duplicated (shared core already carries)
- **Agent Surface** — "dispatched via Agent tool" / "through Codex multi-agent" / "from pi-subagents" → shared core: "Read the assignment context provided by the first officer."
- **File-pointer** — "read the named dispatch file" (codex, pi) → shared core: same instruction.
- **Worktree / split-root** — "stay on repo root / keep under worktree / state-checkout path" (pi) → shared core: entire "Working" + "Split-Root State Contract" sections.
- **Frontmatter** — "do not modify YAML frontmatter" (pi) → shared core: "Do NOT modify YAML frontmatter."
- **Path-scoped commit** — "commit path-scoped in state checkout" (pi) → shared core: "Concurrency-safe state commits" + "Multi-writer sync."
- **Feedback interaction** — "FO routes fixes; ensign re-checks" (claude, codex) → shared core / FO-side policy.

### Genuinely host-specific (the only content that should remain)
- **Clarification tool** — `SendMessage(to="team-lead")` (claude) / worker thread (codex) / `contact_supervisor` with `need_decision`/`progress_update` (pi). ~3-8 lines each.
- **Completion signal mechanism** — `SendMessage(to=team-lead, "Done:...")` (claude) / final worker-thread message (codex) / subagent return (pi). ~3-5 lines each.
- **Captain communication** — claude's Shift+Up/Down + codex's conversation text. ~3 lines each. Pi doesn't have this.
- **Shutdown response protocol** — claude's structured JSON `SendMessage` response + codex's plain-text acknowledge. Pi doesn't have this (no mailbox shutdown).
- **Codex's `«context-budget»` / `«addressable-worker»` notes** — a few lines that are really FO-side policy leaking into the ensign; should move to the FO adapter or be dropped.

### Token impact (v0.22.0 → head, from the session debrief)
- `pi-ensign-runtime.md`: 1,768c → 2,838c (+1,070c; grew from the model-stamping + completion-signal additions, but most is still duplication).
- `codex-ensign-runtime.md`: 2,390c → 2,847c (+457c).
- `claude-ensign-runtime.md`: 2,556c → 2,556c (unchanged — not yet touched by any sweep).
- Total ensign adapters: 6,714c → 8,241c. A binding-block restructure should cut each to ~500-800c (just the host-specific bindings), saving ~4,000-5,000c (~1,000-1,250 tokens) across the three files.

## Approach

Mirror `t0g`'s FO-adapter restructure for the ensign adapters:

1. **Audit `ensign-shared-core.md`** — confirm it already carries every discipline the per-host files duplicate (worktree, split-root, frontmatter, path-scoped commit, proof, assignment reading, feedback routing). If any gap exists, fill it in the shared core first (not in the per-host files).

2. **Restructure each ensign adapter to a compact binding block:**
   - Short intro (one line: "How the shared ensign core executes on <host>.")
   - `## Ensign implementation` (or `## Runtime implementation` — match the FO adapter's heading for consistency):
     - `Clarification` → `<host tool>` (SendMessage / worker thread / contact_supervisor)
     - `Completion signal` → `<host mechanism>` (SendMessage Done: / final thread message / subagent return)
     - `Captain communication` → (claude/codex only; pi omits)
     - `Shutdown response` → (claude/codex only; pi omits)
   - Remove all duplicated sections (Agent Surface, Dispatch, Awaiting Completion, worktree rules, split-root, frontmatter, path-scoped commit, feedback interaction).

3. **Consider absorbing into the FO adapter's `→` lines:** the ensign's completion signal + clarification tool are things the FO interprets, not things the ensign controls. If the FO adapter's `«completion-signal»` and `«addressable-worker»` `→` lines already carry these facts, the ensign adapter may not need them at all — the ensign just does what the shared core says, and the FO knows how to interpret the result. Evaluate this; if the ensign adapter can be fully absorbed, eliminate the per-host ensign files entirely and wire the `SKILL.md` runtime-adapter section to the shared core only.

4. **Encode the ensign binding-block shape in `docs/runtime-support.md`** — add the ensign adapter to the "Runtime implementation" section `t0g` added for FO adapters. Same shape: binding block keyed by the host-specific concern (clarification, completion, captain comms, shutdown), not lifecycle re-narration.

5. **Update contractlint guards** — the `TestPiEnsignRuntimeAvoidsNegativeHostContrast` guard (#417) should be updated if the pi ensign file is restructured (the banned phrases may no longer exist; the positive-binding assertions should re-target the new binding-block bullets). Same for any codex ensign guard.

## Scope

In scope:
- All three ensign adapters (claude, codex, pi) restructured to binding blocks (or absorbed into the FO adapter if the ensign-specific content is fully FO-side).
- `ensign-shared-core.md` gaps filled if any duplicated content isn't already there.
- `docs/runtime-support.md` updated to cover the ensign adapter shape.
- Contractlint guards updated.

Out of scope:
- The FO adapters (`t0g` already shipped #418).
- The shared ensign core's discipline (it's the authority — don't trim it).
- New ensign capabilities or contract changes — this is a prose-shape restructure, not a behavior change.

## Acceptance criteria (provisional — ideation finalizes)

**AC-1 — Each ensign adapter is a compact binding block (or absorbed).**
Verified by: a structural review that each file carries only host-specific content (clarification tool, completion signal, captain comms, shutdown); no shared-core duplication. If absorbed into the FO adapter, the ensign files are deleted and `SKILL.md` wires to the shared core only.

**AC-2 — The shared ensign core covers every discipline the per-host files previously duplicated.**
Verified by: a cross-check that `ensign-shared-core.md` carries worktree, split-root, frontmatter, path-scoped commit, proof, assignment reading, feedback routing — all the content the per-host files carried redundantly.

**AC-3 (VALUE / the gate) — absolute `wc -c` <= the v0.22.0 baseline.**
The gate is the absolute byte count against the v0.22.0 baseline (the independent ref that moved the wrong way: codex +457, pi +1070), NOT a before/after reduction: `codex-ensign-runtime.md` <= 2390 B (now 2847) AND `pi-ensign-runtime.md` <= 1768 B (now 2838). `claude-ensign-runtime.md` (2556) and `ensign-shared-core.md` (8829) are AT baseline → net-zero guardrail: they must not grow. Report each file's absolute `wc -c` beside its baseline (`git show v0.22.0:skills/ensign/references/<file> | wc -c`). The ~500-800c per-adapter figure (saving ~4,000-5,000c total) remains the aspirational target, but the pass/fail GATE is the absolute v0.22.0 number above.

**AC-4 — Contractlint guards updated + green.**
Verified by: `go test ./internal/contractlint/...` green; the #417 pi ensign guard (and any codex ensign guard) updated to pin the new binding-block shape.

**AC-5 — Gates green.**
Verified by: `go test ./...`; `gofmt -l ./cmd ./internal`.

**AC-6 — docs/runtime-support.md alignment.**
The ensign binding-block shape this cut adopts is documented in `docs/runtime-support.md`'s Runtime-implementation section (extending the FO-adapter shape t0g added), aligned with WHICHEVER final approach lands — compact per-host binding blocks OR full absorption into the FO adapter (per approach step 3). No drift between the authority doc and the shipped ensign files (or their removal, if absorbed).

## Test plan

- Structural review (AC-1, AC-2): binding blocks carry only host-specific content; shared core covers the rest.
- Token count (AC-3): before/after char/word proxy.
- `go test ./...` (AC-4, AC-5): contractlint + full suite green.
- This is a shipped-contract change (high-stakes surface). Detached adversarial audit at validation.

## Related

- `t0g` `codex-pi-runtime-binding-block-cleanup` (#418, merged) — the FO-adapter restructure this mirrors.
- `docs/runtime-support.md` "Runtime implementation" section (added by `t0g`) — the authority to extend for ensign adapters.
- `ensign-shared-core.md` — the shared discipline (the authority; don't trim).
- `pi-ensign-runtime.md`, `claude-ensign-runtime.md`, `codex-ensign-runtime.md` — the files to restructure.
- `TestPiEnsignRuntimeAvoidsNegativeHostContrast` (#417) — the guard to update.
- `skills/ensign/SKILL.md` — wires the runtime adapter; update if ensign files are absorbed.
