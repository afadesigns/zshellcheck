// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1575(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aws configure set region us-east-1",
			input:    `aws configure set region us-east-1`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — aws configure set aws_secret_access_key VALUE",
			input: `aws configure set aws_secret_access_key AKIAEXAMPLEKEYXYZ`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1575",
					Message: "`aws configure set aws_secret_access_key …` puts the secret in ps / history. Use IAM-role auth or import from stdin / 0600 file.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — aws configure set aws_session_token",
			input: `aws configure set aws_session_token FwoGZXIvYXdzEXAMPLE`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1575",
					Message: "`aws configure set aws_session_token …` puts the secret in ps / history. Use IAM-role auth or import from stdin / 0600 file.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1575")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
