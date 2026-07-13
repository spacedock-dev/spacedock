# Proposal: Help Spacedock users adopt Roborev validation

## Decision

Ship a user-invocable `spacedock-roborev-setup` skill that helps a
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

The user invokes `$spacedock-roborev-setup` from the project root and asks
to add Roborev to the current Spacedock workflow.

The skill then:

1. Discovers the workflow README, stage names, base branch, worktree policy,
   split state checkout, and current Roborev state.
2. Explains the integration boundary and proposed workflow edits before making
   changes.
3. Directs setup through Roborev's current official documentation:
   - [Quick Start](https://roborev.io/quickstart/)
   - [Configuration](https://roborev.io/configuration/)
   - [Subagent Review Panels](https://roborev.io/advanced/subagent-review-panels/)
4. Points the user to those documents to choose available review agents and
   configure an advisory `quick` panel plus a required `code_completion` panel in
   `.roborev.toml`. The Roborev documentation remains the command and schema
   authority.
5. Runs or proposes `roborev init` from the canonical code repository outside
   Safehouse, subject to the user's approval for installation and configuration
   changes.
6. Adds the split state checkout's actual state branch to that checkout's
   `excluded_branches` during setup. This is configuration, not a live enqueue
   probe. The integration never deliberately queues a state-only review merely
   to prove the exclusion.
7. Updates the workflow's implementation and validation stage definitions:
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

## Why the required panel belongs at implementation exit

Running `code_completion` before implementation completes keeps routine
review-fix-review iterations with the implementation worker. Roborev's panel
members remain independent reviewers; the maker only receives and addresses
their synthesized findings.

This placement avoids dispatching a fresh behavioral validator merely to relay
code-review findings back to implementation. It also keeps ordinary Roborev
iterations out of Spacedock's validation feedback-cycle count.

The implementation worker records one disposition for each synthesized
finding:

- `fix`: change the code and strengthen proof;
- `rebut`: cite concrete repository evidence and ask a replacement independent
  panel to decide whether the finding clears;
- `needs decision`: stop for a real product, contract, or policy fork.

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
- disposition of every synthesized finding.

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
2. Configure and initialize Roborev from its official setup documentation.
3. Add the state-branch exclusion during setup.
4. Update the implementation and validation stage definitions.
5. Pilot on the next entity leaving implementation.
6. Prove a failing panel keeps the entity in implementation, a fixing commit
   invalidates the old result, and a replacement panel can pass.
7. Confirm the fresh validator verifies the stored exact-head evidence and
   independently reproduces the entity's ACs.

The pilot succeeds when Safehouse commits receive advisory reviews when the
daemon is healthy, the required branch panel produces one persisted synthesis
parent, failures stay inside implementation, and the entity retains a durable,
self-sufficient evidence trail.
