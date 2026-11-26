package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type installProgress struct {
	elapsedTrace
}

type installSpan struct {
	progress *installProgress
	label    string
	started  time.Time
}

func (p *installProgress) start(label string) installSpan {
	if p == nil {
		return installSpan{}
	}
	now := p.now()
	p.at(now, "Start: "+label)
	return installSpan{progress: p, label: label, started: now}
}

func (s installSpan) end(err error, continues bool) {
	if s.progress == nil {
		return
	}
	status, detail := "Done", ""
	if err != nil {
		status = "Failed"
		if continues {
			status = "Warning"
		}
		var childExit *exec.ExitError
		if errors.As(err, &childExit) {
			if code := childExit.ExitCode(); code >= 0 {
				detail = fmt.Sprintf("; child exit %d", code)
			} else {
				status, detail = "Interrupted", "; child ended by signal"
			}
		}
		if continues {
			detail += "; continuing"
		}
	}
	s.outcome(status, detail)
}

func (s installSpan) outcome(status, detail string) {
	if s.progress == nil {
		return
	}
	now := s.progress.now()
	s.progress.at(now, fmt.Sprintf("%s: %s (%s%s)", status, s.label, durationSeconds(now.Sub(s.started)), detail))
}

func (p *installProgress) skip(label, reason string) {
	if p != nil {
		p.at(p.now(), fmt.Sprintf("Skipped: %s (%s)", label, reason))
	}
}

func (p *installProgress) finish(code int) {
	if p != nil {
		now := p.now()
		p.at(now, fmt.Sprintf("Total: %s (exit %d)", durationSeconds(now.Sub(p.zero)), code))
	}
}

func runInstallCommand(ctx context.Context, args, env []string, dir string, stdout, stderr io.Writer, now func() time.Time) int {
	setup, code := parsePiSetupArgs("install", args, stderr)
	if code != 0 {
		return code
	}
	var progress *installProgress
	if setup.verbose {
		progress = &installProgress{newElapsedTrace(stderr, now)}
	}
	code = runInstallSetup(ctx, args, setup, execHost{installProgress: progress}, execPiRuntimeOps{}, env, stdout, stderr, progress)
	if code == 0 {
		if versionsWanted(args) {
			collect := progress.start("Collect installed versions")
			reportInstalledVersions(env, dir, stdout, progress)
			collect.end(nil, false)
		} else {
			progress.skip("Collect installed versions", "not used for this install path")
		}
	}
	if code != 2 {
		progress.finish(code)
	}
	return code
}
