package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1459(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run --cap-add=NET_BIND_SERVICE",
			input:    `docker run --cap-add=NET_BIND_SERVICE alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --cap-add=SYS_ADMIN (equals form)",
			input: `docker run --cap-add=SYS_ADMIN alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1459",
					Message: "Dangerous Linux capability granted — breaks the container's security boundary. Prefer `--cap-drop=ALL` and add back only minimum needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker run --cap-add SYS_PTRACE (space form)",
			input: `docker run --cap-add SYS_PTRACE alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1459",
					Message: "Dangerous Linux capability granted — breaks the container's security boundary. Prefer `--cap-drop=ALL` and add back only minimum needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run --cap-add=ALL",
			input: `podman run --cap-add=ALL alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1459",
					Message: "Dangerous Linux capability granted — breaks the container's security boundary. Prefer `--cap-drop=ALL` and add back only minimum needed.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1459")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
