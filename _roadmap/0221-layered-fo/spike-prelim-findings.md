# 0205 spike — preliminary run findings (banked 2026-06-16)

The `haiku-loop-spike` (`w4`) was partially executed live this session as a de-risking probe,
then banked. Two of its ACs ran against a bare `claude --model haiku` launch (authed via
`~/.claude/benchmark-token` -> `CLAUDE_CODE_OAUTH_TOKEN`, the operator-local `internal/ensigncycle`
recipe: isolated HOME, drop `CLAUDECODE`, spacedock on PATH). The full dispatch -> gate -> merge
drive is deferred to a credentialed lane or a Commander.

## What ran

- **AC-1 (launch path) — PASS** (~$0.047). Bare `claude --model haiku-4-5` boots in an isolated
  HOME, authenticates via the benchmark-token (`apiKeySource: none`), emits the full stream-json
  event stream, returns the smoke token. The staff-review-flagged-unproven bare-launch path is
  proven runnable here (it is NOT the `spacedock claude` front door, which would force the full
  contract).
- **AC-2 (loop-substitution, smallest slice) — PASS, N=1** (~$0.029). Given only the simplified
  hand-loop (`boot -> next -> advance -> commit -> stop`) and NOT the FO contract, the Haiku FO
  followed the prose-functions exactly in order over a throwaway single-root fixture: it advanced
  `probe-task` backlog -> review via `status --set`, committed, then stopped. Durable-state graded
  (entity at `review`, commit `advance: probe-task -> review` landed). It did not improvise, did not
  do the stage work itself, did not skip the commit. 4 turns.

Total spend ~$0.08.

## The load-bearing caveat (captain, 2026-06-16)

**Haiku reliability is the open question, not capability.** In earlier days Haiku passed these
tests; at some point it became FLAKY. A single pass (AC-2 N=1) establishes only that Haiku CAN
follow the loop once, not that it HOLDS it. The spike's real product is a reliability measurement,
not a capability demo.

This sharpens the AC-4 must-build decision rule: **the must-build signal is FLAKINESS, not just
always-breaks.** A loop step Haiku follows *sometimes* and botches *other times* is exactly the
unreliability a binary verb removes, so any step that is flaky across N>=3 (ideally tracked across
haiku versions) is must-build, not only steps that break every time. "Holds reliably" requires
clean across all N, and even then stays provisional and version-sensitive given the documented
historical regression.

## Still owed (the full spike)

- The smallest slice exercised no `«dispatch»` (nested worker), no `«gate»` (route-to-L3), no
  `«merge»`. The high-value failure modes — bare-dispatch, auto-approve-gate, idle-vs-completion —
  are NOT yet exercised.
- N>=3 per the staff-review fold (N=1 here), and reliability tracked across haiku versions per the
  flakiness caveat.
- The final drive on the real `spacedock claude` shape (the `haiku-drive-validation` member), not
  the bare hand-loop proxy.

## Disposition

Banked. `w4` stays at `ideation`, gate-pending. The full drive is a credentialed-lane / Commander
exercise once 0.20.3 (2y) and 0.20.4 land. These results + the flakiness caveat are its input.
