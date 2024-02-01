package main

import (
	"fmt"
	"github.com/shreyassanthu77/cisp/lexer"
	"github.com/shreyassanthu77/cisp/parser"
)

func main() {
	input := `
factorial[n] {
    @if expr($n) {
        @return 1;
    }

}
	`
	lex := lexer.New(input)
	par := parser.New(lex)

	ast, err := par.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", ast)
}
