package dispatch

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// semanticName reserves eight characters for the host's retry/cycle suffix.
func semanticName(slug, stage string) (string, error) {
	for _, part := range []string{slug, stage} {
		if !namePattern.MatchString(part) && !(len(part) == 1 && strings.Contains("abcdefghijklmnopqrstuvwxyz0123456789", part)) {
			return "", fmt.Errorf("invalid name component %q: must match %s", part, namePattern)
		}
	}
	if len(slug)+1+len(stage) <= 56 {
		return slug + "-" + stage, nil
	}
	budget := 56 - len(stage) - 10
	if budget < 8 {
		return "", fmt.Errorf("stage %q leaves insufficient readable name budget", stage)
	}
	head := strings.TrimRight(slug[:budget], "-")
	if len(head) < 8 {
		return "", fmt.Errorf("slug %q leaves insufficient readable name budget", slug)
	}
	digest := sha256.Sum256([]byte(slug))
	return fmt.Sprintf("%s-%x-%s", head, digest[:4], stage), nil
}

func validWorkerSuffix(s string) bool {
	if len(s) > 8 {
		return false
	}
	return s == "" || s == "-retry" || (strings.HasPrefix(s, "-cycle") && isAllDigits(s[6:])) || (strings.HasPrefix(s, "-") && isAllDigits(s[1:]))
}

// resolveWorker enumerates both generations before interpreting suffixes. A
// collision, including two possible stage boundaries, never acquires ownership.
func resolveWorker(name string, stages []string, records ...map[string]entityRecord) decomposeResult {
	matches := map[decomposeResult]bool{}
	for _, records := range records {
		for slug, rec := range records {
			for _, stage := range stages {
				semantic, _ := semanticName(slug, stage)
				candidates := []string{semantic, "spacedock-ensign-" + slug + "-" + stage}
				for _, style := range []string{"sd-b32", "sequential", "slug"} {
					candidates = append(candidates, capWorkerName("spacedock-ensign", slug, stage, rec.id, style))
				}
				for _, base := range candidates {
					if base == "" || !strings.HasPrefix(name, base) {
						continue
					}
					suffix := strings.TrimPrefix(name, base)
					if validWorkerSuffix(suffix) {
						matches[decomposeResult{slug: slug, stage: stage, cycle: strings.TrimPrefix(suffix, "-"), ok: true}] = true
					}
				}
			}
		}
	}
	return uniqueWorker(matches)
}

func uniqueWorker(matches map[decomposeResult]bool) decomposeResult {
	if len(matches) == 1 {
		for d := range matches {
			return d
		}
	}
	return decomposeResult{}
}

func validateWorkerName(workflowDir, entityPath, stage string) (string, error) {
	slug := status.EntitySlug(entityPath)
	name, err := semanticName(slug, stage)
	if err != nil {
		return "", err
	}
	root := splitRootStateCheckout(workflowDir)
	if root == "" {
		root = workflowDir
	}
	active, archived := loadEntityFrontmatter(root), loadEntityFrontmatter(filepath.Join(root, "_archive"))
	active[slug] = entityRecord{slug: slug, id: status.ParseFrontmatter(entityPath)["id"]}
	stages := readStageNames(workflowDir)
	stages = append(stages, stage)
	d := resolveWorker(name, stages, active, archived)
	if !d.ok || d.slug != slug || d.stage != stage {
		return "", fmt.Errorf("ambiguous generated worker name %q; choose a distinct task slug", name)
	}
	return name, nil
}

// runName exposes canonical identity without build, artifact or stamp effects.
func runName(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("dispatch name", flag.ContinueOnError)
	flags.SetOutput(stderr)
	wd := flags.String("workflow-dir", "", "workflow definition directory")
	ep := flags.String("entity-path", "", "canonical entity file")
	stage := flags.String("stage", "", "declared workflow stage")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 || *wd == "" || *ep == "" || *stage == "" {
		return buildError(stderr, 2, "dispatch name requires --workflow-dir, --entity-path and --stage, with no positional arguments")
	}
	entityPath, err := filepath.Abs(*ep)
	if err != nil {
		return buildError(stderr, 1, "%v", err)
	}
	if strings.Contains(filepath.ToSlash(entityPath), "/.worktrees/") {
		return buildError(stderr, 1, "entity_path must refer to the project-root entity, not a worktree copy")
	}
	if _, err := os.ReadFile(entityPath); err != nil {
		return buildError(stderr, 1, "entity file not readable: %v", err)
	}
	readme, err := os.ReadFile(filepath.Join(*wd, "README.md"))
	if err != nil {
		return buildError(stderr, 1, "workflow README not readable: %v", err)
	}
	stages, _ := status.ParseStagesWithDefaultsData(readme)
	declared := false
	for _, s := range stages {
		if s.Name == *stage {
			declared = true
		}
	}
	if !declared {
		return buildError(stderr, 1, "stage %q not declared in workflow", *stage)
	}
	name, err := validateWorkerName(*wd, entityPath, *stage)
	if err != nil {
		return buildError(stderr, 1, "%v", err)
	}
	fmt.Fprintln(stdout, name)
	return 0
}
