package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1951(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — `ceph osd pool ls detail`",
			input:    `ceph osd pool ls detail`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — `ceph -s` (health)",
			input:    `ceph -s`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — `ceph osd pool delete rbd rbd --yes-i-really-really-mean-it now` (mangled)",
			input: `ceph osd pool delete rbd rbd --yes-i-really-really-mean-it now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1951",
					Message: "`ceph … --yes-i-really-really-mean-it` automates the double-safety phrase — a typo or stale loop silently deletes production pools. Run deletions interactively, or spell the pool name in a runbook commit.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — `ceph config-key rm key --yes-i-really-mean-it now`",
			input: `ceph config-key rm key --yes-i-really-mean-it now`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1951",
					Message: "`ceph … --yes-i-really-really-mean-it` automates the double-safety phrase — a typo or stale loop silently deletes production pools. Run deletions interactively, or spell the pool name in a runbook commit.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1951")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
