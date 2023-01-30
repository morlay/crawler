package crawler

import (
	"net/url"

	"github.com/morlay/crawler/pkg/directives/strfmt"
)

type Source struct {
	Name       string
	Origin     string
	Operations map[string]Operation
}

type Operation struct {
	Path       *strfmt.StringOrTemplate
	Parameters map[string]Parameter
	Response   Selection
}

func (o *Operation) RequestURI(source *Source, params map[string]string) (*url.URL, error) {
	paramsInPath := map[string]any{}
	query := url.Values{}

	for name, p := range o.Parameters {
		v, ok := params[name]

		if !ok && p.Default != nil {
			v = *p.Default
		}

		if p.In == "query" {
			query.Add(name, v)
			continue
		}

		paramsInPath[name] = v
	}

	p, err := o.Path.Execute(func(key string) (any, error) {
		v, ok := paramsInPath[key]
		if ok {
			return v, nil
		}
		return "", nil
	})

	u, err := url.Parse(fixURI(string(p), source.Origin))
	if err != nil {
		return nil, err
	}

	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	return u, nil
}

type Parameter struct {
	In      string
	Default *string
	Enum    map[string]string
}
