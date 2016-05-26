// Copyright 2016 Alexander Palaistras. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package html

import (
	// Standard library.
	"bytes"
	"fmt"
	"io"

	// Internal packages.
	"github.com/deuill/farsight/parser"

	// Third-party packages.
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

// HTMLParser represents a parser and tokeniser for HTML documents.
type HTMLParser struct{}

// Parse reads an HTML document from the reader passed, and returns a document
// containing a single parent node. An error is returned if parsing fails.
func (h *HTMLParser) Parse(r io.Reader) (parser.Document, error) {
	n, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return &HTMLDocument{nodes: []*html.Node{n}}, nil
}

// HTMLDocument represents a collection of nodes under a single parent container.
type HTMLDocument struct {
	nodes []*html.Node
}

// Filter traverses the document tree and attempts to match elements against
// the provided CSS selector. On success, a new document is returned, containing
// a list of all matched elements. An error is returned if the CSS selector is
// malformed, or no elements were matched.
func (h *HTMLDocument) Filter(attr string) (parser.Document, error) {
	sel, err := cascadia.Compile(attr)
	if err != nil {
		return nil, err
	}

	sub := &HTMLDocument{nodes: []*html.Node{}}
	for _, n := range h.nodes {
		sub.nodes = append(sub.nodes, sel.MatchAll(n)...)
	}

	if len(sub.nodes) == 0 {
		return nil, fmt.Errorf("Attribute '%s' matched no elements", attr)
	}

	return sub, nil
}

// List decomposes the target HTMLDocument into a slice of HTMLDocument types,
// each containing a single node from the parent's list of nodes.
func (h *HTMLDocument) List() []parser.Document {
	var docs []parser.Document

	for _, n := range h.nodes {
		docs = append(docs, &HTMLDocument{nodes: []*html.Node{n}})
	}

	return docs
}

// Returns the document contents by traversing the tree and concatenating all
// data contained within text nodes.
func (h *HTMLDocument) String() string {
	var buf bytes.Buffer

	for _, n := range h.nodes {
		buf.WriteString(getNodeText(n))
	}

	return buf.String()
}

// Traverse document tree and return the first text node's contents as a string.
func getNodeText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.FirstChild != nil {
		var buf bytes.Buffer

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			buf.WriteString(getNodeText(c))
		}

		return buf.String()
	}

	return ""
}

func init() {
	// Register HTML language parser for later use.
	parser.Register("html", &HTMLParser{})
}
