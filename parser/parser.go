package parser

import "fmt"
import "github.com/goj/golog/lexer"

func Parse(l lexer.Lexer) Program {
	return parseProgram(l)
}

// Program -> ClauseOrFact* EOF
func parseProgram(l lexer.Lexer) (result Program) {
	for l.Peek().Typ != lexer.TknEOF {
		result = append(result, parseClauseOrFact(l))
	}
	return
}

// ClauseOrFact -> ClauseHead (`.` | `:-` Term `.`)
func parseClauseOrFact(l lexer.Lexer) (result Clause) {
	result.Head = parseClauseHead(l)
	switch l.Next().Typ {
	case lexer.TknDot:
		result.Body = Pred{Name: "true", Args: []Term{}}
	case lexer.TknColonDash:
		result.Body = parseTerm(l)
		nextToken(l, lexer.TknDot, "clause body")
	}
	return
}

// ClauseHead -> atom [ArgList]
func parseClauseHead(l lexer.Lexer) (result ClauseHead) {
	tkn := nextToken(l, lexer.TknAtom, "clause head")
	result.Name = tkn.Val
	result.Args = []Term{}
	if l.Peek().Typ == lexer.TknOpenParen {
		result.Args = parseArgList(l)
	}
	return
}

// ArgList -> [`(` Term (`,` Term) `)`]
func parseArgList(l lexer.Lexer) (result []Term) {
	nextToken(l, lexer.TknOpenParen, "argument list")
	result = []Term{}
	for {
		result = append(result, parseTerm(l))
		tkn := l.Next()
		switch tkn.Typ {
		case lexer.TknComma:
			continue
		case lexer.TknCloseParen:
			return
		default:
			synErrExpected(l, "`,` or `(`", "argument list", tkn)
		}
	}
	return // never happens, as default panics
}

func parseTerm(l lexer.Lexer) Term {
	tkn := l.Peek()
	if tkn.Typ == lexer.TknVariable {
		return parseVariable(l)
	}
	return parseClauseHead(l)
}

func parseVariable(l lexer.Lexer) Variable {
	tkn := nextToken(l, lexer.TknVariable, "a term")
	return Variable{Name: tkn.Val}
}

func nextToken(l lexer.Lexer, expectedType lexer.TokenType, what string) lexer.Token {
	tkn := l.Next()
	if tkn.Typ != expectedType {
		synErrExpected(l, prettyType(expectedType), what, tkn)
	}
	return tkn
}

func synErrExpected(l lexer.Lexer, typeName, what string, tkn lexer.Token) {
	err := fmt.Sprintf("expected %s when parsing %s, got %v", typeName, what, tkn)
	l.Highlight(tkn, err)
	panic(err)
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
		return "`:=`"
	case lexer.TknComma:
		return "a comma"
	case lexer.TknDot:
		return "a dot"
	}
	return fmt.Sprintf("???(%d)", t)
}
