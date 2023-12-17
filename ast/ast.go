package ast

import "fmt"

type Location struct {
	Start Span
	End   Span
}

func (l Location) String() string {
	return fmt.Sprintf("%d:%d-%d:%d", l.Start.Line, l.Start.Col, l.End.Line, l.End.Col)
}

type Span struct {
	Pos  int
	Line int
	Col  int
}

type Identifier struct {
	Loc   Location
	Value string
}

func (i Identifier) String() string {
	return fmt.Sprintf("ID(%s)", i.Value)
}

type ValueType string

const (
	STRING        ValueType = "STRING"
	NUMBER        ValueType = "NUMBER"
	HEX           ValueType = "HEX"
	IDENT         ValueType = "IDENT"
	FUNCTION_CALL ValueType = "FUNCTION_CALL"
	EXPRESSION    ValueType = "EXPRESSION"
)

type FunctionCall struct {
	Loc  Location
	Name Identifier
	Args []Value
}

func (f FunctionCall) String() string {
	return fmt.Sprintf("fn %s(%v)", f.Name.Value, f.Args)
}

type Expression struct {
	Loc   Location
	Left  Value
	Op    string
	Right Value
}

func (e Expression) String() string {
	return fmt.Sprintf("%v %s %v", e.Left, e.Op, e.Right)
}

type Value struct {
	Loc  Location
	Type ValueType
	Data interface{}
}

func (v Value) String() string {
	switch v.Type {
	case STRING:
		return fmt.Sprintf("/%s/", v.Data.(string))
	case NUMBER:
		return fmt.Sprintf("%s", v.Data.(string))
	case HEX:
		return fmt.Sprintf("#%s", v.Data.(string))
	default:
		return fmt.Sprintf("%v", v.Data)
	}
}

type Stylesheet struct {
	Loc   Location
	Rules []RuleNode
}

type RuleType string

const (
	AT_RULE    RuleType = "AT_RULE"    // @media { ... }
	STYLE_RULE RuleType = "STYLE_RULE" // div { ... }
	DECL_RULE  RuleType = "DECL_RULE"  // color: red;
)

type RuleNode struct {
	Loc  Location
	Type RuleType
	Rule interface{}
}

type AtRule struct {
	Name   Identifier
	Params []Value
	Rules  []RuleNode
}

type StyleRule struct {
	Selectors []Selector
	Rules     []RuleNode
}

type SelectorType string

const (
	TYPE_SELECTOR      SelectorType = "TYPE_SELECTOR"      // div
	CLASS_SELECTOR     SelectorType = "CLASS_SELECTOR"     // .foo
	ID_SELECTOR        SelectorType = "ID_SELECTOR"        // #foo
	ATTRIBUTE_SELECTOR SelectorType = "ATTRIBUTE_SELECTOR" // [foo]
	PSEUDO_SELECTOR    SelectorType = "PSEUDO_SELECTOR"    // :foo
)

type Selector struct {
	Loc        Location
	Type       SelectorType
	Parts      []SelectorPart
	Attributes map[string]Value
	Combinator string
	Next       *Selector
}

func (s Selector) String() string {
	format := "S("
	args := []interface{}{}
	first := true
	for _, part := range s.Parts {
		if first {
			first = false
		} else {
			format += "_"
		}
		format += "%s"
		args = append(args, part)
	}

	if s.Next != nil {
		format += " %s %s"
		args = append(args, s.Combinator, *s.Next)
	}
	format += ")"
	return fmt.Sprintf(format, args...)
}

type SelectorPart struct {
	Loc   Location
	Type  SelectorType
	Value Value
}

func (s SelectorPart) String() string {
	format := ""
	switch s.Type {
	case TYPE_SELECTOR:
		format += "%s"
	case CLASS_SELECTOR:
		format += ".%s"
	case ID_SELECTOR:
		format += "#%s"
	case ATTRIBUTE_SELECTOR:
		format += "[%s]"
	case PSEUDO_SELECTOR:
		format += ":%s"
	}
	return fmt.Sprintf(format, s.Value.Data.(Identifier).Value)
}

type Declaration struct {
	Name  Identifier
	Value Value
}
