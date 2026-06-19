// ABOUTME: fotier resolution tests — the five-case tier+arming fixture (AC-2),
// ABOUTME: the capable-FO inert assertion (AC-4), and the alias-normalization table.
package fotier

import "testing"

// TestResolveTierArmingFiveCases is AC-2: the tier-resolution + arming decision
// over the five model/gated-stage combinations, asserting {Tier, RouteGateVerdicts}
// against each known input. The unset and garbage rows are the safety-critical
// ones — an un-provable session MUST default to level-2-only and arm routing, never
// silently capable (which would let a Haiku FO self-approve a verdict).
func TestResolveTierArmingFiveCases(t *testing.T) {
	cases := []struct {
		name          string
		model         string
		hasGatedStage bool
		wantTier      Tier
		wantRoute     bool
	}{
		{"haiku + gated stage → armed", "haiku", true, Level2Only, true},
		{"haiku + no gated stage → inert", "haiku", false, Level2Only, false},
		{"opus + gated stage → inert (capable)", "opus", true, Level3Capable, false},
		{"sonnet + gated stage → inert (capable)", "sonnet", true, Level3Capable, false},
		{"unset + gated stage → armed (fail-safe)", "", true, Level2Only, true},
		{"garbage model + gated stage → armed (fail-safe)", "gpt-4-turbo", true, Level2Only, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Resolve(tc.model, tc.hasGatedStage)
			if got.Tier != tc.wantTier {
				t.Errorf("Resolve(%q, %v).Tier = %q, want %q", tc.model, tc.hasGatedStage, got.Tier, tc.wantTier)
			}
			if got.RouteGateVerdicts != tc.wantRoute {
				t.Errorf("Resolve(%q, %v).RouteGateVerdicts = %v, want %v",
					tc.model, tc.hasGatedStage, got.RouteGateVerdicts, tc.wantRoute)
			}
		})
	}
}

// TestResolveCapableFOIsInert is AC-4: a level-3-capable FO (unset is fail-safe
// level-2-only, so the capable arms are an explicit opus/sonnet) returns
// RouteGateVerdicts=false even at a gated stage — the gate flow is the unmodified
// present-gate path, and the mechanism adds no cost to the common case.
func TestResolveCapableFOIsInert(t *testing.T) {
	for _, model := range []string{"opus", "sonnet"} {
		t.Run(model, func(t *testing.T) {
			got := Resolve(model, true)
			if got.Tier != Level3Capable {
				t.Errorf("Resolve(%q, true).Tier = %q, want %q", model, got.Tier, Level3Capable)
			}
			if got.RouteGateVerdicts {
				t.Errorf("Resolve(%q, true).RouteGateVerdicts = true, want false (capable FO is inert)", model)
			}
		})
	}
}

// TestResolveGateRoute pins the gate-route target the armed FO sends to. It is
// the standing teammate name, constant across cases.
func TestResolveGateRoute(t *testing.T) {
	if got := Resolve("haiku", true); got.GateRoute != "level-3-judge" {
		t.Errorf("armed GateRoute = %q, want level-3-judge", got.GateRoute)
	}
}

// TestNormalizeModel covers the launcher's alias normalization: the `--model`
// value the launcher reads off the passthrough is normalized to a canonical model
// name (haiku/sonnet/opus) that Resolve maps to a tier, or "" for an unresolvable
// value (which Resolve then treats fail-safe as level-2-only). `default` is the
// alias for the host's default model — opus — and the full model ids and [1m]
// suffix normalize to their family.
func TestNormalizeModel(t *testing.T) {
	cases := []struct {
		raw  string
		want string
	}{
		{"haiku", "haiku"},
		{"sonnet", "sonnet"},
		{"opus", "opus"},
		{"default", "opus"},
		{"", ""},
		{"Haiku", "haiku"},
		{"  sonnet  ", "sonnet"},
		{"sonnet[1m]", "sonnet"},
		{"claude-haiku-4-5-20251001", "haiku"},
		{"claude-opus-4-8[1m]", "opus"},
		{"claude-sonnet-4-6", "sonnet"},
		{"gpt-4-turbo", ""},
		{"random-garbage", ""},
	}
	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			if got := NormalizeModel(tc.raw); got != tc.want {
				t.Errorf("NormalizeModel(%q) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}

// TestResolveNormalizesBeforeMapping wires the two together: Resolve accepts a raw
// model string (the launcher may export the canonical name, but a hand-launched
// `--model claude-sonnet-4-6` could land in the var verbatim), so Resolve must
// normalize before mapping. A canonical capable id resolves capable; an
// unrecognized id falls through to the fail-safe level-2-only.
func TestResolveNormalizesBeforeMapping(t *testing.T) {
	if got := Resolve("claude-sonnet-4-6", true); got.Tier != Level3Capable || got.RouteGateVerdicts {
		t.Errorf("Resolve(claude-sonnet-4-6, true) = %+v, want capable+inert", got)
	}
	if got := Resolve("claude-haiku-4-5-20251001", true); got.Tier != Level2Only || !got.RouteGateVerdicts {
		t.Errorf("Resolve(claude-haiku-4-5, true) = %+v, want level-2-only+armed", got)
	}
}
