---
title: Keep live-run artifacts out of the checkout
status: backlog
source: "FO write-scope violation during the paseo-spacedock plugin prototype: a live run with a relative SPACEDOCK_LIVE_ARTIFACT_DIR wrote .live-artifacts-<label> run directories into the repository root."
score:
started:
completed:
verdict:
worktree:
issue:
id: y7wntrqp0d31zxa65kkzbve0
---

`internal/ensigncycle` has three copies of the same artifact-root block:
`codexLiveArtifactDir`, `claudeLiveArtifactDir`, and `piLiveArtifactDir`. Each
returns `t.TempDir()` when `SPACEDOCK_LIVE_ARTIFACT_DIR` is empty. Each
otherwise returns `filepath.Join(root, name)`.

A relative `root` joins against the process working directory. That is the
repository root for the documented invocation. A run with
`SPACEDOCK_LIVE_ARTIFACT_DIR=.live-artifacts-<label>` then writes whole run
directories of agent transcripts, final messages, and topology dumps into the
checkout. Seven such directories accumulated in the repository root before they
were removed by hand.

The workflow README documents retention as `SPACEDOCK_LIVE_ARTIFACT_DIR=<dir>`
and does not say whether `<dir>` must be absolute. CI sets an absolute path
inside the workspace, so CI depends on an absolute value being honored as
written.

**Desired end state:** a live run never writes into the checkout unless the
operator gives an absolute path that deliberately points there.

A stashed candidate implementation exists on `main`
(`git stash list` entry "FO violation: live-artifact tmp default"): one shared
resolver, three delegating one-liners, and four focused tests. It was written by
the first officer outside write scope. Treat it as reference only; the stage
sequence below owns the change.
