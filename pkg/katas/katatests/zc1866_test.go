package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1866(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker exec web bash`",
			input:    `docker exec web bash`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker exec -u app web bash`",
			input:    `docker exec -u app web bash`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker exec -u 0 web bash`",
			input: `docker exec -u 0 web bash`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1866",
					Message: "`docker exec -u 0` drops a root shell — bypasses the image's non-root `USER` and, without userns remap, equals host root. Keep execs as the container user.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `podman exec --user=root app sh`",
			input: `podman exec --user=root app sh`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1866",
					Message: "`podman exec -u root` drops a root shell — bypasses the image's non-root `USER` and, without userns remap, equals host root. Keep execs as the container user.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1866")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
