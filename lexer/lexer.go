package lexer

import (
	"github.com/Gonzih/go-interpreter/token"
)

type Lexer struct {
	input string
	// current char pos
	position int
	// next char pos
	readPossition int
	// current char
	ch byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPossition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPossition]
	}
	l.position = l.readPossition
	l.readPossition++
}

var tokenTable = map[byte]token.TokenType{
	'=': token.ASSIGN,
	';': token.SEMICOLON,
	'(': token.LPAREN,
	')': token.RPAREN,
	'{': token.LBRACE,
	'}': token.RBRACE,
	'+': token.PLUS,
	'-': token.MINUS,
	'!': token.BANG,
	'/': token.SLASH,
	'*': token.ASTERISK,
	'<': token.LT,
	'>': token.GT,
	',': token.COMMA,
	0:   token.EOF,
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaces()

	// handling != and ==
	if (l.ch == '=' || l.ch == '!') && l.peekChar() == '=' {
		ch := l.ch
		l.readChar()
		tok.Literal = string(ch) + string(l.ch)
		tok.Type = token.LookupIdent(tok.Literal)
	} else if tt, ok := tokenTable[l.ch]; ok {
		tok = newToken(tt, l.ch)
	} else if isLetter(l.ch) {
		tok.Literal = l.readIdentifier()
		tok.Type = token.LookupIdent(tok.Literal)
		return tok
	} else if isDigit(l.ch) {
		tok.Type = token.INT
		tok.Literal = l.readNumber()
		return tok
	} else {
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPossition >= len(l.input) {
		return 0
	}

	return l.input[l.readPossition]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '?'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	tok := token.Token{Type: tokenType, Literal: string(ch)}
	if tokenType == token.EOF {
		tok.Literal = ""
	}

	return tok
}
