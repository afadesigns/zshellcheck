package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1491(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — export PATH=/usr/bin",
			input:    `export PATH=/usr/bin`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — export LD_PRELOAD=/tmp/evil.so",
			input: `export LD_PRELOAD=/tmp/evil.so`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1491",
					Message: "`export LD_PRELOAD=...` forces every subsequent binary to load a custom library — classic privesc/persistence. Scope to a single invocation if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — export LD_LIBRARY_PATH=/opt/untrusted",
			input: `export LD_LIBRARY_PATH=/opt/untrusted`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1491",
					Message: "`export LD_LIBRARY_PATH=...` forces every subsequent binary to load a custom library — classic privesc/persistence. Scope to a single invocation if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — export LD_AUDIT=/tmp/audit.so",
			input: `export LD_AUDIT=/tmp/audit.so`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1491",
					Message: "`export LD_AUDIT=...` forces every subsequent binary to load a custom library — classic privesc/persistence. Scope to a single invocation if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1491")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
