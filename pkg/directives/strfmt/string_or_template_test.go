package strfmt

import (
	"testing"

	testingx "github.com/octohelm/x/testing"
)

func TestStringOrTemplate(t *testing.T) {
	st := &StringOrTemplate{}
	if err := st.UnmarshalText([]byte("a[href='#{id}']")); err != nil {
		t.Fatal(err)
	}

	ret, err := st.Execute(func(id string) (any, error) {
		return "book", nil
	})

	testingx.Expect(t, err, testingx.Be[error](nil))
	testingx.Expect(t, string(ret), testingx.Be("a[href='#book']"))
}
