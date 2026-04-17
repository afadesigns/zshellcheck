package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1458(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run --user nobody",
			input:    `docker run --user 1000 alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --user root",
			input: `docker run --user root alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1458",
					Message: "Explicit root UID inside a container lets container-escape bugs become host root. Use a non-root USER in the image.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker run --user 0",
			input: `docker run --user 0 alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1458",
					Message: "Explicit root UID inside a container lets container-escape bugs become host root. Use a non-root USER in the image.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1458")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
