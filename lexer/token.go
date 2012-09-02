package lexer

import "fmt"

type Token struct {
	Typ TokenType
	Val string
	Pos int
}

type TokenType int

const (
	TknError TokenType = iota
	TknOpenParen
	TknCloseParen
	TknDot
	TknComma
	TknColonDash
	TknAtom
	TknVariable
	TknNumber
	TknEOF
)

func (t TokenType) String() string {
	switch t {
	case TknError:
		return "ERROR"
	case TknOpenParen:
		return "open paren"
	case TknCloseParen:
		return "close paren"
	case TknDot:
		return "`.`"
	case TknComma:
		return "a comma"
	case TknColonDash:
		return "`:-`"
	case TknAtom:
		return "an atom"
	case TknVariable:
		return "a variable"
	case TknNumber:
		return "a number"
	case TknEOF:
		return "EOF"
	}
	return fmt.Sprintf("other token (%d)", t)
}

func (t Token) String() string {
	switch t.Typ {
	case TknError:
		return "ERROR: " + t.Val
	case TknEOF:
		return "EOF"
	}
	if len(t.Val) > 10 {
		return fmt.Sprintf("%s %.7q...", t.Typ, t.Val)
	}
	return fmt.Sprintf("%s %q", t.Typ, t.Val)
}
