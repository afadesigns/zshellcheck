package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1682(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — npm ci (no unsafe-perm)",
			input:    `npm ci`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — yarn install",
			input:    `yarn install`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — npm install --unsafe-perm",
			input: `npm install --unsafe-perm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1682",
					Message: "`npm --unsafe-perm` keeps root for every lifecycle script — a compromised dep executes as root. Build in a dedicated builder container or run as a non-root user.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — npm install --unsafe-perm=true",
			input: `npm install --unsafe-perm=true`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1682",
					Message: "`npm --unsafe-perm=true` keeps root for every lifecycle script — a compromised dep executes as root. Build in a dedicated builder container or run as a non-root user.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1682")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
