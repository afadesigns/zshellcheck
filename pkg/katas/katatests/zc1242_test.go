package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1242(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid tar with -C",
			input:    `tar xzf archive.tar.gz -C /opt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid tar create",
			input:    `tar czf archive.tar.gz dir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid tar extract without -C",
			input: `tar xzf archive.tar.gz`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1242",
					Message: "Use `tar -C dir` to specify extraction directory. Without `-C`, tar extracts into the current directory which may overwrite files.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1242")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
