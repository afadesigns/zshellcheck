package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1864(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `mount -o remount,noexec /tmp` (tightening)",
			input:    `mount -o remount,noexec /tmp`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `mount -o remount,rw /` (unrelated)",
			input:    `mount -o remount,rw /`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `mount -o remount,exec /tmp`",
			input: `mount -o remount,exec /tmp`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1864",
					Message: "`mount -o remount,exec` re-enables `exec` on a `noexec`/`nosuid`/`nodev`-hardened mount — dropped payloads suddenly execute. Pair with a `trap ... EXIT` that restores the original flags or skip the remount.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `mount -o remount,rw,suid /var`",
			input: `mount -o remount,rw,suid /var`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1864",
					Message: "`mount -o remount,rw,suid` re-enables `suid` on a `noexec`/`nosuid`/`nodev`-hardened mount — dropped payloads suddenly execute. Pair with a `trap ... EXIT` that restores the original flags or skip the remount.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1864")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
