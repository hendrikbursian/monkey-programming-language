package evaluator

import (
	"errors"
	"fmt"
	"github.com/hendrikbursian/monkey-programming-language/lexer"
	"github.com/hendrikbursian/monkey-programming-language/object"
	"github.com/hendrikbursian/monkey-programming-language/parser"
	"strconv"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5+5+5+5-10", 10},
		{"2*2*2*2*2", 32},
		{"-50+100+-50", 0},
		{"5*2+10", 20},
		{"5+2*10", 25},
		{"20+2*-10", 0},
		{"50/2*2+10", 60},
		{"2*(5+10)", 30},
		{"3*3*3+10", 37},
		{"(5+10*2+15/3)*2+-10", 50},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			evaluated := testEval(test.input)
			testIntegerObject(t, evaluated, test.expected)
		})
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1<2", true},
		{"1>2", false},
		{"1<1", false},
		{"1>1", false},
		{"1==1", true},
		{"1!=1", false},
		{"1==2", false},
		{"1!=2", true},
		{"true==false", false},
		{"false==false", true},
		{"true==false", false},
		{"true!=false", true},
		{"false!=true", true},
		{"(1<2)==true", true},
		{"(1<2)==false", false},
		{"(1>2)==true", false},
		{"(1>2)==false", true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
	}

	for _, tt := range tests {
		res := testEval(tt.input)
		testBooleanObject(t, res, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", 10},
		{"if false { 10 }", nil},
		{"if 1 < 2 { 10 }", 10},
		{"if 1 > 2 { 10 }", nil},
		{"if 1 > 2 { 10 } else { 20 }", 20},
		{"if 1 < 2 { 10 } else { 20 }", 10},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			evaluated := testEval(test.input)
			integer, ok := test.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				if evaluated != nil {
					t.Errorf("evaluated is not nil. got=%T (%+v)", evaluated, evaluated)
				}
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2*5; 9;", 10},
		{"9; return 2*5; 9;", 10},
		{`if 10 > 1 {
            if 10 > 1 {
              return 10;
            }

            return 1
          }`, 10},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			evaluated := testEval(test.input)
			testIntegerObject(t, evaluated, test.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
		expectedLine    int
		expectedColumn  int
	}{
		{
			"5+true;",
			"type mismatch: INTEGER + BOOLEAN, expecting: INTEGER",
			1, 3,
		},
		{
			"10;5+true; 5;",
			"type mismatch: INTEGER + BOOLEAN, expecting: INTEGER",
			1, 6,
		},
		{
			"-true;",
			"unknown operator: -BOOLEAN",
			1, 1,
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
			1, 6,
		},
		{
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
			1, 9,
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
			1, 20,
		},
		{
			"\"hello\" - \"world\"",
			"unknown operator: STRING - STRING",
			1, 9,
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}`,
			"unknown operator: BOOLEAN + BOOLEAN",
			4, 17,
		},
		{
			"foobar",
			"identifier not found: foobar",
			1, 1,
		},
		{
			`
let test = fn(x, y, z){return x*y*z};
test(3)`,
			"missing parameters \"y, z\" in function call",
			3, 1,
		},
		{
			`
let test = fn(x, y){return x*y*z};
test(3)`,
			"missing parameters \"y\" in function call",
			3, 1,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, test.input), func(t *testing.T) {
			evaluated := testEval(test.input)
			errorObject, ok := evaluated.(*object.Error)

			if !ok {
				t.Errorf("no error object returned. got=%T", evaluated)
			}

			if errorObject.Message != test.expectedMessage {
				t.Errorf("incorrect error message. want=%q, got=%q", test.expectedMessage, errorObject.Message)

			}

			if errorObject.Line != test.expectedLine {
				t.Errorf("incorrect line. want=%d, got=%d", test.expectedLine, errorObject.Line)
			}

			if errorObject.Column != test.expectedColumn {
				t.Errorf("incorrect column. want=%d, got=%d", test.expectedColumn, errorObject.Column)
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5*5; a;", 25},
		{"let a = 5; let b = a; b", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c", 15},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			evaluated := testEval(test.input)
			testIntegerObject(t, evaluated, test.expected)
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("parameter length is incorrect. want=%d got=%d", 1, len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "{ (x + 2) }"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x){x;}; identity(5)", 5},
		{"let identity = fn(x){return x;}; identity(5)", 5},
		{"let double = fn(x){x*2;}; double(5)", 10},
		{"let add = fn(x, y){x+y}; add(5, 3)", 8},
		{"let add = fn(x, y){x+y}; add(add(5, 3), 5+5)", 18},
		{"fn(x){x;}(5)", 5},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			evaluated := testEval(test.input)
			testIntegerObject(t, evaluated, test.expected)
		})
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
    fn(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(2);
`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 4)
}

func TestEvalStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\"hello world\"", "hello world"},
		{"\"helloworld\"", "helloworld"},
		{"\"hello\" + \" \" + \"world\"", "hello world"},
		{"\"hello\" * 3", "hellohellohello"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			evaluated := testEval(test.input)
			testStringObject(t, evaluated, test.expected)
		})
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`l("")`, 0},
		{`l("four")`, 4},
		{`l("hello world")`, 11},
		{`l(1)`, "argument to `l` not supported. got=INTEGER"},
		{`l("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`l(["hello", "world", [], ["hello"]])`, 4},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			evaluated := testEval(test.input)

			switch expected := test.expected.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case string:
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Errorf("object is not an Error. got=%T (%+v)", evaluated, evaluated)
					return
				}

				if errObj.Message != expected {
					t.Errorf("errObj.message is not %q. got=%q", expected, errObj.Message)
				}
			}
		})
	}

}

func TestArrayLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{
			input:    `[1]`,
			expected: []interface{}{1},
		},
		{
			input:    `["test"]`,
			expected: []interface{}{"test"},
		},
		{
			input:    `["hello"+"world", 2+2, false]`,
			expected: []interface{}{"helloworld", 4, false},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			evaluated := testEval(tt.input)

			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("evaluated object is not an array. got=%T", evaluated)
			}

			if len(arr.Elements) != len(tt.expected) {
				t.Errorf("arr.Elements doesnt have length %d. got=%d", len(tt.expected), len(arr.Elements))
			}

			for j, el := range arr.Elements {
				switch el.(type) {
				case *object.Integer:
					if !testIntegerObject(t, el, int64(tt.expected[j].(int))) {
						t.Errorf("object is not %d. got=%d", tt.expected[j], el)
					}
				case *object.String:
					if !testStringObject(t, el, tt.expected[j].(string)) {
						t.Errorf("object is not %s. got=%s", tt.expected[j], el)
					}
				case *object.Boolean:
					if !testBooleanObject(t, el, tt.expected[j].(bool)) {
						t.Errorf("object is not %T. got=%T", tt.expected[j], el)
					}
				default:
					if !testBooleanObject(t, el, tt.expected[j].(bool)) {
						t.Errorf("dont know what is happening. want %T (%+v), got=%T (%+v)", tt.expected[j], tt.expected[j], el, el)
					}
				}

			}
		})
	}
}

func TestEvalIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`["test"][0]`,
			"test",
		},
		{
			`["test", "world"][1]`,
			"world",
		},
		{
			`[1, "test", "world"][0]`,
			1,
		},
		{
			`[["hello"][0], "test", "world"][1-1]`,
			"hello",
		},
		{
			`let test = fn(){"hello"}; [[test()][fn(){if true {return 0}}()], "test", "world"][1-1]`,
			"hello",
		},
		{
			`["hello"][2]`,
			errors.New("index 2 out of bounds (array length: 1)"),
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			evaluated := testEval(test.input)

			switch expected := test.expected.(type) {
			case string:
				testStringObject(t, evaluated, expected)
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case bool:
				testBooleanObject(t, evaluated, expected)
			case error:
				testErrorObject(t, evaluated, expected.Error())
			}
		})
	}
}
func testEval(code string) object.Object {
	lexer := lexer.New(code)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not of type String. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. want=%s, got=%s", expected, result.Value)
		return false
	}

	return true
}

func testErrorObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("object is not of type Error. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Message != expected {
		t.Errorf("error has wrong Message. want=%s, got=%s", expected, result.Message)
		return false
	}

	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not of type Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. want=%d, got=%d", expected, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not of type Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. want=%t, got=%t", expected, result.Value)
		return false
	}

	return true
}
