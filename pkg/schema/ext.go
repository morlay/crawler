package schema

import (
	"context"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

type DiscardExt struct {
}

func (e *DiscardExt) Init(ctx context.Context, params *graphql.Params) context.Context {
	return ctx
}

func (e *DiscardExt) Name() string {
	return ""
}

func (e *DiscardExt) ExecutionDidStart(ctx context.Context) (context.Context, graphql.ExecutionFinishFunc) {
	return ctx, func(result *graphql.Result) {}
}

func (e *DiscardExt) ResolveFieldDidStart(ctx context.Context, info *graphql.ResolveInfo) (context.Context, graphql.ResolveFieldFinishFunc) {
	return ctx, func(i interface{}, err error) {}
}

func (e *DiscardExt) ParseDidStart(ctx context.Context) (context.Context, graphql.ParseFinishFunc) {
	return ctx, func(err error) {}
}

func (e *DiscardExt) ValidationDidStart(ctx context.Context) (context.Context, graphql.ValidationFinishFunc) {
	return ctx, func(errors []gqlerrors.FormattedError) {}
}

func (e *DiscardExt) HasResult() bool {
	return false
}

func (e *DiscardExt) GetResult(ctx context.Context) interface{} {
	return nil
}
