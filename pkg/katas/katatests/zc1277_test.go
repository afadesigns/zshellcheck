package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// ZC1277 was retired as a duplicate of ZC1108. Kept as a no-op stub so
// legacy `disabled_katas` lists keep parsing; the detection runs under
// ZC1108 now.

func TestZC1277Stub(t *testing.T) {
	cases := []string{
		"echo hi",
		"ls",
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			v := testutil.Check(in, "ZC1277")
			testutil.AssertViolations(t, in, v, []katas.Violation{})
		})
	}
}
