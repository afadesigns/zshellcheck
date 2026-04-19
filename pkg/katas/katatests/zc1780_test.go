package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1780(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `sysctl fs.protected_symlinks=1`",
			input:    `sysctl fs.protected_symlinks=1`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `sysctl -a` (list all)",
			input:    `sysctl -a`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `sysctl -w fs.protected_symlinks=0`",
			input: `sysctl -w fs.protected_symlinks=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1780",
					Message: "`sysctl fs.protected_symlinks=0` disables symlink follow protection in sticky dirs — re-opens a TOCTOU race in sticky dirs. Leave the default unless you have a specific reason; otherwise scope the change to a mount namespace.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `sysctl fs.protected_hardlinks=0`",
			input: `sysctl fs.protected_hardlinks=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1780",
					Message: "`sysctl fs.protected_hardlinks=0` disables hardlink creation protection in sticky dirs — re-opens a TOCTOU race in sticky dirs. Leave the default unless you have a specific reason; otherwise scope the change to a mount namespace.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1780")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
