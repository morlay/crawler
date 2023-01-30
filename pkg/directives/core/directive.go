package core

import (
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
)

type Directive interface {
	directiveName() *DirectiveName
}

type DirectiveName struct {
	name string
}

func (n *DirectiveName) directiveName() *DirectiveName {
	return n
}

type directive struct {
	name  string
	rtype reflect.Type
	gql   *graphql.Directive
}

func (d *directive) New(params map[string]any) (Directive, error) {
	inst := reflect.New(d.rtype).Interface()

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     inst,
		DecodeHook: mapstructure.TextUnmarshallerHookFunc(),
	})
	if err != nil {
		return nil, err
	}

	if err := dec.Decode(params); err != nil {
		return nil, err
	}

	dt := inst.(Directive)
	dt.directiveName().name = d.name
	return dt, nil
}
