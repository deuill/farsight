// Copyright 2016 Alexander Palaistras. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package parser

import (
	// Standard library.
	"fmt"
	"io"
)

type Parser interface {
	Parse(io.Reader) (Document, error)
}

type Document interface {
	Filter(attr string) (Document, error)
	List() []Document
	String() string
}

// A map of all registered parsers.
var parsers map[string]Parser

// Register a parser under a unique name.
func Register(name string, rcvr Parser) error {
	if _, exists := parsers[name]; exists {
		return fmt.Errorf("Parser '%s' already registered, refusing to overwrite", name)
	}

	parsers[name] = rcvr
	return nil
}

// Read and parse document by calling the appropriate concrete parser for the kind
// passed. Returns a parsed document, which can then be queried against, or an
// error if parsing fails.
func Parse(kind string, src io.Reader) (Document, error) {
	if _, exists := parsers[kind]; !exists {
		return nil, fmt.Errorf("Parser for '%s' not found", kind)
	}

	return parsers[kind].Parse(src)
}

func init() {
	parsers = make(map[string]Parser)
}
