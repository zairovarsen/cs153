package token

import (
	"fmt"
	"strconv"
)

// Represents a lexical token
type Tag int

// Enums representing all possible token types
const (
	// Special tokens
	ILLEGAL Tag = iota
	EOF
	COMMENT

	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	CHAR   // 'a'
	STRING // "abc"

	// Bool
	TRUE
	FALSE

	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND // &
	OR  // |
	XOR // ^
	SHL // <<
	SHR // >>

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN // &=
	OR_ASSIGN  // |=
	XOR_ASSIGN // ^=
	SHL_ASSIGN // <<=
	SHR_ASSIGN // >>=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // =
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	keyword_start
	// Keywords
	NIL
	BREAK
	CONST
	CONTINUE

	ELSE
	FOR
	WHILE

	FUNC
	IF

	RETURN

	CLASS
	VAR
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND: "&",
	OR:  "|",
	XOR: "^",
	SHL: "<<",
	SHR: ">>",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN: "&=",
	OR_ASSIGN:  "|=",
	XOR_ASSIGN: "^=",
	SHL_ASSIGN: "<<=",
	SHR_ASSIGN: ">>=",

	LAND:  "&&",
	LOR:   "||",
	ARROW: "<-",
	INC:   "++",
	DEC:   "--",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	DEFINE:   ":=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	BREAK:    "break",
	CONST:    "const",
	CONTINUE: "continue",

	ELSE: "else",
	FOR:  "for",

	FUNC: "func",
	IF:   "if",

	RETURN: "return",

	VAR:   "var",
	CLASS: "class",
}

func (t Tag) String() string {
	s := ""
	if 0 <= t && t < Tag(len(tokens)) {
		s = tokens[t]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

func TokenTypeToName(t Tag) string {
	return tokens[t]
}

var keywords = map[string]Tag{}

func initKeywords() {
	keywords := make(map[string]Tag, keyword_end-(keyword_start)+1)
	for i := keyword_start + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(word string) Tag {
	v, ok := keywords[word]
	if ok {
		return v
	}
	return IDENT
}

type Token struct {
	Tag Tag
	Loc Loc
	Pos Position
}

func NewToken(tokenType Tag, loc Loc, pos Position) *Token {
	return &Token{
		tokenType,
		loc,
		pos,
	}
}

func (t Token) Lexeme(source string) string {
	return source[t.Loc.Start:t.Loc.End]
}

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
func (t Tag) Precedence() Precedence {
	switch t {
	case DEFINE:
		return PREC_ASSIGNMENT
	case LOR:
		return PREC_LOR
	case LAND:
		return PREC_LAND
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return PREC_EQUALITY_COMPARISON
	case ADD, SUB, OR, XOR:
		return PREC_OPER_LOW
	case MUL, QUO, REM, SHL, SHR, AND:
		return PREC_OPER_HIGH
	}
	return PREC_NONE
}

type Precedence int

const (
	PREC_NONE Precedence = iota
	PREC_ASSIGNMENT
	PREC_LOR
	PREC_LAND
	PREC_EQUALITY_COMPARISON // neq, lss, leq, gtr, geq
	PREC_OPER_LOW            // add, sub, or, xor
	PREC_OPER_HIGH           // mul, quo, rem, shl, shr, and, and_not
	PREC_UNARY
	PREC_CALL
	PREC_PRIMARY
)

// Position of a token.
type Position struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}

func (p Position) String() string {
	filename := p.Filename
	if filename == "" {
		return fmt.Sprintf("%d:%d", p.Line, p.Column)
	}
	return fmt.Sprintf("%s:%d:%d", filename, p.Line, p.Column)
}

type Loc struct {
	Start int
	End   int
}

func NewLoc(start int, end int) Loc {
	return Loc{start, end}
}
