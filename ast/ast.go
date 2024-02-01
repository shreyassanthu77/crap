package ast

import "fmt"

type Value interface {
	isValue()
}

type NilValue struct{}

func (n NilValue) isValue() {}

func (n NilValue) String() string {
	return "Nil"
}

type Identifier struct {
	Name string
}

func (i Identifier) isValue() {}

func (i Identifier) String() string {
	return fmt.Sprintf("Id(%s)", i.Name)
}

type String struct {
	Value string
}

func (s String) isValue() {}

func (s String) String() string {
	return fmt.Sprintf("'%s'", s.Value)
}

type Int struct {
	Value int64
}

func (n Int) isValue() {}

func (n Int) String() string {
	return fmt.Sprintf("%d", n.Value)
}

type Float struct {
	Value float64
}

func (n Float) isValue() {}

func (n Float) String() string {
	return fmt.Sprintf("%f", n.Value)
}

type Boolean struct {
	Value bool
}

func (b Boolean) isValue() {}

func (b Boolean) String() string {
	return fmt.Sprintf("Bool(%t)", b.Value)
}

type VarianleDerefValue struct {
	Variable Identifier
}

func (v VarianleDerefValue) isValue() {}

func (v VarianleDerefValue) String() string {
	return fmt.Sprintf("var(%s)", v.Variable.Name)
}

type Statement interface {
	isStatement()
}

type Declaration struct {
	Property Identifier
	Values   []Value
}

func (d Declaration) isStatement() {}

type FunctionCall struct {
	Name       string
	Parameters []Value
}

func (f FunctionCall) isValue() {}

func (f FunctionCall) String() string {
	return fmt.Sprintf("Call(%s, %v)", f.Name, f.Parameters)
}

type UnaryOp struct {
	Op    string
	Value Value
}

func (u UnaryOp) isValue() {}

type BinaryOp struct {
	Left  Value
	Op    string
	Right Value
}

func (b BinaryOp) isValue() {}

type AtRule struct {
	Name       Identifier
	Params     []Value
	Statements []Statement
}

func (r AtRule) isRule() {}

func (r AtRule) isStatement() {}

type Selector struct {
	Identifier Identifier
	Atrributes map[Identifier]Value
}

type IRule interface {
	isRule()
}

type Rule struct {
	Selector   Selector
	Statements []Statement
}

func (r Rule) isRule() {}

func (r Rule) isStatement() {}

type Program struct {
	Rules []IRule
}
