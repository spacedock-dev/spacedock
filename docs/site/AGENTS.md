# Authoring directive for `docs/site/`

**Always apply this when creating or modifying any content under `docs/site/`.** It is the Spacedock adaptation of Recce's documentation-writing standard ([`writing-content/reference/doc.md`](https://github.com/DataRecce/recce-team/blob/main/recce-team/skills/writing-content/reference/doc.md)). It governs both the **shape** of a page (structure and simplicity) and its **voice** (terminology, register, and the tells to avoid). When a question falls outside it, fall back to the `elements-of-style:writing-clearly-and-concisely` (Strunk) skill.

## The first rule: simple beats complete

Spacedock is not complicated; wordy docs make it *feel* complicated. Every edit should make the reader's job smaller, not bigger.

- **Lead with the problem or the payoff, not the definition.** Open each page with why the reader is here and what they will be able to do, then name the mechanism.
  - ✗ "Spacedock is a multi-agent orchestrator."
  - ✓ "You have work that needs doing in stages, with a human sign-off before anything ships. Spacedock runs that for you."
- **Introduce the fewest terms possible, as late as possible.** Don't front-load a glossary. Define a term on first real use, gloss it once, then just use it. If a page introduces more than a handful of new terms, cut or defer some.
- **Lead with what the user sees and must know; keep the how-it-works light.** Name the visible behavior and the required concepts first. Internal mechanics (scheduling and reuse conditions, file/branch naming templates, parser internals, query plumbing) get at most a sentence or a link to the source. If a paragraph reads as protocol documentation, compress it or cut it.
- **Product anatomy is not a reader concept.** Do not name internal components (launcher, plugin, host) unless the reader must type or choose between them. The reader installs Spacedock and launches Spacedock; how it decomposes is not their problem.
- **Write from the reader's seat, never the maintainer's.** If a sentence's subject is the build, the script, or the parser, recast it around what the reader gets or does.
  - ✗ "The build emits a curated `llms.txt` index at the site root."
  - ✓ "Start from `llms.txt`, the curated index of these pages."
- **A well-named command needs no caption.** "Run `spacedock doctor`." is complete; explaining that it diagnoses problems repeats the name. Likewise, never pre-document what a tool prints interactively (the install script already prints its own `PATH` note).
- **A heading is a promise; the section delivers exactly that.** Codex setup is not "Install Spacedock"; a tab the reader chose ("no Homebrew") must not explain the thing they opted out of. Headings also take the reader's angle: name what the reader does or gets ("Turn the report into a workflow"), not what the product emits ("The commission offer").
- **Serve the typical reader; route edge audiences elsewhere.** A Get-started page carries only the path a typical user walks. Contributor and from-source material lives in the repo, not inline, not even as one sentence.
- **Document the durable value, not the output inventory.** Do not enumerate an artifact's current sections or fields; they drift with releases. Name what the reader gets out of it (the survey report is "the four things you learn", not a list of its sections).
- **A real example with sample output beats description.** One concrete command plus the output the reader will see teaches faster than paragraphs. Compose samples from the product's real templates, and never restate in prose what the sample already says ("Accept this design, or tell me what to change" makes a "nothing is generated until you accept" sentence dead weight).
- **Keep agent-facing surfaces off user pages.** Commands meant for the agents (`spacedock status`, `spacedock dispatch`) are not taught in Get started, even though their output is real.
- **The docs may defer to the agent.** Spacedock runs inside a coding agent; "ask the agent if anything is unclear" is a legitimate close. Pages cover what the reader needs before and between sessions, not every contingency within one.
- **Cut, don't pad.** If a sentence still carries its meaning with a clause removed, remove the clause. If a paragraph repeats the page above it, delete it and link instead.
- **Don't repeat content across pages.** One page owns each idea; others link to it. Duplicated explanation is the main thing that makes the docs feel long. An audience split ("new-user view" vs "operator view" of the same feature) is not a reason for a second page; it is the same topic.
- **Two levels of structure only:** section → page. No deeper nesting; no page that exists only to hold sub-pages.

## First impression: "simple," not "enterprise"

The first screen a reader sees should make them feel *"there is one idea here, and I already get it,"* not *"this is a complex system with a manual."* Front-loading vocabulary or role tables is what makes simple software read as enterprise software.

- **No glossary or role table in the first screen.** Introduce a term only when the next sentence needs it; defer the rest to the page that owns them and link. A roles or personas table near the top reads as enterprise structure. Link to it, do not lead with it.
- **Name the one idea, then say everything else is detail.** Give the reader permission to feel they understand the core before they read on. For Spacedock that idea is: work moves through stages, and nothing crosses a gate without a decision the reader owns.
- **The docs home connects forward, it does not re-pitch.** A reader who arrived from the landing page has already seen the problem and the promise. Do not restate the pain. State the one idea and point into the docs.

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
Optimize for lookup in seconds. Use **tables** for flags, fields, and options (Name · Type/required · Description · Default). Give copy-paste examples. Mark required vs optional; show defaults; note edge cases. When the title and a table carry everything, ship the title and the table: no framing intro, no closing caveat. If a "reference" page is really an explanation, it belongs in Concepts; if it's really a procedure, it belongs in a how-to.

## Link instead of repeat
- Cross-link liberally with **relative** internal links (so `mkdocs build --strict` resolves them). Concept → tutorial + reference; tutorial → concept + reference; reference → concept.
- Link the first mention of a load-bearing term (`workflow`, `gate`, `ensign`, `commission`) to the page that owns it, rather than re-explaining it.
- Descriptive link text, never "click here".
- **Link by payoff, not by contents.** A cross-link or Next entry names what the reader gains there ("read about the survey report to understand your usage pattern"), not what the page contains ("covers what survey reports").
- **External tools get a gloss and a link to their own site.** A dependency that is not Spacedock (`agentsview`, `safehouse`) is a real thing the reader may install: name it, gloss it in one clause, link it on first mention.

## Voice

This voice is grounded in two real signals, not invented: the root `README.md` (Spacedock's public voice: direct, anti-hype, claim-first, concrete over abstract, second person) and the `comm-officer` mod (`docs/dev/_mods/comm-officer.md`, the project's Strunk-based prose discipline, which defers to this directive when one exists).

- **Precise, honest, technical.** State what is true and what a thing does. Spacedock's own pitch leads with the claim, not the adjective.
- **No marketing or hype adjectives.** Avoid "powerful", "seamless", "revolutionary", "effortless", "blazing-fast", "game-changing". If a sentence still carries its meaning with the adjective removed, remove it. The product ethos, evidence over assertion, is the writing ethos.
- **Concrete over abstract.** Name the command, the file, the outcome. Prefer "`spacedock status --next` lists the items ready to dispatch" over "Spacedock surfaces actionable work."
- **Claim-first, then support.** Lead a section with the load-bearing sentence; follow with the detail. Mirrors the README's "what's different" bullets (bold claim, then the mechanism).

## Avoid the AI-writing tells

Generated prose has a recognizable texture: padded, impersonal, and the fastest way to make simple software feel like a manual. Cut these on sight:

- **Em-dashes.** Do not use `—`. Rewrite as a period, a comma, a colon, or parentheses. A sentence that needs an em-dash usually wants to be two sentences. The exemption is narrow: output reproduced verbatim (a captured transcript, a rendered template, an error string) must match what the tool prints. A sample you compose is prose: format it cleanly (newlines, colons) within the tool's real shape.
- **The "not just X, but Y" frame** and its cousins ("it's not only…", "more than just…"). State what the thing is. Drop the contrast scaffolding.
- **Rule-of-three padding.** Three parallel adjectives or clauses where one carries the meaning ("clear, simple, and easy to follow"). Keep the load-bearing one.
- **Hollow transitions and hedges.** "That said," "It's worth noting that," "In order to," "It's important to understand." Delete them; start with the content.
- **Empty intensifiers.** "very", "really", "quite", "actually", "simply", "just" (when it adds nothing), "leverage", "utilize". Prefer the plain verb.
- **Throat-clearing openers.** A page or section that opens by restating its own title or announcing what it will cover ("This page covers…", "In this section, we will…"). Open with the content.

Read it back and remove every word that survives removal. If a sentence sounds like it is performing thoroughness rather than saying something, rewrite it.

## Tone and register per audience

- **New-user pages (Get started): welcoming and encouraging, still precise.** Assume no prior Spacedock knowledge; define a term on first use; show the command and the output to expect. Confidence-building, never breezy.
- **Operator pages (Concepts, Running workflows): direct and operational.** The reader is doing the work; tell them the steps and the decision points plainly.
- **Reference pages: exact and unembellished.** Precision outranks warmth. Name the contract, the test, the failure mode.
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
- **Show expected output when the reader must check it.** Name what a command prints when the reader acts on the result (e.g. `spacedock status --next` prints the dispatchable set). A well-named command needs neither caption nor sample output.
- **Headings.** Sentence case ("Get started", "Your first workflow"), not Title Case. One `#` h1 per page (the page title); section headings start at `##`.
- **Links.** Descriptive link text, never "click here" / "this link". Internal links are relative so `mkdocs build --strict` can resolve them. GitHub links and install commands point at `main`, never a development branch.
- **Lists and emphasis.** Bullets for parallel items; bold for the load-bearing claim of a bullet (the README pattern). Use emphasis sparingly; if everything is bold, nothing is.

## Before you commit a docs change
- [ ] Page opens with the problem/payoff, not a definition.
- [ ] Introduces only the terms this page actually needs; each defined on first use.
- [ ] Re-read every touched page in reader-journey order: no term appears before the page that introduces it, including terms your own rewrite brought back.
- [ ] Matches its type's shape (concept / how-to / reference): no steps in concepts, no essays in reference.
- [ ] No content duplicated from another page; shared ideas are linked, not repeated.
- [ ] Internal links are relative and descriptive; first mentions of key terms link out.
- [ ] Ordered lists actually count up (watch for fenced blocks resetting the counter to 1).
- [ ] No AI-writing tells (em-dashes in prose, "not just X but Y", rule-of-three padding, hollow transitions, empty intensifiers, throat-clearing openers).
- [ ] Voice and terminology follow this directive.
- [ ] `mkdocs build --strict` passes.
- [ ] Read it back once and cut every sentence that survives removal.
