package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1523(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — tar xf foo.tar -C /tmp/stage",
			input:    `tar xf foo.tar -C /tmp/stage`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — tar xf foo.tar -C /",
			input: `tar xf foo.tar -C /`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1523",
					Message: "`tar -C /` extracts into the filesystem root — overwrites any path that happens to be inside the archive. Stage, inspect, then copy.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1523")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
