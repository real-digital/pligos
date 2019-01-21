package pathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type config struct {
	A string
	B string `filepath:"resolve"`
	C subConfig
	D []subConfig
	E []string `filepath:"resolve"`
}

type subConfig struct {
	A string `filepath:"resolve"`
}

func Test_Resolve(t *testing.T) {
	root := "/home/gopher"

	input := config{
		A: "foobar",
		B: "go",
		C: subConfig{A: "documents"},
		D: []subConfig{{A: "movies"}, {A: "images"}},
		E: []string{"music"},
	}

	expected := config{
		A: "foobar",
		B: "/home/gopher/go",
		C: subConfig{A: "/home/gopher/documents"},
		D: []subConfig{{A: "/home/gopher/movies"}, {A: "/home/gopher/images"}},
		E: []string{"/home/gopher/music"},
	}

	Resolve(&input, root)

	assert.Equal(t, expected, input)

}
