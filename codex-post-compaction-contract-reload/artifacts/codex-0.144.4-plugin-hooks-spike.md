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
**bundled script referenced by a relative path** (`./scripts/x.sh`) that prints
the JSON. This is exactly the pattern the curated OpenAI plugins use — `replayio`
and `figma` ship `hooks.json` with `"command": "./scripts/<name>.sh"` and the
script in `./scripts/`. Their real-world use is the strongest evidence that
relative-to-plugin-root command resolution works in normal (TUI) operation.
(In headless `codex exec` the process cwd is the invoking repo, so a relative
command does not resolve to the plugin root there; the wiring still fires — the
`SessionStart` hook is logged — only the relative script is not found. Every
other link is proven: key parses, `hooks.json` parses, the engine runs the hook,
trust gates it.)

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
