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
			&object.Integer{
				Value: 6,
			},
		},
		{
			"true",
			&object.Boolean{
				Value: true,
			},
		},
		{
			"false",
			&object.Boolean{
				Value: false,
			},
		},
		{
			"",
			&object.Null{},
		},
		{
			"-6",
			&object.Integer{
				Value: -6,
			},
		},
		{
			"!true",
			&object.Boolean{
				Value: false,
			},
		},
		{
			"!false",
			&object.Boolean{
				Value: true,
			},
		},
		{
			"!!true",
			&object.Boolean{
				Value: true,
			},
		},
		{
			"2 + 2",
			&object.Integer{
				Value: 4,
			},
		},
		{
			"6 * 6",
			&object.Integer{
				Value: 36,
			},
		},
		{
			"44 - 2",
			&object.Integer{
				Value: 42,
			},
		},
		{
			"9 / 3",
			&object.Integer{
				Value: 3,
			},
		},
		{
			"(100 - 20) / 4 * 2 + 2",
			&object.Integer{
				Value: 42,
			},
		},
	}

	for _, testCase := range testCases {
		parser := parser.New(lexer.New(testCase.input))
		_, program := parser.ParseProgram()
		actual := Eval(program)

		assert.Equal(t, testCase.expected, actual)
	}
}
