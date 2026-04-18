package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1685(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sleep 30",
			input:    `sleep 30`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — sleep variable",
			input:    `sleep $timeout`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sleep infinity",
			input: `sleep infinity`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1685",
					Message: "`sleep infinity` does not trap SIGTERM — the orchestrator hangs until SIGKILL. Use `exec tail -f /dev/null` or front with `tini` / `dumb-init`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1685")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
