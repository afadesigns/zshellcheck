package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1259(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid docker pull with tag",
			input:    `docker pull alpine:3.19`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid docker pull without tag",
			input: `docker pull nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1259",
					Message: "Pin Docker image to a specific tag instead of defaulting to `:latest`. Untagged pulls are non-reproducible and may break unexpectedly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1259")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
