// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

// TestArgValueAfter exercises the dispatch helper used by katas
// matching `--flag value` style options.
func TestArgValueAfter(t *testing.T) {
	mkArg := func(s string) ast.Expression {
		return &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: s}, Value: s}
	}
	cmd := &ast.SimpleCommand{Arguments: []ast.Expression{mkArg("--bind"), mkArg("0.0.0.0"), mkArg("--port"), mkArg("8080")}}
	if got := ArgValueAfter(cmd, "--bind"); got != "0.0.0.0" {
		t.Errorf("--bind: got %q want 0.0.0.0", got)
	}
	if got := ArgValueAfter(cmd, "--port"); got != "8080" {
		t.Errorf("--port: got %q want 8080", got)
	}
	if got := ArgValueAfter(cmd, "--missing"); got != "" {
		t.Errorf("--missing: got %q want \"\"", got)
	}
	tail := &ast.SimpleCommand{Arguments: []ast.Expression{mkArg("--bind")}}
	if got := ArgValueAfter(tail, "--bind"); got != "" {
		t.Errorf("tail-key: got %q want \"\"", got)
	}
}

// TestIsAlphaNumeric exercises the byte-class predicate behind the
// kata1071 self-reference detector.
func TestIsAlphaNumeric(t *testing.T) {
	cases := []struct {
		b    byte
		want bool
	}{
		{'a', true},
		{'Z', true},
		{'0', true},
		{'9', true},
		{'_', false},
		{'-', false},
		{' ', false},
		{'\n', false},
	}
	for _, tc := range cases {
		if got := isAlphaNumeric(tc.b); got != tc.want {
			t.Errorf("isAlphaNumeric(%q): got %v want %v", tc.b, got, tc.want)
		}
	}
}

// TestIsDevNullAndStringValue exercises the ZC1053 string extraction
// over StringLiteral and ConcatenatedExpression node shapes, plus the
// quoted-literal `/dev/null` detector.
func TestIsDevNullAndStringValue(t *testing.T) {
	bare := &ast.StringLiteral{Value: "/dev/null"}
	if !isDevNull(bare) {
		t.Errorf("isDevNull(bare /dev/null): expected true")
	}
	quoted := &ast.StringLiteral{Value: "\"/dev/null\""}
	if !isDevNull(quoted) {
		t.Errorf("isDevNull(quoted /dev/null): expected true")
	}
	other := &ast.StringLiteral{Value: "/tmp/log"}
	if isDevNull(other) {
		t.Errorf("isDevNull(/tmp/log): expected false")
	}
	concat := &ast.ConcatenatedExpression{Parts: []ast.Expression{
		&ast.StringLiteral{Value: "/dev"},
		&ast.StringLiteral{Value: "/null"},
	}}
	if got := getStringValueZC1053(concat); got != "/dev/null" {
		t.Errorf("getStringValueZC1053(concat): got %q want /dev/null", got)
	}
	if got := getStringValueZC1053(&ast.Identifier{Value: "x"}); got != "" {
		t.Errorf("getStringValueZC1053(non-string): got %q want \"\"", got)
	}
}

// TestZC1796HasPgArg covers every flag in the pg_dump / pg_restore
// indicator set plus the negative case.
func TestZC1796HasPgArg(t *testing.T) {
	mk := func(args ...string) *ast.SimpleCommand {
		exprs := make([]ast.Expression, 0, len(args))
		for _, a := range args {
			exprs = append(exprs, &ast.Identifier{Value: a})
		}
		return &ast.SimpleCommand{Arguments: exprs}
	}
	flags := []string{"-d", "--dbname", "-F", "--format", "-U", "--username", "--if-exists", "--no-owner", "--no-acl"}
	for _, f := range flags {
		if !zc1796HasPgArg(mk(f, "value")) {
			t.Errorf("zc1796HasPgArg(%s): expected true", f)
		}
	}
	if zc1796HasPgArg(mk("--unrelated")) {
		t.Errorf("zc1796HasPgArg(--unrelated): expected false")
	}
	if zc1796HasPgArg(mk()) {
		t.Errorf("zc1796HasPgArg(empty): expected false")
	}
}

// TestZC1960IsGcloudSshCmd covers the head-prefix gate and the
// trailing `--command` / `--command=` variants.
func TestZC1960IsGcloudSshCmd(t *testing.T) {
	mk := func(args ...string) *ast.SimpleCommand {
		exprs := make([]ast.Expression, 0, len(args))
		for _, a := range args {
			exprs = append(exprs, &ast.Identifier{Value: a})
		}
		return &ast.SimpleCommand{Arguments: exprs}
	}
	if !zc1960IsGcloudSshCmd(mk("compute", "ssh", "host", "--command", "uptime")) {
		t.Errorf("expected true for --command separate")
	}
	if !zc1960IsGcloudSshCmd(mk("compute", "ssh", "host", "--command=uptime")) {
		t.Errorf("expected true for --command=")
	}
	if zc1960IsGcloudSshCmd(mk("compute", "ssh", "host")) {
		t.Errorf("expected false without --command")
	}
	if zc1960IsGcloudSshCmd(mk("compute", "scp", "f1", "f2")) {
		t.Errorf("expected false for compute scp")
	}
	if zc1960IsGcloudSshCmd(mk("compute")) {
		t.Errorf("expected false for too few args")
	}
}

// TestZC1045StringHasSub exercises every branch of the embedded-
// substitution detector over double-quoted string literals.
func TestZC1045StringHasSub(t *testing.T) {
	cases := []struct {
		name string
		val  string
		want bool
	}{
		{"empty", ``, false},
		{"single-byte", `"`, false},
		{"unquoted", `hello`, false},
		{"plain quoted", `"hello"`, false},
		{"backtick sub", "\"echo `date`\"", true},
		{"dollar paren sub", `"echo $(date)"`, true},
		{"escaped backtick", "\"echo \\` not sub\"", false},
		{"escaped dollar", `"\$(not-sub)"`, false},
		{"dollar without paren", `"$var"`, false},
	}
	for _, tc := range cases {
		if got := zc1045StringHasSub(tc.val); got != tc.want {
			t.Errorf("%s: got %v want %v", tc.name, got, tc.want)
		}
	}
}

// TestZC1045ConcatHasSub exercises the array-walk path. A
// CommandSubstitution node anywhere in Parts trips the detector.
func TestZC1045ConcatHasSub(t *testing.T) {
	plain := &ast.ConcatenatedExpression{Parts: []ast.Expression{
		&ast.StringLiteral{Value: "a"},
		&ast.StringLiteral{Value: "b"},
	}}
	if zc1045ConcatHasSub(plain) {
		t.Errorf("expected false for plain literal concat")
	}
	withSub := &ast.ConcatenatedExpression{Parts: []ast.Expression{
		&ast.StringLiteral{Value: "prefix"},
		&ast.CommandSubstitution{Command: &ast.SimpleCommand{}},
	}}
	if !zc1045ConcatHasSub(withSub) {
		t.Errorf("expected true when CommandSubstitution present")
	}
}

// TestZC1071SelfReferences exercises every node-type branch of the
// self-reference detector used by ZC1071.
func TestZC1071SelfReferences(t *testing.T) {
	id := &ast.Identifier{Value: "foo"}
	if !zc1071SelfReferences(&ast.ArrayAccess{Left: id}, "foo") {
		t.Errorf("ArrayAccess(foo) referencing foo: expected true")
	}
	if zc1071SelfReferences(&ast.ArrayAccess{Left: id}, "bar") {
		t.Errorf("ArrayAccess(foo) referencing bar: expected false")
	}
	if !zc1071SelfReferences(&ast.Identifier{Value: "$foo"}, "foo") {
		t.Errorf("Identifier($foo) referencing foo: expected true")
	}
	if !zc1071SelfReferences(&ast.Identifier{Value: "${foo}"}, "foo") {
		t.Errorf("Identifier(${foo}) referencing foo: expected true")
	}
	prefix := &ast.PrefixExpression{Operator: "$", Right: id}
	if !zc1071SelfReferences(prefix, "foo") {
		t.Errorf("PrefixExpression($foo) referencing foo: expected true")
	}
	wrong := &ast.PrefixExpression{Operator: "!", Right: id}
	if zc1071SelfReferences(wrong, "foo") {
		t.Errorf("PrefixExpression(!) referencing foo: expected false")
	}
	if zc1071SelfReferences(&ast.IntegerLiteral{Value: 1}, "foo") {
		t.Errorf("IntegerLiteral: expected false")
	}
}

// TestCommandIdentifier covers the head-identifier helper used by
// every Check entry point.
func TestCommandIdentifier(t *testing.T) {
	cmd := &ast.SimpleCommand{Name: &ast.Identifier{Value: "echo"}}
	if got := CommandIdentifier(cmd); got != "echo" {
		t.Errorf("ident head: got %q want echo", got)
	}
	notIdent := &ast.SimpleCommand{Name: &ast.StringLiteral{Value: "x"}}
	if got := CommandIdentifier(notIdent); got != "" {
		t.Errorf("non-ident head: got %q want \"\"", got)
	}
}
