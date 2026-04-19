package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1794(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `cosign verify registry.example.com/app:1.2.3`",
			input:    `cosign verify registry.example.com/app:1.2.3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `cosign sign --key cosign.key registry.example.com/app:1.2.3`",
			input:    `cosign sign --key cosign.key registry.example.com/app:1.2.3`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `cosign verify --insecure-ignore-tlog img`",
			input: `cosign verify --insecure-ignore-tlog img`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1794",
					Message: "`cosign --insecure-ignore-tlog` removes a rung of the signature chain (transparency log / SCT / TLS / HTTPS-only registry). Drop the flag and fix the underlying trust anchor.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `cosign sign --allow-insecure-registry img`",
			input: `cosign sign --allow-insecure-registry img`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1794",
					Message: "`cosign --allow-insecure-registry` removes a rung of the signature chain (transparency log / SCT / TLS / HTTPS-only registry). Drop the flag and fix the underlying trust anchor.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1794")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
