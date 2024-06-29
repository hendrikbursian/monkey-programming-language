package parser

import (
	"fmt"
	"github.com/hendrikbursian/monkey-programming-language/ast"
	"github.com/hendrikbursian/monkey-programming-language/lexer"
	"github.com/hendrikbursian/monkey-programming-language/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	FUNCTION_CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQUAL:               EQUALS,
	token.NOT_EQUAL:           EQUALS,
	token.LESS_THAN:           LESSGREATER,
	token.GREATER_THAN:        LESSGREATER,
	token.PLUS:                SUM,
	token.MINUS:               SUM,
	token.SLASH:               PRODUCT,
	token.ASTERISK:            PRODUCT,
	token.LEFT_PAREN:          FUNCTION_CALL,
	token.LEFT_SQUARE_BRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token
	errors       []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.INTEGER, parser.parseIntegerLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.LEFT_PAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.IF, parser.parseIfStatement)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)
	parser.registerPrefix(token.STRING, parser.parseStringLiteral)
	parser.registerPrefix(token.LEFT_SQUARE_BRACKET, parser.parseArrayLiteral)

	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQUAL, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQUAL, parser.parseInfixExpression)
	parser.registerInfix(token.LESS_THAN, parser.parseInfixExpression)
	parser.registerInfix(token.GREATER_THAN, parser.parseInfixExpression)
	parser.registerInfix(token.LEFT_PAREN, parser.parseCallExpression)
	parser.registerInfix(token.LEFT_SQUARE_BRACKET, parser.parseIndexExpression)

	parser.nextToken()
	parser.nextToken()

	return parser
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekError(tokenType token.TokenType) {
	message := fmt.Sprintf("In line %d column %d expected next token to be '%s' got '%s' instead.", parser.peekToken.Line, parser.peekToken.Column, tokenType, parser.peekToken.Type)
	parser.errors = append(parser.errors, message)

}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	parser.infixParseFns[tokenType] = fn
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.currentToken.Type != token.EOF {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		} else {
			fmt.Printf("Parser error: Token '%s' cannot be parsed\n", parser.currentToken.Type)
		}
		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{
		Token: parser.currentToken,
	}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Identifier = &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	parser.nextToken()

	statement.Value = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: parser.currentToken,
	}

	parser.nextToken()

	statement.ReturnValue = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))

	statement := &ast.ExpressionStatement{
		Token: parser.currentToken,
	}

	statement.Expression = parser.parseExpression(LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))

	prefix := parser.prefixParseFns[parser.currentToken.Type]
	if prefix == nil {
		parser.noPrefixParseFnError(parser.currentToken)
		return nil
	}

	leftExpression := prefix()

	for !parser.peekTokenIs(token.SEMICOLON) &&
		!parser.peekTokenIs(token.LEFT_BRACKET) &&
		!parser.peekTokenIs(token.RIGHT_BRACKET) &&
		precedence < parser.peekPrecedence() {

		infix := parser.infixParseFns[parser.peekToken.Type]
		if infix == nil {
			return leftExpression
		}

		parser.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.nextToken()

	expression := parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RIGHT_PAREN) {
		return nil
	}

	return expression
}

func (parser *Parser) parseIfStatement() ast.Expression {
	defer untrace(trace("parseIfStatement"))
	expression := &ast.IfExpression{
		Token: parser.currentToken,
	}

	parser.nextToken()
	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.LEFT_BRACKET) {
		return nil
	}

	expression.Consequence = parser.parseBlockStatement()

	if parser.peekTokenIs(token.ELSE) {
		parser.nextToken()

		if !parser.expectPeek(token.LEFT_BRACKET) {
			return nil
		}

		expression.Alternative = parser.parseBlockStatement()
	}

	return expression
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	defer untrace(trace("parseBlockStatement"))
	block := &ast.BlockStatement{
		Token:      parser.currentToken,
		Statements: []ast.Statement{},
	}

	parser.nextToken()
	for !parser.currentTokenIs(token.RIGHT_BRACKET) && !parser.currentTokenIs(token.EOF) {
		statement := parser.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		parser.nextToken()
	}

	return block
}

func (parser *Parser) noPrefixParseFnError(token token.Token) {
	message := fmt.Sprintf("no prefix parse function for %s at %d:%d found", token.Type, token.Line, token.Column)
	parser.errors = append(parser.errors, message)
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
	integerLiteral := &ast.IntegerLiteral{Token: parser.currentToken}

	value, err := strconv.ParseInt(parser.currentToken.Literal, 10, 64)
	if err != nil {
		message := fmt.Sprintf("Could not parse %q as integer at %d:%d", parser.currentToken.Literal, parser.currentToken.Line, parser.currentToken.Column)
		parser.errors = append(parser.errors, message)
		return nil
	}

	integerLiteral.Value = value

	return integerLiteral
}

func (parser *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	}
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	function := &ast.FunctionLiteral{
		Token: parser.currentToken,
	}

	if !parser.expectPeek(token.LEFT_PAREN) {
		return nil
	}

	function.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.LEFT_BRACKET) {
		return nil
	}

	function.Body = parser.parseBlockStatement()

	return function
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if parser.peekTokenIs(token.RIGHT_PAREN) {
		parser.nextToken()
		return identifiers
	}

	parser.nextToken()

	identifiers = append(identifiers, &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal,
	})

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		identifiers = append(identifiers, &ast.Identifier{
			Token: parser.currentToken,
			Value: parser.currentToken.Literal,
		})
	}

	if !parser.expectPeek(token.RIGHT_PAREN) {
		return nil
	}

	return identifiers
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
	}

	parser.nextToken()

	expression.Right = parser.parseExpression(PREFIX)

	return expression
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
		Left:     left,
	}

	precedence := parser.currentPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	call := &ast.CallExpression{Function: function}
	call.Arguments = parser.parseExpressionList(token.RIGHT_PAREN)
	return call
}

func (parser *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: parser.currentToken,
		Value: parser.currentTokenIs(token.TRUE),
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.currentToken}
	arr.Elements = p.parseExpressionList(token.RIGHT_SQUARE_BRACKET)

	return arr
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.currentToken,
		Left:  left,
	}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_SQUARE_BRACKET) {
		return nil
	}

	return exp
}

func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
	if parser.peekToken.Type == tokenType {
		parser.nextToken()
		return true
	} else {
		parser.peekError(tokenType)
		return false
	}
}

func (parser *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return parser.currentToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return parser.peekToken.Type == tokenType
}

func (parser *Parser) currentPrecedence() int {
	if precedence, ok := precedences[parser.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (parser *Parser) peekPrecedence() int {
	if precedence, ok := precedences[parser.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}
