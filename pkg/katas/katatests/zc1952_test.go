// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1952(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `zfs set sync=standard tank/data`",
			input:    `zfs set sync=standard tank/data`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `zfs get sync tank` (read only)",
			input:    `zfs get sync tank`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `zfs set sync=disabled tank/pg`",
			input: `zfs set sync=disabled tank/pg`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1952",
					Message: "`zfs set sync=disabled` turns `fsync()` into a no-op — DBs (PostgreSQL/MariaDB/etcd) lose committed transactions on crash. Leave sync at `standard`; use a SLOG vdev if latency is the concern.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `zfs set sync=disabled $POOL`",
			input: `zfs set sync=disabled $POOL`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1952",
					Message: "`zfs set sync=disabled` turns `fsync()` into a no-op — DBs (PostgreSQL/MariaDB/etcd) lose committed transactions on crash. Leave sync at `standard`; use a SLOG vdev if latency is the concern.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1952")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
