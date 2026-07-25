# First Officer gate-lifecycle topology audit

Reviewed WIP commit: `cabdef33`

Classification: bounded topology reset; lifecycle ACs unchanged.

The six-event lifecycle behavior is useful, but its boot-resident placement is invalid:

- shared core grows from 26,092 to 32,289 bytes, exactly +6,197;
- the hard guard requires `<26,755`, leaving only 662 bytes of original headroom;
- Claude, Codex, and Pi prompt-load ratchets each fail by the same 6,197 bytes;
- fitting the procedure into the available headroom would require roughly 89% compression and would remove load-bearing provenance, failure, stale, and resume rules.

The correct topology is a deferred, non-user-invocable `fo-gate-lifecycle` skill loaded on the first engaged gate. The boot core retains only a small trigger and routes every gate entry—already-gated startup, engage, worker completion, and resume—through it before package mutation, validation, presentation, or decision handling. The existing boot-core ceiling stays unchanged; the deferred module must be honestly included in worst-case host prompt accounting.

Additional material gaps in the checkpoint:

- the headless/no-conn path still invokes `present-gate` directly before bind/open validation;
- the shipped-skill mutant is static rather than the required live mutated-skill replay;
- actor-swap and raw-review-dump mutants are absent;
- the full AC-5 refusal and AC-7 stale/resume matrices are incomplete;
- workflow-discovery before/after candidate equality is not proved;
- Codex consumes approval but produces no observed successor spawn/handle/worker output.

Claude and Pi authentication failures occurred before workflow behavior and remain external validation conditions. They require later green evidence but do not justify changing the design or ACs.
