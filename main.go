package main

import (
	"fmt"

	"github.com/shreyassanthu77/cisp/interpreter"
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

customPrint[var] {
	print: $var;
}

main {
	--msg: "Hello, World!";
	print: 1 + 2*3 - 4 / 2;
	customPrint: $msg + " " + "This is a test";
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
