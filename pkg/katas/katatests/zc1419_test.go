package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1419(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — chmod 755",
			input:    `chmod 755 script.sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chmod 777",
			input: `chmod 777 file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1419",
					Message: "Avoid `chmod 777`/`a+rwx` — grants world-writable access. Prefer restrictive modes (755, 644, 700, 600) matched to the actual file purpose.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chmod a+rwx",
			input: `chmod a+rwx dir`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1419",
					Message: "Avoid `chmod 777`/`a+rwx` — grants world-writable access. Prefer restrictive modes (755, 644, 700, 600) matched to the actual file purpose.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1419")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
