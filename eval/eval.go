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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.LetStatement:
		return evalLetStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Value, env)

	case *ast.ReturnStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)

	case *ast.InfixExpression:
		return evalInfixExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)

	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL
	for _, statement := range program.Statements {
		result = Eval(statement, env)

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

func evalBlockStatement(blockStatement *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	innerEnv := object.NewEnclosedEnvironment(env)

	for _, statement := range blockStatement.Statements {
		result = Eval(statement, innerEnv)

		if result != nil && (result.Type() == object.ERROR_OBJECT || result.Type() == object.RETURN_VALUE_OBJECT) {
			return result
		}
	}
	return result
}

func evalLetStatement(letStatement *ast.LetStatement, env *object.Environment) object.Object {
	value := Eval(letStatement.Value, env)
	if isError(value) {
		return value
	}

	env.Set(letStatement.Name.Value, value)

	return nil
}

func evalIdentifier(identifier *ast.Identifier, env *object.Environment) object.Object {
	value, ok := env.Get(identifier.Value)
	if !ok {
		return newError("identifier not found %s", identifier.Value)
	}
	return value
}

func evalPrefixExpression(expression *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(expression.Right, env)
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
	return newError(
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

func evalInfixExpression(expression *ast.InfixExpression, env *object.Environment) object.Object {
	operator := expression.Operator
	left := Eval(expression.Left, env)
	if isError(left) {
		return left
	}

	right := Eval(expression.Right, env)
	if isError(right) {
		return right
	}

	if left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT {
		return evalIntegerInfixExpression(operator, left, right)
	}

	if left.Type() != right.Type() {
		return newError(
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
		return newError(
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
		return newError(
			"unknown integer operation %s %s %s",
			left.Type(), operator, right.Type(),
		)
	}
}

func evalIfExpression(expression *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(expression.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(expression.Consequence, env)
	}

	if expression.Alternative != nil {
		return Eval(expression.Alternative, env)
	}

	return nil
}

func evalFunctionLiteral(expression *ast.FunctionLiteral, env *object.Environment) object.Object {
	function := &object.Function{
		Parameters: expression.Parameters,
		Body:       expression.Body,
		Env:        env,
	}

	return function
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
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

func newError(format string, a ...any) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(value object.Object) bool {
	return value != nil && value.Type() == object.ERROR_OBJECT
}
