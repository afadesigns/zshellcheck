package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1239(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid kubectl exec without -it",
			input:    `kubectl exec mypod -- ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid kubectl exec -it",
			input: `kubectl exec -it mypod -- bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1239",
					Message: "Avoid `kubectl exec -it` in scripts — TTY allocation hangs without a terminal. Use `kubectl exec pod -- cmd` for non-interactive execution.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1239")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
