// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"fmt"
	"sort"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

// Severity defines the severity of a violation.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
	SeverityStyle   Severity = "style"
)

// Violation represents a found violation in the code.
type Violation struct {
	KataID  string
	Message string
	Line    int
	Column  int
	Level   Severity
}

// FixEdit is a single text replacement applied by the auto-fixer.
// Coordinates are 1-based; Length is the number of bytes of source to
// replace starting at Line:Column. Replace may be empty (deletion) and
// may span multiple lines. KataID records the kata that produced the edit
// so the CLI can withhold behavior-changing (unsafe) fixes by default.
type FixEdit struct {
	Line    int
	Column  int
	Length  int
	Replace string
	KataID  string
	// Group ties every edit produced by one Fix call to the violation it
	// repairs. A rewrite such as `[[ a -eq b ]]` to `(( a == b ))` spans
	// three edits that only make sense together, so the applier keeps or
	// discards a group as a unit and a filter drops a group as a unit.
	Group int
}

// safeFixKatas lists the katas whose auto-fix is purely syntactic and
// value-preserving — applying it cannot change runtime behavior or drop a
// comment. Every other fix may alter behavior (a command or flag swap, a
// scope or glob-qualifier change) and is applied only under -unsafe-fixes.
var safeFixKatas = map[string]bool{
	"ZC1001": true, // $arr[i] -> ${arr[i]}
	"ZC1002": true, // `cmd` -> $(cmd)
	"ZC1073": true, // drop redundant $ inside (( … ))
	"ZC1086": true, // function f { } -> f() { }
	"ZC1411": true, // enable -n name -> disable name
	"ZC1502": true, // insert `--` end-of-options guard
	"ZC1637": true, // readonly x -> typeset -r x
	"ZC1643": true, // $(cat f) -> $(<f)
}

// IsSafeFix reports whether the kata's auto-fix is value-preserving and so
// applied by `-fix` without `-unsafe-fixes`.
func (kr *KatasRegistry) IsSafeFix(id string) bool {
	return safeFixKatas[id]
}

// Kata represents a single linting rule. Fix is optional; when non-nil
// the auto-fixer invokes it with the AST node, the violation, and the
// full file source (byte slice) so the fix can inspect a span around
// the violation before producing edits. Katas with no safe
// deterministic fix leave Fix nil and the fixer skips them.
type Kata struct {
	ID          string
	Title       string
	Description string
	Severity    Severity
	Check       func(node ast.Node) []Violation
	Fix         func(node ast.Node, v Violation, source []byte) []FixEdit
}

// KatasRegistry is a registry for all available Katas.
type KatasRegistry struct {
	KatasByType map[string][]Kata
	KatasByID   map[string]Kata
}

// NewKatasRegistry creates a new KatasRegistry.
func NewKatasRegistry() *KatasRegistry {
	return &KatasRegistry{
		KatasByType: make(map[string][]Kata),
		KatasByID:   make(map[string]Kata),
	}
}

// RegisterKata registers a new Kata.
func (kr *KatasRegistry) RegisterKata(nodeType ast.Node, kata Kata) {
	if kata.Severity == "" {
		kata.Severity = SeverityWarning
	}
	key := fmt.Sprintf("%T", nodeType)
	kr.KatasByType[key] = append(kr.KatasByType[key], kata)
	kr.KatasByID[kata.ID] = kata
}

// GetKata returns a Kata by its ID.
func (kr *KatasRegistry) GetKata(id string) (Kata, bool) {
	kata, ok := kr.KatasByID[id]
	return kata, ok
}

// IsFixable reports whether the kata with the given ID ships a
// deterministic auto-fix.
func (kr *KatasRegistry) IsFixable(id string) bool {
	kata, ok := kr.KatasByID[id]
	return ok && kata.Fix != nil
}

// KatasByNodeType returns all registered Katas grouped by node type.
func (kr *KatasRegistry) KatasByNodeType() map[string][]Kata {
	return kr.KatasByType
}

// AllKatas returns every registered kata sorted by ID. It backs the
// `--list-rules` and `--explain` CLI surfaces and any tooling that needs
// a stable, deduplicated enumeration of the rule set.
func (kr *KatasRegistry) AllKatas() []Kata {
	out := make([]Kata, 0, len(kr.KatasByID))
	for _, k := range kr.KatasByID {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (kr *KatasRegistry) Check(node ast.Node, disabledKatas []string) []Violation {
	var violations []Violation
	key := fmt.Sprintf("%T", node)
	if katasForNode, ok := kr.KatasByType[key]; ok {
		for _, kata := range katasForNode {
			// Check if disabled
			disabled := false
			for _, d := range disabledKatas {
				if d == kata.ID {
					disabled = true
					break
				}
			}
			if !disabled {
				vs := kata.Check(node)
				for i := range vs {
					if vs[i].Level == "" {
						vs[i].Level = kata.Severity
					}
				}
				violations = append(violations, vs...)
			}
		}
	}
	return violations
}

// FixesFor invokes the Fix function of the kata that produced the
// given violation, if any, and returns the resulting edits. An empty
// slice means the kata has no deterministic auto-fix; callers should
// skip the violation in fix mode.
func (kr *KatasRegistry) FixesFor(node ast.Node, v Violation, source []byte) []FixEdit {
	kata, ok := kr.KatasByID[v.KataID]
	if !ok || kata.Fix == nil {
		return nil
	}
	// A single violation's edits form one group.
	return stampKataID(kata.Fix(node, v, source), kata.ID, 0)
}

// stampKataID records the producing kata on each edit so the CLI can
// filter behavior-changing fixes without re-deriving their origin.
func stampKataID(edits []FixEdit, id string, group int) []FixEdit {
	for i := range edits {
		edits[i].KataID = id
		edits[i].Group = group
	}
	return edits
}

// CheckAndFix runs Check for every kata registered against the node
// type and, when a kata also declares a Fix, invokes that Fix for each
// emitted violation. Returns the violations (including ones without a
// fix) and the concatenated edits. Use this from the CLI fix mode so
// each node is visited exactly once.
func (kr *KatasRegistry) CheckAndFix(node ast.Node, disabledKatas []string, source []byte) ([]Violation, []FixEdit) {
	var violations []Violation
	var edits []FixEdit
	key := fmt.Sprintf("%T", node)
	katasForNode, ok := kr.KatasByType[key]
	if !ok {
		return nil, nil
	}
	for _, kata := range katasForNode {
		skip := false
		for _, d := range disabledKatas {
			if d == kata.ID {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		vs := kata.Check(node)
		for i := range vs {
			if vs[i].Level == "" {
				vs[i].Level = kata.Severity
			}
			if kata.Fix != nil {
				// Group is the violation's index in the slice this call
				// returns; the caller shifts it by the number of violations
				// already collected so every group is unique per file.
				group := len(violations) + i
				edits = append(edits, stampKataID(kata.Fix(node, vs[i], source), kata.ID, group)...)
			}
		}
		violations = append(violations, vs...)
	}
	return violations, edits
}

// DefaultKatasRegistry is the default registry.
var DefaultKatasRegistry = NewKatasRegistry()

// RegisterKata registers a new Kata with the default registry.
func RegisterKata(nodeType ast.Node, kata Kata) {
	DefaultKatasRegistry.RegisterKata(nodeType, kata)
}

// Registry is the global registry.
var Registry = DefaultKatasRegistry
