package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
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
	identifier := parser.parseIdentifier()

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

func (parser *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
}

func (parser *Parser) parseExpression(precedence precedencePriority) ast.Expression {
	if parser.currentTokenIs(token.IDENT) {
		return parser.parseIdentifier()
	}

	if parser.currentTokenIs(token.INT) {
		if parser.nextTokenIs(token.PLUS) {
			return parser.parseAddExpression()
		}

		return parser.parseInt()
	}

	return nil
}

func (parser *Parser) parseAddExpression() ast.Expression {
	left := parser.parseInt()

	parser.advanceToExpectedToken(token.PLUS)

	tok := parser.currentToken
	parser.advanceTokens()

	right := parser.parseExpression(LOWEST)

	return &ast.AddExpression{
		Token: tok,
		Left:  left,
		Right: right,
	}
}

func (parser *Parser) parseInt() ast.Expression {
	return &ast.Int{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
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
