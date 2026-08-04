# Backlog review: refresh the v1 pilot manifest after archive moves

Recommend approval. The current focused test fails exactly seven paths that normal workflow completion moved under `_archive/`. The repair updates those seven bindings and the independent archive-count assertion from 15 to 22.

The task changes two test-oracle files by an expected +8/-8. It does not change runtime code, state history, gate schema, or the 3d candidate. Its value gate is current-checkout `go test ./...` and `go test ./... -race` returning green without snapshot overrides.
