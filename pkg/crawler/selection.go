package crawler

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/morlay/crawler/pkg/directives/core"
	"github.com/morlay/crawler/pkg/schema"
	"github.com/morlay/crawler/pkg/schema/gqlutil"
)

func SelectionFor(gqlType graphql.Type, selectionSet *ast.SelectionSet, directives ...*ast.Directive) Selection {
	if selectionSet == nil {
		return &BasicSelection{
			Type:       gqlType.(*graphql.Scalar),
			Directives: SelectionDirectivesFromAST(directives...),
		}
	}

	switch t := gqlType.(type) {
	case *graphql.List:
		return &ListSelection{
			Type:       t,
			Item:       SelectionFor(t.OfType, selectionSet),
			Directives: SelectionDirectivesFromAST(directives...),
		}
	case *graphql.Object:
		o := &ObjectSelection{
			Type:   t,
			Fields: map[string]Selection{},
		}
		for _, f := range selectionSet.Selections {
			switch x := f.(type) {
			case *ast.Field:
				fieldName := x.Name.Value
				o.Fields[fieldName] = SelectionFor(
					t.Fields()[fieldName].Type,
					x.GetSelectionSet(), x.Directives...,
				)
			}
		}
		return o
	}

	return nil
}

type Selection interface {
	Extract(selection *goquery.Selection, keyPath ...any) (any, error)
}

func SelectionDirectivesFromAST(directives ...*ast.Directive) *SelectionDirectives {
	sd := &SelectionDirectives{}

	for i := range directives {
		d := directives[i]
		directiveName := d.Name.Value

		params := map[string]any{}

		for _, arg := range d.Arguments {
			paramName := arg.Name.Value
			params[paramName] = gqlutil.GoValueFromAST(arg.Value)
		}

		directive, err := core.New(directiveName, params)
		if err != nil {
			panic(err)
		}

		if pd, ok := directive.(core.PipeDirective); ok {
			sd.Pipes = append(sd.Pipes, pd)
		} else {
			sd.Metadata = append(sd.Metadata, directive)
		}
	}

	return sd
}

type SelectionDirectives struct {
	Metadata []core.Directive
	Pipes    []core.PipeDirective
}

type ListSelection struct {
	Type       *graphql.List
	Item       Selection
	Directives *SelectionDirectives
}

func (s *ListSelection) Extract(selection *goquery.Selection, keyPath ...any) (any, error) {
	if directives := s.Directives; directives != nil {
		ret, err := core.Pipe(context.Background(), selection, directives.Pipes...)
		if err != nil {
			return nil, err
		}

		switch x := ret.(type) {
		case *goquery.Selection:
			selection = x
		}
	}

	var err error

	list := make([]any, 0)

	selection.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		v, e := s.Item.Extract(selection, append(keyPath, i)...)
		if e != nil {
			err = e
			return false
		}
		list = append(list, v)
		return true
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

type ObjectSelection struct {
	Type       *graphql.Object
	Fields     map[string]Selection
	Directives *SelectionDirectives
}

func (s *ObjectSelection) Extract(selection *goquery.Selection, keyPath ...any) (any, error) {
	if directives := s.Directives; directives != nil {
		ret, err := core.Pipe(context.Background(), selection, directives.Pipes...)
		if err != nil {
			return nil, err
		}

		switch x := ret.(type) {
		case *goquery.Selection:
			selection = x
		}
	}

	o := map[string]any{}

	for k := range s.Fields {
		newKeyPath := append(keyPath, k)

		f := s.Fields[k]
		v, err := f.Extract(selection, newKeyPath...)
		if err != nil {
			return nil, err
		}
		o[k] = v
	}

	return o, nil
}

type BasicSelection struct {
	Type       *graphql.Scalar
	Directives *SelectionDirectives
}

func (s *BasicSelection) Extract(selection *goquery.Selection, keyPath ...any) (any, error) {
	var output any

	switch s.Type.Name() {
	case "Int":
		output = 0
	case "Float":
		output = float64(0)
	case "String", "ID":
		output = ""
	case "Boolean":
		output = false
	}

	if directives := s.Directives; directives != nil {
		ret, err := core.Pipe(context.Background(), selection, directives.Pipes...)
		if err != nil {
			return nil, errors.Wrapf(err, "%s：extract pipe execute failed", KeyPath(keyPath))
		}

		switch x := ret.(type) {
		case *goquery.Selection:
			if err := mapstructure.WeakDecode(strings.TrimSpace(x.Text()), &output); err != nil {
				return nil, errors.Wrapf(err, "%s：extract decode failed", KeyPath(keyPath))
			}
		default:
			if err := mapstructure.WeakDecode(x, &output); err != nil {
				return nil, errors.Wrapf(err, "%s：extract decode failed", KeyPath(keyPath))
			}
		}

		return output, nil
	}

	if err := mapstructure.Decode(selection.Text(), &output); err != nil {
		return nil, errors.Wrapf(err, "%s：extract decode failed", KeyPath(keyPath))
	}

	return output, nil
}

func SelectionFromGraphQL(src []byte) (Selection, error) {
	s, requestAst, err := schema.Load(src)
	if err != nil {
		return nil, errors.Wrap(err, "Load failed")
	}

	b := &selectionParser{}

	s.AddExtensions(b)
	if ret := graphql.Execute(graphql.ExecuteParams{Schema: *s, AST: requestAst}); ret.HasErrors() {
		return nil, ret.Errors[0]
	}

	return b.selection, nil
}

type selectionParser struct {
	schema.DiscardExt
	selection Selection
}

func (c *selectionParser) Name() string {
	return "selectionParser"
}

func (c *selectionParser) ResolveFieldDidStart(ctx context.Context, info *graphql.ResolveInfo) (context.Context, graphql.ResolveFieldFinishFunc) {
	c.selection = SelectionFor(
		info.Schema.QueryType(),
		info.Operation.GetSelectionSet(),
	)
	return ctx, func(i interface{}, err error) {}
}

type KeyPath []any

func (k KeyPath) String() string {
	b := &strings.Builder{}

	for i := range k {
		switch x := k[i].(type) {
		case int:
			_, _ = fmt.Fprintf(b, "[%d]", x)
		default:
			if i > 0 {
				b.WriteString(".")
			}
			_, _ = fmt.Fprintf(b, "%s", x)
		}

	}

	return b.String()
}
