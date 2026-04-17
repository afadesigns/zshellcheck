package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1435(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — killall plain",
			input:    `killall myproc`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — killall -9",
			input: `killall -9 myproc`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1435",
					Message: "`killall -9 name` force-kills every matching process, including unrelated instances on multi-user or containerized hosts. Start with -TERM, or kill by PID after `pgrep`/`pidof`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1435")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
