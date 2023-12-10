package ast

type Stylesheet struct {
	Rules   []Rule
	Imports []string
}

type Rule struct {
	Selectors    []Selector
	Declarations []Declaration
}

type Selector struct {
	Tag        string
	ID         string
	Class      string
	Atrributes map[string]string
}

type Declaration struct {
	Property string
	Value    Value
}

type ValueType int

const (
	NUMBER ValueType = iota
	STRING
	FUNCTION_CALL
)

type Value struct {
	Type ValueType
	data interface{}
}

type FunctionCall struct {
	Name string
	Args []Value
}

func NumberValue(n float64) Value {
	return Value{NUMBER, n}
}

func StringValue(s string) Value {
	return Value{STRING, s}
}

func FunctionCallValue(name string, args []Value) Value {
	return Value{FUNCTION_CALL, FunctionCall{name, args}}
}

func (v Value) Num() float64 {
	if v.Type != NUMBER {
		panic("not a number")
	}
	return v.data.(float64)
}

func (v Value) Str() string {
	if v.Type != STRING {
		panic("not a string")
	}
	return v.data.(string)
}

func (v Value) Fn() FunctionCall {
	if v.Type != FUNCTION_CALL {
		panic("not a function call")
	}
	return v.data.(FunctionCall)
}

type ValueMapper struct {
	Number       func(float64) interface{}
	String       func(string) interface{}
	FunctionCall func(FunctionCall) interface{}
}

func (v Value) Map(
	mapper ValueMapper,
) interface{} {
	switch v.Type {
	case NUMBER:
		return mapper.Number(v.data.(float64))
	case STRING:
		return mapper.String(v.data.(string))
	case FUNCTION_CALL:
		return mapper.FunctionCall(v.data.(FunctionCall))
	}

	panic("unreachable")
}
