package parser

import (
	"fmt"
	"github.com/hendrikbursian/monkey-programming-language/ast"
	"github.com/hendrikbursian/monkey-programming-language/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {

	tests := []struct {
		code               string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let fooabar = y;", "fooabar", "y"},
	}

	for _, test := range tests {
		lexer := lexer.New(test.code)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements got=%d", len(program.Statements))
		}

		statement := program.Statements[0]
		if !testLetStatement(t, statement, test.expectedIdentifier) {
			return
		}

		value := statement.(*ast.LetStatement).Value
		if !testLiteralExpression(t, value, test.expectedValue) {
			return
		}

	}
}

func testLetStatement(t *testing.T, statement ast.Statement, identifier string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("statement.TokenLiteral should be 'let'. got=%q", statement.TokenLiteral())
		return false
	}

	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement is not *ast.LetStatement. got=%T", statement)
		return false
	}

	if letStatement.Identifier.Value != identifier {
		t.Errorf("letStatement.Identifier.Value is not '%s'. got=%s", identifier, letStatement.Identifier.Value)
		return false
	}

	if letStatement.Identifier.TokenLiteral() != identifier {
		t.Errorf("letStatement.Identifier.TokenLiteral is not '%s'. got=%s", identifier, letStatement.Identifier.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, error := range errors {
		t.Errorf("parser error: %q", error)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		code  string
		value interface{}
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return x;", "x"},
	}

	for _, test := range tests {
		lexer := lexer.New(test.code)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		for _, statement := range program.Statements {
			returnStatement, ok := statement.(*ast.ReturnStatement)
			if !ok {
				t.Errorf("Statement is not *ast.ReturnStatement. got=%T", returnStatement)
				continue
			}

			if returnStatement.TokenLiteral() != "return" {
				t.Errorf("returnStatement.TokenLiteral is not 'return'. got=%s", returnStatement.TokenLiteral())
			}

			if !testLiteralExpression(t, returnStatement.ReturnValue, test.value) {
				return
			}
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	code := `foobar;`

	lexer := lexer.New(code)
	parser := New(lexer)
	program := parser.ParseProgram()

	if len(program.Statements) != 1 {
		t.Errorf("Length of program.Statements is not 1. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statement is not an ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	identifier, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("Statement is not an ast.Identifier. got=%T", statement.Expression)
	}

	if identifier.Value != "foobar" {
		t.Errorf("identifier.Value is not 'foobar' got=%s", identifier.Value)
	}

	if identifier.TokenLiteral() != "foobar" {
		t.Errorf("identifier.TokenLiteral is not 'foobar' got=%s", identifier.TokenLiteral())
	}
}

func TestIntegerExpression(t *testing.T) {
	code := `5;`

	lexer := lexer.New(code)
	parser := New(lexer)
	program := parser.ParseProgram()

	if len(program.Statements) != 1 {
		t.Errorf("Length of program.Statements is not 1. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statement is not an ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	integerLiteral, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Statement is not an ast.IntegerLiteral. got=%T", statement.Expression)
	}

	if integerLiteral.Value != 5 {
		t.Errorf("integerLiteral.Value is not 'foobar' got=%d", integerLiteral.Value)
	}

	if integerLiteral.TokenLiteral() != "5" {
		t.Errorf("integerLiteral.TokenLiteral is not 'foobar' got=%s", integerLiteral.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, test := range prefixTests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Length of program.Statements is not 1. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is not an ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Statement is not an ast.PrefixExpression. got=%T", statement.Expression)
		}
		if expression.Operator != test.operator {
			t.Fatalf("expression.Operator is not '%s' got=%s", test.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, test.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) bool {
	integerLiteral, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expression is not *ast.IntegerLiteral. got=%T", expression)
		return false
	}

	if integerLiteral.Value != value {
		t.Errorf("integerLiteral.Value is not %d. got=%d", value, integerLiteral.Value)
		return false
	}

	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integerLiteral.TokenLiteral is not %q. got=%s", value, integerLiteral.TokenLiteral())
		return false
	}

	return true
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
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, test := range infixTests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Length of program.Statements is not 1. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is not an ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Statement is not an ast.InfixExpression. got=%T", statement.Expression)
		}

		if !testLiteralExpression(t, expression.Left, test.leftValue) {
			return
		}

		if expression.Operator != test.operator {
			t.Fatalf("expression.Operator is not '%s' got=%s", test.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, test.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3>5==false",
			"((3 > 5) == false)",
		},
		{
			"3<5==false",
			"((3 < 5) == false)",
		},
		{
			"1+(2+3)+4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2/(5+5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5+5)",
			"(-(5 + 5))",
		},
		{
			"!(true==true)",
			"(!(true == true))",
		},
		{
			"if a<a - (b+ c) { hello } else {test}",
			"if (a < (a - (b + c))) { hello } else { test }",
		},
		{
			"a + add (b*c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a,b,1,2*3,4+5,add(6,7*8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a+b+c*d/f+g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a *[1,2,3,4][b*c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		actual := program.String()
		if actual != test.expected {
			t.Errorf("\nexpected:\n%q, got=\n%q", test.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := statement.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("expression not *ast.Boolean. got=%T", statement.Expression)
		}
		if boolean.Value != test.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", test.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestIfStatement(t *testing.T) {
	code := `if a<b { hello }`

	lexer := lexer.New(code)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not of type %T. got=%T", &ast.ExpressionStatement{}, program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement.Expression is not of type %T. got=%T", &ast.IfExpression{}, statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "a", "<", "b") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Fatalf("expression.Consequence.Statements does not contain 1 statements. got=%d", len(expression.Consequence.Statements))
	}
	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expression.Consequence.Statements[0] is not of type ast.ExpressionStatement. got=%T", expression.Consequence.Statements[0])
	}
	if !testLiteralExpression(t, consequence.Expression, "hello") {
		return
	}

}

func TestIfElseStatement(t *testing.T) {
	code := `if a<b { hello } else { test }`

	lexer := lexer.New(code)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not of type %T. got=%T", &ast.ExpressionStatement{}, program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("statement.Expression is not of type %T. got=%T", &ast.IfExpression{}, statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "a", "<", "b") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Fatalf("expression.Consequence.Statements does not contain 1 statements. got=%d", len(expression.Consequence.Statements))
	}
	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expression.Consequence.Statements[0] is not of type ast.ExpressionStatement. got=%T", expression.Consequence.Statements[0])
	}
	if !testLiteralExpression(t, consequence.Expression, "hello") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Fatalf("expression.Alternative.Statements does not contain 1 statements. got=%d", len(expression.Alternative.Statements))
	}
	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expression.Alternative.Statements[0] is not of type ast.ExpressionStatement. got=%T", expression.Alternative.Statements[0])
	}
	if !testLiteralExpression(t, alternative.Expression, "test") {
		return
	}

}

func TestFunctionLiteralParsing(t *testing.T) {
	code := `fn(x, y) { x + y; }`

	lexer := lexer.New(code)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not of type %T. got=%T", &ast.ExpressionStatement{}, program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement.Expression is not of type %T. got=%T", &ast.FunctionLiteral{}, statement.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function.Parameters does not contain 2 statements. got=%d", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements does not contain 1 statements. got=%d", len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("bodyStatement is not of type %T. got=%T", &ast.ExpressionStatement{}, bodyStatement)
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		code           string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(a) {};", []string{"a"}},
		{"fn(a,b) {};", []string{"a", "b"}},
		{"fn(a,b,c) {};", []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		lexer := lexer.New(test.code)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(test.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d", len(test.expectedParams), len(function.Parameters))
		}

		for i, identifier := range test.expectedParams {
			testLiteralExpression(t, function.Parameters[i], identifier)
		}
	}

}

func TestCallExpressionParsing(t *testing.T) {
	code := "add(1, 1*2, 4+5);"

	lexer := lexer.New(code)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statement length wrong. want %d, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not of type %T. got=%T", &ast.ExpressionStatement{}, program.Statements[0])
	}

	call, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("statement.Exression not of type %T. got=%T (%+v)", &ast.CallExpression{}, statement.Expression, statement.Expression)
	}

	if !testIdentifier(t, call.Function, "add") {
		return
	}

	if len(call.Arguments) != 3 {
		t.Fatalf("length parameters is not %d. got=%d", 3, len(call.Arguments))
	}

	testLiteralExpression(t, call.Arguments[0], 1)
	testInfixExpression(t, call.Arguments[1], 1, "*", 2)
	testInfixExpression(t, call.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestStringLiteral(t *testing.T) {
	code := `"hello world"`
	parser, program := testParse(code)
	checkParserErrors(t, parser)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expressions is not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value is not %q. got=%q", "hello world", literal.Value)
	}
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression not *ast.Identifier. got=%T", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value not %s. got=%s", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral not %s. got=%s", value, identifier.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch value := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(value))
	case int64:
		return testIntegerLiteral(t, expression, value)
	case bool:
		return testBooleanLiteral(t, expression, value)
	case string:
		return testIdentifier(t, expression, value)
	}

	t.Errorf("type of expresion not handled. got=%T", expression)
	return false
}

func testInfixExpression(t *testing.T, expression ast.Expression, left interface{}, operator string, right interface{}) bool {
	operatorExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression is not ast.OperatorExpression. got=%T(%s)", expression, expression)
		return false
	}

	if !testLiteralExpression(t, operatorExpression.Left, left) {
		return false
	}

	if operatorExpression.Operator != operator {
		t.Errorf("operatorExpression.Operator is not '%s'. got=%s", operator, operatorExpression.Operator)
		return false
	}

	if !testLiteralExpression(t, operatorExpression.Right, right) {
		return false
	}

	return true

}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	boolean, ok := expression.(*ast.Boolean)

	if !ok {
		t.Errorf("expression not *ast.Boolean. got=%T", expression)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral not %t. got=%s", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func testParse(code string) (*Parser, *ast.Program) {
	lexer := lexer.New(code)
	parser := New(lexer)
	program := parser.ParseProgram()

	return parser, program
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1,2*2,3+3]"

	p, program := testParse(input)
	checkParserErrors(t, p)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	arr, ok := statement.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expresison is not an ArrayLiteral. got=%T", statement.Expression)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("array has wrong length. expected=%d, got=%d", 3, len(arr.Elements))
	}

	testIntegerLiteral(t, arr.Elements[0], 1)
	testInfixExpression(t, arr.Elements[1], 2, "*", 2)
	testInfixExpression(t, arr.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := `myarray[1+1]`

	p, program := testParse(input)
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	idx, ok := stmt.Expression.(*ast.IndexExpression)

	if !ok {
		t.Fatalf("stmt.expresion is not of type ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, idx.Left, "myarray") {
		return
	}

	if !testInfixExpression(t, idx.Index, 1, "+", 1) {
		return
	}
}

func TestHashLiteral(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	p, program := testParse(input)
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Errorf("stmt.Expresison is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("length of hash.Pairs is incorrect. got=%d, want=%d", len(hash.Pairs), 3)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not a stringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.Value]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 2/2, "two": 1+1, "three": 9/3}`

	p, program := testParse(input)
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Errorf("stmt.Expresison is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("length of hash.Pairs is incorrect. got=%d, want=%d", len(hash.Pairs), 3)
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 2, "/", 2)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 1, "+", 1)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 9, "/", 3)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.Value]
		if !ok {
			t.Errorf("key %s does not exist", literal.String())
			continue
		}

		testFunc(value)
	}
}

func TestEmptyHashLiteral(t *testing.T) {
	input := `{}`

	p, program := testParse(input)
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Errorf("stmt.Expresison is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("length of hash.Pairs is incorrect. got=%d, want=%d", len(hash.Pairs), 0)
	}
}
