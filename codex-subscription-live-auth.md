---
title: Use Codex subscription auth in live CI
status: backlog
source: captain-request
started:
completed:
verdict:
worktree:
id: hrercm3ff4ww94rnqhqbqkyp
---
# Use Codex subscription auth in live CI

Configure the Codex and Pi live CI lanes to use Codex subscription OAuth instead of an API key.

The implementation must preserve isolated homes and must not print OAuth secrets.

## Acceptance criteria

**AC-1** Codex live CI authenticates from a dedicated environment secret containing the Codex `~/.codex/auth.json` payload, and the Codex lane verifies login status without `OPENAI_API_KEY`.

**AC-2** Pi live CI authenticates from a separate dedicated environment secret containing Pi's `openai-codex` OAuth record, selects an exact provider-qualified Codex model, and preserves isolated Pi homes.

**AC-3** The two secret formats and the GitHub Environment setup are documented, including rotation and refresh-token handling.

**AC-4** Workflow and authentication tests prove that the selected secrets, provider/model, isolation, and no-secret-output boundaries remain correct.
