// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1961(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `gcloud iam service-accounts list`",
			input:    `gcloud iam service-accounts list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `gcloud auth print-access-token` (short-lived)",
			input:    `gcloud auth print-access-token`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `gcloud iam service-accounts keys create key.json --iam-account=$SA`",
			input: `gcloud iam service-accounts keys create key.json --iam-account=$SA`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1961",
					Message: "`gcloud iam service-accounts keys create` mints a long-lived JSON key — no auto-rotate, no refresh. Prefer Workload Identity Federation, `--impersonate-service-account`, or the attached service account.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:     "valid — `gcloud iam service-accounts keys list` (read only)",
			input:    `gcloud iam service-accounts keys list`,
			expected: []katas.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1961")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
