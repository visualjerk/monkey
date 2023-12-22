package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	char         byte // current char under examination
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readNextChar()
	return lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var tokenType token.TokenType
	tokenLiteral := string(lexer.char)

	switch lexer.char {
	case '=':
		tokenType = token.ASSIGN
	case ';':
		tokenType = token.SEMICOLON
	case '(':
		tokenType = token.LPAREN
	case ')':
		tokenType = token.RPAREN
	case '{':
		tokenType = token.LBRACE
	case '}':
		tokenType = token.RBRACE
	case ',':
		tokenType = token.COMMA
	case '+':
		tokenType = token.PLUS
	case 0:
		tokenType = token.EOF
		tokenLiteral = ""
	}

	lexer.readNextChar()

	return token.Token{
		Type:    tokenType,
		Literal: tokenLiteral,
	}
}

func (lexer *Lexer) readNextChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.char = 0
	} else {
		lexer.char = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}
