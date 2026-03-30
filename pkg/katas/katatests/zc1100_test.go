package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1100(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid basename with suffix flag",
			input:    `basename -s .txt file.txt`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid basename with multiple args",
			input:    `basename /path/to/file .txt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid simple dirname",
			input: `dirname /path/to/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1100",
					Message: "Use `${var%/*}` instead of `dirname` to extract the directory path. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid simple basename",
			input: `basename /path/to/file`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1100",
					Message: "Use `${var##*/}` instead of `basename` to extract the filename. Parameter expansion avoids spawning an external process.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1100")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
