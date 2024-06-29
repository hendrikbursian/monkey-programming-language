package ast

import (
	"fmt"
	"github.com/hendrikbursian/monkey-programming-language/token"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		program        Program
		expectedString string
	}{
		{
			program: Program{
				Statements: []Statement{
					&LetStatement{
						Token: token.Token{
							Type:    token.LET,
							Line:    1,
							Column:  1,
							Literal: "let",
						},
						Identifier: &Identifier{
							Token: token.Token{
								Type:    token.IDENTIFIER,
								Literal: "myVar",
								Line:    1,
								Column:  5,
							},
							Value: "myVar",
						},
						Value: &Identifier{
							Token: token.Token{
								Type:    token.IDENTIFIER,
								Literal: "anotherVar",
								Line:    1,
								Column:  14,
							},
							Value: "anotherVar",
						},
					},
				},
			},
			expectedString: "let myVar = anotherVar;",
		},
		{
			program: Program{
				Statements: []Statement{

					&LetStatement{
						Token: token.Token{
							Type:    token.LET,
							Line:    1,
							Column:  1,
							Literal: "let",
						},
						Identifier: &Identifier{
							Token: token.Token{
								Type:    token.IDENTIFIER,
								Literal: "myVar",
								Line:    1,
								Column:  5,
							},
							Value: "myVar",
						},
						Value: &ArrayLiteral{
							Token: token.Token{
								Type:    token.LEFT_SQUARE_BRACKET,
								Line:    1,
								Column:  1,
								Literal: "[",
							},
							Values: []Expression{
								&IntegerLiteral{
									Token: token.Token{
										Type:    token.INTEGER,
										Literal: "4",
										Line:    1,
										Column:  2,
									},
									Value: 4,
								},
								&StringLiteral{
									Token: token.Token{
										Type:    token.STRING,
										Literal: "\"hallo\"",
										Line:    1,
										Column:  2,
									},
									Value: "hallo",
								},
							},
						},
					},
				},
			},
			expectedString: `let myVar = [4, "hallo"];`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {

		})

		if tt.program.String() != tt.expectedString {
			t.Errorf("program.String() wrong, expected: %q. got=%q", tt.expectedString, tt.program.String())
		}
	}
}