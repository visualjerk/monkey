package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProgram(t *testing.T) {
	input := `
	let x = 5;
	let y = 10 + 5;

	return 5;
`

	expected := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &ast.Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "x",
					},
					Value: "x",
				},
				Value: &ast.Int{
					Token: token.Token{
						Type:    token.INT,
						Literal: "5",
					},
					Value: "5",
				},
			},
			&ast.LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &ast.Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "y",
					},
					Value: "y",
				},
				Value: &ast.AddExpression{
					Token: token.Token{
						Type:    token.PLUS,
						Literal: "+",
					},
					Left: &ast.Int{
						Token: token.Token{
							Type:    token.INT,
							Literal: "10",
						},
						Value: "10",
					},
					Right: &ast.Int{
						Token: token.Token{
							Type:    token.INT,
							Literal: "5",
						},
						Value: "5",
					},
				},
			},
			&ast.ReturnStatement{
				Token: token.Token{
					Type:    token.RETURN,
					Literal: "return",
				},
				Value: &ast.Int{
					Token: token.Token{
						Type:    token.INT,
						Literal: "5",
					},
					Value: "5",
				},
			},
		},
	}

	parser := New(lexer.New(input))
	actual := parser.ParseProgram()

	assert.Equal(t, expected, actual)
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	expected := &ast.Program{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "foobar",
				},
				Value: &ast.Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "foobar",
					},
					Value: "foobar",
				},
			},
		},
	}

	parser := New(lexer.New(input))
	actual := parser.ParseProgram()

	assert.Equal(t, expected, actual)
}

func TestParserErrors(t *testing.T) {
	input := `
	let x 5;
	let = 10 + 5;
	let 10;
`

	expected := []string{
		"expected next token to be =, got INT instead",
		"expected next token to be IDENT, got = instead",
		"expected next token to be IDENT, got INT instead",
		"expected next token to be =, got INT instead",
	}

	parser := New(lexer.New(input))
	parser.ParseProgram()
	actual := parser.GetErrors()

	assert.Equal(t, expected, actual)
}
