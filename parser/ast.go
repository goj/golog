package parser

type Program []Clause

type Clause struct {
	Head ClauseHead
	Body Term
}

type ClauseHead struct {
	Name string
	Args []Term
}

type Term interface{}

type Variable struct{
	Name string
}

type Pred struct {
	Name string
	Args []Term
}
