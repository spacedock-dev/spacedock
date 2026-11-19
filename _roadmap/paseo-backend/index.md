# Sprint paseo-backend

**Goal.** Let a Claude first officer dispatch a fresh Codex ensign on the same Paseo daemon and verify the work from Git.

**Deliverable.** An opt-in Paseo dispatch binding and a live Claude-to-Codex acceptance test. The existing `spacedock dispatch build --host codex --stamp` supplies the assignment. The binding adopts the stamped checkout, launches the worker and follows it through verified completion. [Staff report](staff-report.md)

## Why this sprint exists

A Codex first officer drove a worker through this sequence by hand. The ensign committed the receipt on its task branch and the report in the canonical state checkout. A Claude-marked CLI subprocess produced the same Codex assignment bytes. The remaining integration risk is whether a live Claude FO selects the Paseo binding and receives the worker's completion signal. [Git proof](run-log.md#completion-and-durable-proof) · [Cross-host envelope comparison](run-log.md#cross-host-envelope)

## First version

Fresh Codex workers run on the FO's daemon. Spacedock owns the worktree and the assignment's write scope. Paseo owns the workspace and worker handles. The FO applies the existing checklist, report verification and one-retry rule. A worker can finish and the FO can steer it without a child daemon grant. Parent-directed mid-run questions are a separate capability with an ACL dependency. [Assignment](run-log.md#fixture-and-dispatch) · [Steering](run-log.md#steering-and-roster) · [Completion](run-log.md#completion-and-durable-proof) · [Retry contract](run-log.md#retry-and-budget-contracts) · [Child return channel](run-log.md#child-return-channel)

## Sequence

1. **Exercise the live Claude-to-Codex handoff first.** Write the acceptance fixture. Drive the existing build and Paseo calls from a real Claude FO. Check the Codex provider/model, worker cwd and unsolicited finish event. A host-marker substitution covers only the build. [Cross-host comparison](run-log.md#cross-host-envelope) · [Live Codex baseline](run-log.md#completion-and-durable-proof)
2. **Bind the proven calls at the FO dispatch point.** Add a skill smoke test for selecting Paseo, then add `references/paseo-fo-dispatch.md`. Use the existing Codex pointer, stamped entity and workspace adoption. Route the acceptance fixture through this binding. [Envelope](run-log.md#fixture-and-dispatch) · [Workspace adoption](run-log.md#workspace-and-spawn) · [Structured worktree lookup](run-log.md#child-return-channel)
3. **Cover completion, recovery and closure.** Verify report freshness and clean path-scoped commits. Fixture-test duplicate notifications, ambiguous creation, an ended worker without a report and timeout during progress. Confirm that a second failure holds only its entity. [Completion](run-log.md#completion-and-durable-proof) · [Timeout](run-log.md#steering-and-roster) · [Retry](run-log.md#retry-and-budget-contracts) · [Closure](run-log.md#shutdown-and-cleanup)
4. **Qualify child questions when the Paseo daemon ACL is available.** The return path already worked with the daemon grant. Enable parent-directed messaging for ordinary workers after testing the scoped permission in a disposable fixture. This dependency applies to child-initiated messages. [Labelled child and matching message receipts](run-log.md#child-return-channel)

## Definition of done

The first three tasks deliver the fresh-worker lane. Through the installed binding, a real Claude FO launches one Codex worker, steers it during the run, receives its finish event, verifies the code and split-root report commits, and observes the worker closed. The acceptance record includes the requested and reported Paseo modes. Both probe launches requested `auto` and returned `full-access`; this is a Paseo-side finding. [First launch](run-log.md#workspace-and-spawn) · [Second launch](run-log.md#child-return-channel) · [Commit assertions](run-log.md#completion-and-durable-proof) · [Shutdown assertion](run-log.md#shutdown-and-cleanup)

## Out of scope

Remote daemons, worker reuse across stages, other worker runtimes and implementation of the Paseo daemon ACL. The acceptance test covers same-daemon shared paths and fresh workers. [Fixture paths](run-log.md#fixture-and-dispatch) · [Cross-host test boundary](run-log.md#cross-host-envelope)
