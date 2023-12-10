package main

import (
	"fmt"
)

const input = `:root {
	--port: 8080;
}

.a[href="/"][method="GET"] {
	content: "Hello World";
}

`

// TODO: Add support for functions, expressions and objects
// `
// a[href="/data"][method="GET"] {
// 	content: json({
// 		"message": "Hello World"
// 	});
// }
// `

func debugPrint(i interface{}) {
	fmt.Printf("%+v\n", i)
}
func main() {
	// p := parser.Parse(input)
	//
	// debugPrint(p)
}
