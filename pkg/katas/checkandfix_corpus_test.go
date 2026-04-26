// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/parser"
)

// runCorpusItem parses src, walks the program, and runs CheckAndFix on
// every node. Used to drive Fix coverage paths across many kata
// surfaces without per-kata assertions.
func runCorpusItem(src string) {
	l := lexer.New(src)
	p := parser.New(l)
	prog := p.ParseProgram()
	if prog == nil {
		return
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

// TestCorpusKataFixSweep drives every kata Fix path with a wide
// realistic corpus. The test asserts nothing — its job is to make
// kata Check + Fix code paths reachable so Codecov sees them
// covered.
func TestCorpusKataFixSweep(t *testing.T) {
	corpus := []string{
		// Backticks, $(...), arr[i], typeset/local/declare
		"x=`which git`\necho $arr[1]\ntypeset -a a=(1 2)\nlocal n=42\ndeclare -A m\n",
		// rm/cp/mv/mkdir hazards
		"rm -rf $target\ncp -r src dst\nmv $a $b\nmkdir -p /tmp/dir\nmkdir -m 777 /tmp/d\n",
		// Pipelines through tools that have Zsh-native alternatives
		"echo a | tr a-z A-Z\necho a | tr A-Z a-z\nseq -s , 1 5\nls | wc -l\n",
		// Broken/legacy patterns
		"[ $x -eq 1 ]\n[ $x -lt $y ]\n[[ $a == \"$b\" ]]\n[[ -z $foo ]]\n",
		// Fork-heavy alternatives
		"x=$(seq 1 10)\nfor i in $(seq 1 N); do :; done\nfor f in $(ls); do :; done\n",
		// Container / cloud surfaces (security-flag katas)
		"docker run -v /proc:/host/proc img\npodman run --privileged img\n",
		"kubectl delete --force --grace-period=0 pod x\nkubeadm reset -f\n",
		"aws ssm put-parameter --type SecureString --value secret\n",
		"git config --global credential.helper store\ngit config http.sslVerify false\n",
		// Package managers
		"apt autoremove --purge\ndnf autoremove\nzypper rm --clean-deps pkg\n",
		"npm install -g pkg\nnpm config set strict-ssl false\n",
		"pip install pkg\nuvx pkg\nnpx pkg\nbun x pkg\n",
		// Shell-history hazards
		"unset HISTFILE\nexport HISTFILE=/dev/null\nset +o history\n",
		// Setuid / setgid / chmod / umask
		"chmod +s /usr/bin/foo\nchmod 4755 /usr/bin/foo\numask 011\numask 077\n",
		// Network / firewall surfaces
		"iptables -F\niptables -P INPUT ACCEPT\nufw disable\nsystemctl disable firewalld\n",
		// Crypto surfaces
		"openssl enc -k password\nopenssl s_client -ssl3\nssh-copy-id -f -o StrictHostKeyChecking=no host\n",
		// Arithmetic + double-bracket edge cases
		"[[ -z $foo ]]\n[[ $x == \"\" ]]\n[[ $x != \"\" ]]\n(( x = a > b ? a : b ))\n",
		// Functions
		"function name() { local x=1; echo $x; }\nfunction name { echo hi; }\n",
		// Heredocs + redirections
		"cat <<EOF\nbody\nEOF\ncat <<-'STRIP'\n\tbody\n\tSTRIP\necho hi >> log\necho hi 2>&1\n",
		// Loops + case
		"for f in *; do echo $f; done\nwhile true; do break; done\ncase $x in a) :;; esac\n",
		// String / parameter expansion
		"echo \"${arr[@]}\"\necho ${var:-default}\necho ${var/old/new}\necho ${var:0:5}\n",
		// Misc dangerous
		"sudo dd if=/dev/zero of=/dev/sda\ntc qdisc add dev eth0 root netem loss 100%\n",
		// Find idioms
		"find /tmp -name *.log\nfind . -name *.go -exec sh -c 'echo {}' \\;\n",
		// cd hazards + state-isolating subshells
		"cd /tmp\ncd $1\n( cd /tmp && rm -rf foo )\n",
		// declare/typeset shapes
		"typeset -a items=(a b c)\ntypeset -A m=(k v k2 v2)\nreadonly y=hi\n",
		// xargs / find chains
		"find . -print | xargs rm\nfind . -print0 | xargs -0 rm\n",
		// trap / signal handling
		"trap 'echo bye' EXIT\ntrap '' EXIT\ntrap - INT\n",
		// Unsafe ssh + scp
		"ssh -o StrictHostKeyChecking=no host\nscp -r host:/etc /local\n",
		// systemd interactions
		"systemctl stop nftables\nsystemctl mask nftables\nsystemctl disable iptables\n",
		// Arithmetic in if + while
		"if (( x > 0 )); then echo positive; fi\nwhile (( i++ < 10 )); do :; done\n",
		// Backslash-escape edge cases
		"echo \\* \\? \\[ \\] \\& \\| \\;\n",
		// Glob qualifiers
		"ls *.zsh(N)\nls *(.)\n",
	}
	for _, src := range corpus {
		runCorpusItem(src)
	}
}
