package main

import (
	"fmt"
	"github.com/shreyassanthu77/cisp/lexer"
	"github.com/shreyassanthu77/cisp/parser"
)

func main() {
	input := `@hello world {
			hi {
				foo: bar;
		}
	}
	@hello world2;
	hello {
		foo: bar;
	}`
	lex := lexer.New(input)
	par := parser.New(lex)

	ast, err := par.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", ast)
}
