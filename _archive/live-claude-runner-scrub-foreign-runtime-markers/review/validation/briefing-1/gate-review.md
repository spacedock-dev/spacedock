# Validation gate: isolate live runtime host markers

## Capability delivered

Each supported live host child receives its own runtime detector and credentials
exactly once while foreign runtime-family markers are removed. Production still
rejects genuinely mixed host markers. The recorded Claude gate journey now proves
one flag-free, ambiguity-free, successful dispatch build followed by its committed
successor effect.

## Exact candidate

Candidate `6ea4cdc80ddc1be9b1bf7233211abb127bf56c90` is based on
`4ff98d8cd97ebcf17b6a583070ce69234e24fc87`. Its aggregate diff changes six
`internal/ensigncycle` test/harness files and no production, skill, workflow, or
provider file.

## Validation evidence

- Exact-slice Claude, Codex, and Pi child-environment tests fail on any retained
  foreign marker, duplicate target marker, or altered credential, PATH, HOME, or
  target-host state.
- Both unchanged production ambiguity controls still reject Codex+Claude and
  Codex+Pi mixed markers.
- The recorded-gate observer requires exactly one dispatch-build attempt and one
  byte-matching `exit=0`; a planted `exit=1` log retains the successor effect and is
  rejected.
- Focused default and live-tag tests, `go test ./...`, and `go test ./... -race`
  passed. Roborev job 2342 returned no findings.
- The same independent validator reproduced all four acceptance criteria and
  recommends PASSED.

## Deferred risk

The Pi builder is not selected by every future live CI shape. That becomes material
only if such CI coverage is promised or a supported Pi journey exhibits marker
contamination; the current explicit live-tag proof passes.

## Recommendation

Approve the exact candidate for the merge and terminalization path.

## Decision

Approve to merge candidate `6ea4cdc8`; revise to return a material finding to
implementation; or hold before merge.
