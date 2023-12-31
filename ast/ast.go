package ast

import (
	"bytes"
	"monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (programm *Program) String() string {
	var out bytes.Buffer

	for _, statement := range programm.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
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

func (letStatement *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(letStatement.TokenLiteral() + " ")
	out.WriteString(letStatement.Name.String() + " = ")
	out.WriteString(letStatement.Value.String() + ";")

	return out.String()
}

type ReturnStatement struct {
	Token token.Token // the token.RETURN token
	Value Expression
}

func (returnStatement *ReturnStatement) statementNode() {}
func (returnStatement *ReturnStatement) TokenLiteral() string {
	return returnStatement.Token.Literal
}

func (returnStatement *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(returnStatement.TokenLiteral() + "")
	out.WriteString(returnStatement.Value.String() + ";")

	return out.String()
}

type ExpressionStatement struct {
	Token token.Token // the first token of the expression
	Value Expression
}

func (expressionStatement *ExpressionStatement) statementNode() {}
func (expressionStatement *ExpressionStatement) TokenLiteral() string {
	return expressionStatement.Token.Literal
}
func (expressionStatement *ExpressionStatement) String() string {
	return expressionStatement.Value.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (identifier *Identifier) expressionNode() {}
func (identifier *Identifier) TokenLiteral() string {
	return identifier.Token.Literal
}
func (identifier *Identifier) String() string {
	return identifier.Value
}

type Int struct {
	Token token.Token // the token.INT token
	Value string
}

func (int *Int) expressionNode() {}
func (int *Int) TokenLiteral() string {
	return int.Token.Literal
}
func (int *Int) String() string {
	return int.Value
}

type AddExpression struct {
	Token token.Token // the token.PLUS token
	Left  Expression
	Right Expression
}

func (addExpression *AddExpression) expressionNode() {}
func (addExpression *AddExpression) TokenLiteral() string {
	return addExpression.Token.Literal
}
func (addExpression *AddExpression) String() string {
	var out bytes.Buffer

	out.WriteString(addExpression.Left.String() + " ")
	out.WriteString(addExpression.TokenLiteral() + " ")
	out.WriteString(addExpression.Right.String())

	return out.String()
}
