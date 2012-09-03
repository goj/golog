package parser

import (
	"fmt"
	"github.com/goj/golog/lexer"
	"os"
	"strings"
	"unicode/utf8"
)

type parser struct {
	name      string             // used only for error reports.
	input     string             // the string being scanned.
	debug     bool               // is debug tracing enabled
	trace     []string           // debug trace
	peekedToken *lexer.Token        // current token
	tokens    <-chan lexer.Token // channel of scanned tokens.
}

type ParserFlags int

const (
	Debug ParserFlags = 1 << iota
)

func ParseString(filename, content string, flags ParserFlags) Program {
	tokens := lexer.Tokens(content)
	p := parser{
		name:      filename,
		input:     content,
		debug:     flags & Debug != 0,
		tokens:    tokens,
	}
	return p.parseProgram()
}

func prettyType(t lexer.TokenType) string {
	switch t {
	case lexer.TknAtom:
		return "an atom"
	case lexer.TknVariable:
		return "a variable"
	case lexer.TknOpenParen:
		return "`(`"
	case lexer.TknCloseParen:
		return "`)`"
	case lexer.TknColonDash:
		return "`:-`"
	case lexer.TknComma:
		return "a comma"
	case lexer.TknDot:
		return "a dot"
	}
	return fmt.Sprintf("???(%d)", t)
}

func (p *parser) nextToken() (result lexer.Token) {
	if p.peekedToken != nil {
		result = *p.peekedToken
		p.peekedToken = nil
	} else {
		result = <-p.tokens
	}
	return
}

func (p *parser) peekToken() lexer.Token {
	if p.peekedToken == nil {
		lookahead := <-p.tokens
		p.peekedToken = &lookahead
	}
	return *p.peekedToken
}

func (p *parser) expectToken(expectedType lexer.TokenType) lexer.Token {
	tkn := p.nextToken()
	if tkn.Typ != expectedType {
		p.synErrExpected(prettyType(expectedType), tkn)
	}
	return tkn
}

func (p *parser) synErrExpected(typeName string, tkn lexer.Token) {
	whatsParsed := p.trace[len(p.trace)-1]
	err := fmt.Sprintf("expected %s when parsing %s, got %v", typeName, whatsParsed, tkn)
	p.syntaxError(tkn, err)
	panic(err) // FIXME: don't panic on syntax errors
}

func (p *parser) syntaxError(t lexer.Token, err string) {
	line, lno, col := findLineOf(p.input, t.Pos)
	fmt.Printf("%s:%d:%d: %s\n", p.name, lno+1, col+1, err)
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

func trace(p *parser, msg string) *parser {
	if p.debug {
		fmt.Fprintf(os.Stderr, "%s+%s\n", indentFor(p.trace), msg)
	}
	p.trace = append(p.trace, msg)
	return p
}

func un(p *parser) {
	last := p.trace[len(p.trace)-1]
	p.trace = p.trace[:len(p.trace)-1]
	if p.debug {
		fmt.Fprintf(os.Stderr, "%s-%s\n", indentFor(p.trace), last)
	}
}

func indentFor(s []string) string {
	return strings.Repeat(" ", len(s))
}
