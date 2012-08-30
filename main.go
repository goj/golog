package main

import (
	"github.com/goj/golog/lexer"
	"io/ioutil"
	"log"
	"fmt"
)

func main() {
	data, err := ioutil.ReadFile("test/socrates.pl")
	if err != nil {
		log.Panicf("couldn't open the source code file: %v\n", err)
	}
	src := string(data)
	tokens := lexer.LexAll("<stdin>", src)
	fmt.Printf("got the following tokens:\n%v\n\n"+
		"for the following program:\n\n%s", tokens, src)
}
