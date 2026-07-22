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
    material = ($2 ~ /^supported:/ && $3 ~ /^present:/ && $4 ~ /^(value-ac|boundary):/ && $5 ~ /^supported:/)
    recorded_material = ($6 == "material")
    check = (($6 == "material" || $6 == "correct-but-disproportionate") && material == recorded_material) ? "accept" : "reject"

    if ($7 != check) {
        printf "FAIL %s: got %s, want %s\n", $1, check, $7 > "/dev/stderr"
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
