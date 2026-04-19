package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1805(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `aws dynamodb describe-table --table-name mytbl`",
			input:    `aws dynamodb describe-table --table-name mytbl`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `aws cloudformation list-stacks`",
			input:    `aws cloudformation list-stacks`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `aws cloudformation delete-stack --stack-name prod`",
			input: `aws cloudformation delete-stack --stack-name prod`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1805",
					Message: "`aws cloudformation delete-stack` removes every resource the stack manages, no rollback. Stage a confirmation, pin IDs via `--cli-input-json`, and export a backup first where the service supports one.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `aws kms schedule-key-deletion --key-id k`",
			input: `aws kms schedule-key-deletion --key-id k`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1805",
					Message: "`aws kms schedule-key-deletion` queues CMK deletion — ciphertext becomes unreadable after the grace window. Stage a confirmation, pin IDs via `--cli-input-json`, and export a backup first where the service supports one.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1805")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
