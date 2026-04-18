package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1627(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — crontab from /etc",
			input:    `crontab /etc/cron.install.conf`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — crontab $HOME path",
			input:    `crontab $HOME/.crontab`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — crontab /tmp/newcron",
			input: `crontab /tmp/newcron`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1627",
					Message: "`crontab /tmp/newcron` reads cron rules from a world-traversable path — a concurrent local user can substitute the file between write and read. Stage the file in `$XDG_RUNTIME_DIR/` or `mktemp -d`, or pipe via `crontab -`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — crontab -u bob /tmp/evil",
			input: `crontab -u bob /tmp/evil`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1627",
					Message: "`crontab /tmp/evil` reads cron rules from a world-traversable path — a concurrent local user can substitute the file between write and read. Stage the file in `$XDG_RUNTIME_DIR/` or `mktemp -d`, or pipe via `crontab -`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1627")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
