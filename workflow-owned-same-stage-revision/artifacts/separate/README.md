---
commissioned-by: spacedock@0.28.0-pre3
id-style: slug
stages:
  states:
    - name: implementation
      initial: true
    - name: validation
      gate: true
      fresh: true
      feedback-to: implementation
    - name: done
      terminal: true
---
# Dev topology control
### `implementation`
Correct plan and commit.
### `validation`
Independent validator compares plan with frozen input and commits verdict before re-gating.
