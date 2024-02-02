package interpreter

import "github.com/shreyassanthu77/cisp/ast"

func evalBinaryOp(op ast.BinaryOp, env *Environment) (ast.Value, error) {
	// left, err := evalValue(op.Left, env)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// right, err := evalValue(op.Right, env)
	// if err != nil {
	// 	return nil, err
	// }
	//

	switch op.Op {
	default:
		panic("Not implemented")
	}
}
