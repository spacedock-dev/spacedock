// ABOUTME: The Spacedock release lookup behind the front door's daily update hint:
// ABOUTME: the installer-matched User-Agent, allowed redirects and answer parsing.
package updatecheck

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// maxAnswer bounds a release answer; a larger body is refused, not truncated.
const maxAnswer = 64 * 1024

// Match install.sh lookup_user_agent's accepted version and marker vocabulary.
var versionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$`)

// UserAgent builds `spacedock/<version> (<os>; <arch>[; ci][; internal])` with
// install.sh's version check, CI truthiness and internal marker.
func UserAgent(version, goos, goarch string, getenv func(string) string) string {
	version = strings.TrimPrefix(version, "v")
	if len(version) > 32 || !versionPattern.MatchString(version) {
		version = "unknown"
	}
	ua := fmt.Sprintf("spacedock/%s (%s; %s", version, goos, goarch)
	switch strings.ToLower(strings.Trim(getenv("CI"), " \t\n\r\v\f")) {
	case "1", "true", "yes", "on":
		ua += "; ci"
	}
	if getenv("SPACEDOCK_INTERNAL") == "1" {
		ua += "; internal"
	}
	return ua + ")"
}

// Latest returns the newest Spacedock release on channel ("stable" or "edge"),
// or "" on any failure. It follows only the route's api.github.com redirect and
// sends only the edge release selector, never an identity.
func Latest(ctx context.Context, client *http.Client, channel, ua string) string {
	c := *client
	c.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	target := "https://spacedock.md/latest-release/spacedock"
	if channel == "edge" {
		target += "?prerelease=1"
	}
	for redirects := 0; redirects < 5; redirects++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
		if err != nil {
			return ""
		}
		req.Header.Set("User-Agent", ua)
		res, err := c.Do(req)
		if err != nil {
			return ""
		}
		if res.StatusCode >= 300 && res.StatusCode < 400 {
			location, err := url.Parse(res.Header.Get("Location"))
			res.Body.Close()
			if err != nil || !allowedRedirect(location, channel) {
				return ""
			}
			target = location.String()
			continue
		}
		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			return ""
		}
		data, err := io.ReadAll(io.LimitReader(res.Body, maxAnswer+1))
		res.Body.Close()
		if err != nil || len(data) > maxAnswer {
			return ""
		}
		return newestRelease(data, channel)
	}
	return ""
}

// GitHub answers a release object for /releases/latest (stable) and a
// one-element array for /releases?per_page=1 (edge).
func newestRelease(data []byte, channel string) string {
	type releaseAnswer struct {
		Tag        string `json:"tag_name"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
	}
	var answers []releaseAnswer
	if json.Unmarshal(data, &answers) != nil {
		var answer releaseAnswer
		if json.Unmarshal(data, &answer) != nil {
			return ""
		}
		answers = []releaseAnswer{answer}
	}
	for _, answer := range answers {
		if answer.Draft || (channel == "stable" && answer.Prerelease) {
			continue
		}
		if version := releaseVersion(answer.Tag, channel); version != "" {
			return version
		}
	}
	return ""
}

// The route sends `spacedock/` clients only to Spacedock's GitHub API releases.
// The path is compared as sent: RawPath is set for non-canonical escaping.
func allowedRedirect(u *url.URL, channel string) bool {
	if u.Scheme != "https" || u.User != nil || u.Fragment != "" || u.RawPath != "" || u.ForceQuery {
		return false
	}
	if u.Host != "api.github.com" {
		return false
	}
	const base = "/repos/spacedock-dev/spacedock/releases"
	if channel == "edge" {
		return u.Path == base && u.RawQuery == "per_page=1"
	}
	return u.Path == base+"/latest" && u.RawQuery == ""
}

func releaseVersion(tag, channel string) string {
	if !strings.HasPrefix(tag, "v") {
		return ""
	}
	v := strings.TrimPrefix(tag, "v")
	if !versionPattern.MatchString(v) || strings.Contains(v, "+") || (channel == "stable" && strings.Contains(v, "-")) {
		return ""
	}
	return v
}
