package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1608(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find -exec rm {}",
			input:    `find . -type f -exec rm {} \;`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — find -exec sh -c 'positional'",
			input:    `find . -exec sh -c 'grep pat "$1"' _ {} \;`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — find -exec sh -c 'echo {}'`,
			input: `find . -exec sh -c 'echo {}' \;`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1608",
					Message: "`find -exec sh -c '... {} ...'` interpolates filenames into the shell script — metacharacters break out. Pass `{}` as a positional arg: `sh -c '... \"$1\"' _ {} \\;`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — find -exec bash -c "grep X {}"`,
			input: `find . -exec bash -c "grep X {}" \;`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1608",
					Message: "`find -exec sh -c '... {} ...'` interpolates filenames into the shell script — metacharacters break out. Pass `{}` as a positional arg: `sh -c '... \"$1\"' _ {} \\;`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1608")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
