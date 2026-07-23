---
id: frze3yqm9da0vp0r53qqdc8t
title: Record advisory review rounds without gate semantics
status: backlog
source: "02av deferred round-recorder plumbing and 3j jobs 592/594/597 incident, 2026-07-23"
started:
completed:
verdict:
score:
worktree:
issue:
---

Provide one owned write surface for correction-round Briefings, reviewer Annotations and advisory Resolution, and the worker's triage Resolution, without selecting a gate or advancing workflow state.

## Problem

02av requires a correct-but-disproportionate finding to become a durable decline on the reviewed round, but deliberately deferred Roborev ingestion, room creation, ordered-log append, frontmatter pointer, and Feedback Cycles projection. The 3j incident consequently retained jobs 592/594/597 and the final disposition only in prose; its downstream task had no round room, and a gate-record backfill was both semantically wrong and mechanically unavailable.

## Minimum value demonstration seed

In a disposable workflow fixture, supply exact reviewed commit `90aea55`, retained reviewer outputs corresponding to jobs 592 and 594, and the worker's duplicate-member decline. One recording operation writes a digest-bound round Briefing, reviewer finding Annotation, reviewer advisory Resolution, worker decline Annotation, worker advisory Resolution, and a readable Feedback Cycles projection. A read-back proves both Resolutions are advisory and that gate selection, application, status, and the reviewed snapshot are unchanged. As the red control, corrupt the reviewed-snapshot digest: the operation must fail and leave every tracked byte unchanged.

## Boundary

This task owns generic advisory-round persistence and an ergonomic one-off backfill path. It does not launch or poll Roborev, decide materiality, apply suggestions, create binding Decisions, select or consume gates, or stop independent work.
