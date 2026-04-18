package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1674(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — docker run default",
			input:    `docker run alpine`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — docker run --oom-score-adj 0",
			input:    `docker run --oom-score-adj 0 alpine`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — docker run --oom-kill-disable",
			input: `docker run --oom-kill-disable alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1674",
					Message: "`--oom-kill-disable` shifts OOM pressure onto the rest of the host — cap memory with `--memory=<limit>` instead of rigging the OOM score.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — podman run --oom-score-adj=-1000",
			input: `podman run --oom-score-adj=-1000 alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1674",
					Message: "`--oom-score-adj=-1000` shifts OOM pressure onto the rest of the host — cap memory with `--memory=<limit>` instead of rigging the OOM score.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1674")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
