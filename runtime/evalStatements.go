package runtime

import (
	f "github.com/Mstr0A/a0-lang/frontend"
)

// Evaling the Program //
func evalProgram(program f.Program, env *Environment) (RuntimeVal, error) {
	var lastEvaluated RuntimeVal
	var err error

	for _, statement := range program.Body {
		lastEvaluated, err = Evaluate(statement, env)
		if err != nil {
			return nil, err
		}
	}

	return lastEvaluated, nil
}

// Evaluating Variable Declarations //
func evalVarDeclaration(declaration f.VarDeclaration, env *Environment) (RuntimeVal, error) {
	value := declaration.Value
	if value == nil {
		return env.DeclareVar(declaration.Identifier, NadaVal{}, declaration.Constant)
	} else {
		evaluatedValue, err := Evaluate(declaration.Value, env)
		if err != nil {
			return nil, err
		}

		return env.DeclareVar(declaration.Identifier, evaluatedValue, declaration.Constant)
	}
}

// Evaluating Variable Declarations //
func evalFunctionDeclaration(declaration f.FunctionDeclaration, env *Environment) (RuntimeVal, error) {
	fn := UserFunctionValue{
		Name:           declaration.Name,
		Parameters:     declaration.Parameters,
		DeclarationEnv: env,
		Body:           declaration.Body,
	}

	return env.DeclareVar(declaration.Name, fn, true)
}

// Evaluating If Statements //
func evalIfStmt(stmt f.IfStmt, env *Environment) (RuntimeVal, error) {
	condVal, err := Evaluate(stmt.Condition, env)
	if err != nil {
		return nil, err
	}

	boolCond, ok := condVal.(BoolVal)
	if !ok {
		return nil, &InterpretingError{Message: "If statement condition must be a boolean"}
	}

	if boolCond.Value {
		var lastEvaluated RuntimeVal = NadaVal{}
		for _, s := range stmt.Body {
			lastEvaluated, err = Evaluate(s, env)
			if err != nil {
				return nil, err
			}
		}
		return lastEvaluated, nil
	}

	return NadaVal{}, nil
}

// Evaluating While Loops //
func evalWhileStmt(stmt f.WhileStmt, env *Environment) (RuntimeVal, error) {
	var result RuntimeVal = NadaVal{}

	for {
		condVal, err := Evaluate(stmt.Condition, env)
		if err != nil {
			return nil, err
		}

		boolCond, ok := condVal.(BoolVal)
		if !ok {
			return nil, &InterpretingError{Message: "While loop condition must be a boolean"}
		}

		if !boolCond.Value {
			break
		}

		for _, innerStmt := range stmt.Body {
			result, err = Evaluate(innerStmt, env)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

// Evaluating For Loops //
func evalForStmt(stmt f.ForStmt, env *Environment) (RuntimeVal, error) {
	countVal, err := Evaluate(stmt.Condition, env)
	if err != nil {
		return nil, err
	}

	numVal, ok := countVal.(NumberVal)
	if !ok {
		return nil, &InterpretingError{Message: "For loop count must evaluate to a number"}
	}

	var lastEvaluated RuntimeVal
	for i := 0; i < int(numVal.Value); i++ {
		for _, s := range stmt.Body {
			lastEvaluated, err = Evaluate(s, env)
			if err != nil {
				return nil, err
			}
		}
	}

	return lastEvaluated, nil
}

// Evaluating Return Statements //
func evalReturnStmt(stmt f.ReturnStmt, env *Environment) (RuntimeVal, error) {
	val, err := Evaluate(stmt.Value, env)
	if err != nil {
		return nil, err
	}
	return ReturnValue{Value: val}, nil
}
