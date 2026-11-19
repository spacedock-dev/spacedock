# Suggested phrasing: a fork of paseo, managed by Spacedock

This file gives the words and commands that produce a working result today. `answer-today.md` has the evidence, and `run-log.md` has the transcript.

## Why the plain sentence is not enough

"Clone paseo and maintain our fork with changes managed by spacedock" does not work as an instruction to the first officer (FO), for two reasons:

1. The FO runs workflows that already exist. With no workflow, it reports that and stops (`answer-today.md`, Path 1, steps 1 and 4). Clone, fork and setup come before the FO.
2. The sentence leaves open the two choices that decide the result:
   - which repository is `origin`;
   - where the workflow and its state live.

   With `origin` set to the upstream, the loop fails at the first `state ready` (Path 1, step 7).

## The sentence to use

Say it to a setup agent or do it yourself, before you start the FO. Fill in every slot:

> Set up a **maintained fork** of `getpaseo/paseo`. Clone our fork `us/paseo` so that `origin` is the fork, and add the upstream as the remote `upstream`. Keep the Spacedock workflow **{in the tree | off the tree}**. Put the state in **{the fork | `us/paseo-state`}**. Our changes merge as pull requests into the fork's `main`. When the upstream moves, **{rebase | merge}** `main` on `upstream/main`.

For a **contribution**, replace the first and the last two sentences with: "Each change goes to `getpaseo/paseo` as a pull request from a branch on our fork, based on `upstream/main`."

After setup, start the FO from the fork root:

```sh
cd paseo && spacedock claude
```

`paseo`, `us/paseo` and `us/paseo-state` are examples. CL has not decided whether paseo is the real first target. The commands do not depend on it.

## Rules that the run proved

- **`origin` must be a remote you can write to.** Clone the fork, or clone with `git clone -o upstream <upstream-url>` and add the fork as `origin`. In the asr shape (`origin` = upstream, fork = `fork`), state can still work if the state lives in its own repository (Path 3). But the shipped pr-merge hook pushes to `origin` and fails (Path 3, steps 7 and 11).
- **Base a contribution on `upstream/main`.** The worktree starts from the commit that the code checkout has checked out (Path 4, step 4). A branch from the fork's `main` carries every fork commit, and every workflow file if the workflow is in the tree (Path 2, step 11).
- **Open the upstream pull request by hand.** The shipped hook always opens the PR in the `origin` repository. For an upstream PR, run `gh pr create --repo <upstream-owner>/<repo> --base main --head <fork-owner>:<branch>` (Path 4, step 4).
- **Give the `pr` field as `owner/repo#N`** for an upstream PR. The hook text uses `--repo` for that form. The `state sweep` probe in the binary passes the value to `gh` with no `--repo`. Real `gh` was not tested with that form (Path 3, step 9).

## Option 1 of 2 for "where does the workflow live": in the tree

The workflow README and the `.gitignore` entries are committed to the fork's `main`. The state is an orphan branch on the fork. This is what commission does by default.

```sh
# Clone the fork with the upstream as a second remote
git clone git@github.com:us/paseo.git && cd paseo
git remote add upstream https://github.com/getpaseo/paseo.git
```

Then run the commission skill (`/spacedock:commission`) in that clone, with a worktree stage and the pr-merge mod. After that, the full maintained-fork loop runs with no manual step, through `merge guard` (Path 2, steps 2–8).

Costs:

- The fork's `main` differs from the upstream by the workflow files (Path 2, step 10).
- The state history is a public branch on the fork.
- Contributions need a branch from `upstream/main`, and they must not contain the workflow commit.

## Option 2 of 2 for "where does the workflow live": off the tree

`.spacedock/` is excluded through `.git/info/exclude`. `.spacedock/repo` is a separate repository that holds the spec on `main` and the state on the orphan branch `dev-state`, and its remote is the state repository. The fork's tree stays the same as the upstream's, except for our code changes.

No spacedock command creates this layout:

- `state new` puts the state branch in the code repository (Path 3, step 2).
- `state init` fetches from the code repository's `origin` (Path 3, step 10).

The script below did both jobs in the run: first setup, and resume on a fresh machine (Path 4, steps 1 and 3). It is the cargento `setup.sh` pattern, with the URLs and the README as arguments.

```sh
#!/bin/sh
# overlay-setup.sh — create or resume an off-tree Spacedock overlay on a fork clone.
# Usage: overlay-setup.sh UPSTREAM_URL FORK_URL STATE_URL DIR [README]
#   DIR     where to clone the fork (the code checkout)
#   README  workflow README for the first setup (omit when STATE_URL already has it)
# Layout: DIR/.spacedock/repo = spec repo (origin = STATE_URL, branch main)
#         DIR/.spacedock/dev/README.md -> ../repo/README.md
#         DIR/.spacedock/dev/.spacedock-state = linked worktree of the spec repo, branch dev-state
set -e
UPSTREAM=$1 FORK=$2 STATE=$3 DIR=$4 README=$5
BRANCH=dev-state

[ -d "$DIR/.git" ] || git clone -q "$FORK" "$DIR"
cd "$DIR"
git remote get-url upstream >/dev/null 2>&1 || git remote add upstream "$UPSTREAM"
for p in .spacedock/ .worktrees/; do
  grep -qxF "$p" .git/info/exclude || echo "$p" >> .git/info/exclude
done

if [ ! -d .spacedock/repo/.git ]; then
  if git ls-remote --exit-code --heads "$STATE" main >/dev/null 2>&1; then
    git clone -q "$STATE" .spacedock/repo                      # resume
  else
    git init -q -b main .spacedock/repo                        # first setup
    cp "$README" .spacedock/repo/README.md
    git -C .spacedock/repo add README.md
    git -C .spacedock/repo commit -qm "Add workflow spec"
    git -C .spacedock/repo remote add origin "$STATE"
    git -C .spacedock/repo push -q origin main
    git -C .spacedock/repo checkout -q --orphan "$BRANCH"
    git -C .spacedock/repo rm -rfq --cached .
    git -C .spacedock/repo commit -q --allow-empty -m "seed state"
    git -C .spacedock/repo push -q origin "$BRANCH"
    git -C .spacedock/repo checkout -qf main
  fi
fi

mkdir -p .spacedock/dev
ln -sfn ../repo/README.md .spacedock/dev/README.md
[ -d .spacedock/repo/_mods ] && ln -sfn ../repo/_mods .spacedock/dev/_mods
if [ ! -d .spacedock/dev/.spacedock-state ]; then
  git -C .spacedock/repo fetch -q origin "$BRANCH"
  git -C .spacedock/repo show-ref --verify --quiet "refs/heads/$BRANCH" ||
    git -C .spacedock/repo branch -q "$BRANCH" "origin/$BRANCH"
  git -C .spacedock/repo worktree add -q ../dev/.spacedock-state "$BRANCH"
fi
"${SPACEDOCK_BIN:-spacedock}" state ready --workflow-dir .spacedock/dev
```

The README needs `state: .spacedock-state` and `state-branch: dev-state` in its frontmatter. Commission has no off-tree journey, so do not run it inside the fork: its split-root journey writes to the code repository. Write the README by hand, or copy it from a commission run in a scratch directory and set those two keys. The run used a hand-written README. A commission run in a scratch directory was not tested.

Costs:

- A second repository to create and give access to.
- One script to run on each new machine.
- `--workflow-dir .spacedock/dev` on commands when the clone also has another workflow (Path 4, step 5).

## Option 1 of 2 for "when the upstream moves": merge

```sh
git fetch upstream && git checkout main && git merge upstream/main && git push origin main
```

The push is a fast-forward. Every commit SHA that a stage report or an archived entity cites stays reachable from `main` (Path 2, step 9b). The cost is a merge commit for each sync, and a history that is harder to send upstream commit by commit.

## Option 2 of 2 for "when the upstream moves": rebase

```sh
git fetch upstream && git checkout main && git rebase upstream/main && git push --force-with-lease origin main
```

Our changes stay as a linear series on top of the upstream, which is the form an upstream PR needs. The cost: the push rewrites the fork's `main`. Each SHA that a stage report cites is no longer on `main` (Path 2, step 9a). Every open worktree branch based on the old `main` must be rebased too.

No spacedock command does either sync. `state ready`, `state sweep` and the boot record do not report that `main` is behind the upstream (Path 2, step 9).
