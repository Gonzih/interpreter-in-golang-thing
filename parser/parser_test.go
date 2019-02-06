package parser

import (
	"testing"

	"github.com/Gonzih/go-interpreter/ast"
	"github.com/Gonzih/go-interpreter/lexer"
	"github.com/stretchr/testify/assert"
)

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

func testLetStatement(t *testing.T, name string, s ast.Statement) {
	assert.Equal(t, "let", s.TokenLiteral())
	assert.IsType(t, &ast.LetStatement{}, s)

	letStmt := s.(*ast.LetStatement)

	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
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

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	assert.Len(t, errors, 0)

	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}

	if len(errors) > 0 {
		t.FailNow()
	}
}
