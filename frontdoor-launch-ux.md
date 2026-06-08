---
id: yqf0amtyecjcft0vsw6nbqtk
title: spacedock codex/claude front-door launch UX — auto-install messaging, pre-launch info banner, neutral bootstrap prompt
status: ideation
source: "captain live-test of 0.19.8 (2026-06-08). The first real `spacedock codex` run surfaced three front-door UX issues (A/B/D below). The pre-cut audit verified z9's auto-install SEAM via tests but nobody ran the front door end-to-end — so a real `spacedock codex` live drive is this task's load-bearing proof."
started: 2026-06-08T21:00:22Z
completed:
verdict:
score:
worktree:
issue:
group: binary-ux
sprint: 0199-pre-flip-mechanics
sprint-readiness: ready
---

Make the `spacedock codex` (and `claude`) launch experience honest and useful: don't tell the user to install manually right before silently auto-installing, show a short pre-launch info banner, and ship a neutral bootstrap prompt.

## Problem

- **A — the auto-install message is self-contradicting.** On `NoPluginFound`, `gateHost` (`frontdoor.go:124-128`) prints the *manual-install remedy* — "no installed codex plugin found. Run `spacedock install --host codex` (or `spacedock claude --skip-contract-check` to bootstrap)" — and then `runCodex`/`runClaude` **auto-install anyway** (`ops.Install`, silently — `execHost.Install` writes no progress). So the tool says "install it yourself," then quietly does it. Two defects: the manual remedy fires on the auto-install path, and it hardcodes a `spacedock claude` hint that is wrong in a codex run.
- **B — no useful pre-launch info.** `spacedock codex`/`claude` launches straight into the host with no Spacedock context (version, which workflow it detected, anything actionable).
- **D — the shipped bootstrap prompts carry personal flavor text.** Both `bootstrapPrompt` (`frontdoor.go:24`) and `codexBootstrapPrompt` (`:289`) literally read "You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage." — so the FO tries to *relay* that at an empty-team session start ("no team to relay your note to"). The product default should be neutral.

## Proposed approach

Ideation firms. Sketch:
- **A:** on the auto-install path, suppress the gate's manual remedy and print "Installing the {host} plugin…" before `ops.Install`; make any bootstrap hint host-aware (no `spacedock claude` in a codex run). Keep the manual remedy ONLY on the `--no-install` / hard-fail branch.
- **B:** a short pre-launch banner — `Spacedock vX` · workflow detected (`docs/dev` or none) · 1–2 useful lines — before the host launches.
- **D:** neutralize the shipped bootstrap prompts: drop the personal/relay text, keep the functional `Assume $spacedock:first-officer for the entire session.` (codex) and the equivalent FO-select for claude.

## Out of scope

- The codex plugin install MECHANISM (z9 — works; this is messaging/UX around it).

## Acceptance criteria

Ideation/implementation fills in. Sketch — all behavioral halves proven by a LIVE DRIVE (front-door surface; a grep does not satisfy):
- **AC-A** A real `spacedock codex` with no plugin prints "Installing the codex plugin…" and NO `spacedock claude` hint; the `--no-install` path still prints the manual remedy. (live drive + argv/message oracle.)
- **AC-B** The pre-launch banner renders the version + detected-workflow line. (live drive.)
- **AC-D** The launched inner prompt contains NO personal/relay text ("I love you", "tell all subagents"). (a launch-argv oracle over the bootstrap const + a live drive.)

## Test plan

argv/message oracles for the install-path message + the bootstrap-prompt const, PLUS a real `spacedock codex` (and `claude`) **live drive** observing the banner, the "Installing…" line, and the neutral prompt — the end-to-end front-door run nobody did before z9 shipped. Front-door = high-stakes → detached adversarial audit at validation. Watch overlap with `th` (also edits `frontdoor.go`, disjoint funcs: th = launchEnv/argv-prefix on wrap; this = gateHost message + banner + bootstrap const).
