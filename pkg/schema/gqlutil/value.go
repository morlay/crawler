package gqlutil

import (
	"strconv"

	"github.com/graphql-go/graphql/language/ast"
)

func GoValueFromAST(v ast.Value) any {
	switch x := v.(type) {
	case *ast.ObjectValue:
		values := map[string]any{}
		for _, f := range x.Fields {
			values[f.Name.Value] = GoValueFromAST(f.Value)
		}
		return values
	case *ast.ListValue:
		values := make([]any, len(x.Values))
		for i := range x.Values {
			values[i] = GoValueFromAST(x.Values[i])
		}
		return values
	case *ast.IntValue:
		ret, _ := strconv.ParseInt(x.Value, 10, 64)
		return ret
	case *ast.FloatValue:
		ret, _ := strconv.ParseFloat(x.Value, 64)
		return ret
	case *ast.StringValue:
		return x.Value
	case *ast.BooleanValue:
		return x.Value
	}
	return nil
}
