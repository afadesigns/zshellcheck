// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1777(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cp lib.so /usr/local/lib/` (unrelated file)",
			input:    `cp lib.so /usr/local/lib/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cat /etc/ld.so.preload` (read only)",
			input:    `cat /etc/ld.so.preload`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `tee /etc/ld.so.preload`",
			input: `tee /etc/ld.so.preload`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1777",
					Message: "`tee /etc/ld.so.preload` writes `/etc/ld.so.preload` — linker force-loads each listed library into every process. Audit for unexpected entries; for a scoped preload use `LD_PRELOAD=` on a single invocation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cp /tmp/x.so /etc/ld.so.preload`",
			input: `cp /tmp/x.so /etc/ld.so.preload`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1777",
					Message: "`cp /etc/ld.so.preload` writes `/etc/ld.so.preload` — linker force-loads each listed library into every process. Audit for unexpected entries; for a scoped preload use `LD_PRELOAD=` on a single invocation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1777")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
