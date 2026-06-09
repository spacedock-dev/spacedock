# Adding a runtime

A host is supported when a live or fixture-backed run launches it as a first officer, dispatches an ensign through that host's native agent mechanism, and verifies durable workflow state: process exit, entity body, state-checkout git log, and clean status. A host is not supported because its instructions mention Spacedock, and a substring search over code or prose is not proof of behavior. Spacedock ships `spacedock claude` and `spacedock codex` as proven front doors; adding a host means earning the same level of proof.

This page is the contributor's orientation. The full procedure, the exact checklists, and the worked Pi example live in [Multi-host support](../reference/multi-host.md). Read that before you write code.

## What "supported" means

The supported claim has four observable parts. The runtime must be able to:

- **Launch** the host as the first officer.
- **Delegate** to an ensign through the host's own subagent or team mechanism.
- **Observe** that the ensign completed.
- **Verify** durable state: process exit, entity body, state-checkout git log, and clean status.

If you cannot demonstrate all four against a real or fixture-backed run, the host is a spike, not a runtime.

## Add support in layers

Add support in small layers, each with its own proof at its own abstraction level. The order matters: every later layer assumes the earlier one is proven.

1. **Skill adapters.** Add `skills/first-officer/references/<host>-first-officer-runtime.md` and `skills/ensign/references/<host>-ensign-runtime.md`, and wire both from the matching `SKILL.md` runtime-adapter section. The adapter must name the host's native mechanism. Do not emulate Claude `Agent`, `SendMessage`, `TeamCreate`, or `TeamDelete` unless the host actually provides those tools.
2. **Dispatch host mode.** Teach `spacedock dispatch build` to accept `host: "<host>"` when the assignment shape differs by host. Keep entity and worktree paths explicit, especially for split-root workflows (`state: .spacedock-state`). Test the positive shape and the banned-tool negative case.
3. **Runtime contracts and registries.** For long-lived workers, define the minimum worker record (label, substrate, run/session handle, entity, stage, state, completion epoch) and reject stale completion evidence: a previous completion must never satisfy a later assignment. For a host team API, adapt Spacedock lifecycle intents to the host's native action schema.
4. **Launch/install UX.** Add `spacedock <host>` only after a manual or live harness proves the runtime path. Add `spacedock install --host <host>` only when the install path is known and checkable without mutating unrelated global host state. Add `spacedock doctor --host <host>` when there is a manifest, package, or runtime health check to verify.
5. **Live runner.** Prove the host with a live-gated test when the claim is runtime integration. Use a temp workflow fixture, isolated host config and session directories, and copied credentials rather than global host state. Assert process exit, entity content, git log, and clean status. Never pass on transcript phrasing.

## Match the proof to the claim

Use the smallest proof at the same abstraction level as the claim:

- **Text claim** (an adapter mentions the right tool): parse or inspect the real instruction files.
- **Dispatch shape claim:** run `spacedock dispatch build` with a fixture and inspect the emitted JSON or body.
- **Adapter claim:** table-test lifecycle intents to exact host-native payloads.
- **Registry claim:** unit-test persistence and stale-epoch rejection.
- **Runtime claim:** a live-gated host run that mutates a temp workflow and verifies durable state.

The failure mode to guard against: declaring a host supported because its prose looks right. Substring presence is acceptable proof only when the claim itself is about text being present or absent.

## Manifesting from void

When a runtime looks unsupported on first contact, do not read setup friction as proof the product path is impossible. A missing `auth.json`, an extension not auto-discovered in a temp home, or a subagent tool schema that differs from Claude's is harness work, not a blocker. A real blocker is a proven inability to launch, delegate, observe completion, or verify durable state *after* the harness is correct.

Run the first implementation and validation loop under a deliberate "assume it works" prompt so the loop irons out auth, package paths, and tool-shape mismatches before anyone declares a wall:

```text
Assume <runtime> support is supposed to work. Do not treat missing polish, auth setup friction, or tool-shape mismatch as proof the runtime is impossible. In first-officer capacity, iron out the frictions:

- if auth is missing in an isolated harness, copy/reuse the existing host auth file correctly;
- if the dispatch substrate needs a local package/extension path, wire it explicitly;
- if the host tool shape differs from Claude/Codex, adapt to the host-native contract rather than emulating Claude tools;
- if a live test fails due to harness setup, fix the harness and rerun;
- only stop for a real product/design blocker, not for first-contact setup friction.
```

The prompt earns its place by changing the default interpretation of a failure: harness work gets fixed in-loop, and only a proven product or design blocker stops the work.

The worked Pi runtime is in [Multi-host support](../reference/multi-host.md): the live smoke mechanism, the exact parent prompt, the install and doctor surface, and the full acceptance checklist.
