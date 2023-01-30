package crawler

import (
	"context"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/stretchr/objx"

	"github.com/morlay/crawler/pkg/directives/strfmt"
	"github.com/morlay/crawler/pkg/schema"
	"github.com/morlay/crawler/pkg/schema/gqlutil"
)

func NewCrawler(src []byte) (Crawler, error) {
	s, requestAst, err := schema.Load(src)
	if err != nil {
		return nil, err
	}
	b := &crawlerParser{}

	s.AddExtensions(b)
	if ret := graphql.Execute(graphql.ExecuteParams{Schema: *s, AST: requestAst}); ret.HasErrors() {
		return nil, ret.Errors[0]
	}

	return &b.crawler, nil
}

type crawlerParser struct {
	schema.DiscardExt
	crawler
}

func (c *crawlerParser) Name() string {
	return "crawlerParser"
}

func (c *crawlerParser) ResolveFieldDidStart(ctx context.Context, info *graphql.ResolveInfo) (context.Context, graphql.ResolveFieldFinishFunc) {
	c.info = info

	selections := info.Operation.GetSelectionSet()

	for _, s := range selections.Selections {
		switch x := s.(type) {
		case *ast.Field:
			if x.Name.Value == "from" {
				source := &Source{}
				operationTypes := info.ReturnType.(*graphql.Object).Fields()

				for _, arg := range x.Arguments {
					switch arg.Name.Value {
					case "origin":
						source.Origin = gqlutil.GoValueFromAST(arg.Value).(string)
					case "name":
						source.Name = gqlutil.GoValueFromAST(arg.Value).(string)
					}
				}

				source.Operations = map[string]Operation{}

				for _, s := range s.GetSelectionSet().Selections {
					switch x := s.(type) {
					case *ast.Field:
						operationName := x.Name.Value

						operation := Operation{}

						for _, arg := range x.Arguments {
							switch arg.Name.Value {
							case "path":
								operation.Path = &strfmt.StringOrTemplate{}
								if err := operation.Path.UnmarshalText([]byte(gqlutil.GoValueFromAST(arg.Value).(string))); err != nil {
									panic(err)
								}
							case "parameters":
								operation.Parameters = map[string]Parameter{}

								for _, value := range gqlutil.GoValueFromAST(arg.Value).([]any) {
									p := objx.New(value)
									param := Parameter{}

									if v := p.Get("default"); !v.IsNil() {
										str := v.Str("")
										param.Default = &str
									}

									param.In = p.Get("in").Str("path")

									if p.Has("enum") {
										param.Enum = map[string]string{}

										p.Get("enum").EachObjxMap(func(i int, m objx.Map) bool {
											param.Enum[m.Get("v").Str()] = m.Get("n").Str()
											return true
										})
									}

									operation.Parameters[p.Get("name").Str()] = param
								}
							}
						}

						operation.Response = SelectionFor(operationTypes[operationName].Type, x.SelectionSet)

						source.Operations[operationName] = operation
					}
				}

				c.source = source
			}
		}
	}

	return ctx, func(i interface{}, err error) {}
}
