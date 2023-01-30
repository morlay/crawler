package core

import (
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
)

func toGqlDirective(t reflect.Type) *graphql.Directive {
	dc := graphql.DirectiveConfig{
		Args: map[string]*graphql.ArgumentConfig{},
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Name == "DirectiveName" {
			if n, ok := f.Tag.Lookup("name"); ok {
				dc.Name = n
			}

			if location, ok := f.Tag.Lookup("on"); ok {
				locations := strings.Split(location, ",")
				for _, l := range locations {
					if l != "repeatable" {
						dc.Locations = append(dc.Locations, l)
					} else {
						dc.Description = l
					}
				}
			}

			continue
		}

		if argTag, ok := f.Tag.Lookup("json"); ok {
			parts := strings.Split(argTag, ",")
			name := parts[0]
			if name == "" {
				name = f.Name
			}

			gglType := fromGoType(f.Type)
			if !strings.Contains(argTag, "omitempty") {
				dc.Args[name] = &graphql.ArgumentConfig{Type: graphql.NewNonNull(gglType)}
			} else {
				dc.Args[name] = &graphql.ArgumentConfig{Type: gglType}
			}
		}
	}

	return graphql.NewDirective(dc)
}
