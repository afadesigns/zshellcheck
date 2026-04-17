package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1437(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — dmesg -T (human time)",
			input:    `dmesg -T`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — dmesg -c",
			input: `dmesg -c`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1437",
					Message: "`dmesg -c`/`-C` clears the kernel ring buffer — subsequent debugging loses earlier messages. Use plain `dmesg` or `journalctl -k`.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1437")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
