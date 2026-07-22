# Debrief 2026-07-22-02 — FO session (0.26.0 release cut + post-release fixes)

Scope: cutting the 0.26.0 stable release and shipping the post-replay/post-release fix wave. Continues the 0260-proportionality sprint's theme into practice — the captain repeatedly steered fixes *smaller* and priced over-engineering.

## Captain directives (binding on future sessions)

- **CI-lane meaningfulness — "loads vs drives."** Merge on the lanes that actually *drive* the changed code path, not merely *load* a changed file. A live lane that compiles a touched file but never exercises the change adds no coverage → waive it. Applied this session to upgrade-remedy (too-old remedy fires only on version mismatch, never in a live drive), ab (discovery-dir prune), and auto-pre0 (release.yml runs only on a real tag). Captain set this via "tell me which ci lanes are meaningful, approve those." This refines the dev-README proof-policy rule "required lanes are a function of the diff." Candidate to fold into the proof policy prose.
- **Scope to the reported bug; ideation must not over-scope.** upgrade-remedy's ideation built all four install sources (+416 LOC) for a one-source (`@next`) bug; captain trimmed it to `@next`-only (+109). Expand past the reported bug only on explicit direction.
- **A confirmed tautological test: consider delete before give-teeth.** mp gave teeth to four tautological tests; #1 (`activeSessionFile_would_flip_to_teammate`) tested a test-local *port*, not production, and should have been deleted. Ask "should this test exist / does it guard production?" before adding an oracle. (Deleted as a follow-up, #554.)
- **`main` is not branch-protected** — required CI is FO policy, not a GitHub gate. Merges are the FO's judgment call against the diff→lane mapping, not a hard block.

## State — what shipped

- **0.26.0 released** (v0.26.0 stable Latest; v0.27.0-pre0 edge; stable cask 0.26.0; edge cask 0.27.0-pre0). The e2e-gate was WAIVED via `SPACEDOCK_E2E_GATE_WAIVER` (release commit = green-validated #546 content + stamps/docs; no FO-behavior delta) and the waiver **removed after both the stable and the auto-pre0 gates passed** — gate re-armed. The auto-cut `v0.27.0-pre0` did not fire its run (the known `auto-pre0-tag-push-release-trigger` bug); recovered via the documented delete+replay.
- **Post-release fixes, all merged to main:**
  - #548 mp — four tautological tests given independent oracles.
  - #549 joint q4+9q4 — retire legacy TeamCreate path + replace session-wide Degraded Mode with a per-entity bounded retry rung + durable `### Dispatch Retries` ledger; rename "back-channel" → "inter-agent communication"; widen-scope retirement of the orphaned legacy teardown-grade machinery (captain-directed). claude FO load −15.8 KB.
  - #552 ab — relocate the three prose-skill fixtures from top-level `fixtures/` to `skills/integration/testdata/` (Go-idiomatic package-adjacent; no Go consumer, captain chose the home); prune `testdata` from discovery.
  - #553 auto-pre0 — push the pre0 tag over a repo-scoped SSH **deploy key** (`EDGE_RELEASE_DEPLOY_KEY`, provisioned this session) + a verify-or-fail poll on `edge-advance`.
  - #554 mp-cleanup — drop the weak `activeSessionFile` characterization subtest.
  - #556 upgrade-remedy — emit `brew upgrade spacedock@next` for `@next` edge-cask installs (the `@next`-only trim).

## State — mechanisms & learnings

- **auto-pre0 deploy-key trigger is spike-proven.** A runner-embedded deploy-key SSH tag push *creates* a `release.yml` run (observer run 29884082483); a `GITHUB_TOKEN` push creates zero. The credential was the sole variable. The next real stable cut self-proves it (the verify-or-fail guard reds `edge-advance` if the pre0 run doesn't appear). Deferred hardening filed (below).
- **SSH push WORKS this session** (`git@github.com:...`) and is REQUIRED for any branch touching `.github/workflows/**` — the local osxkeychain HTTPS credential is a stale OAuth token without `workflow` scope (it silently blocked auto-pre0 until switched to SSH). **This corrects the 2026-07-20-02 note "gh's ssh protocol setting has no working key"** — the SSH key pushed #553 and #556 fine.
- **Transport-stall → nudge validated.** The upgrade-remedy trim ensign died mid-response with "API Error: Connection closed mid-response"; a resume nudge recovered it (no degradation). Live confirmation that 9q4's retry-rung design (shipped in #549) is correct — an API transport failure is nudge-recoverable, not a degrade trigger.
- **Live-lane flakes are model-specific.** The joint's `claude-live (sonnet)` failed `TestLiveDefaultHeadlessStopsAtGate` (FO drove the draft but left status `draft` instead of advancing to `review`) while opus/codex/pi passed the identical host-neutral scenario → model variance, not a contract regression; re-run-to-green.
- **Local `main` is behind `origin/main`** (left untouched per the captain's parallel durable-decisions/0270 work). The local checkout still carries the pre-ab-move `fixtures/`, so `spacedock new`/discovery hit the multi-workflow ambiguity locally — pass `--workflow-dir docs/dev`. New worktrees based off `origin/main` are clean.
- **Minor:** a git-ignored goreleaser `dist/` snapshot raced contractlint's repo-root FS walk → a transient `-race` flake; settled on re-run. Consider cleaning stale `dist/`.

## Follow-ups filed (backlog)

- `commit-doctor-wire-test-edge-cask-remedy` (c0gaq7cc) — the #556 `gateHost → remedy` wire has no committed test (value is live-proven, the wire is not); add a committed doctor test staging a fake `Caskroom/spacedock@next/` path.
- `harden-edge-advance-verify-poll-actions-read` (bd94c7te) — add `actions: read` to the `edge-advance` job so the verify poll's Actions read doesn't rely on the repo being public (~1 line; push over SSH).
- `classify-control-plane-work-plane` (j7st8gr3, issue #555) — add a `work.classify` step before `write.classify` so control-plane work is not forced into meta-entities to obtain write/dispatch authority.
