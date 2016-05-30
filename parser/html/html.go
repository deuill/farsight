// Copyright 2016 Alexander Palaistras. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package html

import (
	// Standard library.
	"bytes"
	"fmt"
	"io"
	"strings"

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
func (h *HTMLDocument) Filter(sel string) (parser.Document, error) {
	var attr string

	// Parse optional attribute selector.
	idx := strings.LastIndex(sel, "/")
	if idx > 0 {
		attr = sel[(idx + 1):]
		sel = sel[:idx]
	}

	s, err := cascadia.Compile(sel)
	if err != nil {
		return nil, err
	}

	sub := &HTMLDocument{nodes: []*html.Node{}}
	for _, n := range h.nodes {
		sub.nodes = append(sub.nodes, s.MatchAll(n)...)
	}

	if len(sub.nodes) == 0 {
		return nil, fmt.Errorf("Selector '%s' matched no elements", sel)
	}

	// Loop through node attributes and attempt to match requested attribute.
	// If a matching attribute key is matched, replace current node with a
	// TextNode containing only the attribute value.
	if attr != "" {
		for i, n := range sub.nodes {
			var found bool

			for _, a := range n.Attr {
				if a.Key == attr {
					sub.nodes[i] = &html.Node{Type: html.TextNode, Data: a.Val}
					found = true
				}
			}

			if !found {
				return nil, fmt.Errorf("Unable to find attribute '%s' for selector '%s'", attr, sel)
			}
		}
	}

	return sub, nil
}

// Slice decomposes the target HTMLDocument into a slice of HTMLDocument types,
// each containing a single node from the parent's list of nodes.
func (h *HTMLDocument) Slice() []parser.Document {
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

		return strings.TrimSpace(buf.String())
	}

	return ""
}

func init() {
	// Register HTML language parser for later use.
	parser.Register("html", &HTMLParser{})
}
