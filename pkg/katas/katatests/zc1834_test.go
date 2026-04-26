// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1834(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `tc qdisc add dev eth0 root netem loss 5%` (partial)",
			input:    `tc qdisc add dev eth0 root netem loss 5%`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `tc qdisc del dev eth0 root` (cleanup)",
			input:    `tc qdisc del dev eth0 root`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tc qdisc add … netem loss 100%`",
			input: `tc qdisc add dev eth0 root netem loss 100%`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1834",
					Message: "`tc qdisc add … netem loss 100%` black-holes every packet on the target interface — remote SSH dies instantly. Stage on a secondary dev or wrap in a timed recovery (`at now + N minutes … tc qdisc del …`).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `tc qdisc replace … netem corrupt 100%`",
			input: `tc qdisc replace dev eth0 root netem corrupt 100%`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1834",
					Message: "`tc qdisc replace … netem corrupt 100%` black-holes every packet on the target interface — remote SSH dies instantly. Stage on a secondary dev or wrap in a timed recovery (`at now + N minutes … tc qdisc del …`).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1834")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
