package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1880(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `kubectl annotate pod/foo key=val`",
			input:    `kubectl annotate pod/foo key=val`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `kubectl label pod/foo role=app`",
			input:    `kubectl label pod/foo role=app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `kubectl annotate pod/foo --overwrite key=val`",
			input: `kubectl annotate pod/foo --overwrite key=val`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1880",
					Message: "`kubectl annotate --overwrite` silently replaces an existing controller signal — cert-manager, external-dns, HPA watchers reconcile on the new value. Inspect first; drop `--overwrite` so conflicts error.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `kubectl label node/bar --overwrite role=worker`",
			input: `kubectl label node/bar --overwrite role=worker`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1880",
					Message: "`kubectl label --overwrite` silently replaces an existing controller signal — cert-manager, external-dns, HPA watchers reconcile on the new value. Inspect first; drop `--overwrite` so conflicts error.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1880")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
