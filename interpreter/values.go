package interpreter

import (
	"fmt"

	"github.com/shreyassanthu77/cisp/ast"
)

func evalUnaryOp(op ast.UnaryOp, env *Environment) (ast.Value, error) {
	val, err := evalValue(op.Value, env)
	if err != nil {
		return nil, err
	}

	switch op.Op {
	case "+":
		return val, nil
	case "-":
		switch val := val.(type) {
		case ast.Int:
			return ast.Int{Value: -val.Value}, nil
		case ast.Float:
			return ast.Float{Value: -val.Value}, nil
		}
	case "!":
		switch val := val.(type) {
		case ast.Boolean:
			return ast.Boolean{Value: !val.Value}, nil
		default:
			return nil, fmt.Errorf("invalid type for unary operator !: %T", val)
		}
	}

	return nil, fmt.Errorf("invalid unary operator %s", op.Op)
}

func evalValue(value ast.Value, env *Environment) (ast.Value, error) {
	switch value := value.(type) {
	case ast.FunctionCall:
		return evalFnCall(value, env)
	case ast.Identifier:
		_, err := env.getVar(value.Name)
		if err != nil {
			_, err := env.genFn(value.Name)
			if err == nil {
				return nil, fmt.Errorf("You cannot use a function as a value use %s() instead of %s if you want to call it", value.Name, value.Name)
			}
			return nil, fmt.Errorf("Literal Identifiers are not allowed use $variable if you want to use a variable")
		}
		return nil, fmt.Errorf("Literal Identifiers are not allowed use $%s instead of %s", value.Name, value.Name)
	case ast.Int, ast.Float, ast.String, ast.Boolean, ast.NilValue:
		return value, nil
	case ast.UnaryOp:
		return evalUnaryOp(value, env)
	case ast.BinaryOp:
		return evalBinaryOp(value, env)
	case ast.VarianleDerefValue:
		val, err := env.getVar(value.Variable.Name)
		if err != nil {
			return nil, err
		}
		return val, nil
	}

	return value, nil
}

func evalFnCall(fnCall ast.FunctionCall, env *Environment) (ast.Value, error) {
	fn, err := env.genFn(fnCall.Fn.Name)
	if err != nil {
		return nil, err
	}

	for i, param := range fnCall.Parameters {
		fnCall.Parameters[i], err = evalValue(param, env)
		if err != nil {
			return nil, err
		}
	}

	return evalRule(fn, fnCall.Parameters, env)
}
