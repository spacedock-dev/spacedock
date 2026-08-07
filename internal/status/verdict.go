// ABOUTME: conventional-value case boundary — lowercase on the CLI, the schema's
// ABOUTME: spelling in frontmatter, and case-insensitive on every read back.
package status

import "strings"

// canonicalConventional maps value to the conventional spelling it matches
// case-insensitively, and returns it UNCHANGED otherwise. The allowed set is
// read from the embedded entity mdschema — never a Go literal — so editing the
// schema changes what is canonicalised on write exactly as it changes what
// fieldViolation accepts on read.
//
// Unmatched values pass through untouched, deliberately. A conventional list is
// advisory (`invalid_severity: warn`), so a field may legitimately carry a value
// outside it; case-folding such a value would be an unasked-for edit to state
// the caller wrote on purpose. `verdict` is the only conventional enum in the
// entity schema today, but nothing here is verdict-specific.
func canonicalConventional(field, value string) string {
	spec, ok := loadEntitySchema().fields[field]
	if !ok || len(spec.Conventional) == 0 {
		return value
	}
	trimmed := strings.TrimSpace(value)
	for _, allowed := range spec.Conventional {
		if strings.EqualFold(trimmed, allowed) {
			return allowed
		}
	}
	return value
}

// storedVerdict maps a CLI verdict (`--verdict passed|rejected`, lowercase, as
// the help text documents) to the value written into frontmatter, where the
// entity schema declares the conventional set as [PASSED REJECTED]. The CLI
// surface deliberately does NOT change, since every caller, skill and doc passes
// lowercase.
//
// This is the guard for the ONE writer that does not go through
// updateFrontmatter: gates.FinalizeTerminalApproval owns its own locked write,
// so finalize() must hand it the stored spelling. Every other write — the
// user-facing `--set`, finalize's gate-less path, archive — is canonicalised
// inside updateFrontmatter and needs no caller-side normalisation.
//
// An empty verdict stays empty: it matches no conventional value, so it passes
// through, and a verdict-less transition writes nothing rather than an
// empty-string uppercase.
func storedVerdict(verdict string) string {
	return canonicalConventional("verdict", verdict)
}

// isRejectedVerdict reports whether a verdict READ back from frontmatter is the
// rejected one, ignoring case. Several guards carve out rejected entities (a
// rejected entity never ran the merge ceremony, so the merge-hook and
// terminal-status requirements are vacuous for it). Those carve-outs must hold
// for the `REJECTED` this binary writes, for the `rejected` older binaries wrote,
// and for whatever case a hand-edit left behind — a guard that refused on case
// alone would strand exactly the entities it is meant to exempt.
func isRejectedVerdict(verdict string) bool {
	return strings.EqualFold(strings.TrimSpace(verdict), "rejected")
}
