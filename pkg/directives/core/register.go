package core

import (
	"fmt"
	"io"
	"reflect"
	"sort"

	"github.com/pkg/errors"

	"github.com/graphql-go/graphql"
)

var directives = map[string]*directive{}

func New(name string, params map[string]any) (Directive, error) {
	if d, ok := directives[name]; ok {
		return d.New(params)
	}
	return nil, errors.Errorf("directive %q is not defined", name)
}

func RegisteredDirectives() []*graphql.Directive {
	names := make([]string, 0)
	for n := range directives {
		names = append(names, n)
	}
	sort.Strings(names)

	list := make([]*graphql.Directive, len(names))

	for i := range list {
		list[i] = directives[names[i]].gql
	}
	return list
}

func PrintRegisteredDirectives(w io.Writer) {
	for _, d := range RegisteredDirectives() {
		_, _ = fmt.Fprintf(w, `directive @%s(`, d.Name)

		for i := range d.Args {
			if i > 0 {
				_, _ = fmt.Fprintf(w, ", ")
			}
			a := d.Args[i]
			_, _ = fmt.Fprintf(w, "%s: %s", a.Name(), a.Type)
		}

		_, _ = fmt.Fprintf(w, ")")

		if d.Description == "repeatable" {
			_, _ = fmt.Fprintf(w, " repeatable")
		}

		_, _ = fmt.Fprintf(w, " on ")

		for i, l := range d.Locations {
			if i > 0 {
				_, _ = fmt.Fprintf(w, " | ")
			}
			_, _ = fmt.Fprintf(w, "%s", l)
		}

		_, _ = fmt.Fprintf(w, "\n")
	}
}

func RegisterPipeDirective(d PipeDirective) {
	Register(d)
}

func Register(d Directive) {
	t := reflect.TypeOf(d)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	dd := &directive{}
	dd.rtype = t
	dd.gql = toGqlDirective(dd.rtype)
	dd.name = dd.gql.Name
	directives[dd.name] = dd
}

func fromGoType(t reflect.Type) graphql.Type {
	switch t.Kind() {
	case reflect.Slice:
		return graphql.NewList(fromGoType(t.Elem()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return graphql.Int
	case reflect.Float32, reflect.Float64:
		return graphql.Float
	case reflect.Bool:
		return graphql.Boolean
	}
	return graphql.String
}
