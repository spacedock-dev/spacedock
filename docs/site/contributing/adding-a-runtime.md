# Adding a runtime

A host is *supported* when a live or fixture-backed run launches it as a first officer, dispatches an ensign through its native mechanism, and verifies the resulting workflow state. A host is not supported just because its instructions mention Spacedock.

## Add support in layers

Each layer has its own proof:

1. **Skill adapters** — add a first-officer and an ensign runtime adapter for the host, wired from the corresponding skill. The adapter must name the host's native dispatch mechanism, not emulate Claude's tools.
2. **Dispatch host mode** — teach `spacedock dispatch build` the host's assignment shape, keeping entity and worktree paths explicit. Test both the positive shape and the banned-tool negative case.
3. **Runtime contracts and registries** — define the minimum worker record if the host has long-lived workers, and reject stale completion evidence so a previous completion never satisfies a later assignment.
4. **Launch/install UX** — add `spacedock <host>`, `install --host <host>`, and `doctor --host <host>` only after the underlying mechanism is proven.
5. **Live runner** — prove the host with a live-gated test that asserts process exit, entity content, git log, and clean state — never transcript phrasing.

## Manifesting from void

When a runtime seems unsupported on first contact, don't treat setup friction as proof the path is impossible. A missing auth file, an undiscovered extension, or a different subagent tool schema is harness work, not a product blocker. Use the deliberate "assume it works" operating prompt to force the implementation loop to iron out auth, package paths, and tool-shape mismatches before declaring a blocker.

## The full guide

The complete runtime contract — the layer-by-layer acceptance checklist, the "manifesting from void" prompt, and the Pi live-smoke mechanism — is on the [Multi-host support](../reference/multi-host.md) page.
