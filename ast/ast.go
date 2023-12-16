package ast

type Location struct {
	Start Span
	End   Span
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

type Expression struct {
	Loc   Location
	Left  Value
	Op    string
	Right Value
}

type Value struct {
	Loc  Location
	Type ValueType
	Data interface{}
}

type Stylesheet struct {
	Loc     Location
	Rules   []Rule   // content: "Hello World";
	AtRules []AtRule // @media { ... }
}

type Rule struct {
	Loc          Location
	Selectors    []Selector    // ".a", "div", "a[href="/"]" etc.
	Declarations []Declaration // "content: "Hello World";"
	Rules        []Rule        // &.a { ... }, @media { ... }
}

type Selector struct {
	Loc        Location
	Elements   []SelectorElement    // div > a[href="/"] etc.
	Attributes map[Identifier]Value // [method="GET"] etc.
}

type SelectorElement struct {
	Loc        Location
	Combinator string           // " ", ">", "+", "~"
	Identifier Identifier       // "div", "a", "p" etc.
	Next       *SelectorElement // "div > a" -> "a"
}

type Declaration struct {
	Loc      Location
	Property Identifier // "content"
	Value    Value      // "Hello World"
}

type AtRule struct {
	Loc    Location
	Name   Identifier // "import", "media" etc.
	Params []Value    // "./other.css", "screen" etc.
	Rules  []Rule     // @media { ... }
}
