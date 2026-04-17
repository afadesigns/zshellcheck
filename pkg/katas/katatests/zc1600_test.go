package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1600(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — bwrap for sandboxed exec",
			input:    `bwrap --ro-bind / / --unshare-user --uid 1000 /bin/sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — no chroot call",
			input:    `mount --bind /src /dst`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — chroot /var/sandbox /bin/sh",
			input: `chroot /var/sandbox /bin/sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1600",
					Message: "`chroot` without `--userspec=` runs the inner command as uid 0. Pass `--userspec=USER:GROUP` to drop privileges, or use `bwrap` / `firejail` for user-namespace isolation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — chroot $ROOT sh",
			input: `chroot $ROOT sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1600",
					Message: "`chroot` without `--userspec=` runs the inner command as uid 0. Pass `--userspec=USER:GROUP` to drop privileges, or use `bwrap` / `firejail` for user-namespace isolation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1600")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
