package frontend

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

////////////
// Tokens //
////////////

type Token int

const (
	// Data Types
	EOF Token = iota
	ILLEGAL
	IDENT
	RETURN
	INT
	FLOAT
	STRING
	CHAR

	// Reserved Characters
	VAR
	CONST
	OPENCURLY
	CLOSECURLY
	OPENPAREN
	CLOSEPAREN
	OPENBRACKET
	CLOSEBRACKET
	ADD
	SUB
	MUL
	DIV
	MOD
	NOT   // !, not
	COLON // :
	COMMA // ,
	DOT   // .
	DE    // ==
	NE    // !=
	GT    // >
	LT    // <
	GTE   // >=
	LTE   // <=

	// Reserved Words (Key Words)
	IF
	FOR
	WHILE
	FUN
	AND // and, &&
	OR  // or, ||

	// Equals
	EQUALS // =
)

var TokensList = []string{
	// Data Types
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	RETURN:  "RETURN",
	INT:     "INT",
	FLOAT:   "FLOAT",
	STRING:  "STRING",
	CHAR:    "CHAR",

	// Reserved Characters
	VAR:          "VAR",
	CONST:        "CONST",
	OPENCURLY:    "OPENCURLY",    // {
	CLOSECURLY:   "CLOSECURLY",   // }
	OPENPAREN:    "OPENPAREN",    // (
	CLOSEPAREN:   "CLOSEPAREN",   // )
	OPENBRACKET:  "OPENBRACKET",  // [
	CLOSEBRACKET: "CLOSEBRACKET", // ]
	ADD:          "ADD",
	SUB:          "SUB",
	MUL:          "MUL",
	DIV:          "DIV",
	MOD:          "MOD",
	NOT:          "NOT",   // !
	COLON:        "COLON", // :
	COMMA:        "COMMA", // ,
	DOT:          "DOT",   // .
	DE:           "DE",    // ==
	NE:           "NE",    // !=
	GT:           "GT",    // >
	LT:           "LT",    // <
	GTE:          "GTE",   // >=
	LTE:          "LTE",   // <=

	// Reserved Words (Key Words)
	IF:    "IF",
	FOR:   "FOR",
	WHILE: "WHILE",
	FUN:   "FUN",
	AND:   "AND", // and, &&
	OR:    "OR",  // or, ||

	// Assignment
	EQUALS: "EQUALS", // =
}

type TokenItem struct {
	pos       Position
	tokenType Token
	value     string
}

////////////
// Lexing //////////////

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() ([]TokenItem, error) {
	tokenList := []TokenItem{}
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				EOFPos := Position{line: l.pos.line, column: l.pos.column}
				tokenList = append(tokenList, TokenItem{EOFPos, EOF, ""})
				return tokenList, nil
			}
			// if it finds an error while reading that is not EOF
			return nil, err
		}

		l.pos.column++

		switch r {
		case '\n':
			l.resetPosition()
			continue
		case '+':
			tokenList = append(tokenList, TokenItem{l.pos, ADD, "+"})
		case '-':
			tokenList = append(tokenList, TokenItem{l.pos, SUB, "-"})
		case '*':
			tokenList = append(tokenList, TokenItem{l.pos, MUL, "*"})
		case '/':
			tokenList = append(tokenList, TokenItem{l.pos, DIV, "/"})
		case '%':
			tokenList = append(tokenList, TokenItem{l.pos, MOD, "%"})
		case '=':
			equalPos := l.pos

			err := l.goBack()
			if err != nil {
				return nil, err
			}

			lit, equalType, err := l.lexEquals()
			if err != nil {
				return nil, err
			}

			tokenList = append(tokenList, TokenItem{equalPos, equalType, lit})
		case '(':
			tokenList = append(tokenList, TokenItem{l.pos, OPENPAREN, "("})
		case ')':
			tokenList = append(tokenList, TokenItem{l.pos, CLOSEPAREN, ")"})
		case '{':
			tokenList = append(tokenList, TokenItem{l.pos, OPENCURLY, "{"})
		case '}':
			tokenList = append(tokenList, TokenItem{l.pos, CLOSECURLY, "}"})
		case '[':
			tokenList = append(tokenList, TokenItem{l.pos, OPENBRACKET, "["})
		case ']':
			tokenList = append(tokenList, TokenItem{l.pos, CLOSEBRACKET, "]"})
		case '!':
			notPos := l.pos

			err := l.goBack()
			if err != nil {
				return nil, err
			}

			lit, notType, err := l.lexNot()
			if err != nil {
				return nil, err
			}

			tokenList = append(tokenList, TokenItem{notPos, notType, lit})
		case ':':
			tokenList = append(tokenList, TokenItem{l.pos, COLON, ":"})
		case ',':
			tokenList = append(tokenList, TokenItem{l.pos, COMMA, ","})
		case '.':
			tokenList = append(tokenList, TokenItem{l.pos, DOT, "."})
		case '&':
			andPos := l.pos

			err := l.goBack()
			if err != nil {
				return nil, err
			}

			lit, andType, err := l.lexAnd()
			if err != nil {
				return nil, err
			}

			tokenList = append(tokenList, TokenItem{andPos, andType, lit})
		case '|':
			orPos := l.pos

			err := l.goBack()
			if err != nil {
				return nil, err
			}

			lit, orType, err := l.lexOr()
			if err != nil {
				return nil, err
			}

			tokenList = append(tokenList, TokenItem{orPos, orType, lit})
		case '<':
			ltPos := l.pos

			err := l.goBack()
			if err != nil {
				return nil, err
			}

			lit, ltType, err := l.lexLessThan()
			if err != nil {
				return nil, err
			}

			tokenList = append(tokenList, TokenItem{ltPos, ltType, lit})
		case '>':
			gtPos := l.pos

			err := l.goBack()
			if err != nil {
				return nil, err
			}

			lit, gtType, err := l.lexGreaterThan()
			if err != nil {
				return nil, err
			}

			tokenList = append(tokenList, TokenItem{gtPos, gtType, lit})
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				intPos := l.pos

				err := l.goBack()
				if err != nil {
					return nil, err
				}

				lit, varType, err := l.lexNum()
				if err != nil {
					return nil, err
				}

				tokenList = append(tokenList, TokenItem{intPos, varType, lit})
			} else if unicode.IsLetter(r) {
				letterPos := l.pos

				err := l.goBack()
				if err != nil {
					return nil, err
				}

				lit, err := l.lexIdent()
				if err != nil {
					return nil, err
				}

				switch lit {
				case "func", "fun", "fn", "funky", "def":
					tokenList = append(tokenList, TokenItem{letterPos, FUN, lit})
				case "if", "❓":
					tokenList = append(tokenList, TokenItem{letterPos, IF, lit})
				case "for":
					tokenList = append(tokenList, TokenItem{letterPos, FOR, lit})
				case "while", "loop", "forever":
					tokenList = append(tokenList, TokenItem{letterPos, WHILE, lit})
				case "var", "val", "define", "let":
					tokenList = append(tokenList, TokenItem{letterPos, VAR, lit})
				case "const":
					tokenList = append(tokenList, TokenItem{letterPos, CONST, lit})
				case "and", "plus":
					tokenList = append(tokenList, TokenItem{letterPos, AND, lit})
				case "or", "perhaps":
					tokenList = append(tokenList, TokenItem{letterPos, OR, lit})
				case "not":
					tokenList = append(tokenList, TokenItem{letterPos, NOT, lit})
				case "return":
					tokenList = append(tokenList, TokenItem{letterPos, RETURN, lit})
				default:
					tokenList = append(tokenList, TokenItem{letterPos, IDENT, lit})
				}
			} else if r == '"' {
				stringPos := l.pos

				err := l.goBack()
				if err != nil {
					return nil, err
				}

				lit, varType, err := l.lexString()
				if err != nil {
					return nil, err
				}

				tokenList = append(tokenList, TokenItem{stringPos, varType, lit})
			} else {
				tokenList = append(tokenList, TokenItem{l.pos, ILLEGAL, string(r)})
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) goBack() error {
	l.pos.column--
	if err := l.reader.UnreadRune(); err != nil {
		return err
	}
	return nil
}

func (l *Lexer) lexNum() (string, Token, error) {
	var literal string
	varType := INT
	dotCount := 0
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return literal, varType, nil
			}
			return "", ILLEGAL, err
		}

		l.pos.column++
		if unicode.IsDigit(r) {
			literal += string(r)
		} else if r == '.' {
			if dotCount == 0 {
				varType = FLOAT
			} else {
				varType = ILLEGAL
			}
			dotCount++
			literal += string(r)
		} else {
			err := l.goBack()
			if err != nil {
				return "", ILLEGAL, err
			}

			return literal, varType, nil
		}
	}
}

func (l *Lexer) lexIdent() (string, error) {
	var literal string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return literal, nil
			}
			return "", err
		}

		l.pos.column++
		if unicode.IsLetter(r) ||
			unicode.IsDigit(r) ||
			r == '_' {
			literal += string(r)
		} else {
			err := l.goBack()
			if err != nil {
				return "", err
			}

			return literal, nil
		}
	}
}

func (l *Lexer) lexString() (string, Token, error) {
	var literal string

	// Skip the opening quote
	r, _, err := l.reader.ReadRune()
	if err != nil {
		return "", ILLEGAL, err
	}
	l.pos.column++

	if r != '"' {
		return "", ILLEGAL, nil
	}

	// Read until closing quote
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// Unterminated string
				return literal, ILLEGAL, nil
			}
			return literal, ILLEGAL, err
		}

		l.pos.column++

		if r == '"' {
			// Found closing quote, we're done
			break
		}

		literal += string(r)
	}

	// Always return STRING type for quoted content
	// A single character in quotes is still a string, not a char
	return literal, STRING, nil
}

func (l *Lexer) lexEquals() (string, Token, error) {
	var equalType Token
	equalCount := 0
	var lit strings.Builder

readLoop:
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break readLoop
			}
			return "", ILLEGAL, err
		}

		l.pos.column++

		switch r {
		case '=':
			lit.WriteRune(r)
			equalCount++
		default:
			if err := l.goBack(); err != nil {
				return lit.String(), ILLEGAL, err
			}
			break readLoop
		}
	}

	switch equalCount {
	case 1:
		equalType = EQUALS
	case 2:
		equalType = DE
	default:
		equalType = ILLEGAL
	}

	return lit.String(), equalType, nil
}

func (l *Lexer) lexNot() (string, Token, error) {
	var notType Token
	notCount := 0
	equalCount := 0
	var lit strings.Builder

readLoop:
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break readLoop
			}
			return "", ILLEGAL, err
		}

		l.pos.column++

		switch r {
		case '!':
			lit.WriteRune(r)
			notCount++
		case '=':
			lit.WriteRune(r)
			equalCount++
		default:
			err := l.goBack()
			if err != nil {
				return lit.String(), ILLEGAL, err
			}
			break readLoop
		}
	}

	if notCount == 1 && equalCount == 1 {
		notType = NE
	} else if notCount == 1 && equalCount == 0 {
		notType = NOT
	} else {
		notType = ILLEGAL
	}

	return lit.String(), notType, nil
}

func (l *Lexer) lexAnd() (string, Token, error) {
	var andType Token
	andCount := 0
	var lit strings.Builder
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", ILLEGAL, err
		}

		l.pos.column++

		if r == '&' {
			lit.WriteRune(r)
			andCount++
		} else {
			err := l.goBack()
			if err != nil {
				return lit.String(), ILLEGAL, err
			}
			break
		}
	}

	if andCount == 2 {
		andType = AND
	} else {
		andType = ILLEGAL
	}

	return lit.String(), andType, nil
}

func (l *Lexer) lexOr() (string, Token, error) {
	var orType Token
	orCount := 0
	var lit strings.Builder
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", ILLEGAL, err
		}

		l.pos.column++

		if r == '|' {
			lit.WriteRune(r)
			orCount++
		} else {
			err := l.goBack()
			if err != nil {
				return lit.String(), ILLEGAL, err
			}
			break
		}
	}

	if orType == 2 {
		orType = OR
	} else {
		orType = ILLEGAL
	}

	return lit.String(), orType, nil
}

func (l *Lexer) lexLessThan() (string, Token, error) {
	var ltType Token
	var lit strings.Builder
	ltCount := 0
	equalCount := 0

readLoop:
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break readLoop
			}
			return "", ILLEGAL, err
		}

		l.pos.column++

		switch r {
		case '<':
			lit.WriteRune(r)
			ltCount++
		case '=':
			lit.WriteRune(r)
			equalCount++
		default:
			err := l.goBack()
			if err != nil {
				return lit.String(), ILLEGAL, err
			}
			break readLoop
		}
	}

	if ltCount == 1 && equalCount == 1 {
		ltType = LTE
	} else if ltCount == 1 && equalCount == 0 {
		ltType = LT
	} else {
		ltType = ILLEGAL
	}

	return lit.String(), ltType, nil
}

func (l *Lexer) lexGreaterThan() (string, Token, error) {
	var gtType Token
	var lit strings.Builder
	gtCount := 0
	equalCount := 0

readLoop:
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break readLoop
			}
			return "", ILLEGAL, err
		}

		l.pos.column++

		switch r {
		case '>':
			lit.WriteRune(r)
			gtCount++
		case '=':
			lit.WriteRune(r)
			equalCount++
		default:
			err := l.goBack()
			if err != nil {
				return lit.String(), ILLEGAL, err
			}
			break readLoop
		}
	}

	if gtCount == 1 && equalCount == 1 {
		gtType = GTE
	} else if gtCount == 1 && equalCount == 0 {
		gtType = GT
	} else {
		gtType = ILLEGAL
	}

	return lit.String(), gtType, nil
}
