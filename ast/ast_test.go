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
			input: &ReturnStatement{
				Token: token.Token{
					Type:    token.RETURN,
					Literal: "return",
				},
				Value: &IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "5",
					},
					Value: 5,
				},
			},
			expected: "return 5;",
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
		{
			input: &ExpressionStatement{
				Token: token.Token{
					Type:    token.TRUE,
					Literal: "true",
				},
				Value: &Boolean{
					Token: token.Token{
						Type:    token.TRUE,
						Literal: "true",
					},
					Value: true,
				},
			},
			expected: "true",
		},
		{
			input: &ExpressionStatement{
				Token: token.Token{
					Type:    token.IF,
					Literal: "if",
				},
				Value: &IfExpression{
					Token: token.Token{
						Type:    token.IF,
						Literal: "if",
					},
					Condition: &Boolean{
						Token: token.Token{
							Type:    token.TRUE,
							Literal: "true",
						},
						Value: true,
					},
					Consequence: &BlockStatement{
						Token: token.Token{
							Type:    token.LBRACE,
							Literal: "{",
						},
						Statements: []Statement{
							&ExpressionStatement{
								Token: token.Token{
									Type:    token.INT,
									Literal: "5",
								},
								Value: &IntegerLiteral{
									Token: token.Token{
										Type:    token.INT,
										Literal: "5",
									},
									Value: 5,
								},
							},
						},
					},
					Alternative: &BlockStatement{
						Token: token.Token{
							Type:    token.LBRACE,
							Literal: "{",
						},
						Statements: []Statement{
							&ExpressionStatement{
								Token: token.Token{
									Type:    token.INT,
									Literal: "3",
								},
								Value: &IntegerLiteral{
									Token: token.Token{
										Type:    token.INT,
										Literal: "3",
									},
									Value: 3,
								},
							},
						},
					},
				},
			},
			expected: "if true { 5 } else { 3 }",
		},
		{
			input: &ExpressionStatement{
				Token: token.Token{
					Type:    token.FUNCTION,
					Literal: "fn",
				},
				Value: &FunctionLiteral{
					Token: token.Token{
						Type:    token.FUNCTION,
						Literal: "fn",
					},
					Arguments: []*Identifier{
						{
							Token: token.Token{
								Type:    token.IDENT,
								Literal: "a",
							},
							Value: "a",
						},
						{
							Token: token.Token{
								Type:    token.IDENT,
								Literal: "b",
							},
							Value: "b",
						},
					},
					Body: &BlockStatement{
						Token: token.Token{
							Type:    token.LBRACE,
							Literal: "{",
						},
						Statements: []Statement{
							&ReturnStatement{
								Token: token.Token{
									Type:    token.RETURN,
									Literal: "return",
								},
								Value: &IntegerLiteral{
									Token: token.Token{
										Type:    token.INT,
										Literal: "3",
									},
									Value: 3,
								},
							},
						},
					},
				},
			},
			expected: "fn(a, b) { return 3; }",
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
