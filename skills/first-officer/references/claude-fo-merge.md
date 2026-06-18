# First Officer Merge Module (Claude)

The Claude part of Merge-and-Cleanup step 10 (fo-merge-core.md), read alongside the core at the terminal boundary: the per-name worker teardown.

## Terminal Teardown (Merge-and-Cleanup step 10)

10. **Teardown workers at terminal.** Derive the entity's worker cohort from the live roster — every worker whose handle decomposes to this entity's slug (roster and decomposition are the adapter's). Shut each one down per-name with the cooperative `SendMessage({"type":"shutdown_request"})` and drop them from session memory; the auto-team `members[]` prunes each terminated member. There is no bulk team-delete, no settle interval, and no attempt cap on this path. On a legacy host the teardown is the bounded `TeamDelete` settle-and-cap procedure in `Skill(skill="spacedock:using-legacy-claude-team")`, loaded only on the probe match.
