package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1873(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt ERR_RETURN` (explicit default)",
			input:    `unsetopt ERR_RETURN`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt NOMATCH` (unrelated)",
			input:    `setopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt ERR_RETURN`",
			input: `setopt ERR_RETURN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1873",
					Message: "`setopt ERR_RETURN` returns from every function on first non-zero exit — probing helpers (`test -f`, `grep -q`) bail before the fallback branch. Scope inside a `LOCAL_OPTIONS` function if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_ERR_RETURN`",
			input: `unsetopt NO_ERR_RETURN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1873",
					Message: "`unsetopt NO_ERR_RETURN` returns from every function on first non-zero exit — probing helpers (`test -f`, `grep -q`) bail before the fallback branch. Scope inside a `LOCAL_OPTIONS` function if needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1873")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
