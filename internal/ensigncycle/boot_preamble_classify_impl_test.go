// ABOUTME: Pure classifier folding the wrong-root and broad-search boot detectors
// ABOUTME: into one ordered check claudeLiveRunner.run consults on its failure path.
package ensigncycle

// classifyBootPreambleFailure runs the boot-preamble detectors over a captured
// stream in priority order — wrong-root wander first (the most specific, earliest
// diagnosis: the FO operated outside the fixture entirely), then a broad-search
// filesystem sweep (the FO stayed on the fixture but hunted the filesystem for a
// workflow or a contract file instead of proceeding on what it already read) — and
// returns the first legible diagnosis it finds. A nil return means neither preamble
// class fired, so the caller's normal stall/assertion path applies unchanged.
//
// Extracted as a pure function (stream + workflowRoot in, error out) so the
// classification claudeLiveRunner.run performs on every launch is unit-testable
// offline, with no subprocess and no model spend.
func classifyBootPreambleFailure(stream, workflowRoot string) error {
	if wrongRoot := detectWrongRootBoot(stream, workflowRoot); wrongRoot != nil {
		return wrongRoot
	}
	return detectBroadSearchAtBoot(stream, workflowRoot)
}
