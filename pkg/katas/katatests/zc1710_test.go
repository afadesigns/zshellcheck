package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1710(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — journalctl --vacuum-time=2weeks (real retention)",
			input:    `journalctl -q --vacuum-time=2weeks`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — journalctl --vacuum-size=500M",
			input:    `journalctl -q --vacuum-size=500M`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — journalctl --vacuum-size=1 (wipe)",
			input: `journalctl -q --vacuum-size=1`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1710",
					Message: "`journalctl --vacuum-size=1` flushes the systemd journal — classic audit-clear shape. Set retention in `/etc/systemd/journald.conf` (`SystemMaxUse=`, `MaxRetentionSec=`) instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — journalctl --vacuum-time=1s",
			input: `journalctl -m --vacuum-time=1s`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1710",
					Message: "`journalctl --vacuum-time=1s` flushes the systemd journal — classic audit-clear shape. Set retention in `/etc/systemd/journald.conf` (`SystemMaxUse=`, `MaxRetentionSec=`) instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1710")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
