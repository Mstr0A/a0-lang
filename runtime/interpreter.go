package runtime

import (
	"fmt"

	f "github.com/Mstr0A/a0-lang/frontend"
)

type InterpretingError struct {
	Message string
}

func (e *InterpretingError) Error() string {
	return fmt.Sprintf("Interpretation Error: %s", e.Message)
}

// Main Eval //
func Evaluate(astNode f.Stmt, env *Environment) (RuntimeVal, error) {
	switch castedNode := astNode.(type) {
	case f.Program:
		return evalProgram(castedNode, env)
	case f.NumericLiteral:
		return NumberVal{Value: castedNode.Value}, nil
	case f.StringLiteral:
		return StringVal{Value: castedNode.Value}, nil
	case f.Identifier:
		return evalIdentifier(castedNode, env)
	case f.ObjectLiteral:
		return evalObjectExpr(castedNode, env)
	case f.MemberExpr:
		return evalMemberExpr(castedNode, env)
	case f.BinaryExpr:
		return evalBinaryExpr(castedNode, env)
	case f.UnaryExpr:
		return evalUnaryExpr(castedNode, env)
	case f.VarDeclaration:
		return evalVarDeclaration(castedNode, env)
	case f.FunctionDeclaration:
		return evalFunctionDeclaration(castedNode, env)
	case f.AssignmentExpr:
		return evalAssignmentExpr(castedNode, env)
	case f.CallExpr:
		return evalCallExpr(castedNode, env)
	case f.LogicalExpr:
		return evalLogicalExpr(castedNode, env)
	case f.IfStmt:
		return evalIfStmt(castedNode, env)
	case f.WhileStmt:
		return evalWhileStmt(castedNode, env)
	case f.ForStmt:
		return evalForStmt(castedNode, env)
	case f.ReturnStmt:
		return evalReturnStmt(castedNode, env)
	default:
		errorMessage := fmt.Sprintf("AST Node has not been added for interpretation: %v", castedNode)
		err := &InterpretingError{Message: errorMessage}
		return nil, err
	}
}
