package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1651(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — localhost bind",
			input:    `docker run -p 127.0.0.1:8080:80 nginx`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — implicit port (not explicit 0.0.0.0)",
			input:    `docker run -p 8080:80 nginx`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — 0.0.0.0 explicit",
			input: `docker run -p 0.0.0.0:8080:80 nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1651",
					Message: "`docker run -p 0.0.0.0:8080:80` publishes to every interface. Bind to `127.0.0.1:HOST:CONT` and put nginx / caddy in front for external access.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — [::] IPv6",
			input: `podman run -p [::]:8080:80 nginx`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1651",
					Message: "`podman run -p [::]:8080:80` publishes to every interface. Bind to `127.0.0.1:HOST:CONT` and put nginx / caddy in front for external access.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1651")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
