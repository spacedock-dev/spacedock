// ABOUTME: Schema-driven per-field conformance — drives type/pattern/enum checks
// ABOUTME: from the embedded entity.mdschema.yml at the field's declared severity.
package status

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	spacedock "github.com/spacedock-dev/spacedock"
	"gopkg.in/yaml.v3"
)

// schemaField is one entry of the entity mdschema's frontmatter.fields map. Only
// the keys the field-conformance loop reads are surfaced; unknown keys (semantics,
// coerce_*, etc.) are ignored by yaml.v3.
type schemaField struct {
	Type            string   `yaml:"type"`
	Pattern         string   `yaml:"pattern"`
	Conventional    []string `yaml:"conventional"`
	InvalidSeverity string   `yaml:"invalid_severity"`
	UnknownSeverity string   `yaml:"unknown_severity"`
}

// entitySchema is the parsed entity mdschema's field declarations — the SSOT the
// validator drives per-field checks from.
type entitySchema struct {
	fields map[string]schemaField
}

var (
	loadedSchema     *entitySchema
	loadedSchemaOnce sync.Once
)

// loadEntitySchema parses the embedded entity.mdschema.yml once and returns its
// field declarations. The schema is JSON (a YAML subset), so yaml.v3 parses it
// directly. A parse failure yields an empty field set — field conformance is a
// warn-tier advisory and must never block the binary; a corrupt embedded schema
// is a build-time problem the schema-driven test (AC-4) catches.
func loadEntitySchema() *entitySchema {
	loadedSchemaOnce.Do(func() {
		var doc struct {
			Frontmatter struct {
				Fields map[string]schemaField `yaml:"fields"`
			} `yaml:"frontmatter"`
		}
		s := &entitySchema{fields: map[string]schemaField{}}
		if err := yaml.Unmarshal(spacedock.EntityMDSchema, &doc); err == nil {
			s.fields = doc.Frontmatter.Fields
		}
		loadedSchema = s
	})
	return loadedSchema
}

// isoLayouts are the timestamp shapes the spacedock writer emits (second and
// microsecond precision, UTC Z). A value matching any layout is a valid
// iso8601; the date-only RFC3339 prefix is accepted for fields the FO may set
// by hand.
var isoLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05.000000Z",
	"2006-01-02",
}

// isISO8601 reports whether s parses as one of the accepted ISO-8601 layouts.
func isISO8601(s string) bool {
	for _, layout := range isoLayouts {
		if _, err := time.Parse(layout, s); err == nil {
			return true
		}
	}
	return false
}

// fieldConformanceWarnings returns warn-tier diagnostics for every active
// entity whose frontmatter violates a schema field's declared
// pattern / conventional enum / type, at the field's declared severity.
// Archived entities are skipped: archived scope is publish-only, so a
// non-conformant field there can never be cleared by a tool-mediated write. Only
// fields whose severity is `warn` are surfaced here (no schema field is `error`
// today; an `error`-severity field would route through the structural-error path
// instead). Each line names the field and the violated rule plus the entity
// evidence so the FO can locate the entity.
func fieldConformanceWarnings(entities []*entity, workflowDir string) []string {
	schema := loadEntitySchema()
	var warns []string
	for _, e := range entities {
		if e.scope != "active" {
			continue
		}
		for name, spec := range schema.fields {
			if !isWarnSeverity(spec) {
				continue
			}
			value, present := e.fields[name]
			if !present || value == "" {
				continue
			}
			problem := fieldViolation(name, spec, value)
			if problem == "" {
				continue
			}
			warns = append(warns, entityEvidenceLine("Warning", e, workflowDir, problem, e.displayID))
		}
	}
	return warns
}

// isWarnSeverity reports whether a schema field carries a `warn` severity for
// its conformance check (either invalid_severity or unknown_severity). Fields
// with a pattern/conventional/type check but no explicit severity default to
// warn — the schema marks every per-field check warn today, and a field-level
// conformance gap must never harden into an exit-1 gate without an explicit
// `error` severity in the schema.
func isWarnSeverity(spec schemaField) bool {
	switch spec.InvalidSeverity {
	case "warn":
		return true
	case "error":
		return false
	}
	if spec.UnknownSeverity == "error" {
		return false
	}
	// No explicit invalid_severity: a checkable field (pattern/enum/type) with no
	// declared severity is advisory (warn). Fields with nothing to check are
	// skipped by fieldViolation returning "".
	return spec.Pattern != "" || len(spec.Conventional) > 0 || spec.Type == "iso8601" || spec.Type == "numeric_string"
}

// fieldViolation returns a human-readable problem string when value violates the
// field's schema-declared rule, or "" when it conforms / has nothing to check.
// The rule (pattern, enum, type) is read from the schema spec — never a Go
// literal — so editing the schema changes what is enforced.
func fieldViolation(name string, spec schemaField, value string) string {
	if spec.Pattern != "" {
		re, err := regexp.Compile(spec.Pattern)
		if err != nil {
			return ""
		}
		if !re.MatchString(value) {
			return fmt.Sprintf("field '%s' value %q does not match pattern %s", name, value, spec.Pattern)
		}
		return ""
	}
	if len(spec.Conventional) > 0 {
		// Case-insensitive: a conventional enum fixes which token was chosen, not
		// its casing. Entities written before the verdict writer normalised case
		// (`verdict: passed`) are semantically conformant, and warning on them
		// would demand hand-editing already-archived terminal entities to silence
		// a diagnostic about a value the tool itself wrote.
		for _, allowed := range spec.Conventional {
			if strings.EqualFold(value, allowed) {
				return ""
			}
		}
		return fmt.Sprintf("field '%s' value %q is not one of [%s]", name, value, strings.Join(spec.Conventional, " "))
	}
	switch spec.Type {
	case "numeric_string":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Sprintf("field '%s' value %q is not numeric", name, value)
		}
	case "iso8601":
		if !isISO8601(value) {
			return fmt.Sprintf("field '%s' value %q is not a valid ISO-8601 timestamp", name, value)
		}
	}
	return ""
}
