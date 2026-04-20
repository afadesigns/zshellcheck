package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1958(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `helm upgrade myapp bitnami/nginx`",
			input:    `helm upgrade myapp bitnami/nginx`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `helm upgrade myapp ./chart --atomic --wait`",
			input:    `helm upgrade myapp ./chart --atomic --wait`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `helm upgrade myapp ./chart --force`",
			input: `helm upgrade myapp ./chart --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1958",
					Message: "`helm upgrade --force` is delete+create — pods die, PodDisruptionBudget is bypassed, Services reset their `clusterIP`. Use plain `helm upgrade` (three-way merge) or `--atomic`/`--wait` for a supervised roll.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `helm3 upgrade --install myapp ./chart --force`",
			input: `helm3 upgrade --install myapp ./chart --force`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1958",
					Message: "`helm upgrade --force` is delete+create — pods die, PodDisruptionBudget is bypassed, Services reset their `clusterIP`. Use plain `helm upgrade` (three-way merge) or `--atomic`/`--wait` for a supervised roll.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1958")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
