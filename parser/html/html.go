package html

import (
	// Standard library.
	"bytes"
	"io"

	// Internal packages.
	"github.com/deuill/farsight/parser"

	// Third-party packages.
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type HTMLParser struct{}

func (h *HTMLParser) Parse(r io.Reader) (parser.Document, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return &HTMLDocument{nodes: []*html.Node{doc}}, nil
}

type HTMLDocument struct {
	nodes []*html.Node
}

func (h *HTMLDocument) Filter(attr string) (parser.Document, error) {
	sel, err := cascadia.Compile(attr)
	if err != nil {
		return nil, err
	}

	sub := &HTMLDocument{nodes: []*html.Node{}}
	for _, n := range h.nodes {
		sub.nodes = append(sub.nodes, sel.MatchAll(n)...)
	}

	return sub, nil
}

func (h *HTMLDocument) String() string {
	var buf bytes.Buffer

	for _, n := range h.nodes {
		buf.WriteString(getNodeText(n))
	}

	return buf.String()
}

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
