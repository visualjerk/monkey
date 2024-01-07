package ast

import (
	"monkey/token"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newIntegerLiteral(value int64) *IntegerLiteral {
	return &IntegerLiteral{
		Token: token.Token{
			Type:    token.INT,
			Literal: strconv.FormatInt(value, 10),
		},
		Value: value,
	}
}

func newIntegerExpression(value int64) *ExpressionStatement {
	return &ExpressionStatement{
		Token: token.Token{
			Type:    token.INT,
			Literal: strconv.FormatInt(value, 10),
		},
		Value: newIntegerLiteral(value),
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		input    Statement
		expected string
	}{
		{
			&LetStatement{
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
				Value: newIntegerLiteral(5),
			},
			"let foo = 5;",
		},
		{
			&ExpressionStatement{
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
					Right:    newIntegerLiteral(5),
				},
			},
			"(+5)",
		},
		{
			&ReturnStatement{
				Token: token.Token{
					Type:    token.RETURN,
					Literal: "return",
				},
				Value: newIntegerLiteral(5),
			},
			"return 5;",
		},
		{
			&ExpressionStatement{
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
					Left:     newIntegerLiteral(5),
					Right:    newIntegerLiteral(5),
				},
			},
			"(5 + 5)",
		},
		{
			&ExpressionStatement{
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
			"true",
		},
		{
			&ExpressionStatement{
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
							newIntegerExpression(5),
						},
					},
					Alternative: &BlockStatement{
						Token: token.Token{
							Type:    token.LBRACE,
							Literal: "{",
						},
						Statements: []Statement{
							newIntegerExpression(3),
						},
					},
				},
			},
			"if true { 5 } else { 3 }",
		},
		{
			&ExpressionStatement{
				Token: token.Token{
					Type:    token.FUNCTION,
					Literal: "fn",
				},
				Value: &FunctionLiteral{
					Token: token.Token{
						Type:    token.FUNCTION,
						Literal: "fn",
					},
					Parameters: []*Identifier{
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
								Value: newIntegerLiteral(3),
							},
						},
					},
				},
			},
			"fn(a, b) { return 3; }",
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
