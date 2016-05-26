package farsight

import (
	// Standard library.
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	// Internal packages.
	"github.com/deuill/farsight/source"
)

type TestSource struct {
	data map[string]TestCase
}

func (t *TestSource) Fetch(src string) (io.Reader, error) {
	if _, exists := t.data[src]; !exists {
		return nil, fmt.Errorf("Unknown source data requested")
	}

	return strings.NewReader(t.data[src].Content), nil
}

type TestCase struct {
	Content  string
	Actual   interface{}
	Expected interface{}
}

var fetchTests = map[string]TestCase{
	// Fetch and set ID attribute.
	"html://id-test": {
		`<html><div id="hello">Hello World</div></html>`,
		&struct {
			Hello string `farsight:"#hello"`
		}{},
		&struct {
			Hello string `farsight:"#hello"`
		}{
			"Hello World",
		},
	},
}

func TestFetch(t *testing.T) {
	// Register mock source.
	source.Register("html", &TestSource{data: fetchTests})

	// Execute tests sequentially.
	for k, v := range fetchTests {
		if err := Fetch(k, v.Actual, "html"); err != nil {
			t.Errorf("Fetch failed for '%s': %s", k, err)
		}

		if reflect.DeepEqual(v.Actual, v.Expected) == false {
			t.Errorf("Testing '%s' failed: expected '%v', actual '%v'\n", k, v.Expected, v.Actual)
		}
	}
}
