package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1856(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unset arr` (delete whole variable)",
			input:    `unset arr`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unset FOO` (scalar)",
			input:    `unset FOO`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unset arr[0]`",
			input: `unset arr[0]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1856",
					Message: "`unset (arr[0])` is a Bash idiom — in Zsh it tries to unset a parameter literally named `(arr[0])` and leaves the array untouched. Use `arr[N]=()` or rebuild with `arr=(\"${(@)arr:#pattern}\")`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unset myarray[3]`",
			input: `unset myarray[3]`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1856",
					Message: "`unset (myarray[3])` is a Bash idiom — in Zsh it tries to unset a parameter literally named `(myarray[3])` and leaves the array untouched. Use `arr[N]=()` or rebuild with `arr=(\"${(@)arr:#pattern}\")`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1856")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
