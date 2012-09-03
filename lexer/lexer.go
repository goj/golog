package lexer

import (
	"unicode"
	"unicode/utf8"
)

type lexer struct {
	input   string // the string being scanned.
	start   int    // start position of this token.
	pos     int    // current position in the input.
	tokens  chan Token // channel of scanned tokens.
}

type stateFn func(*lexer) stateFn

var simpleTokens = map[rune]TokenType{
	'.': TknDot,
	'(': TknOpenParen,
	')': TknCloseParen,
}

var inVarAtom = []*unicode.RangeTable{
	unicode.Letter,
	unicode.Number,
	unicode.Pc, // connector punctuation like '_'
}

func lexTopLevel(l *lexer) stateFn {
	for l.pos < len(l.input) {
		r, rlen := l.runeAt(l.start)
		l.pos += rlen
		simpleType, isSimple := simpleTokens[r]
		switch {
		case isSimple:
			l.emit(simpleType)
		case unicode.IsSpace(r):
			l.start += rlen
		case unicode.IsUpper(r):
			return lexVariable
		case unicode.IsLower(r):
			return lexSimpleAtom
		case r == ':':
			next, _ := l.runeAt(l.pos)
			if next == '-' {
				l.pos += 1
				l.emit(TknColonDash)
				continue
			}
			fallthrough
		}
	}
	return nil
}

func lexVariable(l *lexer) stateFn {
	l.forward(inVarAtom)
	l.emit(TknVariable)
	return lexTopLevel
}

func lexSimpleAtom(l *lexer) stateFn {
	l.forward(inVarAtom)
	l.emit(TknAtom)
	return lexTopLevel
}

func (l *lexer) runeAt(pos int) (rune, int) {
	return utf8.DecodeRuneInString(l.input[pos:])
}

func (l *lexer) forward(tables []*unicode.RangeTable) {
	for {
		r, rlen := utf8.DecodeRuneInString(l.input[l.pos:])
		if !unicode.IsOneOf(tables, r) {
			break
		}
		l.pos += rlen
	}
}

func (l *lexer) run() {
	for state := lexTopLevel; state != nil; {
		state = state(l)
	}
	l.emit(TknEOF)
	close(l.tokens) // No more tokens will be delivered.
}

// emit passes an Token back to the client.
func (l *lexer) emit(t TokenType) {
	l.tokens <- Token{t, l.input[l.start:l.pos], l.start}
	l.start = l.pos
}
