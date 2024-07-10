package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENTIFIER = "IDENTIFIER" // add, foobar, x, y, ...
	INTEGER    = "INTEGER"    // 1343456
	STRING     = "STRING"

	// Operators
	ASSIGN       = "="
	PLUS         = "+"
	MINUS        = "-"
	SLASH        = "/"
	ASTERISK     = "*"
	LESS_THAN    = "<"
	GREATER_THAN = ">"
	BANG         = "!"
	EQUAL        = "=="
	NOT_EQUAL    = "!="

	// Delimiters
	COMMA                = ","
	SEMICOLON            = ";"
	COLON                = ":"
	LEFT_PAREN           = "("
	RIGHT_PAREN          = ")"
	LEFT_CURLY_BRACE     = "{"
	RIGHT_CURLY_BRACE    = "}"
	LEFT_SQUARE_BRACKET  = "["
	RIGHT_SQUARE_BRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

func GetTokenType(identifier string) TokenType {
	if tokenType, ok := keywords[identifier]; ok {
		return tokenType
	}
	return IDENTIFIER
}
