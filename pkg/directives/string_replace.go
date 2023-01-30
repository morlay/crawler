package directives

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/morlay/crawler/pkg/directives/core"
	"github.com/morlay/crawler/pkg/directives/strfmt"
)

func init() {
	core.RegisterPipeDirective(&StringReplace{})
}

type StringReplace struct {
	core.DirectiveName `name:"_string_replace" on:"FIELD,repeatable"`

	From  *strfmt.StringOrTemplate `json:"from,omitempty"`
	FromR *strfmt.RegexString      `json:"fromR,omitempty"`
	To    *strfmt.StringOrTemplate `json:"to"`
}

func (s *StringReplace) Execute(ctx context.Context, input any) (any, error) {
	switch x := input.(type) {
	case string:
		to, err := mayResolve(ctx, s.To)
		if err != nil {
			return nil, err
		}
		if s.FromR != nil {
			return s.FromR.ReplaceAll(x, to), nil
		}
		from, err := mayResolve(ctx, s.From)
		if err != nil {
			return nil, err
		}
		return strings.ReplaceAll(x, from, to), nil
	}
	return nil, errors.Wrap(core.InvalidInput, "should input with value as type string")
}
