// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1455(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run isolated",
			input:    `docker run -p 8080:80 alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --net=host",
			input: `docker run --net=host alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1455",
					Message: "`--net=host` / `--network=host` lets the container reach host-local services. Use `-p` for explicit publishes or dedicated container networks.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1455")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
