---
commissioned-by: spacedock@0.28.0-pre3
id-style: slug
stages:
  states:
    - name: plan
      initial: true
      gate: true
      feedback-to: plan
    - name: done
      terminal: true
---
# Triage-like plan workflow
### `plan`
Correct the local plan against frozen-input.txt, commit it, and present for captain review. No separate reviewer or correction-round publication is declared.
