package eval

import (
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Value)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		return evalPrefixExpression(node)

	}

	return NULL
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object = NULL
	for _, statement := range statements {
		result = Eval(statement)
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(expression *ast.PrefixExpression) object.Object {
	right := Eval(expression.Right)

	switch expression.Operator {
	case "-":
		return evalMinusOperator(right)
	case "!":
		return evalBangOperator(right)
	default:
		return NULL
	}
}

func evalMinusOperator(right object.Object) object.Object {
	if integer, ok := right.(*object.Integer); ok {
		return &object.Integer{
			Value: -integer.Value,
		}
	}
	// TODO: Handle unexpected usage of "-"
	return NULL
}

func evalBangOperator(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}
