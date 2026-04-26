// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

func TestZC1472(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []katas.Violation
	}{
		{
			name:     "valid — aws s3 cp file bucket without ACL",
			input:    `aws s3 cp file s3://bucket/`,
			expected: []katas.Violation{},
		},
		{
			name:     "valid — aws s3 cp --acl private",
			input:    `aws s3 cp file s3://bucket/ --acl private`,
			expected: []katas.Violation{},
		},
		{
			name:  "invalid — aws s3 cp --acl public-read",
			input: `aws s3 cp file s3://bucket/ --acl public-read`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1472",
					Message: "Canned ACL `public-read` makes the object (often the bucket) world-readable. Use a scoped bucket policy instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — aws s3 sync --acl=public-read-write",
			input: `aws s3 sync ./ s3://bucket/ --acl=public-read-write`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1472",
					Message: "Canned ACL `public-read-write` makes the object (often the bucket) world-readable. Use a scoped bucket policy instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
		{
			name:  "invalid — aws s3api put-bucket-acl --acl public-read",
			input: `aws s3api put-bucket-acl --bucket foo --acl public-read`,
			expected: []katas.Violation{
				{
					KataID:  "ZC1472",
					Message: "Canned ACL `public-read` makes the object (often the bucket) world-readable. Use a scoped bucket policy instead.",
					Line:    1,
					Column:  1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := testutil.Check(tt.input, "ZC1472")
			testutil.AssertViolations(t, tt.input, violations, tt.expected)
		})
	}
}
