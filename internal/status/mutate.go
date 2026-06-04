// ABOUTME: Mutation engine — update_frontmatter (--set) and run_archive
// ABOUTME: (--archive) parse the FM slice into yaml.Node, mutate, re-marshal.
package status

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// timestampFields auto-fill now() when given bare (no value). Matches
// TIMESTAMP_FIELDS.
var timestampFields = map[string]bool{"started": true, "completed": true}

// fieldUpdate is one --set field operation. hasValue distinguishes an explicit
// value (including the empty clear "") from a bare timestamp field (value=None
// in the oracle, auto-filled with now()).
type fieldUpdate struct {
	field    string
	value    string
	hasValue bool // false => bare timestamp field
}

// resolvedFields is the small insertion-order tracker the --set narration
// emits in (field: old -> new lines, JSON changes, quiet pairs). Reader-side
// ordering is the yaml.Node Content list; this struct's only job is to remember
// which --set updates actually resolved (a bare timestamp on an already-set
// field is skipped) and in what order they were requested.
type resolvedFields struct {
	order  []string
	values map[string]string
}

func (r *resolvedFields) set(key, val string) {
	if _, ok := r.values[key]; !ok {
		r.order = append(r.order, key)
	}
	r.values[key] = val
}

func (r *resolvedFields) get(key string) (string, bool) {
	v, ok := r.values[key]
	return v, ok
}

func (r *resolvedFields) keys() []string { return r.order }

// nowTimestamp returns the YYYY-MM-DDTHH:MM:SSZ stamp the oracle writes.
func nowTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

// updateFrontmatter parses the in-fence YAML body into a yaml.Node, mutates
// the target scalar node(s) or appends new K,V pairs for genuinely-missing
// fields, and re-marshals. The yaml.v3 emitter is canonical, so the contract
// is FIELD-exact, not byte-exact: unknown fields and key order survive
// --set/--archive intact (the AC-2 contract); the rewritten YAML body may
// renormalize visual whitespace (trailing whitespace, empty-value lines,
// pre-comment whitespace) that does not change any value. The 5 documented
// emitter normalizations are accepted observed behavior. Bare timestamp
// fields skip when already set. The write is atomic and preserves the file's
// EOF-newline state byte-for-byte: the splice replaces only the in-fence
// body, leaving the opening/close fence lines and the post-fence body —
// including the file's terminal newline state — untouched.
func updateFrontmatter(path string, updates []fieldUpdate) (*resolvedFields, error) {
	if !hasOpeningFence(path) {
		return nil, fmt.Errorf("No frontmatter found in %s", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Universal-newline normalize (so a CRLF entity comes out LF), preserve
	// the file's trailing-newline state through split/join on "\n".
	lines := strings.Split(normalizeNewlines(string(data)), "\n")

	fmStart, fmEnd := -1, -1
	inFM := false
	for i, line := range lines {
		if strings.TrimRight(line, " \t") == "---" {
			if !inFM {
				inFM = true
				fmStart = i
			} else {
				fmEnd = i
				break
			}
		}
	}
	if fmStart < 0 || fmEnd < 0 {
		return nil, fmt.Errorf("No frontmatter found in %s", path)
	}

	// Slice the in-fence YAML body and parse it into a yaml.Node we can
	// mutate. Unchanged scalars survive FIELD-exact through re-marshal — same
	// value, same key order — while the yaml.v3 emitter is free to re-emit
	// the same value in its canonical form (the 5 documented normalizations).
	// That is the seam's load-bearing property under AC-2's field-exact
	// contract; pinned by TestUpdateFrontmatterNodeRoundTrip + the 5
	// divergence tests.
	fmBody := strings.Join(lines[fmStart+1:fmEnd], "\n") + "\n"
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(fmBody), &doc); err != nil {
		return nil, fmt.Errorf("parse frontmatter %s: %w", path, err)
	}
	var mapping *yaml.Node
	if len(doc.Content) > 0 && doc.Content[0].Kind == yaml.MappingNode {
		mapping = doc.Content[0]
	} else {
		// Empty or non-mapping frontmatter — bootstrap a fresh mapping so we
		// can append the resolved updates as a brand-new K,V list.
		mapping = &yaml.Node{Kind: yaml.MappingNode}
		doc = yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{mapping}}
	}

	// Current values support skip-if-set for bare timestamp fields. Read the
	// top-level scalars off the node; nested mappings render as empty (the
	// indented-lines-ignored semantic).
	current := map[string]string{}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		k := mapping.Content[i]
		v := mapping.Content[i+1]
		if k.Kind != yaml.ScalarNode {
			continue
		}
		if v.Kind == yaml.ScalarNode {
			current[k.Value] = v.Value
		} else {
			current[k.Value] = ""
		}
	}

	now := nowTimestamp()
	resolved := &resolvedFields{values: map[string]string{}}
	for _, u := range updates {
		if !u.hasValue {
			if current[u.field] != "" {
				continue
			}
			resolved.set(u.field, now)
		} else {
			resolved.set(u.field, u.value)
		}
	}

	// Mutate matching scalar nodes; append new K,V pairs for genuinely
	// missing fields before re-marshaling.
	written := map[string]bool{}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		k := mapping.Content[i]
		v := mapping.Content[i+1]
		if k.Kind != yaml.ScalarNode {
			continue
		}
		val, ok := resolved.get(k.Value)
		if !ok {
			continue
		}
		setScalarValue(v, val)
		written[k.Value] = true
	}
	for _, key := range resolved.keys() {
		if written[key] {
			continue
		}
		val, _ := resolved.get(key)
		newKey := &yaml.Node{Kind: yaml.ScalarNode, Value: key}
		newVal := &yaml.Node{Kind: yaml.ScalarNode}
		setScalarValue(newVal, val)
		mapping.Content = append(mapping.Content, newKey, newVal)
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return nil, fmt.Errorf("marshal frontmatter %s: %w", path, err)
	}
	// yaml.Marshal terminates with "\n"; the in-fence splice expects each
	// rewritten line as its own element of `lines`, with the join restoring
	// the file's EOF state via the trailing empty element when present.
	newBody := strings.TrimRight(string(out), "\n")
	newLines := strings.Split(newBody, "\n")
	rebuilt := make([]string, 0, len(lines)-(fmEnd-fmStart-1)+len(newLines))
	rebuilt = append(rebuilt, lines[:fmStart+1]...)
	rebuilt = append(rebuilt, newLines...)
	rebuilt = append(rebuilt, lines[fmEnd:]...)

	if err := atomicWrite(path, []byte(strings.Join(rebuilt, "\n"))); err != nil {
		return nil, err
	}
	return resolved, nil
}

// setScalarValue writes val into a scalar node and picks its quoting style
// per the writer policy: a value YAML would otherwise parse as a mapping
// (`: `) or a comment (` #`, `\t#`, leading `#`) is wrapped in double quotes
// — the divergence #1 writer side plus the legacy option-C quoting — so the
// reader round-trips it whole. Other values let yaml.v3 auto-decide: a plain
// string emits plain, an existing scalar's style is RESET to the policy's
// choice so a mutation does not drag along a stale quote shape. The tag is
// cleared so yaml.v3 re-infers it from the value.
func setScalarValue(node *yaml.Node, val string) {
	node.Kind = yaml.ScalarNode
	node.Value = val
	node.Tag = ""
	if needsExplicitQuoting(val) {
		node.Style = yaml.DoubleQuotedStyle
	} else {
		node.Style = 0
	}
}

// needsExplicitQuoting is the writer policy: a value containing a colon-space
// (`: `), a space-or-tab-then-`#` (` #`, `\t#`), or starting with `#` is
// wrapped in double quotes. The colon-space case is divergence #1's writer
// side — yaml.v3 would otherwise parse the unquoted form as a nested mapping.
// The `#` cases preserve option-C round-trip: an unquoted whitespace-preceded
// `#` is treated as a comment by yaml.v3 (so the value would be truncated on
// a re-read), and a leading-`#` value is a comment-only line outright.
func needsExplicitQuoting(val string) bool {
	if val == "" {
		return false
	}
	if strings.Contains(val, ": ") {
		return true
	}
	if strings.Contains(val, " #") || strings.Contains(val, "\t#") {
		return true
	}
	if val[0] == '#' {
		return true
	}
	return false
}

// atomicWrite writes data to a temp file in the same directory and renames it
// into place, so a reader never observes a half-written entity. The principle
// stated in decision B applied to --set, --archive's stamp, and --new.
func atomicWrite(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".status-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return nil
}

// runArchive archives an entity (flat or folder form), stamping archived: before
// the move and printing `archived: {dest}`. Enforces the source-missing,
// already-archived, mod-block, and merge-hook guards. Matches run_archive.
// entityDir is the absolute entity root for I/O; definitionDir is the README root
// used to resolve merge-hook mods (definitionDir/_mods/); spellingDir is the
// as-passed spelling used for the printed dest (so a relative --workflow-dir
// renders a relative `archived:` path, matching the oracle's literal
// os.path.join).
func runArchive(definitionDir, entityDir, spellingDir, slug string, force, quiet, asJSON bool, stdout, stderr io.Writer) int {
	flatPath := filepath.Join(entityDir, slug+".md")
	folderRoot := filepath.Join(entityDir, slug)
	folderIndex := filepath.Join(folderRoot, "index.md")
	flatExists := isRegularFile(flatPath)
	folderExists := isRegularFile(folderIndex)

	var isFolder bool
	var sourcePath string
	switch {
	case folderExists:
		if flatExists {
			fmt.Fprintf(stderr,
				"Warning: entity '%s' has both %s and %s; preferring folder form. "+
					"Remove the flat file to silence this warning.\n",
				slug, flatPath, folderIndex)
		}
		isFolder = true
		sourcePath = folderIndex
	case flatExists:
		isFolder = false
		sourcePath = flatPath
	default:
		fmt.Fprintf(stderr, "Error: entity not found: %s\n", slug)
		return 1
	}

	fields := ParseFrontmatter(sourcePath)
	modBlock := strings.TrimSpace(fields["mod-block"])
	pr := strings.TrimSpace(fields["pr"])
	verdict := strings.TrimSpace(fields["verdict"])
	if modBlock != "" {
		if !force {
			fmt.Fprintf(stderr, "Error: entity %s has pending mod-block (%s). Use --force to override.\n", slug, modBlock)
			return 1
		}
		fmt.Fprintf(stderr, "Warning: --force overriding mod-block (%s) on entity %s\n", modBlock, slug)
	}

	// Verdict gate: archiving a TERMINAL entity finalizes it, so it requires a
	// verdict — the same gate the --set finalize path enforces (handlers.go
	// runSet), closed here so finalization-by-archive cannot route around it.
	// Scoped to entities whose status is a declared terminal stage (data-driven
	// over the README, not hardcoded 'done'); archiving a non-terminal entity is
	// not a finalize and is unaffected. --force bypasses.
	if verdict == "" && !force {
		readme := filepath.Join(definitionDir, "README.md")
		if fileExists(readme) {
			status := strings.TrimSpace(fields["status"])
			for _, s := range parseStagesBlock(readme) {
				if s.terminal && s.Name == status {
					fmt.Fprintf(stderr, "Error: entity %s cannot be archived from terminal stage '%s' without a verdict. Set verdict first, or use --force.\n", slug, status)
					return 1
				}
			}
		}
	}

	// Merge-hook invariant: archival is terminal. Refuse unless the hook ran
	// (pr set), is in flight (mod-block set, handled above), or --force. Under
	// `merge: local` the pr-requirement is exempted — the workflow declared it
	// merges locally, so an empty pr is expected at archival; the mod-block guard
	// above (policy-independent) still catches an in-flight ceremony. A
	// `verdict: rejected` entity is also exempted: it never ran the merge ceremony
	// (no PR to require, no merge to gate on), so the requirement is vacuous — this
	// matches the --set finalize path (runSet), keeping reject-then-archive on the
	// happy path without --force.
	policy, perr := resolveMergePolicy(definitionDir)
	if perr != nil {
		fmt.Fprintf(stderr, "Error: %s\n", perr)
		return 1
	}
	if !force && policy != mergeLocal && verdict != "rejected" && modBlock == "" && pr == "" {
		mergeHooks := scanMods(definitionDir)["merge"]
		if len(mergeHooks) > 0 {
			fmt.Fprintf(stderr,
				"Error: entity %s cannot be archived — workflow has merge hook(s) [%s] "+
					"that have not run (pr field is empty and mod-block is empty). "+
					"Invoke the hook first, or use --force to bypass.\n",
				slug, strings.Join(mergeHooks, ", "))
			return 1
		}
	}

	archiveDir := filepath.Join(entityDir, "_archive")
	var destPath, destSpelling string
	if isFolder {
		destPath = filepath.Join(archiveDir, slug)
		destSpelling = PyJoin(spellingDir, "_archive", slug)
		if fileExists(destPath) {
			fmt.Fprintf(stderr, "Error: already archived: %s/\n", slug)
			return 1
		}
	} else {
		destPath = filepath.Join(archiveDir, slug+".md")
		destSpelling = PyJoin(spellingDir, "_archive", slug+".md")
		if fileExists(destPath) {
			fmt.Fprintf(stderr, "Error: already archived: %s.md\n", slug)
			return 1
		}
	}

	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "Error: %s\n", err)
		return 1
	}

	if _, err := updateFrontmatter(sourcePath, []fieldUpdate{{field: "archived", value: nowTimestamp(), hasValue: true}}); err != nil {
		fmt.Fprintf(stderr, "Error: %s\n", err)
		return 1
	}

	var moveErr error
	if isFolder {
		moveErr = os.Rename(folderRoot, destPath)
	} else {
		moveErr = os.Rename(sourcePath, destPath)
	}
	if moveErr != nil {
		fmt.Fprintf(stderr, "Error: %s\n", moveErr)
		return 1
	}
	switch {
	case asJSON:
		emitJSON(stdout, newJSONObj().set("command", "archive").set("slug", slug))
	case quiet:
		fmt.Fprintf(stdout, "archived slug=%s\n", slug)
	default:
		fmt.Fprintf(stdout, "archived: %s\n", destSpelling)
	}
	return 0
}

// PyJoin concatenates path components the way Python's os.path.join does: it
// joins with the OS separator without cleaning a leading "." (so
// PyJoin(".", "_archive", "x.md") == "./_archive/x.md", unlike filepath.Join
// which would collapse it to "_archive/x.md"). An absolute later component
// resets the accumulated path, matching os.path.join.
func PyJoin(parts ...string) string {
	sep := string(filepath.Separator)
	result := ""
	for _, p := range parts {
		switch {
		case result == "":
			result = p
		case filepath.IsAbs(p):
			result = p
		case strings.HasSuffix(result, sep):
			result += p
		default:
			result += sep + p
		}
	}
	return result
}

// scanMods scans definitionDir/_mods/*.md for `## Hook:` headings, returning
// hookPoint -> sorted mod names. Mods are workflow definition (lifecycle hooks
// declared for the workflow, kin to the README stages), so they live next to the
// README in definitionDir/_mods/, not in the entity/state checkout. In
// single-root (definitionDir == entityDir) this is the same dir the oracle's
// scan_mods reads, so it matches byte-for-byte; under split-root it reads the
// definition dir, never the state checkout.
func scanMods(definitionDir string) map[string][]string {
	modsDir := filepath.Join(definitionDir, "_mods")
	info, err := os.Stat(modsDir)
	if err != nil || !info.IsDir() {
		return map[string][]string{}
	}
	entries, err := os.ReadDir(modsDir)
	if err != nil {
		return map[string][]string{}
	}
	var files []string
	for _, ent := range entries {
		if !ent.IsDir() && strings.HasSuffix(ent.Name(), ".md") {
			files = append(files, ent.Name())
		}
	}
	// glob returns sorted order.
	sort.Strings(files)
	hooks := map[string][]string{}
	for _, name := range files {
		modName := strings.TrimSuffix(name, ".md")
		data, err := os.ReadFile(filepath.Join(modsDir, name))
		if err != nil {
			continue
		}
		for _, line := range splitLines(string(data)) {
			if strings.HasPrefix(line, "## Hook:") {
				point := strings.TrimSpace(strings.TrimPrefix(line, "## Hook:"))
				if point != "" {
					hooks[point] = append(hooks[point], modName)
				}
			}
		}
	}
	for k := range hooks {
		sort.Strings(hooks[k])
	}
	return hooks
}
