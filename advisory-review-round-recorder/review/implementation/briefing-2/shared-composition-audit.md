# Advisory-round shared-composition audit

Classification: bounded in-design correction; ACs and architecture unchanged.

The corrected checkpoint genuinely shares entity mutation, room publication, the round loader, and landed Resolution validation. It is not another mechanism reset. However, 683 net production LOC before CLI is not acceptable against the declared 365 stop, and the checkpoint still has material defects:

- canonical Annotation validation is not yet shared with the provider path;
- multiple authorized `actor:ensign` triage Resolutions are not rejected;
- three top-level rebuilders duplicate frontmatter/range/CRLF splicing;
- Feedback Cycles matching is not scoped to the owned section;
- malformed artifact URIs can be treated as remote and entity-root symlink resolution repeats per artifact;
- record/validate retain avoidable digest/summary work and path-derived identity duplication;
- the direct publisher tests do not yet prove whole-operation stale-room/rollback cleanup.

The smallest credible surface is 500-520 LOC before CLI and 555-575 total. The authorized hard stops are therefore 540 pre-CLI and 600 total. This is a measured correction to a bad estimate, not permission for more mechanism.

Required pass:

1. introduce and reuse one canonical Annotation validator in provider and review-log paths;
2. reject more than one authorized worker triage Resolution and retain the second-reviewer falsifier;
3. share the top-level replacement builder and remove unused Status-only expectation breadth;
4. make Feedback Cycles insertion section-scoped and single-pass;
5. reject malformed artifact URIs, resolve the entity root once, remove repeated digest/summary/identity work;
6. exercise stale-room and rollback through the complete round operation with whole-tree cleanup assertions;
7. stop again before CLI if net production exceeds 540 or any AC guard weakens.
