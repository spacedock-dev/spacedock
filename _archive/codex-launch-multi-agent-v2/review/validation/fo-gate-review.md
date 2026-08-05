# z5 validation gate review

Entity: `codex-launch-multi-agent-v2`
Stage: `validation`

The fresh Sol/medium validator completed the required full, race, focused,
formatting, credential-gated, and detached-audit checks against cycle-3
commit `1a8a64cdb`. The report has five DONE items and two FAILED items.
AC-1 through AC-4 are cited, but cycle 3 found one material evidence defect
and the workflow requires Captain escalation.

Findings:

- The corrected lifecycle test uses the built `spacedock codex` front door,
  typed ordered records, exact cardinality, v2 context, and a zero-event
  disabled control. However, raw byte search combines identity fields from
  unrelated JSON: an actual `other-parent/worker-b` child passes when a
  separate record carries the expected parent-thread and agent-path text.
  Same-worker identity is still not proven.

Science Officer advisory: HOLD and present to Captain. Cycle 3 escalation
forbids another automatic feedback round.

Recommendation: HOLD pending Captain disposition. The proposed bounded fix is
to atomically decode each child identity record, require `{parent_thread_id,
agent_path}` together in the same typed record with exactly one matching child,
retain explicit and targetless misbinding negatives, and rerun validation. An
alternative is a formal Captain scope reset that changes the same-worker AC;
do not weaken it silently.
