package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1084(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "quoted glob",
			input:    `find . -name "*.txt"`,
			expected: []katas.Violation{},
		},
		{
			name:     "single quoted glob",
			input:    `find . -name '*.txt'`,
			expected: []katas.Violation{},
		},
				{
					name:  "unquoted star glob",
					input: `find . -name *.txt`,
					expected: []katas.Violation{
						{
							KataID:  "ZC1084",
							Message: "Quote globs in `find` commands. `(*)txt` is subject to shell expansion.",
							Line:    1,
							Column:  14,
						},
					},
				},
				{
					name:  "unquoted question glob",
					input: `find . -name file?.txt`,
					expected: []katas.Violation{
						{
							KataID:  "ZC1084",
							Message: "Quote globs in `find` commands. `file(?)txt` is subject to shell expansion.",
							Line:    1,
							Column:  14,
						},
					},
				},
				{
					name:     "unquoted bracket glob (merged)",
					input:    `find . -name [a-z]`,
					expected: []katas.Violation{
						{
							KataID:  "ZC1084",
							Message: "Quote globs in `find` commands. `-(name[a-z])` contains unquoted brackets.",
							Line:    1,
							Column:  8, // Points to -name
						},
					},
				},
				{
					name:     "unquoted bracket glob (partial)",
					input:    `find . -name file[a-z]`,
					expected: []katas.Violation{
						{
							KataID:  "ZC1084",
							Message: "Quote globs in `find` commands. `file[a-z]` is subject to shell expansion.",
							Line:    1,
							Column:  18, // Points to [
						},
					},
				},
				{
					name:     "escaped glob",
					input:    `find . -name \*.txt`,
					expected: []katas.Violation{},
				},
				{
					name:     "escaped question",
					input:    `find . -name file\?.txt`,
					expected: []katas.Violation{},
				},
				{
					name:     "double backslash (escaped backslash + glob)",
					input:    `find . -name \\*.txt`,
					expected: []katas.Violation{
						{
							KataID:  "ZC1084",
							Message: "Quote globs in `find` commands. `\\\\(*)txt` is subject to shell expansion.",
							Line:    1,
							Column:  14,
						},
					},
				},
		
		{
			name:     "other flag (ignore)",
			input:    `find . -type f`,
			expected: []katas.Violation{},
		},
		{
			name:     "find with exec (ignore args)",
			input:    `find . -exec echo {} \;`,
			expected: []katas.Violation{},
		},
		{
			name:     "quoted bracket glob",
			input:    `find . -name '[a-z]'`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vioations := testutil.Check(tt.input, "ZC1084")
			testutil.AssertViolations(t, tt.input, vioations, tt.expected)
		})
	}
}
