// ABOUTME: per-project review mode (local or remote) kept outside the repository,
// ABOUTME: and the remote-review dispatch that refuses to upload without it.
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

// reviewSettingPath keys the setting by project root under the user's config
// directory, so a choice is never committed and never inherited through a clone.
func reviewSettingPath(env []string, dir string) (project, path string) {
	project = dir
	if out, err := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel").Output(); err == nil {
		project = strings.TrimSpace(string(out))
	}
	vars := envMap(env)
	base := vars["XDG_CONFIG_HOME"]
	if base == "" {
		base = filepath.Join(vars["HOME"], ".config")
	}
	sum := sha256.Sum256([]byte(project))
	return project, filepath.Join(base, "spacedock", "projects", hex.EncodeToString(sum[:8])+".json")
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
	setting.Project = project
	return setting, path, nil
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
			if err != nil {
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
				if err := os.MkdirAll(filepath.Dir(path), 0o700); err == nil {
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

func newRemoteReviewCommand(env []string, dir string, stderr io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:                "remote-review <file> [publish flags]",
		Short:              "Publish an artifact to the relay for review, if this project allows it",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			setting, _, err := readReviewSetting(env, dir)
			if err != nil {
				fmt.Fprintln(stderr, "spacedock remote-review:", err)
				return exitCodeError{1}
			}
			if setting.Mode != "remote" {
				fmt.Fprintf(stderr, "spacedock remote-review: %s has not chosen remote review, so nothing was uploaded.\n"+
					"Ask the user first; if they agree, run `spacedock review-mode set remote` and retry.\n", setting.Project)
				return exitCodeError{3}
			}
			bin, err := exec.LookPath(extensionPrefix + "remote-review")
			if err != nil {
				fmt.Fprintln(stderr, "spacedock remote-review: spacedock-remote-review is not installed; reinstall spacedock to get Subspace")
				return exitCodeError{127}
			}
			code, err := execHost{}.Launch(append([]string{bin, "publish"}, args...), env)
			if err != nil {
				fmt.Fprintln(stderr, "spacedock remote-review:", err)
			}
			if code != 0 {
				return exitCodeError{code}
			}
			return nil
		},
	}
}
