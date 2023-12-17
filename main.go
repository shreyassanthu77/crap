package main

import (
	"fmt"

	"github.con/shreyascodes-tech/sss-lang/ast"
	"github.con/shreyascodes-tech/sss-lang/parser"
)

const input = `/* This is a comment
	 that spans multiple lines */
@import "other.css";
	:root a > b.c {
	color: red;
	
	b {
		color: blue;
	}
}

/* .a[href="/"][method="GET"] {
	content: "Hello World";
}

a[href="/data"][method="GET"] {
	content: json({
		"message": "Hello World"
	});
} */
`

func debugPrint(i interface{}) {
	fmt.Printf("%+v\n", i)
}

func main() {
	// l := lexer.New(input)
	// for {
	// 	t := l.Next()
	// 	debugPrint(t)
	// 	if t.Type == lexer.EOF {
	// 		break
	// 	}
	// }
	p := parser.New(input)
	stylesheet := p.Parse()

	printRules(stylesheet.Rules)
}

func printRules(rules []ast.RuleNode) {
	for _, rule := range rules {
		switch rule.Type {
		case ast.AT_RULE:
			rule := rule.Rule.(ast.AtRule)
			fmt.Printf("@%s %v {", rule.Name, rule.Params)
			if len(rule.Rules) > 0 {
				fmt.Println()
			}
			printRules(rule.Rules)
			fmt.Println("}")
		case ast.STYLE_RULE:
			rule := rule.Rule.(ast.StyleRule)
			format := ""
			args := []interface{}{}
			for _, selector := range rule.Selectors {
				format += "%s"
				args = append(args, selector)
			}
			fmt.Printf(format+" {", args...)
			if len(rule.Rules) > 0 {
				fmt.Println()
			}
			printRules(rule.Rules)
			fmt.Println("}")
		case ast.DECL_RULE:
			rule := rule.Rule.(ast.Declaration)
			fmt.Printf("%s: %v;\n", rule.Name, rule.Value)
		}
	}
}
