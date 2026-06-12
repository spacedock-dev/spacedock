# Voice & tone

This is the writing-style guide for Spacedock's public documentation. It is grounded in two real voice signals, not invented:

- **The root `README.md`**: Spacedock's actual public voice. Direct, anti-hype, claim-first ("Spacedock is a multi-agent orchestrator where nothing ships without a decision"), concrete over abstract, second person addressed to the reader.
- **The `comm-officer` mod** (`docs/dev/_mods/comm-officer.md`): the project's prose discipline. It applies the `elements-of-style:writing-clearly-and-concisely` skill (Strunk) and is **light-touch by default**: "Preserve the caller's voice, rhythm, and technical vocabulary. Cut empty words, tighten sentences, fix clear grammar errors." Critically, it **defers to a project voice guide when one exists**. This page is that guide, so the comm-officer and any doc-writer apply it.

## Voice

- **Precise, honest, technical.** State what is true and what a thing does. Spacedock's own pitch leads with the claim, not the adjective.
- **No marketing or hype adjectives.** Avoid "powerful", "seamless", "revolutionary", "effortless", "blazing-fast", "game-changing". If a sentence still carries its meaning with the adjective removed, remove it. The product ethos, evidence over assertion, is the writing ethos.
- **Concrete over abstract.** Name the command, the file, the outcome. Prefer "`spacedock status --next` lists the items ready to dispatch" over "Spacedock surfaces actionable work."
- **Claim-first, then support.** Lead a section with the load-bearing sentence; follow with the detail. Mirrors the README's "what's different" bullets (bold claim, then the mechanism).

## Avoid the AI-writing tells

Generated prose has a recognizable texture: padded, impersonal, and the fastest way to make simple software feel like a manual. Cut these on sight:

- **Em-dashes.** Do not use `—`. Rewrite as a period, a comma, a colon, or parentheses. A sentence that needs an em-dash usually wants to be two sentences.
- **The "not just X, but Y" frame** and its cousins ("it's not only…", "more than just…"). State what the thing is. Drop the contrast scaffolding.
- **Rule-of-three padding.** Three parallel adjectives or clauses where one carries the meaning ("clear, simple, and easy to follow"). Keep the load-bearing one.
- **Hollow transitions and hedges.** "That said," "It's worth noting that," "In order to," "It's important to understand." Delete them; start with the content.
- **Empty intensifiers.** "very", "really", "quite", "actually", "simply", "just" (when it adds nothing), "leverage", "utilize". Prefer the plain verb.
- **Throat-clearing openers.** A page or section that opens by restating its own title or announcing what it will cover ("This page covers…", "In this section, we will…"). Open with the content.

Read it back and remove every word that survives removal. If a sentence sounds like it is performing thoroughness rather than saying something, rewrite it.

## Tone and register per audience

- **New-user pages (Get started): welcoming and encouraging, still precise.** Assume no prior Spacedock knowledge; define a term on first use; show the command and the output to expect. Confidence-building, never breezy.
- **Operator pages (Concepts, Running workflows): direct and operational.** The reader is doing the work; tell them the steps and the decision points plainly.
- **Contributor pages (Contributing, Reference): exact and unembellished.** Precision outranks warmth. Name the contract, the test, the failure mode.
- **Person and tense.** Second person ("you run", "you approve") and present tense for how-to and instructions, the README's register. Imperative for steps ("Run `spacedock doctor`."). Describe the system in the present tense ("the first officer dispatches an ensign"), not the future.

## Canonical terminology and capitalization

Pinned from how the README and skills actually use these forms, not an imposed guess. Use them consistently.

| Term | Form | Notes (grounded in real usage) |
|------|------|--------------------------------|
| Spacedock | `Spacedock` | The product. Always capitalized. `spacedock` (lowercase, code font) only as the literal command/binary. |
| Captain | `Captain` (role) / `the captain` (prose) | README roles table capitalizes the role name; running prose uses lowercase "the captain" (skills use `{captain}`). The human operator. |
| First Officer | `First Officer` (role) / `the first officer` (prose) | Same pattern: Title Case in the roles table / when naming the role; lowercase in running prose ("the first officer reads the README"). The orchestrator agent. |
| Ensign | `Ensign` (role) / `the ensign`, `ensigns` (prose) | Same pattern. The worker agent that moves one item through one stage. |
| workflow | `workflow` | Common noun, lowercase. A directory of markdown entities + a README. |
| entity | `entity` | Common noun, lowercase. One work item (a markdown file or folder). The README also says "work item"; prefer "entity" in docs, gloss it as "work item" on first use for new users. |
| stage | `stage` | Common noun, lowercase. backlog → ideation → implementation → validation → done. |
| gate | `gate` | Common noun, lowercase. The decision point at the end of a stage. |
| sprint | `sprint` | Common noun, lowercase. A grouped set of entities driven to a deliverable. |
| worktree, mod, safehouse | lowercase | Common nouns. `safehouse` is also the sandbox profile filename `.safehouse`. |

Rule of thumb: capitalize a **role** when you name it as a role (the roles table, a definition); use lowercase for the same word in ordinary running prose; never capitalize the common-noun primitives (workflow, entity, stage, gate, sprint).

## Formatting conventions

- **Commands and code.** Inline code font for commands, flags, filenames, and identifiers: `spacedock claude`, `--strict`, `mkdocs.yml`. Multi-line commands and config in fenced code blocks with a language tag (` ```bash `, ` ```yaml `).
- **Show expected output.** For a command a reader runs, name the result, as `install.md` does ("Prints the installed version, e.g. `spacedock 0.20.0`").
- **Headings.** Sentence case ("Get started", "Your first workflow"), not Title Case. One `#` h1 per page (the page title); section headings start at `##`.
- **Links.** Descriptive link text, never "click here" / "this link". Internal links are relative so `mkdocs build --strict` can resolve them.
- **Lists and emphasis.** Bullets for parallel items; bold for the load-bearing claim of a bullet (the README pattern). Use emphasis sparingly; if everything is bold, nothing is.

## How this guide is applied

The `comm-officer` mod loads a project voice guide on first use and defers to it over plain Strunk. Once this page ships, it is that guide: doc contributions and comm-officer polish follow these rules. When this guide is silent on a question, fall back to Strunk via `elements-of-style:writing-clearly-and-concisely`.
