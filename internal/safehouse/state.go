// ABOUTME: The shared three-way sandbox-state render strings (enabled /
// ABOUTME: available-not-enabled / unavailable) the startup surfaces all source.
package safehouse

// State renders the three-way sandbox posture from the two inputs the launch
// decision already turns on: `selected` is whether a launch would be wrapped (a
// .safehouse profile is present, or a --safehouse* flag forced it), and
// `available` is whether the safehouse binary resolves on PATH. The launcher
// banner, `status --boot`, and `--version` all read these strings so the posture
// reads identically across surfaces.
//
//   - unavailable — the binary is not on PATH; nothing can wrap the launch, so a
//     present profile cannot take effect. This dominates even when selected.
//   - enabled — available and selected: the launch will be wrapped through safehouse.
//   - available, not enabled — available but nothing selects it: the binary is
//     installed, this launch is not sandboxed.
func State(selected, available bool) string {
	if !available {
		return "unavailable (safehouse not on PATH)"
	}
	if selected {
		return "enabled (safehouse)"
	}
	return "available, not enabled (no .safehouse profile)"
}
