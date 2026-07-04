---
id: bd4tpadcw0efgvjx6tavym0g
title: "Persist the science-officer as a standing workflow mod so future FO sessions spawn it automatically"
status: backlog
source: "FO session 2026-07-04: the science-officer (an opus standing teammate that turns raw operational friction into filable backlog entities) was spawned ad hoc via a hand-assembled Agent() call mirroring _mods/comm-officer.md, because no mod file declares it. Confirmed: docs/dev/_mods/ holds only comm-officer.md and pr-merge.md, and dispatch spawn-standing-all --workflow-dir docs/dev returns exactly one spec (comm-officer). Drafted by the science-officer teammate itself."
started:
completed:
verdict:
score: 0.28
worktree:
issue:
---

The science-officer role currently exists only as a manually re-assembled Agent() call each session, invisible to dispatch spawn-standing-all and lost at teardown. Add docs/dev/_mods/science-officer.md mirroring _mods/comm-officer.md's shape (standing: true, Hook: startup with model: opus, Hook: shutdown, a Routing guidance section, and an Agent Prompt carrying the charter this session used: read the workflow README schema/stages/proof-policy, skim backlog+ideation for overlap before drafting, produce per-issue entity-seed blocks the FO can file, stay standing). Mod files are off-limits for direct FO edit, so this goes through a dispatched worker in a worktree at implementation stage like any other product change; low blast radius (workflow-local mod, not shipped skills/ scaffolding) so no detached adversarial audit is owed, but a live spawn from the emitted mod spec should prove the prompt content actually boots into the role, not just that the file parses.
