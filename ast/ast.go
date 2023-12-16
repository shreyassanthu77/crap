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
	Loc     Location
	Rules   []Rule   // content: "Hello World";
	AtRules []AtRule // @media { ... }
}

func (s Stylesheet) String() string {
	format := ""
	args := []interface{}{}

	format += "AtRules {"
	for _, atRule := range s.AtRules {
		format += "\n%s"
		args = append(args, atRule)
	}
	if len(s.AtRules) > 0 {
		format += "\n"
	}

	format += "}\nRules {"
	for _, rule := range s.Rules {
		format += "\n%s"
		args = append(args, rule)
	}
	if len(s.Rules) > 0 {
		format += "\n"
	}
	format += "}"

	return fmt.Sprintf(format, args...)
}

type Rule struct {
	Loc          Location
	Selectors    []Selector    // ".a", "div", "a[href="/"]" etc.
	Declarations []Declaration // "content: "Hello World";"
	Rules        []Rule        // &.a { ... }, @media { ... }
}

func (r Rule) String() string {
	return fmt.Sprintf(`Rule{
Selectors: %v,
Declarations: %v,
}`, r.Selectors, r.Declarations)
}

type Selector struct {
	Loc        Location
	Elements   []SelectorElement    // div > a[href="/"] etc.
	Attributes map[Identifier]Value // [method="GET"] etc.
}

func (s Selector) String() string {
	return fmt.Sprintf("%v %v", s.Elements, s.Attributes)
}

type SelectorElement struct {
	Loc        Location
	Combinator string           // " ", ">", "+", "~"
	Identifier Identifier       // "div", "a", "p" etc.
	Next       *SelectorElement // "div > a" -> "a"
}

func (s SelectorElement) String() string {
	return fmt.Sprintf("%s %s %v", s.Identifier.Value, s.Combinator, *s.Next)
}

type Declaration struct {
	Loc      Location
	Property Identifier // "content"
	Value    Value      // "Hello World"
}

func (d Declaration) String() string {
	return fmt.Sprintf("%s: %v;", d.Property.Value, d.Value)
}

type AtRule struct {
	Loc   Location
	Name  Identifier // "import", "media" etc.
	Param Value      // "./other.css", "screen" etc.
	Rules []Rule     // @media { ... }
}

func (r AtRule) String() string {
	format := "@%s %s {"
	args := []interface{}{r.Name.Value, r.Param}
	for _, rule := range r.Rules {
		format += "\n\t%s"
		args = append(args, rule)
	}
	if len(r.Rules) > 0 {
		format += "\n"
	}
	format += "}"
	return fmt.Sprintf(format, args...)
}
