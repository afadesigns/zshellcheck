// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1738(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `aws rds delete-db-instance` with explicit final snapshot",
			input:    `aws rds delete-db-instance --db-instance-identifier mydb --final-db-snapshot-identifier mydb-final`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws rds describe-db-instances`",
			input:    `aws rds describe-db-instances`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `aws rds delete-db-instance ... --skip-final-snapshot`",
			input: `aws rds delete-db-instance --db-instance-identifier mydb --skip-final-snapshot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1738",
					Message: "`aws rds delete-db-instance --skip-final-snapshot` deletes the database with no recovery snapshot. Drop the flag or pass `--final-db-snapshot-identifier <name>` so the snapshot is explicit and verifiable.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `aws rds delete-db-cluster ... --skip-final-snapshot`",
			input: `aws rds delete-db-cluster --db-cluster-identifier mycluster --skip-final-snapshot`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1738",
					Message: "`aws rds delete-db-cluster --skip-final-snapshot` deletes the database with no recovery snapshot. Drop the flag or pass `--final-db-snapshot-identifier <name>` so the snapshot is explicit and verifiable.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1738")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
