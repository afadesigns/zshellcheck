package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1262(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid chmod -R 755",
			input:    `chmod -R 755 /var/www`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid chmod 777 non-recursive",
			input:    `chmod 777 /tmp/test`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid chmod -R 777",
			input: `chmod -R 777 /var/www`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1262",
					Message: "Never use `chmod -R 777` — it makes everything world-writable. Use `find -type d -exec chmod 755` and `find -type f -exec chmod 644` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1262")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
