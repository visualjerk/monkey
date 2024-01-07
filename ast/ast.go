package ast

import (
	"bytes"
	"monkey/token"
	"strings"
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

	out.WriteString(returnStatement.TokenLiteral() + " ")
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

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (blockStatement *BlockStatement) statementNode() {}
func (blockStatement *BlockStatement) TokenLiteral() string {
	return blockStatement.Token.Literal
}
func (blockStatement *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	for _, statement := range blockStatement.Statements {
		out.WriteString(statement.String())
	}
	out.WriteString(" }")

	return out.String()
}

type IfExpression struct {
	Token       token.Token // the if token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExpression *IfExpression) expressionNode() {}
func (ifExpression *IfExpression) TokenLiteral() string {
	return ifExpression.Token.Literal
}
func (ifExpression *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ifExpression.Condition.String())
	out.WriteString(" ")
	out.WriteString(ifExpression.Consequence.String())

	if ifExpression.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ifExpression.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // the fn token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (functionLiteral *FunctionLiteral) expressionNode() {}
func (functionLiteral *FunctionLiteral) TokenLiteral() string {
	return functionLiteral.Token.Literal
}
func (functionLiteral *FunctionLiteral) String() string {
	var out bytes.Buffer

	parameters := []string{}
	for _, argument := range functionLiteral.Parameters {
		parameters = append(parameters, argument.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") ")
	out.WriteString(functionLiteral.Body.String())

	return out.String()
}

type PrefixExpression struct {
	Token    token.Token // the prefix token
	Operator string
	Right    Expression
}

func (prefixExpression *PrefixExpression) expressionNode() {}
func (prefixExpression *PrefixExpression) TokenLiteral() string {
	return prefixExpression.Token.Literal
}
func (prefixExpression *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(prefixExpression.Operator)
	out.WriteString(prefixExpression.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // the infix token
	Operator string
	Left     Expression
	Right    Expression
}

func (infixExpression *InfixExpression) expressionNode() {}
func (infixExpression *InfixExpression) TokenLiteral() string {
	return infixExpression.Token.Literal
}
func (infixExpression *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(infixExpression.Left.String())
	out.WriteString(" " + infixExpression.Operator + " ")
	out.WriteString(infixExpression.Right.String())
	out.WriteString(")")

	return out.String()
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

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (int *IntegerLiteral) expressionNode() {}
func (int *IntegerLiteral) TokenLiteral() string {
	return int.Token.Literal
}
func (int *IntegerLiteral) String() string {
	return int.Token.Literal
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (boolean *Boolean) expressionNode() {}
func (boolean *Boolean) TokenLiteral() string {
	return boolean.Token.Literal
}
func (boolean *Boolean) String() string {
	return boolean.Token.Literal
}
