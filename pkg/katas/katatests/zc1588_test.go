package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1588(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — nsenter without target 1",
			input:    `nsenter -t 4242 -m sh`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — nsenter on arbitrary pid",
			input:    `nsenter -t 8123 -m sh`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — nsenter -t 1 -m -u -i -n -p sh",
			input: `nsenter -t 1 -m -u -i -n -p sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1588",
					Message: "`nsenter --target 1` joins the host init namespaces — classic container-escape primitive. Use `docker exec` / `kubectl exec` for legitimate debugging.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — nsenter -t1 -m sh",
			input: `nsenter -t1 -m sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1588",
					Message: "`nsenter --target 1` joins the host init namespaces — classic container-escape primitive. Use `docker exec` / `kubectl exec` for legitimate debugging.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1588")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
