package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1606(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — mkdir -m 700",
			input:    `mkdir -m 700 /root/data`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mkdir -m 755",
			input:    `mkdir -m 755 /opt/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — mkdir -m 1777 (sticky)",
			input:    `mkdir -m 1777 /tmp/shared`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — mkdir -m 777",
			input: `mkdir -m 777 /tmp/shared`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1606",
					Message: "`mkdir -m 777` creates a world-writable path without the sticky bit — TOCTOU symlink-attack ground. Use `-m 700` / `-m 750`, or `-m 1777` if a shared sticky dir is actually needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — install -m 666",
			input: `install -m 666 file /tmp/x`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1606",
					Message: "`install -m 666` creates a world-writable path without the sticky bit — TOCTOU symlink-attack ground. Use `-m 700` / `-m 750`, or `-m 1777` if a shared sticky dir is actually needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1606")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
