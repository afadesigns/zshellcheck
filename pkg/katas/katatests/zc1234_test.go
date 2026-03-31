package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1234(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid docker run --rm",
			input:    `docker run --rm alpine echo hello`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid docker run -d",
			input:    `docker run -d nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid docker run without --rm",
			input: `docker run alpine echo hello`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1234",
					Message: "Use `docker run --rm` to auto-remove containers after exit. Without `--rm`, stopped containers accumulate on disk.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1234")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
