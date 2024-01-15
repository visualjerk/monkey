package object

import (
	"fmt"
)

type ObjectType string

const (
	INTEGER_OBJECT      ObjectType = "INTEGER"
	BOOLEAN_OBJECT      ObjectType = "BOOLEAN"
	NULL_OBJECT         ObjectType = "NULL"
	RETURN_VALUE_OBJECT ObjectType = "RETURN_VALUE"
	ERROR_OBJECT        ObjectType = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (integer *Integer) Type() ObjectType { return INTEGER_OBJECT }
func (integer *Integer) Inspect() string  { return fmt.Sprintf("%d", integer.Value) }

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Type() ObjectType { return BOOLEAN_OBJECT }
func (boolean *Boolean) Inspect() string  { return fmt.Sprintf("%t", boolean.Value) }

type Null struct{}

func (null *Null) Type() ObjectType { return NULL_OBJECT }
func (null *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (returnValue *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJECT }
func (returnValue *ReturnValue) Inspect() string {
	return returnValue.Value.Inspect()
}

type Error struct {
	Message string
}

func (errorObject *Error) Type() ObjectType { return ERROR_OBJECT }
func (errorObject *Error) Inspect() string {
	return "Error: " + errorObject.Message
}
