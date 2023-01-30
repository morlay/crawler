package core

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

var InvalidInput = errors.New("Invalid input")

type PipeDirective interface {
	Directive

	Execute(ctx context.Context, input any) (any, error)
}

func Pipe(ctx context.Context, input any, directives ...PipeDirective) (any, error) {
	final := input

	fState := &flowState{}

	// store root input for multi-step pipe
	fState.Store("$", input)

	ctx = ContextWithFlowContext(ctx, fState)

	for i := range directives {
		ret, err := directives[i].Execute(ctx, final)
		if err != nil {
			return nil, err
		}
		final = ret
	}

	return final, nil
}

type FlowState interface {
	Load(key string) (any, bool)
	Store(key string, value any)
}

type flowStateContext struct{}

func FlowStateFromContext(ctx context.Context) FlowState {
	return ctx.Value(flowStateContext{}).(FlowState)
}

func ContextWithFlowContext(ctx context.Context, flowState FlowState) context.Context {
	return context.WithValue(ctx, flowStateContext{}, flowState)
}

type flowState struct {
	store sync.Map
}

func (f *flowState) Load(key string) (any, bool) {
	return f.store.Load(key)
}

func (f *flowState) Store(key string, value any) {
	f.store.Store(key, value)
}
