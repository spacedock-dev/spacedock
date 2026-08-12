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

Configure the Codex and Pi live CI lanes to prefer Codex subscription OAuth.

Keep `OPENAI_API_KEY` as a fallback during the migration. A lane must fail only when neither its OAuth secret nor `OPENAI_API_KEY` is available.

The implementation must preserve isolated homes and must not print OAuth secrets.

## Scope

Included:

- Codex CI authentication selection and isolated `CODEX_HOME` seeding.
- Pi CI authentication selection and isolated Pi-home seeding.
- GitHub Environment secret instructions for both OAuth file formats.
- Workflow and authentication tests.
- Runtime documentation for rotation and API-key fallback.

Excluded:

- Changes to the live model or thinking level.
- Changes to local developer authentication.
- Automatic persistence of refreshed OAuth tokens back to GitHub secrets.

## Acceptance criteria

**AC-1** Codex live CI prefers a dedicated environment secret containing the Codex `~/.codex/auth.json` payload, falls back to `OPENAI_API_KEY` when that secret is unavailable, and verifies login status without printing either credential.

**AC-2** Pi live CI prefers a separate dedicated environment secret containing Pi's `openai-codex` OAuth record, falls back to `OPENAI_API_KEY` when that secret is unavailable, selects an exact provider-qualified Codex model, and preserves isolated Pi homes.

**AC-3** The two secret formats and the GitHub Environment setup are documented, including rotation, refresh-token handling, and the API-key fallback.

**AC-4** Workflow and authentication tests prove OAuth preference, API-key fallback, provider/model selection, isolation, and no-secret-output boundaries.
