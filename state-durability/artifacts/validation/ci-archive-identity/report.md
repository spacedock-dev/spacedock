# CI archive fixture correction

Candidate: `6a9e4c839d2d4285f1078c35a5bfef4a2ab6ecd7`, above `c8c6d112553dc094a6da393cec97a1537d312320`. The correction changes only `skills/integration/launcher_smoke_test.go`: 5 additions / 1 deletion. The full durability diff above naming `70165e9b5` remains 22 files, 635 additions / 44 deletions, +591 net.

## Finding and disposition

- Released user and normal workflow: runtime CI's offline launcher smoke exercises explicit durable retirement of its folder-form state entity.
- Observable harm: runs `35056628848` and `35056631117` at stack tip `7365c74ce` failed `TestLauncherListSetArchive`, archive exit 1 / commit exit 128. The fixture file was identical at that tip and durability `c8c6d1125`.
- Authority: `value-ac[AC-2]` real-CLI retirement evidence must exercise a committed archive move with a self-contained fixture.
- Trigger evidence: `gitInitFixture` supplies identity only to its seed command. The later archive commit lacks persistent repository identity. The [red log](identity-required.log) reproduces the same wrapper failure; [Git Trace2](required-trace.jsonl) ties `git commit -q -m "archive pilot-entity (retirement)"` to `no email was given and auto-detection is disabled`, exit 128.
- Proposed classification: Material evidence defect, owned by this task's newly created state fixture. FO separately authorized the one-file fix. Product behavior, shared helpers, gate authority and approval remain unchanged.

## Completed work

- DONE: Reproduce and identify the precise cause of the CI archive fixture failure without candidate edits.
  Clearing user Git configuration and author/committer variables alone still passed on macOS because Git could auto-detect identity. Adding isolated `user.useConfigOnly=true` made the unchanged fixture fail at its real archive commit. The exact underlying diagnostic comes from the controlled local trace; the retained CI wrapper does not expose Git's inner stderr.
- DONE: Propose the smallest task-owned correction with falsifiable local proof and unchanged product semantics.
  The state fixture now uses existing `testgit.InitRepo`, then the same add/seed commit. This persists identity for the real launcher commit. The exact identity-required test changed from red to [green](identity-required-green.log), with author/committer/config environment overrides scrubbed. Removing persistent identity makes that same test fail again.

## Checks on the committed correction

- Identity-required `TestLauncherListSetArchive`, `-count=1`: passed, 1.375s. The red trace and green trace are retained beside this report.
- Ordinary `go test ./skills/integration -count=1`: passed, 3.490s; [log](integration.log).
- Required `go test ./...`: exit 1 solely at `TestCodexResolveManifestAgainstInstalledHost`, `codex_resolve_test.go:44`; all other packages passed. Integration passed in 5.805s; [normal log](normal.log).
- Required `go test ./... -race`: exit 1 solely at the same resolver assertion; all other packages passed. Integration passed in 13.086s; no data-race report appeared; [race log](race.log).
- Exact retained baseline: the host reports `spacedock@spacedock` absent while resolution returns the installed `spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`. This existing environment-dependent failure remains out of scope.
- Required `gofmt -w ./cmd ./internal` completed, and the changed integration file was formatted separately. Unrelated preexisting release-fixture whitespace was preserved.

## Reproduce the identity-required smoke

Run from any directory, supplying the candidate checkout as the argument. This changes only a temporary Git configuration and test fixtures, not global configuration.

```bash
python3 - /path/to/candidate-checkout <<'PYCODE'
import os, pathlib, subprocess, sys, tempfile
with tempfile.TemporaryDirectory() as temp:
    cfg = pathlib.Path(temp) / "gitconfig"
    cfg.write_text("[user]\n\tuseConfigOnly = true\n")
    env = os.environ.copy()
    for key in list(env):
        if key.startswith(("GIT_AUTHOR_", "GIT_COMMITTER_", "GIT_CONFIG_")) or key == "EMAIL":
            env.pop(key, None)
    env["GIT_CONFIG_GLOBAL"] = str(cfg)
    env["GIT_CONFIG_NOSYSTEM"] = "1"
    subprocess.run(["go", "test", "./skills/integration", "-run",
                    "^TestLauncherListSetArchive$", "-count=1", "-v"],
                   cwd=sys.argv[1], env=env, check=True)
PYCODE
```

The correction is ready for the retained validator's recheck. No approved report, entity frontmatter, gate, push, rebase, or CI run changed during this correction.
