package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1258(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid rsync --delete",
			input:    `rsync -az --delete src/ dst/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid rsync single file",
			input:    `rsync -az file host:path`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rsync dir without --delete",
			input: `rsync -az src/ dst/`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1258",
					Message: "Consider `rsync --delete` for directory sync. Without `--delete`, files removed from source remain on the destination.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1258")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
