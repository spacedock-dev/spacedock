# Mods & standing teammates

Mods extend a workflow with lifecycle hooks and standing agents — behavior layered on top of the base stage primitives, without changing the workflow's stage definitions.

## Lifecycle hooks

A mod registers behavior at a point in the workflow lifecycle — startup, idle, or merge. The first officer runs registered hooks at the matching lifecycle event. For example, a startup hook can spawn a standing teammate before the normal event loop begins; a merge hook can guard terminalization so an entity cannot be closed while a merge is pending.

Mods live under the workflow's `_mods/` directory. They are convention-driven prose the first officer reads and acts on, not a separate binary.

## Standing teammates

A standing teammate is an agent kept alive for the session, available on demand rather than dispatched per stage. The **comm-officer** is the canonical example: a prose-polishing teammate spawned by a startup hook and kept resident for the captain session.

The comm-officer polishes drafts about to be presented to the captain — PR bodies, gate review summaries, entity narrative sections — and the entity body prose before it's committed. It does not polish live chat replies, short operational statuses, or commit messages. Polish is best-effort and non-blocking: callers proceed with un-polished text if no reply arrives within two minutes.

The comm-officer is light-touch by default and defers to a project voice guide when one applies. For Spacedock's docs, that guide is the [Voice & tone](../contributing/voice-and-tone.md) page.
