package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1930(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `setopt HASH_CMDS` (explicit default)",
			input:    `setopt HASH_CMDS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `unsetopt NOMATCH` (unrelated)",
			input:    `unsetopt NOMATCH`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unsetopt HASH_CMDS`",
			input: `unsetopt HASH_CMDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1930",
					Message: "`unsetopt HASH_CMDS` re-walks `$PATH` on every call — tens to hundreds of ms per command on slow filesystems. Keep it on; use `rehash` or `hash -r` to invalidate the cache after a targeted binary swap.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `setopt NO_HASH_CMDS`",
			input: `setopt NO_HASH_CMDS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1930",
					Message: "`setopt NO_HASH_CMDS` re-walks `$PATH` on every call — tens to hundreds of ms per command on slow filesystems. Keep it on; use `rehash` or `hash -r` to invalidate the cache after a targeted binary swap.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1930")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
