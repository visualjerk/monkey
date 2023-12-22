package lexer

import (
	"monkey/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	results := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for _, result := range results {
		actualToken := lexer.NextToken()
		expectedToken := token.Token{
			Type:    result.expectedType,
			Literal: result.expectedLiteral,
		}

		assert.Equal(t, actualToken, expectedToken, "Tokens should be equal.")
	}
}
