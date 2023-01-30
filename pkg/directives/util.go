package directives

import (
	"context"

	"github.com/pkg/errors"

	"github.com/morlay/crawler/pkg/directives/core"
	"github.com/morlay/crawler/pkg/directives/strfmt"
)

func mayResolve(ctx context.Context, str *strfmt.StringOrTemplate) (string, error) {
	if str.IsTemplate() {
		f := core.FlowStateFromContext(ctx)

		b, err := str.Execute(func(key string) (any, error) {
			v, ok := f.Load(key)
			if !ok {
				return nil, errors.Errorf("%s is undefined", key)
			}
			return v, nil
		})

		return string(b), err
	}
	return str.String(), nil
}
