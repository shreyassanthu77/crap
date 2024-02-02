package ast

import "fmt"

type Value interface {
	IsValue()
}

type NilValue struct{}

func (n NilValue) IsValue() {}

func (n NilValue) String() string {
	return "Nil"
}

type Identifier struct {
	Name string
}

func (i Identifier) IsValue() {}

func (i Identifier) String() string {
	return fmt.Sprintf("Id(%s)", i.Name)
}

type String struct {
	Value string
}

func (s String) IsValue() {}

func (s String) String() string {
	return fmt.Sprintf("'%s'", s.Value)
}

type Int struct {
	Value int64
}

func (n Int) IsValue() {}

func (n Int) String() string {
	return fmt.Sprintf("%d", n.Value)
}

type Float struct {
	Value float64
}

func (n Float) IsValue() {}

func (n Float) String() string {
	return fmt.Sprintf("%f", n.Value)
}

type Boolean struct {
	Value bool
}

func (b Boolean) IsValue() {}

func (b Boolean) String() string {
	return fmt.Sprintf("Bool(%t)", b.Value)
}

type VarianleDerefValue struct {
	Variable Identifier
}

func (v VarianleDerefValue) IsValue() {}

func (v VarianleDerefValue) String() string {
	return fmt.Sprintf("var(%s)", v.Variable.Name)
}

type Statement interface {
	IsStatement()
}

type Declaration struct {
	Property   Identifier
	Parameters []Value
}

func (d Declaration) IsStatement() {}

type FunctionCall struct {
	Fn         Identifier
	Parameters []Value
}

func (f FunctionCall) IsValue() {}

func (f FunctionCall) String() string {
	return fmt.Sprintf("Call(%s, %v)", f.Fn, f.Parameters)
}

type UnaryOp struct {
	Op    string
	Value Value
}

func (u UnaryOp) IsValue() {}

type BinaryOp struct {
	Left  Value
	Op    string
	Right Value
}

func (b BinaryOp) IsValue() {}

type AtRule struct {
	Name       string
	Parameters []Value
	Body       []Statement
}

func (r AtRule) isRule() {}

func (r AtRule) IsStatement() {}

type Attreibute struct {
	Name    Identifier
	Default Value
}

type Selector struct {
	Identifier Identifier
	Atrributes []Attreibute
}

type IRule interface {
	isRule()
}

type Rule struct {
	Selector Selector
	Body     []Statement
}

func (r Rule) isRule() {}

func (r Rule) IsStatement() {}

type Program struct {
	Rules []IRule
}
