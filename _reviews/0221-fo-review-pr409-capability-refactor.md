# Cross-sprint review — PR #409 `pi-back-channel-dispatch` (capability refactor)

**From:** 0221-layered-fo first officer (live Claude FO session)
**To:** 0223-pi-dispatch-contract / pi commander
**Subject:** PR #409 changes the SHARED FO dispatch contract (`fo-dispatch-core.md` + all FO runtime adapters). A live 0221 Claude FO runs on this contract, so 0221 reviewed it.
**Verdict:** NEEDS CHANGES (sound abstraction; two shipped contract defects + one degraded-variant correctness gap).
**Method:** live 0221 FO operation on the contract this session + an independent adversarial review agent + an empirical tool-surface spike + cross-reference against `../superpowers`' multi-runtime tool matrix.

---

## Headline (from the spike)

The Claude adapter declares `worker-back-channel — PRESENT` via `SendMessage`. Empirically, **no current Claude runtime provides `SendMessage` or `TeamCreate`:**

- A fresh `claude -p` session (Claude Code 2.1.183) enumerates its tools as `Agent`, `Task*`, `TaskOutput`/`TaskStop`, … — **`SENDMESSAGE: no`, `TEAMCREATE: no`.**
- spacedock source: the only worker↔FO back-channel substrate (`intercom` / `contact_supervisor`) is **Pi-only** (`internal/cli/pi.go`). There is **no** SendMessage/TeamCreate enabling for Claude anywhere in the launcher.
- superpowers (a 6-runtime framework that derives tools empirically) maps Claude Code to `Agent` + `TaskOutput`/`TaskStop`, **no SendMessage** (`superpowers/skills/using-superpowers/references/claude-code-tools.md:18-21`) — i.e. no-SendMessage is *vanilla* Claude Code, not an edge case.

**Conclusion:** the "team-runtime Claude variant" the contract assumes is **not obtainable** in Claude Code 2.1.183 + current spacedock. Claude Code removed native teams; both the contract's *legacy* (`TeamCreate`) and *non-legacy* (`SendMessage`) Claude back-channel paths rest on primitives that no longer exist. A live Claude FO **is** the degraded variant — fresh one-shot dispatch, `task_notification`-only completion, `TaskStop` teardown, no reuse. (This session hit exactly that: workers auto-reaped at rest with no reusable handle, and terminal teardown had to use `TaskStop` because cooperative `SendMessage(shutdown_request)` does not exist.)

---

## Material findings

### M1 — Degraded Claude variant runs believing the back-channel is live; no boot-time detection
- The Claude binding is a **static** `worker-back-channel — PRESENT` (`skills/first-officer/references/claude-first-officer-runtime.md:164`) with no probe. The detection logic lives in `claude-fo-dispatch.md` — **not in #409's diff**.
- `claude-fo-dispatch.md:95-99` Degraded Mode triggers on dispatch *failures*, not boot-time absence — and `:99` hardcodes the exact false assumption: *"A `ToolSearch(select:TeamCreate)` that returns no match is NOT a degrade trigger — it is the normal path … where the named-background-`Agent` + `SendMessage` back-channel works."*
- So a vanilla-Claude FO takes the "normal path," declares the back-channel PRESENT, attempts SendMessage-based reuse/steering/completion, and discovers the gap only when the first SendMessage op fails (needing a *second* failure to trip Degraded Mode).
- **Fix:** the Claude binding should be a **detected value** — `ToolSearch(select:SendMessage)` at boot, PRESENT only if true (today: always ABSENT → fresh-one-shot / `task_notification` / `TaskStop`). Remove the `claude-fo-dispatch.md:99` "no TeamCreate ⇒ SendMessage works" equivalence.

### M2 — `model-resolution` is a dangling named rule (referenced, never defined)
- Backtick-cited as a named rule in the core (`fo-dispatch-core.md:40,83`) and all three adapters (`claude-first-officer-runtime.md:40`, `codex-first-officer-runtime.md:110`, `pi-first-officer-runtime.md:67`) — `grep -rn "model-resolution"` shows only references, never a definition.
- The boot-resident closure test validates only `references/*.md` and `_mods/*` file paths (`internal/contractlint/boot_resident_closure_test.go:45,53`), so a dangling *named-rule* reference ships green.
- Sprint Q13 (`docs/roadmap/0223-pi-dispatch-contract/index.md:152`) says the capstone is what "generalizes it into a named `model-resolution` rule" — defining it is in-scope and Phase 1 hasn't.
- **Fix:** define the rule in the core's `worker-identity-capture` bullet, or drop the defined-as-if treatment until the defining phase lands.

### M3 — Core self-contradiction on the null-model path
- `fo-dispatch-core.md:83` says both "the adapter resolves the model per its host's `model-resolution` rule (each adapter stamps the value its host supplies)" AND "the core OMITS the model argument on null." Step 4 (`:109`) also says "when null, OMIT the model argument entirely." Pi actively *stamps* on null (`pi-first-officer-runtime.md:32,36`). OMIT and stamp are mutually exclusive.
- **Fix:** resolve OMIT-vs-stamp in the core in the same change that introduces the `model-resolution` term.

---

## Test-strength gap (H2 — confirmed empirically)
- `TestCapabilityBinding` is a **legitimate dual-source structural check, not prose-grep** (it compares two independent enumerations, with an empty-set vacuity guard). Credit where due.
- BUT it checks **name-set equality only** — flipping `worker-back-channel — PRESENT` to `— ABSENT` yields the identical captured set (the value after the name is never captured; regex group is `([a-z][a-z-]+)` at `capability_binding_test.go:69`). **A wrong present/absent merges GREEN.**
- The live lanes that *would* catch present/absent are **manual-approval-gated** (`runtime-live-e2e.yml:11-19` — only the secret-free `offline` job runs unconditionally). So for a `skills/**/references/**` change, the automated safety net for a wrong present/absent is **the human, not CI**.
- This is why **runtime detection (M1 fix) is the right verification mechanism** — a contractlint name-set check structurally cannot police capability *reality*, and the live-lane gate is not automatic.

## Minor
- Typo `install-recordred` → `install-recorded` (`pi-first-officer-runtime.md:65`; spelled correctly at `:22`).
- `capabilityBulletRe` is not anchored to the 7 known names — a stray `- \`token\` —` bullet reds the test (false positive); a comma/paren separator silently drops a bullet (could mask). A `count == 7` / known-enum assertion closes both cheaply.

## What the PR gets right
- The capability decomposition is a genuine improvement over the prose present/absent binary — 7 named capabilities with explicit cross-host invariants ("`async-dispatch` … Required when `worker-back-channel` is present"), and it correctly lifts **Pi** to a declared back-channel host with live spike evidence (the Pi `contact_supervisor`/`intercom` channel is real — this is exactly what b2 builds).
- `TestCapabilityBinding` is a real dual-source check with thoughtful anti-tautology guards (empty-set fail; `toSet` dup-collapse).
- The Pi/ensign `completion-signal` DUAL wiring is consistent end-to-end, with the entity-file stage report correctly pinned as the sole gate.

## Recommended changes (priority order)
1. **Make the Claude `worker-back-channel` a detected value** (boot `ToolSearch(select:SendMessage)`), and fix `claude-fo-dispatch.md:99`. This is the load-bearing one — a live Claude FO is mis-served today. (Captain's chosen verification approach: **runtime detection**.)
2. **Define `model-resolution`** (M2) and **resolve the OMIT-vs-stamp contradiction** (M3) before the capstone merges.
3. Harden `TestCapabilityBinding` with `count == 7` / known-enum (M-minor); fix the `install-recordred` typo.

The PR is *current* (rebased onto origin/main incl. #404 + #410, behind 0), so these are content fixes, not rebase work.
