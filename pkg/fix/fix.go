// Package fix applies FixEdit sets produced by kata Fix functions to
// source files. It handles offset math (1-based Line/Column to absolute
// byte offsets), sorts edits bottom-up so earlier offsets stay valid,
// and renders either a rewritten source string or a unified diff.
package fix

import (
	"fmt"
	"sort"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// Apply returns the source rewritten with every edit in edits applied.
//
// Overlap handling: when two edits overlap (for example an outer
// “ `which git` “ -> `$(which git)` fix and an inner `which` ->
// `whence` fix), the outer edit wins and the inner is dropped. Running
// the fixer a second time picks up the previously-suppressed inner
// edit. This keeps each run idempotent and avoids the surprise of a
// partial rewrite inside an already-rewritten span.
//
// Non-overlapping edits are applied in descending start-offset order
// so that earlier splices do not invalidate the offsets of later ones.
func Apply(source string, edits []katas.FixEdit) (string, error) {
	if len(edits) == 0 {
		return source, nil
	}

	resolved, err := resolveOffsets(source, edits)
	if err != nil {
		return "", err
	}

	resolved = resolveConflicts(resolved)

	// Sort by start offset descending so we splice from the end
	// backwards.
	sort.SliceStable(resolved, func(i, j int) bool {
		return resolved[i].start > resolved[j].start
	})

	out := source
	for _, e := range resolved {
		out = out[:e.start] + e.replace + out[e.start+e.length:]
	}
	return out, nil
}

// resolveConflicts drops edits that overlap the span of a surviving
// edit. The policy: prefer the edit with the earlier start; when two
// edits share a start offset, prefer the longer. Returned edits are
// guaranteed pairwise-disjoint and preserve the input order for
// disjoint edits so Apply stays deterministic.
func resolveConflicts(edits []resolvedEdit) []resolvedEdit {
	if len(edits) < 2 {
		return edits
	}
	// Sort ascending by start; break ties by descending length so the
	// longer span wins the single-pass overlap check.
	sort.SliceStable(edits, func(i, j int) bool {
		if edits[i].start != edits[j].start {
			return edits[i].start < edits[j].start
		}
		return edits[i].length > edits[j].length
	})

	kept := edits[:0]
	var lastEnd int
	for _, e := range edits {
		if len(kept) > 0 && e.start < lastEnd {
			// Overlaps with an already-kept edit; drop it.
			continue
		}
		kept = append(kept, e)
		lastEnd = e.start + e.length
	}
	// Return a copy so callers cannot depend on the in-place shuffle.
	out := make([]resolvedEdit, len(kept))
	copy(out, kept)
	return out
}

// Overlap returns true when two edits share any source bytes. Exposed
// for callers that want to pre-filter fix candidates before calling
// Apply (for example a CLI that wants to report how many edits were
// suppressed by nesting).
func Overlap(a, b katas.FixEdit) bool {
	// Same line, overlapping column ranges.
	if a.Line == b.Line {
		return a.Column < b.Column+b.Length && b.Column < a.Column+a.Length
	}
	// Different lines: non-overlapping by definition (edits don't
	// span lines via the 1-D column metric alone).
	return false
}

// Diff returns a unified-format diff between the original source and
// the source after applying edits. The diff uses the filename on both
// the "---" and "+++" lines with a "(fixed)" suffix on the new side.
func Diff(filename, source string, edits []katas.FixEdit) (string, error) {
	fixed, err := Apply(source, edits)
	if err != nil {
		return "", err
	}
	if fixed == source {
		return "", nil
	}
	return unifiedDiff(filename, source, fixed), nil
}

type resolvedEdit struct {
	start   int
	length  int
	replace string
}

func resolveOffsets(source string, edits []katas.FixEdit) ([]resolvedEdit, error) {
	// Pre-compute line start offsets so each edit can map Line:Column to
	// a byte offset in one lookup. lineStarts[i] is the 0-based offset of
	// the first byte of line (i+1); line numbers are 1-based.
	lineStarts := []int{0}
	for i := 0; i < len(source); i++ {
		if source[i] == '\n' {
			lineStarts = append(lineStarts, i+1)
		}
	}

	out := make([]resolvedEdit, 0, len(edits))
	for idx, e := range edits {
		if e.Line < 1 || e.Line > len(lineStarts) {
			return nil, fmt.Errorf("fix: edit #%d has out-of-range line %d (source has %d lines)",
				idx, e.Line, len(lineStarts))
		}
		if e.Column < 1 {
			return nil, fmt.Errorf("fix: edit #%d has non-positive column %d", idx, e.Column)
		}
		if e.Length < 0 {
			return nil, fmt.Errorf("fix: edit #%d has negative length %d", idx, e.Length)
		}
		start := lineStarts[e.Line-1] + (e.Column - 1)
		if start > len(source) {
			return nil, fmt.Errorf("fix: edit #%d starts past end of source (offset %d > len %d)",
				idx, start, len(source))
		}
		if start+e.Length > len(source) {
			return nil, fmt.Errorf("fix: edit #%d ends past end of source (end %d > len %d)",
				idx, start+e.Length, len(source))
		}
		out = append(out, resolvedEdit{
			start:   start,
			length:  e.Length,
			replace: e.Replace,
		})
	}
	return out, nil
}

// unifiedDiff produces a minimal unified diff. It is not a full
// implementation of diff3 — it emits a single hunk per contiguous
// block of changes with three lines of context on each side. Good
// enough for CLI display of auto-fix previews.
func unifiedDiff(filename, a, b string) string {
	al := splitLines(a)
	bl := splitLines(b)

	// Myers-style LCS would be ideal; for short files the line-by-line
	// longest-common-subsequence table below is plenty fast.
	lcs := lcsTable(al, bl)

	var hunks []hunk
	i, j := 0, 0
	for i < len(al) || j < len(bl) {
		// Matching run — accumulate context but do not start a hunk yet.
		for i < len(al) && j < len(bl) && al[i] == bl[j] {
			i++
			j++
		}
		if i >= len(al) && j >= len(bl) {
			break
		}
		// Divergence: walk forward until lines match again.
		h := hunk{aStart: i, bStart: j}
		for i < len(al) || j < len(bl) {
			if i < len(al) && j < len(bl) && al[i] == bl[j] {
				break
			}
			// Choose the side with more remaining lcs; fallback to
			// deletion then insertion to keep output stable when the
			// table is square.
			switch {
			case j >= len(bl) || (i < len(al) && lcs[i+1][j] >= lcs[i][j+1]):
				h.edits = append(h.edits, diffEdit{op: '-', text: al[i]})
				i++
			default:
				h.edits = append(h.edits, diffEdit{op: '+', text: bl[j]})
				j++
			}
		}
		h.aEnd = i
		h.bEnd = j
		hunks = append(hunks, h)
	}

	var out strings.Builder
	fmt.Fprintf(&out, "--- %s\n", filename)
	fmt.Fprintf(&out, "+++ %s (fixed)\n", filename)
	for _, h := range hunks {
		aContextStart := maxInt(0, h.aStart-3)
		aContextEnd := minInt(len(al), h.aEnd+3)
		bContextStart := maxInt(0, h.bStart-3)
		bContextEnd := minInt(len(bl), h.bEnd+3)

		fmt.Fprintf(&out, "@@ -%d,%d +%d,%d @@\n",
			aContextStart+1, aContextEnd-aContextStart,
			bContextStart+1, bContextEnd-bContextStart,
		)
		// Leading context
		for k := aContextStart; k < h.aStart; k++ {
			fmt.Fprintf(&out, " %s\n", al[k])
		}
		// The hunk edits in order
		for _, e := range h.edits {
			fmt.Fprintf(&out, "%c%s\n", e.op, e.text)
		}
		// Trailing context
		for k := h.aEnd; k < aContextEnd; k++ {
			fmt.Fprintf(&out, " %s\n", al[k])
		}
	}
	return out.String()
}

type hunk struct {
	aStart, aEnd int
	bStart, bEnd int
	edits        []diffEdit
}

type diffEdit struct {
	op   byte
	text string
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	// Drop the trailing empty element produced by a trailing newline —
	// callers round-trip through Apply which preserves the newline.
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func lcsTable(a, b []string) [][]int {
	// lcs[i][j] is the length of the longest common subsequence of
	// a[i:] and b[j:].
	t := make([][]int, len(a)+1)
	for i := range t {
		t[i] = make([]int, len(b)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			switch {
			case a[i] == b[j]:
				t[i][j] = t[i+1][j+1] + 1
			case t[i+1][j] >= t[i][j+1]:
				t[i][j] = t[i+1][j]
			default:
				t[i][j] = t[i][j+1]
			}
		}
	}
	return t
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
