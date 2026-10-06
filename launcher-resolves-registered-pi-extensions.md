---
title: The Pi launcher demands an extension path that Pi itself already resolves from settings
status: ideation
score: 0.8
source: "FO review of the Pi stack, 2026-10-04: internal/cli/pi.go:645-648 resolves pi-subagents from the env var or HOME, while Pi resolves a registered npm entry from the agent directory."
id: 7qksxxpbdcqz96mh0basxxzb
---

## Problem

The launcher does not use the same resolution as the host it launches.

`piRuntimeConfigFromEnv` in `internal/cli/pi.go` reads `PI_SUBAGENTS_PACKAGE_ROOT`,
and when that is unset it falls back to
`join(HOME, ".pi/agent/npm/node_modules/pi-subagents")`. Pi itself resolves the
`npm:pi-subagents` entry in the agent directory's `settings.json` to
`<agentDir>/npm/node_modules/pi-subagents`. Those are different roots whenever
`HOME` is not the operator's real home.

So an operator with an isolated or sandboxed home, or with
`PI_CODING_AGENT_DIR` set, gets the setup-incomplete block and its step 5 telling
them to export the variable, although Pi can already find the packages. The
doctor probes were fixed for this class in #817; the launch path was not, and the
hint is still printed from `internal/cli/pi.go:472`.

## Visible value

An operator whose home is not the default runs `spacedock pi` and both
extensions resolve from the registration Pi already reads, with neither
package-root variable exported.

## Out of scope

- Pi's own resolution rules.
- The isolated-home harness contract, which is folded into the Pi lane work.
- Installing the packages, and the plugin and skill paths.

## Expected surface and tolerance

`internal/cli/pi.go`: the launch-path resolution and its source label.
`internal/cli/pi_frontdoor_test.go`: the resolution tests. Estimate net +80,
across 2 files. Tolerance net +40 to +140, at most 3 files.

## Acceptance criteria

**AC-1 (VALUE) - Both roots resolve from the registration with no variable exported.**
Verified by: a test that registers both npm entries in a settings.json inside a
temporary agent directory, points `HOME` elsewhere, exports neither variable, and
requires both roots to resolve under that agent directory's npm path.
Falsifier: remove the registration read and that test must fail.

**AC-2 - An exported override still wins.**
Verified by: the existing override tests. Falsifier: ignore the variable and they
must fail.

**AC-3 - A missing install still fails closed.**
Verified by: the setup-incomplete path, with no registration and no variables.
Falsifier: return a path that does not exist and the test must fail.

**AC-4 (no-regression) -** The doctor's probes and the setup message keep their
current behaviour, and the offline suite passes.

## Test plan

Deterministic tests only. No live run. Every test must be able to fail for a
reason other than editing its own expectation.
