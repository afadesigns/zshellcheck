package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1566(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — gem install rails",
			input:    `gem install rails`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — gem install -P HighSecurity",
			input:    `gem install -P HighSecurity rails`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — gem install -P NoSecurity",
			input: `gem install -P NoSecurity rails`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1566",
					Message: "`gem -P NoSecurity` skips signature verification — MITM or account compromise becomes RCE at install. Use HighSecurity.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — gem install --trust-policy LowSecurity",
			input: `gem install --trust-policy LowSecurity rails`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1566",
					Message: "`gem -P LowSecurity` skips signature verification — MITM or account compromise becomes RCE at install. Use HighSecurity.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1566")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
