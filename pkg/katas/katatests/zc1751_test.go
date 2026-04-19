package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1751(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `rpm -e libfoo`",
			input:    `rpm -e libfoo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `dnf remove libfoo`",
			input:    `dnf remove libfoo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `rpm -q libfoo --nodeps` (query, not erase)",
			input:    `rpm -q libfoo --nodeps`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `rpm -e --nodeps libfoo`",
			input: `rpm -e --nodeps libfoo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1751",
					Message: "`rpm ... --nodeps` removes the package without the dependency solver — dependents break (libc, openssl, systemd units). Resolve the conflict explicitly instead of bypassing.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `dnf remove --nodeps libfoo`",
			input: `dnf remove --nodeps libfoo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1751",
					Message: "`dnf ... --nodeps` removes the package without the dependency solver — dependents break (libc, openssl, systemd units). Resolve the conflict explicitly instead of bypassing.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1751")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
