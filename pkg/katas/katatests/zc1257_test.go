package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1257(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid docker stop -t",
			input:    `docker stop -t 5 mycontainer`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid docker stop without -t",
			input: `docker stop mycontainer`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1257",
					Message: "Use `docker stop -t N` to set an explicit shutdown timeout. The default 10s may be too long or too short for your use case.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1257")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
