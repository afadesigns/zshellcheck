package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1649(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — -days 90 (Let's Encrypt style)",
			input:    `openssl req -x509 -days 90 -nodes`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — -days 398 (1-year max)",
			input:    `openssl req -x509 -days 398 -nodes`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — -days 3650 (10 years)",
			input: `openssl req -x509 -days 3650 -nodes -newkey rsa:2048`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1649",
					Message: "`openssl req -days 3650` issues a cert with a long validity. Keep leaf certs under 398 days and automate rotation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — openssl x509 -days 1095 (3 years)",
			input: `openssl x509 -req -days 1095 -signkey key -in csr`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1649",
					Message: "`openssl x509 -days 1095` issues a cert with a long validity. Keep leaf certs under 398 days and automate rotation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1649")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
