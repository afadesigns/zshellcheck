// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package fix

import (
	"os"
	"regexp"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// goldenTestFiles hold the per-kata golden fix-output tests. The
// coverage guard scans them for the kata IDs they exercise.
var goldenTestFiles = []string{
	"integration_test.go",
	"integration_gap_test.go",
}

// TestEveryFixableKataHasGoldenTest enumerates every kata that ships a
// Fix and asserts a golden fix-output test exists for it. A golden test
// applies the fix and pins its exact rewrite, so a kata without one can
// change its output — or break it — undetected. The guard reads the
// registry (the source of truth for what is fixable) and the golden test
// files (the source of truth for what is covered), so it can neither
// drift from a stale list nor pass on a kata whose test was deleted.
// When a new fixable kata is added, this test fails until its golden
// test lands, extending the project's "no new public API without a test"
// rule to the auto-fixer.
func TestEveryFixableKataHasGoldenTest(t *testing.T) {
	idRe := regexp.MustCompile(`ZC\d{4}`)
	covered := make(map[string]bool)
	for _, name := range goldenTestFiles {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read golden test file %s: %v", name, err)
		}
		for _, id := range idRe.FindAllString(string(data), -1) {
			covered[id] = true
		}
	}

	missing := make([]string, 0)
	for _, k := range katas.Registry.AllKatas() {
		if k.Fix == nil {
			continue
		}
		if !covered[k.ID] {
			missing = append(missing, k.ID)
		}
	}

	if len(missing) != 0 {
		t.Errorf("fixable katas with no golden fix-output test (%d): %v\n"+
			"add a golden test to one of %v that applies the fix and asserts its exact output",
			len(missing), missing, goldenTestFiles)
	}
}
