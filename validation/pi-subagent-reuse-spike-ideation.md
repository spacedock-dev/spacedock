# Pi subagent reuse spike ideation summary

Updated `docs/dev/.spacedock-state/pi-subagent-reuse-spike/index.md` with spike findings and recommendation.

Key result: local Pi setup was safe enough for a temp workflow experiment. A persisted child session produced an initial marker and a resume-equivalent revived/forked follow-up produced a second marker; a fresh comparison produced a third marker. Direct `subagent({ action: "resume" })` was documented but not cleanly proven as a visible non-interactive tool call.

Recommendation recorded in the entity: keep fresh redispatch as the default for Pi validation `feedback-to: implementation`; reserve resume for manual/debug or future opt-in tooling with registry epochs, durable handles, and child-isolated token/cost evidence.
