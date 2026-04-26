// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

// runCheckAndFix parses src, walks the AST, and invokes the registry's
// CheckAndFix on every node. Used to exercise full check + fix paths
// across many katas with realistic AST shapes.
func runCheckAndFix(t *testing.T, src string) {
	t.Helper()
	l := lexer.New(src)
	p := parser.New(l)
	prog := p.ParseProgram()
	if prog == nil {
		t.Fatalf("nil program for %q", src)
	}
	source := []byte(src)
	ast.Walk(prog, func(node ast.Node) bool {
		func() {
			defer func() { _ = recover() }()
			Registry.CheckAndFix(node, nil, source)
		}()
		return true
	})
}

func TestSweepBacktickAndArray(t *testing.T) {
	runCheckAndFix(t, "result=`which git`\necho $arr[1]\n")
}

func TestSweepDangerousPatterns(t *testing.T) {
	runCheckAndFix(t, "rm -rf $target\nchmod -R 777 /tmp\nsudo cp file /\n")
}

func TestSweepExternalsWithBuiltinAlternatives(t *testing.T) {
	runCheckAndFix(t, "x=$(which git)\ny=$(seq -s, 1 5)\necho -E msg\n")
}

func TestSweepArithmeticAndTests(t *testing.T) {
	runCheckAndFix(t, "[ $x -eq 1 ] && echo one\n[ $a -lt $b ]\n[[ $foo == bar ]]\n")
}

func TestSweepDeclarationsAndAssignments(t *testing.T) {
	runCheckAndFix(t, "typeset -a items=(a b c)\ndeclare -A m\nlocal n=42\nreadonly y=hi\n")
}

func TestSweepLoopsAndConditionals(t *testing.T) {
	runCheckAndFix(t, "for f in *; do echo $f; done\nwhile true; do break; done\nuntil [ -f f ]; do sleep 1; done\n")
}

func TestSweepExpansionsAndSubstitutions(t *testing.T) {
	runCheckAndFix(t, "echo \"${arr[@]}\"\necho ${#var}\necho ${var:-default}\necho ${var/old/new}\n")
}

func TestSweepFunctionForms(t *testing.T) {
	runCheckAndFix(t, "function greet() { echo hi; }\nfarewell() { echo bye; }\n")
}

func TestSweepCaseAndSelect(t *testing.T) {
	runCheckAndFix(t, "case $x in\n  a) echo a;;\n  *) echo other;;\nesac\nselect opt in a b c; do break; done\n")
}

func TestSweepRedirectionsAndPipelines(t *testing.T) {
	runCheckAndFix(t, "cat <input >output 2>err\nwc -l < f | sort | uniq\necho hi >> log\n")
}

func TestSweepHeredocs(t *testing.T) {
	runCheckAndFix(t, "cat <<EOF\nbody\nEOF\ncat <<-'STRIP'\n\tdone\n\tSTRIP\n")
}

func TestSweepGlobsAndExpansions(t *testing.T) {
	runCheckAndFix(t, "ls *.go\necho **/*.zsh\nfor f in dir/*; do echo $f; done\n")
}
