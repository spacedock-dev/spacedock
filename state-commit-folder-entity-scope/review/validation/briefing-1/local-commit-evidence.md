# Candidate commit evidence

- Base: `07811d27`
- Candidate: `d4d39ed616f10021d2737f5f919eb243ba62eae0`
- Parent: `030b9c3c77fe6608f00d66c67190303db9c369ab`
- Subject: `docs state commit scoped add behavior`
- Author date: `2026-07-22T23:06:52+08:00`
- Candidate worktree porcelain at package assembly: empty

Exact base-to-candidate name-status:

```text
M	docs/site/reference/command-reference.md
M	internal/cli/state_commit_test.go
M	internal/cli/state_sync.go
```

Exact base-to-candidate diff stat:

```text
 docs/site/reference/command-reference.md |   9 +
 internal/cli/state_commit_test.go        | 363 ++++++++++++++++++++++++++++++-
 internal/cli/state_sync.go               |  62 ++++--
 3 files changed, 413 insertions(+), 21 deletions(-)
```

The exact commands and outcomes for focused, full, race, formatting, cleanliness, and Roborev verification are retained unchanged in `validator-report.md`.

