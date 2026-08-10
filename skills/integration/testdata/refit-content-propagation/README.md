# Refit content-propagation fixture

`site-workflow/` is a dev-shape workflow README as the commission skill emitted it at
`spacedock@0.25.0`, before the template carried the anti-over-engineering scar tissue.
It was produced by a dispatched agent driving that revision of the commission skill for
the mission "a personal-site development workflow", so it is what commission actually
emits rather than a hand-written approximation: it carries the
`Verified by: {grep / ...}` example and none of the scar-tissue content.

Drive `skills/refit/SKILL.md` Phase 3b against `site-workflow/` and read the emitted
README diff. Against a template that has gained content, the diff carries those content
hunks, including the gated stages' `Gate content` instructions; against a template that
has not, it carries only the `commissioned-by:` stamp line.
That contrast is the observation the fixture exists for.
