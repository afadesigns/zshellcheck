package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1453(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — sudo apt-get install",
			input:    `sudo apt-get install -y curl`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — sudo pip install",
			input: `sudo pip install requests`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1453",
					Message: "`sudo pip` runs a language package manager as root. Prefer `--user`, a virtualenv/venv, or a version manager (nvm/pyenv/rbenv).",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — sudo npm install -g",
			input: `sudo npm install -g ts-node`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1453",
					Message: "`sudo npm` runs a language package manager as root. Prefer `--user`, a virtualenv/venv, or a version manager (nvm/pyenv/rbenv).",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1453")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
