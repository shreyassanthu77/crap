package interpreter

import "github.com/shreyassanthu77/cisp/ast"

func (i *Interpreter) applyArgs(rule ast.Rule, args []ast.Value) {
	attributes := rule.Selector.Atrributes
	params := args
	if len(attributes) != len(params) {
		i.throwAtSpan(rule.Span, "expected %d parameters, got %d", len(attributes), len(params))
	}

	for idx, attr := range attributes {
		param := params[idx]
		if isNilValue(param) {
			if attr.Default != nil {
				param = attr.Default
			} else {
				i.throwAtSpan(args[idx].GetSpan(), "parameter %s is required", attr.Name.Name)
			}
		}
		i.env.setVar(attr.Name.Name, param)
	}

}

func (i *Interpreter) evalRule(rule ast.Rule, args []ast.Value) ast.Value {
	return i.fork(func() ast.Value {
		return ast.NilValue{}
	})
}
