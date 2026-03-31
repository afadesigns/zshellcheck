package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1214(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid sudo -u",
			input:    `sudo -u postgres psql`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid su",
			input: `su -c "service nginx restart" www-data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1214",
					Message: "Avoid `su` in scripts — it prompts for a password interactively. Use `sudo -u user cmd` for non-interactive privilege switching.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1214")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
