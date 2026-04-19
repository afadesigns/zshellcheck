package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1843(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `docker run ubuntu` (no cgroup-parent)",
			input:    `docker run ubuntu`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `docker run --cgroup-parent=custom app` (non-host slice)",
			input:    `docker run --cgroup-parent=custom app`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `docker run --cgroup-parent=/ ubuntu`",
			input: `docker run --cgroup-parent=/ ubuntu`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1843",
					Message: "`docker run --cgroup-parent=/` puts the container under a host-managed slice — the engine's memory/CPU caps no longer apply. Drop the flag or pass `--memory`/`--cpus`/`--pids-limit` directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `podman run --cgroup-parent /system.slice alpine`",
			input: `podman run --cgroup-parent /system.slice alpine`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1843",
					Message: "`podman run --cgroup-parent=/system.slice` puts the container under a host-managed slice — the engine's memory/CPU caps no longer apply. Drop the flag or pass `--memory`/`--cpus`/`--pids-limit` directly.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1843")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
