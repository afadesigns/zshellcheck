package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1215(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid source os-release",
			input:    `source /etc/os-release`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid cat os-release",
			input: `cat /etc/os-release`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1215",
					Message: "Source `/etc/os-release` directly with `. /etc/os-release` instead of parsing with `cat`. It exports variables like `$ID` and `$VERSION_ID`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1215")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
