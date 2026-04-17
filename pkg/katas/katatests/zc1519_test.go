package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1519(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ulimit -u 4096",
			input:    `ulimit -u 4096`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ulimit -n unlimited (different limit)",
			input:    `ulimit -n unlimited`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ulimit -u unlimited",
			input: `ulimit -u unlimited`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1519",
					Message: "`ulimit -u unlimited` removes the user process cap — fork bomb surface. Pick a realistic number or set it via /etc/security/limits.d/.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1519")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
