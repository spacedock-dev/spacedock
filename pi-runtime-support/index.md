---
title: Pi runtime support — adapt Spacedock to pi-native teams/subagents
status: backlog
score: "0.36"
source: captain (2026-06-03) — dogfood pi support from the prior PR #155 and evaluate pi-agent-teams / pi-subagents as usable ensign constructs
issue: spacedock-dev/spacedock#155
id: s9kcdyb9r5t8addppnnce54j
---
# Pi runtime support — adapt Spacedock to pi-native teams/subagents

Bring forward the useful parts of the old PR #155 Pi runtime compatibility baseline into Spacedock v1, but design it for the current Go launcher, split-root state model, and current pi ecosystem.

The starting assumptions for this ideation pass are:

- pi support is allowed to be real, not theoretical;
- `pi-subagents` can act as a Spacedock ensign dispatch construct for dogfooding;
- `pi-agent-teams` may be usable as a higher-level Pi team substrate, but its tool signatures are not Claude-compatible and need adapter design;
- old PR #155 used a lower-level session-registry helper model that may still be valuable.

## Seed questions

- What is the smallest v1-compatible Pi runtime adapter that lets the first officer dispatch, wait, route follow-up, and shut down an ensign?
- Should the first slice target `pi-subagents`, `pi-agent-teams`, a local helper like PR #155, or an abstraction that can support more than one Pi substrate?
- What changes belong in skill runtime docs, `spacedock dispatch build`, a new `spacedock pi` front door, package manifests, or tests?
- What live/dogfood scenario proves this is real without relying on prose grep or a fake task?

## Out of scope

- Rewriting existing Claude or Codex runtime behavior.
- Adding PR merge or mod behavior specifically for Pi.
- Treating Claude `Agent` / `SendMessage` tool signatures as available in Pi unless an adapter explicitly provides them.
