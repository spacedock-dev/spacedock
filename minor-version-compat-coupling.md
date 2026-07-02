---
title: Replace the contract-integer compatibility gate with minor-version coupling (self-maintaining skill↔binary compat)
status: ideation
score: 0.6
source: "captain proposal during the 0.23.0 cut (2026-06-22). The contract integer stayed dead at 1 across 0.20/0.22/0.23 and was never bumped when the interface broke (the v0.23.0 skill-skew bug). contract-version-bump-v2 (7h3) is the tactical integer bump this would replace."
id: kr7s9efxas4fqhtcj94fs4br
started: 2026-07-02T09:55:07Z
---

Replace the manual contract-integer gate (binary reports `contract <N>`; plugin declares `requires-contract: >=X,<Y`) with minor-version coupling: 0.X.* skills require a >=0.X.0 binary (same-minor compatible; patches interchangeable). The integer is a hand-maintained value that got forgotten across three minors — minor-coupling moves the floor automatically each release, so the skill-skew class (old binary + new skills booting clean, then breaking on missing verbs) cannot recur.

## Problem

The contract integer decouples interface-version from release-version IN THEORY, but in PRACTICE it stayed `1` across 0.20/0.22/0.23 even when the skill↔binary interface broke (v0.23.0 added hard deps on `state commit/ready/sweep` + `merge guard`, absent in contract-1 binaries). The gate gave a false-green and the FO broke cryptically. `contract-version-bump-v2` (7h3) ships the integer bump to 2 for v0.23.0 as a tactical fix; this task replaces the mechanism so the bump is never forgotten again.

## Proposed approach (design pass needed)

- The plugin declares a binary version floor (e.g. `requires-binary: >=0.23.0`, or compatibility keyed on the minor) instead of `requires-contract`.
- The boot gate parses the binary's SEMVER from `--version` (already reported: `spacedock 0.23.0`) and checks `binary >= floor` (and/or same-minor), replacing the `contract <N>` integer check.
- Decide the exact semantics: minor-exact-match vs version-floor; forward-compat for a newer binary; how a patch-only skills update interacts.
- Migrate `internal/contract` (the gate + Compare math), the plugin manifest schema, the migration/range-bracketing test (#410), and the FO/ensign contract text (shared-core Startup step 1) off the integer.
- Provide a clean migration from `requires-contract: >=2,<3` (what 0.23.0 ships) to the minor-version floor.

## Acceptance criteria

- **AC-1 (VALUE)** — an old-minor binary is rejected by new-minor skills with the clean "binary too old, upgrade" abort, WITHOUT any hand-maintained integer to bump: verified by a test where bumping only the release version (no separate integer edit) already moves the compatibility floor.
- **AC-2** — patch-level skew stays compatible (0.X.0 skills ↔ 0.X.1 binary), minor-level skew is rejected; proven against fixtures and ideally the real published binaries of two minors.
- **AC-3** — the contract-integer mechanism is fully removed (no `CONTRACT_VERSION` integer, no `requires-contract` range) with a documented migration; the FO/ensign contract text describes the minor-version gate.

## Test plan

- Go tests for the semver-floor gate (reject old minor, accept same/newer); patch-interchange fixtures; the divergeable floor check.
- Real-binary cross-minor check (published v0.22.0 vs v0.23.0).
- Detached adversarial audit (contract/release surface).
- Needs an ideation pass on the exact semantics (minor-exact vs floor, forward-compat) before implementation.
