# Approve PR #782 for merge

## Recommendation

Approve. The combined stack tip `284299e3b` includes cleanup commit `de61c1baf`, passes independent behavior validation, and stays within its approved surface estimate.

## What ships

The Pi front door pins release-shaped launchers to their own Spacedock package ref. It makes at most one repair attempt, rechecks the installed source, and refuses launch after an error or remaining mismatch. Development builds and user-managed package sources remain untouched.

## Proof that matters

- AC-1: unpinned and wrong-line starts each performed one pinned install, wrote the binary's own `@v0.27.2` source, and launched once.
- AC-4: install failure and ineffective repair each exited 1 after one attempt, launched zero times, preserved settings, and printed an actionable refusal.
- AC-5: `file:`, dirty development, `--plugin-dir`, and `SPACEDOCK_REPO_ROOT` cases made zero installs and preserved settings bytes.
- Focused derivation, classifier, and repair tables passed. Isolated `go test ./...` and `go test ./... -race` passed.
- The earlier false-green boundary was exercised directly: the independent reviewer reproduced both failures before the correction and repeated the same scenarios green afterward.

## Scope and stack

The PR diff is +310/−1, net +309 across three files, against +295±80 net across 3±1 files. Cleanup `de61c1baf` and release base `af70297dd` are ancestors. No removed fake-seam test, dedicated live journey, or registry wiring returned.

## Residuals

Two accepted deferred risks remain: unsupported older Pi package managers might not honor `@ref`, and simultaneous Pi launches can race a single-user settings rewrite. Promote either if it becomes supported or is observed to leave persistent wrong-line state.
