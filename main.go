package main

import (
	"fmt"

	"github.com/shreyassanthu77/cisp/interpreter"
	"github.com/shreyassanthu77/cisp/lexer"
	"github.com/shreyassanthu77/cisp/parser"
)

func main() {
	input := `
fibonacci.rec[n][a=0][b=1] {
	@if $n == 0 {
		@return $a;
	}
	print: $a;
	@return fibonacci.rec($n - 1, $b, $a + $b);
}

fibonacci[n] {
	@return fibonacci.rec($n, (), ());
}

main {
	fibonacci: 10;
}`
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
