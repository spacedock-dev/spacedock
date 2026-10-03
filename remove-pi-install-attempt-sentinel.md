---
title: The Pi install offer is suppressed in every session after one failed attempt
status: ideation
score: 0.8
source: "FO review of task ekw, 2026-10-03: the install-attempt sentinel has no session key, and its only measured effect is to suppress the offer."
sprint: pi-ux
sprint-readiness: ready
id: 271f46crset81jwf2rgast0c
started: 2026-10-03T17:32:28Z
gates:
    version: 1
    records:
        - id: gate:271f46crset81jwf2rgast0c:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:271f46crset81jwf2rgast0c-ideation-1
              briefing:
                id: briefing:271f46crset81jwf2rgast0c:ideation:attempt-1:revision-1
                digest: sha256:b5468d7dede594b3b955763bdff5f636257252272941a0c80b931f2ae33d2b1f
                room-ref: '@review/ideation/briefing-1'
---

## Problem

The First Officer install reference uses `${TMPDIR:-/tmp}/spacedock-install-attempted`
as a machine-wide failure latch. It has no host, session or project key. It asks
the First Officer to skip the offer if the file exists and to create it before
installation, so a failed attempt is intended to suppress later offers everywhere.
The user must discover and remove a temp file to recover the automatic path.

This is not a Pi-only feature. `skills/first-officer/SKILL.md:29-38` loads the
shared core for claude, codex and Pi; `first-officer-shared-core.md:10,47` loads
`fo-install.md`, the one install reference for all three hosts. Step 5 already
says "Fall back, never loop" and prohibits a second install attempt. Removing the
attempt file needs neither a replacement marker nor session identity. Task ekw
is superseded; do not restore its session-key design.

## Visible value and evidence limit

The intended value is that an earlier failed install no longer suppresses a
later session's offer in any host or project. The checkable change is removal of
the instruction's on-disk attempt gate and write: one existence gate and one
write directive become zero. This removes the documented machine-wide latch.

No check on HEAD exercises the install-offer flow. There is no deterministic
baseline observing a failed install followed by a new session's offer, and no
test that this removal can turn from red to green. A new session actually getting
the offer after failure, or a host actually stopping after one attempt, cannot
be observed without a live First Officer run. That is a residual behavioural
limit, not proven evidence; no live lane is authorized for this task. Source
inspection proves the instruction change only.

## Current-tree inventory and corrected baseline

Audited code baseline: `4436ec14c`. The original five/six-file surface was based
on stale worktrees, not this HEAD. The FO confirmed that error during ideation
and approved the one-file correction and narrower, honest proof claims.

Every current sentinel-bearing instruction is in
`skills/first-officer/references/fo-install.md`:

- Line 28: the only literal path reference (twice on that line), existence
  reader, skip-offer branch, manual `rm` advice, and fallback-message reference.
- Line 29: "Sentinel absent" precondition on OFFER.
- Line 30: `touch` before install and lifetime attempt-count instruction.
- Line 32: indirect fallback to the sentinel-dependent step-1 message.
- Line 31 has no sentinel dependency; only its step number changes.

The four other entries originally cited are obsolete for this baseline:

| Original location | Current-tree finding and disposition |
| --- | --- |
| `internal/contractlint/version_gate_smoke_test.go:79` | No sentinel guard on HEAD; line 79 closes the registry loop. The obsolete guard belongs to stale worktrees. Preserve this file. |
| `skills/integration/testdata/version_gate_flow.sh:17` | File absent from HEAD; stale-worktree fixture, not shipped. Do not restore it. |
| `skills/integration/version_gate_fixture_test.go:126` | File absent from HEAD; stale-worktree test, not a current proof owner. Do not restore it. |
| `docs/site/get-started/install.md` | No sentinel prose or "re-enables the offer" sentence on HEAD; that version was in stale worktrees. No current line to delete. Preserve the document and manual commands. |

`skills/integration/install_hint_channel_test.go` neither names the sentinel nor
requires persistent suppression. No test on HEAD asserts that the file exists or
that a failed install stays suppressed; there is no such test to remove or
convert. Do not bring back the obsolete tests to delete them again.

Search scope was the tracked current code tree (`git grep`), supplemented by
searching install/attempt/sentinel and flow wording in skills and contractlint.
This excludes nested old worktrees and the separate workflow-state checkout.
The broad legacy-name search also found
`internal/ensigncycle/testdata/claude_live_auto_continue_run31915540750_sonnet.stream.jsonl:16,60`:
it records `fo-install-gate.md`, not the sentinel path. This is archived trace
evidence, not another shipped install reference; leave it unchanged. No other
current tracked file names the sentinel path. The design body intentionally
names it for audit; workflow state is not shipped product text.

## Proposed approach and exact instruction change

Delete the step-1 existence check, the persistent one-attempt rule, the
sentinel-exists explanation and `rm` advice. Remove "Sentinel absent" from the
offer heading and the `touch` instruction from approval. Move the useful manual
hint and stale-`SPACEDOCK_BIN` advice into the terminal fallback so deleting step
1 does not delete recovery guidance or leave a dangling step reference.

Replace `fo-install.md:28-32` with the following four steps. This is the concrete
before/after design: the before is the five numbered lines inventoried above;
the after is exact text below. Keep lines 1-27 byte-identical.

```markdown
1. **OFFER:** "Run the install for you and resume startup once the binary lands?" On decline: ABORT with the manual command; once `spacedock` is on PATH, start the first officer with the host's spacedock launcher command.
2. **On approval:** run the command exactly once.
3. **Converge (session-scoped repoint).** `install.sh` prints `install.sh: installed spacedock <version> to <dir>/spacedock` to stderr and warns when the dir is off PATH. Parse the installed path from that stderr line; if absent, probe `$HOME/.local/bin/spacedock`. If a path resolves, set `SPACEDOCK_BIN` to it **for this session only — never persist it to a shell profile** — and re-check `${SPACEDOCK_BIN:-spacedock} --version`. If line 1 parses to a compatible version, that repointing IS the gate's one launcher resolution (blessed by the shared-core invariant): resume Startup. Next session: if the install dir is on PATH, bare `spacedock` resolves with no override; if not, the gate fails again and the fallback hint names the exact installed path and tells the human to add the dir to PATH (or launch with `SPACEDOCK_BIN=<path>`).
4. **Fall back, never loop.** If the re-check still fails, or no path resolves, print the channel-correct command to run manually and ABORT. If `SPACEDOCK_BIN` was set-but-stale, name its session-scoped replacement. **No second install attempt, no proceeding without the re-check**.
```

Do not delete old temp files on the user's machine; after removal the instruction
neither consults nor writes them. No session marker, session identity, retry
mechanism, migration, new helper or new fixture is introduced. The simplest
alternative is keeping the existing terminal fallback, which already bounds the
flow; a replacement session marker adds state without serving the value.

## Documentation disposition

The instruction replacement above is the entire documentation change.
`docs/site/get-started/install.md` already contains neither the sentinel nor its
removal advice. Its concrete before/after is identical: keep the Homebrew commands
at lines 9-12 and Linux stable/edge commands at lines 25 and 36 untouched. Do not
add speculative site prose or fabricate a deletion against a stale version.

## Expected surface and tolerance

- Implementation: only `skills/first-officer/references/fo-install.md`, lines
  28-32. No test currently names the sentinel, so no check must move.
- Estimate net LOC change: -1, across 1 file (4 insertions, 5 deletions).
- Tolerance: net -5 to 0, at most 1 implementation file. Any additional file or
  mechanism needs a revised gate. The state-checkout design/report is the
  ideation artifact, not part of the product implementation count.
- Semantic change: offer eligibility no longer consults stored attempt state.
  Within the existing supported-OS, outside-sandbox startup path, the offer
  depends on binary absence, not earlier attempts. User approval still controls
  execution; wrong-version, sandbox and unsupported-OS handling are unchanged.
- Unchanged: CLI grammar, install commands, classifier, persisted formats,
  mutation authority, convergence, version re-check and terminal no-retry rule.
  Do not alter the session-scoped launcher repoint; that is not an attempt marker.

## Acceptance criteria

**AC-1 (VALUE, instruction-level) — Shipped instructions have no persistent
attempt gate or write.** The current baseline has one path-naming line, one
existence gate and one write directive; after the change there are zero. No
shipped file names or consults `spacedock-install-attempted`, so the instructed
offer decision no longer depends on earlier on-disk attempt state. Proof:
`git grep -n -I -F spacedock-install-attempted` in the product checkout returns
no matches (exit 1), plus inspect the four-step diff for any replacement marker.
Falsifier: restore the old check or write. This measures text, not host execution.

**AC-2 — The required install content is unchanged.** The channel classifier,
per-OS commands, sandbox arm, offer question/decline behaviour, convergence and
version re-check remain intact; the two casks are still `spacedock` and
`spacedock@next`. Proof: compare lines 1-27 and the preserved offer/convergence
text against `4436ec14c`, and run the focused existing checks below. Falsifiers:
change a preserved command, classifier arm, sandbox instruction or convergence
step. Test coverage gaps are stated below rather than attributed to tests.

**AC-3 — The instruction still terminates after one approved attempt.** Approval
still says "exactly once"; failure prints the channel-correct manual command,
retains stale-launcher advice and aborts with no second attempt or unchecked
resume. Proof: independent review of the exact replacement against the old
steps 1/3/5. Falsifier: add retry/proceed language or lose the manual recovery
message. No automated or live behavioural observation is claimed.

## Test plan and proof owners

No spike needed: this removes a prose-only mechanism, uses existing deferred
loading and keeps the existing terminal fallback; it introduces no parser,
runtime identity or unverified handoff. There is no current primary behavioural
proof owner for the offer flow. Do not invent a model fixture that proves its own
copy of these instructions, and do not claim a prose search executes the flow.

The existing no-regression owners are:

- `internal/contractlint/version_gate_smoke_test.go:33`
  `TestVersionGateDeferredTrigger`: shared-core deferral and nontrivial reference;
  deleting the deferred-load entry fails it.
- Same file, line 57, `TestVersionGateSandboxRegistry`: shared registry and
  human outside-sandbox message; losing the latter fails it. This is a structural
  guard, not execution of sandbox decisions.
- `internal/contractlint/install_hint_drift_test.go:56`
  `TestInstallHintNoDrift`: Linux stable and Homebrew tap/stable tokens match
  published docs; changing one side fails it.
- `skills/integration/install_hint_channel_test.go:95`
  `TestChannelClassifierTable`: executes the shipped classifier; removing an
  edge-classifying arm fails the corresponding independent path example.
- Same file, line 153, `TestStableHintMatchesPublishedDoc`: stable Linux command
  matches the published form; drifting the reference command fails extraction
  or equality.
- Same file, line 176, `TestEdgeHintDeliversChannelToScript`: executes the **doc's**
  edge pipeline with a local fetch substitute; binding the channel to curl
  instead of sh fails. It does not independently execute the reference's edge
  command; byte-preservation review covers that command.
- Same file, line 286, `TestContractCasks` (`names`, `edge-satisfies-pin`): extracts
  both hinted casks and, when Homebrew/tap is available, checks token existence
  and edge minor compatibility. Removing a cask fails the count; an unknown
  cask or too-old edge version fails the tool-gated subtests. `-short` skips tap
  resolution, so it is not evidence of cask availability. Keep this test unchanged.

Implementation validation is bounded to text review/search and these existing
checks (seconds, no new scaffolding):

```sh
go test ./internal/contractlint -run '^(TestVersionGateDeferredTrigger|TestVersionGateSandboxRegistry|TestInstallHintNoDrift)$' -count=1
go test ./skills/integration -short -run '^(TestChannelClassifierTable|TestStableHintMatchesPublishedDoc|TestEdgeHintDeliversChannelToScript|TestContractCasks)$' -count=1 -v
```

Require all executed checks to remain green; report the cask subtests as skipped,
not passed. The same commands were run during ideation: three contractlint and
three integration tests passed; cask resolution skipped under `-short`. No new
or updated tests. No live lane, repository-wide suite, race suite or formatting
sweep is warranted for this design-only round. Independent review must accept
the explicit behavioural limit; this plan cannot prove cross-session offers or
cross-host execution without separately authorized live evidence.

## Stage Report: ideation

- DONE: Name every place that carries the sentinel today, with the exact path and line, and report the full list rather than the five the body already names. Include skills/first-officer/references/fo-install.md, internal/contractlint/version_gate_smoke_test.go, skills/integration/testdata/version_gate_flow.sh, skills/integration/version_gate_fixture_test.go and docs/site/get-started/install.md, and search for others.
  Current-tree inventory records fo-install.md:28-30,32, four obsolete stale-worktree entries, and the two archived legacy-name trace lines; FO approved the corrected HEAD baseline.
- DONE: Find out whether any test requires this rule: a failed install stays suppressed. Name each such test and say what it must become. If a test asserts the sentinel exists, it protects the fault and must be removed with it.
  None on HEAD; both originally named flow fixtures are absent, and the remaining smoke tests have no sentinel assertion. No tests need removal or replacement.
- DONE: Confirm that the install commands, the channel classifier and the two brew casks in fo-install.md are untouched by this change, and name the check that pins them.
  Exact replacement preserves lines 1-27; test-plan inventory names classifier, command-drift, edge-pipeline and cask checks, including their coverage limits and falsifiers.
- DONE: Record the design in the body: the exact deletions, the semantic change, the per-criterion proof plan, and the statement that the changed file is shared by claude, codex and Pi. Keep the surface within the body's tolerance, or say why it must grow.
  FO-approved corrected surface is one instruction file, +4/-5 lines (net -1); body supplies exact replacement and instruction-level ACs with no runtime-proof claim.
- DONE: Run bounded existing no-regression checks as an ideation baseline.
  Three contractlint and three integration checks passed; their independent falsifiers are named in the test plan, not treated as proof of the offer flow.
- SKIPPED: Live install-offer observation and repository-wide suite.
  Explicitly prohibited; cross-session offer/no-retry behaviour remains unobserved. TestContractCasks tap resolution also skipped under -short.

### Summary

Designed removal of the machine-wide attempt file for claude, codex and Pi without a replacement marker, preserving the one-shot fallback and all install commands. Corrected stale-worktree assumptions against HEAD 4436ec14c with FO approval; no current test exercises the offer flow, so the design explicitly limits proof to shipped instructions and existing no-regression checks. Only this state-checkout body/report changed in ideation; product implementation awaits the review gate.
