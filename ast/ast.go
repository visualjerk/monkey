package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (programm *Program) TokenLiteral() string {
	if len(programm.Statements) > 0 {
		return programm.Statements[0].TokenLiteral()
	}
	return ""
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (letStatement *LetStatement) statementNode() {}
func (letStatement *LetStatement) TokenLiteral() string {
	return letStatement.Token.Literal
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (identifier *Identifier) expressionNode() {}
func (identifier *Identifier) TokenLiteral() string {
	return identifier.Token.Literal
}

type Int struct {
	Token token.Token // the token.INT token
	Value string
}

func (int *Int) expressionNode() {}
func (int *Int) TokenLiteral() string {
	return int.Token.Literal
}
