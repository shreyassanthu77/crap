package main

import (
	"encoding/json"
	"fmt"

	"github.com/shreyassanthu77/cisp/lexer"
	"github.com/shreyassanthu77/cisp/parser"
)

func main() {
	input := `
factorial[n] {
	@if n == 1 {
			@return 1;
	}
	@return $n * factorial($n - 1);
}

main {
	print: factorial(5);
}
	`
	lex := lexer.New(input)
	par := parser.New(lex)

	ast, err := par.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	jsoned, err := json.MarshalIndent(ast, "", "  ")
	fmt.Println(string(jsoned))
}
