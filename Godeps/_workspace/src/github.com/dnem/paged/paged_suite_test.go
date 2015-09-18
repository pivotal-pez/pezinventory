package paged_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unrolled/render"

	"testing"
)

func TestPaged(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Paged Test Suite")
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
