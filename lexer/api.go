package lexer

func Tokens(input string) (<-chan Token) {
	l := &lexer{
		input:  input,
		tokens: make(chan Token),
	}
	go l.run() // Concurrently run state machine.
	return l.tokens
}

func LexAll(input string) []Token {
	ret := []Token{}
	toks := Tokens(input)
	for t := <-toks; t.Typ != TknEOF; t = <-toks {
		ret = append(ret, t)
	}
	return ret
}

