---
id: 008h7wr55c7fn5x3r2wk26yz
title: Ship the Roborev workflow setup skill
status: ideation
source: Captain decision after spacedock-subspace Roborev adoption pilot, 2026-07-13
started: 2026-07-13T16:11:06Z
completed:
verdict:
score:
worktree:
issue:
---

Ship Roborev adoption as a first-party, user-invocable `spacedock:roborev-setup` skill in the main Spacedock plugin. The skill helps a user decide whether Roborev fits a code workflow and configures the Spacedock workflow boundary without making Roborev a dependency of ordinary Spacedock use.

## Reference artifacts

- [Approved Spacedock–Roborev integration proposal](artifacts/roborev-spacedock-integration-proposal.md)
- [Draft `roborev-setup` skill](artifacts/roborev-setup-skill/SKILL.md)
- [Draft skill UI metadata](artifacts/roborev-setup-skill/agents/openai.yaml)

These files preserve the reviewed proposal, the user feedback folded into it, and the validated setup-skill package that seeded this task. Ideation may revise them, but should retain the decision history they capture.

## Problem

Spacedock users can benefit from independent Roborev code-review evidence, but adoption currently requires knowing how to combine Roborev panels, implementation-exit ownership, fresh validation, split state checkouts, daemon placement, and Safehouse access. A separate integration plugin would make the setup entry point harder to discover and would create a one-skill packaging boundary without independent runtime code or release needs.

Skill behavior is harder to test than command behavior. Repository policy bans prose-grep as behavioral proof, requires a live drive for a skill change, and requires a detached adversarial audit for shipped contract or scaffolding changes. The design must therefore produce fixture-backed durable evidence rather than treating skill wording or transcript phrasing as proof.

## Proposed approach

- Add `skills/roborev-setup/SKILL.md` to the main `spacedock` plugin with `name: roborev-setup` and `user-invocable: true`, exposed as `spacedock:roborev-setup`.
- Keep the skill setup-only. Do not load it from the first-officer or ensign core and do not make Roborev mandatory for normal workflows.
- Point users to Roborev's official setup and panel documentation instead of owning a duplicate Roborev configuration schema.
- Configure the Spacedock workflow convention: `quick` is a cheap post-commit cost gate during implementation, the exact-head `code_completion` panel is the authoritative implementation-exit evidence, and fresh validation verifies the stored synthesis-parent evidence before behavioral validation.
- Order the two panels explicitly because Roborev does not express cross-panel prerequisites in `.roborev.toml`: after the final candidate commit, wait for that exact tip's existing `quick` review; fix and recommit any Medium-or-higher finding; let Low findings remain advisory; launch `code_completion` only after `quick` clears the severity floor. Put this invariant in the generated implementation-stage text or dispatch checklist, not only in setup explanation.
- Add the split state checkout's actual branch to `excluded_branches` during setup without enqueueing a probe review.
- Detect Safehouse or another sandbox during setup. When present, advise the user about the smallest read-only runtime access needed for the external daemon; omit sandbox-specific advice otherwise.
- Pilot the behavior against the existing `spacedock-subspace` adoption experience and retain the approved proposal as ideation input.

## Out of scope

- A general `spacedock-integrations` plugin or dedicated Roborev plugin.
- Roborev daemon installation, service management, agent authentication, or an independently versioned adapter.
- Roborev `fix`, `refine`, or agent-hook ownership of code changes or workflow routing.
- A Spacedock lifecycle mod or new stage-exit hook.
- Tests that claim behavioral coverage by matching prose in `SKILL.md`.

## Acceptance criteria

**AC-1 (VALUE) - A user with an eligible commissioned code workflow can invoke `spacedock:roborev-setup` and finish with a workflow that requires exact-head `code_completion` evidence at implementation exit while preserving fresh behavioral validation afterward.**
Verified by: a live-gated, isolated host drive over a temporary commissioned workflow; assertions inspect the resulting workflow README, process exit, state-checkout git history, and clean repository state rather than transcript phrasing. The generated implementation contract must wait for the exact-tip `quick` result before launching `code_completion`, block on Medium-or-higher findings, and allow Low findings without another round.

**AC-2 - The skill is discoverable as a user command in the installed main Spacedock plugin without loading Roborev instructions into normal first-officer or ensign operation.**
Verified by: structural contractlint coverage for valid frontmatter, directory/name agreement, `user-invocable: true`, published skill-surface discovery, and absence from FO/ensign load closure; install smoke coverage resolves `spacedock:roborev-setup` from the current plugin checkout.

**AC-3 - Split workflow state is excluded during setup using its actual branch, and setup never queues a state-only Roborev review as a probe.**
Verified by: the isolated live fixture uses a split state checkout on a non-default branch and a recording fake Roborev executable; assertions inspect the resulting checkout configuration and command journal, fail if a state `review`, `post-commit`, `fix`, or `refine` call occurs, and confirm the state-checkout commit.

**AC-4 - Sandbox advice is conditional and does not commit machine-local permissions or broad Roborev data access into the project.**
Verified by: ideation must define a durable behavior oracle for sandbox-present and sandbox-absent fixture runs. Transcript/prose matching is explicitly insufficient. The oracle must also demonstrate that no project file gains a broad `~/.roborev` mount or machine-local Safehouse configuration.

**AC-5 - The setup skill retains Spacedock's ownership boundary: Roborev supplies review evidence, while implementation owns fixes and Spacedock owns validation and routing.**
Verified by: the live fixture first queues an exact-tip `quick` review with a Medium finding and proves no `code_completion` job starts; implementation fixes and commits, waits for the replacement tip's `quick` result, then launches `code_completion`. A Low-only quick result does not force another round. Durable workflow state remains in implementation until the authoritative panel passes, then fresh validation verifies the stored parent and independently exercises its acceptance check.

**AC-6 - The test suite can refute a materially broken setup rather than merely certify the shipped text.**
Verified by: a detached adversarial audit on a throwaway checkout mutates at least the state-branch exclusion or exact-head requirement and demonstrates that the selected fixture/live oracle turns red, then restores the checkout and records a clean audit or routes findings back to implementation.

## Test plan

- Offline structural checks: extend `internal/contractlint` only for frontmatter validity, user-command discovery, name/path agreement, and load-boundary invariants. These tests do not claim setup behavior.
- Hermetic fixture: create a temporary commissioned code workflow with implementation and fresh-validation stages, a split state checkout on a deliberately non-default branch, and a recording fake `roborev` command. Seed documented Roborev choices so the live prompt does not depend on network access or interactive agent selection.
- Live behavior: add one live-gated setup drive through a supported host, invoking the installed local plugin and granting the reversible setup edits. Assert process exit, resulting workflow files, state-checkout branch/config and git log, fake-command journal, and clean status. Ideation must decide whether one representative host plus shared install smoke is sufficient or whether the user-command surface requires more host lanes.
- Sandbox matrix: exercise sandbox-present and sandbox-absent fixtures only after ideation defines an on-disk or process-state oracle that does not rely on transcript wording.
- Replacement-review journey: use deterministic fake synthesis records to prove failure stays in implementation, a fixing commit invalidates prior evidence, and fresh validation consumes the passing exact-head replacement.
- Panel-order journey: the recording fake Roborev command must prove `wait HEAD` completes before `code_completion` starts, a Medium-or-higher quick finding prevents the expensive panel, the fixing commit causes a new exact-tip quick wait, and a Low finding remains advisory. A negative fixture that starts both panels independently must fail.
- Required repository gates: `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
- Before merge, run the proof-policy detached adversarial audit because this adds shipped skill/scaffolding.

## Ideation questions

- What durable artifact or process-state change proves conditional sandbox advice without turning machine-local Safehouse settings into committed workflow configuration?
- Can the live setup harness drive the user skill directly with a fake Roborev command, or does a small deterministic setup command need to own the testable mutations while the skill owns judgment and explanation?
- Which single host should carry the behavioral live drive, and what install/discovery checks are required to justify host-neutral availability?

## Captain feedback: `quick` is a cost gate

Roborev treats the configured post-commit `quick` panel and an explicitly invoked `code_completion` panel as independent jobs; `.roborev.toml` has no cross-panel prerequisite. The generated workflow must therefore encode the order explicitly: wait on the final candidate commit's existing quick review, route Medium-or-higher findings back through implementation, leave Low findings advisory, and launch the expensive authoritative panel only after the exact-tip quick review clears. `quick` schedules and saves cost; it never substitutes for `code_completion` as implementation-exit evidence.
