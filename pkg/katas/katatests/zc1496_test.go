package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1496(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — strings /bin/ls",
			input:    `strings /bin/ls`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — strings /dev/mem",
			input: `strings /dev/mem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1496",
					Message: "Reading `/dev/mem` leaks kernel / physical memory. Use kdump + crash on a crash-kernel image if you need a dump.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — xxd /dev/port",
			input: `xxd /dev/port`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1496",
					Message: "Reading `/dev/port` leaks kernel / physical memory. Use kdump + crash on a crash-kernel image if you need a dump.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — cat /dev/kmem",
			input: `cat /dev/kmem`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1496",
					Message: "Reading `/dev/kmem` leaks kernel / physical memory. Use kdump + crash on a crash-kernel image if you need a dump.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1496")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
