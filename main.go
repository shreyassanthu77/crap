package main

import (
	"fmt"

	"github.com/shreyassanthu77/cisp/interpreter"
	"github.com/shreyassanthu77/cisp/lexer"
	"github.com/shreyassanthu77/cisp/parser"
)

func main() {
	input := `

ex[n] {
	@return $n;
}

factorial[n] {
	@if $n == 0 || $n == 1 {
			--result: 1;
			@return $result;
	} @else {
		print: "else block";
	}
	@return $n * factorial($n - 1);
}

main {
	--result: factorial(5);
	@return $result;
}
	`
	lex := lexer.New(input)
	par := parser.New(lex)

	ast, err := par.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := interpreter.Eval(ast)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Result:", res)

	// jsoned, err := json.MarshalIndent(ast, "", "  ")
	// fmt.Println(string(jsoned))
}
