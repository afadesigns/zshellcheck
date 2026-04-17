package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1542(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — snap install firefox",
			input:    `snap install firefox`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — snap install --dangerous local.snap",
			input: `snap install --dangerous ./local.snap`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1542",
					Message: "`snap install --dangerous` installs an assertion-unverified snap — any .snap on disk can register system services. Use --devmode or the store.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1542")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
