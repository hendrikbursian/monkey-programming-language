package ast

import (
	"bytes"
	"github.com/hendrikbursian/monkey-programming-language/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
	Line() int
	Column() int
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (program *Program) String() string {
	var out bytes.Buffer

	for _, statement := range program.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}
func (program *Program) Line() int {
	if len(program.Statements) == 0 {
		return 1
	}
	return program.Statements[0].Line()
}
func (program *Program) Column() int {
	if len(program.Statements) == 0 {
		return 1
	}
	return program.Statements[0].Column()
}

type LetStatement struct {
	Token      token.Token
	Identifier *Identifier
	Value      Expression
}

func (statement *LetStatement) statementNode()       {}
func (statement *LetStatement) TokenLiteral() string { return statement.Token.Literal }
func (statement *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(statement.TokenLiteral() + " ")
	out.WriteString(statement.Identifier.Value)
	out.WriteString(" = ")

	if statement.Value != nil {
		out.WriteString(statement.Value.String())
	}

	out.WriteString(";")

	return out.String()
}
func (statement *LetStatement) Line() int   { return statement.Token.Line }
func (statement *LetStatement) Column() int { return statement.Token.Column }

type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }
func (identifier *Identifier) String() string       { return identifier.Value }
func (identifier *Identifier) Line() int            { return identifier.Token.Line }
func (identifier *Identifier) Column() int          { return identifier.Token.Column }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (IntegerLiteral *IntegerLiteral) expressionNode()      {}
func (IntegerLiteral *IntegerLiteral) TokenLiteral() string { return IntegerLiteral.Token.Literal }
func (IntegerLiteral *IntegerLiteral) String() string       { return IntegerLiteral.Token.Literal }
func (IntegerLiteral *IntegerLiteral) Line() int            { return IntegerLiteral.Token.Line }
func (IntegerLiteral *IntegerLiteral) Column() int          { return IntegerLiteral.Token.Column }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (statement *ReturnStatement) statementNode()       {}
func (statement *ReturnStatement) TokenLiteral() string { return statement.Token.Literal }
func (statement *ReturnStatement) Line() int            { return statement.Token.Line }
func (statement *ReturnStatement) Column() int          { return statement.Token.Column }
func (statement *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(statement.Token.Literal + " ")

	if statement.ReturnValue != nil {
		out.WriteString(statement.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// TODO: dont use ExpressionStatements solo
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (statement *ExpressionStatement) statementNode()       {}
func (statement *ExpressionStatement) TokenLiteral() string { return statement.TokenLiteral() }
func (statement *ExpressionStatement) Line() int            { return statement.Token.Line }
func (statement *ExpressionStatement) Column() int          { return statement.Token.Column }
func (statement *ExpressionStatement) String() string {
	if statement.Expression == nil {
		return ""
	}

	return statement.Expression.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (expression *PrefixExpression) expressionNode()      {}
func (expression *PrefixExpression) TokenLiteral() string { return expression.Token.Literal }
func (expression *PrefixExpression) Line() int            { return expression.Token.Line }
func (expression *PrefixExpression) Column() int          { return expression.Token.Column }
func (expression *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(expression.Operator)
	out.WriteString(expression.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (expression *InfixExpression) expressionNode()      {}
func (expression *InfixExpression) TokenLiteral() string { return expression.Token.Literal }
func (expression *InfixExpression) Line() int            { return expression.Token.Line }
func (expression *InfixExpression) Column() int          { return expression.Token.Column }
func (expression *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(expression.Left.String())
	out.WriteString(" " + expression.Operator + " ")
	out.WriteString(expression.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (expression *Boolean) expressionNode()      {}
func (expression *Boolean) TokenLiteral() string { return expression.Token.Literal }
func (expression *Boolean) String() string       { return expression.Token.Literal }
func (expression *Boolean) Line() int            { return expression.Token.Line }
func (expression *Boolean) Column() int          { return expression.Token.Column }

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (block *BlockStatement) statementNode()       {}
func (block *BlockStatement) TokenLiteral() string { return block.Token.Literal }
func (block *BlockStatement) Line() int            { return block.Token.Line }
func (block *BlockStatement) Column() int          { return block.Token.Column }
func (block *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	for _, s := range block.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")

	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (expression *IfExpression) expressionNode()      {}
func (expression *IfExpression) TokenLiteral() string { return expression.Token.Literal }
func (expression *IfExpression) Line() int            { return expression.Token.Line }
func (expression *IfExpression) Column() int          { return expression.Token.Column }
func (expression *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(expression.Condition.String())
	out.WriteString(" ")
	out.WriteString(expression.Consequence.String())

	if expression.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(expression.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (expression *FunctionLiteral) expressionNode()      {}
func (expression *FunctionLiteral) TokenLiteral() string { return expression.Token.Literal }
func (expression *FunctionLiteral) Line() int            { return expression.Token.Line }
func (expression *FunctionLiteral) Column() int          { return expression.Token.Column }
func (expression *FunctionLiteral) String() string {
	var out bytes.Buffer

	parameters := []string{}
	for _, parameter := range expression.Parameters {
		parameters = append(parameters, parameter.String())
	}

	out.WriteString(expression.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ","))
	out.WriteString(") ")
	out.WriteString(expression.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression // identifier or function literal
	Arguments []Expression
}

func (expression *CallExpression) expressionNode()      {}
func (expression *CallExpression) TokenLiteral() string { return expression.Token.Literal }
func (expression *CallExpression) Line() int            { return expression.Function.Line() }
func (expression *CallExpression) Column() int          { return expression.Function.Column() }
func (expression *CallExpression) String() string {
	var out bytes.Buffer

	arguments := []string{}
	for _, argument := range expression.Arguments {
		arguments = append(arguments, argument.String())
	}

	out.WriteString(expression.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(arguments, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("\"")
	out.WriteString(sl.Token.Literal)
	out.WriteString("\"")

	return string(out.Bytes())
}

func (sl *StringLiteral) Line() int   { return sl.Token.Line }
func (sl *StringLiteral) Column() int { return sl.Token.Column }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (ar *ArrayLiteral) expressionNode()      {}
func (ar *ArrayLiteral) TokenLiteral() string { return ar.Token.Literal }
func (ar *ArrayLiteral) Line() int            { return ar.Token.Line }
func (ar *ArrayLiteral) Column() int          { return ar.Token.Column }
func (ar *ArrayLiteral) String() string {
	var buf bytes.Buffer

	elements := []string{}
	for _, value := range ar.Elements {
		elements = append(elements, value.String())
	}

	buf.WriteString("[")
	buf.WriteString(strings.Join(elements, ", "))
	buf.WriteString("]")

	return buf.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) Line() int            { return ie.Token.Line }
func (ie *IndexExpression) Column() int          { return ie.Token.Column }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) Line() int            { return hl.Token.Line }
func (hl *HashLiteral) Column() int          { return hl.Token.Column }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String(), ": ", value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
