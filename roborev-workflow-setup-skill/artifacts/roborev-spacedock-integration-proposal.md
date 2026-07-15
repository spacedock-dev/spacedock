# Proposal: Help Spacedock users adopt Roborev validation

## Decision

Ship a user-invocable `roborev-setup` skill, exposed by the plugin as
`spacedock:roborev-setup`, that helps a
Spacedock user add Roborev to an existing development workflow when independent
code review would strengthen the workflow.

Roborev supplies code-review evidence. Spacedock retains authority over stage
state, implementation work, behavioral validation, rejection routing, and
human gates.

The integration uses two review layers:

1. Roborev's Git post-commit hook queues advisory, asynchronous reviews during
   implementation.
2. The implementation worker runs a required, multi-reviewer `code_completion`
   panel against the complete implementation branch before reporting the stage
   complete.

The fresh Spacedock validator verifies that exact-head Roborev evidence, then
reproduces the entity's acceptance criteria. Roborev cleanliness never replaces
behavioral validation.

## Why setup starts with `roborev quickstart`

The first draft told the setup agent to follow three official web pages before
it named any local setup command. In the adoption run this caused an avoidable
search detour: the agent went looking for setup guidance even though
`roborev quickstart` already prints the repository's **Current state** and a
version-matched **Configuration playbook**.

The revised skill runs `roborev quickstart` as its first Roborev command and
works from those two sections. Local `roborev <command> --help` is the next
source for syntax. A directly named official page is a targeted fallback only
when those local sources leave a specific configuration gap; broad search is
not part of normal setup.

The pilot also produced an over-broad `review_guidelines` block by copying
repository procedures and test commands. Roborev's own repository uses this
field to calibrate judgments the reviewer cannot infer on its own: trust
boundaries, intentional compatibility posture, false-positive suppressions,
and review-focus boundaries. The setup skill follows that narrower pattern,
never duplicates `AGENTS.md`, workflow or stage instructions, component
procedures, or required developer commands, and shows the proposed calibration
to the user before writing it.

The setup also combined strict TDD with an automatic post-commit review hook
without saying whether the expected-red state belonged in Git. That ambiguity
caused implementers to commit intentionally failing or non-buildable tests,
which immediately consumed reviewer time and tokens and polluted pass-rate and
correction data with a known failure. Excluding those commits from convergence
accounting avoids false escalation but not wasted review work. The generated
contract now keeps RED in the working tree, records the pre-fix command, exact
failure, and predicted reason in the Stage Report, and permits a product commit
only after the test and minimal implementation pass the focused check and
relevant suite together. The post-commit hook remains enabled for that green
candidate.

The adoption pilot also produced a valuable semantic adversarial checklist for
the maker. That guidance belongs in `AGENTS.md` or generated implementation-stage
instructions, not in `review_guidelines`: trace changed values and events across
representations and lifecycle phases; matrix adjacent variants and boundaries;
validate full records atomically; inspect scaling and implicit limits; and ask
how tests could pass while observable behavior is wrong.

## When the skill should recommend adoption

Recommend this integration when the workflow:

- develops code in Git branches or worktrees;
- has a distinct implementation stage followed by fresh validation;
- routes validation defects back to implementation;
- benefits from persistent review history or multiple reviewer perspectives;
- runs trusted code, or places the whole Roborev daemon inside an appropriate
  VM or container.

Do not force it onto prose-only workflows, workflows without a meaningful code
diff, or projects whose trust boundary cannot tolerate an external local review
daemon and its agent processes.

## User adoption journey

The user invokes `$spacedock:roborev-setup` from the project root and asks
to add Roborev to the current Spacedock workflow.

The skill then:

1. Discovers the workflow README, stage names, base branch, worktree policy,
   and split state checkout.
2. Runs `roborev quickstart` from the canonical code repository as the first
   Roborev command. Its Current state and Configuration playbook are the setup
   authority.
3. Explains the integration boundary, the missing items that matter, the
   intentionally omitted fix/refine and agent-hook features, and the proposed
   workflow edits before making changes. Any proposed `review_guidelines`
   calibration is reviewed separately by the user before it is written.
4. Uses local `--help` for command syntax and opens only the directly relevant
   official configuration or panel page when quickstart leaves a named gap.
5. With approval, applies only the missing setup needed for an advisory `quick`
   panel and a required `code_completion` panel, then reruns `quickstart` to
   verify the result.
6. Adds the split state checkout's actual state branch to that checkout's
   `excluded_branches` during setup. This is configuration, not a live enqueue
   probe. The integration never deliberately queues a state-only review merely
   to prove the exclusion.
7. Updates the workflow's implementation and validation stage definitions:
   - expected RED is an uncommitted working-tree state evidenced by the pre-fix
     command, exact failure, and predicted reason in the Stage Report;
   - red-only and non-buildable product commits are forbidden, and the test and
     minimal implementation are committed together only after focused and
     relevant-suite checks pass while the post-commit hook stays enabled;
   - implementation must produce passing exact-head `code_completion` evidence;
   - validation must verify that evidence without rerunning an unchanged panel;
   - a fixing commit invalidates the old evidence and returns responsibility to
     implementation.
8. Checks daemon reachability and records any remaining operator-only setup.
9. Detects whether Spacedock workers run inside Safehouse or another sandbox.
   When they do, the skill explains the boundary and advises the user about the
   minimum runtime access needed. It omits sandbox-specific advice when workers
   run directly on the host.

No first-officer core-contract change and no new workflow schema are required.
The integration is a workflow-owned implementation-exit discipline.

## Exact-tip cost-gate order

Roborev's post-commit `quick` job and an explicitly requested
`code_completion` panel are independent. The generated implementation-stage
contract therefore orders them:

1. Produce test-first RED in the working tree, record its command, exact
   failure, and predicted reason in the Stage Report, and leave `HEAD`
   unchanged.
2. Add the minimal implementation, pass the focused check and relevant suite,
   then commit the test and implementation together. Never create a red-only or
   non-buildable product commit; keep the post-commit hook enabled.
3. Record the final green candidate commit.
4. Wait for that exact commit's already-enqueued quick review with `roborev wait
   <head>`, then inspect `roborev show <head> --json` and verify the quick panel.
5. Keep Medium-or-higher findings in implementation. Fix, rerun the relevant
   checks, and recommit only green, then
   wait for the replacement tip. Low findings remain advisory.
6. Only after quick clears, launch `code_completion` for the complete branch and
   preserve its synthesis parent.

The two jobs never start in parallel. Missing or failed final-tip quick
coverage is repaired before the expensive panel starts. Any fixing commit
invalidates both results and restarts the sequence at quick. Quick saves cost;
it never satisfies the implementation-exit evidence requirement.

The `wait` exit code alone is not the Low-versus-Medium decision: a Low-only
review may still have a failed verdict. The implementation contract classifies
the stored exact-tip findings and permits only the Low-only case to continue.

## Implementation-owned semantic pass

Before the final candidate commit and the quick wait, the generated
implementation guidance requires a semantic adversarial pass:

- trace identity, cardinality, order, exact bytes, attribution, authority, and
  terminal state across every representation and lifecycle phase;
- matrix empty/terminal, repeated/out-of-order, every input path, Unicode, EOF,
  size, visibility, and layout variants;
- use canonical validators or atomic validation of the complete record;
- inspect hot paths and readers for multiplicative work, blocking I/O,
  unbounded allocation, and implicit size limits, with scaling/over-limit proof;
- assert exact observable results plus failure and cleanup behavior, especially
  where the old tests could pass despite wrong behavior.

Implementation owns this work and any fixes. It is pre-review preparation, not
Roborev reviewer calibration or Roborev-owned remediation.

## Why the required panel belongs at implementation exit

Running `code_completion` before implementation completes keeps routine
review-fix-review iterations with the implementation worker. Roborev's panel
members remain independent reviewers; the maker only receives and addresses
their synthesized findings.

This placement avoids dispatching a fresh behavioral validator merely to relay
code-review findings back to implementation. It also keeps ordinary Roborev
iterations out of Spacedock's validation feedback-cycle count.

Before assigning a disposition, the implementation worker records this
release-scope triage for each synthesized finding:

1. **Released user and normal workflow:** identify who encounters the finding
   in a released product and the normal workflow that triggers it.
2. **Observable harm:** state what that user loses or cannot complete.
3. **Protected value:** name the affected value AC or a non-negotiable safety,
   security, data-integrity, or compatibility boundary.
4. **Trigger evidence:** cite evidence that the trigger is common or explicitly
   promised.

All four fields are required before implementation assigns `fix`, `rebut`, or
`needs decision`. Without completed triage, a finding cannot be classified as
release-blocking. Completed triage supports classification; it does not by
itself make a finding release-blocking.

After triage, the implementation worker records one disposition for each
finding:

- `fix`: change the code and strengthen proof;
- `rebut`: cite concrete repository evidence and ask a replacement independent
  panel to decide whether the finding clears;
- `needs decision`: stop for a real product, contract, or policy fork.

The automatic review-fix-review loop allows at most two remediation rounds
after the initial required panel. A remediation round starts when implementation
acts on the current panel's findings and ends when the replacement
`code_completion` panel completes. Required quick-gate reruns are part of that
round, not additional rounds. If findings remain after the second remediation
round, implementation stops and sends the First Officer and captain every
remaining finding, its completed four-field triage, its disposition, and the
supporting evidence. No worker may make another code change or launch another
panel until the First Officer or captain directs the next action.

The implementation stage remains incomplete until the current head has passing
required evidence. Operational failure also blocks completion; absent output is
never a pass.

## Evidence contract

Roborev freezes a submitted branch range to commit SHAs. The implementation
worker fetches the synthesis parent with `roborev show --job <id> --json` and
records:

- Roborev version and synthesis parent ID;
- exact frozen `merge-base..head` range;
- current implementation head SHA;
- panel name and configured member outcomes;
- synthesis-parent verdict;
- completed four-field release-scope triage and disposition of every
  synthesized finding.

The synthesis parent is authoritative. Required member jobs must finish
successfully and appear exactly once, but an individual member verdict is not a
second gate: synthesis may reject a member's false positive. The required
result is a passing synthesis parent over the current implementation head.

Keep the compact evidence above in the implementation Stage Report. When the
full synthesis contains findings or is too long for a clear Stage Report, store
it as a folder-form entity artifact and link it from the report. Do not duplicate
the full text in both places.

Because Roborev's database is machine-local, the Stage Report and any linked
artifact must remain sufficient to understand what was reviewed and why the
stage advanced.

## Validation journey

The fresh validator:

1. Resolves the current immutable implementation head.
2. Fetches the recorded synthesis parent and verifies its frozen range, panel,
   required member execution, and passing parent verdict.
3. Confirms that no code commit has landed since that review.
4. Reproduces every acceptance criterion and the workflow's focused, full, and
   live checks.
5. Routes any newly discovered code defect back to implementation.

The validator does not rerun an unchanged passing panel. Any fixing commit
invalidates the old result; implementation must adjudicate the defect and
produce a passing replacement before validation resumes.

## Post-commit and daemon failure behavior

Post-commit reviews are advisory and fail open. If the daemon cannot be reached
or started, the commit succeeds and the quick review may be missing. A later
commit does not necessarily backfill that job.

The complete `code_completion` panel is the hard requirement. It covers the full
base-to-head range, including commits that missed quick review, and fails closed
on daemon, database, disk, required-member, or synthesis failure.

The skill should encourage a user-level daemon service and check `roborev
status` during setup or engagement, but it must not interpret missing advisory
coverage as a clean review.

## Agent-hook decision

Do not install Roborev's agent hook for a Spacedock-managed workflow. That hook
tracks session activity in a second local daemon and tells the active Codex or
Claude session to invoke `roborev-fix`. Its next-action authority and
code-writing fix loop compete with Spacedock's implementation ownership and
feedback routing.

The agent hook remains useful for projects using Roborev without Spacedock.
Reconsider it here only if a future repository-scoped notification can route
findings into Spacedock without invoking a competing fixer.

`roborev quickstart` may also report Roborev fix/refine skills as missing. Leave
them missing for this integration: implementation owns code changes and the
First Officer owns routing.

## Runtime and sandbox topology

When the skill detects Safehouse or another sandbox, it advises the user to run
one machine-local Roborev daemon outside that sandbox and install and
authenticate the reviewer agents there. Sandboxed Spacedock sessions act as
clients over loopback.

The skill advises exposing only what the sandboxed client needs, typically:

```text
add-dirs-ro=~/.roborev/runtime
```

Prefer a system or Homebrew binary. If the binary lives in a sandbox-hidden
home directory, grant its directory read-only access. Do not expose all of
`~/.roborev` or relocate `ROBOREV_DATA_DIR` into the project.

The external daemon and its child agents do not inherit Spacedock's Safehouse
boundary. Use this topology only for trusted repositories; sandbox the whole
daemon for untrusted code.

## Pilot in `spacedock-subspace`

After the current working-tree changes settle:

1. Invoke the skill from the canonical project root.
2. Confirm `roborev quickstart` is the first Roborev command and configure only
   the missing setup it identifies; use direct official pages only for a named
   panel-schema gap.
3. Add the state-branch exclusion during setup.
4. Update the implementation and validation stage definitions.
5. Inspect the generated implementation-stage contract for working-tree RED,
   durable pre-fix evidence, green-before-commit, combined test-and-implementation
   commits, the red-only/non-buildable prohibition, and an enabled post-commit
   hook. Drive it once to prove RED leaves `HEAD` unchanged and enqueues no job,
   while the combined green candidate does enqueue `quick`.
6. Prove the final-tip quick job finishes before `code_completion` starts;
   Medium-or-higher findings hold implementation while Low findings do not.
7. Prove every finding gets all four release-scope fields before disposition or
   release-blocking classification. Show that a fixing commit invalidates the
   old result, a replacement panel can pass, and findings left after two
   remediation rounds stop for First Officer/captain direction with their
   completed triage before another change or panel.
8. Confirm the fresh validator verifies the stored exact-head evidence and
   independently reproduces the entity's ACs.

The pilot succeeds when Safehouse commits receive advisory reviews when the
daemon is healthy, the required branch panel produces one persisted synthesis
parent, failures stay inside implementation, and the entity retains a durable,
self-sufficient evidence trail.
