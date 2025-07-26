package frontend

///////////////
// NodeTypes //
///////////////

// NodeType //
type NodeType string

const (
	// Statements
	ProgramNode             NodeType = "Program"
	VarDeclarationNode      NodeType = "VarDeclaration"
	FunctionDeclarationNode NodeType = "FunctionDeclaration"

	// Expressions
	AssignmentExpressionNode NodeType = "AssignmentExpr"
	MemberExpressionNode     NodeType = "MemberExpr"
	CallExpressionNode       NodeType = "CallExpr"

	// Literals
	ObjectLiteralNode     NodeType = "Object"
	PropertyNode          NodeType = "Property"
	NumericLiteralNode    NodeType = "NumericLiteral"
	StringLiteralNode     NodeType = "StringLiteral"
	IdentifierNode        NodeType = "Identifier"
	LogicalExpressionNode NodeType = "LogicalExpr"
	BinaryExpressionNode  NodeType = "BinaryExpr"
	UnaryExpressionNode   NodeType = "UnaryExpr"

	// Keywords
	IfStmtNode     NodeType = "IfStmt"
	WhileStmtNode  NodeType = "WhileStmt"
	ForStmtNode    NodeType = "ForStmt"
	ReturnStmtNode NodeType = "ReturnStmt"
)

// Base Types //
type Stmt interface {
	NodeType() NodeType
}

type Expr interface {
	Stmt // Expr embeds Stmt so now Expr is also as Stmt
}

// Statements //

type Program struct {
	Body []Stmt
}

func (p Program) NodeType() NodeType {
	return ProgramNode
}

type VarDeclaration struct {
	Constant   bool
	Identifier string
	Value      Expr
}

func (v VarDeclaration) NodeType() NodeType {
	return VarDeclarationNode
}

type FunctionDeclaration struct {
	Name       string
	Parameters []string
	Body       []Stmt
}

func (f FunctionDeclaration) NodeType() NodeType {
	return FunctionDeclarationNode
}

type IfStmt struct {
	Condition Expr
	Body      []Stmt
}

func (i IfStmt) NodeType() NodeType {
	return IfStmtNode
}

type WhileStmt struct {
	Condition Expr
	Body      []Stmt
}

func (w WhileStmt) NodeType() NodeType {
	return WhileStmtNode
}

type ForStmt struct {
	Condition Expr
	Body      []Stmt
}

func (f ForStmt) NodeType() NodeType {
	return ForStmtNode
}

type ReturnStmt struct {
	Value Expr
}

func (r ReturnStmt) NodeType() NodeType {
	return ReturnStmtNode
}

// Expressions //

type AssignmentExpr struct {
	Assignee Expr
	Value    Expr
}

func (a AssignmentExpr) NodeType() NodeType {
	return AssignmentExpressionNode
}

type CallExpr struct {
	Args   []Expr
	Caller Expr
}

func (c CallExpr) NodeType() NodeType {
	return CallExpressionNode
}

type MemberExpr struct {
	Object   Expr
	Property Expr
	Computed bool
}

func (m MemberExpr) NodeType() NodeType {
	return MemberExpressionNode
}

// Literals //
type LogicalExpr struct {
	Left     Expr
	Right    Expr
	Operator string
}

func (l LogicalExpr) NodeType() NodeType {
	return LogicalExpressionNode
}

type BinaryExpr struct {
	Left     Expr
	Right    Expr
	Operator string
}

func (b BinaryExpr) NodeType() NodeType {
	return BinaryExpressionNode
}

type UnaryExpr struct {
	Operant  Expr
	Operator string
}

func (b UnaryExpr) NodeType() NodeType {
	return UnaryExpressionNode
}

type NumericLiteral struct {
	Value float64
}

func (n NumericLiteral) NodeType() NodeType {
	return NumericLiteralNode
}

type StringLiteral struct {
	Value string
}

func (s StringLiteral) NodeType() NodeType {
	return StringLiteralNode
}

type Identifier struct {
	Symbol string
}

func (i Identifier) NodeType() NodeType {
	return IdentifierNode
}

type Property struct {
	Key   string
	Value Expr
}

func (p Property) NodeType() NodeType {
	return PropertyNode
}

type ObjectLiteral struct {
	Properties []Property
}

func (o ObjectLiteral) NodeType() NodeType {
	return ObjectLiteralNode
}
