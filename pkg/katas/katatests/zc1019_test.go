package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// ZC1019 was retired as a duplicate of ZC1005. Kept as a no-op stub so
// legacy `disabled_katas` lists keep parsing; the detection runs under
// ZC1005 now.

func TestZC1019Stub(t *testing.T) {
	cases := []string{
		"echo hi",
		"ls",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			v := testutil.Check(in, "ZC1019")
			testutil.AssertViolations(t, in, v, []katas.Violation{})
		})
	}
}
