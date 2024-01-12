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
	}

	for _, testCase := range testCases {
		parser := parser.New(lexer.New(testCase.input))
		_, program := parser.ParseProgram()
		actual := Eval(program)

		assert.Equal(t, testCase.expected, actual)
	}
}
