# Re-ideation gate review: deferred First Officer gate lifecycle

## Capability and change

Move the detailed six-event gate lifecycle out of the boot-resident First Officer core into one deferred, non-user-invocable `fo-gate-lifecycle` skill. Retain only a small common entry trigger covering every engaged gate route before mutation, validation, presentation, decision, or resume handling.

The lifecycle, provenance, fail-closed routes, package rules, prompt ceiling, and strict observed-spawn ACs are unchanged.

## Evidence and reviewed snapshot

- WIP `cabdef33` proves much of the behavior but adds 6,197 bytes to a core with only 662 bytes of headroom.
- Extracting the 6,222-byte procedure leaves the core around 26,067 bytes before triggers; final core is capped at 26,754 and deferred skill at 6,600 bytes.
- The deferred skill is registered in address lint and every host's honest worst-case prompt set; the boot-core guard is not rebaselined.
- One route table covers already-gated startup, headless conn/no-conn, interactive engage, worker completion, and open/pending/revise/hold/stale/consumed resume. Non-gated greet stays load-free.
- The checkpoint's headless direct-presenter bypass is explicitly deleted.
- Codex consume-without-spawn remains an implementation blocker. Claude/Pi authentication remains an external validation condition requiring later green runs.
- Missing live shipped-skill, actor-swap/raw-dump, refusal, resume, discovery, and successor-attribution controls are explicitly required.

The exact revised entity, landed gate contract, and prior topology audit are identified by URI and SHA in `briefing.json`.

## Recommendation and decision

Recommendation: **approve re-ideation and proceed to fresh implementation** on the preserved branch. Do not enter validation without a real successor spawn/handle/worker output and all prompt/load guards green.

Decision requested: approve, revise with a concrete remaining topology gap, or hold for a named prerequisite.
