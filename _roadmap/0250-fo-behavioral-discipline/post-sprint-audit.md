# 0250 FO behavioral discipline — independent pre-cut audit

**Audited commit:** `aeeadf451c80748215de344ca8b814eca4e81afa` (`v0.25.0-pre2`)

## Verdict: BLOCKED

Do not cut `v0.25.0` yet. Two independently reproduced Codex front-door
failures prevent `spacedock codex` from reliably starting the first officer.
They sit in the release-critical launcher path, so neither can wait for a
follow-up sprint.

This was an adversarial review of `v0.24.0..aeeadf45`: the assembled FO
contract and host adapters, launcher/bootstrap behavior, Runtime Live E2E
evidence, release gates, leanness cap, and the `7v` coordination. I used source
traces, the real installed Codex parser, focused Go checks, and GitHub Actions
state; I did not treat contract wording or Contractlint presence as behavioral
proof.

## Ship blockers

### B1 — a forwarded Codex positional creates an invalid two-prompt invocation

`runCodex` appends the entire forwarded argv and then appends
`codexBootstrapPrompt` whenever it does not classify the launch as a resume
([`internal/cli/frontdoor.go` at `aeeadf45`, lines 606–608](../../../internal/cli/frontdoor.go)).
Codex accepts one optional positional prompt, not two. A caller who forwards a
valid positional prompt or subcommand therefore receives a parser error instead
of a first-officer session.

This is a 0.25 release-delta regression: `d0082e3d` replaced the preceding
nonempty-post-fence opaque posture with bootstrap for every non-`resume` argv.

Reproduction against the real Codex CLI at the audited commit:

```console
$ go run ./cmd/spacedock codex --skip-compat-check -- --model gpt-5.6-sol opaque-argument
spacedock 0.25.0-pre2+dev · launching codex as your first officer
Workflow: docs/dev
Sandbox: unavailable (safehouse not on PATH)
codex is your first officer — ask it for the queue and next steps.
error: unexpected argument 'You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Assume $spacedock:first-officer for the entire session.' found

Usage: codex [OPTIONS] [PROMPT]
       codex [OPTIONS] <COMMAND> [ARGS]

For more information, try '--help'.
exit status 2
```

The user-facing reference promises that material after `--` forwards verbatim
to the host ([`docs/site/reference/command-reference.md`, line 27](../../site/reference/command-reference.md)). Existing tests instead accept the invalid
constructed argv: `TestCodexPostFenceWithoutResumeBootstrapsFirstOfficer` and
`TestCodexFrontDoorBootstrapsPostFenceWithoutResumeThroughSafehouse` use a fake
launcher, so neither exercises Codex's single-prompt grammar.

**Smallest corrective action:** preserve bootstrap for an options-only fresh
launch, but do not append a second prompt after a forwarded positional or
subcommand. Restore a command-position-aware classifier (or reject the ambiguous
shape with an actionable diagnostic), and add regression coverage that invokes
the real Codex argument parser or an equivalent grammar-enforcing shim. Cover a
model-only launch, a forwarded positional, `exec`, and `resume`.

### B2 — `--profile resume` is mistaken for the `resume` subcommand

The new resume detector is `slices.Contains(fd.passthrough, "resume")`
([`internal/cli/frontdoor.go` at `aeeadf45`, line 578](../../../internal/cli/frontdoor.go)).
It treats every exact token named `resume` as the subcommand. Codex documents
`--profile <CONFIG_PROFILE_V2>`, so a valid profile named `resume` is a fresh
launch, not a resumed session.

I placed a harmless fake `codex` executable first on `PATH`; it printed the
argv and exited. This direct front-door command produced no banner, default
approval flag, or first-officer bootstrap prompt:

```console
$ PATH="$audit_bin:$PATH" go run ./cmd/spacedock codex --skip-compat-check -- --profile resume
argv: <--profile> <resume>
```

The expected fresh-launch argv includes both
`--ask-for-approval on-request` and the `$spacedock:first-officer` bootstrap
prompt. An isolated package-test overlay independently failed with the same
observed argv:

```text
launch argv = [codex --profile resume], want
[codex --ask-for-approval on-request --profile resume <bootstrap prompt>]
```

**Smallest corrective action:** determine whether `resume` occupies Codex's
subcommand position, rather than searching all forwarded tokens. The same
grammar-aware regression suite must prove that a profile value named `resume`
boots a fresh FO while `codex --model <model> resume <id>` remains prompt-free.

## Evidence that held

- The exact-SHA Runtime Live E2E predicate passed for the audited commit:

  ```console
  $ go run ./cmd/spacedock-release e2e-gate aeeadf451c80748215de344ca8b814eca4e81afa
  green Runtime Live E2E run 29221296098 matches release commit aeeadf451c80748215de344ca8b814eca4e81afa
  ```

  [Run 29221296098](https://github.com/spacedock-dev/spacedock/actions/runs/29221296098)
  completed successfully for that SHA. Its offline, Claude (Sonnet and Opus),
  Codex, and Pi jobs passed. This supports the assembled contract and runtime
  adapters, but it does not cover these blockers: the Codex live runner invokes
  `codex exec --json` directly
  ([`internal/ensigncycle/codex_live_runner_test.go`, lines 18–23 and 402](../../../internal/ensigncycle/codex_live_runner_test.go)),
  rather than the `spacedock codex` front door.

- The release mechanism fails closed on exact SHA. The focused release checks
  passed:

  ```console
  $ go test ./internal/release ./cmd/spacedock-release -run \
    '^(TestE2EGatePassesForGreenRunOnReleaseCommit|TestE2EGateBlocksGreenRunOnWrongCommit|TestReleaseWorkflowGatesGoreleaserOnE2E|TestReleaseWorkflowGatesGoreleaserOnManifestTag|TestReleaseWorkflowEdgeAdvanceDecisionGates|TestReleaseWorkflowAlwaysCutPre0|TestEdgeAdvancePatchDoesNotRegressNext|TestManifestTagGatePassesWhenStamped|TestE2EGateCommandPassesOnGreenRun|TestManifestTagGateCommandPassesOnMatch|TestEdgeAdvanceDecisionCommandPrintsVerdict)$' -count=1
  ok   github.com/spacedock-dev/spacedock/internal/release
  ok   github.com/spacedock-dev/spacedock/cmd/spacedock-release
  ```

  The workflow gates Goreleaser on both the matching Runtime Live E2E result and
  manifest/prose agreement. The current `0.25.0-pre2` manifests are deliberately
  not eligible for a `v0.25.0` tag; the normal future stamp → exact-SHA E2E →
  annotated-tag sequence remains required. That is an expected release step,
  not a release-gate defect.

- The 0250 resident-contract cap holds. Against the documented `v0.24.0`
  baseline, `first-officer-shared-core.md` plus `present-gate/SKILL.md` grew
  from 27,000 to 30,215 bytes: **+3,215 bytes**, below the 5,600-byte ceiling.

- The `7v` sequencing alternative holds: `spacedock status --workflow-dir
  docs/dev --where id=7v` reports `pi-bootstrap-prompt-parity` as `backlog`.
  No started stale-target task can reintroduce the removed flourish.

## Non-blocking follow-ups

1. **Single-root FO state paths need an explicit classifier rule.**
   [`fo-write-core.md`, line 10](../../../skills/first-officer/references/fo-write-core.md)
   permits `allowed-state` only beneath `.spacedock-state/**` or
   `{workflow_dir}/_archive/**`, while the same document authorizes entity
   frontmatter and `new` writes. The binary supports direct entity files in a
   single-root workflow:

   ```console
   $ go run ./cmd/spacedock status --workflow-dir internal/status/testdata/seq-workflow --boot --identify --json
   ... "state_backend":"single-root", ...
   ... "entity_dir":".../internal/status/testdata/seq-workflow", ...
   ```

   Define flat `{workflow_dir}/<entity>.md` and
   `{workflow_dir}/<entity>/index.md` paths as allowed state, or define a shared
   state-root abstraction. The current binary behavior works, so this is a
   contract clarification, not a cut blocker.

2. **Tighten release-runbook wording.** `docs/releasing.md` calls Runtime Live
   E2E `workflow_dispatch`-only, but `.github/workflows/runtime-live-e2e.yml`
   also has a `pull_request` trigger. The safety claim remains true: fresh main
   commits do not run it automatically. State that narrower fact. Also have the
   runbook select or verify the dispatched run by `REL_SHA`, rather than merely
   the latest main-branch run; the CI gate already fails closed, so this is an
   operator-friction improvement rather than a release-safety defect.

## Required re-audit boundary

After both blockers are fixed, run the direct Codex-parser probes above, the
focused front-door tests, and a new Runtime Live E2E on the final stamped
`v0.25.0` commit. Confirm its `headSha` equals the stamped commit before creating
the annotated tag. This report does not authorize a tag on the current pre2
commit.
