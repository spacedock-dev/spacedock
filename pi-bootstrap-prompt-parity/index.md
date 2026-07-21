---
title: Pi bootstrap prompt parity — match claude/codex warmth; "Use $spacedock:first-officer" is the cold outlier
status: validation
source: "Captain (2026-06-20): the pi bootstrap prompt (internal/cli/pi.go:20) is 'Use $spacedock:first-officer for this whole Pi session.' — a bare mechanism trigger. claude (frontdoor.go:25) and codex (frontdoor.go:434) get 'You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage.' (codex appends 'Assume $spacedock:first-officer for the entire session.'). Pi is the cold outlier — pure mechanism, zero warmth. The skill is the contract (single source of truth), but the launch moment is the one chance to frame the commissioning, and pi's reads like a config line."
score:
started: 2026-07-20T23:43:13Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-bootstrap-prompt-parity
issue:
sprint:
sprint-readiness:
id: 7vtn8yda8vn0p7y8am3f43c8
mod-block: merge:pr-merge
---

# Pi bootstrap prompt parity

## End value

`spacedock pi` launches the FO with a bootstrap prompt that matches claude/codex's warmth and commissioning framing — not the current bare mechanism trigger. The launch moment frames the officer's posture (you got this; love your crew; engage) and triggers the skill, instead of reading like a config line. Pi is no longer the cold outlier among the three hosts.

## Problem — root cause already determined

The three bootstrap prompts (verified in source):

- **claude** (`internal/cli/frontdoor.go:25`): `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage."`
- **codex** (`internal/cli/frontdoor.go:434`): `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage. Assume $spacedock:first-officer for the entire session."`
- **pi** (`internal/cli/pi.go:20`): `"Use $spacedock:first-officer for this whole Pi session."`

Pi's prompt is pure mechanism (trigger the skill) with zero warmth or commissioning framing. Claude/codex share a warm core ("You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage.") — codex appends the skill-trigger clause; claude's skill-trigger comes from the `--agent spacedock:first-officer` flag instead. Pi has neither the warm core nor the flag-based trigger, so it falls back to the bare sentence.

## Approach

Bring pi's prompt to parity. Two options:

- **(a) Mirror codex exactly:** `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage. Assume $spacedock:first-officer for the entire session."` — byte-identical to codex (codex already carries the skill-trigger clause, which pi needs since pi has no `--agent` flag). Recommend — true parity, zero drift risk, the warm core + the trigger in one.
- **(b) Pi-specific warm prompt:** a pi-flavored variant of the warm core. More character, but introduces a third prompt to maintain and risks drift. Reject unless the captain wants a distinct pi persona.

Pick (a). The prompt is a single `const` swap in `internal/cli/pi.go:20`; no other change. The skill remains the contract (single source of truth); the prompt is the launch-moment framing + trigger, exactly as on codex.

## Scope

In scope:
- Replace `piBootstrapPrompt`'s value with the codex prompt (the warm core + `Assume $spacedock:first-officer for the entire session.`).

Out of scope:
- Changing the skill or the contract — the prompt is the trigger, not the contract.
- claude/codex prompts — they're already the warm reference.
- The `--agent` vs skill-trigger mechanism difference — pi has no `--agent` flag, so the prompt carries the trigger (as codex does).

## Acceptance criteria (provisional — ideation finalizes; proof = behavior)

**AC-1 — `piBootstrapPrompt` matches codex's prompt (warm core + skill-trigger clause).**
Verified by: a Go test asserting `piBootstrapPrompt == codexBootstrapPrompt` (byte-identical) — OR, if a distinct pi persona is chosen, asserting the warm core is present + the skill-trigger clause is present. (Binding two independent values — the two prompt constants — that can diverge: legitimate structural check, not prose-grep.)

**AC-2 — `spacedock pi` launches with the warm prompt (the smoke / a live launch observes it).**
Verified by: the pi-live smoke or a live launch capturing the appended prompt and confirming the warm core + skill-trigger appear in the launch argv. Behavior-bound, not a prose claim.

## Test plan

- Go test (AC-1): `piBootstrapPrompt == codexBootstrapPrompt` (or the warm-core + trigger-clause assertion).
- Live/harness (AC-2): the smoke's captured launch argv contains the warm core + skill-trigger.
- `pi-live` lane (touches `internal/cli/pi.go` — pi-only surface).

## Related

- `internal/cli/pi.go:20` (`piBootstrapPrompt`) — the source of truth.
- `internal/cli/frontdoor.go:25,434` (`bootstrapPrompt`, `codexBootstrapPrompt`) — the warm reference.
- The 0223 Shaping FO debrief — surfaced the "without love" observation during the fnm-fix smoke.

## Re-scope (2026-06-20) — extension context-hook commissioning

The original approach (option a: swap `piBootstrapPrompt` to mirror codex) landed as a **STOPGAP** — worktree commit `081d6819` on branch `spacedock-ensign/pi-bootstrap-prompt-stopgap` warms `piBootstrapPrompt` to byte-identical to `codexBootstrapPrompt`. That fixes the cold-outlier state immediately (one-line `const` change, build green). It is NOT the architecturally correct fix.

### The real design — pi extension `context`-hook commissioning (the superpowers pattern)

Pi has no `--agent` flag (verified: `pi --help` has none; pi's main-session persona is skills + extensions + system prompt, not named agents). But the pi extension API exposes a `context` hook that can inject a bootstrap message into the session context — exactly the pattern `obra/superpowers` uses (`.pi/extensions/superpowers.ts`):

```ts
pi.on("session_start", () => { injectBootstrap = true; });
pi.on("session_compact", () => { injectBootstrap = true; });
pi.on("agent_end", () => { injectBootstrap = false; });
pi.on("context", (event) => {
  if (!injectBootstrap) return;
  if (event.messages.some(messageContainsBootstrap)) return;
  const bootstrap = getBootstrapContent();  // the warm commissioning + FO contract pointer + tool mapping
  const bootstrapMessage = { role: "user", content: [{ type: "text", text: bootstrap }], timestamp: Date.now() };
  const insertAt = firstNonCompactionSummaryIndex(event.messages);
  return { messages: [...event.messages.slice(0, insertAt), bootstrapMessage, ...event.messages.slice(insertAt)] };
});
```

The injected bootstrap is a rich `<EXTREMELY_IMPORTANT>`-framed message: the warm commissioning ("You totally got this… I love you… Engage.") + the FO contract pointer + a pi tool mapping. It **re-injects on `session_compact`** (survives compaction — the prompt stopgap doesn't) and **suppresses on `agent_end`** (doesn't double-inject).

### Why this supersedes the prompt stopgap

| | prompt stopgap (landed) | extension context-hook (this task's real scope) |
|---|---|---|
| warmth at launch | ✅ (immediate) | ✅ |
| survives compaction | ❌ (prompt is a launch arg, lost on compact) | ✅ (re-injected on `session_compact`) |
| versioned with the package | ❌ (in the binary) | ✅ (in `.pi/extensions/spacedock.ts`, ships with `pi install`) |
| structured content | ❌ (one string) | ✅ (framed message + tool mapping) |
| single source of truth | shared with codex (two constants) | the extension (one place) |

### Revised scope

In scope (this task, now):
- Add a `pi.on("context")` handler to `.pi/extensions/spacedock.ts` (shipped by `eq` #406) that injects a warm FO bootstrap at `session_start` and re-injects on `session_compact`, suppressing on `agent_end` and de-duplicating (the superpowers pattern).
- The bootstrap content: the warm commissioning core + a pointer to load `$spacedock:first-officer` + a pi tool mapping (pi's lowercase tools: `read`/`write`/`edit`/`bash`/`grep`/`find`/`ls`; no `Skill` tool — load skills via `read`; no standard subagent tool — `subagent` from `pi-subagents` if available; no standard task-list tool — plan files / `TODO.md`).
- The launcher prompt (`piBootstrapPrompt`): **demote to a minimal trigger** (or drop if the extension hook makes it redundant). Record the decision — does the extension hook fully replace the launcher prompt, or does the launcher still need a one-line trigger? Ideation confirms.
- Coexistence with the stopgap: the stopgap (commit `081d6819`) lands first as immediate warmth; this task supersedes it. When this task lands, the stopgap's prompt change is either kept (minimal trigger) or reverted (hook fully owns commissioning) — record the decision.

Out of scope:
- Defining a pi "custom agent" / `--agent` flag — pi has no such mechanism; the context hook is the lever.
- Changing the FO skill/contract — the hook injects a *pointer* to the skill, not the contract.
- claude/codex prompts — they're the warm reference; this task brings pi to their warmth via the extension route.

### Revised ACs

- **AC-1** — `.pi/extensions/spacedock.ts` injects a warm FO bootstrap at `session_start` (a live pi session observes the injected message in its first context). Behavior-bound, not prose-grep.
- **AC-2** — the bootstrap re-injects on `session_compact` (a compacted session observes the bootstrap still present after compaction). The property the prompt stopgap can't provide.
- **AC-3** — the bootstrap de-duplicates (doesn't double-inject when already present) and suppresses after `agent_end`.
- **AC-4** — the launcher prompt is either demoted to a minimal trigger or dropped (decision recorded); the extension hook owns the commissioning.

### Relationship to the stopgap

The stopgap (commit `081d6819`, branch `spacedock-ensign/pi-bootstrap-prompt-stopgap`) is an immediate one-line warmth fix that should merge independently of this task — it closes the cold-outlier state now. This task (extension context-hook) is the architecturally correct supersession; it lands later and either keeps or reverts the stopgap's prompt change per AC-4's decision. Do NOT block the stopgap on this task.

### Related (added)

- `.pi/extensions/spacedock.ts` — the extension file (shipped by `eq` #406) that hosts the context hook.
- `obra/superpowers` `.pi/extensions/superpowers.ts` — the reference pattern (session_start/session_compact/agent_end/context hooks; bootstrap injection; de-duplication).
- `spacedock-ensign/pi-bootstrap-prompt-stopgap` branch (commit `081d6819) — the stopgap this task supersedes.

## Spike validation (2026-07-20, FO session — mechanism proven, pitfalls documented)

A throwaway extension (`/tmp/pi-ext-spike/compact-reinject.ts`, driven headlessly over RPC mode) validated every mechanism this task's design depends on. This is spike evidence, not the deliverable; the implementation dispatches fresh from these findings.

**Proven live:**

- `session_start` arms, the `context` hook injects the bootstrap, `agent_end` suppresses (the entity's exact pattern).
- `ctx.compact()` runs `session_before_compact` → compaction completes → `session_compact` fires (`tokensBefore=11363` observed).
- **AC-2's core property holds:** after compaction, the `context` hook re-injects the bootstrap. Verified structurally (`before_provider_request` payload contains the marker) and behaviorally (the post-compaction agent quoted the bootstrap verbatim).
- Extension edits between `--continue` runs took effect immediately (reload-on-resume); `/reload` covers live sessions (docs).

**Implementation pitfalls (each cost a spike iteration — do not rediscover):**

1. **De-dup structurally, never by raw substring.** The compaction *summary text mentions the marker string*, so `messages.some(m => JSON.stringify(m).includes(marker))` suppresses re-injection in exactly the case it is needed. De-dup on message shape instead: role `user`, text starts with `<EXTREMELY_IMPORTANT>`, contains the marker.
2. **Compaction needs `keepRecentTokens` (default 20000) of history** or `prepareCompaction` returns nothing ("Nothing to compact (session too small)"). Live tests must set a small `compaction.keepRecentTokens` in project `.pi/settings.json` and pass `--approve` — non-interactive `-p` ignores project settings without a trust decision.
3. **Headless timing:** a `-p` run exits at turn end and aborts in-flight compaction ("Compaction cancelled"); compacting while the agent streams aborts the active turn ("Request was aborted"). **RPC mode (`--mode rpc`) is the correct live harness** — the process stays alive across prompts so turn-end-triggered compaction completes.
4. **Verify via `before_provider_request`, not model self-report** — the model once answered "no marker" with the marker present in payload (flaky self-inspection; structural check is deterministic).
5. `session_compact` re-trigger right after a completed compaction hits "Already compacted"/"session too small" — the test driver must add history between compactions.

**Test assets from the spike (reusable as the live-test scaffold):** the extension pattern + an RPC driver script (`spawn pi -a --mode rpc -e <ext> -c`, JSONL prompts, wait on `agent_end` events).

## Dispatch scope (captain-directed fast track, 2026-07-21)

Captain: "dispatch the parity so we don't need to hand inject this extension" — ideation rode this dispatch (design spike-validated above; same precedent as fix-pi-live-lane).

- Deliverable: extend the repo's `.pi/extensions/spacedock.ts` (currently `resources_discover` only) with the designed behavior — arm on `session_start`, re-arm on `session_compact`, suppress on `agent_end`, inject via the `context` hook, structural de-dup, insert after the leading compaction summary.
- Reference implementation (dogfooded, verified live in this FO session): `~/.pi/agent/extensions/spacedock-compact-reinject.ts`. Spike artifacts: `/tmp/pi-ext-spike/` (`compact-reinject.ts`, `rpc-driver2.mjs`, `injection.log`). Live-cycle proof: `/tmp/spacedock-compact-reinject.log`.
- Bootstrap prompt mirrors the shipped FO contract per the closed questions; the stopgap text is a starting point, not gospel.
- Test via the RPC harness approach (pitfalls above) — never `pi -p` for compaction, never compact mid-stream. Do NOT grep instruction prose in tests (same trap class as the wiped workflow-guard prose-greps); verify extension effects.
- Hand-over-hand: stage report must note that the captain deletes `~/.pi/agent/extensions/spacedock-compact-reinject.ts` once the shipped version is verified loaded.

## Stage Report: implementation

- DONE: The shipped extension re-injects the FO bootstrap after a real compaction, proven by execution evidence (RPC harness before_provider_request payload or equivalent), not by reading the code.
  Evidence: live `pi --approve --mode rpc -e .pi/extensions/spacedock.ts -e /tmp/spacedock-pi-proof-probe.ts` completed compaction (`session_before_compact`, `session_compact`, `compact_complete tokensBefore=28069`) and the next `before_provider_request` logged `marker=true bootstrapCount=1`.
- DONE: De-dup is structural and proven: when the compaction summary mentions the marker text, exactly one bootstrap exists in context — no double-inject, no false-skip.
  Evidence: `go test ./internal/piruntime -run TestSpacedockPiExtensionBootstrapBehavior` fails if summary-marker prose suppresses reinjection, if structural duplicate messages double-inject, or if `agent_end` fails to suppress; live proof also logged `markerOccurrences=3` but `bootstrapCount=1` after compact.
- DONE: Extension loads clean alongside the repo's existing extensions via /reload, and the bootstrap text mirrors the shipped first-officer contract rather than a hand-tuned variant.
  Evidence: RPC `/reload` driver loaded `.pi/extensions/spacedock.ts` with `sawReload=true` and empty stderr; bootstrap text uses `SPACEDOCK-FO-BOOTSTRAP-v1` contract text with the Pi tool mapping from the shipped FO bootstrap.

### Summary

Implemented `.pi/extensions/spacedock.ts` context-hook commissioning: `session_start` and `session_compact` arm FO bootstrap injection, `agent_end` suppresses it, and structural de-dup avoids false skips when compaction summaries mention the marker. Added a Node-backed Go behavior test for the extension and left `piBootstrapPrompt` as the minimal launcher trigger while the extension owns durable commissioning; captain should delete `~/.pi/agent/extensions/spacedock-compact-reinject.ts` after verifying the shipped extension is loaded. Code commit: `fed85b03`; validation gates run: `go test ./...`, `go test ./... -race`, focused extension behavior test, live RPC compaction proof, and RPC `/reload` smoke.

## Stage Report: validation

- DONE: Independently re-proves the core claim with a fresh live run (RPC compaction → bootstrap present post-compact, exactly once), rather than citing the implementation report's evidence.
  Evidence: fresh `node /tmp/spacedock-validation-rpc-driver.mjs` drove `go run ./cmd/spacedock pi --plugin-dir <worktree> -- --approve --mode rpc`; probe log shows `session_compact`, `compact_complete tokensBefore=22000`, then next `before_provider_request structuralBootstrapCount=1 markerOccurrences=2` (one structural bootstrap, marker also in summary).
- DONE: Attacks the structural de-dup for edge cases: false-skip on summary variants, double-inject with a pre-existing bootstrap, suppression failure after agent_end.
  Evidence: fresh `node /tmp/spacedock-validation-edge-harness.mjs` fails if compaction-summary marker prose suppresses reinjection, if pre-existing structural bootstrap double-injects, or if `agent_end` does not suppress.
- DONE: Confirms launcher trigger and extension hook coexist without double-injection in a real session (AC-4 coherence).
  Evidence: fresh `go run ./cmd/spacedock pi --plugin-dir <worktree> -- --approve --print --no-tools --no-session -e /tmp/spacedock-validation-pi-bootstrap-probe.ts` logged first provider request with `structuralBootstrapCount=1 markerOccurrences=1 firstLauncherMsg=1`.

### Summary

Validation PASSED. The shipped extension behavior was re-proven with a fresh live Pi RPC compaction run and a separate real `spacedock pi --print` launch-coexistence run; repository gates `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` also completed cleanly. Deferred note: an adversarial phrase outside the observed/canonical "compaction summary" wording (`compacted summary`) would insert before rather than after that summary, but still injects exactly once and does not violate the current promised Pi compaction-summary path.
