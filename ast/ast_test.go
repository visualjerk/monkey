package ast

import (
	"monkey/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	testCases := []struct {
		input    Statement
		expected string
	}{
		{
			input: &LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "foo",
					},
					Value: "foo",
				},
				Value: &IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "5",
					},
					Value: 5,
				},
			},
			expected: "let foo = 5;",
		},
		{
			input: &ExpressionStatement{
				Token: token.Token{
					Type:    token.BANG,
					Literal: "!",
				},
				Value: &PrefixExpression{
					Token: token.Token{
						Type:    token.BANG,
						Literal: "!",
					},
					Operator: "+",
					Right: &IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "5",
						},
						Value: 5,
					},
				},
			},
			expected: "(+5)",
		},
		{
			input: &ExpressionStatement{
				Token: token.Token{
					Type:    token.INT,
					Literal: "5",
				},
				Value: &InfixExpression{
					Token: token.Token{
						Type:    token.PLUS,
						Literal: "+",
					},
					Operator: "+",
					Left: &IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "5",
						},
						Value: 5,
					},
					Right: &IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "5",
						},
						Value: 5,
					},
				},
			},
			expected: "(5 + 5)",
		},
	}

	for _, testCase := range testCases {
		program := &Program{
			Statements: []Statement{
				testCase.input,
			},
		}
		assert.Equal(t, testCase.expected, program.String())
	}
}
