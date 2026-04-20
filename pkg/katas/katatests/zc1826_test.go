package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1826(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `install -m 0755 src /usr/local/bin/app`",
			input:    `install -m 0755 src /usr/local/bin/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `install -d /opt/app` (no mode)",
			input:    `install -d /opt/app`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `install -m 4755 …` (numeric setuid is owned by ZC1892)",
			input:    `install -m 4755 src /usr/local/bin/app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `install -m u+s src /usr/local/bin/app`",
			input: `install -m u+s src /usr/local/bin/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1826",
					Message: "`install -m u+s` applies a symbolic setuid/setgid bit — easy to miss in review. Use `0755` and grant narrow caps with `setcap` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `install -m ug+s src /usr/local/bin/app`",
			input: `install -m ug+s src /usr/local/bin/app`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1826",
					Message: "`install -m ug+s` applies a symbolic setuid/setgid bit — easy to miss in review. Use `0755` and grant narrow caps with `setcap` instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1826")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
