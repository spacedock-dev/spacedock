# Split-root state

Split-root state keeps workflow progress commits off your code branch. Workflow rules stay with the code; work records live in a separate Git checkout.

Ask your agent to configure this layout. On a fresh clone, restore the state checkout with:

```bash
spacedock state init
```

Conflicting edits to the same work item stop for your decision.

<a id="concurrent-writers"></a>
