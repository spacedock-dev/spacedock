// ABOUTME: fotier resolves the FO's tier from its launcher-set model and decides
// ABOUTME: whether to route gate verdicts to a stronger standing teammate.
package fotier

import "strings"

// Tier is the FO's self-identified capability level. A level-2-only FO runs the
// contract by the book but routes every judgment call to a stronger teammate; a
// level-3-capable FO makes those calls itself.
type Tier string

const (
	// Level2Only drives the contract but does not adjudicate — gate verdicts,
	// scope calls, and conflict-recovery decisions route to level-3.
	Level2Only Tier = "level-2-only"
	// Level3Capable adjudicates judgment calls directly; the routing table is inert.
	Level3Capable Tier = "level-3-capable"
)

// gateRoute is the standing teammate an armed FO sends gate verdicts to.
const gateRoute = "level-3-judge"

// capableModels are the models that resolve to level-3-capable. Every other
// value — unset, garbage, an unrecognized id — falls through to the fail-safe
// level-2-only tier. Keeping the capable set as the explicit opt-in (rather than
// listing the weak models) is the safety property: a model can only earn
// adjudication authority by being recognized as one of these.
var capableModels = map[string]bool{
	"sonnet": true,
	"opus":   true,
}

// Result is what Resolve returns: the resolved tier, whether gate-verdict routing
// is armed for this session, and the teammate that authors the verdict when it is.
type Result struct {
	Tier              Tier
	RouteGateVerdicts bool
	GateRoute         string
}

// Resolve maps the launcher-set model (SPACEDOCK_FO_MODEL) to a tier and decides
// the gate-verdict arming. The model is normalized first, so a raw `--model`
// value that survived into the var verbatim (a hand-launched session, a resume)
// still resolves. The default is fail-SAFE: an unset, unrecognized, or garbage
// model resolves to level-2-only, never silently capable — a weak FO that lost the
// var by accident must route its verdicts rather than self-approve. Gate-verdict
// routing arms only when a level-2-only FO meets a gate:true stage; a capable FO
// arms nothing and the gate flow is unchanged.
func Resolve(model string, hasGatedStage bool) Result {
	tier := Level2Only
	if capableModels[NormalizeModel(model)] {
		tier = Level3Capable
	}
	r := Result{Tier: tier, GateRoute: gateRoute}
	if tier == Level2Only && hasGatedStage {
		r.RouteGateVerdicts = true
	}
	return r
}

// NormalizeModel canonicalizes a `--model` value to a bare model name (haiku /
// sonnet / opus) or "" when it cannot be resolved. It folds case, trims space,
// drops a trailing [..] extended-context suffix (e.g. [1m]), and matches a full
// `claude-{family}-…` id by its family segment. `default` is the alias for the
// host's default model (opus). An unrecognized value returns "" so Resolve treats
// it fail-safe as level-2-only rather than guessing a tier.
func NormalizeModel(raw string) string {
	m := strings.ToLower(strings.TrimSpace(raw))
	if i := strings.Index(m, "["); i >= 0 {
		m = m[:i]
	}
	switch m {
	case "":
		return ""
	case "default":
		return "opus"
	case "haiku", "sonnet", "opus":
		return m
	}
	for _, family := range []string{"haiku", "sonnet", "opus"} {
		if strings.HasPrefix(m, "claude-"+family+"-") {
			return family
		}
	}
	return ""
}
