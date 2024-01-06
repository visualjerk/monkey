package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type operatorPrecedence int

const (
	_ operatorPrecedence = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(x)
)

var precedences = map[token.TokenType]operatorPrecedence{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

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
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.PLUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)

	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)

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

func (parser *Parser) parseExpression(precedence operatorPrecedence) ast.Expression {
	prefixFn := parser.prefixParseFns[parser.currentToken.Type]

	if prefixFn == nil {
		msg := fmt.Sprintf("no prefix parse expression for %q found", parser.currentToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}

	leftExpression := prefixFn()

	for !parser.nextTokenIs(token.SEMICOLON) && precedence < parser.nextPrecedence() {
		infixFn := parser.infixParseFns[parser.nextToken.Type]

		if infixFn == nil {
			return leftExpression
		}

		parser.advanceTokens()

		leftExpression = infixFn(leftExpression)
	}

	return leftExpression
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	prefix := parser.currentToken
	parser.advanceTokens()
	right := parser.parseExpression(PREFIX)

	return &ast.PrefixExpression{
		Token:    prefix,
		Operator: prefix.Literal,
		Right:    right,
	}
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infix := parser.currentToken
	precedence := parser.currentPrecedence()
	parser.advanceTokens()
	right := parser.parseExpression(precedence)

	return &ast.InfixExpression{
		Token:    infix,
		Operator: infix.Literal,
		Left:     left,
		Right:    right,
	}
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

func (parser *Parser) parseBoolean() ast.Expression {
	value, err := strconv.ParseBool(parser.currentToken.Literal)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as boolean", parser.currentToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}

	return &ast.Boolean{
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

func (parser *Parser) nextPrecedence() operatorPrecedence {
	if priority, ok := precedences[parser.nextToken.Type]; ok {
		return priority
	}
	return LOWEST
}

func (parser *Parser) currentPrecedence() operatorPrecedence {
	if priority, ok := precedences[parser.currentToken.Type]; ok {
		return priority
	}
	return LOWEST
}
