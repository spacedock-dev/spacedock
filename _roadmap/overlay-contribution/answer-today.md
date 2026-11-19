# What happens today: "clone paseo and maintain our fork with changes managed by spacedock"

This file says what happens when a user gives that sentence to a Spacedock first officer (FO). The evidence is one run in `/tmp/overlay-run`. The run used local bare repositories for the upstream, the fork, and a state remote, plus a stub for `gh`. `run-log.md` has the full transcript. Each "Evidence" cell names a section of that log.

The run used `spacedock 0.29.0-pre0+dev`, built from `db4728199` on `main`. No model ran. For each FO or ensign step, the run executed the exact command that the skill text names.

## Short answer

The FO does nothing with the sentence. It finds no workflow, says so, and stops. Clone, fork, remote setup and workflow placement are all manual, and no skill text covers them.

When a person does those steps by hand, the result depends on two choices that the sentence does not make:

1. **Which repository `origin` is.** If `origin` is the upstream (the literal "clone paseo"), the workflow loop fails at its first convergence step. If `origin` is the fork, the loop runs to the end.
2. **Where the workflow lives.** The default commission layout puts the workflow README and `.gitignore` edits on the fork's `main`, and the state branch on the fork. The off-tree overlay (cargento, asr) keeps the fork clean, but no command creates it, and `state new` and `state init` fail in it. With `origin = fork` and a setup script, the overlay runs the maintained-fork loop end to end (Path 4).

No spacedock command does anything when the upstream moves.

Marks: **auto** = a spacedock command or skill does it and exits 0. **manual** = a person or the agent must do it with plain `git`/`gh`, and no skill text says how. **fails** = the step exits non-zero or leaves broken state.

## Path 1: the sentence as given

The FO receives the sentence through `spacedock claude -- "clone paseo and maintain our fork ..."`.

| # | Step | What happens | Mark | Evidence |
|---|------|--------------|------|----------|
| 1 | FO boot in the launch directory | `status --boot --identify --json` exits 1: `no commissioned Spacedock workflow found ... report this and stop; do NOT search the filesystem`. The FO contract (`«state.boot»`, zero discovery) says: report and STOP. | auto (stops) | A0b |
| 2 | Clone paseo | No spacedock command or skill text clones. Done literally, `git clone <upstream>` gives `origin = upstream`. | manual | A1, S1 |
| 3 | Fork paseo | No spacedock command or skill text forks. A search of the commission and FO skills for `clone\|fork\|upstream\|overlay\|info/exclude` gives no relevant hit. | manual | S1 |
| 4 | Boot again in the clone | Exits 1 with the same "no workflow" message. The FO stops again. | auto (stops) | A2 |
| 5 | Commission the workflow | Commission has two journeys: split-root in the same repository, and inline. It has no overlay journey. The default directory is `docs/<slug>/`, and commission appends `.worktrees/` to the tracked `.gitignore`. | manual | A2, S1 |
| 6 | `state new` | Exits 0. It creates the orphan branch in the code repository and appends `docs/paseo-fork/.spacedock-state/` to the tracked `.gitignore`. It prints `Warning: could not push spacedock-state/paseo-fork to origin`. `git status` shows `?? .gitignore` and `?? docs/paseo-fork/README.md`. | auto, with a warning | A3 |
| 7 | FO `engage`: `state ready` | Exits 1: `fatal: couldn't find remote ref refs/heads/spacedock-state/paseo-fork`. `engage` runs this first, so the FO halts here. | **fails** | A4 |
| 8 | `new` + `state commit` | `new` exits 0. `state commit` exits 1: `state publication failed; the new local commit remains recoverable`. | **fails** | A4 |
| 9 | `dispatch build --stamp` | Prints `state publication failed; the stamp is already durable locally`. It writes `started:` and `worktree: .worktrees/spacedock-ensign-add-banner` to the entity, then stops. `git worktree list` shows no code worktree. A retry exits 1 with the same error. The worker never gets a worktree. | **fails** | A5, A5b |

The CAUTION in `index.md` is still correct at HEAD: an `origin` that you cannot write to breaks the loop at `state ready`, `state commit`, and `dispatch build --stamp`.

## Path 2: clone the fork, workflow in the tree

This is what a person gets after doing steps 2 and 3 "right": clone the fork, add the upstream as a second remote, and let commission use its default layout.

| # | Step | What happens | Mark | Evidence |
|---|------|--------------|------|----------|
| 1 | Clone the fork, `git remote add upstream ...` | `origin = fork`, `upstream = upstream`. | manual | B1 |
| 2 | Commission default + `state new` | Exits 0. The orphan branch reaches the fork (`refs/heads/spacedock-state/paseo-fork` in `git ls-remote fork.git`). A person commits the README and `.gitignore` on `main` and pushes them to the fork, as commission says. | auto | B2 |
| 3 | `state ready`, `new`, `state commit` | All exit 0: `Committed and pushed add-banner to spacedock-state/paseo-fork`. | auto | B3 |
| 4 | `dispatch build --stamp` | Exits 0. The worktree is `.worktrees/spacedock-ensign-add-banner` on branch **`add-banner`**, based on the checked-out `main` (`17a9a80`). | auto | B4 |
| 5 | Ensign commit + stage report + `state commit` | Exits 0. | auto | B5 |
| 6 | `gate prepare` + `gate record --consume` | Exit 0. Route `approved-awaiting-merge`. A review artifact must be committed. An uncommitted file in the code worktree is refused (`selected source is not the exact committed file`). | auto | B6 |
| 7 | pr-merge merge hook | The hook resolves `CODE_REPO` from the worktree's `origin`: `us/paseo`, which is the fork. It pushes the branch to the fork and runs `gh pr create --repo us/paseo --base main --head add-banner`. The pull request opens **in the fork**, not upstream. The audit link resolves the state repository from the state checkout's `origin`, which is also the fork. | auto (fork PR) | B7 |
| 8 | Merge in the fork, sentinel, `merge guard` | Exits 0: `finalized: add-banner -> done (verdict passed), archived`. | auto | B8 |
| 9 | The upstream moves | `state ready` and `state sweep` exit 0 and do not mention the upstream. The boot record has no upstream field. Nothing tells the FO that `main` is now 1 commit behind `upstream/main`. | manual | B9 |
| 9a | Sync by rebase | `git rebase upstream/main` replays both fork commits. The SHA that the stage report cites (`0c9d01f`) is no longer on the rebased `main`. The push to the fork is rejected as non-fast-forward, so it needs a force push. | manual | B9a |
| 9b | Sync by merge | `git merge upstream/main` keeps `0c9d01f` as an ancestor. The push to the fork is a fast-forward. | manual | B9b |
| 10 | What the fork carries on top of upstream | `git diff --stat upstream/main...main`: `.gitignore` (+2), `docs/paseo-fork/README.md` (+52), `src/banner.txt` (+1). The fork also has the public branch `spacedock-state/paseo-fork`. | result | B10 |
| 11 | Contribution from this clone | The worktree branch descends from the fork's `main`, so an upstream PR from `add-banner` carries the commit `Add spacedock workflow` and the same three files. To send it upstream, the hook needs `--repo <upstream> --head <fork-owner>:<branch>`. The shipped hook does not do this. | **fails** for contribution | B11 |

Path 2 is a working maintained fork if spacedock files on the fork's `main` are acceptable. It does not give a clean upstream contribution.

## Path 3: overlay off the tree (the cargento and asr layout)

The workflow lives under `.spacedock/`, which is in `.git/info/exclude`. `.spacedock/repo` is its own repository, and its `origin` is the state remote. `.spacedock/dev/README.md` is a symlink into it. The state checkout is a linked worktree of `.spacedock/repo` on the orphan branch `dev-state`. The code remotes follow asr: `origin = upstream`, `fork = fork`.

| # | Step | What happens | Mark | Evidence |
|---|------|--------------|------|----------|
| 1 | Clone, add `fork` remote, write `info/exclude` | | manual | C1 |
| 2 | `state new --workflow-dir .spacedock/dev` | Exits 0, but in the wrong repository. `git rev-parse` climbs from `.spacedock/dev` to the code repository. So `dev-state` is created in the code repository, the push goes to the upstream (and fails with a warning), and a new `.gitignore` appears (`?? .gitignore`). | **fails** for this layout | C2 |
| 3 | Build the overlay by hand (cargento `setup.sh` pattern) | `git init` the spec repository, push `main` and an orphan `dev-state` to the state remote, symlink the README, `git worktree add` the state checkout. The code repository's `git status --porcelain --untracked-files=all` is empty. | manual | C3 |
| 4 | Boot from the code root, `state ready`, `new`, `state commit` | All exit 0. Boot finds `.spacedock/dev` (`state_backend: split-root`). State commits go to `state.git` on `dev-state`. | auto | C4 |
| 5 | `dispatch build --stamp` | Exits 0. The worktree lands in the code repository (`.worktrees/spacedock-ensign-add-banner`, branch `add-banner`), based on the checked-out `main`. | auto | C5 |
| 6 | Ensign + `gate prepare` + `gate record --consume` | Exit 0, `approved-awaiting-merge`. | auto | C6 |
| 7 | pr-merge merge hook | `git remote get-url origin` in the worktree is the upstream, so `CODE_REPO = getpaseo/paseo`. The branch push to `origin` is rejected (`[remote rejected] ... add-banner`). The hook text then falls back to a local merge. | **fails** | C7 |
| 8 | Push to `fork`, open the upstream PR (by hand) | A push to `fork` exits 0. `gh pr create --repo getpaseo/paseo --head us:add-banner` is the command that a fork contribution needs. The upstream diff is one file with one added line (`src/banner.txt`). The code repository stays clean. | manual | C8, C8b |
| 8b | Branch name collision | The first push to `fork` was rejected (`fetch first`). Scenario B had already pushed a branch `add-banner` to the same fork. The worktree branch is the bare slug, so two workflows or two clones that use one fork collide on equal slugs. | observation | C8, C8b |
| 9 | Merged probe with `pr: getpaseo/paseo#1` | `state sweep` reports `PR merged`. `gh.log` shows that the binary called `gh pr view getpaseo/paseo#1 --json state` with no `--repo`. The stub accepted that form. This run did not test whether real `gh` accepts `owner/repo#N` as the positional argument. `merge guard` then exits 0. | auto against the stub; real `gh` untested | C9 |
| 10 | Fresh machine: `state init` | Exits 1: `git fetch origin dev-state failed: couldn't find remote ref dev-state`. `state init` reads the declared `state-branch: dev-state` correctly. But it fetches from the **code** repository's `origin` (the upstream), not from the spec repository. | **fails** | C10 |
| 10b | The same step by hand (`setup.sh` step 4) | `git -C .spacedock/repo fetch origin dev-state && git -C .spacedock/repo worktree add ...` then `state ready` exits 0. The archived entity is present. | manual | C10b |
| 11 | Maintained fork with asr remotes | With local `main` set to `fork/main`, `dispatch trunk` returns `main`. The hook's trunk push is `git push origin main`, so it targets the upstream (rejected here as non-fast-forward). The correct target is `fork`. | **fails** | C11 |

## Path 4: overlay off the tree, `origin = fork`

This is Path 3 with the roadmap's remote naming (`origin = fork`, `upstream` as the second remote). It uses the setup script in `suggested-phrasing.md`.

| # | Step | What happens | Mark | Evidence |
|---|------|--------------|------|----------|
| 1 | `overlay-setup.sh` first run | Clones the fork, adds `upstream`, writes `info/exclude`, creates the spec repository and `dev-state` on the state remote, and runs `state ready` (exit 0). The code repository is clean. | manual (one script) | D1 |
| 2 | Maintained fork loop | `new`, `state commit`, `dispatch build --stamp`, `gate prepare`, and `gate record --consume` all exit 0. The merge hook resolves `CODE_REPO = us/paseo` and runs `gh pr create --repo us/paseo --head add-footer`. The audit link resolves the state repository as `us/paseo-state`, not the fork. The fork gets only the branch `add-footer`. The code repository stays clean. | auto | D2 |
| 3 | Fresh machine | The same script with no README argument clones the spec repository, attaches `dev-state`, and runs `state ready` (exit 0). The resumed clone lists `add-footer`. | manual (one script) | D3, D3b |
| 4 | Contribution from the same clone | Check out `upstream/main` before dispatch. The worktree is then based on `dd96437 upstream: v2`. Pushing the branch to `origin` and running `gh pr create --repo getpaseo/paseo --head us:fix-typo` gives an upstream diff of one file (`README.md`) and one commit. The shipped hook does not do this step: it always uses `origin` as the PR repository, and `dispatch trunk` still returns `main`. | manual | D4 |
| 5 | Two workflows visible | The fork's `main` still had scenario B's in-tree workflow, so `new` with no `--workflow-dir` refused: `multiple commissioned Spacedock workflows found`. This is correct behavior, but a fork that switches layouts must remove the old README. | observation | D2 (first attempt) |

Path 4 is the shape that works today with the fewest manual steps for both cases. The contribution PR is still manual.

## Corrections to `index.md` from this run

- **Branch-name leak: fixed.** `index.md` step 6 says the branch is `spacedock-ensign/<slug>`. At HEAD the branch is the bare slug (`add-banner`, B4 and C5). Only the local worktree directory keeps the `spacedock-ensign-` prefix, and it never reaches a remote.
- **`state init` and `state-branch`.** Cargento's `DETACH.md` says `state init` fetches the default branch name. At HEAD it fetches the declared `state-branch`. It fails in the overlay because it fetches from the wrong repository (C10).
- **Probe `--repo`.** The shipped `mods/pr-merge.md` now passes `--repo` for qualified and unqualified references. The two binary probes (`state sweep` and boot `PR_STATE`) still run `gh pr view <ref>` with no `--repo` (C9).
- **The CAUTION about `origin`** still holds (Path 1, steps 7–9).

## What this run did not cover

- The asr checkout. My sandbox cannot read `~/git/asr/WhisperLiveKit/.spacedock/repo` (`Operation not permitted`). The asr remote shape in Path 3 comes from the mandate's description, not from the files.
- The Hermes fork. CL has given no URL, clone, or access.
- Real `gh`, real GitHub, and a live model. A real FO might improvise past a "manual" step. This file reports only what the binary and skill text do.
- A paseo checkout. The upstream here is a two-file stand-in.
