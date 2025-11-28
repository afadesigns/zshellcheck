package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1039(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid usage",
			input:    `rm /tmp/file`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid usage recursive",
			input:    `rm -rf /tmp/dir`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rm root",
			input: `rm /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1039",
					Message: "Avoid `rm` on the root directory `/`. This is highly dangerous.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "invalid rm root quoted",
			input: `rm "/"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1039",
					Message: "Avoid `rm` on the root directory `/`. This is highly dangerous.",
					Line:    1,
					Column:  4,
				},
			},
		},
		{
			name:  "invalid rm root single quoted",
			input: `rm '/'`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1039",
					Message: "Avoid `rm` on the root directory `/`. This is highly dangerous.",
					Line:    1,
					Column:  4,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1039")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
