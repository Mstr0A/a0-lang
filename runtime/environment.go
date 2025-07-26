package runtime

import (
	"fmt"
	"strings"
)

func setupGlobalScope(env *Environment) {
	// Default global variables
	env.DeclareVar("nada", NadaVal{}, true)
	env.DeclareVar("true", BoolVal{Value: true}, true)
	env.DeclareVar("false", BoolVal{Value: false}, true)

	// Defining native global functions
	env.DeclareVar("print", NativeFunctionValue{
		Name: "print",
		Call: func(args []RuntimeVal, env *Environment) RuntimeVal {
			var builder strings.Builder
			for i, arg := range args {
				if i > 0 {
					builder.WriteString("")
				}
				builder.WriteString(arg.String())
			}
			fmt.Println(builder.String())
			return NadaVal{}
		},
	}, true)
}

type Environment struct {
	global    bool
	parent    *Environment
	variables map[string]RuntimeVal
	constants map[string]struct{}
}

func NewEnvironment(parentEnv *Environment) *Environment {
	e := &Environment{
		global:    parentEnv == nil,
		parent:    parentEnv,
		variables: make(map[string]RuntimeVal),
		constants: make(map[string]struct{}),
	}

	if e.global {
		setupGlobalScope(e)
	}

	return e
}

func (env *Environment) setVar(name string, value RuntimeVal) {
	env.variables[name] = value
}

func (env *Environment) DeclareVar(varName string, value RuntimeVal, constant bool) (RuntimeVal, error) {
	_, exists := env.variables[varName]
	if exists {
		errorMessage := fmt.Sprintf("Variable %v already defined, cannot redeclare", varName)
		return nil, &InterpretingError{Message: errorMessage}
	}
	env.setVar(varName, value)

	if constant {
		env.constants[varName] = struct{}{}
	}

	return value, nil
}

func (env *Environment) AssignVal(varName string, value RuntimeVal) (RuntimeVal, error) {
	resolvedEnv, err := env.resolve(varName)
	if err != nil {
		return nil, err
	}

	if _, exists := resolvedEnv.constants[varName]; exists {
		errorMessage := fmt.Sprintf("Cannot assign to constant variable: %v", varName)
		return nil, &InterpretingError{Message: errorMessage}
	}

	resolvedEnv.setVar(varName, value)
	return value, nil
}

func (env *Environment) LookupVar(varName string) (RuntimeVal, error) {
	resolvedEnv, err := env.resolve(varName)
	if err != nil {
		return nil, err
	}
	return resolvedEnv.variables[varName], nil
}

func (env *Environment) resolve(varName string) (*Environment, error) {
	_, exists := env.variables[varName]
	if exists {
		return env, nil
	}
	if env.parent == nil {
		errorMessage := fmt.Sprintf("Variable %v does not exist", varName)
		panic(errorMessage)
	}
	return env.parent.resolve(varName)
}
