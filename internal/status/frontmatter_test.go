// ABOUTME: Frontmatter reader table tests — the yaml.v3-backed reader keeps
// ABOUTME: the convergent cases and surfaces the documented divergence on bad quotes.
package status

import (
	"reflect"
	"testing"
)

func TestContentHasOpeningFence(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"plain fence", "---\nid: 1\n---\n", true},
		{"leading blank lines skipped", "\n\n---\nid: 1\n---\n", true},
		{"whitespace-first-line disqualifies", "   \n---\n", false},
		{"no fence", "# heading\n", false},
		{"bom then fence", utf8BOM + "---\nid: 1\n---\n", true},
		{"empty file", "", false},
		{"text before fence", "hello\n---\n", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := contentHasOpeningFence([]byte(tc.in)); got != tc.want {
				t.Fatalf("contentHasOpeningFence(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestParseFrontmatterContent(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want map[string]string
	}{
		{
			name: "basic fields",
			in:   "---\nid: \"001\"\ntitle: Hello\nscore: 0.8\n---\nbody\n",
			want: map[string]string{"id": "001", "title": "Hello", "score": "0.8"},
		},
		{
			name: "empty value yields empty string",
			in:   "---\nscore:\n---\n",
			want: map[string]string{"score": ""},
		},
		{
			name: "matched single quotes stripped",
			in:   "---\ntitle: 'Quoted'\n---\n",
			want: map[string]string{"title": "Quoted"},
		},
		// Divergence #4 retired: the prior `mismatched quotes preserved` case
		// (Python line parser kept `"half'` literally) is dropped — the
		// yaml.v3-backed reader raises a parse error on a malformed quote,
		// which is the loud-failure mode we want. The migration check
		// confirms no live entity triggers it. See
		// TestParseFrontmatterMalformedQuoteIsDivergence below.
		{
			name: "nested indented lines ignored",
			in:   "---\nstages:\n  defaults:\n    worktree: false\nid: 1\n---\n",
			want: map[string]string{"stages": "", "id": "1"},
		},
		{
			name: "last top-level key wins",
			in:   "---\nstatus: a\nstatus: b\n---\n",
			want: map[string]string{"status": "b"},
		},
		{
			name: "no opening fence yields empty",
			in:   "# heading\nid: 1\n",
			want: map[string]string{},
		},
		{
			name: "value with colon splits on first colon only",
			in:   "---\nurl: http://x:8080\n---\n",
			want: map[string]string{"url": "http://x:8080"},
		},
		{
			name: "leading bom on first key line",
			in:   utf8BOM + "---\nid: 1\n---\n",
			want: map[string]string{"id": "1"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseFrontmatterContent([]byte(tc.in))
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("parseFrontmatterContent(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// TestParseFrontmatterMalformedQuoteIsDivergence pins divergence #4: the
// yaml.v3-backed reader raises a parse error on a mismatched / unterminated
// quote in a value (e.g. `title: "half'`), and the public reader surfaces
// that as an empty field map rather than silently preserving the malformed
// bytes (the retired Python-line-parser quirk). The migration check asserts
// no live entity trips this; the surface here is the unit-test contract.
func TestParseFrontmatterMalformedQuoteIsDivergence(t *testing.T) {
	got := parseFrontmatterContent([]byte("---\ntitle: \"half'\n---\n"))
	if len(got) != 0 {
		t.Fatalf("malformed-quote input should surface as empty (parse error), got %v", got)
	}
}
