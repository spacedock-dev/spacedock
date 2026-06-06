// ABOUTME: Pi-subagents assignment wrappers for Spacedock stage dispatches.
// ABOUTME: Keeps Pi transport fields additive around dispatch-build artifacts.
package piruntime

// SubagentDispatch is the Spacedock-owned subset of the pi-subagents stage
// dispatch wrapper. The canonical assignment remains the dispatch-build artifact;
// this wrapper only adds Pi transport metadata.
type SubagentDispatch struct {
	Task    string `json:"task"`
	Context string `json:"context"`
	Phase   string `json:"phase,omitempty"`
	Label   string `json:"label,omitempty"`
}

// SubagentStageDispatch wraps an already-built dispatch artifact prompt/content
// for pi-subagents. It intentionally has no acceptance field: stage acceptance is
// owned by the canonical assignment/checklist and independent Spacedock validation.
func SubagentStageDispatch(assignment, phase, label string) SubagentDispatch {
	return SubagentDispatch{
		Task:    assignment,
		Context: "fresh",
		Phase:   phase,
		Label:   label,
	}
}
