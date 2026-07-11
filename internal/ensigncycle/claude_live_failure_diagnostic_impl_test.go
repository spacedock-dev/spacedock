// ABOUTME: Selects Claude live detector evidence and reports it only after another failure.
// ABOUTME: Keeps wrong-root and broad-search observations diagnostic-only and silent on success.
package ensigncycle

type claudeLiveDiagnosticReporter interface {
	Cleanup(func())
	Failed() bool
	Logf(string, ...any)
}

func detectClaudeLiveFailureDiagnostic(stream, workflowRoot string) error {
	if wrongRoot := detectWrongRootBoot(stream, workflowRoot); wrongRoot != nil {
		return wrongRoot
	}
	return detectBroadSearchAtBoot(stream, workflowRoot)
}

func registerClaudeLiveFailureDiagnostic(reporter claudeLiveDiagnosticReporter, diagnostic error) {
	if diagnostic == nil {
		return
	}
	reporter.Cleanup(func() {
		if reporter.Failed() {
			reporter.Logf("Additional diagnostic: %v", diagnostic)
		}
	})
}
