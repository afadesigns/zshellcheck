package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1265(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid systemctl enable --now",
			input:    `systemctl enable --now nginx`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid systemctl start",
			input:    `systemctl start nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid systemctl enable without --now",
			input: `systemctl enable nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1265",
					Message: "Use `systemctl enable --now` to enable and start the service immediately. Without `--now`, the service only starts on next boot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1265")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
