# Validation gate: gate prepare resolves operator-supplied artifact paths without doubling

## Recommendation

Approve validation. This accepts commit `a20385e8e` for delivery in the 0.27.x fix stack.

## What you need to know

- All four acceptance criteria passed with behavioral evidence.
- Artifact and Reference paths now resolve consistently in absolute, launch-relative, and state-relative forms without creating duplicate bindings or rooms.
- Ambiguous, absent, non-ENOENT, symlink, non-regular, unreadable, and foreign-repository paths fail before mutation with actionable flag and path evidence.
- Omitted launch context preserves existing wrong-root refusal, and `gitsource.Inspect` remains the sole repository-safety guard.
- Focused tests, `go test ./...`, `go test ./... -race`, and formatting checks passed. A detached adversarial audit found no material findings.

## Estimate check

| Surface | Approved estimate | Actual | Variance | Result |
| --- | ---: | ---: | ---: | --- |
| Net LOC | +151, tolerance +110 to +215 | +183 | +32 (+21.2%) | Within estimate |
| Files | 6, maximum 7 | 6 | 0 | Within estimate |

Gross change is +222/−39 lines. The First Officer skill remains within its byte cap at 7,697 bytes.

## Decision effect

Approval closes validation and makes this ticket eligible for the next stacked PR above #767. Rejection returns it to implementation with the stated finding.
