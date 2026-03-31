package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1202(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ip addr",
			input:    `ip addr show`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid ifconfig",
			input: `ifconfig eth0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1202",
					Message: "Avoid `ifconfig` — it is deprecated on modern Linux. Use `ip addr`, `ip link`, or `ip route` from iproute2.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1202")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
