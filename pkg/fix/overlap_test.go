// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package fix

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestOverlapSameLineOverlapping(t *testing.T) {
	a := katas.FixEdit{Line: 1, Column: 1, Length: 5, Replace: ""}
	b := katas.FixEdit{Line: 1, Column: 3, Length: 4, Replace: ""}
	if !Overlap(a, b) {
		t.Errorf("expected overlap on same line")
	}
}

func TestOverlapSameLineDisjoint(t *testing.T) {
	a := katas.FixEdit{Line: 1, Column: 1, Length: 2, Replace: ""}
	b := katas.FixEdit{Line: 1, Column: 5, Length: 2, Replace: ""}
	if Overlap(a, b) {
		t.Errorf("disjoint edits reported as overlapping")
	}
}

func TestOverlapDifferentLines(t *testing.T) {
	a := katas.FixEdit{Line: 1, Column: 1, Length: 5, Replace: ""}
	b := katas.FixEdit{Line: 2, Column: 1, Length: 5, Replace: ""}
	if Overlap(a, b) {
		t.Errorf("edits on different lines reported as overlapping")
	}
}
