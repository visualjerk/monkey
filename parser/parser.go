package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	nextToken    token.Token

	errors []string
}

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

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	default:
		return nil
	}
}

func (parser *Parser) parseLetStatement() ast.Statement {
	tok := parser.currentToken

	parser.advanceToExpectedToken(token.IDENT)
	identifier := parser.parseIdentifier()

	parser.advanceToExpectedToken(token.ASSIGN)

	parser.advanceTokens()
	expression := parser.parseExpression()

	return &ast.LetStatement{
		Token: tok,
		Name:  identifier,
		Value: expression,
	}
}

func (parser *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
}

func (parser *Parser) parseExpression() ast.Expression {
	if parser.currentTokenIs(token.INT) {
		if parser.nextTokenIs(token.PLUS) {
			return parser.parseAddExpression()
		}

		return parser.parseInt()
	}

	// parser.expectError("unknown expression!")
	return nil
}

func (parser *Parser) parseAddExpression() ast.Expression {
	left := parser.parseInt()

	parser.advanceToExpectedToken(token.PLUS)

	tok := parser.currentToken
	parser.advanceTokens()

	right := parser.parseExpression()

	return &ast.AddExpression{
		Token: tok,
		Left:  left,
		Right: right,
	}
}

func (parser *Parser) parseInt() ast.Expression {
	if !parser.currentTokenIs(token.INT) {
		// parser.expectError("not an int!")
	}

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
	parser.tokenTypeError(tokenType)
	return false
}

func (parser *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return parser.currentToken.Type == tokenType
}

func (parser *Parser) nextTokenIs(tokenType token.TokenType) bool {
	return parser.nextToken.Type == tokenType
}

func (parser *Parser) tokenTypeError(tokenType token.TokenType) {
	message := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		tokenType,
		parser.nextToken.Type,
	)
	parser.errors = append(parser.errors, message)
}
