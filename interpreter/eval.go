package interpreter

import (
	"fmt"

	"github.com/shreyassanthu77/cisp/ast"
)

func evalVarDeclaration(decl ast.Declaration, env *Environment) (ast.Value, error) {
	name := decl.Property.Name[2:] // remove the -- from the name

	if len(decl.Parameters) != 1 {
		return nil, fmt.Errorf("variable declaration should have exactly one value")
	}
	value, err := evalValue(decl.Parameters[0], env)

	val, err := evalValue(value, env)
	if err != nil {
		return nil, err
	}
	env.setVar(name, val)
	return val, nil
}

func evalStmt(stmt ast.Statement, env *Environment) (ast.Value, error) {
	switch stmt := stmt.(type) {
	case ast.Rule:
		err := env.setFn(stmt)
		if err != nil {
			return nil, err
		}
		return ast.NilValue{}, nil
	case ast.AtRule:
		panic("at-rules not supported yet")
	case ast.Declaration:
		if len(stmt.Property.Name) > 2 && stmt.Property.Name[:2] == "--" {
			return evalVarDeclaration(stmt, env)
		}
		fnCall := ast.FunctionCall{
			Fn:         stmt.Property,
			Parameters: stmt.Parameters,
		}
		return evalFnCall(fnCall, env)
	case NativeFnCall:
		return stmt.Handler(env)
	}

	return nil, nil
}

func verifyAndAddParamsToEnv(attributes []ast.Attreibute, params []ast.Value, env *Environment) error {
	if len(attributes) != len(params) {
		return fmt.Errorf("expected %d parameters, got %d", len(attributes), len(params))
	}

	for i, attr := range attributes {
		env.Vars[attr.Name.Name] = params[i]
	}

	return nil
}

func evalRule(rule ast.Rule, params []ast.Value, parent *Environment) (ast.Value, error) {
	env := parent.fork()
	err := verifyAndAddParamsToEnv(rule.Selector.Atrributes, params, env)
	if err != nil {
		return nil, err
	}

	var res ast.Value = ast.NilValue{}
	for _, stmt := range rule.Body {
		_res, err := evalStmt(stmt, env)
		if err != nil {
			return nil, err
		}
		res = _res
	}
	return res, nil
}
