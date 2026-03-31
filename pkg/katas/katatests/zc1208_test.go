package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1208(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid visudo -c",
			input:    `visudo -c -f /etc/sudoers.d/myconfig`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid visudo",
			input: `visudo -f /etc/sudoers`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1208",
					Message: "Avoid `visudo` in scripts — it opens an interactive editor. Write to `/etc/sudoers.d/` drop-in files and validate with `visudo -c`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1208")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
