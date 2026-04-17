package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1505(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dpkg -i pkg.deb",
			input:    `dpkg -i pkg.deb`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dpkg -i --force-confnew pkg.deb",
			input: `dpkg -i --force-confnew pkg.deb`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1505",
					Message: "`--force-confnew` silently picks maintainer or local conffile — legit /etc changes disappear or new defaults are ignored. Use ucf/etckeeper.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — dpkg -i pkg.deb --force-confold",
			input: `dpkg -i pkg.deb --force-confold`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1505",
					Message: "`--force-confold` silently picks maintainer or local conffile — legit /etc changes disappear or new defaults are ignored. Use ucf/etckeeper.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1505")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
