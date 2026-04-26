// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1217(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid systemctl",
			input:    `systemctl restart nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid service command",
			input: `service nginx restart`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1217",
					Message: "Avoid `service` — it is a SysVinit compatibility wrapper. Use `systemctl` directly on systemd systems.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1217")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
