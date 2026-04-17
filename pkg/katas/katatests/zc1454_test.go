package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1454(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run alpine",
			input:    `docker run alpine echo hi`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --privileged",
			input: `docker run --privileged alpine sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1454",
					Message: "`--privileged` disables container isolation — effectively host root. Use `--cap-add` + `--device` for narrow permissions.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1454")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
