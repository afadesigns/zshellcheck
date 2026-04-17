package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1526(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — wipefs --no-act",
			input:    `wipefs --no-act $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — wipefs -a disk",
			input: `wipefs -a $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1526",
					Message: "`wipefs -a` erases every filesystem signature — unrecoverable. Run with `--no-act` first, or use `sgdisk --zap-all` for scoped deletion.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — wipefs -af disk",
			input: `wipefs -af $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1526",
					Message: "`wipefs -a` erases every filesystem signature — unrecoverable. Run with `--no-act` first, or use `sgdisk --zap-all` for scoped deletion.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1526")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
