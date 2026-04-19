package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1826(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `install -m 0755 src /usr/local/bin/app`",
			input:    `install -m 0755 src /usr/local/bin/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `install -d /opt/app` (no mode)",
			input:    `install -d /opt/app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `install -m 4755 src /usr/local/bin/myapp`",
			input: `install -m 4755 src /usr/local/bin/myapp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1826",
					Message: "`install -m 4755` drops a setuid/setgid binary in one step. If DEST is on `$PATH`, every local user can invoke the elevated binary. Only install trusted builds, and prefer narrow `setcap` capabilities over setuid.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `install -m 2755 src /usr/local/bin/myapp`",
			input: `install -m 2755 src /usr/local/bin/myapp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1826",
					Message: "`install -m 2755` drops a setuid/setgid binary in one step. If DEST is on `$PATH`, every local user can invoke the elevated binary. Only install trusted builds, and prefer narrow `setcap` capabilities over setuid.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1826")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
