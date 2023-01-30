package directives

import (
	"bytes"
	"os"
	"testing"

	"github.com/morlay/crawler/pkg/directives/core"
)

func TestFormatDirective(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	core.PrintRegisteredDirectives(buf)
	_ = os.WriteFile("directives.gql", buf.Bytes(), 0755)
}
