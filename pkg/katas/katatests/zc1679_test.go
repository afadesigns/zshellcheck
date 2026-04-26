// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1679(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — gcloud add-iam-policy-binding roles/viewer",
			input:    `gcloud projects add-iam-policy-binding PROJ --member=user:a@ex.com --role=roles/viewer`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — gcloud iam service-accounts create",
			input:    `gcloud iam service-accounts create foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — add-iam-policy-binding roles/owner",
			input: `gcloud projects add-iam-policy-binding PROJ --member=user:a@ex.com --role=roles/owner`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1679",
					Message: "`gcloud ... add-iam-policy-binding --role=roles/owner` grants primitive / IAM-admin — use a predefined role with the minimum scope or a custom role, and apply admin changes via Terraform.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — add-iam-policy-binding roles/iam.serviceAccountTokenCreator",
			input: `gcloud projects add-iam-policy-binding PROJ --member=user:a@ex.com --role=roles/iam.serviceAccountTokenCreator`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1679",
					Message: "`gcloud ... add-iam-policy-binding --role=roles/iam.serviceAccountTokenCreator` grants primitive / IAM-admin — use a predefined role with the minimum scope or a custom role, and apply admin changes via Terraform.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1679")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
