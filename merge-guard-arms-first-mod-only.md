---
id: wpbsnwn2ngyy4attg53y2qgh
title: Reconcile merge-hook arming with the contract that runs every registered merge mod
status: backlog
source: FO research session 2026-08-18 - overlay-contribution probe (workflow run wf_40ea0f6e-aa8)
started:
completed:
verdict:
score:
worktree:
issue:
sprint: overlay-contribution
group: merge-hooks
---

`merge guard` arms exactly one merge mod while the FO contract runs all of them, so a workflow with two registered merge hooks gets one `mod-block` naming one mod and two hooks executing.

## Problem

`internal/status/merge.go:188` arms `mergeHooks[0]` - the alphabetically first registered merge mod:

    case modBlock == "" && hookRegistered:
        return arm(roots, slug, mergeHooks[0], quiet, asJSON, stdout, stderr)

The shared FO contract at `skills/first-officer/references/first-officer-shared-core.md:151` specifies the opposite for the run side: «hooks.run»(point) "run[s] the registered point alphabetically by mod filename, exactly once" - that is, every registered mod for the point.

Observed in a fixture, with `merge guard --verdict passed` run three times against the same workflow:

- only `fork-pr.md` installed: `mods.merge = ["fork-pr"]`, armed `merge:fork-pr`, exit 0.
- `fork-pr.md` plus the shipped `pr-merge.md`: `mods.merge = ["fork-pr","pr-merge"]`, armed `merge:fork-pr` only, exit 0.
- the same two with the custom mod renamed `zz-fork-pr.md`: `mods.merge = ["pr-merge","zz-fork-pr"]`, armed `merge:pr-merge` only, exit 0.

Nothing warns. For a merge hook whose body opens a pull request, the consequence of the second and third cases is two `gh pr create` calls from one merge boundary. The defect is independent of the overlay-contribution work that surfaced it; any workflow with two merge mods is exposed, and the current workaround - install exactly one merge mod - is undocumented.

The right resolution is a design question this seed does not settle: arm every registered merge mod and let `mod-block` carry a set, or declare single-merge-hook a validated constraint and fail `status --validate` when a workflow registers more than one. Ideation owns the choice.

## Out of scope

The content of any particular merge mod, and the overlay-contribution shape that surfaced this.

## Acceptance criteria

**AC-1 (VALUE) - A workflow registering two merge mods either runs both hooks with both armed, or is refused at validation; it never silently arms one and runs two.**
Verified by: a behavior fixture with two registered merge mods asserting that the armed set matches the set the contract runs, or that `status --validate` exits non-zero with a message naming both mods. Fails today: `mods.merge` lists two, `arm` names one, exit 0 and silent.

**AC-2 - The arming behavior and the contract's «hooks.run» text describe the same rule.**
Verified by: the fixture in AC-1 exercising the behavior; the contract edit alone is not evidence. Paired with AC-1 per the workflow's mechanism-serves-value rule.

## Test plan

Go unit tests over the arming branch at `merge.go:179-188` with one, two, and zero registered merge mods. A command-level behavior fixture for AC-1 driving `merge guard` against a workflow with two `_mods/` files. If the resolution changes the shipped FO contract text, the Claude live lane is required per the workflow's path-to-lane mapping.
