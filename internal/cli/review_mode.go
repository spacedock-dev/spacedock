// ABOUTME: per-project review mode (local or remote), kept beside Subspace's relay
// ABOUTME: device credentials and outside the repository.
package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type reviewSetting struct {
	Project   string `json:"project"`
	Mode      string `json:"mode"`
	DecidedAt string `json:"decided_at,omitempty"`
}

// reviewSettingPath keys the setting by project root beside Subspace's relay
// device credentials ($XDG_STATE_HOME/subspace), outside the repository, so a
// choice is never committed and never inherited through a clone.
func reviewSettingPath(env []string, dir string) (project, path string) {
	project = dir
	if out, err := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel").Output(); err == nil {
		project = strings.TrimSpace(string(out))
	}
	vars := envMap(env)
	base := vars["XDG_STATE_HOME"]
	if base == "" {
		base = filepath.Join(vars["HOME"], ".local", "state")
	}
	sum := sha256.Sum256([]byte(project))
	return project, filepath.Join(base, "subspace", "review-mode", hex.EncodeToString(sum[:8])+".json")
}

func readReviewSetting(env []string, dir string) (reviewSetting, string, error) {
	project, path := reviewSettingPath(env, dir)
	setting := reviewSetting{Project: project}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return setting, path, nil
	}
	if err != nil {
		return setting, path, err
	}
	if err := json.Unmarshal(data, &setting); err != nil {
		return reviewSetting{Project: project}, path, fmt.Errorf("unreadable review setting %s: %w", path, err)
	}
	if setting.Mode != "" && setting.Mode != "local" && setting.Mode != "remote" {
		return reviewSetting{Project: project}, path, fmt.Errorf("unreadable review setting %s: unknown mode %q", path, setting.Mode)
	}
	setting.Project = project
	return setting, path, nil
}

// reviewChoiceCheck is the review choice's row in the doctor version table. A choice
// that cannot be read is a problem, because the review skill stops on it.
func reviewChoiceCheck(in versionInputs) ([]pieceRow, []string) {
	setting, path, err := readReviewSetting(in.env, in.dir)
	switch {
	case err != nil:
		return []pieceRow{{"review choice", "unreadable", path}},
			[]string{"UNREADABLE: the review choice for this directory. Run `spacedock review-mode clear` to be asked again."}
	case setting.Mode == "":
		return []pieceRow{{"review choice", "not chosen", "for " + setting.Project}}, nil
	}
	return []pieceRow{{"review choice", setting.Mode, "chosen " + setting.DecidedAt + " for " + setting.Project}}, nil
}

func newReviewModeCommand(env []string, dir string, stdout, stderr io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:                "review-mode get [--json] | set local|remote | clear",
		Short:              "Show or record how this project's artifacts are reviewed",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || wantsHelp(args) {
				fmt.Fprintln(stdout, "Usage: spacedock review-mode get [--json] | set local|remote | clear")
				return nil
			}
			setting, path, err := readReviewSetting(env, dir)
			if err != nil && args[0] == "get" {
				fmt.Fprintln(stderr, "spacedock review-mode:", err)
				return exitCodeError{1}
			}
			switch {
			case args[0] == "get":
				if len(args) > 1 && args[1] == "--json" {
					out, _ := json.Marshal(setting)
					fmt.Fprintln(stdout, string(out))
				} else if setting.Mode == "" {
					fmt.Fprintf(stdout, "no review mode chosen for %s\n", setting.Project)
				} else {
					fmt.Fprintf(stdout, "%s (chosen %s) for %s\n", setting.Mode, setting.DecidedAt, setting.Project)
				}
				return nil
			case args[0] == "set" && len(args) == 2 && (args[1] == "local" || args[1] == "remote"):
				setting.Mode = args[1]
				setting.DecidedAt = time.Now().UTC().Format(time.RFC3339)
				out, _ := json.Marshal(setting)
				err = os.MkdirAll(filepath.Dir(path), 0o700)
				if err == nil {
					err = os.WriteFile(path, append(out, '\n'), 0o600)
				}
				if err != nil {
					fmt.Fprintln(stderr, "spacedock review-mode:", err)
					return exitCodeError{1}
				}
				fmt.Fprintf(stdout, "review mode for %s is now %s\n", setting.Project, setting.Mode)
				return nil
			case args[0] == "clear" && len(args) == 1:
				if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
					fmt.Fprintln(stderr, "spacedock review-mode:", err)
					return exitCodeError{1}
				}
				fmt.Fprintf(stdout, "review mode for %s cleared\n", setting.Project)
				return nil
			}
			fmt.Fprintln(stderr, "Usage: spacedock review-mode get [--json] | set local|remote | clear")
			return exitCodeError{2}
		},
	}
}
