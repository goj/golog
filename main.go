package main

import (
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
	ast, errs := parser.Parse(filename, src, parser.Debug)
	for _, err := range errs {
		fmt.Println(err)
	}
	if ast != nil {
		fmt.Printf("got the following AST:\n%# v\n\n"+
			"for the following program:\n\n%s", pretty.Formatter(ast), src)
	}
}
