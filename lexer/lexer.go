package lexer

import (
	"monkey/token"
	"slices"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	char         byte // current char under examination
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

func (lexer *Lexer) NextToken() token.Token {
	lexer.skipWhitespace()

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
	default:
		if isLetter(lexer.char) {
			tokenLiteral = lexer.readIdentifier()
			tokenType = token.LookupIdentifier(tokenLiteral)

			return token.Token{
				Type:    tokenType,
				Literal: tokenLiteral,
			}
		}
		if isDigit(lexer.char) {
			tokenLiteral = lexer.readNumber()
			tokenType = token.INT

			return token.Token{
				Type:    tokenType,
				Literal: tokenLiteral,
			}
		}

		tokenType = token.ILLEGAL
		tokenLiteral = ""
	}

	lexer.readChar()

	return token.Token{
		Type:    tokenType,
		Literal: tokenLiteral,
	}
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.char = 0
	} else {
		lexer.char = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lexer *Lexer) readIdentifier() string {
	position := lexer.position

	for isLetter(lexer.char) {
		lexer.readChar()
	}

	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position

	for isDigit(lexer.char) {
		lexer.readChar()
	}

	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) skipWhitespace() {
	for isWhitespace(lexer.char) {
		lexer.readChar()
	}
}

func isWhitespace(char byte) bool {
	return slices.Contains([]byte{' ', '\t', '\n', '\r'}, char)
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
