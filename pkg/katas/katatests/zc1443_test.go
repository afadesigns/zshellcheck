package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1443(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — terraform plan",
			input:    `terraform plan`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — terraform destroy",
			input: `terraform destroy`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1443",
					Message: "`terraform destroy` without `-target` removes every resource in state. Scope with `-target=...` or gate behind interactive confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1443")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
