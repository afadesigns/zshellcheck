package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1957(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `lvchange -ay data/home` (activate)",
			input:    `lvchange -ay data/home`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `vgchange -ay data`",
			input:    `vgchange -ay data`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `lvchange -an data/home`",
			input: `lvchange -an data/home`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1957",
					Message: "`lvchange -an` deactivates the LV/VG — unflushed writes on a mounted fs may be lost, open fds see EIO. Umount and stop holders first, verify with `lsof`/`fuser`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `vgchange -an data`",
			input: `vgchange -an data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1957",
					Message: "`vgchange -an` deactivates the LV/VG — unflushed writes on a mounted fs may be lost, open fds see EIO. Umount and stop holders first, verify with `lsof`/`fuser`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1957")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
