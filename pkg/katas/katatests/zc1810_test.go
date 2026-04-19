package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1810(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `wget https://example.com/file.tar.gz`",
			input:    `wget https://example.com/file.tar.gz`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `wget -r --level=2 https://example.com/`",
			input:    `wget -r --level=2 https://example.com/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `wget -r -l3 https://example.com/`",
			input:    `wget -r -l3 https://example.com/`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `wget -r https://example.com/`",
			input: `wget -r https://example.com/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1810",
					Message: "`wget -r` / `--mirror` without `--level=N` follows links to arbitrary depth — the crawl can exhaust disk and climb into parent paths. Pin `--level=3`, add `--no-parent`, and cap with `--quota=1G`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `wget --mirror https://example.com/`",
			input: `wget --mirror https://example.com/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1810",
					Message: "`wget -r` / `--mirror` without `--level=N` follows links to arbitrary depth — the crawl can exhaust disk and climb into parent paths. Pin `--level=3`, add `--no-parent`, and cap with `--quota=1G`.",
					Line:    1,
					Column:  8,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1810")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
