package runtime

import (
	"fmt"
	"strconv"

	f "github.com/Mstr0A/a0-lang/frontend"
)

// Logical expression eval //
func evalLogicalExpr(logicOp f.LogicalExpr, env *Environment) (RuntimeVal, error) {
	leftSide, err := Evaluate(logicOp.Left, env)
	if err != nil {
		return nil, err
	}

	rightSide, err := Evaluate(logicOp.Right, env)
	if err != nil {
		return nil, err
	}

	switch logicOp.Operator {
	case "and":
		return BoolVal{isTruthy(leftSide) && isTruthy(rightSide)}, nil
	case "or":
		return BoolVal{isTruthy(leftSide) || isTruthy(rightSide)}, nil
	case "not":
		return BoolVal{!isTruthy(leftSide)}, nil
	case "==":
		return BoolVal{deepEqual(leftSide, rightSide)}, nil
	case "!=":
		return BoolVal{!deepEqual(leftSide, rightSide)}, nil
	case "<":
		return BoolVal{lessThan(leftSide, rightSide)}, nil
	case "<=":
		return BoolVal{lessEqual(leftSide, rightSide)}, nil
	case ">":
		return BoolVal{greaterThan(leftSide, rightSide)}, nil
	case ">=":
		return BoolVal{greaterEqual(leftSide, rightSide)}, nil
	default:
		return nil, fmt.Errorf("unknown logical operator: %s", logicOp.Operator)
	}
}

func isTruthy(val RuntimeVal) bool {
	switch v := val.(type) {
	case BoolVal:
		return v.Value
	case NumberVal:
		return v.Value != 0
	case NadaVal:
		return false
	default:
		return val != nil
	}
}

func deepEqual(a, b RuntimeVal) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch a := a.(type) {
	case NumberVal:
		if b, ok := b.(NumberVal); ok {
			return a.Value == b.Value
		}
	case BoolVal:
		if b, ok := b.(BoolVal); ok {
			return a.Value == b.Value
		}
	case NadaVal:
		if _, ok := b.(NadaVal); ok {
			return true
		}
	case ObjectVal:
		if b, ok := b.(ObjectVal); ok {
			return objectsEqual(a.Properties, b.Properties)
		}
	}

	return false
}

func objectsEqual(a, b map[string]RuntimeVal) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valA := range a {
		valB, ok := b[key]
		if !ok || !deepEqual(valA, valB) {
			return false
		}
	}

	return true
}

func lessThan(a, b RuntimeVal) bool {
	if aNum, ok := a.(NumberVal); ok {
		if bNum, ok := b.(NumberVal); ok {
			return aNum.Value < bNum.Value
		}
	}
	return false
}

func lessEqual(a, b RuntimeVal) bool {
	if aNum, ok := a.(NumberVal); ok {
		if bNum, ok := b.(NumberVal); ok {
			return aNum.Value <= bNum.Value
		}
	}
	return false
}

func greaterThan(a, b RuntimeVal) bool {
	if aNum, ok := a.(NumberVal); ok {
		if bNum, ok := b.(NumberVal); ok {
			return aNum.Value > bNum.Value
		}
	}
	return false
}

func greaterEqual(a, b RuntimeVal) bool {
	if aNum, ok := a.(NumberVal); ok {
		if bNum, ok := b.(NumberVal); ok {
			return aNum.Value >= bNum.Value
		}
	}
	return false
}

// Binary expression eval //
func evalBinaryExpr(binOp f.BinaryExpr, env *Environment) (RuntimeVal, error) {
	leftSide, err := Evaluate(binOp.Left, env)
	if err != nil {
		return nil, err
	}

	rightSide, err := Evaluate(binOp.Right, env)
	if err != nil {
		return nil, err
	}

	if leftNum, ok1 := leftSide.(NumberVal); ok1 {
		if rightNum, ok2 := rightSide.(NumberVal); ok2 {
			return evalNumericBinaryExpr(leftNum, rightNum, binOp.Operator)
		}
	}

	return NadaVal{}, nil
}

func evalNumericBinaryExpr(leftSide NumberVal, rightSide NumberVal, operator string) (NumberVal, error) {
	var result float64

	switch operator {
	case "+":
		result = leftSide.Value + rightSide.Value
	case "-":
		result = leftSide.Value - rightSide.Value
	case "*":
		result = leftSide.Value * rightSide.Value
	case "/":
		if rightSide.Value == 0 {
			result = 0
		} else {
			result = leftSide.Value / rightSide.Value
		}
	case "%":
		leftInt := int(leftSide.Value)
		rightInt := int(rightSide.Value)
		result = float64(leftInt % rightInt)
	default:
		errorMessage := fmt.Sprintf("Unknown operator %v", operator)
		return NumberVal{}, &InterpretingError{Message: errorMessage}
	}

	return NumberVal{Value: result}, nil
}

// Unary expression eval //
func evalUnaryExpr(uOp f.UnaryExpr, env *Environment) (RuntimeVal, error) {
	operant, err := Evaluate(uOp.Operant, env)
	if err != nil {
		return nil, err
	}

	if operantNum, ok := operant.(NumberVal); ok {
		return evalNumericUnaryExpr(operantNum, uOp.Operator), nil
	}

	return NadaVal{}, nil
}

func evalNumericUnaryExpr(operant NumberVal, operator string) RuntimeVal {
	var result float64

	switch operator {
	case "-":
		result = -operant.Value
	case "!":
		if result == 0 {
			result = 1
		} else {
			result = 0
		}
	default:
		return operant
	}

	return NumberVal{Value: result}
}

// Evaluating Identifiers //
func evalIdentifier(ident f.Identifier, env *Environment) (RuntimeVal, error) {
	value, err := env.LookupVar(ident.Symbol)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func evalObjectExpr(obj f.ObjectLiteral, env *Environment) (RuntimeVal, error) {
	var err error
	object := ObjectVal{Properties: make(map[string]RuntimeVal)}

	for _, property := range obj.Properties {
		key := property.Key
		value := property.Value

		var runtimeVal RuntimeVal
		if value == nil {
			runtimeVal, err = env.LookupVar(key)
			if err != nil {
				return nil, err
			}
		} else {
			runtimeVal, err = Evaluate(value, env)
			if err != nil {
				return nil, err
			}
		}

		object.Properties[key] = runtimeVal
	}

	return object, err
}

func evalMemberExpr(expr f.MemberExpr, env *Environment) (RuntimeVal, error) {
	objVal, err := Evaluate(expr.Object, env)
	if err != nil {
		return nil, err
	}

	obj, ok := objVal.(ObjectVal)
	if !ok {
		return nil, fmt.Errorf("Attempted to access property of non-object value: %v", objVal)
	}

	var key string

	if expr.Computed {
		propVal, err := Evaluate(expr.Property, env)
		if err != nil {
			return nil, err
		}

		switch k := propVal.(type) {
		case StringVal:
			key = k.Value
		case NumberVal:
			key = strconv.FormatFloat(k.Value, 'f', -1, 64)
		default:
			return nil, fmt.Errorf("Invalid computed property key type: %T", propVal)
		}

	} else {
		ident, ok := expr.Property.(f.Identifier)
		if !ok {
			return nil, fmt.Errorf("Expected Identifier for non-computed property, got %T", expr.Property)
		}
		key = ident.Symbol
	}

	val, exists := obj.Properties[key]
	if !exists {
		return NadaVal{}, nil
	}

	return val, nil
}

// Evaluating Assignment Expression //
func evalAssignmentExpr(node f.AssignmentExpr, env *Environment) (RuntimeVal, error) {
	if node.Assignee.NodeType() != f.IdentifierNode {
		errorMessage := fmt.Sprintf("Invalid left side of assignemt: %v", node.Assignee)
		panic(errorMessage)
	}

	assigneeName := node.Assignee.(f.Identifier).Symbol
	assigneeValue, err := Evaluate(node.Value, env)
	if err != nil {
		return nil, err
	}

	valueToReturn, err := env.AssignVal(assigneeName, assigneeValue)
	if err != nil {
		return nil, err
	}

	return valueToReturn, nil
}

func evalCallExpr(expr f.CallExpr, env *Environment) (RuntimeVal, error) {
	var err error
	args := make([]RuntimeVal, len(expr.Args))
	for i, arg := range expr.Args {
		args[i], err = Evaluate(arg, env)
		if err != nil {
			return nil, err
		}
	}

	fn, err := Evaluate(expr.Caller, env)
	if err != nil {
		return nil, err
	}

	switch callableFn := fn.(type) {
	case NativeFunctionValue:
		result := callableFn.Call(args, env)
		return result, nil

	case UserFunctionValue:
		scope := NewEnvironment(callableFn.DeclarationEnv)

		// Creates the variables for the paremeters list
		if len(callableFn.Parameters) != len(args) {
			errorMessage := fmt.Sprintf("Args do not match amount of parameters in function call for: %s", callableFn.Name)
			return nil, &InterpretingError{Message: errorMessage}
		}
		for i := 0; i < len(callableFn.Parameters); i++ {
			varName := callableFn.Parameters[i]
			scope.DeclareVar(varName, args[i], false)
		}

		var result RuntimeVal = NadaVal{}
		for _, stmt := range callableFn.Body {
			result, err = Evaluate(stmt, scope)
			if err != nil {
				return nil, err
			}

			if ret, ok := result.(ReturnValue); ok {
				return ret.Value, nil
			}
		}

		return NadaVal{}, nil

	default:
		errorMessage := fmt.Sprintf("Cannot call value that is not a function: %v", fn)
		return nil, &InterpretingError{Message: errorMessage}
	}
}
