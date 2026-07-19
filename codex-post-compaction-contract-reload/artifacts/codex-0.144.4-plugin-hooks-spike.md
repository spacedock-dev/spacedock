# Codex 0.144.4 plugin-bundled hooks spike

Date: 2026-07-18. Implementation's "first spike, ahead of the fixture matrix":
does `.codex-plugin/plugin.json` declare a hooks key that Codex loads and trusts?

Answer: YES. Reproduced with a live Codex 0.144.4 CLI against an isolated
`CODEX_HOME`, a local marketplace, and an installed test plugin. Observability
was `RUST_LOG=codex_core_plugins=trace,codex_hooks=trace codex exec` (session
init loads plugins and fires `SessionStart` hooks before the model call; no auth
was available, so firing was observed via marker-file side effects and the
`hook: SessionStart Completed` trace line, not a model turn).

## What loads and fires

- `.codex-plugin/plugin.json` with **`"hooks": "./hooks.json"`** (a string path,
  the form documented in Codex's bundled `plugin-creator` `plugin-json-spec.md`)
  plus a `hooks.json` file at the plugin root **loads, trusts, and fires**.
  Proven reproducibly: a `SessionStart` hook running `touch <abs>` created the
  marker under `--dangerously-bypass-hook-trust`.
- `hooks.json` schema is the Claude-settings shape:
  `{"hooks":{"<Event>":[{"matcher"?:"...","hooks":[{"type":"command","command":"..."}]}]}}`.
  `PostCompact` with `"matcher":"manual|auto"` is accepted.
- Trust gate: the hook fires **only when trusted** (bypass flag, or a persisted
  trust decision made interactively in the TUI). Without trust it is inert and
  writes nothing — the harmless-absence behavior AC-4 requires.

## Execution model (constrains the command)

The hook `command` string is **whitespace-split into argv and exec'd directly —
no shell, no quote-stripping**. Verified: `touch /abs/path` fires; `sh -c "..."`,
`... | tee ...`, `echo ... > file`, `touch "/quoted path"`, and a `command`
array all fail (quotes stay literal, shell operators are inert). Therefore a
hook that must emit multi-word prose cannot be an inline command; it must be a
**bundled script**. That script MUST be referenced by a plugin-root-absolute path
via the `${PLUGIN_ROOT}` token Codex substitutes (see the next section), NOT by a
bare relative path. The curated OpenAI plugins (`replayio`, `figma`) ship
`"command": "./scripts/<name>.sh"`, but a bare relative command resolves against
the **session cwd** — the operator's project — not the plugin root; it fires only
when the session cwd happens to be the plugin directory (as it is not when the FO
operates on any other repo, the normal shipped usage). Proven below: from an
unrelated cwd `./hooks/x.sh` does NOT fire while `${PLUGIN_ROOT}/hooks/x.sh` does.

## Hook output contract

The command's stdout is parsed as the hook result JSON. For the reminder:
`{"systemMessage":"Spacedock: compaction completed. ..."}`. The 0.144.4 probe
(`codex-0.144.4-hook-probe.md`) already proved a `PostCompact` `systemMessage`
surfaces as one visible warning and does NOT enter model context — so the bundled
hook is a captain-facing UI cue, failure-open, with a manual-cue fallback.

## Important caveat for the deliverable

Codex's bundled `plugin-creator` skill and its `validate_plugin.py` linter
declare the manifest `hooks` field **unsupported** and strip it from generated
manifests. That linter is stale relative to the runtime: `codex_core_plugins`
parses `hooks` (accepting string / string-array / object / object-array), and
OpenAI's own curated plugins ship hooks. So the Spacedock plugin must ship the
`hooks` surface directly and must not depend on `validate_plugin.py` accepting
it. This is consistent with the failure-open design: if a future Codex removes
or gates plugin hooks, the binding degrades to the manual captain cue and nothing
is blocked.

## Command path resolution: `${PLUGIN_ROOT}` (M1 fix, verified codex-cli 0.144.6)

The validation review (round 1) rejected the first cut because the shipped
`hooks.json` referenced the script by the cwd-relative path
`./hooks/codex_post_compact_notice.sh`. Codex exec's a hook command from the
**session cwd** (the operator's project the FO is working on), so the relative
path fails to resolve whenever the FO operates on any repo other than the plugin's
own — the normal, shipped usage. The earlier offline test masked this by resolving
`./hooks/...` against the plugin repo root itself.

Reproduced against a live `codex-cli 0.144.6` CLI, marketplace-installed local
plugin, SessionStart hook fired via `codex exec --dangerously-bypass-hook-trust`
from an UNRELATED cwd (a temp dir that is not the plugin and has no `./hooks/`).
The bundled script recorded `argv0` and the plugin-root env vars Codex exposes:

| hook `command`                     | result from unrelated cwd |
|------------------------------------|---------------------------|
| `<abs>/hooks/mark.sh` (control)    | FIRED                     |
| `./hooks/mark.sh`                  | no-fire (the defect)      |
| `${PLUGIN_ROOT}/hooks/mark.sh`     | FIRED                     |
| `$PLUGIN_ROOT/hooks/mark.sh`       | no-fire (no shell — bare `$VAR` is not expanded) |
| `${CODEX_PLUGIN_ROOT}/hooks/mark.sh` | no-fire (not a substituted token; env var unset) |
| `${CLAUDE_PLUGIN_ROOT}/hooks/mark.sh` | FIRED                  |

Findings:

- Codex **template-substitutes the brace forms `${PLUGIN_ROOT}` and
  `${CLAUDE_PLUGIN_ROOT}`** in the command string before the no-shell exec,
  replacing them with the materialized plugin directory
  (`~/.codex/plugins/cache/<marketplace>/<plugin>/<version>`). The bare `$VAR`
  form is not expanded (there is no shell).
- The child process environment has `PLUGIN_ROOT` and `CLAUDE_PLUGIN_ROOT` set to
  the plugin dir; `CODEX_PLUGIN_ROOT` is unset.
- `${PLUGIN_ROOT}` is the neutral, documented token. The shipped hook now uses
  `${PLUGIN_ROOT}/hooks/codex_post_compact_notice.sh`, which resolves to an
  absolute path under the plugin install dir independent of session cwd.

The offline fixture (`codex_post_compact_hook_test.go`) was corrected to mirror
this: `resolveHookCommand` requires the `${PLUGIN_ROOT}/` prefix and substitutes
it — it FAILS on a bare cwd-relative command — and
`TestCodexPostCompactHookFiresFromUnrelatedCwdViaPluginRoot` runs the resolved
command from an unrelated temp-dir cwd (with `PLUGIN_ROOT` set) and asserts the
`systemMessage`, while proving the `./hooks/<script>` form does not resolve there.
