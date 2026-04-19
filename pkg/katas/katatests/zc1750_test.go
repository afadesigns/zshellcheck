package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1750(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl proxy --port 8001` (loopback default)",
			input:    `kubectl proxy --port 8001`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl proxy --address 127.0.0.1`",
			input:    `kubectl proxy --address 127.0.0.1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl proxy --address 0.0.0.0`",
			input: `kubectl proxy --address 0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1750",
					Message: "`kubectl proxy --address 0.0.0.0` exposes the cluster-admin API tunnel to every reachable interface. Keep the loopback default and tunnel over SSH, or restrict `--address` to a firewalled interface.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kubectl proxy --address=0.0.0.0`",
			input: `kubectl proxy --address=0.0.0.0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1750",
					Message: "`kubectl proxy --address=0.0.0.0` exposes the cluster-admin API tunnel to every reachable interface. Keep the loopback default and tunnel over SSH, or restrict `--address` to a firewalled interface.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1750")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
