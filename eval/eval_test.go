package eval

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalInteger(t *testing.T) {
	input := "6"
	expected := &object.Integer{
		Value: 6,
	}

	parser := parser.New(lexer.New(input))
	_, program := parser.ParseProgram()
	actual := Eval(program)

	assert.Equal(t, expected, actual)
}
