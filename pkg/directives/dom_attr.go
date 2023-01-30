package directives

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/morlay/crawler/pkg/directives/core"
	"github.com/pkg/errors"
)

func init() {
	core.RegisterPipeDirective(&DomAttr{})
}

type DomAttr struct {
	core.DirectiveName `name:"_dom_attr" on:"FIELD,repeatable"`

	Name string `json:"name"`
}

func (d *DomAttr) Execute(ctx context.Context, input any) (any, error) {
	switch x := input.(type) {
	case *goquery.Selection:
		if d.Name == "text" {
			return strings.TrimSpace(x.Text()), nil
		}
		return strings.TrimSpace(x.AttrOr(d.Name, "")), nil
	}
	return nil, errors.Wrap(core.InvalidInput, "should input with *goquery.Selection")
}
