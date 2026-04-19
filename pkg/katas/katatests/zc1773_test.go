package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1773(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `xargs -r rm` (guard present)",
			input:    `xargs -r rm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `xargs --no-run-if-empty rm`",
			input:    `xargs --no-run-if-empty rm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `xargs -0r rm` (combined short flags)",
			input:    `xargs -0r rm`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `xargs` alone (no command, probably a noop / listing)",
			input:    `xargs`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `xargs rm` (no guard)",
			input: `xargs rm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1773",
					Message: "`xargs` without `-r` / `--no-run-if-empty` runs the child once with no arguments when stdin is empty — a destructive surprise for `xargs rm`, `xargs kill`, etc. Add `-r` or switch to `find ... -exec cmd {} +`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `xargs -0 kill -9`",
			input: `xargs -0 kill -9`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1773",
					Message: "`xargs` without `-r` / `--no-run-if-empty` runs the child once with no arguments when stdin is empty — a destructive surprise for `xargs rm`, `xargs kill`, etc. Add `-r` or switch to `find ... -exec cmd {} +`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1773")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
