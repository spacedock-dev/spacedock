# Offline CI failure evidence

- Pull request: `spacedock-dev/spacedock#565`
- Workflow run: `30112268591`
- Job: `89544401796` (`offline`)
- Command: `go test ./...`
- Failing package: `internal/ensigncycle`
- Root error: `spacedock state commit: git commit failed: Author identity unknown`
- Affected controls: real CLI replay, terminal consume, AC-5 refusal matrix, and AC-7 resume matrix.
- Passed in the same run: build, macOS install, and Ubuntu install.

The clean runner had no global Git author identity. The test fixture must own the local identity needed by the product process it invokes.
