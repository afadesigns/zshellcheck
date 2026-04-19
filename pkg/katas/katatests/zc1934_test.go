package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1934(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unsetopt AUTO_NAME_DIRS` (explicit default)",
			input:    `unsetopt AUTO_NAME_DIRS`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `setopt AUTO_CD` (unrelated)",
			input:    `setopt AUTO_CD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `setopt AUTO_NAME_DIRS`",
			input: `setopt AUTO_NAME_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1934",
					Message: "`setopt AUTO_NAME_DIRS` auto-registers every absolute-path parameter as a named dir — `foo=/srv/data` makes `~foo` expand, `%~` prompts surface names the user never picked. Keep off; use `hash -d name=/path`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unsetopt NO_AUTO_NAME_DIRS`",
			input: `unsetopt NO_AUTO_NAME_DIRS`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1934",
					Message: "`unsetopt NO_AUTO_NAME_DIRS` auto-registers every absolute-path parameter as a named dir — `foo=/srv/data` makes `~foo` expand, `%~` prompts surface names the user never picked. Keep off; use `hash -d name=/path`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1934")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
