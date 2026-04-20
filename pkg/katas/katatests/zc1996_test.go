package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1996(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `unshare -m $CMD` (mount namespace only)",
			input:    `unshare -m $CMD`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `bwrap --unshare-all $CMD` (rootless runtime)",
			input:    `bwrap --unshare-all $CMD`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `unshare -Ur $CMD`",
			input: `unshare -Ur $CMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1996",
					Message: "`unshare -Ur` opens a user namespace and maps the caller to uid 0 inside it — also the standard opening move for many kernel-LPE chains. Route legit rootless runtimes through `bwrap`/`podman --rootless` so the intent is clear.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unshare -U $CMD`",
			input: `unshare -U $CMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1996",
					Message: "`unshare -U` opens a user namespace and maps the caller to uid 0 inside it — also the standard opening move for many kernel-LPE chains. Route legit rootless runtimes through `bwrap`/`podman --rootless` so the intent is clear.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `unshare -Urm $CMD` (short bundle)",
			input: `unshare -Urm $CMD`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1996",
					Message: "`unshare -Urm` opens a user namespace and maps the caller to uid 0 inside it — also the standard opening move for many kernel-LPE chains. Route legit rootless runtimes through `bwrap`/`podman --rootless` so the intent is clear.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1996")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
