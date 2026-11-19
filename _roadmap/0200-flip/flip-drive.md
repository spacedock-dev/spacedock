# Drive the 0.20.0 main-flip (`pj`) — cold-boot driving prompt

Goal: cut **0.20.0 on `main`** and flip the marketplace to serve the stable plugin from `main` (keep `next` as the dev/edge channel). **NO contract change.** This is **outward-facing — the captain approves EVERY outward step** (archive, force-push, tag, marketplace); you execute and pause at each.

**Read first — the authoritative plan:** `docs/dev/.spacedock-state/main-flip-0200-marketplace/index.md` — the 9-step runbook, 8 ACs, and the flip checklist. Follow it exactly; this prompt is just the spine.

**Boot:**
```bash
git fetch origin next && git reset --hard origin/next && go build -o ./spacedock ./cmd/spacedock
git -C docs/dev/.spacedock-state pull --rebase origin spacedock-state/dev
security find-generic-password -s "Claude Code-credentials" -w | python3 -c "import sys,json;print(json.load(sys.stdin)['claudeAiOauth']['accessToken'])" > ~/.claude/benchmark-token
./spacedock status --workflow-dir docs/dev --boot
```
Then `Skill(spacedock:first-officer)`, TeamCreate. Advance `pj` to `implementation` as you begin.

**Precondition (met):** the pre-flip set is merged on `next` — `nzb` (e2e-gate), `k6d` (two-channel devBranch), `cmx` (front-door banner), all PASSED. Confirm `go test ./...` is green from the repo root before anything outward.

**Spine (per pj's runbook — 📡 = captain-gated outward step):**
1. Pre-cut antipattern audit on the prepared `next` line.
2. Fold all 0.20.0 content onto the `next` tip: marketplace `ref` `next`→`main` **+ the paired `skills/integration/marketplace_manifest_test.go` edit in ONE commit** (green-by-construction); `AGENTS.md` + `docs/releasing.md` reconciliation.
3. **Freeze `next`**; record `OLD_MAIN=$(git rev-parse origin/main)` and `$PREPARED=$(git rev-parse origin/next)`.
4. 📡 Live e2e (`workflow_dispatch`) on `$PREPARED` → wait green → **re-verify the green run's `headSha == $PREPARED`**.
5. 📡 Archive `origin/main` as **`archive/v0`** — a NON-`v*` ref (`v0-archived` matches the release glob and would fire goreleaser).
6. 📡 Guarded flip: `git push --force-with-lease=main:$OLD_MAIN origin $PREPARED:main`; verify `origin/next` still resolves.
7. 📡 Tag `v0.20.0` (annotated) + push → `release.yml` fires; its `e2e-gate` matches the green run from step 4 → goreleaser publishes; `k6d`'s stamp commits on `main` post-tag.
8. Post-flip: calendar-bump on `main`; verify a fresh-HOME install resolves `0.20.0` from `main`; run the upgrade journeys (incl. the `next`-pinned-user one); release the `next` freeze.

**Watch:**
- **Tag-the-green-tip invariant:** the green-run commit, the force-pushed `main` tip, and the tagged commit are ALL ONE SHA. Nothing lands on `next` between step 3 (freeze) and step 7 (tag).
- **Node-20 CI:** if today is ≥ ~2026-06-16, bump `actions/checkout@v4` / `setup-go@v5` / `goreleaser-action@v6` BEFORE the cut runs `release.yml` (cleanups-0199 item #5) — the flip's cut depends on a healthy `release.yml`.
- **NO contract change** (stays 1; 0.20.0 is the first real stable release). The captain authorizes each 📡 step.
