// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/testutil"
)

// TestZC1188DirectionAwareAdvice pins the PATH-direction advice. A
// prepend (`dir:$PATH`) must be rewritten as `path=(dir $path)`; the
// append form `path+=(dir)` would move the directory to the end and
// reverse command precedence, so the kata must not recommend it there.
func TestZC1188DirectionAwareAdvice(t *testing.T) {
	prepend := testutil.Check("export PATH=/opt/bin:$PATH", "ZC1188")
	if len(prepend) != 1 {
		t.Fatalf("ZC1188 should fire once on a prepend, got %d", len(prepend))
	}
	if !strings.Contains(prepend[0].Message, "Prepend with `path=(dir $path)`") {
		t.Errorf("ZC1188 prepend should advise the prepend form, got %q", prepend[0].Message)
	}

	appendCase := testutil.Check("export PATH=$PATH:/opt/bin", "ZC1188")
	if len(appendCase) != 1 {
		t.Fatalf("ZC1188 should fire once on an append, got %d", len(appendCase))
	}
	if !strings.Contains(appendCase[0].Message, "Append with `path+=(dir)`") {
		t.Errorf("ZC1188 append should advise the append form, got %q", appendCase[0].Message)
	}

	// A replacement that is neither a prepend nor an append gets both forms.
	other := testutil.Check("export PATH=/usr/local/bin", "ZC1188")
	if len(other) != 1 || !strings.Contains(other[0].Message, "to prepend") {
		t.Errorf("ZC1188 should give both forms for a plain PATH assignment")
	}
}
