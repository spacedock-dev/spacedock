// ABOUTME: Package-level policy for the instruction-file-read quarantine.
// ABOUTME: Structural checks live here; prose-grep and code-bound behavior substitutes do not.
package contractlint

// This package is the instruction-file-read quarantine for tests. Tests outside
// this package must not read skill, contract, runtime-adapter, or agent markdown.
// Tests inside it may read those files only for structural facts a machine can
// verify without interpreting prose: reference closure, frontmatter validity,
// structural absence, deduplication, line/count floors, and portability markers.
//
// Do not add prose-grep checks here. Do not add prose-to-code consistency checks
// as behavior substitutes. If behavior matters, test it by running the behavior;
// if no behavior test exists yet, delete the read and report the owed test.
