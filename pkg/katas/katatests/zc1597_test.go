package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1597(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — systemd-run as non-root user",
			input:    `systemd-run -p User=www-data /usr/bin/cleanup`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — systemd-run without user property",
			input:    `systemd-run /usr/bin/cleanup`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — systemd-run -p User=root",
			input: `systemd-run -p User=root sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1597",
					Message: "`systemd-run -p User=root` runs arbitrary commands as root via systemd — bypasses the `sudo` audit path. Prefer explicit `sudo` or a fixed systemd unit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — systemd-run -p User=0",
			input: `systemd-run -p User=0 sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1597",
					Message: "`systemd-run -p User=0` runs arbitrary commands as root via systemd — bypasses the `sudo` audit path. Prefer explicit `sudo` or a fixed systemd unit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1597")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
