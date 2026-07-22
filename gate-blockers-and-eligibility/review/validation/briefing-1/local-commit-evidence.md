# Candidate commit evidence

- Base: `fa240a76cd67fc0ea2552901824722ca8bfa1c73`
- Candidate: `c7612661cce857e90ae2073ac861f5b8b32b72c0`
- Candidate worktree porcelain at package assembly: empty

Commit sequence:

```text
b1513f808159ffcfc5be14ce1b9040da92996d99  Add one-use gate application lifecycle
6c41403a080483a48ad4486dea1750b3231d4bd9  Harden gate consumption eligibility
2cba20e0c5d79b250e796236197896580df6467d  Align gate eligibility with workflow discovery
c7612661cce857e90ae2073ac861f5b8b32b72c0  Strengthen gate eligibility evidence
```

Exact base-to-candidate paths:

```text
M  docs/site/reference/command-reference.md
M  docs/site/reference/frontmatter-contract.md
M  docs/specs/gate-resolution-frontmatter-contract.md
M  internal/cli/cli.go
M  internal/cli/gate_test.go
M  internal/cli/help.go
A  internal/gates/application.go
A  internal/gates/application_test.go
M  internal/gates/gates_test.go
M  internal/gates/model.go
M  internal/gates/operation.go
M  internal/status/discover.go
M  internal/status/gates_coexist_test.go
M  internal/status/handlers.go
```

Exact base-to-candidate stat: 14 files, 1,146 insertions, 52 deletions. The full per-file stat and exact command outcomes are preserved in `validator-report.md` and `validation-execution-evidence.md`.

