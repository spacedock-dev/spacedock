---
title: Install the checked-out candidate into a neutral bin for Claude live E2E
status: ideation
source: "Captain direction 2026-07-13 after Opus treated the checkout-shaped SPACEDOCK_BIN path as a workflow root."
score: 0.85
started: 2026-07-13T04:20:06Z
completed:
verdict:
worktree:
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
