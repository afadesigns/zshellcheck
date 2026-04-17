package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1407(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — regular file path",
			input:    `cat /etc/hosts`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — /dev/tcp redirect",
			input: `echo hi > /dev/tcp/1.2.3.4/80`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1407",
					Message: "`/dev/tcp/...` and `/dev/udp/...` are Bash-only virtual files. In Zsh load `zsh/net/tcp` and use `ztcp host port` / `ztcp -c $fd` for TCP.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1407")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
