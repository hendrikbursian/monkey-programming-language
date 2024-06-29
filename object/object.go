package object

import (
	"bytes"
	"fmt"
	"github.com/hendrikbursian/monkey-programming-language/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJECT      = "INTEGER"
	BOOLEAN_OBJECT      = "BOOLEAN"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION"
	STRING_OBJECT       = "STRING"
	BUILTIN_OBJECT      = "BUILTIN"
	ARRAY_OBJECT        = "ARRAY"
)

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{
		store: store,
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (env *Environment) Get(name string) (Object, bool) {
	object, ok := env.store[name]
	if !ok && env.outer != nil {
		object, ok = env.outer.Get(name)
	}

	return object, ok
}

func (env *Environment) Set(name string, value Object) Object {
	env.store[name] = value
	return value
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJECT }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJECT }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJECT }
func (rv *ReturnValue) Inspect() string  { return fmt.Sprintf("%s", rv.Value.Inspect()) }

type Error struct {
	Message string
	Line    int
	Column  int
}

func (err *Error) Type() ObjectType { return ERROR_OBJECT }
func (err *Error) Inspect() string {
	return fmt.Sprintf("Error at position %d:%d - %s", err.Line, err.Column, err.Message)
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJECT }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range f.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJECT }
func (s *String) Inspect() string {
	var buf bytes.Buffer
	buf.WriteString("\"")
	buf.WriteString(s.Value)
	buf.WriteString("\"")

	return buf.String()
}

type BuiltinFunction func(args ...Object) (Object, error)

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJECT }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJECT }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
