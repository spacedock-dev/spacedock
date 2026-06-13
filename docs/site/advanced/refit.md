# Refit a workflow

When you upgrade Spacedock, run `/spacedock:refit path/to/workflow` to bring the workflow's generated files up to date while leaving your edits in place. Nothing is auto-replaced: you see a diff and decide, file by file, and if a schema change affects your entities, refit proposes the migration and waits for your approval. Git is the safety net; ask the agent about anything else.
