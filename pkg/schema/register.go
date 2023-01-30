package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/morlay/crawler/pkg/schema/gqlutil"
)

func newRegister(doc *ast.Document) *register {
	return &register{doc: doc, types: map[string]graphql.Type{}}
}

type register struct {
	doc   *ast.Document
	types map[string]graphql.Type
}

func (r *register) NamedType(name string) (tpe graphql.Type) {
	if t, ok := r.types[name]; ok {
		return t
	}

	defer func() {
		r.types[name] = tpe
	}()

	for i := range r.doc.Definitions {
		switch x := r.doc.Definitions[i].(type) {
		case *ast.ObjectDefinition:
			if name == x.Name.Value {
				return r.TypeFromAST(x)
			}
		case *ast.InputObjectDefinition:
			if name == x.Name.Value {
				return r.TypeFromAST(x)
			}
		case *ast.EnumDefinition:
			if name == x.Name.Value {
				return r.TypeFromAST(x)
			}
		}
	}

	return
}

func (r *register) TypeFromAST(node ast.Node) graphql.Type {
	switch x := node.(type) {
	case *ast.NonNull:
		return graphql.NewNonNull(r.TypeFromAST(x.Type))
	case *ast.Named:
		switch typeOrRef := x.Name.Value; typeOrRef {
		case "ID":
			return graphql.ID
		case "String":
			return graphql.String
		case "Boolean":
			return graphql.Boolean
		case "Int":
			return graphql.Int
		case "Float":
			return graphql.String
		default:
			return r.NamedType(typeOrRef)
		}
	case *ast.List:
		return graphql.NewList(r.TypeFromAST(x.Type))
	case *ast.ObjectDefinition:
		fields := graphql.Fields{}

		for _, f := range x.Fields {
			field := &graphql.Field{
				Name: f.Name.Value,
				Type: r.TypeFromAST(f.Type),
			}

			if len(f.Arguments) > 0 {
				field.Args = graphql.FieldConfigArgument{}

				for _, arg := range f.Arguments {
					field.Args[arg.Name.Value] = &graphql.ArgumentConfig{
						Type:         r.TypeFromAST(arg.Type),
						DefaultValue: gqlutil.GoValueFromAST(arg.DefaultValue),
					}
				}
			}

			fields[field.Name] = field
		}

		return graphql.NewObject(graphql.ObjectConfig{
			Name:   x.Name.Value,
			Fields: fields,
		})
	case *ast.InputObjectDefinition:
		inputs := graphql.InputObjectConfigFieldMap{}

		for _, f := range x.Fields {
			inputs[f.Name.Value] = &graphql.InputObjectFieldConfig{
				Type:         r.TypeFromAST(f.Type),
				DefaultValue: gqlutil.GoValueFromAST(f.DefaultValue),
			}
		}

		return graphql.NewInputObject(graphql.InputObjectConfig{
			Name:   x.Name.Value,
			Fields: inputs,
		})
	case *ast.EnumDefinition:
		enums := graphql.EnumValueConfigMap{}

		for _, v := range x.Values {
			e := &graphql.EnumValueConfig{
				Value: v.Name.Value,
			}
			enums[v.Name.Value] = e
		}
		return graphql.NewEnum(graphql.EnumConfig{
			Name:   x.Name.Value,
			Values: enums,
		})
	default:

	}

	return nil
}
