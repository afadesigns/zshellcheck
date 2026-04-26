// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1959(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `trivy image alpine:3.20`",
			input:    `trivy image alpine:3.20`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `trivy image --download-db-only`",
			input:    `trivy image --download-db-only`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `trivy image alpine:3.20 --skip-db-update`",
			input: `trivy image alpine:3.20 --skip-db-update`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1959",
					Message: "`trivy --skip-db-update` scans against the cached DB — every CVE disclosed since last refresh is missed. Keep the default download, or run `trivy --download-db-only` once per day in a scheduled job.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `trivy image alpine:3.20 --skip-update`",
			input: `trivy image alpine:3.20 --skip-update`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1959",
					Message: "`trivy --skip-update` scans against the cached DB — every CVE disclosed since last refresh is missed. Keep the default download, or run `trivy --download-db-only` once per day in a scheduled job.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1959")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
