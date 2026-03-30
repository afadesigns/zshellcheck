package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1167(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "invalid timeout command",
			input: `timeout 5 cmd`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1167",
					Message: "Avoid `timeout` — it's unavailable on macOS. Use Zsh `TMOUT` variable or `zmodload zsh/sched` for portable timeout functionality.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1167")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
