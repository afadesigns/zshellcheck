package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1560(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — pip install foo",
			input:    `pip install foo`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — pip install --index-url https://x foo",
			input:    `pip install --index-url https://pypi.example.com foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — pip install --trusted-host pypi.example.com foo",
			input: `pip install --trusted-host pypi.example.com foo`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1560",
					Message: "`pip --trusted-host` skips TLS verification and allows plain-HTTP for that index. Fix the CA trust and keep --index-url on https://.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — pip3 install foo --trusted-host pypi.org",
			input: `pip3 install foo --trusted-host pypi.org`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1560",
					Message: "`pip --trusted-host` skips TLS verification and allows plain-HTTP for that index. Fix the CA trust and keep --index-url on https://.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1560")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
