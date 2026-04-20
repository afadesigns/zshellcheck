package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1966(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `zpool import $POOL` (no force)",
			input:    `zpool import $POOL`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zpool export $POOL`",
			input:    `zpool export $POOL`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `zpool import -f $POOL`",
			input: `zpool import -f $POOL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1966",
					Message: "`zpool import -f` bypasses hostid/txg safety — forced import of a pool already online elsewhere (SAN/HA) corrupts it; forced export drops in-flight txgs. `zfs unmount -a` first, then plain `zpool export`/`import`.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `zpool export -f $POOL`",
			input: `zpool export -f $POOL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1966",
					Message: "`zpool export -f` bypasses hostid/txg safety — forced import of a pool already online elsewhere (SAN/HA) corrupts it; forced export drops in-flight txgs. `zfs unmount -a` first, then plain `zpool export`/`import`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1966")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
