package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1436(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl -p (reload from config)",
			input:    `sysctl -p`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl -w (ephemeral)",
			input: `sysctl -w vm.swappiness=10`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1436",
					Message: "`sysctl -w` setting is lost on reboot. Persist in `/etc/sysctl.d/*.conf` and reload with `sysctl --system`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1436")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
