package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1756(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `chmod 660 /var/run/docker.sock` (group only)",
			input:    `chmod 660 /var/run/docker.sock`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `chmod 666 /tmp/file` (not a runtime socket)",
			input:    `chmod 666 /tmp/file`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `chmod 666 /var/run/docker.sock`",
			input: `chmod 666 /var/run/docker.sock`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1756",
					Message: "`chmod 666 /var/run/docker.sock` grants every local user access to a root-equivalent container-runtime socket. Keep `0660` owned by the runtime group (`root:docker` etc.) and restrict membership.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `chmod 777 /run/containerd/containerd.sock`",
			input: `chmod 777 /run/containerd/containerd.sock`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1756",
					Message: "`chmod 777 /run/containerd/containerd.sock` grants every local user access to a root-equivalent container-runtime socket. Keep `0660` owned by the runtime group (`root:docker` etc.) and restrict membership.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1756")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
