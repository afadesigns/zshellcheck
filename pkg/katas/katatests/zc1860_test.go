package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1860(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `hostnamectl status`",
			input:    `hostnamectl status`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `hostname -f` (read-only query)",
			input:    `hostname -f`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `hostnamectl set-hostname worker-42`",
			input: `hostnamectl set-hostname worker-42`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1860",
					Message: "`hostnamectl set-hostname worker-42` updates the kernel hostname live, but running services keep the old `gethostname()` — syslog tags, Prometheus labels, TLS SANs stay stale. Apply at provisioning or reboot.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `hostname worker-42`",
			input: `hostname worker-42`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1860",
					Message: "`hostname worker-42` updates the kernel hostname live, but running services keep the old `gethostname()` — syslog tags, Prometheus labels, TLS SANs stay stale. Apply at provisioning or reboot.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1860")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
