#!/bin/sh
set -eu

fixture=${1:-"$(dirname "$0")/testdata/finding-triage-materiality.tsv"}

awk -F '\t' '
NR == 1 {
    if ($1 != "id" || NF != 7) {
        print "invalid fixture header" > "/dev/stderr"
        exit 2
    }
    next
}
{
    exact_columns = (NF == 7)
    cited = ($4 ~ /^(value-ac\[AC-[1-9][0-9]*\]|captain-ruling\[[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]\]|contract\[[^][[:space:]#]+#[^][[:space:]]+\]):[[:space:]]+[^[:space:]].*$/)
    no_boundary = ($4 ~ /^none:[[:space:]]+[^[:space:]].*$/)
    valid_boundary = (cited || no_boundary)
    material = ($2 ~ /^supported:/ && $3 ~ /^present:/ && cited && $5 ~ /^supported:/)
    recorded_material = ($6 == "material")
    known_class = ($6 == "material" || $6 == "correct-but-disproportionate")
    check = (exact_columns && valid_boundary && known_class && material == recorded_material) ? "accept" : "reject"
    expected = $7

    if (expected != check) {
        printf "FAIL %s: got %s, want %s\n", $1, check, expected > "/dev/stderr"
        failures++
    } else {
        printf "%s %s\n", toupper(check), $1
    }
}
END {
    if (NR == 1) {
        print "fixture has no cases" > "/dev/stderr"
        exit 2
    }
    exit failures > 0
}
' "$fixture"
