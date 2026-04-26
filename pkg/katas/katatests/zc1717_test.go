// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1717(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker pull` (no bypass)",
			input:    `docker pull nginx:1.27`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker push` (no bypass)",
			input:    `docker push myorg/app:1.2.3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker version --disable-content-trust` (not pull/push subcmd)",
			input:    `docker version --disable-content-trust`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker pull --disable-content-trust`",
			input: `docker pull --disable-content-trust myorg/app:1.2.3`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1717",
					Message: "`docker pull --disable-content-trust` overrides `DOCKER_CONTENT_TRUST=1` — unsigned image moves into the registry or local store. Sign the artifact (`docker trust sign`) instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `docker push --disable-content-trust`",
			input: `docker push --disable-content-trust myorg/app:1.2.3`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1717",
					Message: "`docker push --disable-content-trust` overrides `DOCKER_CONTENT_TRUST=1` — unsigned image moves into the registry or local store. Sign the artifact (`docker trust sign`) instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1717")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
