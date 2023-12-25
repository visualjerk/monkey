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
}

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lex}

	// Read two tokens, so currentToken and nextToken are set initially
	parser.advanceTokens()
	parser.advanceTokens()

	return parser
}

func (parser *Parser) advanceTokens() {
	parser.currentToken = parser.nextToken
	parser.nextToken = parser.lexer.GetNextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for parser.currentToken.Type != token.EOF {
		statement := parser.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		parser.advanceTokens()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	if parser.currentToken.Type == token.LET {
		return parser.parseLetStatement()
	}

	return nil
}

func (parser *Parser) parseLetStatement() ast.Statement {
	tok := parser.currentToken
	parser.advanceTokens()

	identifier := parser.parseIdentifier()
	parser.advanceTokens()

	if parser.currentToken.Type != token.ASSIGN {
		parser.handleError("missing equal sign!")
	}
	parser.advanceTokens()

	expression := parser.parseExpression()

	return &ast.LetStatement{
		Token: tok,
		Name:  identifier,
		Value: expression,
	}
}

func (parser *Parser) parseIdentifier() *ast.Identifier {
	if parser.currentToken.Type != token.IDENT {
		parser.handleError("missing identifier!")
	}
	return &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
}

func (parser *Parser) parseExpression() ast.Expression {
	if parser.currentToken.Type != token.INT {
		parser.handleError("unknown expression!")
	}
	return &ast.Int{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
}

func (parser *Parser) handleError(message string) {
	fmt.Printf("unexpected token %+v\n", parser.currentToken)
	panic(message)
}
