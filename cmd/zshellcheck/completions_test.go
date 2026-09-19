// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"bytes"
	"flag"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestRunCompletions(t *testing.T) {
	want := readTestFile(t, "../../completions/zsh/_zshellcheck")
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, ".zshellcheckrc"), "invalid configuration\n")
	t.Chdir(dir)
	for _, name := range []string{"-completions", "--completions"} {
		code, out, errOut := runCaptured(t, name, "-cpuprofile", filepath.Join(dir, "missing", "cpu.out"))
		if code != 0 || out != want || errOut != "" {
			t.Errorf("%s: exit=%d stdout matches=%v stderr=%q", name, code, out == want, errOut)
		}
	}
}

func TestPrintCompletionsFailure(t *testing.T) {
	var errOut bytes.Buffer
	if code := printCompletions(failingWriter{}, &errOut); code != 1 || !strings.Contains(errOut.String(), "Error writing completions") {
		t.Errorf("failed export: exit=%d stderr=%q", code, errOut.String())
	}
}

func TestCompletionFlagParity(t *testing.T) {
	oldFlags := flag.CommandLine
	t.Cleanup(func() { flag.CommandLine = oldFlags })
	resetFlags()
	registerRunFlags()
	zsh := shellOutput(t, "zsh", `fpath=("$1" $fpath); autoload -Uz _zshellcheck; _arguments() { print -rl -- "$@"; }; _zshellcheck`, "../../completions/zsh")
	bash := bashCompletions(t, "-")
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		for _, prefix := range []string{"-", "--"} {
			name := prefix + f.Name
			if !slices.Contains(bash, name) {
				t.Errorf("Bash completion missing %s", name)
			}
			if !strings.Contains(zsh, "\n"+name+"[") && !strings.Contains(zsh, "\n"+name+"=[") {
				t.Errorf("Zsh completion missing %s", name)
			}
		}
	})
}

func TestBashCompletionValues(t *testing.T) {
	tests := []struct {
		words []string
		want  []string
	}{
		{[]string{"--format", ""}, []string{"text", "json", "sarif"}},
		{[]string{"-format", "j"}, []string{"json"}},
		{[]string{"--format=j"}, []string{"json"}},
		{[]string{"--severity", "error,w"}, []string{"error,warning"}},
		{[]string{"-severity=error,w"}, []string{"error,warning"}},
		{[]string{"--baseline", ""}, []string{"<file>"}},
		{[]string{"--cpuprofile="}, []string{"<file>"}},
		{[]string{"--", "--format", ""}, []string{"<file>"}},
		{[]string{"script.zsh", "--format", ""}, []string{"<file>"}},
		{[]string{"-", "--format", ""}, []string{"<file>"}},
		{[]string{"--baseline", "--", "--format", "j"}, []string{"json"}},
		{[]string{"--explain", ""}, nil},
		{[]string{"--unknown="}, nil},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.words, " "), func(t *testing.T) {
			if got := bashCompletions(t, tt.words...); !slices.Equal(got, tt.want) {
				t.Errorf("completion=%q, want %q", got, tt.want)
			}
		})
	}
}

func bashCompletions(t *testing.T, words ...string) []string {
	t.Helper()
	return bashCompletionsAt(t, "", words...)
}

func bashCompletionsAt(t *testing.T, replacement string, words ...string) []string {
	t.Helper()
	// Supply the bash-completion callbacks without requiring a system
	// installation; the completion function itself runs unchanged.
	script := `source "$1"; replacement=$2; shift 2
_init_completion() { words=("${COMP_WORDS[@]}"); cword=$COMP_CWORD; cur=${words[cword]}; }
_filedir() { COMPREPLY=('<file>'); }
COMP_WORDS=(zshellcheck "$@"); COMP_CWORD=$((${#COMP_WORDS[@]} - 1))
_zshellcheck zshellcheck "$replacement" zshellcheck
((${#COMPREPLY[@]})) && printf '%s\n' "${COMPREPLY[@]}"
exit 0`
	args := append([]string{"../../completions/bash/zshellcheck-completion.bash", replacement}, words...)
	out := strings.TrimSpace(shellOutput(t, "bash", script, args...))
	if out == "" {
		return nil
	}
	return strings.Split(out, "\n")
}

func TestBashCompletionReplacementWord(t *testing.T) {
	for _, tt := range []struct {
		input, replacement, want string
	}{
		{"--format=j", "--format=j", "--format=json"},
		{"--format=j", "j", "json"},
		{`--baseline="file w`, "file w", "<file>"},
		{`--baseline='file w`, "file w", "<file>"},
		{`--baseline="file"`, `--baseline="file"`, "--baseline=<file>"},
	} {
		t.Run(tt.input+"/"+tt.replacement, func(t *testing.T) {
			got := bashCompletionsAt(t, tt.replacement, tt.input)
			if !slices.Equal(got, []string{tt.want}) {
				t.Errorf("replacement=%q, want %q", got, tt.want)
			}
		})
	}
}

func TestZshCompletionOptionBoundary(t *testing.T) {
	script := `fpath=("$1" $fpath); shift; autoload -Uz _zshellcheck
_arguments() { print options; }
_files() { print files; }
words=(zshellcheck "$@"); CURRENT=${#words}; _zshellcheck`
	for _, tt := range []struct {
		words []string
		want  string
	}{
		{[]string{"--", "--fi"}, "files\n"},
		{[]string{"--baseline", "--", "--fi"}, "options\n"},
		{[]string{"--explain", "--", "--fi"}, "options\n"},
		{[]string{"script.zsh", "--fi"}, "files\n"},
		{[]string{"-", "--fi"}, "files\n"},
		{[]string{"--format=json", "--fi"}, "options\n"},
	} {
		t.Run(strings.Join(tt.words, " "), func(t *testing.T) {
			args := append([]string{"../../completions/zsh"}, tt.words...)
			if got := shellOutput(t, "zsh", script, args...); got != tt.want {
				t.Errorf("completion dispatch=%q, want %q", got, tt.want)
			}
		})
	}
}

func shellOutput(t *testing.T, shell, script string, args ...string) string {
	t.Helper()
	bin, err := exec.LookPath(shell)
	if err != nil {
		t.Skipf("%s is unavailable", shell)
	}
	// Keep quotes in test inputs out of Windows command-line argument parsing.
	var input strings.Builder
	input.WriteString("set --")
	for _, arg := range args {
		input.WriteString(" '")
		input.WriteString(strings.ReplaceAll(arg, "'", `'\''`))
		input.WriteByte('\'')
	}
	input.WriteByte('\n')
	input.WriteString(script)
	command := exec.Command(bin, "-s")
	command.Stdin = strings.NewReader(input.String())
	out, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v\n%s", shell, err, out)
	}
	return string(out)
}

func TestShellOutputArguments(t *testing.T) {
	args := []string{"", `--baseline="file"`, "file's name", `path\with\slashes\`, "$literal", "two\nlines"}
	want := strings.Join(args, "\x00") + "\x00"
	if got := shellOutput(t, "bash", `printf '%s\000' "$@"`, args...); got != want {
		t.Errorf("shell arguments=%q, want %q", got, want)
	}
}
