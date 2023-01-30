package directives

import (
	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

	"github.com/morlay/crawler/pkg/directives/core"
	"github.com/morlay/crawler/pkg/directives/strfmt"
)

func init() {
	core.RegisterPipeDirective(&DomCloset{})
}

type DomCloset struct {
	core.DirectiveName `name:"_dom_closet" on:"FIELD,repeatable"`
	Select             *strfmt.StringOrTemplate `json:"select"`
}

func (d *DomCloset) Execute(ctx context.Context, input any) (any, error) {
	switch x := input.(type) {
	case *goquery.Selection:
		sel, err := mayResolve(ctx, d.Select)
		if err != nil {
			return nil, err
		}
		return x.Closest(sel), nil
	}
	return nil, errors.Wrap(core.InvalidInput, "should input with *goquery.Selection")
}
