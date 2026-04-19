package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1905(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ssh -L 8080:target:80 host` (no -g)",
			input:    `ssh -L 8080:target:80 host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ssh -g host` (no forward)",
			input:    `ssh -g host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ssh -g -L 8080:target:80 host`",
			input: `ssh -g -L 8080:target:80 host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1905",
					Message: "`ssh -g` with `-L`/`-D` binds the forward on `0.0.0.0` — anyone on the same LAN segment can ride the tunnel. Drop `-g` or pin `bind_address:port` in the forward spec.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ssh -gD 1080 host`",
			input: `ssh -gD 1080 host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1905",
					Message: "`ssh -g` with `-L`/`-D` binds the forward on `0.0.0.0` — anyone on the same LAN segment can ride the tunnel. Drop `-g` or pin `bind_address:port` in the forward spec.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1905")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
