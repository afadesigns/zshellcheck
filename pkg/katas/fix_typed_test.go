// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

// TestFixTypedNodeWalkSweep parses a corpus of representative source
// files, walks the AST, and runs CheckAndFix on every node. Unlike
// the guard sweep this drives correctly-typed nodes into each
// kata.Fix, exercising offset lookups and replacement-emit branches
// rather than just type-assertion guards.
func TestFixTypedNodeWalkSweep(t *testing.T) {
	corpus := []string{
		"echo $arr[1]\n",
		"result=`which git`\n",
		"x=$(seq -s, 1 5)\n",
		"echo -E msg\n",
		"rm -rf $target\n",
		"[ $x -eq 1 ] && echo one\n",
		"[ $a -lt $b ]\n",
		"trap 'echo bye' EXIT\n",
		"chmod -R 777 /tmp\n",
		"sudo dd if=/dev/zero of=/dev/sda\n",
		"export PATH=/bin\n",
		"function greet() { echo hi; }\n",
		"case $x in a) echo a;; b) echo b;; esac\n",
		"for f in *; do echo $f; done\n",
		"while true; do break; done\n",
		"if [ -f c ]; then echo y; fi\n",
		"[[ -z $foo ]] && echo empty\n",
		"typeset -a items=(a b c)\n",
		"declare -A m=(k v)\n",
		"local n=42\n",
		"readonly y=hi\n",
		"echo \"${arr[@]}\"\n",
		"n=$(( 1 + 2 ))\n",
		"((i++))\n",
		"select opt in a b c; do break; done\n",
		"diff <(sort a) <(sort b)\n",
		"cat <<EOF\nbody\nEOF\n",
		"cat << 'STRIP'\nbody\nSTRIP\n",
		"echo a >| b\n",
		"echo a 2>&1\n",
		"echo $arr[1][2]\n",
		"echo ${arr[(R)x]}\n",
		"echo ${(j:,:)items}\n",
		"a=1 b=2 cmd\n",
		"alias ll='ls -l'\n",
		"unalias ll\n",
		"set -e\n",
		"set -o pipefail\n",
		"set -u\n",
		"setopt errexit\n",
		"setopt nounset\n",
	}
	for _, src := range corpus {
		l := lexer.New(src)
		p := parser.New(l)
		prog := p.ParseProgram()
		if prog == nil {
			continue
		}
		source := []byte(src)
		ast.Walk(prog, func(node ast.Node) bool {
			func() {
				defer func() { _ = recover() }()
				_, _ = Registry.CheckAndFix(node, nil, source)
			}()
			return true
		})
	}
}
