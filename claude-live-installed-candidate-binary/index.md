---
title: Install the checked-out candidate into a neutral bin for Claude live E2E
status: validation
source: "Captain direction 2026-07-13 after Opus treated the checkout-shaped SPACEDOCK_BIN path as a workflow root."
score: 0.85
started: 2026-07-13T04:20:06Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-claude-live-installed-candidate-binary
issue:
id: 2h786rdtw1f26q7q98v5bp43
---

The Claude live job currently builds `./spacedock`, exports `SPACEDOCK_BIN=$(pwd)/spacedock`, and adds the checkout to PATH. The front door then deliberately re-exports its own executable path to the first officer. In a real Opus failure, that repository-shaped path became the apparent project root even though the host process started in a fixture.

## Problem

The live candidate must remain the exact checked-out release SHA, but its binary path should not visually identify the repository as a workflow root. Public package installation, Homebrew, and the release install script cannot prove the candidate commit. A symlink is also invalid because the launcher resolves it back to the checkout.

## Proposed approach

Change only the Claude live job's candidate-binary setup and its existing binary
artifact entry:

1. Create an empty `$RUNNER_TEMP/spacedock-live-bin` and fail if its expected
   `spacedock` output already exists (including as a dangling symlink).
2. Run `GOBIN="$candidate_dir" go install ./cmd/spacedock` from the checked-out
   repository. This compiles the local package; it does not consult a module
   proxy, Homebrew, or a release installer.
3. Export that physical `$candidate_dir/spacedock` through both
   `SPACEDOCK_BIN` and `GITHUB_PATH`. The live runners therefore resolve the
   same candidate explicitly, and subsequent CI steps can resolve it by PATH.
4. Before exporting it, fail closed unless all of these hold:
   - the output is executable and not a symlink;
   - its canonical path is outside the checkout's canonical path; and
   - `go version -m` reports `vcs.revision` equal to `git rev-parse HEAD` and
     `vcs.modified=false`.
5. Write the observed values (candidate path, canonical path, checkout
   revision, embedded revision, plus modified state) to
   `live-artifacts/claude/<model>/candidate-binary-provenance.txt`, which the
   existing artifact upload already retains. Replace the obsolete `./spacedock`
   artifact with the installed candidate binary, preserving the prior ability
   to inspect the exact launcher after a failure.

The existing `livePluginDir` path remains unchanged: it stages the checkout's
`.claude-plugin`, `skills`, and optional `agents` bytes into a workflow-free
temporary plugin directory, then the live runner passes that directory via
`--plugin-dir`. This task changes where the launcher executable lives, not
which plugin bytes Claude receives.

This is harness isolation, not a substitute for the separate
launch-cwd-authority task. It removes a repository-shaped executable path from
the model's environment; it does not make a semantic promise that a model can
never choose a wrong cwd.

## Out of scope

- Public/brew/release-script installation.
- Symlink-based relocation.
- Removing `SPACEDOCK_BIN` propagation or changing the launcher invariant.
- Treating this as a semantic guarantee that a model will never choose a wrong cwd.

## Acceptance criteria

- **AC-1 (VALUE):** Claude live E2E exercises a physical candidate binary outside the checkout while continuing to run the exact workflow-dispatch SHA and current local plugin bytes.
  - Verified by: a final-SHA Claude live job whose provenance artifact binds the
    installed binary's embedded revision to the checked-out revision, plus the
    existing durable live scenarios (fixture entity state, state-checkout log,
    and clean-state assertions) passing through that binary.
- **AC-2:** The host-visible binary path no longer names the repository checkout, and no symlink resolves it back there.
  - Verified by: the live setup's executable/not-symlink and canonical
    candidate-outside-canonical-workspace checks, recorded in the provenance
    artifact—not instruction-text matching.
- **AC-3:** The change does not weaken the separate wrong-root detector or claim to solve launch-cwd authority.
  - Verified by: bounded-diff review and the existing wrong-root controls;
    `internal/ensigncycle` detector code and the first-officer contract remain
    untouched.

## Riskiest-mechanism spike

No unverified installation mechanism remains. On 2026-07-13, from the clean
checked-out `aeeadf451c80748215de344ca8b814eca4e81afa`, a local
`GOBIN=<temp>/spacedock-live-bin go install ./cmd/spacedock` produced a regular
executable outside the checkout. `go version -m` reported
`vcs.revision=aeeadf451c80748215de344ca8b814eca4e81afa` and
`vcs.modified=false`, matching `git rev-parse HEAD`.

The spike also rejected a tempting but wrong provenance check: on macOS,
`realpath /tmp/...` is `/private/tmp/...`. The durable invariant is therefore
that the canonical candidate path lies outside the canonical workspace path,
not that canonicalization preserves the original spelling. That comparison is
portable to the Ubuntu live runner and catches a symlinked route back into the
checkout.

## Test plan

No new Go or prose/config test is warranted: the behavior is GitHub Actions'
local candidate install, so a static workflow-shape test would only re-state the
shell. The setup step itself parses the binary's independently emitted build
metadata and fails before model spend if provenance or placement is wrong.

Run the affected Claude live lane on the exact candidate SHA and inspect:

- `candidate-binary-provenance.txt` and the uploaded installed binary;
- the existing `TestLiveDefaultHeadlessStopsAtGate` durable fixture result and
  wrong-root controls; and
- the normal full Runtime Live E2E result for release certification.

No public documentation diff is needed: this is an internal CI harness change,
and the local developer recipe remains intentionally free to use a local build.

## Stage Report: ideation

- DONE: Narrow the implementation to the Claude job's candidate install,
  provenance artifact, and existing binary artifact path; preserve staged local
  plugin bytes and every runtime/frontdoor semantic.
- DONE: Spike local-source `GOBIN` installation against the actual checkout.
  The candidate was a regular executable outside the checkout; embedded
  `vcs.revision` exactly matched `aeeadf451c80748215de344ca8b814eca4e81afa`
  and `vcs.modified=false`.
- DONE: Define behavior-level evidence without Contractlint or a prose-grep.
  The live setup compares structured `go version -m` metadata and canonical
  paths, writes an uploaded provenance artifact, and the existing durable live
  scenarios remain the runtime proof.

### Summary

The smallest repair is a local `go install` into a neutral runner-temp bin for
the Claude job only. It preserves the current checkout's staged plugin while
making the launcher path non-checkout-shaped, with an artifact that proves the
installed binary embeds the actual checked-out revision. The design deliberately
does not alter the wrong-root detector, first-officer wording, or launch-cwd
authority work.

## Stage Report: implementation

- DONE: Install the checked-out candidate with local go install into a physical runner-temp bin outside the checkout, failing closed on executable, symlink, canonical-path, revision, and modified-state checks.
  Commit `5ff92f9b` installs with `GOBIN=$RUNNER_TEMP/spacedock-live-bin`; focused tests exercise success plus dangling-symlink, inside-checkout, dirty-state, and revision-mismatch rejection.
- DONE: Export the same installed candidate through SPACEDOCK_BIN and GITHUB_PATH, preserve staged local plugin bytes, and replace the obsolete binary artifact with provenance plus the installed executable.
  The workflow exports the verified path through both GitHub environment files, writes `candidate-binary-provenance.txt`, uploads the runner-temp executable, and leaves plugin staging unchanged.
- DONE: Run the affected checks, commit the bounded CI change, and report evidence that the embedded revision equals the checkout SHA while wrong-root controls remain untouched.
  Focused tests, `go test ./...`, and `go test ./... -race` passed; a clean standalone checkout produced matching `5ff92f9b58e2aae918a9cbbda560eb7394dae413` revisions with `vcs.modified=false`, and the diff contains no `internal/ensigncycle` or first-officer contract files.

### Summary

Claude live E2E now installs the exact checked-out candidate into a physical runner-temp bin, verifies its placement and Go VCS metadata before model spend, and retains both provenance and executable artifacts. Behavioral workflow tests cover the positive path and the principal fail-closed boundaries without changing staged plugin bytes or wrong-root detection.

## Stage Report: validation

- DONE: Independently audit commit 5ff92f9b and reproduce the focused positive and fail-closed candidate-install tests plus full and race suites; verify the branch is clean and the diff stays bounded to the Claude live harness.
  Focused install tests passed 5/5; `go test ./...`, `go test ./... -race`, and `git diff --check` passed; the clean two-file diff is only the Claude workflow and its behavioral test.
- FAILED: AC-1 — Claude live E2E exercises a physical candidate binary outside the checkout while continuing to run the exact workflow-dispatch SHA and current local plugin bytes.
  No origin branch, final-SHA GitHub run, durable scenario result, or uploaded live artifact exists; an exact-`5ff92f9b` standalone clone proved local install/provenance, but absent live evidence is not green.
- DONE: AC-2 — The host-visible binary path no longer names the repository checkout, and no symlink resolves it back there.
  Runner-equivalent standalone execution produced a regular executable under `/private/tmp/.../runner/spacedock-live-bin`, outside `/private/tmp/.../checkout`, with no symlink target.
- DONE: AC-3 — The change does not weaken the separate wrong-root detector or claim to solve launch-cwd authority.
  The diff touches no `internal/ensigncycle` or first-officer files, and all wrong-root detector cases passed independently.
- DONE: Verify AC-1 through AC-3 with exact evidence: physical non-symlink binary outside the canonical checkout, embedded clean revision equals the reviewed head, local plugin bytes remain current, and wrong-root/first-officer behavior is untouched.
  Exact-SHA standalone provenance recorded checkout and embedded revision `5ff92f9b58e2aae918a9cbbda560eb7394dae413` with `vcs.modified=false`; local plugin staging code is unchanged, but its required live-path proof remains the AC-1 failure above.
- FAILED: Run or inspect the required exact-SHA Claude live E2E evidence and retained provenance/executable artifacts; return PASSED or REJECTED and name any infrastructure blocker without treating a skipped or absent live run as green.
  REJECTED: commit `5ff92f9b` is not on origin and has no workflow run or uploaded artifacts; external push/live dispatch was withheld pending captain approval.

### Summary

The implementation passes all local behavioral, full, race, bounded-diff, exact-revision, placement, and wrong-root checks in a fresh standalone checkout matching the CI topology. Validation recommends REJECTED only because the acceptance contract requires a final-SHA Claude live run and retained CI artifacts, and none currently exists; rerun validation after that evidence is available.

## Stage Report: validation (cycle 2)

- DONE: Run corrected-guideline Roborev first on exact range 557f8df3..5ff92f9b and inspect its stored range/prompt; on any finding stop without push or CI.
  Roborev job `895` returned `No issues found`; stored `git_ref` is exact `557f8df3e6a62d34987edda70533375fc48ba8f6..5ff92f9b58e2aae918a9cbbda560eb7394dae413`, and its prompt contains all four configured guideline sections.
- FAILED: Only after Roborev PASS push exact head 5ff92f9b, verify the remote branch, and dispatch the exact-SHA Claude live E2E required by AC-1 with retained provenance and executable artifacts.
  GitHub rejected the exact-head push because the HTTPS OAuth credential lacks workflow scope; the remote branch remains absent and no CI was dispatched.
- SKIPPED: Inspect the actual live result and artifacts, rerun local focused/full/race evidence if the head changes, and return PASSED or REJECTED without treating waiting/skipped/absent CI as green.
  REJECTED: no remote head or live run exists after the authorized push failed; local head stayed exact `5ff92f9b`, so the already-passing focused/full/race evidence did not require rerun.

### Summary

The corrected Roborev exact-range gate passed cleanly and verified the intended guideline injection. Validation remains REJECTED on an explicit infrastructure blocker: Git cannot publish the workflow-changing commit without an OAuth credential carrying workflow scope, so the required exact-SHA live evidence cannot yet be produced.

## Stage Report: validation (cycle 3)

- DONE: Run corrected-guideline Roborev first on exact range 557f8df3..5ff92f9b and inspect its stored range/prompt; on any finding stop without push or CI.
  Roborev job `902` returned `No issues found`; its stored `git_ref` is exact `557f8df3e6a62d34987edda70533375fc48ba8f6..5ff92f9b58e2aae918a9cbbda560eb7394dae413`, and the prompt contains all four configured guideline sections.
- DONE: Only after Roborev PASS push exact head 5ff92f9b, verify the remote branch, and dispatch the exact-SHA Claude live E2E required by AC-1 with retained provenance and executable artifacts.
  The remote branch equals exact `5ff92f9b`; run `29303305720` is bound to that SHA, and retained Sonnet artifact `8299521188` plus Opus artifact `8299603800` contain provenance and installed executables through 2026-10-12.
- DONE: Inspect the actual live result and artifacts, rerun local focused/full/race evidence if the head changes, and return PASSED or REJECTED without treating waiting/skipped/absent CI as green.
  PASSED for this task: both Claude jobs completed `success`; head stayed exact and clean, so prior focused/full/race evidence did not require rerun. The overall workflow is honestly red because the independent Codex keep-moving scenario failed its durable dispatch assertion.
- DONE: AC-1 (VALUE) — Claude live E2E exercises a physical candidate binary outside the checkout while continuing to run the exact workflow-dispatch SHA and current local plugin bytes.
  Sonnet and Opus passed the live gate, ensign-cycle, shared-scenario, pty, and merged-team lanes; both uploaded binaries embed exact `5ff92f9b`, report `vcs.modified=false`, and ran with the checkout's unchanged staged plugin path.
- DONE: AC-2 — The host-visible binary path no longer names the repository checkout, and no symlink resolves it back there.
  Both artifacts record `/home/runner/work/_temp/spacedock-live-bin/spacedock`, outside `/home/runner/work/spacedock/spacedock`; downloaded binaries are regular executable ELF files, not symlinks, with identical SHA-256 `bea2fbfec20213a1dc677add0ddb2485a9d0c4c3c620618c2611e1faecb15745`.
- DONE: AC-3 — The change does not weaken the separate wrong-root detector or claim to solve launch-cwd authority.
  The exact two-file diff remains limited to the Claude workflow install/upload path and its behavioral release test; no `internal/ensigncycle` or first-officer contract file changed.

### Summary

Exact-SHA live evidence now closes the prior external-proof gap: both Claude variants exercised the neutral installed candidate successfully and retained independently inspectable provenance and executable artifacts. Validation recommends PASSED for `2h`; run `29303305720` is not a green release-certification run because an unrelated Codex live assertion failed, and that failure remains visible for separate triage.
