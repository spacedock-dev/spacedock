# Release 0.28.1: plan and defect findings

Status: draft for CL review, 2026-11-17. Nothing in this plan is pushed, tagged or published.

CL decisions (2026-11-17): include #884. The channel is sticky: a pre-release install stays on edge, so defect 1 is closed as working as documented. Skip defect 2. Fix defect 3. Defects 4 and 7 wait for CL after a clarification.

All findings are against `origin/main` at `647d55163`. The local `main` checkout is at `b2524f043`, 22 commits behind `origin/main`.

## Changes on main after v0.28.0

`v0.28.0` and the auto-cut `v0.29.0-pre0` both point to `2fcbd10bc`. `origin/main` adds these changes. The version-stamp commits are not listed.

1. PR #875, install counting: `install.sh` looks up releases through `spacedock.md/latest-release/...` with a `spacedock-install/<version>` User-Agent. The PR also bounds the lookups and validates tags. `install.sh` is served from `main`, so this change is already live for curl installs.
2. Install guide (`ad22447e1`): the guide now installs the stable channel. The install hint test follows it (`eff994eff`).
3. PR #881 (Kent): the update and install path for a channel. `install` adds the channel marketplace before it removes the Claude Code plugin. `update` checks each host that holds a Spacedock piece, fails when a piece or a Pi extension package is missing, and shows the refusal from Claude Code. `doctor` adds a MIXED line, and the brew upgrade has one route. Also `34bc3813e`: `update` leaves a Codex `--plugin-dir` dev install alone.
4. PR #884 (Jared): new feature. Before an interactive launch, the binary checks once a day for a newer Spacedock through `spacedock.md`. The disclosure names GitHub. New package: `internal/updatecheck`.
5. Test changes only: PRs #878 and #879 change the live-journey graders and add a Haiku XFAIL.

## Defect findings

| # | Defect | Real on origin/main? | Where | Fix |
|---|--------|----------------------|-------|-----|
| 1 | `update` takes a pre5 install to 0.29.0-pre0, not 0.28.0 | No. This is the documented edge behavior. | A pre-release binary is an edge build (`init.go:139`, the channel is stamped in the build). The edge channel resolves the newest release, including pre-releases (`install.sh:170`). The release CI cuts `v0.29.0-pre0` on the v0.28.0 commit (`docs/releasing.md:29`). The two binaries contain the same code. | No code fix. If pre-release users must land on stable, `update` needs a way to change the channel. That is a new decision (question 2). |
| 2 | Compass binds to loopback on a remote VM | Yes, in the skill text of this repo. The bind default belongs to the Compass repo. | `skills/compass/SKILL.md:29` never passes `--host`. Line 32 gives `http://127.0.0.1:4553/` with no remote case. | Keep the loopback bind, because Compass has no authentication. When `SSH_CONNECTION` or `SSH_TTY` is set, tell the user to run `ssh -L 4553:127.0.0.1:4553 <host>`. A different default or authentication is a change in `spacedock-dev/compass`. |
| 3 | The closing doctor of the installer prints MISSING before the PATH export | Yes | `install.sh:392` runs `"$INSTALL_DIR/spacedock" install --host claude` with the PATH of the installer. The doctor in that command looks up `subspace` and `compass` on that PATH (`versions.go:152,175`). | Run it as `PATH="$INSTALL_DIR:$PATH" "$INSTALL_DIR/spacedock" install --host claude`. This fix goes live from `main` without a release. |
| 4 | Commission omits `## Acceptance criteria`, so `--ac-scan` fails | Yes, for the seed entities | `skills/commission/SKILL.md:505-521`: the seed body has no AC section. The task template has one (line 479). `--ac-scan` exits non-zero by design (`internal/status/gate_extract.go:366`). | Add the `## Acceptance criteria` / AC-1 stub to the seed body. Also check the `experiment` and `refinement` templates, which have no AC heading. |
| 5 | `pair --confirm` does not answer a prompt that waits | Not in this repo: `spacedock-dev/subspace` (`subspace-relay-client pair`, also `rr`) | This repo only installs it and documents it (`docs/site/review/reference.md:63`). | Subspace. |
| 6 | A stopped `pair` leaves a state file that breaks the next one | Not in this repo: `spacedock-dev/subspace` | No pair state code is in this repo. | Subspace. Remove the file on a signal and at exit. Treat a file whose PID has stopped as stale. |
| 7 | An old local Subspace plugin hid `remote-review`, and the error named the incorrect cause | Partly. The error text is in this repo. Subspace does not load the skill. | `skills/review/SKILL.md:184` checks only the relay binary. When the skill is missing, lines 207-208 use again the remedy for a binary failure ("name the program and its version from that line"). Thus the agent names the relay binary, which is current. `review_mode.go` only stores the choice. | Give lines 207-208 a separate message: the relay is current, but the Subspace plugin is missing or old, so update or remove the local copy, then run `/reload-plugins`. |
| 8 | Installs are possibly not counted on the Homebrew path | Yes. The fix is not in the release. | Only `install.sh` makes the counted lookup. The goreleaser casks in `spacedock-dev/homebrew-tap` download directly from GitHub release assets (`.goreleaser.yaml:111-139`). The daily check of #884 sends `spacedock/<ver>`. It counts active binaries, not installs. | Change the cask URL in `.goreleaser.yaml` or `spacedock-dev/homebrew-tap`, or count the daily check as the signal for Homebrew. This is a decision about the analytics. It does not have to be in 0.28.1. |

## Proposed content of 0.28.1

Use the procedure in `docs/releasing.md` "Cutting a Patch Release". 0.28 is the newest stable line, so a patch on it is permitted.

1. Branch `release/0.28.1` from `v0.28.0` in `.worktrees/release-0.28.1`.
2. Cherry-pick, in first-parent order, all non-stamp commits from `v0.28.0..origin/main`: #875, `ad22447e1`, `eff994eff`, #881, #884, and the test-only #878 and #879. Do not take the three `release: bump/stamp` commits. #881 does not apply cleanly without #875 (conflict in `internal/release/install_extras_test.go`). Thus keep the order.
3. Add the new fix for defect 3, and the fixes for defects 4 and 7 if CL approves them. Each fix has a focused test or a skill smoke test. Land each fix on `main` first by a PR, then cherry-pick it.
4. Stamp 0.28.1 (`go run ./cmd/spacedock-release stamp-version 0.28.1 ...`). Run `go test ./...`, `go test ./... -race` and `gofmt`.
5. CL gates: push the branch, dispatch Runtime Live E2E on the SHA, make the annotated tag, and push it.

Items left out: defect 1 (working as documented), defect 2 (CL: skip), defects 5 and 6 (Subspace), and the analytics for defect 8.

Effect on edge users: a patch tag does not cut a pre-release (`edge-advance` skips). Users on `0.29.0-pre0` get these fixes only from a new `0.29.0-pre1`.

## Not verified yet

- I simulated only the order of the cherry-picks. The sandbox refused `git commit-tree`, so I did not run the full sequence. I did not run the test suite on the result.
- The `experiment` and `refinement` commission templates (defect 4) need one more read.
