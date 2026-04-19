package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1943(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `systemd-nspawn -D /srv/container /bin/sh` (not booting)",
			input:    `systemd-nspawn -D /srv/container /bin/sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `machinectl start web`",
			input:    `machinectl start web`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `systemd-nspawn -b -D /srv/container`",
			input: `systemd-nspawn -b -D /srv/container`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1943",
					Message: "`systemd-nspawn -b` runs the rootfs's `/sbin/init` with minimal isolation — init scripts execute first and can probe the host. Use `-U`, drop caps with `--capability=`, pair with `--private-network`, prefer `machinectl start`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `systemd-nspawn --boot -D $ROOT` (mangled)",
			input: `systemd-nspawn --boot -D $ROOT`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1943",
					Message: "`systemd-nspawn --boot` runs the rootfs's `/sbin/init` with minimal isolation — init scripts execute first and can probe the host. Use `-U`, drop caps with `--capability=`, pair with `--private-network`, prefer `machinectl start`.",
					Line:    1,
					Column:  18,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1943")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
