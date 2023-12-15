package main

import (
	"fmt"

	"github.con/shreyascodes-tech/sss-lang/lexer"
)

const input = `:root {
	--port: 8080;
}

.a[href="/"][method="GET"] {
	content: "Hello World";
}

a[href="/data"][method="GET"] {
	content: json({
		"message": "Hello World"
	});
}
`

func debugPrint(i interface{}) {
	fmt.Printf("%+v\n", i)
}

func main() {
	l := lexer.New(input)
	for {
		t := l.Next()
		debugPrint(t)
		if t.Type == lexer.EOF {
			break
		}
	}
}
