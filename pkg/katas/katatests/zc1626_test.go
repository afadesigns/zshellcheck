// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1626(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — --set non-secret",
			input:    `helm install myapp chart --set replicas=3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — -f values.yaml",
			input:    `helm install myapp chart -f /secure/values.yaml`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — --set-file points at path",
			input:    `helm install myapp chart --set-file db.password=/run/secrets/db`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — helm install --set password=...",
			input: `helm install myapp chart --set password=s3cret`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1626",
					Message: "`helm install --set password=s3cret` places a secret value in argv — readable via `ps`. Use `-f values.yaml` or `--set-file password=PATH`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — helm upgrade --set-string token=$TOKEN",
			input: `helm upgrade myapp chart --set-string token=$TOKEN`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1626",
					Message: "`helm upgrade --set-string token=$TOKEN` places a secret value in argv — readable via `ps`. Use `-f values.yaml` or `--set-file token=PATH`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1626")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
