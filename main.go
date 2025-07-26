package main

import (
	"flag"
	"fmt"
	"os"

	f "github.com/Mstr0A/a0-lang/frontend"
	r "github.com/Mstr0A/a0-lang/runtime"
)

///////////////////
// Main Function //
///////////////////

// handles statement nodes
func printStmt(node f.Stmt, indent string, isLast bool) {
	// tree symbols
	branch, nextIndent := "├── ", indent+"│   "
	if isLast {
		branch, nextIndent = "└── ", indent+"    "
	}

	switch n := node.(type) {
	case f.Program:
		fmt.Println(indent + branch + "Program")
		for i, stmt := range n.Body {
			printStmt(stmt, nextIndent, i == len(n.Body)-1)
		}

	case f.VarDeclaration:
		fmt.Printf("%s%sVarDeclaration: Name: %s | Constant: %t\n",
			indent, branch,
			n.Identifier,
			n.Constant,
		)
		if n.Value != nil {
			printExpr(n.Value, nextIndent, true)
		}

	case f.AssignmentExpr:
		fmt.Printf("%s%sAssignmentExpr\n", indent, branch)
		printExpr(n.Assignee, nextIndent, false)
		printExpr(n.Value, nextIndent, true)

	case f.FunctionDeclaration:
		fmt.Printf("%s%sFunctionDeclaration\n", indent, branch)

		// Name
		fmt.Printf("%s%sName: %s\n",
			nextIndent, "└── ", n.Name,
		)

		// Parameters
		fmt.Printf("%s├── Parameters\n", nextIndent)
		for i, param := range n.Parameters {
			pBranch := "│   ├── "
			if i == len(n.Parameters)-1 {
				pBranch = "│   └── "
			}
			fmt.Printf("%s%sIdentifier (%s)\n",
				nextIndent, pBranch, param,
			)
		}

		// Body
		bodyIndent := nextIndent + "    "
		fmt.Printf("%s└── Body\n", nextIndent)
		for i, stmt := range n.Body {
			printStmt(stmt, bodyIndent, i == len(n.Body)-1)
		}

	case f.CallExpr:
		// Treat bare CallExpr as a statement
		fmt.Printf("%s%sCallExpr\n", indent, branch)
		printExpr(n.Caller, nextIndent, false)
		for i, arg := range n.Args {
			printExpr(arg, nextIndent, i == len(n.Args)-1)
		}

	case f.ObjectLiteral:
		fmt.Printf("%s%sObjectLiteral\n", indent, branch)
		for i, prop := range n.Properties {
			propBranch := "├── "
			if i == len(n.Properties)-1 {
				propBranch = "└── "
			}
			fmt.Printf("%s%sProperty: Key: %s\n",
				nextIndent, propBranch, prop.Key,
			)
			// property value is an Expr
			printExpr(prop.Value, nextIndent+"│   ", i == len(n.Properties)-1)
		}

	default:
		fmt.Printf("%s%sUnknown stmt node of type %T\n", indent, branch, n)
	}
}

// handles expression nodes
func printExpr(node f.Expr, indent string, isLast bool) {
	branch, nextIndent := "├── ", indent+"│   "
	if isLast {
		branch, nextIndent = "└── ", indent+"    "
	}

	switch n := node.(type) {
	case f.Identifier:
		fmt.Printf("%s%sIdentifier (%s)\n", indent, branch, n.Symbol)

	case f.NumericLiteral:
		fmt.Printf("%s%sNumericLiteral (%f)\n", indent, branch, n.Value)

	case f.BinaryExpr:
		fmt.Printf("%s%sBinaryExpr (Operator: %s)\n", indent, branch, n.Operator)
		printExpr(n.Left, nextIndent, false)
		printExpr(n.Right, nextIndent, true)

	case f.LogicalExpr:
		fmt.Printf("%s%sLogicalExpr (Operator: %s)\n", indent, branch, n.Operator)
		printExpr(n.Left, nextIndent, false)
		printExpr(n.Right, nextIndent, true)

	case f.UnaryExpr:
		fmt.Printf("%s%sUnaryExpr (Operator: %s)\n", indent, branch, n.Operator)
		printExpr(n.Operant, nextIndent, true)

	case f.CallExpr:
		fmt.Printf("%s%sCallExpr\n", indent, branch)
		printExpr(n.Caller, nextIndent, false)
		for i, arg := range n.Args {
			printExpr(arg, nextIndent, i == len(n.Args)-1)
		}

	default:
		fmt.Printf("%s%sUnknown expr node of type %T\n", indent, branch, n)
	}
}

func printAST(root f.Stmt) {
	printStmt(root, "", true)
}

func main() {
	///////////
	// Flags //
	///////////

	showTokens := flag.Bool("tokens", false, "Print the token list")
	showAst := flag.Bool("ast", false, "Print the AST")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: yourlang [options] <file>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	//////////
	// File //
	//////////

	filePath := flag.Args()[0]
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	///////////
	// Lexer //
	///////////

	lexer := f.NewLexer(file)
	tokenList, err := lexer.Lex()
	if err != nil {
		fmt.Println(err)
		return
	}
	if *showTokens {
		fmt.Println("Tokens:")
		for _, tok := range tokenList {
			fmt.Println(tok)
		}
	}

	//////////////////////////
	// Parser & Interpreter //
	//////////////////////////

	parser := f.NewParser(tokenList)
	program, err := parser.ProduceAst()
	if err != nil {
		fmt.Println(err)
		return
	}
	if *showAst {
		fmt.Println("AST:")
		printAST(program)
	}

	if *showAst || *showTokens {
		return
	}

	env := r.NewEnvironment(nil)
	_, err = r.Evaluate(program, env)
	if err != nil {
		fmt.Println(err)
		return
	}
}
