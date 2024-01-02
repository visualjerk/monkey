package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type precedencePriority int

const (
	_ precedencePriority = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(x)
)

type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	nextToken    token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lex}

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)

	// Read two tokens, so currentToken and nextToken are set initially
	parser.advanceTokens()
	parser.advanceTokens()

	return parser
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for !parser.currentTokenIs(token.EOF) {
		statement := parser.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		parser.advanceTokens()
	}

	return program
}

func (parser *Parser) GetErrors() []string {
	return parser.errors
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	parser.infixParseFns[tokenType] = fn
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	case token.SEMICOLON:
		return nil
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	tok := parser.currentToken

	parser.advanceToExpectedToken(token.IDENT)
	identifier := &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}

	parser.advanceToExpectedToken(token.ASSIGN)

	parser.advanceTokens()
	expression := parser.parseExpression(LOWEST)

	return &ast.LetStatement{
		Token: tok,
		Name:  identifier,
		Value: expression,
	}
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	tok := parser.currentToken

	parser.advanceTokens()
	expression := parser.parseExpression(LOWEST)

	return &ast.ReturnStatement{
		Token: tok,
		Value: expression,
	}
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	tok := parser.currentToken

	expression := parser.parseExpression(LOWEST)

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.advanceTokens()
	}

	return &ast.ExpressionStatement{
		Token: tok,
		Value: expression,
	}
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
}

func (parser *Parser) parseExpression(precedence precedencePriority) ast.Expression {
	prefixFn := parser.prefixParseFns[parser.currentToken.Type]

	if prefixFn == nil {
		return nil
	}

	leftExpression := prefixFn()

	return leftExpression
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parser.currentToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}

	return &ast.IntegerLiteral{
		Token: parser.currentToken,
		Value: value,
	}
}

func (parser *Parser) advanceTokens() {
	parser.currentToken = parser.nextToken
	parser.nextToken = parser.lexer.GetNextToken()
}

func (parser *Parser) advanceToExpectedToken(tokenType token.TokenType) bool {
	if parser.nextTokenIs(tokenType) {
		parser.advanceTokens()
		return true
	}
	parser.nextTokenError(tokenType)
	return false
}

func (parser *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return parser.currentToken.Type == tokenType
}

func (parser *Parser) nextTokenIs(tokenType token.TokenType) bool {
	return parser.nextToken.Type == tokenType
}

func (parser *Parser) nextTokenError(tokenType token.TokenType) {
	message := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		tokenType,
		parser.nextToken.Type,
	)
	parser.errors = append(parser.errors, message)
}
