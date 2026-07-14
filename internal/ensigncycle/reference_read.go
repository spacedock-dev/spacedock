package ensigncycle

import "strings"

// ReferenceReadSucceeded accepts a reference body only when the host reports a
// successful read and every document-specific anchor is present. Incidental prose
// that resembles a filesystem error is not a transport failure.
func ReferenceReadSucceeded(transportSucceeded bool, output string, anchors ...string) bool {
	if !transportSucceeded {
		return false
	}
	body := strings.TrimSpace(strings.ToLower(output))
	if body == "" {
		return false
	}
	for _, anchor := range anchors {
		if anchor == "" || !strings.Contains(body, strings.ToLower(anchor)) {
			return false
		}
	}
	return true
}
