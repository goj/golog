package lexer

import (
	"unicode/utf8"
	"strings"
	"fmt"
)

type Lexer interface {
	Next() Token
	Peek() Token
	Highlight(t Token, err string)
}

func NewLexer(name, input string) Lexer {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan Token),
	}
	go l.run() // Concurrently run state machine.
	return l
}

func LexAll(name, input string) []Token {
	ret := []Token{}
	l := NewLexer(name, input)
	for t := l.Next(); t.Typ != TknEOF; t = l.Next() {
		ret = append(ret, t)
	}
	return ret
}

func (l *lexer) Next() (result Token) {
	if l.hasPeek {
		result = l.peeked
		l.hasPeek = false
	} else {
		result = <-l.tokens
	}
	return
}

func (l *lexer) Peek() Token {
	if !l.hasPeek {
		l.peeked = <-l.tokens
		l.hasPeek = true
	}
	return l.peeked
}

func (l *lexer) Highlight(t Token, err string) {
	line, lno, col := findLineOf(l.input, t.pos)
	fmt.Printf("%s:%d:%d: %s\n", l.name, lno+1, col+1, err)
	fmt.Print(line)
	fmt.Print(strings.Repeat(" ", col))
	fmt.Println(strings.Repeat("~", utf8.RuneCountInString(t.Val)))
}

func findLineOf(s string, pos int) (string, int, int) {
	lines := strings.SplitAfter(s, "\n")
	for lno, line := range lines {
		if pos < len(line) {
			return line, lno, pos
		}
		pos -= len(line)
	}
	return "", 0, 0
}
