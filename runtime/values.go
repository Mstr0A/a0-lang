package runtime

import (
	"fmt"
	"strconv"

	f "github.com/Mstr0A/a0-lang/frontend"
)

////////////////
// ValueTypes //
////////////////

// ValueType //
type ValueType string

const (
	NumberType         ValueType = "Number"
	StringType         ValueType = "String"
	NadaType           ValueType = "Nada"
	BoolType           ValueType = "Bool"
	ObjectType         ValueType = "Object"
	NativeFunctionType ValueType = "NativeFunction"
	UserFunctionType   ValueType = "UserFunction"
	ReturnSignalType   ValueType = "ReturnSignal"
)

// Runtime Value //
type RuntimeVal interface {
	ValueType() ValueType
	String() string
}

// Number Value //
type NumberVal struct {
	Value float64
}

func (n NumberVal) ValueType() ValueType {
	return NumberType
}

func (n NumberVal) String() string {
	return strconv.FormatFloat(n.Value, 'f', -1, 64)
}

// Number Value //
type StringVal struct {
	Value string
}

func (s StringVal) ValueType() ValueType {
	return NumberType
}

func (s StringVal) String() string {
	return s.Value
}

// Nada Value //
type NadaVal struct{}

func (n NadaVal) ValueType() ValueType {
	return NadaType
}

func (n NadaVal) String() string {
	return "nada"
}

// Bool Value //
type BoolVal struct {
	Value bool
}

func (b BoolVal) ValueType() ValueType {
	return BoolType
}

func (b BoolVal) String() string {
	return strconv.FormatBool(b.Value)
}

// Object Value //
type ObjectVal struct {
	Properties map[string]RuntimeVal
	ObjectName string
}

func (o ObjectVal) ValueType() ValueType {
	return ObjectType
}

func (o ObjectVal) String() string {
	return fmt.Sprintf("User Object (%s)", o.ObjectName)
}

// Function Value //
type FunctionCall func(args []RuntimeVal, env *Environment) RuntimeVal

type NativeFunctionValue struct {
	Call FunctionCall
	Name string
}

func (nf NativeFunctionValue) ValueType() ValueType {
	return NativeFunctionType
}

func (nf NativeFunctionValue) String() string {
	return fmt.Sprintf("Native Function (%s)", nf.Name)
}

type UserFunctionValue struct {
	Name           string
	Parameters     []string
	DeclarationEnv *Environment
	Body           []f.Stmt
}

func (uf UserFunctionValue) ValueType() ValueType {
	return UserFunctionType
}

func (uf UserFunctionValue) String() string {
	return fmt.Sprintf("User Function (%s)", uf.Name)
}

// Return Value //
type ReturnValue struct {
	Value RuntimeVal
}

func (r ReturnValue) ValueType() ValueType {
	return ReturnSignalType
}

func (r ReturnValue) String() string {
	if r.Value == nil {
		return "return <nil>"
	}
	return fmt.Sprintf("return %v", r.Value)
}

// I don't know why this specifically needs the error interface
func (r ReturnValue) Error() string {
	if r.Value == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", r.Value)
}
