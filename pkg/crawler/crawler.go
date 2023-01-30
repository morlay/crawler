package crawler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	_ "github.com/morlay/crawler/pkg/directives"
)

type Crawler interface {
	Source() *Source
	Do(ctx context.Context, action string, params map[string]string) (Result, error)
}

type Result interface {
	Scan(ctx context.Context, target any) error
}

type crawler struct {
	source *Source
	schema graphql.Schema
	info   *graphql.ResolveInfo
}

func (c *crawler) Source() *Source {
	return c.source
}

var defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

func (c *crawler) Do(ctx context.Context, action string, params map[string]string) (Result, error) {
	p, ok := c.source.Operations[action]
	if !ok {
		return nil, errors.Errorf("operation `%s` not found", action)
	}

	u, err := p.RequestURI(c.source, params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", defaultUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		d, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, err
		}
		data, err := p.Response.Extract(d.Selection)
		if err != nil {
			return nil, err
		}

		return &result{data}, nil
	}
	return nil, errors.Errorf("%s: request failed", resp.Status)
}

type result struct {
	data any
}

func (r *result) Scan(ctx context.Context, target any) error {
	switch x := target.(type) {
	case io.Writer:
		d := json.NewEncoder(x)
		d.SetIndent("", "  ")
		return d.Encode(r.data)
	default:
		return mapstructure.Decode(r.data, target)
	}
}
