package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1203(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ss",
			input:    `ss -tulnp`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid netstat",
			input: `netstat -tulnp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1203",
					Message: "Avoid `netstat` — it is deprecated on modern Linux. Use `ss` from iproute2 for faster, more detailed socket statistics.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1203")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
