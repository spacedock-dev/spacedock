# Second independent ideation review: issue 813

Recommendation: **approve ideation** at state commit `09aa784d4e96930d63c5ed949c13059c2f14e870`.
The revision closes R1 in the design, proof plan, and proposed documentation. Candidate behavior remains the implementation and validation owners' obligation.

## Review basis

- Current committed [design](index.md), compared with the first review's unchanged finding package and the first officer's authorized correction boundary.
- Accepted wrapper and native Codex observations in [the first review](ideation-review.md), state commit `68bc1d8574158c35919f69f41272e70b805d3eed`.
- The first officer's source comparison found no product diff between `661c607ae34620caa9579f2f65b53e4a5d242344` and `15e0baea39a6d47a802a01f45cdcfd8063c6edfd` in the examined packages.

This review retained that evidence without repeating the live observations or unrelated full baselines.
A byte comparison found the five approved seed value statements, frontmatter, and first review artifact unchanged.

## R1 correction assessment

| Required boundary | Committed revision and assessment |
|---|---|
| Claude selection | The table names `CLAUDE_CONFIG_DIR` and `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS`. Both use the selected-host allowance. Complete. |
| Codex selection | The table names `CODEX_HOME`. The native settings observation distinguishes selected `on-request` from fallback `never`. Complete. |
| Pi selection | The table names `PI_CODING_AGENT_DIR`. The candidate plan requires native state from the selected settings. Complete. |
| Presence and exact values | Presence selects forwarding, including empty values and explicit Claude teams `0`. Absent values receive neither automatic allowances nor synthesized values. Complete. |
| Host and mode limits | Only the selected host's present canonical names receive automatic allowances in wrapped passthrough. Ordinary mode retains its existing behavior. Complete. |
| Caller setup | The candidate runs without caller-supplied environment allowances. The design adds no caller setup obligation. Complete. |
| Filesystem and environment limits | Existing Safehouse rules still control path access. The correction adds neither filesystem grants nor general environment forwarding. Complete. |
| Existing mechanism | The revision composes the existing named allowance. A new transport, supervisor, parser, or enforcement mechanism is unnecessary. Complete. |

R1 retains its **Material** classification and task ownership under `value-ac[AC-2]`.
The supported wrapped-provider workflow previously selected fallback permissions despite an explicit caller settings directory.
The first review's native observation established `on-request` becoming `never` without the named allowance.
The design owner closed that omission within the authorized boundary. Implementation and validation retain ownership of the candidate proof.

## Value, proof, and surface

| Value criterion | Proof coverage after revision |
|---|---|
| AC-1 | Existing literal-argv fixtures and candidate protocol startup retain clean stdout and native protocol responses. |
| AC-2 | The environment matrix covers every selected host, wrapped/unwrapped passthrough, and ordinary mode. Values include absent, empty, `0`, spaces, and `=`. |
| AC-2, effective settings | Real-wrapper candidate observations retain all canonical sentinels and selected native settings. Codex uses independently different fallback and selected model/approval pairs. |
| AC-3 | Existing sandbox precedence, binary propagation, scrubbing-wrapper, and Pi stable-path owners remain intact. The new allowance does not replace those requirements. |
| AC-4 | Common compatibility-flag acceptance and ordinary-launch regression owners remain intact. The recorded unsuccessful full-suite baseline still requires implementation and validation evidence. |
| AC-5 | The concrete documentation names all four variables, presence behavior, empty values, explicit `0`, filesystem limits, and ordinary-mode stability. |

Expected names and values remain literals independent of production lists.
The plan identifies omitted allowances, empty-value filtering, synthesized defaults, wrong-host allowances, and ordinary-mode changes as distinct falsifying edits.
Byte-exact inner argv alone cannot satisfy the effective-settings proof.
The existing startup observations own that additional evidence, with no new standing check or live journey.
The inherited detached audit and applicable live CI obligations remain in force.

The surface remains +400 net LOC across 13 product files, from 470 insertions and 70 deletions.
The tolerance remains +280 through +520 net LOC, with at most 15 product files.
The listed helper, host functions, environment fixtures, help, and documentation already cover R1. Implementation must count its actual diff.
The approved value statements and ordinary-behavior boundary remain unchanged.

## Remaining findings and recommendation

No new finding requires Material, Deferred risk, Polish, or Needs decision triage.
R1 is closed for ideation. Its effective-behavior proof remains pending with implementation and validation.
The recommendation is approve ideation with the original value criteria, surface, and tolerance.
