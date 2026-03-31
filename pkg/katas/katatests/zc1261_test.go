package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1261(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid base64 encode",
			input:    `base64 file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid base64 -d",
			input: `base64 -d encoded.txt`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1261",
					Message: "Inspect `base64 -d` output before piping to execution. Blindly executing decoded content is a code injection vector.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1261")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
