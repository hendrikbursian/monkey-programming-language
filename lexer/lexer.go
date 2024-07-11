package lexer

import (
	"fmt"
	"github.com/hendrikbursian/monkey-programming-language/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
	line         int
	column       int
}

func New(input string) *Lexer {
	lexer := &Lexer{
		input: input,
		line:  1,
	}
	lexer.readChar()
	return lexer
}

func newToken(tokenType token.TokenType, char byte, line int, column int) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(char),
		Line:    line,
		Column:  column,
	}
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.char {
	case ';':
		tok = newToken(token.SEMICOLON, lexer.char, lexer.line, lexer.column)
		break
	case '=':
		if lexer.peekChar() == '=' {
			tok.Type = token.EQUAL
			tok.Line = lexer.line
			tok.Column = lexer.column

			char := lexer.char
			lexer.readChar()
			tok.Literal = string(char) + string(lexer.char)
		} else {
			tok = newToken(token.ASSIGN, lexer.char, lexer.line, lexer.column)
		}
		break
	case '(':
		tok = newToken(token.LEFT_PAREN, lexer.char, lexer.line, lexer.column)
		break
	case ')':
		tok = newToken(token.RIGHT_PAREN, lexer.char, lexer.line, lexer.column)
		break
	case '{':
		tok = newToken(token.LEFT_CURLY_BRACE, lexer.char, lexer.line, lexer.column)
		break
	case '}':
		tok = newToken(token.RIGHT_CURLY_BRACE, lexer.char, lexer.line, lexer.column)
		break
	case '[':
		tok = newToken(token.LEFT_SQUARE_BRACKET, lexer.char, lexer.line, lexer.column)
		break
	case ']':
		tok = newToken(token.RIGHT_SQUARE_BRACKET, lexer.char, lexer.line, lexer.column)
		break
	case ',':
		tok = newToken(token.COMMA, lexer.char, lexer.line, lexer.column)
		break
	case '+':
		tok = newToken(token.PLUS, lexer.char, lexer.line, lexer.column)
		break
	case '-':
		tok = newToken(token.MINUS, lexer.char, lexer.line, lexer.column)
		break
	case '/':
		tok = newToken(token.SLASH, lexer.char, lexer.line, lexer.column)
		break
	case '*':
		tok = newToken(token.ASTERISK, lexer.char, lexer.line, lexer.column)
		break
	case '<':
		tok = newToken(token.LESS_THAN, lexer.char, lexer.line, lexer.column)
		break
	case '>':
		tok = newToken(token.GREATER_THAN, lexer.char, lexer.line, lexer.column)
		break
	case ':':
		tok = newToken(token.COLON, lexer.char, lexer.line, lexer.column)
		break
	case '.':
		tok = newToken(token.DOT, lexer.char, lexer.line, lexer.column)
		break
	case '!':
		if lexer.peekChar() == '=' {
			tok.Type = token.NOT_EQUAL
			tok.Line = lexer.line
			tok.Column = lexer.column

			char := lexer.char
			lexer.readChar()

			tok.Literal = string(char) + string(lexer.char)
		} else {
			tok = newToken(token.BANG, lexer.char, lexer.line, lexer.column)
		}
		break
	case '"':
		tok.Type = token.STRING
		tok.Line = lexer.line
		tok.Column = lexer.column
		tok.Literal = lexer.readString()
	case 0:
		tok = token.Token{
			Type:    token.EOF,
			Literal: "",
			Line:    lexer.line,
			Column:  lexer.column,
		}
		break
	default:
		if isLetter(lexer.char) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.GetTokenType(tok.Literal)
			tok.Line = lexer.line
			tok.Column = lexer.column - len(tok.Literal)
			return tok
		} else if isNumber(lexer.char) {
			tok.Type = token.INTEGER
			tok.Literal = lexer.readNumber()
			tok.Line = lexer.line
			tok.Column = lexer.column - len(tok.Literal)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.char, lexer.line, lexer.column)
			fmt.Printf("WARN: '%s'(byte: %d) at %d:%d is not a legal token", tok.Literal, tok.Literal[0], tok.Line, tok.Column)
		}
	}

	lexer.readChar()
	return tok
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.char == ' ' || lexer.char == '\t' || lexer.char == '\n' || lexer.char == '\r' {
		if lexer.char == '\n' {
			lexer.line++
			lexer.column = 0
		}
		lexer.readChar()
	}
}

func (lexer *Lexer) readString() string {
	position := lexer.position + 1

	for {
		lexer.readChar()

		if lexer.char == '"' || lexer.char == 0 {
			break
		}
	}

	return lexer.input[position:lexer.position]

}

func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	for isLetter(lexer.char) {
		lexer.readChar()
	}

	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isNumber(lexer.char) {
		lexer.readChar()
	}

	return lexer.input[position:lexer.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isNumber(char byte) bool {
	return byte('0') <= char && char <= byte('9')
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.char = 0
	} else {
		lexer.char = lexer.input[lexer.readPosition]
	}

	lexer.position = lexer.readPosition
	lexer.readPosition += 1
	lexer.column += 1
}

func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
}
