package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1548(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — helm install foo chart",
			input:    `helm install foo bitnami/nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — helm install foo chart --disable-openapi-validation",
			input: `helm install foo bitnami/nginx --disable-openapi-validation`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1548",
					Message: "`helm --disable-openapi-validation` hides bad manifests until the controller crashes. Fix the schema deviation.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — helm upgrade foo chart --disable-openapi-validation",
			input: `helm upgrade foo bitnami/nginx --disable-openapi-validation`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1548",
					Message: "`helm --disable-openapi-validation` hides bad manifests until the controller crashes. Fix the schema deviation.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1548")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
