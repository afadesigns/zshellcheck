package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1669(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — git gc default",
			input:    `git gc`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — git gc --auto",
			input:    `git gc --auto`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — git gc --prune=now",
			input: `git gc --prune=now --aggressive`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1669",
					Message: "`git gc --prune=now` bulldozes the reflog / prune recovery window — keep the default cadence unless you are actively purging leaked secrets, and mirror the dropped history off-box first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — git reflog expire --expire=now --all",
			input: `git reflog expire --expire=now --all`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1669",
					Message: "`git reflog expire --expire=now` bulldozes the reflog / prune recovery window — keep the default cadence unless you are actively purging leaked secrets, and mirror the dropped history off-box first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1669")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
