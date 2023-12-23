package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT TokenType = "IDENT" // add, foobar, x, y, ...
	INT   TokenType = "INT"

	// Operators
	ASSIGN TokenType = "="
	PLUS   TokenType = "+"

	EQ TokenType = "=="

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	// Keywords
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
)

var oneCharTokens = map[byte]TokenType{
	0:   EOF,
	'=': ASSIGN,
	'+': PLUS,
	',': COMMA,
	';': SEMICOLON,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
}

func LookupOneCharToken(char byte) (TokenType, bool) {
	if tokenType, ok := oneCharTokens[char]; ok {
		return tokenType, true
	}

	return ILLEGAL, false
}

var twoCharTokens = map[string]TokenType{
	"==": EQ,
}

func LookupTwoCharToken(chars string) (TokenType, bool) {
	if tokenType, ok := twoCharTokens[chars]; ok {
		return tokenType, true
	}

	return ILLEGAL, false
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdentifier(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}
	return IDENT
}
