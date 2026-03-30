package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1174(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid paste with files",
			input:    `paste -sd, file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid paste -sd in pipeline",
			input: `paste -s -d,`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1174",
					Message: "Use Zsh `${(j:delim:)array}` to join array elements instead of `paste -sd`. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1174")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
