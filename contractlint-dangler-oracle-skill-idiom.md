---
title: "contractlint dangler-oracle does not resolve the Skill() cross-file pointer idiom (referenceProsePointerDanglers lags the nt ref->skill move)"
source: "0240 pre-cut antipattern audit (2026-07-01, lens 2). internal/contractlint/boot_resident_closure_test.go referenceProsePointerDanglers resolves a cross-file prose pointer to a watched section ONLY via an intra-file heading or a references/*.md path token (bodyReferenceRe). nt made the correct cross-file idiom for fo-write-core/fo-status-viewer the Skill(skill=spacedock:...) form and deleted the references/*.md paths (now dead per deadDeferredReferencePaths). A future contract line pointing at one of these sections via the Skill() form while naming the section phrase would FALSE-FAIL as dangling. Latent (no such line today); the guard's resolution model lags the idiom it introduced. Non-blocking; captain fast-tracked."
status: ideation
group: tooling
id: fezpp9ddsaz6bhek6ra2se69
started: 2026-07-02T01:37:56Z
---

## Problem
`referenceProsePointerDanglers` (internal/contractlint/boot_resident_closure_test.go ~L288) resolves a cross-file prose pointer to a watched section (e.g. "FO Write Scope", "Status Viewer") only via an intra-file heading or a `references/*.md` path token. nt's ref->skill move made the correct cross-file idiom `Skill(skill="spacedock:fo-write-core")` and deleted the reference paths — so a future line pointing at a watched section via `Skill()` while naming the section phrase would false-fail as dangling.

## Desired fix
Teach the oracle to ALSO resolve a watched-section prose pointer when the same context names the owning skill via the `Skill(skill="spacedock:<name>")` token — mapping the skill to its SKILL.md section anchors (the deferredSkillCores map already keys skill path -> anchors). The guard's resolution model then matches the Skill() idiom nt introduced.

## Rough acceptance sketch (ideation tightens into measured ACs + test)
- RED-then-GREEN: a contract line pointing at a watched section (e.g. "FO Write Scope") while naming `Skill(skill="spacedock:fo-write-core")` RESOLVES (no false-fail); a genuine dangling watched-section pointer (no intra-file heading, no path token, no matching Skill token) STILL reds.
- The existing dangling-target + prose-pointer controls stay RED-capable (non-weakening).
- go test ./internal/contractlint/ green.
