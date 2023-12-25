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
		},
	}

	parser := New(lexer.New(input))
	actual := parser.ParseProgram()

	assert.Equal(t, expected, actual)
}
