---
id: k74gt0qv3j4b86knvy2rhsta
title: Lighten the interactive boot greet — managed-workflow list + an explicit "engage" verb; no forced workflow pick, no gate render at the greet
status: backlog
source: "FO session 2026-07-04: interactive boot ran ~5 minutes; original framing was Startup step 8 rendering a full present-gate review per ready gate before stopping. Captain's subspace-tui review (2026-07-04) redirected the scope wider: the launcher bootstrap prompt (frontdoor.go bootstrapPrompt/codexBootstrapPrompt, both literally ending '...Engage.') should drop that flourish since 'engage' becomes a real captain-invoked verb; Startup step 3's multi-workflow pick goes away in favor of just listing managed workflows; gate assembly moves behind the engage verb. science-officer flagged 4 open questions (is engage a binary command or an FO interaction verb; what 'sweep in-flight tasks' means precisely; single-workflow vs all-managed-workflows scope; confirm the frontdoor.go referent) plus a concrete collision: 7v (pi-bootstrap-prompt-parity) pins pi byte-identical to the CURRENT codexBootstrapPrompt text including 'Engage.' — must coordinate/sequence."
started:
completed:
verdict:
score: 0.35
worktree:
issue:
---

The interactive FO boot is too heavy: the launcher bootstrap prompt (`internal/cli/frontdoor.go`) ends with a flourish "Engage.", and Startup step 8 directs the greet to render a full `present-gate` review for every ready gate before stopping, while step 3 asks the captain to pick when multiple workflows are discovered. Redirect (captain, 2026-07-04): make the greet light — just list the managed workflow(s) and hint "Use engage"; do not force a workflow pick; do not render gates at the greet. Move the in-flight sweep (gate assembly / ready-task processing) behind an explicit captain-invoked `engage` verb. OPEN (ideation must resolve WITH the captain before designing): is `engage` a new binary command or an FO interaction verb built on the existing `status --next` sweep; what exactly "sweep in-flight tasks" does (dispatch / present gates / advance); and whether `engage` operates on one workflow or across all managed ones. Headless/single-entity mode still needs a resolved single workflow — the no-pick change is interactive-only. Do NOT dispatch ideation until the captain settles engage's nature; the entity boundary (one entity vs. an engage-command + greet pair) depends on it.
