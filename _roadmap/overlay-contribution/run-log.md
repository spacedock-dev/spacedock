# Run log: overlay-contribution temporary-directory run

Full transcript of the run that `answer-today.md` cites. Each section header (A0, B7, C10, ...) is the label that `answer-today.md` uses.

- Date: 2026-11-18. Binary: `spacedock 0.29.0-pre0+dev`, built with `go build ./cmd/spacedock` from `db4728199` (main).
- Root: `/tmp/overlay-run` (removed after the run).
- `remotes/upstream.git`: stands in for the paseo upstream. Read-only (`chmod -R a-w`), so a push to it fails as a push to an upstream without write access does.
- `remotes/fork.git`: a bare clone of the upstream; stands in for our fork.
- `remotes/state.git`: an empty bare repository; stands in for a separate state remote.
- `bin/gh`: a stub. It logs every call to `gh.log`, maps a URL that contains `fork` to `us/paseo` and `upstream` to `getpaseo/paseo`, and answers `pr view` from a file.
- The workflow README in every scenario has three stages: `backlog` (initial), `implementation` (worktree, gate), `done` (terminal); `state: .spacedock-state`. Scenario C adds `state-branch: dev-state`.
- No model ran. Each first-officer or ensign step is the exact command that the skill text names.
- Helper: `r CMD` prints `$ CMD`, runs it, and prints `[exit N]`.

## Scenario A: the sentence taken literally (clone the upstream)

```text
### A0 launch directory is not a git repository

$ spacedock status --boot --identify --json | head -c 600
Error: no commissioned Spacedock workflow found in /tmp/overlay-run/work/A — report this and stop; do NOT search the filesystem for one. If a workflow exists elsewhere, point at it with --workflow-dir <dir>.
[exit 0]
### A1 'clone paseo' taken literally

$ git clone -q /tmp/overlay-run/remotes/upstream.git paseo
[exit 0]

$ git -C paseo remote -v
origin	/tmp/overlay-run/remotes/upstream.git (fetch)
origin	/tmp/overlay-run/remotes/upstream.git (push)
[exit 0]

$ spacedock status --boot --identify --json | head -c 600
Error: no commissioned Spacedock workflow found in /tmp/overlay-run/work/A/paseo — report this and stop; do NOT search the filesystem for one. If a workflow exists elsewhere, point at it with --workflow-dir <dir>.
[exit 0]

$ spacedock status --boot --identify --json >/dev/null
Error: no commissioned Spacedock workflow found in /tmp/overlay-run/work/A/paseo — report this and stop; do NOT search the filesystem for one. If a workflow exists elsewhere, point at it with --workflow-dir <dir>.
[exit 1]
### A2 commission default: docs/<slug>/README.md + tracked .gitignore (commission SKILL.md lines 197, 269-273, 289-317)

$ mkdir -p docs/paseo-fork && cp ../../../scripts/README.workflow.md docs/paseo-fork/README.md
[exit 0]

$ grep -qxF '.worktrees/' .gitignore 2>/dev/null || echo '.worktrees/' >> .gitignore
[exit 0]

$ spacedock status --workflow-dir docs/paseo-fork --validate
Error: --validate requires state checkout directory declared by README state: .spacedock-state at /tmp/overlay-run/work/A/paseo/docs/paseo-fork/.spacedock-state
[exit 1]
### A3 state new (orphan state branch in the same repository, push to origin)

$ spacedock state new --workflow-dir docs/paseo-fork
Birthed split-root state at /tmp/overlay-run/work/A/paseo/docs/paseo-fork/.spacedock-state (branch spacedock-state/paseo-fork).
Next: commit the .gitignore edit on your code branch (git add .gitignore && git commit).
Warning: could not push spacedock-state/paseo-fork to origin (no remote or push failed). Push the orphan branch later to share state: git -C /tmp/overlay-run/work/A/paseo/docs/paseo-fork/.spacedock-state push origin spacedock-state/paseo-fork.
[exit 0]

$ git status --porcelain --untracked-files=all
?? .gitignore
?? docs/paseo-fork/README.md
[exit 0]

$ git branch -a
* main
+ spacedock-state/paseo-fork
  remotes/origin/HEAD -> origin/main
  remotes/origin/main
[exit 0]

$ cat .gitignore
.worktrees/
docs/paseo-fork/.spacedock-state/
[exit 0]
### A4 boot, new entity, state commit with origin = read-only upstream

$ spacedock status --boot --identify --json | head -c 400; echo
{"command":"boot","mods":{},"id_style":"slug","next_id":"n/a (id-style: slug)","orphans":[],"pr_state":{"status":"none","entries":[]},"dispatchable":[],"team_state":{"present":"false","hint":"dispatch a named background Agent for team-mode messaging (no TeamCreate needed)"},"state_backend":"split-root","definition_dir":"/private/tmp/overlay-run/work/A/paseo/docs/paseo-fork","entity_dir":"/private/
[exit 0]

$ printf -- '---\ntitle: Add banner\nstatus: backlog\n---\n\nAdd a banner to the app.\n' | spacedock new add-banner
created: /private/tmp/overlay-run/work/A/paseo/docs/paseo-fork/.spacedock-state/add-banner.md id=add-banner
[exit 0]

$ spacedock state ready --workflow-dir docs/paseo-fork
spacedock state ready: state synchronization failed:
fatal: couldn't find remote ref refs/heads/spacedock-state/paseo-fork

[exit 1]

$ spacedock state commit add-banner --workflow-dir docs/paseo-fork
spacedock state commit: state publication failed; the new local commit remains recoverable:
fatal: couldn't find remote ref refs/heads/spacedock-state/paseo-fork

[exit 1]

$ git -C docs/paseo-fork/.spacedock-state log --oneline
4163125 state: update add-banner
a4ecfab seed state
[exit 0]
### A5 dispatch with origin = read-only upstream

$ spacedock status --workflow-dir docs/paseo-fork --set add-banner status=implementation
status: backlog -> implementation
[exit 0]

$ echo 'add the banner' | spacedock dispatch build --workflow-dir docs/paseo-fork --entity-path docs/paseo-fork/.spacedock-state/add-banner.md --stage implementation --checklist-file - --stamp --host claude | head -c 300; echo
dispatch build --stamp: state publication failed; the stamp is already durable locally:
fatal: couldn't find remote ref refs/heads/spacedock-state/paseo-fork


[exit 0]

$ grep -E "^(status|worktree|started):" docs/paseo-fork/.spacedock-state/add-banner.md
status: implementation
started: 2026-11-18T18:25:15Z
worktree: .worktrees/spacedock-ensign-add-banner
[exit 0]

$ git worktree list
/private/tmp/overlay-run/work/A/paseo                                   89634c4 [main]
/private/tmp/overlay-run/work/A/paseo/docs/paseo-fork/.spacedock-state  cfd6d48 [spacedock-state/paseo-fork]
[exit 0]
### A5b retry the same dispatch

$ echo 'add the banner' | spacedock dispatch build --workflow-dir docs/paseo-fork --entity-path docs/paseo-fork/.spacedock-state/add-banner.md --stage implementation --checklist-file - --stamp --host claude >/dev/null
dispatch build --stamp: state publication failed; the stamp is already durable locally:
fatal: couldn't find remote ref refs/heads/spacedock-state/paseo-fork

[exit 1]

$ git worktree list
/private/tmp/overlay-run/work/A/paseo                                   89634c4 [main]
/private/tmp/overlay-run/work/A/paseo/docs/paseo-fork/.spacedock-state  cfd6d48 [spacedock-state/paseo-fork]
[exit 0]
### A0b exit code at a non-repository launch directory

$ spacedock status --boot --identify --json >/dev/null
Error: no commissioned Spacedock workflow found in /tmp/overlay-run/work/A — report this and stop; do NOT search the filesystem for one. If a workflow exists elsewhere, point at it with --workflow-dir <dir>.
```

## Scenario B: clone the fork, upstream as second remote, default commission layout (workflow in the tree)

```text
### B1 clone the fork, add upstream as a second remote (roadmap step 1)

$ git clone -q /tmp/overlay-run/remotes/fork.git paseo && git -C paseo remote add upstream /tmp/overlay-run/remotes/upstream.git && git -C paseo remote -v
origin	/tmp/overlay-run/remotes/fork.git (fetch)
origin	/tmp/overlay-run/remotes/fork.git (push)
upstream	/tmp/overlay-run/remotes/upstream.git (fetch)
upstream	/tmp/overlay-run/remotes/upstream.git (push)
[exit 0]
### B2 commission default layout, then commit the README and .gitignore on the code branch as commission says

$ mkdir -p docs/paseo-fork && cp ../../../scripts/README.workflow.md docs/paseo-fork/README.md
[exit 0]

$ echo '.worktrees/' >> .gitignore
[exit 0]

$ spacedock state new --workflow-dir docs/paseo-fork
Birthed split-root state at /tmp/overlay-run/work/B/paseo/docs/paseo-fork/.spacedock-state (branch spacedock-state/paseo-fork).
Next: commit the .gitignore edit on your code branch (git add .gitignore && git commit).
[exit 0]

$ git add .gitignore docs/paseo-fork/README.md && git commit -qm "Add spacedock workflow" && git push -q origin main
[exit 0]

$ git status --porcelain --untracked-files=all
[exit 0]

$ git ls-remote ../../../remotes/fork.git
17a9a80fd60a7465adcf1ddc31d662db428dbd5b	HEAD
17a9a80fd60a7465adcf1ddc31d662db428dbd5b	refs/heads/main
18e95ce6bb8bdae9e5abe9f2530d6ef12353a052	refs/heads/spacedock-state/paseo-fork
[exit 0]
### B3 engage converge + entity

$ spacedock state ready --workflow-dir docs/paseo-fork
State checkout ready — integrated peers' state from spacedock-state/paseo-fork.
[exit 0]

$ printf -- '---\ntitle: Add banner\nstatus: backlog\n---\n\nAdd a banner to the app.\n' | spacedock new add-banner
created: /private/tmp/overlay-run/work/B/paseo/docs/paseo-fork/.spacedock-state/add-banner.md id=add-banner
[exit 0]

$ spacedock state commit add-banner --workflow-dir docs/paseo-fork
Committed and pushed add-banner to spacedock-state/paseo-fork.
[exit 0]

$ spacedock status --workflow-dir docs/paseo-fork --set add-banner status=implementation && spacedock state commit add-banner --workflow-dir docs/paseo-fork
status: backlog -> implementation
Committed and pushed add-banner to spacedock-state/paseo-fork.
[exit 0]
### B4 dispatch

$ echo 'add the banner' | spacedock dispatch build --workflow-dir docs/paseo-fork --entity-path docs/paseo-fork/.spacedock-state/add-banner.md --stage implementation --checklist-file - --stamp --host claude >/dev/null
[exit 0]

$ git worktree list
/private/tmp/overlay-run/work/B/paseo                                         17a9a80 [main]
/private/tmp/overlay-run/work/B/paseo/.worktrees/spacedock-ensign-add-banner  17a9a80 [add-banner]
/private/tmp/overlay-run/work/B/paseo/docs/paseo-fork/.spacedock-state        6dea5e7 [spacedock-state/paseo-fork]
[exit 0]
### B5 ensign: code commit in the worktree, stage report in the entity

$ printf 'BANNER\n' > .worktrees/spacedock-ensign-add-banner/src/banner.txt && git -C .worktrees/spacedock-ensign-add-banner add src/banner.txt && git -C .worktrees/spacedock-ensign-add-banner commit -qm 'Add banner' && git -C .worktrees/spacedock-ensign-add-banner log --oneline -1
0c9d01f Add banner
[exit 0]

$ printf '\n## Stage Report: implementation\n\n- DONE: add the banner\n  commit on branch add-banner adds src/banner.txt\n\n### Summary\n\nAdded the banner file.\n' >> docs/paseo-fork/.spacedock-state/add-banner.md && spacedock state commit add-banner --workflow-dir docs/paseo-fork
Committed and pushed add-banner to spacedock-state/paseo-fork.
[exit 0]
### B6 gate

$ printf '# Review add-banner\n\nThe branch adds src/banner.txt.\n' > .worktrees/spacedock-ensign-add-banner/REVIEW.md
[exit 0]

$ spacedock gate prepare add-banner --workflow-dir docs/paseo-fork --question 'Ship the banner?' --artifact .worktrees/spacedock-ensign-add-banner/REVIEW.md --summary 'banner added' | tail -5
Error: --artifact /tmp/overlay-run/work/B/paseo/.worktrees/spacedock-ensign-add-banner/REVIEW.md: selected source is not the exact committed file; commit the exact source before preparation
[exit 0]

$ spacedock gate record add-banner --workflow-dir docs/paseo-fork --decision approve --actor person:captain --consume | tail -4
Error: entity has no gates record
[exit 0]
### B6 gate (artifact = the committed entity file in the state checkout)

$ spacedock gate prepare add-banner --workflow-dir docs/paseo-fork --question 'Ship the banner?' --artifact docs/paseo-fork/.spacedock-state/add-banner.md --summary 'banner added' | tail -3
briefing=briefing:add-banner:implementation:attempt-1:revision-1
digest=sha256:652c35819e50b1cccacc493425d06a4fbbbc3c05af9f31beabbe264477f2c0d6
state=open
[exit 0]

$ spacedock gate record add-banner --workflow-dir docs/paseo-fork --decision approve --actor person:captain --consume | tail -3
recorded gate=gate:add-banner:implementation attempt=gate-attempt:add-banner-implementation-1 state=closed briefing=briefing:add-banner:implementation:attempt-1:revision-1 resolution=resolution:spacedock:add-banner:implementation:1 decision=approve
sync=pushed phase=record
gate=gate:add-banner:implementation attempt=gate-attempt:add-banner-implementation-1 application=advance/pending condition=approved-pending eligible=true consumed=false target-stage=done route=approved-awaiting-merge
[exit 0]
### B7 pr-merge merge hook, step by step as mods/pr-merge.md writes it

$ BASE=$(spacedock dispatch trunk --workflow-dir docs/paseo-fork); echo "BASE=$BASE"
BASE=main
[exit 0]

$ git -C .worktrees/spacedock-ensign-add-banner remote get-url origin
/tmp/overlay-run/remotes/fork.git
[exit 0]

$ gh repo view $(git -C .worktrees/spacedock-ensign-add-banner remote get-url origin) --json nameWithOwner --jq .nameWithOwner
us/paseo
[exit 0]

$ git -C docs/paseo-fork/.spacedock-state remote get-url origin
/tmp/overlay-run/remotes/fork.git
[exit 0]

$ git -C .worktrees/spacedock-ensign-add-banner push -q origin main
[exit 0]

$ git -C .worktrees/spacedock-ensign-add-banner push -q origin $(git -C .worktrees/spacedock-ensign-add-banner rev-parse HEAD):refs/heads/add-banner
[exit 0]

$ gh pr create --repo us/paseo --base main --head add-banner --title 'Add banner' --body-file /dev/null
https://github.example/us/paseo/pull/1
[exit 0]

$ cat /tmp/overlay-run/gh.log
gh repo view /tmp/overlay-run/remotes/fork.git --json nameWithOwner --jq .nameWithOwner
gh pr create --repo us/paseo --base main --head add-banner --title Add banner --body-file /dev/null
[exit 0]
### B8 PR merged in the fork; startup hook finalizes

$ git push -q origin add-banner:main 2>&1; git -C ../../../remotes/fork.git log --oneline main
0c9d01f Add banner
17a9a80 Add spacedock workflow
89634c4 upstream: initial
[exit 0]

$ echo MERGED > /tmp/overlay-run/gh.prstate; spacedock status --workflow-dir docs/paseo-fork --set add-banner pr=#1 >/dev/null; gh pr view 1 --repo us/paseo --json state --jq .state
MERGED
[exit 0]

$ spacedock status --workflow-dir docs/paseo-fork --set add-banner pr=pr-merge:1 && spacedock state commit add-banner --workflow-dir docs/paseo-fork
pr: #1 -> pr-merge:1
Committed and pushed add-banner to spacedock-state/paseo-fork.
[exit 0]

$ spacedock merge guard add-banner --workflow-dir docs/paseo-fork --verdict passed | tail -4
finalized: add-banner -> done (verdict passed), archived. State durability: pushed to the split-root origin.
Next: push; remove the worktree (`git worktree remove .worktrees/spacedock-ensign-add-banner`, no --force — if it refuses on untracked files, audit them with the operator before any --force); delete the local branch (`git branch -d`); keep the remote branch while a PR references it; tear down the entity's workers per your runtime adapter.
[exit 0]

$ spacedock status --workflow-dir docs/paseo-fork --archived 2>&1 | tail -3
ID     SLUG                           STATUS               TITLE                          SCORE
--     ----                           ------               -----                          -----
add-banner add-banner                     done                 Add banner                     
[exit 0]
### B9 the upstream moves (maintainer pushed 'upstream: v2')

$ spacedock state ready --workflow-dir docs/paseo-fork
State checkout ready — integrated peers' state from spacedock-state/paseo-fork.
[exit 0]

$ spacedock state sweep --workflow-dir docs/paseo-fork
0 entity(ies) merged but not yet terminalized.
[exit 0]

$ spacedock status --boot --identify --json | grep -o "upstream" | head -1; echo "(no upstream mention above means none)"
(no upstream mention above means none)
[exit 0]

$ git fetch -q upstream && git log --oneline --graph --all -- . | head -0; git rev-list --left-right --count main...upstream/main
head: illegal line count -- 0
2	1
[exit 0]
### B9a option rebase

$ git branch -f rebase-try main && git rebase -q upstream/main rebase-try && git log --oneline rebase-try
a99e5c4 Add banner
5f4ef2d Add spacedock workflow
dd96437 upstream: v2
89634c4 upstream: initial
[exit 0]

$ git merge-base --is-ancestor 0c9d01f rebase-try; echo "0c9d01f (cited in the stage report) is ancestor of rebased main: $?"
0c9d01f (cited in the stage report) is ancestor of rebased main: 1
[exit 0]

$ git push --dry-run origin rebase-try:main
To /tmp/overlay-run/remotes/fork.git
 ! [rejected]        rebase-try -> main (non-fast-forward)
error: failed to push some refs to '/tmp/overlay-run/remotes/fork.git'
hint: Updates were rejected because a pushed branch tip is behind its remote
hint: counterpart. If you want to integrate the remote changes, use 'git pull'
hint: before pushing again.
hint: See the 'Note about fast-forwards' in 'git push --help' for details.
[exit 1]
### B9b option merge

$ git checkout -q main && git branch -f merge-try main && git checkout -q merge-try && git merge -q --no-edit upstream/main && git log --oneline --graph merge-try
*   de007c2 Merge remote-tracking branch 'upstream/main' into merge-try
|\  
| * dd96437 upstream: v2
* | 0c9d01f Add banner
* | 17a9a80 Add spacedock workflow
|/  
* 89634c4 upstream: initial
[exit 0]

$ git merge-base --is-ancestor 0c9d01f merge-try; echo "0c9d01f is ancestor of merged main: $?"
0c9d01f is ancestor of merged main: 0
[exit 0]

$ git push --dry-run origin merge-try:main
To /tmp/overlay-run/remotes/fork.git
   0c9d01f..de007c2  merge-try -> main
[exit 0]

$ git checkout -q main
[exit 0]
### B10 what the fork shows on top of upstream (the diff a maintainer would see in a PR from fork main)

$ git diff --stat upstream/main...main
 .gitignore                |  2 ++
 docs/paseo-fork/README.md | 52 +++++++++++++++++++++++++++++++++++++++++++++++
 src/banner.txt            |  1 +
 3 files changed, 55 insertions(+)
[exit 0]

$ git ls-remote ../../../remotes/fork.git
0c9d01f6785f5e83a499e790fa212900efd29c22	HEAD
0c9d01f6785f5e83a499e790fa212900efd29c22	refs/heads/add-banner
0c9d01f6785f5e83a499e790fa212900efd29c22	refs/heads/main
750f3519d9754fa23eb4c65dea6a6f10310e087a	refs/heads/spacedock-state/paseo-fork
[exit 0]
### B11 contribution case from the same clone: what an upstream PR from branch add-banner would carry

$ git diff --stat upstream/main...add-banner
 .gitignore                |  2 ++
 docs/paseo-fork/README.md | 52 +++++++++++++++++++++++++++++++++++++++++++++++
 src/banner.txt            |  1 +
 3 files changed, 55 insertions(+)
[exit 0]

$ git log --oneline upstream/main..add-banner
0c9d01f Add banner
17a9a80 Add spacedock workflow
[exit 0]
```

## Scenario C: overlay off the tree (cargento/asr pattern), asr remote shape (origin = upstream, fork = second remote)

```text
### C1 code clone, asr remote shape: origin = upstream, fork = second remote

$ git clone -q /tmp/overlay-run/remotes/upstream.git paseo && git -C paseo remote add fork /tmp/overlay-run/remotes/fork.git && git -C paseo remote -v
fork	/tmp/overlay-run/remotes/fork.git (fetch)
fork	/tmp/overlay-run/remotes/fork.git (push)
origin	/tmp/overlay-run/remotes/upstream.git (fetch)
origin	/tmp/overlay-run/remotes/upstream.git (push)
[exit 0]

$ printf '.spacedock/\n.worktrees/\n' >> .git/info/exclude
[exit 0]
### C2 try state new on the overlay layout (.spacedock/dev is a plain directory inside the code repo)

$ mkdir -p .spacedock/dev && sed "s/^state: .spacedock-state$/state: .spacedock-state\nstate-branch: dev-state/" ../../../scripts/README.workflow.md > .spacedock/dev/README.md && grep -n "state" .spacedock/dev/README.md | head -3
7:state: .spacedock-state
8:state-branch: dev-state
13:  states:
[exit 0]

$ spacedock state new --workflow-dir .spacedock/dev
Birthed split-root state at /tmp/overlay-run/work/C/paseo/.spacedock/dev/.spacedock-state (branch dev-state).
Next: commit the .gitignore edit on your code branch (git add .gitignore && git commit).
Warning: could not push dev-state to origin (no remote or push failed). Push the orphan branch later to share state: git -C /tmp/overlay-run/work/C/paseo/.spacedock/dev/.spacedock-state push origin dev-state.
[exit 0]

$ git status --porcelain --untracked-files=all
?? .gitignore
[exit 0]

$ git branch -a
+ dev-state
* main
  remotes/origin/HEAD -> origin/main
  remotes/origin/main
[exit 0]
### C3 manual overlay (cargento setup.sh pattern): spec repo .spacedock/repo -> state remote; per-file symlinks; dev-state orphan worktree

$ git init -q -b main .spacedock/repo && sed 's/^state: .spacedock-state$/state: .spacedock-state\nstate-branch: dev-state/' ../../../scripts/README.workflow.md > .spacedock/repo/README.md && git -C .spacedock/repo add README.md && git -C .spacedock/repo commit -qm spec && git -C .spacedock/repo remote add origin /tmp/overlay-run/remotes/state.git && git -C .spacedock/repo push -q origin main
[exit 0]

$ git -C .spacedock/repo checkout -q --orphan dev-state && git -C .spacedock/repo rm -rfq --cached . && git -C .spacedock/repo commit -q --allow-empty -m "seed state" && git -C .spacedock/repo push -q origin dev-state && git -C .spacedock/repo checkout -qf main
[exit 0]

$ mkdir -p .spacedock/dev && ln -s ../repo/README.md .spacedock/dev/README.md && git -C .spacedock/repo worktree add -q ../dev/.spacedock-state dev-state
[exit 0]

$ git status --porcelain --untracked-files=all
[exit 0]
### C4 boot and converge from the code root

$ spacedock status --boot --identify --json | head -c 300; echo
{"command":"boot","mods":{},"id_style":"slug","next_id":"n/a (id-style: slug)","orphans":[],"pr_state":{"status":"none","entries":[]},"dispatchable":[],"team_state":{"present":"false","hint":"dispatch a named background Agent for team-mode messaging (no TeamCreate needed)"},"state_backend":"split-ro
[exit 0]

$ spacedock state ready --workflow-dir .spacedock/dev
State checkout ready — integrated peers' state from dev-state.
[exit 0]

$ printf -- '---\ntitle: Add banner\nstatus: backlog\n---\n\nAdd a banner.\n' | spacedock new add-banner
created: /private/tmp/overlay-run/work/C/paseo/.spacedock/dev/.spacedock-state/add-banner.md id=add-banner
[exit 0]

$ spacedock state commit add-banner --workflow-dir .spacedock/dev
Committed and pushed add-banner to dev-state.
[exit 0]

$ spacedock status --workflow-dir .spacedock/dev --set add-banner status=implementation && spacedock state commit add-banner --workflow-dir .spacedock/dev
status: backlog -> implementation
Committed and pushed add-banner to dev-state.
[exit 0]
### C5 dispatch

$ echo 'add the banner' | spacedock dispatch build --workflow-dir .spacedock/dev --entity-path .spacedock/dev/.spacedock-state/add-banner.md --stage implementation --checklist-file - --stamp --host claude >/dev/null
[exit 0]

$ git worktree list
/private/tmp/overlay-run/work/C/paseo                                         dd96437 [main]
/private/tmp/overlay-run/work/C/paseo/.worktrees/spacedock-ensign-add-banner  dd96437 [add-banner]
[exit 0]

$ git ls-remote ../../../remotes/state.git
7aa2ea0679bf33ce89f35935804f3a1d041d70a3	HEAD
43efc286ee58a1045cadcfcb207375308e5457b2	refs/heads/dev-state
7aa2ea0679bf33ce89f35935804f3a1d041d70a3	refs/heads/main
[exit 0]
### C6 ensign + gate

$ printf 'BANNER\n' > .worktrees/spacedock-ensign-add-banner/src/banner.txt && git -C .worktrees/spacedock-ensign-add-banner add src/banner.txt && git -C .worktrees/spacedock-ensign-add-banner commit -qm 'Add banner' && git -C .worktrees/spacedock-ensign-add-banner log --oneline -1
223cd89 Add banner
[exit 0]

$ printf '\n## Stage Report: implementation\n\n- DONE: add the banner\n  commit on branch add-banner adds src/banner.txt\n\n### Summary\n\nAdded the banner file.\n' >> .spacedock/dev/.spacedock-state/add-banner.md && spacedock state commit add-banner --workflow-dir .spacedock/dev
Committed and pushed add-banner to dev-state.
[exit 0]

$ spacedock gate prepare add-banner --workflow-dir .spacedock/dev --question 'Ship?' --artifact .spacedock/dev/.spacedock-state/add-banner.md --summary 'banner added' | grep -E 'state=|git-root'
state=open
[exit 0]

$ spacedock gate record add-banner --workflow-dir .spacedock/dev --decision approve --actor person:captain --consume | tail -1
gate=gate:add-banner:implementation attempt=gate-attempt:add-banner-implementation-1 application=advance/pending condition=approved-pending eligible=true consumed=false target-stage=done route=approved-awaiting-merge
[exit 0]
### C7 pr-merge merge hook resolution with origin = upstream

$ git -C .worktrees/spacedock-ensign-add-banner remote get-url origin
/tmp/overlay-run/remotes/upstream.git
[exit 0]

$ gh repo view $(git -C .worktrees/spacedock-ensign-add-banner remote get-url origin) --json nameWithOwner --jq .nameWithOwner
getpaseo/paseo
[exit 0]

$ git -C .worktrees/spacedock-ensign-add-banner push -q origin main
[exit 0]

$ git -C .worktrees/spacedock-ensign-add-banner push -q origin $(git -C .worktrees/spacedock-ensign-add-banner rev-parse HEAD):refs/heads/add-banner
error: remote unpack failed: unable to create temporary object directory
To /tmp/overlay-run/remotes/upstream.git
 ! [remote rejected] 223cd89df5fdf85a421054893c50f9137654b627 -> add-banner (unpacker error)
error: failed to push some refs to '/tmp/overlay-run/remotes/upstream.git'
[exit 1]

$ git -C .spacedock/dev/.spacedock-state remote get-url origin
/tmp/overlay-run/remotes/state.git
[exit 0]
### C8 what the commands that do work need instead (by hand): push to the fork, PR to upstream with owner:branch

$ git -C .worktrees/spacedock-ensign-add-banner push -q fork $(git -C .worktrees/spacedock-ensign-add-banner rev-parse HEAD):refs/heads/add-banner
To /tmp/overlay-run/remotes/fork.git
 ! [rejected]        223cd89df5fdf85a421054893c50f9137654b627 -> add-banner (fetch first)
error: failed to push some refs to '/tmp/overlay-run/remotes/fork.git'
hint: Updates were rejected because the remote contains work that you do not
hint: have locally. This is usually caused by another repository pushing to
hint: the same ref. If you want to integrate the remote changes, use
hint: 'git pull' before pushing again.
hint: See the 'Note about fast-forwards' in 'git push --help' for details.
[exit 1]

$ gh pr create --repo getpaseo/paseo --base main --head us:add-banner --title 'Add banner' --body-file /dev/null
https://github.example/getpaseo/paseo/pull/1
[exit 0]

$ git diff --stat origin/main...add-banner
 src/banner.txt | 1 +
 1 file changed, 1 insertion(+)
[exit 0]

$ git status --porcelain --untracked-files=all
[exit 0]

$ git ls-remote ../../../remotes/fork.git | grep add-banner
0c9d01f6785f5e83a499e790fa212900efd29c22	refs/heads/add-banner
[exit 0]
### C8b the fork already had branch add-banner from scenario B (same slug). Delete it, then push again

$ git push -q fork :add-banner
[exit 0]

$ git -C .worktrees/spacedock-ensign-add-banner push -q fork $(git -C .worktrees/spacedock-ensign-add-banner rev-parse HEAD):refs/heads/add-banner
[exit 0]

$ git ls-remote ../../../remotes/fork.git | grep add-banner
223cd89df5fdf85a421054893c50f9137654b627	refs/heads/add-banner
[exit 0]
### C9 qualified pr field, merged upstream, finalize

$ spacedock status --workflow-dir .spacedock/dev --set add-banner pr=getpaseo/paseo#1 && spacedock state commit add-banner --workflow-dir .spacedock/dev
pr:  -> getpaseo/paseo#1
Committed and pushed add-banner to dev-state.
[exit 0]

$ echo MERGED > /tmp/overlay-run/gh.prstate; gh pr view 1 --repo getpaseo/paseo --json state --jq .state
MERGED
[exit 0]

$ spacedock state sweep --workflow-dir .spacedock/dev
1 entity(ies) merged but not yet terminalized.
next: advance each per the workflow's startup-hook advancement (_mods/).
  add-banner (PR getpaseo/paseo#1): PR merged but status=implementation
[exit 0]

$ spacedock status --workflow-dir .spacedock/dev --set add-banner pr=pr-merge:1 && spacedock state commit add-banner --workflow-dir .spacedock/dev && spacedock merge guard add-banner --workflow-dir .spacedock/dev --verdict passed | head -1
pr: getpaseo/paseo#1 -> pr-merge:1
Committed and pushed add-banner to dev-state.
finalized: add-banner -> done (verdict passed), archived. State durability: pushed to the split-root origin.
[exit 0]

$ cat /tmp/overlay-run/gh.log
gh pr view 1 --repo getpaseo/paseo --json state --jq .state
gh pr view getpaseo/paseo#1 --json state --jq .state
[exit 0]
### C10 fresh machine: code clone + spec clone + symlink, then state init

$ git clone -q /tmp/overlay-run/remotes/upstream.git paseo && cd paseo && git remote add fork /tmp/overlay-run/remotes/fork.git && printf '.spacedock/\n.worktrees/\n' >> .git/info/exclude
[exit 0]
(eval):cd:6: no such file or directory: paseo

$ git clone -q /tmp/overlay-run/remotes/state.git .spacedock/repo && mkdir -p .spacedock/dev && ln -s ../repo/README.md .spacedock/dev/README.md
[exit 0]

$ spacedock state init --workflow-dir .spacedock/dev
spacedock state init: git fetch origin dev-state failed:
fatal: couldn't find remote ref dev-state

Manual fallback: git fetch origin dev-state && git worktree add /tmp/overlay-run/work/C2/paseo/.spacedock/dev/.spacedock-state dev-state
[exit 1]

$ git branch -a
* main
  remotes/origin/HEAD -> origin/main
  remotes/origin/main
[exit 0]
### C10b same step done by hand (setup.sh step 4)

$ git -C .spacedock/repo fetch -q origin dev-state && git -C .spacedock/repo worktree add -q ../dev/.spacedock-state dev-state && spacedock state ready --workflow-dir .spacedock/dev
State checkout ready — integrated peers' state from dev-state.
[exit 0]

$ spacedock status --workflow-dir .spacedock/dev --archived | tail -1
add-banner add-banner                     done                 Add banner                     
[exit 0]
### C11 maintained fork with asr remotes: local main follows fork/main; what the merge hook's trunk push does

$ git fetch -q fork && git checkout -q -B main fork/main && git log --oneline -3
0c9d01f Add banner
17a9a80 Add spacedock workflow
89634c4 upstream: initial
[exit 0]

$ spacedock dispatch trunk --workflow-dir .spacedock/dev
main
[exit 0]

$ git push --dry-run origin main
To /tmp/overlay-run/remotes/upstream.git
 ! [rejected]        main -> main (non-fast-forward)
error: failed to push some refs to '/tmp/overlay-run/remotes/upstream.git'
hint: Updates were rejected because the tip of your current branch is behind
hint: its remote counterpart. If you want to integrate the remote changes,
hint: use 'git pull' before pushing again.
hint: See the 'Note about fast-forwards' in 'git push --help' for details.
[exit 1]

$ git push --dry-run fork main
Everything up-to-date
[exit 0]
```

## Scenario D: overlay off the tree with origin = fork, created and resumed by `overlay-setup.sh`

The script is in `suggested-phrasing.md`. `state2.git` is a second empty state remote. The fork's `main` still carries scenario B's in-tree workflow, so `new` needs `--workflow-dir`.

```text
### D1 first setup with overlay-setup.sh (origin = fork)

$ /tmp/overlay-run/scripts/overlay-setup.sh /tmp/overlay-run/remotes/upstream.git /tmp/overlay-run/remotes/fork.git /tmp/overlay-run/remotes/state2.git paseo /tmp/overlay-run/scripts/README.overlay.md
State checkout ready — integrated peers' state from dev-state.
[exit 0]

$ git remote -v; git status --porcelain --untracked-files=all
origin	/tmp/overlay-run/remotes/fork.git (fetch)
origin	/tmp/overlay-run/remotes/fork.git (push)
upstream	/tmp/overlay-run/remotes/upstream.git (fetch)
upstream	/tmp/overlay-run/remotes/upstream.git (push)
[exit 0]
### D2 maintained fork loop: entity, dispatch, ensign, gate, merge hook (origin = fork)

$ git checkout -q main && git log --oneline -1
0c9d01f Add banner
[exit 0]

$ printf -- '---\ntitle: Add footer\nstatus: implementation\n---\n\nAdd a footer.\n' | spacedock new add-footer && spacedock state commit add-footer --workflow-dir .spacedock/dev
Error: multiple commissioned Spacedock workflows found under /private/tmp/overlay-run/work/D/paseo; pass --workflow-dir <dir> to pick one:
  /private/tmp/overlay-run/work/D/paseo/.spacedock/dev
  /private/tmp/overlay-run/work/D/paseo/docs/paseo-fork
[exit 1]

$ echo 'add the footer' | spacedock dispatch build --workflow-dir .spacedock/dev --entity-path .spacedock/dev/.spacedock-state/add-footer.md --stage implementation --checklist-file - --stamp --host claude >/dev/null && git worktree list
dispatch build --stamp: entity file not readable at '/tmp/overlay-run/work/D/paseo/.spacedock/dev/.spacedock-state/add-footer.md'
[exit 1]

$ printf 'FOOTER\n' > .worktrees/spacedock-ensign-add-footer/src/footer.txt && git -C .worktrees/spacedock-ensign-add-footer add src/footer.txt && git -C .worktrees/spacedock-ensign-add-footer commit -qm 'Add footer' && printf '\n## Stage Report: implementation\n\n- DONE: add the footer\n  commit adds src/footer.txt\n\n### Summary\n\nAdded the footer.\n' >> .spacedock/dev/.spacedock-state/add-footer.md && spacedock state commit add-footer --workflow-dir .spacedock/dev
(eval):1: no such file or directory: .worktrees/spacedock-ensign-add-footer/src/footer.txt
[exit 1]

$ spacedock gate prepare add-footer --workflow-dir .spacedock/dev --question 'Ship?' --artifact .spacedock/dev/.spacedock-state/add-footer.md --summary footer | grep state= && spacedock gate record add-footer --workflow-dir .spacedock/dev --decision approve --actor person:captain --consume | tail -1
Error: Error: unknown reference: add-footer
[exit 1]

$ CODE_REPO=$(gh repo view $(git -C .worktrees/spacedock-ensign-add-footer remote get-url origin) --json nameWithOwner --jq .nameWithOwner); echo $CODE_REPO; git -C .worktrees/spacedock-ensign-add-footer push -q origin main && git -C .worktrees/spacedock-ensign-add-footer push -q origin $(git -C .worktrees/spacedock-ensign-add-footer rev-parse HEAD):refs/heads/add-footer && gh pr create --repo $CODE_REPO --base main --head add-footer --title 'Add footer' --body-file /dev/null
fatal: cannot change to '.worktrees/spacedock-ensign-add-footer': No such file or directory
unknown/unknown
fatal: cannot change to '.worktrees/spacedock-ensign-add-footer': No such file or directory
[exit 128]

$ gh repo view $(git -C .spacedock/dev/.spacedock-state remote get-url origin) --json nameWithOwner --jq .nameWithOwner
us/paseo-state
[exit 0]

$ git status --porcelain --untracked-files=all; git ls-remote ../../../remotes/fork.git
0c9d01f6785f5e83a499e790fa212900efd29c22	HEAD
223cd89df5fdf85a421054893c50f9137654b627	refs/heads/add-banner
0c9d01f6785f5e83a499e790fa212900efd29c22	refs/heads/main
750f3519d9754fa23eb4c65dea6a6f10310e087a	refs/heads/spacedock-state/paseo-fork
[exit 0]
### D3 resume on a fresh machine with the same script (no README argument)

$ /tmp/overlay-run/scripts/overlay-setup.sh /tmp/overlay-run/remotes/upstream.git /tmp/overlay-run/remotes/fork.git /tmp/overlay-run/remotes/state2.git paseo2
State checkout ready — integrated peers' state from dev-state.
[exit 0]

$ spacedock status --workflow-dir paseo2/.spacedock/dev | tail -1
--     ----                           ------               -----                          -----
[exit 0]
### D2 (rerun) maintained fork loop with --workflow-dir (the fork's main still has scenario B's in-tree workflow, so new sees two)

$ printf -- '---\ntitle: Add footer\nstatus: implementation\n---\n\nAdd a footer.\n' | spacedock new add-footer --workflow-dir .spacedock/dev && spacedock state commit add-footer --workflow-dir .spacedock/dev
created: /tmp/overlay-run/work/D/paseo/.spacedock/dev/.spacedock-state/add-footer.md id=add-footer
Committed and pushed add-footer to dev-state.
[exit 0]

$ echo 'add the footer' | spacedock dispatch build --workflow-dir .spacedock/dev --entity-path .spacedock/dev/.spacedock-state/add-footer.md --stage implementation --checklist-file - --stamp --host claude >/dev/null && git worktree list
/private/tmp/overlay-run/work/D/paseo                                         0c9d01f [main]
/private/tmp/overlay-run/work/D/paseo/.worktrees/spacedock-ensign-add-footer  0c9d01f [add-footer]
[exit 0]

$ printf 'FOOTER\n' > .worktrees/spacedock-ensign-add-footer/src/footer.txt && git -C .worktrees/spacedock-ensign-add-footer add src/footer.txt && git -C .worktrees/spacedock-ensign-add-footer commit -qm 'Add footer' && printf '\n## Stage Report: implementation\n\n- DONE: add the footer\n  commit adds src/footer.txt\n\n### Summary\n\nAdded the footer.\n' >> .spacedock/dev/.spacedock-state/add-footer.md && spacedock state commit add-footer --workflow-dir .spacedock/dev
Committed and pushed add-footer to dev-state.
[exit 0]

$ spacedock gate prepare add-footer --workflow-dir .spacedock/dev --question 'Ship?' --artifact .spacedock/dev/.spacedock-state/add-footer.md --summary footer | grep state= && spacedock gate record add-footer --workflow-dir .spacedock/dev --decision approve --actor person:captain --consume | tail -1
state=open
gate=gate:add-footer:implementation attempt=gate-attempt:add-footer-implementation-1 application=advance/pending condition=approved-pending eligible=true consumed=false target-stage=done route=approved-awaiting-merge
[exit 0]

$ CODE_REPO=$(gh repo view $(git -C .worktrees/spacedock-ensign-add-footer remote get-url origin) --json nameWithOwner --jq .nameWithOwner); echo $CODE_REPO; git -C .worktrees/spacedock-ensign-add-footer push -q origin main && git -C .worktrees/spacedock-ensign-add-footer push -q origin $(git -C .worktrees/spacedock-ensign-add-footer rev-parse HEAD):refs/heads/add-footer && gh pr create --repo $CODE_REPO --base main --head add-footer --title 'Add footer' --body-file /dev/null
us/paseo
https://github.example/us/paseo/pull/1
[exit 0]

$ gh repo view $(git -C .spacedock/dev/.spacedock-state remote get-url origin) --json nameWithOwner --jq .nameWithOwner
us/paseo-state
[exit 0]

$ git status --porcelain --untracked-files=all
[exit 0]

$ git ls-remote ../../../remotes/fork.git | grep -v spacedock-state; git ls-remote ../../../remotes/state2.git
0c9d01f6785f5e83a499e790fa212900efd29c22	HEAD
223cd89df5fdf85a421054893c50f9137654b627	refs/heads/add-banner
2429f7ab043c470213d4dcceaccf68a79fa87fbf	refs/heads/add-footer
0c9d01f6785f5e83a499e790fa212900efd29c22	refs/heads/main
3731a2d8905bc153ce3fefb7d67222ee74e01552	HEAD
efbfe1948a716b6ab04e24220937dab030bfc65f	refs/heads/dev-state
3731a2d8905bc153ce3fefb7d67222ee74e01552	refs/heads/main
[exit 0]

$ cat /tmp/overlay-run/gh.log
gh repo view /tmp/overlay-run/remotes/fork.git --json nameWithOwner --jq .nameWithOwner
gh pr create --repo us/paseo --base main --head add-footer --title Add footer --body-file /dev/null
gh repo view /tmp/overlay-run/remotes/state2.git --json nameWithOwner --jq .nameWithOwner
[exit 0]
### D3b the resumed clone sees the entity

$ spacedock state ready --workflow-dir ../paseo2/.spacedock/dev && spacedock status --workflow-dir ../paseo2/.spacedock/dev | tail -1
State checkout ready — integrated peers' state from dev-state.
add-footer add-footer                     implementation       Add footer                     
[exit 0]
### D4 contribution from the same overlay: check out upstream/main before dispatch, PR to upstream by hand

$ git fetch -q upstream && git checkout -q -B upstream-main upstream/main && git log --oneline -1
dd96437 upstream: v2
[exit 0]

$ printf -- '---\ntitle: Fix typo\nstatus: implementation\n---\n\nFix a typo.\n' | spacedock new fix-typo --workflow-dir .spacedock/dev && spacedock state commit fix-typo --workflow-dir .spacedock/dev
created: /tmp/overlay-run/work/D/paseo/.spacedock/dev/.spacedock-state/fix-typo.md id=fix-typo
Committed and pushed fix-typo to dev-state.
[exit 0]

$ echo 'fix the typo' | spacedock dispatch build --workflow-dir .spacedock/dev --entity-path .spacedock/dev/.spacedock-state/fix-typo.md --stage implementation --checklist-file - --stamp --host claude >/dev/null && git -C .worktrees/spacedock-ensign-fix-typo log --oneline -1
dd96437 upstream: v2
[exit 0]

$ printf 'paseo!\n' > .worktrees/spacedock-ensign-fix-typo/README.md && git -C .worktrees/spacedock-ensign-fix-typo commit -qam 'Fix typo'
[exit 0]

$ git -C .worktrees/spacedock-ensign-fix-typo push -q origin HEAD:refs/heads/fix-typo && gh pr create --repo getpaseo/paseo --base main --head us:fix-typo --title 'Fix typo' --body-file /dev/null
https://github.example/getpaseo/paseo/pull/1
[exit 0]

$ git diff --stat upstream/main...fix-typo; git log --oneline upstream/main..fix-typo
 README.md | 2 +-
 1 file changed, 1 insertion(+), 1 deletion(-)
f19262a Fix typo
[exit 0]

$ spacedock dispatch trunk --workflow-dir .spacedock/dev
main
[exit 0]

$ git status --porcelain --untracked-files=all; git checkout -q main
[exit 0]
```

## Skill text search

```text
[exit 1]
### S1 skill text: does commission or the FO contract mention clone, fork, upstream, or an overlay?

$ grep -n -i -E 'clone|fork|upstream|overlay|info/exclude' skills/commission/SKILL.md skills/first-officer/SKILL.md skills/first-officer/references/*.md | grep -v -i 'cloned split-root\|fresh clone' | head
skills/first-officer/references/claude-fo-dispatch.md:15:The spawn call (fo-dispatch-core.md `## Dispatch Adapter`) is the Agent tool. **Use Agent() for initial dispatch** — SendMessage is only for advancing a reused agent to its next stage. **NEVER use `subagent_type="first-officer"`** — that clones yourself instead of dispatching a worker. `«async-dispatch»` binds ASYNC here: `Agent(name=…, run_in_background=true)` returns immediately; single-entity/bare mode uses the blocking call and dispatches one entity at a time.
skills/first-officer/references/codex-first-officer-runtime.md:9:- `«worker.spawn»`: `spawn_agent(task_name,message,fork_turns)`
skills/first-officer/references/codex-first-officer-runtime.md:19:- `«worker.spawn»`: A successful `«dispatch.build»` is not a dispatch by itself. For every ready entity, including one just advanced by gate approval, call `spawn_agent(task_name,message,fork_turns="none")` with the intact helper-emitted pointer as `message`, optionally followed by ordinary FO scope instructions. Every spawn is a fresh dispatch; deliberate continuity uses `followup_task` with the existing handle. sanitize the helper-emitted `name` to a lowercase digit/underscore `task_name`; retain the mapping from helper identity to sanitized task path, entity slug, stage, and cycle.
skills/first-officer/references/codex-first-officer-runtime.md:20:  Build fresh Codex dispatches in the named shape: omit `--bare-mode`. A zero-exit build must carry `name` and `prompt`; map them through `CodexMultiAgentV2SpawnInput` to `spawn_agent(task_name,message,fork_turns="none")`.
[exit 0]

$ grep -n 'Journey' skills/commission/SKILL.md
289:#### Journey 1 — Scaffold a split-root workflow (orphan-branch state, same repo)
319:#### Journey 2 — Scaffold an inline workflow
[exit 0]
```
