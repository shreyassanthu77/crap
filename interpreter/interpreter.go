package interpreter

import (
	"fmt"
	"os"

	"github.com/shreyassanthu77/cisp/ast"
	"github.com/shreyassanthu77/cisp/lexer"
	"github.com/shreyassanthu77/cisp/parser"
)

type Interpreter struct {
	env     *Environment
	program ast.Program
}

func (i *Interpreter) throw(format string, args ...any) {
	inEr := fmt.Sprintf(format, args...)
	fmt.Printf("ERROR: %s\n", inEr)
	os.Exit(1)
}

func (i *Interpreter) throwAtSpan(span lexer.Span, format string, args ...any) {
	fmt.Printf("%d:%d: ", span.Start.Line, span.Start.Col)
	i.throw(format, args...)
}

func New(source string) (*Interpreter, error) {
	lexer := lexer.New(source)
	parser := parser.New(lexer)

	program, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	root := NewRootEnv()
	return &Interpreter{
		env:     root,
		program: program,
	}, nil
}

func (i *Interpreter) fork(cb func() ast.Value) ast.Value {
	i.env = i.env.fork()
	res := cb()
	i.env = i.env.Parent
	return res
}

func (i *Interpreter) Run() ast.Value {
	root := i.env
	var mainRule ast.Rule
	for _, rule := range i.program.Rules {
		switch rule := rule.(type) {
		case ast.AtRule:
			i.throw("global at-rules not supported yet")
		case ast.Rule:
			if existing, err := root.setFn(rule); err != nil {
				i.throwAtSpan(rule.Span, "function %s already defined at %d:%d", existing.Selector.Identifier.Name, existing.Span.Start.Line, existing.Span.Start.Col)
			}
			if rule.Selector.Identifier.Name == "main" {
				mainRule = rule
			}
		}
	}

	return i.evalRule(mainRule, []ast.Value{})
}
