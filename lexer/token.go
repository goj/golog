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
	TknColonDash
	TknComma
	TknAtom
	TknVariable
	TknNumber
	TknEOF
)

func (t Token) String() string {
	switch t.Typ {
	case TknError:
		return "ERROR: " + t.Val
	case TknEOF:
		return "EOF"
	case TknAtom:
		return fmt.Sprintf("'%s'", cut(t.Val))
	}
	return cut(t.Val)
}

func cut(str string) string {
	if len(str) > 10 {
		return fmt.Sprintf("%.7s...", str)
	}
	return str
}
