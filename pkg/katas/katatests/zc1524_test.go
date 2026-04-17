package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1524(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl -p",
			input:    `sysctl -p /etc/sysctl.d/99-hardening.conf`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl -e -p",
			input: `sysctl -e -p /etc/sysctl.d/99-hardening.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1524",
					Message: "`sysctl -e` suppresses error output — typos in sysctl.d/ conffiles silently skip. Remove and surface the real error.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sysctl -q",
			input: `sysctl -q -p /etc/sysctl.d/99-hardening.conf`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1524",
					Message: "`sysctl -q` suppresses error output — typos in sysctl.d/ conffiles silently skip. Remove and surface the real error.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1524")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
