# Sprint 0201 — post-flip release model (Model B)

**Goal:** make the stable channel serve a **pinned release tag** (not a moving branch HEAD),
so landing work on `main` is decoupled from shipping it — buying a **fast release cadence**
without "main HEAD must always be release-quality" as a standing tax, and keeping
experimental work on the edge channel.

**Deliverable:** a separate marketplace repo serving two pinned channels — **stable**
(`spacedock`, `ref: v0.X.Y`) and **edge** (`spacedock-edge`, `next` HEAD) — plus a
**stamp-then-tag** release ritual whose tagged commit's manifest matches its tag. NO contract
change (stays at 1).

> **Status: scope-locked 2026-06-13 (captain).** 0.20.1 ships as **decouple-first**: land
> the separate marketplace repo + pinned channels (`gpvg`, leading with its decoupling spike)
> and defer **stamp-then-tag** (`ezn`) to a follow-up once the repo split proves out. `w6`
> (marketplace-repo-decouple) is **folded into `gpvg`** — it is the same decouple. `tw`
> (next-independent-release-line) and `qp` (steady-state-release-runbook) **defer** to the
> follow-up. 0.20.1 also carries a **UX-cleanup cluster** tracked alongside: `gj`
> (startup/`--version` sandbox state + runtime detection), `zrc` (non-sandboxed auto-mode,
> couples to `gj`), `8p` (brew cask deps). `te`/`tes` (install-refresh) is **both** the
> migration dependency for `gpvg` and a standalone UX fix — sequence it early.

## The decision — Model B (captain 2026-06-09)

Three models were weighed; **Model B** was chosen, grounded by the plugin-distribution
research (`wcdgsgd88`, 9 agents, adversarially verified):

- **A — keep serving branch HEAD, `main` as trunk.** Rejected: a git source with no version
  pin means *"every new commit is treated as a new version"*
  ([plugins-reference](https://code.claude.com/docs/en/plugins-reference)) — i.e. every
  `main` commit is a release. No mature ecosystem's stable channel works this way.
- **B — serve a pinned tag/version for stable (chosen).** Matches what every mature
  ecosystem (npm `latest`, VS Code pre-release, Homebrew stable, Chrome channels) and
  Anthropic's own curated `claude-plugins-official` (sha-pins 125/125 entries) actually do.
- **C — promotion (`next` = trunk, promote to `main` at a cut).** Cleanest *flow*, but still
  serves a branch HEAD on whichever branch is "stable," inheriting A's HEAD-is-release
  coupling; also reintroduces the `source.ref` divergence unless the marketplace repo lands
  first (B's prerequisite).

**Why B fits the stated goals:** it gives the *direct-to-`main` ergonomics* the captain
wanted (cheap PRs land on `main`) **minus** the "every commit ships" tax — shipping is a
tag bump. Experimental stays isolated on the edge entry. Cadence is cheap: ship == tag + one
live-e2e.

## The settled design

- **Branches.** `main` is the **trunk** — normal/UX work lands via PR, gated by the offline
  suite + install smoke (NOT live-e2e). `next` = `main` + in-flight experimental, served at
  HEAD (the moving edge channel, like superpowers' `-dev`). Experimental **graduates** via a
  PR to `main` when ready; previewable on edge meanwhile.
- **Channels.** `spacedock` (stable) pinned `ref: v0.X.Y`; `spacedock-edge` on `next` HEAD.
  Two entries of one source, each resolving to a distinct version (the docs' channel
  convention + its "each channel must resolve to a different version" caveat).
- **Release ritual = stamp-then-tag.** Bump `plugin.json` → commit → one live-e2e (the
  existing `e2e-gate`) → **tag that commit** `v0.X.Y` → repoint stable's `ref`. The *same*
  tag drives the binary release (unchanged) and the stable plugin repoint — one tag, two
  consumers.
- **Separate marketplace repo** holds `marketplace.json`; plugin branches carry no manifest
  → kills the permanent `main`/`next` `source.ref` divergence and makes `next → main` a
  clean fast-forward.
- **Migration = seamless via binary update.** The next binary carries the new marketplace
  source; `install` / the front door auto-repoints existing users. Coupled to the
  install-refresh fix (`tes`).

## Proposed members

Membership is the query, not this table:

```bash
spacedock status --workflow-dir docs/dev --where sprint=0201-post-flip-release-model
```

| Entity | group | what it delivers |
|--------|-------|------------------|
| `marketplace-repo-and-pinned-channels` (`gpvg`) | release-model | separate marketplace repo + stable(tag)/edge(next) entries + install repoint + channel-stamp mapping; **decoupling spike first** |
| `stamp-then-tag-release-ritual` (`ezn`) | release-model | invert release ordering so the tagged commit's manifest matches its tag; rewrite `releasing.md`; delete the dead `next` marketplace `version` field |

**Dependency (not a member):** `tes` (install-refresh-and-upgrade-hint) — the seamless
migration rides its install-refresh fix; sequence `tes` before the migration step of `gpvg`.

## Definition of done

- The stable channel serves a tag-pinned plugin; a no-op `main` commit does **not** update
  stable while edge advances (the decoupling exercise passes).
- A release-machinery guard fails when a tagged commit's `plugin.json` version ≠ its tag.
- Plugin branches carry no marketplace manifest; `next → main` fast-forwards cleanly.
- `docs/releasing.md` describes the stamp-then-tag recurring ritual.

## Out of scope

- The broader 0.21.x (decision abstraction) and 0.22.x (dynamic workflow UI via mcp-app)
  bands — separate shaping cycles.
- The 0.19.8 install-refresh correctness + upgrade hint beyond what migration needs (`tes`).
- Contract-version semantics — unchanged at 1.

## Provenance

- Decision + grounding: plugin-distribution research `wcdgsgd88`.
- Refuted claim (do NOT build on): the research aside that *"only Homebrew exposes a
  branch-HEAD model"* — npm git-URL deps also resolve a branch tip. The load-bearing part
  (every ecosystem's *stable* channel serves a pinned artifact) survives.
