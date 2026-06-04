// ABOUTME: Runtime-observable AC classifier (lead-with-`live` convention) and
// ABOUTME: the three-way ci-run:/session: live-run citation resolver.
package status

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// LiveACFlag is one runtime-observable AC finding: the AC header line and the
// citation text trailing the `live` lead (e.g. `ci-run:123`, `session:/p.jsonl`,
// or a placeholder like `<no artifact yet>`). The gate resolves the citation to
// decide refuse / pass / indeterminate.
type LiveACFlag struct {
	Header   string
	Citation string
}

// liveLeadRe matches a proof clause that leads with the explicit `live`
// convention — `Verified by: live <ref>` — and captures the trailing citation
// text. Case-insensitive on `live` so `Live`/`LIVE` are honored; the lead must
// follow the marker so a passing mention of "live" mid-clause does not classify
// (the convention is "lead WITH live", an authored property of the clause).
var liveLeadRe = regexp.MustCompile(`(?is)^(?:verified by|oracle:|proof:|end state[:.])\s*:?\s*live\b\s*(.*)$`)

// classifyLiveACs walks the entity body's `**AC-N` blocks (the same block-walk
// shape ClassifyEntityACs uses), isolates each block's proof clause from its
// marker, and returns a flag for every clause that leads with `live`. The
// captured Citation is the trimmed text after `live` — what the gate resolves.
// Pure (no I/O) so tests drive it with literal bodies.
func classifyLiveACs(body string) []LiveACFlag {
	var flags []LiveACFlag

	lines := splitLines(body)
	var current []string
	var currentHeader string

	flush := func() {
		if currentHeader == "" {
			return
		}
		block := strings.Join(current, "\n")
		clause := strings.TrimSpace(isolateProofClause(block))
		if m := liveLeadRe.FindStringSubmatch(clause); m != nil {
			flags = append(flags, LiveACFlag{
				Header:   currentHeader,
				Citation: strings.TrimSpace(m[1]),
			})
		}
	}

	for _, line := range lines {
		if acHeaderRe.MatchString(line) {
			flush()
			currentHeader = line
			current = nil
			continue
		}
		if strings.HasPrefix(line, "## ") {
			flush()
			currentHeader = ""
			current = nil
			continue
		}
		if currentHeader != "" {
			current = append(current, line)
		}
	}
	flush()

	return flags
}

// liveResolutionKind is the three-way outcome of resolving a live-run citation.
type liveResolutionKind int

const (
	// citedAndReal: the cited live run exists (ci-run id resolves, or session
	// .jsonl is present on disk). The gate passes.
	citedAndReal liveResolutionKind = iota
	// definitivelyAbsent: the citation is a placeholder, or the cited run
	// genuinely does not exist (ci-run 404, session path absent). The gate
	// refuses — this is the yy slip.
	definitivelyAbsent
	// indeterminate: the citation COULD NOT BE CHECKED (connectivity/auth
	// error reaching GitHub). The gate surfaces a tooling error, NOT a refusal,
	// so a network blip never masquerades as a missing live run.
	indeterminate
)

// liveResolution carries the resolution kind plus the underlying tooling error
// for the indeterminate case (surfaced verbatim so an operator sees the real
// `gh` failure).
type liveResolution struct {
	kind liveResolutionKind
	err  error
}

// errConnectivity is the sentinel a ciRunResolver returns when it could not
// reach GitHub (connectivity/auth failure) — distinct from a definitive 404.
// On this error the resolution is indeterminate, never a refusal.
var errConnectivity = errors.New("live-run citation could not be checked")

// ciRunResolver answers "does CI run <id> exist?" with a three-way signal:
// (true, nil) → exists; (false, nil) → definitively absent (404); (_, err) →
// could not check (connectivity/auth). The default shells to `gh`; tests inject
// a stub for offline determinism.
type ciRunResolver func(id string) (bool, error)

// placeholderCiteRe matches a citation that is a placeholder (angle-bracketed
// `<…>`) rather than a real artifact ref. A placeholder is definitivelyAbsent.
var placeholderCiteRe = regexp.MustCompile(`^<.*>$`)

// resolveLiveCitation resolves a live-AC citation to a three-way outcome. An
// empty or `<…>` placeholder is definitivelyAbsent. A `ci-run:<id>` is resolved
// by the injected resolver (default ghCiRunResolver) and mapped three-way. A
// `session:<path>` resolves by checking the .jsonl is present on disk (offline,
// no network). Any unrecognized citation shape is definitivelyAbsent — an
// unresolvable ref is not a live run.
func resolveLiveCitation(citation string, resolver ciRunResolver) liveResolution {
	cite := strings.TrimSpace(citation)
	if cite == "" || placeholderCiteRe.MatchString(cite) {
		return liveResolution{kind: definitivelyAbsent}
	}

	switch {
	case strings.HasPrefix(cite, "ci-run:"):
		id := strings.TrimSpace(strings.TrimPrefix(cite, "ci-run:"))
		if id == "" || placeholderCiteRe.MatchString(id) {
			return liveResolution{kind: definitivelyAbsent}
		}
		if resolver == nil {
			resolver = ghCiRunResolver
		}
		exists, err := resolver(id)
		if err != nil {
			return liveResolution{kind: indeterminate, err: err}
		}
		if exists {
			return liveResolution{kind: citedAndReal}
		}
		return liveResolution{kind: definitivelyAbsent}
	case strings.HasPrefix(cite, "session:"):
		path := strings.TrimSpace(strings.TrimPrefix(cite, "session:"))
		if path == "" || placeholderCiteRe.MatchString(path) {
			return liveResolution{kind: definitivelyAbsent}
		}
		if info, err := os.Stat(path); err == nil && info.Mode().IsRegular() {
			return liveResolution{kind: citedAndReal}
		}
		return liveResolution{kind: definitivelyAbsent}
	default:
		return liveResolution{kind: definitivelyAbsent}
	}
}

// notFoundRe matches gh's DEFINITIVE 404 token — the PARENTHESIZED `(HTTP 404)`
// form gh emits on a genuine Not-Found (`gh: Not Found (HTTP 404)`). Anchoring to
// the parenthesized form is load-bearing: a bare `Not Found` substring matches a
// proxy/gateway 5xx body that merely contains the words, which would FALSE-REFUSE
// a valid entity (the exact network-blip-masquerades-as-missing-run regression
// this guard exists to prevent). Only the parenthesized token is gh's own 404.
var notFoundRe = regexp.MustCompile(`\(HTTP 404\)`)

// classifyGhRunOutput is the pure three-way discriminator over a `gh api` run
// probe's (exit-error, combined stdout+stderr). It is the crux of the
// connectivity-vs-404 distinction, factored out of the shell so it is testable
// with literal gh error strings:
//   - exit 0 → citedAndReal (the run resolved).
//   - a genuine parenthesized `(HTTP 404)` → definitivelyAbsent (refuse). NOTE:
//     a private/unscoped-token repo where the run EXISTS also yields a masked
//     404 (GitHub returns 404 not 403 to avoid leaking existence) — we still
//     refuse (cannot prove existence) but the gate diagnostic names the
//     token-scope remediation so a scope failure isn't sent chasing a wrong fix.
//   - any OTHER non-zero exit (auth 401, rate-limit 403, DNS/connect failure, a
//     5xx gateway body that merely contains the words "Not Found") →
//     indeterminate: a tooling error, NOT a refusal, so a blip never masquerades
//     as a missing live run.
func classifyGhRunOutput(runErr error, combined []byte) liveResolutionKind {
	if runErr == nil {
		return citedAndReal
	}
	if notFoundRe.Match(combined) {
		return definitivelyAbsent
	}
	return indeterminate
}

// ghCiRunResolver shells to `gh api repos/{owner}/{repo}/actions/runs/<id>` and
// maps the result through classifyGhRunOutput. A `gh`-not-on-PATH is itself a
// connectivity (indeterminate) condition. The (bool, error) seam keeps the
// ciRunResolver contract: a definitive 404 → (false, nil); anything else
// non-zero → (false, errConnectivity); exit 0 → (true, nil).
func ghCiRunResolver(id string) (bool, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return false, errConnectivity
	}
	cmd := exec.Command("gh", "api", "repos/{owner}/{repo}/actions/runs/"+id)
	combined, err := cmd.CombinedOutput()
	switch classifyGhRunOutput(err, combined) {
	case citedAndReal:
		return true, nil
	case definitivelyAbsent:
		return false, nil
	default:
		return false, errConnectivity
	}
}
