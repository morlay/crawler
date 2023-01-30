package directives

import (
	"context"

	"github.com/morlay/crawler/pkg/directives/core"
)

func init() {
	core.RegisterPipeDirective(&Def{})
}

type Def struct {
	core.DirectiveName `name:"_def" on:"FIELD,repeatable"`
	Var                string `json:"var"`
}

func (d *Def) Execute(ctx context.Context, input any) (any, error) {
	f := core.FlowStateFromContext(ctx)
	f.Store(d.Var, input)
	rootInput, _ := f.Load("$")
	return rootInput, nil
}
