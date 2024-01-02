package ast

import (
	"monkey/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
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
				Value: &IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "5",
					},
					Value: 5,
				},
			},
		},
	}

	expected := "let foo = 5;"

	assert.Equal(t, expected, program.String())
}
