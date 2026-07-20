#!/bin/sh
# ABOUTME: Codex PostCompact hook that emits the Spacedock post-compaction reload
# ABOUTME: reminder as a systemMessage; failure-open captain-facing UI cue, no writes.
cat <<'JSON'
{"systemMessage":"Spacedock: compaction completed. Before continuing, ask the first officer to reread the authoritative `spacedock:first-officer` contract and reconcile durable workflow and live worker state."}
JSON
