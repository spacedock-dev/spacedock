package ensigncycle

type codexAuthMode int

const (
	codexAuthSkip codexAuthMode = iota
	codexAuthRun
	codexAuthFatal
)

type codexLiveAuthDecision struct {
	mode    codexAuthMode
	message string
}

func decideCodexLiveAuth(openAIAPIKey, required string) codexLiveAuthDecision {
	if openAIAPIKey != "" {
		return codexLiveAuthDecision{mode: codexAuthRun}
	}
	if required != "" {
		return codexLiveAuthDecision{
			mode:    codexAuthFatal,
			message: "OPENAI_API_KEY is required for the approval-gated codex-live lane",
		}
	}
	return codexLiveAuthDecision{
		mode:    codexAuthSkip,
		message: "no live Codex auth available: set OPENAI_API_KEY to run the live Codex smoke",
	}
}
