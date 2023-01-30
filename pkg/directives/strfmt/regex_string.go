package strfmt

import "regexp"

type RegexString regexp.Regexp

func (r *RegexString) ReplaceAll(s string, to string) string {
	return (*regexp.Regexp)(r).ReplaceAllString(s, to)
}

func (r *RegexString) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	re, err := regexp.Compile(string(text))
	if err != nil {
		return err
	}
	*r = RegexString(*re)
	return nil
}

func (r RegexString) MarshalText() (text []byte, err error) {
	return []byte((*regexp.Regexp)(&r).String()), nil
}
