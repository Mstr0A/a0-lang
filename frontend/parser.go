package frontend

import (
	"fmt"
	"strconv"
)

///////////////////
// Parsing Error //
///////////////////

type ParsingError struct {
	Message string
	Pos     Position
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("Parse Error at (%d, %d): %s", e.Pos.line, e.Pos.column, e.Message)
}

////////////
// Parser //
////////////

type Parser struct {
	tokens       []TokenItem
	tokenIndex   int
	currentToken TokenItem
}

func TokenToFloat(token TokenItem) float64 {
	stringValue := token.value
	floatValue, _ := strconv.ParseFloat(stringValue, 64)
	return floatValue
}

func NewParser(tokens []TokenItem) *Parser {
	p := Parser{
		tokens:     tokens,
		tokenIndex: -1,
	}
	p.advance()
	return &p
}

func (p *Parser) eat() TokenItem {
	prev := p.currentToken
	p.advance()
	return prev
}

func (p *Parser) expect(expectedType Token, errMsg string) (TokenItem, error) {
	token := p.eat()
	if token.tokenType != expectedType {
		return TokenItem{}, &ParsingError{
			Message: fmt.Sprintf("Parsing Error: %s", errMsg),
			Pos:     token.pos,
		}
	}
	return token, nil
}

func (p *Parser) ProduceAst() (Program, error) {
	program := Program{}

	for {
		stmt, err := p.parseStmt()
		if err != nil {
			return Program{}, err
		}
		program.Body = append(program.Body, stmt)
		if p.currentToken.tokenType == EOF {
			break
		}
	}

	return program, nil
}

func (p *Parser) advance() {
	p.tokenIndex++
	if p.tokenIndex < len(p.tokens) {
		p.currentToken = p.tokens[p.tokenIndex]
	}
}

func (p *Parser) parseStmt() (Stmt, error) {
	switch p.currentToken.tokenType {
	case VAR, CONST:
		return p.parseVarDeclaration()
	case FUN:
		return p.parseFunctionDeclaration()
	case IF:
		return p.parseIfStmt()
	case WHILE:
		return p.parseWhileStmt()
	case FOR:
		return p.parseForStmt()
	case RETURN:
		return p.parseReturnStmt()
	default:
		return p.parseExpr()
	}
}

// Parsing Expressions
func (p *Parser) parseExpr() (Expr, error) {
	return p.parseAssignmentExpr()
}

func (p *Parser) parseAdditive() (Expr, error) {
	left, err := p.parseMulti()
	if err != nil {
		return nil, err
	}

	for p.currentToken.tokenType == ADD || p.currentToken.tokenType == SUB {
		operator := p.eat().value
		right, err := p.parseMulti()
		if err != nil {
			return nil, err
		}
		left = BinaryExpr{
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}
	return left, nil
}

func (p *Parser) parseMulti() (Expr, error) {
	left, err := p.parseCallMemberExpr()
	if err != nil {
		return nil, err
	}

	for p.currentToken.tokenType == MUL || p.currentToken.tokenType == DIV || p.currentToken.tokenType == MOD {
		operator := p.eat().value
		right, err := p.parseCallMemberExpr()
		if err != nil {
			return nil, err
		}
		left = BinaryExpr{
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}
	return left, nil
}

func (p *Parser) parsePrimary() (Expr, error) {
	tokenType := p.currentToken.tokenType

	if tokenType == NOT {
		p.eat()
		expr, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}

		return UnaryExpr{
			Operator: "!",
			Operant:  expr,
		}, nil
	}

	switch tokenType {
	case IDENT:
		token := p.eat()
		return Identifier{Symbol: token.value}, nil
	case INT, FLOAT:
		token := p.eat()
		return NumericLiteral{Value: TokenToFloat(token)}, nil
	case STRING:
		token := p.eat()
		return StringLiteral{Value: token.value}, nil
	case OPENPAREN:
		p.eat() // Skip '('
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(CLOSEPAREN, "Expected closing parenthesis")
		if err != nil {
			return nil, err
		}

		return value, nil
	case OPENCURLY:
		return p.parseObjectExpr()
	case EOF, CLOSEPAREN, CLOSECURLY, COMMA:
		return nil, &ParsingError{
			Message: "Expected an expression or value but found none",
			Pos:     p.currentToken.pos,
		}
	case ILLEGAL:
		return nil, &ParsingError{
			Message: fmt.Sprintf("Illegal token passed \"%v\"", p.currentToken.value),
			Pos:     p.currentToken.pos,
		}
	default:
		return nil, &ParsingError{
			Message: fmt.Sprintf("Unrecognized Primary Token (Type: %s, Value: %s)", TokensList[p.currentToken.tokenType], p.currentToken.value),
			Pos:     p.currentToken.pos,
		}
	}
}

// Parsing Variable Declarations
func (p *Parser) parseVarDeclaration() (Stmt, error) {
	isConstant := p.currentToken.tokenType == CONST
	p.eat()

	identifier, err := p.expect(IDENT, "Expected identifier name after var | const keyword")
	if err != nil {
		return nil, err
	}

	if p.currentToken.tokenType != EQUALS {
		if isConstant {
			return nil, &ParsingError{
				Message: "Uninitialized constant",
				Pos:     p.currentToken.pos,
			}
		}
		return VarDeclaration{
			Constant:   isConstant,
			Identifier: identifier.value,
			Value:      nil,
		}, nil
	}

	p.eat()
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return VarDeclaration{
		Constant:   isConstant,
		Identifier: identifier.value,
		Value:      value,
	}, nil
}

func (p *Parser) parseAssignmentExpr() (Expr, error) {
	expr, err := p.parseLogicalExpr()
	if err != nil {
		return nil, err
	}

	if p.currentToken.tokenType == EQUALS {
		p.eat() // consume the '=' token

		value, err := p.parseAssignmentExpr()
		if err != nil {
			return nil, err
		}

		return AssignmentExpr{
			Assignee: expr,
			Value:    value,
		}, nil
	}

	return expr, nil // If no assignment, return the expression as-is
}

// Parsing Objects
func (p *Parser) parseObjectExpr() (Expr, error) {
	if p.currentToken.tokenType != OPENCURLY {
		return p.parseAdditive()
	}
	p.eat() // Skip the open brace
	properties := []Property{}

	for p.currentToken.tokenType != EOF && p.currentToken.tokenType != CLOSECURLY {
		object, err := p.expect(IDENT, "Object missing identifier")
		if err != nil {
			return nil, err
		}
		key := object.value

		// Handle shorthand properties { foo }
		if p.currentToken.tokenType == COMMA || p.currentToken.tokenType == CLOSECURLY {
			properties = append(properties, Property{Key: key, Value: nil})
			if p.currentToken.tokenType == COMMA {
				p.eat() // Skip comma
			}
			continue
		}

		// Expect colon for normal key-value pair
		_, err = p.expect(COLON, "Missing colon after identifier")
		if err != nil {
			return nil, err
		}

		// Handle nested objects { key: { ... } }
		var value Expr
		if p.currentToken.tokenType == OPENCURLY {
			value, err = p.parseObjectExpr() // Recursively parse nested object
			if err != nil {
				return nil, err
			}
		} else {
			value, err = p.parseExpr() // Parse other value types
			if err != nil {
				return nil, err
			}
		}

		properties = append(properties, Property{Key: key, Value: value})

		// Expect comma or closing brace
		if p.currentToken.tokenType != CLOSECURLY {
			_, err = p.expect(COMMA, "Expected comma or closing brace after property")
			if err != nil {
				return nil, err
			}
		}
	}

	_, err := p.expect(CLOSECURLY, "Object literal missing closing brace")
	if err != nil {
		return nil, err
	}

	return ObjectLiteral{Properties: properties}, nil
}

// Parsing Member Calls
func (p *Parser) parseCallMemberExpr() (Expr, error) {
	member, err := p.parseMemberExpr()
	if err != nil {
		return nil, err
	}

	if p.currentToken.tokenType == OPENPAREN {
		return p.parseCallExpr(member)
	}

	return member, nil
}

// Parsing Calls
func (p *Parser) parseCallExpr(caller Expr) (Expr, error) {
	arguments, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	callExpr := CallExpr{Caller: caller, Args: arguments}

	if p.currentToken.tokenType == OPENPAREN {
		return p.parseCallExpr(callExpr)
	}

	return callExpr, nil
}

func (p *Parser) parseArguments() ([]Expr, error) {
	args := []Expr{}

	_, err := p.expect(OPENPAREN, "Expected \"(\"")
	if err != nil {
		return nil, err
	}

	if p.currentToken.tokenType == CLOSEPAREN {
		p.eat()
		return args, nil
	}

	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.currentToken.tokenType != COMMA {
			break
		}
		p.eat() // Skip comma
	}

	_, err = p.expect(CLOSEPAREN, "Expected \")\"")
	if err != nil {
		return nil, err
	}

	return args, nil
}

func (p *Parser) parseMemberExpr() (Expr, error) {
	object, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for p.currentToken.tokenType == DOT || p.currentToken.tokenType == OPENBRACKET {
		operator := p.eat()
		var property Expr
		var computed bool

		// Non-computed values (dot values obj.expr)
		if operator.tokenType == DOT {
			computed = false
			property, err = p.parsePrimary()
			if err != nil {
				return nil, err
			}

			if property.NodeType() != IdentifierNode {
				return nil, &ParsingError{
					Pos:     p.currentToken.pos,
					Message: "Cannot use dot operator without having an identifier after it",
				}
			}
		} else { // this allows chaining
			computed = true
			property, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
			p.expect(CLOSEBRACKET, "Expected \"]\"")
		}

		object = MemberExpr{
			Object:   object,
			Property: property,
			Computed: computed,
		}
	}

	return object, nil
}

// Parsing Function Declarations
func (p *Parser) parseFunctionDeclaration() (Stmt, error) {
	p.eat() // Skip the fun keyword

	name, err := p.expect(IDENT, "Expected function name after keyword \"fun\"")
	if err != nil {
		return nil, err
	}

	args, err := p.parseArguments()
	if err != nil {
		return nil, err
	}

	params := []string{}
	for _, arg := range args {
		if arg.NodeType() != IdentifierNode {
			return nil, &ParsingError{
				Message: "Expected parameter inside function declaration",
				Pos:     name.pos,
			}
		}
		params = append(params, arg.(Identifier).Symbol)
	}

	_, err = p.expect(OPENCURLY, "Expected \"{\"")
	if err != nil {
		return nil, err
	}

	body := []Stmt{}
	for p.currentToken.tokenType != EOF && p.currentToken.tokenType != CLOSECURLY {
		statement, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		body = append(body, statement)
	}

	_, err = p.expect(CLOSECURLY, "Expected \"}\"")
	if err != nil {
		return nil, err
	}

	return FunctionDeclaration{
		Name:       name.value,
		Parameters: params,
		Body:       body,
	}, nil
}

func (p *Parser) parseLogicalExpr() (Expr, error) {
	left, err := p.parseEqualityExpr()
	if err != nil {
		return nil, err
	}

	for p.currentToken.tokenType == AND || p.currentToken.tokenType == OR {
		operator := p.eat().value

		right, err := p.parseEqualityExpr()
		if err != nil {
			return nil, err
		}

		left = LogicalExpr{
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}

	return left, nil
}

func (p *Parser) parseEqualityExpr() (Expr, error) {
	left, err := p.parseRelationalExpr()
	if err != nil {
		return nil, err
	}

	for p.currentToken.tokenType == DE || p.currentToken.tokenType == NE {
		operator := p.eat().value

		right, err := p.parseRelationalExpr()
		if err != nil {
			return nil, err
		}

		left = LogicalExpr{
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}

	return left, nil
}

func (p *Parser) parseRelationalExpr() (Stmt, error) {
	left, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}

	for p.currentToken.tokenType == LT || p.currentToken.tokenType == GT ||
		p.currentToken.tokenType == LTE || p.currentToken.tokenType == GTE {

		operator := p.eat().value

		right, err := p.parseAdditive()
		if err != nil {
			return nil, err
		}

		left = LogicalExpr{
			Left:     left,
			Right:    right,
			Operator: operator,
		}
	}

	return left, nil
}

// Parsing if statements
func (p *Parser) parseIfStmt() (Stmt, error) {
	_, err := p.expect(IF, "Expected 'if' keyword")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(OPENPAREN, "Expected '(' after 'if'")
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.expect(CLOSEPAREN, "Expected ')' after if condition")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(OPENCURLY, "Expected '{' to begin if statement body")
	if err != nil {
		return nil, err
	}

	body := []Stmt{}
	for p.currentToken.tokenType != EOF && p.currentToken.tokenType != CLOSECURLY {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	_, err = p.expect(CLOSECURLY, "Expected '}' to close if statement body")
	if err != nil {
		return nil, err
	}

	return IfStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

// Parsing while loops
func (p *Parser) parseWhileStmt() (Stmt, error) {
	_, err := p.expect(WHILE, "Expected 'while' keyword")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(OPENPAREN, "Expected '(' after 'while'")
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.expect(CLOSEPAREN, "Expected ')' after while condition")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(OPENCURLY, "Expected '{' to begin while loop body")
	if err != nil {
		return nil, err
	}

	body := []Stmt{}
	for p.currentToken.tokenType != EOF && p.currentToken.tokenType != CLOSECURLY {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	_, err = p.expect(CLOSECURLY, "Expected '}' to close while loop body")
	if err != nil {
		return nil, err
	}

	return WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

// Parsing for loops
func (p *Parser) parseForStmt() (Stmt, error) {
	_, err := p.expect(FOR, "Expected 'for' keyword")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(OPENPAREN, "Expected '(' after 'for'")
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	_, err = p.expect(CLOSEPAREN, "Expected ')' after for condition")
	if err != nil {
		return nil, err
	}

	_, err = p.expect(OPENCURLY, "Expected '{' to begin for loop body")
	if err != nil {
		return nil, err
	}

	body := []Stmt{}
	for p.currentToken.tokenType != EOF && p.currentToken.tokenType != CLOSECURLY {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		body = append(body, stmt)
	}

	_, err = p.expect(CLOSECURLY, "Expected '}' to close while loop body")
	if err != nil {
		return nil, err
	}

	return ForStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

// Parsing Return Statements
func (p *Parser) parseReturnStmt() (Stmt, error) {
	_, err := p.expect(RETURN, "Expected 'return' keyword")
	if err != nil {
		return nil, err
	}

	// If next token is close curly or EOF, no return value
	if p.currentToken.tokenType == CLOSECURLY || p.currentToken.tokenType == EOF {
		return ReturnStmt{Value: nil}, nil
	}

	// Otherwise parse expression for return value
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return ReturnStmt{Value: expr}, nil
}
