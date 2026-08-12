---
title: Use Codex subscription auth in live CI
status: ideation
source: captain-request
started:
completed:
verdict:
worktree:
id: hrercm3ff4ww94rnqhqbqkyp
gates:
    version: 1
    records:
        - id: gate:hrercm3ff4ww94rnqhqbqkyp:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:hrercm3ff4ww94rnqhqbqkyp-backlog-1
              briefing:
                id: briefing:hrercm3ff4ww94rnqhqbqkyp:backlog:attempt-1:revision-1
                digest: sha256:4725aa1d24ef0107e1d244fd1e666415c73cded25f72011060f8069526344c87
                request-digest: sha256:2f934a126b11db2e1a45db3445cf1ed198de1992864f2c63f7b69a8c2b39fa7a
                room-ref: ./codex-subscription-live-auth/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:hrercm3ff4ww94rnqhqbqkyp:backlog:1
                briefing: briefing:hrercm3ff4ww94rnqhqbqkyp:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-12T14:24:42.50027Z"
                decision: approve
                reason: Captain approved the seed. A backlog stage report is not required for this gate.
              application:
                target-stage: ideation
                state: consumed
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
