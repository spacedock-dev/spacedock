---
id: rbjkna5jem4vgj3vtv072gzq
title: spacedock claude auto-installs the plugin when absent (--no-install opt-out)
status: ideation
source: "captain (2026-06-05) — friction F8: `spacedock claude` refuses to launch with no installed plugin, forcing the user to pass --skip-contract-check (a launch blocker). Captain direction: 'we need something simpler [than task 44]. maybe just install the plugin unless --no-install is specified in spacedock claude.' Interim relief; task 44 (bundle-into-binary) is the eventual structural fix but is deferred."
score: "0.36"
started: 2026-06-06T03:06:22Z
completed:
verdict:
worktree:
issue:
---

`spacedock claude` runs a fail-fast contract gate before launch (`internal/cli/frontdoor.go:167-171`, `gateHost(ops, "claude", stderr)`) that denies when no plugin is installed AND no `--plugin-dir` is passed. With no installed plugin, `ResolveManifest` returns empty and the gate prints "no installed claude plugin found. Run `spacedock install --host claude` (or --skip-contract-check to bootstrap)" and exits 1. No-plugin is the most common first-run state, yet the only advertised escape is a niche bootstrap flag or a separate install round-trip — so a fresh user is blocked from the one command they tried.

## Direction (for ideation)

- On `spacedock claude` when no plugin is resolvable (the NoPluginFound / empty-manifest case ONLY — NOT a version mismatch), AUTO-RUN the install (`spacedock install --host claude`, which already exists programmatically at `internal/cli/host_exec.go` installArgvSequence) and then proceed to launch — so the single command the user typed yields a working FO session.
- Gate it with a `--no-install` opt-out for users who want the old refuse-and-instruct behavior.
- This is verdict-scoped: a real version mismatch (TooOldBinary / TooOldPlugin / MalformedRange) must STILL fail fast — do not auto-install over an incompatibility. Have gateHost return the verdict (not just a bool) so runClaude can distinguish no-plugin (auto-install) from incompatible (hard-fail).
- Distinct from option A (silent launch with NO plugin -> broken session, rejected): here we INSTALL the plugin so the launch actually works.

## Out of scope

Bundling the plugin into the binary + injecting --plugin-dir (that is task 44 `bundle-asset-distribution`, deferred — it eventually makes the no-plugin case impossible and supersedes this auto-install round-trip). The binary-ABSENT journey (that is qa `spacedock-binary-missing-install-journey` — different root cause: here the binary is present, the plugin is missing). Codex (no --plugin-dir equivalent; this task is Claude's no-plugin launch).

## Acceptance criteria

**AC-1 — `spacedock claude` with no installed plugin and no flags auto-installs the plugin then launches a working session; `--no-install` preserves the refuse-and-instruct behavior.**
Verified by: a frontdoor test (mirroring TestClaudeFrontDoor* in `internal/cli/frontdoor_test.go`) where, with a no-plugin manifest resolver, runClaude invokes the install sequence then proceeds (observable: install called + launch reached); and with --no-install set, it does NOT install and exits with the instruct message. Independent source: the install-invocation + launch-reached are observed behaviors the test asserts, not a string in the file.

**AC-2 — a real version mismatch still fails fast even without --no-install (auto-install does not paper over incompatibility).**
Verified by: a frontdoor test with a TooOldBinary/TooOldPlugin manifest asserting runClaude still exits non-zero (no auto-install, no launch) — gateHost's verdict-scoped branch.

## Test plan

Go frontdoor unit tests in `internal/cli/frontdoor_test.go` over the no-plugin / --no-install / version-mismatch verdicts (the existing TestClaudeFrontDoor* harness already stubs manifest resolvers). Front-door launcher is a high-stakes surface -> detached adversarial audit before merge.
