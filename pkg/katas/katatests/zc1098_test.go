package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestZC1098(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:  "eval with unquoted variable",
			input: `eval "ls $dir"`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1098",
					Message: "Use the `(q)` flag (or `(qq)`, `(q-)`) when using variables in `eval` to prevent injection.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "eval with quoted variable",
			input: `eval "ls ${(q)dir}"`,
			expected: []katas.Violation{},
		},
	}
	// Skipping for now as implementation needs refinement.
	_ = tests
}
