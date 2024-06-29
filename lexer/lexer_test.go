package lexer

import (
	"fmt"
	"testing"

	"github.com/hendrikbursian/monkey-programming-language/token"
)

func TestNextToken(t *testing.T) {

	code := `let five = 5;
let ten = 10;

let add = fn(x, y) {
    x + y;
};

let result = add(five, ten);

!-/*5;
5 < 10 > 5;

if 5 < 10 {
    return true;
} else {
    return false;
}

true == true
true != false

"foobar"
"foo bar"

[2, "hallo"]
    `
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.LET, "let", 1, 1},
		{token.IDENTIFIER, "five", 1, 5},
		{token.ASSIGN, "=", 1, 10},
		{token.INTEGER, "5", 1, 12},
		{token.SEMICOLON, ";", 1, 13},

		{token.LET, "let", 2, 1},
		{token.IDENTIFIER, "ten", 2, 5},
		{token.ASSIGN, "=", 2, 9},
		{token.INTEGER, "10", 2, 11},
		{token.SEMICOLON, ";", 2, 13},

		{token.LET, "let", 4, 1},
		{token.IDENTIFIER, "add", 4, 5},
		{token.ASSIGN, "=", 4, 9},
		{token.FUNCTION, "fn", 4, 11},
		{token.LEFT_PAREN, "(", 4, 13},
		{token.IDENTIFIER, "x", 4, 14},
		{token.COMMA, ",", 4, 15},
		{token.IDENTIFIER, "y", 4, 17},
		{token.RIGHT_PAREN, ")", 4, 18},
		{token.LEFT_BRACKET, "{", 4, 20},

		{token.IDENTIFIER, "x", 5, 5},
		{token.PLUS, "+", 5, 7},
		{token.IDENTIFIER, "y", 5, 9},
		{token.SEMICOLON, ";", 5, 10},

		{token.RIGHT_BRACKET, "}", 6, 1},
		{token.SEMICOLON, ";", 6, 2},

		{token.LET, "let", 8, 1},
		{token.IDENTIFIER, "result", 8, 5},
		{token.ASSIGN, "=", 8, 12},
		{token.IDENTIFIER, "add", 8, 14},
		{token.LEFT_PAREN, "(", 8, 17},
		{token.IDENTIFIER, "five", 8, 18},
		{token.COMMA, ",", 8, 22},
		{token.IDENTIFIER, "ten", 8, 24},
		{token.RIGHT_PAREN, ")", 8, 27},
		{token.SEMICOLON, ";", 8, 28},

		// !-/*5;
		{token.BANG, "!", 10, 1},
		{token.MINUS, "-", 10, 2},
		{token.SLASH, "/", 10, 3},
		{token.ASTERISK, "*", 10, 4},
		{token.INTEGER, "5", 10, 5},
		{token.SEMICOLON, ";", 10, 6},

		// 5 < 10 > 5;
		{token.INTEGER, "5", 11, 1},
		{token.LESS_THAN, "<", 11, 3},
		{token.INTEGER, "10", 11, 5},
		{token.GREATER_THAN, ">", 11, 8},
		{token.INTEGER, "5", 11, 10},
		{token.SEMICOLON, ";", 11, 11},

		// if (5 < 10) {
		{token.IF, "if", 13, 1},
		{token.INTEGER, "5", 13, 4},
		{token.LESS_THAN, "<", 13, 6},
		{token.INTEGER, "10", 13, 8},
		{token.LEFT_BRACKET, "{", 13, 11},

		//     return true;
		{token.RETURN, "return", 14, 5},
		{token.TRUE, "true", 14, 12},
		{token.SEMICOLON, ";", 14, 16},

		// } else {
		{token.RIGHT_BRACKET, "}", 15, 1},
		{token.ELSE, "else", 15, 3},
		{token.LEFT_BRACKET, "{", 15, 8},

		//     return false;
		{token.RETURN, "return", 16, 5},
		{token.FALSE, "false", 16, 12},
		{token.SEMICOLON, ";", 16, 17},

		// }
		{token.RIGHT_BRACKET, "}", 17, 1},

		// true == true
		{token.TRUE, "true", 19, 1},
		{token.EQUAL, "==", 19, 6},
		{token.TRUE, "true", 19, 9},

		// true != false
		{token.TRUE, "true", 20, 1},
		{token.NOT_EQUAL, "!=", 20, 6},
		{token.FALSE, "false", 20, 9},

		{token.STRING, "foobar", 22, 1},
		{token.STRING, "foo bar", 23, 1},

		// [2, "hallo"]
		{token.LEFT_SQUARE_BRACKET, "[", 25, 1},
		{token.INTEGER, "2", 25, 2},
		{token.COMMA, ",", 25, 3},
		{token.STRING, "hallo", 25, 5},
		{token.RIGHT_SQUARE_BRACKET, "]", 25, 12},

		{token.EOF, "", 26, 5},
	}

	l := New(code)

	for i, tt := range tests {
		tok := l.NextToken()
		t.Run(fmt.Sprintf("test[%d] token: %s", i, tt.expectedType), func(t *testing.T) {
			if tok.Type != tt.expectedType {
				t.Errorf("test [%d] TokenType %s not correct. Should be: %s", i, tok.Type, tt.expectedType)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("test [%d] TokenLiteral %s not correct. Should be: %s", i, tok.Literal, tt.expectedLiteral)
			}

			if tok.Line != tt.expectedLine {
				t.Errorf("test [%d] TokenLine %d not correct. Should be: %d", i, tok.Line, tt.expectedLine)
			}

			if tok.Column != tt.expectedColumn {
				t.Errorf("test [%d] TokenColumn %d not correct. Should be: %d", i, tok.Column, tt.expectedColumn)
			}

		})
	}

}
