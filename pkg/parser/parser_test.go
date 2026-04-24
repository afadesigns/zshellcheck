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
		{"let _ = -a * b", "let _ = ((-a) * b);"},
		{"let _ = !-a", "let _ = (!(-a));"},
		{"let _ = a + b + c", "let _ = ((a + b) + c);"},
		{"let _ = a + b - c", "let _ = ((a + b) - c);"},
		{"let _ = a * b * c", "let _ = ((a * b) * c);"},
		{"let _ = a * b / c", "let _ = ((a * b) / c);"},
		{"let _ = a + b / c", "let _ = (a + (b / c));"},
		{"let _ = a + b * c + d / e - f", "let _ = (((a + (b * c)) + (d / e)) - f);"},
		{"let _ = 3 + 4; let _ = -5 * 5", "let _ = (3 + 4);let _ = ((-5) * 5);"},
		// No, -5 is Prefix(-, 5). - is command starter?
		// p.peekTokenIs(token.MINUS) -> SimpleCommand.
		// So -5 * 5 -> Command -5 args * 5.
		// So this one might still fail if checked as ExprStmt.
		// But `3 + 4` is inside `let`.
		// `let _ = 3 + 4; -5 * 5`.
		// Stmt 1: Let. Stmt 2: -5 * 5.
		// We can wrap the second part too? `let _ = -5 * 5`?
		// Input: "let _ = 3 + 4; let _ = -5 * 5".

		{"let _ = 5 > 4 == 3 < 4", "let _ = ((5 > 4) == (3 < 4));"},
		{"let _ = 5 < 4 != 3 > 4", "let _ = ((5 < 4) != (3 > 4));"},
		{"let _ = 3 + 4 * 5 == 3 * 1 + 4 * 5", "let _ = ((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));"},
		{"let _ = true", "let _ = true;"},
		{"let _ = false", "let _ = false;"},
		{"let _ = 3 > 5 == false", "let _ = ((3 > 5) == false);"},
		{"let _ = 3 < 5 == true", "let _ = ((3 < 5) == true);"},
		{"let _ = 1 + (2 + 3) + 4", "let _ = ((1 + ((2 + 3))) + 4);"},
		{"let val = (5 + 5) * 2", "let val = (((5 + 5)) * 2);"},
		{"let _ = 2 / (5 + 5)", "let _ = (2 / ((5 + 5)));"},
		{"let _ = -(5 + 5)", "let _ = (-((5 + 5)));"},
		{"let _ = a + add( (b * c) ) + d", "let _ = ((a + add(((b * c)))) + d);"},
		{"let _ = add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8) )", "let _ = add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)));"},
		{"let _ = add(a + b + c * d / f + g)", "let _ = add((((a + b) + ((c * d) / f)) + g));"},
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

	if len(stmt.Condition.(*ast.BlockStatement).Statements) != 1 {
		t.Fatalf("stmt.Condition.(*ast.BlockStatement).Statements does not contain 1 statement. got=%d",
			len(stmt.Condition.(*ast.BlockStatement).Statements))
	}

	condStmt, ok := stmt.Condition.(*ast.BlockStatement).Statements[0].(*ast.ArithmeticCommand)
	if !ok {
		t.Fatalf("stmt.Condition.(*ast.BlockStatement).Statements[0] is not ast.ArithmeticCommand. got=%T",
			stmt.Condition.(*ast.BlockStatement).Statements[0])
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

	if len(stmt.Condition.(*ast.BlockStatement).Statements) != 1 {
		t.Fatalf("stmt.Condition.(*ast.BlockStatement).Statements does not contain 1 statement. got=%d",
			len(stmt.Condition.(*ast.BlockStatement).Statements))
	}

	condStmt, ok := stmt.Condition.(*ast.BlockStatement).Statements[0].(*ast.ArithmeticCommand)
	if !ok {
		t.Fatalf("stmt.Condition.(*ast.BlockStatement).Statements[0] is not ast.ArithmeticCommand. got=%T",
			stmt.Condition.(*ast.BlockStatement).Statements[0])
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

	if len(stmt.Alternative.(*ast.BlockStatement).Statements) != 1 {
		t.Fatalf("alternative is not 1 statement. got=%d", len(stmt.Alternative.(*ast.BlockStatement).Statements))
	}

	alternative, ok := stmt.Alternative.(*ast.BlockStatement).Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Alternative.Statements[0] is not ast.ReturnStatement. got=%T", stmt.Alternative.(*ast.BlockStatement).Statements[0])
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

	if len(function.Params) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Params))
	}

	testLiteralExpression(t, function.Params[0], "x")
	testLiteralExpression(t, function.Params[1], "y")

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

		if len(function.Params) != len(tt.expectedParams) {
			t.Errorf("len parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Params))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Params[i], ident)
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

/*
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
*/

func TestIndexExpression(t *testing.T) {
	tests := []struct {
		input string

		expectedLeft string

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
		input string

		expectedLeft string

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

func TestWhileLoopStatement(t *testing.T) {
	input := "while [[ $x -lt 10 ]]; do echo $x; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileLoopStatement)
	if !ok {
		t.Fatalf("not ast.WhileLoopStatement. got=%T", program.Statements[0])
	}
	if stmt.TokenLiteral() != "while" {
		t.Errorf("TokenLiteral not 'while', got %q", stmt.TokenLiteral())
	}
	if stmt.Body == nil {
		t.Fatal("Body is nil")
	}
}

func TestSelectStatement(t *testing.T) {
	input := "select item in a b c; do echo $item; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.SelectStatement)
	if !ok {
		t.Fatalf("not ast.SelectStatement. got=%T", program.Statements[0])
	}
	if stmt.Name.Value != "item" {
		t.Errorf("Name not 'item', got %q", stmt.Name.Value)
	}
	if len(stmt.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(stmt.Items))
	}
}

func TestSelectStatementNoIn(t *testing.T) {
	input := "select item; do echo $item; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.SelectStatement)
	if !ok {
		t.Fatalf("not ast.SelectStatement. got=%T", program.Statements[0])
	}
	if stmt.Name.Value != "item" {
		t.Errorf("Name not 'item', got %q", stmt.Name.Value)
	}
}

func TestCoprocStatement(t *testing.T) {
	input := "coproc myproc { echo hello; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.CoprocStatement)
	if !ok {
		t.Fatalf("not ast.CoprocStatement. got=%T", program.Statements[0])
	}
	if stmt.Name != "myproc" {
		t.Errorf("Name not 'myproc', got %q", stmt.Name)
	}
	if stmt.Command == nil {
		t.Fatal("Command is nil")
	}
}

func TestCoprocStatementNoName(t *testing.T) {
	input := "coproc { echo hello; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.CoprocStatement)
	if !ok {
		t.Fatalf("not ast.CoprocStatement. got=%T", program.Statements[0])
	}
	if stmt.Name != "" {
		t.Errorf("expected empty Name, got %q", stmt.Name)
	}
}

func TestCaseStatement(t *testing.T) {
	input := `case $x in
a) echo a;;
b|c) echo bc;;
esac`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.CaseStatement)
	if !ok {
		t.Fatalf("not ast.CaseStatement. got=%T", program.Statements[0])
	}
	if len(stmt.Clauses) != 2 {
		t.Fatalf("expected 2 clauses, got %d", len(stmt.Clauses))
	}
	// Second clause should have 2 patterns
	if len(stmt.Clauses[1].Patterns) != 2 {
		t.Errorf("expected 2 patterns in second clause, got %d", len(stmt.Clauses[1].Patterns))
	}
}

func TestDeclarationStatement(t *testing.T) {
	input := "typeset -a myarray"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.DeclarationStatement)
	if !ok {
		t.Fatalf("not ast.DeclarationStatement. got=%T", program.Statements[0])
	}
	if stmt.Command != "typeset" {
		t.Errorf("Command not 'typeset', got %q", stmt.Command)
	}
}

func TestDeclarationStatementWithAssignment(t *testing.T) {
	input := "declare x=5"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.DeclarationStatement)
	if !ok {
		t.Fatalf("not ast.DeclarationStatement. got=%T", program.Statements[0])
	}
	if len(stmt.Assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(stmt.Assignments))
	}
	if stmt.Assignments[0].Name.Value != "x" {
		t.Errorf("expected assignment name 'x', got %q", stmt.Assignments[0].Name.Value)
	}
}

func TestDeclarationStatementWithAppend(t *testing.T) {
	input := "declare x+=5"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.DeclarationStatement)
	if !ok {
		t.Fatalf("not ast.DeclarationStatement. got=%T", program.Statements[0])
	}
	if len(stmt.Assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(stmt.Assignments))
	}
	if !stmt.Assignments[0].IsAppend {
		t.Error("expected IsAppend to be true")
	}
}

func TestDeclarationStatementWithArrayValue(t *testing.T) {
	input := "declare x=(a b c)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.DeclarationStatement)
	if !ok {
		t.Fatalf("not ast.DeclarationStatement. got=%T", program.Statements[0])
	}
	if len(stmt.Assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(stmt.Assignments))
	}
}

func TestArithmeticCommand(t *testing.T) {
	input := "((x + 1))"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ArithmeticCommand)
	if !ok {
		t.Fatalf("not ast.ArithmeticCommand. got=%T", program.Statements[0])
	}
	if stmt.Expression == nil {
		t.Fatal("Expression is nil")
	}
}

func TestArithmeticParamExistenceBare(t *testing.T) {
	// `(( $+name ))` and `(( $+name[key] ))` — bare `$+...` inside arithmetic.
	// Regression test for issue #1047.
	cases := []string{
		`(( $+commands[ls] ))`,
		`(( $+foo ))`,
		`(( $+arr[key] ))`,
	}
	for _, input := range cases {
		l := lexer.New(input)
		p := New(l)
		_ = p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Errorf("%s: unexpected parser errors: %v", input, errs)
		}
	}
}

func TestArithmeticLogicalChain(t *testing.T) {
	// `(( A )) && (( B ))`, `(( A )) || (( B ))`, and mixed chains.
	// Regression test for issue #1047.
	cases := []string{
		`(( 1 )) && (( 2 ))`,
		`(( 1 )) || (( 2 ))`,
		`(( $+commands[ls] )) && (( $+commands[eza] ))`,
		`(( 1 )) && (( 2 )) || (( 3 ))`,
	}
	for _, input := range cases {
		l := lexer.New(input)
		p := New(l)
		_ = p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Errorf("%s: unexpected parser errors: %v", input, errs)
		}
	}
}

func TestSpecialParametersAndLengthOperator(t *testing.T) {
	// Zsh single-char special parameters: $?, $@, $$, $_. Previously
	// they lexed as DOLLAR + ILLEGAL (or DOLLAR + QUESTION) and the
	// parser's prefix table rejected the tail. They now lex as a
	// single VARIABLE token.
	//
	// $#name is the length-of-parameter operator; the identifier
	// after the `#` is part of the expression, not a fresh token
	// to be left for the caller to stumble over. Real code like
	// `(( $#BUFFER > 0 ))` from zsh-autosuggestions widgets.zsh
	// required both fixes to parse without errors.
	inputs := []string{
		"x=$?",
		"y=$@",
		"z=$$",
		"w=$_",
		"retval=$?",
		"if (( $#BUFFER )); then echo yes; fi",
		"if (( $#BUFFER > 0 )); then echo yes; fi",
		"len=$#ARRAY",
	}
	for _, input := range inputs {
		l := lexer.New(input)
		p := New(l)
		_ = p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Errorf("%s:\n  unexpected parser errors: %v", input, errs)
		}
	}
}

func TestArraySubscriptFlags(t *testing.T) {
	// Zsh subscript flags before the index value: arr[(R)value]
	// (reverse match), arr[(r)pat] (key match), arr[(I)i] (integer),
	// arr[(ri)pat] (composed). These let code search arrays by
	// content; every Zsh plugin in the curated corpus uses at least
	// one form. Before this fix the `(…)` was consumed as a grouped
	// expression and the subject IDENT was rejected with
	// "expected next token to be ], got VARIABLE".
	inputs := []string{
		"a=${arr[(R)val]}",
		"b=${arr[(r)pat]}",
		"c=${arr[(I)i]}",
		"d=${arr[(ri)pat]}",
		"e=${arr[(R)$~pattern]}",
		"f=${history_match_keys[1]}",
		"g=${arr[1]}",
		"h=${arr[(R)${nested}]}",
	}
	for _, input := range inputs {
		l := lexer.New(input)
		p := New(l)
		_ = p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Errorf("%s:\n  unexpected parser errors: %v", input, errs)
		}
	}
}

func TestBraceParamExpansionModifiersAreOpaque(t *testing.T) {
	// Every Zsh parameter-expansion modifier that comes after the
	// subject is accepted without a structured AST yet (tracked in
	// issue #129). The parser walks through the modifier body with
	// brace-depth tracking so the closing `}` is found correctly,
	// unblocking common plugin code (oh-my-zsh, zsh-autosuggestions,
	// prezto) until the richer parameter-expansion node ships.
	inputs := []string{
		"a=${var#prefix}",
		"b=${var##longest}",
		"c=${var%suffix}",
		"d=${var%%longest}",
		"e=${var/old/new}",
		"f=${var//all/repl}",
		"g=${var:-default}",
		"h=${var:=default}",
		"i=${var:+alt}",
		"j=${var:?err}",
		"k=${var:3}",
		"l=${var:3:5}",
		"m=${var:#glob}",
		"n=${nested:-${other:-fallback}}",
		"o=${arr[1]#*:}",
	}
	for _, input := range inputs {
		l := lexer.New(input)
		p := New(l)
		_ = p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Errorf("%s:\n  unexpected parser errors: %v", input, errs)
		}
	}
}

func TestBraceParamExpansionSingleCharFlags(t *testing.T) {
	// Zsh single-char pre-flags inside `${X name}` that modify the
	// expansion: `=` (split), `~` (glob interpret), `^` (rc-style).
	// Real-world code relies on all three; without explicit handling
	// the parser tried to run `=` / `^` through the generic prefix
	// path, which has no entry for ASSIGN / CARET, and produced
	// "no prefix parse function for =" on zsh-autosuggestions.
	inputs := []string{
		`a=(${=VAR})`,
		`b=(${~PATTERN})`,
		`c=(${^ARRAY})`,
		`d=(${=ZSH_AUTOSUGGEST_STRATEGY})`,
	}
	for _, input := range inputs {
		l := lexer.New(input)
		p := New(l)
		_ = p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Errorf("%s:\n  unexpected parser errors: %v", input, errs)
		}
	}
}

func TestDeclarationFollowedByIfOnNextLine(t *testing.T) {
	// Regression: `typeset -g A` directly followed by an `if … then …
	// fi` on the next line used to swallow the `if` keyword, leaving
	// the parser unable to consume the subsequent `then` / `fi`. The
	// body of any well-formed Zsh script inside common plugins hits
	// this pattern (zsh-autosuggestions src/async.zsh among others).
	inputs := []string{
		`foo() {
  typeset -g A B
  if [[ -n "$x" ]]; then echo ok; fi
}`,
		`foo() {
  typeset A
  if [[ 1 -eq 1 ]]; then print hi; fi
}`,
		`foo() {
  local VAR=value
  if [[ -n "$VAR" ]]; then echo set; fi
}`,
		`foo() {
  readonly FLAG
  if true; then echo go; fi
}`,
	}
	for _, input := range inputs {
		l := lexer.New(input)
		p := New(l)
		_ = p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			t.Errorf("%s:\n  unexpected parser errors: %v", input, errs)
		}
	}
}

func TestArithmeticInsideIfWithLogicalChain(t *testing.T) {
	// The original repro from #1047.
	input := `if (( $+commands[ls] )) && (( $+commands[eza] )); then
  echo "ok"
fi`
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
	if errs := p.Errors(); len(errs) != 0 {
		t.Errorf("unexpected parser errors: %v", errs)
	}
}

func TestSubshellStatement(t *testing.T) {
	input := "(echo hello)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.Subshell)
	if !ok {
		t.Fatalf("not ast.Subshell. got=%T", program.Statements[0])
	}
	if stmt.Command == nil {
		t.Fatal("Command is nil")
	}
}

func TestSimpleCommandStatement(t *testing.T) {
	input := "echo hello world"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	cmd, ok := stmt.Expression.(*ast.SimpleCommand)
	if !ok {
		t.Fatalf("not ast.SimpleCommand. got=%T", stmt.Expression)
	}
	if cmd.Name.String() != "echo" {
		t.Errorf("Name not 'echo', got %q", cmd.Name.String())
	}
	if len(cmd.Arguments) != 2 {
		t.Errorf("expected 2 arguments, got %d", len(cmd.Arguments))
	}
}

func TestTestCommandStatement(t *testing.T) {
	input := "test -f /tmp/foo"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	cmd, ok := stmt.Expression.(*ast.SimpleCommand)
	if !ok {
		t.Fatalf("not ast.SimpleCommand. got=%T", stmt.Expression)
	}
	if cmd.Name.String() != "test" {
		t.Errorf("Name not 'test', got %q", cmd.Name.String())
	}
}

func TestBraceGroupStatement(t *testing.T) {
	input := "{ echo hello; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("not ast.BlockStatement. got=%T", program.Statements[0])
	}
}

func TestPipelineStatement(t *testing.T) {
	input := "echo hello | grep foo"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	infix, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("not ast.InfixExpression. got=%T", stmt.Expression)
	}
	if infix.Operator != "|" {
		t.Errorf("Operator not '|', got %q", infix.Operator)
	}
}

func TestCommandListAndOr(t *testing.T) {
	input := "cmd1 && cmd2 || cmd3"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestNegatedPipeline(t *testing.T) {
	input := "! grep foo /tmp/bar"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestRedirectionInCommand(t *testing.T) {
	input := "echo hello > /tmp/out"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestRedirectionAppend(t *testing.T) {
	input := "echo hello >> /tmp/out"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestDollarParenCommandSubstitution(t *testing.T) {
	input := "$(echo hello)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestDollarParenArithmetic(t *testing.T) {
	input := "$((1 + 2))"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestProcessSubstitution(t *testing.T) {
	input := "diff <(sort a) <(sort b)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestForEachLoop(t *testing.T) {
	input := "for item in a b c; do echo $item; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForLoopStatement)
	if !ok {
		t.Fatalf("not ast.ForLoopStatement. got=%T", program.Statements[0])
	}
	if stmt.Name == nil {
		t.Fatal("Name is nil")
	}
	if stmt.Name.Value != "item" {
		t.Errorf("Name not 'item', got %q", stmt.Name.Value)
	}
	if len(stmt.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(stmt.Items))
	}
}

func TestForLoopNoIn(t *testing.T) {
	input := "for item; do echo $item; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForLoopStatement)
	if !ok {
		t.Fatalf("not ast.ForLoopStatement. got=%T", program.Statements[0])
	}
	if stmt.Name == nil {
		t.Fatal("Name is nil")
	}
}

func TestFunctionDefinitionNameParens(t *testing.T) {
	input := "myfunc() { echo hello; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestFunctionDefinitionKeyword(t *testing.T) {
	input := "function myfunc { echo hello; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestDoubleBracketExpression(t *testing.T) {
	input := "[[ -f /tmp/foo ]]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestCommandSubstitutionBacktick(t *testing.T) {
	input := "`echo hello`"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	_, ok = stmt.Expression.(*ast.CommandSubstitution)
	if !ok {
		t.Fatalf("not ast.CommandSubstitution. got=%T", stmt.Expression)
	}
}

func TestShebangStatement(t *testing.T) {
	input := "#!/bin/zsh\necho hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.Shebang)
	if !ok {
		t.Fatalf("not ast.Shebang. got=%T", program.Statements[0])
	}
}

func TestDollarSpecialVariables(t *testing.T) {
	tests := []string{
		"$#",
		"$0",
		"$*",
		"$!",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			l := lexer.New(input)
			p := New(l)
			program := p.ParseProgram()
			_ = program
			// These may produce parser errors for some special forms, but should not panic.
		})
	}
}

func TestGroupedExpressionInArithmetic(t *testing.T) {
	input := "let x = (y + z)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestArrayAccessWithLengthOperator(t *testing.T) {
	input := "${#arr}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestArrayAccessWithFlags(t *testing.T) {
	input := "${(f)foo}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestInvalidIntegerLiteral(t *testing.T) {
	input := "let x = 99999999999999999999999999999999999"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Error("expected parser error for invalid integer")
	}
}

func TestCaseStatementWithLeadingParen(t *testing.T) {
	input := `case $x in
(a) echo a;;
esac`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.CaseStatement)
	if !ok {
		t.Fatalf("not ast.CaseStatement. got=%T", program.Statements[0])
	}
	if len(stmt.Clauses) != 1 {
		t.Fatalf("expected 1 clause, got %d", len(stmt.Clauses))
	}
}

func TestCommandWithRedirections(t *testing.T) {
	tests := []string{
		"echo hello > /dev/null",
		"echo hello >> /tmp/log",
		"cat << EOF",
		"cmd 2>&1",
		"cmd <&3",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			l := lexer.New(input)
			p := New(l)
			program := p.ParseProgram()
			_ = program
		})
	}
}

func TestSimpleCommandWithDollarArgs(t *testing.T) {
	input := "echo $HOME $USER"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestSimpleCommandWithStringArgs(t *testing.T) {
	input := `echo "hello" "world"`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestSimpleCommandStartingWithSlash(t *testing.T) {
	input := "/usr/bin/echo hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestSimpleCommandStartingWithDot(t *testing.T) {
	input := ". /etc/profile"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestSimpleCommandStartingWithColon(t *testing.T) {
	input := ": noop"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestSemicolonSkipping(t *testing.T) {
	input := ";echo hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// Semicolons should be skipped; at least echo hello should parse
	found := false
	for _, s := range program.Statements {
		if s != nil {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one non-nil statement")
	}
}

func TestParseCurPrecedence(t *testing.T) {
	// Test curPrecedence with tokens that have precedence
	input := "let x = a && b"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestEqTildeOperator(t *testing.T) {
	input := "[[ $x =~ foo ]]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestNumericComparisonOperators(t *testing.T) {
	tests := []string{
		"[[ $a -eq $b ]]",
		"[[ $a -ne $b ]]",
		"[[ $a -lt $b ]]",
		"[[ $a -le $b ]]",
		"[[ $a -gt $b ]]",
		"[[ $a -ge $b ]]",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			l := lexer.New(input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
			}
		})
	}
}

func TestWhileInPipeline(t *testing.T) {
	input := "echo a | while true; do echo b; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestFunctionDefinitionInSingleCommand(t *testing.T) {
	input := "greet() { echo hi; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestProcessSubstitutionOutputDirection(t *testing.T) {
	input := "tee >(sort > sorted.txt)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestCommentSkipping(t *testing.T) {
	input := "# this is a comment\necho hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	found := false
	for _, s := range program.Statements {
		if s != nil {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one non-nil statement after comment")
	}
}

func TestDoubleParenInExpression(t *testing.T) {
	input := "let x = ((1 + 2))"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestInfixRedirection(t *testing.T) {
	input := "let x = a >> b"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestInfixRedirectionLeftShift(t *testing.T) {
	input := "let x = a << b"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestInfixRedirectionGtAmp(t *testing.T) {
	input := "let x = a >& b"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestInfixRedirectionLtAmp(t *testing.T) {
	input := "let x = a <& b"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestInvalidArrayAccess(t *testing.T) {
	input := "$arr[0]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestDollarWithSpecialTokens(t *testing.T) {
	tests := []string{
		"echo $#",
		"echo $0",
		"echo $*",
		"echo $!",
		"echo $-",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			l := lexer.New(input)
			p := New(l)
			program := p.ParseProgram()
			_ = program
		})
	}
}

func TestFunctionDefinitionKeywordWithParens(t *testing.T) {
	input := "function myfunc() { echo hello; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestCoprocStatementSimpleCommand(t *testing.T) {
	input := "coproc echo hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestCommandWithGTRedirection(t *testing.T) {
	input := "> /tmp/out echo hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestCommandWithLTRedirection(t *testing.T) {
	input := "< /tmp/in cat"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestCommandWithAmpersand(t *testing.T) {
	input := "& echo hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestDeclarationStatementMultipleAssignments(t *testing.T) {
	input := "typeset -r x=1 y=2"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.DeclarationStatement)
	if !ok {
		t.Fatalf("not ast.DeclarationStatement. got=%T", program.Statements[0])
	}
	if len(stmt.Assignments) < 2 {
		t.Errorf("expected at least 2 assignments, got %d", len(stmt.Assignments))
	}
}

func TestDeclarationStatementNoAssignment(t *testing.T) {
	input := "declare -a arr"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestForLoopArithmeticWithSemicolons(t *testing.T) {
	input := "for ((;i<10;)) do echo $i; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForLoopStatement)
	if !ok {
		t.Fatalf("not ast.ForLoopStatement. got=%T", program.Statements[0])
	}
	if stmt.Init != nil {
		t.Error("Init should be nil for empty init")
	}
}

func TestForLoopArithmeticEmptyPost(t *testing.T) {
	input := "for ((i=0;i<10;)) do echo $i; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestSingleCommandWithFunctionNotParenRparen(t *testing.T) {
	input := "echo (hello)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestEqProcessSubstitution(t *testing.T) {
	input := "diff =(sort a) =(sort b)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestBlockStatementWithSemicolons(t *testing.T) {
	input := "{ ;echo hello; ;echo world; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestIfStatementWithElse(t *testing.T) {
	input := "if [[ -f /tmp/foo ]]; then echo yes; else echo no; fi"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("not ast.IfStatement. got=%T", program.Statements[0])
	}
	if stmt.Alternative == nil {
		t.Fatal("expected Alternative to not be nil")
	}
}

func TestCommandConcatenation(t *testing.T) {
	input := `echo ${HOME}/.config`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestCommandWithVariableArgs(t *testing.T) {
	input := "echo ${var} $HOME ${arr[0]}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestExpressionStatementFallback(t *testing.T) {
	input := "true"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestExpressionCallNoBody(t *testing.T) {
	input := "myfunc();"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestMultipleRedirections(t *testing.T) {
	input := "echo hello > out.txt 2>&1"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseSimpleCommandWithTilde(t *testing.T) {
	input := "cd ~/projects"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseSimpleCommandWithGlob(t *testing.T) {
	input := "echo *.txt"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseSimpleCommandWithBang(t *testing.T) {
	input := "echo !ref"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseSimpleCommandWithBrace(t *testing.T) {
	input := "echo {a,b,c}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseIdentWithDollarLparen(t *testing.T) {
	input := "echo $(date)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseIdentWithDollarBrace(t *testing.T) {
	input := "echo ${HOME}"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseIdentFollowedByIdent(t *testing.T) {
	input := "cmd arg1 arg2"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseIdentFollowedByString(t *testing.T) {
	input := `cmd "arg1"`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseIdentFollowedByInt(t *testing.T) {
	input := "cmd 42"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseIdentFollowedByMinus(t *testing.T) {
	input := "cmd -f"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseIdentFollowedByDot(t *testing.T) {
	input := "cmd .hidden"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseIdentFollowedByVariable(t *testing.T) {
	input := "cmd $var"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestParseIdentFollowedByDollar(t *testing.T) {
	input := "cmd $#"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestSingleBracketCommand(t *testing.T) {
	input := "[ -f /tmp/foo ]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestCommandWithSlashPrefix(t *testing.T) {
	input := "/bin/echo hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestCommandWithGtgtPrefix(t *testing.T) {
	input := ">> /tmp/log echo hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestCommandWithLtltPrefix(t *testing.T) {
	input := "<< EOF"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestCommandWithGtampPrefix(t *testing.T) {
	input := ">& /dev/null echo hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestCommandWithLtampPrefix(t *testing.T) {
	input := "<& 3 echo hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestPipelineWithMultipleRedirections(t *testing.T) {
	input := "echo hello > /tmp/a >> /tmp/b < /tmp/c"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestPipelineWithNegation(t *testing.T) {
	input := "! echo hello | grep foo"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestSingleCommandFuncDefInPipeline(t *testing.T) {
	// Test the path where single command sees () and becomes function def
	input := "myfunc() { echo hello; } | grep hello"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestSingleCommandWithParensNotFunc(t *testing.T) {
	// name ( non-rparen -- triggers the else branch in parseSingleCommand
	input := "cmd (arg1 arg2)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestParserErrorOnMissing(t *testing.T) {
	// Test noPrefixParseFnError
	input := "let x = "
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Error("expected parser errors")
	}
}

func TestParserPeekError(t *testing.T) {
	// Test peekError path - missing expected token
	input := "let = 5"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Error("expected parser errors")
	}
}

func TestDollarParenDoubleRparen(t *testing.T) {
	// Test $(( expr )) with DoubleRparen token
	input := "$((x + y))"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestDollarParenWithDoubleParenSeparated(t *testing.T) {
	// Test $((expr)) where )) is two separate ) tokens
	input := "$(( 1 + 2 ))"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestForLoopArithmeticOptionalSemicolonBeforeDo(t *testing.T) {
	input := "for ((i=0;i<10;i++)); do echo $i; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestDeclarationStatementWithFlagIdent(t *testing.T) {
	// Test ident starting with - being treated as flag
	input := "declare -x PATH"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestDeclarationStatementUnexpectedToken(t *testing.T) {
	// Test unexpected token in declaration
	input := "declare 42"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
	// Should not panic
}

func TestDeclarationValue_ArrayLiteral(t *testing.T) {
	// Test parseDeclarationValue with array
	input := "declare arr=(one two three)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}
}

func TestWhileLoopMissingDo(t *testing.T) {
	input := "while true; echo foo; done"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
	// Should not panic even with missing do
}

func TestCaseStatementEmpty(t *testing.T) {
	input := "case $x in\nesac"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestSubshellMissingRparen(t *testing.T) {
	input := "(echo hello"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
	// Should not panic
}

func TestDoubleBracketMissingClose(t *testing.T) {
	input := "[[ -f foo"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
	// Should not panic, produces error
}

func TestLBracketSpaceAfterIdentBreaksLoop(t *testing.T) {
	// Test the path where LBRACKET with preceding space stops infix parsing
	input := "arr [0]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestProcessSubstitutionInExpression(t *testing.T) {
	input := "<(sort file.txt)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestReturnStatementNoValue(t *testing.T) {
	input := "return;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	_ = program
}

func TestParseSingleCommand_FuncDefDetailed(t *testing.T) {
	input := "greet() { echo hi; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}

	if stmt, ok := program.Statements[0].(*ast.ExpressionStatement); ok {
		if funcDef, ok := stmt.Expression.(*ast.FunctionDefinition); ok {
			if funcDef.Name.Value != "greet" {
				t.Errorf("expected name 'greet', got %q", funcDef.Name.Value)
			}
		}
	}
}

func TestCommandPipelineMultiRedir(t *testing.T) {
	tests := []string{
		"echo hello > out.txt",
		"echo hello >> out.txt",
		"cat < input.txt",
		"cmd << EOF",
		"cmd 2>&1",
		"cmd <&3",
		"echo a > out1 >> out2",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			l := lexer.New(input)
			p := New(l)
			_ = p.ParseProgram()
		})
	}
}

func TestCaseMultiplePatterns(t *testing.T) {
	input := "case $x in\na|b|c) echo matched;;\n*) echo default;;\nesac"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.CaseStatement)
	if !ok {
		t.Fatalf("not ast.CaseStatement. got=%T", program.Statements[0])
	}
	if len(stmt.Clauses) != 2 {
		t.Fatalf("expected 2 clauses, got %d", len(stmt.Clauses))
	}
	if len(stmt.Clauses[0].Patterns) != 3 {
		t.Errorf("expected 3 patterns in first clause, got %d", len(stmt.Clauses[0].Patterns))
	}
}

func TestForLoopEmptyInitAndPost(t *testing.T) {
	// Empty init and empty post (condition present)
	input := "for ((;i<10;)) do echo loop; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ForLoopStatement)
	if !ok {
		t.Fatalf("not ast.ForLoopStatement. got=%T", program.Statements[0])
	}
	if stmt.Init != nil {
		t.Error("Init should be nil")
	}
	if stmt.Condition == nil {
		t.Error("Condition should not be nil")
	}
	if stmt.Post != nil {
		t.Error("Post should be nil")
	}
}

func TestIfNoFi(t *testing.T) {
	input := "if true; then echo yes"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestIfMissingThen(t *testing.T) {
	input := "if true; echo yes; fi"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestSingleCmdMultiArgs(t *testing.T) {
	input := "ls -la /tmp /var /home"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestCmdWordConcatenation(t *testing.T) {
	input := "echo ${HOME}/.config/file"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	cmd := stmt.Expression.(*ast.SimpleCommand)
	if len(cmd.Arguments) < 1 {
		t.Fatal("expected at least 1 argument")
	}
}

func TestExprOrFuncDef_NotCall(t *testing.T) {
	input := "true"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestExprOrFuncDef_CallWithArgs(t *testing.T) {
	input := "myfunc(1, 2);"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestFuncLiteralNamedWithParens(t *testing.T) {
	input := "function myfunc() { echo hello; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	fl, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("not FunctionLiteral. got=%T", stmt.Expression)
	}
	if fl.Name.Value != "myfunc" {
		t.Errorf("expected name 'myfunc', got %q", fl.Name.Value)
	}
}

func TestFuncLiteralNoParens(t *testing.T) {
	input := "function myfunc { echo hello; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	fl, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("not FunctionLiteral. got=%T", stmt.Expression)
	}
	if len(fl.Params) != 0 {
		t.Errorf("expected 0 params, got %d", len(fl.Params))
	}
}

func TestDollarParenSeparateRparens(t *testing.T) {
	input := "$(( 1 + 2 ))"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestDeclLocal(t *testing.T) {
	input := "local x y z"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestCmdWordLiterals(t *testing.T) {
	tests := []string{
		"echo *",
		"echo ?",
		"echo ~user",
		"echo .hidden",
		"echo ,list",
		"echo :path",
	}
	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			l := lexer.New(input)
			p := New(l)
			_ = p.ParseProgram()
		})
	}
}

func TestSelectSemiBeforeDo(t *testing.T) {
	input := "select opt in a b c; do echo $opt; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.SelectStatement)
	if !ok {
		t.Fatalf("not ast.SelectStatement. got=%T", program.Statements[0])
	}
	if len(stmt.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(stmt.Items))
	}
}

func TestCoprocNamedBrace(t *testing.T) {
	input := "coproc worker { echo hello; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.CoprocStatement)
	if !ok {
		t.Fatalf("not ast.CoprocStatement. got=%T", program.Statements[0])
	}
	if stmt.Name != "worker" {
		t.Errorf("expected name 'worker', got %q", stmt.Name)
	}
}

func TestArithCmdMissingClose(t *testing.T) {
	input := "((x + 1"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestCaseMissingIn(t *testing.T) {
	input := "case $x a) echo a;; esac"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestForLoopMissingDoKw(t *testing.T) {
	input := "for i in a b c; echo $i; done"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
	if len(p.Errors()) == 0 {
		t.Error("expected parser errors for missing 'do'")
	}
}

func TestSelectMissingDoKw(t *testing.T) {
	input := "select i in a b c; echo $i; done"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
	if len(p.Errors()) == 0 {
		t.Error("expected parser errors for missing 'do'")
	}
}

func TestGroupedExprArithMissing(t *testing.T) {
	input := "let x = (a + b"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestFuncParamsWithComma(t *testing.T) {
	input := "function(a, b, c) { echo; }"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	fl := stmt.Expression.(*ast.FunctionLiteral)
	if len(fl.Params) != 3 {
		t.Errorf("expected 3 params, got %d", len(fl.Params))
	}
}

func TestProcessSubOutput(t *testing.T) {
	input := "tee >(sort > sorted.txt)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestEqProcessSub(t *testing.T) {
	input := "diff =(sort a) =(sort b)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) < 1 {
		t.Fatalf("expected at least 1 statement, got=%d", len(program.Statements))
	}
}

func TestBacktickMissingClose(t *testing.T) {
	input := "`echo hello"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestIndexExprMissingClose(t *testing.T) {
	input := "arr[0"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestArrayAccessMissingRbrace(t *testing.T) {
	input := "${arr[0]"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestProcessSubMissingClose(t *testing.T) {
	input := "<(sort a"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestFuncLitMissingLbrace(t *testing.T) {
	input := "function myfunc echo hello"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestDollarParenMissingClose(t *testing.T) {
	input := "$(echo hello"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestCaseMissingRparen(t *testing.T) {
	input := "case $x in\na echo a;;\nesac"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestWhileInPipeRHS(t *testing.T) {
	input := "echo a | while true; do read line; done"
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()
}

func TestForLoopArithNoInit(t *testing.T) {
	input := "for ((;i<10;i++)) do echo $i; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ForLoopStatement)
	if stmt.Init != nil {
		t.Error("Init should be nil")
	}
	if stmt.Condition == nil {
		t.Error("Condition should not be nil")
	}
	if stmt.Post == nil {
		t.Error("Post should not be nil")
	}
}

func TestForLoopArithNoPost(t *testing.T) {
	input := "for ((i=0;i<10;)) do echo $i; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ForLoopStatement)
	if stmt.Init == nil {
		t.Error("Init should not be nil")
	}
	if stmt.Post != nil {
		t.Error("Post should be nil")
	}
}

func TestForLoopArithWithSemicolonBeforeDo(t *testing.T) {
	// Tests the optional semicolon before DO path
	input := "for ((i=0;i<10;i++)); do echo $i; done"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ForLoopStatement)
	if stmt.Init == nil {
		t.Error("Init should not be nil")
	}
}

func TestDeclarationMultiFlags(t *testing.T) {
	input := "typeset -A -g mymap"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.DeclarationStatement)
	if len(stmt.Flags) < 2 {
		t.Errorf("expected at least 2 flags, got %d", len(stmt.Flags))
	}
}

func TestDeclarationPlusEq(t *testing.T) {
	input := "declare x+=5"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.DeclarationStatement)
	if len(stmt.Assignments) < 1 {
		t.Fatalf("expected at least 1 assignment, got %d", len(stmt.Assignments))
	}
	if !stmt.Assignments[0].IsAppend {
		t.Error("expected IsAppend to be true")
	}
}

func TestCaseStatementSemicolonsInBody(t *testing.T) {
	input := "case $x in\n  a) echo one; echo two;;\nesac"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.CaseStatement)
	if len(stmt.Clauses) != 1 {
		t.Fatalf("expected 1 clause, got %d", len(stmt.Clauses))
	}
}

// TestParser_DDCommandWithIfOfKeywords is the regression test for
// https://github.com/afadesigns/zshellcheck/issues/435. `dd if=foo of=bar`
// contains identifiers that happen to match the `if`/`of` keyword prefix
// — the lexer must demote them to IDENT when `=` follows so the parser
// does not try to open an if-statement.
func TestParser_DDCommandWithIfOfKeywords(t *testing.T) {
	inputs := []string{
		`dd if=src of=dst`,
		`dd if=/tmp/x of=/tmp/y bs=4M`,
		`some_tool if=arg`,
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			l := lexer.New(input)
			p := New(l)
			p.ParseProgram()
			if errs := p.Errors(); len(errs) > 0 {
				t.Fatalf("unexpected parser errors on %q: %v", input, errs)
			}
		})
	}
}

// TestParser_ForLoopAndOr is the regression test for
// https://github.com/afadesigns/zshellcheck/issues/347. `for ... in ...;
// do ...; done` and `||` previously raised parser errors; both must
// parse cleanly now.
func TestParser_ForLoopAndOr(t *testing.T) {
	inputs := []string{
		`for x in a b c; do echo $x; done`,
		`true || echo fail`,
		`cmd1 || cmd2 || cmd3`,
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			l := lexer.New(input)
			p := New(l)
			p.ParseProgram()
			if errs := p.Errors(); len(errs) > 0 {
				t.Fatalf("unexpected parser errors on %q: %v", input, errs)
			}
		})
	}
}

// TestParser_IfElifElseFi is the regression test for the elif branch of
// https://github.com/afadesigns/zshellcheck/issues/126.
func TestParser_IfElifElseFi(t *testing.T) {
	input := "if [[ $a == 1 ]]; then\n  echo one\nelif [[ $a == 2 ]]; then\n  echo two\nelse\n  echo other\nfi\n"
	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		t.Fatalf("unexpected parser errors: %v", errs)
	}
}
