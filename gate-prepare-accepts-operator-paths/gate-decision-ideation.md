# Design gate: resolve operator-supplied gate paths safely

## Recommendation

Approve implementation. The revised design fixes the reproduced doubled-path failure without guessing caller intent or weakening selected-source authority, and independent review found no remaining material issue.

## Selected design

- The CLI passes its explicit launch directory into the existing selected-source resolver.
- Direct package callers that omit launch context retain existing behavior; no launch candidate is invented.
- Relative spellings probe de-duplicated cleaned lexical paths under the launch and state roots.
- Exactly one present path proceeds through the existing `gitsource.Inspect` safety owner.
- Multiple present paths refuse before mutation and require an absolute spelling.
- Non-ENOENT `Lstat` failures refuse immediately with flag/path attribution and no fallback.

## Proof and safety

- AC-1 measures the stalled split-root prepare changing from exit 1/no binding to exit 0/one open binding.
- AC-2 proves absolute, launch-relative, and state-relative forms resolve to the same immutable Git source without broadening omitted-context callers.
- AC-3 proves ambiguity, absence, non-ENOENT errors, symlinks, non-regular files, unreadable files, and foreign repositories fail before mutation.
- AC-4 pins exact operator wording, unchanged grammar/help, and the fixed FO skill byte cap.

## Surface

Estimate: **+151 net LOC (+189/−38) across 6 files**. Tolerance: **+110 to +215 net LOC and at most 7 files**. The proposed FO skill is **7,697 bytes**, three below its 7,700-byte cap.

## Delivery proof owed

Implementation must run focused gates/CLI/gitsource/contractlint tests, formatting, full and race suites, and report exact insertions, deletions, net LOC, files, semantics, and cap bytes. Validation must reject any overrun or AC without behavioral/state evidence.

