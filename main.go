package main

import (
	"github.com/goj/golog/lexer"
	"github.com/goj/golog/parser"
	"github.com/kr/pretty"
	"io/ioutil"
	"log"
	"fmt"
)

func main() {
	filename := "test/socrates.pl"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("couldn't open the source code file: %v\n", err)
	}
	src := string(data)
	// toks := lexer.LexAll(filename, src)
	lex := lexer.NewLexer(filename, src)
	ast := parser.Parse(lex)
	fmt.Printf("got the following AST:\n%# v\n\n"+
		"for the following program:\n\n%s\n", pretty.Formatter(ast), src)
}
