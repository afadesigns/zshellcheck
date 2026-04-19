package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1835(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `smartctl -s on $DISK` (default)",
			input:    `smartctl -s on $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `smartctl -a $DISK` (just report)",
			input:    `smartctl -a $DISK`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `smartctl -s off $DISK`",
			input: `smartctl -s off $DISK`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1835",
					Message: "`smartctl -s off` disables the drive's SMART attribute collection — `smartctl -H` keeps reporting PASSED until the disk falls off the bus. Leave it `on` and configure `smartd.conf` for proactive alerts.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1835")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
