package lexer

import (
	"unicode"
	"fmt"
	"unicode/utf8"
)

type lexer struct {
	name   string     // used only for error reports.
	input  string     // the string being scanned.
	start  int        // start position of this token.
	pos    int        // current position in the input.
	width  int        // width of last rune read from input.
	tokens chan Token // channel of scanned tokens.
}

type stateFn func(*lexer) stateFn

var simpleTokens = map[rune]tokenType {
	'.': TknDot,
	'(': TknOpenParen,
	')': TknCloseParen,
}

func lexTopLevel(l *lexer) stateFn {
	for l.pos < len(l.input) {
		fmt.Printf("lexTopLevel: %d\n", l.start)
		r, rlen := utf8.DecodeRuneInString(l.input[l.start:])
		simpleType, isSimple := simpleTokens[r]
		switch {
		case isSimple:
			l.emit(simpleType, rlen)
		case unicode.IsUpper(r):
			return nil // lexVariable
		}
	}
	return nil
}

func (l *lexer) run() {
	for state := lexTopLevel; state != nil; {
		state = state(l)
	}
	l.emit(TknEOF, 0)
	close(l.tokens) // No more tokens will be delivered.
}

func Lex(name, input string) (*lexer, chan Token) {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan Token),
	}
	go l.run() // Concurrently run state machine.
	return l, l.tokens
}

// emit passes an Token back to the client.
func (l *lexer) emit(t tokenType, extraLen int) {
	l.pos += extraLen
	l.tokens <- Token{t, l.input[l.start:l.pos]}
	fmt.Printf("emitting %v %v\n", t, l.input[l.start:l.pos])
	l.start = l.pos
}
