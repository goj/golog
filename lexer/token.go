package lexer

import "fmt"

type Token struct {
	Typ tokenType
	Val string
}

type tokenType int

const (
	TknError tokenType = iota
	TknOpenParen
	TknCloseParen
	TknDot
	TknColonDash
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
	}
	if len(t.Val) > 10 {
		return fmt.Sprintf("%.7q...", t.Val)
	}
	return fmt.Sprintf("%q", t.Val)
}
