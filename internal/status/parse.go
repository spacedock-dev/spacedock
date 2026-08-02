// ABOUTME: Argv parsers matching the oracle's parse_* helpers — --workflow-dir,
// ABOUTME: --set, --where, --fields, --archive/--resolve/--short-id, id material.
package status

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// parseWorkflowDir extracts --workflow-dir and returns (dir, remaining, err).
// Matches parse_workflow_dir.
func parseWorkflowDir(args []string) (string, []string, error) {
	var remaining []string
	workflowDir := ""
	i := 0
	for i < len(args) {
		if args[i] == "--workflow-dir" {
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--workflow-dir requires an argument")
			}
			workflowDir = args[i+1]
			i += 2
			continue
		}
		remaining = append(remaining, args[i])
		i++
	}
	return workflowDir, remaining, nil
}

// parseRootArg extracts --root. Matches parse_root_arg.
func parseRootArg(args []string) (string, error) {
	for i := 0; i < len(args); i++ {
		if args[i] == "--root" {
			if i+1 >= len(args) {
				return "", fmt.Errorf("--root requires an argument")
			}
			return args[i+1], nil
		}
	}
	return "", nil
}

// parseSingleArg parses a flag taking a single non---- argument (--resolve,
// --short-id). label names the argument in the error. Matches parse_resolve_arg
// / parse_short_id_arg.
func parseSingleArg(args []string, flag, label string) (string, error) {
	for i := 0; i < len(args); i++ {
		if args[i] == flag {
			if i+1 >= len(args) {
				return "", fmt.Errorf("%s requires a %s argument", flag, label)
			}
			ref := args[i+1]
			if ref == "" || strings.HasPrefix(ref, "--") {
				return "", fmt.Errorf("%s requires a %s argument", flag, label)
			}
			return ref, nil
		}
	}
	return "", nil
}

// parseArchiveArg parses --archive <slug>. Matches parse_archive_arg.
func parseArchiveArg(args []string) (string, error) {
	for i := 0; i < len(args); i++ {
		if args[i] == "--archive" {
			if i+1 >= len(args) {
				return "", fmt.Errorf("--archive requires a slug argument")
			}
			slug := args[i+1]
			if slug == "" || strings.HasPrefix(slug, "--") {
				return "", fmt.Errorf("--archive requires a slug argument")
			}
			return slug, nil
		}
	}
	return "", nil
}

// parseFieldsOptions parses --fields and --all-fields. Returns (explicitFields,
// allFieldsFlag, err). Matches parse_fields_options.
func parseFieldsOptions(args []string) ([]string, bool, error) {
	var explicitFields []string
	explicitSet := false
	allFieldsFlag := false
	i := 0
	for i < len(args) {
		switch args[i] {
		case "--fields":
			if i+1 >= len(args) {
				return nil, false, fmt.Errorf("--fields requires a comma-separated list of field names")
			}
			raw := args[i+1]
			if strings.TrimSpace(raw) == "" {
				return nil, false, fmt.Errorf("--fields requires a comma-separated list of field names")
			}
			var parsed []string
			seen := map[string]bool{}
			for _, name := range strings.Split(raw, ",") {
				name = strings.TrimSpace(name)
				if name == "" || seen[name] {
					continue
				}
				seen[name] = true
				parsed = append(parsed, name)
			}
			if len(parsed) == 0 {
				return nil, false, fmt.Errorf("--fields requires a comma-separated list of field names")
			}
			explicitFields = parsed
			explicitSet = true
			i += 2
			continue
		case "--all-fields":
			allFieldsFlag = true
			i++
			continue
		}
		i++
	}
	if explicitSet && allFieldsFlag {
		return nil, false, fmt.Errorf("--fields and --all-fields are mutually exclusive")
	}
	if !explicitSet {
		return nil, allFieldsFlag, nil
	}
	return explicitFields, allFieldsFlag, nil
}

// parseIDMaterialOptions parses --id-seed/--id-actor. Returns (seed, actor,
// flagsSeen, err). Matches parse_id_material_options. flagsSeen lists which
// flags appeared, used for the "only with --next-id" guard.
func parseIDMaterialOptions(args []string) (string, string, []string, error) {
	values := map[string]string{"--id-seed": "", "--id-actor": ""}
	seen := map[string]bool{}
	var order []string
	i := 0
	for i < len(args) {
		arg := args[i]
		if _, ok := values[arg]; ok {
			if seen[arg] {
				return "", "", nil, fmt.Errorf("%s may only be supplied once", arg)
			}
			if i+1 >= len(args) {
				return "", "", nil, fmt.Errorf("%s requires an argument", arg)
			}
			value := args[i+1]
			if strings.HasPrefix(value, "--") {
				return "", "", nil, fmt.Errorf("%s requires an argument", arg)
			}
			values[arg] = value
			seen[arg] = true
			order = append(order, arg)
			i += 2
			continue
		}
		i++
	}
	return values["--id-seed"], values["--id-actor"], order, nil
}

// whereFilter is a parsed --where clause. value==nil means presence/absence.
type whereFilter struct {
	field string
	op    string // "=" or "!="
	value *string
}

const whereSyntaxHelp = "--where requires an operator: use 'field = value', 'field != value', 'field !=' (non-empty), or 'field =' (empty)"

// whereOperatorRe matches one --where operator. `!=` is listed first so the
// leftmost-first alternation consumes it whole at a given position instead of
// matching its trailing `=` as a second, spurious operator.
var whereOperatorRe = regexp.MustCompile(`!=|=`)

// countWhereOperators counts the operators in a single --where argument. A
// well-formed clause has exactly one; more than one means the argument is the
// compound-in-one-string shape (e.g. "sprint=A sprint-readiness!=defer") rather
// than a single clause — Check A of the #314 fix (see docs/dev/.spacedock-state
// status-where-robust-and-discoverable.md). All four compound shapes named in
// AC-1 count 2 here, regardless of which operator each half uses.
func countWhereOperators(whereArg string) int {
	return len(whereOperatorRe.FindAllStringIndex(whereArg, -1))
}

// parseWhereFilters parses all --where clauses. Matches parse_where_filters.
func parseWhereFilters(args []string) ([]whereFilter, error) {
	var filters []whereFilter
	i := 0
	for i < len(args) {
		if args[i] == "--where" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--where requires an argument")
			}
			whereArg := args[i+1]
			if strings.TrimSpace(whereArg) == "" {
				return nil, fmt.Errorf("--where argument cannot be empty")
			}
			if n := countWhereOperators(whereArg); n > 1 {
				return nil, fmt.Errorf(
					"--where takes one clause per flag; %q has %d operators — repeat --where to AND clauses instead of combining them in one string: --where 'field=value' --where 'field2!=value2'",
					whereArg, n)
			}
			var op, fieldPart, valuePart string
			if strings.Contains(whereArg, "!=") {
				op = "!="
				fieldPart, valuePart, _ = strings.Cut(whereArg, "!=")
			} else if strings.Contains(whereArg, "=") {
				op = "="
				fieldPart, valuePart, _ = strings.Cut(whereArg, "=")
			} else {
				return nil, fmt.Errorf("%s", whereSyntaxHelp)
			}
			field := strings.TrimSpace(fieldPart)
			if field == "" {
				return nil, fmt.Errorf("%s", whereSyntaxHelp)
			}
			value := strings.TrimSpace(valuePart)
			var valPtr *string
			if value != "" {
				v := value
				valPtr = &v
			}
			filters = append(filters, whereFilter{field: field, op: op, value: valPtr})
			i += 2
			continue
		}
		i++
	}
	return filters, nil
}

// derivedWhereFields are the field names --where/--fields can reference that
// are computed by a materialize* pass rather than read from frontmatter or the
// entity schema: materializeGateEligibility and materializeSuppressedBy only
// write these keys into e.fields when a filter or explicit field already names
// them (discover.go's referenced guard, format.go's materializeSuppressedBy),
// so collecting known field names off scanned entities after materialization
// would silently drop a name whenever nothing else in the same invocation
// referenced it first. Listed statically instead so validation never depends on
// that ordering.
var derivedWhereFields = []string{"gate-condition", "gate-eligible", "gate-readiness", suppressedByField}

// knownWhereFields returns the field names a --where clause may reference: the
// union of every scanned entity's frontmatter keys (covers workflow-specific
// fields like sprint/sprint-readiness, which the schema's permissive_additions
// intentionally leaves uncanonicalized), the canonical entity schema's declared
// field names (covers verdict/completed/archived/mod-block/pr/issue — canonical
// but absent from an active-only corpus where every archived member carrying
// them has been filtered out), and the derived names above.
func knownWhereFields(entities []*entity) map[string]bool {
	known := map[string]bool{}
	for _, e := range entities {
		for k := range e.fields {
			known[k] = true
		}
	}
	for k := range loadEntitySchema().fields {
		known[k] = true
	}
	for _, k := range derivedWhereFields {
		known[k] = true
	}
	return known
}

// validateWhereFields rejects a --where clause naming a field absent from
// knownWhereFields — Check B of the #314 fix. A misspelled or unknown field
// (e.g. "spint") reads as the empty string in applyFilters rather than erroring,
// which silently returns the wrong row set; this makes it a loud exit-1 instead.
// Skipped on a zero-entity scan: the result is empty either way, and every
// --where would otherwise error on a fresh workflow before any entity exists to
// populate the corpus half of the union.
func validateWhereFields(entities []*entity, filters []whereFilter) error {
	if len(entities) == 0 || len(filters) == 0 {
		return nil
	}
	known := knownWhereFields(entities)
	for _, f := range filters {
		if !known[f.field] {
			names := make([]string, 0, len(known))
			for name := range known {
				names = append(names, name)
			}
			sort.Strings(names)
			return fmt.Errorf("--where: unknown field %q — known fields: %s", f.field, strings.Join(names, ", "))
		}
	}
	return nil
}

// applyFilters keeps entities matching all --where clauses. Matches
// apply_filters.
func applyFilters(entities []*entity, filters []whereFilter) []*entity {
	if len(filters) == 0 {
		return entities
	}
	var out []*entity
	for _, e := range entities {
		match := true
		for _, f := range filters {
			fieldVal := e.fields[f.field]
			if f.op == "!=" {
				if f.value == nil {
					if fieldVal == "" {
						match = false
					}
				} else if fieldVal == *f.value {
					match = false
				}
			} else { // "="
				if f.value == nil {
					if fieldVal != "" {
						match = false
					}
				} else if fieldVal != *f.value {
					match = false
				}
			}
			if !match {
				break
			}
		}
		if match {
			out = append(out, e)
		}
	}
	return out
}

// setUpdate is a parsed --set target.
type setUpdate struct {
	slug    string
	updates []fieldUpdate
}

// parseSetArgs parses --set <slug> field=value... Matches parse_set_args. A
// token starting with -- terminates the field list (truncation).
func parseSetArgs(args []string) (*setUpdate, error) {
	for i := 0; i < len(args); i++ {
		if args[i] == "--set" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--set requires a slug argument")
			}
			slug := args[i+1]
			var updates []fieldUpdate
			j := i + 2
			for j < len(args) && !strings.HasPrefix(args[j], "--") {
				arg := args[j]
				if strings.Contains(arg, "=") {
					field, value, _ := strings.Cut(arg, "=")
					updates = append(updates, fieldUpdate{field: field, value: value, hasValue: true})
				} else if timestampFields[arg] {
					updates = append(updates, fieldUpdate{field: arg, hasValue: false})
				} else {
					return nil, fmt.Errorf(
						"bare field '%s' requires a value (use %s=value); only completed, started support auto-fill",
						arg, arg)
				}
				j++
			}
			if len(updates) == 0 {
				return nil, fmt.Errorf("--set requires at least one field=value argument")
			}
			return &setUpdate{slug: slug, updates: updates}, nil
		}
	}
	return nil, nil
}

// parsePageLimitArgs parses --page N and --limit N for the default status
// listing. page defaults to 1, limit defaults to defaultPageLimit; --limit 0
// disables pagination. pageSet/limitSet report whether the flag was
// explicitly supplied (so callers can reject it on non-listing commands and
// detect the --page-with-limit-0 contradiction — an explicit page selection
// makes no sense once --limit 0 asks for every row on one unbounded page).
func parsePageLimitArgs(args []string) (page, limit int, pageSet, limitSet bool, err error) {
	page = 1
	limit = defaultPageLimit
	i := 0
	for i < len(args) {
		switch args[i] {
		case "--page":
			if i+1 >= len(args) {
				return 0, 0, false, false, fmt.Errorf("--page requires an integer argument")
			}
			n, perr := strconv.Atoi(args[i+1])
			if perr != nil || n < 1 {
				return 0, 0, false, false, fmt.Errorf("--page must be a positive integer, got %q", args[i+1])
			}
			page = n
			pageSet = true
			i += 2
			continue
		case "--limit":
			if i+1 >= len(args) {
				return 0, 0, false, false, fmt.Errorf("--limit requires an integer argument")
			}
			n, perr := strconv.Atoi(args[i+1])
			if perr != nil || n < 0 {
				return 0, 0, false, false, fmt.Errorf("--limit must be a non-negative integer, got %q", args[i+1])
			}
			limit = n
			limitSet = true
			i += 2
			continue
		}
		i++
	}
	if pageSet && limit == 0 {
		return 0, 0, false, false, fmt.Errorf("--page cannot be combined with --limit 0 (--limit 0 disables pagination and returns every row on one page)")
	}
	return page, limit, pageSet, limitSet, nil
}

// parseNewArg parses --new [--folder] <slug>. Returns slug or "". The optional
// --folder flag may sit between --new and the slug; the slug is the next token
// and must not itself be --prefixed. A --new with no slug is an error.
func parseNewArg(args []string) (string, error) {
	for i := 0; i < len(args); i++ {
		if args[i] == "--new" {
			j := i + 1
			if j < len(args) && args[j] == "--folder" {
				j++
			}
			if j >= len(args) {
				return "", fmt.Errorf("--new requires a slug argument")
			}
			slug := args[j]
			if slug == "" || strings.HasPrefix(slug, "--") {
				return "", fmt.Errorf("--new requires a slug argument")
			}
			return slug, nil
		}
	}
	return "", nil
}
