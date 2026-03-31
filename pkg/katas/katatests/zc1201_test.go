package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1201(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid ssh",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid rsh",
			input: `rsh host command`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1201",
					Message: "Avoid `rsh` — it is an insecure legacy protocol. Use `ssh`/`scp`/`rsync` for encrypted remote operations.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid rcp",
			input: `rcp file host:/path`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1201",
					Message: "Avoid `rcp` — it is an insecure legacy protocol. Use `ssh`/`scp`/`rsync` for encrypted remote operations.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1201")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
