package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1668(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aws iam attach-user-policy ReadOnlyAccess",
			input:    `aws iam attach-user-policy --user-name foo --policy-arn arn:aws:iam::aws:policy/ReadOnlyAccess`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — aws iam create-access-key",
			input:    `aws iam create-access-key --user-name foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — attach-user-policy AdministratorAccess",
			input: `aws iam attach-user-policy --user-name foo --policy-arn arn:aws:iam::aws:policy/AdministratorAccess`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1668",
					Message: "`aws iam attach-user-policy ... arn:aws:iam::aws:policy/AdministratorAccess` grants sweeping admin — use a scoped inline policy (`put-user-policy`) or a customer-managed policy with the minimum `Action`/`Resource` set.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — attach-role-policy PowerUserAccess",
			input: `aws iam attach-role-policy --role-name r --policy-arn arn:aws:iam::aws:policy/PowerUserAccess`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1668",
					Message: "`aws iam attach-role-policy ... arn:aws:iam::aws:policy/PowerUserAccess` grants sweeping admin — use a scoped inline policy (`put-user-policy`) or a customer-managed policy with the minimum `Action`/`Resource` set.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1668")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
