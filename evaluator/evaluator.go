package evaluator

import (
	"fmt"
	"github.com/hendrikbursian/monkey-programming-language/ast"
	"github.com/hendrikbursian/monkey-programming-language/object"
	"strings"
)

var (
	TRUE        = &object.Boolean{Value: true}
	FALSE       = &object.Boolean{Value: false}
	EMPTY_MAYBE = object.Maybe{Value: nil}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.PropertyExpression:
		return evalPropertyExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.FunctionLiteral:
		return evalFunction(node, env)
	case *ast.ArrayLiteral:
		return evalArray(node, env)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.Boolean:
		if node.Value {
			return TRUE
		} else {
			return FALSE
		}
	default:
		return nil
	}
}

func evalReturnStatement(node *ast.ReturnStatement, env *object.Environment) object.Object {
	value := Eval(node.ReturnValue, env)
	if isError(value) {
		return value
	}
	return &object.ReturnValue{Value: value}
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	switch node.Operator {
	case "!":
		return evalBangOperatorExpression(node, right)
	case "-":
		return evalMinusOperatorExpression(node, right)
	default:
		return newError(node.Line(), node.Column(), "unknown operator: %s%s", node.Operator, right.Type())
	}
}

func evalBangOperatorExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default:
		return newError(node.Line(), node.Column(), "unknown operator: %s%s", node.Operator, right.Type())
	}
}

func evalMinusOperatorExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJECT {
		return newError(node.Line(), node.Column(), "unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: -value}
}

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalIntegerInfixExpression(node, left, right)
	case left.Type() == object.STRING_OBJECT && right.Type() == object.STRING_OBJECT:
		return evalStringInfixExpression(node, left, right)
	case left.Type() == object.STRING_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalStringIntegerInfixExpression(node, left, right)
	case node.Operator == "==":
		return getBooleanObject(left == right)
	case node.Operator == "!=":
		return getBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError(node.Right.Line(), node.Right.Column(), "type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
	default:
		return newError(node.Line(), node.Column(), "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalStringInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch node.Operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return getBooleanObject(leftValue == rightValue)
	case "!=":
		return getBooleanObject(leftValue != rightValue)
	default:
		return newError(node.Line(), node.Column(), "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalStringIntegerInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	string := left.(*object.String).Value
	integer := right.(*object.Integer).Value

	switch node.Operator {
	case "*":
		return &object.String{Value: strings.Repeat(string, int(integer))}
	default:
		return newError(node.Line(), node.Column(), "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalIntegerInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch node.Operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return getBooleanObject(leftValue < rightValue)
	case ">":
		return getBooleanObject(leftValue > rightValue)
	case "==":
		return getBooleanObject(leftValue == rightValue)
	case "!=":
		return getBooleanObject(leftValue != rightValue)
	default:
		return newError(node.Line(), node.Column(), "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	switch condition {
	case TRUE:
		value := Eval(ie.Consequence, env)
		if isError(value) {
			return value
		}

		return wrapMaybe(value)
	case FALSE:
		if ie.Alternative == nil {
			return &EMPTY_MAYBE
		}

		value := Eval(ie.Alternative, env)
		if isError(value) {
			return value
		}

		return wrapMaybe(value)
	default:
		return &EMPTY_MAYBE
	}
}

func wrapMaybe(value object.Object) object.Object {
	switch obj := value.(type) {
	case *object.ReturnValue:
		obj.Value = wrapMaybe(obj.Value)
		return obj
	case *object.Maybe:
		return obj
	default:
		return &object.Maybe{Value: value}
	}
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJECT || rt == object.ERROR_OBJECT {
				return result
			}
		}
	}

	return result
}

func evalLetStatement(node *ast.LetStatement, env *object.Environment) object.Object {
	value := Eval(node.Value, env)
	if isError(value) {
		return value
	}

	env.Set(node.Identifier.Value, value)
	return nil
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if identifier, ok := env.Get(node.Value); ok {
		return identifier
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError(node.Line(), node.Column(), "identifier not found: %s", node.Value)
}

func evalFunction(node *ast.FunctionLiteral, env *object.Environment) object.Object {
	return &object.Function{
		Parameters: node.Parameters,
		Env:        env,
		Body:       node.Body,
	}
}

func evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	function := Eval(node.Function, env)
	if isError(function) {
		return function
	}

	args := evalExpressions(node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	switch fn := function.(type) {
	case *object.Function:
		extendedEnv := object.NewEnclosedEnvironment(fn.Env)
		if len(args) < len(fn.Parameters) {
			missingParameters := []string{}

			for i := len(args); i < len(fn.Parameters); i++ {
				missingParameters = append(missingParameters, fn.Parameters[i].Value)
			}

			return newError(node.Line(), node.Column(), "missing parameters %q in function call", strings.Join(missingParameters, ", "))
		}
		for paramIdx, param := range fn.Parameters {
			env.Set(param.Value, args[paramIdx])
		}
		evaluated := Eval(fn.Body, extendedEnv)

		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		result, err := fn.Fn(args...)
		if err != nil {
			return newError(node.Line(), node.Column(), err.Error())
		}
		return result
	default:
		return newError(node.Line(), node.Column(), "not a function: %s", function)
	}
}

func evalArray(node *ast.ArrayLiteral, env *object.Environment) object.Object {
	obj := &object.Array{}
	obj.Elements = evalExpressions(node.Elements, env)
	return obj
}

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	index := Eval(node.Index, env)
	if isError(index) {
		return index
	}

	switch left.Type() {
	case object.ARRAY_OBJECT:
		if index.Type() != object.INTEGER_OBJECT {
			return newError(node.Index.Line(), node.Index.Column(), "cannot use %s as index for array", index.Type())
		}
		arrObj := left.(*object.Array)
		idxValue := index.(*object.Integer).Value

		if int(idxValue) < 0 || int(idxValue) >= len(arrObj.Elements) {
			return &EMPTY_MAYBE
		}

		return wrapMaybe(arrObj.Elements[idxValue])
	case object.HASH_OBJECT:
		hashableIndex, ok := index.(object.Hashable)
		if !ok {
			return newError(node.Index.Line(), node.Index.Column(), "can not use index of type %s for hash", index.Type())
		}

		hash := left.(*object.Hash)
		value, ok := hash.Pairs[hashableIndex.HashKey()]
		if !ok {
			return &EMPTY_MAYBE
		}

		return wrapMaybe(value.Value)
	default:
		return newError(node.Line(), node.Column(), "cannot use index of %s", left.Type())
	}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for key, value := range node.Pairs {
		keyObj := Eval(key, env)

		if isError(keyObj) {
			return keyObj
		}

		hashableKey, ok := keyObj.(object.Hashable)
		if !ok {
			return newError(value.Line(), value.Column(), "cannot use type %s as key for hash", keyObj.Type())
		}

		valueObj := Eval(value, env)
		if isError(valueObj) {
			return valueObj
		}

		hashed := hashableKey.HashKey()
		pairs[hashed] = object.HashPair{Key: keyObj, Value: valueObj}
	}

	return &object.Hash{Pairs: pairs}
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		evaluated := Eval(expression, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}
func evalPropertyExpression(prop *ast.PropertyExpression, env *object.Environment) object.Object {
	subject := Eval(prop.Subject, env)
	if isError(subject) {
		return subject
	}

	switch sub := subject.(type) {
	case *object.Maybe:
		switch prop.Name.Value {
		case "hasValue":
			return getBooleanObject(sub.Value != nil)
		case "value":
			if sub.Value == nil {
				return newError(prop.Line(), prop.Column(), "%q has no value! check before with \"hasValue\"!", prop.String())
			}

			return sub.Value
		}
	}

	return newError(prop.Line(), prop.Column(), "%s has no property %q.", subject.Type(), prop.Name.TokenLiteral())
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func getBooleanObject(value bool) object.Object {
	if value {
		return TRUE
	}

	return FALSE
}

func newError(line, column int, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
		Line:    line,
		Column:  column,
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJECT
	}

	return false
}
