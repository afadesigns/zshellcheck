package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1531(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — wget -t 5 https://host",
			input:    `wget -t 5 https://host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — wget -t 0 https://host",
			input: `wget -t 0 https://host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1531",
					Message: "`wget -t 0` retries forever — script hangs on dead endpoint. Use finite `-t 5` plus `--timeout=<seconds>`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1531")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
