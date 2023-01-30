package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/morlay/crawler/pkg/directives/core"
	"github.com/pkg/errors"
)

func Load(src []byte) (*graphql.Schema, *ast.Document, error) {
	params := parser.ParseParams{
		Source: &source.Source{
			Name: "Query",
			Body: src,
		},
		Options: parser.ParseOptions{
			NoSource:   true,
			NoLocation: false,
		},
	}

	doc, err := parser.Parse(params)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parse failed")
	}

	r := newRegister(doc)
	q := r.NamedType("Query")

	schemaConfig := graphql.SchemaConfig{
		Query:      q.(*graphql.Object),
		Directives: core.RegisteredDirectives(),
	}

	for n := range r.types {
		t := r.types[n]
		schemaConfig.Types = append(schemaConfig.Types, t)
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, nil, err
	}

	for _, n := range doc.Definitions {
		switch x := n.(type) {
		case *ast.OperationDefinition:
			return &schema, &ast.Document{
				Loc:         x.Loc,
				Definitions: []ast.Node{x},
			}, nil
		}
	}

	return &schema, nil, nil
}
