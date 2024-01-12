package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJECT ObjectType = "INTEGER"
	BOOLEAN_OBJECT ObjectType = "BOOLEAN"
	NULL_OBJECT    ObjectType = "NULL"
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
