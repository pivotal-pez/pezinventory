package wrapper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unrolled/render"

	"testing"
)

func TestWrapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cfmgo/wrapper test suite")
}

var formatter *render.Render

//Formatter returns the address for a global response formatter
//realized in the `github.com/unrolled/render` package.
func Formatter() *render.Render {
	if formatter == nil {
		formatter = render.New(render.Options{
			IndentJSON: true,
		})
	}
	return formatter
}
