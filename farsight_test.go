// Copyright 2016 Alexander Palaistras. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

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

// Test cases for `farsight.Fetch` function.
var fetchTests = map[string]TestCase{
	"html://string": {
		`<html><div id="hello">Hello World</div></html>`,
		&struct {
			Text string `farsight:"#hello"`
		}{},
		&struct {
			Text string `farsight:"#hello"`
		}{
			"Hello World",
		},
	},
	"html://slice": {
		`<body><ul id="g"><li>Hello</li><li>World</li></ul></body>`,
		&struct {
			List []string `farsight:"#g li"`
		}{},
		&struct {
			List []string `farsight:"#g li"`
		}{
			[]string{"Hello", "World"},
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
