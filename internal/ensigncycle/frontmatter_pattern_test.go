package ensigncycle

import (
	"regexp"
	"strings"
	"testing"
)

// frontmatterPattern matches only inside an entity's leading frontmatter block,
// so a body line such as a stage report's "Verdict: PASS" can neither satisfy
// nor trip an assertion about a frontmatter field.
type frontmatterPattern struct{ re *regexp.Regexp }

func frontmatterRegexp(expr string) frontmatterPattern {
	return frontmatterPattern{re: regexp.MustCompile(expr)}
}

func (p frontmatterPattern) MatchString(entity string) bool {
	return p.re.MatchString(entityFrontmatter(entity))
}

// entityFrontmatter returns the lines between the opening and closing `---`,
// or "" when the entity has no frontmatter or it is unclosed.
func entityFrontmatter(entity string) string {
	rest, ok := strings.CutPrefix(entity, "---\n")
	if !ok {
		return ""
	}
	if strings.HasPrefix(rest, "---\n") || rest == "---" {
		return ""
	}
	if end := strings.Index(rest, "\n---\n"); end >= 0 {
		return rest[:end+1]
	}
	if strings.HasSuffix(rest, "\n---") {
		return rest[:len(rest)-len("---")]
	}
	return ""
}

func TestFrontmatterPattern(t *testing.T) {
	status := frontmatterRegexp(`(?im)^status:\s*done\s*$`)
	for name, tc := range map[string]struct {
		entity string
		want   bool
	}{
		"frontmatter match":  {"---\ntitle: T\nstatus: done\n---\n# T\n", true},
		"body-only decoy":    {"---\ntitle: T\nstatus: validation\n---\n# T\n\nStatus: done\n", false},
		"no frontmatter":     {"# T\n\nstatus: done\n", false},
		"unclosed":           {"---\ntitle: T\nstatus: done\n", false},
		"closed at file end": {"---\nstatus: done\n---", true},
	} {
		if got := status.MatchString(tc.entity); got != tc.want {
			t.Errorf("%s: MatchString = %v, want %v", name, got, tc.want)
		}
	}
}
