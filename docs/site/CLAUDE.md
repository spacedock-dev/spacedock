# Authoring directive for `docs/site/`

**Always apply this when creating or modifying any content under `docs/site/`.** It is the Spacedock adaptation of Recce's documentation-writing standard ([`writing-content/reference/doc.md`](https://github.com/DataRecce/recce-team/blob/main/recce-team/skills/writing-content/reference/doc.md)). It governs **structure and simplicity**; [`contributing/voice-and-tone.md`](contributing/voice-and-tone.md) governs **voice, terminology, and capitalization**. When they overlap, follow voice-and-tone for wording and this file for shape. When both are silent, fall back to the `elements-of-style:writing-clearly-and-concisely` (Strunk) skill.

## The first rule: simple beats complete

Spacedock is not complicated; wordy docs make it *feel* complicated. Every edit should make the reader's job smaller, not bigger.

- **Lead with the problem or the payoff, not the definition.** Open each page with why the reader is here and what they will be able to do, then name the mechanism.
  - ✗ "Spacedock is a multi-agent orchestrator."
  - ✓ "You have work that needs doing in stages, with a human sign-off before anything ships. Spacedock runs that for you."
- **Introduce the fewest terms possible, as late as possible.** Don't front-load a glossary. Define a term on first real use, gloss it once, then just use it. If a page introduces more than a handful of new terms, cut or defer some.
- **Cut, don't pad.** If a sentence still carries its meaning with a clause removed, remove the clause. If a paragraph repeats the page above it, delete it and link instead.
- **Don't repeat content across pages.** One page owns each idea; others link to it. Duplicated explanation is the main thing that makes the docs feel long.
- **Two levels of structure only:** section → page. No deeper nesting; no page that exists only to hold sub-pages.

## First impression: "simple," not "enterprise"

The first screen a reader sees should make them feel *"there is one idea here, and I already get it,"* not *"this is a complex system with a manual."* Front-loading vocabulary or role tables is what makes simple software read as enterprise software.

- **No glossary or role table in the first screen.** Introduce a term only when the next sentence needs it; defer the rest to the page that owns them and link. A roles or personas table near the top reads as enterprise structure. Link to it, do not lead with it.
- **Name the one idea, then say everything else is detail.** Give the reader permission to feel they understand the core before they read on. For Spacedock that idea is: work moves through stages, and nothing crosses a gate without a decision the reader owns.
- **The docs home connects forward, it does not re-pitch.** A reader who arrived from the landing page has already seen the problem and the promise. Do not restate the pain. State the one idea and point into the docs.

## Write it plain: cut the AI tells

The fastest way to make the docs feel padded is generated-prose texture. The voice rules that catch it (no em-dashes, no "not just X but Y", no rule-of-three padding, no hollow transitions or empty intensifiers) live in [`contributing/voice-and-tone.md`](contributing/voice-and-tone.md) under "Avoid the AI-writing tells." Apply them on every edit.

## Pick the page type, then follow its shape

Every page is one of three types. The nav already sorts roughly this way; match the type to the reader's question.

| Type | Reader asks | Spacedock sections | Must NOT contain |
|------|-------------|--------------------|------------------|
| **Concept** | "What is this / why does it matter?" | Concepts | Step-by-step instructions (link to the how-to instead) |
| **Tutorial / how-to** | "How do I do X?" | Get started, Running workflows | Long conceptual digressions (link to the concept instead) |
| **Reference** | "What are the options / the exact contract?" | Reference | Narrative prose, tutorial steps |

### Concept pages
Opening (1–2 sentences: what + why) → how it works (2–4 short paragraphs, ideally one diagram or worked snippet) → when to use / when not → **Related** links (tutorial + reference). Describe the system; don't instruct. If you're numbering steps, you're in the wrong page type.

### Tutorial / how-to pages
State the goal in one sentence. List prerequisites up front (never bury them mid-page). Numbered steps, **one action per step**, each showing the expected result (name the command *and* the output, e.g. ``spacedock status --next`` then what it prints). End with verification and "next steps" links. Push command-grammar and theory to the bottom or to a linked concept; the reader wants the outcome first.

### Reference pages
Optimize for lookup in seconds. Use **tables** for flags, fields, and options (Name · Type/required · Description · Default). Give copy-paste examples. Mark required vs optional; show defaults; note edge cases. If a "reference" page is really an explanation, it belongs in Concepts; if it's really a procedure, it belongs in a how-to.

## Link instead of repeat
- Cross-link liberally with **relative** internal links (so `mkdocs build --strict` resolves them). Concept → tutorial + reference; tutorial → concept + reference; reference → concept.
- Link the first mention of a load-bearing term (`workflow`, `gate`, `ensign`, `commission`) to the page that owns it, rather than re-explaining it.
- Descriptive link text, never "click here".

## Before you commit a docs change
- [ ] Page opens with the problem/payoff, not a definition.
- [ ] Introduces only the terms this page actually needs; each defined on first use.
- [ ] Matches its type's shape (concept / how-to / reference): no steps in concepts, no essays in reference.
- [ ] No content duplicated from another page; shared ideas are linked, not repeated.
- [ ] Internal links are relative and descriptive; first mentions of key terms link out.
- [ ] Ordered lists actually count up (watch for fenced blocks resetting the counter to 1).
- [ ] Voice and terminology follow [`contributing/voice-and-tone.md`](contributing/voice-and-tone.md).
- [ ] `mkdocs build --strict` passes.
- [ ] Read it back once and cut every sentence that survives removal.
