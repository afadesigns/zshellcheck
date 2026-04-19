package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1769(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `vagrant destroy` (prompt kept)",
			input:    `vagrant destroy`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `vagrant halt`",
			input:    `vagrant halt`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `vagrant destroy --force`",
			input: `vagrant destroy --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1769",
					Message: "`vagrant destroy --force` skips the prompt and drops the VM (and any un-exported data inside). Drop the flag, or use `vagrant halt` + `vagrant up` to cycle without destroy.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `vagrant destroy -f myvm`",
			input: `vagrant destroy -f myvm`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1769",
					Message: "`vagrant destroy -f` skips the prompt and drops the VM (and any un-exported data inside). Drop the flag, or use `vagrant halt` + `vagrant up` to cycle without destroy.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1769")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
