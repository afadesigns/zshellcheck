package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1708(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — find without -L",
			input:    `find /var/log -mtime +30 -delete`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — find -L without destructive action",
			input:    `find -L /opt -name '*.bak'`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — find -L … -delete",
			input: `find -L /var/log -mtime +30 -delete`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1708",
					Message: "`find -L … -delete/-exec` follows symlinks into unintended trees — drop `-L`, add `-xdev`, or scope the walk explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — find -L … -exec rm",
			input: `find -L /var/log -mtime +30 -exec rm -f`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1708",
					Message: "`find -L … -delete/-exec` follows symlinks into unintended trees — drop `-L`, add `-xdev`, or scope the walk explicitly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1708")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
