package eval

import (
	"fmt"
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
		return evalProgram(node)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.ExpressionStatement:
		return Eval(node.Value)

	case *ast.ReturnStatement:
		value := Eval(node.Value)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		return evalPrefixExpression(node)

	case *ast.InfixExpression:
		return evalInfixExpression(node)

	case *ast.IfExpression:
		return evalIfExpression(node)

	}

	return NULL
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object = NULL
	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.Error:
			return result
		case *object.ReturnValue:
			// Unwrap return value
			return result.Value
		}
	}
	return result
}

func evalBlockStatement(blockStatement *ast.BlockStatement) object.Object {
	var result object.Object = NULL
	for _, statement := range blockStatement.Statements {
		result = Eval(statement)

		if result.Type() == object.ERROR_OBJECT || result.Type() == object.RETURN_VALUE_OBJECT {
			return result
		}
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
	if isError(right) {
		return right
	}

	switch expression.Operator {
	case "-":
		return evalPrefixMinusOperator(right)
	case "!":
		return evalBangOperator(right)
	default:
		return NULL
	}
}

func evalPrefixMinusOperator(right object.Object) object.Object {
	if integer, ok := right.(*object.Integer); ok {
		return &object.Integer{
			Value: -integer.Value,
		}
	}
	return createError(
		"unknown operation -%s", right.Type(),
	)
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

func evalInfixExpression(expression *ast.InfixExpression) object.Object {
	operator := expression.Operator
	left := Eval(expression.Left)
	if isError(left) {
		return left
	}

	right := Eval(expression.Right)
	if isError(right) {
		return right
	}

	if left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT {
		return evalIntegerInfixExpression(operator, left, right)
	}

	if left.Type() != right.Type() {
		return createError(
			"type mismatch %s %s %s",
			left.Type(), operator, right.Type(),
		)
	}

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return createError(
			"unknown operation %s %s %s",
			left.Type(), operator, right.Type(),
		)
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return createError(
			"unknown integer operation %s %s %s",
			left.Type(), operator, right.Type(),
		)
	}
}

func evalIfExpression(expression *ast.IfExpression) object.Object {
	condition := Eval(expression.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(expression.Consequence)
	}

	if expression.Alternative != nil {
		return Eval(expression.Alternative)
	}

	return nil
}

func isTruthy(value object.Object) bool {
	switch value {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func createError(format string, a ...any) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(value object.Object) bool {
	return value != nil && value.Type() == object.ERROR_OBJECT
}
