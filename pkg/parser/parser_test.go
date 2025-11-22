package parser

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements)) 
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `return 5;
return 10;
return 993322;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ht, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ht.Value != "foobar" {
		t.Errorf("ht.Value not %s. got=%s", "foobar", ht.Value)
	}

	if ht.TokenLiteral() != "foobar" {
		t.Errorf("ht.TokenLiteral not %s. got=%s", "foobar",
			ht.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{" -15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
        // Changed test to be valid Zsh arithmetic syntax
        {"let val = (5 + 5) * 2", "let val = ((5 + 5) * 2);"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		        {"-(5 + 5)", "(-(5 + 5))"},
		        		{"a + add( (b * c) ) + d", "((a + add((b * c))) + d)"},
		        		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8) )", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		        		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, tt := range tests {
		t.Logf("Testing input: %q", tt.input)
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfStatement(t *testing.T) {
	input := "if ((1 < 2)); then return true; fi"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.Condition.Statements) != 1 {
		t.Fatalf("stmt.Condition.Statements does not contain 1 statement. got=%d",
			len(stmt.Condition.Statements))
	}

	condStmt, ok := stmt.Condition.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Condition.Statements[0] is not ast.ExpressionStatement. got=%T",
			stmt.Condition.Statements[0])
	}

	if !testInfixExpression(t, condStmt.Expression, 1, "<", 2) {
		return
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Fatalf("consequence is not 1 statement. got=%d", len(stmt.Consequence.Statements))
	}

	consequence, ok := stmt.Consequence.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.ReturnStatement. got=%T", stmt.Consequence.Statements[0])
	}

	if !testLiteralExpression(t, consequence.ReturnValue, true) {
		return
	}

	if stmt.Alternative != nil {
		t.Errorf("stmt.Alternative was not nil. got=%+v", stmt.Alternative)
	}
}

func TestIfElseStatement(t *testing.T) {
	input := "if ((1 > 2)); then return true; else return false; fi"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.Condition.Statements) != 1 {
		t.Fatalf("stmt.Condition.Statements does not contain 1 statement. got=%d",
			len(stmt.Condition.Statements))
	}

	condStmt, ok := stmt.Condition.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Condition.Statements[0] is not ast.ExpressionStatement. got=%T",
			stmt.Condition.Statements[0])
	}

	if !testInfixExpression(t, condStmt.Expression, 1, ">", 2) {
		return
	}
	if len(stmt.Consequence.Statements) != 1 {
		t.Fatalf("consequence is not 1 statement. got=%d", len(stmt.Consequence.Statements))
	}

	consequence, ok := stmt.Consequence.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.ReturnStatement. got=%T", stmt.Consequence.Statements[0])
	}

	if !testLiteralExpression(t, consequence.ReturnValue, true) {
		return
	}

	if len(stmt.Alternative.Statements) != 1 {
		t.Fatalf("alternative is not 1 statement. got=%d", len(stmt.Alternative.Statements))
	}

	alternative, ok := stmt.Alternative.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Alternative.Statements[0] is not ast.ReturnStatement. got=%T", stmt.Alternative.Statements[0])
	}

	if !testLiteralExpression(t, alternative.ReturnValue, false) {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "function(x, y) { x + y; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"function() {}", []string{}},
		{"function(x) {}", []string{"x"}},
		{"function(x, y, z) {}", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("len parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong number of arguments. want=3, got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCommandSubstitutionWithArrayAccess(t *testing.T) {

	input := "`${my_array[1]}`"

	l := lexer.New(input)

	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {

		t.Fatalf("program.Statements does not contain 1 statement. got=%d",

			len(program.Statements))

	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",

			program.Statements[0])

	}

	cs, ok := stmt.Expression.(*ast.CommandSubstitution)

	if !ok {

		t.Fatalf("stmt.Expression is not ast.CommandSubstitution. got=%T",

			stmt.Expression)

	}

	idxExp, ok := cs.Command.(*ast.ArrayAccess)

	if !ok {

		t.Fatalf("cs.Command is not ast.ArrayAccess. got=%T", cs.Command)

	}

	if !testIdentifier(t, idxExp.Left, "my_array") {

		return

	}

	if !testIntegerLiteral(t, idxExp.Index, 1) {

		return

	}

}

func TestIndexExpression(t *testing.T) {

	tests := []struct {

		input         string

		expectedLeft  string

		expectedIndex interface{}
	}{

		{"my_array[1]", "my_array", 1},

		{"users[id]", "users", "id"},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)

		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {

			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",

				program.Statements[0])

		}

		idxExp, ok := stmt.Expression.(*ast.IndexExpression)

		if !ok {

			t.Fatalf("stmt.Expression is not ast.IndexExpression. got=%T",

				stmt.Expression)

		}

		if !testIdentifier(t, idxExp.Left, tt.expectedLeft) {

			return

		}

		if !testLiteralExpression(t, idxExp.Index, tt.expectedIndex) {

			return

		}

	}

}

func TestArrayAccessDollarLbrace(t *testing.T) {

	tests := []struct {

		input         string

		expectedLeft  string

		expectedIndex interface{}
	}{

		{"${my_array[1]}", "my_array", 1},

		{"${users[id]}", "users", "id"},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)

		p := New(l)

		program := p.ParseProgram()

		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {

			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",

				program.Statements[0])

		}

		aa, ok := stmt.Expression.(*ast.ArrayAccess)

		if !ok {

			t.Fatalf("stmt.Expression is not ast.ArrayAccess. got=%T",

				stmt.Expression)

		}

		if !testIdentifier(t, aa.Left, tt.expectedLeft) {

			return

		}

		if !testLiteralExpression(t, aa.Index, tt.expectedIndex) {

			return

		}

	}

}

func TestForLoopStatement(t *testing.T) {

	input := "for ((i=0;i<10;i++)) do echo $i; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {

		t.Fatalf("program.Statements does not contain 1 statement. got=%d",

			len(program.Statements))

	}

	stmt, ok := program.Statements[0].(*ast.ForLoopStatement)

	if !ok {

		t.Fatalf("program.Statements[0] is not ast.ForLoopStatement. got=%T",

			program.Statements[0])

	}

	if stmt.TokenLiteral() != "for" {

		t.Errorf("stmt.TokenLiteral not 'for', got %q", stmt.TokenLiteral())

	}

	// Test Init expression

	testInfixExpression(t, stmt.Init, "i", "=", 0)

	// Test Condition expression

	testInfixExpression(t, stmt.Condition, "i", "<", 10)

	// Test Post expression

	if !testPostfixExpression(t, stmt.Post, "i", "++") {

		return

	}

	// Test Body statement

	if len(stmt.Body.Statements) != 1 {

		for i, s := range stmt.Body.Statements {

			t.Logf("stmt %d: %T %s", i, s, s.String())

		}

		t.Fatalf("ForLoopStatement.Body.Statements does not contain 1 statement. got=%d", len(stmt.Body.Statements))

	}

	bodyStmt, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {

		t.Fatalf("ForLoopStatement.Body.Statements[0] is not ast.ExpressionStatement. got=%T", bodyStmt)

	}

	cmd, ok := bodyStmt.Expression.(*ast.SimpleCommand)

	if !ok {

		t.Fatalf("bodyStmt.Expression is not *ast.SimpleCommand. got=%T", bodyStmt.Expression)

	}

	if cmd.Name.String() != "echo" {

		t.Errorf("cmd.Name is not 'echo'. got=%q", cmd.Name.String())

	}

	if len(cmd.Arguments) != 1 {

		t.Fatalf("cmd.Arguments does not contain 1 argument. got=%d", len(cmd.Arguments))

	}

	if cmd.Arguments[0].String() != "$i" {

		t.Errorf("cmd.Arguments[0] is not '$i'. got=%q", cmd.Arguments[0].String())

	}

}

func testPostfixExpression(t *testing.T, exp ast.Expression, left string, operator string) bool {
	postfixExp, ok := exp.(*ast.PostfixExpression)
	if !ok {
		t.Errorf("exp not *ast.PostfixExpression. got=%T", exp)
		return false
	}

	if !testIdentifier(t, postfixExp.Left, left) {
		return false
	}

	if postfixExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%s", operator, postfixExp.Operator)
		return false
	}

	return true
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("exp not *ast.PrefixExpression. got=%T", exp)
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%s", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestForLoopStatementStub(t *testing.T) { // Renamed from TestForLoopStatement, now a stub for later use.

	input := "for ((i=0;i<10;i++)) do echo $i; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {

		t.Fatalf("program.Statements does not contain 1 statement. got=%d",

			len(program.Statements))

	}

	stmt, ok := program.Statements[0].(*ast.ForLoopStatement)

	if !ok {

		t.Fatalf("program.Statements[0] is not ast.ForLoopStatement. got=%T",

			program.Statements[0])

	}

	if stmt.TokenLiteral() != "for" {

		t.Errorf("stmt.TokenLiteral not 'for', got %q", stmt.TokenLiteral())

	}

}
