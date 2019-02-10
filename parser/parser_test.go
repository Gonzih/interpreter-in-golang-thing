package parser

import (
	"fmt"
	"testing"

	"github.com/Gonzih/go-interpreter/ast"
	"github.com/Gonzih/go-interpreter/lexer"
	"github.com/stretchr/testify/assert"
)

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}

	assert.Len(t, errors, 0)

	if len(errors) > 0 {
		t.FailNow()
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) {
	assert.NotNil(t, exp)
	if exp == nil {
		t.FailNow()
	}

	ident, ok := exp.(*ast.Identifier)
	assert.NotNil(t, ident)
	assert.True(t, ok)
	if ident == nil {
		t.FailNow()
	}

	assert.Equal(t, value, ident.Value)
	assert.Equal(t, value, ident.TokenLiteral())
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(v))
	case int64:
		testIntegerLiteral(t, exp, v)
	case string:
		testIdentifier(t, exp, v)
	case bool:
		testBoolean(t, exp, v)
	default:
		t.Errorf("Canot handle type %T", exp)
		t.FailNow()
	}
}

func testBoolean(t *testing.T, inExp ast.Expression, value bool) {
	assert.NotNil(t, inExp)
	if inExp == nil {
		t.FailNow()
	}

	exp, ok := inExp.(*ast.Boolean)
	assert.NotNil(t, exp)
	assert.True(t, ok)
	if exp == nil {
		t.FailNow()
	}

	assert.Equal(t, value, exp.Value)
	assert.Equal(t, fmt.Sprintf("%v", value), exp.TokenLiteral())
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) {
	assert.NotNil(t, il)
	if il == nil {
		t.FailNow()
	}

	integ, ok := il.(*ast.IntegerLiteral)
	assert.NotNil(t, integ)
	assert.True(t, ok)
	if integ == nil {
		t.FailNow()
	}

	assert.Equal(t, value, integ.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), integ.TokenLiteral())
}

func testLetStatement(t *testing.T, name string, s ast.Statement) {
	assert.Equal(t, "let", s.TokenLiteral())
	assert.IsType(t, &ast.LetStatement{}, s)

	letStmt := s.(*ast.LetStatement)

	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

func testInfixExpression(t *testing.T, inExp ast.Expression, left interface{}, operator string, right interface{}) {
	exp, ok := inExp.(*ast.InfixExpression)
	assert.NotNil(t, exp)
	assert.True(t, ok)
	if exp == nil {
		t.FailNow()
	}

	testLiteralExpression(t, exp.Left, left)
	assert.Equal(t, exp.Operator, operator)
	testLiteralExpression(t, exp.Right, right)
}

func TestLetStatements(t *testing.T) {
	input := ` let x = 5;
	let y = 10;
	let foobar = 838383;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.Len(t, program.Statements, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(t, tt.expectedIdentifier, stmt)
	}
}

func TestReturnStatements(t *testing.T) {
	input := `return 5;
	return 10;
	return 1911;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.Len(t, program.Statements, 3)

	for _, stmt := range program.Statements {
		assert.IsType(t, &ast.ReturnStatement{}, stmt)
		returnStmt := stmt.(*ast.ReturnStatement)
		assert.Equal(t, "return", returnStmt.TokenLiteral())
	}
}

func TestIdentifierExpressions(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.Len(t, program.Statements, 1)

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	ident, ok := expStmt.Expression.(*ast.Identifier)
	assert.NotNil(t, ident)
	assert.True(t, ok)
	if ident == nil {
		t.FailNow()
	}

	testLiteralExpression(t, ident, "foobar")
}

func TestIntegerExpressions(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.Len(t, program.Statements, 1)

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	ident, ok := expStmt.Expression.(*ast.IntegerLiteral)
	assert.NotNil(t, ident)
	assert.True(t, ok)
	if ident == nil {
		t.FailNow()
	}

	testLiteralExpression(t, ident, 5)
}

func TestParsePrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		assert.Len(t, program.Statements, 1)

		expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		exp, ok := expStmt.Expression.(*ast.PrefixExpression)
		assert.NotNil(t, exp)
		assert.True(t, ok)
		if exp == nil {
			t.FailNow()
		}

		assert.Equal(t, tt.operator, exp.Operator)
		testLiteralExpression(t, exp.Right, tt.value)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 6;", 5, "+", 6},
		{"5 - 6;", 5, "-", 6},
		{"5 * 6;", 5, "*", 6},
		{"5 / 6;", 5, "/", 6},
		{"5 > 6;", 5, ">", 6},
		{"5 < 6;", 5, "<", 6},
		{"5 == 6;", 5, "==", 6},
		{"5 != 6;", 5, "!=", 6},
		{"true == true;", true, "==", true},
		{"true != false", true, "!=", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		assert.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a + b", "(a + b)"},
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
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		assert.Equal(t, tt.expected, program.String())
	}
}

func TestBoolExpressions(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"true", true},
		{"true;", true},
		{"false", false},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		assert.Len(t, program.Statements, 1)

		expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		testLiteralExpression(t, expStmt.Expression, tt.value)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.NotNil(t, stmt)
	assert.True(t, ok)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	assert.NotNil(t, exp)
	assert.True(t, ok)

	testInfixExpression(t, exp.Condition, "x", "<", "y")
	assert.Len(t, exp.Consequence.Statements, 1)

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.NotNil(t, consequence)
	assert.True(t, ok)

	testIdentifier(t, consequence.Expression, "x")

	assert.Nil(t, exp.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.NotNil(t, stmt)
	assert.True(t, ok)

	exp, ok := stmt.Expression.(*ast.IfExpression)
	assert.NotNil(t, exp)
	assert.True(t, ok)

	testInfixExpression(t, exp.Condition, "x", "<", "y")
	assert.Len(t, exp.Consequence.Statements, 1)

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.NotNil(t, consequence)
	assert.True(t, ok)
	testIdentifier(t, consequence.Expression, "x")

	assert.NotNil(t, exp.Alternative)
	if exp.Alternative == nil {
		t.FailNow()
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	assert.NotNil(t, alternative)
	assert.True(t, ok)
	testIdentifier(t, alternative.Expression, "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.NotNil(t, stmt)
	assert.True(t, ok)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	assert.NotNil(t, function)
	assert.True(t, ok)

	assert.Len(t, function.Parameters, 2)
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	body, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	assert.NotNil(t, body)
	assert.True(t, ok)

	testInfixExpression(t, body.Expression, "x", "+", "y")
}

func TestFunctionParametersParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		assert.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		assert.NotNil(t, stmt)
		assert.True(t, ok)

		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		assert.NotNil(t, function)
		assert.True(t, ok)

		assert.Len(t, function.Parameters, len(tt.expectedParams))

		for i, param := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], param)
		}
	}
}
