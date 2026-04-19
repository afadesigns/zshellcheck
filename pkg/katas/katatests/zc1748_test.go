package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1748(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `helm repo add bitnami https://charts.bitnami.com/bitnami`",
			input:    `helm repo add bitnami https://charts.bitnami.com/bitnami`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `helm repo update`",
			input:    `helm repo update`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `helm repo add myrepo http://internal/charts`",
			input: `helm repo add myrepo http://internal/charts`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1748",
					Message: "`helm repo add myrepo http://internal/charts` fetches charts over plaintext HTTP — any MITM swaps the chart and its referenced images. Use `https://` and `helm install --verify` for provenance.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1748")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
