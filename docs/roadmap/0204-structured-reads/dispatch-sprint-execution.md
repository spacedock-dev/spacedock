# 0204 — dispatch-sprint-execution (cold-boot Commander package)

The self-contained package a **Commander** runs from a cold boot to drive the 0.20.4 "structured reads" sprint to its deliverable. Per `docs/roadmap/README.md`: the shaping FO packages the sprint; the Commander boots `spacedock:first-officer`, creates its own team, and drives. This is a cold-boot brief, not a context transfer or a state tracker.

You are the **Commander** for `0204-structured-reads` in the `docs/dev` workflow. Split-root: state in `.spacedock-state` (branch `spacedock-state/dev`); product on `main` (post-flip Model B trunk). **Do NOT change the main worktree's branch.** Boot `spacedock claude` → the FO contract → `status --boot --json`.

## Your charge
Drive the `0204-structured-reads` members to the sprint DoD and cut `v0.20.4`.

- **Durable strategy = the index.** `docs/roadmap/0204-structured-reads/index.md` (goal / scope / DoD / deliverable). Do not turn it into a tracker.
- **Membership + state = the query, never a static list.** `spacedock status --workflow-dir docs/dev --where sprint=0204-structured-reads`. Fresh-pickup queue: `--where sprint=0204-structured-reads --where 'sprint-readiness != defer'`. The backbone `e6a` is `sprint-readiness=in-progress`; the four ready members are `sprint-readiness=ready`.
- **Per-task needs + dependencies live in each entity.** Read the entity body/report — there is no roster here by design.

## Recommended drive order (re-derive from the entities)
1. **Unblock first — `5h0` (`ban-readme-substring-assertions`).** Not a sprint member, but it BLOCKS both `e6a`'s merge and the README-slim push. The dev-README slim (`a9e669ae`, local `main`, unpushed) broke two prose-grep doc-contract guards (`TestSharedScenarioDocsContract`, `TestCodexForegroundWaitWatchdogDocsContract`); 5h0 carries the trigger note. Retarget the guards to `docs/runtime-live-ci.md` or convert per the proof policy. Until then local `main` is RED; `origin/main` is green (slim unpushed).
2. **`e6a` (backbone)** — parked at `status: implementation` (worktree `.worktrees/spacedock-ensign-status-section-reader`, branch `spacedock-ensign/status-section-reader`, commits `ebbd9ba3` + `e75fede2`, local). Validation owes the **AC6 live `--read` drive** at the `:105` FO completion-gate site + the **detached adversarial audit**. Merge after 5h0.
3. **The ready set, by score** — `0q` (drop the default SOURCE render; ~7.2k boot saving) → `48` (template restructure; coordinate the `development.md` overlap with `ey`) → `j7` (**BIMODAL**: spike the staleness-echo mechanism FIRST; if harness-inherent, ship a roadmap decision, not code) → `6r` (CI-log summary).

## Proof discipline (this workflow)
No prose-grep over instruction files. Detached adversarial audit before merging a high-stakes surface (shipped contract/scaffolding, `status` paths, launcher, CI/release). A contract or skill change is PASSED only when a live drive observed the behavior.

## Coordination
This checkout is shared by other Claude sessions: a **shaping FO** (filed + ideated 0204) and a separate **0203 Commander**. **Stay out of 0203 members.** Commit state path-scoped; `git -C docs/dev/.spacedock-state pull --rebase` before writes. Cross-session notify = append an unread message to the peer team's `~/.claude/teams/{team}/inboxes/team-lead.json`.

## Boot
`spacedock claude` → FO contract → `status --boot --json` → run the membership query → read the entities → drive.
