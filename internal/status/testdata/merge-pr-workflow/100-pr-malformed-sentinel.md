---
id: "100"
title: Entity carrying a malformed merge sentinel and no mod-block
status: implementation
score: "0.50"
source: roadmap
pr: pr-merge:abc
---
# Entity carrying a malformed merge sentinel and no mod-block

A merge sentinel whose suffix is NOT a positive integer (`pr: pr-merge:abc`),
with an EMPTY mod-block. The prefix matches the pr-merge sentinel, but the suffix
is garbage — no PR number `abc` could name. `merge guard` must NOT treat this as a
landed merge: it must NOT finalize or archive. Only a well-formed sentinel
(`pr-merge:{positive-int}` / `local-merge:{sha}`) finalizes. This pins the
fail-open hole where a bare prefix match drove a full finalize+archive.
