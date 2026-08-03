// ABOUTME: Builds the per-command --json envelopes (status/next/resolve/boot/
// ABOUTME: next-id/validate) as ordered string objects over the same data the table renders.
package status

import (
	"sort"
	"strconv"
	"strings"
)

// resolveJSONFields returns the ordered key set for a status-shape JSON object.
// Strict in JSON mode: explicit --fields is exactly the named order; --all-fields
// is defaults followed by sorted non-empty non-underscore non-default keys; no
// flag is the defaults. Mirrors resolveExtraFields' --all-fields scan but yields
// the full key set (defaults + extras), since JSON keys ARE the projection.
func resolveJSONFields(entities []*entity, explicitFields []string, allFields bool, defaultFields []string) []string {
	if explicitFields != nil {
		return append([]string(nil), explicitFields...)
	}
	if allFields {
		defaults := map[string]bool{}
		for _, f := range defaultFields {
			defaults[f] = true
		}
		seen := map[string]bool{}
		for _, e := range entities {
			for key, val := range e.fields {
				if strings.HasPrefix(key, "_") || defaults[key] || val == "" {
					continue
				}
				seen[key] = true
			}
		}
		extras := make([]string, 0, len(seen))
		for k := range seen {
			extras = append(extras, k)
		}
		sort.Strings(extras)
		return append(append([]string(nil), defaultFields...), extras...)
	}
	return append([]string(nil), defaultFields...)
}

// entityJSONObj projects one entity to an ordered object over fields, reading
// each value from e.fields (id is already the display id post-applyEffectiveIDs).
func entityJSONObj(e *entity, fields []string) *jsonObj {
	o := newJSONObj()
	for _, f := range fields {
		o.set(f, e.fields[f])
	}
	return o
}

// statusJSON builds the {"command":"status","entities":[...],"pagination":{...}}
// envelope for the default / --archived / --where reads. Array order is
// sortDefault; entities is bounded to the win slice.
func statusJSON(entities []*entity, stages []Stage, fields []string, win paginationWindow) *jsonObj {
	sorted := sortDefault(entities, stages)
	page := sorted[win.start:win.end]
	arr := make(jsonArr, 0, len(page))
	for _, e := range page {
		arr = append(arr, entityJSONObj(e, fields))
	}
	return newJSONObj().set("command", "status").setValue("entities", arr).setValue("pagination", paginationJSONObj(win))
}

// paginationJSONObj renders the pagination window as the all-strings object
// beside entities: page, limit, total, start, end (1-based inclusive, (0,0)
// for an empty window), has_next.
func paginationJSONObj(w paginationWindow) *jsonObj {
	start, end := w.display()
	return newJSONObj().
		set("page", strconv.Itoa(w.page)).
		set("limit", strconv.Itoa(w.limit)).
		set("total", strconv.Itoa(w.total)).
		set("start", strconv.Itoa(start)).
		set("end", strconv.Itoa(end)).
		set("has_next", strconv.FormatBool(w.hasNext))
}

// nextFixedFields are the always-present --next keys: id, slug, plus the three
// computed dispatch columns. --fields adds frontmatter keys after these.
var nextFixedFields = []string{"id", "slug", "current", "next", "worktree"}

// nextJSON builds the {"command":"next","dispatchable":[...]} envelope. The
// fixed five are always present; explicit/--all-fields frontmatter keys are
// additive after them (the computed columns are not projectable, per spike).
func nextJSON(entities []*entity, stages []Stage, explicitFields []string, allFields bool) *jsonObj {
	disp := computeDispatchable(entities, stages)
	extras := resolveNextExtras(entities, explicitFields, allFields)
	arr := dispatchableJSONArr(disp)
	for i, d := range disp {
		for _, f := range extras {
			arr[i].set(f, d.e.fields[f])
		}
	}
	return newJSONObj().set("command", "next").setValue("dispatchable", arr)
}

// resolveNextExtras returns the additive frontmatter keys for --next JSON:
// explicit --fields minus any that collide with the fixed five, or the
// --all-fields scan minus the fixed five.
func resolveNextExtras(entities []*entity, explicitFields []string, allFields bool) []string {
	fixed := map[string]bool{}
	for _, f := range nextFixedFields {
		fixed[f] = true
	}
	if explicitFields != nil {
		var out []string
		for _, f := range explicitFields {
			if !fixed[f] {
				out = append(out, f)
			}
		}
		return out
	}
	if allFields {
		seen := map[string]bool{}
		for _, e := range entities {
			for key, val := range e.fields {
				if strings.HasPrefix(key, "_") || fixed[key] || val == "" {
					continue
				}
				seen[key] = true
			}
		}
		out := make([]string, 0, len(seen))
		for k := range seen {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}
	return nil
}

// dispatchableJSONArr builds the dispatchable array shared by --next and --boot:
// the fixed five keys per row, no projected extras (boot takes no --fields).
func dispatchableJSONArr(disp []dispatchable) jsonArr {
	arr := make(jsonArr, 0, len(disp))
	for _, d := range disp {
		o := newJSONObj()
		o.set("id", d.e.fields["id"])
		o.set("slug", d.e.fields["slug"])
		o.set("current", d.e.fields["status"])
		o.set("next", d.next)
		o.set("worktree", d.nextWorktree)
		arr = append(arr, o)
	}
	return arr
}

// readyGatesJSONArr builds the identify-only scheduling index. Its four fixed
// keys identify the entity and the reduced durable lifecycle state.
func readyGatesJSONArr(entities []*entity) jsonArr {
	arr := make(jsonArr, 0, len(entities))
	for _, e := range entities {
		arr = append(arr, newJSONObj().
			set("id", e.fields["id"]).
			set("slug", e.fields["slug"]).
			set("current", e.fields["status"]).
			set("readiness", e.fields["gate-readiness"]))
	}
	return arr
}

// bootJSON builds the nested {"command":"boot",...} envelope from gathered boot
// data. mods is an object of point->[mods] (empty -> {}); min_prefix is present
// only for sd-b32; team_state.present is the string "true"/"false".
func bootJSON(d *bootData) *jsonObj {
	out := newJSONObj().set("command", "boot")

	mods := newJSONObj()
	points := make([]string, 0, len(d.hooks))
	for p := range d.hooks {
		points = append(points, p)
	}
	sort.Strings(points)
	for _, p := range points {
		vals := append([]string(nil), d.hooks[p]...)
		sort.Strings(vals)
		mods.setValue(p, jsonStrArr(vals))
	}
	out.setValue("mods", mods)

	out.set("id_style", d.idStyle)
	out.set("next_id", d.nextID)
	if d.idStyle == "sd-b32" {
		out.set("min_prefix", strconv.Itoa(sdB32MinPrefix))
	}

	orphans := make(jsonArr, 0, len(d.orphans))
	for _, o := range d.orphans {
		orphans = append(orphans, newJSONObj().
			set("id", o.id).set("slug", o.slug).set("worktree", o.worktree).
			set("dir_exists", o.dirExists).set("branch_exists", o.branchExists))
	}
	out.setValue("orphans", orphans)

	prState := newJSONObj().set("status", d.prStatus)
	prEntries := make(jsonArr, 0, len(d.prResults))
	for _, r := range d.prResults {
		prEntries = append(prEntries, newJSONObj().
			set("id", r.id).set("slug", r.slug).set("pr", r.pr).set("state", r.state))
	}
	prState.setValue("entries", prEntries)
	out.setValue("pr_state", prState)

	out.setValue("dispatchable", dispatchableJSONArr(d.dispatchable))

	// present/absent hints are both resolved in gatherBoot so this renderer carries
	// no host-specific string; present is the string-bool team_state convention.
	team := newJSONObj()
	if d.teamPresent {
		team.set("present", "true")
	} else {
		team.set("present", "false")
	}
	team.set("hint", d.teamHint)
	out.setValue("team_state", team)

	// State backend keys, appended AFTER team_state so the FO's key-order parse of
	// every prior key is preserved. entity_dir_present is the string "true"/"false"
	// matching team_state.present's string-bool convention.
	out.set("state_backend", d.stateBackend)
	out.set("definition_dir", d.definitionDir)
	out.set("entity_dir", d.entityDir)
	out.set("entity_dir_present", strconv.FormatBool(d.entityDirPresent))
	// state_remote ("origin"/"none") is appended AFTER entity_dir_present, present
	// only under split-root where remote sync applies. Single-root omits it so the
	// envelope carries no remote concept where there is none.
	if d.stateRemote != "" {
		out.set("state_remote", d.stateRemote)
	}

	// sandbox: the three-way safehouse posture, appended AFTER the state-backend keys
	// so every existing key's relative order is preserved for the FO's key-order parse.
	out.set("sandbox", d.sandbox)

	// Identify mode folds the workflow discovery result and the stage taxonomy into
	// the record, appended AFTER every existing key so the un-flagged --boot key set
	// and its pinned order are byte-for-byte unchanged. Absent otherwise.
	if d.identify {
		out.setValue("discovery", jsonStrArr(d.discovery))
		out.setValue("stages", stagesJSONArr(d.stages))
		out.setValue("ready_gates", readyGatesJSONArr(d.readyGates))
	}

	return out
}

// resolveJSON builds the single {"command":"resolve",...} object. id is the
// display id (uniform with every other command); stored_id is the full stored
// form the text resolve line carries, so cross-command round-trip holds.
func resolveJSON(workflowDir string, e *entity) *jsonObj {
	return newJSONObj().
		set("command", "resolve").
		set("workflow", realpathOf(workflowDir)).
		set("scope", scopeOf(e)).
		set("slug", e.slug).
		set("id", e.fields["id"]).
		set("stored_id", e.storedID).
		set("path", e.path)
}

// singletonJSON builds a {"command":cmd,key:value} object for the single-token
// read outputs (--next-id, --short-id, --validate) when --json is explicit.
func singletonJSON(cmd, key, value string) *jsonObj {
	return newJSONObj().set("command", cmd).set(key, value)
}

// readJSON builds the {"command":"read",...} envelope for --read: the realpath'd
// path, the file's total_lines, the frontmatter as a nested object, and the
// ordered headings array (text/level/offset/lines, every value stringified to
// hold the all-strings contract). With fields non-nil, the frontmatter object is
// projected to exactly those keys in their named order (a missing key is the
// empty string), mirroring entityJSONObj's --where/--next semantics; with no
// fields the whole map is emitted with keys sorted for byte stability (since
// ParseFrontmatter returns an unordered map).
func readJSON(path string, sr sectionRead, fields []string) *jsonObj {
	fm := newJSONObj()
	if fields != nil {
		// Named projection: exactly the requested keys, in request order, missing
		// key -> empty string. resolveJSONFields with explicit fields returns them
		// verbatim, so the projected order is the caller's order. This mirrors
		// entityJSONObj's semantics over sr.frontmatter (readJSON has no *entity).
		for _, k := range resolveJSONFields(nil, fields, false, nil) {
			fm.set(k, sr.frontmatter[k])
		}
	} else {
		keys := make([]string, 0, len(sr.frontmatter))
		for k := range sr.frontmatter {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fm.set(k, sr.frontmatter[k])
		}
	}

	headings := make(jsonArr, 0, len(sr.headings))
	for _, h := range sr.headings {
		headings = append(headings, newJSONObj().
			set("text", h.text).
			set("level", strconv.Itoa(h.level)).
			set("offset", strconv.Itoa(h.offset)).
			set("lines", strconv.Itoa(h.lines)))
	}

	obj := newJSONObj().
		set("command", "read").
		set("path", realpathOf(path)).
		set("total_lines", strconv.Itoa(sr.totalLines)).
		setValue("frontmatter", fm).
		setValue("headings", headings)

	// A file carrying a stages: taxonomy (the workflow README) surfaces it as a
	// structured sibling array — every leaf a string, like headings — so a reader
	// consumes the stage names/ordering and the per-stage flags machine-readably
	// instead of re-deriving them from the flattened "stages":"" scalar. A file
	// with no stages: block emits no array (the array is keyed on the block).
	if stages := parseStagesBlock(path); len(stages) > 0 {
		obj.setValue("stages", stagesJSONArr(stages))
	}

	return obj
}

// stagesJSONArr renders the resolved workflow stages as an ordered array of
// ordered objects, every leaf a string. The typed gate/terminal/initial/worktree
// flags are always present (a stage always has a resolved bool); the optional
// feedback-to/agent/fresh/model keys appear only when the stage declares them,
// matching the stages: block's own presence semantics.
func stagesJSONArr(stages []Stage) jsonArr {
	arr := make(jsonArr, 0, len(stages))
	for _, s := range stages {
		o := newJSONObj().
			set("name", s.Name).
			set("worktree", strconv.FormatBool(s.Worktree)).
			set("gate", strconv.FormatBool(s.gate)).
			set("terminal", strconv.FormatBool(s.terminal)).
			set("initial", strconv.FormatBool(s.initial))
		for _, of := range []string{"feedback-to", "agent", "fresh", "model"} {
			if v, ok := s.optional[of]; ok {
				o.set(of, v)
			}
		}
		arr = append(arr, o)
	}
	return arr
}
