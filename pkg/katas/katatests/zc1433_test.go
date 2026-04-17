package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1433(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — userdel plain",
			input:    `userdel alice`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — userdel -f",
			input: `userdel -f alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1433",
					Message: "`userdel -f`/`-r` forcibly removes user (kills processes, deletes home). Check for active sessions first with `who -u` / `loginctl list-users`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — userdel -rf",
			input: `userdel -rf alice`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1433",
					Message: "`userdel -f`/`-r` forcibly removes user (kills processes, deletes home). Check for active sessions first with `who -u` / `loginctl list-users`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1433")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
