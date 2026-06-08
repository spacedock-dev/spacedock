---
id: 4tqngghpq91dsc17g0njh25m
title: survey — agentsview install-detection false-negatives under sandbox (asks to install when `agentsview --version` works)
status: backlog
source: "captain manual install-path test (2026-06-08) — running /spacedock:survey under sandbox, the skill detected agentsview as NOT installed and asked to `brew install --cask agentsview`, but `agentsview --version` → v0.32.1 (installed). The detection is sandbox-flawed."
score: "0.27"
started:
completed:
verdict:
worktree:
issue:
sprint: 0198-pre-flip-hardening
group: survey
sprint-readiness: ready
---

Under sandbox, the survey skill's agentsview-installed detection gives a false negative — it reports agentsview missing (and prompts to install) even though `agentsview --version` succeeds.

## Problem

- `/spacedock:survey` checks whether agentsview is installed before reading the local session history.
- Under sandbox, that check returns "not installed" → the skill asks to `brew install --cask agentsview` — but `agentsview --version` reports `v0.32.1` (it IS installed).
- The detection mechanism (whatever it probes) does not survive the sandbox, while the actual `agentsview` invocation does. So the skill nags to install an already-present tool and can derail a read-only survey.

## Proposed approach (ideation firms)

Make the detection robust under sandbox: probe agentsview the way the survey actually uses it (e.g. the `agentsview --version` exit code, or the real read path) rather than whatever sandbox-fragile check is used now. Ground the fix against the actual sandbox behavior — what the current check does vs what `--version` does.

## Acceptance criteria (sketch)

- Under sandbox, with agentsview installed, the survey detects it as present and does NOT prompt to install — verified by a live drive under sandbox (the survey proceeds without the install prompt), since the failure only manifests sandboxed.

## Notes

Survey-skill bug; same area as xn / 1p27 / 69rk. Surfaced by captain manual install-path testing. 0198 survey group.
