package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1931(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ip netns list`",
			input:    `ip netns list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ip netns exec red ping host`",
			input:    `ip netns exec red ping host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ip netns delete red`",
			input: `ip netns delete red`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1931",
					Message: "`ip netns delete` tears down every interface, veth, tunnel, and WireGuard peer inside the namespace. Stop the workloads first and verify `ip -n $NS link` is empty before deleting.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ip netns del $NS`",
			input: `ip netns del $NS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1931",
					Message: "`ip netns del` tears down every interface, veth, tunnel, and WireGuard peer inside the namespace. Stop the workloads first and verify `ip -n $NS link` is empty before deleting.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1931")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
