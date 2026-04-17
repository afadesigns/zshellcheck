package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1538(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — zpool list",
			input:    `zpool list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — zfs list",
			input:    `zfs list`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — zfs destroy mydataset (no -r)",
			input:    `zfs destroy tank/data/old`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — zpool destroy -f tank",
			input: `zpool destroy -f tank`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1538",
					Message: "`zpool destroy -f` irrecoverably destroys the ZFS pool/dataset and every snapshot on it. Require explicit target confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — zfs destroy -rR tank/data",
			input: `zfs destroy -rR tank/data`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1538",
					Message: "`zfs destroy -rR` irrecoverably destroys the ZFS pool/dataset and every snapshot on it. Require explicit target confirmation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1538")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
