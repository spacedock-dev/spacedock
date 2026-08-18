---
id: dnk9880z04q7qm9gv4vb983c
title: Name and support the overlay-contribution shape for external repositories
status: backlog
source: FO research session 2026-08-18 - overlay-contribution probe (workflow run wf_40ea0f6e-aa8)
started:
completed:
verdict:
score:
worktree:
issue:
sprint: overlay-contribution
group: shape-and-scaffolding
---

A contributor who wants to send pull requests to a repository they do not own can already run a full Spacedock workflow against it, by overlaying an untracked workflow directory inside their fork clone. The mechanism works and was proven end to end. Nothing names the shape, nothing creates it, and two spacedock-authored strings still reach the upstream maintainer. This task makes the shape first-class.

## Problem

The shape is: fork the upstream, clone the fork, put the workflow at `<fork>/spacedock/flow` excluded through `.git/info/exclude`, use split-root state, and open the upstream pull request from a merge-hook mod.

The load-bearing discovery is that overlay **placement** is what makes cross-repository gating work. `internal/gates/prepare.go:137` builds `gitsource.Roots{Main: workflowDir, State: filepath.Dir(entityPath)}`. Overlaying the workflow directory inside the code repository makes `Main` resolve to that repository, so code artifacts classify as `main`, which leaves `State` free to be a clone of a wholly separate private repository. An A/B control isolates the causation: identical state repository and identical artifact, with the workflow directory in a third repository, returns `selected source is not owned by a workflow Git root` at exit 1; with the workflow directory overlaid inside the code repository the same command returns `git-root://main/...` at exit 0.

### End-to-end journey, and its current status

| Step | What it looks like | Status today |
|---|---|---|
| 1. Fork and clone | `gh repo fork <upstream> --clone`, then `git remote add upstream <url>` | Works. Nothing in spacedock knows what a fork is, so the contributor does this by hand. |
| 2. Exclude the overlay | `printf 'spacedock/\n.worktrees/\n' >> .git/info/exclude` | Works. Discovery still finds the workflow: the downward scan prunes only `discoverIgnoreDirs` (`internal/status/handlers.go:579-583`) and tracked-gitignore directory patterns, never `info/exclude`, and does not skip untracked or hidden directories. |
| 3. Commission the workflow | Author `<fork>/spacedock/flow/README.md` | Works, by hand. No commission journey covers this shape. |
| 4. Create state | `spacedock state new` | **Gap.** Births the orphan in the code repo and dirties the tracked `.gitignore`. Tracked by `overlay-contribution-state-remote`. |
| 5. Boot | `spacedock status --boot --identify --json` from the fork root | Works. Fails from a greenroom-style non-git wrapper directory, because the downward walk prunes every child holding `.git`. |
| 6. Dispatch | `spacedock dispatch build --stamp` | Works. Creates the worktree in the fork. Branch is named `spacedock-ensign/<slug>` - see the leak below. |
| 7. Implement and report | ensign commits code, files the stage report | Works. The stage-report completeness guard (`internal/status/entered_stage.go:41-65`) blocks advance until a complete report is committed. |
| 8. Gate | `spacedock gate prepare --artifact ...` | Works, including a code artifact in the fork and an entity artifact in a separate private state repository resolved in one briefing. |
| 9. Open the PR | merge-hook mod runs `gh pr create` | Works today with a user-authored `_mods/` merge hook and no Go change: the binary never issues `gh pr create`, the mod body does. Mod discovery is structural (`internal/status/mutate.go:523-561` globs `_mods/*.md` and registers any `## Hook:` line, with no allowlist), and `refit` leaves custom-named mods alone. Nothing is shipped for it. |
| 10. Track the PR to merge | pr-merge startup/idle hooks poll, then `merge guard` finalizes | **Gap.** The probes drop `--repo`; tracked by `pr-probe-repo-qualifier`. `merge guard` correctly parks on an open PR indefinitely with no drift (`internal/status/merge.go:172-177`). |

Confidentiality at the content layer is airtight and was measured: across every ref in the fork, the repository held only `.gitignore`, `README.md`, and the source file, and grepping every blob for `spacedock`, `Stage Report`, and `briefing` returned zero hits. The maintainer-visible diff was exactly one file and four insertions, with `git status --porcelain --untracked-files=all` empty.

### The two remaining leaks

- **Branch name.** `internal/dispatch/stamp.go:98` and `:143` derive the worktree branch as `strings.ReplaceAll(subagentType,":","-") + "/" + slug`, so the pull request head ref reads `spacedock-ensign/<slug>` and the maintainer sees it. Declaring `agent: contrib` on the worktree stage yields `contrib/<slug>` but also sets `subagent_type: contrib` in the dispatch envelope, so that agent must exist - a workaround, not a fix.
- **Audit link.** `mods/pr-merge.md:95` marks a root-relative link to the contributor's state branch **Required** in the pull request body template, so it renders in the text the maintainer reads.

### Not required to fork first

A contributor can start before forking, but only if the clone has no remote named `origin`. `internal/statesync/publish.go:70` defines "no remote" as "`git remote get-url origin` fails", not "publication is impossible", so `git clone -o upstream <url>` degrades every state verb cleanly to local-only and the whole loop runs at exit 0; adding `origin` later publishes the accumulated orphan history on the next `state commit`. A naive `git clone <upstream>` is dead on arrival: `state ready` (the FO's first boot call), `state commit`, `gate record --consume`, and `merge guard` all exit 1, and `dispatch build --stamp` deadlocks - it writes `worktree:` into the frontmatter, then dies before creating the worktree, after which plain `dispatch build` refuses with `worktree path ... does not exist`.

## Out of scope

The state creation and resume gap (`overlay-contribution-state-remote`), the PR-probe qualifier gap (`pr-probe-repo-qualifier`), and the merge-hook arming defect (`merge-guard-arms-first-mod-only`). This task owns naming the shape, closing the two leaks, and shipping the scaffolding a contributor needs.

## Acceptance criteria

**AC-1 (VALUE) - A contributor reaches an upstream-ready pull request through documented spacedock steps, and the maintainer sees no spacedock-authored string in the diff, the branch name, or the pull request body.**
Verified by: a behavior fixture over local bare repos that drives the shape to the point of the `gh pr create` arguments (gh stubbed), then asserts three things - the diff against `upstream/main` contains only the intended source change, the head ref matches no `spacedock` or `ensign` substring, and the rendered body contains no link into the state repository. Each of the three fails today.

**AC-2 - The worktree branch prefix is declarable and independent of `subagent_type`.**
Verified by: a unit test asserting that a declared prefix produces the branch name while the dispatch envelope's `subagent_type` is unchanged. The falsifying change is re-deriving the branch from `subagent_type` at `stamp.go:98`.

**AC-3 - A shipped merge-hook mod opens a fork-to-upstream pull request, and installing it does not double-fire alongside `pr-merge`.**
Verified by: a fixture asserting the mod's `gh pr create` arguments carry `--repo <upstream>` and `--head <owner>:<branch>`, plus an assertion that only one merge hook fires for the shape. Depends on `merge-guard-arms-first-mod-only`.

**AC-4 - The shape has a commission journey that produces a working overlay without manual repair.**
Verified by: following the journey in a fixture and asserting `git status --porcelain --untracked-files=all` is empty at every step. Depends on `overlay-contribution-state-remote` for AC-1 of that task.

## Test plan

Command-level behavior fixtures over local bare repos, with gh stubbed on PATH; the probe run established this substrate is sufficient and needs no network or forge account. Go unit tests for the branch-prefix derivation. The commission journey needs a skill-integration test in the `skills/integration` shape (extract the shipped runnable block, execute it against independent fixture conditions, assert on-disk state), not a prose grep over the skill text. Sequence after its two dependencies land, since AC-3 and AC-4 both rest on them.
