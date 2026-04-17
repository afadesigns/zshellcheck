package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1460(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run with no security-opt",
			input:    `docker run alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --security-opt=no-new-privileges",
			input:    `docker run --security-opt=no-new-privileges alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --security-opt seccomp=unconfined (space form)",
			input: `docker run --security-opt seccomp=unconfined alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1460",
					Message: "Disabling seccomp or AppArmor removes the main syscall/MAC filter that blocks container escapes. Keep the default profile.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — docker run --security-opt=apparmor=unconfined (equals form)",
			input: `docker run --security-opt=apparmor=unconfined alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1460",
					Message: "Disabling seccomp or AppArmor removes the main syscall/MAC filter that blocks container escapes. Keep the default profile.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1460")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
