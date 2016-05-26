package http

import (
	// Standard library.
	"bytes"
	"io"
	"net/http"

	// Internal packages.
	"github.com/deuill/farsight/source"
)

type HTTP struct{}

func (h *HTTP) Fetch(src string) (io.Reader, error) {
	// Attempt to fetch resource from source endpoint.
	resp, err := http.Get(src)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var buffer bytes.Buffer

	// Fetch and copy body content locally. This incurs some extra overhead, but
	// avoids having to pass responsibility for closing the Reader to the caller.
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}

func init() {
	h := &HTTP{}

	// Register HTTP source for both "http" and "https" endpoints.
	source.Register("http", h)
	source.Register("https", h)
}
