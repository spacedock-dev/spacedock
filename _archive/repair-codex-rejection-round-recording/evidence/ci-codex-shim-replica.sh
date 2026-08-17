#!/usr/bin/env bash
# Local replica of the CI codex shim (.github/workflows/runtime-live-e2e.yml).
# The repair loop must measure the same model+effort the stack PR's full matrix
# will judge AC-2 against; bare `codex` defaults are a different configuration.
set -euo pipefail
real_codex="${SPACEDOCK_CODEX_REAL_BIN:?SPACEDOCK_CODEX_REAL_BIN is not set}"
case "${1:-}" in
  exec)
    exec "$real_codex" exec --model gpt-5.6-luna -c 'model_reasoning_effort="max"' "${@:2}"
    ;;
  --ask-for-approval)
    if [ "${2:-}" = "on-request" ] && [ "${3:-}" = "exec" ]; then
      exec "$real_codex" "${@:1:3}" --model gpt-5.6-luna -c 'model_reasoning_effort="max"' "${@:4}"
    fi
    ;;
  --dangerously-bypass-approvals-and-sandbox)
    if [ "${2:-}" = "exec" ]; then
      exec "$real_codex" "${@:1:2}" --model gpt-5.6-luna -c 'model_reasoning_effort="max"' "${@:3}"
    fi
    ;;
esac
exec "$real_codex" "$@"
