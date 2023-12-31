package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stringTestCase = struct {
	input    string
	expected string
}

func runStringTestCases(t *testing.T, testCases []stringTestCase) {
	for _, testCase := range testCases {
		parser := New(lexer.New(testCase.input))
		errors, program := parser.ParseProgram()

		assert.Nil(t, errors)
		assert.Equal(t, testCase.expected, program.String())
	}
}

func TestParseProgram(t *testing.T) {
	input := `
	let x = 5;

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
				Value: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "5",
					},
					Value: 5,
				},
			},
			&ast.ReturnStatement{
				Token: token.Token{
					Type:    token.RETURN,
					Literal: "return",
				},
				Value: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "5",
					},
					Value: 5,
				},
			},
		},
	}

	parser := New(lexer.New(input))
	errors, actual := parser.ParseProgram()

	assert.Nil(t, errors)
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
	actual, _ := parser.ParseProgram()
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
	errors, actual := parser.ParseProgram()

	assert.Nil(t, errors)
	assert.Equal(t, expected, actual)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	expected := &ast.Program{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Token: token.Token{
					Type:    token.INT,
					Literal: "5",
				},
				Value: &ast.IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "5",
					},
					Value: 5,
				},
			},
		},
	}

	parser := New(lexer.New(input))
	errors, actual := parser.ParseProgram()

	assert.Nil(t, errors)
	assert.Equal(t, expected, actual)
}

func TestPrefixExpression(t *testing.T) {
	testCases := []struct {
		input        string
		prefixToken  token.TokenType
		operator     string
		integerValue int64
	}{
		{"+5;", token.PLUS, "+", 5},
		{"!10;", token.BANG, "!", 10},
		{"-10;", token.MINUS, "-", 10},
	}

	for _, testCase := range testCases {
		expected := &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{
						Type:    testCase.prefixToken,
						Literal: testCase.operator,
					},
					Value: &ast.PrefixExpression{
						Token: token.Token{
							Type:    testCase.prefixToken,
							Literal: testCase.operator,
						},
						Operator: testCase.operator,
						Right: &ast.IntegerLiteral{
							Token: token.Token{
								Type:    token.INT,
								Literal: fmt.Sprintf("%d", testCase.integerValue),
							},
							Value: testCase.integerValue,
						},
					},
				},
			},
		}

		parser := New(lexer.New(testCase.input))
		errors, actual := parser.ParseProgram()

		assert.Nil(t, errors)
		assert.Equal(t, expected, actual)
	}
}

func TestInfixExpression(t *testing.T) {
	testCases := []struct {
		input      string
		infixToken token.TokenType
		operator   string
		leftValue  int64
		rightValue int64
	}{
		{"5 + 5;", token.PLUS, "+", 5, 5},
		{"5 - 5;", token.MINUS, "-", 5, 5},
		{"5 == 5;", token.EQ, "==", 5, 5},
		{"5 != 5;", token.NOT_EQ, "!=", 5, 5},
		{"5 < 5;", token.LT, "<", 5, 5},
		{"5 > 5;", token.GT, ">", 5, 5},
		{"5 * 5;", token.ASTERISK, "*", 5, 5},
		{"5 / 5;", token.SLASH, "/", 5, 5},
	}

	for _, testCase := range testCases {
		expected := &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{
						Type:    token.INT,
						Literal: fmt.Sprintf("%d", testCase.leftValue),
					},
					Value: &ast.InfixExpression{
						Token: token.Token{
							Type:    testCase.infixToken,
							Literal: testCase.operator,
						},
						Operator: testCase.operator,
						Left: &ast.IntegerLiteral{
							Token: token.Token{
								Type:    token.INT,
								Literal: fmt.Sprintf("%d", testCase.leftValue),
							},
							Value: testCase.leftValue,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{
								Type:    token.INT,
								Literal: fmt.Sprintf("%d", testCase.rightValue),
							},
							Value: testCase.rightValue,
						},
					},
				},
			},
		}

		parser := New(lexer.New(testCase.input))
		errors, actual := parser.ParseProgram()

		assert.Nil(t, errors)
		assert.Equal(t, expected, actual)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	runStringTestCases(t, []stringTestCase{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b * c",
			"(a + (b * c))",
		},
		{
			"a - b / c",
			"(a - (b / c))",
		},
		{
			"3 - 4 * 5 == 3 * 1 + 4 * 5",
			"((3 - (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"a < b == false",
			"((a < b) == false)",
		},
		{
			"5 > 3 == true",
			"((5 > 3) == true)",
		},
		{
			"(a - b) / c",
			"((a - b) / c)",
		},
		{
			"-(3 + 2)",
			"(-(3 + 2))",
		},
	})
}

func TestIfExpressions(t *testing.T) {
	runStringTestCases(t, []stringTestCase{
		{
			"if (x < y) { x }",
			"if (x < y) { x }",
		},
		{
			"if (x < y) { x } else { y }",
			"if (x < y) { x } else { y }",
		},
	})
}

func TestFunctionLiteral(t *testing.T) {
	runStringTestCases(t, []stringTestCase{
		{
			"fn() { return 5; }",
			"fn() { return 5; }",
		},
		{
			"fn(a) { return a; }",
			"fn(a) { return a; }",
		},
		{
			"fn(a, b, c) { return a + b + c; }",
			"fn(a, b, c) { return ((a + b) + c); }",
		},
	})
}

func TestFunctionLiteralParserErrors(t *testing.T) {
	input := `fn(a b) return a;`

	expected := []string{
		"expected next token to be ), got IDENT instead",
		"expected next token to be {, got IDENT instead",
		"no prefix parse expression for ) found",
	}

	parser := New(lexer.New(input))
	actual, _ := parser.ParseProgram()

	assert.Equal(t, expected, actual)
}

func TestCallExpression(t *testing.T) {
	runStringTestCases(t, []stringTestCase{
		{
			"foo()",
			"foo()",
		},
		{
			"foo(a)",
			"foo(a)",
		},
		{
			"foo(a, b)",
			"foo(a, b)",
		},
	})
}
