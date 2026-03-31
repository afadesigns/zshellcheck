package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1226(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid dmesg -T",
			input:    `dmesg -T`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid dmesg without -T",
			input: `dmesg -l err`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1226",
					Message: "Use `dmesg -T` for human-readable timestamps instead of raw kernel boot-seconds. Or use `--time-format=iso` for ISO 8601.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1226")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
