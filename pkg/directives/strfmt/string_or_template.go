package strfmt

import (
	"bytes"
	"regexp"

	"github.com/octohelm/x/encoding"
)

type StringOrTemplate struct {
	parts    [][]byte
	paramIdx map[string]int
}

var re = regexp.MustCompile("({[^}]+})")

func (t *StringOrTemplate) IsTemplate() bool {
	return t.paramIdx != nil
}

func (t *StringOrTemplate) Execute(resolve func(key string) (any, error)) ([]byte, error) {
	if t.IsTemplate() {
		n := &StringOrTemplate{parts: t.parts[:]}

		for name := range t.paramIdx {
			v, err := resolve(name)
			if err != nil {
				return nil, err
			}

			p, err := encoding.MarshalText(v)
			if err != nil {
				return nil, err
			}
			n.parts[t.paramIdx[name]] = p
		}

		return n.MarshalText()
	}
	return t.MarshalText()
}

func (t *StringOrTemplate) UnmarshalText(text []byte) error {
	if re.Match(text) {
		holders := re.FindAllIndex(text, -1)

		idx := 0

		t.paramIdx = map[string]int{}

		appendPart := func(p []byte, isVar bool) {
			t.parts = append(t.parts, p)
			if isVar {
				t.paramIdx[string(bytes.TrimSpace(p[1:len(p)-1]))] = idx
			}
			idx++
		}

		l := 0
		for _, rng := range holders {
			varName := text[rng[0]:rng[1]]

			appendPart(text[l:rng[0]], false)
			appendPart(varName, true)
			l = rng[1]
		}

		if l < len(text) {
			appendPart(text[l:], false)
		}

		return nil
	}
	t.parts = [][]byte{text}
	return nil
}

func (t *StringOrTemplate) String() string {
	return string(bytes.Join(t.parts, []byte("")))
}

func (t StringOrTemplate) MarshalText() ([]byte, error) {
	return bytes.Join(t.parts, []byte("")), nil
}
