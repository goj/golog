package parser

import (
	"fmt"
	"log"
	"os"
	"github.com/goj/golog/lexer"
	"strings"
	"unicode/utf8"
)

type SyntaxError struct {
	lead, msg, ctx, sample string
}

type parser struct {
	filename string
	input    string
	tokens   <-chan lexer.Token
	next     lexer.Token
	trace    []string
	debug    bool
	errors   []SyntaxError
}

type ParserFlags int

const (
	Debug ParserFlags = 1 << iota
)

func newParser(fname, src string) parser {
	tokens := lexer.Tokens(fname, src)
	return parser{
		filename: fname,
		input:    src,
		tokens:   tokens,
		next:     <-tokens,
	}
}

func Parse(fname, src string, flags ParserFlags) (result Program, errors []SyntaxError) {
	p := newParser(fname, src)
	p.debug = flags & Debug != 0
	return p.parseProgram(), p.errors
}

// Program -> ClauseOrFact* EOF
func (p *parser) parseProgram() (result Program) {
	defer un(trace(p, "program"))
	result = Program{}
	for p.next.Typ != lexer.TknEOF {
		result = append(result, p.parseClauseOrFact())
	}
	return
}

// ClauseOrFact -> ClauseHead ClauseBody
func (p *parser) parseClauseOrFact() (result Clause) {
	defer un(trace(p, "clause or fact"))
	result.Head = p.parseClauseHead()
	result.Body = p.parseClauseBody()
	return
}

// ClauseHead -> atom [ArgList]
func (p *parser) parseClauseHead() (result ClauseHead) {
	defer un(trace(p, "clause head"))
	p.assertTokenType(lexer.TknAtom)
	result.Name = p.next.Val
	result.Args = []Term{}
	p.nextToken()
	if p.next.Typ == lexer.TknOpenParen {
		result.Args = p.parseArgList()
	}
	return
}

// ArgList -> `(` Term (`,` Term) `)`
func (p *parser) parseArgList() (result []Term) {
	defer un(trace(p, "argument list"))
	p.assertTokenType(lexer.TknOpenParen)
	result = []Term{}
	for {
		p.nextToken()
		result = append(result, p.parseTerm())
		switch p.next.Typ {
		case lexer.TknComma:
			continue
		case lexer.TknCloseParen:
			p.nextToken()
			return
		default:
			p.syntaxError("expected `,` or `)`, got %s", p.next.Typ)
		}
	}
	return // never happens, as default panics
}


// ClauseBody -> '.' | Term
func (p *parser) parseClauseBody() (result Term) {
	defer un(trace(p, "clause body"))
	switch p.next.Typ {
	case lexer.TknDot:
		result = Pred{Name: "true", Args: []Term{}}
	case lexer.TknColonDash:
		p.nextToken()
		result = p.parseTerm()
	}
	p.nextToken()
	return
}

func (p *parser) parseTerm() Term {
	defer un(trace(p, "term"))
	if p.next.Typ == lexer.TknVariable {
		return p.parseVariable()
	}
	return p.parseClauseHead()
}

func (p *parser) parseVariable() Variable {
	defer un(trace(p, "variable"))
	p.assertTokenType(lexer.TknVariable)
	defer p.nextToken()
	return Variable{Name: p.next.Val}
}

func (p *parser) nextToken() {
	p.next = <- p.tokens
	if p.debug {
		fmt.Fprintf(os.Stderr, "%s~ %s\n", indentFor(p.trace), p.next)
	}
}

func (p *parser) nextTokenAs(typ lexer.TokenType) {
	p.nextToken()
	p.assertTokenType(typ)
}

func (p *parser) assertTokenType(typ lexer.TokenType) {
	if p.next.Typ != typ {
		p.syntaxError("expected %s, got %s", typ, p.next.Typ)
	}
}

func un(p *parser) {
	last := p.trace[len(p.trace)-1]
	p.trace = p.trace[:len(p.trace)-1]
	if p.debug {
		fmt.Fprintf(os.Stderr, "%s- %s\n", indentFor(p.trace), last)
	}
}

func trace(p *parser, msg string) *parser {
	if p.debug {
		fmt.Fprintf(os.Stderr, "%s+ %s // %s\n", indentFor(p.trace), msg, p.next)
	}
	p.trace = append(p.trace, msg)
	return p
}

func indentFor(s []string) string {
	return strings.Repeat("  ", len(s))
}

func (p *parser) syntaxError(format string, args ...interface{}) {
	line, lno, col := findLineOf(p.input, p.next.Pos)
	err := SyntaxError{
		lead:   fmt.Sprintf("%s:%d:%d", p.filename, lno+1, col+1),
		msg:    fmt.Sprintf(format, args...),
		ctx:    p.trace[len(p.trace)-1],
		sample: line + underline(col, p.next),
	}
	p.errors = append(p.errors, err)
	if p.debug {
		log.Panicf("%s\n", err)
	}
}

func underline(col int, tkn lexer.Token) string {
	// TODO: support multi-line tokens
	whatLen := utf8.RuneCountInString(tkn.Val)
	return strings.Repeat(" ", col) + strings.Repeat("~", whatLen)
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

func (se SyntaxError) String() string {
	return fmt.Sprintf("%s: %s when parsing %s:\n%s\n", se.lead, se.msg, se.ctx, se.sample)
}
