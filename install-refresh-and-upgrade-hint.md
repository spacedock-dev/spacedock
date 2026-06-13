---
id: tes9th8ncq1p01am9qk7eex4
title: install refresh leaves a stale plugin + no upgrade path is surfaced (the 0.19.8 thing)
status: ideation
source: "Captain field report 2026-06-09 — first real 0.20.0 install. `spacedock install --host codex` on a tag-fresh 0.20.0 binary returned `OK: spacedock binary 0.20.0 and plugin 0.19.8 are compatible.` The plugin stayed at 0.19.8 (older than BOTH main HEAD 0.20.0 and next HEAD 0.19.9), and nothing told the user a newer plugin exists, how to get it, or whether the front door upgrades for them."
started: 2026-06-13T04:08:47Z
completed:
verdict:
score:
worktree:
issue:
---

A user with a tag-fresh **0.20.0 binary** ran `spacedock install --host codex` to upgrade, and ended up still on plugin **0.19.8** with a `compatible` message and no path forward. Two separable facets, both real, surfaced by the first 0.20.0 install.

## Problem

**Facet 1 — the install refresh did not take (correctness).** `Install` (`internal/cli/host_exec.go`) runs a cleanup-then-pin sequence (`plugin remove` / `marketplace remove` — both tolerated — then `marketplace add <source> [--ref <branch>]`, then `plugin add`). After running it against a 0.20.0 binary, the installed codex plugin was still **0.19.8** — older than both `main` HEAD (0.20.0, the stable `source.ref`) and `next` HEAD (0.19.9). So the refresh that was supposed to pull the current branch HEAD did not land. Root cause is unknown and must be investigated, not assumed — candidate hypotheses:
  - a tolerated cleanup step failed, leaving the old plugin in place, and the `plugin add` no-op'd on an already-present entry;
  - the codex marketplace add resolved a stale cache rather than re-fetching the branch HEAD;
  - the binary's resolved install branch/source pointed somewhere other than `main`;
  - `install` short-circuited a reinstall because it judged the versions already "compatible".

**Facet 2 — no upgrade is ever hinted (UX).** The `Compatible` verdict (`internal/contract/contract.go`) compares CONTRACT versions only (binary contract `1` vs plugin `requires-contract: ">=1,<2"`), not display semvers — so "binary 0.20.0, plugin 0.19.8 are compatible" is correct by design. But that is exactly what confuses: it prints two different version numbers side by side and then says "you're fine," with no nudge that a newer plugin (matching the binary) is available, how to get it, or that `spacedock claude`/`codex` could refresh it. A contract-compatible-but-behind plugin should surface an honest, opt-in upgrade hint.

## Proposed direction (ideation firms; may split into two tasks)

- **Facet 1:** root-cause why the refresh left a stale plugin, and make `spacedock install` reliably re-pull the current branch HEAD (or report honestly when it cannot). Reproduce on a fixture/controlled install where the installed plugin is behind the source HEAD and assert the post-install on-disk manifest advances.
- **Facet 2:** when binary and plugin display semvers differ while contract-compatible, surface an opt-in upgrade hint (a newer plugin is available — run `spacedock install --host {host}`), and/or have the front door offer/perform the refresh. Keep the contract-based `compatible` verdict; add the hint as an additional, honest line — never a false "you must upgrade".

## Out of scope

- Changing the contract-compatibility semantics themselves (contract `1` plugin with a contract `1` binary IS compatible — that stays).
- The branch/release-model decision (HEAD-vs-tag serving, trunk direction) — tracked separately in the post-flip roadmap shaping; this task is the install/upgrade correctness + hint regardless of model.

## Acceptance criteria

(Ideation firms. Each verified by command output / on-disk plugin manifest state on a controlled install — never a prose-grep of the skill or contract.)

**AC-1 (sketch) — install actually advances a behind plugin.** Verified by: a controlled install where the installed plugin manifest is behind the source branch HEAD; after `spacedock install`, the installed plugin manifest version equals the source HEAD (or the command reports an explicit, testable failure) — checked against the resulting on-disk manifest, not the command's own claim alone.

**AC-2 (sketch) — a contract-compatible-but-behind plugin surfaces an opt-in upgrade hint.** Verified by: the message rendered when binary and plugin display semvers differ but are contract-compatible names that a newer plugin is available and how to get it — observed in the rendered output of a fixture with that version skew, not a prose match.

## Test plan

(Ideation/implementation firms.) Facet 1: a host-exec seam test or controlled install smoke that drives `Install` against a behind-source state and asserts the resulting installed manifest advances. Facet 2: a unit/golden test over the verdict-rendering path feeding a contract-compatible binary/plugin semver skew and asserting the hint line; the expected hint text's trigger (semver skew) comes from the version inputs, not from the file under test.

## Notes

First real-world 0.20.0 install feedback. Lands in the 0.20.x cleanup/UX band. Related: the post-flip branch-model decision (separate concern). The "binary follows tags, plugin follows branch HEAD" asymmetry is the backdrop — a tag-fresh binary against a stale branch-HEAD plugin is exactly the skew a user hits.
