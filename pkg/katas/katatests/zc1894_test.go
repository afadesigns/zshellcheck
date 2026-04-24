package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1894(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `conntrack -L` (list)",
			input:    `conntrack -L`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `conntrack -D -s 10.0.0.5` (narrow delete)",
			input:    `conntrack -D -s 10.0.0.5`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `conntrack -F`",
			input: `conntrack -F`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1894",
					Message: "`conntrack -F` wipes every tracked flow — stateful `ctstate ESTABLISHED` allowances drop, running SSH sessions lose their entry. Gate with `at now + N min` or narrow to one flow with `conntrack -D -s <ip>`.",
					Line:    1,
					Column:  11,
				},
			},
		},
		{
			name:  "invalid — `conntrack --flush conntrack`",
			input: `conntrack --flush conntrack`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1894",
					Message: "`conntrack -F` wipes every tracked flow — stateful `ctstate ESTABLISHED` allowances drop, running SSH sessions lose their entry. Gate with `at now + N min` or narrow to one flow with `conntrack -D -s <ip>`.",
					Line:    1,
					Column:  12,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1894")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
