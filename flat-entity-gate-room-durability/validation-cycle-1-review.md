# Validation cycle 1: malformed canonical ref bypass

## Rejected snapshot

Candidate `add2d0632` implemented canonical `@review/...` refs and direct flat
advisory rounds at **+249 net LOC across 13 files**, within the approved tolerance.

## Material finding

After an ordinary flat gate prepare, replacing the persisted room ref with
`@review/a/../b` made direct canonical resolution fail. Gate close nevertheless
returned success and wrote Resolution/application bytes because prepared-room
classification converted the resolver error into “not prepared” and retained-authority
validation skipped the binding.

This violates AC-5: malformed reserved refs must fail without spending authority or
changing entity bytes.

## Authorized disposition

Fix the task-owned bypass narrowly: preserve canonical resolver errors through
prepared-room classification and gate-close validation, then add a byte-clean
observable regression. Decline the unrelated pre-existing repository formatting
polish.

## Correction evidence

Correction `e4cd8afec` propagates the malformed-ref error and adds
`TestMalformedCanonicalRoomRefCannotSpendGateAuthority`. Focused, isolated slow-package,
full, and race suites pass. Final surface is **+250 net LOC across 13 files**, exactly
the approved upper LOC bound and within the file tolerance.

## Advisory decision

Record this validation cycle as revise if the rejected snapshot, material finding,
authorized disposition, and correction evidence are accurately represented.
