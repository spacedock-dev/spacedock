#!/bin/sh
# Deterministic fresh-worker stand-in for the Codex dispatch-boundary fixture.
# Every input path is supplied by the parent test; this probe never discovers
# files through ~/.codex, /tmp, the network, or the current working directory.
set -eu

prompt_path=$1
artifact_path=$2
plugin_root=$3
entity_path=$4
stage=$5
report_shape=$6

prompt=$(cat "$prompt_path")
prefix='$spacedock:ensign; then Read '
case "$prompt" in
  "$prefix"*) ;;
  *)
    printf '%s\n' 'bootstrap edge missing; refusing to read the ensign contract or dispatch artifact' >&2
    exit 41
    ;;
esac
printf '%s\n' 'bootstrap edge accepted'

for contract in \
  skills/ensign/SKILL.md \
  skills/ensign/references/ensign-shared-core.md \
  skills/ensign/references/codex-ensign-runtime.md
do
  contract_path=$plugin_root/$contract
  test -s "$contract_path"
  printf 'contract-read:%s\n' "$contract"
done

artifact=$(cat "$artifact_path")
case "$artifact" in
  *'### Completion checklist'*) ;;
  *)
    printf '%s\n' 'dispatch artifact is missing its completion checklist' >&2
    exit 42
    ;;
esac
printf '%s\n' 'artifact-read'

case "$report_shape" in
  done)
    {
      printf '\n## Stage Report: %s\n\n' "$stage"
      printf '%s\n' '- DONE: Codex fresh bootstrap exercised'
      printf '%s\n' '  The supplied child probe loaded the ensign contract before reading the dispatch artifact.'
      printf '\n%s\n\n' '### Summary'
      printf '%s\n' 'The deterministic child followed the fresh Codex bootstrap boundary.'
    } >> "$entity_path"
    ;;
  checkbox)
    {
      printf '\n## Stage Report\n\n'
      printf '%s\n' '- [x] Codex fresh bootstrap exercised'
      printf '%s\n' '  The generic checkbox shape is intentionally invalid.'
      printf '\n%s\n\n' '### Summary'
      printf '%s\n' 'The negative fixture deliberately emitted a non-compliant report.'
    } >> "$entity_path"
    ;;
  *)
    printf 'unknown report shape: %s\n' "$report_shape" >&2
    exit 43
    ;;
esac
printf 'report-written:%s\n' "$report_shape"
