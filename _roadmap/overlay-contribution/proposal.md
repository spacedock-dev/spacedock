# Proposal: revise the overlay-contribution roadmap entry

This file proposes changes to `index.md` and to its members. It is based on the run in `answer-today.md` (transcript: `run-log.md`). It does not edit `index.md`, and it does not bind a release. Those decisions belong to CL.

## Summary

1. **Cover two shapes, not one.** `index.md` covers only the contribution shape: pull requests to an upstream that we do not own. CL's request ("maintain our fork") is the second shape: our own changes kept on top of an upstream that moves. The two shapes share the layout and the state mechanism. They differ in four things: PR target, worktree base, upstream sync, and which strings must stay hidden.
2. **Use the layout that is in use.** Replace `<fork>/spacedock/flow/` with `.spacedock/<workflow>/` plus a spec-and-state repository at `.spacedock/repo`. Cargento and asr use this layout. In the run it kept the fork's tree clean and ran the maintained-fork loop end to end (`answer-today.md`, Path 4).
3. **Replace "clone your fork, not the upstream" with remote roles.** In the overlay layout, state does not use the code repository's `origin`. State worked with `origin = upstream` (Path 3, steps 4–6). Only the merge hook needs a code remote that it can push to. It also needs to know where to open the PR. Declare both, and accept both `origin` namings.
4. **Update the statuses.** The branch-name leak is fixed at HEAD. `state init` now reads `state-branch` but fetches from the wrong repository. The binary PR probes still omit `--repo`.
5. **Add three tasks:** remote roles in the merge hook, a declared worktree base, and an upstream-sync procedure for the maintained fork.

## The two shapes

| | Contribution | Maintained fork |
|---|---|---|
| What the user wants | Send changes upstream as pull requests. | Keep our changes on the fork's `main` while the upstream moves. |
| Worktree base | `upstream/main` | the fork's `main` |
| PR target | the upstream repository, head `<fork-owner>:<branch>` | the fork, base `main` |
| Upstream sync | none: each branch is based on the current `upstream/main` | a repeating step: bring `upstream/main` into the fork's `main` (rebase or merge, open decision 1) |
| Must not reach the other party | any spacedock string in the diff, branch name or PR body that the upstream maintainer reads | nothing; the fork is ours. Spacedock files on the fork's `main` are allowed only if CL chooses the in-tree layout (open decision 2) |
| Works today | everything except the PR itself; the PR is manual (Path 4, step 4) | everything, with `origin = fork` (Path 2 in the tree, Path 4 off the tree) |

One clone can serve both shapes. The run did both from one overlay (Path 4, steps 2 and 4). The only difference was the branch checked out at dispatch time.

## Proposed text changes to `index.md`

### Title, Goal and Deliverable

Replace with:

> **Goal:** Let a user run a Spacedock workflow on a fork of a repository that they do not own, in two shapes: **contribution** (send changes upstream as pull requests) and **maintained fork** (keep our changes on top of an upstream that moves).
>
> **Deliverable:** One named and documented layout, one setup command that creates and resumes it, and a merge hook that opens the PR in the right repository for each shape. In the contribution shape, no spacedock-authored string reaches the upstream maintainer.

Keep "Target release train: not bound."

### "Why this sprint exists"

Keep the `prepare.go:137` paragraph and the A/B result. Replace the opening with:

> The mechanism works today when it is set up by hand. A run with local bare repositories drove the maintained-fork loop from `new` to `merge guard` with the code repository clean. It also opened a contribution PR whose upstream diff was one file. Four things are missing:
>
> 1. No command creates or resumes the layout.
> 2. The merge hook pushes to `origin` and opens the PR in `origin`, so it cannot open an upstream PR, and it fails when `origin` is the upstream.
> 3. Nothing tells the user when the fork's `main` falls behind the upstream.
> 4. The first officer stops with "no workflow" when it receives a fork request, and gives no setup path.

### The layout (new section, replaces step 2–3 wording and the `<fork>/spacedock/flow` path)

> ```
> <fork-clone>/                         code repository; .git/info/exclude has .spacedock/ and .worktrees/
> ├── .worktrees/<slug>/                worker worktrees (branch <slug>)
> └── .spacedock/                       plain directory, not a repository
>     ├── repo/                         spec-and-state repository; origin = the state remote
>     │   └── README.md, _mods/         on branch main
>     └── <workflow>/                   plain directory, not a repository
>         ├── README.md -> ../repo/README.md
>         ├── _mods     -> ../repo/_mods
>         └── .spacedock-state/         linked worktree of repo/, orphan branch <workflow>-state
> ```
>
> Two constraints make this work. Do not change them without a test:
>
> - `.spacedock/` and `.spacedock/<workflow>/` must not be repositories. `FindGitRoot` climbs from the workflow directory and stops at the first `.git`. Worktrees must land in the code repository, and `gate prepare` must classify code artifacts as `main`.
> - `state:` cannot contain `..`, so the state checkout sits under the workflow directory.

Why this replaces `<fork>/spacedock/flow/`: that path keeps the README untracked in the fork, so the spec has no history and no way to reach a second machine. `.spacedock/repo` versions the spec and the state in one repository, outside the fork. The leading dot keeps it out of directory listings.

### Remotes (replaces the "CAUTION: do not clone the upstream directly" section)

> The layout uses three remote roles. The names of the remotes are free.
>
> | Role | Used by | Default |
> |---|---|---|
> | state remote | `state ready`, `state commit`, gate verbs, `merge guard` | `origin` of `.spacedock/repo` |
> | code push remote | the merge hook: branch push and trunk push | `origin` of the code repository |
> | PR repository | the merge hook: `gh pr create --repo`, and the `pr` field | the repository of the code push remote |
>
> A maintained fork that follows the roadmap naming (`origin = fork`) needs no declaration. The asr naming (`origin = upstream`, fork as `fork`) declares the code push remote as `fork`. A contribution declares the PR repository as the upstream, and the hook then passes `--head <fork-owner>:<branch>`.
>
> CAUTION: In the layout where the workflow is in the tree, state uses the code repository's `origin`. If you cannot write to `origin`, `state ready`, `state commit` and `dispatch build --stamp` fail. `dispatch build --stamp` writes `worktree:` and then stops before it creates the worktree.

The CAUTION is still true at HEAD (Path 1, steps 7–9). Keep its `publish.go:70` explanation.

### The journey table

Replace the "Status today" column with the statuses from the run:

| # | Step | Status today |
|---|------|--------------|
| 1 | Fork, clone, add the second remote | Manual. No skill text mentions clone, fork or upstream. |
| 2 | Create the overlay (exclude, spec repository, symlinks, state worktree) | Manual. `state new` creates the state branch in the code repository and adds `?? .gitignore`. A short setup script (in `suggested-phrasing.md`) does the whole step. |
| 3 | Write the README | Manual. Commission has no overlay journey, and its split-root journey writes to the code repository. |
| 4 | Boot from the fork root | Works. |
| 5 | `state ready`, `new`, `state commit` | Works. State goes to the state remote, also when `origin` is the upstream. |
| 6 | `dispatch build --stamp` | Works. Branch name is the bare slug. The worktree starts from the checked-out commit. |
| 7 | Ensign commit and stage report | Works. |
| 8 | `gate prepare`, `gate record --consume` | Works. The artifact must be committed. |
| 9 | Merge hook opens the PR | Works for the maintained fork when `origin` is the fork. **Gap** for contribution: the PR always opens in `origin`. **Fails** when `origin` is the upstream. |
| 10 | Track the PR to merge, `merge guard` | `merge guard` works. **Gap:** `state sweep` and boot `PR_STATE` call `gh pr view <owner/repo#N>` with no `--repo`. |
| 11 | Resume on a fresh machine | **Gap.** `state init` fetches `state-branch` from the code repository's `origin`. The setup script works. |
| 12 | The upstream moves (maintained fork) | **Gap.** No command notices or acts on it. Sync is plain git. |

### "What the maintainer sees today"

Replace the two-strings paragraph with:

> One spacedock string can still reach the upstream maintainer, and only in the contribution shape: the audit link in the PR body, which `mods/pr-merge.md` marks as required. The branch name is the bare slug at HEAD. In the maintained-fork shape the PR opens in our own fork, so the audit link is wanted there.

### Definition of Done

1. In the contribution shape, a user reaches an upstream PR through documented steps. The diff against `upstream/main`, the head ref, and the PR body contain no spacedock-authored string. A behavior fixture proves all three. *(Revised: the head-ref part passes at HEAD and stays as a regression check.)*
2. In the maintained-fork shape, the loop from `new` to `merge guard` runs with the code repository clean, for both remote namings. A behavior fixture proves it. *(New.)*
3. One spacedock command creates the layout. On a fresh machine, the same command or `state init` resumes it, and `state ready` exits 0 with entity history intact. *(Replaces old DoD 2 and 3.)*
4. A `pr` field that carries `owner/repo#N` resolves to the real PR state from any launch directory. *(Unchanged, old DoD 4.)*
5. A probe error during `state sweep` produces a diagnostic. *(Unchanged.)*
6. Two registered merge mods are either both armed or refused at validation. *(Unchanged.)*
7. The shipped merge hook opens the PR in the declared PR repository, pushes to the declared code push remote, and does not fire twice. *(Revised old DoD 7.)*
8. When the fork's `main` is behind the upstream, the FO says so at engage, and the documented sync procedure runs. *(New. Its form depends on open decision 1.)*
9. A commission journey or the setup command produces a working overlay with no manual repair. *(Unchanged, old DoD 8.)*
10. `go test ./...` is green.

### Out of scope

Keep "Greenroom", "A fork concept in the binary", and "Survey and debrief in a foreign repository".

Revise "A second remote model" to:

> PR creation stays in the mod body. The remote roles are README frontmatter that the mod body reads. The binary reads only the worktree base (task 6 below).

## Tasks

Existing members, with proposed changes:

1. **`overlay-contribution-state-remote` (dq)**:
   - Restate the problem for the `.spacedock/repo` layout.
   - The run found that `state init` reads `state-branch: dev-state` correctly, but runs `git fetch origin` in the code repository (`fatal: couldn't find remote ref dev-state`). Resolve the state repository from the spec repository (the README symlink target's git root) or from a `state-remote:` key.
   - AC-1's fixture path changes from `<repo>/spacedock/flow` to `.spacedock/dev`.
   - `state new` must create the orphan branch in the spec repository, not in the code repository.
2. **`pr-probe-repo-qualifier` (12c)**: no change. The run confirms that `state sweep` calls `gh pr view getpaseo/paseo#1 --json state` with no `--repo`. `mods/pr-merge.md:159` still records `#57`, unqualified.
3. **`merge-guard-arms-first-mod-only` (wpb)**: no change.
4. **`overlay-contribution-shape` (dn)**:
   - Rename it to cover both shapes.
   - Drop the branch-name leak from the problem statement: it is fixed. Keep AC-2 (a declarable branch prefix) only if CL wants a namespace. The run shows one reason: two workflows or two clones on one fork collide on equal slugs (Path 3, step 8b).
   - Add the setup command and an overlay commission journey.
   - Change the FO's no-workflow stop so that it names the setup path.

New members:

5. **Remote roles in the merge hook.** The pr-merge hook reads a code push remote and a PR repository from the README, with the defaults in the Remotes table. When the PR repository is not the push remote's repository, it passes `--head <owner>:<branch>` and skips the trunk push. AC: a fixture with a stubbed `gh` shows the right `--repo` and `--head` for the three cases in the run: `origin = fork`, `origin = upstream` with a `fork` remote, and contribution. The audit link is left out when the PR repository is not ours. Closes Path 3 steps 7 and 11, and Path 4 step 4.
6. **Declared worktree base.** `dispatch build --stamp` runs `git worktree add -b <slug> <path>` with no start point, so the base is whatever the code checkout has checked out. Add a declared base ref (for example `base: upstream/main`) and pass it as the start point. AC: a contribution worktree starts at `upstream/main` while the code checkout stays on the fork's `main`. This is a binary change.
7. **Upstream sync for the maintained fork.** At engage, report how far the fork's `main` is behind `upstream/main`. Ship the sync procedure as a mod hook or a documented FO step. Its form depends on open decision 1. With rebase, the procedure must also handle open worktree branches and the SHAs cited in stage reports. With merge, a conflict goes to the captain.
8. **Real-fork trial** *(only if CL chooses option 1 of open decision 3)*: drive the shape on a paseo fork, after CL approves the clone and the push access.

Sequence: tasks 1, 2, 3 and 5 are independent and can run first. Task 6 comes before the contribution part of task 4. Task 7 needs open decision 1. Task 4 runs last.

## Open decisions for CL

### 1. When the upstream moves: rebase or merge?

- **Option 1, merge** `upstream/main` into the fork's `main`. The push is a fast-forward. Every SHA that a stage report cites stays on `main` (Path 2, step 9b). Our changes are mixed with merge commits, so sending them upstream later needs a cherry-pick.
- **Option 2, rebase** the fork's `main` on `upstream/main`. Our changes stay as a clean series, ready to send upstream. Each sync needs a force push, open worktree branches need a rebase, and cited SHAs leave `main` (Path 2, step 9a). The state history keeps the old SHAs, but they are then unreachable from `main`.

### 2. Spacedock files: in the tree of the fork, or off the tree?

- **Option 1, in the tree.** This is the commission default and needs no new command. The loop runs today with no manual step after the clone (Path 2). The fork's `main` carries the README and `.gitignore` lines, and the state is a public branch on the fork. Contributions must be based on `upstream/main` so that they do not carry them.
- **Option 2, off the tree** (`.spacedock/` with its own spec-and-state repository). The fork stays equal to the upstream except for our code. The state can be private. It needs the setup command (task 1 and task 4), or the script in `suggested-phrasing.md` until then (Path 4).

This proposal is written for option 2, because it serves both shapes from one clone. With option 1, the Layout section becomes "commission default", and task 1 shrinks to the contribution case.

### 3. Paseo: the real first target, or only an example?

- **Option 1, real first target.** Add task 8. The fixtures stay bare-repository fixtures, but the shape is accepted only after a run on a paseo fork. This needs CL's approval for the clone and for push access to a fork.
- **Option 2, only an example.** Bare-repository fixtures prove the shape. The first real case is the Hermes fork that the exe.dev bot maintains, when CL gives its URL or access. `paseo` stays only as the example name in the docs.
