package eval

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	testCases := []struct {
		input    string
		expected object.Object
	}{
		{
			"6",
			&object.Integer{Value: 6},
		},
		{
			"true",
			&object.Boolean{Value: true},
		},
		{
			"false",
			&object.Boolean{Value: false},
		},
		{
			"",
			&object.Null{},
		},
		{
			"-6",
			&object.Integer{Value: -6},
		},
		{
			"!true",
			&object.Boolean{Value: false},
		},
		{
			"!false",
			&object.Boolean{Value: true},
		},
		{
			"!!true",
			&object.Boolean{Value: true},
		},
		{
			"2 + 2",
			&object.Integer{Value: 4},
		},
		{
			"6 * 6",
			&object.Integer{Value: 36},
		},
		{
			"44 - 2",
			&object.Integer{Value: 42},
		},
		{
			"9 / 3",
			&object.Integer{Value: 3},
		},
		{
			"(100 - 20) / 4 * 2 + 2",
			&object.Integer{Value: 42},
		},
		{
			"9 > 3",
			&object.Boolean{Value: true},
		},
		{
			"9 < 3",
			&object.Boolean{Value: false},
		},
		{
			"3 == 3",
			&object.Boolean{Value: true},
		},
		{
			"10 == 3",
			&object.Boolean{Value: false},
		},
		{
			"10 != 3",
			&object.Boolean{Value: true},
		},
		{
			"3 != 3",
			&object.Boolean{Value: false},
		},
		{
			"true == true",
			&object.Boolean{Value: true},
		},
		{
			"true == false",
			&object.Boolean{Value: false},
		},
		{
			"true != false",
			&object.Boolean{Value: true},
		},
		{
			"false != false",
			&object.Boolean{Value: false},
		},
		{
			"if (true) { 5; }",
			&object.Integer{Value: 5},
		},
		{
			"if (2 < 1) { 5; }",
			nil,
		},
		{
			"if (10 == 5) { 5; } else { 10; }",
			&object.Integer{Value: 10},
		},
		{
			"if (5 - 2 > 20) { true; } else { false; }",
			&object.Boolean{Value: false},
		},
		{
			"return 10;",
			&object.Integer{Value: 10},
		},
		{
			"5; return 7; 5;",
			&object.Integer{Value: 7},
		},
		{
			"return 2; return 3;",
			&object.Integer{Value: 2},
		},
		{
			`if (10 > 5) {
				if (10 > 5) {
					return 10;
				}
				return 5;
			}`,
			&object.Integer{Value: 10},
		},
		{
			"let x = 5; x;",
			&object.Integer{Value: 5},
		},
		{
			"let x = 5; let y = 2; x * y;",
			&object.Integer{Value: 10},
		},
		{
			"let x = 5;",
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			parser := parser.New(lexer.New(testCase.input))
			_, program := parser.ParseProgram()
			actual := Eval(program, object.NewEnvironment())

			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestEvalErrors(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"true + true;",
			"Error: unknown operation BOOLEAN + BOOLEAN",
		},
		{
			"true + true; true;",
			"Error: unknown operation BOOLEAN + BOOLEAN",
		},
		{
			"-true;",
			"Error: unknown operation -BOOLEAN",
		},
		{
			"5 + true;",
			"Error: type mismatch INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Error: type mismatch INTEGER + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"Error: unknown operation BOOLEAN + BOOLEAN",
		},
		{
			"let a = 5; b;",
			"Error: identifier not found b",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			parser := parser.New(lexer.New(testCase.input))
			_, program := parser.ParseProgram()
			actual := Eval(program, object.NewEnvironment())

			assert.Equal(t, testCase.expected, actual.Inspect())
		})
	}
}
