package interpreter

import (
	"fmt"

	"github.com/shreyassanthu77/cisp/ast"
)

type ReturnValue struct {
	Value ast.Value
}

func (r ReturnValue) IsValue() {}

func isReturnValue(v ast.Value) bool {
	_, ok := v.(ReturnValue)
	return ok
}

func evalIfRule(env *Environment, rule ast.AtRule) (ast.Value, error) {
	if len(rule.Parameters) != 1 {
		return ast.NilValue{}, fmt.Errorf("if rules should have exactly one parameter")
	}

	condition, err := evalValue(rule.Parameters[0], env)
	if err != nil {
		return ast.NilValue{}, err
	}

	conditionResult, ok := condition.(ast.Boolean)
	if !ok {
		return ast.NilValue{}, fmt.Errorf("if rule condition must evaluate to a boolean")
	}

	if conditionResult.Value {
		return evalStatementList(rule.Body, env)
	}

	return ast.NilValue{}, nil
}

func evalReturnRule(env *Environment, rule ast.AtRule) (ast.Value, error) {
	if len(rule.Parameters) != 1 {
		return ast.NilValue{}, fmt.Errorf("return rules should have exactly one parameter")
	}

	value, err := evalValue(rule.Parameters[0], env)
	if err != nil {
		return ast.NilValue{}, err
	}

	return ReturnValue{Value: value}, nil
}

func evalAtRule(env *Environment, rule ast.AtRule) (ast.Value, error) {
	switch rule.Name {
	case "if":
		return evalIfRule(env, rule)
	case "return":
		return evalReturnRule(env, rule)
	}

	return ast.NilValue{}, fmt.Errorf("at rules are not supported yet")
}
