package directives

import (
	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

	"github.com/morlay/crawler/pkg/directives/core"
	"github.com/morlay/crawler/pkg/directives/strfmt"
)

func init() {
	core.RegisterPipeDirective(&DomQuery{})
}

type DomQuery struct {
	core.DirectiveName `name:"_dom_query" on:"FIELD,repeatable"`

	Select *strfmt.StringOrTemplate `json:"select"`
}

func (d *DomQuery) Execute(ctx context.Context, input any) (any, error) {
	switch x := input.(type) {
	case *goquery.Selection:
		sel, err := mayResolve(ctx, d.Select)
		if err != nil {
			return nil, err
		}
		return x.Find(sel), nil
	}
	return nil, errors.Wrap(core.InvalidInput, "should input with *goquery.Selection")
}
