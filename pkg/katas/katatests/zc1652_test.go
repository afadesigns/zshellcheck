package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1652(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — ssh without X11",
			input:    `ssh user@host`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — ssh -X (untrusted)",
			input:    `ssh -X user@host firefox`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — ssh -Y user@host",
			input: `ssh -Y user@host xclock`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1652",
					Message: "`ssh -Y` enables trusted X11 forwarding — remote clients get full access to the local X server. Use `-X` (untrusted) or drop X11 forwarding entirely.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — ssh -i key -Y user@host",
			input: `ssh -i key -Y user@host`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1652",
					Message: "`ssh -Y` enables trusted X11 forwarding — remote clients get full access to the local X server. Use `-X` (untrusted) or drop X11 forwarding entirely.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1652")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
