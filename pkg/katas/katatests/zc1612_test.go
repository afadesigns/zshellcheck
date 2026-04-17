package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1612(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sysctl tightening ptrace_scope",
			input:    `sysctl -w kernel.yama.ptrace_scope=3`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — unrelated sysctl",
			input:    `sysctl -w vm.swappiness=10`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sysctl -w kernel.yama.ptrace_scope=0",
			input: `sysctl -w kernel.yama.ptrace_scope=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1612",
					Message: "`sysctl ... kernel.yama.ptrace_scope=0` disables YAMA ptrace scope (lets any process attach) — defense-in-depth loss. Leave the default unless a measured need justifies it.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sysctl -w kernel.kptr_restrict=0",
			input: `sysctl -w kernel.kptr_restrict=0`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1612",
					Message: "`sysctl ... kernel.kptr_restrict=0` disables kernel pointer restriction (leaks kptrs to /proc) — defense-in-depth loss. Leave the default unless a measured need justifies it.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1612")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
