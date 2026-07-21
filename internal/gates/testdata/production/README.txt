These eight fixtures are exact frontmatter snapshots of the eight 0260 production
entities named by the recorder design. They were copied from the spacedock-state/dev
checkout at 9594033dab6c8f0ebc90438cb115cd1e2523e992. The test rewrites each fixture through
the production recorder and compares the complete decoded gates document before and
after, including opaque application and unknown historical fields.

The source entity slugs match the fixture filenames. The snapshots intentionally retain
legacy records with absent digests, first-retained sequence numbers greater than one,
multiple logical gates, application subtrees, and historical extension fields.
