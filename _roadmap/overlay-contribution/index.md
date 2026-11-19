# Sprint overlay-contribution

**Goal:** Let a contributor run a Spacedock workflow against a repository that they
do not own, and send a clean pull request upstream.

**Deliverable:** A named and documented shape. The contributor forks a repository,
overlays a workflow on the clone, and opens an upstream pull request. No
spacedock-authored string reaches the maintainer.

**Target release train:** not bound. The captain binds this line at scope-lock.

## Why this sprint exists

The mechanism works today. A live test drove the full loop against a fork. The
maintainer-visible diff was one file and four insertions.

Three things are missing. Nothing names the shape. No command creates it. Two
spacedock strings still reach the maintainer.

The position of the workflow directory is the load-bearing fact.
`internal/gates/prepare.go:137` builds two roots: `Main` from the workflow directory,
and `State` from the entity directory. An overlay inside the code repository makes
`Main` resolve to the code repository. Code artifacts then classify correctly. `State`
is then free to be a separate private repository.

An A/B test proved the cause. The state repository and the artifact stayed the same in
both runs. With the workflow directory in a third repository, `gate prepare` exits 1
with `selected source is not owned by a workflow Git root`. With the workflow directory
overlaid inside the code repository, the same command exits 0 and returns
`git-root://main/...`.

## The journey

This is the shape the sprint delivers. Each step gives the contributor action and the
status today.

| # | Step | Status today |
|---|------|--------------|
| 1 | Fork the upstream repository. Then clone your fork. Add the upstream as a second remote. | Manual. Spacedock has no fork concept. A repo-wide search for `fork`, `owner:branch`, and `pr-target` returns no real hit. |
| 2 | Add `spacedock/` and `.worktrees/` to `.git/info/exclude`. | Works. Discovery prunes `discoverIgnoreDirs` and tracked-gitignore patterns only. It does not prune `info/exclude`, and it does not skip untracked directories. |
| 3 | Write the workflow README at `<fork>/spacedock/flow/README.md`. | Works, by hand. No commission journey covers this shape. `status --validate` accepts unknown frontmatter keys today. |
| 4 | Create the state checkout. | **Gap.** `state new` births the orphan branch in the code repository and writes to the tracked `.gitignore`. |
| 5 | Run `spacedock status --boot --identify --json` from the fork root. | Works. Discovery fails from a parent directory that is not a git repository. |
| 6 | Run `spacedock dispatch build --stamp`. | Works. The worktree lands in the fork. **Leak:** the branch name is `spacedock-ensign/<slug>`. |
| 7 | The ensign commits the code and files the stage report. | Works. `internal/status/entered_stage.go:41-65` refuses to advance until a complete report is committed. |
| 8 | Run `spacedock gate prepare <slug> --artifact <file>.md`. | Works across two repositories. One briefing resolved a code artifact in the fork and an entity artifact in a private state repository. |
| 9 | A merge-hook mod runs `gh pr create` against the upstream. | Works today with no Go change. The binary never runs `gh pr create`. The mod body runs it. Nothing is shipped for this. |
| 10 | The mod hooks poll the pull request. `merge guard` finalizes on merge. | **Gap.** Both probes omit `--repo`. `merge guard` itself is correct and parks on an open pull request without drift. |

### CAUTION: do not clone the upstream directly

Clone your fork, or clone with `git clone -o upstream <url>`. A remote named `origin`
that you cannot write to breaks the loop. `state ready`, `state commit`,
`gate record --consume`, and `merge guard` all exit 1. `dispatch build --stamp`
deadlocks: it writes the `worktree:` field, then exits before it creates the worktree.

The cause is one predicate. `internal/statesync/publish.go:70` defines "no remote" as
"the command `git remote get-url origin` fails". It does not define it as "publication
is not possible". A clone with no remote named `origin` therefore degrades to
local-only and the full loop exits 0. Add `origin` later, and the next `state commit`
publishes the full local history.

### What the maintainer sees today

The content layer is already clean, and a live test measured it. Every ref in the fork
held three files only. A search of every blob for `spacedock`, `Stage Report`, and
`briefing` returned no match. `git status --porcelain --untracked-files=all` was empty.

Two strings still reach the maintainer. The first is the branch name, from
`internal/dispatch/stamp.go:98`. The second is the audit link in the pull request body,
which `mods/pr-merge.md:95` marks as required.

## Definition of Done

1. A contributor reaches an upstream-ready pull request through documented spacedock
   steps. The diff, the branch name, and the pull request body contain no
   spacedock-authored string. A behavior fixture proves all three.
2. `state new` leaves the code repository clean. `git status --porcelain
   --untracked-files=all` is empty with no manual repair.
3. A state branch in a separate repository can be created, and then resumed on a fresh
   clone. `state ready` exits 0 with the entity history intact.
4. A `pr` field that carries `owner/repo#N` resolves to the real state of the pull
   request. The same entity reports the same state from any launch directory.
5. A probe error during `state sweep` produces a diagnostic. It is never a silent skip.
6. A workflow with two registered merge mods either arms every mod that the contract
   runs, or the workflow is refused at validation. It never arms one and runs two.
7. A shipped merge-hook mod opens a fork-to-upstream pull request. Installation of that
   mod does not fire a second pull request.
8. A commission journey produces a working overlay with no manual repair.
9. `go test ./...` is green from the repository root.

Every fixture in this sprint uses local bare repositories and a stub for `gh`. The
probe run proved that this substrate is sufficient. No network account is needed.

## Members

Membership is a query. This file never lists the members.

```bash
spacedock status --workflow-dir docs/dev --where sprint=overlay-contribution
```

The drivable set, which excludes deferred members:

```bash
spacedock status --workflow-dir docs/dev --where sprint=overlay-contribution --where 'sprint-readiness != defer'
```

### Sequence

Two members are independent and can run first: the state placement work and the
pull-request probe work. The merge-hook arming defect is independent of both.

The shape member runs last. Its acceptance criteria depend on the state placement work
and on the merge-hook arming work.

## Out of scope

- **Greenroom.** A separate private repository for notes solves the problem of the
  repository owner, not the problem of the external contributor. Split-root state
  already keeps the reasoning off the pull request, because the orphan branch has no
  common ancestor with the trunk.
- **A fork concept in the binary.** The contributor forks by hand. This sprint does not
  add `gh repo fork` to any command.
- **A second remote model.** `gh pr create` stays in the mod body. This sprint does not
  move pull request creation into Go.
- **Survey and debrief in a foreign repository.** Both are broken outside the spacedock
  checkout. They have their own backlog entities.

## Sprint lifecycle checklist

**Shape — Shaping FO**
- [ ] **Scope-lock** with the captain. Bind the target release train in this file *(captain decides)*
- [ ] **Carve** — stamp `sprint` / `group` / `sprint-readiness` on the members. Write this `index.md`
- [ ] **Ideate** each gated member. Exercise the riskiest mechanism first, or record "no spike needed"
- [ ] **⚠️ Preflight staff review (sprint-wide)** — one independent reviewer refutes the sprint as a whole → `staff-review.md`
- [ ] **Present ideation gates** — checklist accounting and AC cross-check per member *(captain decides)*
- [ ] **Package** — write `dispatch-sprint-execution.md`

**Drive — Commander (a separate, cold-booted session)**
- [ ] **Implementation → validation → done** per member. Run the detached adversarial audit at validation for every high-stakes surface
- [ ] **Merge** each member to `main`
- [ ] **⚠️ Pre-cut antipattern audit** — one independent reviewer over the assembled sprint, before the tag fires
- [ ] **Cut the release** — `go test ./...` green, then follow [`docs/releasing.md`](../../releasing.md) *(captain authorizes)*

**Close — Shaping FO**
- [ ] **Seed the next sprint** — fold the recorded non-blocking findings into the next backlog
