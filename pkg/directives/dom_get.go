package directives

import (
	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

	"github.com/morlay/crawler/pkg/directives/core"
)

func init() {
	core.RegisterPipeDirective(&DomGet{})
}

type DomGet struct {
	core.DirectiveName `name:"_dom_get" on:"FIELD,repeatable"`

	Idx int `json:"idx,omitempty"`
}

func (fn *DomGet) Execute(ctx context.Context, input any) (any, error) {
	switch x := input.(type) {
	case *goquery.Selection:
		return x.Eq(fn.Idx), nil
	}
	return nil, errors.Wrap(core.InvalidInput, "should input with *goquery.Selection")
}
