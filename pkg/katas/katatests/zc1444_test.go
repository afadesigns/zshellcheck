package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1444(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — redis-cli GET",
			input:    `redis-cli GET foo`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — redis-cli FLUSHALL",
			input: `redis-cli FLUSHALL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1444",
					Message: "`redis-cli FLUSHALL`/`FLUSHDB` wipes Redis data. Disable via `rename-command` in redis.conf on production, or require explicit confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — redis-cli FLUSHDB",
			input: `redis-cli FLUSHDB`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1444",
					Message: "`redis-cli FLUSHALL`/`FLUSHDB` wipes Redis data. Disable via `rename-command` in redis.conf on production, or require explicit confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1444")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
