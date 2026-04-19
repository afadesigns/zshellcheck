package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1749(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `virsh undefine mydomain` (config only)",
			input:    `virsh undefine mydomain`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `virsh list --all`",
			input:    `virsh list --all`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `virsh undefine mydomain --remove-all-storage`",
			input: `virsh undefine mydomain --remove-all-storage`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1749",
					Message: "`virsh undefine ... --remove-all-storage` deletes every disk image the domain references — no soft-delete, no recycle bin. Back up first (`qemu-img convert`), `undefine` without the flag, then `virsh vol-delete` deliberately.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `virsh undefine mydomain --wipe-storage`",
			input: `virsh undefine mydomain --wipe-storage`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1749",
					Message: "`virsh undefine ... --wipe-storage` deletes every disk image the domain references — no soft-delete, no recycle bin. Back up first (`qemu-img convert`), `undefine` without the flag, then `virsh vol-delete` deliberately.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1749")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
