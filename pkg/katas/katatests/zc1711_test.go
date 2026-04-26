// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1711(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — etcdctl del --prefix /app/",
			input:    `etcdctl del --prefix /app/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — etcdctl get --prefix \"\" (read-only)",
			input:    `etcdctl get --prefix ""`,
			expected: []katas.Violation{},
		},
		{
			name:  `invalid — etcdctl del --prefix ""`,
			input: `etcdctl del --prefix ""`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1711",
					Message: "`etcdctl del --prefix \"\"` deletes the entire etcd keyspace (including kube-apiserver state) — scope to a specific namespace prefix and review with `get --prefix --keys-only` first.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  `invalid — etcdctl del --from-key ""`,
			input: `etcdctl del --from-key ""`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1711",
					Message: "`etcdctl del --from-key \"\"` deletes the entire etcd keyspace (including kube-apiserver state) — scope to a specific namespace prefix and review with `get --prefix --keys-only` first.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1711")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
